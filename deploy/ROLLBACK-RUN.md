# A migration that fails part-way, and the rollback by snapshot: results

This document is the permanent record of UPD-02 of phase 34. It records an update from the build of this release (B0, migration level 69) to a build whose first added migration succeeds and whose second fails, played with the *Updating* blocks of `deploy/DEPLOY.md`: that the service stops with the failing migration named in the journal, that the restart loop writes one pre-upgrade snapshot before the first attempt and keeps it on every restart, that every `restore with:` line names that one file, that the update's own check line names it too, that the previous binary refuses the half-migrated database with the schema guard's sentence and both levels, and that *Rolling back* block 1, played as written, puts the previous binary and the snapshot back and starts the previous binary on the state it had before the update, row for row.

The record rests on one evidence set: a single `upgrade-verify.sh run` from fresh containers, made once against the committed text of every bound file (see *How this was run*); this record uses its stage `rb`, the same set as `deploy/MIGRATION-RESTORE-RUN.md`. Its raw output is the Appendix, pasted unchanged, and `bash deploy/test/upgrade-verify.sh record-check` proves every block against its scratch file byte for byte. Every number in this document is a line of that set and cites its file in backticks as `DIRECTORY/FILE`, with `DIRECTORY` the one scratch directory `amd64-t34-rb`; the run files beside it are named by their path below `deploy/test/out`. A file that is not pasted is named as `FILE in DIRECTORY`.

The record says what was done, how many times, and what was seen. The failing build is a throwaway build of this tree, not a release, and the failure is a migration that names a table that does not exist; other failures are named in *Not tested*. What the container cannot show is a named deviation or a line of *Not tested*, never a gap in the table.

## Summary

| Requirement | Result | Platform | Date, commit |
|-------------|--------|----------|--------------|
| UPD-02 | PASS: `UPD-02 PASS` stands in `amd64-t34-summary.txt` and in `amd64-t34-rb/rb.txt`. The update to F stopped with `STOP-MESSAGE: goose up: partial migration error (type:sql,version:71): SQL logic error: no such table: t34_no_such_table (1)` after `APPLIED-BEFORE-FAILURE: 70`; the restart loop ran to its end without a reset (`STARTS: 5`, `ACTIVE-STATE: failed`, `START-LIMIT-LINE: yes`, `HEALTHZ-DURING-LOOP: never`) and wrote one snapshot at level 69 that every later start kept (`JOURNAL-SNAPSHOT-WRITTEN: 1`, `JOURNAL-SNAPSHOT-KEPT: 4`, `WRITTEN-LEVEL: 69`), named by all five `restore with:` lines (`JOURNAL-RESTORE-WITH: 5`, `RESTORE-WITH-DISTINCT-PATHS: 1`, `RESTORE-WITH-NAMES-WRITTEN: yes`) and by the update's check line (`UPDATE-CHECK-SNAPSHOT-LINE: yes`); the previous binary refused the database at level 70 (`GUARD-LINE: yes`, `GUARD-DB-LEVEL: 70`, `GUARD-BINARY-LEVEL: 69`, `GUARD-SENTENCE: yes`, `GUARD-HEALTHZ: never`); *Rolling back* block 1 ended with every line at status 0 (`ROLLBACK-RC: 0 0 0 0 0 0 0 0 0 0 0`), the previous build answered `/readyz` with `status` and `integrity` `ok` at `LEVEL-AFTER-ROLLBACK: 69`, and the state after equals the state before (`COMPARE-S0-S2: bad=0`), all in `amd64-t34-rb/rb.txt`. | Debian 12 container, systemd PID 1, linux/amd64 under emulation on an arm64 host (`amd64-t34-rb/versions.txt`) | 2026-11-21, 00854be |

Evidence commit: 00854be

The previous binary, B0, is the one built from `NEW_REF` (pinned by `run` to the short sha of HEAD, `00854be`, read from `deploy/test/out/amd64-t34-evidence-commit.txt`) with the Build line of `deploy/DEPLOY.md`: `BINARY: holzcloud v0.1.4-255-g00854be (00854be)` in `amd64-t34-summary.txt`, `B0-VERSION: v0.1.4-255-g00854be` in `amd64-t34-rb/rb.txt`, and the line `holzcloud v0.1.4-255-g00854be (00854be)` that `/opt/holzcloud/holzcloud version` printed in *Updating* block 1 (`amd64-t34-rb/rb-update.out`). The failing build, F, is `F-BUILD: v0.1.4-256-g4c07745 (4c07745) migrations 71 sha256 af493d4a045ef7f0f54c4cb9021a9b561a8b0492828b0d9e8bd8389e1bb3b3d9` (`amd64-t34-rb/rb.txt`, the same line as `holzcloud-f:` in `amd64-t34-rb/versions.txt`): a throwaway build of this tree with one local commit that adds two migrations, made in a scratch clone and deleted with it (`DEVIATION[synthetic-build]`), so its commit `4c07745` exists nowhere now. The levels are computed from the migration files, never typed: `LEVEL-EXPECTED-B0: 69` and `LEVEL-EXPECTED-F: 71`. The summary carries `STAGES: mig upd rb` and no `REHEARSAL` line.

## How this was run

`deploy/test/upgrade-verify.sh` is a sibling of `deploy/test/run.sh` in the shape of `deploy/test/backup-verify.sh`; it edits neither and reaches the harness through `run.sh`'s public subcommands only. `run` refuses to start unless `PLATFORM=linux/amd64` and `TARGET=systemd-base`, refuses while a bound file or the binary's source differs from HEAD or is untracked, pins `NEW_REF` to the short sha of HEAD so that the binary and the scripts are one commit, and plays the stages `mig`, `upd` and `rb`, each in its own fresh containers. Every line the stage `rb` plays from *Updating* and *Rolling back* is read out of the fenced bash blocks of `deploy/DEPLOY.md` at that commit; a line the driver cannot translate ends the stage. The verdict is computed from the `KEY: value` file the stage wrote, each key exactly once, so a trace can never pass for the set. Stage `rb` used one container, `holzcloud-t34-rb`.

The set was made once, from fresh containers, on the clean committed tree of `00854be` without `REHEARSAL`, started 2026-11-21T13:13:55Z and ended 13:25:09Z on the host clock (11 minutes 14 seconds for all three stages), exit status 0 (`deploy/test/out/amd64-t34-pass1.started`, `deploy/test/out/amd64-t34-pass1.finished` and `deploy/test/out/amd64-t34-pass1.status`). Before it the worktree held no earlier output of the driver and no `holzcloud-t34-*` container existed, and after it `bash deploy/test/upgrade-verify.sh down` removed the four containers. In the same window and at the same commit, `PLATFORM=linux/amd64 TARGET=systemd-base bash deploy/test/backup-verify.sh run` made the set of phase 32 again (13:25:27Z to 13:27:56Z, status 0). No second set was made against `00854be`, and no other set is cited. An earlier set was made against `36b67cc`. Its stage `rb` read `UPD-02 PASS`, but the update's check line of *Updating* step 3 printed `no migration, no snapshot` after the failed update, because it matched `pre-upgrade snapshot written` only and the restarts since `6e86cdf` log `pre-upgrade snapshot kept`. The commit `00854be` fixed that line in *Updating* and in *Rolling back*, and made the driver record it (`UPDATE-CHECK-SNAPSHOT-LINE`) and require it for the verdict. So this record was written from this set, and nothing of the earlier set is quoted. The command was:

```text
PLATFORM=linux/amd64 TARGET=systemd-base bash deploy/test/upgrade-verify.sh run
```

The knob not set on that line had its default: `STAGES=mig upd rb`. The log of the run is `deploy/test/out/amd64-t34-pass1.log`, which is not pasted; its first line reads `upgrade-verify: platform linux/amd64, NEW_REF pinned to 00854be, v0.1.1 is 5bcae9b448db814701310c5dfab26031d68c0d92, stages: mig upd rb`.

No bound file changed after the evidence commit. The bound files for this set are `deploy/DEPLOY.md`, `deploy/backup.sh`, `deploy/restore.sh`, `deploy/holzcloud.service`, `deploy/test/run.sh`, `deploy/test/Dockerfile`, `deploy/test/upgrade-verify.sh` and the five plugin archives under `plugins/`, together with the binary built from `cmd`, `internal`, `go.mod`, `go.sum`, `changelog.go` and the `CHANGELOG.md` it embeds. This command exited 0 when the record was written:

```text
git diff --quiet 00854be HEAD -- deploy/DEPLOY.md deploy/backup.sh deploy/restore.sh deploy/holzcloud.service deploy/test/run.sh deploy/test/Dockerfile deploy/test/upgrade-verify.sh plugins/bestellung/bestellung.zip plugins/jahreszahl/jahreszahl.zip plugins/kontaktformular/kontaktformular.zip plugins/nicht-gefunden/nicht-gefunden.zip plugins/suche/suche.zip cmd internal go.mod go.sum changelog.go CHANGELOG.md
```

