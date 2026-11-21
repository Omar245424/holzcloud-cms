# A migration that fails part-way, and the rollback by snapshot: results

This document is the permanent record of UPD-02 of phase 34. It records an update from the build of this release (B0, migration level 69) to a build whose first added migration succeeds and whose second fails, played with the *Updating* blocks of `deploy/DEPLOY.md`: that the service stops with the failing migration named in the journal, that the restart loop writes one pre-upgrade snapshot before the first attempt and keeps it on every restart, that every `restore with:` line names that one file, that the update's own check line names it too, that the previous binary refuses the half-migrated database with the schema guard's sentence and both levels, and that *Rolling back* block 1, played as written, puts the previous binary and the snapshot back and starts the previous binary on the state it had before the update, row for row.

The record rests on one evidence set: a single `upgrade-verify.sh run` from fresh containers, made once against the committed text of every bound file (see *How this was run*); this record uses its stage `rb`, the same set as `deploy/MIGRATION-RESTORE-RUN.md`. Its raw output is the Appendix, pasted unchanged, and `bash deploy/test/upgrade-verify.sh record-check` proves every block against its scratch file byte for byte. Every number in this document is a line of that set and cites its file in backticks as `DIRECTORY/FILE`, with `DIRECTORY` the one scratch directory `amd64-t34-rb`; the run files beside it are named by their path below `deploy/test/out`. A file that is not pasted is named as `FILE in DIRECTORY`.

The record says what was done, how many times, and what was seen. The failing build is a throwaway build of this tree, not a release, and the failure is a migration that names a table that does not exist; other failures are named in *Not tested*. What the container cannot show is a named deviation or a line of *Not tested*, never a gap in the table.

## Summary

| Requirement | Result | Platform | Date, commit |
|-------------|--------|----------|--------------|
| UPD-02 | PASS: `UPD-02 PASS` stands in `amd64-t34-summary.txt` and in `amd64-t34-rb/rb.txt`. The update to F stopped with `STOP-MESSAGE: goose up: partial migration error (type:sql,version:71): SQL logic error: no such table: t34_no_such_table (1)` after `APPLIED-BEFORE-FAILURE: 70`; the restart loop ran to its end without a reset (`STARTS: 5`, `ACTIVE-STATE: failed`, `START-LIMIT-LINE: yes`, `HEALTHZ-DURING-LOOP: never`) and wrote one snapshot at level 69 that every later start kept (`JOURNAL-SNAPSHOT-WRITTEN: 1`, `JOURNAL-SNAPSHOT-KEPT: 4`, `WRITTEN-LEVEL: 69`), named by all five `restore with:` lines (`JOURNAL-RESTORE-WITH: 5`, `RESTORE-WITH-DISTINCT-PATHS: 1`, `RESTORE-WITH-NAMES-WRITTEN: yes`) and by the update's check line (`UPDATE-CHECK-SNAPSHOT-LINE: yes`); the previous binary refused the database at level 70 (`GUARD-LINE: yes`, `GUARD-DB-LEVEL: 70`, `GUARD-BINARY-LEVEL: 69`, `GUARD-SENTENCE: yes`, `GUARD-HEALTHZ: never`); *Rolling back* block 1 ended with every line at status 0 (`ROLLBACK-RC: 0 0 0 0 0 0 0 0 0 0 0`), the previous build answered `/readyz` with `status` and `integrity` `ok` at `LEVEL-AFTER-ROLLBACK: 69`, and the state after equals the state before (`COMPARE-S0-S2: bad=0`), all in `amd64-t34-rb/rb.txt`. | Debian 12 container, systemd PID 1, linux/amd64 under emulation on an arm64 host (`amd64-t34-rb/versions.txt`) | 2026-11-21, cefe4ac |

Evidence commit: cefe4ac

The previous binary, B0, is the one built from `NEW_REF` (pinned by `run` to the short sha of HEAD, `cefe4ac`, read from `deploy/test/out/amd64-t34-evidence-commit.txt`) with the Build line of `deploy/DEPLOY.md`: `BINARY: holzcloud v0.1.4-261-gcefe4ac (cefe4ac)` in `amd64-t34-summary.txt`, `B0-VERSION: v0.1.4-261-gcefe4ac` in `amd64-t34-rb/rb.txt`, and the line `holzcloud v0.1.4-261-gcefe4ac (cefe4ac)` that `/opt/holzcloud/holzcloud version` printed in *Updating* block 1 (`amd64-t34-rb/rb-update.out`). The failing build, F, is `F-BUILD: v0.1.4-262-gd8a027e (d8a027e) migrations 71 sha256 1bf74347445452df41d1b5425bc0a50a286e77a866ffced35f9de352121bdda9` (`amd64-t34-rb/rb.txt`, the same line as `holzcloud-f:` in `amd64-t34-rb/versions.txt`): a throwaway build of this tree with one local commit that adds two migrations, made in a scratch clone and deleted with it (`DEVIATION[synthetic-build]`), so its commit `d8a027e` exists nowhere now. The levels are computed from the migration files, never typed: `LEVEL-EXPECTED-B0: 69` and `LEVEL-EXPECTED-F: 71`. The summary carries `STAGES: mig upd rb` and no `REHEARSAL` line.

## How this was run

`deploy/test/upgrade-verify.sh` is a sibling of `deploy/test/run.sh` in the shape of `deploy/test/backup-verify.sh`; it edits neither and reaches the harness through `run.sh`'s public subcommands only. `run` refuses to start unless `PLATFORM=linux/amd64` and `TARGET=systemd-base`, refuses while a bound file or the binary's source differs from HEAD or is untracked, pins `NEW_REF` to the short sha of HEAD so that the binary and the scripts are one commit, and plays the stages `mig`, `upd` and `rb`, each in its own fresh containers. Every line the stage `rb` plays from *Updating* and *Rolling back* is read out of the fenced bash blocks of `deploy/DEPLOY.md` at that commit; a line the driver cannot translate ends the stage. The verdict is computed from the `KEY: value` file the stage wrote, each key exactly once, so a trace can never pass for the set. Stage `rb` used one container, `holzcloud-t34-rb`.

The set was made once, from fresh containers, on the clean committed tree of `cefe4ac` without `REHEARSAL`, started 2026-11-21T14:13:05Z and ended 14:24:59Z on the host clock (11 minutes 54 seconds for all three stages), exit status 0 (`deploy/test/out/amd64-t34-pass1.started`, `deploy/test/out/amd64-t34-pass1.finished` and `deploy/test/out/amd64-t34-pass1.status`). Before it the worktree held no earlier output of the driver and no `holzcloud-t34-*` container existed, and after it `bash deploy/test/upgrade-verify.sh down` removed the four containers. In the same window and at the same commit, `PLATFORM=linux/amd64 TARGET=systemd-base bash deploy/test/backup-verify.sh run` made the set of phase 32 again (14:25:11Z to 14:27:08Z, status 0). No second set was made against `cefe4ac`, and no other set is cited. Two earlier sets were superseded. The first was made against `36b67cc`. Its stage `rb` read `UPD-02 PASS`, but the update's check line of *Updating* step 3 printed `no migration, no snapshot` after the failed update, because it matched `pre-upgrade snapshot written` only and the restarts since `6e86cdf` log `pre-upgrade snapshot kept`. The commit `00854be` fixed that line in *Updating* and in *Rolling back*, and made the driver record it (`UPDATE-CHECK-SNAPSHOT-LINE`) and require it for the verdict. The second was made against `00854be`, and this record was first written from it. Its four verdicts read PASS, but the browser record of phase 34 then found that a plugin switched on in the admin did not start when the browser dropped the request that switched it on, because the plugin was loaded with that request's context and closed with it. The commit `cefe4ac` fixed that in `internal/plugin/runtime.go`, so the binary changed after the set of `00854be`, and the set was made again. So this record was written from this set, and nothing of either earlier set is quoted. The command was:

```text
PLATFORM=linux/amd64 TARGET=systemd-base bash deploy/test/upgrade-verify.sh run
```

The knob not set on that line had its default: `STAGES=mig upd rb`. The log of the run is `deploy/test/out/amd64-t34-pass1.log`, which is not pasted; its first line reads `upgrade-verify: platform linux/amd64, NEW_REF pinned to cefe4ac, v0.1.1 is 5bcae9b448db814701310c5dfab26031d68c0d92, stages: mig upd rb`.

No bound file changed after the evidence commit. The bound files for this set are `deploy/DEPLOY.md`, `deploy/backup.sh`, `deploy/restore.sh`, `deploy/holzcloud.service`, `deploy/test/run.sh`, `deploy/test/Dockerfile`, `deploy/test/upgrade-verify.sh` and the five plugin archives under `plugins/`, together with the binary built from `cmd`, `internal`, `go.mod`, `go.sum`, `changelog.go` and the `CHANGELOG.md` it embeds. This command exited 0 when the record was written:

```text
git diff --quiet cefe4ac HEAD -- deploy/DEPLOY.md deploy/backup.sh deploy/restore.sh deploy/holzcloud.service deploy/test/run.sh deploy/test/Dockerfile deploy/test/upgrade-verify.sh plugins/bestellung/bestellung.zip plugins/jahreszahl/jahreszahl.zip plugins/kontaktformular/kontaktformular.zip plugins/nicht-gefunden/nicht-gefunden.zip plugins/suche/suche.zip cmd internal go.mod go.sum changelog.go CHANGELOG.md
```

Before the run, `bash deploy/test/upgrade-verify.sh selftest` ended `SELFTEST OK (409 cases)`; its output is not part of the set.

The stage in order: B0 was installed by the document's Install text and seeded, and its state captured as `s0`. The driver then waited until the last start of the unit was older than its start-limit interval (`START-WINDOW-WAIT-S: 274` in `amd64-t34-rb/rb.txt`), a wait and not a reset, so that the restart loop had its full five starts as on a server where an update comes long after the last start. *Updating* blocks 1 to 3 were played line by line (`amd64-t34-rb/rb-update.out`), and no `reset-failed` ran between the failing start and the end of its restart loop. The journal of the loop was kept (`amd64-t34-rb/rb-journal.txt`) and the snapshots listed with their levels (`amd64-t34-rb/rb-snapshots.txt`). Then `SCENARIO[previous-on-half]` started the previous binary on the half-migrated database (`amd64-t34-rb/rb-guard.txt`). *Rolling back* block 1 was played as written (`amd64-t34-rb/rb-rollback.out`), and the state after it, `s2`, was compared with `s0` line by line (`amd64-t34-rb/rb-state.txt`).

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

Each finding has the class of the earlier records' rule: (1) a bug of the harness, fixed in the driver; (2) a limit of the container, a named deviation or a *Not tested* line; (3) a measured property of a shipped file that matches what it and its document say, recorded and not edited; (4) a defect, a shipped file doing less than it or its document says, with a fix in a commit of its own. No finding of this set is class 1 or class 4, and no shipped file was edited after the evidence commit. The defect that the earlier set at `36b67cc` showed in the check line of *Updating* step 3 was fixed in `00854be` before this set (*How this was run*), and F-6 is that line observed on the fixed text. The defect that the browser record found at `00854be`, a plugin load ended by a dropped request, was fixed in `cefe4ac` before this set.

The test for class 4 was whether a result contradicts something `deploy/DEPLOY.md` or the binary says at `cefe4ac`, or leaves the operator worse off than they say. The document says that a failed update is restarted every five seconds until the start limit, that from 0.1.5 the restarts write no further snapshot and log `pre-upgrade snapshot kept` with the same `path`, that every `restore with:` line names the file written before the first attempt, that an older binary refuses a newer database with `db_level` and `binary_level`, and that *Rolling back* block 1 restores the snapshot with the old binary and ends with `no migration, no snapshot`. F-1 to F-7 are those sentences observed.

**F-1 (class 3): The restart loop wrote one snapshot, before the first attempt, and kept it on every restart; the five `restore with:` lines all name it.** The first start of F logged `pre-upgrade snapshot written` with `"level":69`, and the four starts after it each logged `pre-upgrade snapshot kept` with the same `path`, `"level":69`, `"db_level":70` and the reason (`amd64-t34-rb/rb-journal.txt`; `JOURNAL-SNAPSHOT-WRITTEN: 1`, `JOURNAL-SNAPSHOT-KEPT: 4`, `JOURNAL-RESTORE-WITH: 5`, `RESTORE-WITH-DISTINCT-PATHS: 1`, `RESTORE-WITH-NAMES-WRITTEN: yes` in `amd64-t34-rb/rb.txt`). After the loop the data directory held that snapshot at level 69 and the snapshot of the install's first start at level 0, and none at level 70 (`SNAPSHOTS-AFTER-LOOP` in `amd64-t34-rb/rb.txt`, `amd64-t34-rb/rb-snapshots.txt`); while the loop ran, `pre-upgrade.pending` named the snapshot and its level (`MARKER-DURING: pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite 69`). This is what the fix `6e86cdf` ("restarts after a failed migration keep the snapshot before it instead of snapshotting the half-migrated database") does, and it answers the probe of finding F-1 of phase 34's research, in which the code before it wrote a further snapshot at level 70 on each of starts 2 to 5, and named those half-migrated copies in its `restore with:` lines. The line of the first start and the line of a restart (`amd64-t34-rb/rb-journal.txt`):

```text
{"time":"2026-11-21T14:23:28.439671593Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite","level":69}
{"time":"2026-11-21T14:23:33.726664137Z","level":"WARN","msg":"pre-upgrade snapshot kept","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite","level":69,"db_level":70,"reason":"an earlier start of this update applied part of its migrations and did not finish; the snapshot from before it stays the way back, so no snapshot of the half-migrated database is written"}
```

**F-2 (class 3): The journal names migration 70 as applied before the failure of 71, so the half-migrated level is in the journal and not only in the database.** The first start's lines after the snapshot (`amd64-t34-rb/rb-journal.txt`; `APPLIED-BEFORE-FAILURE: 70`, `LEVEL-HALF: 70` in `amd64-t34-rb/rb.txt`):

```text
{"time":"2026-11-21T14:23:28.487366218Z","level":"INFO","msg":"migration applied","version":70,"type":"sql","duration_ms":4}
{"time":"2026-11-21T14:23:28.488909052Z","level":"ERROR","msg":"migrations failed","err":"goose up: partial migration error (type:sql,version:71): SQL logic error: no such table: t34_no_such_table (1)"}
{"time":"2026-11-21T14:23:28.489365427Z","level":"ERROR","msg":"restore with: systemctl stop holzcloud && cp /opt/holzcloud/data/pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite /opt/holzcloud/data/holzcloud.sqlite && systemctl start holzcloud"}
```

**F-3 (class 3): systemd's start limit ended the loop after five starts, and the unit stayed failed with nothing answering.** The journal counts the restarts to 5 and ends with `Start request repeated too quickly.` (`amd64-t34-rb/rb-journal.txt`); `STARTS: 5`, `NRESTARTS: 5`, `RESULT: exit-code`, `ACTIVE-STATE: failed`, `START-LIMIT-LINE: yes` and `HEALTHZ-DURING-LOOP: never` (`amd64-t34-rb/rb.txt`). The starts came about five seconds apart, from 14:23:28Z to 14:23:50Z.

