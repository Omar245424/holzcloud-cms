# Two-factor enrolment, a shop order and a form submission through a real browser: results

This document is the permanent record of UAT-01, UAT-02 and UAT-03 of phase 34. It records three flows played through a real browser, headless Chromium driven by Playwright, against a dev server of this release on the Mac with a data directory of its own: a second administrator created in the admin and forced through the two-step enrolment at the first sign-in, with ten recovery codes, one of them accepted once and refused the second time (UAT-01); two shop orders placed by visitors, one with invoice and one with prepayment, each confirmed and then listed in the admin with its payment method and state (UAT-02); a contact message that arrives in the admin, and a submission that fills the hidden field, gets the same answer and lands in quarantine with its reason (UAT-03).

The record rests on one evidence set: a single `browser-verify.sh run`, made once against the committed text of every bound file (see *How this was run*). Its raw output is the Appendix, pasted unchanged, and `bash deploy/test/browser-verify.sh record-check` proves every block against its scratch file byte for byte. Every value in this document is a line of that set and cites its file in backticks as `DIRECTORY/FILE`, with `DIRECTORY` the scratch directory `mac-t34-uat`; the run files beside it are named by their path below `deploy/test/out`.

The Playwright script that played the flows is not part of the repository: the stack allows no JavaScript in the tree beyond `cmd/holzcloud/assets/htmx.min.js`, and the script is JavaScript. `run` copied it byte for byte into the set, and its full text is the appendix block `mac-t34-uat/script.txt`, proven byte for byte like every other block; its sha256 is the `script sha256:` line of `mac-t34-uat/versions.txt`. No secret was printed: the script makes its passwords at run time and prints counts, shapes and names, never a password, a TOTP secret, a code, a recovery code, a cookie value or an order token, and the driver's leak scan checked its output in `run` and checks this whole record outside the script block in `record-check`. The screenshots the script took stay in the scratch directory and are not pasted.

The record says what was done, how many times, and what was seen. It played each flow once, in one browser, with JavaScript on, against a dev server on loopback. What that cannot show is a named deviation or a line of *Not tested*, never a gap in the table.

## Summary

| Requirement | Result | Platform | Date, commit |
|-------------|--------|----------|--------------|
| UAT-01 | PASS: `UAT-01 PASS` stands in `mac-t34-uat/summary.txt`. The first administrator, made on the setup form, was sent to the enrolment and enrolled (`ADMIN1-ENROLLED: yes`); creating a second administrator asked for the password again (`PASSWORD-QUESTION: yes`, `ADMIN2-CREATED: yes`); at her first sign-in she was forced to the enrolment (`FORCED-REDIRECT: /admin/2fa/einrichten`, `SECRET-SHOWN: yes`), received ten recovery codes (`RECOVERY-CODES: 10`, `RECOVERY-CODE-SHAPE: 5-5 lower-case alphanumeric`) and saw the factor on her account (`KONTO-FACTOR: on`, `KONTO-CODES-LEFT: 10`); the next sign-in asked for a code and let her in (`RELOGIN-ASKS-CODE: yes`, `RELOGIN-OK: yes`), and a recovery code was accepted once and refused the second time (`RECOVERY-CODE-ONCE: yes (accepted once, refused the second time)`, `KONTO-CODES-LEFT-AFTER-ONE: 9`), all in `mac-t34-uat/uat.txt`. | macOS on arm64, dev server built from the evidence commit, headless Chromium (Playwright 1.59.1), German locale, JavaScript on (`mac-t34-uat/versions.txt`) | 2026-11-21, 6e26d6c |
| UAT-02 | PASS: `UAT-02 PASS` stands in `mac-t34-uat/summary.txt`. The product page answered without setting a cookie (`PRODUCT-PAGE: 200`, `COOKIE-BEFORE-ACTION: none`), the basket cookie came with the first action (`BASKET-COOKIE: hc_warenkorb`), the checkout offered invoice and prepayment and no Payrexx (`PAYMENT-OFFERED: invoice, prepay`, `PAYREXX-OFFERED: no`); two orders were confirmed on a page of their own (`ORDER-1-URL-SHAPE` and `ORDER-2-URL-SHAPE: matches /bestellung/[0-9a-f]{32}`, `ORDER-1-NUMBER: 2026-0001`, `ORDER-2-NUMBER: 2026-0002`) and both are listed in the admin with their payment method and state (`ORDER-1-ADMIN-LISTED: yes`, `ORDER-2-ADMIN-LISTED: yes`, `ORDER-1-ADMIN-PAYMENT: Rechnung`, `ORDER-2-ADMIN-PAYMENT: Vorauskasse`, `ORDER-1-ADMIN-STATUS: Neu`, `ORDER-2-ADMIN-STATUS: Neu`), all in `mac-t34-uat/uat.txt`. | same | 2026-11-21, 6e26d6c |
| UAT-03 | PASS: `UAT-03 PASS` stands in `mac-t34-uat/summary.txt`. A visitor sent a contact message after a wait of 4.5 seconds (`FORM-PAGE: 200`, `FORM-WAIT-MS: 4502`) and got the answer `formular=gesendet, Danke, die Nachricht ist angekommen. Wir melden uns.`, and the message is in the admin (`MESSAGE-IN-ADMIN: yes`); a second visitor filled the hidden field (`SPAM-FIELD: website`) and got the same answer (`SPAM-ANSWER-SAME: yes`), and the submission is in quarantine with its reason and not among the messages (`SPAM-IN-QUARANTINE: yes`, `QUARANTINE-REASON: Das versteckte Feld war ausgefüllt`, `SPAM-NOT-IN-MESSAGES: yes`); every visitor request went to one origin (`SAME-ORIGIN: yes (42 requests, 1 origin)`), all in `mac-t34-uat/uat.txt`. | same | 2026-11-21, 6e26d6c |

Evidence commit: 6e26d6c

The binary under test is the one built from `NEW_REF` (pinned by `run` to the short sha of HEAD, `6e26d6c`, read from `mac-t34-evidence-commit.txt`) with `CGO_ENABLED=0 go build -trimpath` and the version stamped by `git describe`: `BINARY: holzcloud v0.1.4-266-g6e26d6c (6e26d6c)` in `mac-t34-uat/summary.txt`, the same as `binary:`, `NEW_REF:` and `evidence commit:` in `mac-t34-uat/versions.txt`, and the `version` and `commit` of the server's first log line in `mac-t34-uat/server-excerpt.txt`. The same file names the tools: `go version go1.27.4 darwin/arm64`, Node `v22.23.5`, Playwright `1.59.1` and `chromium 147.0.7727.15 (chromium-1217)`, on `Darwin 27.2.0 arm64, macOS 27.2.1`. The summary carries no `REHEARSAL` line. The commit `6e26d6c` holds the code of `cefe4ac` and, on top of it, only the records of the container sets and `docs/beta.md`; the binary's source is the same at both.

## How this was run

