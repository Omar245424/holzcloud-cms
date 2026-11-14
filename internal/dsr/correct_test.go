package dsr

import (
	"context"
	"errors"
	"testing"

	"github.com/holzcloud/holzcloud-cms/internal/db"
)

// Correction (Phase 33.2, DSR-02). Every helper of this file starts with
// correcttest, so the tests beside it keep their names.

// correcttestOrder reads one column of an order and its updated_at.
func correcttestOrder(t *testing.T, database *db.DB, websiteID int64, number, column string) (value, updated string) {
	t.Helper()
	// The column is one of the test's own literals, never input.
	if err := database.Read.QueryRow(`SELECT `+column+`, updated_at FROM orders WHERE website_id = $1 AND number = $2`,
		websiteID, number).Scan(&value, &updated); err != nil {
		t.Fatalf("read order %s: %v", number, err)
	}
	return value, updated
}

func TestAnOrdersEmailIsCorrected(t *testing.T) {
	_, database := newTestStore(t)
	ctx := context.Background()
	a := website(t, database, "A")
	b := website(t, database, "B")
	findOrder(t, database, a, "2026-0001", "ada@example.com")
	findOrder(t, database, a, "2026-0002", "bea@example.com")
	findOrder(t, database, b, "2026-0003", "ada@example.com")
	s := NewService(database)

	ch, err := s.Correct(ctx, a, "ADA@example.com", "orders", "2026-0001", "email", "  ada.new@example.com ")
	if err != nil {
		t.Fatalf("Correct: %v", err)
	}
	if ch.Old != "ada@example.com" || ch.New != "ada.new@example.com" || ch.Invoice || ch.Field != "email" ||
		ch.Source != "orders" || ch.Key != "2026-0001" {
		t.Errorf("change = %+v", ch)
	}
	email, updated := correcttestOrder(t, database, a, "2026-0001", "email")
	if email != "ada.new@example.com" {
		t.Errorf("email = %q, want the new address", email)
	}
	if updated == "2026-03-01T10:00:00Z" {
		t.Error("updated_at did not move")
	}
	res, err := s.Find(ctx, a, "ada.new@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Records) != 1 || res.Records[0].Key != "2026-0001" {
		t.Errorf("the new address finds %+v, want the order", res.Records)
	}

	// A value that is not one address changes nothing.
	if _, err := s.Correct(ctx, a, "ada.new@example.com", "orders", "2026-0001", "email", "Ada <ada@example.com>"); !errors.Is(err, ErrBadEmail) {
		t.Errorf("an address with a name: %v, want ErrBadEmail", err)
	}
	if email, _ := correcttestOrder(t, database, a, "2026-0001", "email"); email != "ada.new@example.com" {
		t.Errorf("a refused value changed the order: %q", email)
	}

	// Another customer's order, and the person's order of another website,
	// cannot be reached through this address on this website.
	for _, tc := range []struct {
		website int64
		number  string
	}{{a, "2026-0002"}, {a, "2026-0003"}, {b, "2026-0002"}} {
		if _, err := s.Correct(ctx, tc.website, "ada.new@example.com", "orders", tc.number, "email", "x@example.com"); !errors.Is(err, ErrNotFound) {
			t.Errorf("website %d order %s: %v, want ErrNotFound", tc.website, tc.number, err)
		}
	}
	if email, _ := correcttestOrder(t, database, a, "2026-0002", "email"); email != "bea@example.com" {
		t.Errorf("another customer's order changed: %q", email)
	}
	if email, _ := correcttestOrder(t, database, b, "2026-0003", "email"); email != "ada@example.com" {
		t.Errorf("the order of website B changed: %q", email)
	}

	// The phone is a contact field too.
	ch, err = s.Correct(ctx, a, "ada.new@example.com", "orders", "2026-0001", "phone", "+41 79 111 11 11")
	if err != nil || ch.Invoice || ch.Old != "+41 79 000 00 00" {
		t.Errorf("phone: %+v, %v", ch, err)
	}
}
