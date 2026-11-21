#!/usr/bin/env bash
# Phase 34: UAT-01, UAT-02 and UAT-03 played in a real browser against a dev server
# with its own data directory.
#
# A sibling of deploy/test/backup-verify.sh in its shape (bound files, an evidence
# commit pinned to HEAD, verdicts read from KEY: value files, never from a trace),
# but without a container: the three flows are a browser's, and they run on the Mac.
# It builds this tree and tools/totpcode into UAT_DIR, starts the binary there on a
# data directory of its own (loopback, mail queued to an unreachable port, never
# sent), runs a Playwright script against it, and computes the verdicts from the
# KEY: value lines the script prints.
#
#   browser-verify.sh selftest      the pure functions below, a table of cases, no browser
#   browser-verify.sh run           the evidence: build, dev server, the script, the raw
#                                   files below deploy/test/out/mac-t34-uat
#   browser-verify.sh appendix      the appendix blocks of the fixed paste list, from the files below E
#   browser-verify.sh record-check  the proof of deploy/BROWSER-RUN.md (R) against the scratch
#                                   root (E): every block byte for byte, the script text
#                                   included, every cited file pasted, one Evidence commit,
#                                   pinned and unchanged since
#
# R and E default to deploy/BROWSER-RUN.md and deploy/test/out.
#
# Why the script is not in the tree: the stack allows no JavaScript beyond
# cmd/holzcloud/assets/htmx.min.js (CLAUDE.md), and a Playwright script is
# JavaScript. It lives outside the repository (UAT_JS); `run` copies it byte for
# byte to mac-t34-uat/script.txt, and the record pastes that copy, so its text and
# its output reach the repository only as the record's appendix. `run` and
# `record-check` refuse while the tracked JavaScript is anything but the htmx file.
#
# The flows (the script prints the keys; this file judges them):
#   UAT-01  a second administrator, created on /admin/users/new behind the password
#           question, is forced through two-step enrolment at the first sign-in,
#           sees ten recovery codes and the factor on /admin/konto, and is asked for
#           a code at the next sign-in; a recovery code works once
#   UAT-02  a visitor orders with invoice and, in a second order, with prepayment;
#           both orders are listed in the admin with their payment method (Payrexx
#           is never configured)
#   UAT-03  a contact message sent after a wait of at least four seconds arrives;
#           a submission that fills the honeypot gets the same answer and lands in
#           quarantine with its reason; every visitor request goes to one origin
#
# Nothing secret leaves UAT_DIR: the script makes its passwords at run time and
# prints counts, shapes and names, never a TOTP secret, a code, a recovery code, a
# password, a cookie value or an order token. leak_scan reads uat.txt and
# server-excerpt.txt in `run` and the whole record in `record-check`.
#
# The bound files of the evidence set (BOUND_FILES below): this file,
# tools/totpcode/main.go and plugins/kontaktformular/kontaktformular.zip (the
# archive the script uploads), together with the binary built from cmd, internal,
# go.mod and go.sum of the evidence commit. `run` refuses to start while one of
# them differs from HEAD or is untracked, and pins NEW_REF to the short sha of HEAD.
#
# Environment:
#   UAT_DIR     absolute, outside the repository (and outside the main checkout of a
#               worktree); the driver creates it. Refused when it exists and holds
#               anything but what an earlier run of this driver left there.
#   UAT_JS      absolute path of the Playwright script, outside the repository; must exist
#   PW          the Playwright module directory (holds package.json); nothing is installed
#   PORT        18793, the dev server's port on 127.0.0.1; refused when something answers there
#   REHEARSAL   1 lets `run` start while a bound file differs from HEAD; summary.txt
#               then says REHEARSAL, and such a set is never evidence
#
# Names fixed for the phase: evidence directory deploy/test/out/mac-t34-uat (summary.txt,
# versions.txt, uat.txt, script.txt, server-excerpt.txt); run files
# deploy/test/out/mac-t34-evidence-commit.txt and deploy/test/out/mac-t34-pass1.{started,finished,status,log}.
# The server log, the script's stderr and its screenshots stay in UAT_DIR.
#
# The server: started by this file in the background, its pid kept; that pid and no
# other is stopped (SIGTERM, then a wait) at the end and on every exit path.
#
# Exit codes: 0 ok, 1 failure (a verdict that is not PASS included), 2 usage or a
# refused guard, 3 the platform cannot give the evidence (no go, node or curl).
#
# Written for bash 3.2 (macOS): no arrays, no mapfile, no case conversion. Run it
# with bash explicitly.

