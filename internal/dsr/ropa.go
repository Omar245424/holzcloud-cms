package dsr

// The record of processing activities of a website (Phase 33.2, DSR-05,
// D-02).
//
// GDPR Art. 30(1) and the Swiss DSG Art. 12(2) ask for nearly the same
// contents, and the struct below carries both:
//
//   - the controller and its contact details, with any representative and
//     data protection officer: Art. 30(1)(a), DSG 12(2)(a) — Controller;
//   - the purposes: Art. 30(1)(b), DSG 12(2)(b) — Activity.Purpose;
//   - the categories of data subjects and of personal data: Art. 30(1)(c),
//     DSG 12(2)(c) — Activity.Subjects and Activity.Data;
//   - the categories of recipients: Art. 30(1)(d), DSG 12(2)(d) —
//     Activity.Recipients;
//   - transfers to third countries and their safeguards: Art. 30(1)(e),
//     DSG 12(2)(g) — Activity.Transfers;
//   - the time limits for erasure: Art. 30(1)(f), DSG 12(2)(e) —
//     Activity.Retention;
//   - the security measures: Art. 30(1)(g), DSG 12(2)(f) — Security.
//
// The record is generated, not typed. It follows from what the website
// uses (privacy.Features), the plugins that keep data for it, the frameworks
// it follows and the periods it keeps, the same facts the published privacy
// notice is composed from. A record kept by hand drifts from the notice the
// day somebody switches the shop on; this one cannot (DSR-05).
//
// What the CMS cannot know is a fill-in and never a guess: the hosting
// provider, the mail provider, the backup target, the country the server
// stands in, a data protection officer, the TLS of the proxy in front. The
// one provider the record names is Payrexx, because the CMS integrates that
// provider and no other. Three statements of the research are assumptions,
// and the record words them as sentences for the operator to check rather
// than as facts: that a small EU shop is not exempt under Art. 30(5) because
// its processing is not occasional (A5), that a transfer to Payrexx in
// Switzerland rests on the Commission's adequacy decision (A6), and that the
// UK's Art. 30 is unchanged (A11).
//
// Every sentence is a catalogue key, translated when the record is built.
// The admin shows the record in the operator's language; the access export
// carries its categories in the website's.

import (
	"strings"
	"time"

	"github.com/holzcloud/holzcloud-cms/internal/i18n"
	"github.com/holzcloud/holzcloud-cms/internal/legal"
	"github.com/holzcloud/holzcloud-cms/internal/page"
	"github.com/holzcloud/holzcloud-cms/internal/privacy"
)

// RecordInput is everything the record is generated from. The admin fills
// it from the website; BuildRecord reads nothing else.
type RecordInput struct {
	WebsiteID int64
	// The controller, as the website's settings name it.
	Name, Street, PostalCode, City, Country, Email, Phone string
	Set                                                   legal.Set
	Features                                              privacy.Features
	// Holders are the plugins that keep data for the website.
	Holders []Holder
	// Periods are the website's retention periods (Service.Periods).
	Periods []Period
	// AccountMonths is the installation's account period, LogDays the
	// server log's, BackupDays how long the nightly backups are kept.
	AccountMonths, LogDays, BackupDays int
	// OrderPeriod is the law of the seller's country, OrderYears the period
	// in force (Service.OrderYears).
	OrderPeriod legal.OrderPeriod
	OrderYears  int
	// MailConfigured is true when mail about the website is delivered.
	MailConfigured bool
	// Lang is the language the record is written in.
	Lang string
	Now  time.Time
}

// Activity is one processing activity. Every cell is text in the record's
// language.
type Activity struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Purpose    string `json:"purpose"`
	Subjects   string `json:"data_subjects"`
	Data       string `json:"data_categories"`
	LegalBasis string `json:"legal_basis"`
	Recipients string `json:"recipients"`
	Transfers  string `json:"transfers"`
	Retention  string `json:"retention"`
	// Kinds are the retention kinds of the records the activity produces
	// (Record.Kind), so the access export can say which activity a record
	// it hands over belongs to.
	Kinds []string `json:"kinds"`
}

// Exemption says whether a framework excuses a small organisation from
// keeping the record.
type Exemption struct {
	Framework string `json:"framework"`
	Text      string `json:"text"`
}