Before the run, `bash deploy/test/upgrade-verify.sh selftest` ended `SELFTEST OK (409 cases)`; its output is not part of the set.

The stage in order: B0 was installed by the document's Install text and seeded, and its state captured as `s0`. The driver then waited until the last start of the unit was older than its start-limit interval (`START-WINDOW-WAIT-S: 276` in `amd64-t34-rb/rb.txt`), a wait and not a reset, so that the restart loop had its full five starts as on a server where an update comes long after the last start. *Updating* blocks 1 to 3 were played line by line (`amd64-t34-rb/rb-update.out`), and no `reset-failed` ran between the failing start and the end of its restart loop. The journal of the loop was kept (`amd64-t34-rb/rb-journal.txt`) and the snapshots listed with their levels (`amd64-t34-rb/rb-snapshots.txt`). Then `SCENARIO[previous-on-half]` started the previous binary on the half-migrated database (`amd64-t34-rb/rb-guard.txt`). *Rolling back* block 1 was played as written (`amd64-t34-rb/rb-rollback.out`), and the state after it, `s2`, was compared with `s0` line by line (`amd64-t34-rb/rb-state.txt`).

## Deviations named

Everything here precedes every result: it is what the container does differently from a real server and what the harness did that the document does not contain. The first column is the name written to `amd64-t34-rb/deviations.txt`; the kinds of the whole run are the `DEVIATIONS:` line of `amd64-t34-summary.txt`, of which `old-plugins`, `old-script`, `backup-name`, `browser` and `transfer` belong to the stage `mig` only.

| Name | Why | What a real server does instead |
|------|-----|---------------------------------|
| DEVIATION[emulation] | The Docker host is aarch64 and the linux/amd64 container runs under emulation (`host architecture: aarch64`, `container uname -m: x86_64` in `amd64-t34-rb/versions.txt`), so every wall time of this set, the restarts of the unit among them, is an emulation figure | A native amd64 CPU and disk |
| DEVIATION[synthetic-build] | The binary of the update is F, a throwaway build of this tree with two added migrations: `00070_t34_ok.sql` creates the table `t34_ok` and succeeds, `00071_t34_fail.sql` creates `t34_half` and then inserts into a table that does not exist. It was made in a scratch clone with one local commit and the Build line of the document, the clone deleted afterwards; nothing of it is committed or pushed, and the Build line of *Updating* block 2 is skipped for it (`SKIP[synthetic-build]` in `amd64-t34-rb/rb-update.out`). Two files, because goose runs each file in its own transaction: one failing file alone is rolled back whole, the database stays at 69, and the previous binary would have nothing to refuse. Released binaries give no update whose migration fails part-way while the previous binary carries the schema guard | A release whose migration fails part-way, for whatever reason it fails |
| SCENARIO[previous-on-half] | After the restart loop ended, `systemctl reset-failed holzcloud`, `cp -p /opt/holzcloud/holzcloud.prev /opt/holzcloud/holzcloud` and `systemctl start holzcloud`: the previous binary was started on the half-migrated database, which *Rolling back* says not to do, so that the record shows the schema guard refusing it | An operator who takes *Rolling back* block 1 and never starts the old binary on the migrated database |
| DEVIATION[start-limit] | `systemctl reset-failed holzcloud` ran once outside the document, before that scenario. The restart loop of the failing update ran to its end without one, and the `reset-failed` of *Rolling back* block 1 is the document's own line | No reset outside the document |
| DEVIATION[scp] | Four lines: `scp holzcloud user@your-server:` played as `docker cp` of B0 at the install and of F at the update to `/root/holzcloud`, `scp -r deploy/ user@your-server:` at the install as a tar stream of the committed `deploy` directory into `/root`, and `scp deploy/backup.sh deploy/restore.sh user@your-server:` at the update as `docker cp` of the two committed scripts to `/root/` | `scp` from the build machine into the operator's home directory |
| DEVIATION[fixture] | The five plugin archives of this tree were copied into `/opt/holzcloud/plugin-zips` before `run.sh seed`; with the website, page, AI key and test image that `run.sh seed` makes, they stand in for operator work | The operator's own content and plugins |
| DEVIATION[dropin] | `/etc/systemd/system/holzcloud.service.d/container.conf` sets `ProcSubset=all`, because the document's unchanged start crashed in the container (`unit-start-unchanged.txt` in `amd64-t34-rb`) | The unit as shipped |
| DEVIATION[mount] | The unit's `RequiresMountsFor` names `/opt/holzcloud/data`, a plain directory on the container's root filesystem, so that line was not exercised | A data directory on its own mount |

## Findings

Each finding has the class of the earlier records' rule: (1) a bug of the harness, fixed in the driver; (2) a limit of the container, a named deviation or a *Not tested* line; (3) a measured property of a shipped file that matches what it and its document say, recorded and not edited; (4) a defect, a shipped file doing less than it or its document says, with a fix in a commit of its own. No finding of this set is class 1 or class 4, and no shipped file was edited after the evidence commit. The defect that the earlier set at `36b67cc` showed in the check line of *Updating* step 3 was fixed in `00854be` before this set (*How this was run*), and F-6 is that line observed on the fixed text.

The test for class 4 was whether a result contradicts something `deploy/DEPLOY.md` or the binary says at `00854be`, or leaves the operator worse off than they say. The document says that a failed update is restarted every five seconds until the start limit, that from 0.1.5 the restarts write no further snapshot and log `pre-upgrade snapshot kept` with the same `path`, that every `restore with:` line names the file written before the first attempt, that an older binary refuses a newer database with `db_level` and `binary_level`, and that *Rolling back* block 1 restores the snapshot with the old binary and ends with `no migration, no snapshot`. F-1 to F-7 are those sentences observed.

**F-1 (class 3): The restart loop wrote one snapshot, before the first attempt, and kept it on every restart; the five `restore with:` lines all name it.** The first start of F logged `pre-upgrade snapshot written` with `"level":69`, and the four starts after it each logged `pre-upgrade snapshot kept` with the same `path`, `"level":69`, `"db_level":70` and the reason (`amd64-t34-rb/rb-journal.txt`; `JOURNAL-SNAPSHOT-WRITTEN: 1`, `JOURNAL-SNAPSHOT-KEPT: 4`, `JOURNAL-RESTORE-WITH: 5`, `RESTORE-WITH-DISTINCT-PATHS: 1`, `RESTORE-WITH-NAMES-WRITTEN: yes` in `amd64-t34-rb/rb.txt`). After the loop the data directory held that snapshot at level 69 and the snapshot of the install's first start at level 0, and none at level 70 (`SNAPSHOTS-AFTER-LOOP` in `amd64-t34-rb/rb.txt`, `amd64-t34-rb/rb-snapshots.txt`); while the loop ran, `pre-upgrade.pending` named the snapshot and its level (`MARKER-DURING: pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite 69`). This is what the fix `6e86cdf` ("restarts after a failed migration keep the snapshot before it instead of snapshotting the half-migrated database") does, and it answers the probe of finding F-1 of phase 34's research, in which the code before it wrote a further snapshot at level 70 on each of starts 2 to 5, and named those half-migrated copies in its `restore with:` lines. The line of the first start and the line of a restart (`amd64-t34-rb/rb-journal.txt`):

```text
{"time":"2026-11-21T13:23:41.823526554Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite","level":69}
{"time":"2026-11-21T13:23:47.221236417Z","level":"WARN","msg":"pre-upgrade snapshot kept","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite","level":69,"db_level":70,"reason":"an earlier start of this update applied part of its migrations and did not finish; the snapshot from before it stays the way back, so no snapshot of the half-migrated database is written"}
```

**F-2 (class 3): The journal names migration 70 as applied before the failure of 71, so the half-migrated level is in the journal and not only in the database.** The first start's lines after the snapshot (`amd64-t34-rb/rb-journal.txt`; `APPLIED-BEFORE-FAILURE: 70`, `LEVEL-HALF: 70` in `amd64-t34-rb/rb.txt`):

```text
{"time":"2026-11-21T13:23:41.866263013Z","level":"INFO","msg":"migration applied","version":70,"type":"sql","duration_ms":5}
{"time":"2026-11-21T13:23:41.867951888Z","level":"ERROR","msg":"migrations failed","err":"goose up: partial migration error (type:sql,version:71): SQL logic error: no such table: t34_no_such_table (1)"}
{"time":"2026-11-21T13:23:41.868448221Z","level":"ERROR","msg":"restore with: systemctl stop holzcloud && cp /opt/holzcloud/data/pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite /opt/holzcloud/data/holzcloud.sqlite && systemctl start holzcloud"}
```

**F-3 (class 3): systemd's start limit ended the loop after five starts, and the unit stayed failed with nothing answering.** The journal counts the restarts to 5 and ends with `Start request repeated too quickly.` (`amd64-t34-rb/rb-journal.txt`); `STARTS: 5`, `NRESTARTS: 5`, `RESULT: exit-code`, `ACTIVE-STATE: failed`, `START-LIMIT-LINE: yes` and `HEALTHZ-DURING-LOOP: never` (`amd64-t34-rb/rb.txt`). The starts came about five seconds apart, from 13:23:41Z to 13:24:03Z.

