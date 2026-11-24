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