// ProcessingRecord is the record of one website.
type ProcessingRecord struct {
	GeneratedAt string `json:"generated_at"`
	// Controller holds name, street, postal_code, city, country, email,
	// phone, representative and data_protection_officer.
	Controller map[string]string `json:"controller"`
	Frameworks []string          `json:"frameworks"`
	Activities []Activity        `json:"activities"`
	Security   []string          `json:"security"`
	Exemptions []Exemption       `json:"exemptions"`
	// Lang is the language the record is written in: the one asked for, or
	// English when this build has no catalogue for it.
	Lang string `json:"lang"`
}

// fillIn marks what only the operator knows.
var fillIn = i18n.N("(to be filled in)")

// The activities' names.
var (
	nameServerLog    = i18n.N("Delivering the website (server log)")
	nameContactForms = i18n.N("Contact forms")
	nameOrders       = i18n.N("Shop orders")
	nameFarmShop     = i18n.N("Orders through the order form")
	namePayment      = i18n.N("Online payment")
	nameProtected    = i18n.N("Protected pages")
	nameRequests     = i18n.N("Requests about personal data and opt-outs")
	nameStaff        = i18n.N("Staff accounts and activity log")
	nameMail         = i18n.N("E-mail sent by the website")
	nameBackups      = i18n.N("Backups")
)

// The purposes.
var (
	purposeServerLog    = i18n.N("Delivering the pages of the website and keeping it secure.")
	purposeContactForms = i18n.N("Answering messages sent through a form on the website.")
	purposeOrders       = i18n.N("Carrying out orders, and keeping the records the law requires.")
	purposeFarmShop     = i18n.N("Carrying out orders placed through the order form, and keeping the records the law requires.")
	purposePayment      = i18n.N("Taking payment for an order online.")
	purposeProtected    = i18n.N("Showing pages only to people who know their password.")
	purposeRequests     = i18n.N("Answering requests about personal data and opt-outs, and being able to show how they were answered.")
	purposeStaff        = i18n.N("Letting people work on the website, and tracing who changed what.")
	purposeMail         = i18n.N("Sending order confirmations, notifications about messages and requests, and sign-in links.")
	purposeBackups      = i18n.N("Restoring the website after a failure.")
)

// The data subjects.
var (
	subjectsVisitors  = i18n.N("Visitors of the website")
	subjectsSenders   = i18n.N("People who write through a form")
	subjectsCustomers = i18n.N("Customers")
	subjectsPayers    = i18n.N("Customers who pay online")
	subjectsUnlockers = i18n.N("Visitors who enter a page password")
	subjectsRequests  = i18n.N("People who make a request about their data")
	subjectsStaff     = i18n.N("People who work on the website")
	subjectsMail      = i18n.N("Customers, people who write through a form or make a request, and staff")
	subjectsAll       = i18n.N("Everybody named in the activities above")
)

// The categories of data.
var (
	dataServerLog    = i18n.N("Shortened IP address, the address requested, time, status and size of the response; the full IP address of a failed sign-in to the admin area.")
	dataContactForms = i18n.N("Name, e-mail address, subject, message, the fields the operator added, the consent text where one is asked for, attachments.")
	dataOrders       = i18n.N("Name, company, VAT number, address, e-mail address, telephone number, note, the items ordered, the choice of delivery and payment.")
	dataPayment      = i18n.N("Amount, currency and order number, with name, e-mail address, street, postcode, town and country.")
	dataProtected    = i18n.N("A signed cookie with the page number and an expiry time; nothing about the person beyond the server log.")
	dataRequests     = i18n.N("Name, e-mail address, the request, how it arrived and how it was answered, the language of the page; for an opt-out whether the browser sent Global Privacy Control.")
	dataStaff        = i18n.N("Name, e-mail address, password as an Argon2id hash, second factor, sessions, sign-in times, the changes made to the website.")
	dataMail         = i18n.N("Recipient address, subject and text of the mail.")
	dataBackups      = i18n.N("A copy of the database and the files, with everything listed above.")
)

// The legal bases, one sentence per basis and framework. An activity names
// the bases it rests on, and the record writes the sentences of the
// frameworks the website follows.
const (
	basisContract   = "contract"
	basisObligation = "obligation"
	basisInterest   = "interest"
)

