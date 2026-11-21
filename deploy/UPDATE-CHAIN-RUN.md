# Six migrating updates on one data directory, and the pruning of the pre-upgrade snapshots: results

This document is the permanent record of UPD-01 of phase 34. It records six updates in a row on one server, each played with the *Updating* blocks of `deploy/DEPLOY.md` and each to a build with one migration more than the build before it, so that every update had a migration pending and wrote a pre-upgrade snapshot. It records that after every update the data directory held exactly the newest five of all the snapshots written there, the one of the install's first start included; that each snapshot that left the directory was removed by the update that wrote a newer one, with a `removed old pre-upgrade snapshot` line in the journal of that update; and that each snapshot carries the version of the binary that wrote it in its name and holds the migration level from before that update. What `deploy/DEPLOY.md` says under *Updating* is what was played against: "It is written only when migrations are pending", "Only the newest five are kept; older ones are removed at the next update that writes one." and "VERSION is the version of the **new** binary, the one that wrote it."

The six binaries of the updates are not releases, and this is said before anything else: they are throwaway builds of this tree, build k with k no-op migrations added, made in a scratch clone and deleted with it (`DEVIATION[synthetic-build]`). Of the released tags only 0.1.2 to 0.1.3 migrates, so an update chain of released binaries writes one snapshot and never reaches the sixth, and the pruning can only be seen with builds made for it.

The record rests on one evidence set: a single `upgrade-verify.sh run` from fresh containers, made once against the committed text of every bound file (see *How this was run*); this record uses its stage `upd`, the same set as `deploy/MIGRATION-RESTORE-RUN.md` and `deploy/ROLLBACK-RUN.md`. Its raw output is the Appendix, pasted unchanged, and `bash deploy/test/upgrade-verify.sh record-check` proves every block against its scratch file byte for byte. Every number in this document is a line of that set and cites its file in backticks as `DIRECTORY/FILE`, with `DIRECTORY` the one scratch directory `amd64-t34-upd`; the run files beside it are named by their path below `deploy/test/out`. A file that is not pasted is named as `FILE in DIRECTORY`.

The record says what was done, how many times, and what was seen. It played one chain of six updates, one after another in one stage, once; it does not say what updates weeks apart or released binaries do. What the container cannot show is a named deviation or a line of *Not tested*, never a gap in the table.

## Summary

| Requirement | Result | Platform | Date, commit |
|-------------|--------|----------|--------------|
| UPD-01 | PASS: `UPD-01 PASS` stands in `amd64-t34-summary.txt` and in `amd64-t34-upd/upd.txt`. Each of the six updates was served by its own build (`READYZ-1: v0.1.4-256-gca2154f` to `READYZ-6: v0.1.4-261-g79ce726`, each the version of `BUILD-1` to `BUILD-6`), wrote one snapshot (`WRITTEN-1` to `WRITTEN-6`) and left the database at the level its migrations give (`LEVEL-EXPECTED-B0: 69`, `LEVEL-AFTER-1: 70` to `LEVEL-AFTER-6: 75`); after every update the snapshots in the data directory were exactly the newest five of all written there (`SNAPSHOTS-AFTER-1` to `SNAPSHOTS-AFTER-6` equal to `EXPECTED-1` to `EXPECTED-6`, `MATCH-1: yes` to `MATCH-6: yes`); the fifth and the sixth update each removed the oldest snapshot (`REMOVED-5: pre-upgrade-v0.1.4-255-g00854be-20261121-131604.sqlite`, `REMOVED-6: pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite`), and every snapshot that left the directory was logged (`REMOVALS-TOTAL: 2`, `REMOVALS-LOGGED: 2`), all in `amd64-t34-upd/upd.txt`. | Debian 12 container, systemd PID 1, linux/amd64 under emulation on an arm64 host (`amd64-t34-upd/versions.txt`) | 2026-11-21, 00854be |

Evidence commit: 00854be

The binary installed before the first update, B0, is the one built from `NEW_REF` (pinned by `run` to the short sha of HEAD, `00854be`, read from `deploy/test/out/amd64-t34-evidence-commit.txt`) with the Build line of `deploy/DEPLOY.md`: `BINARY: holzcloud v0.1.4-255-g00854be (00854be)` in `amd64-t34-summary.txt`, the heading line `after the install and run.sh seed (holzcloud v0.1.4-255-g00854be (00854be)), oldest first:` in `amd64-t34-upd/upd-snapshots.txt`, and the line `holzcloud v0.1.4-255-g00854be (00854be)` that `/opt/holzcloud/holzcloud version` printed in *Updating* block 1 of the first update (`amd64-t34-upd/upd-1.out`). Its level is computed from the migration files, never typed: `LEVEL-EXPECTED-B0: 69` (`amd64-t34-upd/upd.txt`). The six builds of the updates (`amd64-t34-upd/upd-builds.txt`, the same as the lines `holzcloud-syn-1:` to `holzcloud-syn-6:` of `amd64-t34-upd/versions.txt`):

```text
build 1: version v0.1.4-256-gca2154f commit ca2154f migrations 70 sha256 e7ec75b55567c4a2fc87d5e56a88c83c86b871259e45459de3d52c5486abbbf6
build 2: version v0.1.4-257-g03507cb commit 03507cb migrations 71 sha256 dda6bf4273abc6067c54e665429fe503455a3c29c27bbb8810224bd6aa8abe8d
build 3: version v0.1.4-258-g841f633 commit 841f633 migrations 72 sha256 6018b53cd4e4a38856393ca575bfb3ecb73612d88a3f509d44796a30018f3bc3
build 4: version v0.1.4-259-g185eda3 commit 185eda3 migrations 73 sha256 67e66e9a424a61d6cff5554637cc8cc1382dd5e25a25428dd4a05b038c1b8ed7
build 5: version v0.1.4-260-g4dd752d commit 4dd752d migrations 74 sha256 331adf0a1ab75ce0ca2a78ee11d17d0d2b93c94d8223fd66262b30bea838645d
build 6: version v0.1.4-261-g79ce726 commit 79ce726 migrations 75 sha256 db655b4532fc4f2ceb7c6e3a13c0688a60f8c34363151f0bcfc49f5cb58a0047
```

Each is `00854be` with one local commit per build on top, made in one scratch clone and deleted with it, so the six commits exist nowhere now. After the sixth update the container ran build 6 (`running binary: holzcloud v0.1.4-261-g79ce726 (79ce726)` in `amd64-t34-upd/versions.txt`). The summary carries `STAGES: mig upd rb` and no `REHEARSAL` line.

## How this was run

`deploy/test/upgrade-verify.sh` is a sibling of `deploy/test/run.sh` in the shape of `deploy/test/backup-verify.sh`; it edits neither and reaches the harness through `run.sh`'s public subcommands only. `run` refuses to start unless `PLATFORM=linux/amd64` and `TARGET=systemd-base`, refuses while a bound file or the binary's source differs from HEAD or is untracked, pins `NEW_REF` to the short sha of HEAD so that the binary and the scripts are one commit, and plays the stages `mig`, `upd` and `rb`, each in its own fresh containers. Every line the stage `upd` plays from *Updating* is read out of the fenced bash blocks of `deploy/DEPLOY.md` at that commit; a line the driver cannot translate ends the stage. The verdict is computed from the `KEY: value` file the stage wrote, each key exactly once, so a trace can never pass for the set. Stage `upd` used one container, `holzcloud-t34-upd`.

The set was made once, from fresh containers, on the clean committed tree of `00854be` without `REHEARSAL`, started 2026-11-21T13:13:55Z and ended 13:25:09Z on the host clock (11 minutes 14 seconds for all three stages), exit status 0 (`deploy/test/out/amd64-t34-pass1.started`, `deploy/test/out/amd64-t34-pass1.finished` and `deploy/test/out/amd64-t34-pass1.status`). Before it the worktree held no earlier output of the driver and no `holzcloud-t34-*` container existed, and after it `bash deploy/test/upgrade-verify.sh down` removed the four containers. In the same window and at the same commit, `PLATFORM=linux/amd64 TARGET=systemd-base bash deploy/test/backup-verify.sh run` made the set of phase 32 again (13:25:27Z to 13:27:56Z, status 0; that set is recorded in `deploy/BACKUP-VERIFY-RUN.md`). No second set was made against `00854be`, and no other set is cited. An earlier set was made against `36b67cc`; the commit `00854be` changed `deploy/DEPLOY.md` and `deploy/test/upgrade-verify.sh` after it, so this record was written from this set and nothing of that set is quoted. The command was:

