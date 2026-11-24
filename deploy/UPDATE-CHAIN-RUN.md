# Six migrating updates on one data directory, and the pruning of the pre-upgrade snapshots: results

This document is the permanent record of UPD-01 of phase 34. It records six updates in a row on one server, each played with the *Updating* blocks of `deploy/DEPLOY.md` and each to a build with one migration more than the build before it, so that every update had a migration pending and wrote a pre-upgrade snapshot. It records that after every update the data directory held exactly the newest five of all the snapshots written there, the one of the install's first start included; that each snapshot that left the directory was removed by the update that wrote a newer one, with a `removed old pre-upgrade snapshot` line in the journal of that update; and that each snapshot carries the version of the binary that wrote it in its name and holds the migration level from before that update. What `deploy/DEPLOY.md` says under *Updating* is what was played against: "It is written only when migrations are pending", "Only the newest five are kept; older ones are removed at the next update that writes one." and "VERSION is the version of the **new** binary, the one that wrote it."

The six binaries of the updates are not releases, and this is said before anything else: they are throwaway builds of this tree, build k with k no-op migrations added, made in a scratch clone and deleted with it (`DEVIATION[synthetic-build]`). Of the released tags only 0.1.2 to 0.1.3 migrates, so an update chain of released binaries writes one snapshot and never reaches the sixth, and the pruning can only be seen with builds made for it.

The record rests on one evidence set: a single `upgrade-verify.sh run` from fresh containers, made once against the committed text of every bound file (see *How this was run*); this record uses its stage `upd`, the same set as `deploy/MIGRATION-RESTORE-RUN.md` and `deploy/ROLLBACK-RUN.md`. Its raw output is the Appendix, pasted unchanged, and `bash deploy/test/upgrade-verify.sh record-check` proves every block against its scratch file byte for byte. Every number in this document is a line of that set and cites its file in backticks as `DIRECTORY/FILE`, with `DIRECTORY` the one scratch directory `amd64-t34-upd`; the run files beside it are named by their path below `deploy/test/out`. A file that is not pasted is named as `FILE in DIRECTORY`.

The record says what was done, how many times, and what was seen. It played one chain of six updates, one after another in one stage, once; it does not say what updates weeks apart or released binaries do. What the container cannot show is a named deviation or a line of *Not tested*, never a gap in the table.

## Summary

| Requirement | Result | Platform | Date, commit |
|-------------|--------|----------|--------------|
| UPD-01 | PASS: `UPD-01 PASS` stands in `amd64-t34-summary.txt` and in `amd64-t34-upd/upd.txt`. Each of the six updates was served by its own build (`READYZ-1: v0.1.4-262-g7283b76` to `READYZ-6: v0.1.4-267-g456235e`, each the version of `BUILD-1` to `BUILD-6`), wrote one snapshot (`WRITTEN-1` to `WRITTEN-6`) and left the database at the level its migrations give (`LEVEL-EXPECTED-B0: 69`, `LEVEL-AFTER-1: 70` to `LEVEL-AFTER-6: 75`); after every update the snapshots in the data directory were exactly the newest five of all written there (`SNAPSHOTS-AFTER-1` to `SNAPSHOTS-AFTER-6` equal to `EXPECTED-1` to `EXPECTED-6`, `MATCH-1: yes` to `MATCH-6: yes`); the fifth and the sixth update each removed the oldest snapshot (`REMOVED-5: pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141543.sqlite`, `REMOVED-6: pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite`), and every snapshot that left the directory was logged (`REMOVALS-TOTAL: 2`, `REMOVALS-LOGGED: 2`), all in `amd64-t34-upd/upd.txt`. | Debian 12 container, systemd PID 1, linux/amd64 under emulation on an arm64 host (`amd64-t34-upd/versions.txt`) | 2026-11-21, cefe4ac |

Evidence commit: cefe4ac

The binary installed before the first update, B0, is the one built from `NEW_REF` (pinned by `run` to the short sha of HEAD, `cefe4ac`, read from `deploy/test/out/amd64-t34-evidence-commit.txt`) with the Build line of `deploy/DEPLOY.md`: `BINARY: holzcloud v0.1.4-261-gcefe4ac (cefe4ac)` in `amd64-t34-summary.txt`, the heading line `after the install and run.sh seed (holzcloud v0.1.4-261-gcefe4ac (cefe4ac)), oldest first:` in `amd64-t34-upd/upd-snapshots.txt`, and the line `holzcloud v0.1.4-261-gcefe4ac (cefe4ac)` that `/opt/holzcloud/holzcloud version` printed in *Updating* block 1 of the first update (`amd64-t34-upd/upd-1.out`). Its level is computed from the migration files, never typed: `LEVEL-EXPECTED-B0: 69` (`amd64-t34-upd/upd.txt`). The six builds of the updates (`amd64-t34-upd/upd-builds.txt`, the same as the lines `holzcloud-syn-1:` to `holzcloud-syn-6:` of `amd64-t34-upd/versions.txt`):

```text
build 1: version v0.1.4-262-g7283b76 commit 7283b76 migrations 70 sha256 472900f6cdcf8d8eb0853feea43bc4344a6b3655e134f740af4cbcf36daec7e9
build 2: version v0.1.4-263-g5083f30 commit 5083f30 migrations 71 sha256 8788209711aedd74a4c084824f48e28bbd79bfc80ca653551c0cb274d3a0bac8
build 3: version v0.1.4-264-g1f4590b commit 1f4590b migrations 72 sha256 7d4b92156ecf064f4a827f90540c667cee3a22e72c1ed320e5d3b9f192f72483
build 4: version v0.1.4-265-g5f8e3f6 commit 5f8e3f6 migrations 73 sha256 18583839aed44478bc2dfeb8cabaaf9252655839c01223afd64a71285219a9ae
build 5: version v0.1.4-266-g10e0ef7 commit 10e0ef7 migrations 74 sha256 96c268a92b027850121d84a5b8ea781d46ff92a543dddb4415cabda197a2bcad
build 6: version v0.1.4-267-g456235e commit 456235e migrations 75 sha256 c4c38f8db381583ff5c5243d7d5bb3d4e1a8e7b69afbd1ea8462d84311c3322f
```

Each is `cefe4ac` with one local commit per build on top, made in one scratch clone and deleted with it, so the six commits exist nowhere now. After the sixth update the container ran build 6 (`running binary: holzcloud v0.1.4-267-g456235e (456235e)` in `amd64-t34-upd/versions.txt`). The summary carries `STAGES: mig upd rb` and no `REHEARSAL` line.

## How this was run

`deploy/test/upgrade-verify.sh` is a sibling of `deploy/test/run.sh` in the shape of `deploy/test/backup-verify.sh`; it edits neither and reaches the harness through `run.sh`'s public subcommands only. `run` refuses to start unless `PLATFORM=linux/amd64` and `TARGET=systemd-base`, refuses while a bound file or the binary's source differs from HEAD or is untracked, pins `NEW_REF` to the short sha of HEAD so that the binary and the scripts are one commit, and plays the stages `mig`, `upd` and `rb`, each in its own fresh containers. Every line the stage `upd` plays from *Updating* is read out of the fenced bash blocks of `deploy/DEPLOY.md` at that commit; a line the driver cannot translate ends the stage. The verdict is computed from the `KEY: value` file the stage wrote, each key exactly once, so a trace can never pass for the set. Stage `upd` used one container, `holzcloud-t34-upd`.

