#!/usr/bin/env bash
# Phase 34: MIG-01, MIG-02, UPD-01 and UPD-02 played in fresh systemd-base containers.
#
# A sibling of deploy/test/run.sh in the shape of deploy/test/backup-verify.sh. It
# edits neither run.sh, nightly.sh nor backup-verify.sh and reaches run.sh through
# its public subcommands only (build up install seed transfer check); what this
# file adds is the scenarios of phase 34, the files they write and the rule that
# reads a verdict from those files, never from a trace.
#
#   upgrade-verify.sh selftest       the pure functions below, a table of cases, no Docker
#   upgrade-verify.sh run            the evidence: the stages of STAGES, each in its own
#                                    fresh containers, the raw files below deploy/test/out
#   upgrade-verify.sh down           remove the four holzcloud-t34-* containers, nothing else
#   upgrade-verify.sh appendix STAGE the appendix blocks of the fixed paste list of the
#                                    record of STAGE (mig, upd or rb), from the files below E
#   upgrade-verify.sh record-check   the proof of the record R against the scratch root E:
#                                    every block byte for byte, every cited file pasted, one
#                                    Evidence commit, pinned and unchanged since
#
# Subcommands can be chained in one call and run in order.
#
# The bound files of the evidence set (BOUND_FILES below): deploy/DEPLOY.md,
# deploy/backup.sh, deploy/restore.sh, deploy/holzcloud.service, deploy/test/run.sh,
# deploy/test/Dockerfile, this file and the five plugin archives, together with the
# binary built from the evidence commit (its source: cmd, internal, go.mod, go.sum).
# The old release is bound by the commit sha of its tag, written into mig.txt
# (tags never move). `run` refuses to start while a bound file differs from HEAD or
# is untracked, and pins NEW_REF to the short sha of HEAD, so the binary and the
# scripts are one commit.
#
# Stages (STAGES, default "mig upd rb"; a subset is a rehearsal, and a stage left
# out reads NOT-RUN, so a subset can never read all PASS):
#   mig   v0.1.1 (migration level 60) installed by today's Install text, seeded with
#         v0.1.1's own plugin archives, backed up by v0.1.1's own backup.sh and
#         binary; the backup restored with the document's restore line into a
#         second container that runs this tree (level 69); the levels read from
#         the database files, the journal of the first start, the 17 checks of
#         run.sh check on the target (MIG-01, MIG-02)
#   upd   added by plan 34-04 (UPD-01)
#   rb    added by plan 34-04 (UPD-02)
#
# Environment:
#   PLATFORM    must be linux/amd64 (on this arm64 Mac it runs under emulation)
#   TARGET      must be systemd-base
#   STAGES      see above
#   REHEARSAL   1 lets `run` start while a bound file differs from HEAD; the summary
#               then says REHEARSAL, and such a set is never evidence
#   R           record-check: the record, one of deploy/MIGRATION-RESTORE-RUN.md,
#               deploy/UPDATE-CHAIN-RUN.md or deploy/ROLLBACK-RUN.md; its paste list
#               follows from its name
#   E           appendix and record-check: the scratch root, default deploy/test/out
#
# Names fixed for the phase: containers holzcloud-t34-mig-src, holzcloud-t34-mig-tgt,
# holzcloud-t34-upd and holzcloud-t34-rb; directories deploy/test/out/amd64-t34-mig-src,
# amd64-t34-mig-tgt, amd64-t34-upd and amd64-t34-rb; run files
# deploy/test/out/amd64-t34-evidence-commit.txt, amd64-t34-summary.txt and
# amd64-t34-pass1.{started,finished,status,log}; records deploy/MIGRATION-RESTORE-RUN.md
# (mig), deploy/UPDATE-CHAIN-RUN.md (upd) and deploy/ROLLBACK-RUN.md (rb).
#
# Exit codes: 0 ok, 1 failure (a verdict that is not PASS included), 2 usage,
# 3 the platform cannot give the evidence (no Docker daemon).
#
# Written for bash 3.2 (macOS): no arrays, no mapfile, no case conversion. GNU-only
# work runs inside the container through docker exec.

set -euo pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
cd "$ROOT"

PLATFORM=${PLATFORM:-}
TARGET=${TARGET:-}
STAGES=${STAGES:-mig upd rb}
REHEARSAL=${REHEARSAL:-}