```text
PLATFORM=linux/amd64 TARGET=systemd-base bash deploy/test/upgrade-verify.sh run
```

The knob not set on that line had its default: `STAGES=mig upd rb`. The log of the run is `deploy/test/out/amd64-t34-pass1.log`, which is not pasted; its first line reads `upgrade-verify: platform linux/amd64, NEW_REF pinned to 00854be, v0.1.1 is 5bcae9b448db814701310c5dfab26031d68c0d92, stages: mig upd rb`.

No bound file changed after the evidence commit. The bound files for this set are `deploy/DEPLOY.md`, `deploy/backup.sh`, `deploy/restore.sh`, `deploy/holzcloud.service`, `deploy/test/run.sh`, `deploy/test/Dockerfile`, `deploy/test/upgrade-verify.sh` and the five plugin archives under `plugins/`, together with the binary built from `cmd`, `internal`, `go.mod`, `go.sum`, `changelog.go` and the `CHANGELOG.md` it embeds. This command exited 0 when the record was written:

```text
git diff --quiet 00854be HEAD -- deploy/DEPLOY.md deploy/backup.sh deploy/restore.sh deploy/holzcloud.service deploy/test/run.sh deploy/test/Dockerfile deploy/test/upgrade-verify.sh plugins/bestellung/bestellung.zip plugins/jahreszahl/jahreszahl.zip plugins/kontaktformular/kontaktformular.zip plugins/nicht-gefunden/nicht-gefunden.zip plugins/suche/suche.zip cmd internal go.mod go.sum changelog.go CHANGELOG.md
```

Before the run, `bash deploy/test/upgrade-verify.sh selftest` ended `SELFTEST OK (409 cases)`; its output is not part of the set.

The stage in order: B0 was installed by the document's Install text (`run.sh build up install postinstall`), the five plugin archives of this tree were copied in and `run.sh seed` ran; the snapshots then in the data directory were listed with their levels (`amd64-t34-upd/upd-snapshots.txt`, `SNAPSHOTS-AFTER-0` in `amd64-t34-upd/upd.txt`). The driver made the six builds in one scratch clone on the host (`amd64-t34-upd/upd-builds.txt`). Then, for k from 1 to 6: `systemctl reset-failed holzcloud` (`DEVIATION[start-limit]`); *Updating* blocks 1 to 3 played line by line with build k as the new binary, each line stamped with the container clock (`amd64-t34-upd/upd-1.out` to `amd64-t34-upd/upd-6.out`); `/readyz` read for its version (`READYZ-k`); the journal since the start of that update kept (`upd-k-journal.lst` in `amd64-t34-upd`, not pasted) and its lines `pre-upgrade snapshot written` and `removed old pre-upgrade snapshot` read (`WRITTEN-k`, `REMOVED-k`); the snapshots in the data directory listed oldest first by the time stamp in the name (`SNAPSHOTS-AFTER-k`) and compared with the newest five of everything there before the first update and written since (`EXPECTED-k`, `MATCH-k`); the snapshots that had left the directory since the update before named (`LEFT-k`); the level of the live database read through a read-only connection (`LEVEL-AFTER-k`); and the snapshot list written with the level read from each file (`amd64-t34-upd/upd-snapshots.txt`). Last the driver counted the snapshots that had left the directory and those of them a `removed old pre-upgrade snapshot` line of the same update named (`REMOVALS-TOTAL`, `REMOVALS-LOGGED`).

## Deviations named

Everything here precedes every result: it is what the container does differently from a real server and what the harness did that the document does not contain. The first column is the name written to `amd64-t34-upd/deviations.txt`; the kinds of the whole run are the `DEVIATIONS:` line of `amd64-t34-summary.txt`, of which `old-plugins`, `old-script`, `backup-name`, `browser` and `transfer` belong to the stage `mig` only.

| Name | Why | What a real server does instead |
|------|-----|---------------------------------|
| DEVIATION[emulation] | The Docker host is aarch64 and the linux/amd64 container runs under emulation (`host architecture: aarch64`, `container uname -m: x86_64` in `amd64-t34-upd/versions.txt`), so every wall time of this set, the starts of the unit among them, is an emulation figure | A native amd64 CPU and disk |
| DEVIATION[synthetic-build] | The binaries of the six updates are throwaway builds of this tree, build k with the no-op migrations `00070_t34_noop_1.sql` to `00075_t34_noop_6.sql` up to the kth, each of which creates one empty table. They were made in one scratch clone with one local commit per build, so that each build's version differs and the snapshot names tell the builds apart, and with the Build line of the document; the clone was deleted afterwards, and nothing of it is committed or pushed. The Build line of *Updating* block 2 is skipped for them (`SKIP[synthetic-build]` in `amd64-t34-upd/upd-1.out` to `amd64-t34-upd/upd-6.out`). Released binaries alone never prune: of the tags only 0.1.2 to 0.1.3 migrates, so they give no five updates that each write a snapshot | Releases, each with the migrations of its own changes |
| DEVIATION[start-limit] | `systemctl reset-failed holzcloud` ran before each of the six updates: six updates within minutes, each with a stop and a start, meet `StartLimitBurst` (five starts in 300 seconds). Every update here came up at its first start (one `starting` line in each `upd-k-journal.lst` in `amd64-t34-upd`), so the reset cleared the count of earlier starts that had succeeded and hid no restart loop. The stage `rb` of the same set does not reset before its update and waits for the start window instead (`deploy/ROLLBACK-RUN.md`), because there the restart loop running to its own end is what is recorded | Updates days apart, which never meet the limit |
| DEVIATION[scp] | Four lines: `scp holzcloud user@your-server:` played as `docker cp` of B0 at the install and of build k at update k to `/root/holzcloud`, `scp -r deploy/ user@your-server:` at the install as a tar stream of the committed `deploy` directory into `/root`, and `scp deploy/backup.sh deploy/restore.sh user@your-server:` at each update as `docker cp` of the two committed scripts to `/root/` | `scp` from the build machine into the operator's home directory |
| DEVIATION[fixture] | The five plugin archives of this tree were copied into `/opt/holzcloud/plugin-zips` before `run.sh seed`; with the website, page, AI key and test image that `run.sh seed` makes, they stand in for operator work | The operator's own content and plugins |
| DEVIATION[dropin] | `/etc/systemd/system/holzcloud.service.d/container.conf` sets `ProcSubset=all`, because the document's unchanged start crashed in the container (`unit-start-unchanged.txt` in `amd64-t34-upd`) | The unit as shipped |
| DEVIATION[mount] | The unit's `RequiresMountsFor` names `/opt/holzcloud/data`, a plain directory on the container's root filesystem, so that line was not exercised, and the snapshots were written to and removed from that plain directory | A data directory on its own mount |

## Findings

Each finding has the class of the earlier records' rule: (1) a bug of the harness, fixed in the driver; (2) a limit of the container, a named deviation or a *Not tested* line; (3) a measured property of a shipped file that matches what it and its document say, recorded and not edited; (4) a defect, a shipped file doing less than it or its document says, with a fix in a commit of its own. No finding of this set is class 1 or class 4, and no shipped file was edited after the evidence commit.

The test for class 4 was whether a result contradicts something `deploy/DEPLOY.md` or the binary says at `00854be`, or leaves the operator worse off than they say. The document says that a snapshot is written only when migrations are pending, that it is named `pre-upgrade-VERSION-TIMESTAMP.sqlite` with the version of the new binary, that it holds the old schema and the journal line `pre-upgrade snapshot written` names its `path` and `level`, and that only the newest five are kept, older ones removed at the next update that writes one. F-1 to F-4 are those sentences observed; F-5 is what the check line of *Updating* step 3 shows of it.

**F-1 (class 3): After every one of the six updates the data directory held exactly the newest five of all the snapshots written there, ordered by the time stamp in their names.** Before the first update it held one snapshot, the one of the install's first start (`SNAPSHOTS-AFTER-0: pre-upgrade-v0.1.4-255-g00854be-20261121-131604.sqlite`); updates 1 to 4 each added one and removed none (`REMOVED-1: none` to `REMOVED-4: none`), and the fourth left five; updates 5 and 6 each added one and removed the oldest, so each left five again. Every observed set equals the expected one name for name (`SNAPSHOTS-AFTER-k`, `EXPECTED-k`, `MATCH-1: yes` to `MATCH-6: yes` in `amd64-t34-upd/upd.txt`). The list after update 6, with the level read from each file (`amd64-t34-upd/upd-snapshots.txt`):