var (
	basisEU = map[string]string{
		basisContract:   i18n.N("Art. 6(1)(b) GDPR: performing a contract, or steps taken at the person's request before one."),
		basisObligation: i18n.N("Art. 6(1)(c) GDPR: a legal obligation."),
		basisInterest:   i18n.N("Art. 6(1)(f) GDPR: a legitimate interest in the purpose named."),
	}
	basisUK = map[string]string{
		basisContract:   i18n.N("Art. 6(1)(b) UK GDPR: performing a contract, or steps taken at the person's request before one."),
		basisObligation: i18n.N("Art. 6(1)(c) UK GDPR: a legal obligation."),
		basisInterest:   i18n.N("Art. 6(1)(f) UK GDPR: a legitimate interest in the purpose named."),
	}
	basisCH = map[string]string{
		basisContract:   i18n.N("Art. 31(2)(a) FADP: directly connected with the conclusion or performance of a contract."),
		basisObligation: i18n.N("Art. 31(1) FADP: required by law."),
		basisInterest:   i18n.N("Art. 31(1) FADP: an overriding private interest in the purpose named."),
	}
)

// The recipients.
var (
	recipientsNone     = i18n.N("Nobody outside this server.")
	recipientsHosting  = i18n.N("The hosting provider of the server:")
	recipientsMail     = i18n.N("The mail provider that delivers the website's mail:")
	recipientsNoMail   = i18n.N("No mail server is set up for this website, so no mail leaves the server.")
	recipientsPayrexx  = i18n.N("Payrexx, the payment provider, for an order paid online.")
	recipientsPayment  = i18n.N("Payrexx, the payment provider. Name, address and the date of the data processing agreement:")
	recipientsBackups  = i18n.N("The place the backups are copied to:")
	transfersServer    = i18n.N("The country the server stands in, and the safeguards if it is outside the EU, the EEA, Switzerland or the UK:")
	transfersPayrexx   = i18n.N("Payrexx processes the payment in Switzerland. For a website in the EU or the UK, check that the adequacy decision for Switzerland covers it (Art. 45 GDPR).")
	retentionPayment   = i18n.N("Payrexx keeps the payment under its own terms; this website keeps only the payment status with the order.")
	retentionLog       = i18n.N("%d days, enforced by the system journal.")
	retentionForms     = i18n.N("%d months after it arrived, then deleted by the retention job.")
	retentionOrderSays = i18n.N("%d years from the end of the year of the order, then deleted by the retention job.")
	retentionRequests  = i18n.N("%d months after the request or its answer, whichever is later.")
	retentionUnlock    = i18n.N("The cookie expires after %d hours.")
	retentionStaff     = i18n.N("An account is deleted after %d months without a sign-in; administrators are never deleted automatically. The activity log is kept until an administrator clears it.")
	retentionMail      = i18n.N("Sent mail for 90 days; mail that could not be sent for 90 days after it was queued.")
	retentionBackups   = i18n.N("%d days for the nightly backups; the five newest snapshots taken before an update are kept with no time limit.")
)

// The security measures, each one the repository enforces and
// docs/security.md documents.
var securityMeasures = []string{
	i18n.N("Passwords are stored only as Argon2id hashes."),
	i18n.N("Administrators must sign in with a second factor."),
	i18n.N("Sessions are stored on the server and renewed at every sign-in."),
	i18n.N("Every request that changes something is protected against cross-site request forgery."),
	i18n.N("A Content-Security-Policy of default-src 'self' on every response; nothing is loaded from other servers."),
	i18n.N("The password is asked again before irreversible actions, such as an erasure."),
	i18n.N("The backups the shipped backup script writes carry SHA-256 checksums."),
	i18n.N("IP addresses are shortened before the request log writes them."),
}

// securityTLS is the measure the CMS cannot see: the encryption the proxy in
// front of it does.
var securityTLS = i18n.N("Encryption in transit (TLS) at the reverse proxy in front of the CMS:")

// The exemptions per framework (33.2-RESEARCH §B.4).
var exemptionTexts = map[legal.Framework]string{
	legal.EU: i18n.N("Art. 30(5) GDPR: an organisation with fewer than 250 employees is exempt only if its processing is occasional, unlikely to result in a risk and includes no special categories of data. A website that takes orders or messages all the time is generally not occasional, so keep this record."),
	legal.UK: i18n.N("Art. 30(5) UK GDPR: an organisation with fewer than 250 employees is exempt only if its processing is occasional, unlikely to result in a risk and includes no special categories of data. A website that takes orders or messages all the time is generally not occasional, so keep this record."),
	legal.CH: i18n.N("Art. 24 of the Data Protection Ordinance: a business with fewer than 250 employees on 1 January is exempt, unless it processes sensitive personal data on a large scale or carries out high-risk profiling. Keeping the record anyway does no harm."),
	legal.US: i18n.N("The US state privacy laws ask for no record of processing activities. The record still helps to answer requests and to check the privacy notice."),
}