`deploy/test/browser-verify.sh` builds this tree and `tools/totpcode` into a scratch directory outside the repository, starts the binary there on a data directory of its own, on `127.0.0.1` only, with mail pointed at an unreachable port so that it is queued and never sent, and runs the Playwright script against it. `run` refuses while a bound file or the binary's source differs from HEAD or is untracked, or while the tracked JavaScript is anything but the htmx file, and pins `NEW_REF` to the short sha of HEAD. The script prints `KEY: value` lines; the driver computes the verdict of each requirement from those lines (`mac-t34-uat/uat.txt`), never from a trace, and a flow whose keys are missing reads FAIL or NOT-RUN.

The set was made once, on the clean committed tree of `6e26d6c` without `REHEARSAL`, started 2026-11-21T14:37:32Z and ended 14:38:48Z on the host clock (1 minute 16 seconds), exit status 0 (`deploy/test/out/mac-t34-pass1.started`, `deploy/test/out/mac-t34-pass1.finished` and `mac-t34-pass1.status`). Before it, the worktree held no earlier output of the driver and nothing answered on the port; after it the driver had stopped the server it started, and nothing answered on the port. No second set was made against `6e26d6c`, and no other set is cited. One earlier set was superseded. It was made against `bef9191`, whose binary is built from the same source as `00854be`, and it did not pass: the contact-form plugin, switched on in the admin, did not start, because the plugin was loaded with the context of the admin request that switched it on and closed with it when that request was dropped. No record was written from that set. The commit `cefe4ac` fixed the load in `internal/plugin/runtime.go`, so the binary changed after it, and the set was made again. This record was written from this set, and nothing of the earlier set is quoted. The command was:

```text
UAT_DIR=$SCRATCH/t34-uat/evidence2 UAT_JS=$SCRATCH/t34-uat/uat.js PW=/Users/holz/.npm/_npx/e41f203b7505f1fb/node_modules/playwright bash deploy/test/browser-verify.sh run
```

`$SCRATCH` stands for the scratch directory of the session that made the set, outside the repository; `UAT_DIR` was a fresh directory in it, and `UAT_JS` the script whose text is the block `mac-t34-uat/script.txt`. The knob not set on that line had its default: `PORT=18793`. The log of the run is `deploy/test/out/mac-t34-pass1.log`, which is not pasted; its first line reads `browser-verify: NEW_REF pinned to 6e26d6c, UAT_DIR …, port 18793`, with the scratch path left out here. The server log, the script's stderr (empty) and the screenshots stay in `UAT_DIR` and are not pasted; `mac-t34-uat/server-excerpt.txt` holds the lines of the server log a record cites.

No bound file changed after the evidence commit. The bound files for this set are `deploy/test/browser-verify.sh`, `tools/totpcode/main.go` and `plugins/kontaktformular/kontaktformular.zip` (the archive the script uploads), together with the binary built from `cmd`, `internal`, `go.mod`, `go.sum`, `changelog.go` and the `CHANGELOG.md` it embeds. This command exited 0 when the record was written:

```text
git diff --quiet 6e26d6c HEAD -- deploy/test/browser-verify.sh tools/totpcode/main.go plugins/kontaktformular/kontaktformular.zip cmd internal go.mod go.sum changelog.go CHANGELOG.md
```

Before the run, `bash deploy/test/browser-verify.sh selftest` ended `SELFTEST OK (137 cases)`; its output is not part of the set.

## Setup

Everything the three flows needed was made through the admin screens, in the browser, by the first administrator: the website `Holzladen T34` on the form for a new website, with a seller address; its domain `laden.t34.test` on the website's page; the shop switched on with currency CHF, an order address and payment details; the product `Holzbrett`, published, with a price of 25.00 and a stock of 10; the plugin `kontaktformular` uploaded from `plugins/kontaktformular/kontaktformular.zip`, bound to the website and switched on; and the published page `Kontakt` holding `[[formular]]`. No row was written by direct SQL (`FIXTURE-SQL: none (website, domain, shop, product, plugin and page made through the admin screens)` in `mac-t34-uat/uat.txt`). The script's notes in `deploy/test/out/mac-t34-pass1.log` name each step: the website created on `/admin/websites/new`, the domain listed, the product listed, the plugin acting on the website, and the page saved. The server log has the plugin loaded (`"msg":"plugin loaded","plugin":"kontaktformular","version":"1.6.0"` in `mac-t34-uat/server-excerpt.txt`), and the screenshot of the plugin list, which is not pasted, shows it running.

The two administrator passwords were made at run time from 18 random bytes each, held in the script's memory and never printed; the TOTP secrets, codes and recovery codes the script read from the pages were held the same way, and every such value is replaced by `***` in anything the script prints. The visitors' names, addresses and messages are fixed text of the script.

## UAT-01: a new administrator is led from the first sign-in to the recovery codes

**What was done.** The first administrator was made on `/admin/setup`; the setup landed on `/admin/2fa/einrichten`, where the script read the secret from the page, computed a code with `tools/totpcode` and confirmed it (`ADMIN1-ENROLLED: yes`; `"msg":"two factor enabled","user_id":1` in `mac-t34-uat/server-excerpt.txt`). She then filled the form on `/admin/users/new` for a second administrator; the admin asked for her password first on `/admin/bestaetigen` (`PASSWORD-QUESTION: yes`), and after it the account was made and is on the user list (`ADMIN2-CREATED: yes`). In a second browser context the second administrator signed in for the first time and asked for `/admin/`; she was sent to `/admin/2fa/einrichten` instead (`FORCED-REDIRECT: /admin/2fa/einrichten`), where the page showed a secret of at least sixteen characters (`SECRET-SHOWN: yes`). She confirmed a code and was shown ten recovery codes of one shape (`RECOVERY-CODES: 10`, `RECOVERY-CODE-SHAPE: 5-5 lower-case alphanumeric`; `"user_id":2` in the second `two factor enabled` line). The link of the codes page led to `/admin/konto`, which showed the factor on and ten codes left (`KONTO-FACTOR: on`, `KONTO-CODES-LEFT: 10`). She signed out and in again: the sign-in asked for a code on `/admin/2fa` and, with a fresh code, let her in (`RELOGIN-ASKS-CODE: yes`, `RELOGIN-OK: yes`). Then, signed out once more, she signed in with the first recovery code in place of a code, which let her in and left nine codes (`KONTO-CODES-LEFT-AFTER-ONE: 9`), and the same recovery code at the next sign-in was refused (`RECOVERY-CODE-ONCE: yes (accepted once, refused the second time)`), all in `mac-t34-uat/uat.txt`.

*What it shows.* An administrator created in the admin, behind the password question, was forced at her first sign-in through the two-step enrolment in a real browser, was shown ten recovery codes, saw the factor on her account, was asked for a code at the next sign-in, and could use a recovery code once and not twice.

## UAT-02: an order through basket and checkout, with invoice and with prepayment