set -euo pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
cd "$ROOT"
# The main checkout when this tree is a git worktree (the same directory otherwise).
MAIN_ROOT=$(cd "$(git rev-parse --git-common-dir)/.." && pwd)

UAT_DIR=${UAT_DIR:-}
UAT_JS=${UAT_JS:-}
PW=${PW:-}
PORT=${PORT:-18793}
REHEARSAL=${REHEARSAL:-}

BASE_OUT=deploy/test/out
T34_OUT=$BASE_OUT/mac-t34-uat
RUN_PREFIX=$BASE_OUT/mac-t34-pass1
EVIDENCE_FILE=$BASE_OUT/mac-t34-evidence-commit.txt
BOUND_FILES="deploy/test/browser-verify.sh tools/totpcode/main.go plugins/kontaktformular/kontaktformular.zip"
BINARY_SOURCE="cmd internal go.mod go.sum"
HTMX_JS=cmd/holzcloud/assets/htmx.min.js
# What a run leaves in UAT_DIR, and the mark that says a run of this file made it.
UAT_MARK=.browser-verify-t34
UAT_ENTRIES="$UAT_MARK holzcloud totpcode data server.log node.out node.err browser.txt shots"

say() { printf '%s\n' "$*" >&2; }
die() {
  say "browser-verify.sh: $*"
  exit 1
}
refuse() {
  say "browser-verify.sh: $*"
  exit 2
}
usage() {
  awk 'NR > 1 && /^$/ { exit } NR > 1 { print }' "${BASH_SOURCE[0]}" >&2
  exit 2
}

# ---------------------------------------------------------------------------
# Pure functions. No browser, no server, no network: arguments or files in, one
# answer out. `selftest` holds a table of cases, the failing ones included.
# ---------------------------------------------------------------------------

