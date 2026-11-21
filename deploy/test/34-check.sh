#!/usr/bin/env bash
# The gate of phase 34: the promises of 0.1.5 that rest on the records of phase 34, and
# the evidence those records bind.
#
#   bash deploy/test/34-check.sh RULE...
#
# RULE is one of quotes, selftest:
#
# quotes    every required row of docs/beta.md (REQUIRED_ROWS) exists once, names one of
#           the records of 0.1.5 (RECORDS) in its Record column, and every phrase of its
#           Cited words column is found in that record by grep -F, byte for byte. A phrase
#           is a sentence between straight double quotes or a backticked token outside
#           them. A phrase that spans a line break of the record, or that holds a backtick
#           and is cited as a backticked token, is not found, and the rule fails.
# selftest  one negative case per rule, each beside a control that passes, on copies in a
#           temporary directory; the real files are never written.
#
# Why a second gate. deploy/test/32-check.sh is the gate of phase 32 and checks the rows
# that rest on deploy/BACKUP-VERIFY-RUN.md; it stays as it is, because it is the gate of
# record of that phase. Each phase's gate checks its own rows. This one is its sibling: the
# same shape, the same reading of a row and of a phrase, for the rows that rest on the
# records phase 34 made.
#
# Environment:
#   ROOT     the repository; default: the repository that holds this file. Every git
#            command of this file runs as `git -C "$ROOT"`, so the gate works from a
#            worktree and from a copy of the file elsewhere.
#   REC_ROOT where the records are read from; default ROOT. The names in RECORDS are
#            relative to it.
#   BETA     the page of promises; default $ROOT/docs/beta.md
#
# The last line is `34-CHECK OK (N rules)` or `34-CHECK FAILED (K of N)`; the selftest ends
# `SELFTEST OK (N cases)` or `SELFTEST FAILED (K of N)`. Exit status: 0 every rule passed,
# 1 a rule failed, 2 usage or a rule could not run.
#
# bash 3.2 and BWK awk (macOS): no mapfile, no associative arrays, no gawk extensions.
set -u
export LC_ALL=C

HERE=$(cd "$(dirname "$0")" && pwd)
if [ -z "${ROOT:-}" ]; then
  ROOT=$(git -C "$HERE" rev-parse --show-toplevel 2> /dev/null) || ROOT=$(cd "$HERE/../.." && pwd)
fi
ME="$HERE/$(basename "$0")"
REC_ROOT=${REC_ROOT:-$ROOT}
BETA=${BETA:-$ROOT/docs/beta.md}

# The records of 0.1.5, relative to REC_ROOT.
RECORDS="deploy/MIGRATION-RESTORE-RUN.md deploy/UPDATE-CHAIN-RUN.md deploy/ROLLBACK-RUN.md deploy/BROWSER-RUN.md deploy/BACKUP-VERIFY-RUN.md"
# The rows of docs/beta.md that must exist and cite a record of 0.1.5.
REQUIRED_ROWS="P-27"

TMPROOT=$(mktemp -d) || exit 2
trap 'rm -rf "$TMPROOT"' EXIT

RULE_FAIL=0
SKIP=77
fail() { echo "FAIL: $*"; RULE_FAIL=$((RULE_FAIL + 1)); }
ok() { echo "ok: $*"; }
need() { [ -f "$1" ] || { fail "missing file: $1"; return 2; }; }

# ------------------------------------------------------------- helpers ----

