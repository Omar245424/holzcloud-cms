package dsr

import (
	"encoding/json"
	"regexp"
	"slices"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/holzcloud/holzcloud-cms/internal/i18n"
	"github.com/holzcloud/holzcloud-cms/internal/legal"
	"github.com/holzcloud/holzcloud-cms/internal/privacy"
)

// The record of processing activities (Phase 33.2, DSR-05, D-02).

// recordInput is a website with an address, the frameworks EU and CH, the
// default periods of a Swiss seller and nothing switched on.
func recordInput() RecordInput {
	return RecordInput{
		WebsiteID:  1,
		Name:       "Hofladen Muster",
		Street:     "Dorfstrasse 1",
		PostalCode: "3000",
		City:       "Bern",
		Country:    "CH",
		Email:      "hallo@hofladen-muster.example",
		Phone:      "+41 31 000 00 00",
		Set:        legal.ParseSet("eu,ch"),
		Periods: []Period{
			{Kind: RetentionFormSubmission, Months: 12},
			{Kind: RetentionRequest, Months: 24},
			{Kind: RetentionOrder, Months: 120},
		},
		AccountMonths: 24,
		LogDays:       14,
		BackupDays:    30,
		OrderPeriod:   legal.OrderPeriodFor("CH"),
		OrderYears:    10,
		Lang:          i18n.Source,
		Now:           time.Date(2026, 11, 20, 9, 0, 0, 0, time.UTC),
	}
}

func activityIDs(rec ProcessingRecord) []string {
	var out []string
	for _, a := range rec.Activities {
		out = append(out, a.ID)
	}
	sort.Strings(out)
	return out
}

func activityOf(t *testing.T, rec ProcessingRecord, id string) Activity {
	t.Helper()
	for _, a := range rec.Activities {
		if a.ID == id {
			return a
		}
	}
	t.Fatalf("the record has no activity %q: %v", id, activityIDs(rec))
	return Activity{}
}

func sorted(ids ...string) []string {
	sort.Strings(ids)
	return ids
}

func TestTheRecordFollowsWhatTheWebsiteUses(t *testing.T) {
	always := []string{"server_log", "requests", "staff_accounts", "mail", "backups"}

	in := recordInput()
	if got, want := activityIDs(BuildRecord(in)), sorted(always...); !slices.Equal(got, want) {
		t.Errorf("nothing switched on: activities %v, want %v", got, want)
	}

	in.Features = privacy.Features{Shop: true, Payment: true, Forms: true}
	in.Holders = []Holder{{PluginID: "kontaktformular", Name: "Kontaktformular", Entries: 3, Covered: true}}
	if got, want := activityIDs(BuildRecord(in)), sorted(append([]string{"contact_forms", "orders", "payment"}, always...)...); !slices.Equal(got, want) {
		t.Errorf("shop, payment and forms: activities %v, want %v", got, want)
	}

	in.Features.ProtectedPages = true
	if got := activityIDs(BuildRecord(in)); !slices.Contains(got, "protected_pages") {
		t.Errorf("protected pages: activities %v lack protected_pages", got)
	}

	// The farm shop's orders when its plugin answers the privacy hook.
	in.Holders = append(in.Holders, Holder{PluginID: "bestellung", Name: "Bestellung", Entries: 1, Covered: true})
	if got := activityIDs(BuildRecord(in)); !slices.Contains(got, "farm_shop_orders") {
		t.Errorf("a covered bestellung: activities %v lack farm_shop_orders", got)
	}

	// Payment without a shop is no payment.
	noShop := recordInput()
	noShop.Features = privacy.Features{Payment: true}
	if got := activityIDs(BuildRecord(noShop)); slices.Contains(got, "payment") || slices.Contains(got, "orders") {
		t.Errorf("payment without a shop: activities %v", got)
	}

	// A plugin the CMS cannot ask: its own activity, with nothing claimed
	// about it.
	in.Holders = append(in.Holders, Holder{PluginID: "gaestebuch", Name: "Gästebuch", Entries: 7, Reason: reasonNoHook})
	rec := BuildRecord(in)
	g := activityOf(t, rec, "plugin:gaestebuch")
	if g.Name != "Gästebuch" {
		t.Errorf("the plugin's activity is named %q, want the plugin's name", g.Name)
	}
	fill := i18n.T(in.Lang, fillIn)
	for what, cell := range map[string]string{
		"purpose": g.Purpose, "subjects": g.Subjects, "data": g.Data, "legal basis": g.LegalBasis,
		"recipients": g.Recipients, "transfers": g.Transfers, "retention": g.Retention,
	} {
		if cell != fill {
			t.Errorf("gaestebuch's %s = %q, want the fill-in %q", what, cell, fill)
		}
	}

	// Every activity says something in every cell.
	for _, a := range rec.Activities {
		for what, cell := range map[string]string{
			"name": a.Name, "purpose": a.Purpose, "subjects": a.Subjects, "data": a.Data,
			"legal basis": a.LegalBasis, "recipients": a.Recipients, "transfers": a.Transfers, "retention": a.Retention,
		} {
			if strings.TrimSpace(cell) == "" {
				t.Errorf("%s: the %s is empty", a.ID, what)
			}
		}
	}
}