# inside P D: success when P is D or lies below it (textually).
inside() {
  case "$1" in
    "$2" | "$2"/*) return 0 ;;
  esac
  return 1
}

# phys_path P: P with its longest existing ancestor resolved to a physical path
# (symlinks followed), the rest appended unchanged. /tmp becomes /private/tmp.
phys_path() {
  local p=$1 rest="" d
  while [ ! -d "$p" ] && [ "$p" != / ] && [ -n "$p" ]; do
    rest="/$(basename "$p")$rest"
    p=$(dirname "$p")
  done
  d=$(cd -P "$p" 2> /dev/null && pwd) || d=$p
  [ "$d" = / ] && d=""
  printf '%s%s\n' "$d" "$rest"
}

# outside_repo P: success when P, logical and physical, is neither inside the
# repository nor inside the main checkout, and is no parent of either.
outside_repo() {
  local p=$1 q r
  q=$(phys_path "$p")
  for r in "$ROOT" "$MAIN_ROOT" "$(phys_path "$ROOT")" "$(phys_path "$MAIN_ROOT")"; do
    if inside "$p" "$r" || inside "$q" "$r" || inside "$r" "$p" || inside "$r" "$q"; then
      return 1
    fi
  done
  return 0
}

# safe_uat_dir P: an absolute path, no parent step, not the root, not inside the
# repository or the main checkout and no parent of them.
safe_uat_dir() {
  local p=${1:-}
  case "$p" in
    /*) ;;
    *) return 1 ;;
  esac
  case "$p" in
    / | */.. | */../* | */. | */./*) return 1 ;;
  esac
  [ "$p" != "${HOME:-/nonexistent}" ] || return 1
  outside_repo "$p"
}

# safe_script P: an absolute path of an existing file outside the repository.
safe_script() {
  local p=${1:-}
  case "$p" in
    /*) ;;
    *) return 1 ;;
  esac
  [ -f "$p" ] || return 1
  outside_repo "$p"
}

# uat_dir_left P: success when P is absent, empty, or holds the mark of this file
# and nothing but the entries a run of it leaves (UAT_ENTRIES).
uat_dir_left() {
  local p=$1 e
  [ -e "$p" ] || return 0
  [ -d "$p" ] || return 1
  [ -n "$(ls -A "$p")" ] || return 0
  [ -f "$p/$UAT_MARK" ] || return 1
  for e in $(ls -A "$p"); do
    case " $UAT_ENTRIES " in
      *" $e "*) ;;
      *) return 1 ;;
    esac
  done
  return 0
}

# kv_lines: from stdin, only the lines "KEY: value" (an upper-case key of letters,
# digits and dashes, a value that is not empty), unchanged and in order.
kv_lines() {
  grep -E '^[A-Z][A-Z0-9-]*: .' || true
}

# leak_scan FILE: success when no line of FILE holds what must never be pasted: a
# base32 run of 16 or more (a TOTP secret), a run of exactly 32 hex digits (an
# order token; a sha256 or a commit sha is longer), a basket or session cookie with
# its value, or a recovery code (five and five of its alphabet). Prints the number
# and kind of every such line, never the line.
leak_scan() {
  local f=${1:-} bad=0 n
  [ -f "$f" ] || {
    say "leak_scan: no file $f"
    return 1
  }
  for n in $(grep -nE '[A-Z2-7]{16,}' "$f" | cut -d: -f1); do
    say "leak_scan: $f line $n: a base32 run of 16 or more"
    bad=1
  done
  for n in $(grep -nE '(^|[^0-9a-f])[0-9a-f]{32}([^0-9a-f]|$)' "$f" | cut -d: -f1); do
    say "leak_scan: $f line $n: a run of 32 hex digits"
    bad=1
  done
  for n in $(grep -nE '(hc_warenkorb|holzcloud_session)=[^;[:space:]]' "$f" | cut -d: -f1); do
    say "leak_scan: $f line $n: a cookie with its value"
    bad=1
  done
  for n in $(grep -nE '(^|[^0-9a-z])[2-9b-df-hj-km-np-tv-xz]{5}-[2-9b-df-hj-km-np-tv-xz]{5}([^0-9a-z]|$)' "$f" | cut -d: -f1); do
    say "leak_scan: $f line $n: a recovery code"
    bad=1
  done
  [ "$bad" -eq 0 ]
}

# kv FILE KEY: the value of the one line "KEY: value"; empty when absent. kv_count
# FILE KEY: how many such lines there are.
kv() { sed -n "s/^$2: //p" "$1" 2> /dev/null | head -n 1; }
kv_count() { grep -c "^$2: " "$1" 2> /dev/null || true; }

# kv_is FILE KEY VALUE: the key occurs exactly once and has exactly that value.
kv_is() {
  [ "$(kv_count "$1" "$2")" = 1 ] && [ "$(kv "$1" "$2")" = "$3" ]
}
# kv_starts FILE KEY PREFIX: the key occurs exactly once and its value begins with PREFIX.
kv_starts() {
  [ "$(kv_count "$1" "$2")" = 1 ] || return 1
  case "$(kv "$1" "$2")" in
    "$3"*) return 0 ;;
  esac
  return 1
}

# verdict_uat01 FILE: PASS, FAIL or NOT-RUN (no file) for the UAT-01 keys of uat.txt.
verdict_uat01() {
  local f=${1:-}
  if [ ! -f "$f" ]; then
    echo NOT-RUN
    return 0
  fi
  if kv_is "$f" ADMIN1-ENROLLED yes &&
    kv_is "$f" ADMIN2-CREATED yes &&
    kv_is "$f" PASSWORD-QUESTION yes &&
    kv_is "$f" FORCED-REDIRECT /admin/2fa/einrichten &&
    kv_is "$f" SECRET-SHOWN yes &&
    kv_is "$f" RECOVERY-CODES 10 &&
    kv_is "$f" KONTO-FACTOR on &&
    kv_is "$f" RELOGIN-ASKS-CODE yes &&
    kv_is "$f" RELOGIN-OK yes &&
    kv_starts "$f" RECOVERY-CODE-ONCE "yes "; then
    echo PASS
  else
    echo FAIL
  fi
}

# js_tracked_ok: the tracked JavaScript of the tree is exactly the htmx file.
js_tracked_ok() {
  [ "$(git ls-files '*.js' '*.mjs' '*.cjs')" = "$HTMX_JS" ]
}

# safe_out PATH: a path directly below deploy/test/out whose name begins with mac-t34-.
safe_out() {
  local p=${1:-} rest
  case "$p" in
    deploy/test/out/mac-t34-?*) ;;
    *) return 1 ;;
  esac
  rest=${p#deploy/test/out/}
  case "$rest" in
    */* | *..*) return 1 ;;
  esac
  return 0
}

# ---------------------------------------------------------------------------
# The run
# ---------------------------------------------------------------------------

# bound_check: the bound files are tracked and equal to HEAD. Prints one line per
# problem and fails when there is one.
bound_check() {
  local f bad=0
  for f in $BOUND_FILES; do
    if ! git ls-files --error-unmatch -- "$f" > /dev/null 2>&1; then
      say "browser-verify.sh: bound file $f is not tracked"
      bad=1
    elif ! git diff --quiet HEAD -- "$f"; then
      say "browser-verify.sh: bound file $f differs from HEAD"
      bad=1
    fi
  done
  if ! git diff --quiet HEAD -- $BINARY_SOURCE || [ -n "$(git ls-files --others --exclude-standard -- $BINARY_SOURCE)" ]; then
    say "browser-verify.sh: cmd, internal, go.mod or go.sum differ from HEAD"
    bad=1
  fi
  [ "$bad" -eq 0 ]
}

