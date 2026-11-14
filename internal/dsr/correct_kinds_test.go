package dsr

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

// Correction by the rule of each record kind (Phase 33.2, DSR-02).

func TestAnInvoiceFieldIsCorrectedWithItsOriginalKept(t *testing.T) {
	_, database := newTestStore(t)
	ctx := context.Background()
	a := website(t, database, "A")
	findOrder(t, database, a, "2026-0001", "ada@example.com")
	s := NewService(database)

	for _, tc := range []struct{ field, old, value, want string }{
		{"street", "Bahnhofstrasse 1", "Neue Gasse 2", "Neue Gasse 2"},
		{"vat_number", "", "CHE-123.456.789 MWST", "CHE-123.456.789 MWST"},
		{"country", "CH", "de", "DE"},
		{"name", "Ada Lovelace", "Ada King", "Ada King"},
	} {
		ch, err := s.Correct(ctx, a, "ada@example.com", "orders", "2026-0001", tc.field, tc.value)
		if err != nil {
			t.Fatalf("%s: %v", tc.field, err)
		}
		if !ch.Invoice || ch.Old != tc.old || ch.New != tc.want {
			t.Errorf("%s: change %+v, want old %q, new %q, invoice", tc.field, ch, tc.old, tc.want)
		}
		if got, _ := correcttestOrder(t, database, a, "2026-0001", tc.field); got != tc.want {
			t.Errorf("%s = %q, want %q", tc.field, got, tc.want)
		}
	}
	// What the checkout refuses, the correction refuses.
	for _, tc := range []struct {
		field, value string
		want         error
	}{
		{"city", "", ErrEmptyValue},
		{"country", "Schweiz", ErrBadValue},
		{"street", strings.Repeat("x", 201), ErrTooLong},
	} {
		if _, err := s.Correct(ctx, a, "ada@example.com", "orders", "2026-0001", tc.field, tc.value); !errors.Is(err, tc.want) {
			t.Errorf("%s = %q: %v, want %v", tc.field, tc.value, err, tc.want)
		}
	}
	if got, _ := correcttestOrder(t, database, a, "2026-0001", "city"); got != "Zürich" {
		t.Errorf("a refused value changed the city: %q", got)
	}
}

func TestAnErasedOrderCannotBeCorrected(t *testing.T) {
	_, database := newTestStore(t)
	ctx := context.Background()
	a := website(t, database, "A")
	findOrder(t, database, a, "2026-0001", "ada@example.com")
	// Erased, and the address left in place so the search still finds it:
	// the guard has to be the erasure, not the missing address.
	findExec(t, database, `UPDATE orders SET erased_at = '2026-11-01T00:00:00Z' WHERE number = '2026-0001'`)
	s := NewService(database)

	res, err := s.Find(ctx, a, "ada@example.com")
	if err != nil || len(res.Records) != 1 {
		t.Fatalf("find: %v, %+v", err, res)
	}
	if !res.Records[0].Erased {
		t.Error("the erased order is not marked erased")
	}
	if f := Correctable(res.Records[0]); len(f) != 0 {
		t.Errorf("an erased order offers %+v", f)
	}
	for _, field := range []string{"email", "street"} {
		if _, err := s.Correct(ctx, a, "ada@example.com", "orders", "2026-0001", field, "x@example.com"); !errors.Is(err, ErrErased) {
			t.Errorf("%s: %v, want ErrErased", field, err)
		}
	}
	if got, _ := correcttestOrder(t, database, a, "2026-0001", "street"); got != "Bahnhofstrasse 1" {
		t.Errorf("the erased order changed: %q", got)
	}
}

func TestALegacyFormMessageIsCorrected(t *testing.T) {
	_, database := newTestStore(t)
	ctx := context.Background()
	a := website(t, database, "A")
	id := findExec(t, database, `INSERT INTO form_messages (website_id, name, email, subject, body, created_at)
		VALUES ($1, 'Ada', 'ada@example.com', 'Frage', 'Habt ihr Eichenbretter?', '2026-02-01T09:00:00Z')`, a)
	key := itoa(id)
	s := NewService(database)

	ch, err := s.Correct(ctx, a, "ada@example.com", "form_messages", key, "name", "Ada Lovelace")
	if err != nil || ch.Old != "Ada" || ch.New != "Ada Lovelace" || ch.Invoice {
		t.Fatalf("name: %+v, %v", ch, err)
	}
	if _, err := s.Correct(ctx, a, "ada@example.com", "form_messages", key, "email", "ada@lovelace.example"); err != nil {
		t.Fatalf("email: %v", err)
	}
	if _, err := s.Correct(ctx, a, "ada@lovelace.example", "form_messages", key, "body", "Etwas anderes"); !errors.Is(err, ErrNotCorrectable) {
		t.Errorf("body: %v, want ErrNotCorrectable", err)
	}
	var name, email, body string
	if err := database.Read.QueryRow(`SELECT name, email, body FROM form_messages WHERE id = $1`, id).Scan(&name, &email, &body); err != nil {
		t.Fatal(err)
	}
	if name != "Ada Lovelace" || email != "ada@lovelace.example" || body != "Habt ihr Eichenbretter?" {
		t.Errorf("row = %q, %q, %q", name, email, body)
	}
}