**F-4 (class 3): The previous binary, started on the half-migrated database, refused it with the schema guard's sentence and both levels, five times, and never answered.** The sentence is the one `internal/db/schema_guard.go` builds, and it names the way back (`amd64-t34-rb/rb-guard.txt`; `GUARD-LINE: yes`, `GUARD-DB-LEVEL: 70`, `GUARD-BINARY-LEVEL: 69`, `GUARD-SENTENCE: yes`, `GUARD-STARTS: 5`, `GUARD-HEALTHZ: never`, `GUARD-ACTIVE-STATE: failed` in `amd64-t34-rb/rb.txt`):

```text
{"time":"2026-11-21T13:24:18.899698418Z","level":"ERROR","msg":"refusing to start: the database schema is newer than this binary","db_level":70,"binary_level":69,"err":"the database is at migration level 70, and this binary knows migrations up to 69: a newer holzcloud has migrated it, and this binary refuses to run against a schema it does not know. The way back: start the newer binary again, or restore the pre-upgrade snapshot that the newer binary wrote before it migrated, together with the older binary, as deploy/DEPLOY.md describes under \"Rolling back\""}
```

**F-5 (class 3): The first two lines of *Rolling back* block 1 listed both snapshots with their levels, and the restore line took the one at level 69.** The level list (`amd64-t34-rb/rb-rollback.out`):

```text
1795267492.257012545 /opt/holzcloud/data/pre-upgrade-v0.1.4-255-g00854be-20261121-131835.sqlite level 0
1795267492.367905295 /opt/holzcloud/data/pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite level 69
```

The file at level 0 is the snapshot of the install's first start, 12288 bytes by the `ls -l` line above it, the "snapshot of a first start, an empty database" that the document warns of; the newest by modification time is the update's own, at the level of the old binary, so the `ls -t` of the restore line took the right file.

**F-6 (class 3): After the failed update, the check line of *Updating* step 3 printed the `pre-upgrade snapshot kept` line of the last start, naming the snapshot, as the document now says.** Line 1 of the step failed after its retries with status 7, as it should (`UPDATE-CHECK-RC: 7`), and line 2 printed (`amd64-t34-rb/rb-update.out`; `UPDATE-CHECK-SNAPSHOT-LINE: yes` in `amd64-t34-rb/rb.txt`):

```text
1795267452.275755429 {"time":"2026-11-21T13:24:03.15234055Z","level":"WARN","msg":"pre-upgrade snapshot kept","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite","level":69,"db_level":70,"reason":"an earlier start of this update applied part of its migrations and did not finish; the snapshot from before it stays the way back, so no snapshot of the half-migrated database is written"}
```

It shows the line of the last start only, because `ActiveEnterTimestamp` names the last start of the loop; the `written` line of the first start is earlier in the journal (F-1).

**F-7 (class 3): *Rolling back* block 1 played as written put the previous build back on the state before the update, and kept the migrated database aside.** Every line ended with status 0 (`ROLLBACK-RC: 0 0 0 0 0 0 0 0 0 0 0`), `/readyz` showed the previous build with `status` and `integrity` `ok` (`ROLLBACK-READYZ`), and the last line printed `no migration, no snapshot` (`ROLLBACK-CHECK-LINE`). The database is at 69 without the table of migration 70 (`LEVEL-AFTER-ROLLBACK: 69`, `T34-OK-TABLE: absent`), `holzcloud.sqlite.after-update` holds level 70 (`AFTER-UPDATE-LEVEL: 70`), and the marker is gone (`MARKER-AFTER: absent`), all in `amd64-t34-rb/rb.txt`. The state before the update and after the rollback, schema, tables, the rows of five tables and the hashes of the seeded page and picture, compared line by line with no difference (`amd64-t34-rb/rb-state.txt`, `COMPARE-S0-S2: bad=0`).

## UPD-02: a migration that fails part-way, and the way back

**What was done.** B0 was installed and seeded, and its state `s0` read through a read-only connection (`amd64-t34-rb/rb-state.txt`). *Updating* block 1 copied the binary aside as `holzcloud.prev`, printed its version and made the backup with the document's line, ending `=== Backup complete: /opt/holzcloud/backups/20261121-132338 ===`; block 2 skipped the Build line (`DEVIATION[synthetic-build]`), carried F and the two scripts over, stopped the service, copied F into place, installed the scripts and started it, each line with status 0 (`amd64-t34-rb/rb-update.out`). The first start of F wrote the snapshot at level 69, applied 70, failed at 71 and printed its `restore with:` line (F-1, F-2); systemd restarted it four times, each restart keeping the snapshot and failing at 71 again, until the start limit (F-3; `amd64-t34-rb/rb-journal.txt`). Block 3 found nothing listening (`UPDATE-CHECK-RC: 7`) and its journal line named the kept snapshot (F-6). The previous binary was then started on the half-migrated database and refused it five times (`SCENARIO[previous-on-half]`, F-4; `amd64-t34-rb/rb-guard.txt`). Last, *Rolling back* block 1 was played as written, its own `reset-failed` included: the two listing lines (F-5), the stop, the previous binary back, the write-ahead files removed, the migrated database set aside, the snapshot copied in as `holzcloud`, the reset, the start and the two check lines (`amd64-t34-rb/rb-rollback.out`):

```text
1795267505.333816593   "status": "ok",
1795267505.350889218   "version": "v0.1.4-255-g00854be",
1795267505.367354843   "commit": "00854be",
1795267505.383830551   "integrity": "ok"
1795267505.718159301 no migration, no snapshot
```

Then `s2` was read and compared with `s0` (F-7).

*What it shows.* When an update to a build whose migration 71 failed after migration 70 had run, the service stopped with the failing migration named in the journal and was restarted by systemd until its start limit, every restart kept the one pre-upgrade snapshot written before the first attempt at level 69 and every `restore with:` line named it, the previous binary refused the database at level 70 with the schema guard's sentence, and *Rolling back* block 1, played as written, started the previous binary on a database whose schema, tables, rows and files equal those before the update.

**What it does not show.** A migration that fails for another reason, a release as the failing build, an operator who starts the failing build again, and the files the snapshot does not hold (*Not tested*).

## Measurements

The figures the harness wrote. Every one is an emulation figure on an aarch64 host (`amd64-t34-rb/versions.txt`); none is a promise of speed.

| Figure | Value | Where |
|--------|-------|-------|
| The wait until the last start was outside the start-limit interval | 276 s | `START-WINDOW-WAIT-S` in `amd64-t34-rb/rb.txt` |
| The backup of *Updating* block 1 | 1795267418.642721636 to 1795267419.953274095 wall, 21M | the `WALL` line and the size line of `amd64-t34-rb/rb-update.out` |
| The first start of F, from `starting` to `restore with:` | 13:23:41.665134596Z to 13:23:41.868448221Z | `amd64-t34-rb/rb-journal.txt` |
| The restart loop, first to fifth start | 13:23:41.665134596Z to 13:24:03.09254805Z | `amd64-t34-rb/rb-journal.txt` |
| The readiness check of *Updating* step 3, 31 attempts | 1795267421.676714137 to 1795267451.982999554 wall | `amd64-t34-rb/rb-update.out` |
| The refusals of the previous binary, first to fifth start | 13:24:18.836233877Z to 13:24:39.840858803Z | `amd64-t34-rb/rb-guard.txt` |
| *Rolling back* block 1, first line to last | 1795267491.887905045 to 1795267505.740274635 wall | `amd64-t34-rb/rb-rollback.out` |
| The snapshot of the update and that of the first start | 765952 bytes and 12288 bytes | the `ls -l` lines of `amd64-t34-rb/rb-rollback.out` |
| The whole run, all three stages | 13:13:55Z to 13:25:09Z | `deploy/test/out/amd64-t34-pass1.started`, `deploy/test/out/amd64-t34-pass1.finished` |

## Container vs. server

| Aspect | In the container | On a real server | Measured or assumed |
|--------|------------------|------------------|---------------------|
| Speed and timing | amd64 emulated on aarch64 | A native server | Not measured natively |
| The failing build | A throwaway build of this tree, two added migrations, the second inserting into a missing table | A release whose migration fails for its own reason | Measured for this failure; others assumed to take the same path in the binary |
| The start window | A wait of 276 s so that the loop got its five starts | An update long after the last start | Measured here |
| The previous binary on the half-migrated database | Started on purpose, outside the document | Not started, if the operator follows *Rolling back* | Measured for the refusal |
| The build and transfer of F | Built in a scratch clone, carried by `docker cp` | Built on the development machine, carried by `scp` | Named deviations |
| `holzcloud.service` | With the container-only drop-in `ProcSubset=all` | As shipped | Named deviation |

## Not tested

