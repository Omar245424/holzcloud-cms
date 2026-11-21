# A backup of v0.1.1 restored into 0.1.5, across nine migrations: results

This document is the permanent record of MIG-01 and MIG-02 of phase 34. It records a backup of the oldest beta, v0.1.1 at migration level 60, made on a server of that release by its own binary and its own `backup.sh`, carried to a second server, and restored there into the build of this release, at level 69, with the restore line of `deploy/DEPLOY.md`. It records what `restore.sh` printed, that the first start after the restore wrote one pre-upgrade snapshot of the restored database at level 60 and then applied the nine migrations 61 to 69, and that the restored site then passed the 17 assertions of `deploy/test/run.sh check`, with the seeded page, the uploaded picture and the old login still there (MIG-01); and it names both releases with their commits and both levels, each read from the database files and computed from the migration files of the two commits (MIG-02).

The record rests on one evidence set: a single `upgrade-verify.sh run` from fresh containers, made once against the committed text of every bound file (see *How this was run*); this record uses its stage `mig`. Its raw output is the Appendix, pasted unchanged, and `bash deploy/test/upgrade-verify.sh record-check` proves every block against its scratch file byte for byte. Every number in this document is a line of that set and cites its file in backticks as `DIRECTORY/FILE`, with `DIRECTORY` the scratch directory `amd64-t34-mig-src` (the source server, and the files the stage collected about the target); the run files beside it are named by their path below `deploy/test/out`. A file that is not pasted is named as `FILE in DIRECTORY`.

The record says what was done, how many times, and what was seen. It does not say that every older release restores this way: it played one pair, v0.1.1 into this release, once. What the containers cannot show is a named deviation or a line of *Not tested*, never a gap in the table.

## Summary

| Requirement | Result | Platform | Date, commit |
|-------------|--------|----------|--------------|
| MIG-01 | PASS: `MIG-01 PASS` stands in `amd64-t34-summary.txt` and in `amd64-t34-mig-src/mig.txt`. The backup made by v0.1.1's own `backup.sh` and binary ended `BACKUP-RC: 0` at `LEVEL-BEFORE: 60`; the document's restore line ended `RESTORE-RC: 0` with `RESTORE-NOTICE: yes` and `RESTORE-COMPLETE: yes`; the first start wrote one snapshot at level 60 (`SNAPSHOTS-AFTER: 1`, `SNAPSHOT-LEVEL: 60`, `JOURNAL-SNAPSHOT-LINES: 1`) and applied nine migrations (`JOURNAL-APPLIED-LINES: 9`, `JOURNAL-APPLIED-VERSIONS: 61 62 63 64 65 66 67 68 69`) to `LEVEL-AFTER: 69`; the old login, the seeded page and the same picture came back (`LOGIN-AFTER: 303 /admin/`, `PAGE-AFTER: 200`, `MEDIA-SHA-EQUAL: yes`), and `run.sh check` read `CHECK: 17 of 17` with `CHECK-RC: 0`, all in `amd64-t34-mig-src/mig.txt`. | Debian 12 container, systemd PID 1, linux/amd64 under emulation on an arm64 host (`amd64-t34-mig-src/versions.txt`) | 2026-11-21, cefe4ac |
| MIG-02 | PASS: `MIG-02 PASS` stands in `amd64-t34-summary.txt` and in `amd64-t34-mig-src/mig.txt`. The two releases are named with their commits (`OLD-RELEASE: v0.1.1 (5bcae9b)`, `OLD-COMMIT: 5bcae9b448db814701310c5dfab26031d68c0d92`, `NEW-RELEASE: v0.1.4-261-gcefe4ac (cefe4ac)`), both levels with what the migration files of the two commits give (`LEVEL-BEFORE: 60`, `LEVEL-EXPECTED-BEFORE: 60`, `LEVEL-AFTER: 69`, `LEVEL-EXPECTED-AFTER: 69`), and the backup's own `MANIFEST` names its maker (`MANIFEST-VERSION: holzcloud v0.1.1 (5bcae9b)`); what `restore.sh` printed is `amd64-t34-mig-src/mig-restore.out` and what the first start wrote is `amd64-t34-mig-src/mig-journal.txt`. | same | 2026-11-21, cefe4ac |

Evidence commit: cefe4ac

The binary under test is the one built from `NEW_REF` (pinned by `run` to the short sha of HEAD, `cefe4ac`, read from `deploy/test/out/amd64-t34-evidence-commit.txt`) with the Build line of `deploy/DEPLOY.md`: `BINARY: holzcloud v0.1.4-261-gcefe4ac (cefe4ac)` in `amd64-t34-summary.txt`, the same as `target binary:`, `NEW_REF:` and `evidence commit:` in `amd64-t34-mig-src/versions.txt`, and the `version` and `commit` of the `/readyz` answer that `restore.sh` printed (`amd64-t34-mig-src/mig-restore.out`). The old binary is v0.1.1, built by the same Build line from the commit of its tag: `source binary: holzcloud v0.1.1 (5bcae9b)` and `v0.1.1 commit: 5bcae9b448db814701310c5dfab26031d68c0d92` in `amd64-t34-mig-src/versions.txt`, `OLD-VERSION-LINE: holzcloud v0.1.1 (5bcae9b)` in `amd64-t34-mig-src/mig.txt`. The summary carries `STAGES: mig upd rb` and no `REHEARSAL` line.

## How this was run