# port_answers: success when something answers HTTP on 127.0.0.1:PORT.
port_answers() {
  local rc=0
  curl -s -o /dev/null --max-time 2 "http://127.0.0.1:$PORT/" || rc=$?
  [ "$rc" -ne 7 ] && [ "$rc" -ne 28 ]
}

SRV_PID=""
# start_server: the binary of UAT_DIR on UAT_DIR/data, loopback, a clean environment.
start_server() {
  : > "$UAT_DIR/server.log"
  env -i PATH="$PATH" HOME="$HOME" TMPDIR="${TMPDIR:-/tmp}" \
    HOLZCLOUD_DATA_DIR="$UAT_DIR/data" \
    HOLZCLOUD_PORT="$PORT" \
    HOLZCLOUD_LISTEN=127.0.0.1 \
    HOLZCLOUD_SMTP_HOST=127.0.0.1 \
    HOLZCLOUD_SMTP_PORT=2 \
    HOLZCLOUD_SMTP_FROM=cms@example.test \
    "$UAT_DIR/holzcloud" < /dev/null >> "$UAT_DIR/server.log" 2>&1 &
  SRV_PID=$!
  say "browser-verify: server started, pid $SRV_PID, port $PORT, data $UAT_DIR/data"
}

# stop_server: the pid start_server kept, and no other. SIGTERM, a wait of up to
# twenty seconds, then SIGKILL to the same pid.
stop_server() {
  local i=0
  [ -n "$SRV_PID" ] || return 0
  if kill -0 "$SRV_PID" 2> /dev/null; then
    kill -TERM "$SRV_PID" 2> /dev/null || true
    while kill -0 "$SRV_PID" 2> /dev/null && [ "$i" -lt 40 ]; do
      sleep 0.5
      i=$((i + 1))
    done
    if kill -0 "$SRV_PID" 2> /dev/null; then
      say "browser-verify: server pid $SRV_PID did not stop on SIGTERM; SIGKILL"
      kill -KILL "$SRV_PID" 2> /dev/null || true
    fi
  fi
  wait "$SRV_PID" 2> /dev/null || true
  say "browser-verify: server pid $SRV_PID stopped"
  SRV_PID=""
}

# wait_ready: /admin/login answers below 500 within thirty seconds.
wait_ready() {
  local i=0 code
  while [ "$i" -lt 150 ]; do
    kill -0 "$SRV_PID" 2> /dev/null || die "the server exited before it answered (see $UAT_DIR/server.log)"
    code=$(curl -s -o /dev/null -w '%{http_code}' --max-time 2 "http://127.0.0.1:$PORT/admin/login" || true)
    case "$code" in
      2?? | 3?? | 4??) return 0 ;;
    esac
    sleep 0.2
    i=$((i + 1))
  done
  die "the server did not answer on port $PORT within thirty seconds"
}

# excerpt_of LOG DIR: the server log lines a record cites, from the JSON log LOG:
# the start, the listening address, one line counting the applied migrations with
# their lowest and highest version, the pre-upgrade snapshot, every plugin line,
# every second factor switched on, every warning and error, the shutdown; DIR is
# written as $UAT_DIR and each line is cut to 300 characters. Last, the line count.
excerpt_of() {
  local log=$1 dir=$2
  awk -v dir="$dir" '
    function out(s,   i, r) {
      r = ""
      while (dir != "" && (i = index(s, dir)) > 0) { r = r substr(s, 1, i - 1) "$UAT_DIR"; s = substr(s, i + length(dir)) }
      print substr(r s, 1, 300)
    }
    /"msg":"migration applied"/ {
      n++
      v = $0; sub(/.*"version":/, "", v); sub(/[^0-9].*/, "", v); v += 0
      if (lo == "" || v < lo) lo = v
      if (v > hi) hi = v
      if (!said) { said = 1; mark = NR }
      next
    }
    /"msg":"(starting|server starting|pre-upgrade snapshot written|shutting down|shutdown complete|two factor enabled)"/ ||
    /"msg":"[^"]*plugin/ || /"level":"(WARN|ERROR)"/ { out($0) }
    END {
      if (n) printf "migrations applied at this start: %d (versions %d to %d)\n", n, lo, hi
      else print "migrations applied at this start: 0"
      printf "server log lines: %d\n", NR
    }
  ' "$log"
}

