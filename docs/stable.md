# What 1.0 will promise, and what it still lacks

0.1.5 is a beta, and `docs/beta.md` says what it promises. This page says what 1.0 will promise and, for each area, what is still missing, as milestone v0.1.5 leaves it. A promise on this page is not given until 1.0 ships: until then it is the owner's decision of what to promise, written down in advance, and nothing here holds a release to it. Each area has two parts. *What 1.0 promises* is written in the future tense, and what it states as present cites its source. *What is still missing* is a checklist, and each item names the requirement, the owner step, the record or the file it waits on. An item is checked only when a record shows it, and no item is checked yet.

## Template contract

### What 1.0 promises

A theme written for template contract version 1 will keep working in every 1.x release. Within 1.x the contract will only grow: fields will be added, and none will be removed or renamed, and the meaning of a field will not change. A break will come only with 2.0. `holzcloud template check` will report an incompatibility before an upload is accepted. The data contract that templates see is part of this promise: the structs in `internal/template/loader.go`, the fixtures `SampleData` and `MinimalData` in `internal/template/sample.go`, and the specification `internal/tmplspec/TEMPLATE-SPEC.md`, which the tests of `internal/template` tie to one another.

What stands today is smaller. `ContractVersion` in `internal/template/contract.go` is `1.0`, written as major.minor, and a theme may say which contract it was written for in `theme.json`, as `{"contract": "1.0"}` (`internal/tmplspec/TEMPLATE-SPEC.md` §2.4). `holzcloud template check` prints the contract it reads a theme against and the version the theme declares, adds a note when the two differ, and refuses nothing; the comment on `ContractVersion` says that the package draws no conclusion from a difference. None of the promise above is in force in 0.1.5: `docs/beta.md` says that the template contract may still change or go away from one release to the next.

### What is still missing

- [ ] `holzcloud template check` refusing, before an upload is accepted, a theme whose declared major version differs from `ContractVersion`: today it prints a note and refuses nothing (`internal/template/contract.go`, `cmd/holzcloud/cli_template.go`).
- [ ] The compatibility rules beyond declaring a version, that is, what a minor change may break, written into `internal/tmplspec/TEMPLATE-SPEC.md` and held by a test (the Future list of `REQUIREMENTS.md`, which STB-02 left open).
- [ ] A check that compares the fields of one release with those of the release before it, so that removing or renaming a field fails a build: the tests of `internal/template` tie `loader.go`, `sample.go` and the specification to one another within one tree and compare nothing across releases (`internal/template/sample_test.go`).
- [ ] A theme written for contract version 1.0 played unchanged against the first release after 1.0, which cannot exist before 1.0 does (STB-01).

## Data contract

### What 1.0 promises

Every 1.x release will take over a backup or a data directory of every earlier 1.x release and of 0.1.5. A backup is what `deploy/backup.sh` writes: a verified snapshot of the database (`holzcloud.sqlite`), `csrf.key`, the five directories `media`, `templates`, `plugins`, `sprachen` and `marke`, a `MANIFEST` and, from 0.1.5, a `SHA256SUMS` over every other file (`deploy/DEPLOY.md`, *Backups*). To take a backup over means that `deploy/restore.sh` restores it and the first start of the release leaves a site whose pages, media, users and plugins are those of the backup. To take a data directory over means that the release starts on it and migrates it. A backup or a data directory of 0.1.1 to 0.1.4 is not taken over by 1.x directly: it reaches 1.x through 0.1.5 (see *Migrations*).

What stands today is one run of each kind, in containers. A backup that v0.1.1 made at migration level 60, with its own `backup.sh` and binary, was restored with the `restore.sh` of 0.1.5 into 0.1.5 at level 69, and the restored site passed the 17 assertions of `run.sh check` (`deploy/MIGRATION-RESTORE-RUN.md`). A backup made before 0.1.5 has no checksum file, and `restore.sh` restores it with a notice that its files were not verified (`deploy/BACKUP-VERIFY-RUN.md`). Neither run is a promise of 0.1.5 for other releases: `docs/beta.md` lists what the two did not play.