The set was made once, from fresh containers, on the clean committed tree of `cefe4ac` without `REHEARSAL`, started 2026-11-21T14:13:05Z and ended 14:24:59Z on the host clock (11 minutes 54 seconds for all three stages), exit status 0 (`deploy/test/out/amd64-t34-pass1.started`, `deploy/test/out/amd64-t34-pass1.finished` and `deploy/test/out/amd64-t34-pass1.status`). Before it the worktree held no earlier output of the driver and no `holzcloud-t34-*` container existed, and after it `bash deploy/test/upgrade-verify.sh down` removed the four containers. In the same window and at the same commit, `PLATFORM=linux/amd64 TARGET=systemd-base bash deploy/test/backup-verify.sh run` made the set of phase 32 again (14:25:11Z to 14:27:08Z, status 0; that set is recorded in `deploy/BACKUP-VERIFY-RUN.md`). No second set was made against `cefe4ac`, and no other set is cited. Two earlier sets were superseded. The first was made against `36b67cc`; the commit `00854be` changed `deploy/DEPLOY.md` and `deploy/test/upgrade-verify.sh` after it. The second was made against `00854be`, and this record was first written from it. Its four verdicts read PASS, but the browser record of phase 34 then found that a plugin switched on in the admin did not start when the browser dropped the request that switched it on, because the plugin was loaded with that request's context and closed with it. The commit `cefe4ac` fixed that in `internal/plugin/runtime.go`, so the binary changed after the set of `00854be`, and the set was made again. So this record was written from this set, and nothing of either earlier set is quoted. The command was:

```text
PLATFORM=linux/amd64 TARGET=systemd-base bash deploy/test/upgrade-verify.sh run
```

The knob not set on that line had its default: `STAGES=mig upd rb`. The log of the run is `deploy/test/out/amd64-t34-pass1.log`, which is not pasted; its first line reads `upgrade-verify: platform linux/amd64, NEW_REF pinned to cefe4ac, v0.1.1 is 5bcae9b448db814701310c5dfab26031d68c0d92, stages: mig upd rb`.

No bound file changed after the evidence commit. The bound files for this set are `deploy/DEPLOY.md`, `deploy/backup.sh`, `deploy/restore.sh`, `deploy/holzcloud.service`, `deploy/test/run.sh`, `deploy/test/Dockerfile`, `deploy/test/upgrade-verify.sh` and the five plugin archives under `plugins/`, together with the binary built from `cmd`, `internal`, `go.mod`, `go.sum`, `changelog.go` and the `CHANGELOG.md` it embeds. This command exited 0 when the record was written:

```text
git diff --quiet cefe4ac HEAD -- deploy/DEPLOY.md deploy/backup.sh deploy/restore.sh deploy/holzcloud.service deploy/test/run.sh deploy/test/Dockerfile deploy/test/upgrade-verify.sh plugins/bestellung/bestellung.zip plugins/jahreszahl/jahreszahl.zip plugins/kontaktformular/kontaktformular.zip plugins/nicht-gefunden/nicht-gefunden.zip plugins/suche/suche.zip cmd internal go.mod go.sum changelog.go CHANGELOG.md
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

Each finding has the class of the earlier records' rule: (1) a bug of the harness, fixed in the driver; (2) a limit of the container, a named deviation or a *Not tested* line; (3) a measured property of a shipped file that matches what it and its document say, recorded and not edited; (4) a defect, a shipped file doing less than it or its document says, with a fix in a commit of its own. No finding of this set is class 1 or class 4, and no shipped file was edited after the evidence commit. The defect that the browser record found at `00854be`, a plugin load ended by a dropped request, was fixed in `cefe4ac` before this set (*How this was run*).

The test for class 4 was whether a result contradicts something `deploy/DEPLOY.md` or the binary says at `cefe4ac`, or leaves the operator worse off than they say. The document says that a snapshot is written only when migrations are pending, that it is named `pre-upgrade-VERSION-TIMESTAMP.sqlite` with the version of the new binary, that it holds the old schema and the journal line `pre-upgrade snapshot written` names its `path` and `level`, and that only the newest five are kept, older ones removed at the next update that writes one. F-1 to F-4 are those sentences observed; F-5 is what the check line of *Updating* step 3 shows of it.

**F-1 (class 3): After every one of the six updates the data directory held exactly the newest five of all the snapshots written there, ordered by the time stamp in their names.** Before the first update it held one snapshot, the one of the install's first start (`SNAPSHOTS-AFTER-0: pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141543.sqlite`); updates 1 to 4 each added one and removed none (`REMOVED-1: none` to `REMOVED-4: none`), and the fourth left five; updates 5 and 6 each added one and removed the oldest, so each left five again. Every observed set equals the expected one name for name (`SNAPSHOTS-AFTER-k`, `EXPECTED-k`, `MATCH-1: yes` to `MATCH-6: yes` in `amd64-t34-upd/upd.txt`). The list after update 6, with the level read from each file (`amd64-t34-upd/upd-snapshots.txt`):

```text
after update 6 (v0.1.4-267-g456235e), oldest first:
  pre-upgrade-v0.1.4-263-g5083f30-20261121-141644.sqlite level 70
  pre-upgrade-v0.1.4-264-g1f4590b-20261121-141659.sqlite level 71
  pre-upgrade-v0.1.4-265-g5f8e3f6-20261121-141715.sqlite level 72
  pre-upgrade-v0.1.4-266-g10e0ef7-20261121-141731.sqlite level 73
  pre-upgrade-v0.1.4-267-g456235e-20261121-141747.sqlite level 74