`deploy/test/upgrade-verify.sh` is a sibling of `deploy/test/run.sh` in the shape of `deploy/test/backup-verify.sh`; it edits neither and reaches the harness through `run.sh`'s public subcommands only (`build`, `up`, `install`, `seed`, `transfer`, `check`). `run` refuses to start unless `PLATFORM=linux/amd64` and `TARGET=systemd-base`, refuses while a bound file or the binary's source differs from HEAD or is untracked, pins `NEW_REF` to the short sha of HEAD so that the binary and the scripts are one commit, and plays the stages `mig`, `upd` and `rb`, each in its own fresh containers. The verdict of each requirement is computed from the `KEY: value` files the stages wrote, and a stage that did not run reads NOT-RUN, so a trace can never pass for the set. Stage `mig` used two containers, `holzcloud-t34-mig-src` for v0.1.1 and `holzcloud-t34-mig-tgt` for this release; the stages `upd` and `rb` of the same set are recorded in `deploy/UPDATE-CHAIN-RUN.md` and `deploy/ROLLBACK-RUN.md`.

The set was made once, from fresh containers, on the clean committed tree of `cefe4ac` without `REHEARSAL`, started 2026-11-21T14:13:05Z and ended 14:24:59Z on the host clock (11 minutes 54 seconds for all three stages), exit status 0 (`deploy/test/out/amd64-t34-pass1.started`, `deploy/test/out/amd64-t34-pass1.finished` and `deploy/test/out/amd64-t34-pass1.status`). Before it the worktree held no earlier output of the driver and no `holzcloud-t34-*` container existed, and after it `bash deploy/test/upgrade-verify.sh down` removed the four containers. In the same window and at the same commit, `PLATFORM=linux/amd64 TARGET=systemd-base bash deploy/test/backup-verify.sh run` made the set of phase 32 again (14:25:11Z to 14:27:08Z, status 0, all three verdicts PASS; that set is recorded in `deploy/BACKUP-VERIFY-RUN.md`). No second set was made against `cefe4ac`, and no other set is cited. Two earlier sets were superseded. The first was made against `36b67cc`; the commit `00854be` changed `deploy/DEPLOY.md` and `deploy/test/upgrade-verify.sh` after it. The second was made against `00854be`, and this record was first written from it. Its four verdicts read PASS, but the browser record of phase 34 then found that a plugin switched on in the admin did not start when the browser dropped the request that switched it on, because the plugin was loaded with that request's context and closed with it. The commit `cefe4ac` fixed that in `internal/plugin/runtime.go`, so the binary changed after the set of `00854be`, and the set was made again. So this record was written from this set, and nothing of either earlier set is quoted. The command was:

```text
PLATFORM=linux/amd64 TARGET=systemd-base bash deploy/test/upgrade-verify.sh run
```

The knob not set on that line had its default: `STAGES=mig upd rb`. The log of the run is `deploy/test/out/amd64-t34-pass1.log`, which is not pasted; its first line reads `upgrade-verify: platform linux/amd64, NEW_REF pinned to cefe4ac, v0.1.1 is 5bcae9b448db814701310c5dfab26031d68c0d92, stages: mig upd rb`.

No bound file changed after the evidence commit. The bound files for this set are `deploy/DEPLOY.md`, `deploy/backup.sh`, `deploy/restore.sh`, `deploy/holzcloud.service`, `deploy/test/run.sh`, `deploy/test/Dockerfile`, `deploy/test/upgrade-verify.sh` and the five plugin archives under `plugins/`, together with the binary built from `cmd`, `internal`, `go.mod`, `go.sum`, `changelog.go` and the `CHANGELOG.md` it embeds; the old release is bound by the commit of its tag, which `amd64-t34-mig-src/mig.txt` names. This command exited 0 when the record was written:

```text
git diff --quiet cefe4ac HEAD -- deploy/DEPLOY.md deploy/backup.sh deploy/restore.sh deploy/holzcloud.service deploy/test/run.sh deploy/test/Dockerfile deploy/test/upgrade-verify.sh plugins/bestellung/bestellung.zip plugins/jahreszahl/jahreszahl.zip plugins/kontaktformular/kontaktformular.zip plugins/nicht-gefunden/nicht-gefunden.zip plugins/suche/suche.zip cmd internal go.mod go.sum changelog.go CHANGELOG.md
```

Before the run, `bash deploy/test/upgrade-verify.sh selftest` ended `SELFTEST OK (409 cases)` and `bash deploy/test/backup-verify.sh selftest` ended `SELFTEST OK (117 cases)`; neither output is part of the set.

## Deviations named

Everything here precedes every result: it is what the containers do differently from two real servers and what the harness did that the document does not contain. The first column is the name written to `amd64-t34-mig-src/deviations.txt` (the source's lines first, then the target's); the kinds of the whole run are the `DEVIATIONS:` line of `amd64-t34-summary.txt`, and `synthetic-build` belongs to the stages `upd` and `rb` only.

