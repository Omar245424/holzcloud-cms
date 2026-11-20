package admin

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/holzcloud/holzcloud-cms/internal/activity"
	"github.com/holzcloud/holzcloud-cms/internal/db"
	"github.com/holzcloud/holzcloud-cms/internal/user"
)

// Correction from the Privacy screen (Phase 33.2, plan 09). Every helper of
// this file starts with correct, so the plans beside it in this package keep
// theirs apart.

func correctAddress(websiteID int64) string {
	return "/admin/websites/" + strconv.FormatInt(websiteID, 10) + "/datenschutz/berichtigen"
}

// correctPost sends one correction form.
func correctPost(c *eraseClient, websiteID int64, form url.Values, htmx bool) (rec *http.Response, bad, good string) {
	c.t.Helper()
	req := albumPost(correctAddress(websiteID), form, websiteRoute(websiteID), htmx)
	r, b, g := c.do(c.h.HandlePrivacyCorrect, req)
	return r.Result(), b, g
}

// correctEntries are the privacy.correct entries of the activity log, their
// metadata decoded.
func correctEntries(t *testing.T, database *db.DB) []map[string]any {
	t.Helper()
	rows, err := database.Read.Query(`SELECT metadata FROM activity_log WHERE action = $1 ORDER BY id`, activity.ActionPrivacyCorrect)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var out []map[string]any
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err != nil {
			t.Fatal(err)
		}
		m := map[string]any{}
		if err := json.Unmarshal([]byte(raw), &m); err != nil {
			t.Fatal(err)
		}
		out = append(out, m)
	}
	return out
}

func correctOrderColumn(t *testing.T, database *db.DB, number, column string) string {
	t.Helper()
	var v string
	// The column is one of the test's own literals, never input.
	if err := database.Read.QueryRow(`SELECT `+column+` FROM orders WHERE number = $1`, number).Scan(&v); err != nil {
		t.Fatal(err)
	}
	return v
}

// correctForms are the correction forms of a page: their action, their
// hx-post and the field they correct.
var correctFormRe = regexp.MustCompile(`(?s)<form method="post" action="([^"]*)" hx-post="([^"]*)"[^>]*>(.*?)</form>`)

func correctFormFor(body, field string) (action, hxPost string, ok bool) {
	for _, m := range correctFormRe.FindAllStringSubmatch(body, -1) {
		if strings.Contains(m[3], `name="field" value="`+field+`"`) {
			return m[1], m[2], true
		}
	}
	return "", "", false
}

func TestCorrectingAnEmailFromThePrivacyScreen(t *testing.T) {
	for _, htmx := range []bool{false, true} {
		h, sm, database, ws := eraseAdmin(t)
		privacyOrder(t, database, ws.ID, "2026-0042", "ada@example.com", "Ada Lovelace")
		c := &eraseClient{t: t, h: h, sm: sm, fresh: true}

		body := c.get(h.HandlePrivacyPerson, privacyScreen(ws.ID)+"?email=ada%40example.com", ws.ID).Body.String()
		if !strings.Contains(body, retentionSays("Correct…")) {
			t.Fatal("the screen offers no correction")
		}
		action, hxPost, ok := correctFormFor(body, "email")
		if !ok {
			t.Fatal("the order has no form for its e-mail address")
		}
		if action != correctAddress(ws.ID) || hxPost != action {
			t.Errorf("form action %q, hx-post %q, want both %q", action, hxPost, correctAddress(ws.ID))
		}

		rec, bad, good := correctPost(c, ws.ID, url.Values{
			"email": {"ada@example.com"}, "anfrage": {"7"}, "source": {"orders"}, "key": {"2026-0042"},
			"field": {"email"}, "value": {"ada.new@example.com"},
		}, htmx)
		want := privacyScreen(ws.ID) + "?anfrage=7&email=ada.new%40example.com"
		if htmx {
			if rec.StatusCode != http.StatusOK || rec.Header.Get("HX-Redirect") != want {
				t.Errorf("htmx: %d, HX-Redirect %q, want %q", rec.StatusCode, rec.Header.Get("HX-Redirect"), want)
			}
		} else if rec.StatusCode != http.StatusSeeOther || rec.Header.Get("Location") != want {
			t.Errorf("%d, Location %q, want %q", rec.StatusCode, rec.Header.Get("Location"), want)
		}
		if bad != "" || good != "Corrected." {
			t.Errorf("flash %q / %q, want Corrected.", bad, good)
		}
		if got := correctOrderColumn(t, database, "2026-0042", "email"); got != "ada.new@example.com" {
			t.Errorf("email = %q", got)
		}
		after := c.get(h.HandlePrivacyPerson, want, ws.ID).Body.String()
		if !strings.Contains(after, "2026-0042") {
			t.Error("the record is not found under the new address")
		}

		entries := correctEntries(t, database)
		if len(entries) != 1 {
			t.Fatalf("%d privacy.correct entries, want 1", len(entries))
		}
		e := entries[0]
		if e["source"] != "orders" || e["key"] != "2026-0042" || e["field"] != "email" || e["invoice"] != false {
			t.Errorf("entry %v", e)
		}
		raw, _ := json.Marshal(e)
		if strings.Contains(string(raw), "ada") || len(e) != 4 {
			t.Errorf("the entry of a contact field carries more than source, key, field and invoice: %s", raw)
		}
	}
}

func TestAnEditorWithoutThePublishingRightCannotCorrect(t *testing.T) {
	h, sm, database, ws := eraseAdmin(t)
	privacyOrder(t, database, ws.ID, "2026-0042", "ada@example.com", "Ada Lovelace")
	editor := editorWithoutPublishing(t, database, ws.ID)
	c := &eraseClient{t: t, h: h, sm: sm, user: editor, role: user.RoleEditor, fresh: true}

	req := albumPost(correctAddress(ws.ID), url.Values{
		"email": {"ada@example.com"}, "source": {"orders"}, "key": {"2026-0042"},
		"field": {"email"}, "value": {"evil@example.com"},
	}, websiteRoute(ws.ID), false)
	rec, bad, _ := c.do(h.HandlePrivacyCorrect, req)
	wantRefused(t, "correction", rec, bad, privacyScreen(ws.ID)+"?email=ada%40example.com")
	if got := correctOrderColumn(t, database, "2026-0042", "email"); got != "ada@example.com" {
		t.Errorf("an editor without the publishing right changed the order: %q", got)
	}
	if n := len(correctEntries(t, database)); n != 0 {
		t.Errorf("%d privacy.correct entries after a refusal", n)
	}

	body := c.get(h.HandlePrivacyPerson, privacyScreen(ws.ID)+"?email=ada%40example.com", ws.ID).Body.String()
	if !strings.Contains(body, "2026-0042") {
		t.Fatal("the editor does not see the order")
	}
	if strings.Contains(body, "datenschutz/berichtigen") || strings.Contains(body, retentionSays("Correct…")) {
		t.Error("the screen offers the editor a correction that is refused")
	}
}