```

**F-2 (class 3): Each snapshot that left the directory was removed by the update that wrote the newer one, at its first start, and the journal of that update names it right after the line of the snapshot it wrote.** The snapshots that left are those the removal lines name (`LEFT-5` equal to `REMOVED-5`, `LEFT-6` equal to `REMOVED-6`, `REMOVALS-TOTAL: 2`, `REMOVALS-LOGGED: 2` in `amd64-t34-upd/upd.txt`). The lines of update 5 (`upd-5-journal.lst` in `amd64-t34-upd`):

```text
{"time":"2026-11-21T14:17:31.935809095Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-266-g10e0ef7-20261121-141731.sqlite","level":73}
{"time":"2026-11-21T14:17:31.938603803Z","level":"INFO","msg":"removed old pre-upgrade snapshot","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141543.sqlite"}
{"time":"2026-11-21T14:17:31.976282595Z","level":"INFO","msg":"migration applied","version":74,"type":"sql","duration_ms":4}
```

and of update 6 (`upd-6-journal.lst` in `amd64-t34-upd`):

```text
{"time":"2026-11-21T14:17:47.915038213Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-267-g456235e-20261121-141747.sqlite","level":74}
{"time":"2026-11-21T14:17:47.917850088Z","level":"INFO","msg":"removed old pre-upgrade snapshot","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite"}
{"time":"2026-11-21T14:17:47.967620505Z","level":"INFO","msg":"migration applied","version":75,"type":"sql","duration_ms":4}
```

The removal comes after the new snapshot is written and before the first migration runs, so in the fifth and the sixth update the directory held six snapshots for a moment and never fewer than five.

**F-3 (class 3): The snapshot of the install's first start, a copy of the empty database at level 0, counts among the five and was the first to be removed.** The first start of the unit after the install, at 14:15:43 (`unit-start.txt` in `amd64-t34-upd`), found every migration pending and wrote `pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141543.sqlite` at level 0 (`amd64-t34-upd/upd-snapshots.txt`). It stayed through four updates, one of the five kept after the fourth, and the fifth update removed it (`REMOVED-5`). So on a server installed fresh, the first four updates keep that empty copy as one of their five; the newest snapshot by time is the update's own, which is the one *Rolling back* takes (`deploy/ROLLBACK-RUN.md`, F-5).

**F-4 (class 3): Every snapshot carries in its name the version of the new binary, the one that wrote it, and holds the level of the database before that update.** The first update, from B0 to build 1, wrote `pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite` at level 69, the level of B0, while its block 1 printed the previous binary as `holzcloud v0.1.4-261-gcefe4ac (cefe4ac)` (`amd64-t34-upd/upd-1.out`, `amd64-t34-upd/upd-snapshots.txt`). The same holds for each update: the snapshot of update k is named with the version of `BUILD-k` and holds the level the database had before it, `LEVEL-EXPECTED-B0: 69` for the first and the `LEVEL-AFTER` of the update before for the others, up to 74 (`amd64-t34-upd/upd.txt`, `amd64-t34-upd/upd-snapshots.txt`). An operator who reads the names sees the version that was being installed, not the one to go back to.

**F-5 (class 3): The check line of *Updating* step 3 showed, for every update, the line of the snapshot and the one migration, and does not show the removal.** Its `grep` matches `pre-upgrade snapshot (written|kept)` and `migration applied`, so update 6 showed (`amd64-t34-upd/upd-6.out`):

```text
1795270679.387092594 {"time":"2026-11-21T14:17:47.915038213Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-267-g456235e-20261121-141747.sqlite","level":74}
1795270679.404573927 {"time":"2026-11-21T14:17:47.967620505Z","level":"INFO","msg":"migration applied","version":75,"type":"sql","duration_ms":4}
```

The `removed old pre-upgrade snapshot` line of the same start is in the journal (F-2) and not among these lines. The document promises the snapshot line and the migration lines there, and says that older snapshots are removed; it does not say that the check line shows the removal, so this is no defect. Which snapshot went is in the full journal of the start.

## UPD-01: six migrating updates, and the newest five snapshots kept

**What was done.** B0 was installed by the document's Install text and seeded; its first start had written one snapshot, at level 0 (`SNAPSHOTS-AFTER-0` in `amd64-t34-upd/upd.txt`, `amd64-t34-upd/upd-snapshots.txt`). Then six updates were played, each with *Updating* blocks 1 to 3 as written and each to the next build: block 1 copied the running binary aside as `holzcloud.prev`, printed its version and made the backup with the document's line; block 2 skipped the Build line (`DEVIATION[synthetic-build]`), carried the build and the two scripts over, stopped the service, copied the build into place, installed the scripts and started it; block 3 waited for `/readyz`, which showed the new build with `status` and `integrity` `ok`, and printed the journal lines of the start. Every one of the twelve lines of the three blocks ended with `RC 0` in each of the six files (`amd64-t34-upd/upd-1.out` to `amd64-t34-upd/upd-6.out`). The first update's block 3 (`amd64-t34-upd/upd-1.out`):

```text
1795270599.559283626 {
1795270599.576125834   "status": "ok",
1795270599.593030959   "version": "v0.1.4-262-g7283b76",
1795270599.609502876   "commit": "7283b76",
1795270599.626427501   "integrity": "ok"
1795270599.642735709 }
```

```text
1795270599.949804251 {"time":"2026-11-21T14:16:28.500797343Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite","level":69}
1795270599.967086626 {"time":"2026-11-21T14:16:28.518327343Z","level":"INFO","msg":"migration applied","version":70,"type":"sql","duration_ms":4}
```

The six updates, each value a line of `amd64-t34-upd/upd.txt` and the levels of the snapshots those of `amd64-t34-upd/upd-snapshots.txt`; `SNAPSHOTS-AFTER-k` equals `EXPECTED-k` for every k, so the snapshots after are given once:

| k | `READYZ-k` | `WRITTEN-k` | `REMOVED-k` | Levels of the snapshots after, oldest first | `MATCH-k` | `LEVEL-AFTER-k` |
|---|------------|-------------|-------------|----------------------------------------------|-----------|-----------------|
| 1 | v0.1.4-262-g7283b76 | pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite | none | 0 69 | yes | 70 |
| 2 | v0.1.4-263-g5083f30 | pre-upgrade-v0.1.4-263-g5083f30-20261121-141644.sqlite | none | 0 69 70 | yes | 71 |
| 3 | v0.1.4-264-g1f4590b | pre-upgrade-v0.1.4-264-g1f4590b-20261121-141659.sqlite | none | 0 69 70 71 | yes | 72 |
| 4 | v0.1.4-265-g5f8e3f6 | pre-upgrade-v0.1.4-265-g5f8e3f6-20261121-141715.sqlite | none | 0 69 70 71 72 | yes | 73 |
| 5 | v0.1.4-266-g10e0ef7 | pre-upgrade-v0.1.4-266-g10e0ef7-20261121-141731.sqlite | pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141543.sqlite | 69 70 71 72 73 | yes | 74 |
| 6 | v0.1.4-267-g456235e | pre-upgrade-v0.1.4-267-g456235e-20261121-141747.sqlite | pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite | 70 71 72 73 74 | yes | 75 |

The removals were counted and each was found in the journal of its update (`REMOVALS-TOTAL: 2`, `REMOVALS-LOGGED: 2`, F-2), and the verdict line reads `UPD-01 PASS` (`amd64-t34-upd/upd.txt`).

*What it shows.* Six updates in a row on one data directory, each played with the *Updating* blocks of `deploy/DEPLOY.md` and each to a build with one migration more, each wrote one pre-upgrade snapshot named with the new build's version and holding the level from before it, and after every update the data directory held exactly the newest five snapshots, the two older ones each removed by the update that wrote a newer one and named in that update's journal.

**What it does not show.** Released binaries, updates weeks apart, an update without a pending migration between migrating ones, and a migration applied by a subcommand before the service started (*Not tested*).

## Measurements

The figures the harness wrote. Every one is an emulation figure on an aarch64 host (`amd64-t34-upd/versions.txt`); none is a promise of speed.

| Figure | Value | Where |
|--------|-------|-------|
| The first update, *Updating* blocks 1 to 3, first line to last | 1795270585.103843758 to 1795270599.988490210 wall | the `WALL` lines of `amd64-t34-upd/upd-1.out` |
| The sixth update, blocks 1 to 3, first line to last | 1795270664.807393670 to 1795270679.425482469 wall | the `WALL` lines of `amd64-t34-upd/upd-6.out` |
| The readiness check of block 3, first update | 1795270588.392953010 to 1795270599.664453793 wall | the `WALL` line of the `curl` line in `amd64-t34-upd/upd-1.out` |
| The readiness check of block 3, sixth update | 1795270667.819584588 to 1795270679.106112302 wall | the `WALL` line of the `curl` line in `amd64-t34-upd/upd-6.out` |
| The snapshot, the removal and the migration of update 5 | 14:17:31.935809095Z to 14:17:31.976282595Z | `upd-5-journal.lst` in `amd64-t34-upd` |
| The snapshot, the removal and the migration of update 6 | 14:17:47.915038213Z to 14:17:47.967620505Z | `upd-6-journal.lst` in `amd64-t34-upd` |
| The backups of block 1, first and sixth update | `21M` each | the last lines of the backup in `amd64-t34-upd/upd-1.out` and `amd64-t34-upd/upd-6.out` |
| The whole run, all three stages | 14:13:05Z to 14:24:59Z | `deploy/test/out/amd64-t34-pass1.started`, `deploy/test/out/amd64-t34-pass1.finished` |

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
BINARY: holzcloud v0.1.4-261-gcefe4ac (cefe4ac)
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
running binary:        holzcloud v0.1.4-267-g456235e (456235e)
holzcloud-syn-1:  v0.1.4-262-g7283b76 (7283b76) migrations 70 sha256 472900f6cdcf8d8eb0853feea43bc4344a6b3655e134f740af4cbcf36daec7e9
holzcloud-syn-2:  v0.1.4-263-g5083f30 (5083f30) migrations 71 sha256 8788209711aedd74a4c084824f48e28bbd79bfc80ca653551c0cb274d3a0bac8
holzcloud-syn-3:  v0.1.4-264-g1f4590b (1f4590b) migrations 72 sha256 7d4b92156ecf064f4a827f90540c667cee3a22e72c1ed320e5d3b9f192f72483
holzcloud-syn-4:  v0.1.4-265-g5f8e3f6 (5f8e3f6) migrations 73 sha256 18583839aed44478bc2dfeb8cabaaf9252655839c01223afd64a71285219a9ae
holzcloud-syn-5:  v0.1.4-266-g10e0ef7 (10e0ef7) migrations 74 sha256 96c268a92b027850121d84a5b8ea781d46ff92a543dddb4415cabda197a2bcad
holzcloud-syn-6:  v0.1.4-267-g456235e (456235e) migrations 75 sha256 c4c38f8db381583ff5c5243d7d5bb3d4e1a8e7b69afbd1ea8462d84311c3322f
NEW_REF:               cefe4ac
evidence commit:       cefe4ac
```