| Name | Why | What a real server does instead |
|------|-----|---------------------------------|
| DEVIATION[emulation] | The Docker host is aarch64 and the linux/amd64 containers run under emulation (`host architecture: aarch64`, `container uname -m: x86_64` in `amd64-t34-mig-src/versions.txt`), so every wall time of this set, the first start after the restore among them, is an emulation figure | A native amd64 CPU and disk |
| DEVIATION[old-plugins] | The v0.1.1 source was seeded with the five plugin archives of v0.1.1 itself (`git show v0.1.1:plugins/ID/ID.zip`, their sha256 in the `OLD-PLUGIN:` lines of `amd64-t34-mig-src/mig.txt`), because v0.1.1 refuses today's `kontaktformular` and `bestellung` archives: their manifests declare the hook `privacy`, which its manifest check does not know. After the restore the host of this release runs those v0.1.1 modules | A server on v0.1.1 has the plugins of its own time, and an update to 0.1.5 keeps them until the operator installs newer ones |
| DEVIATION[old-script] | The `backup.sh` that today's Install text copied was replaced by the `backup.sh` of v0.1.1 (`git show v0.1.1:deploy/backup.sh`, root, mode 755) before the manual line was played; its sha256 is `BACKUP-SCRIPT` in `amd64-t34-mig-src/mig.txt` | A server still on v0.1.1 has that script from its own install |
| DEVIATION[fixture] | Two lines. The five plugin archives were copied into `/opt/holzcloud/plugin-zips` before `run.sh seed`, and the login user was made with the CLI of v0.1.1; with the website, page, AI key and test image that `run.sh seed` makes, they stand in for operator work. The login password file, the AI key file and the test image cross from container to container by exec-pipe and are part of no backup. No order row is written (`ORDER: none` in `amd64-t34-mig-src/mig.txt`) | The operator's own content, users and plugins |
| DEVIATION[browser] | The logins before the backup and after the restore were played with `curl`, with the CSRF field read out of the page and, because the shipped unit sets `HOLZCLOUD_SECURE=true`, a Referer with the https origin on a plain http post | A browser behind the reverse proxy |
| DEVIATION[backup-name] | The sample timestamp of the document's restore line (`sudo /opt/holzcloud/deploy/restore.sh /opt/holzcloud/backups/20260802-031500`) was replaced by the transferred directory `/opt/holzcloud/backups/20261121-141407`, and `yes` was typed on stdin | The operator types the name of the backup they restore |
| DEVIATION[transfer] | The backup travelled as a tar stream through the host (`docker cp SOURCE:dir -` into `docker cp - TARGET:/opt/holzcloud/backups`), proven lossless by the sha256 of every file on both sides (`TRANSFER-RC: 0`) | `rsync` or `scp` between two machines, which the document leaves to the operator |
| DEVIATION[scp] | Three lines: `scp holzcloud user@your-server:` played as `docker cp` of the v0.1.1 binary to the source and of this release's binary to the target, and `scp -r deploy/ user@your-server:` played as a tar stream of the committed `deploy` directory (without `deploy/test/out`) into `/root` | `scp` from the build machine into the operator's home directory |
| DEVIATION[start-limit] | `systemctl reset-failed holzcloud` ran before the restore and before `run.sh check`, because the install start, the restore's start and check's restart come closer together than on a server | One restore, by hand, which never meets the limit |
| DEVIATION[dropin] | `/etc/systemd/system/holzcloud.service.d/container.conf` sets `ProcSubset=all`, because the document's unchanged start crashed in the container (`unit-start-unchanged.txt` in `amd64-t34-mig-src`) | The unit as shipped |
| DEVIATION[mount] | The unit's `RequiresMountsFor` names `/opt/holzcloud/data`, a plain directory on the container's root filesystem, so that line was not exercised, and the restore renamed the data directory aside (`Moving current data to /opt/holzcloud/data.before-restore-20261121-141448` in `amd64-t34-mig-src/mig-restore.out`) | A data directory on its own mount |

## Findings

Each finding has the class of the earlier records' rule: (1) a bug of the harness, fixed in the driver; (2) a limit of the container, a named deviation or a *Not tested* line; (3) a measured property of a shipped file that matches what it and its document say, recorded and not edited; (4) a defect, a shipped file doing less than it or its document says, with a fix in a commit of its own. No finding of this set is class 1 or class 4, and no shipped file was edited after the evidence commit. The defect that the browser record found at `00854be`, a plugin load ended by a dropped request, was fixed in `cefe4ac` before this set (*How this was run*).

The test for class 4 was whether a result contradicts something `deploy/restore.sh`, `deploy/DEPLOY.md` or the binary says at `cefe4ac`, or leaves the operator worse off than they say. The document says that a backup made before 0.1.5 is restored with a notice, that migrations run automatically on startup, that the first start writes one `pre-upgrade snapshot written` line with the level it holds and then one `migration applied` per migration, and that `restore.sh` waits for `/readyz` for `READY_SECONDS`. F-1 to F-4 are those sentences observed; F-5 is what the pair shows about plugins.

**F-1 (class 3): A backup made by v0.1.1 has no checksum file, and restore.sh took the path the document gives it: a notice that its files were not verified, then the snapshot check.** The backup directory holds no `SHA256SUMS` (`SUMS-IN-BACKUP: no` in `amd64-t34-mig-src/mig.txt`; the v0.1.1 script wrote no `Writing checksums...` line, `amd64-t34-mig-src/mig-backup.out`). The line (`amd64-t34-mig-src/mig-restore.out`):

```text
1795270488.423968589 NOTICE: /opt/holzcloud/backups/20261121-141407 has no checksum file (made before 0.1.5); its files were not verified
```