// The plugins whose data the record describes by name. Any other plugin
// that keeps data gets an activity of fill-ins.
const (
	pluginContactForm = "kontaktformular"
	pluginOrderForm   = "bestellung"
)

// recordLang is the language the record is written in: the one asked for
// when this build can show it, English otherwise.
func recordLang(lang string) string {
	if lang == "" {
		return i18n.Source
	}
	if i18n.Known(lang) || (i18n.Base(lang) != "" && i18n.Known(i18n.Base(lang))) {
		return lang
	}
	return i18n.Source
}

// BuildRecord generates the record. It is pure: the same input gives the
// same record.
func BuildRecord(in RecordInput) ProcessingRecord {
	lang := recordLang(in.Lang)
	t := func(s string) string { return i18n.T(lang, s) }
	fill := t(fillIn)
	withFill := func(s string) string { return t(s) + " " + fill }
	months := map[string]int{}
	for _, p := range in.Periods {
		months[p.Kind] = p.Months
	}
	bases := func(kinds ...string) string {
		var lines []string
		for _, f := range in.Set.List() {
			var table map[string]string
			switch f {
			case legal.EU:
				table = basisEU
			case legal.UK:
				table = basisUK
			case legal.CH:
				table = basisCH
			}
			for _, k := range kinds {
				if s, ok := table[k]; ok {
					lines = append(lines, t(s))
				}
			}
		}
		if len(lines) == 0 {
			return fill
		}
		return strings.Join(lines, "\n")
	}
	mailRecipient := t(recipientsNone)
	mailOut := t(recipientsNoMail)
	if in.MailConfigured {
		mailRecipient = withFill(recipientsMail)
		mailOut = mailRecipient
	}
	transfers := withFill(transfersServer)

	covered := map[string]bool{}
	for _, h := range in.Holders {
		if h.Covered {
			covered[h.PluginID] = true
		}
	}
	shop := in.Features.Shop
	payment := shop && in.Features.Payment

	var acts []Activity
	acts = append(acts, Activity{
		ID: "server_log", Name: t(nameServerLog), Purpose: t(purposeServerLog),
		Subjects: t(subjectsVisitors), Data: t(dataServerLog), LegalBasis: bases(basisInterest),
		Recipients: withFill(recipientsHosting), Transfers: transfers,
		Retention: i18n.Tf(lang, retentionLog, in.LogDays),
	})
	if in.Features.Forms || covered[pluginContactForm] {
		acts = append(acts, Activity{
			ID: "contact_forms", Name: t(nameContactForms), Purpose: t(purposeContactForms),
			Subjects: t(subjectsSenders), Data: t(dataContactForms), LegalBasis: bases(basisContract, basisInterest),
			Recipients: mailRecipient, Transfers: transfers,
			Retention: i18n.Tf(lang, retentionForms, months[RetentionFormSubmission]),
			Kinds:     []string{RetentionFormSubmission},
		})
	}
	orderRetention := i18n.Tf(lang, retentionOrderSays, in.OrderYears) + " " + t(in.OrderPeriod.Basis())
	if shop {
		// "Nobody outside this server" and Payrexx in one cell would
		// contradict each other: without mail, Payrexx stands alone.
		recipients := mailRecipient
		if payment && in.MailConfigured {
			recipients += "\n" + t(recipientsPayrexx)
		} else if payment {
			recipients = t(recipientsPayrexx)
		}
		acts = append(acts, Activity{
			ID: "orders", Name: t(nameOrders), Purpose: t(purposeOrders),
			Subjects: t(subjectsCustomers), Data: t(dataOrders), LegalBasis: bases(basisContract, basisObligation),
			Recipients: recipients, Transfers: transfers, Retention: orderRetention,
			Kinds: []string{RetentionOrder},
		})
	}
	if covered[pluginOrderForm] {
		acts = append(acts, Activity{
			ID: "farm_shop_orders", Name: t(nameFarmShop), Purpose: t(purposeFarmShop),
			Subjects: t(subjectsCustomers), Data: t(dataOrders), LegalBasis: bases(basisContract, basisObligation),
			Recipients: mailRecipient, Transfers: transfers, Retention: orderRetention,
			Kinds: []string{RetentionOrder},
		})
	}
	if payment {
		acts = append(acts, Activity{
			ID: "payment", Name: t(namePayment), Purpose: t(purposePayment),
			Subjects: t(subjectsPayers), Data: t(dataPayment), LegalBasis: bases(basisContract),
			Recipients: withFill(recipientsPayment),
			Transfers:  t(transfersPayrexx) + "\n" + transfers,
			Retention:  t(retentionPayment),
		})
	}
	if in.Features.ProtectedPages {
		acts = append(acts, Activity{
			ID: "protected_pages", Name: t(nameProtected), Purpose: t(purposeProtected),
			Subjects: t(subjectsUnlockers), Data: t(dataProtected), LegalBasis: bases(basisInterest),
			Recipients: t(recipientsNone), Transfers: transfers,
			Retention: i18n.Tf(lang, retentionUnlock, int(page.UnlockLifetime/time.Hour)),
		})
	}
	acts = append(acts,
		Activity{
			ID: "requests", Name: t(nameRequests), Purpose: t(purposeRequests),
			Subjects: t(subjectsRequests), Data: t(dataRequests), LegalBasis: bases(basisObligation),
			Recipients: mailRecipient, Transfers: transfers,
			Retention: i18n.Tf(lang, retentionRequests, months[RetentionRequest]),
			Kinds:     []string{RetentionRequest},
		},
		Activity{
			ID: "staff_accounts", Name: t(nameStaff), Purpose: t(purposeStaff),
			Subjects: t(subjectsStaff), Data: t(dataStaff), LegalBasis: bases(basisContract, basisInterest),
			Recipients: t(recipientsNone), Transfers: transfers,
			Retention: i18n.Tf(lang, retentionStaff, in.AccountMonths),
			Kinds:     []string{RetentionAccount},
		},
		Activity{
			ID: "mail", Name: t(nameMail), Purpose: t(purposeMail),
			Subjects: t(subjectsMail), Data: t(dataMail), LegalBasis: bases(basisContract, basisInterest),
			Recipients: mailOut, Transfers: transfers, Retention: t(retentionMail),
			Kinds: []string{"mail"},
		},
		Activity{
			ID: "backups", Name: t(nameBackups), Purpose: t(purposeBackups),
			Subjects: t(subjectsAll), Data: t(dataBackups), LegalBasis: bases(basisInterest),
			Recipients: withFill(recipientsBackups), Transfers: transfers,
			Retention: i18n.Tf(lang, retentionBackups, in.BackupDays),
		},
	)
	// A plugin the record cannot describe: its name, and nothing claimed.
	for _, h := range in.Holders {
		if h.Covered && (h.PluginID == pluginContactForm || h.PluginID == pluginOrderForm) {
			continue
		}
		acts = append(acts, Activity{
			ID: "plugin:" + h.PluginID, Name: h.Name, Purpose: fill, Subjects: fill, Data: fill,
			LegalBasis: fill, Recipients: fill, Transfers: fill, Retention: fill,
		})
	}

	controller := map[string]string{
		"name": in.Name, "street": in.Street, "postal_code": in.PostalCode, "city": in.City,
		"country": in.Country, "email": in.Email, "phone": in.Phone,
	}
	for k, v := range controller {
		if strings.TrimSpace(v) == "" {
			controller[k] = fill
		}
	}
	controller["representative"] = fill
	controller["data_protection_officer"] = fill

	rec := ProcessingRecord{
		GeneratedAt: in.Now.UTC().Format(time.RFC3339),
		Controller:  controller,
		Frameworks:  []string{},
		Activities:  acts,
		Exemptions:  []Exemption{},
		Lang:        lang,
	}
	for _, f := range in.Set.List() {
		rec.Frameworks = append(rec.Frameworks, string(f))
		if s, ok := exemptionTexts[f]; ok {
			rec.Exemptions = append(rec.Exemptions, Exemption{Framework: string(f), Text: t(s)})
		}
	}
	for _, s := range securityMeasures {
		rec.Security = append(rec.Security, t(s))
	}
	rec.Security = append(rec.Security, withFill(securityTLS))
	for i := range rec.Activities {
		if rec.Activities[i].Kinds == nil {
			rec.Activities[i].Kinds = []string{}
		}
	}
	return rec
}