#### amd64-t34-upd/deviations.txt

```text
DEVIATION[fixture]: the five plugin archives of this tree are copied into /opt/holzcloud/plugin-zips before run.sh seed; with the website, page, AI key and test image run.sh seed makes, they stand in for operator work in the admin and are no step of the document
DEVIATION[synthetic-build]: the six binaries of the updates are throwaway builds of this tree, build k with the no-op migrations 00070_t34_noop_1.sql to 00075_t34_noop_6.sql up to the kth (each creates one empty table), made in one scratch clone with one local commit per build (so that each build's version differs) and the Build line of the document, the clone deleted afterwards; nothing of it is committed or pushed. The Build line of Updating block 2 is skipped for them. Released binaries alone never prune: of the tags only 0.1.2 to 0.1.3 migrates, so they give no five updates that each write a snapshot
DEVIATION[start-limit]: systemctl reset-failed holzcloud runs before each of the six updates: six updates within minutes, each with a stop and a start, meet StartLimitBurst (five starts in 300 seconds), which updates days apart never do
DEVIATION[scp]: 'scp holzcloud user@your-server:' was played as docker cp of the synthetic build of this update to /root/holzcloud in the container (there is no ssh and no second machine here)
DEVIATION[scp]: 'scp deploy/backup.sh deploy/restore.sh user@your-server:' was played as docker cp of the working tree's deploy/backup.sh deploy/restore.sh to /root/ in the container (there is no ssh and no second machine here)
DEVIATION[emulation]: the Docker host is aarch64; the linux/amd64 container runs under emulation, so every wall time of this set (the restarts of the unit among them) is an emulation figure, not a server's
DEVIATION[scp]: 'scp holzcloud user@your-server:' was played as docker cp of deploy/test/out/bin/holzcloud-new-cefe4ac to /root/holzcloud in the container (there is no ssh and no second machine here)
DEVIATION[scp]: 'scp -r deploy/ user@your-server:' was played as a tar stream of the working tree's deploy directory (without deploy/test/out) into /root in the container, by docker cp -
DEVIATION[dropin]: /etc/systemd/system/holzcloud.service.d/container.conf sets ProcSubset=all (container only; the document's unchanged start crashed, see unit-start-unchanged.txt)
DEVIATION[mount]: the unit's RequiresMountsFor names /opt/holzcloud/data, a plain directory on the container's root filesystem; no separate mount exists here, so that line was not exercised
```

#### amd64-t34-upd/upd.txt

```text
LEVEL-EXPECTED-B0: 69
SNAPSHOTS-AFTER-0: pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141543.sqlite
BUILD-1: v0.1.4-262-g7283b76 (7283b76) migrations 70
READYZ-1: v0.1.4-262-g7283b76
WRITTEN-1: pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite
REMOVED-1: none
SNAPSHOTS-AFTER-1: pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141543.sqlite pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite
EXPECTED-1: pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141543.sqlite pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite
MATCH-1: yes
LEFT-1: none
LEVEL-AFTER-1: 70
BUILD-2: v0.1.4-263-g5083f30 (5083f30) migrations 71
READYZ-2: v0.1.4-263-g5083f30
WRITTEN-2: pre-upgrade-v0.1.4-263-g5083f30-20261121-141644.sqlite
REMOVED-2: none
SNAPSHOTS-AFTER-2: pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141543.sqlite pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite pre-upgrade-v0.1.4-263-g5083f30-20261121-141644.sqlite
EXPECTED-2: pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141543.sqlite pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite pre-upgrade-v0.1.4-263-g5083f30-20261121-141644.sqlite
MATCH-2: yes
LEFT-2: none
LEVEL-AFTER-2: 71
BUILD-3: v0.1.4-264-g1f4590b (1f4590b) migrations 72
READYZ-3: v0.1.4-264-g1f4590b
WRITTEN-3: pre-upgrade-v0.1.4-264-g1f4590b-20261121-141659.sqlite
REMOVED-3: none
SNAPSHOTS-AFTER-3: pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141543.sqlite pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite pre-upgrade-v0.1.4-263-g5083f30-20261121-141644.sqlite pre-upgrade-v0.1.4-264-g1f4590b-20261121-141659.sqlite
EXPECTED-3: pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141543.sqlite pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite pre-upgrade-v0.1.4-263-g5083f30-20261121-141644.sqlite pre-upgrade-v0.1.4-264-g1f4590b-20261121-141659.sqlite
MATCH-3: yes
LEFT-3: none
LEVEL-AFTER-3: 72
BUILD-4: v0.1.4-265-g5f8e3f6 (5f8e3f6) migrations 73
READYZ-4: v0.1.4-265-g5f8e3f6
WRITTEN-4: pre-upgrade-v0.1.4-265-g5f8e3f6-20261121-141715.sqlite
REMOVED-4: none
SNAPSHOTS-AFTER-4: pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141543.sqlite pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite pre-upgrade-v0.1.4-263-g5083f30-20261121-141644.sqlite pre-upgrade-v0.1.4-264-g1f4590b-20261121-141659.sqlite pre-upgrade-v0.1.4-265-g5f8e3f6-20261121-141715.sqlite
EXPECTED-4: pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141543.sqlite pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite pre-upgrade-v0.1.4-263-g5083f30-20261121-141644.sqlite pre-upgrade-v0.1.4-264-g1f4590b-20261121-141659.sqlite pre-upgrade-v0.1.4-265-g5f8e3f6-20261121-141715.sqlite
MATCH-4: yes
LEFT-4: none
LEVEL-AFTER-4: 73
BUILD-5: v0.1.4-266-g10e0ef7 (10e0ef7) migrations 74
READYZ-5: v0.1.4-266-g10e0ef7
WRITTEN-5: pre-upgrade-v0.1.4-266-g10e0ef7-20261121-141731.sqlite
REMOVED-5: pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141543.sqlite
SNAPSHOTS-AFTER-5: pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite pre-upgrade-v0.1.4-263-g5083f30-20261121-141644.sqlite pre-upgrade-v0.1.4-264-g1f4590b-20261121-141659.sqlite pre-upgrade-v0.1.4-265-g5f8e3f6-20261121-141715.sqlite pre-upgrade-v0.1.4-266-g10e0ef7-20261121-141731.sqlite
EXPECTED-5: pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite pre-upgrade-v0.1.4-263-g5083f30-20261121-141644.sqlite pre-upgrade-v0.1.4-264-g1f4590b-20261121-141659.sqlite pre-upgrade-v0.1.4-265-g5f8e3f6-20261121-141715.sqlite pre-upgrade-v0.1.4-266-g10e0ef7-20261121-141731.sqlite
MATCH-5: yes
LEFT-5: pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141543.sqlite
LEVEL-AFTER-5: 74
BUILD-6: v0.1.4-267-g456235e (456235e) migrations 75
READYZ-6: v0.1.4-267-g456235e
WRITTEN-6: pre-upgrade-v0.1.4-267-g456235e-20261121-141747.sqlite
REMOVED-6: pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite
SNAPSHOTS-AFTER-6: pre-upgrade-v0.1.4-263-g5083f30-20261121-141644.sqlite pre-upgrade-v0.1.4-264-g1f4590b-20261121-141659.sqlite pre-upgrade-v0.1.4-265-g5f8e3f6-20261121-141715.sqlite pre-upgrade-v0.1.4-266-g10e0ef7-20261121-141731.sqlite pre-upgrade-v0.1.4-267-g456235e-20261121-141747.sqlite
EXPECTED-6: pre-upgrade-v0.1.4-263-g5083f30-20261121-141644.sqlite pre-upgrade-v0.1.4-264-g1f4590b-20261121-141659.sqlite pre-upgrade-v0.1.4-265-g5f8e3f6-20261121-141715.sqlite pre-upgrade-v0.1.4-266-g10e0ef7-20261121-141731.sqlite pre-upgrade-v0.1.4-267-g456235e-20261121-141747.sqlite
MATCH-6: yes
LEFT-6: pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite
LEVEL-AFTER-6: 75
REMOVALS-TOTAL: 2
REMOVALS-LOGGED: 2
UPD-01 PASS
```