**F-4 (class 3): The previous binary, started on the half-migrated database, refused it with the schema guard's sentence and both levels, five times, and never answered.** The sentence is the one `internal/db/schema_guard.go` builds, and it names the way back (`amd64-t34-rb/rb-guard.txt`; `GUARD-LINE: yes`, `GUARD-DB-LEVEL: 70`, `GUARD-BINARY-LEVEL: 69`, `GUARD-SENTENCE: yes`, `GUARD-STARTS: 5`, `GUARD-HEALTHZ: never`, `GUARD-ACTIVE-STATE: failed` in `amd64-t34-rb/rb.txt`):

```text
{"time":"2026-11-21T14:24:05.267885888Z","level":"ERROR","msg":"refusing to start: the database schema is newer than this binary","db_level":70,"binary_level":69,"err":"the database is at migration level 70, and this binary knows migrations up to 69: a newer holzcloud has migrated it, and this binary refuses to run against a schema it does not know. The way back: start the newer binary again, or restore the pre-upgrade snapshot that the newer binary wrote before it migrated, together with the older binary, as deploy/DEPLOY.md describes under \"Rolling back\""}
```

**F-5 (class 3): The first two lines of *Rolling back* block 1 listed both snapshots with their levels, and the restore line took the one at level 69.** The level list (`amd64-t34-rb/rb-rollback.out`):

```text
1795271079.539383168 /opt/holzcloud/data/pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141823.sqlite level 0
1795271080.005354459 /opt/holzcloud/data/pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite level 69
```

The file at level 0 is the snapshot of the install's first start, 12288 bytes by the `ls -l` line above it, the "snapshot of a first start, an empty database" that the document warns of; the newest by modification time is the update's own, at the level of the old binary, so the `ls -t` of the restore line took the right file.

**F-6 (class 3): After the failed update, the check line of *Updating* step 3 printed the `pre-upgrade snapshot kept` line of the last start, naming the snapshot, as the document now says.** Line 1 of the step failed after its retries with status 7, as it should (`UPDATE-CHECK-RC: 7`), and line 2 printed (`amd64-t34-rb/rb-update.out`; `UPDATE-CHECK-SNAPSHOT-LINE: yes` in `amd64-t34-rb/rb.txt`):

```text
1795271038.903466885 {"time":"2026-11-21T14:23:50.217708714Z","level":"WARN","msg":"pre-upgrade snapshot kept","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite","level":69,"db_level":70,"reason":"an earlier start of this update applied part of its migrations and did not finish; the snapshot from before it stays the way back, so no snapshot of the half-migrated database is written"}
```

It shows the line of the last start only, because `ActiveEnterTimestamp` names the last start of the loop; the `written` line of the first start is earlier in the journal (F-1).

**F-7 (class 3): *Rolling back* block 1 played as written put the previous build back on the state before the update, and kept the migrated database aside.** Every line ended with status 0 (`ROLLBACK-RC: 0 0 0 0 0 0 0 0 0 0 0`), `/readyz` showed the previous build with `status` and `integrity` `ok` (`ROLLBACK-READYZ`), and the last line printed `no migration, no snapshot` (`ROLLBACK-CHECK-LINE`). The database is at 69 without the table of migration 70 (`LEVEL-AFTER-ROLLBACK: 69`, `T34-OK-TABLE: absent`), `holzcloud.sqlite.after-update` holds level 70 (`AFTER-UPDATE-LEVEL: 70`), and the marker is gone (`MARKER-AFTER: absent`), all in `amd64-t34-rb/rb.txt`. The state before the update and after the rollback, schema, tables, the rows of five tables and the hashes of the seeded page and picture, compared line by line with no difference (`amd64-t34-rb/rb-state.txt`, `COMPARE-S0-S2: bad=0`).

## UPD-02: a migration that fails part-way, and the way back

**What was done.** B0 was installed and seeded, and its state `s0` read through a read-only connection (`amd64-t34-rb/rb-state.txt`). *Updating* block 1 copied the binary aside as `holzcloud.prev`, printed its version and made the backup with the document's line, ending `=== Backup complete: /opt/holzcloud/backups/20261121-142325 ===`; block 2 skipped the Build line (`DEVIATION[synthetic-build]`), carried F and the two scripts over, stopped the service, copied F into place, installed the scripts and started it, each line with status 0 (`amd64-t34-rb/rb-update.out`). The first start of F wrote the snapshot at level 69, applied 70, failed at 71 and printed its `restore with:` line (F-1, F-2); systemd restarted it four times, each restart keeping the snapshot and failing at 71 again, until the start limit (F-3; `amd64-t34-rb/rb-journal.txt`). Block 3 found nothing listening (`UPDATE-CHECK-RC: 7`) and its journal line named the kept snapshot (F-6). The previous binary was then started on the half-migrated database and refused it five times (`SCENARIO[previous-on-half]`, F-4; `amd64-t34-rb/rb-guard.txt`). Last, *Rolling back* block 1 was played as written, its own `reset-failed` included: the two listing lines (F-5), the stop, the previous binary back, the write-ahead files removed, the migrated database set aside, the snapshot copied in as `holzcloud`, the reset, the start and the two check lines (`amd64-t34-rb/rb-rollback.out`):

```text
1795271095.581295050   "status": "ok",
1795271095.598087425   "version": "v0.1.4-261-gcefe4ac",
1795271095.616134800   "commit": "cefe4ac",
1795271095.633862675   "integrity": "ok"
1795271095.984705384 no migration, no snapshot
```

Then `s2` was read and compared with `s0` (F-7).

*What it shows.* When an update to a build whose migration 71 failed after migration 70 had run, the service stopped with the failing migration named in the journal and was restarted by systemd until its start limit, every restart kept the one pre-upgrade snapshot written before the first attempt at level 69 and every `restore with:` line named it, the previous binary refused the database at level 70 with the schema guard's sentence, and *Rolling back* block 1, played as written, started the previous binary on a database whose schema, tables, rows and files equal those before the update.

**What it does not show.** A migration that fails for another reason, a release as the failing build, an operator who starts the failing build again, and the files the snapshot does not hold (*Not tested*).

## Measurements

The figures the harness wrote. Every one is an emulation figure on an aarch64 host (`amd64-t34-rb/versions.txt`); none is a promise of speed.

| Figure | Value | Where |
|--------|-------|-------|
| The wait until the last start was outside the start-limit interval | 274 s | `START-WINDOW-WAIT-S` in `amd64-t34-rb/rb.txt` |
| The backup of *Updating* block 1 | 1795271005.759755842 to 1795271006.874063384 wall, 21M | the `WALL` line and the size line of `amd64-t34-rb/rb-update.out` |
| The first start of F, from `starting` to `restore with:` | 14:23:28.328734635Z to 14:23:28.489365427Z | `amd64-t34-rb/rb-journal.txt` |
| The restart loop, first to fifth start | 14:23:28.328734635Z to 14:23:50.157396673Z | `amd64-t34-rb/rb-journal.txt` |
| The readiness check of *Updating* step 3, 31 attempts | 1795271008.343704052 to 1795271038.606147885 wall | `amd64-t34-rb/rb-update.out` |
| The refusals of the previous binary, first to fifth start | 14:24:05.207931888Z to 14:24:26.407722926Z | `amd64-t34-rb/rb-guard.txt` |
| *Rolling back* block 1, first line to last | 1795271077.928139042 to 1795271096.007424134 wall | `amd64-t34-rb/rb-rollback.out` |
| The snapshot of the update and that of the first start | 765952 bytes and 12288 bytes | the `ls -l` lines of `amd64-t34-rb/rb-rollback.out` |
| The whole run, all three stages | 14:13:05Z to 14:24:59Z | `deploy/test/out/amd64-t34-pass1.started`, `deploy/test/out/amd64-t34-pass1.finished` |