func TestAPluginRecordIsCorrectedThroughTheHook(t *testing.T) {
	_, database := newTestStore(t)
	ctx := context.Background()
	a := website(t, database, "A")
	s, _, seed, contact := erasetestPlugins(t, database, a)
	msgKey, msg := expiretestMessage(time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC), "aaaaaaaa", 0)
	seed(contact, a, msgKey, msg)
	const farm = "bestellung:20260301T090000-cccccccc"
	seed("bestellung", a, farm, map[string]any{
		"id": "20260301T090000-cccccccc", "eingegangen": "2026-03-01T09:00:00Z",
		"name": "Ada Lovelace", "email": "ada@example.com", "telefon": "+41 79 000 00 00",
		"adresse": "Bahnhofstrasse 1", "waehrung": "CHF",
	})
	const erased = "bestellung:20260302T090000-dddddddd"
	seed("bestellung", a, erased, map[string]any{
		"id": "20260302T090000-dddddddd", "eingegangen": "2026-03-02T09:00:00Z",
		"name": "Ada Lovelace", "email": "ada@example.com", "adresse": "Bahnhofstrasse 1",
		"waehrung": "CHF", "geloescht": "2026-11-01T00:00:00Z", "aufbewahren_bis": "2035-01-01T00:00:00Z",
	})

	ch, err := s.Correct(ctx, a, "ada@example.com", "plugin:"+contact, msgKey, "name", "Ada L.")
	if err != nil || ch.Old != "Ada Lovelace" || ch.New != "Ada L." || ch.Invoice {
		t.Fatalf("contact-form name: %+v, %v", ch, err)
	}
	if _, err := s.Correct(ctx, a, "ada@example.com", "plugin:"+contact, msgKey, "text", "anders"); !errors.Is(err, ErrNotCorrectable) {
		t.Errorf("contact-form text: %v, want ErrNotCorrectable before the plugin is asked", err)
	}
	ch, err = s.Correct(ctx, a, "ada@example.com", "plugin:bestellung", farm, "phone", "+41 79 111 11 11")
	if err != nil || ch.Old != "+41 79 000 00 00" || ch.Invoice {
		t.Fatalf("farm-shop phone: %+v, %v", ch, err)
	}
	ch, err = s.Correct(ctx, a, "ada@example.com", "plugin:bestellung", farm, "address", "Neue Gasse 2")
	if err != nil || ch.Old != "Bahnhofstrasse 1" || !ch.Invoice {
		t.Fatalf("farm-shop address: %+v, %v (the address an order keeps is kept on record)", ch, err)
	}
	if _, err := s.Correct(ctx, a, "ada@example.com", "plugin:bestellung", erased, "phone", "+41 79 111 11 11"); !errors.Is(err, ErrErased) {
		t.Errorf("erased farm-shop order: %v, want ErrErased", err)
	}

	res, err := s.Find(ctx, a, "ada@example.com")
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]string{}
	for _, rec := range res.Records {
		for _, f := range rec.Fields {
			seen[rec.Key+"/"+f.Name] = f.Value
		}
	}
	if seen[msgKey+"/name"] != "Ada L." {
		t.Errorf("the message's name is %q after the correction", seen[msgKey+"/name"])
	}
	if seen[farm+"/phone"] != "+41 79 111 11 11" || seen[farm+"/address"] != "Neue Gasse 2" {
		t.Errorf("the farm-shop order reads %q / %q after the correction", seen[farm+"/phone"], seen[farm+"/address"])
	}
}

func TestLogsAndAccountsAreNotCorrectedHere(t *testing.T) {
	_, database := newTestStore(t)
	ctx := context.Background()
	a := website(t, database, "A")
	findSeed(t, database, a, "2026-0001", "ada@example.com")
	findUser(t, database, "ada@example.com", "Ada", "admin", false)
	s := NewService(database)

	res, err := s.Find(ctx, a, "ada@example.com")
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, rec := range res.Records {
		switch rec.Source {
		case "requests", "mail", "accounts":
			seen[rec.Source] = true
			if f := Correctable(rec); len(f) != 0 {
				t.Errorf("%s %s offers %+v", rec.Source, rec.Key, f)
			}
			if _, err := s.Correct(ctx, a, "ada@example.com", rec.Source, rec.Key, "name", "Ada King"); !errors.Is(err, ErrNotCorrectable) {
				t.Errorf("%s %s: %v, want ErrNotCorrectable", rec.Source, rec.Key, err)
			}
		}
	}
	if len(seen) != 3 {
		t.Fatalf("the seed did not reach every kind: %v", seen)
	}
}
