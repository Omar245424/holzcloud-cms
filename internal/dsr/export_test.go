package dsr

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/holzcloud/holzcloud-cms/internal/i18n"
	"github.com/holzcloud/holzcloud-cms/internal/legal"
	"github.com/holzcloud/holzcloud-cms/internal/privacy"
)

func TestTheExportHoldsTheRecords(t *testing.T) {
	_, database := newTestStore(t)
	a := website(t, database, "Holzwerk")
	number := findOrder(t, database, a, "2026-0003", "Ada@Example.COM")
	res, err := NewService(database).Find(context.Background(), a, "ada@example.com")
	if err != nil {
		t.Fatalf("Find: %v", err)
	}

	now := time.Date(2026, 11, 14, 9, 30, 0, 0, time.UTC)
	raw, err := json.Marshal(BuildExport(res, a, "Holzwerk", now))
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var got struct {
		GeneratedAt string `json:"generated_at"`
		Website     struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"website"`
		Email   string `json:"email"`
		Records map[string][]struct {
			Key    string            `json:"key"`
			Kind   string            `json:"kind"`
			Time   string            `json:"time"`
			Fields map[string]string `json:"fields"`
		} `json:"records"`
		NotCovered json.RawMessage `json:"not_covered"`
		Incomplete json.RawMessage `json:"incomplete"`
	}
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("Unmarshal: %v\n%s", err, raw)
	}
	if _, err := time.Parse(time.RFC3339, got.GeneratedAt); err != nil {
		t.Errorf("generated_at = %q is not RFC 3339", got.GeneratedAt)
	}
	if got.Website.ID != a || got.Website.Name != "Holzwerk" {
		t.Errorf("website = %+v", got.Website)
	}
	if got.Email != "ada@example.com" {
		t.Errorf("email = %q", got.Email)
	}
	orders := got.Records["orders"]
	if len(orders) != 1 || orders[0].Key != number {
		t.Fatalf("records.orders = %+v, want the order %s", orders, number)
	}
	if orders[0].Fields["name"] != "Ada Lovelace" {
		t.Errorf("fields.name = %q", orders[0].Fields["name"])
	}
	if orders[0].Kind != "order" || orders[0].Time == "" {
		t.Errorf("kind/time = %q/%q", orders[0].Kind, orders[0].Time)
	}
	if string(got.NotCovered) != "[]" {
		t.Errorf("not_covered = %s, want []", got.NotCovered)
	}
	if string(got.Incomplete) != "[]" {
		t.Errorf("incomplete = %s, want []", got.Incomplete)
	}
}

// T-33.2-03: the file holds the person's data and no credential. Every
// secret the database keeps about the account and the order is searched for
// in the JSON text, so a field added to a source later cannot carry one past
// this test.
func TestTheExportCarriesNoSecret(t *testing.T) {
	_, database := newTestStore(t)
	a := website(t, database, "A")
	userID := findUser(t, database, "ada@example.com", "Ada Lovelace", "editor", false)
	findExec(t, database, `INSERT INTO user_recovery_codes (user_id, code_hash) VALUES ($1, 'RECOVERY-HASH-1'), ($1, 'RECOVERY-HASH-2')`, userID)
	findExec(t, database, `INSERT INTO user_tokens (user_id, token_hash, purpose, expires_at)
		VALUES ($1, 'TOKEN-HASH-RESET', 'reset', '2030-01-01T00:00:00Z')`, userID)
	number := findOrder(t, database, a, "2026-0601", "ADA@example.com")

	secrets := map[string]string{}
	var password, totp, pending, token string
	if err := database.Read.QueryRow(`SELECT password, totp_secret, totp_pending_secret FROM users WHERE id = $1`, userID).
		Scan(&password, &totp, &pending); err != nil {
		t.Fatal(err)
	}
	if err := database.Read.QueryRow(`SELECT token FROM orders WHERE number = $1`, number).Scan(&token); err != nil {
		t.Fatal(err)
	}
	secrets["password"] = password
	secrets["totp_secret"] = totp
	secrets["totp_pending_secret"] = pending
	secrets["recovery code"] = "RECOVERY-HASH-1"
	secrets["second recovery code"] = "RECOVERY-HASH-2"
	secrets["token hash"] = "TOKEN-HASH-RESET"
	secrets["order token"] = token

	res, err := NewService(database).Find(context.Background(), a, "ada@example.com")
	if err != nil {
		t.Fatalf("Find: %v", err)
	}
	raw, err := json.Marshal(BuildExport(res, a, "A", time.Now()))
	if err != nil {
		t.Fatal(err)
	}
	text := string(raw)
	for what, secret := range secrets {
		if secret == "" {
			t.Fatalf("the fixture stored no %s; the test would prove nothing", what)
		}
		if strings.Contains(text, secret) {
			t.Errorf("the export carries the %s", what)
		}
	}
	for _, want := range []string{"Ada Lovelace", "ada@example.com", number} {
		if !strings.Contains(text, want) {
			t.Errorf("the export lacks %q", want)
		}
	}
	if len(res.BySource()["accounts"]) != 1 || len(res.BySource()["orders"]) != 1 {
		t.Errorf("the export should hold the account and the order: %s", text)
	}
}

// Art. 15(1) GDPR and Art. 25(2) FADP: the answer carries, beside the data,
// who holds it, why, on what basis, for whom, for how long, where it came
// from and what the person can do. All of it comes from the record of
// processing, so the export and the record cannot disagree.
func TestTheExportDescribesWhatItHolds(t *testing.T) {
	_, database := newTestStore(t)
	a := website(t, database, "Holzwerk")
	findOrder(t, database, a, "2026-0801", "ada@example.com")
	findExec(t, database, `INSERT INTO form_messages (website_id, name, email, subject, body, created_at)
		VALUES ($1, 'Ada', 'ada@example.com', 'Frage', 'Habt ihr Eichenbretter?', '2026-02-01T09:00:00Z')`, a)
	res, err := NewService(database).Find(context.Background(), a, "ada@example.com")
	if err != nil {
		t.Fatalf("Find: %v", err)
	}

	in := recordInput()
	in.Name = "Holzwerk"
	in.Lang = "de"
	in.Set = legal.ParseSet("eu,ch")
	in.Features = privacy.Features{Shop: true, Forms: true}
	rec := BuildRecord(in)

	exp := BuildExport(res, a, "Holzwerk", time.Now())
	exp.Describe(rec)
	raw, err := json.Marshal(exp)
	if err != nil {
		t.Fatal(err)
	}
	var got struct {
		Controller map[string]string `json:"controller"`
		Frameworks []string          `json:"frameworks"`
		Categories []struct {
			Kind       string `json:"kind"`
			Purpose    string `json:"purpose"`
			LegalBasis string `json:"legal_basis"`
			Recipients string `json:"recipients"`
			Retention  string `json:"retention"`
			Source     string `json:"source"`
		} `json:"categories"`
		Rights []struct {
			Framework string `json:"framework"`
			Text      string `json:"text"`
		} `json:"rights"`
		Lang string `json:"lang"`
	}
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("%v\n%s", err, raw)
	}
	if got.Controller["name"] != "Holzwerk" || got.Controller["street"] != in.Street || got.Controller["city"] != in.City {
		t.Errorf("controller %v", got.Controller)
	}
	if strings.Join(got.Frameworks, ",") != "eu,ch" {
		t.Errorf("frameworks %v", got.Frameworks)
	}
	if got.Lang != "de" {
		t.Errorf("lang %q", got.Lang)
	}
	kinds := map[string]int{}
	for i, c := range got.Categories {
		kinds[c.Kind] = i
		for what, cell := range map[string]string{"purpose": c.Purpose, "legal basis": c.LegalBasis, "recipients": c.Recipients, "retention": c.Retention} {
			if strings.TrimSpace(cell) == "" {
				t.Errorf("category %s: the %s is empty", c.Kind, what)
			}
		}
		if c.Source != i18n.T("de", sourceProvided) {
			t.Errorf("category %s: source %q, want the German %q", c.Kind, c.Source, i18n.T("de", sourceProvided))
		}
	}
	if len(got.Categories) != 2 {
		t.Fatalf("categories %+v, want order and form_submission only", got.Categories)
	}
	order, okOrder := kinds["order"]
	_, okForm := kinds["form_submission"]
	if !okOrder || !okForm {
		t.Fatalf("categories %+v, want order and form_submission", got.Categories)
	}
	if !strings.Contains(got.Categories[order].Retention, i18n.T("de", legal.OrderPeriodFor("CH").Basis())) {
		t.Errorf("the order's retention %q does not name the seller's period", got.Categories[order].Retention)
	}
	if got.Categories[order].Purpose != i18n.T("de", purposeOrders) || got.Categories[order].Purpose == purposeOrders {
		t.Errorf("the order's purpose %q is not the German record's", got.Categories[order].Purpose)
	}
	rights := map[string]string{}
	for _, r := range got.Rights {
		rights[r.Framework] += r.Text + " "
	}
	if !strings.Contains(rights["eu"], "Aufsichtsbehörde") || !strings.Contains(rights["eu"], "Art. 15") {
		t.Errorf("EU rights %q do not name access and the supervisory authority in German", rights["eu"])
	}
	if !strings.Contains(rights["ch"], "EDÖB") || !strings.Contains(rights["ch"], "Art. 25") {
		t.Errorf("CH rights %q do not name access and the FDPIC in German", rights["ch"])
	}
	if _, ok := rights["us"]; ok {
		t.Error("rights of a framework the website does not follow")
	}

	// Nothing found: no categories, still the controller and the rights.
	none, err := NewService(database).Find(context.Background(), a, "nobody@example.com")
	if err != nil {
		t.Fatal(err)
	}
	empty := BuildExport(none, a, "Holzwerk", time.Now())
	empty.Describe(rec)
	raw, _ = json.Marshal(empty)
	if !strings.Contains(string(raw), `"categories":[]`) {
		t.Errorf("an empty export: categories are not an empty list: %s", raw)
	}
	if empty.Controller["name"] != "Holzwerk" || len(empty.Rights) == 0 {
		t.Errorf("an empty export lacks the controller or the rights: %+v %+v", empty.Controller, empty.Rights)
	}
}