## Container vs. server

| Aspect | In the container | On a real server | Measured or assumed |
|--------|------------------|------------------|---------------------|
| Speed and timing | amd64 emulated on aarch64 | A native server | Not measured natively |
| The failing build | A throwaway build of this tree, two added migrations, the second inserting into a missing table | A release whose migration fails for its own reason | Measured for this failure; others assumed to take the same path in the binary |
| The start window | A wait of 274 s so that the loop got its five starts | An update long after the last start | Measured here |
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
BINARY: holzcloud v0.1.4-261-gcefe4ac (cefe4ac)
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
running binary:        holzcloud v0.1.4-261-gcefe4ac (cefe4ac)
holzcloud-f:  v0.1.4-262-gd8a027e (d8a027e) migrations 71 sha256 1bf74347445452df41d1b5425bc0a50a286e77a866ffced35f9de352121bdda9
NEW_REF:               cefe4ac
evidence commit:       cefe4ac
```

#### amd64-t34-rb/deviations.txt

```text
DEVIATION[fixture]: the five plugin archives of this tree are copied into /opt/holzcloud/plugin-zips before run.sh seed; with the website, page, AI key and test image run.sh seed makes, they stand in for operator work in the admin and are no step of the document
DEVIATION[synthetic-build]: the binary of this update is F (v0.1.4-262-gd8a027e), a throwaway build of this tree with two added migrations, 00070_t34_ok.sql (creates t34_ok) and 00071_t34_fail.sql (creates t34_half, then inserts into a table that does not exist), made in a scratch clone with one local commit and the Build line of the document, the clone deleted afterwards; nothing of it is committed or pushed. The Build line of Updating block 2 is skipped for it. Released binaries give no update whose migration fails part-way while the previous binary carries the schema guard
DEVIATION[scp]: 'scp holzcloud user@your-server:' was played as docker cp of the synthetic build of this update to /root/holzcloud in the container (there is no ssh and no second machine here)
DEVIATION[scp]: 'scp deploy/backup.sh deploy/restore.sh user@your-server:' was played as docker cp of the working tree's deploy/backup.sh deploy/restore.sh to /root/ in the container (there is no ssh and no second machine here)
SCENARIO[previous-on-half]: after the restart loop ended, systemctl reset-failed holzcloud, cp -p /opt/holzcloud/holzcloud.prev /opt/holzcloud/holzcloud and systemctl start holzcloud: the previous binary is started on the half-migrated database, which Rolling back says not to do, so that the record shows the schema guard refusing it
DEVIATION[start-limit]: systemctl reset-failed holzcloud runs once outside the document, before the previous binary is started on the half-migrated database (SCENARIO[previous-on-half]); the restart loop of the failing update ran to its end without one, and the reset-failed of Rolling back block 1 is the document's own line
DEVIATION[emulation]: the Docker host is aarch64; the linux/amd64 container runs under emulation, so every wall time of this set (the restarts of the unit among them) is an emulation figure, not a server's
DEVIATION[scp]: 'scp holzcloud user@your-server:' was played as docker cp of deploy/test/out/bin/holzcloud-new-cefe4ac to /root/holzcloud in the container (there is no ssh and no second machine here)
DEVIATION[scp]: 'scp -r deploy/ user@your-server:' was played as a tar stream of the working tree's deploy directory (without deploy/test/out) into /root in the container, by docker cp -
DEVIATION[dropin]: /etc/systemd/system/holzcloud.service.d/container.conf sets ProcSubset=all (container only; the document's unchanged start crashed, see unit-start-unchanged.txt)
DEVIATION[mount]: the unit's RequiresMountsFor names /opt/holzcloud/data, a plain directory on the container's root filesystem; no separate mount exists here, so that line was not exercised
```

#### amd64-t34-rb/rb.txt

```text
B0-VERSION: v0.1.4-261-gcefe4ac
LEVEL-EXPECTED-B0: 69
LEVEL-EXPECTED-F: 71
F-BUILD: v0.1.4-262-gd8a027e (d8a027e) migrations 71 sha256 1bf74347445452df41d1b5425bc0a50a286e77a866ffced35f9de352121bdda9
START-WINDOW-WAIT-S: 274
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
WRITTEN-PATH: /opt/holzcloud/data/pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite
WRITTEN-LEVEL: 69
STOP-MESSAGE: goose up: partial migration error (type:sql,version:71): SQL logic error: no such table: t34_no_such_table (1)
APPLIED-BEFORE-FAILURE: 70
START-LIMIT-LINE: yes
LEVEL-HALF: 70
MARKER-DURING: pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite 69
SNAPSHOTS-AFTER-LOOP: pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141823.sqlite:0 pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite:69
GUARD-HEALTHZ: never
GUARD-ACTIVE-STATE: failed
GUARD-STARTS: 5
GUARD-LINE: yes
GUARD-DB-LEVEL: 70
GUARD-BINARY-LEVEL: 69
GUARD-SENTENCE: yes
ROLLBACK-RC: 0 0 0 0 0 0 0 0 0 0 0
ROLLBACK-READYZ: {"status":"ok","version":"v0.1.4-261-gcefe4ac","commit":"cefe4ac","integrity":"ok"}
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
WALL 1795271005.289418550 1795271005.417645550
RC 0
$ /opt/holzcloud/holzcloud version
1795271005.639548842 holzcloud v0.1.4-261-gcefe4ac (cefe4ac)
WALL 1795271005.514869342 1795271005.661790884
RC 0
$ sudo -u holzcloud /opt/holzcloud/deploy/backup.sh /opt/holzcloud/backups
1795271005.922937842 === holzcloud-CMS backup: 20261121-142325 ===
1795271005.998759467 Backing up database...
1795271006.123143259 2026/11/21 14:23:26 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/data/holzcloud.sqlite
1795271006.168087426 2026/11/21 14:23:26 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/backups/20261121-142325/holzcloud.sqlite
1795271006.189036634 2026/11/21 14:23:26 INFO backup verified path=/opt/holzcloud/backups/20261121-142325/holzcloud.sqlite tables=60
1795271006.207558134 backup written and verified: /opt/holzcloud/backups/20261121-142325/holzcloud.sqlite
1795271006.224851342 Backing up CSRF key...
1795271006.241393967 Backing up media...
1795271006.319977384 Backing up plugins...
1795271006.457613342 Writing checksums...
1795271006.778485426 Pruning backups older than 30 days...
1795271006.835959843 === Backup complete: /opt/holzcloud/backups/20261121-142325 ===
1795271006.852496718 21M	/opt/holzcloud/backups/20261121-142325
WALL 1795271005.759755842 1795271006.874063384
RC 0
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -X main.Version=$(git describe --tags --always) -X main.Commit=$(git rev-parse --short HEAD)" ./cmd/holzcloud
SKIP[synthetic-build] the binary of this update was built by this line in a scratch clone before the update
RC 0
$ scp holzcloud user@your-server:
RC 0
$ scp deploy/backup.sh deploy/restore.sh user@your-server:
RC 0
$ sudo systemctl stop holzcloud
WALL 1795271007.365147259 1795271007.523004510
RC 0
$ sudo cp ~/holzcloud /opt/holzcloud/holzcloud
WALL 1795271007.628325218 1795271007.758131718
RC 0
$ sudo install -m 755 -o root -g root ~/backup.sh ~/restore.sh /opt/holzcloud/deploy/
WALL 1795271007.858881260 1795271007.987906718
RC 0
$ sudo systemctl start holzcloud
WALL 1795271008.084570926 1795271008.232905010
RC 0
$ curl -fsS --retry 30 --retry-connrefused --retry-delay 1 http://127.0.0.1:8080/readyz | jq '{status, version, commit, integrity}'
1795271008.437830677 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 4 ms: Couldn't connect to server
1795271009.447789802 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271010.455905261 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271011.459014178 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271012.467906720 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271013.475775554 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271014.479492179 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271015.488924263 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271016.486705972 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271017.488249083 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271018.494244167 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271019.508094001 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271020.515284668 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271021.528320002 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271022.537970794 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271023.535050044 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271024.539577003 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271025.543875420 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271026.549626546 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271027.550792713 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271028.555272089 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271029.559777256 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271030.565207131 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271031.566571423 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271032.577326132 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271033.573544799 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271034.590451675 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271035.580413009 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271036.582194842 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271037.583167385 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271038.584065760 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
WALL 1795271008.343704052 1795271038.606147885
RC 7
$ sudo journalctl -u holzcloud --since "$(systemctl show -p ActiveEnterTimestamp --value holzcloud)" -o cat | grep -E 'pre-upgrade snapshot (written|kept)|migration applied' || echo 'no migration, no snapshot'
1795271038.903466885 {"time":"2026-11-21T14:23:50.217708714Z","level":"WARN","msg":"pre-upgrade snapshot kept","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite","level":69,"db_level":70,"reason":"an earlier start of this update applied part of its migrations and did not finish; the snapshot from before it stays the way back, so no snapshot of the half-migrated database is written"}
WALL 1795271038.709792052 1795271038.926154635
RC 0
```

#### amd64-t34-rb/rb-journal.txt

```text
Stopping holzcloud.service - holzcloud-CMS...
{"time":"2026-11-21T14:23:27.480607676Z","level":"INFO","msg":"shutting down","timeout_s":10}
{"time":"2026-11-21T14:23:27.491617093Z","level":"INFO","msg":"database maintenance","checkpoint_busy":0,"db_size_bytes":770048,"wal_size_bytes":0,"page_count":188,"freelist_count":0}
{"time":"2026-11-21T14:23:27.492122134Z","level":"INFO","msg":"shutdown complete"}
holzcloud.service: Deactivated successfully.
Stopped holzcloud.service - holzcloud-CMS.
holzcloud.service: Consumed 14.088s CPU time.
Started holzcloud.service - holzcloud-CMS.
{"time":"2026-11-21T14:23:28.328734635Z","level":"INFO","msg":"starting","version":"v0.1.4-262-gd8a027e","commit":"d8a027e"}
{"time":"2026-11-21T14:23:28.335276718Z","level":"INFO","msg":"configuration loaded","config":{"port":"8080","data_dir":"/opt/holzcloud/data","db_path":"/opt/holzcloud/data/holzcloud.sqlite","log_level":"INFO","secure":true,"max_template_size":10485760,"max_media_size":5242880,"max_video_size":67108864,"max_megapixels":24,"argon2_memory_kb":65536,"argon2_iterations":1,"argon2_parallelism":2,"min_free_bytes":536870912,"smtp_host":"","smtp_from":"","smtp_configured":false,"payrexx_instance":"","payrexx_configured":false,"trusted_proxies":"127.0.0.1/32,::1/128","listen":"127.0.0.1","sso_enabled":false,"sso_configured":false,"sso_provision":false,"sso_admin_group":"","sso_default_website":0,"sso_website_groups":0,"sso_sign_out_path":"/outpost.goauthentik.io/sign_out","oidc_enabled":false,"oidc_issuer":"","oidc_client_id":"","oidc_redirect_url":"","oidc_keys":0,"oidc_client_secret_set":false,"smtp":"aus","secret_key_set":false}}
{"time":"2026-11-21T14:23:28.379422843Z","level":"INFO","msg":"sqlite opened","journal_mode":"wal","path":"/opt/holzcloud/data/holzcloud.sqlite"}
{"time":"2026-11-21T14:23:28.424425093Z","level":"INFO","msg":"sqlite opened","journal_mode":"wal","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite"}
{"time":"2026-11-21T14:23:28.437394552Z","level":"INFO","msg":"backup verified","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite","tables":62}
{"time":"2026-11-21T14:23:28.439671593Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite","level":69}
{"time":"2026-11-21T14:23:28.487366218Z","level":"INFO","msg":"migration applied","version":70,"type":"sql","duration_ms":4}
{"time":"2026-11-21T14:23:28.488909052Z","level":"ERROR","msg":"migrations failed","err":"goose up: partial migration error (type:sql,version:71): SQL logic error: no such table: t34_no_such_table (1)"}
{"time":"2026-11-21T14:23:28.489365427Z","level":"ERROR","msg":"restore with: systemctl stop holzcloud && cp /opt/holzcloud/data/pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite /opt/holzcloud/data/holzcloud.sqlite && systemctl start holzcloud"}
holzcloud.service: Main process exited, code=exited, status=1/FAILURE
holzcloud.service: Failed with result 'exit-code'.
holzcloud.service: Scheduled restart job, restart counter is at 1.
Stopped holzcloud.service - holzcloud-CMS.
Started holzcloud.service - holzcloud-CMS.
{"time":"2026-11-21T14:23:33.666108429Z","level":"INFO","msg":"starting","version":"v0.1.4-262-gd8a027e","commit":"d8a027e"}
{"time":"2026-11-21T14:23:33.668789721Z","level":"INFO","msg":"configuration loaded","config":{"port":"8080","data_dir":"/opt/holzcloud/data","db_path":"/opt/holzcloud/data/holzcloud.sqlite","log_level":"INFO","secure":true,"max_template_size":10485760,"max_media_size":5242880,"max_video_size":67108864,"max_megapixels":24,"argon2_memory_kb":65536,"argon2_iterations":1,"argon2_parallelism":2,"min_free_bytes":536870912,"smtp_host":"","smtp_from":"","smtp_configured":false,"payrexx_instance":"","payrexx_configured":false,"trusted_proxies":"127.0.0.1/32,::1/128","listen":"127.0.0.1","sso_enabled":false,"sso_configured":false,"sso_provision":false,"sso_admin_group":"","sso_default_website":0,"sso_website_groups":0,"sso_sign_out_path":"/outpost.goauthentik.io/sign_out","oidc_enabled":false,"oidc_issuer":"","oidc_client_id":"","oidc_redirect_url":"","oidc_keys":0,"oidc_client_secret_set":false,"smtp":"aus","secret_key_set":false}}
{"time":"2026-11-21T14:23:33.712891637Z","level":"INFO","msg":"sqlite opened","journal_mode":"wal","path":"/opt/holzcloud/data/holzcloud.sqlite"}
{"time":"2026-11-21T14:23:33.726664137Z","level":"WARN","msg":"pre-upgrade snapshot kept","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite","level":69,"db_level":70,"reason":"an earlier start of this update applied part of its migrations and did not finish; the snapshot from before it stays the way back, so no snapshot of the half-migrated database is written"}
{"time":"2026-11-21T14:23:33.789560804Z","level":"ERROR","msg":"migrations failed","err":"goose up: partial migration error (type:sql,version:71): SQL logic error: no such table: t34_no_such_table (1)"}
{"time":"2026-11-21T14:23:33.790071179Z","level":"ERROR","msg":"restore with: systemctl stop holzcloud && cp /opt/holzcloud/data/pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite /opt/holzcloud/data/holzcloud.sqlite && systemctl start holzcloud"}
holzcloud.service: Main process exited, code=exited, status=1/FAILURE
holzcloud.service: Failed with result 'exit-code'.
holzcloud.service: Scheduled restart job, restart counter is at 2.
Stopped holzcloud.service - holzcloud-CMS.
Started holzcloud.service - holzcloud-CMS.
{"time":"2026-11-21T14:23:39.168879417Z","level":"INFO","msg":"starting","version":"v0.1.4-262-gd8a027e","commit":"d8a027e"}
{"time":"2026-11-21T14:23:39.171635876Z","level":"INFO","msg":"configuration loaded","config":{"port":"8080","data_dir":"/opt/holzcloud/data","db_path":"/opt/holzcloud/data/holzcloud.sqlite","log_level":"INFO","secure":true,"max_template_size":10485760,"max_media_size":5242880,"max_video_size":67108864,"max_megapixels":24,"argon2_memory_kb":65536,"argon2_iterations":1,"argon2_parallelism":2,"min_free_bytes":536870912,"smtp_host":"","smtp_from":"","smtp_configured":false,"payrexx_instance":"","payrexx_configured":false,"trusted_proxies":"127.0.0.1/32,::1/128","listen":"127.0.0.1","sso_enabled":false,"sso_configured":false,"sso_provision":false,"sso_admin_group":"","sso_default_website":0,"sso_website_groups":0,"sso_sign_out_path":"/outpost.goauthentik.io/sign_out","oidc_enabled":false,"oidc_issuer":"","oidc_client_id":"","oidc_redirect_url":"","oidc_keys":0,"oidc_client_secret_set":false,"smtp":"aus","secret_key_set":false}}
{"time":"2026-11-21T14:23:39.215604334Z","level":"INFO","msg":"sqlite opened","journal_mode":"wal","path":"/opt/holzcloud/data/holzcloud.sqlite"}
{"time":"2026-11-21T14:23:39.229536292Z","level":"WARN","msg":"pre-upgrade snapshot kept","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite","level":69,"db_level":70,"reason":"an earlier start of this update applied part of its migrations and did not finish; the snapshot from before it stays the way back, so no snapshot of the half-migrated database is written"}
{"time":"2026-11-21T14:23:39.283416626Z","level":"ERROR","msg":"migrations failed","err":"goose up: partial migration error (type:sql,version:71): SQL logic error: no such table: t34_no_such_table (1)"}
{"time":"2026-11-21T14:23:39.283987626Z","level":"ERROR","msg":"restore with: systemctl stop holzcloud && cp /opt/holzcloud/data/pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite /opt/holzcloud/data/holzcloud.sqlite && systemctl start holzcloud"}
holzcloud.service: Main process exited, code=exited, status=1/FAILURE
holzcloud.service: Failed with result 'exit-code'.
holzcloud.service: Scheduled restart job, restart counter is at 3.
Stopped holzcloud.service - holzcloud-CMS.
Started holzcloud.service - holzcloud-CMS.
{"time":"2026-11-21T14:23:44.658448587Z","level":"INFO","msg":"starting","version":"v0.1.4-262-gd8a027e","commit":"d8a027e"}
{"time":"2026-11-21T14:23:44.660989462Z","level":"INFO","msg":"configuration loaded","config":{"port":"8080","data_dir":"/opt/holzcloud/data","db_path":"/opt/holzcloud/data/holzcloud.sqlite","log_level":"INFO","secure":true,"max_template_size":10485760,"max_media_size":5242880,"max_video_size":67108864,"max_megapixels":24,"argon2_memory_kb":65536,"argon2_iterations":1,"argon2_parallelism":2,"min_free_bytes":536870912,"smtp_host":"","smtp_from":"","smtp_configured":false,"payrexx_instance":"","payrexx_configured":false,"trusted_proxies":"127.0.0.1/32,::1/128","listen":"127.0.0.1","sso_enabled":false,"sso_configured":false,"sso_provision":false,"sso_admin_group":"","sso_default_website":0,"sso_website_groups":0,"sso_sign_out_path":"/outpost.goauthentik.io/sign_out","oidc_enabled":false,"oidc_issuer":"","oidc_client_id":"","oidc_redirect_url":"","oidc_keys":0,"oidc_client_secret_set":false,"smtp":"aus","secret_key_set":false}}
{"time":"2026-11-21T14:23:44.705235462Z","level":"INFO","msg":"sqlite opened","journal_mode":"wal","path":"/opt/holzcloud/data/holzcloud.sqlite"}
{"time":"2026-11-21T14:23:44.71930917Z","level":"WARN","msg":"pre-upgrade snapshot kept","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite","level":69,"db_level":70,"reason":"an earlier start of this update applied part of its migrations and did not finish; the snapshot from before it stays the way back, so no snapshot of the half-migrated database is written"}
{"time":"2026-11-21T14:23:44.78457292Z","level":"ERROR","msg":"migrations failed","err":"goose up: partial migration error (type:sql,version:71): SQL logic error: no such table: t34_no_such_table (1)"}
{"time":"2026-11-21T14:23:44.78506842Z","level":"ERROR","msg":"restore with: systemctl stop holzcloud && cp /opt/holzcloud/data/pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite /opt/holzcloud/data/holzcloud.sqlite && systemctl start holzcloud"}
holzcloud.service: Main process exited, code=exited, status=1/FAILURE
holzcloud.service: Failed with result 'exit-code'.
holzcloud.service: Scheduled restart job, restart counter is at 4.
Stopped holzcloud.service - holzcloud-CMS.
Started holzcloud.service - holzcloud-CMS.
{"time":"2026-11-21T14:23:50.157396673Z","level":"INFO","msg":"starting","version":"v0.1.4-262-gd8a027e","commit":"d8a027e"}
{"time":"2026-11-21T14:23:50.159976381Z","level":"INFO","msg":"configuration loaded","config":{"port":"8080","data_dir":"/opt/holzcloud/data","db_path":"/opt/holzcloud/data/holzcloud.sqlite","log_level":"INFO","secure":true,"max_template_size":10485760,"max_media_size":5242880,"max_video_size":67108864,"max_megapixels":24,"argon2_memory_kb":65536,"argon2_iterations":1,"argon2_parallelism":2,"min_free_bytes":536870912,"smtp_host":"","smtp_from":"","smtp_configured":false,"payrexx_instance":"","payrexx_configured":false,"trusted_proxies":"127.0.0.1/32,::1/128","listen":"127.0.0.1","sso_enabled":false,"sso_configured":false,"sso_provision":false,"sso_admin_group":"","sso_default_website":0,"sso_website_groups":0,"sso_sign_out_path":"/outpost.goauthentik.io/sign_out","oidc_enabled":false,"oidc_issuer":"","oidc_client_id":"","oidc_redirect_url":"","oidc_keys":0,"oidc_client_secret_set":false,"smtp":"aus","secret_key_set":false}}
{"time":"2026-11-21T14:23:50.203393256Z","level":"INFO","msg":"sqlite opened","journal_mode":"wal","path":"/opt/holzcloud/data/holzcloud.sqlite"}
{"time":"2026-11-21T14:23:50.217708714Z","level":"WARN","msg":"pre-upgrade snapshot kept","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite","level":69,"db_level":70,"reason":"an earlier start of this update applied part of its migrations and did not finish; the snapshot from before it stays the way back, so no snapshot of the half-migrated database is written"}
{"time":"2026-11-21T14:23:50.306649298Z","level":"ERROR","msg":"migrations failed","err":"goose up: partial migration error (type:sql,version:71): SQL logic error: no such table: t34_no_such_table (1)"}
{"time":"2026-11-21T14:23:50.307208673Z","level":"ERROR","msg":"restore with: systemctl stop holzcloud && cp /opt/holzcloud/data/pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite /opt/holzcloud/data/holzcloud.sqlite && systemctl start holzcloud"}
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
/opt/holzcloud/data/pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141823.sqlite level 0
/opt/holzcloud/data/pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite level 69
```

#### amd64-t34-rb/rb-guard.txt

```text
Started holzcloud.service - holzcloud-CMS.
{"time":"2026-11-21T14:24:05.207931888Z","level":"INFO","msg":"starting","version":"v0.1.4-261-gcefe4ac","commit":"cefe4ac"}
{"time":"2026-11-21T14:24:05.210586888Z","level":"INFO","msg":"configuration loaded","config":{"port":"8080","data_dir":"/opt/holzcloud/data","db_path":"/opt/holzcloud/data/holzcloud.sqlite","log_level":"INFO","secure":true,"max_template_size":10485760,"max_media_size":5242880,"max_video_size":67108864,"max_megapixels":24,"argon2_memory_kb":65536,"argon2_iterations":1,"argon2_parallelism":2,"min_free_bytes":536870912,"smtp_host":"","smtp_from":"","smtp_configured":false,"payrexx_instance":"","payrexx_configured":false,"trusted_proxies":"127.0.0.1/32,::1/128","listen":"127.0.0.1","sso_enabled":false,"sso_configured":false,"sso_provision":false,"sso_admin_group":"","sso_default_website":0,"sso_website_groups":0,"sso_sign_out_path":"/outpost.goauthentik.io/sign_out","oidc_enabled":false,"oidc_issuer":"","oidc_client_id":"","oidc_redirect_url":"","oidc_keys":0,"oidc_client_secret_set":false,"smtp":"aus","secret_key_set":false}}
{"time":"2026-11-21T14:24:05.256774513Z","level":"INFO","msg":"sqlite opened","journal_mode":"wal","path":"/opt/holzcloud/data/holzcloud.sqlite"}
{"time":"2026-11-21T14:24:05.267885888Z","level":"ERROR","msg":"refusing to start: the database schema is newer than this binary","db_level":70,"binary_level":69,"err":"the database is at migration level 70, and this binary knows migrations up to 69: a newer holzcloud has migrated it, and this binary refuses to run against a schema it does not know. The way back: start the newer binary again, or restore the pre-upgrade snapshot that the newer binary wrote before it migrated, together with the older binary, as deploy/DEPLOY.md describes under \"Rolling back\""}
holzcloud.service: Main process exited, code=exited, status=1/FAILURE
holzcloud.service: Failed with result 'exit-code'.
holzcloud.service: Scheduled restart job, restart counter is at 1.
Stopped holzcloud.service - holzcloud-CMS.
Started holzcloud.service - holzcloud-CMS.
{"time":"2026-11-21T14:24:10.437203585Z","level":"INFO","msg":"starting","version":"v0.1.4-261-gcefe4ac","commit":"cefe4ac"}
{"time":"2026-11-21T14:24:10.440025043Z","level":"INFO","msg":"configuration loaded","config":{"port":"8080","data_dir":"/opt/holzcloud/data","db_path":"/opt/holzcloud/data/holzcloud.sqlite","log_level":"INFO","secure":true,"max_template_size":10485760,"max_media_size":5242880,"max_video_size":67108864,"max_megapixels":24,"argon2_memory_kb":65536,"argon2_iterations":1,"argon2_parallelism":2,"min_free_bytes":536870912,"smtp_host":"","smtp_from":"","smtp_configured":false,"payrexx_instance":"","payrexx_configured":false,"trusted_proxies":"127.0.0.1/32,::1/128","listen":"127.0.0.1","sso_enabled":false,"sso_configured":false,"sso_provision":false,"sso_admin_group":"","sso_default_website":0,"sso_website_groups":0,"sso_sign_out_path":"/outpost.goauthentik.io/sign_out","oidc_enabled":false,"oidc_issuer":"","oidc_client_id":"","oidc_redirect_url":"","oidc_keys":0,"oidc_client_secret_set":false,"smtp":"aus","secret_key_set":false}}
{"time":"2026-11-21T14:24:10.485861585Z","level":"INFO","msg":"sqlite opened","journal_mode":"wal","path":"/opt/holzcloud/data/holzcloud.sqlite"}
{"time":"2026-11-21T14:24:10.497312543Z","level":"ERROR","msg":"refusing to start: the database schema is newer than this binary","db_level":70,"binary_level":69,"err":"the database is at migration level 70, and this binary knows migrations up to 69: a newer holzcloud has migrated it, and this binary refuses to run against a schema it does not know. The way back: start the newer binary again, or restore the pre-upgrade snapshot that the newer binary wrote before it migrated, together with the older binary, as deploy/DEPLOY.md describes under \"Rolling back\""}
holzcloud.service: Main process exited, code=exited, status=1/FAILURE
holzcloud.service: Failed with result 'exit-code'.
holzcloud.service: Scheduled restart job, restart counter is at 2.
Stopped holzcloud.service - holzcloud-CMS.
Started holzcloud.service - holzcloud-CMS.
{"time":"2026-11-21T14:24:15.657467462Z","level":"INFO","msg":"starting","version":"v0.1.4-261-gcefe4ac","commit":"cefe4ac"}
{"time":"2026-11-21T14:24:15.660118962Z","level":"INFO","msg":"configuration loaded","config":{"port":"8080","data_dir":"/opt/holzcloud/data","db_path":"/opt/holzcloud/data/holzcloud.sqlite","log_level":"INFO","secure":true,"max_template_size":10485760,"max_media_size":5242880,"max_video_size":67108864,"max_megapixels":24,"argon2_memory_kb":65536,"argon2_iterations":1,"argon2_parallelism":2,"min_free_bytes":536870912,"smtp_host":"","smtp_from":"","smtp_configured":false,"payrexx_instance":"","payrexx_configured":false,"trusted_proxies":"127.0.0.1/32,::1/128","listen":"127.0.0.1","sso_enabled":false,"sso_configured":false,"sso_provision":false,"sso_admin_group":"","sso_default_website":0,"sso_website_groups":0,"sso_sign_out_path":"/outpost.goauthentik.io/sign_out","oidc_enabled":false,"oidc_issuer":"","oidc_client_id":"","oidc_redirect_url":"","oidc_keys":0,"oidc_client_secret_set":false,"smtp":"aus","secret_key_set":false}}
{"time":"2026-11-21T14:24:15.704774504Z","level":"INFO","msg":"sqlite opened","journal_mode":"wal","path":"/opt/holzcloud/data/holzcloud.sqlite"}
{"time":"2026-11-21T14:24:15.715531587Z","level":"ERROR","msg":"refusing to start: the database schema is newer than this binary","db_level":70,"binary_level":69,"err":"the database is at migration level 70, and this binary knows migrations up to 69: a newer holzcloud has migrated it, and this binary refuses to run against a schema it does not know. The way back: start the newer binary again, or restore the pre-upgrade snapshot that the newer binary wrote before it migrated, together with the older binary, as deploy/DEPLOY.md describes under \"Rolling back\""}
holzcloud.service: Main process exited, code=exited, status=1/FAILURE
holzcloud.service: Failed with result 'exit-code'.
holzcloud.service: Scheduled restart job, restart counter is at 3.
Stopped holzcloud.service - holzcloud-CMS.
Started holzcloud.service - holzcloud-CMS.
{"time":"2026-11-21T14:24:20.914200631Z","level":"INFO","msg":"starting","version":"v0.1.4-261-gcefe4ac","commit":"cefe4ac"}
{"time":"2026-11-21T14:24:20.916877298Z","level":"INFO","msg":"configuration loaded","config":{"port":"8080","data_dir":"/opt/holzcloud/data","db_path":"/opt/holzcloud/data/holzcloud.sqlite","log_level":"INFO","secure":true,"max_template_size":10485760,"max_media_size":5242880,"max_video_size":67108864,"max_megapixels":24,"argon2_memory_kb":65536,"argon2_iterations":1,"argon2_parallelism":2,"min_free_bytes":536870912,"smtp_host":"","smtp_from":"","smtp_configured":false,"payrexx_instance":"","payrexx_configured":false,"trusted_proxies":"127.0.0.1/32,::1/128","listen":"127.0.0.1","sso_enabled":false,"sso_configured":false,"sso_provision":false,"sso_admin_group":"","sso_default_website":0,"sso_website_groups":0,"sso_sign_out_path":"/outpost.goauthentik.io/sign_out","oidc_enabled":false,"oidc_issuer":"","oidc_client_id":"","oidc_redirect_url":"","oidc_keys":0,"oidc_client_secret_set":false,"smtp":"aus","secret_key_set":false}}
{"time":"2026-11-21T14:24:20.965085798Z","level":"INFO","msg":"sqlite opened","journal_mode":"wal","path":"/opt/holzcloud/data/holzcloud.sqlite"}
{"time":"2026-11-21T14:24:20.97616909Z","level":"ERROR","msg":"refusing to start: the database schema is newer than this binary","db_level":70,"binary_level":69,"err":"the database is at migration level 70, and this binary knows migrations up to 69: a newer holzcloud has migrated it, and this binary refuses to run against a schema it does not know. The way back: start the newer binary again, or restore the pre-upgrade snapshot that the newer binary wrote before it migrated, together with the older binary, as deploy/DEPLOY.md describes under \"Rolling back\""}
holzcloud.service: Main process exited, code=exited, status=1/FAILURE
holzcloud.service: Failed with result 'exit-code'.
holzcloud.service: Scheduled restart job, restart counter is at 4.
Stopped holzcloud.service - holzcloud-CMS.
Started holzcloud.service - holzcloud-CMS.
{"time":"2026-11-21T14:24:26.407722926Z","level":"INFO","msg":"starting","version":"v0.1.4-261-gcefe4ac","commit":"cefe4ac"}
{"time":"2026-11-21T14:24:26.410440717Z","level":"INFO","msg":"configuration loaded","config":{"port":"8080","data_dir":"/opt/holzcloud/data","db_path":"/opt/holzcloud/data/holzcloud.sqlite","log_level":"INFO","secure":true,"max_template_size":10485760,"max_media_size":5242880,"max_video_size":67108864,"max_megapixels":24,"argon2_memory_kb":65536,"argon2_iterations":1,"argon2_parallelism":2,"min_free_bytes":536870912,"smtp_host":"","smtp_from":"","smtp_configured":false,"payrexx_instance":"","payrexx_configured":false,"trusted_proxies":"127.0.0.1/32,::1/128","listen":"127.0.0.1","sso_enabled":false,"sso_configured":false,"sso_provision":false,"sso_admin_group":"","sso_default_website":0,"sso_website_groups":0,"sso_sign_out_path":"/outpost.goauthentik.io/sign_out","oidc_enabled":false,"oidc_issuer":"","oidc_client_id":"","oidc_redirect_url":"","oidc_keys":0,"oidc_client_secret_set":false,"smtp":"aus","secret_key_set":false}}
{"time":"2026-11-21T14:24:26.514800717Z","level":"INFO","msg":"sqlite opened","journal_mode":"wal","path":"/opt/holzcloud/data/holzcloud.sqlite"}
{"time":"2026-11-21T14:24:26.526263676Z","level":"ERROR","msg":"refusing to start: the database schema is newer than this binary","db_level":70,"binary_level":69,"err":"the database is at migration level 70, and this binary knows migrations up to 69: a newer holzcloud has migrated it, and this binary refuses to run against a schema it does not know. The way back: start the newer binary again, or restore the pre-upgrade snapshot that the newer binary wrote before it migrated, together with the older binary, as deploy/DEPLOY.md describes under \"Rolling back\""}
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
1795271078.295218292 -rw------- 1 holzcloud holzcloud  12288 Nov 21 14:18 /opt/holzcloud/data/pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141823.sqlite
1795271078.367294209 -rw------- 1 holzcloud holzcloud 765952 Nov 21 14:23 /opt/holzcloud/data/pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite
WALL 1795271077.928139042 1795271078.487994375
RC 0
$ for f in /opt/holzcloud/data/pre-upgrade-*.sqlite; do echo "$f level $(sudo -u holzcloud sqlite3 "file:$f?immutable=1" 'select coalesce(max(version_id), 0) from goose_db_version where is_applied = 1')"; done
1795271079.539383168 /opt/holzcloud/data/pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141823.sqlite level 0
1795271080.005354459 /opt/holzcloud/data/pre-upgrade-v0.1.4-262-gd8a027e-20261121-142328.sqlite level 69
WALL 1795271078.805275667 1795271080.092815043
RC 0
$ sudo systemctl stop holzcloud
WALL 1795271080.421343710 1795271081.039029335
RC 0
$ sudo cp -p /opt/holzcloud/holzcloud.prev /opt/holzcloud/holzcloud
WALL 1795271081.323341627 1795271081.835661085
RC 0
$ sudo rm -f /opt/holzcloud/data/holzcloud.sqlite-wal /opt/holzcloud/data/holzcloud.sqlite-shm
WALL 1795271082.119627377 1795271082.668473669
RC 0
$ sudo mv /opt/holzcloud/data/holzcloud.sqlite /opt/holzcloud/data/holzcloud.sqlite.after-update
WALL 1795271082.955751836 1795271083.209816794
RC 0
$ sudo -u holzcloud cp "$(ls -t /opt/holzcloud/data/pre-upgrade-*.sqlite | head -n 1)" /opt/holzcloud/data/holzcloud.sqlite
WALL 1795271083.389940753 1795271083.720947170
RC 0
$ sudo systemctl reset-failed holzcloud
WALL 1795271083.897052753 1795271084.044354503
RC 0
$ sudo systemctl start holzcloud
WALL 1795271084.143992295 1795271084.296704962
RC 0
$ curl -fsS --retry 30 --retry-connrefused --retry-delay 1 http://127.0.0.1:8080/readyz | jq '{status, version, commit, integrity}'
1795271084.506778503 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 4 ms: Couldn't connect to server
1795271085.510954629 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271086.509672129 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271087.511445046 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271088.511979255 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271089.516512047 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271090.519477423 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271091.520122757 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271092.523465174 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271093.523593549 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271094.525819050 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795271095.563998467 {
1795271095.581295050   "status": "ok",
1795271095.598087425   "version": "v0.1.4-261-gcefe4ac",
1795271095.616134800   "commit": "cefe4ac",
1795271095.633862675   "integrity": "ok"
1795271095.650558842 }
WALL 1795271084.409148795 1795271095.671961842
RC 0
$ sudo journalctl -u holzcloud --since "$(systemctl show -p ActiveEnterTimestamp --value holzcloud)" -o cat | grep -E 'pre-upgrade snapshot (written|kept)|migration applied' || echo 'no migration, no snapshot'
1795271095.984705384 no migration, no snapshot
WALL 1795271095.779127550 1795271096.007424134
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
cefe4ac
```

#### amd64-t34-pass1.status

```text
0
```