**What was done.** A visitor, in a browser context of her own, opened `http://laden.t34.test:18793/shop/holzbrett`; the page answered and had set no cookie (`PRODUCT-PAGE: 200`, `COOKIE-BEFORE-ACTION: none`). She put the product in the basket, and only then a cookie was set, the basket's (`BASKET-COOKIE: hc_warenkorb`). On `/kasse` the checkout offered invoice and prepayment and no Payrexx (`PAYMENT-OFFERED: invoice, prepay`, `PAYREXX-OFFERED: no`). She filled the address, chose invoice, ticked the confirmation and sent the order; the browser landed on a confirmation page under `/bestellung/` followed by the order token (`ORDER-1-URL-SHAPE: matches /bestellung/[0-9a-f]{32}`), headed with the order number `2026-0001` (`ORDER-1-NUMBER`). A second visitor did the same with prepayment and got `2026-0002` (`ORDER-2-URL-SHAPE`, `ORDER-2-NUMBER`). The token itself was never printed; only the shape of the address was checked. In the admin, the first administrator's order list holds both numbers (`ORDER-1-ADMIN-LISTED: yes`, `ORDER-2-ADMIN-LISTED: yes`), and each order's page names its payment method and its state (`ORDER-1-ADMIN-PAYMENT: Rechnung`, `ORDER-2-ADMIN-PAYMENT: Vorauskasse`, `ORDER-1-ADMIN-STATUS: Neu`, `ORDER-2-ADMIN-STATUS: Neu`), all in `mac-t34-uat/uat.txt`. Each order queued two mails, one to the customer and one to the operator; all four met the unreachable mail port and stayed queued (the four `mail not sent` lines with `order_customer` and `order_operator` in `mac-t34-uat/server-excerpt.txt`).

*What it shows.* Two visitors ordered through the product page, the basket and the checkout of a real browser, one with invoice and one with prepayment, each landed on a confirmation page with an order number, and both orders are in the admin with their payment method and the state `Neu`; no cookie was set before the visitor acted, and Payrexx was not offered.

Payrexx was not played: it needs the owner's account with a third party, and the shop of this set had no Payrexx configured.

## UAT-03: a form message in the admin, a honeypot submission in quarantine

**What was done.** A visitor opened `http://laden.t34.test:18793/kontakt` (`FORM-PAGE: 200`), filled name, address, subject and message, and sent the form 4502 milliseconds after the page was drawn (`FORM-WAIT-MS: 4502`). The wait is longer than the plugin's minimum of three seconds by more than a second, because the plugin's time trap compares whole seconds. The answer was `formular=gesendet, Danke, die Nachricht ist angekommen. Wir melden uns.` (`MESSAGE-ANSWER`). A second visitor filled the same form and, in addition, the hidden field `website` (`SPAM-FIELD: website`), and sent it after the same wait (`SPAM-WAIT-MS: 4501`); she got exactly the same answer (`SPAM-ANSWER-SAME: yes`), and the server logged `"message":"honeypot triggered"` (`mac-t34-uat/server-excerpt.txt`). In the admin, the plugin's message list of the website holds the first message's subject and not the second's (`MESSAGE-IN-ADMIN: yes`, `SPAM-NOT-IN-MESSAGES: yes`), and its quarantine view (`?ansicht=abgewiesen`) holds the second submission with the reason `Das versteckte Feld war ausgefüllt` (`SPAM-IN-QUARANTINE: yes`, `QUARANTINE-REASON`). Every request of the four visitors' browser contexts, the two orders included, went to the site's own origin (`SAME-ORIGIN: yes (42 requests, 1 origin)`), all in `mac-t34-uat/uat.txt`.

*What it shows.* A message sent through the contact form of a real browser arrived in the admin, and a submission that filled the hidden field got the same answer as a real one, was not among the messages and was listed in quarantine with the reason "Das versteckte Feld war ausgefüllt".

## Deviations named

Everything here precedes every result: it is what this set did differently from an operator's server and an operator's visitors.

| Name | Why | What a real server does instead |
|------|-----|---------------------------------|
| A dev server on the Mac | The binary of the evidence commit ran on macOS on arm64, started by the driver in the background with `HOLZCLOUD_LISTEN=127.0.0.1` and plain http on port 18793, without systemd, without the unit's settings and without `HOLZCLOUD_SECURE` | linux/amd64 under the shipped `holzcloud.service`, behind Caddy over https, with `HOLZCLOUD_SECURE=true` |
| Headless Chromium | One browser, Chromium 147 driven by Playwright 1.59.1, without a window, locale `de-DE`, a viewport of 1280 by 900, JavaScript on (`JS: on` in `mac-t34-uat/uat.txt`) | Whatever browsers the operator and the visitors use |
| Host names by the browser's resolver rule | The public site's host `laden.t34.test` was mapped to `127.0.0.1` by Chromium's `--host-resolver-rules`, so the domain resolution through DNS was not exercised | A domain in DNS pointing at the server |
| Mail queued, never sent | `HOLZCLOUD_SMTP_HOST=127.0.0.1` with `HOLZCLOUD_SMTP_PORT=2`, where nothing listens; the four order mails were queued and failed their first attempt (`mac-t34-uat/server-excerpt.txt`) | A mail server that accepts the mails |
| Codes from `tools/totpcode` | The second factor's codes were computed by `tools/totpcode` from the secret the script read off the enrolment page; no authenticator app was used, and the QR code was not scanned | An authenticator app on the administrator's phone |
| One script, one run | The flows were played by a script, with the waits the script names (a new code window for each sign-in, four and a half seconds before each form submission) | People, at their own pace |

## Findings

Each finding has the class of the earlier records' rule: (1) a bug of the harness, fixed in the driver; (2) a limit of the setup, a named deviation or a *Not tested* line; (3) a measured property of a shipped file that matches what it and its document say, recorded and not edited; (4) a defect, a shipped file doing less than it or its document says, with a fix in a commit of its own. No finding of this set is class 1 or class 4, and no shipped file was edited after the evidence commit. The defect that the superseded set at `bef9191` found, a plugin load ended by a dropped request, was fixed in `cefe4ac` before this set (*How this was run*), and F-1 is the same step observed on the fixed code.

**F-1 (class 3): The contact-form plugin, switched on in the admin, loaded and ran.** The step that switched it on is followed in the server log by `"msg":"plugin loaded","plugin":"kontaktformular","version":"1.6.0"`, and no error line stands in the excerpt (`mac-t34-uat/server-excerpt.txt`); the form was on the page and the flow of UAT-03 passed.

**F-2 (class 3): The shop set no cookie before the visitor acted, and the basket cookie came with the first action.** `COOKIE-BEFORE-ACTION: none` and `BASKET-COOKIE: hc_warenkorb` (`mac-t34-uat/uat.txt`).

**F-3 (class 3): The order mails were queued and their failed first attempt was logged, one line each, with the kind of mail.** Four `WARN` lines `mail not sent` with `order_customer` and `order_operator` and `attempt` 1 (`mac-t34-uat/server-excerpt.txt`). They were never delivered, by design of the set (*Deviations named*).