```text
after update 6 (v0.1.4-261-g79ce726), oldest first:
  pre-upgrade-v0.1.4-257-g03507cb-20261121-131656.sqlite level 70
  pre-upgrade-v0.1.4-258-g841f633-20261121-131711.sqlite level 71
  pre-upgrade-v0.1.4-259-g185eda3-20261121-131727.sqlite level 72
  pre-upgrade-v0.1.4-260-g4dd752d-20261121-131743.sqlite level 73
  pre-upgrade-v0.1.4-261-g79ce726-20261121-131800.sqlite level 74
```

**F-2 (class 3): Each snapshot that left the directory was removed by the update that wrote the newer one, at its first start, and the journal of that update names it right after the line of the snapshot it wrote.** The snapshots that left are those the removal lines name (`LEFT-5` equal to `REMOVED-5`, `LEFT-6` equal to `REMOVED-6`, `REMOVALS-TOTAL: 2`, `REMOVALS-LOGGED: 2` in `amd64-t34-upd/upd.txt`). The lines of update 5 (`upd-5-journal.lst` in `amd64-t34-upd`):

```text
{"time":"2026-11-21T13:17:43.790102764Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-260-g4dd752d-20261121-131743.sqlite","level":73}
{"time":"2026-11-21T13:17:43.793014764Z","level":"INFO","msg":"removed old pre-upgrade snapshot","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-255-g00854be-20261121-131604.sqlite"}
{"time":"2026-11-21T13:17:43.821419555Z","level":"INFO","msg":"migration applied","version":74,"type":"sql","duration_ms":4}
```

and of update 6 (`upd-6-journal.lst` in `amd64-t34-upd`):

```text
{"time":"2026-11-21T13:18:00.095105674Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-261-g79ce726-20261121-131800.sqlite","level":74}
{"time":"2026-11-21T13:18:00.098101424Z","level":"INFO","msg":"removed old pre-upgrade snapshot","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite"}
{"time":"2026-11-21T13:18:00.141569882Z","level":"INFO","msg":"migration applied","version":75,"type":"sql","duration_ms":4}
```

The removal comes after the new snapshot is written and before the first migration runs, so in the fifth and the sixth update the directory held six snapshots for a moment and never fewer than five.

**F-3 (class 3): The snapshot of the install's first start, a copy of the empty database at level 0, counts among the five and was the first to be removed.** The first start of the unit after the install, at 13:16:04 (`unit-start.txt` in `amd64-t34-upd`), found every migration pending and wrote `pre-upgrade-v0.1.4-255-g00854be-20261121-131604.sqlite` at level 0 (`amd64-t34-upd/upd-snapshots.txt`). It stayed through four updates, one of the five kept after the fourth, and the fifth update removed it (`REMOVED-5`). So on a server installed fresh, the first four updates keep that empty copy as one of their five; the newest snapshot by time is the update's own, which is the one *Rolling back* takes (`deploy/ROLLBACK-RUN.md`, F-5).

**F-4 (class 3): Every snapshot carries in its name the version of the new binary, the one that wrote it, and holds the level of the database before that update.** The first update, from B0 to build 1, wrote `pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite` at level 69, the level of B0, while its block 1 printed the previous binary as `holzcloud v0.1.4-255-g00854be (00854be)` (`amd64-t34-upd/upd-1.out`, `amd64-t34-upd/upd-snapshots.txt`). The same holds for each update: the snapshot of update k is named with the version of `BUILD-k` and holds the level the database had before it, `LEVEL-EXPECTED-B0: 69` for the first and the `LEVEL-AFTER` of the update before for the others, up to 74 (`amd64-t34-upd/upd.txt`, `amd64-t34-upd/upd-snapshots.txt`). An operator who reads the names sees the version that was being installed, not the one to go back to.

**F-5 (class 3): The check line of *Updating* step 3 showed, for every update, the line of the snapshot and the one migration, and does not show the removal.** Its `grep` matches `pre-upgrade snapshot (written|kept)` and `migration applied`, so update 6 showed (`amd64-t34-upd/upd-6.out`):

```text
1795267091.537544929 {"time":"2026-11-21T13:18:00.095105674Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-261-g79ce726-20261121-131800.sqlite","level":74}
1795267091.554630804 {"time":"2026-11-21T13:18:00.141569882Z","level":"INFO","msg":"migration applied","version":75,"type":"sql","duration_ms":4}
```

The `removed old pre-upgrade snapshot` line of the same start is in the journal (F-2) and not among these lines. The document promises the snapshot line and the migration lines there, and says that older snapshots are removed; it does not say that the check line shows the removal, so this is no defect. Which snapshot went is in the full journal of the start.

## UPD-01: six migrating updates, and the newest five snapshots kept

**What was done.** B0 was installed by the document's Install text and seeded; its first start had written one snapshot, at level 0 (`SNAPSHOTS-AFTER-0` in `amd64-t34-upd/upd.txt`, `amd64-t34-upd/upd-snapshots.txt`). Then six updates were played, each with *Updating* blocks 1 to 3 as written and each to the next build: block 1 copied the running binary aside as `holzcloud.prev`, printed its version and made the backup with the document's line; block 2 skipped the Build line (`DEVIATION[synthetic-build]`), carried the build and the two scripts over, stopped the service, copied the build into place, installed the scripts and started it; block 3 waited for `/readyz`, which showed the new build with `status` and `integrity` `ok`, and printed the journal lines of the start. Every one of the twelve lines of the three blocks ended with `RC 0` in each of the six files (`amd64-t34-upd/upd-1.out` to `amd64-t34-upd/upd-6.out`). The first update's block 3 (`amd64-t34-upd/upd-1.out`):

```text
1795267012.474889670 {
1795267012.492438795   "status": "ok",
1795267012.510234837   "version": "v0.1.4-256-gca2154f",
1795267012.534073170   "commit": "ca2154f",
1795267012.550964087   "integrity": "ok"
1795267012.568266712 }
```

```text
1795267012.888691920 {"time":"2026-11-21T13:16:41.393479221Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite","level":69}
1795267012.909595629 {"time":"2026-11-21T13:16:41.440296096Z","level":"INFO","msg":"migration applied","version":70,"type":"sql","duration_ms":4}
```

The six updates, each value a line of `amd64-t34-upd/upd.txt` and the levels of the snapshots those of `amd64-t34-upd/upd-snapshots.txt`; `SNAPSHOTS-AFTER-k` equals `EXPECTED-k` for every k, so the snapshots after are given once:

| k | `READYZ-k` | `WRITTEN-k` | `REMOVED-k` | Levels of the snapshots after, oldest first | `MATCH-k` | `LEVEL-AFTER-k` |
|---|------------|-------------|-------------|----------------------------------------------|-----------|-----------------|
| 1 | v0.1.4-256-gca2154f | pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite | none | 0 69 | yes | 70 |
| 2 | v0.1.4-257-g03507cb | pre-upgrade-v0.1.4-257-g03507cb-20261121-131656.sqlite | none | 0 69 70 | yes | 71 |
| 3 | v0.1.4-258-g841f633 | pre-upgrade-v0.1.4-258-g841f633-20261121-131711.sqlite | none | 0 69 70 71 | yes | 72 |
| 4 | v0.1.4-259-g185eda3 | pre-upgrade-v0.1.4-259-g185eda3-20261121-131727.sqlite | none | 0 69 70 71 72 | yes | 73 |
| 5 | v0.1.4-260-g4dd752d | pre-upgrade-v0.1.4-260-g4dd752d-20261121-131743.sqlite | pre-upgrade-v0.1.4-255-g00854be-20261121-131604.sqlite | 69 70 71 72 73 | yes | 74 |
| 6 | v0.1.4-261-g79ce726 | pre-upgrade-v0.1.4-261-g79ce726-20261121-131800.sqlite | pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite | 70 71 72 73 74 | yes | 75 |

The removals were counted and each was found in the journal of its update (`REMOVALS-TOTAL: 2`, `REMOVALS-LOGGED: 2`, F-2), and the verdict line reads `UPD-01 PASS` (`amd64-t34-upd/upd.txt`).

*What it shows.* Six updates in a row on one data directory, each played with the *Updating* blocks of `deploy/DEPLOY.md` and each to a build with one migration more, each wrote one pre-upgrade snapshot named with the new build's version and holding the level from before it, and after every update the data directory held exactly the newest five snapshots, the two older ones each removed by the update that wrote a newer one and named in that update's journal.

**What it does not show.** Released binaries, updates weeks apart, an update without a pending migration between migrating ones, and a migration applied by a subcommand before the service started (*Not tested*).

## Measurements

The figures the harness wrote. Every one is an emulation figure on an aarch64 host (`amd64-t34-upd/versions.txt`); none is a promise of speed.