func TestTheRecordStatesThePeriodsItKeeps(t *testing.T) {
	in := recordInput()
	in.Features = privacy.Features{Shop: true, Forms: true}
	in.Periods = []Period{
		{Kind: RetentionFormSubmission, Months: 6},
		{Kind: RetentionRequest, Months: 24},
		{Kind: RetentionOrder, Months: 120},
	}
	in.AccountMonths = 24
	rec := BuildRecord(in)

	if r := activityOf(t, rec, "contact_forms").Retention; !strings.Contains(r, "6 months") {
		t.Errorf("contact forms retention %q does not say 6 months", r)
	}
	if r := activityOf(t, rec, "requests").Retention; !strings.Contains(r, "24 months") {
		t.Errorf("requests retention %q does not say 24 months", r)
	}
	orders := activityOf(t, rec, "orders").Retention
	if !strings.Contains(orders, "10 years") {
		t.Errorf("orders retention %q does not say 10 years", orders)
	}
	if basis := legal.OrderPeriodFor("CH").Basis(); !strings.Contains(orders, basis) {
		t.Errorf("orders retention %q lacks the Swiss basis %q", orders, basis)
	}
	if r := activityOf(t, rec, "staff_accounts").Retention; !strings.Contains(r, "24 months") {
		t.Errorf("staff accounts retention %q does not say 24 months", r)
	}
	if r := activityOf(t, rec, "server_log").Retention; !strings.Contains(r, "14 days") {
		t.Errorf("server log retention %q does not say 14 days", r)
	}
	if r := activityOf(t, rec, "backups").Retention; !strings.Contains(r, "30 days") {
		t.Errorf("backups retention %q does not say 30 days", r)
	}

	// The operator who keeps orders longer than the law sees their own
	// period.
	in.OrderYears = 12
	if r := activityOf(t, BuildRecord(in), "orders").Retention; !strings.Contains(r, "12 years") {
		t.Errorf("orders retention %q does not follow the website's 12 years", r)
	}

	// The record maps the kinds of the search to their activity.
	for id, kind := range map[string]string{
		"contact_forms": RetentionFormSubmission, "orders": RetentionOrder, "requests": RetentionRequest,
		"staff_accounts": RetentionAccount, "mail": "mail",
	} {
		if !slices.Contains(activityOf(t, rec, id).Kinds, kind) {
			t.Errorf("%s does not cover the kind %s", id, kind)
		}
	}
}

func TestTheRecordNamesTheExemptions(t *testing.T) {
	in := recordInput()
	in.Set = legal.ParseSet("eu,ch,uk,us")
	rec := BuildRecord(in)
	got := map[string]string{}
	for _, e := range rec.Exemptions {
		got[e.Framework] = e.Text
	}
	if len(rec.Exemptions) != 4 {
		t.Fatalf("%d exemptions, want 4: %+v", len(rec.Exemptions), rec.Exemptions)
	}
	if !strings.Contains(got["eu"], "Art. 30(5) GDPR") || !strings.Contains(got["eu"], "250") || !strings.Contains(got["eu"], "not occasional") {
		t.Errorf("EU exemption %q", got["eu"])
	}
	if !strings.Contains(got["uk"], "Art. 30(5) UK GDPR") || !strings.Contains(got["uk"], "250") {
		t.Errorf("UK exemption %q", got["uk"])
	}
	if !strings.Contains(got["ch"], "Art. 24 of the Data Protection Ordinance") || !strings.Contains(got["ch"], "250") {
		t.Errorf("CH exemption %q", got["ch"])
	}
	if !strings.Contains(got["us"], "no record of processing") {
		t.Errorf("US exemption %q", got["us"])
	}
	if !slices.Equal(rec.Frameworks, []string{"eu", "ch", "uk", "us"}) {
		t.Errorf("frameworks %v", rec.Frameworks)
	}

	in.Set = legal.ParseSet("ch")
	rec = BuildRecord(in)
	if len(rec.Exemptions) != 1 || rec.Exemptions[0].Framework != "ch" {
		t.Errorf("set ch: exemptions %+v, want the Swiss one only", rec.Exemptions)
	}
}