- A migration that fails for another reason, a full disk or a process killed in the middle of a migration, was not played; the failure here is an insert into a table that does not exist.
- A binary released before 0.1.5 was not played as the failing build; F is a build of this tree, so it carries the marker and the `kept` line.
- An operator who runs `systemctl reset-failed` and starts the failing build again after the loop was not played; the marker keeps naming the snapshot while the database is above its level, which the tests of phase 34 show and this set does not.
- The first migration of an update failing, so that the database stays at the old level and each restart writes a snapshot of that level, was not played.
- Media, templates, plugins, language files and the logo added after the update were not played: the snapshot holds the database only, and the seeded files were not changed between `s0` and `s2`.
- *Rolling back* block 2, the update that wrote no snapshot, was not played here.
- Undoing the rollback with `holzcloud.sqlite.after-update` was not played; the file was only read for its level.
- A rollback on a real server (SRV-02) was not played; the container ran on the Mac.
- Native amd64 timing was not measured.

## Owner steps

None. Every step of this record ran in a container on the Mac, by the owner's decision; nothing here has to be done on the real server before the record counts.

## Hand-over to docs/beta.md

The sentences of this record that the page may quote, named by place and not copied here, so that each stays in exactly one place: the title sentences of F-1, F-2, F-3, F-4, F-6 and F-7 and the sentence after *What it shows.* in *UPD-02*; the limits, the bullets of *Not tested*.

The limit of the page this record answers, named by its first words: "A rollback after a failed migration was not played" (a migration that failed part-way was played, and *Rolling back* block 1 brought the previous state back, *UPD-02*). The line under *The product itself* that begins "The schema and the migrations are not promised stable" stays as it is: F-4 shows how this release keeps the second half of it, and the first half is not changed by this record.

Nothing in this record is a promise until the page cites it byte for byte. A change to a bound file after the evidence commit means the set is made again before the record is cited for the changed text. `docs/beta.md` was not edited by this plan.

## Appendix: raw output

Each block is the scratch file of the evidence set, pasted unchanged by `bash deploy/test/upgrade-verify.sh appendix rb`, and `bash deploy/test/upgrade-verify.sh record-check` compares every block with its scratch file byte for byte. The files are in `deploy/test/out` (gitignored); the `step-NN.txt` files, `seed.txt`, `http.txt`, `install.txt`, `paths.txt`, the `.lst` files, the unit-start files, the synthetic build under `syn` and the log are not pasted.

#### amd64-t34-summary.txt

```text
MIG-01 PASS
MIG-02 PASS
UPD-01 PASS
UPD-02 PASS
STAGES: mig upd rb
DEVIATIONS: backup-name browser dropin emulation fixture mount old-plugins old-script scp start-limit synthetic-build transfer
BINARY: holzcloud v0.1.4-255-g00854be (00854be)
```

#### amd64-t34-rb/versions.txt

```text
host architecture:     aarch64
host operating system: Docker Desktop
host kernel:           7.0.14-linuxkit
container uname:       Linux 7.0.14-linuxkit x86_64
container systemd:     systemd 252 (252.39-1~deb12u2)
container os:          PRETTY_NAME="Debian GNU/Linux 12 (bookworm)"
host uname -m:         arm64
container uname -m:    x86_64
container systemctl:   systemd 252 (252.39-1~deb12u2)
running binary:        holzcloud v0.1.4-255-g00854be (00854be)
holzcloud-f:  v0.1.4-256-g4c07745 (4c07745) migrations 71 sha256 af493d4a045ef7f0f54c4cb9021a9b561a8b0492828b0d9e8bd8389e1bb3b3d9
NEW_REF:               00854be
evidence commit:       00854be
```

#### amd64-t34-rb/deviations.txt

```text
DEVIATION[fixture]: the five plugin archives of this tree are copied into /opt/holzcloud/plugin-zips before run.sh seed; with the website, page, AI key and test image run.sh seed makes, they stand in for operator work in the admin and are no step of the document
DEVIATION[synthetic-build]: the binary of this update is F (v0.1.4-256-g4c07745), a throwaway build of this tree with two added migrations, 00070_t34_ok.sql (creates t34_ok) and 00071_t34_fail.sql (creates t34_half, then inserts into a table that does not exist), made in a scratch clone with one local commit and the Build line of the document, the clone deleted afterwards; nothing of it is committed or pushed. The Build line of Updating block 2 is skipped for it. Released binaries give no update whose migration fails part-way while the previous binary carries the schema guard
DEVIATION[scp]: 'scp holzcloud user@your-server:' was played as docker cp of the synthetic build of this update to /root/holzcloud in the container (there is no ssh and no second machine here)
DEVIATION[scp]: 'scp deploy/backup.sh deploy/restore.sh user@your-server:' was played as docker cp of the working tree's deploy/backup.sh deploy/restore.sh to /root/ in the container (there is no ssh and no second machine here)
SCENARIO[previous-on-half]: after the restart loop ended, systemctl reset-failed holzcloud, cp -p /opt/holzcloud/holzcloud.prev /opt/holzcloud/holzcloud and systemctl start holzcloud: the previous binary is started on the half-migrated database, which Rolling back says not to do, so that the record shows the schema guard refusing it
DEVIATION[start-limit]: systemctl reset-failed holzcloud runs once outside the document, before the previous binary is started on the half-migrated database (SCENARIO[previous-on-half]); the restart loop of the failing update ran to its end without one, and the reset-failed of Rolling back block 1 is the document's own line
DEVIATION[emulation]: the Docker host is aarch64; the linux/amd64 container runs under emulation, so every wall time of this set (the restarts of the unit among them) is an emulation figure, not a server's
DEVIATION[scp]: 'scp holzcloud user@your-server:' was played as docker cp of deploy/test/out/bin/holzcloud-new-00854be to /root/holzcloud in the container (there is no ssh and no second machine here)
DEVIATION[scp]: 'scp -r deploy/ user@your-server:' was played as a tar stream of the working tree's deploy directory (without deploy/test/out) into /root in the container, by docker cp -
DEVIATION[dropin]: /etc/systemd/system/holzcloud.service.d/container.conf sets ProcSubset=all (container only; the document's unchanged start crashed, see unit-start-unchanged.txt)
DEVIATION[mount]: the unit's RequiresMountsFor names /opt/holzcloud/data, a plain directory on the container's root filesystem; no separate mount exists here, so that line was not exercised
```

#### amd64-t34-rb/rb.txt

```text
B0-VERSION: v0.1.4-255-g00854be
LEVEL-EXPECTED-B0: 69
LEVEL-EXPECTED-F: 71
F-BUILD: v0.1.4-256-g4c07745 (4c07745) migrations 71 sha256 af493d4a045ef7f0f54c4cb9021a9b561a8b0492828b0d9e8bd8389e1bb3b3d9
START-WINDOW-WAIT-S: 276
UPDATE-CHECK-RC: 7
UPDATE-CHECK-SNAPSHOT-LINE: yes
NRESTARTS: 5
RESULT: exit-code
ACTIVE-STATE: failed
HEALTHZ-DURING-LOOP: never
STARTS: 5
JOURNAL-SNAPSHOT-WRITTEN: 1
JOURNAL-SNAPSHOT-KEPT: 4
JOURNAL-MIGRATIONS-FAILED: 5
JOURNAL-RESTORE-WITH: 5
RESTORE-WITH-DISTINCT-PATHS: 1
RESTORE-WITH-NAMES-WRITTEN: yes
WRITTEN-PATH: /opt/holzcloud/data/pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite
WRITTEN-LEVEL: 69
STOP-MESSAGE: goose up: partial migration error (type:sql,version:71): SQL logic error: no such table: t34_no_such_table (1)
APPLIED-BEFORE-FAILURE: 70
START-LIMIT-LINE: yes
LEVEL-HALF: 70
MARKER-DURING: pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite 69
SNAPSHOTS-AFTER-LOOP: pre-upgrade-v0.1.4-255-g00854be-20261121-131835.sqlite:0 pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite:69
GUARD-HEALTHZ: never
GUARD-ACTIVE-STATE: failed
GUARD-STARTS: 5
GUARD-LINE: yes
GUARD-DB-LEVEL: 70
GUARD-BINARY-LEVEL: 69
GUARD-SENTENCE: yes
ROLLBACK-RC: 0 0 0 0 0 0 0 0 0 0 0
ROLLBACK-READYZ: {"status":"ok","version":"v0.1.4-255-g00854be","commit":"00854be","integrity":"ok"}
ROLLBACK-CHECK-LINE: no migration, no snapshot
LEVEL-AFTER-ROLLBACK: 69
T34-OK-TABLE: absent
AFTER-UPDATE-LEVEL: 70
MARKER-AFTER: absent
COMPARE-S0-S2: bad=0
UPD-02 PASS
```

#### amd64-t34-rb/rb-update.out