**F-2 (class 3): The snapshot check of restore.sh ran the nine migrations on a temporary copy of the backup before it stopped the service, so the crossing was tried once on a copy before the live data was touched.** Between `Verifying snapshot...` and `integrity: ok`, the binary of this release applied the versions 61 to 69 to `/tmp/tmp.G0PVKINhfy/holzcloud.sqlite` and the row check read `files: 3 rows checked, 0 missing`; only then came `Stopping holzcloud...` (`amd64-t34-mig-src/mig-restore.out`). The live database was migrated again, by the service, at its first start (F-3).

**F-3 (class 3): The first start after the restore wrote one pre-upgrade snapshot of the restored database at level 60 and applied the migrations 61 to 69, one journal line each.** The two ends of the journal of that start (`amd64-t34-mig-src/mig-journal.txt`):

```text
{"time":"2026-11-21T14:14:49.876991714Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141449.sqlite","level":60}
{"time":"2026-11-21T14:14:49.934072298Z","level":"INFO","msg":"migration applied","version":69,"type":"sql","duration_ms":5}
```

The snapshot is the only one in the data directory after the start and holds level 60 when read from the file (`SNAPSHOTS-AFTER: 1`, `SNAPSHOT-FILE`, `SNAPSHOT-LEVEL: 60` in `amd64-t34-mig-src/mig.txt`), so the database as the backup held it is still on the server after the crossing.

**F-4 (class 3): The first start fitted inside the readiness wait of restore.sh with room to spare, under emulation.** The first start took 13.3 seconds against the 30 seconds `restore.sh` waits (`FIRST-START-S: 13.3`, `READY-SECONDS: 30` in `amd64-t34-mig-src/mig.txt`), and `Waiting for readiness...` to `=== Restore complete ===` read 1795270489.548681922 to 1795270502.749580304 (`amd64-t34-mig-src/mig-restore.out`).

**F-5 (class 3): The plugin modules of v0.1.1 load, run and answer on the host of this release after the restore, and the plugin archives of this release do not install on v0.1.1.** After the restore, `run.sh check` read `PASS loaded`, `PASS running` and `PASS invoke` for all five plugins, `PASS no-load-failures` and `PASS stable` (`amd64-t34-mig-src/mig-check.txt`), with the plugins the backup carried, the archives of v0.1.1 (`DEVIATION[old-plugins]`). The other direction does not hold: v0.1.1 refuses today's `kontaktformular` and `bestellung` archives because their manifests declare the hook `privacy` (`amd64-t34-mig-src/deviations.txt`). That is the plugin interface changing between releases, which is why it stays outside what 1.0 promises.

## MIG-01: a backup of an older release restores into this one, and the site passes its checks

**What was done.** On the source container, today's Install text installed v0.1.1, built from the commit of its tag; the driver seeded it with v0.1.1's own plugin archives and a login user made by v0.1.1's CLI (`SEED-RC: 0`), and read the login, the seeded page and the uploaded picture over HTTP (`LOGIN-BEFORE: 303 /admin/`, `PAGE-BEFORE: 200`, `MEDIA-BEFORE: HTTP 200 925d26e12cb8ed51e272c6cadba3f6fb5d32f4af8bf06d1d2123ddbf8a176577` in `amd64-t34-mig-src/mig.txt`). The document's manual backup line was then played by its exact text with v0.1.1's `backup.sh` and binary (`BACKUP-RC: 0`, `BACKUP-DIR: /opt/holzcloud/backups/20261121-141407`); it ends `=== Backup complete: /opt/holzcloud/backups/20261121-141407 ===` and `20M` (`amd64-t34-mig-src/mig-backup.out`). The backup was carried to the target container (`TRANSFER-RC: 0`), which runs this release, and restored there with the document's restore line, `yes` on stdin, every line stamped with the container clock (`amd64-t34-mig-src/mig-restore.out`): the header names the backup's version, `version=holzcloud v0.1.1 (5bcae9b)`; the notice of F-1; the snapshot check with the nine migrations on a copy (F-2), `integrity: ok` and `files: 3 rows checked, 0 missing`; the stop, the move of the old data aside, the restore of the database, the CSRF key, the media and the plugins; the start, `Waiting for readiness...` and:

```text
1795270502.749580304 === Restore complete ===
1795270502.802876720 {"status":"ok","version":"v0.1.4-261-gcefe4ac","commit":"cefe4ac","uptime_s":1,"integrity":"ok","disk_free_bytes":1584480804864,"database":{"freelist_count":0,"page_count":180,"size_bytes":667648,"wal_bytes":362592}}
```

The first start wrote the snapshot and applied the nine migrations (F-3, `amd64-t34-mig-src/mig-journal.txt`). Then the site was read again: the setup form is gone (`SETUP-AFTER: 404`), the user of v0.1.1 logs in with the old password (`LOGIN-AFTER: 303 /admin/`) and is sent on to set up a second factor (`LOGIN-NEXT-AFTER: 303 /admin/2fa/einrichten`), the seeded page answers (`PAGE-AFTER: 200`), the picture answers with the same bytes (`MEDIA-AFTER` and `MEDIA-SHA-EQUAL: yes`), and `run.sh check` passed all of its assertions (`CHECK-RC: 0`, `CHECK: 17 of 17` in `amd64-t34-mig-src/mig.txt`; the 17 `PASS` lines in `amd64-t34-mig-src/mig-check.txt`).

*What it shows.* A backup made on v0.1.1 at migration level 60 by that release's own `backup.sh` and binary was restored into this release with the document's restore line and ended with status 0; the first start wrote one pre-upgrade snapshot at level 60, applied the nine migrations 61 to 69, and the restored site passed all 17 checks of `run.sh check`, with the old login, the seeded page and the same picture.