#### amd64-t34-upd/upd-builds.txt

```text
six synthetic builds of cefe4ac (v0.1.4-261-gcefe4ac), made in a scratch clone and deleted with it; build k adds the no-op migrations 00070 to 00075 up to the kth:
build 1: version v0.1.4-262-g7283b76 commit 7283b76 migrations 70 sha256 472900f6cdcf8d8eb0853feea43bc4344a6b3655e134f740af4cbcf36daec7e9
build 2: version v0.1.4-263-g5083f30 commit 5083f30 migrations 71 sha256 8788209711aedd74a4c084824f48e28bbd79bfc80ca653551c0cb274d3a0bac8
build 3: version v0.1.4-264-g1f4590b commit 1f4590b migrations 72 sha256 7d4b92156ecf064f4a827f90540c667cee3a22e72c1ed320e5d3b9f192f72483
build 4: version v0.1.4-265-g5f8e3f6 commit 5f8e3f6 migrations 73 sha256 18583839aed44478bc2dfeb8cabaaf9252655839c01223afd64a71285219a9ae
build 5: version v0.1.4-266-g10e0ef7 commit 10e0ef7 migrations 74 sha256 96c268a92b027850121d84a5b8ea781d46ff92a543dddb4415cabda197a2bcad
build 6: version v0.1.4-267-g456235e commit 456235e migrations 75 sha256 c4c38f8db381583ff5c5243d7d5bb3d4e1a8e7b69afbd1ea8462d84311c3322f
```

#### amd64-t34-upd/upd-1.out

```text
$ sudo cp -p /opt/holzcloud/holzcloud /opt/holzcloud/holzcloud.prev
WALL 1795270585.103843758 1795270585.233585550
RC 0
$ /opt/holzcloud/holzcloud version
1795270585.458321425 holzcloud v0.1.4-261-gcefe4ac (cefe4ac)
WALL 1795270585.333095259 1795270585.480051425
RC 0
$ sudo -u holzcloud /opt/holzcloud/deploy/backup.sh /opt/holzcloud/backups
1795270585.738803050 === holzcloud-CMS backup: 20261121-141625 ===
1795270585.817088550 Backing up database...
1795270585.945932675 2026/11/21 14:16:25 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/data/holzcloud.sqlite
1795270586.002512967 2026/11/21 14:16:25 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/backups/20261121-141625/holzcloud.sqlite
1795270586.022391092 2026/11/21 14:16:26 INFO backup verified path=/opt/holzcloud/backups/20261121-141625/holzcloud.sqlite tables=60
1795270586.043188842 backup written and verified: /opt/holzcloud/backups/20261121-141625/holzcloud.sqlite
1795270586.063194134 Backing up CSRF key...
1795270586.080799717 Backing up media...
1795270586.165732884 Backing up plugins...
1795270586.293998051 Writing checksums...
1795270586.621497967 Pruning backups older than 30 days...
1795270586.681791801 === Backup complete: /opt/holzcloud/backups/20261121-141625 ===
1795270586.698679884 21M	/opt/holzcloud/backups/20261121-141625
WALL 1795270585.575884842 1795270586.720642509
RC 0
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -X main.Version=$(git describe --tags --always) -X main.Commit=$(git rev-parse --short HEAD)" ./cmd/holzcloud
SKIP[synthetic-build] the binary of this update was built by this line in a scratch clone before the update
RC 0
$ scp holzcloud user@your-server:
RC 0
$ scp deploy/backup.sh deploy/restore.sh user@your-server:
RC 0
$ sudo systemctl stop holzcloud
WALL 1795270587.220412801 1795270587.580838260
RC 0
$ sudo cp ~/holzcloud /opt/holzcloud/holzcloud
WALL 1795270587.677130343 1795270587.808481260
RC 0
$ sudo install -m 755 -o root -g root ~/backup.sh ~/restore.sh /opt/holzcloud/deploy/
WALL 1795270587.906874801 1795270588.036244135
RC 0
$ sudo systemctl start holzcloud
WALL 1795270588.133476260 1795270588.282441260
RC 0
$ curl -fsS --retry 30 --retry-connrefused --retry-delay 1 http://127.0.0.1:8080/readyz | jq '{status, version, commit, integrity}'
1795270588.493516427 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 4 ms: Couldn't connect to server
1795270589.490756260 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270590.496693928 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270591.496083053 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270592.499422095 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270593.503350429 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270594.511653555 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270595.530036472 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270596.524182333 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270597.526791292 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270598.529488417 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270599.559283626 {
1795270599.576125834   "status": "ok",
1795270599.593030959   "version": "v0.1.4-262-g7283b76",
1795270599.609502876   "commit": "7283b76",
1795270599.626427501   "integrity": "ok"
1795270599.642735709 }
WALL 1795270588.392953010 1795270599.664453793
RC 0
$ sudo journalctl -u holzcloud --since "$(systemctl show -p ActiveEnterTimestamp --value holzcloud)" -o cat | grep -E 'pre-upgrade snapshot (written|kept)|migration applied' || echo 'no migration, no snapshot'
1795270599.949804251 {"time":"2026-11-21T14:16:28.500797343Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite","level":69}
1795270599.967086626 {"time":"2026-11-21T14:16:28.518327343Z","level":"INFO","msg":"migration applied","version":70,"type":"sql","duration_ms":4}
WALL 1795270599.761534501 1795270599.988490210
RC 0
```

