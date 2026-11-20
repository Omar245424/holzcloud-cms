package media

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
)

// DeleteOwned is how the retention job and the erasure remove the files of a
// record a plugin names (Phase 33.2). Every helper here starts with owned.

// ownedPlugin is the plugin the files of these tests are stored for.
const ownedPlugin = "kontaktformular"

// ownedFile stores a media row and its file on disk, as a plugin form's
// attachment is stored: for ownedPlugin.
func ownedFile(t *testing.T, s *Store, websiteID int64, dataDir, name string) *Media {
	t.Helper()
	return ownedStored(t, s, websiteID, ownedPlugin, dataDir, name)
}

// ownedStored stores a media row and its file on disk for the plugin, ""
// for the operator's own upload to the media library.
func ownedStored(t *testing.T, s *Store, websiteID int64, pluginID, dataDir, name string) *Media {
	t.Helper()
	m, err := s.CreateForPlugin(context.Background(), websiteID, pluginID, name, name, "application/pdf", 4, "hash-"+name)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	dir := filepath.Join(dataDir, "media", strconv.FormatInt(websiteID, 10))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte("%PDF"), 0o644); err != nil {
		t.Fatal(err)
	}
	return m
}

func ownedOnDisk(dataDir string, websiteID int64, name string) bool {
	_, err := os.Stat(filepath.Join(dataDir, "media", strconv.FormatInt(websiteID, 10), name))
	return err == nil
}

func TestDeleteOwnedRefusesAnotherWebsitesFile(t *testing.T) {
	s, a := newTestStore(t)
	ctx := context.Background()
	res, err := s.DB.Write.Exec(`INSERT INTO websites (name, description) VALUES ('B', '')`)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := res.LastInsertId()
	dir := t.TempDir()
	m := ownedFile(t, s, b, dir, "brief.pdf")

	if err := s.DeleteOwned(ctx, a, ownedPlugin, m.ID, dir); !errors.Is(err, ErrNotFound) {
		t.Fatalf("DeleteOwned of B's file for A = %v, want ErrNotFound", err)
	}
	if got, _ := s.GetByID(ctx, m.ID); got == nil {
		t.Error("B's row is gone")
	}
	if !ownedOnDisk(dir, b, "brief.pdf") {
		t.Error("B's file is gone")
	}
	if err := s.DeleteOwned(ctx, a, ownedPlugin, 9999, dir); !errors.Is(err, ErrNotFound) {
		t.Errorf("DeleteOwned of no file = %v, want ErrNotFound", err)
	}
}

func TestDeleteOwnedReportsAFileThatStays(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("a read-only directory does not stop root (or Windows) from removing a file; the case cannot be made here")
	}
	s, ws := newTestStore(t)
	ctx := context.Background()
	dir := t.TempDir()
	m := ownedFile(t, s, ws, dir, "anhang.pdf")
	folder := filepath.Join(dir, "media", strconv.FormatInt(ws, 10))
	if err := os.Chmod(folder, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(folder, 0o755) })

	err := s.DeleteOwned(ctx, ws, ownedPlugin, m.ID, dir)
	if !errors.Is(err, ErrFileRemains) {
		t.Fatalf("DeleteOwned with a file that cannot go = %v, want ErrFileRemains", err)
	}
	if got, _ := s.GetByID(ctx, m.ID); got != nil {
		t.Error("the row is still there; the record it belonged to is gone")
	}
	if !ownedOnDisk(dir, ws, "anhang.pdf") {
		t.Fatal("the premise failed: the file was removed from a read-only directory")
	}

	// A file that is already gone is what was asked for, not an error.
	os.Chmod(folder, 0o755)
	gone := ownedFile(t, s, ws, dir, "weg.pdf")
	if err := os.Remove(filepath.Join(folder, "weg.pdf")); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteOwned(ctx, ws, ownedPlugin, gone.ID, dir); err != nil {
		t.Errorf("DeleteOwned of a file already missing = %v, want nil", err)
	}
}

func TestDeleteOwnedDeletesAFileInUse(t *testing.T) {
	s, ws := newTestStore(t)
	ctx := context.Background()
	dir := t.TempDir()
	m := ownedFile(t, s, ws, dir, "plan.pdf")
	pageID := seedPage(t, s, ws, "Anfahrt", "anfahrt")
	if err := s.ReplaceUsage(ctx, pageID, []int64{m.ID}); err != nil {
		t.Fatalf("ReplaceUsage: %v", err)
	}

	if err := s.DeleteOwned(ctx, ws, ownedPlugin, m.ID, dir); err != nil {
		t.Fatalf("DeleteOwned of a file in use = %v; the record it belonged to expired", err)
	}
	if got, _ := s.GetByID(ctx, m.ID); got != nil {
		t.Error("the row is still there")
	}
	if ownedOnDisk(dir, ws, "plan.pdf") {
		t.Error("the file is still on disk")
	}
}

// A plugin's claim reaches only the files stored for that plugin: the
// operator's own upload — the logo, an image on a page — and another
// plugin's file are refused, and nothing changes.
func TestDeleteOwnedRefusesAFileThePluginDidNotStore(t *testing.T) {
	s, ws := newTestStore(t)
	ctx := context.Background()
	dir := t.TempDir()
	logo := ownedStored(t, s, ws, "", dir, "logo.pdf")
	other := ownedStored(t, s, ws, "bestellung", dir, "fremd.pdf")
	for _, m := range []*Media{logo, other} {
		if err := s.DeleteOwned(ctx, ws, ownedPlugin, m.ID, dir); !errors.Is(err, ErrNotOwned) {
			t.Errorf("DeleteOwned of %s for %s = %v, want ErrNotOwned", m.Filename, ownedPlugin, err)
		}
		if got, _ := s.GetByID(ctx, m.ID); got == nil || !ownedOnDisk(dir, ws, m.Filename) {
			t.Errorf("%s is gone", m.Filename)
		}
	}
	if err := s.DeleteOwned(ctx, ws, "bestellung", other.ID, dir); err != nil {
		t.Errorf("DeleteOwned of the plugin's own file = %v, want nil", err)
	}
}