# write_excerpt: excerpt_of over the run's server log.
write_excerpt() {
  excerpt_of "$UAT_DIR/server.log" "$UAT_DIR" > "$T34_OUT/server-excerpt.txt"
}

# write_versions: host, tools, the binary and the script.
write_versions() {
  local f=$T34_OUT/versions.txt pwv
  pwv=$(sed -n 's/^ *"version": *"\([^"]*\)".*/\1/p' "$PW/package.json" | head -n 1)
  {
    echo "host:                  $(uname -srm), macOS $(sw_vers -productVersion 2> /dev/null || echo '?')"
    echo "go:                    $(go version)"
    echo "node:                  $(node --version)"
    echo "playwright:            ${pwv:-?}"
    echo "browser:               $(cat "$UAT_DIR/browser.txt" 2> /dev/null || echo '?')"
    echo "binary:                $("$UAT_DIR/holzcloud" version 2> /dev/null | head -n 1)"
    echo "NEW_REF:               $NEW_REF"
    echo "evidence commit:       $(cat "$EVIDENCE_FILE")"
    echo "script sha256:         $(shasum -a 256 "$UAT_JS" | cut -d' ' -f1)"
  } > "$f"
  cat "$f" >&2
}

# write_summary: the verdicts from uat.txt.
write_summary() {
  local f=$T34_OUT/summary.txt
  {
    echo "---- summary (mac, dev server, Chromium) ----"
    echo "UAT-01 $(verdict_uat01 "$T34_OUT/uat.txt")"
    echo "UAT-02 NOT-RUN"
    echo "UAT-03 NOT-RUN"
    echo "BINARY: $("$UAT_DIR/holzcloud" version 2> /dev/null | head -n 1)"
    if [ -n "$REHEARSAL" ]; then
      echo "REHEARSAL: yes (bound files were allowed to differ from HEAD; this set is not evidence)"
    fi
  } > "$f"
  cat "$f" >&2
}

# run_body: everything `run` does after its guards; its output goes to the log.
run_body() {
  local rc=0 zip
  trap 'stop_server' EXIT
  trap 'stop_server; exit 130' INT TERM
  say "browser-verify: NEW_REF pinned to $NEW_REF, UAT_DIR $UAT_DIR, port $PORT"
  rm -rf "${UAT_DIR:?}"
  mkdir -p "$UAT_DIR"
  : > "$UAT_DIR/$UAT_MARK"
  rm -rf "${T34_OUT:?}"
  mkdir -p "$T34_OUT"
  say "browser-verify: building the binary and tools/totpcode"
  CGO_ENABLED=0 go build -trimpath \
    -ldflags="-X main.Version=$(git describe --tags --always) -X main.Commit=$NEW_REF" \
    -o "$UAT_DIR/holzcloud" ./cmd/holzcloud
  go build -o "$UAT_DIR/totpcode" ./tools/totpcode
  cp "$UAT_JS" "$T34_OUT/script.txt"
  cmp -s "$UAT_JS" "$T34_OUT/script.txt" || die "script.txt is not a byte copy of $UAT_JS"
  port_answers && die "something answers on 127.0.0.1:$PORT already"
  start_server
  wait_ready
  zip="$ROOT/plugins/kontaktformular/kontaktformular.zip"
  say "browser-verify: running the script"
  set +e
  S="$UAT_DIR" PW="$PW" PORT="$PORT" ZIP="$zip" node "$UAT_JS" > "$UAT_DIR/node.out" 2> "$UAT_DIR/node.err"
  rc=$?
  set -e
  say "browser-verify: the script ended with status $rc ($(wc -l < "$UAT_DIR/node.err" | tr -d ' ') lines on stderr, kept in $UAT_DIR/node.err)"
  stop_server
  kv_lines < "$UAT_DIR/node.out" > "$T34_OUT/uat.txt"
  write_excerpt
  if ! leak_scan "$T34_OUT/uat.txt" || ! leak_scan "$T34_OUT/server-excerpt.txt"; then
    mv "$T34_OUT/uat.txt" "$UAT_DIR/uat.txt.refused"
    mv "$T34_OUT/server-excerpt.txt" "$UAT_DIR/server-excerpt.txt.refused"
    say "browser-verify.sh: leak_scan refused the output; it is moved to $UAT_DIR and nothing is pasted"
    write_versions
    write_summary
    return 1
  fi
  grep -E '^([A-Z][A-Z0-9-]*: |note: )' "$UAT_DIR/node.out" >&2 || true
  write_versions
  write_summary
  [ "$rc" -eq 0 ] || return 1
  grep -Eq '^UAT-0[0-9] (FAIL|NOT-RUN)$' "$T34_OUT/summary.txt" && return 1
  return 0
}