OLD_TAG=v0.1.1
MIG_SRC=holzcloud-t34-mig-src
MIG_TGT=holzcloud-t34-mig-tgt
UPD_NAME=holzcloud-t34-upd
RB_NAME=holzcloud-t34-rb
BASE_OUT=deploy/test/out
MIG_SRC_OUT=$BASE_OUT/amd64-t34-mig-src
MIG_TGT_OUT=$BASE_OUT/amd64-t34-mig-tgt
UPD_OUT=$BASE_OUT/amd64-t34-upd
RB_OUT=$BASE_OUT/amd64-t34-rb
RUN_PREFIX=$BASE_OUT/amd64-t34-pass1
EVIDENCE_FILE=$BASE_OUT/amd64-t34-evidence-commit.txt
SUMMARY_FILE=$BASE_OUT/amd64-t34-summary.txt
BOUND_FILES="deploy/DEPLOY.md deploy/backup.sh deploy/restore.sh deploy/holzcloud.service deploy/test/run.sh deploy/test/Dockerfile deploy/test/upgrade-verify.sh plugins/bestellung/bestellung.zip plugins/jahreszahl/jahreszahl.zip plugins/kontaktformular/kontaktformular.zip plugins/nicht-gefunden/nicht-gefunden.zip plugins/suche/suche.zip"
DEPLOY_DOC=deploy/DEPLOY.md
# The document's manual backup line (DEPLOY.md, Backups), played by exact text.
MANUAL_LINE='sudo -u holzcloud /opt/holzcloud/deploy/backup.sh /opt/holzcloud/backups'
PLUGIN_IDS="bestellung jahreszahl kontaktformular nicht-gefunden suche"
NEW_REF=""
OLD_SHA=""

say() { printf '%s\n' "$*" >&2; }
die() {
  say "upgrade-verify.sh: $*"
  exit 1
}
usage() {
  awk 'NR > 1 && /^$/ { exit } NR > 1 { print }' "${BASH_SOURCE[0]}" >&2
  exit 2
}

# ---------------------------------------------------------------------------
# Pure functions. No Docker, no clock, no network: arguments or files in, one
# answer out. `selftest` holds a table of cases, the failing ones included.
# ---------------------------------------------------------------------------

# safe_name NAME: success for the four containers this file makes, and no other.
safe_name() {
  case "${1:-}" in
    holzcloud-t34-mig-src | holzcloud-t34-mig-tgt | holzcloud-t34-upd | holzcloud-t34-rb) return 0 ;;
  esac
  return 1
}

# safe_out PATH: success for the four directories this file writes and its run
# files, each named exactly; no other path, no path below them, no parent step.
safe_out() {
  case "${1:-}" in
    deploy/test/out/amd64-t34-mig-src | deploy/test/out/amd64-t34-mig-tgt | deploy/test/out/amd64-t34-upd \
      | deploy/test/out/amd64-t34-rb) return 0 ;;
    deploy/test/out/amd64-t34-evidence-commit.txt | deploy/test/out/amd64-t34-summary.txt) return 0 ;;
    deploy/test/out/amd64-t34-pass1.started | deploy/test/out/amd64-t34-pass1.finished \
      | deploy/test/out/amd64-t34-pass1.status | deploy/test/out/amd64-t34-pass1.log) return 0 ;;
  esac
  return 1
}

# stages_valid_of TEXT: TEXT is not empty and names only mig, upd and rb.
# stages_valid: the same for STAGES.
stages_valid_of() {
  local s n=0
  for s in ${1:-}; do
    case "$s" in
      mig | upd | rb) n=$((n + 1)) ;;
      *) return 1 ;;
    esac
  done
  [ "$n" -gt 0 ]
}
stages_valid() { stages_valid_of "$STAGES"; }

# stage_wanted NAME: NAME is one of the words of STAGES.
stage_wanted() {
  case " $STAGES " in
    *" $1 "*) return 0 ;;
  esac
  return 1
}

# is_rehearsal: REHEARSAL is set, or STAGES leaves out one of the three stages.
is_rehearsal() {
  [ -n "$REHEARSAL" ] && return 0
  stage_wanted mig && stage_wanted upd && stage_wanted rb && return 1
  return 0
}

# level_of_output TEXT: the one decimal number sqlite3 printed (a trailing carriage
# return dropped); fails on empty output, text, two lines or a sign.
level_of_output() {
  local t
  t=$(printf '%s' "${1:-}" | tr -d '\r')
  case "$t" in
    "" | *[!0-9]*) return 1 ;;
  esac
  printf '%s\n' "$t"
}

# kv FILE KEY: the value of the one line "KEY: value"; empty when absent. kv_count
# FILE KEY: how many such lines there are.
kv() { sed -n "s/^$2: //p" "$1" 2> /dev/null | head -n 1; }
kv_count() { grep -c "^$2: " "$1" 2> /dev/null || true; }

# is_count TEXT: success for a decimal number without a sign.
is_count() {
  case "${1:-}" in
    "" | *[!0-9]*) return 1 ;;
  esac
  return 0
}

# kv_is FILE KEY VALUE: the key occurs exactly once and has exactly that value.
kv_is() {
  [ "$(kv_count "$1" "$2")" = 1 ] && [ "$(kv "$1" "$2")" = "$3" ]
}
# kv_num FILE KEY: the key occurs exactly once and is a number.
kv_num() {
  [ "$(kv_count "$1" "$2")" = 1 ] && is_count "$(kv "$1" "$2")"
}
# kv_same FILE KEY1 KEY2: both keys occur once, are numbers and are equal.
kv_same() {
  kv_num "$1" "$2" && kv_num "$1" "$3" && [ "$(kv "$1" "$2")" -eq "$(kv "$1" "$3")" ]
}