# cell ROW N : the N-th cell of a table row, trimmed.
cell() { printf '%s\n' "$1" | awk -F'|' -v n="$(($2 + 1))" '{ x = $n; gsub(/^ +| +$/, "", x); print x }'; }
# quotes_of TEXT : each "quoted" span, one per line; nonzero when the quotes do not pair up.
quotes_of() {
  printf '%s\n' "$1" | awk '{
    n = split($0, a, "\"")
    if (n % 2 == 0) bad = 1
    for (i = 2; i < n; i += 2) print a[i]
  } END { exit bad ? 1 : 0 }'
}
# unquoted_of TEXT : what is outside the quotes.
unquoted_of() { printf '%s\n' "$1" | awk '{ n = split($0, a, "\""); s = ""; for (i = 1; i <= n; i += 2) s = s " " a[i]; print s }'; }
# ticks_of TEXT : each backticked token without its backticks.
ticks_of() { printf '%s\n' "$1" | grep -o '`[^`]*`' | sed 's/^`//; s/`$//'; }
# phrases_of TEXT : the quoted sentences, then the backticked tokens outside them.
phrases_of() {
  quotes_of "$1"
  ticks_of "$(unquoted_of "$1")"
}
short() { printf '%s' "$1" | cut -c1-100; }
# is_record NAME : NAME is one of RECORDS.
is_record() { case " $RECORDS " in *" $1 "*) return 0 ;; esac; return 1; }

# --------------------------------------------------------------- rules ----

rule_quotes() {
  need "$BETA" || return 2
  local id row rec cited p n np=0 nr=0
  for id in $REQUIRED_ROWS; do
    n=$(grep -c "^| $id |" "$BETA")
    if [ "$n" != 1 ]; then fail "quotes: $id: $n rows in $BETA, want exactly one"; continue; fi
    nr=$((nr + 1))
    row=$(grep "^| $id |" "$BETA")
    rec=$(cell "$row" 3 | tr -d '` ')
    if ! is_record "$rec"; then fail "quotes: $id: its Record column names '$rec', which is not a record of 0.1.5 ($RECORDS)"; continue; fi
    if [ ! -f "$REC_ROOT/$rec" ]; then fail "quotes: $id: the record $rec is missing under $REC_ROOT"; continue; fi
    cited=$(cell "$row" 4)
    if ! quotes_of "$cited" > /dev/null; then fail "quotes: $id: the double quotes of its Cited words do not pair up"; continue; fi
    n=0
    while IFS= read -r p; do
      [ -n "$p" ] || continue
      n=$((n + 1)); np=$((np + 1))
      grep -qF -- "$p" "$REC_ROOT/$rec" || fail "quotes: $id: not in $rec, byte for byte: $(short "$p")"
    done < <(phrases_of "$cited")
    [ "$n" -gt 0 ] || fail "quotes: $id: no cited words"
  done
  [ "$RULE_FAIL" = 0 ] && ok "quotes: $nr rows ($REQUIRED_ROWS) cite a record of 0.1.5, $np phrases found byte for byte"
}

# ------------------------------------------------------------- selftest ----