cmd_run() {
  local rc p
  [ -n "$UAT_DIR" ] || refuse "UAT_DIR is not set"
  safe_uat_dir "$UAT_DIR" || refuse "UAT_DIR '$UAT_DIR' must be absolute, outside the repository and the main checkout, and no parent of them"
  uat_dir_left "$UAT_DIR" || refuse "UAT_DIR '$UAT_DIR' exists and holds something an earlier run of this driver did not leave there"
  [ -n "$UAT_JS" ] || refuse "UAT_JS is not set"
  safe_script "$UAT_JS" || refuse "UAT_JS '$UAT_JS' must be an existing file outside the repository"
  [ -n "$PW" ] && [ -f "$PW/package.json" ] || refuse "PW '$PW' is not a Playwright module directory"
  case "$PORT" in
    "" | *[!0-9]*) refuse "PORT '$PORT' is not a number" ;;
  esac
  [ "$PORT" -ge 1024 ] && [ "$PORT" -le 65535 ] || refuse "PORT $PORT is out of range"
  for p in "$T34_OUT" "$RUN_PREFIX.log" "$EVIDENCE_FILE"; do
    safe_out "$p" || refuse "refusing path $p"
  done
  js_tracked_ok || die "the tracked JavaScript is not exactly $HTMX_JS: $(git ls-files '*.js' '*.mjs' '*.cjs' | tr '\n' ' ')"
  if ! bound_check; then
    if [ -n "$REHEARSAL" ]; then
      say "browser-verify.sh: REHEARSAL=1: going on although a bound file differs from HEAD; this set is not evidence"
    else
      die "a bound file differs from HEAD or is untracked; commit first (the evidence is bound to one commit)"
    fi
  fi
  for p in go node curl shasum; do
    command -v "$p" > /dev/null 2>&1 || {
      say "browser-verify.sh: $p is not installed"
      exit 3
    }
  done
  port_answers && refuse "something answers on 127.0.0.1:$PORT; choose another PORT"
  NEW_REF=$(git rev-parse --short HEAD) || die "HEAD is not a commit"
  mkdir -p "$BASE_OUT"
  rm -f "$RUN_PREFIX.started" "$RUN_PREFIX.finished" "$RUN_PREFIX.status" "$RUN_PREFIX.log"
  printf '%s\n' "$NEW_REF" > "$EVIDENCE_FILE"
  date -u +%Y-%m-%dT%H:%M:%SZ > "$RUN_PREFIX.started"
  set +e
  (set -e; run_body) 2>&1 | tee "$RUN_PREFIX.log"
  rc=${PIPESTATUS[0]}
  set -e
  date -u +%Y-%m-%dT%H:%M:%SZ > "$RUN_PREFIX.finished"
  printf '%s\n' "$rc" > "$RUN_PREFIX.status"
  say "browser-verify: run ended with status $rc"
  exit "$rc"
}

# ---------------------------------------------------------------------------
# selftest
# ---------------------------------------------------------------------------

ST_N=0
ST_BAD=0
ST_TMP=""
st_eq() { # NAME EXPECTED ACTUAL
  ST_N=$((ST_N + 1))
  if [ "$2" != "$3" ]; then
    ST_BAD=$((ST_BAD + 1))
    say "SELFTEST BAD $1: expected '$2', got '$3'"
  fi
}
st_ok() { # NAME cmd args...: the command must return 0
  local n=$1
  shift
  ST_N=$((ST_N + 1))
  if ! "$@" > /dev/null 2>&1; then
    ST_BAD=$((ST_BAD + 1))
    say "SELFTEST BAD $n: expected success from $*"
  fi
}
st_fail() { # NAME cmd args...: the command must return non-zero
  local n=$1
  shift
  ST_N=$((ST_N + 1))
  if "$@" > /dev/null 2>&1; then
    ST_BAD=$((ST_BAD + 1))
    say "SELFTEST BAD $n: expected failure from $*"
  fi
}

