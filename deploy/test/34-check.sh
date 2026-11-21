#!/usr/bin/env bash
# The gate of phase 34: the promises of 0.1.5 that rest on the records of phase 34, and
# the evidence those records bind.
#
#   bash deploy/test/34-check.sh RULE...
#
# RULE is one of quotes, coverage, ledger, nojs, selftest:
#
# quotes    every required row of docs/beta.md (REQUIRED_ROWS) exists once, names one of
#           the records of 0.1.5 (RECORDS) in its Record column, and every phrase of its
#           Cited words column is found in that record by grep -F, byte for byte. A phrase
#           is a sentence between straight double quotes or a backticked token outside
#           them. A phrase that spans a line break of the record, or that holds a backtick
#           and is cited as a backticked token, is not found, and the rule fails.
# coverage  every limit of docs/beta.md marked (0.1.5) names exactly one record with
#           `Record: deploy/X`, and X is one of RECORDS, so no limit of 0.1.5 rests on
#           nothing; and every sentence such a limit quotes between straight double quotes
#           is found in the record it names, byte for byte. deploy/test/32-check.sh checks
#           the quotes of the limits that name deploy/BACKUP-VERIFY-RUN.md as well; since
#           phase 34 it no longer refuses a limit that names no record, and this rule does.
# ledger    the four limits of 0.1.1 that the records of phase 34 answer (LIMIT_1 to
#           LIMIT_4) each begin the first cell of a row of the ledger table of 0.1.5
#           (header `| 0.1.1 said | 0.1.5 says | Promise rows |`), each of those rows names
#           promise rows that exist, and no bullet of the section *What X does not promise*
#           still begins with one of them; the bullet of 0.1.1 that said an older binary
#           must not be started against a migrated database is gone, *The product itself*
#           names deploy/ROLLBACK-RUN.md, and *What beta means here* names P-30.
# nojs      the tracked JavaScript of the repository (`git ls-files` of *.js, *.mjs and
#           *.cjs) is exactly JS_ALLOWED, the htmx file: the browser record of phase 34
#           ran a Playwright script, and that script is pasted into the record, never
#           committed as a file.
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
#   JS_ALLOWED  the tracked JavaScript allowed (nojs only); default
#            cmd/holzcloud/assets/htmx.min.js. Only the selftest sets it to empty.
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
JS_ALLOWED=${JS_ALLOWED-cmd/holzcloud/assets/htmx.min.js}

# The records of 0.1.5, relative to REC_ROOT.
RECORDS="deploy/MIGRATION-RESTORE-RUN.md deploy/UPDATE-CHAIN-RUN.md deploy/ROLLBACK-RUN.md deploy/BROWSER-RUN.md deploy/BACKUP-VERIFY-RUN.md"
# The rows of docs/beta.md that must exist and cite a record of 0.1.5.
REQUIRED_ROWS="P-27 P-28 P-29 P-30 P-31 P-32 P-33 E-9 E-10 E-11 E-12"

# The four limits of 0.1.1 that the records of phase 34 answer, by their first words.
LIMIT_1='A restore that crosses a migration was not played'
LIMIT_2='Five or more successive updates and the pruning of the pre-upgrade snapshots were not played'
LIMIT_3='A rollback after a failed migration was not played'
LIMIT_4='The compulsory two-factor enrolment and the shop and form flows through a browser were not played'
# The bullet of 0.1.1 under *The product itself* that phase 34 rewrote, by its first sentence.
SCHEMA_OLD='- The schema and the migrations are not promised stable, and an older binary must not be started against a migrated database.'
LEDGER_HEAD='| 0.1.1 said | 0.1.5 says | Promise rows |'

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

# section2 HEADING-ERE < FILE : the lines under a level-two heading, up to the next one.
section2() {
  awk -v re="$1" '
    /^```/ { if (on) print; fence = !fence; next }
    fence { if (on) print; next }
    /^## / { if (on) exit; if ($0 ~ re) { on = 1; next } }
    on { print }
  '
}
# section3 HEADING-ERE < FILE : the lines under a level-three heading, up to the next
# heading of level two or three.
section3() {
  awk -v re="$1" '
    /^```/ { if (on) print; fence = !fence; next }
    fence { if (on) print; next }
    /^###? / { if (on) exit; if ($0 ~ re) { on = 1; next } }
    on { print }
  '
}
# ledger_rows < FILE : the data rows of the table whose header is LEDGER_HEAD.
ledger_rows() {
  H="$LEDGER_HEAD" awk 'BEGIN { h = ENVIRON["H"] }
    $0 == h { on = 1; n = 0; next }
    on && !/^\|/ { exit }
    on { n++; if (n > 1) print }'
}

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