# verdict_mig01 FILE: PASS, FAIL or NOT-RUN (no file) for the key lines of mig.txt.
verdict_mig01() {
  local f=${1:-}
  if [ ! -f "$f" ]; then
    echo NOT-RUN
    return 0
  fi
  if kv_is "$f" RESTORE-RC 0 && kv_same "$f" LEVEL-BEFORE LEVEL-EXPECTED-BEFORE \
    && kv_same "$f" LEVEL-AFTER LEVEL-EXPECTED-AFTER; then
    echo PASS
  else
    echo FAIL
  fi
}

# ---------------------------------------------------------------------------
# The containers. Every docker call names a container that safe_name has
# accepted, and every run.sh call a directory that safe_out has accepted.
# ---------------------------------------------------------------------------

# dxc NAME CMD...: docker exec in one of the four containers.
dxc() {
  local n=$1
  shift
  safe_name "$n" || die "refusing container name $n"
  docker exec "$n" "$@"
}

# run_sh NAME OUT SUBCOMMAND...: run.sh through its public subcommands, on that
# container and that directory only, the new binary of NEW_REF. run_sh_old: the
# same with INSTALL_FROM=old and OLD_TAG (run.sh builds the tag from a clone with
# today's Build line and plays today's Install text).
run_sh() {
  local n=$1 o=$2
  shift 2
  safe_name "$n" || die "refusing container name $n"
  safe_out "$o" || die "refusing path $o"
  env PLATFORM="$PLATFORM" TARGET=systemd-base NEW_REF="$NEW_REF" NAME="$n" OUT_DIR="$o" KEEP=1 \
    INSTALL_FROM=new bash "$ROOT/deploy/test/run.sh" "$@"
}
run_sh_old() {
  local n=$1 o=$2
  shift 2
  safe_name "$n" || die "refusing container name $n"
  safe_out "$o" || die "refusing path $o"
  env PLATFORM="$PLATFORM" TARGET=systemd-base NEW_REF="$NEW_REF" NAME="$n" OUT_DIR="$o" KEEP=1 \
    INSTALL_FROM=old OLD_TAG="$OLD_TAG" bash "$ROOT/deploy/test/run.sh" "$@"
}

# run_sh_transfer SRC SRC_OUT TGT TGT_OUT DIR: run.sh transfer of the backup
# directory DIR of SRC into TGT (with the three fixtures run.sh carries along).
run_sh_transfer() {
  local n
  for n in "$1" "$3"; do safe_name "$n" || die "refusing container name $n"; done
  for n in "$2" "$4"; do safe_out "$n" || die "refusing path $n"; done
  env PLATFORM="$PLATFORM" TARGET=systemd-base NEW_REF="$NEW_REF" NAME="$1" OUT_DIR="$2" \
    TGT_NAME="$3" TGT_OUT_DIR="$4" BACKUP_DIR_PATH="$5" KEEP=1 bash "$ROOT/deploy/test/run.sh" transfer
}

# DEV_FILE: where deviation() writes; set by each stage to its own directory.
DEV_FILE=""
# deviation TEXT: one line of the driver's deviations of the current stage.
deviation() {
  say "$*"
  [ -n "$DEV_FILE" ] || die "deviation before a stage directory exists: $*"
  printf '%s\n' "$*" >> "$DEV_FILE"
}

# put FILE KEY VALUE: one key line, also on stderr for the log.
put() {
  printf '%s: %s\n' "$2" "$3" >> "$1"
  say "$2: $3"
}

# bound_check: the bound files are tracked and equal to HEAD. Prints one line per
# problem and fails when there is one.
bound_check() {
  local f bad=0
  for f in $BOUND_FILES; do
    if ! git ls-files --error-unmatch -- "$f" > /dev/null 2>&1; then
      say "upgrade-verify.sh: bound file $f is not tracked"
      bad=1
    elif ! git diff --quiet HEAD -- "$f"; then
      say "upgrade-verify.sh: bound file $f differs from HEAD"
      bad=1
    fi
  done
  [ "$bad" -eq 0 ]
}

# Run INSIDE the container. STAMP_SH: $1 is the line to play, stdin the answer.
# Every output line is prefixed with the container's clock (run.sh's
# RESTORE_RUN_SH technique), a last line "WALL t0 t1", and the exit status is the
# line's own (PIPESTATUS), not that of the loop.
STAMP_SH='set -o pipefail
t0=$(date +%s.%N)
bash -c "$1" 2>&1 | while IFS= read -r l; do printf "%s %s\n" "$(date +%s.%N)" "$l"; done
rc=${PIPESTATUS[0]}
t1=$(date +%s.%N)
echo "WALL $t0 $t1"
exit "$rc"'

# unstamped FILE: the lines of a stamped file without the clock.
unstamped() { sed 's/^[0-9][0-9.]* //' "$1"; }

# restore_line DIR: the document's restore line with DIR in place of its example
# (DEVIATION[backup-name]).
restore_line() {
  local l
  l=$(grep -E '^sudo /opt/holzcloud/deploy/restore\.sh /opt/holzcloud/backups/[0-9]{8}-[0-9]{6}$' "$DEPLOY_DOC" | head -n 1) || l=""
  [ -n "$l" ] || die "the document and this file have drifted apart: $DEPLOY_DOC has no restore line of the shape 'sudo /opt/holzcloud/deploy/restore.sh /opt/holzcloud/backups/TIMESTAMP'"
  printf '%s %s\n' "${l% *}" "$1"
}