ST_N=0
ST_K=0
st_pass() { ST_N=$((ST_N + 1)); echo "ok $ST_N $1"; }
st_bad() { ST_N=$((ST_N + 1)); ST_K=$((ST_K + 1)); echo "FAIL $ST_N $1"; }
# st_run DESC WANT-EXIT MESSAGE NAME=VALUE... -- RULE : run this file on one rule with the
# environment given, each NAME=VALUE one argument (a value may hold spaces and newlines);
# the exit status must be WANT-EXIT and the output must hold MESSAGE.
st_run() {
  local desc="$1" want="$2" msg="$3" out rc
  local -a envs=()
  shift 3
  while [ "$#" -gt 0 ] && [ "$1" != "--" ]; do envs[${#envs[@]}]="$1"; shift; done
  shift
  out=$(env ${envs[@]+"${envs[@]}"} bash "$ME" "$@" 2>&1)
  rc=$?
  if [ "$rc" != "$want" ]; then st_bad "$desc: exit $rc, want $want; last line: $(printf '%s\n' "$out" | tail -n 1)"; return; fi
  if [ -n "$msg" ] && ! printf '%s\n' "$out" | grep -qF -- "$msg"; then st_bad "$desc: the output does not say '$msg'"; return; fi
  st_pass "$desc"
}
# t_row_append ID TOKEN < FILE : append a backticked token to the Cited words of row ID.
t_row_append() {
  ID="$1" TOK="$2" awk 'BEGIN { id = "| " ENVIRON["ID"] " |"; t = ENVIRON["TOK"] }
    index($0, id) == 1 { sub(/ \|$/, ", `" t "` |") } { print }'
}
# t_row_drop ID < FILE : the file without the row ID.
t_row_drop() { ID="$1" awk 'BEGIN { id = "| " ENVIRON["ID"] " |" } index($0, id) != 1 { print }'; }
# t_row_sub ID OLD NEW < FILE : in row ID, the first OLD replaced by NEW.
t_row_sub() {
  ID="$1" O="$2" N="$3" awk 'BEGIN { id = "| " ENVIRON["ID"] " |"; o = ENVIRON["O"]; n = ENVIRON["N"]; ol = length(o) }
    index($0, id) == 1 { i = index($0, o); if (i > 0) $0 = substr($0, 1, i - 1) n substr($0, i + ol) } { print }'
}

selftest() {
  local T
  T=$(mktemp -d "$TMPROOT/st.XXXXXX")
  cp "$BETA" "$T/good.md"
  # quotes
  st_run 'quotes: control, the page as it is' 0 'ok: quotes' "BETA=$T/good.md" -- quotes
  t_row_append P-27 'NO SUCH PHRASE 34' < "$T/good.md" > "$T/q-absent.md"
  st_run 'quotes: P-27 cites a phrase its record does not hold' 1 'P-27: not in deploy/MIGRATION-RESTORE-RUN.md, byte for byte: NO SUCH PHRASE 34' "BETA=$T/q-absent.md" -- quotes
  t_row_sub P-27 '`deploy/MIGRATION-RESTORE-RUN.md`' '`deploy/BACKUP-RUN.md`' < "$T/good.md" > "$T/q-other.md"
  st_run 'quotes: P-27 names a record that is not a record of 0.1.5' 1 "P-27: its Record column names 'deploy/BACKUP-RUN.md', which is not a record of 0.1.5" "BETA=$T/q-other.md" -- quotes
  t_row_drop P-27 < "$T/good.md" > "$T/q-gone.md"
  st_run 'quotes: P-27 is gone' 1 'P-27: 0 rows' "BETA=$T/q-gone.md" -- quotes
  if [ "$ST_K" = 0 ]; then
    echo "SELFTEST OK ($ST_N cases)"
    return 0
  fi
  echo "SELFTEST FAILED ($ST_K of $ST_N)"
  return 1
}

# ----------------------------------------------------------------- main ----

usage() { echo "usage: 34-check.sh quotes ... | selftest" >&2; exit 2; }
[ "$#" -gt 0 ] || usage
if [ "$1" = selftest ]; then
  [ "$#" = 1 ] || usage
  selftest
  exit $?
fi

RULES=""
for a in "$@"; do
  case "$a" in
    quotes) RULES="$RULES $a" ;;
    *) usage ;;
  esac
done

N=0
K=0
S=0
HARD=0
for rule in $RULES; do
  RULE_FAIL=0
  "rule_$rule"
  rc=$?
  if [ "$rc" = "$SKIP" ] && [ "$RULE_FAIL" = 0 ]; then
    S=$((S + 1))
    continue
  fi
  N=$((N + 1))
  [ "$rc" = 2 ] && HARD=1
  if [ "$RULE_FAIL" != 0 ] || { [ "$rc" != 0 ] && [ "$rc" != 1 ]; }; then K=$((K + 1)); fi
done
if [ "$K" = 0 ]; then
  if [ "$S" = 0 ]; then echo "34-CHECK OK ($N rules)"; else echo "34-CHECK OK ($N rules, $S skipped)"; fi
  exit 0
fi
echo "34-CHECK FAILED ($K of $N)"
[ "$HARD" = 1 ] && exit 2
exit 1