| Figure | Value | Where |
|--------|-------|-------|
| The first update, *Updating* blocks 1 to 3, first line to last | 1795266998.111444219 to 1795267012.932875295 wall | the `WALL` lines of `amd64-t34-upd/upd-1.out` |
| The sixth update, blocks 1 to 3, first line to last | 1795267076.825733172 to 1795267091.574820429 wall | the `WALL` lines of `amd64-t34-upd/upd-6.out` |
| The readiness check of block 3, first update | 1795267001.283683012 to 1795267012.589509003 wall | the `WALL` line of the `curl` line in `amd64-t34-upd/upd-1.out` |
| The readiness check of block 3, sixth update | 1795267079.989624757 to 1795267091.263210512 wall | the `WALL` line of the `curl` line in `amd64-t34-upd/upd-6.out` |
| The snapshot, the removal and the migration of update 5 | 13:17:43.790102764Z to 13:17:43.821419555Z | `upd-5-journal.lst` in `amd64-t34-upd` |
| The snapshot, the removal and the migration of update 6 | 13:18:00.095105674Z to 13:18:00.141569882Z | `upd-6-journal.lst` in `amd64-t34-upd` |
| The backups of block 1, first and sixth update | `21M` each | the last lines of the backup in `amd64-t34-upd/upd-1.out` and `amd64-t34-upd/upd-6.out` |
| The whole run, all three stages | 13:13:55Z to 13:25:09Z | `deploy/test/out/amd64-t34-pass1.started`, `deploy/test/out/amd64-t34-pass1.finished` |

## Container vs. server

| Aspect | In the container | On a real server | Measured or assumed |
|--------|------------------|------------------|---------------------|
| Speed and timing | amd64 emulated on aarch64 | A native server | Not measured natively |
| The builds of the updates | Six throwaway builds of this tree, one no-op migration more each | Releases, each with the migrations of its changes | Measured for the pruning; the binary takes the same path whatever the migration does, assumed |
| The pace | Six updates one after another, the start limit reset before each | Updates days or weeks apart | Measured for the snapshot set; the pace is assumed not to matter, since the order is by the time stamp in the name |
| The data directory | A plain directory on the container's root filesystem | Perhaps its own mount | Named deviation |
| The build and transfer | Built in a scratch clone, carried by `docker cp` | Built on the development machine, carried by `scp` | Named deviations |
| `holzcloud.service` | With the container-only drop-in `ProcSubset=all` | As shipped | Named deviation |

## Not tested

- Released binaries were not played as the chain, neither those of 0.1.x nor any of 1.x; the six updates went to throwaway builds of this tree with no-op migrations (`DEVIATION[synthetic-build]`).
- Updates spread over weeks, with nightly backups and restarts between them, and more than one update on the same day over a longer time were not played; the six updates of this set ran one after another in one stage.
- An update whose migrations a subcommand of the new binary applied before the service started (`user list`, `check`, `backup` or `migrate status`, which `deploy/DEPLOY.md` says write the same pre-upgrade snapshot from 0.1.5) was not played; every snapshot of this set was written by the service.
- An update with no migration pending between migrating ones, which writes no snapshot and removes none, was not played; every update of this set migrated.
- A failed update inside the chain was not played here; `deploy/ROLLBACK-RUN.md` records a migration that fails part-way.
- A restore from one of the kept snapshots was not played here; `deploy/ROLLBACK-RUN.md` records the rollback by snapshot.
- A data directory on its own mount was not played (`DEVIATION[mount]`); the snapshots were written to and removed from a plain directory.
- A snapshot copied out of the data directory to keep it longer, as the document suggests, was not played.
- An update chain on a real server (SRV-02) was not played; the container ran on the Mac.
- Native amd64 timing was not measured.

## Owner steps

None. Every step of this record ran in a container on the Mac, by the owner's decision; nothing here has to be done on the real server before the record counts.

## Hand-over to docs/beta.md

The sentences of this record that the page may quote, named by place and not copied here, so that each stays in exactly one place: the title sentences of F-1, F-2, F-3 and F-4 and the sentence after *What it shows.* in *UPD-01*; the limits, the title sentence of F-5 and the bullets of *Not tested*.

The limit of the page this record answers, named by its first words: "Five or more successive updates and the pruning of the pre-upgrade snapshots were not played" (six migrating updates were played on one data directory, and the pruning to the newest five was seen at the fifth and the sixth, *UPD-01*).

Nothing in this record is a promise until the page cites it byte for byte. A change to a bound file after the evidence commit means the set is made again before the record is cited for the changed text. `docs/beta.md` does not cite this record yet.

## Appendix: raw output

Each block is the scratch file of the evidence set, pasted unchanged by `bash deploy/test/upgrade-verify.sh appendix upd`, and `bash deploy/test/upgrade-verify.sh record-check` compares every block with its scratch file byte for byte. The files are in `deploy/test/out` (gitignored); the `step-NN.txt` files, `seed.txt`, `http.txt`, `install.txt`, `paths.txt`, the `.lst` files with the journals of the updates among them, the unit-start files, the synthetic builds under `syn` and the log are not pasted.

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