### What is still missing

- [ ] A restore of a backup of 0.1.5 into the first 1.x release, played by the method of `deploy/MIGRATION-RESTORE-RUN.md`, which cannot be played before a 1.x release exists (STB-01).
- [ ] A restore on a real server, from a backup made there: every restore so far went into a container (SRV-02, OS-13, `deploy/NATIVE-AMD64-RUN.md`).
- [ ] A database and a media tree of production size, restored and started: the restored database of the record is 667648 bytes (`deploy/MIGRATION-RESTORE-RUN.md`, *Not tested*).
- [ ] An order, form messages and a second factor carried across a restore: the seeded source of the record had none of them (`deploy/MIGRATION-RESTORE-RUN.md`, *Not tested*).
- [ ] A backup that the timer made at night, restored across a migration: the backup of the record was made by the document's manual line (`deploy/MIGRATION-RESTORE-RUN.md`, *Not tested*).

## Migrations

### What 1.0 promises

Migrations will run forward only. A release may add migrations, and none will be undone: the goose `Down` section that each migration file carries is not a rollback route, and `holzcloud migrate` offers `status` and `up` and nothing else (`deploy/DEPLOY.md`, *Rolling back*). The releases 0.1.1 to 0.1.4 will reach 1.x through 0.1.5 as an intermediate step: a data directory or a backup of one of them is brought to 0.1.5 first and then to 1.x. A binary of 1.x will refuse a database that a newer binary migrated, name both migration levels and the way back, and serve nothing.

What stands today is the refusal and one crossing. Since 0.1.5 the binary reads the migration level before it does anything else with the database, and it refuses one that a newer binary migrated; a binary released before 0.1.5 does not (`deploy/DEPLOY.md`, *Rolling back*). The refusal was played once, with a build of 0.1.5 as the previous binary, five times in a row, with the levels 70 and 69 named (`deploy/ROLLBACK-RUN.md`). The crossing of the nine migrations 61 to 69 was played once, from a backup of v0.1.1 (`deploy/MIGRATION-RESTORE-RUN.md`).

### What is still missing

- [ ] The step from 0.1.2, 0.1.3 and 0.1.4 to 0.1.5, played as the source of an update and of a restore: only v0.1.1 was a source (`deploy/MIGRATION-RESTORE-RUN.md`, *Not tested*).
- [ ] The migrations timed on a database of production size: the nine of the record ran on 667648 bytes, and the two of the earlier update ran on test data in milliseconds (`deploy/MIGRATION-RESTORE-RUN.md`, `deploy/DEPLOY-RUN.md`).
- [ ] The step from 0.1.5 to the first 1.x release, played once that release exists, with the data directory and with a backup (STB-01).

## Upgrade and rollback

### What 1.0 promises

An update will go from any 1.x release directly to any newer 1.x release, and skipping releases will be allowed. The operator replaces the binary by the steps of *Updating* in `deploy/DEPLOY.md`, and the service writes a verified snapshot of the database before it applies a migration, naming it in the journal.

Going back will be done only with the pre-upgrade snapshot together with the previous binary, by *Rolling back* in `deploy/DEPLOY.md`. The snapshot holds the database only, so media, templates, plugins, language files and the logo come back from a backup and not from the snapshot, and everything written since the update is lost. There will be no promise that an older binary reads a database that a newer binary migrated; the schema guard refuses it instead.

What stands today is two runs in containers, both with throwaway builds of the tree and not with releases. Six updates in a row on one data directory each wrote one snapshot, and after every update the newest five were kept (`deploy/UPDATE-CHAIN-RUN.md`). An update to a build whose second migration failed stopped with the failure named, its restarts kept the one snapshot from before the update, and *Rolling back* block 1 put the previous binary and that snapshot back on the state before the update (`deploy/ROLLBACK-RUN.md`).

### What is still missing