# LEVEL_SQL: the highest applied goose version, the one definition of a level.
LEVEL_SQL="select coalesce(max(version_id), 0) from goose_db_version where is_applied = 1"

# level_file NAME PATH: the level of a database file that no process writes (a
# backup, a snapshot), read through the immutable URI as the service user, so
# that no -wal or -shm file appears beside it. Never a holzcloud subcommand: every
# one of them migrates, and from 0.1.5 snapshots, before it answers. "?" when the
# answer is not one number.
level_file() {
  local out
  out=$(dxc "$1" runuser -u holzcloud -- sqlite3 "file:$2?immutable=1" "$LEVEL_SQL" 2> /dev/null) || out=""
  level_of_output "$out" || echo "?"
}

# level_live NAME: the level of the live database of the running service. The
# immutable URI would ignore the -wal file, which holds what the service wrote
# since its last checkpoint (the migrations of the first start among it), so the
# live file is read through a read-only connection that sees the log; still
# sqlite3, never a holzcloud subcommand.
level_live() {
  local out
  out=$(dxc "$1" runuser -u holzcloud -- sqlite3 -readonly /opt/holzcloud/data/holzcloud.sqlite ".timeout 5000" "$LEVEL_SQL" 2> /dev/null) || out=""
  level_of_output "$out" || echo "?"
}

# migrations_at REF: the number of migration files at REF (git ls-tree), the level
# a binary of REF migrates to. Never a typed number.
migrations_at() {
  git ls-tree --name-only "$1" internal/db/migrations/ | grep -c '\.sql$' || true
}

# ---------------------------------------------------------------------------
# Stage mig: a v0.1.1 backup, made by its own script and binary, restored into
# this tree across the migrations above its level.
# ---------------------------------------------------------------------------

# mig_source: v0.1.1 installed by today's Install text, seeded with its own plugin
# archives, a login user made by its own CLI.
mig_source() {
  local f=$1 id zip sha ver
  # `up` resets the directory, so the stage's own files begin after it.
  run_sh_old "$MIG_SRC" "$MIG_SRC_OUT" build up install
  : > "$f"
  DEV_FILE=$MIG_SRC_OUT/deviations-driver.lst
  : > "$DEV_FILE"
  ver=$(dxc "$MIG_SRC" /opt/holzcloud/holzcloud version | head -n 1 | tr -d '\r') || ver=""
  case "$ver" in
    *"$OLD_TAG"*) ;;
    *) die "mig: the source runs '$ver', not $OLD_TAG" ;;
  esac
  put "$f" OLD-RELEASE "$OLD_TAG ($(git rev-parse --short "$OLD_SHA"))"
  put "$f" OLD-COMMIT "$OLD_SHA"
  put "$f" OLD-VERSION-LINE "$ver"
  put "$f" NEW-RELEASE "$(git describe --tags --always "$NEW_REF") ($NEW_REF)"

  # v0.1.1's own plugin archives (F-4): its manifest check knows no `privacy`
  # hook, which today's kontaktformular and bestellung declare.
  dxc "$MIG_SRC" mkdir -p /opt/holzcloud/plugin-zips || die "mig: could not create /opt/holzcloud/plugin-zips"
  for id in $PLUGIN_IDS; do
    zip=$MIG_SRC_OUT/v011-$id.zip
    git show "$OLD_SHA:plugins/$id/$id.zip" > "$zip" || die "mig: git show $OLD_TAG:plugins/$id/$id.zip failed"
    [ -s "$zip" ] || die "mig: $OLD_TAG has an empty plugins/$id/$id.zip"
    sha=$(shasum -a 256 "$zip" | awk '{ print $1 }')
    docker cp "$zip" "$MIG_SRC:/opt/holzcloud/plugin-zips/$id.zip" > /dev/null || die "mig: docker cp of $zip failed"
    put "$f" OLD-PLUGIN "$id sha256 $sha"
  done
  deviation "DEVIATION[old-plugins]: the v0.1.1 source is seeded with the plugin archives of v0.1.1 itself (git show v0.1.1:plugins/ID/ID.zip), not with this tree's: v0.1.1 refuses today's kontaktformular and bestellung archives, whose manifests declare the hook 'privacy' that its manifest check does not know (finding F-4). After the restore the host of this tree runs those v0.1.1 modules"
  deviation "DEVIATION[fixture]: the five plugin archives are copied into /opt/holzcloud/plugin-zips before run.sh seed, and the login user is made with the CLI of v0.1.1 (password in a root-only file, given on stdin); with the website, page, AI key and test image run.sh seed makes, they stand in for operator work in the admin and are no step of the document. No order row is written (ORDER: none)"
  local rc=0
  run_sh "$MIG_SRC" "$MIG_SRC_OUT" seed || rc=$?
  put "$f" SEED-RC "$rc"
  [ "$rc" -eq 0 ] || die "mig: run.sh seed failed on the $OLD_TAG source"

  dxc "$MIG_SRC" bash -c 'umask 077; openssl rand -hex 12 > /root/.hc-login-pass; test -s /root/.hc-login-pass' \
    || die "mig: could not create the login password file"
  dxc "$MIG_SRC" bash -c 'runuser -u holzcloud -- env HOLZCLOUD_DATA_DIR=/opt/holzcloud/data /opt/holzcloud/holzcloud user create -email login@example.test -role admin < /root/.hc-login-pass' \
    > /dev/null || die "mig: holzcloud user create failed on the $OLD_TAG binary"
  {
    echo "OLD_TAG: $OLD_TAG"
    echo "ORDER: none (the stage writes no order row; run.sh transfer reads this line)"
  } > "$MIG_SRC_OUT/baseline.txt"
}

