package admin

import (
	"context"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/holzcloud/holzcloud-cms/internal/user"
)

// Correction by the rule of each record kind, on the Privacy screen
// (Phase 33.2, plan 09).

func TestAnInvoiceCorrectionKeepsTheOriginalInTheLog(t *testing.T) {
	h, sm, database, ws := eraseAdmin(t)
	privacyOrder(t, database, ws.ID, "2026-0042", "ada@example.com", "Ada Lovelace")
	c := &eraseClient{t: t, h: h, sm: sm, fresh: true}

	body := c.get(h.HandlePrivacyPerson, privacyScreen(ws.ID)+"?email=ada%40example.com", ws.ID).Body.String()
	for _, field := range []string{"email", "phone", "name", "company", "vat_number", "street", "postal_code", "city", "country"} {
		if _, _, ok := correctFormFor(body, field); !ok {
			t.Errorf("the order has no form for %s", field)
		}
	}

	rec, bad, good := correctPost(c, ws.ID, url.Values{
		"email": {"ada@example.com"}, "source": {"orders"}, "key": {"2026-0042"},
		"field": {"city"}, "value": {"Basel"},
	}, false)
	if rec.Header.Get("Location") != privacyScreen(ws.ID)+"?email=ada%40example.com" {
		t.Errorf("Location %q", rec.Header.Get("Location"))
	}
	if bad != "" || good != "Corrected. An invoice already issued with the old details may need a corrected invoice." {
		t.Errorf("flash %q / %q", bad, good)
	}
	if got := correctOrderColumn(t, database, "2026-0042", "city"); got != "Basel" {
		t.Errorf("city = %q", got)
	}
	entries := correctEntries(t, database)
	if len(entries) != 1 {
		t.Fatalf("%d entries", len(entries))
	}
	e := entries[0]
	if e["field"] != "city" || e["invoice"] != true || e["old"] != "Bern" || e["new"] != "Basel" {
		t.Errorf("entry %v, want field city, invoice, old Bern, new Basel", e)
	}

	// A refused value is said, and changes nothing.
	_, bad, _ = correctPost(c, ws.ID, url.Values{
		"email": {"ada@example.com"}, "source": {"orders"}, "key": {"2026-0042"},
		"field": {"country"}, "value": {"Schweiz"},
	}, false)
	if bad == "" {
		t.Error("a country that is not a code was not refused")
	}
	if got := correctOrderColumn(t, database, "2026-0042", "country"); got != "CH" {
		t.Errorf("country = %q", got)
	}
}

func TestTheScreenOffersWhatEachKindAllows(t *testing.T) {
	h, sm, database, ws := eraseAdmin(t)
	ctx := context.Background()
	privacyOrder(t, database, ws.ID, "2026-0042", "ada@example.com", "Ada Lovelace")
	privacyOrder(t, database, ws.ID, "2026-0043", "ada@example.com", "Ada Lovelace")
	if _, err := database.Write.Exec(`UPDATE orders SET erased_at = '2026-11-01T00:00:00Z' WHERE number = '2026-0043'`); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Write.Exec(`INSERT INTO privacy_requests (website_id, kind, locale, name, email, note, created_at)
		VALUES ($1, 'optout', 'de', 'Ada', 'ada@example.com', '', '2026-02-02T09:00:00Z')`, ws.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Write.Exec(`INSERT INTO mail_outbox (website_id, recipient, subject, body, created_at)
		VALUES ($1, 'ada@example.com', 'Kopie', 'Text', '2026-02-01T09:00:01Z')`, ws.ID); err != nil {
		t.Fatal(err)
	}
	staff, err := user.NewStore(database, cheapHashing).Create(ctx, "Ada", "ada@example.com", "passwort123", user.RoleEditor)
	if err != nil {
		t.Fatal(err)
	}
	editLink := `href="/admin/users/` + strconv.FormatInt(staff, 10) + `/edit"`

	admin := &eraseClient{t: t, h: h, sm: sm, role: user.RoleAdmin, fresh: true}
	body := admin.get(h.HandlePrivacyPerson, privacyScreen(ws.ID)+"?email=ada%40example.com", ws.ID).Body.String()
	if !strings.Contains(body, retentionSays("Edit the account")) || !strings.Contains(body, editLink) {
		t.Error("an administrator is not offered the account's edit screen")
	}
	for _, sentence := range []string{
		"This is a record of a request and is not corrected.",
		"This is a copy of a mail as it was queued and is not corrected.",
		"This order was erased and cannot be corrected.",
	} {
		if !strings.Contains(body, retentionSays(sentence)) {
			t.Errorf("the screen does not say %q", sentence)
		}
	}
	if !strings.Contains(body, retentionSays("Invoice fields keep their old value in the activity log, and an invoice already issued may need a corrected invoice.")) {
		t.Error("the order's correction does not say what an invoice field correction keeps")
	}
	// Forms only for what may be corrected: the open order, and none for the
	// erased one, the account, the request or the mail.
	for _, m := range correctFormRe.FindAllStringSubmatch(body, -1) {
		if m[1] != correctAddress(ws.ID) {
			continue
		}
		if !strings.Contains(m[3], `name="source" value="orders"`) || !strings.Contains(m[3], `name="key" value="2026-0042"`) {
			t.Errorf("a correction form for something else: %s", m[3])
		}
	}

	editor := &eraseClient{t: t, h: h, sm: sm, role: user.RoleEditor, fresh: true}
	eb := editor.get(h.HandlePrivacyPerson, privacyScreen(ws.ID)+"?email=ada%40example.com", ws.ID).Body.String()
	if strings.Contains(eb, retentionSays("Edit the account")) {
		t.Error("somebody who is not an administrator is offered the account's edit screen")
	}
}