#### amd64-t34-upd/upd-2.out

```text
$ sudo cp -p /opt/holzcloud/holzcloud /opt/holzcloud/holzcloud.prev
WALL 1795270600.967011210 1795270601.100035335
RC 0
$ /opt/holzcloud/holzcloud version
1795270601.324468585 holzcloud v0.1.4-262-g7283b76 (7283b76)
WALL 1795270601.199137502 1795270601.347024002
RC 0
$ sudo -u holzcloud /opt/holzcloud/deploy/backup.sh /opt/holzcloud/backups
1795270601.603831377 === holzcloud-CMS backup: 20261121-141641 ===
1795270601.685162919 Backing up database...
1795270601.816657752 2026/11/21 14:16:41 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/data/holzcloud.sqlite
1795270601.880198335 2026/11/21 14:16:41 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/backups/20261121-141641/holzcloud.sqlite
1795270601.899651252 2026/11/21 14:16:41 INFO backup verified path=/opt/holzcloud/backups/20261121-141641/holzcloud.sqlite tables=63
1795270601.923059794 backup written and verified: /opt/holzcloud/backups/20261121-141641/holzcloud.sqlite
1795270601.944822752 Backing up CSRF key...
1795270601.964253419 Backing up media...
1795270602.040912544 Backing up plugins...
1795270602.173053336 Writing checksums...
1795270602.496409461 Pruning backups older than 30 days...
1795270602.553861961 === Backup complete: /opt/holzcloud/backups/20261121-141641 ===
1795270602.572282711 21M	/opt/holzcloud/backups/20261121-141641
WALL 1795270601.442878460 1795270602.593573044
RC 0
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -X main.Version=$(git describe --tags --always) -X main.Commit=$(git rev-parse --short HEAD)" ./cmd/holzcloud
SKIP[synthetic-build] the binary of this update was built by this line in a scratch clone before the update
RC 0
$ scp holzcloud user@your-server:
RC 0
$ scp deploy/backup.sh deploy/restore.sh user@your-server:
RC 0
$ sudo systemctl stop holzcloud
WALL 1795270603.084367128 1795270603.236911586
RC 0
$ sudo cp ~/holzcloud /opt/holzcloud/holzcloud
WALL 1795270603.331998211 1795270603.462062920
RC 0
$ sudo install -m 755 -o root -g root ~/backup.sh ~/restore.sh /opt/holzcloud/deploy/
WALL 1795270603.561621003 1795270603.690584795
RC 0
$ sudo systemctl start holzcloud
WALL 1795270603.786587753 1795270603.936544586
RC 0
$ curl -fsS --retry 30 --retry-connrefused --retry-delay 1 http://127.0.0.1:8080/readyz | jq '{status, version, commit, integrity}'
1795270604.156347503 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 5 ms: Couldn't connect to server
1795270605.162484795 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270606.173051837 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270607.184860296 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270608.189729088 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270609.195174089 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270610.206390173 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270611.208155590 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270612.219916132 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270613.231166924 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270614.242756341 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270615.269231175 {
1795270615.287139050   "status": "ok",
1795270615.303823008   "version": "v0.1.4-263-g5083f30",
1795270615.320177633   "commit": "5083f30",
1795270615.336778134   "integrity": "ok"
1795270615.352628800 }
WALL 1795270604.053734711 1795270615.373736550
RC 0
$ sudo journalctl -u holzcloud --since "$(systemctl show -p ActiveEnterTimestamp --value holzcloud)" -o cat | grep -E 'pre-upgrade snapshot (written|kept)|migration applied' || echo 'no migration, no snapshot'
1795270615.662776050 {"time":"2026-11-21T14:16:44.157677087Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-263-g5083f30-20261121-141644.sqlite","level":70}
1795270615.679977009 {"time":"2026-11-21T14:16:44.195696128Z","level":"INFO","msg":"migration applied","version":71,"type":"sql","duration_ms":5}
WALL 1795270615.472589634 1795270615.701287509
RC 0
```

#### amd64-t34-upd/upd-3.out

```text
$ sudo cp -p /opt/holzcloud/holzcloud /opt/holzcloud/holzcloud.prev
WALL 1795270616.773221634 1795270616.901716426
RC 0
$ /opt/holzcloud/holzcloud version
1795270617.122847176 holzcloud v0.1.4-263-g5083f30 (5083f30)
WALL 1795270616.997354009 1795270617.144907134
RC 0
$ sudo -u holzcloud /opt/holzcloud/deploy/backup.sh /opt/holzcloud/backups
1795270617.405021926 === holzcloud-CMS backup: 20261121-141657 ===
1795270617.483843385 Backing up database...
1795270617.611242051 2026/11/21 14:16:57 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/data/holzcloud.sqlite
1795270617.674007635 2026/11/21 14:16:57 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/backups/20261121-141657/holzcloud.sqlite
1795270617.694345426 2026/11/21 14:16:57 INFO backup verified path=/opt/holzcloud/backups/20261121-141657/holzcloud.sqlite tables=64
1795270617.715621968 backup written and verified: /opt/holzcloud/backups/20261121-141657/holzcloud.sqlite
1795270617.733231676 Backing up CSRF key...
1795270617.750951801 Backing up media...
1795270617.831061760 Backing up plugins...
1795270617.960720801 Writing checksums...
1795270618.285644093 Pruning backups older than 30 days...
1795270618.340292135 === Backup complete: /opt/holzcloud/backups/20261121-141657 ===
1795270618.358149218 21M	/opt/holzcloud/backups/20261121-141657
WALL 1795270617.239773384 1795270618.378731302
RC 0
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -X main.Version=$(git describe --tags --always) -X main.Commit=$(git rev-parse --short HEAD)" ./cmd/holzcloud
SKIP[synthetic-build] the binary of this update was built by this line in a scratch clone before the update
RC 0
$ scp holzcloud user@your-server:
RC 0
$ scp deploy/backup.sh deploy/restore.sh user@your-server:
RC 0
$ sudo systemctl stop holzcloud
WALL 1795270618.856427052 1795270619.009419552
RC 0
$ sudo cp ~/holzcloud /opt/holzcloud/holzcloud
WALL 1795270619.109415260 1795270619.238934844
RC 0
$ sudo install -m 755 -o root -g root ~/backup.sh ~/restore.sh /opt/holzcloud/deploy/
WALL 1795270619.339429844 1795270619.466430302
RC 0
$ sudo systemctl start holzcloud
WALL 1795270619.567218844 1795270619.712616011
RC 0
$ curl -fsS --retry 30 --retry-connrefused --retry-delay 1 http://127.0.0.1:8080/readyz | jq '{status, version, commit, integrity}'
1795270619.919888219 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 4 ms: Couldn't connect to server
1795270620.921323303 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270621.933059887 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270622.944535887 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270623.954120054 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270624.954763638 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270625.963405264 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270626.972894958 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270627.977826375 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270628.981205792 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270629.981249835 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270631.015266918 {
1795270631.034748085   "status": "ok",
1795270631.050577418   "version": "v0.1.4-264-g1f4590b",
1795270631.066638293   "commit": "1f4590b",
1795270631.082691335   "integrity": "ok"
1795270631.099581418 }
WALL 1795270619.824754177 1795270631.120534002
RC 0
$ sudo journalctl -u holzcloud --since "$(systemctl show -p ActiveEnterTimestamp --value holzcloud)" -o cat | grep -E 'pre-upgrade snapshot (written|kept)|migration applied' || echo 'no migration, no snapshot'
1795270631.404027044 {"time":"2026-11-21T14:16:59.926742677Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-264-g1f4590b-20261121-141659.sqlite","level":71}
1795270631.421297919 {"time":"2026-11-21T14:16:59.959061552Z","level":"INFO","msg":"migration applied","version":72,"type":"sql","duration_ms":4}
WALL 1795270631.215507168 1795270631.442304835
RC 0
```