# mig_backup: v0.1.1's backup.sh in place of today's, the document's manual line
# played with v0.1.1's binary; MIG_DIR is the directory it names.
mig_backup() {
  local f=$1 sha rc=0
  git show "$OLD_SHA:deploy/backup.sh" > "$MIG_SRC_OUT/v011-backup.sh" || die "mig: git show $OLD_TAG:deploy/backup.sh failed"
  docker cp "$MIG_SRC_OUT/v011-backup.sh" "$MIG_SRC:/opt/holzcloud/deploy/backup.sh" > /dev/null \
    || die "mig: docker cp of the $OLD_TAG backup.sh failed"
  dxc "$MIG_SRC" bash -c 'chown root:root /opt/holzcloud/deploy/backup.sh && chmod 755 /opt/holzcloud/deploy/backup.sh' \
    || die "mig: could not set the owner of the $OLD_TAG backup.sh"
  sha=$(dxc "$MIG_SRC" sha256sum /opt/holzcloud/deploy/backup.sh | cut -d' ' -f1 | tr -d '\r')
  put "$f" BACKUP-SCRIPT "$OLD_TAG sha256 $sha"
  deviation "DEVIATION[old-script]: /opt/holzcloud/deploy/backup.sh is replaced by the backup.sh of v0.1.1 (git show v0.1.1:deploy/backup.sh, root 755) before the manual line is played, because today's Install text copied this tree's; a server still on v0.1.1 has that script there from its own install"
  grep -Fxq -- "$MANUAL_LINE" "$DEPLOY_DOC" || die "the document and this file have drifted apart: $DEPLOY_DOC has no line '$MANUAL_LINE'"
  say "mig: $MANUAL_LINE"
  docker exec "$MIG_SRC" bash -c "$STAMP_SH" _ "$MANUAL_LINE" < /dev/null > "$MIG_SRC_OUT/mig-backup.out" 2>&1 || rc=$?
  put "$f" BACKUP-RC "$rc"
  MIG_DIR=$(unstamped "$MIG_SRC_OUT/mig-backup.out" | sed -n 's/^=== Backup complete: \(.*\) ===$/\1/p' | tail -n 1)
  case "$MIG_DIR" in
    /opt/holzcloud/backups/[0-9]*-[0-9]*) ;;
    *) die "mig: the manual backup named no directory ('$MIG_DIR')" ;;
  esac
  put "$f" BACKUP-DIR "$MIG_DIR"
  put "$f" MANIFEST-VERSION "$(dxc "$MIG_SRC" sed -n 's/^version=//p' "$MIG_DIR/MANIFEST" | head -n 1 | tr -d '\r')"
  if dxc "$MIG_SRC" test -e "$MIG_DIR/SHA256SUMS"; then
    put "$f" SUMS-IN-BACKUP yes
  else
    put "$f" SUMS-IN-BACKUP no
  fi
  put "$f" LEVEL-BEFORE "$(level_file "$MIG_SRC" "$MIG_DIR/holzcloud.sqlite")"
  put "$f" LEVEL-EXPECTED-BEFORE "$(migrations_at "$OLD_SHA")"
}

# mig_restore: a fresh target with this tree, the backup carried over, the
# document's restore line played on it.
mig_restore() {
  local f=$1 rc=0 line base
  run_sh "$MIG_TGT" "$MIG_TGT_OUT" build up install
  run_sh_transfer "$MIG_SRC" "$MIG_SRC_OUT" "$MIG_TGT" "$MIG_TGT_OUT" "$MIG_DIR" || rc=$?
  put "$f" TRANSFER-RC "$rc"
  [ "$rc" -eq 0 ] || die "mig: run.sh transfer failed"
  base=$(basename "$MIG_DIR")
  line=$(restore_line "/opt/holzcloud/backups/$base")
  deviation "DEVIATION[backup-name]: the sample timestamp of the document's restore line is replaced by the transferred directory /opt/holzcloud/backups/$base; the answer 'yes' is typed on stdin"
  deviation "DEVIATION[start-limit]: systemctl reset-failed holzcloud runs before the restore and before run.sh check; the install start, the restore's start and check's restart come closer together than on a server"
  dxc "$MIG_TGT" systemctl reset-failed holzcloud > /dev/null 2>&1 || true
  say "mig: $line"
  rc=0
  printf 'yes\n' | docker exec -i "$MIG_TGT" bash -c "$STAMP_SH" _ "$line" > "$MIG_SRC_OUT/mig-restore.out" 2>&1 || rc=$?
  put "$f" RESTORE-RC "$rc"
  put "$f" READY-SECONDS "$(sed -n 's/^READY_SECONDS="\${READY_SECONDS:-\([0-9]*\)}"$/\1/p' deploy/restore.sh | head -n 1)"
  put "$f" FIRST-START-S "$(awk '/ Starting holzcloud\.\.\.$/ { s = $1 } / === Restore complete ===$/ { c = $1 } END { if (s != "" && c != "") printf "%.1f\n", c - s; else print "?" }' "$MIG_SRC_OUT/mig-restore.out")"
  put "$f" LEVEL-AFTER "$(level_live "$MIG_TGT")"
  put "$f" LEVEL-EXPECTED-AFTER "$(migrations_at "$NEW_REF")"
}