selftest_guards() {
  local f=$ST_TMP/uat.js d=$ST_TMP/left
  printf 'x\n' > "$f"
  st_fail "safe_uat_dir: a relative path" safe_uat_dir scratch/t34-uat
  st_fail "safe_uat_dir: empty" safe_uat_dir ""
  st_fail "safe_uat_dir: the repository itself" safe_uat_dir "$ROOT"
  st_fail "safe_uat_dir: a path inside the repository" safe_uat_dir "$ROOT/deploy/test/out/uat"
  st_fail "safe_uat_dir: the main checkout" safe_uat_dir "$MAIN_ROOT"
  st_fail "safe_uat_dir: a path inside the main checkout" safe_uat_dir "$MAIN_ROOT/data/uat"
  st_fail "safe_uat_dir: a parent of the repository" safe_uat_dir "$(dirname "$MAIN_ROOT")"
  st_fail "safe_uat_dir: the root" safe_uat_dir /
  st_fail "safe_uat_dir: a parent step" safe_uat_dir "$ST_TMP/a/../b"
  st_ok "safe_uat_dir: an absolute path elsewhere" safe_uat_dir "$ST_TMP/run"
  st_ok "safe_script: an existing file elsewhere" safe_script "$f"
  st_fail "safe_script: a missing file" safe_script "$ST_TMP/absent.js"
  st_fail "safe_script: a relative path" safe_script uat.js
  st_fail "safe_script: a file inside the repository" safe_script "$ROOT/deploy/test/backup-verify.sh"
  st_ok "uat_dir_left: absent" uat_dir_left "$d"
  mkdir -p "$d"
  st_ok "uat_dir_left: empty" uat_dir_left "$d"
  : > "$d/server.log"
  st_fail "uat_dir_left: a run's file without the mark" uat_dir_left "$d"
  : > "$d/$UAT_MARK"
  mkdir -p "$d/data"
  st_ok "uat_dir_left: the mark and a run's entries" uat_dir_left "$d"
  : > "$d/notes.txt"
  st_fail "uat_dir_left: something a run does not leave" uat_dir_left "$d"
  st_ok "js_tracked_ok: the tree holds only the htmx file" js_tracked_ok
}

selftest_excerpt() {
  local log=$ST_TMP/server.log got
  printf '%s\n' \
    '{"time":"t1","level":"INFO","msg":"starting","version":"v9","commit":"abc1234"}' \
    '{"time":"t2","level":"INFO","msg":"sqlite opened","path":"/x+y/data/holzcloud.sqlite"}' \
    '{"time":"t3","level":"INFO","msg":"migration applied","version":1,"type":"sql"}' \
    '{"time":"t4","level":"INFO","msg":"migration applied","version":69,"type":"sql"}' \
    '{"time":"t5","level":"INFO","msg":"pre-upgrade snapshot written","path":"/x+y/data/pre-upgrade.sqlite"}' \
    '{"time":"t6","level":"INFO","msg":"plugin installed","id":"kontaktformular"}' \
    '{"time":"t7","level":"WARN","msg":"mail delivery failed","attempt":1}' \
    '{"time":"t8","level":"INFO","msg":"shutdown complete"}' > "$log"
  got=$(excerpt_of "$log" /x+y)
  st_eq "excerpt_of: the cited lines, the directory named, the migrations counted" '{"time":"t1","level":"INFO","msg":"starting","version":"v9","commit":"abc1234"}
{"time":"t5","level":"INFO","msg":"pre-upgrade snapshot written","path":"$UAT_DIR/data/pre-upgrade.sqlite"}
{"time":"t6","level":"INFO","msg":"plugin installed","id":"kontaktformular"}
{"time":"t7","level":"WARN","msg":"mail delivery failed","attempt":1}
{"time":"t8","level":"INFO","msg":"shutdown complete"}
migrations applied at this start: 2 (versions 1 to 69)
server log lines: 8' "$got"
}

selftest_kv_lines() {
  local got
  got=$(printf '%s\n' 'note: setting up' 'ADMIN2-CREATED: yes' 'lower: no' 'RECOVERY-CODES: 10' \
    '  UAT-01: indented' 'ORDER-1-NUMBER: 2026-0001' 'Mixed-Case: no' 'NOVALUE:' 'ERROR: something' '=== STEP' | kv_lines)
  st_eq "kv_lines: only KEY: value lines" "ADMIN2-CREATED: yes
RECOVERY-CODES: 10
ORDER-1-NUMBER: 2026-0001
ERROR: something" "$got"
  st_eq "kv_lines: nothing in, nothing out" "" "$(printf 'note: x\n' | kv_lines)"
}