**F-4 (class 3): A fresh data directory got a pre-upgrade snapshot at level 0 before its first migrations.** The first start wrote `"msg":"pre-upgrade snapshot written"` with `"level":0` and then applied the migrations 1 to 69 (`migrations applied at this start: 69 (versions 1 to 69)` in `mac-t34-uat/server-excerpt.txt`). `deploy/DEPLOY.md` says the snapshot is written when migrations are pending, which on a new installation they are, and that it holds the schema before them; here that is the empty database.

**F-5 (class 2): The database and the snapshot were not of mode `600` on the Mac.** The dev server ran with the shell's umask, not the unit's `UMask=0077`, so in `UAT_DIR`, which is not pasted, the database and the snapshot had mode `644`, inside a data directory of mode `700`. `deploy/DEPLOY.md` says mode `600` for the server under the shipped unit, which this set did not run (*Deviations named*).

## Measurements

The figures the harness wrote. They are figures of one run on one Mac; none is a promise of speed.

| Figure | Value | Where |
|--------|-------|-------|
| The whole run, build, server, the three flows and the shutdown | 14:37:32Z to 14:38:48Z, 1 minute 16 seconds | `deploy/test/out/mac-t34-pass1.started`, `deploy/test/out/mac-t34-pass1.finished` |
| The server's start, from the first log line to listening, with the 69 migrations of a new database | 15:37:36.513 to 15:37:37.423 local time | the `starting` and `server starting` lines of `mac-t34-uat/server-excerpt.txt` |
| The wait before the legitimate and before the hidden-field submission | 4502 ms and 4501 ms | `FORM-WAIT-MS`, `SPAM-WAIT-MS` in `mac-t34-uat/uat.txt` |
| The visitors' requests | 42, to one origin | `SAME-ORIGIN` in `mac-t34-uat/uat.txt` |
| The server log | 89 lines | `server log lines` in `mac-t34-uat/server-excerpt.txt` |

## Not tested

- A server behind Caddy with https and `HOLZCLOUD_SECURE=true` was not played; the server was a dev server on the Mac, on loopback over plain http.
- JavaScript off was not played. The flows are server-rendered forms, and they were played with JavaScript on.
- Payrexx was not played; it needs the owner's account with a third party.
- A real mail delivery was not played: the order mails were queued to an unreachable port and never sent.
- Other browsers were not played; only headless Chromium 147.
- The quarantine reason for a form sent too fast was not played; the spam case used the hidden field.
- Single sign-on administrators were not played. They are not forced to a second factor by this program: their second factor is settled on the identity provider.
- An authenticator app and the QR code were not played; the codes came from `tools/totpcode`.
- A restore through a browser was not played; this record does not answer the limit of `docs/beta.md` that begins "A restore through a browser was not seen".
- An order's later states (paid, shipped, cancelled) and the stock after the orders were not read.

## Owner steps

None. Every step of this record ran on the Mac, by the owner's decision; nothing here has to be done on the real server before the record counts.

## Hand-over to docs/beta.md

The sentences of this record that the page may quote, named by place and not copied here, so that each stays in exactly one place: UAT-01, the sentence after *What it shows.* in *UAT-01*; UAT-02, the sentence after *What it shows.* in *UAT-02*; UAT-03, the sentence after *What it shows.* in *UAT-03*; the setup, the sentence of *Setup* that begins "No row was written by direct SQL"; the limits, the bullets of *Not tested*.

The limit of the page this record answers, named by its first words: "The compulsory two-factor enrolment and the shop and form flows through a browser were not played" (UAT-01, UAT-02 and UAT-03 above). Its clause "orders were written straight into the database" no longer holds for this set: the two orders were placed by visitors through the checkout, and no row was written by direct SQL (*Setup*).

Nothing in this record is a promise until the page cites it byte for byte. A change to a bound file after the evidence commit means the set is made again before the record is cited for the changed text. `docs/beta.md` was not edited by this plan.

## Appendix: raw output