#### amd64-t34-upd/upd-4.out

```text
$ sudo cp -p /opt/holzcloud/holzcloud /opt/holzcloud/holzcloud.prev
WALL 1795270632.713998711 1795270632.843547961
RC 0
$ /opt/holzcloud/holzcloud version
1795270633.068433044 holzcloud v0.1.4-264-g1f4590b (1f4590b)
WALL 1795270632.941347211 1795270633.090412294
RC 0
$ sudo -u holzcloud /opt/holzcloud/deploy/backup.sh /opt/holzcloud/backups
1795270633.353097628 === holzcloud-CMS backup: 20261121-141713 ===
1795270633.428416586 Backing up database...
1795270633.562989170 2026/11/21 14:17:13 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/data/holzcloud.sqlite
1795270633.612419586 2026/11/21 14:17:13 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/backups/20261121-141713/holzcloud.sqlite
1795270633.630675503 2026/11/21 14:17:13 INFO backup verified path=/opt/holzcloud/backups/20261121-141713/holzcloud.sqlite tables=65
1795270633.650484128 backup written and verified: /opt/holzcloud/backups/20261121-141713/holzcloud.sqlite
1795270633.670884711 Backing up CSRF key...
1795270633.687477795 Backing up media...
1795270633.771238003 Backing up plugins...
1795270633.902644961 Writing checksums...
1795270634.221438295 Pruning backups older than 30 days...
1795270634.281068128 === Backup complete: /opt/holzcloud/backups/20261121-141713 ===
1795270634.297562837 21M	/opt/holzcloud/backups/20261121-141713
WALL 1795270633.187932503 1795270634.318879337
RC 0
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -X main.Version=$(git describe --tags --always) -X main.Commit=$(git rev-parse --short HEAD)" ./cmd/holzcloud
SKIP[synthetic-build] the binary of this update was built by this line in a scratch clone before the update
RC 0
$ scp holzcloud user@your-server:
RC 0
$ scp deploy/backup.sh deploy/restore.sh user@your-server:
RC 0
$ sudo systemctl stop holzcloud
WALL 1795270634.799643003 1795270634.949259629
RC 0
$ sudo cp ~/holzcloud /opt/holzcloud/holzcloud
WALL 1795270635.044393087 1795270635.174006004
RC 0
$ sudo install -m 755 -o root -g root ~/backup.sh ~/restore.sh /opt/holzcloud/deploy/
WALL 1795270635.284517629 1795270635.412564837
RC 0
$ sudo systemctl start holzcloud
WALL 1795270635.506899920 1795270635.653443754
RC 0
$ curl -fsS --retry 30 --retry-connrefused --retry-delay 1 http://127.0.0.1:8080/readyz | jq '{status, version, commit, integrity}'
1795270635.867455046 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 4 ms: Couldn't connect to server
1795270636.861805546 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270637.867194338 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270638.869749214 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270639.875338756 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270640.876512465 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270641.886746799 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270642.889855841 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270643.899908549 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270644.906591217 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270645.908386592 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270646.935802009 {
1795270646.952629718   "status": "ok",
1795270646.969094843   "version": "v0.1.4-265-g5f8e3f6",
1795270646.985275301   "commit": "5f8e3f6",
1795270647.001266468   "integrity": "ok"
1795270647.017236551 }
WALL 1795270635.768846921 1795270647.038312509
RC 0
$ sudo journalctl -u holzcloud --since "$(systemctl show -p ActiveEnterTimestamp --value holzcloud)" -o cat | grep -E 'pre-upgrade snapshot (written|kept)|migration applied' || echo 'no migration, no snapshot'
1795270647.324128301 {"time":"2026-11-21T14:17:15.864979962Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-265-g5f8e3f6-20261121-141715.sqlite","level":72}
1795270647.341334676 {"time":"2026-11-21T14:17:15.919389421Z","level":"INFO","msg":"migration applied","version":73,"type":"sql","duration_ms":4}
WALL 1795270647.136600509 1795270647.362594676
RC 0
```

#### amd64-t34-upd/upd-5.out

```text
$ sudo cp -p /opt/holzcloud/holzcloud /opt/holzcloud/holzcloud.prev
WALL 1795270648.700458802 1795270648.829057344
RC 0
$ /opt/holzcloud/holzcloud version
1795270649.061171385 holzcloud v0.1.4-265-g5f8e3f6 (5f8e3f6)
WALL 1795270648.927608344 1795270649.083498552
RC 0
$ sudo -u holzcloud /opt/holzcloud/deploy/backup.sh /opt/holzcloud/backups
1795270649.453125510 === holzcloud-CMS backup: 20261121-141729 ===
1795270649.523198927 Backing up database...
1795270649.648551427 2026/11/21 14:17:29 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/data/holzcloud.sqlite
1795270649.692853427 2026/11/21 14:17:29 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/backups/20261121-141729/holzcloud.sqlite
1795270649.713357969 2026/11/21 14:17:29 INFO backup verified path=/opt/holzcloud/backups/20261121-141729/holzcloud.sqlite tables=66
1795270649.735211677 backup written and verified: /opt/holzcloud/backups/20261121-141729/holzcloud.sqlite
1795270649.752485719 Backing up CSRF key...
1795270649.768639219 Backing up media...
1795270649.846221594 Backing up plugins...
1795270649.982894344 Writing checksums...
1795270650.301926928 Pruning backups older than 30 days...
1795270650.354835844 === Backup complete: /opt/holzcloud/backups/20261121-141729 ===
1795270650.371951386 21M	/opt/holzcloud/backups/20261121-141729
WALL 1795270649.284665677 1795270650.392376511
RC 0
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -X main.Version=$(git describe --tags --always) -X main.Commit=$(git rev-parse --short HEAD)" ./cmd/holzcloud
SKIP[synthetic-build] the binary of this update was built by this line in a scratch clone before the update
RC 0
$ scp holzcloud user@your-server:
RC 0
$ scp deploy/backup.sh deploy/restore.sh user@your-server:
RC 0
$ sudo systemctl stop holzcloud
WALL 1795270650.876610803 1795270651.027230720
RC 0
$ sudo cp ~/holzcloud /opt/holzcloud/holzcloud
WALL 1795270651.121686678 1795270651.250720970
RC 0
$ sudo install -m 755 -o root -g root ~/backup.sh ~/restore.sh /opt/holzcloud/deploy/
WALL 1795270651.354021886 1795270651.481183845
RC 0
$ sudo systemctl start holzcloud
WALL 1795270651.576186720 1795270651.720424095
RC 0
$ curl -fsS --retry 30 --retry-connrefused --retry-delay 1 http://127.0.0.1:8080/readyz | jq '{status, version, commit, integrity}'
1795270651.929310303 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 4 ms: Couldn't connect to server
1795270652.936624595 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270653.937742304 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270654.943627180 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270655.945284389 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270656.952681666 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270657.960558334 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270658.968233542 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270659.978226418 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270660.984212502 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270661.990338585 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270663.025576961 {
1795270663.041928503   "status": "ok",
1795270663.058242294   "version": "v0.1.4-266-g10e0ef7",
1795270663.075135086   "commit": "10e0ef7",
1795270663.091741586   "integrity": "ok"
1795270663.108048711 }
WALL 1795270651.836267970 1795270663.129262794
RC 0
$ sudo journalctl -u holzcloud --since "$(systemctl show -p ActiveEnterTimestamp --value holzcloud)" -o cat | grep -E 'pre-upgrade snapshot (written|kept)|migration applied' || echo 'no migration, no snapshot'
1795270663.415505003 {"time":"2026-11-21T14:17:31.935809095Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-266-g10e0ef7-20261121-141731.sqlite","level":73}
1795270663.433422003 {"time":"2026-11-21T14:17:31.976282595Z","level":"INFO","msg":"migration applied","version":74,"type":"sql","duration_ms":4}
WALL 1795270663.227556919 1795270663.454263003
RC 0
```