**What it does not show.** A restore on a real server, a database of production size, and the releases 0.1.2 to 0.1.4 as sources (*Not tested*). A second factor set up on v0.1.1 and carried across, which this set did not have: the user was sent to set one up.

## MIG-02: the record names both releases, both levels and what the restore printed

**What was done.** The stage wrote both releases with their commits into `amd64-t34-mig-src/mig.txt`: `OLD-RELEASE: v0.1.1 (5bcae9b)` with `OLD-COMMIT: 5bcae9b448db814701310c5dfab26031d68c0d92` and `OLD-VERSION-LINE: holzcloud v0.1.1 (5bcae9b)` from the old binary, and `NEW-RELEASE: v0.1.4-261-gcefe4ac (cefe4ac)`. Neither level was typed. The level before was read from the backup's database file through the immutable URI (`LEVEL-BEFORE: 60`) and the level after from the live database through a read-only connection (`LEVEL-AFTER: 69`), and each was compared with the number of migration files of its commit (`LEVEL-EXPECTED-BEFORE: 60`, `LEVEL-EXPECTED-AFTER: 69`). The backup's own `MANIFEST` names v0.1.1 as its maker (`MANIFEST-VERSION: holzcloud v0.1.1 (5bcae9b)`), and the backup script is bound by its sha256 (`BACKUP-SCRIPT: v0.1.1 sha256 ac02fac9096d3ed1665794c87a1624f8c896bb980ae615d80be0d4bab7d8de96`). What `restore.sh` printed is pasted whole as `amd64-t34-mig-src/mig-restore.out`, and the journal lines of the first start as `amd64-t34-mig-src/mig-journal.txt`, one `pre-upgrade snapshot written` with `"level":60` and nine `migration applied` with the versions 61 to 69 (`JOURNAL-SNAPSHOT-LINES: 1`, `JOURNAL-APPLIED-LINES: 9`, `JOURNAL-APPLIED-VERSIONS: 61 62 63 64 65 66 67 68 69`).