rule_coverage() {
  need "$BETA" || return 2
  local b n rec p nq nl=0 np=0 per=""
  while IFS= read -r b; do
    nl=$((nl + 1))
    n=$(printf '%s\n' "$b" | grep -o 'Record: deploy/' | grep -c .)
    if [ "$n" != 1 ]; then fail "coverage: a limit marked (0.1.5) names $n records, want exactly one: $(short "$b")"; continue; fi
    rec=$(printf '%s\n' "$b" | sed -n 's/.*Record: \(deploy\/[A-Za-z0-9._-]*\.md\).*/\1/p')
    if ! is_record "$rec"; then fail "coverage: a limit marked (0.1.5) names '$rec', which is not a record of 0.1.5: $(short "$b")"; continue; fi
    if [ ! -f "$REC_ROOT/$rec" ]; then fail "coverage: the record $rec is missing under $REC_ROOT"; continue; fi
    if ! quotes_of "$b" > /dev/null; then fail "coverage: the double quotes of a limit marked (0.1.5) do not pair up: $(short "$b")"; continue; fi
    nq=0
    while IFS= read -r p; do
      [ -n "$p" ] || continue
      nq=$((nq + 1)); np=$((np + 1))
      grep -qF -- "$p" "$REC_ROOT/$rec" || fail "coverage: a limit marked (0.1.5): not in $rec, byte for byte: $(short "$p")"
    done < <(quotes_of "$b")
    [ "$nq" -gt 0 ] || fail "coverage: a limit marked (0.1.5) quotes nothing of $rec: $(short "$b")"
    per="$per ${rec#deploy/}"
  done < <(grep '^- (0\.1\.5) ' "$BETA")
  [ "$nl" -gt 0 ] || fail "coverage: no limit marked (0.1.5) in $BETA"
  # shellcheck disable=SC2086
  [ "$RULE_FAIL" = 0 ] && ok "coverage: $nl limits marked (0.1.5), each naming one record of 0.1.5 ($(printf '%s\n' $per | sort | uniq -c | awk '{ printf "%s%s %s", (NR > 1 ? ", " : ""), $1, $2 }')), $np quoted sentences found byte for byte"
}

rule_ledger() {
  need "$BETA" || return 2
  local np bullets rows lim row c1 c3 p n=0 prod beta
  np=$(section2 '^## What [0-9.]+ does not promise$' < "$BETA")
  [ -n "$np" ] || { fail "ledger: no section 'What X does not promise' in $BETA"; return 1; }
  bullets=$(printf '%s\n' "$np" | grep '^- ')
  rows=$(ledger_rows < "$BETA")
  [ -n "$rows" ] || fail "ledger: no table headed '$LEDGER_HEAD' in $BETA"
  for lim in "$LIMIT_1" "$LIMIT_2" "$LIMIT_3" "$LIMIT_4"; do
    if printf '%s\n' "$bullets" | grep -qF -- "- $lim"; then fail "ledger: still a bullet of the does-not-promise section: $lim"; fi
    row=$(printf '%s\n' "$rows" | grep -F -- "| $lim" | head -n 1)
    if [ -z "$row" ]; then fail "ledger: not quoted in the first column of the 0.1.5 ledger: $lim"; continue; fi
    c1=$(cell "$row" 1)
    case "$c1" in "$lim"*) ;; *) fail "ledger: the first cell does not begin with the limit: $lim" ;; esac
    c3=$(cell "$row" 3)
    printf '%s' "$c3" | grep -qE 'P-[0-9]+' || fail "ledger: the row of '$lim' names no promise row"
    for p in $(printf '%s' "$c3" | grep -oE 'P-[0-9]+'); do
      grep -q "^| $p |" "$BETA" || fail "ledger: the row of '$lim' names $p, which is not a promise row"
    done
    n=$((n + 1))
  done
  if printf '%s\n' "$bullets" | grep -qF -- "$SCHEMA_OLD"; then fail "ledger: the bullet of 0.1.1 on an older binary and a migrated database is still there"; fi
  prod=$(printf '%s\n' "$np" | section3 '^### The product itself$')
  [ -n "$prod" ] || fail "ledger: no section 'The product itself' under the does-not-promise section"
  printf '%s\n' "$prod" | grep -qF 'deploy/ROLLBACK-RUN.md' || fail "ledger: 'The product itself' does not name deploy/ROLLBACK-RUN.md"
  beta=$(section2 '^## What beta means here$' < "$BETA")
  printf '%s\n' "$beta" | grep -qE '(^|[^0-9-])P-30([^0-9]|$)' || fail "ledger: 'What beta means here' does not name P-30"
  [ "$RULE_FAIL" = 0 ] && ok "ledger: $n limits quoted in the 0.1.5 ledger and gone from the bullets, the schema bullet of 0.1.1 rewritten and naming deploy/ROLLBACK-RUN.md"
}