stage_mig() {
  local f=$MIG_SRC_OUT/mig.txt
  MIG_DIR=""
  mig_source "$f"
  mig_backup "$f"
  mig_restore "$f"
  printf 'MIG-01 %s\n' "$(verdict_mig01 "$f")" >> "$f"
}

# mig_deviations: deviations.txt of the source directory, the driver's own lines
# first, then those run.sh wrote for the source and for the target; a line that
# occurs twice (the same install played twice) is kept once.
mig_deviations() {
  local d=$MIG_SRC_OUT arch
  arch=$(docker info --format '{{.Architecture}}')
  case "$arch" in
    x86_64 | amd64) ;;
    *) deviation "DEVIATION[emulation]: the Docker host is $arch; the linux/amd64 containers run under emulation, so every wall time of this set (the first start after the restore among them) is an emulation figure, not a server's" ;;
  esac
  [ ! -f "$d/deviations.txt" ] || mv "$d/deviations.txt" "$d/deviations-runsh-src.lst"
  cat "$d/deviations-driver.lst" "$d/deviations-runsh-src.lst" "$MIG_TGT_OUT/deviations.txt" 2> /dev/null \
    | awk '/^DEVIATION\[/ && !seen[$0]++' > "$d/deviations.txt"
}

# mig_versions: run.sh's versions.txt of the source (written by `up`), then both
# releases and the tools of this set.
mig_versions() {
  local f=$MIG_SRC_OUT/versions.txt
  {
    echo "host uname -m:         $(uname -m)"
    echo "container uname -m:    $(dxc "$MIG_SRC" uname -m | tr -d '\r')"
    echo "container systemctl:   $(dxc "$MIG_TGT" systemctl --version | head -n 1 | tr -d '\r')"
    echo "source binary:         $(dxc "$MIG_SRC" /opt/holzcloud/holzcloud version | head -n 1 | tr -d '\r')"
    echo "target binary:         $(dxc "$MIG_TGT" /opt/holzcloud/holzcloud version | head -n 1 | tr -d '\r')"
    echo "$OLD_TAG commit:         $OLD_SHA"
    echo "NEW_REF:               $NEW_REF"
    echo "evidence commit:       $(cat "$EVIDENCE_FILE")"
  } >> "$f"
}

# head_binary: the version line of the HEAD build, from a container that runs it.
head_binary() {
  local n v
  for n in "$MIG_TGT" "$UPD_NAME" "$RB_NAME"; do
    if docker container inspect "$n" > /dev/null 2>&1; then
      v=$(dxc "$n" /opt/holzcloud/holzcloud version 2> /dev/null | head -n 1 | tr -d '\r') || v=""
      if [ -n "$v" ]; then
        printf '%s\n' "$v"
        return 0
      fi
    fi
  done
  echo "?"
}

# write_summary: the verdicts from the files the stages wrote; a stage that did
# not run reads NOT-RUN whatever lies in its directory.
write_summary() {
  local f=$SUMMARY_FILE kinds m1=NOT-RUN m2=NOT-RUN u1=NOT-RUN u2=NOT-RUN d
  if stage_wanted mig; then
    m1=$(verdict_mig01 "$MIG_SRC_OUT/mig.txt")
    m2=NOT-RUN
  fi
  kinds=$(for d in "$MIG_SRC_OUT" "$UPD_OUT" "$RB_OUT"; do
    [ ! -f "$d/deviations.txt" ] || cat "$d/deviations.txt"
  done | sed -n 's/^DEVIATION\[\([a-z0-9-]*\)\].*/\1/p' | sort -u | tr '\n' ' ' | sed 's/ $//')
  {
    echo "MIG-01 $m1"
    echo "MIG-02 $m2"
    echo "UPD-01 $u1"
    echo "UPD-02 $u2"
    echo "STAGES: $STAGES"
    echo "DEVIATIONS: ${kinds:-none}"
    echo "BINARY: $(head_binary)"
    if is_rehearsal; then
      echo "REHEARSAL: yes"
    fi
  } > "$f"
  cat "$f" >&2
}