// T-33.2-48: what the CMS cannot know is a fill-in, never a guess.
func TestTheRecordInventsNoProvider(t *testing.T) {
	in := recordInput()
	in.Set = legal.ParseSet("eu,ch,uk,us")
	in.Features = privacy.Features{Shop: true, Payment: true, Forms: true, ProtectedPages: true}
	in.Holders = []Holder{{PluginID: "kontaktformular", Name: "Kontaktformular", Covered: true}}
	in.MailConfigured = true
	rec := BuildRecord(in)
	fill := i18n.T(in.Lang, fillIn)

	if r := activityOf(t, rec, "server_log").Recipients; !strings.Contains(r, fill) {
		t.Errorf("server log recipients %q: the hosting provider is not a fill-in", r)
	}
	if r := activityOf(t, rec, "mail").Recipients; !strings.Contains(r, fill) {
		t.Errorf("mail recipients %q: the mail provider is not a fill-in", r)
	}
	if r := activityOf(t, rec, "backups").Recipients; !strings.Contains(r, fill) {
		t.Errorf("backups recipients %q: the backup target is not a fill-in", r)
	}
	if rec.Controller["data_protection_officer"] != fill {
		t.Errorf("data protection officer %q, want the fill-in", rec.Controller["data_protection_officer"])
	}
	tls := false
	for _, s := range rec.Security {
		if strings.Contains(s, "TLS") {
			tls = true
			if !strings.Contains(s, fill) {
				t.Errorf("the TLS line %q is not a fill-in", s)
			}
		}
	}
	if !tls {
		t.Error("the security measures have no TLS line")
	}
	for _, a := range rec.Activities {
		if !strings.Contains(a.Transfers, fill) {
			t.Errorf("%s: transfers %q do not leave the server's country to the operator", a.ID, a.Transfers)
		}
	}

	// Nothing the website did not give: no domain and no company name
	// besides the website's own and Payrexx, the one provider the CMS
	// integrates.
	raw, err := json.Marshal(rec)
	if err != nil {
		t.Fatal(err)
	}
	text := string(raw)
	for _, own := range []string{in.Email, in.Name} {
		text = strings.ReplaceAll(text, own, "")
	}
	if m := regexp.MustCompile(`(?i)\b[a-z0-9-]+\.(com|net|org|io|ch|de|eu|uk|cloud|app)\b`).FindAllString(text, -1); len(m) > 0 {
		t.Errorf("the record names domains the website did not give: %v", m)
	}
	for _, company := range []string{"Hetzner", "Infomaniak", "Hostpoint", "Cloudflare", "Amazon", "AWS", "Google",
		"Microsoft", "Azure", "Mailgun", "SendGrid", "Postmark", "Mailjet", "Brevo", "DigitalOcean", "OVH", "Strato", "IONOS",
		" AG", " GmbH", " Ltd", " Inc"} {
		if strings.Contains(text, company) {
			t.Errorf("the record names %q", company)
		}
	}
	if !strings.Contains(text, "Payrexx") {
		t.Error("the payment activity does not name Payrexx, the provider the CMS integrates")
	}
}

func TestTheRecordIsInTheRequestedLanguage(t *testing.T) {
	in := recordInput()
	in.Lang = "de"
	de := BuildRecord(in)
	purpose := activityOf(t, de, "server_log").Purpose
	if purpose == purposeServerLog || purpose != i18n.T("de", purposeServerLog) {
		t.Errorf("de: the server log's purpose is %q, want the German translation", purpose)
	}
	if de.Lang != "de" {
		t.Errorf("de: lang %q", de.Lang)
	}

	in.Lang = "rm"
	rm := BuildRecord(in)
	if got := activityOf(t, rm, "server_log").Purpose; got != purposeServerLog {
		t.Errorf("rm: the server log's purpose is %q, want the English source", got)
	}
	if rm.Lang != "en" {
		t.Errorf("rm: lang %q, want en", rm.Lang)
	}
}