```text
$ sudo cp -p /opt/holzcloud/holzcloud /opt/holzcloud/holzcloud.prev
WALL 1795267418.140374427 1795267418.279961927
RC 0
$ /opt/holzcloud/holzcloud version
1795267418.519708386 holzcloud v0.1.4-255-g00854be (00854be)
WALL 1795267418.390790969 1795267418.541759553
RC 0
$ sudo -u holzcloud /opt/holzcloud/deploy/backup.sh /opt/holzcloud/backups
1795267418.820559803 === holzcloud-CMS backup: 20261121-132338 ===
1795267418.906907094 Backing up database...
1795267419.070586470 2026/11/21 13:23:39 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/data/holzcloud.sqlite
1795267419.156322845 2026/11/21 13:23:39 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/backups/20261121-132338/holzcloud.sqlite
1795267419.176675886 2026/11/21 13:23:39 INFO backup verified path=/opt/holzcloud/backups/20261121-132338/holzcloud.sqlite tables=60
1795267419.195115261 backup written and verified: /opt/holzcloud/backups/20261121-132338/holzcloud.sqlite
1795267419.216621303 Backing up CSRF key...
1795267419.278972553 Backing up media...
1795267419.360935970 Backing up plugins...
1795267419.483105386 Writing checksums...
1795267419.844164887 Pruning backups older than 30 days...
1795267419.908674762 === Backup complete: /opt/holzcloud/backups/20261121-132338 ===
1795267419.925444720 21M	/opt/holzcloud/backups/20261121-132338
WALL 1795267418.642721636 1795267419.953274095
RC 0
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -X main.Version=$(git describe --tags --always) -X main.Commit=$(git rev-parse --short HEAD)" ./cmd/holzcloud
SKIP[synthetic-build] the binary of this update was built by this line in a scratch clone before the update
RC 0
$ scp holzcloud user@your-server:
RC 0
$ scp deploy/backup.sh deploy/restore.sh user@your-server:
RC 0
$ sudo systemctl stop holzcloud
WALL 1795267420.607059345 1795267420.766180012
RC 0
$ sudo cp ~/holzcloud /opt/holzcloud/holzcloud
WALL 1795267420.932625512 1795267421.066973345
RC 0
$ sudo install -m 755 -o root -g root ~/backup.sh ~/restore.sh /opt/holzcloud/deploy/
WALL 1795267421.167608471 1795267421.298952887
RC 0
$ sudo systemctl start holzcloud
WALL 1795267421.410681179 1795267421.559648137
RC 0
$ curl -fsS --retry 30 --retry-connrefused --retry-delay 1 http://127.0.0.1:8080/readyz | jq '{status, version, commit, integrity}'
1795267421.793477888 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 4 ms: Couldn't connect to server
1795267422.793897513 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267423.795177847 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267424.798272833 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267425.803925375 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267426.818367917 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267427.813608709 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267428.820576710 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267429.822872919 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267430.832650961 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267431.833788336 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267432.845724670 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267433.850604379 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267434.863622338 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267435.871451005 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267436.873200672 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267437.875060506 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267438.880159631 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267439.881904007 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267440.892006716 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267441.903135008 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267442.903081883 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267443.915723342 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267444.917938968 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267445.920704593 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267446.917328760 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267447.923434052 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267448.929446511 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267449.947017095 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267450.958436345 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267451.959907221 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
WALL 1795267421.676714137 1795267451.982999554
RC 7
$ sudo journalctl -u holzcloud --since "$(systemctl show -p ActiveEnterTimestamp --value holzcloud)" -o cat | grep -E 'pre-upgrade snapshot (written|kept)|migration applied' || echo 'no migration, no snapshot'
1795267452.275755429 {"time":"2026-11-21T13:24:03.15234055Z","level":"WARN","msg":"pre-upgrade snapshot kept","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite","level":69,"db_level":70,"reason":"an earlier start of this update applied part of its migrations and did not finish; the snapshot from before it stays the way back, so no snapshot of the half-migrated database is written"}
WALL 1795267452.081597763 1795267452.300160721
RC 0
```

#### amd64-t34-rb/rb-journal.txt