selftest_leak_scan() {
  local f=$ST_TMP/leak.txt
  printf '%s\n' 'SECRET-SHOWN: yes' 'RECOVERY-CODES: 10' 'RECOVERY-CODE-SHAPE: 5-5 lower-case alphanumeric' \
    'ORDER-1-URL-SHAPE: matches /bestellung/[0-9a-f]{32}' 'BASKET-COOKIE: hc_warenkorb' 'SAME-ORIGIN: yes (41 requests, 1 origin)' \
    'script sha256: 0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef' 'NEW_REF: cfce101' > "$f"
  st_ok "leak_scan: the documented key lines pass" leak_scan "$f"
  printf 'SECRET: JBSWY3DPEHPK3PXPJBSWY3DP\n' > "$f"
  st_fail "leak_scan: a base32 run of 16 or more" leak_scan "$f"
  printf 'x JBSW Y3DP EHPK 3PXP\n' > "$f"
  st_ok "leak_scan: four base32 groups of four are not a run" leak_scan "$f"
  printf 'ORDER: /bestellung/0123456789abcdef0123456789abcdef done\n' > "$f"
  st_fail "leak_scan: a 32-hex run" leak_scan "$f"
  printf 'COOKIE: hc_warenkorb=abc123\n' > "$f"
  st_fail "leak_scan: hc_warenkorb= with a value" leak_scan "$f"
  printf 'COOKIE: holzcloud_session=abc123\n' > "$f"
  st_fail "leak_scan: holzcloud_session= with a value" leak_scan "$f"
  printf 'CODE: 2b3c4-d5f6g\n' > "$f"
  st_fail "leak_scan: a recovery code" leak_scan "$f"
  st_fail "leak_scan: no file" leak_scan "$ST_TMP/absent.txt"
}

# fx_uat01 FILE [KEY VALUE]: the key lines of a passing UAT-01, one value replaced
# (the value "-" removes the key).
fx_uat01() {
  local f=$1 k=${2:-} v=${3:-}
  printf '%s\n' \
    'ADMIN1-ENROLLED: yes' \
    'ADMIN2-CREATED: yes' \
    'PASSWORD-QUESTION: yes' \
    'FORCED-REDIRECT: /admin/2fa/einrichten' \
    'SECRET-SHOWN: yes' \
    'RECOVERY-CODES: 10' \
    'RECOVERY-CODE-SHAPE: 5-5 lower-case alphanumeric' \
    'KONTO-FACTOR: on' \
    'RELOGIN-ASKS-CODE: yes' \
    'RELOGIN-OK: yes' \
    'RECOVERY-CODE-ONCE: yes (accepted once, refused the second time)' > "$f"
  if [ -n "$k" ]; then
    awk -v k="$k" -v v="$v" 'index($0, k ": ") == 1 { if (v != "-") print k ": " v; next } { print }' "$f" > "$f.x" && mv "$f.x" "$f"
  fi
}

selftest_verdict_uat01() {
  local f=$ST_TMP/uat.txt kv
  fx_uat01 "$f"
  st_eq "verdict_uat01: every key right" PASS "$(verdict_uat01 "$f")"
  st_eq "verdict_uat01: no file" NOT-RUN "$(verdict_uat01 "$ST_TMP/absent.txt")"
  for kv in 'ADMIN2-CREATED|no' 'FORCED-REDIRECT|/admin/' 'SECRET-SHOWN|no' 'RECOVERY-CODES|9' \
    'KONTO-FACTOR|off' 'RELOGIN-ASKS-CODE|no' 'RELOGIN-OK|no' 'PASSWORD-QUESTION|no' 'ADMIN1-ENROLLED|no' \
    'RECOVERY-CODE-ONCE|no (accepted twice)' 'RECOVERY-CODES|-' 'FORCED-REDIRECT|-'; do
    fx_uat01 "$f" "${kv%%|*}" "${kv#*|}"
    st_eq "verdict_uat01: ${kv%%|*} '${kv#*|}'" FAIL "$(verdict_uat01 "$f")"
  done
  fx_uat01 "$f"
  printf 'RECOVERY-CODES: 10\n' >> "$f"
  st_eq "verdict_uat01: a key twice" FAIL "$(verdict_uat01 "$f")"
}

cmd_selftest() {
  local t=${TMPDIR:-/tmp}
  ST_N=0
  ST_BAD=0
  ST_TMP=$(mktemp -d "${t%/}/browser-verify-selftest.XXXXXX")
  trap 'rm -rf "${ST_TMP:?}"' EXIT
  selftest_guards
  selftest_excerpt
  selftest_kv_lines
  selftest_leak_scan
  selftest_verdict_uat01
  if [ "$ST_BAD" -ne 0 ]; then
    echo "SELFTEST FAILED: $ST_BAD of $ST_N cases"
    exit 1
  fi
  echo "SELFTEST OK ($ST_N cases)"
}

[ "$#" -gt 0 ] || usage
while [ "$#" -gt 0 ]; do
  sub=$1
  shift
  case "$sub" in
    selftest) cmd_selftest ;;
    run) cmd_run ;;
    *) usage ;;
  esac
done
