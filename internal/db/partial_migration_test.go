package db

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log/slog"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/pressly/goose/v3"
)

// Each SQL migration runs in a transaction of its own, so when a later file of
// an update fails, the earlier ones stay committed and the database stands
// between the two levels. Before phase 34 the journal named none of them,
// only the failure (F-2), and the operator could not read the half-migrated
// level from it. The migrations a partial failure applied are logged as
// "migration applied", the same line a successful run writes, before the
// error is returned.
func TestAPartialFailureNamesTheMigrationsItApplied(t *testing.T) {
	database, err := Open(filepath.Join(t.TempDir(), "partial.sqlite"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer database.Close()

	files := fstest.MapFS{
		"00001_ok.sql": {Data: []byte("-- +goose Up\nCREATE TABLE partial_ok (id INTEGER PRIMARY KEY) STRICT;\n\n" +
			"-- +goose Down\nDROP TABLE partial_ok;\n")},
		"00002_fail.sql": {Data: []byte("-- +goose Up\nINSERT INTO partial_no_such_table VALUES (1);\n\n" +
			"-- +goose Down\nSELECT 1;\n")},
	}
	provider, err := goose.NewProvider(goose.DialectSQLite3, database.Write, files)
	if err != nil {
		t.Fatalf("provider: %v", err)
	}

	var buf bytes.Buffer
	saved := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(saved) })

	err = applyMigrations(provider)
	if err == nil {
		t.Fatal("a failing migration returned no error")
	}
	if !strings.Contains(err.Error(), "version:2") {
		t.Errorf("the error does not name the failing migration 2: %v", err)
	}

	var applied []float64
	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		var line map[string]any
		if err := json.Unmarshal(scanner.Bytes(), &line); err != nil {
			t.Fatalf("log line %q: %v", scanner.Text(), err)
		}
		if line["msg"] != "migration applied" {
			continue
		}
		v, _ := line["version"].(float64)
		applied = append(applied, v)
	}
	if len(applied) != 1 || applied[0] != 1 {
		t.Errorf("migration applied lines name versions %v, want exactly [1]", applied)
	}

	level, err := CurrentVersion(database.Write)
	if err != nil {
		t.Fatalf("CurrentVersion: %v", err)
	}
	if level != 1 {
		t.Errorf("level after the partial failure: %d, want 1", level)
	}
}
