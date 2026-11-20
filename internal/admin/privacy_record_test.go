package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/holzcloud/holzcloud-cms/internal/auth"
	"github.com/holzcloud/holzcloud-cms/internal/dsr"
	"github.com/holzcloud/holzcloud-cms/internal/user"
)

// The record of processing activities (Phase 33.2, plan 10).

func recordScreen(websiteID int64) string {
	return privacyScreen(websiteID) + "/verzeichnis"
}

func TestTheRecordScreenAndItsExports(t *testing.T) {
	h, sm, database, ws := privacyAdmin(t)
	if _, err := database.Write.Exec(`UPDATE websites SET shop_base = 'shop', legal_frameworks = 'eu,ch', country = 'CH',
		street = 'Dorfstrasse 1', postal_code = '3000', city = 'Bern' WHERE id = $1`, ws.ID); err != nil {
		t.Fatal(err)
	}

	// The screen.
	rec := privacyGet(t, h, sm, h.HandlePrivacyRecord, recordScreen(ws.ID), ws.ID)
	if rec.Code != http.StatusOK {
		t.Fatalf("screen: status %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{
		optoutSays("Shop orders"), optoutSays("Delivering the website (server log)"), optoutSays("Backups"),
		optoutSays("Purpose"), optoutSays("Legal basis"), optoutSays("Recipients"), optoutSays("Retention"),
		optoutSays("Security measures"), optoutSays("Exemptions"),
		`href="` + recordScreen(ws.ID) + `/drucken"`, `href="` + recordScreen(ws.ID) + `.json"`,
		"Dorfstrasse 1",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("the screen lacks %q", want)
		}
	}
	// The privacy tabs need the navigation middleware's website, which this
	// handler test does not run; the browser proof shows the tab.

	// The print page: the same activities, no admin navigation.
	print := privacyGet(t, h, sm, h.HandlePrivacyRecordPrint, recordScreen(ws.ID)+"/drucken", ws.ID)
	if print.Code != http.StatusOK {
		t.Fatalf("print: status %d, want 200", print.Code)
	}
	pb := print.Body.String()
	for _, want := range []string{optoutSays("Shop orders"), optoutSays("Backups"), "noindex", optoutSays("Record of processing activities")} {
		if !strings.Contains(pb, want) {
			t.Errorf("the print page lacks %q", want)
		}
	}
	for _, nav := range []string{`class="section-tabs"`, `class="sidebar`, "/admin/logout"} {
		if strings.Contains(pb, nav) {
			t.Errorf("the print page carries the admin's navigation (%s)", nav)
		}
	}

	// The JSON.
	js := privacyGet(t, h, sm, h.HandlePrivacyRecordJSON, recordScreen(ws.ID)+".json", ws.ID)
	if js.Code != http.StatusOK {
		t.Fatalf("json: status %d, want 200", js.Code)
	}
	if cd := js.Header().Get("Content-Disposition"); !strings.HasPrefix(cd, "attachment; filename=") {
		t.Errorf("Content-Disposition = %q, want an attachment", cd)
	}
	if cc := js.Header().Get("Cache-Control"); cc != "no-store" {
		t.Errorf("Cache-Control = %q, want no-store", cc)
	}
	if ct := js.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %q", ct)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(js.Body.Bytes(), &raw); err != nil {
		t.Fatalf("the body is not JSON: %v", err)
	}
	for _, key := range []string{"generated_at", "controller", "frameworks", "activities", "security", "exemptions"} {
		if _, ok := raw[key]; !ok {
			t.Errorf("the JSON has no %q", key)
		}
	}
	var got dsr.ProcessingRecord
	if err := json.Unmarshal(js.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	ids := map[string]bool{}
	for _, a := range got.Activities {
		ids[a.ID] = true
	}
	if !ids["orders"] || !ids["server_log"] || ids["payment"] {
		t.Errorf("activities %v: want orders and server_log, and no payment without Payrexx", ids)
	}
	if got.Controller["street"] != "Dorfstrasse 1" {
		t.Errorf("controller %v", got.Controller)
	}
}

// The assistant's tool reads the record through the handler: the website's
// main language unless asked for another, and the same activities as the
// screen.
func TestTheRecordScreenAndTheAssistantsOpAgree(t *testing.T) {
	h, _, database, ws := privacyAdmin(t)
	if _, err := database.Write.Exec(`UPDATE websites SET locale = 'it', shop_base = 'shop' WHERE id = $1`, ws.ID); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	rec, err := h.OpProcessingRecord(ctx, ws.ID, "")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Lang != "it" {
		t.Errorf("lang %q, want the website's it", rec.Lang)
	}
	ids := map[string]bool{}
	for _, a := range rec.Activities {
		ids[a.ID] = true
	}
	if !ids["orders"] || !ids["server_log"] {
		t.Errorf("activities %v", ids)
	}
	de, err := h.OpProcessingRecord(ctx, ws.ID, "de")
	if err != nil || de.Lang != "de" {
		t.Errorf("asked for de: %q, %v", de.Lang, err)
	}
	if _, err := h.OpProcessingRecord(ctx, ws.ID+99, ""); err == nil {
		t.Error("an unknown website gave a record")
	}
}

func TestTheRecordScreenRefusesAnotherWebsite(t *testing.T) {
	h, sm, database, a := privacyAdmin(t)
	ctx := context.Background()
	b, err := h.domains.CreateWebsite(ctx, "Andere", "")
	if err != nil {
		t.Fatal(err)
	}
	users := user.NewStore(database, cheapHashing)
	editorID, err := users.Create(ctx, "Nur B", "nur-b@example.com", "passwort123", user.RoleEditor)
	if err != nil {
		t.Fatal(err)
	}
	if err := users.SetRights(ctx, editorID, user.Rights{MayPublish: true, Websites: []int64{b.ID}}); err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /admin/websites/{id}/datenschutz/verzeichnis", h.ErrHandler(h.HandlePrivacyRecord))
	mux.HandleFunc("GET /admin/websites/{id}/datenschutz/verzeichnis.json", h.ErrHandler(h.HandlePrivacyRecordJSON))
	mux.HandleFunc("GET /admin/websites/{id}/datenschutz/verzeichnis/drucken", h.ErrHandler(h.HandlePrivacyRecordPrint))
	guarded := auth.RequireWebsiteAccess(sm, NewWebsiteAccessLookup(database))(mux)
	drive := func(target string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		sm.LoadAndSave(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sm.Put(r.Context(), auth.SessionKeyUserID, editorID)
			guarded.ServeHTTP(w, r)
		})).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
		return rec
	}
	for _, suffix := range []string{"", ".json", "/drucken"} {
		if rec := drive(recordScreen(a.ID) + suffix); rec.Code != http.StatusForbidden && rec.Code != http.StatusNotFound {
			t.Errorf("%s on A: status %d for an editor of B only, want a refusal", suffix, rec.Code)
		}
		if rec := drive(recordScreen(b.ID) + suffix); rec.Code != http.StatusOK {
			t.Errorf("%s on B: status %d, want 200", suffix, rec.Code)
		}
	}
	if rec := drive("/admin/websites/" + strconv.FormatInt(a.ID+99, 10) + "/datenschutz/verzeichnis"); rec.Code == http.StatusOK {
		t.Error("an unknown website answered 200")
	}
}