#### amd64-t34-upd/versions.txt

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
running binary:        holzcloud v0.1.4-261-g79ce726 (79ce726)
holzcloud-syn-1:  v0.1.4-256-gca2154f (ca2154f) migrations 70 sha256 e7ec75b55567c4a2fc87d5e56a88c83c86b871259e45459de3d52c5486abbbf6
holzcloud-syn-2:  v0.1.4-257-g03507cb (03507cb) migrations 71 sha256 dda6bf4273abc6067c54e665429fe503455a3c29c27bbb8810224bd6aa8abe8d
holzcloud-syn-3:  v0.1.4-258-g841f633 (841f633) migrations 72 sha256 6018b53cd4e4a38856393ca575bfb3ecb73612d88a3f509d44796a30018f3bc3
holzcloud-syn-4:  v0.1.4-259-g185eda3 (185eda3) migrations 73 sha256 67e66e9a424a61d6cff5554637cc8cc1382dd5e25a25428dd4a05b038c1b8ed7
holzcloud-syn-5:  v0.1.4-260-g4dd752d (4dd752d) migrations 74 sha256 331adf0a1ab75ce0ca2a78ee11d17d0d2b93c94d8223fd66262b30bea838645d
holzcloud-syn-6:  v0.1.4-261-g79ce726 (79ce726) migrations 75 sha256 db655b4532fc4f2ceb7c6e3a13c0688a60f8c34363151f0bcfc49f5cb58a0047
NEW_REF:               00854be
evidence commit:       00854be
```

#### amd64-t34-upd/deviations.txt

```text
DEVIATION[fixture]: the five plugin archives of this tree are copied into /opt/holzcloud/plugin-zips before run.sh seed; with the website, page, AI key and test image run.sh seed makes, they stand in for operator work in the admin and are no step of the document
DEVIATION[synthetic-build]: the six binaries of the updates are throwaway builds of this tree, build k with the no-op migrations 00070_t34_noop_1.sql to 00075_t34_noop_6.sql up to the kth (each creates one empty table), made in one scratch clone with one local commit per build (so that each build's version differs) and the Build line of the document, the clone deleted afterwards; nothing of it is committed or pushed. The Build line of Updating block 2 is skipped for them. Released binaries alone never prune: of the tags only 0.1.2 to 0.1.3 migrates, so they give no five updates that each write a snapshot
DEVIATION[start-limit]: systemctl reset-failed holzcloud runs before each of the six updates: six updates within minutes, each with a stop and a start, meet StartLimitBurst (five starts in 300 seconds), which updates days apart never do
DEVIATION[scp]: 'scp holzcloud user@your-server:' was played as docker cp of the synthetic build of this update to /root/holzcloud in the container (there is no ssh and no second machine here)
DEVIATION[scp]: 'scp deploy/backup.sh deploy/restore.sh user@your-server:' was played as docker cp of the working tree's deploy/backup.sh deploy/restore.sh to /root/ in the container (there is no ssh and no second machine here)
DEVIATION[emulation]: the Docker host is aarch64; the linux/amd64 container runs under emulation, so every wall time of this set (the restarts of the unit among them) is an emulation figure, not a server's
DEVIATION[scp]: 'scp holzcloud user@your-server:' was played as docker cp of deploy/test/out/bin/holzcloud-new-00854be to /root/holzcloud in the container (there is no ssh and no second machine here)
DEVIATION[scp]: 'scp -r deploy/ user@your-server:' was played as a tar stream of the working tree's deploy directory (without deploy/test/out) into /root in the container, by docker cp -
DEVIATION[dropin]: /etc/systemd/system/holzcloud.service.d/container.conf sets ProcSubset=all (container only; the document's unchanged start crashed, see unit-start-unchanged.txt)
DEVIATION[mount]: the unit's RequiresMountsFor names /opt/holzcloud/data, a plain directory on the container's root filesystem; no separate mount exists here, so that line was not exercised
```

#### amd64-t34-upd/upd.txt

```text
LEVEL-EXPECTED-B0: 69
SNAPSHOTS-AFTER-0: pre-upgrade-v0.1.4-255-g00854be-20261121-131604.sqlite
BUILD-1: v0.1.4-256-gca2154f (ca2154f) migrations 70
READYZ-1: v0.1.4-256-gca2154f
WRITTEN-1: pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite
REMOVED-1: none
SNAPSHOTS-AFTER-1: pre-upgrade-v0.1.4-255-g00854be-20261121-131604.sqlite pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite
EXPECTED-1: pre-upgrade-v0.1.4-255-g00854be-20261121-131604.sqlite pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite
MATCH-1: yes
LEFT-1: none
LEVEL-AFTER-1: 70
BUILD-2: v0.1.4-257-g03507cb (03507cb) migrations 71
READYZ-2: v0.1.4-257-g03507cb
WRITTEN-2: pre-upgrade-v0.1.4-257-g03507cb-20261121-131656.sqlite
REMOVED-2: none
SNAPSHOTS-AFTER-2: pre-upgrade-v0.1.4-255-g00854be-20261121-131604.sqlite pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite pre-upgrade-v0.1.4-257-g03507cb-20261121-131656.sqlite
EXPECTED-2: pre-upgrade-v0.1.4-255-g00854be-20261121-131604.sqlite pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite pre-upgrade-v0.1.4-257-g03507cb-20261121-131656.sqlite
MATCH-2: yes
LEFT-2: none
LEVEL-AFTER-2: 71
BUILD-3: v0.1.4-258-g841f633 (841f633) migrations 72
READYZ-3: v0.1.4-258-g841f633
WRITTEN-3: pre-upgrade-v0.1.4-258-g841f633-20261121-131711.sqlite
REMOVED-3: none
SNAPSHOTS-AFTER-3: pre-upgrade-v0.1.4-255-g00854be-20261121-131604.sqlite pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite pre-upgrade-v0.1.4-257-g03507cb-20261121-131656.sqlite pre-upgrade-v0.1.4-258-g841f633-20261121-131711.sqlite
EXPECTED-3: pre-upgrade-v0.1.4-255-g00854be-20261121-131604.sqlite pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite pre-upgrade-v0.1.4-257-g03507cb-20261121-131656.sqlite pre-upgrade-v0.1.4-258-g841f633-20261121-131711.sqlite
MATCH-3: yes
LEFT-3: none
LEVEL-AFTER-3: 72
BUILD-4: v0.1.4-259-g185eda3 (185eda3) migrations 73
READYZ-4: v0.1.4-259-g185eda3
WRITTEN-4: pre-upgrade-v0.1.4-259-g185eda3-20261121-131727.sqlite
REMOVED-4: none
SNAPSHOTS-AFTER-4: pre-upgrade-v0.1.4-255-g00854be-20261121-131604.sqlite pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite pre-upgrade-v0.1.4-257-g03507cb-20261121-131656.sqlite pre-upgrade-v0.1.4-258-g841f633-20261121-131711.sqlite pre-upgrade-v0.1.4-259-g185eda3-20261121-131727.sqlite
EXPECTED-4: pre-upgrade-v0.1.4-255-g00854be-20261121-131604.sqlite pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite pre-upgrade-v0.1.4-257-g03507cb-20261121-131656.sqlite pre-upgrade-v0.1.4-258-g841f633-20261121-131711.sqlite pre-upgrade-v0.1.4-259-g185eda3-20261121-131727.sqlite
MATCH-4: yes
LEFT-4: none
LEVEL-AFTER-4: 73
BUILD-5: v0.1.4-260-g4dd752d (4dd752d) migrations 74
READYZ-5: v0.1.4-260-g4dd752d
WRITTEN-5: pre-upgrade-v0.1.4-260-g4dd752d-20261121-131743.sqlite
REMOVED-5: pre-upgrade-v0.1.4-255-g00854be-20261121-131604.sqlite
SNAPSHOTS-AFTER-5: pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite pre-upgrade-v0.1.4-257-g03507cb-20261121-131656.sqlite pre-upgrade-v0.1.4-258-g841f633-20261121-131711.sqlite pre-upgrade-v0.1.4-259-g185eda3-20261121-131727.sqlite pre-upgrade-v0.1.4-260-g4dd752d-20261121-131743.sqlite
EXPECTED-5: pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite pre-upgrade-v0.1.4-257-g03507cb-20261121-131656.sqlite pre-upgrade-v0.1.4-258-g841f633-20261121-131711.sqlite pre-upgrade-v0.1.4-259-g185eda3-20261121-131727.sqlite pre-upgrade-v0.1.4-260-g4dd752d-20261121-131743.sqlite
MATCH-5: yes
LEFT-5: pre-upgrade-v0.1.4-255-g00854be-20261121-131604.sqlite
LEVEL-AFTER-5: 74
BUILD-6: v0.1.4-261-g79ce726 (79ce726) migrations 75
READYZ-6: v0.1.4-261-g79ce726
WRITTEN-6: pre-upgrade-v0.1.4-261-g79ce726-20261121-131800.sqlite
REMOVED-6: pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite
SNAPSHOTS-AFTER-6: pre-upgrade-v0.1.4-257-g03507cb-20261121-131656.sqlite pre-upgrade-v0.1.4-258-g841f633-20261121-131711.sqlite pre-upgrade-v0.1.4-259-g185eda3-20261121-131727.sqlite pre-upgrade-v0.1.4-260-g4dd752d-20261121-131743.sqlite pre-upgrade-v0.1.4-261-g79ce726-20261121-131800.sqlite
EXPECTED-6: pre-upgrade-v0.1.4-257-g03507cb-20261121-131656.sqlite pre-upgrade-v0.1.4-258-g841f633-20261121-131711.sqlite pre-upgrade-v0.1.4-259-g185eda3-20261121-131727.sqlite pre-upgrade-v0.1.4-260-g4dd752d-20261121-131743.sqlite pre-upgrade-v0.1.4-261-g79ce726-20261121-131800.sqlite
MATCH-6: yes
LEFT-6: pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite
LEVEL-AFTER-6: 75
REMOVALS-TOTAL: 2
REMOVALS-LOGGED: 2
UPD-01 PASS
```

#### amd64-t34-upd/upd-builds.txt

```text
six synthetic builds of 00854be (v0.1.4-255-g00854be), made in a scratch clone and deleted with it; build k adds the no-op migrations 00070 to 00075 up to the kth:
build 1: version v0.1.4-256-gca2154f commit ca2154f migrations 70 sha256 e7ec75b55567c4a2fc87d5e56a88c83c86b871259e45459de3d52c5486abbbf6
build 2: version v0.1.4-257-g03507cb commit 03507cb migrations 71 sha256 dda6bf4273abc6067c54e665429fe503455a3c29c27bbb8810224bd6aa8abe8d
build 3: version v0.1.4-258-g841f633 commit 841f633 migrations 72 sha256 6018b53cd4e4a38856393ca575bfb3ecb73612d88a3f509d44796a30018f3bc3
build 4: version v0.1.4-259-g185eda3 commit 185eda3 migrations 73 sha256 67e66e9a424a61d6cff5554637cc8cc1382dd5e25a25428dd4a05b038c1b8ed7
build 5: version v0.1.4-260-g4dd752d commit 4dd752d migrations 74 sha256 331adf0a1ab75ce0ca2a78ee11d17d0d2b93c94d8223fd66262b30bea838645d
build 6: version v0.1.4-261-g79ce726 commit 79ce726 migrations 75 sha256 db655b4532fc4f2ceb7c6e3a13c0688a60f8c34363151f0bcfc49f5cb58a0047
```

#### amd64-t34-upd/upd-1.out

```text
$ sudo cp -p /opt/holzcloud/holzcloud /opt/holzcloud/holzcloud.prev
WALL 1795266998.111444219 1795266998.242485303
RC 0
$ /opt/holzcloud/holzcloud version
1795266998.471518428 holzcloud v0.1.4-255-g00854be (00854be)
WALL 1795266998.339437928 1795266998.497641219
RC 0
$ sudo -u holzcloud /opt/holzcloud/deploy/backup.sh /opt/holzcloud/backups
1795266998.785552969 === holzcloud-CMS backup: 20261121-131638 ===
1795266998.863007178 Backing up database...
1795266998.997914303 2026/11/21 13:16:38 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/data/holzcloud.sqlite
1795266999.058890303 2026/11/21 13:16:39 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/backups/20261121-131638/holzcloud.sqlite
1795266999.082950470 2026/11/21 13:16:39 INFO backup verified path=/opt/holzcloud/backups/20261121-131638/holzcloud.sqlite tables=60
1795266999.108209136 backup written and verified: /opt/holzcloud/backups/20261121-131638/holzcloud.sqlite
1795266999.129550720 Backing up CSRF key...
1795266999.148652511 Backing up media...
1795266999.240758678 Backing up plugins...
1795266999.377962678 Writing checksums...
1795266999.701354262 Pruning backups older than 30 days...
1795266999.756181178 === Backup complete: /opt/holzcloud/backups/20261121-131638 ===
1795266999.772292512 21M	/opt/holzcloud/backups/20261121-131638
WALL 1795266998.610352136 1795266999.794003595
RC 0
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -X main.Version=$(git describe --tags --always) -X main.Commit=$(git rev-parse --short HEAD)" ./cmd/holzcloud
SKIP[synthetic-build] the binary of this update was built by this line in a scratch clone before the update
RC 0
$ scp holzcloud user@your-server:
RC 0
$ scp deploy/backup.sh deploy/restore.sh user@your-server:
RC 0
$ sudo systemctl stop holzcloud
WALL 1795267000.311555637 1795267000.466892429
RC 0
$ sudo cp ~/holzcloud /opt/holzcloud/holzcloud
WALL 1795267000.563960262 1795267000.698437929
RC 0
$ sudo install -m 755 -o root -g root ~/backup.sh ~/restore.sh /opt/holzcloud/deploy/
WALL 1795267000.797666095 1795267000.926581512
RC 0
$ sudo systemctl start holzcloud
WALL 1795267001.020601262 1795267001.166379679
RC 0
$ curl -fsS --retry 30 --retry-connrefused --retry-delay 1 http://127.0.0.1:8080/readyz | jq '{status, version, commit, integrity}'
1795267001.396968971 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 7 ms: Couldn't connect to server
1795267002.399941013 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267003.403854597 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267004.411114055 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267005.422744708 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267006.424451625 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267007.447220876 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267008.426571085 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267009.429707335 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267010.439028377 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267011.441841878 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267012.474889670 {
1795267012.492438795   "status": "ok",
1795267012.510234837   "version": "v0.1.4-256-gca2154f",
1795267012.534073170   "commit": "ca2154f",
1795267012.550964087   "integrity": "ok"
1795267012.568266712 }
WALL 1795267001.283683012 1795267012.589509003
RC 0
$ sudo journalctl -u holzcloud --since "$(systemctl show -p ActiveEnterTimestamp --value holzcloud)" -o cat | grep -E 'pre-upgrade snapshot (written|kept)|migration applied' || echo 'no migration, no snapshot'
1795267012.888691920 {"time":"2026-11-21T13:16:41.393479221Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite","level":69}
1795267012.909595629 {"time":"2026-11-21T13:16:41.440296096Z","level":"INFO","msg":"migration applied","version":70,"type":"sql","duration_ms":4}
WALL 1795267012.693828378 1795267012.932875295
RC 0
```

#### amd64-t34-upd/upd-2.out

```text
$ sudo cp -p /opt/holzcloud/holzcloud /opt/holzcloud/holzcloud.prev
WALL 1795267013.874280921 1795267014.003076212
RC 0
$ /opt/holzcloud/holzcloud version
1795267014.217547379 holzcloud v0.1.4-256-gca2154f (ca2154f)
WALL 1795267014.095495462 1795267014.239457463
RC 0
$ sudo -u holzcloud /opt/holzcloud/deploy/backup.sh /opt/holzcloud/backups
1795267014.495382796 === holzcloud-CMS backup: 20261121-131654 ===
1795267014.569101004 Backing up database...
1795267014.694540046 2026/11/21 13:16:54 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/data/holzcloud.sqlite
1795267014.741224879 2026/11/21 13:16:54 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/backups/20261121-131654/holzcloud.sqlite
1795267014.758169838 2026/11/21 13:16:54 INFO backup verified path=/opt/holzcloud/backups/20261121-131654/holzcloud.sqlite tables=63
1795267014.776896129 backup written and verified: /opt/holzcloud/backups/20261121-131654/holzcloud.sqlite
1795267014.794802921 Backing up CSRF key...
1795267014.811364046 Backing up media...
1795267014.898237296 Backing up plugins...
1795267015.030642296 Writing checksums...
1795267015.350902005 Pruning backups older than 30 days...
1795267015.410396380 === Backup complete: /opt/holzcloud/backups/20261121-131654 ===
1795267015.426478588 21M	/opt/holzcloud/backups/20261121-131654
WALL 1795267014.336613046 1795267015.447573838
RC 0
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -X main.Version=$(git describe --tags --always) -X main.Commit=$(git rev-parse --short HEAD)" ./cmd/holzcloud
SKIP[synthetic-build] the binary of this update was built by this line in a scratch clone before the update
RC 0
$ scp holzcloud user@your-server:
RC 0
$ scp deploy/backup.sh deploy/restore.sh user@your-server:
RC 0
$ sudo systemctl stop holzcloud
WALL 1795267015.920102047 1795267016.066687588
RC 0
$ sudo cp ~/holzcloud /opt/holzcloud/holzcloud
WALL 1795267016.158448588 1795267016.284864380
RC 0
$ sudo install -m 755 -o root -g root ~/backup.sh ~/restore.sh /opt/holzcloud/deploy/
WALL 1795267016.385267630 1795267016.513136839
RC 0
$ sudo systemctl start holzcloud
WALL 1795267016.606085130 1795267016.749949339
RC 0
$ curl -fsS --retry 30 --retry-connrefused --retry-delay 1 http://127.0.0.1:8080/readyz | jq '{status, version, commit, integrity}'
1795267016.953597130 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 4 ms: Couldn't connect to server
1795267017.953816756 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267018.964012173 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267019.969206465 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267020.971424882 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267021.974281841 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267022.977040258 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267023.978134134 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267024.980512051 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267025.989055968 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267027.009917969 {
1795267027.026463427   "status": "ok",
1795267027.042574885   "version": "v0.1.4-257-g03507cb",
1795267027.058707219   "commit": "03507cb",
1795267027.074962719   "integrity": "ok"
1795267027.091689344 }
WALL 1795267016.854322505 1795267027.112468469
RC 0
$ sudo journalctl -u holzcloud --since "$(systemctl show -p ActiveEnterTimestamp --value holzcloud)" -o cat | grep -E 'pre-upgrade snapshot (written|kept)|migration applied' || echo 'no migration, no snapshot'
1795267027.395369469 {"time":"2026-11-21T13:16:56.947001297Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-257-g03507cb-20261121-131656.sqlite","level":70}
1795267027.412323885 {"time":"2026-11-21T13:16:57.000728381Z","level":"INFO","msg":"migration applied","version":71,"type":"sql","duration_ms":4}
WALL 1795267027.208076552 1795267027.432846760
RC 0
```

#### amd64-t34-upd/upd-3.out

```text
$ sudo cp -p /opt/holzcloud/holzcloud /opt/holzcloud/holzcloud.prev
WALL 1795267028.510892053 1795267028.639233219
RC 0
$ /opt/holzcloud/holzcloud version
1795267028.858715053 holzcloud v0.1.4-257-g03507cb (03507cb)
WALL 1795267028.736324719 1795267028.880894553
RC 0
$ sudo -u holzcloud /opt/holzcloud/deploy/backup.sh /opt/holzcloud/backups
1795267029.133504803 === holzcloud-CMS backup: 20261121-131709 ===
1795267029.211715761 Backing up database...
1795267029.337508261 2026/11/21 13:17:09 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/data/holzcloud.sqlite
1795267029.385269053 2026/11/21 13:17:09 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/backups/20261121-131709/holzcloud.sqlite
1795267029.405579095 2026/11/21 13:17:09 INFO backup verified path=/opt/holzcloud/backups/20261121-131709/holzcloud.sqlite tables=64
1795267029.426200511 backup written and verified: /opt/holzcloud/backups/20261121-131709/holzcloud.sqlite
1795267029.446437678 Backing up CSRF key...
1795267029.462450511 Backing up media...
1795267029.539633053 Backing up plugins...
1795267029.677692678 Writing checksums...
1795267029.992978095 Pruning backups older than 30 days...
1795267030.047862887 === Backup complete: /opt/holzcloud/backups/20261121-131709 ===
1795267030.064163262 21M	/opt/holzcloud/backups/20261121-131709
WALL 1795267028.973653095 1795267030.085112970
RC 0
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -X main.Version=$(git describe --tags --always) -X main.Commit=$(git rev-parse --short HEAD)" ./cmd/holzcloud
SKIP[synthetic-build] the binary of this update was built by this line in a scratch clone before the update
RC 0
$ scp holzcloud user@your-server:
RC 0
$ scp deploy/backup.sh deploy/restore.sh user@your-server:
RC 0
$ sudo systemctl stop holzcloud
WALL 1795267030.551785595 1795267030.701363345
RC 0
$ sudo cp ~/holzcloud /opt/holzcloud/holzcloud
WALL 1795267030.792647220 1795267030.920195137
RC 0
$ sudo install -m 755 -o root -g root ~/backup.sh ~/restore.sh /opt/holzcloud/deploy/
WALL 1795267031.014367929 1795267031.140041387
RC 0
$ sudo systemctl start holzcloud
WALL 1795267031.233964179 1795267031.377990762
RC 0
$ curl -fsS --retry 30 --retry-connrefused --retry-delay 1 http://127.0.0.1:8080/readyz | jq '{status, version, commit, integrity}'
1795267031.580469179 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 4 ms: Couldn't connect to server
1795267032.581973263 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267033.585832263 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267034.588706875 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267035.599583125 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267036.603188209 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267037.612651876 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267038.614705168 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267039.616249460 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267040.622160169 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267041.622480670 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267042.647186545 {
1795267042.663736628   "status": "ok",
1795267042.680058753   "version": "v0.1.4-258-g841f633",
1795267042.696739712   "commit": "841f633",
1795267042.713750003   "integrity": "ok"
1795267042.730568420 }
WALL 1795267031.486090846 1795267042.752386087
RC 0
$ sudo journalctl -u holzcloud --since "$(systemctl show -p ActiveEnterTimestamp --value holzcloud)" -o cat | grep -E 'pre-upgrade snapshot (written|kept)|migration applied' || echo 'no migration, no snapshot'
1795267043.044091129 {"time":"2026-11-21T13:17:11.581335512Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-258-g841f633-20261121-131711.sqlite","level":71}
1795267043.062109670 {"time":"2026-11-21T13:17:11.643552512Z","level":"INFO","msg":"migration applied","version":72,"type":"sql","duration_ms":5}
WALL 1795267042.853029379 1795267043.084210795
RC 0
```

#### amd64-t34-upd/upd-4.out

```text
$ sudo cp -p /opt/holzcloud/holzcloud /opt/holzcloud/holzcloud.prev
WALL 1795267044.350164754 1795267044.485011879
RC 0
$ /opt/holzcloud/holzcloud version
1795267044.715293004 holzcloud v0.1.4-258-g841f633 (841f633)
WALL 1795267044.580973421 1795267044.738050463
RC 0
$ sudo -u holzcloud /opt/holzcloud/deploy/backup.sh /opt/holzcloud/backups
1795267045.002915796 === holzcloud-CMS backup: 20261121-131724 ===
1795267045.083251630 Backing up database...
1795267045.212907171 2026/11/21 13:17:25 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/data/holzcloud.sqlite
1795267045.264118338 2026/11/21 13:17:25 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/backups/20261121-131724/holzcloud.sqlite
1795267045.281726463 2026/11/21 13:17:25 INFO backup verified path=/opt/holzcloud/backups/20261121-131724/holzcloud.sqlite tables=65
1795267045.301778921 backup written and verified: /opt/holzcloud/backups/20261121-131724/holzcloud.sqlite
1795267045.319719588 Backing up CSRF key...
1795267045.337282505 Backing up media...
1795267045.411643296 Backing up plugins...
1795267045.545390338 Writing checksums...
1795267045.865405297 Pruning backups older than 30 days...
1795267045.921270422 === Backup complete: /opt/holzcloud/backups/20261121-131724 ===
1795267045.937348088 21M	/opt/holzcloud/backups/20261121-131724
WALL 1795267044.832554921 1795267045.959114005
RC 0
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -X main.Version=$(git describe --tags --always) -X main.Commit=$(git rev-parse --short HEAD)" ./cmd/holzcloud
SKIP[synthetic-build] the binary of this update was built by this line in a scratch clone before the update
RC 0
$ scp holzcloud user@your-server:
RC 0
$ scp deploy/backup.sh deploy/restore.sh user@your-server:
RC 0
$ sudo systemctl stop holzcloud
WALL 1795267046.445100464 1795267046.594466339
RC 0
$ sudo cp ~/holzcloud /opt/holzcloud/holzcloud
WALL 1795267046.691110589 1795267046.822602547
RC 0
$ sudo install -m 755 -o root -g root ~/backup.sh ~/restore.sh /opt/holzcloud/deploy/
WALL 1795267046.962173297 1795267047.110855131
RC 0
$ sudo systemctl start holzcloud
WALL 1795267047.317699797 1795267047.490174797
RC 0
$ curl -fsS --retry 30 --retry-connrefused --retry-delay 1 http://127.0.0.1:8080/readyz | jq '{status, version, commit, integrity}'
1795267047.704221048 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 5 ms: Couldn't connect to server
1795267048.712842590 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267049.712979340 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267050.714981632 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267051.718506133 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267052.720742092 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267053.723913800 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267054.728828551 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267055.740144718 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267056.740462593 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267057.743407302 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267058.770597428 {
1795267058.786678886   "status": "ok",
1795267058.804138678   "version": "v0.1.4-259-g185eda3",
1795267058.820393553   "commit": "185eda3",
1795267058.836750678   "integrity": "ok"
1795267058.852950261 }
WALL 1795267047.609142172 1795267058.874131344
RC 0
$ sudo journalctl -u holzcloud --since "$(systemctl show -p ActiveEnterTimestamp --value holzcloud)" -o cat | grep -E 'pre-upgrade snapshot (written|kept)|migration applied' || echo 'no migration, no snapshot'
1795267059.163950928 {"time":"2026-11-21T13:17:27.707593798Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-259-g185eda3-20261121-131727.sqlite","level":72}
1795267059.180681678 {"time":"2026-11-21T13:17:27.750552464Z","level":"INFO","msg":"migration applied","version":73,"type":"sql","duration_ms":5}
WALL 1795267058.973231261 1795267059.201925386
RC 0
```

#### amd64-t34-upd/upd-5.out

```text
$ sudo cp -p /opt/holzcloud/holzcloud /opt/holzcloud/holzcloud.prev
WALL 1795267060.582289554 1795267060.708561137
RC 0
$ /opt/holzcloud/holzcloud version
1795267060.926688679 holzcloud v0.1.4-259-g185eda3 (185eda3)
WALL 1795267060.804439054 1795267060.948586429
RC 0
$ sudo -u holzcloud /opt/holzcloud/deploy/backup.sh /opt/holzcloud/backups
1795267061.210552637 === holzcloud-CMS backup: 20261121-131741 ===
1795267061.292134721 Backing up database...
1795267061.445453846 2026/11/21 13:17:41 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/data/holzcloud.sqlite
1795267061.496067346 2026/11/21 13:17:41 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/backups/20261121-131741/holzcloud.sqlite
1795267061.516156471 2026/11/21 13:17:41 INFO backup verified path=/opt/holzcloud/backups/20261121-131741/holzcloud.sqlite tables=66
1795267061.539062221 backup written and verified: /opt/holzcloud/backups/20261121-131741/holzcloud.sqlite
1795267061.559118054 Backing up CSRF key...
1795267061.575958221 Backing up media...
1795267061.654389387 Backing up plugins...
1795267061.791134304 Writing checksums...
1795267062.123158346 Pruning backups older than 30 days...
1795267062.178502721 === Backup complete: /opt/holzcloud/backups/20261121-131741 ===
1795267062.196047221 21M	/opt/holzcloud/backups/20261121-131741
WALL 1795267061.043831887 1795267062.217720721
RC 0
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -X main.Version=$(git describe --tags --always) -X main.Commit=$(git rev-parse --short HEAD)" ./cmd/holzcloud
SKIP[synthetic-build] the binary of this update was built by this line in a scratch clone before the update
RC 0
$ scp holzcloud user@your-server:
RC 0
$ scp deploy/backup.sh deploy/restore.sh user@your-server:
RC 0
$ sudo systemctl stop holzcloud
WALL 1795267062.742500055 1795267062.890560430
RC 0
$ sudo cp ~/holzcloud /opt/holzcloud/holzcloud
WALL 1795267062.985887846 1795267063.117635430
RC 0
$ sudo install -m 755 -o root -g root ~/backup.sh ~/restore.sh /opt/holzcloud/deploy/
WALL 1795267063.216020472 1795267063.342517138
RC 0
$ sudo systemctl start holzcloud
WALL 1795267063.437833513 1795267063.580939930
RC 0
$ curl -fsS --retry 30 --retry-connrefused --retry-delay 1 http://127.0.0.1:8080/readyz | jq '{status, version, commit, integrity}'
1795267063.789383222 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 5 ms: Couldn't connect to server
1795267064.793476416 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267065.799719458 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267066.798154501 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267067.801793918 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267068.807655543 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267069.811859877 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267070.821385169 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267071.826997670 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267072.836037254 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267073.836845296 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267074.866899338 {
1795267074.884017338   "status": "ok",
1795267074.901330296   "version": "v0.1.4-260-g4dd752d",
1795267074.923740046   "commit": "4dd752d",
1795267074.945949213   "integrity": "ok"
1795267074.964447338 }
WALL 1795267063.695416888 1795267074.986335713
RC 0
$ sudo journalctl -u holzcloud --since "$(systemctl show -p ActiveEnterTimestamp --value holzcloud)" -o cat | grep -E 'pre-upgrade snapshot (written|kept)|migration applied' || echo 'no migration, no snapshot'
1795267075.305593213 {"time":"2026-11-21T13:17:43.790102764Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-260-g4dd752d-20261121-131743.sqlite","level":73}
1795267075.322813046 {"time":"2026-11-21T13:17:43.821419555Z","level":"INFO","msg":"migration applied","version":74,"type":"sql","duration_ms":4}
WALL 1795267075.102834171 1795267075.344762838
RC 0
```

#### amd64-t34-upd/upd-6.out

```text
$ sudo cp -p /opt/holzcloud/holzcloud /opt/holzcloud/holzcloud.prev
WALL 1795267076.825733172 1795267076.961509089
RC 0
$ /opt/holzcloud/holzcloud version
1795267077.189266172 holzcloud v0.1.4-260-g4dd752d (4dd752d)
WALL 1795267077.061435881 1795267077.211000339
RC 0
$ sudo -u holzcloud /opt/holzcloud/deploy/backup.sh /opt/holzcloud/backups
1795267077.477576631 === holzcloud-CMS backup: 20261121-131757 ===
1795267077.562515964 Backing up database...
1795267077.690813339 2026/11/21 13:17:57 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/data/holzcloud.sqlite
1795267077.741875631 2026/11/21 13:17:57 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/backups/20261121-131757/holzcloud.sqlite
1795267077.762874131 2026/11/21 13:17:57 INFO backup verified path=/opt/holzcloud/backups/20261121-131757/holzcloud.sqlite tables=67
1795267077.780549298 backup written and verified: /opt/holzcloud/backups/20261121-131757/holzcloud.sqlite
1795267077.801744381 Backing up CSRF key...
1795267077.819555756 Backing up media...
1795267077.914192131 Backing up plugins...
1795267078.054572798 Writing checksums...
1795267078.385468423 Pruning backups older than 30 days...
1795267078.445308131 === Backup complete: /opt/holzcloud/backups/20261121-131757 ===
1795267078.463170340 21M	/opt/holzcloud/backups/20261121-131757
WALL 1795267077.313245339 1795267078.483935340
RC 0
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -X main.Version=$(git describe --tags --always) -X main.Commit=$(git rev-parse --short HEAD)" ./cmd/holzcloud
SKIP[synthetic-build] the binary of this update was built by this line in a scratch clone before the update
RC 0
$ scp holzcloud user@your-server:
RC 0
$ scp deploy/backup.sh deploy/restore.sh user@your-server:
RC 0
$ sudo systemctl stop holzcloud
WALL 1795267078.989446923 1795267079.149103840
RC 0
$ sudo cp ~/holzcloud /opt/holzcloud/holzcloud
WALL 1795267079.250200757 1795267079.380141257
RC 0
$ sudo install -m 755 -o root -g root ~/backup.sh ~/restore.sh /opt/holzcloud/deploy/
WALL 1795267079.482706090 1795267079.621892215
RC 0
$ sudo systemctl start holzcloud
WALL 1795267079.726819090 1795267079.875798549
RC 0
$ curl -fsS --retry 30 --retry-connrefused --retry-delay 1 http://127.0.0.1:8080/readyz | jq '{status, version, commit, integrity}'
1795267080.087891424 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 4 ms: Couldn't connect to server
1795267081.085975841 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267082.087574758 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267083.096425967 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267084.100226759 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267085.105888093 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267086.111288677 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267087.113366010 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267088.121971511 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267089.125797636 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267090.135009178 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795267091.158373471 {
1795267091.174848387   "status": "ok",
1795267091.191099804   "version": "v0.1.4-261-g79ce726",
1795267091.208026054   "commit": "79ce726",
1795267091.224775762   "integrity": "ok"
1795267091.241721887 }
WALL 1795267079.989624757 1795267091.263210512
RC 0
$ sudo journalctl -u holzcloud --since "$(systemctl show -p ActiveEnterTimestamp --value holzcloud)" -o cat | grep -E 'pre-upgrade snapshot (written|kept)|migration applied' || echo 'no migration, no snapshot'
1795267091.537544929 {"time":"2026-11-21T13:18:00.095105674Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-261-g79ce726-20261121-131800.sqlite","level":74}
1795267091.554630804 {"time":"2026-11-21T13:18:00.141569882Z","level":"INFO","msg":"migration applied","version":75,"type":"sql","duration_ms":4}
WALL 1795267091.356493721 1795267091.574820429
RC 0
```

#### amd64-t34-upd/upd-snapshots.txt

```text
after the install and run.sh seed (holzcloud v0.1.4-255-g00854be (00854be)), oldest first:
  pre-upgrade-v0.1.4-255-g00854be-20261121-131604.sqlite level 0
after update 1 (v0.1.4-256-gca2154f), oldest first:
  pre-upgrade-v0.1.4-255-g00854be-20261121-131604.sqlite level 0
  pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite level 69
after update 2 (v0.1.4-257-g03507cb), oldest first:
  pre-upgrade-v0.1.4-255-g00854be-20261121-131604.sqlite level 0
  pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite level 69
  pre-upgrade-v0.1.4-257-g03507cb-20261121-131656.sqlite level 70
after update 3 (v0.1.4-258-g841f633), oldest first:
  pre-upgrade-v0.1.4-255-g00854be-20261121-131604.sqlite level 0
  pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite level 69
  pre-upgrade-v0.1.4-257-g03507cb-20261121-131656.sqlite level 70
  pre-upgrade-v0.1.4-258-g841f633-20261121-131711.sqlite level 71
after update 4 (v0.1.4-259-g185eda3), oldest first:
  pre-upgrade-v0.1.4-255-g00854be-20261121-131604.sqlite level 0
  pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite level 69
  pre-upgrade-v0.1.4-257-g03507cb-20261121-131656.sqlite level 70
  pre-upgrade-v0.1.4-258-g841f633-20261121-131711.sqlite level 71
  pre-upgrade-v0.1.4-259-g185eda3-20261121-131727.sqlite level 72
after update 5 (v0.1.4-260-g4dd752d), oldest first:
  pre-upgrade-v0.1.4-256-gca2154f-20261121-131641.sqlite level 69
  pre-upgrade-v0.1.4-257-g03507cb-20261121-131656.sqlite level 70
  pre-upgrade-v0.1.4-258-g841f633-20261121-131711.sqlite level 71
  pre-upgrade-v0.1.4-259-g185eda3-20261121-131727.sqlite level 72
  pre-upgrade-v0.1.4-260-g4dd752d-20261121-131743.sqlite level 73
after update 6 (v0.1.4-261-g79ce726), oldest first:
  pre-upgrade-v0.1.4-257-g03507cb-20261121-131656.sqlite level 70
  pre-upgrade-v0.1.4-258-g841f633-20261121-131711.sqlite level 71
  pre-upgrade-v0.1.4-259-g185eda3-20261121-131727.sqlite level 72
  pre-upgrade-v0.1.4-260-g4dd752d-20261121-131743.sqlite level 73
  pre-upgrade-v0.1.4-261-g79ce726-20261121-131800.sqlite level 74
```

#### amd64-t34-evidence-commit.txt

```text
00854be
```

#### amd64-t34-pass1.status

```text
0
```