```text
Stopping holzcloud.service - holzcloud-CMS...
{"time":"2026-11-21T13:23:40.722594179Z","level":"INFO","msg":"shutting down","timeout_s":10}
{"time":"2026-11-21T13:23:40.734097012Z","level":"INFO","msg":"database maintenance","checkpoint_busy":0,"db_size_bytes":770048,"wal_size_bytes":0,"page_count":188,"freelist_count":0}
{"time":"2026-11-21T13:23:40.734589637Z","level":"INFO","msg":"shutdown complete"}
holzcloud.service: Deactivated successfully.
Stopped holzcloud.service - holzcloud-CMS.
holzcloud.service: Consumed 13.902s CPU time.
Started holzcloud.service - holzcloud-CMS.
{"time":"2026-11-21T13:23:41.665134596Z","level":"INFO","msg":"starting","version":"v0.1.4-256-g4c07745","commit":"4c07745"}
{"time":"2026-11-21T13:23:41.669244096Z","level":"INFO","msg":"configuration loaded","config":{"port":"8080","data_dir":"/opt/holzcloud/data","db_path":"/opt/holzcloud/data/holzcloud.sqlite","log_level":"INFO","secure":true,"max_template_size":10485760,"max_media_size":5242880,"max_video_size":67108864,"max_megapixels":24,"argon2_memory_kb":65536,"argon2_iterations":1,"argon2_parallelism":2,"min_free_bytes":536870912,"smtp_host":"","smtp_from":"","smtp_configured":false,"payrexx_instance":"","payrexx_configured":false,"trusted_proxies":"127.0.0.1/32,::1/128","listen":"127.0.0.1","sso_enabled":false,"sso_configured":false,"sso_provision":false,"sso_admin_group":"","sso_default_website":0,"sso_website_groups":0,"sso_sign_out_path":"/outpost.goauthentik.io/sign_out","oidc_enabled":false,"oidc_issuer":"","oidc_client_id":"","oidc_redirect_url":"","oidc_keys":0,"oidc_client_secret_set":false,"smtp":"aus","secret_key_set":false}}
{"time":"2026-11-21T13:23:41.716869387Z","level":"INFO","msg":"sqlite opened","journal_mode":"wal","path":"/opt/holzcloud/data/holzcloud.sqlite"}
{"time":"2026-11-21T13:23:41.807700138Z","level":"INFO","msg":"sqlite opened","journal_mode":"wal","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite"}
{"time":"2026-11-21T13:23:41.821266554Z","level":"INFO","msg":"backup verified","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite","tables":62}
{"time":"2026-11-21T13:23:41.823526554Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite","level":69}
{"time":"2026-11-21T13:23:41.866263013Z","level":"INFO","msg":"migration applied","version":70,"type":"sql","duration_ms":5}
{"time":"2026-11-21T13:23:41.867951888Z","level":"ERROR","msg":"migrations failed","err":"goose up: partial migration error (type:sql,version:71): SQL logic error: no such table: t34_no_such_table (1)"}
{"time":"2026-11-21T13:23:41.868448221Z","level":"ERROR","msg":"restore with: systemctl stop holzcloud && cp /opt/holzcloud/data/pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite /opt/holzcloud/data/holzcloud.sqlite && systemctl start holzcloud"}
holzcloud.service: Main process exited, code=exited, status=1/FAILURE
holzcloud.service: Failed with result 'exit-code'.
holzcloud.service: Scheduled restart job, restart counter is at 1.
Stopped holzcloud.service - holzcloud-CMS.
Started holzcloud.service - holzcloud-CMS.
{"time":"2026-11-21T13:23:47.159407126Z","level":"INFO","msg":"starting","version":"v0.1.4-256-g4c07745","commit":"4c07745"}
{"time":"2026-11-21T13:23:47.162010876Z","level":"INFO","msg":"configuration loaded","config":{"port":"8080","data_dir":"/opt/holzcloud/data","db_path":"/opt/holzcloud/data/holzcloud.sqlite","log_level":"INFO","secure":true,"max_template_size":10485760,"max_media_size":5242880,"max_video_size":67108864,"max_megapixels":24,"argon2_memory_kb":65536,"argon2_iterations":1,"argon2_parallelism":2,"min_free_bytes":536870912,"smtp_host":"","smtp_from":"","smtp_configured":false,"payrexx_instance":"","payrexx_configured":false,"trusted_proxies":"127.0.0.1/32,::1/128","listen":"127.0.0.1","sso_enabled":false,"sso_configured":false,"sso_provision":false,"sso_admin_group":"","sso_default_website":0,"sso_website_groups":0,"sso_sign_out_path":"/outpost.goauthentik.io/sign_out","oidc_enabled":false,"oidc_issuer":"","oidc_client_id":"","oidc_redirect_url":"","oidc_keys":0,"oidc_client_secret_set":false,"smtp":"aus","secret_key_set":false}}
{"time":"2026-11-21T13:23:47.207099167Z","level":"INFO","msg":"sqlite opened","journal_mode":"wal","path":"/opt/holzcloud/data/holzcloud.sqlite"}
{"time":"2026-11-21T13:23:47.221236417Z","level":"WARN","msg":"pre-upgrade snapshot kept","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite","level":69,"db_level":70,"reason":"an earlier start of this update applied part of its migrations and did not finish; the snapshot from before it stays the way back, so no snapshot of the half-migrated database is written"}
{"time":"2026-11-21T13:23:47.247912334Z","level":"ERROR","msg":"migrations failed","err":"goose up: partial migration error (type:sql,version:71): SQL logic error: no such table: t34_no_such_table (1)"}
{"time":"2026-11-21T13:23:47.249993376Z","level":"ERROR","msg":"restore with: systemctl stop holzcloud && cp /opt/holzcloud/data/pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite /opt/holzcloud/data/holzcloud.sqlite && systemctl start holzcloud"}
holzcloud.service: Main process exited, code=exited, status=1/FAILURE
holzcloud.service: Failed with result 'exit-code'.
holzcloud.service: Scheduled restart job, restart counter is at 2.
Stopped holzcloud.service - holzcloud-CMS.
Started holzcloud.service - holzcloud-CMS.
{"time":"2026-11-21T13:23:52.593836587Z","level":"INFO","msg":"starting","version":"v0.1.4-256-g4c07745","commit":"4c07745"}
{"time":"2026-11-21T13:23:52.596418837Z","level":"INFO","msg":"configuration loaded","config":{"port":"8080","data_dir":"/opt/holzcloud/data","db_path":"/opt/holzcloud/data/holzcloud.sqlite","log_level":"INFO","secure":true,"max_template_size":10485760,"max_media_size":5242880,"max_video_size":67108864,"max_megapixels":24,"argon2_memory_kb":65536,"argon2_iterations":1,"argon2_parallelism":2,"min_free_bytes":536870912,"smtp_host":"","smtp_from":"","smtp_configured":false,"payrexx_instance":"","payrexx_configured":false,"trusted_proxies":"127.0.0.1/32,::1/128","listen":"127.0.0.1","sso_enabled":false,"sso_configured":false,"sso_provision":false,"sso_admin_group":"","sso_default_website":0,"sso_website_groups":0,"sso_sign_out_path":"/outpost.goauthentik.io/sign_out","oidc_enabled":false,"oidc_issuer":"","oidc_client_id":"","oidc_redirect_url":"","oidc_keys":0,"oidc_client_secret_set":false,"smtp":"aus","secret_key_set":false}}
{"time":"2026-11-21T13:23:52.641333045Z","level":"INFO","msg":"sqlite opened","journal_mode":"wal","path":"/opt/holzcloud/data/holzcloud.sqlite"}
{"time":"2026-11-21T13:23:52.655215753Z","level":"WARN","msg":"pre-upgrade snapshot kept","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite","level":69,"db_level":70,"reason":"an earlier start of this update applied part of its migrations and did not finish; the snapshot from before it stays the way back, so no snapshot of the half-migrated database is written"}
{"time":"2026-11-21T13:23:52.695451712Z","level":"ERROR","msg":"migrations failed","err":"goose up: partial migration error (type:sql,version:71): SQL logic error: no such table: t34_no_such_table (1)"}
{"time":"2026-11-21T13:23:52.695971087Z","level":"ERROR","msg":"restore with: systemctl stop holzcloud && cp /opt/holzcloud/data/pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite /opt/holzcloud/data/holzcloud.sqlite && systemctl start holzcloud"}
holzcloud.service: Main process exited, code=exited, status=1/FAILURE
holzcloud.service: Failed with result 'exit-code'.
holzcloud.service: Scheduled restart job, restart counter is at 3.
Stopped holzcloud.service - holzcloud-CMS.
Started holzcloud.service - holzcloud-CMS.
{"time":"2026-11-21T13:23:57.843817506Z","level":"INFO","msg":"starting","version":"v0.1.4-256-g4c07745","commit":"4c07745"}
{"time":"2026-11-21T13:23:57.846386756Z","level":"INFO","msg":"configuration loaded","config":{"port":"8080","data_dir":"/opt/holzcloud/data","db_path":"/opt/holzcloud/data/holzcloud.sqlite","log_level":"INFO","secure":true,"max_template_size":10485760,"max_media_size":5242880,"max_video_size":67108864,"max_megapixels":24,"argon2_memory_kb":65536,"argon2_iterations":1,"argon2_parallelism":2,"min_free_bytes":536870912,"smtp_host":"","smtp_from":"","smtp_configured":false,"payrexx_instance":"","payrexx_configured":false,"trusted_proxies":"127.0.0.1/32,::1/128","listen":"127.0.0.1","sso_enabled":false,"sso_configured":false,"sso_provision":false,"sso_admin_group":"","sso_default_website":0,"sso_website_groups":0,"sso_sign_out_path":"/outpost.goauthentik.io/sign_out","oidc_enabled":false,"oidc_issuer":"","oidc_client_id":"","oidc_redirect_url":"","oidc_keys":0,"oidc_client_secret_set":false,"smtp":"aus","secret_key_set":false}}
{"time":"2026-11-21T13:23:57.898000756Z","level":"INFO","msg":"sqlite opened","journal_mode":"wal","path":"/opt/holzcloud/data/holzcloud.sqlite"}
{"time":"2026-11-21T13:23:57.912108214Z","level":"WARN","msg":"pre-upgrade snapshot kept","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite","level":69,"db_level":70,"reason":"an earlier start of this update applied part of its migrations and did not finish; the snapshot from before it stays the way back, so no snapshot of the half-migrated database is written"}
{"time":"2026-11-21T13:23:57.948455881Z","level":"ERROR","msg":"migrations failed","err":"goose up: partial migration error (type:sql,version:71): SQL logic error: no such table: t34_no_such_table (1)"}
{"time":"2026-11-21T13:23:57.948948006Z","level":"ERROR","msg":"restore with: systemctl stop holzcloud && cp /opt/holzcloud/data/pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite /opt/holzcloud/data/holzcloud.sqlite && systemctl start holzcloud"}
holzcloud.service: Main process exited, code=exited, status=1/FAILURE
holzcloud.service: Failed with result 'exit-code'.
holzcloud.service: Scheduled restart job, restart counter is at 4.
Stopped holzcloud.service - holzcloud-CMS.
Started holzcloud.service - holzcloud-CMS.
{"time":"2026-11-21T13:24:03.09254805Z","level":"INFO","msg":"starting","version":"v0.1.4-256-g4c07745","commit":"4c07745"}
{"time":"2026-11-21T13:24:03.095805217Z","level":"INFO","msg":"configuration loaded","config":{"port":"8080","data_dir":"/opt/holzcloud/data","db_path":"/opt/holzcloud/data/holzcloud.sqlite","log_level":"INFO","secure":true,"max_template_size":10485760,"max_media_size":5242880,"max_video_size":67108864,"max_megapixels":24,"argon2_memory_kb":65536,"argon2_iterations":1,"argon2_parallelism":2,"min_free_bytes":536870912,"smtp_host":"","smtp_from":"","smtp_configured":false,"payrexx_instance":"","payrexx_configured":false,"trusted_proxies":"127.0.0.1/32,::1/128","listen":"127.0.0.1","sso_enabled":false,"sso_configured":false,"sso_provision":false,"sso_admin_group":"","sso_default_website":0,"sso_website_groups":0,"sso_sign_out_path":"/outpost.goauthentik.io/sign_out","oidc_enabled":false,"oidc_issuer":"","oidc_client_id":"","oidc_redirect_url":"","oidc_keys":0,"oidc_client_secret_set":false,"smtp":"aus","secret_key_set":false}}
{"time":"2026-11-21T13:24:03.138597467Z","level":"INFO","msg":"sqlite opened","journal_mode":"wal","path":"/opt/holzcloud/data/holzcloud.sqlite"}
{"time":"2026-11-21T13:24:03.15234055Z","level":"WARN","msg":"pre-upgrade snapshot kept","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite","level":69,"db_level":70,"reason":"an earlier start of this update applied part of its migrations and did not finish; the snapshot from before it stays the way back, so no snapshot of the half-migrated database is written"}
{"time":"2026-11-21T13:24:03.222788425Z","level":"ERROR","msg":"migrations failed","err":"goose up: partial migration error (type:sql,version:71): SQL logic error: no such table: t34_no_such_table (1)"}
{"time":"2026-11-21T13:24:03.223289508Z","level":"ERROR","msg":"restore with: systemctl stop holzcloud && cp /opt/holzcloud/data/pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite /opt/holzcloud/data/holzcloud.sqlite && systemctl start holzcloud"}
holzcloud.service: Main process exited, code=exited, status=1/FAILURE
holzcloud.service: Failed with result 'exit-code'.
holzcloud.service: Scheduled restart job, restart counter is at 5.
Stopped holzcloud.service - holzcloud-CMS.
holzcloud.service: Start request repeated too quickly.
holzcloud.service: Failed with result 'exit-code'.
Failed to start holzcloud.service - holzcloud-CMS.
```

#### amd64-t34-rb/rb-snapshots.txt