*What it shows.* The restore crossed from v0.1.1 (commit 5bcae9b, migration level 60, read from the backup's database file) to this release (commit cefe4ac, level 69, read from the restored database), and `restore.sh` printed the notice that a backup made before 0.1.5 has no checksum file, verified the snapshot, restored it and ended with `=== Restore complete ===` and a `/readyz` answer of this release with `status` and `integrity` `ok`.

## Measurements

The figures the harness wrote. Every one is an emulation figure on an aarch64 host (`amd64-t34-mig-src/versions.txt`); none is a promise of speed.

| Figure | Value | Where |
|--------|-------|-------|
| The manual backup on v0.1.1 | 1795270446.996067666 to 1795270448.263606459 wall, 20M | the `WALL` line and the last line of `amd64-t34-mig-src/mig-backup.out` |
| The whole restore, from its first line to its last | 1795270488.107763797 to 1795270502.877167554 wall | the `WALL` line of `amd64-t34-mig-src/mig-restore.out` |
| The first start after the restore | 13.3 s, against a readiness wait of 30 s | `FIRST-START-S`, `READY-SECONDS` in `amd64-t34-mig-src/mig.txt` |
| The snapshot and the nine migrations of the first start | 14:14:49.876991714Z to 14:14:49.934072298Z | `amd64-t34-mig-src/mig-journal.txt` |
| The tables of the v0.1.1 database | 55 | `backup verified ... tables=55` in `amd64-t34-mig-src/mig-backup.out` |
| The restored database after its first start | 180 pages, 667648 bytes | the `/readyz` answer in `amd64-t34-mig-src/mig-restore.out` |
| The median of ten requests per plugin route after the restore | 0.010560 s to 0.017325 s | the `LATENCY` lines of `amd64-t34-mig-src/mig-check.txt` |
| The whole run, all three stages | 14:13:05Z to 14:24:59Z | `deploy/test/out/amd64-t34-pass1.started`, `deploy/test/out/amd64-t34-pass1.finished` |

## Container vs. server

| Aspect | In the containers | On real servers | Measured or assumed |
|--------|-------------------|-----------------|---------------------|
| Speed and timing | amd64 emulated on aarch64 | A native server | Not measured natively |
| The two servers | Two containers on one host, the backup carried as a tar stream through the host | Two machines, the backup carried by `rsync` or `scp` | Measured for the restore; the transfer is proven lossless by sha256, the tool is assumed |
| The source's backup script | v0.1.1's `backup.sh` put in place of today's after today's Install text | The script of v0.1.1's own install | Measured for the backup; the install path is assumed equivalent |
| The plugins | v0.1.1's own archives | Whatever the operator installed on v0.1.1 | Measured for these five |
| The content | One seeded website, page, picture and user | Years of pages, media and orders | Measured for this size only |
| `holzcloud.service` | With the container-only drop-in `ProcSubset=all` | As shipped | Named deviation |

## Not tested

- A restore on a real server (SRV-02) was not played; both servers were containers on the Mac.
- Native amd64 timing was not measured.
- The releases 0.1.2, 0.1.3 and 0.1.4 were not played as sources; only v0.1.1, the oldest beta, was.
- A backup made by the timer at night was not played; the backup of this set was made by the document's manual line.
- A database of production size was not played; the restored database is 667648 bytes.
- The source's own `postinstall` backup line was not played on v0.1.1: `run.sh postinstall` would run today's `backup.sh` against the v0.1.1 binary.
- An order, form messages and a second factor were not carried across: the seeded source had no order row (`ORDER: none`), and the user had no second factor.
- A restore that fails part-way through its migrations was not played here; `deploy/ROLLBACK-RUN.md` records a failing migration at an update.

## Owner steps

None. Every step of this record ran in containers on the Mac, by the owner's decision; nothing here has to be done on the real server before the record counts.

## Hand-over to docs/beta.md

The sentences of this record that the page may quote, named by place and not copied here, so that each stays in exactly one place: MIG-01, the title sentences of F-1, F-2 and F-3 and the sentence after *What it shows.* in *MIG-01*; MIG-02, the sentence after *What it shows.* in *MIG-02*; the plugin interface, the title sentence of F-5; the limits, the bullets of *Not tested*.

The limit of the page this record answers, named by its first words: "A restore that crosses a migration was not played" (a backup of v0.1.1 at level 60 was restored into this release at level 69, F-3 and *MIG-01*).

Nothing in this record is a promise until the page cites it byte for byte. A change to a bound file after the evidence commit means the set is made again before the record is cited for the changed text. `docs/beta.md` was not edited by this plan.

## Appendix: raw output

Each block is the scratch file of the evidence set, pasted unchanged by `bash deploy/test/upgrade-verify.sh appendix mig`, and `bash deploy/test/upgrade-verify.sh record-check` compares every block with its scratch file byte for byte. The files are in `deploy/test/out` (gitignored); the files of the target container (`amd64-t34-mig-tgt`), the `step-NN.txt` files, `seed.txt`, `baseline.txt`, `http.txt`, `install.txt`, the unit-start files, the v0.1.1 archives and script and the log are not pasted.

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

#### amd64-t34-mig-src/versions.txt

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
source binary:         holzcloud v0.1.1 (5bcae9b)
target binary:         holzcloud v0.1.4-261-gcefe4ac (cefe4ac)
v0.1.1 commit:         5bcae9b448db814701310c5dfab26031d68c0d92
NEW_REF:               cefe4ac
evidence commit:       cefe4ac
```

#### amd64-t34-mig-src/deviations.txt

```text
DEVIATION[old-plugins]: the v0.1.1 source is seeded with the plugin archives of v0.1.1 itself (git show v0.1.1:plugins/ID/ID.zip), not with this tree's: v0.1.1 refuses today's kontaktformular and bestellung archives, whose manifests declare the hook 'privacy' that its manifest check does not know (finding F-4). After the restore the host of this tree runs those v0.1.1 modules
DEVIATION[fixture]: the five plugin archives are copied into /opt/holzcloud/plugin-zips before run.sh seed, and the login user is made with the CLI of v0.1.1 (password in a root-only file, given on stdin); with the website, page, AI key and test image run.sh seed makes, they stand in for operator work in the admin and are no step of the document. No order row is written (ORDER: none)
DEVIATION[browser]: the login before the backup and after the restore is played with curl, not in a browser; the shipped unit sets HOLZCLOUD_SECURE=true, so the post over plain http carries a Referer with the https origin, and the hidden CSRF field is read out of the page
DEVIATION[old-script]: /opt/holzcloud/deploy/backup.sh is replaced by the backup.sh of v0.1.1 (git show v0.1.1:deploy/backup.sh, root 755) before the manual line is played, because today's Install text copied this tree's; a server still on v0.1.1 has that script there from its own install
DEVIATION[backup-name]: the sample timestamp of the document's restore line is replaced by the transferred directory /opt/holzcloud/backups/20261121-141407; the answer 'yes' is typed on stdin
DEVIATION[start-limit]: systemctl reset-failed holzcloud runs before the restore and before run.sh check; the install start, the restore's start and check's restart come closer together than on a server
DEVIATION[emulation]: the Docker host is aarch64; the linux/amd64 containers run under emulation, so every wall time of this set (the first start after the restore among them) is an emulation figure, not a server's
DEVIATION[scp]: 'scp holzcloud user@your-server:' was played as docker cp of deploy/test/out/bin/holzcloud-old-v0.1.1 to /root/holzcloud in the container (there is no ssh and no second machine here)
DEVIATION[scp]: 'scp -r deploy/ user@your-server:' was played as a tar stream of the working tree's deploy directory (without deploy/test/out) into /root in the container, by docker cp -
DEVIATION[dropin]: /etc/systemd/system/holzcloud.service.d/container.conf sets ProcSubset=all (container only; the document's unchanged start crashed, see unit-start-unchanged.txt)
DEVIATION[mount]: the unit's RequiresMountsFor names /opt/holzcloud/data, a plain directory on the container's root filesystem; no separate mount exists here, so that line was not exercised
DEVIATION[scp]: 'scp holzcloud user@your-server:' was played as docker cp of deploy/test/out/bin/holzcloud-new-cefe4ac to /root/holzcloud in the container (there is no ssh and no second machine here)
DEVIATION[transfer]: the backup travels as a tar stream through the host (docker cp SOURCE:dir - into docker cp - TARGET:/opt/holzcloud/backups); on two real machines it would be rsync or scp, which the document leaves to the operator. Proof of a lossless transfer is the sha256 of every file on both sides
DEVIATION[fixture]: the login password file, the AI key file and the test image are fixtures that stand in for what an operator remembers; none of them is part of any backup, and they cross from container to container by exec-pipe, never through a host file
```

#### amd64-t34-mig-src/mig.txt

```text
OLD-RELEASE: v0.1.1 (5bcae9b)
OLD-COMMIT: 5bcae9b448db814701310c5dfab26031d68c0d92
OLD-VERSION-LINE: holzcloud v0.1.1 (5bcae9b)
NEW-RELEASE: v0.1.4-261-gcefe4ac (cefe4ac)
OLD-PLUGIN: bestellung sha256 2ebdb9da08e322ab1c3312e742309222f22496fdd4deb3bf13e21ed6eb75eceb
OLD-PLUGIN: jahreszahl sha256 3bd0ca9de1d077e9ccc5ae49d0b9fad95d5cdbf724efc844f472312d5bbde1d2
OLD-PLUGIN: kontaktformular sha256 062eb61c4841b2e8cddb1321725368a6fd3b380f139550ce9e232ec36e03af9a
OLD-PLUGIN: nicht-gefunden sha256 5863d48f398ce603aee47afc9a5c401e5b0be9772e45c105453ef595ddb8b78a
OLD-PLUGIN: suche sha256 e2d1a5247c8972a83f13e0a41cf58687efc51e174835999270ed7a6a44904f14
SEED-RC: 0
LOGIN-BEFORE: 303 /admin/
PAGE-BEFORE: 200
MEDIA-PATH: /media/1/00b72934f04ab1e713522f67a6c75537.jpg
MEDIA-BEFORE: HTTP 200 925d26e12cb8ed51e272c6cadba3f6fb5d32f4af8bf06d1d2123ddbf8a176577
BACKUP-SCRIPT: v0.1.1 sha256 ac02fac9096d3ed1665794c87a1624f8c896bb980ae615d80be0d4bab7d8de96
BACKUP-RC: 0
BACKUP-DIR: /opt/holzcloud/backups/20261121-141407
MANIFEST-VERSION: holzcloud v0.1.1 (5bcae9b)
SUMS-IN-BACKUP: no
LEVEL-BEFORE: 60
LEVEL-EXPECTED-BEFORE: 60
TRANSFER-RC: 0
RESTORE-RC: 0
RESTORE-NOTICE: yes
RESTORE-COMPLETE: yes
READY-SECONDS: 30
FIRST-START-S: 13.3
SNAPSHOTS-AFTER: 1
SNAPSHOT-FILE: /opt/holzcloud/data/pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141449.sqlite
SNAPSHOT-LEVEL: 60
JOURNAL-SNAPSHOT-LINES: 1
JOURNAL-APPLIED-LINES: 9
JOURNAL-APPLIED-VERSIONS: 61 62 63 64 65 66 67 68 69
LEVEL-AFTER: 69
LEVEL-EXPECTED-AFTER: 69
SETUP-AFTER: 404
LOGIN-AFTER: 303 /admin/
LOGIN-NEXT-AFTER: 303 /admin/2fa/einrichten
PAGE-AFTER: 200
MEDIA-AFTER: HTTP 200 925d26e12cb8ed51e272c6cadba3f6fb5d32f4af8bf06d1d2123ddbf8a176577
MEDIA-SHA-EQUAL: yes
CHECK-RC: 0
CHECK: 17 of 17
ORDER: none (the stage writes no order row; transfer only carries the line)
MIG-01 PASS
MIG-02 PASS
```

#### amd64-t34-mig-src/mig-backup.out

```text
1795270447.194285792 === holzcloud-CMS backup: 20261121-141407 ===
1795270447.328697708 Backing up database...
1795270447.530897083 2026/11/21 14:14:07 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/data/holzcloud.sqlite
1795270447.625190292 2026/11/21 14:14:07 INFO sqlite opened journal_mode=wal path=/opt/holzcloud/backups/20261121-141407/holzcloud.sqlite
1795270447.651378917 2026/11/21 14:14:07 INFO backup verified path=/opt/holzcloud/backups/20261121-141407/holzcloud.sqlite tables=55
1795270447.679014292 backup written and verified: /opt/holzcloud/backups/20261121-141407/holzcloud.sqlite
1795270447.712188500 Backing up CSRF key...
1795270447.741007333 Backing up media...
1795270447.837439917 Backing up plugins...
1795270448.188205459 Pruning backups older than 30 days...
1795270448.209007209 === Backup complete: /opt/holzcloud/backups/20261121-141407 ===
1795270448.230059417 20M	/opt/holzcloud/backups/20261121-141407
WALL 1795270446.996067666 1795270448.263606459
```

#### amd64-t34-mig-src/mig-restore.out

```text
1795270488.302617088 === Restoring from /opt/holzcloud/backups/20261121-141407 ===
1795270488.332649047 timestamp=20261121-141407
1795270488.354866714 version=holzcloud v0.1.1 (5bcae9b)
1795270488.380400880 data_dir=/opt/holzcloud/data
1795270488.402844130 host=d2d1bc24e975
1795270488.423968589 NOTICE: /opt/holzcloud/backups/20261121-141407 has no checksum file (made before 0.1.5); its files were not verified
1795270488.453830714 Verifying snapshot...
1795270488.770740797 2026/11/21 14:14:48 INFO sqlite opened journal_mode=wal path=/tmp/tmp.G0PVKINhfy/holzcloud.sqlite
1795270488.793895630 2026/11/21 14:14:48 INFO migration applied version=61 type=sql duration_ms=18
1795270488.820879130 2026/11/21 14:14:48 INFO migration applied version=62 type=sql duration_ms=3
1795270488.841901380 2026/11/21 14:14:48 INFO migration applied version=63 type=sql duration_ms=8
1795270488.864876047 2026/11/21 14:14:48 INFO migration applied version=64 type=sql duration_ms=1
1795270488.889866089 2026/11/21 14:14:48 INFO migration applied version=65 type=sql duration_ms=0
1795270488.917202172 2026/11/21 14:14:48 INFO migration applied version=66 type=sql duration_ms=0
1795270488.942475505 2026/11/21 14:14:48 INFO migration applied version=67 type=sql duration_ms=4
1795270488.979490380 2026/11/21 14:14:48 INFO migration applied version=68 type=sql duration_ms=4
1795270489.007956422 2026/11/21 14:14:48 INFO migration applied version=69 type=sql duration_ms=5
1795270489.039678256 integrity: ok
1795270489.060549381 files: 3 rows checked, 0 missing
1795270489.086892964 Stopping holzcloud...
1795270489.109451797 Moving current data to /opt/holzcloud/data.before-restore-20261121-141448
1795270489.137811881 Restoring database...
1795270489.172467256 Restoring CSRF key...
1795270489.191605172 Restoring media...
1795270489.230056006 Restoring plugins...
1795270489.491412714 Starting holzcloud...
1795270489.548681922 Waiting for readiness...
1795270502.749580304 === Restore complete ===
1795270502.802876720 {"status":"ok","version":"v0.1.4-261-gcefe4ac","commit":"cefe4ac","uptime_s":1,"integrity":"ok","disk_free_bytes":1584480804864,"database":{"freelist_count":0,"page_count":180,"size_bytes":667648,"wal_bytes":362592}}
1795270502.833908012 
WALL 1795270488.107763797 1795270502.877167554
```

#### amd64-t34-mig-src/mig-journal.txt

```text
{"time":"2026-11-21T14:14:49.876991714Z","level":"INFO","msg":"pre-upgrade snapshot written","path":"/opt/holzcloud/data/pre-upgrade-v0.1.4-261-gcefe4ac-20261121-141449.sqlite","level":60}
{"time":"2026-11-21T14:14:49.932408506Z","level":"INFO","msg":"migration applied","version":61,"type":"sql","duration_ms":10}
{"time":"2026-11-21T14:14:49.934048714Z","level":"INFO","msg":"migration applied","version":62,"type":"sql","duration_ms":0}
{"time":"2026-11-21T14:14:49.934056214Z","level":"INFO","msg":"migration applied","version":63,"type":"sql","duration_ms":4}
{"time":"2026-11-21T14:14:49.934059548Z","level":"INFO","msg":"migration applied","version":64,"type":"sql","duration_ms":0}
{"time":"2026-11-21T14:14:49.934062339Z","level":"INFO","msg":"migration applied","version":65,"type":"sql","duration_ms":0}
{"time":"2026-11-21T14:14:49.934064548Z","level":"INFO","msg":"migration applied","version":66,"type":"sql","duration_ms":0}
{"time":"2026-11-21T14:14:49.934066964Z","level":"INFO","msg":"migration applied","version":67,"type":"sql","duration_ms":9}
{"time":"2026-11-21T14:14:49.934069506Z","level":"INFO","msg":"migration applied","version":68,"type":"sql","duration_ms":4}
{"time":"2026-11-21T14:14:49.934072298Z","level":"INFO","msg":"migration applied","version":69,"type":"sql","duration_ms":5}
```

#### amd64-t34-mig-src/mig-check.txt

```text
PASS loaded bestellung
PASS loaded jahreszahl
PASS loaded kontaktformular
PASS loaded nicht-gefunden
PASS loaded suche
PASS no-load-failures
PASS running bestellung
PASS running jahreszahl
PASS running kontaktformular
PASS running nicht-gefunden
PASS running suche
PASS invoke jahreszahl
PASS invoke kontaktformular
PASS invoke bestellung
PASS invoke suche
PASS invoke nicht-gefunden
PASS stable
engine: compiler (11 anonymous executable mappings)
# hook budget CallTimeout is 2 s (internal/plugin/runtime.go); the medians are of ten requests
LATENCY /plugintest median_s=0.017325 n=10
LATENCY /suche?q=Plugintest median_s=0.015193 n=10
LATENCY /formular median_s=0.010560 n=10
LATENCY /bestellung median_s=0.010583 n=10
Seccomp:	0
Seccomp_filters:	0
NRestarts=0
Result=success
{"time":"2026-11-21T14:15:08.973601417Z","level":"INFO","msg":"plugin loaded","plugin":"kontaktformular","version":"1.5.0"}
{"time":"2026-11-21T14:15:10.929511502Z","level":"INFO","msg":"plugin loaded","plugin":"bestellung","version":"1.0.0"}
{"time":"2026-11-21T14:15:12.370521086Z","level":"INFO","msg":"plugin loaded","plugin":"nicht-gefunden","version":"1.0.0"}
{"time":"2026-11-21T14:15:13.754449461Z","level":"INFO","msg":"plugin loaded","plugin":"suche","version":"1.0.0"}
{"time":"2026-11-21T14:15:15.101226545Z","level":"INFO","msg":"plugin loaded","plugin":"jahreszahl","version":"1.0.0"}
```

#### amd64-t34-evidence-commit.txt

```text
cefe4ac
```

#### amd64-t34-pass1.status

```text
0
```