# run_body: everything `run` does after its guards; its output goes to the log.
run_body() {
  say "upgrade-verify: platform $PLATFORM, NEW_REF pinned to $NEW_REF, $OLD_TAG is $OLD_SHA, stages: $STAGES"
  if stage_wanted mig; then
    stage_mig
    mig_versions
    mig_deviations
  fi
  write_summary
  grep -Eq '^(MIG|UPD)-0[0-9] (FAIL|NOT-RUN)$' "$SUMMARY_FILE" && return 1
  return 0
}

cmd_run() {
  local rc p n
  stages_valid || {
    say "upgrade-verify.sh: STAGES '$STAGES' must name mig, upd or rb"
    exit 2
  }
  for p in upd rb; do
    if stage_wanted "$p"; then
      say "upgrade-verify.sh: the stage $p is added by plan 34-04; run STAGES=mig until then"
      exit 2
    fi
  done
  [ "$PLATFORM" = linux/amd64 ] || {
    say "upgrade-verify.sh: PLATFORM must be linux/amd64 (got '$PLATFORM')"
    exit 2
  }
  [ "$TARGET" = systemd-base ] || {
    say "upgrade-verify.sh: TARGET must be systemd-base (got '$TARGET')"
    exit 2
  }
  for p in "$MIG_SRC_OUT" "$MIG_TGT_OUT" "$UPD_OUT" "$RB_OUT" "$RUN_PREFIX.started" "$RUN_PREFIX.finished" \
    "$RUN_PREFIX.status" "$RUN_PREFIX.log" "$EVIDENCE_FILE" "$SUMMARY_FILE"; do
    safe_out "$p" || die "refusing path $p"
  done
  if ! bound_check; then
    if [ -n "$REHEARSAL" ]; then
      say "upgrade-verify.sh: REHEARSAL=1: going on although a bound file differs from HEAD; this set is not evidence"
    else
      die "a bound file differs from HEAD or is untracked; commit first (the evidence is bound to one commit)"
    fi
  fi
  command -v docker > /dev/null 2>&1 || {
    say "upgrade-verify.sh: docker is not installed"
    exit 3
  }
  docker info > /dev/null 2>&1 || {
    say "upgrade-verify.sh: the Docker daemon does not answer"
    exit 3
  }
  for n in "$MIG_SRC" "$MIG_TGT" "$UPD_NAME" "$RB_NAME"; do
    safe_name "$n" || die "refusing container name $n"
    if docker container inspect "$n" > /dev/null 2>&1; then
      die "container $n exists; run 'upgrade-verify.sh down' first (a stale one is never reused)"
    fi
  done
  NEW_REF=$(git rev-parse --short HEAD) || die "HEAD is not a commit"
  OLD_SHA=$(git rev-parse --verify --quiet "$OLD_TAG^{commit}") || die "the tag $OLD_TAG is not in this repository"
  mkdir -p "$BASE_OUT"
  # The directories of this file only: a stale stage directory never lends a verdict.
  for p in "$MIG_SRC_OUT" "$MIG_TGT_OUT" "$UPD_OUT" "$RB_OUT"; do
    rm -rf "${p:?}"
  done
  rm -f "$RUN_PREFIX.started" "$RUN_PREFIX.finished" "$RUN_PREFIX.status" "$RUN_PREFIX.log" "$SUMMARY_FILE"
  printf '%s\n' "$NEW_REF" > "$EVIDENCE_FILE"
  date -u +%Y-%m-%dT%H:%M:%SZ > "$RUN_PREFIX.started"
  set +e
  (set -e; run_body) 2>&1 | tee "$RUN_PREFIX.log"
  rc=${PIPESTATUS[0]}
  set -e
  date -u +%Y-%m-%dT%H:%M:%SZ > "$RUN_PREFIX.finished"
  printf '%s\n' "$rc" > "$RUN_PREFIX.status"
  say "upgrade-verify: run ended with status $rc"
  exit "$rc"
}

