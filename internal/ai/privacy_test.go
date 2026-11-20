package ai

import (
	"context"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/holzcloud/holzcloud-cms/internal/domain"
	"github.com/holzcloud/holzcloud-cms/internal/dsr"
	"github.com/holzcloud/holzcloud-cms/internal/media"
	"github.com/holzcloud/holzcloud-cms/internal/page"
)

// The record of processing over MCP (Phase 33.2, plan 10). What the record
// says is tested in internal/dsr and internal/admin; here it is the wiring:
// which website a key reaches, the language handed on, and that no tool
// hands out a person's data.

// fakePrivacy stands in for the admin handler's OpProcessingRecord.
type fakePrivacy struct {
	asked []string
}

func (f *fakePrivacy) OpProcessingRecord(_ context.Context, websiteID int64, lang string) (dsr.ProcessingRecord, error) {
	f.asked = append(f.asked, lang)
	if lang == "" {
		// The handler's default: the website's main language.
		lang = "fr"
	}
	return dsr.ProcessingRecord{
		GeneratedAt: "2026-11-20T09:00:00Z",
		Controller:  map[string]string{"name": "Eins"},
		Frameworks:  []string{"eu", "ch"},
		Activities:  []dsr.Activity{{ID: "server_log", Name: "Server log"}, {ID: "orders", Name: "Orders"}},
		Security:    []string{"Argon2id"},
		Exemptions:  []dsr.Exemption{{Framework: "ch", Text: "Art. 24"}},
		Lang:        lang,
	}, nil
}

func privacySetUp(t *testing.T, ops any) (*httptest.Server, int64, int64, siteKeys) {
	t.Helper()
	database := newTestDB(t)
	domains := domain.NewStore(database)
	ctx := context.Background()
	first, _ := domains.CreateWebsite(ctx, "Eins", "")
	second, _ := domains.CreateWebsite(ctx, "Zwei", "")
	tokens := NewStore(database)
	var k siteKeys
	var err error
	if k.read, _, err = tokens.Issue(ctx, "lesen", 0, false, 0); err != nil {
		t.Fatal(err)
	}
	if k.onlyFirst, _, err = tokens.Issue(ctx, "zwei", second.ID, false, 0); err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(NewServer(tokens, "Test", slog.New(slog.DiscardHandler), Tools(Deps{
		Domains: domains, Pages: page.NewStore(database), Media: media.NewStore(database), Ops: ops,
	})))
	t.Cleanup(ts.Close)
	return ts, first.ID, second.ID, k
}

func TestGetProcessingRecordAnswersForTheKeysWebsite(t *testing.T) {
	fake := &fakePrivacy{}
	ts, first, _, k := privacySetUp(t, fake)

	got, failed := callTool(t, ts, k.read, "get_processing_record", map[string]any{"website": first})
	if failed {
		t.Fatalf("get_processing_record: %v", got["text"])
	}
	acts, _ := got["activities"].([]any)
	if len(acts) != 2 {
		t.Errorf("activities %v", got["activities"])
	}
	if got["lang"] != "fr" {
		t.Errorf("lang %v, want the website's main language", got["lang"])
	}
	if len(fake.asked) != 1 || fake.asked[0] != "" {
		t.Errorf("handed on lang %q, want empty for the website's default", fake.asked)
	}

	got, failed = callTool(t, ts, k.read, "get_processing_record", map[string]any{"website": first, "lang": "de"})
	if failed || got["lang"] != "de" {
		t.Errorf("lang de: %v", got)
	}

	// A key for the second website asking for the first.
	got, failed = callTool(t, ts, k.onlyFirst, "get_processing_record", map[string]any{"website": first})
	if !failed || !strings.Contains(got["text"].(string), "another website") {
		t.Errorf("a key for another website got %v", got)
	}
	if len(fake.asked) != 2 {
		t.Errorf("a refused call reached the handler: %v", fake.asked)
	}

	// A build without the handler says so.
	bare, first2, _, k2 := privacySetUp(t, nil)
	got, failed = callTool(t, bare, k2.read, "get_processing_record", map[string]any{"website": first2})
	if !failed || !strings.Contains(got["text"].(string), "not available") {
		t.Errorf("without the handler: %v", got)
	}
}

// T-33.2-47, 33.1 D-16 carried over: a person's data stays in the admin.
func TestNoToolReturnsAPersonsData(t *testing.T) {
	tools := Tools(Deps{})
	found := false
	for _, tool := range tools {
		name := strings.ToLower(tool.Name)
		for _, bad := range []string{"person", "personal_data", "privacy_request", "opt_out", "optout", "erase"} {
			if strings.Contains(name, bad) {
				t.Errorf("tool %q names %q: a person's data is not exposed over MCP", tool.Name, bad)
			}
		}
		if tool.Name != "get_processing_record" {
			continue
		}
		found = true
		for arg := range tool.InputSchema.Properties {
			if strings.Contains(strings.ToLower(arg), "mail") || strings.Contains(strings.ToLower(arg), "address") {
				t.Errorf("get_processing_record takes %q: the record is about the website, not a person", arg)
			}
		}
		if tool.Writes || tool.Admin {
			t.Errorf("get_processing_record: writes %v, admin %v; it only reads what holds no person's data", tool.Writes, tool.Admin)
		}
	}
	if !found {
		t.Error("there is no get_processing_record tool")
	}
}