```text
the pre-upgrade snapshots after the restart loop, in name order, each with its level (immutable read):
/opt/holzcloud/data/pre-upgrade-v0.1.4-255-g00854be-20261121-131835.sqlite level 0
/opt/holzcloud/data/pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite level 69
```

#### amd64-t34-rb/rb-guard.txt

```text
Started holzcloud.service - holzcloud-CMS.
{"time":"2026-11-21T13:24:18.836233877Z","level":"INFO","msg":"starting","version":"v0.1.4-255-g00854be","commit":"00854be"}
{"time":"2026-11-21T13:24:18.839410918Z","level":"INFO","msg":"configuration loaded","config":{"port":"8080","data_dir":"/opt/holzcloud/data","db_path":"/opt/holzcloud/data/holzcloud.sqlite","log_level":"INFO","secure":true,"max_template_size":10485760,"max_media_size":5242880,"max_video_size":67108864,"max_megapixels":24,"argon2_memory_kb":65536,"argon2_iterations":1,"argon2_parallelism":2,"min_free_bytes":536870912,"smtp_host":"","smtp_from":"","smtp_configured":false,"payrexx_instance":"","payrexx_configured":false,"trusted_proxies":"127.0.0.1/32,::1/128","listen":"127.0.0.1","sso_enabled":false,"sso_configured":false,"sso_provision":false,"sso_admin_group":"","sso_default_website":0,"sso_website_groups":0,"sso_sign_out_path":"/outpost.goauthentik.io/sign_out","oidc_enabled":false,"oidc_issuer":"","oidc_client_id":"","oidc_redirect_url":"","oidc_keys":0,"oidc_client_secret_set":false,"smtp":"aus","secret_key_set":false}}
{"time":"2026-11-21T13:24:18.888426418Z","level":"INFO","msg":"sqlite opened","journal_mode":"wal","path":"/opt/holzcloud/data/holzcloud.sqlite"}
{"time":"2026-11-21T13:24:18.899698418Z","level":"ERROR","msg":"refusing to start: the database schema is newer than this binary","db_level":70,"binary_level":69,"err":"the database is at migration level 70, and this binary knows migrations up to 69: a newer holzcloud has migrated it, and this binary refuses to run against a schema it does not know. The way back: start the newer binary again, or restore the pre-upgrade snapshot that the newer binary wrote before it migrated, together with the older binary, as deploy/DEPLOY.md describes under \"Rolling back\""}
holzcloud.service: Main process exited, code=exited, status=1/FAILURE
holzcloud.service: Failed with result 'exit-code'.
holzcloud.service: Scheduled restart job, restart counter is at 1.
Stopped holzcloud.service - holzcloud-CMS.
Started holzcloud.service - holzcloud-CMS.
{"time":"2026-11-21T13:24:24.101556754Z","level":"INFO","msg":"starting","version":"v0.1.4-255-g00854be","commit":"00854be"}
{"time":"2026-11-21T13:24:24.104332004Z","level":"INFO","msg":"configuration loaded","config":{"port":"8080","data_dir":"/opt/holzcloud/data","db_path":"/opt/holzcloud/data/holzcloud.sqlite","log_level":"INFO","secure":true,"max_template_size":10485760,"max_media_size":5242880,"max_video_size":67108864,"max_megapixels":24,"argon2_memory_kb":65536,"argon2_iterations":1,"argon2_parallelism":2,"min_free_bytes":536870912,"smtp_host":"","smtp_from":"","smtp_configured":false,"payrexx_instance":"","payrexx_configured":false,"trusted_proxies":"127.0.0.1/32,::1/128","listen":"127.0.0.1","sso_enabled":false,"sso_configured":false,"sso_provision":false,"sso_admin_group":"","sso_default_website":0,"sso_website_groups":0,"sso_sign_out_path":"/outpost.goauthentik.io/sign_out","oidc_enabled":false,"oidc_issuer":"","oidc_client_id":"","oidc_redirect_url":"","oidc_keys":0,"oidc_client_secret_set":false,"smtp":"aus","secret_key_set":false}}
{"time":"2026-11-21T13:24:24.147377921Z","level":"INFO","msg":"sqlite opened","journal_mode":"wal","path":"/opt/holzcloud/data/holzcloud.sqlite"}
{"time":"2026-11-21T13:24:24.163824129Z","level":"ERROR","msg":"refusing to start: the database schema is newer than this binary","db_level":70,"binary_level":69,"err":"the database is at migration level 70, and this binary knows migrations up to 69: a newer holzcloud has migrated it, and this binary refuses to run against a schema it does not know. The way back: start the newer binary again, or restore the pre-upgrade snapshot that the newer binary wrote before it migrated, together with the older binary, as deploy/DEPLOY.md describes under \"Rolling back\""}
holzcloud.service: Main process exited, code=exited, status=1/FAILURE
holzcloud.service: Failed with result 'exit-code'.
holzcloud.service: Scheduled restart job, restart counter is at 2.
Stopped holzcloud.service - holzcloud-CMS.
Started holzcloud.service - holzcloud-CMS.
{"time":"2026-11-21T13:24:29.336809965Z","level":"INFO","msg":"starting","version":"v0.1.4-255-g00854be","commit":"00854be"}
{"time":"2026-11-21T13:24:29.339474382Z","level":"INFO","msg":"configuration loaded","config":{"port":"8080","data_dir":"/opt/holzcloud/data","db_path":"/opt/holzcloud/data/holzcloud.sqlite","log_level":"INFO","secure":true,"max_template_size":10485760,"max_media_size":5242880,"max_video_size":67108864,"max_megapixels":24,"argon2_memory_kb":65536,"argon2_iterations":1,"argon2_parallelism":2,"min_free_bytes":536870912,"smtp_host":"","smtp_from":"","smtp_configured":false,"payrexx_instance":"","payrexx_configured":false,"trusted_proxies":"127.0.0.1/32,::1/128","listen":"127.0.0.1","sso_enabled":false,"sso_configured":false,"sso_provision":false,"sso_admin_group":"","sso_default_website":0,"sso_website_groups":0,"sso_sign_out_path":"/outpost.goauthentik.io/sign_out","oidc_enabled":false,"oidc_issuer":"","oidc_client_id":"","oidc_redirect_url":"","oidc_keys":0,"oidc_client_secret_set":false,"smtp":"aus","secret_key_set":false}}
{"time":"2026-11-21T13:24:29.383634798Z","level":"INFO","msg":"sqlite opened","journal_mode":"wal","path":"/opt/holzcloud/data/holzcloud.sqlite"}
{"time":"2026-11-21T13:24:29.393967757Z","level":"ERROR","msg":"refusing to start: the database schema is newer than this binary","db_level":70,"binary_level":69,"err":"the database is at migration level 70, and this binary knows migrations up to 69: a newer holzcloud has migrated it, and this binary refuses to run against a schema it does not know. The way back: start the newer binary again, or restore the pre-upgrade snapshot that the newer binary wrote before it migrated, together with the older binary, as deploy/DEPLOY.md describes under \"Rolling back\""}
holzcloud.service: Main process exited, code=exited, status=1/FAILURE
holzcloud.service: Failed with result 'exit-code'.
holzcloud.service: Scheduled restart job, restart counter is at 3.
Stopped holzcloud.service - holzcloud-CMS.
Started holzcloud.service - holzcloud-CMS.
{"time":"2026-11-21T13:24:34.589599551Z","level":"INFO","msg":"starting","version":"v0.1.4-255-g00854be","commit":"00854be"}
{"time":"2026-11-21T13:24:34.592304051Z","level":"INFO","msg":"configuration loaded","config":{"port":"8080","data_dir":"/opt/holzcloud/data","db_path":"/opt/holzcloud/data/holzcloud.sqlite","log_level":"INFO","secure":true,"max_template_size":10485760,"max_media_size":5242880,"max_video_size":67108864,"max_megapixels":24,"argon2_memory_kb":65536,"argon2_iterations":1,"argon2_parallelism":2,"min_free_bytes":536870912,"smtp_host":"","smtp_from":"","smtp_configured":false,"payrexx_instance":"","payrexx_configured":false,"trusted_proxies":"127.0.0.1/32,::1/128","listen":"127.0.0.1","sso_enabled":false,"sso_configured":false,"sso_provision":false,"sso_admin_group":"","sso_default_website":0,"sso_website_groups":0,"sso_sign_out_path":"/outpost.goauthentik.io/sign_out","oidc_enabled":false,"oidc_issuer":"","oidc_client_id":"","oidc_redirect_url":"","oidc_keys":0,"oidc_client_secret_set":false,"smtp":"aus","secret_key_set":false}}
{"time":"2026-11-21T13:24:34.635971551Z","level":"INFO","msg":"sqlite opened","journal_mode":"wal","path":"/opt/holzcloud/data/holzcloud.sqlite"}
{"time":"2026-11-21T13:24:34.646222176Z","level":"ERROR","msg":"refusing to start: the database schema is newer than this binary","db_level":70,"binary_level":69,"err":"the database is at migration level 70, and this binary knows migrations up to 69: a newer holzcloud has migrated it, and this binary refuses to run against a schema it does not know. The way back: start the newer binary again, or restore the pre-upgrade snapshot that the newer binary wrote before it migrated, together with the older binary, as deploy/DEPLOY.md describes under \"Rolling back\""}
holzcloud.service: Main process exited, code=exited, status=1/FAILURE
holzcloud.service: Failed with result 'exit-code'.
holzcloud.service: Scheduled restart job, restart counter is at 4.
Stopped holzcloud.service - holzcloud-CMS.
Started holzcloud.service - holzcloud-CMS.
{"time":"2026-11-21T13:24:39.840858803Z","level":"INFO","msg":"starting","version":"v0.1.4-255-g00854be","commit":"00854be"}
{"time":"2026-11-21T13:24:39.843428387Z","level":"INFO","msg":"configuration loaded","config":{"port":"8080","data_dir":"/opt/holzcloud/data","db_path":"/opt/holzcloud/data/holzcloud.sqlite","log_level":"INFO","secure":true,"max_template_size":10485760,"max_media_size":5242880,"max_video_size":67108864,"max_megapixels":24,"argon2_memory_kb":65536,"argon2_iterations":1,"argon2_parallelism":2,"min_free_bytes":536870912,"smtp_host":"","smtp_from":"","smtp_configured":false,"payrexx_instance":"","payrexx_configured":false,"trusted_proxies":"127.0.0.1/32,::1/128","listen":"127.0.0.1","sso_enabled":false,"sso_configured":false,"sso_provision":false,"sso_admin_group":"","sso_default_website":0,"sso_website_groups":0,"sso_sign_out_path":"/outpost.goauthentik.io/sign_out","oidc_enabled":false,"oidc_issuer":"","oidc_client_id":"","oidc_redirect_url":"","oidc_keys":0,"oidc_client_secret_set":false,"smtp":"aus","secret_key_set":false}}
{"time":"2026-11-21T13:24:39.887708928Z","level":"INFO","msg":"sqlite opened","journal_mode":"wal","path":"/opt/holzcloud/data/holzcloud.sqlite"}
{"time":"2026-11-21T13:24:39.898787595Z","level":"ERROR","msg":"refusing to start: the database schema is newer than this binary","db_level":70,"binary_level":69,"err":"the database is at migration level 70, and this binary knows migrations up to 69: a newer holzcloud has migrated it, and this binary refuses to run against a schema it does not know. The way back: start the newer binary again, or restore the pre-upgrade snapshot that the newer binary wrote before it migrated, together with the older binary, as deploy/DEPLOY.md describes under \"Rolling back\""}
holzcloud.service: Main process exited, code=exited, status=1/FAILURE
holzcloud.service: Failed with result 'exit-code'.
holzcloud.service: Scheduled restart job, restart counter is at 5.
Stopped holzcloud.service - holzcloud-CMS.
holzcloud.service: Start request repeated too quickly.
holzcloud.service: Failed with result 'exit-code'.
Failed to start holzcloud.service - holzcloud-CMS.
```