#### amd64-t34-upd/upd-6.out

```text
$ sudo cp -p /opt/holzcloud/holzcloud /opt/holzcloud/holzcloud.prev
WALL 1795270664.807393670 1795270664.935599254
RC 0
$ /opt/holzcloud/holzcloud version
1795270665.154448712 holzcloud v0.1.4-266-g10e0ef7 (10e0ef7)
WALL 1795270665.032175295 1795270665.176900254
RC 0
$ sudo -u holzcloud /opt/holzcloud/deploy/backup.sh /opt/holzcloud/backups
1795270665.434148129 === holzcloud-CMS backup: 20261121-141745 ===
1795270665.509304879 Backing up database...
1795270665.636314296 2026/11/21 14:17:45 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/data/holzcloud.sqlite
1795270665.681241629 2026/11/21 14:17:45 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/backups/20261121-141745/holzcloud.sqlite
1795270665.700503171 2026/11/21 14:17:45 INFO backup verified path=/opt/holzcloud/backups/20261121-141745/holzcloud.sqlite tables=67
1795270665.720656296 backup written and verified: /opt/holzcloud/backups/20261121-141745/holzcloud.sqlite
1795270665.736715254 Backing up CSRF key...
1795270665.752358546 Backing up media...
1795270665.836443462 Backing up plugins...
1795270665.963257629 Writing checksums...
1795270666.286191963 Pruning backups older than 30 days...
1795270666.340610296 === Backup complete: /opt/holzcloud/backups/20261121-141745 ===
1795270666.358719296 21M	/opt/holzcloud/backups/20261121-141745
WALL 1795270665.272155295 1795270666.380334546
RC 0
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -X main.Version=$(git describe --tags --always) -X main.Commit=$(git rev-parse --short HEAD)" ./cmd/holzcloud
SKIP[synthetic-build] the binary of this update was built by this line in a scratch clone before the update
RC 0
$ scp holzcloud user@your-server:
RC 0
$ scp deploy/backup.sh deploy/restore.sh user@your-server:
RC 0
$ sudo systemctl stop holzcloud
WALL 1795270666.862862671 1795270667.014111213
RC 0
$ sudo cp ~/holzcloud /opt/holzcloud/holzcloud
WALL 1795270667.111465713 1795270667.240802630
RC 0
$ sudo install -m 755 -o root -g root ~/backup.sh ~/restore.sh /opt/holzcloud/deploy/
WALL 1795270667.337106005 1795270667.465050505
RC 0
$ sudo systemctl start holzcloud
WALL 1795270667.561532755 1795270667.708200047
RC 0
$ curl -fsS --retry 30 --retry-connrefused --retry-delay 1 http://127.0.0.1:8080/readyz | jq '{status, version, commit, integrity}'
1795270667.915902005 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 5 ms: Couldn't connect to server
1795270668.918766380 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270669.928866923 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270670.931030631 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270671.935649757 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270672.941541216 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270673.943207383 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270674.951569175 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270675.958285342 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270676.961217176 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270677.963964593 curl: (7) Failed to connect to 127.0.0.1 port 8080 after 0 ms: Couldn't connect to server
1795270678.993537052 {
1795270679.012993677   "status": "ok",
1795270679.033416885   "version": "v0.1.4-267-g456235e",
1795270679.051081052   "commit": "456235e",
1795270679.068262969   "integrity": "ok"
1795270679.084771802 }
WALL 1795270667.819584588 1795270679.106112302
RC 0
$ sudo journalctl -u holzcloud --since "$(systemctl show -p ActiveEnterTimestamp --value holzcloud)" -o cat | grep -E 'pre-upgrade snapshot (written|kept)|migration applied' || echo 'no migration, no snapshot'
1795270679.387092594 {"time":"2026-11-21T14:17:47.915038213Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-267-g456235e-20261121-141747.sqlite","level":74}
1795270679.404573927 {"time":"2026-11-21T14:17:47.967620505Z","level":"INFO","msg":"migration applied","version":75,"type":"sql","duration_ms":4}
WALL 1795270679.198909094 1795270679.425482469
RC 0
```

#### amd64-t34-upd/upd-snapshots.txt

```text
after the install and run.sh seed (holzcloud v0.1.4-261-gcefe4ac (cefe4ac)), oldest first:
  pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141543.sqlite level 0
after update 1 (v0.1.4-262-g7283b76), oldest first:
  pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141543.sqlite level 0
  pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite level 69
after update 2 (v0.1.4-263-g5083f30), oldest first:
  pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141543.sqlite level 0
  pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite level 69
  pre-upgrade-v0.1.4-263-g5083f30-20261121-141644.sqlite level 70
after update 3 (v0.1.4-264-g1f4590b), oldest first:
  pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141543.sqlite level 0
  pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite level 69
  pre-upgrade-v0.1.4-263-g5083f30-20261121-141644.sqlite level 70
  pre-upgrade-v0.1.4-264-g1f4590b-20261121-141659.sqlite level 71
after update 4 (v0.1.4-265-g5f8e3f6), oldest first:
  pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141543.sqlite level 0
  pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite level 69
  pre-upgrade-v0.1.4-263-g5083f30-20261121-141644.sqlite level 70
  pre-upgrade-v0.1.4-264-g1f4590b-20261121-141659.sqlite level 71
  pre-upgrade-v0.1.4-265-g5f8e3f6-20261121-141715.sqlite level 72
after update 5 (v0.1.4-266-g10e0ef7), oldest first:
  pre-upgrade-v0.1.4-262-g7283b76-20261121-141628.sqlite level 69
  pre-upgrade-v0.1.4-263-g5083f30-20261121-141644.sqlite level 70
  pre-upgrade-v0.1.4-264-g1f4590b-20261121-141659.sqlite level 71
  pre-upgrade-v0.1.4-265-g5f8e3f6-20261121-141715.sqlite level 72
  pre-upgrade-v0.1.4-266-g10e0ef7-20261121-141731.sqlite level 73
after update 6 (v0.1.4-267-g456235e), oldest first:
  pre-upgrade-v0.1.4-263-g5083f30-20261121-141644.sqlite level 70
  pre-upgrade-v0.1.4-264-g1f4590b-20261121-141659.sqlite level 71
  pre-upgrade-v0.1.4-265-g5f8e3f6-20261121-141715.sqlite level 72
  pre-upgrade-v0.1.4-266-g10e0ef7-20261121-141731.sqlite level 73
  pre-upgrade-v0.1.4-267-g456235e-20261121-141747.sqlite level 74
```

#### amd64-t34-evidence-commit.txt

```text
cefe4ac
```

#### amd64-t34-pass1.status

```text
0
```