- [ ] An update between released 1.x versions that skips a release, and a chain made of released binaries: each of the six updates went to the next build, and the builds were not releases (`deploy/UPDATE-CHAIN-RUN.md`, *Not tested*).
- [ ] A rollback after a migration that fails for another reason, such as a full disk or a process killed in the middle of a migration, and after an update whose first migration fails (`deploy/ROLLBACK-RUN.md`, *Not tested*).
- [ ] *Rolling back* block 2, for an update that wrote no snapshot, and the undoing of a rollback with `holzcloud.sqlite.after-update` (`deploy/ROLLBACK-RUN.md`, *Not tested*).
- [ ] A reboot of a real server with the service and the backup timer enabled, after an update: a restart of a container stood in for it (SRV-03, OS-14, `deploy/NATIVE-AMD64-RUN.md`).

## Supported platforms

### What 1.0 promises

1.x will be supported on linux/amd64, on Debian 12 and on Ubuntu 24.04, with systemd and the unit that ships (`deploy/holzcloud.service`), behind Caddy (`deploy/Caddyfile.example`). arm64 will not be promised.

What stands today is narrower. The binary is built for linux/amd64, and the container image is built for linux/amd64 only (`deploy/DEPLOY.md`, *Build*, *The container image*). Every record of 0.1.1 and of 0.1.5 ran on Debian 12 with systemd 252 in a container, amd64 under emulation on an arm64 host, except where a row of `docs/beta.md` says otherwise, and the install was played on Debian 12 only (`deploy/DEPLOY-RUN.md`). One native amd64 run exists, on a GitHub-hosted runner (`deploy/NATIVE-AMD64-RUN.md`). Ubuntu 24.04 was not installed, and the proxy was not played.

### What is still missing

- [ ] A native amd64 run on a real server under the host's own systemd, with the seccomp-based directives of the unit enforced there: SRV-01 ran once on a GitHub-hosted amd64 runner, in a privileged Debian 12 container with systemd 252, and exercised the plugin runtime under `MemoryDenyWriteExecute=true` and the backup timings (SRV-01, OS-5, `deploy/NATIVE-AMD64-RUN.md`).
- [ ] A reboot with the service and the timer units enabled, on a real server: a restart of a container stood in for it (SRV-03, OS-14, `deploy/NATIVE-AMD64-RUN.md`).
- [ ] An install on Ubuntu 24.04, whose systemd is 255, and on the image of a provider: the install text was played on Debian 12 only, and only systemd 252 and 257 were measured (`deploy/DEPLOY-RUN.md`, `deploy/SECURITY-SCORE.md`).
- [ ] Certificate issuance, DNS, ACME, the proxy reaching the service and a browser over https: no record plays them (`deploy/Caddyfile.example`, `docs/beta.md`).
- [ ] A data directory on its own mount, so that `RequiresMountsFor` in `deploy/holzcloud.service` is exercised: every record had one filesystem.

## Outside the 1.0 promise

The plugin ABI and the SDK, which are the WASM plugins, and the MCP tools for AI keys are not yet promised by 1.0.

The plugin ABI changed as late as 0.1.5: the hook `privacy` (`HookPrivacy` in `internal/plugin/manifest.go`), the field `RetainUntil` of the privacy call in `sdk/plugin.go` and the column `media.owner_plugin` (migration 00069) were added while `ABIVersion` stayed 1, so a host of 0.1.1 refuses a plugin archive of 0.1.5 that carries the same number (`deploy/MIGRATION-RESTORE-RUN.md`).

No record of 0.1.5 holds the MCP tools (`docs/ai-access.md`) still, and their set and arguments may change in any release until a later release promises them. A later release may promise the plugin ABI, the SDK and the MCP tools separately.

Also not promised are arm64 (see *Supported platforms*) and that an older binary reads a database that a newer binary migrated (see *Upgrade and rollback*).

## Not yet decided

Which 1.x releases receive security fixes, and for how long, is not decided. Nothing is promised about it. This page will change when the owner decides.