rule_nojs() {
  local got want
  got=$(git -C "$ROOT" ls-files -- '*.js' '*.mjs' '*.cjs') || { fail "nojs: git ls-files failed in $ROOT"; return 2; }
  # shellcheck disable=SC2086
  want=$(printf '%s\n' $JS_ALLOWED | grep . | sort)
  got=$(printf '%s\n' "$got" | grep . | sort)
  if [ "$got" != "$want" ]; then
    # shellcheck disable=SC2086
    fail "nojs: the tracked JavaScript is not exactly '$JS_ALLOWED': $(printf '%s ' $got)"
    return 1
  fi
  ok "nojs: the tracked JavaScript is exactly $JS_ALLOWED"
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
# t_after LINE NEW < FILE : NEW inserted after the line that begins with LINE.
t_after() { O="$1" N="$2" awk 'BEGIN { o = ENVIRON["O"]; n = ENVIRON["N"] } { print } index($0, o) == 1 { print n }'; }
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
  t_row_append E-11 'Evidence commit: 0000000' < "$T/good.md" > "$T/q-e11.md"
  st_run 'quotes: E-11 cites an evidence commit its record does not name' 1 'E-11: not in deploy/ROLLBACK-RUN.md, byte for byte: Evidence commit: 0000000' "BETA=$T/q-e11.md" -- quotes
  # coverage
  st_run 'coverage: control, the page as it is' 0 'ok: coverage' "BETA=$T/good.md" -- coverage
  t_after '### The unit and the platform' '- (0.1.5) A limit that rests on nothing. Why: a case of the selftest.' < "$T/good.md" > "$T/c-none.md"
  st_run 'coverage: a limit marked (0.1.5) names no record' 1 'coverage: a limit marked (0.1.5) names 0 records, want exactly one' "BETA=$T/c-none.md" -- coverage
  t_after '### The unit and the platform' '- (0.1.5) A limit of the rollback. Why: a case of the selftest. Record: deploy/ROLLBACK-RUN.md, *Not tested*: "NO SUCH PHRASE 34 in deploy/ROLLBACK-RUN.md."' < "$T/good.md" > "$T/c-absent.md"
  st_run 'coverage: a limit marked (0.1.5) quotes what its record does not say' 1 'coverage: a limit marked (0.1.5): not in deploy/ROLLBACK-RUN.md, byte for byte: NO SUCH PHRASE 34' "BETA=$T/c-absent.md" -- coverage
  t_after '### The unit and the platform' '- (0.1.5) A limit of an older record. Why: a case of the selftest. Record: deploy/BACKUP-RUN.md, *Not tested*: "BAK-01 PASS"' < "$T/good.md" > "$T/c-old.md"
  st_run 'coverage: a limit marked (0.1.5) names a record that is not of 0.1.5' 1 "coverage: a limit marked (0.1.5) names 'deploy/BACKUP-RUN.md', which is not a record of 0.1.5" "BETA=$T/c-old.md" -- coverage
  # ledger
  st_run 'ledger: control, the page as it is' 0 'ok: ledger' "BETA=$T/good.md" -- ledger
  t_after '### Install, update and rollback' "- $LIMIT_3: the update migrated cleanly. Why: kept." < "$T/good.md" > "$T/l-bullet.md"
  st_run 'ledger: a limit left as a bullet' 1 "ledger: still a bullet of the does-not-promise section: $LIMIT_3" "BETA=$T/l-bullet.md" -- ledger
  t_after '### The product itself' "$SCHEMA_OLD Why: kept." < "$T/good.md" > "$T/l-schema.md"
  st_run 'ledger: the schema bullet of 0.1.1 still there' 1 'ledger: the bullet of 0.1.1 on an older binary and a migrated database is still there' "BETA=$T/l-schema.md" -- ledger
  # nojs
  st_run 'nojs: control, the repository as it is' 0 'ok: nojs' -- nojs
  st_run 'nojs: no JavaScript allowed, so the htmx file is one too many' 1 'nojs: the tracked JavaScript is not exactly' "JS_ALLOWED=" -- nojs
  if [ "$ST_K" = 0 ]; then
    echo "SELFTEST OK ($ST_N cases)"
    return 0
  fi
  echo "SELFTEST FAILED ($ST_K of $ST_N)"
  return 1
}

# ----------------------------------------------------------------- main ----

usage() { echo "usage: 34-check.sh quotes|coverage|ledger|nojs ... | selftest" >&2; exit 2; }
[ "$#" -gt 0 ] || usage
if [ "$1" = selftest ]; then
  [ "$#" = 1 ] || usage
  selftest
  exit $?
fi

RULES=""
for a in "$@"; do
  case "$a" in
    quotes | coverage | ledger | nojs) RULES="$RULES $a" ;;
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