Each block is the scratch file of the evidence set, pasted unchanged by `bash deploy/test/browser-verify.sh appendix`, and `bash deploy/test/browser-verify.sh record-check` compares every block with its scratch file byte for byte. The files are in `deploy/test/out` (gitignored); the run log, the started and finished stamps and everything in `UAT_DIR` (the server log, the script's stderr, the screenshots) are not pasted.

#### mac-t34-uat/summary.txt

```text
---- summary (mac, dev server, Chromium) ----
UAT-01 PASS
UAT-02 PASS
UAT-03 PASS
BINARY: holzcloud v0.1.4-266-g6e26d6c (6e26d6c)
```

#### mac-t34-uat/versions.txt

```text
host:                  Darwin 27.2.0 arm64, macOS 27.2.1
go:                    go version go1.27.4 darwin/arm64
node:                  v22.23.5
playwright:            1.59.1
browser:               chromium 147.0.7727.15 (chromium-1217)
binary:                holzcloud v0.1.4-266-g6e26d6c (6e26d6c)
NEW_REF:               6e26d6c
evidence commit:       6e26d6c
script sha256:         9e8f11d6cfbc2b0a53d678e83ec263dd0cbc781f5b5106ae042b9dc3c2de0e69
```

#### mac-t34-uat/uat.txt

```text
ADMIN1-ENROLLED: yes
PASSWORD-QUESTION: yes
ADMIN2-CREATED: yes
FORCED-REDIRECT: /admin/2fa/einrichten
SECRET-SHOWN: yes
RECOVERY-CODES: 10
RECOVERY-CODE-SHAPE: 5-5 lower-case alphanumeric
KONTO-FACTOR: on
KONTO-CODES-LEFT: 10
RELOGIN-ASKS-CODE: yes
RELOGIN-OK: yes
KONTO-CODES-LEFT-AFTER-ONE: 9
RECOVERY-CODE-ONCE: yes (accepted once, refused the second time)
FIXTURE-SQL: none (website, domain, shop, product, plugin and page made through the admin screens)
PRODUCT-PAGE: 200
COOKIE-BEFORE-ACTION: none
BASKET-COOKIE: hc_warenkorb
PAYMENT-OFFERED: invoice, prepay
PAYREXX-OFFERED: no
ORDER-1-URL-SHAPE: matches /bestellung/[0-9a-f]{32}
ORDER-1-NUMBER: 2026-0001
ORDER-2-URL-SHAPE: matches /bestellung/[0-9a-f]{32}
ORDER-2-NUMBER: 2026-0002
ORDER-1-ADMIN-LISTED: yes
ORDER-1-ADMIN-PAYMENT: Rechnung
ORDER-1-ADMIN-STATUS: Neu
ORDER-2-ADMIN-LISTED: yes
ORDER-2-ADMIN-PAYMENT: Vorauskasse
ORDER-2-ADMIN-STATUS: Neu
FORM-PAGE: 200
FORM-WAIT-MS: 4502
MESSAGE-ANSWER: formular=gesendet, Danke, die Nachricht ist angekommen. Wir melden uns.
SPAM-FIELD: website
SPAM-WAIT-MS: 4501
SPAM-ANSWER-SAME: yes
MESSAGE-IN-ADMIN: yes
SPAM-NOT-IN-MESSAGES: yes
SPAM-IN-QUARANTINE: yes
QUARANTINE-REASON: Das versteckte Feld war ausgefüllt
SAME-ORIGIN: yes (42 requests, 1 origin)
JS: on
```

#### mac-t34-uat/script.txt

```text
// Phase 34 (plan 34-03): UAT-01, UAT-02 and UAT-03 in a real browser.
//
// Run by deploy/test/browser-verify.sh, which builds the binary and tools/totpcode
// into S, starts the dev server on 127.0.0.1:PORT with its own data directory and
// stops it again. This script only drives the browser and prints "KEY: value"
// lines (and short "note:" lines); the driver judges the keys. It never prints a
// password, a TOTP secret, a code, a recovery code, a cookie value or an order
// token: only counts, shapes and names.
//
// Not committed: the tree holds no JavaScript beyond htmx. Its text reaches the
// repository as the appendix of deploy/BROWSER-RUN.md.
'use strict';
const { chromium } = require(process.env.PW);
const { execFileSync } = require('child_process');
const crypto = require('crypto');
const fs = require('fs');

const S = process.env.S;
const PORT = process.env.PORT;
const ZIP = process.env.ZIP;
const BASE = 'http://127.0.0.1:' + PORT;
const SHOTS = S + '/shots';
const HOST = 'laden.t34.test';
const PUB = 'http://' + HOST + ':' + PORT;
const ARGS = ['--host-resolver-rules=MAP ' + HOST + ' 127.0.0.1'];

// ---------- output ----------
const secrets = [];
function hide(v) { if (v) secrets.push(String(v)); return v; }
function scrub(t) { let s = String(t); for (const x of secrets) s = s.split(x).join('***'); return s; }
function put(key, value) { console.log(key + ': ' + scrub(String(value).replace(/\s+/g, ' ').trim())); }
function note(text) { console.log('note: ' + scrub(String(text).replace(/\s+/g, ' ').trim()).slice(0, 200)); }
async function shot(page, name) { await page.waitForTimeout(300); await page.screenshot({ path: SHOTS + '/' + name, fullPage: true }); }

// ---------- helpers ----------
function password() { return hide(crypto.randomBytes(18).toString('base64url')); }
const lastWindow = new Map();
async function code(secret) {
  // a code is accepted once per account: every sign-in waits for a window of its own
  while (Math.floor(Date.now() / 30000) === lastWindow.get(secret)) await new Promise(r => setTimeout(r, 1000));
  // not in the last two seconds of a window, so the code is still valid on arrival
  while (Date.now() % 30000 > 28000) await new Promise(r => setTimeout(r, 500));
  lastWindow.set(secret, Math.floor(Date.now() / 30000));
  return hide(execFileSync(S + '/totpcode', [secret]).toString().trim());
}
async function settle(page) { await page.waitForLoadState('networkidle'); await page.waitForTimeout(200); }
function pathOf(page) { return new URL(page.url()).pathname; }
async function context(browser, opts = {}) {
  return browser.newContext({ viewport: { width: 1280, height: 900 }, locale: 'de-DE', ...opts });
}
async function signIn(page, email, pw) {
  await page.goto(BASE + '/admin/login');
  await page.locator('form[action="/admin/login"] #email').first().fill(email);
  await page.locator('form[action="/admin/login"] #password').first().fill(pw);
  await page.locator('form[action="/admin/login"] button[type=submit]').first().click();
  await settle(page);
}
async function signOut(page) {
  if (!pathOf(page).startsWith('/admin') || pathOf(page) === '/admin/login') await page.goto(BASE + '/admin/');
  await settle(page);
  await page.locator('summary.topbar-user').click();
  await page.locator('form[action="/admin/logout"] button[type=submit]').click();
  await settle(page);
}
async function enrol(page) {
  // on /admin/2fa/einrichten: the secret from the page (kept in memory), a code, the codes page
  const secret = hide(((await page.textContent('.totp-secret')) || '').replace(/\s+/g, ''));
  await page.fill('#code', await code(secret));
  await page.click('form[action="/admin/2fa/einrichten"] button[type=submit]');
  await settle(page);
  const codes = await page.locator('ul.recovery-codes li').allInnerTexts();
  codes.forEach(c => hide(c.trim()));
  return { secret, codes: codes.map(c => c.trim()) };
}
// the password question in front of the irreversible
async function confirmIfAsked(page, pw) {
  if (pathOf(page) !== '/admin/bestaetigen') return false;
  await page.fill('form[action="/admin/bestaetigen"] #password', pw);
  await page.click('form[action="/admin/bestaetigen"] button[type=submit]');
  await settle(page);
  return true;
}
// the number on the second tile of /admin/konto's security bar ("%d recovery codes left")
async function codesLeft(page) {
  const t = await page.locator('.security-bar .security-tile').nth(1).locator('strong').innerText().catch(() => '');
  return (t.match(/\d+/) || ['?'])[0];
}
function shape(c) {
  const m = /^([a-z0-9]+)-([a-z0-9]+)$/.exec(c);
  return m ? `${m[1].length}-${m[2].length} lower-case alphanumeric` : 'other';
}

(async () => {
  fs.mkdirSync(SHOTS, { recursive: true });
  const browser = await chromium.launch({ args: ARGS });
  fs.writeFileSync(S + '/browser.txt', 'chromium ' + browser.version() + ' (' + (chromium.executablePath().match(/chromium[a-z_-]*-\d+/) || ['?'])[0] + ')\n');
  note('browser chromium ' + browser.version() + ', locale de-DE, JavaScript on');

  // =====================================================================
  // UAT-01: a second administrator is forced through two-step enrolment
  const A1 = { email: 'admin1@example.test', pw: password() };
  const A2 = { email: 'admin2@example.test', pw: password() };
  const c1 = await context(browser);
  const p1 = await c1.newPage();
  await p1.goto(BASE + '/admin/setup');
  await p1.locator('form[action="/admin/setup"] #email').first().fill(A1.email);
  await p1.locator('form[action="/admin/setup"] #password').first().fill(A1.pw);
  await p1.locator('form[action="/admin/setup"] #password_confirm').first().fill(A1.pw);
  await p1.locator('form[action="/admin/setup"] button[type=submit]').first().click();
  await settle(p1);
  note('admin1 after setup lands on ' + pathOf(p1));
  if (pathOf(p1) !== '/admin/2fa/einrichten') await p1.goto(BASE + '/admin/2fa/einrichten');
  const e1 = await enrol(p1);
  put('ADMIN1-ENROLLED', e1.codes.length === 10 ? 'yes' : 'no');

  // admin2 on /admin/users/new, behind the password question
  async function fillUser(page) {
    await page.goto(BASE + '/admin/users/new');
    await settle(page);
    await page.fill('#name', 'Zweite Administratorin');
    await page.fill('#email', A2.email);
    await page.selectOption('#role', 'admin');
    await page.fill('form[action="/admin/users/new"] #password', A2.pw);
    await page.click('form[action="/admin/users/new"] button[type=submit]');
    await settle(page);
  }
  await fillUser(p1);
  const asked = await confirmIfAsked(p1, A1.pw);
  put('PASSWORD-QUESTION', asked ? 'yes' : 'no');
  if (asked) await fillUser(p1);
  await p1.goto(BASE + '/admin/users');
  await settle(p1);
  const listed = (await p1.innerText('main')).includes(A2.email);
  put('ADMIN2-CREATED', listed ? 'yes' : 'no');
  await shot(p1, 'uat01-users.png');

  // admin2 signs in for the first time
  const c2 = await context(browser);
  const p2 = await c2.newPage();
  await signIn(p2, A2.email, A2.pw);
  await p2.goto(BASE + '/admin/');
  await settle(p2);
  put('FORCED-REDIRECT', pathOf(p2));
  const secretShown = await p2.locator('.totp-secret').count() === 1 && ((await p2.textContent('.totp-secret')) || '').replace(/\s+/g, '').length >= 16;
  put('SECRET-SHOWN', secretShown ? 'yes' : 'no');
  await shot(p2, 'uat01-forced-setup.png');
  const e2 = await enrol(p2);
  put('RECOVERY-CODES', e2.codes.length);
  const shapes = [...new Set(e2.codes.map(shape))];
  put('RECOVERY-CODE-SHAPE', shapes.length === 1 ? shapes[0] : 'mixed: ' + shapes.join(', '));
  await p2.click('.form-actions a[href="/admin/konto"]'); // the link of the codes page, "I have written them down"
  await settle(p2);
  const factorOn = await p2.locator('form[action="/admin/2fa/codes"]').count() === 1 &&
    await p2.locator('.security-tile--ok').first().count() === 1;
  put('KONTO-FACTOR', factorOn ? 'on' : 'off');
  put('KONTO-CODES-LEFT', await codesLeft(p2));
  await shot(p2, 'uat01-konto.png');

  // sign out, sign in: the code question, a fresh code, in
  await signOut(p2);
  await signIn(p2, A2.email, A2.pw);
  const asksCode = pathOf(p2) === '/admin/2fa' && await p2.locator('form[action="/admin/2fa"] #code').count() === 1;
  put('RELOGIN-ASKS-CODE', asksCode ? 'yes' : 'no');
  await p2.fill('#code', await code(e2.secret));
  await p2.click('form[action="/admin/2fa"] button[type=submit]');
  await settle(p2);
  put('RELOGIN-OK', pathOf(p2).startsWith('/admin') && !['/admin/2fa', '/admin/login', '/admin/2fa/einrichten'].includes(pathOf(p2)) ? 'yes' : 'no');

  // a recovery code: accepted once, refused the second time
  async function withRecovery(rc) {
    await signOut(p2);
    await signIn(p2, A2.email, A2.pw);
    await p2.click('a[href="/admin/2fa?wiederherstellung=1"]');
    await settle(p2);
    await p2.fill('#code', rc);
    await p2.click('form[action^="/admin/2fa"] button[type=submit]');
    await settle(p2);
    return !pathOf(p2).startsWith('/admin/2fa') && pathOf(p2) !== '/admin/login';
  }
  const first = await withRecovery(e2.codes[0]);
  let second = null;
  if (first) {
    await p2.goto(BASE + '/admin/konto');
    await settle(p2);
    put('KONTO-CODES-LEFT-AFTER-ONE', await codesLeft(p2));
    second = await withRecovery(e2.codes[0]);
    await shot(p2, 'uat01-recovery-refused.png');
  }
  put('RECOVERY-CODE-ONCE', first && second === false ? 'yes (accepted once, refused the second time)'
    : `no (first ${first ? 'accepted' : 'refused'}, second ${second === null ? 'not tried' : second ? 'accepted' : 'refused'})`);
  await c2.close();

  // =====================================================================
  // Setup for UAT-02 and UAT-03, as admin1, through the admin screens only
  put('FIXTURE-SQL', 'none (website, domain, shop, product, plugin and page made through the admin screens)');
  // a website, with a seller address in CH
  await p1.goto(BASE + '/admin/websites/new');
  await settle(p1);
  await p1.fill('form[action="/admin/websites/new"] #name', 'Holzladen T34');
  await p1.click('form[action="/admin/websites/new"] button[type=submit]');
  await settle(p1);
  const W = (pathOf(p1).match(/^\/admin\/websites\/(\d+)$/) || [])[1];
  if (!W) throw new Error('no website was created (landed on ' + pathOf(p1) + ')');
  note('website ' + W + ' created on /admin/websites/new');
  const wf = `form[action="/admin/websites/${W}"]`;
  for (const [id, v] of [['street', 'Holzgasse 4'], ['postal_code', '3011'], ['city', 'Bern'], ['country', 'CH'], ['contact_email', 'laden@example.test']]) {
    if (await p1.locator(`${wf} #${id}`).count()) await p1.fill(`${wf} #${id}`, v);
  }
  if (await p1.locator(`${wf} input[name=active]`).count()) await p1.check(`${wf} input[name=active]`);
  await p1.click(`${wf} button.btn--primary[type=submit]`);
  await settle(p1);
  // its domain (the form posts through htmx and swaps the list)
  await p1.goto(`${BASE}/admin/websites/${W}`);
  await settle(p1);
  await p1.fill('form.domain-add-form #domain', HOST);
  await p1.check('form.domain-add-form input[name=is_primary]');
  await p1.click('form.domain-add-form button[type=submit]');
  await settle(p1);
  note('domain ' + HOST + ' listed: ' + ((await p1.innerText('#domain-list')).includes(HOST) ? 'yes' : 'no'));
  // the shop on
  await p1.goto(`${BASE}/admin/websites/${W}/shop`);
  await settle(p1);
  await p1.fill('#shop_base', 'shop');
  if (await p1.locator('#currency option[value=CHF]').count()) await p1.selectOption('#currency', 'CHF');
  if (await p1.locator('#order_email').count()) await p1.fill('#order_email', 'bestellungen@example.test');
  if (await p1.locator('#payment_details').count()) await p1.fill('#payment_details', 'Bitte auf das Konto der Holzladen T34 überweisen.');
  await p1.locator('main form button.btn--primary[type=submit]').first().click();
  await settle(p1);
  // one product, published, with a price and stock
  await p1.goto(`${BASE}/admin/websites/${W}/produkte/neu`);
  await settle(p1);
  await p1.fill('#title', 'Holzbrett');
  await p1.fill('#slug', 'holzbrett');
  await p1.fill('#price', '25.00');
  await p1.fill('#stock', '10');
  await p1.selectOption('#status', 'published');
  await p1.locator('main form button.btn--primary[type=submit]').first().click();
  await settle(p1);
  await p1.goto(`${BASE}/admin/websites/${W}/produkte`);
  await settle(p1);
  note('product list holds Holzbrett: ' + ((await p1.innerText('main')).includes('Holzbrett') ? 'yes' : 'no'));
  await shot(p1, 'setup-products.png');
  // the contact form plugin: uploaded, bound to the website, switched on
  await p1.goto(BASE + '/admin/plugins');
  await settle(p1);
  await p1.click('details.plug-install > summary');
  await p1.setInputFiles('form[action="/admin/plugins/upload"] #plugin', ZIP);
  await p1.click('form[action="/admin/plugins/upload"] button[type=submit]');
  await settle(p1);
  if (await confirmIfAsked(p1, A1.pw)) {
    note('the upload asked for the password again');
    await p1.goto(BASE + '/admin/plugins'); await settle(p1);
    await p1.click('details.plug-install > summary');
    await p1.setInputFiles('form[action="/admin/plugins/upload"] #plugin', ZIP);
    await p1.click('form[action="/admin/plugins/upload"] button[type=submit]');
    await settle(p1);
  }
  await p1.goto(BASE + '/admin/plugins');
  await settle(p1);
  const plug = p1.locator('details.plug', { has: p1.locator('code', { hasText: /^kontaktformular$/ }) });
  if (await plug.count() !== 1) throw new Error('kontaktformular is not on /admin/plugins after the upload');
  // both controls post on change through htmx and answer with HX-Redirect to
  // /admin/plugins, a fresh page on which the plugin's panel is closed again
  async function openPlug() { if (await plug.getAttribute('open') === null) await plug.locator('summary').click(); }
  await openPlug();
  await plug.locator(`form[action="/admin/plugins/kontaktformular/websites"] input[name=website][value="${W}"]`).click();
  await settle(p1);
  await openPlug();
  await plug.locator('form[action="/admin/plugins/kontaktformular/enable"] input[name=an]').click();
  await settle(p1);
  if (await confirmIfAsked(p1, A1.pw)) {
    note('switching it on asked for the password again');
    await p1.goto(BASE + '/admin/plugins'); await settle(p1);
    await openPlug();
    await plug.locator('form[action="/admin/plugins/kontaktformular/enable"] input[name=an]').click();
    await settle(p1);
  }
  await openPlug();
  note('kontaktformular acts on website ' + W + ': ' + (await plug.locator(`input[name=website][value="${W}"]`).isChecked() ? 'yes' : 'no'));
  await p1.goto(BASE + '/admin/plugins');
  await settle(p1);
  note('kontaktformular on /admin/plugins: ' + (await plug.locator('summary').innerText()).replace(/\s+/g, ' ').slice(0, 120));
  await shot(p1, 'setup-plugins.png');
  // a published page holding [[formular]]
  await p1.goto(`${BASE}/admin/websites/${W}/pages/new`);
  await settle(p1);
  await p1.fill('form.editor #title', 'Kontakt');
  await p1.fill('form.editor #content_markdown', 'Schreiben Sie uns.\n\n[[formular]]');
  await p1.click('label[for=panel-settings]');
  await p1.selectOption('form.editor #status', 'published');
  await p1.fill('form.editor #slug', 'kontakt');
  await p1.locator('form.editor button.btn--primary[type=submit]').first().click();
  await settle(p1);
  note('page Kontakt saved, editor at ' + pathOf(p1));

  // =====================================================================
  // UAT-02: two orders, invoice and prepayment
  const visits = [];
  async function visitor() {
    const ctx = await context(browser);
    const log = [];
    ctx.on('request', r => log.push(r.url()));
    visits.push(log);
    return ctx;
  }
  async function order(method, who, first) {
    const ctx = await visitor();
    const p = await ctx.newPage();
    const res = await p.goto(PUB + '/shop/holzbrett');
    await settle(p);
    if (first) {
      put('PRODUCT-PAGE', res ? res.status() : '?');
      const before = (await ctx.cookies()).map(c => c.name);
      put('COOKIE-BEFORE-ACTION', before.length ? before.join(', ') : 'none');
    }
    await p.click('form[action$="/warenkorb/hinzufuegen"] button[type=submit]');
    await settle(p);
    if (first) put('BASKET-COOKIE', (await ctx.cookies()).map(c => c.name).join(', ') || 'none');
    await p.goto(PUB + '/kasse');
    await settle(p);
    if (first) {
      const offered = await p.locator('input[name=zahlungsart]').evaluateAll(es => es.map(e => e.value));
      put('PAYMENT-OFFERED', offered.join(', ') || 'none');
      put('PAYREXX-OFFERED', offered.includes('payrexx') ? 'yes' : 'no');
    }
    await p.fill('[name=email]', who.email);
    await p.fill('[name=name]', who.name);
    await p.fill('[name=strasse]', who.street);
    await p.selectOption('[name=land]', 'CH');
    await p.fill('[name=plz]', who.plz);
    await p.fill('[name=ort]', who.city);
    await p.check(`input[name=zahlungsart][value=${method}]`);
    await p.check('input[name=bestaetigt]');
    await p.click('form[action$="/kasse"] button[type=submit]');
    await settle(p);
    const url = pathOf(p);
    const number = (((await p.locator('h1').first().innerText().catch(() => '')) || '').match(/\d{4}-\d{4}/) || ['?'])[0];
    await shot(p, `uat02-order-${method}.png`);
    await ctx.close();
    return { shape: /^\/bestellung\/[0-9a-f]{32}$/.test(url) ? 'matches /bestellung/[0-9a-f]{32}' : 'does not match (' + url.replace(/[0-9a-f]{32}/g, '…') + ')', number };
  }
  const o1 = await order('invoice', { email: 'kundin1@example.test', name: 'Ada Kundin', street: 'Marktgasse 1', plz: '8001', city: 'Zürich' }, true);
  put('ORDER-1-URL-SHAPE', o1.shape);
  put('ORDER-1-NUMBER', o1.number);
  const o2 = await order('prepay', { email: 'kunde2@example.test', name: 'Bruno Kunde', street: 'Seestrasse 2', plz: '6003', city: 'Luzern' }, false);
  put('ORDER-2-URL-SHAPE', o2.shape);
  put('ORDER-2-NUMBER', o2.number);
  await p1.goto(`${BASE}/admin/websites/${W}/bestellungen`);
  await settle(p1);
  const listText = await p1.innerText('main');
  await shot(p1, 'uat02-admin-orders.png');
  for (const [n, o] of [[1, o1], [2, o2]]) {
    put(`ORDER-${n}-ADMIN-LISTED`, o.number !== '?' && listText.includes(o.number) ? 'yes' : 'no');
    await p1.goto(`${BASE}/admin/websites/${W}/bestellungen/${o.number}`);
    await settle(p1);
    const pairs = await p1.locator('dl dt').evaluateAll(ds => ds.map(d => [d.textContent.trim(), (d.nextElementSibling && d.nextElementSibling.textContent || '').trim()]));
    const pay = pairs.find(([k]) => /Zahlungsart|Payment method/i.test(k));
    put(`ORDER-${n}-ADMIN-PAYMENT`, pay ? pay[1].split(/\s+—\s+/)[0].trim() : '?');
    const st = await p1.locator('select#status option:checked').innerText().catch(() => '?');
    put(`ORDER-${n}-ADMIN-STATUS`, st.trim());
    await shot(p1, `uat02-admin-order-${n}.png`);
  }

  // =====================================================================
  // UAT-03: a message after a wait, a submission into the honeypot
  async function contact(subject, trap) {
    const ctx = await visitor();
    const p = await ctx.newPage();
    const res = await p.goto(PUB + '/kontakt');
    await settle(p);
    const drawn = Date.now();
    await p.fill('#cf-name', trap ? 'Robo Spam' : 'Clara Fragt');
    await p.fill('#cf-email', trap ? 'robo@example.test' : 'clara@example.test');
    await p.fill('#cf-subject', subject);
    await p.fill('#cf-body', trap ? 'Billige Uhren, jetzt kaufen.' : 'Gibt es das Holzbrett auch in Eiche?');
    const consent = p.locator('input[name=einwilligung]');
    if (await consent.count()) await consent.check();
    let field = '';
    if (trap) {
      field = (await p.getAttribute('#cf-website', 'name')) || '?';
      await p.locator('#cf-website').fill('https://spam.example', { force: true });
    }
    // the time trap counts whole seconds against a three-second minimum: wait four and a half
    const rest = 4500 - (Date.now() - drawn);
    if (rest > 0) await p.waitForTimeout(rest);
    const waited = Date.now() - drawn;
    await p.click('form.contact-form button[type=submit]');
    await settle(p);
    const q = new URL(p.url()).searchParams.get('formular') || '';
    const notice = ((await p.locator('.contact-form__notice').first().innerText().catch(() => '')) || '').trim();
    await shot(p, trap ? 'uat03-spam-answer.png' : 'uat03-message-answer.png');
    await ctx.close();
    return { status: res ? res.status() : '?', waited, answer: `formular=${q}, ${notice}`, field };
  }
  const SUBJ = 'Frage zum Holzbrett', SPAM = 'Billige Uhren';
  const m = await contact(SUBJ, false);
  put('FORM-PAGE', m.status);
  put('FORM-WAIT-MS', m.waited);
  put('MESSAGE-ANSWER', m.answer);
  const sp = await contact(SPAM, true);
  put('SPAM-FIELD', sp.field);
  put('SPAM-WAIT-MS', sp.waited);
  put('SPAM-ANSWER-SAME', sp.answer === m.answer ? 'yes' : 'no (' + sp.answer + ')');
  await p1.goto(`${BASE}/admin/websites/${W}/plugins/kontaktformular`);
  await settle(p1);
  const inbox = await p1.innerText('main');
  put('MESSAGE-IN-ADMIN', inbox.includes(SUBJ) ? 'yes' : 'no');
  put('SPAM-NOT-IN-MESSAGES', inbox.includes(SPAM) ? 'no' : 'yes');
  await shot(p1, 'uat03-admin-messages.png');
  await p1.goto(`${BASE}/admin/websites/${W}/plugins/kontaktformular?ansicht=abgewiesen`);
  await settle(p1);
  const card = p1.locator('main article.card', { hasText: SPAM });
  put('SPAM-IN-QUARANTINE', await card.count() === 1 ? 'yes' : 'no');
  put('QUARANTINE-REASON', await card.count() ? (await card.locator('p.text-muted strong').first().innerText()).trim() : '');
  await shot(p1, 'uat03-admin-quarantine.png');
  const all = visits.flat();
  const origins = [...new Set(all.map(u => new URL(u).origin))];
  put('SAME-ORIGIN', (origins.length === 1 && origins[0] === PUB ? 'yes' : 'no') + ` (${all.length} requests, ${origins.length} origin${origins.length === 1 ? '' : 's'})`);
  if (origins.length !== 1) note('origins: ' + origins.join(', '));
  put('JS', 'on');
  await c1.close();

  await browser.close();
  process.exit(0);
})().catch(e => {
  console.error(scrub(e && e.stack || e));
  put('SCRIPT-ERROR', scrub(String(e && e.message || e).split('\n')[0]).slice(0, 160));
  process.exit(2);
});
```

#### mac-t34-uat/server-excerpt.txt

```text
{"time":"2026-11-21T15:37:36.513086+01:00","level":"INFO","msg":"starting","version":"v0.1.4-266-g6e26d6c","commit":"6e26d6c"}
{"time":"2026-11-21T15:37:36.518585+01:00","level":"INFO","msg":"pre-upgrade snapshot written","path":"$UAT_DIR/data/pre-upgrade-v0.1.4-266-g6e26d6c-20261121-143736.sqlite","level":0}
{"time":"2026-11-21T15:37:37.423885+01:00","level":"INFO","msg":"server starting","addr":"127.0.0.1:18793","log_level":"INFO"}
{"time":"2026-11-21T15:37:39.354251+01:00","level":"INFO","msg":"two factor enabled","user_id":1}
{"time":"2026-11-21T15:37:46.88525+01:00","level":"INFO","msg":"two factor enabled","user_id":2}
{"time":"2026-11-21T15:38:19.771109+01:00","level":"INFO","msg":"plugin loaded","plugin":"kontaktformular","version":"1.6.0"}
{"time":"2026-11-21T15:38:37.425464+01:00","level":"WARN","msg":"mail not sent","id":1,"kind":"order_customer","attempt":1,"err":"verbindung zum mailserver fehlgeschlagen: dial tcp 127.0.0.1:2: connect: connection refused"}
{"time":"2026-11-21T15:38:37.425923+01:00","level":"WARN","msg":"mail not sent","id":2,"kind":"order_operator","attempt":1,"err":"verbindung zum mailserver fehlgeschlagen: dial tcp 127.0.0.1:2: connect: connection refused"}
{"time":"2026-11-21T15:38:37.426247+01:00","level":"WARN","msg":"mail not sent","id":3,"kind":"order_customer","attempt":1,"err":"verbindung zum mailserver fehlgeschlagen: dial tcp 127.0.0.1:2: connect: connection refused"}
{"time":"2026-11-21T15:38:37.426539+01:00","level":"WARN","msg":"mail not sent","id":4,"kind":"order_operator","attempt":1,"err":"verbindung zum mailserver fehlgeschlagen: dial tcp 127.0.0.1:2: connect: connection refused"}
{"time":"2026-11-21T15:38:44.860275+01:00","level":"INFO","msg":"plugin","plugin":"kontaktformular","message":"honeypot triggered"}
{"time":"2026-11-21T15:38:48.074834+01:00","level":"INFO","msg":"shutting down","timeout_s":10}
{"time":"2026-11-21T15:38:48.078622+01:00","level":"INFO","msg":"shutdown complete"}
migrations applied at this start: 69 (versions 1 to 69)
server log lines: 89
```

#### mac-t34-evidence-commit.txt

```text
6e26d6c
```

#### mac-t34-pass1.status

```text
0
```