cmd_down() {
  local n
  command -v docker > /dev/null 2>&1 || die "docker is not installed"
  for n in "$MIG_SRC" "$MIG_TGT" "$UPD_NAME" "$RB_NAME"; do
    safe_name "$n" || die "refusing container name $n"
    docker rm -f "$n" > /dev/null 2>&1 || true
    say "upgrade-verify: container $n removed"
  done
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
# st_ans NAME EXPECTED cmd args...: the command's stdout (an error reads ERROR).
st_ans() {
  local n=$1 want=$2 got
  shift 2
  got=$("$@" 2> /dev/null) || got=ERROR
  st_eq "$n" "$want" "$got"
}

selftest_names() {
  local n
  for n in holzcloud-t34-mig-src holzcloud-t34-mig-tgt holzcloud-t34-upd holzcloud-t34-rb; do
    st_ok "safe_name: $n" safe_name "$n"
  done
  for n in holzcloud-t32 holzcloud-t34 holzcloud-t34-mig holzcloud-t34-mig-srcx holzcloud-t34-rb2 \
    holzcloud-unit-test-amd64 holzcloud-t34-upd/x ""; do
    st_fail "safe_name: '$n'" safe_name "$n"
  done
  for n in amd64-t34-mig-src amd64-t34-mig-tgt amd64-t34-upd amd64-t34-rb; do
    st_ok "safe_out: $n" safe_out "deploy/test/out/$n"
  done
  for n in amd64-t34-evidence-commit.txt amd64-t34-summary.txt amd64-t34-pass1.started amd64-t34-pass1.finished \
    amd64-t34-pass1.status amd64-t34-pass1.log; do
    st_ok "safe_out: the run file $n" safe_out "deploy/test/out/$n"
  done
  for n in deploy/test/out/amd64-t32-src deploy/test/out/amd64-t34- deploy/test/out/amd64-t34-mig \
    deploy/test/out/amd64-t34-mig-src/x deploy/test/out/amd64-t34-../x deploy/test/out/amd64-t34-other \
    deploy/test/out /deploy/test/out/amd64-t34-upd ./deploy/test/out/amd64-t34-rb . ""; do
    st_fail "safe_out: '$n'" safe_out "$n"
  done
}

selftest_stages() {
  local s
  for s in "mig upd rb" "mig" "upd" "rb" "mig rb" "rb mig"; do
    st_ok "stages_valid: '$s'" stages_valid_of "$s"
  done
  for s in "" "mig bak01" "MIG" "migration" "mig upd rb all"; do
    st_fail "stages_valid: '$s'" stages_valid_of "$s"
  done
}

selftest_level() {
  st_ans "level_of_output: 60" 60 level_of_output 60
  st_ans "level_of_output: 69 with a carriage return" 69 level_of_output "$(printf '69\r')"
  st_ans "level_of_output: empty" ERROR level_of_output ""
  st_ans "level_of_output: text" ERROR level_of_output "Error: unable to open database"
  st_ans "level_of_output: two lines" ERROR level_of_output "$(printf '60\n61')"
  st_ans "level_of_output: a sign" ERROR level_of_output "-1"
  st_ans "level_of_output: a space inside" ERROR level_of_output "6 0"
}

# fx_set FILE KEY VALUE: the line "KEY: ..." of FILE replaced by "KEY: VALUE"; with
# VALUE "-" the line is removed.
fx_set() {
  awk -v k="$2" -v v="$3" 'index($0, k ": ") == 1 { if (v != "-") print k ": " v; next } { print }' "$1" > "$1.x"
  mv "$1.x" "$1"
}

# fx_mig FILE [KEY VALUE]: every key of mig.txt, all good; with KEY and VALUE, that
# one line replaced (or removed for "-").
fx_mig() {
  local f=$1
  {
    echo "OLD-RELEASE: v0.1.1 (5bcae9b)"
    echo "NEW-RELEASE: v0.1.4-60-gcfce101 (cfce101)"
    echo "RESTORE-RC: 0"
    echo "LEVEL-BEFORE: 60"
    echo "LEVEL-EXPECTED-BEFORE: 60"
    echo "LEVEL-AFTER: 69"
    echo "LEVEL-EXPECTED-AFTER: 69"
  } > "$f"
  if [ -n "${2:-}" ]; then fx_set "$f" "$2" "${3:-}"; fi
}

selftest_verdict_mig01() {
  local f=$ST_TMP/mig.txt kv
  st_ans "verdict_mig01: no file" NOT-RUN verdict_mig01 "$ST_TMP/absent.txt"
  fx_mig "$f"
  st_ans "verdict_mig01: every key good" PASS verdict_mig01 "$f"
  for kv in "RESTORE-RC|1" "RESTORE-RC|3" "RESTORE-RC|-" "LEVEL-BEFORE|61" "LEVEL-BEFORE|?" "LEVEL-BEFORE|-" \
    "LEVEL-EXPECTED-BEFORE|-" "LEVEL-AFTER|60" "LEVEL-AFTER|68" "LEVEL-AFTER|-" "LEVEL-EXPECTED-AFTER|70" \
    "LEVEL-EXPECTED-AFTER|-"; do
    fx_mig "$f" "${kv%%|*}" "${kv#*|}"
    st_ans "verdict_mig01: ${kv%%|*} '${kv#*|}'" FAIL verdict_mig01 "$f"
  done
  fx_mig "$f"
  echo "RESTORE-RC: 0" >> "$f"
  st_ans "verdict_mig01: a key twice" FAIL verdict_mig01 "$f"
  fx_mig "$f" LEVEL-BEFORE ""
  fx_set "$f" LEVEL-EXPECTED-BEFORE ""
  st_ans "verdict_mig01: two empty levels are no level" FAIL verdict_mig01 "$f"
}

cmd_selftest() {
  ST_N=0
  ST_BAD=0
  ST_TMP=$(mktemp -d "${TMPDIR:-/tmp}/upgrade-verify-selftest.XXXXXX")
  trap 'rm -rf "${ST_TMP:?}"' EXIT
  selftest_names
  selftest_stages
  selftest_level
  selftest_verdict_mig01
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
    down) cmd_down ;;
    *) usage ;;
  esac
done