#### amd64-t34-rb/rb-rollback.out

```text
$ ls -l /opt/holzcloud/data/pre-upgrade-*.sqlite
1795267491.957829961 -rw------- 1 holzcloud holzcloud  12288 Nov 21 13:18 /opt/holzcloud/data/pre-upgrade-v0.1.4-255-g00854be-20261121-131835.sqlite
1795267491.975298753 -rw------- 1 holzcloud holzcloud 765952 Nov 21 13:23 /opt/holzcloud/data/pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite
WALL 1795267491.887905045 1795267491.996640211
RC 0
$ for f in /opt/holzcloud/data/pre-upgrade-*.sqlite; do echo "$f level $(sudo -u holzcloud sqlite3 "file:$f?immutable=1" 'select coalesce(max(version_id), 0) from goose_db_version where is_applied = 1')"; done
1795267492.257012545 /opt/holzcloud/data/pre-upgrade-v0.1.4-255-g00854be-20261121-131835.sqlite level 0
1795267492.367905295 /opt/holzcloud/data/pre-upgrade-v0.1.4-256-g4c07745-20261121-132341.sqlite level 69
WALL 1795267492.098922795 1795267492.390066170
RC 0
$ sudo systemctl stop holzcloud
WALL 1795267492.485325878 1795267492.626402920
RC 0
$ sudo cp -p /opt/holzcloud/holzcloud.prev /opt/holzcloud/holzcloud
WALL 1795267492.731102670 1795267492.861221795
RC 0
$ sudo rm -f /opt/holzcloud/data/holzcloud.sqlite-wal /opt/holzcloud/data/holzcloud.sqlite-shm
WALL 1795267492.955823004 1795267493.079744754
RC 0
$ sudo mv /opt/holzcloud/data/holzcloud.sqlite /opt/holzcloud/data/holzcloud.sqlite.after-update
WALL 1795267493.176032545 1795267493.299753920
RC 0
$ sudo -u holzcloud cp "$(ls -t /opt/holzcloud/data/pre-upgrade-*.sqlite | head -n 1)" /opt/holzcloud/data/holzcloud.sqlite
WALL 1795267493.415938504 1795267493.574246254
RC 0
$ sudo systemctl reset-failed holzcloud
WALL 1795267493.670746629 1795267493.805904587
RC 0
$ sudo systemctl start holzcloud
WALL 1795267493.899246421 1795267494.048696879
RC 0
$ curl -fsS --retry 30 --retry-connrefused --retry-delay 1 http://127.0.0.1:8080/readyz | jq '{status, version, commit, integrity}'
1795267494.252768171 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 4 ms: Couldn't connect to server
1795267495.248491255 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267496.250912588 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267497.257871297 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267498.255082839 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267499.258521007 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267500.262583882 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267501.263443299 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267502.271857133 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267503.273667050 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267504.278333009 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267505.317362593 {
1795267505.333816593   "status": "ok",
1795267505.350889218   "version": "v0.1.4-255-g00854be",
1795267505.367354843   "commit": "00854be",
1795267505.383830551   "integrity": "ok"
1795267505.400822468 }
WALL 1795267494.151423337 1795267505.422630301
RC 0
$ sudo journalctl -u holzcloud --since "$(systemctl show -p ActiveEnterTimestamp --value holzcloud)" -o cat | grep -E 'pre-upgrade snapshot (written|kept)|migration applied' || echo 'no migration, no snapshot'
1795267505.718159301 no migration, no snapshot
WALL 1795267505.524527426 1795267505.740274635
RC 0
```

#### amd64-t34-rb/rb-state.txt

```text
s0, before the update (B0 running, seeded):
s0 schema: 69
s0 tables: activity_log ai_tokens album_items albums app_settings block_types cart_items carts content_types csv_imports form_messages goose_db_version mail_outbox media media_usage media_variants menu_items menus oauth_clients oauth_codes order_items orders outbox page_field_defs page_revisions page_terms pages pages_fts pages_fts_config pages_fts_data pages_fts_docsize pages_fts_idx plugin_migrations plugin_store plugin_websites plugins privacy_requests product_media product_terms products redirects retention_deletions retention_periods retention_runs saved_views sessions site_values snippets templates terms theme_wording user_recovery_codes user_tokens user_websites users website_domains website_mail website_templates websites
s0 rows pages: 4
s0 rows users: 0
s0 rows media: 1
s0 rows orders: 0
s0 rows plugins: 5
s0 page: 722cb2f1d7d8ab8fb9c875251cf546d261f0225b006bb68d898d0791b1d27092
s0 media: 925d26e12cb8ed51e272c6cadba3f6fb5d32f4af8bf06d1d2123ddbf8a176577
s2, after Rolling back block 1:
s2 schema: 69
s2 tables: activity_log ai_tokens album_items albums app_settings block_types cart_items carts content_types csv_imports form_messages goose_db_version mail_outbox media media_usage media_variants menu_items menus oauth_clients oauth_codes order_items orders outbox page_field_defs page_revisions page_terms pages pages_fts pages_fts_config pages_fts_data pages_fts_docsize pages_fts_idx plugin_migrations plugin_store plugin_websites plugins privacy_requests product_media product_terms products redirects retention_deletions retention_periods retention_runs saved_views sessions site_values snippets templates terms theme_wording user_recovery_codes user_tokens user_websites users website_domains website_mail website_templates websites
s2 rows pages: 4
s2 rows users: 0
s2 rows media: 1
s2 rows orders: 0
s2 rows plugins: 5
s2 page: 722cb2f1d7d8ab8fb9c875251cf546d261f0225b006bb68d898d0791b1d27092
s2 media: 925d26e12cb8ed51e272c6cadba3f6fb5d32f4af8bf06d1d2123ddbf8a176577
comparison line by line:
bad=0
```

#### amd64-t34-evidence-commit.txt

```text
00854be
```

#### amd64-t34-pass1.status

```text
0
```
