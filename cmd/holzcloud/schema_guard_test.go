package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/pressly/goose/v3"

	"github.com/holzcloud/holzcloud-cms/internal/config"
	"github.com/holzcloud/holzcloud-cms/internal/db"
)

// newerDatabaseIn leaves a database in dataDir that a newer holzcloud has
// migrated one level beyond this binary: migrated to this binary's head, then
// one more migration applied through a real goose provider. It returns this
// binary's level, read from the database after RunMigrations, so no level is
// written here as a number.
func newerDatabaseIn(t *testing.T, dataDir string) int64 {
	t.Helper()
	database, err := db.Open(filepath.Join(dataDir, "holzcloud.sqlite"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	defer database.Close()
	if err := db.RunMigrations(database.Write); err != nil {
		t.Fatalf("RunMigrations: %v", err)
	}
	highest, err := db.CurrentVersion(database.Write)
	if err != nil {
		t.Fatalf("CurrentVersion: %v", err)
	}
	files := fstest.MapFS{
		fmt.Sprintf("%05d_newer.sql", highest+1): {Data: []byte(
			"-- +goose Up\nCREATE TABLE newer_binary_table (id INTEGER PRIMARY KEY) STRICT;\n\n" +
				"-- +goose Down\nDROP TABLE newer_binary_table;\n")},
	}
	provider, err := goose.NewProvider(goose.DialectSQLite3, database.Write, files)
	if err != nil {
		t.Fatalf("newer provider: %v", err)
	}
	if _, err := provider.Up(context.Background()); err != nil {
		t.Fatalf("apply newer migration: %v", err)
	}
	return highest
}

func levelIn(t *testing.T, dataDir string) int64 {
	t.Helper()
	database, err := db.Open(filepath.Join(dataDir, "holzcloud.sqlite"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	defer database.Close()
	v, err := db.CurrentVersion(database.Write)
	if err != nil {
		t.Fatalf("CurrentVersion: %v", err)
	}
	return v
}

func TestEveryDatabaseSubcommandRefusesANewerDatabase(t *testing.T) {
	dataDir := t.TempDir()
	t.Setenv("HOLZCLOUD_DATA_DIR", dataDir)
	highest := newerDatabaseIn(t, dataDir)

	for name, run := range map[string]func() error{
		"check":          func() error { return cmdCheck(nil) },
		"backup":         func() error { return cmdBackup([]string{t.TempDir()}) },
		"migrate status": func() error { return cmdMigrate([]string{"status"}) },
	} {
		t.Run(name, func(t *testing.T) {
			err := run()
			var newer *db.NewerSchemaError
			if !errors.As(err, &newer) {
				t.Fatalf("%s on a newer database: %v", name, err)
			}
			if newer.DBLevel != highest+1 || newer.BinaryLevel != highest {
				t.Errorf("levels: database %d, binary %d; want %d and %d",
					newer.DBLevel, newer.BinaryLevel, highest+1, highest)
			}
			if strings.HasPrefix(err.Error(), "migrations:") {
				t.Errorf("the refusal is wrapped: %s", err)
			}
			if !strings.HasPrefix(err.Error(), "the database is at migration level ") {
				t.Errorf("the operator does not read the refusal first: %s", err)
			}
		})
	}
	if v := levelIn(t, dataDir); v != highest+1 {
		t.Fatalf("level after the refusals: %d, want %d", v, highest+1)
	}
}

func snapshotsIn(t *testing.T, dataDir string) []string {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(dataDir, "pre-upgrade-*.sqlite"))
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	return matches
}

func TestTheServerRefusesANewerDatabaseBeforeItSnapshots(t *testing.T) {
	dataDir := t.TempDir()
	t.Setenv("HOLZCLOUD_DATA_DIR", dataDir)
	highest := newerDatabaseIn(t, dataDir)
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}

	database, err := openDatabase(cfg)
	if database != nil {
		database.Close()
	}
	var newer *db.NewerSchemaError
	if !errors.As(err, &newer) {
		t.Fatalf("openDatabase on a newer database: %v", err)
	}
	if newer.DBLevel != highest+1 || newer.BinaryLevel != highest {
		t.Errorf("levels: database %d, binary %d; want %d and %d",
			newer.DBLevel, newer.BinaryLevel, highest+1, highest)
	}
	if got := snapshotsIn(t, dataDir); len(got) != 0 {
		t.Errorf("a refused database got a pre-upgrade snapshot: %v", got)
	}
	if v := levelIn(t, dataDir); v != highest+1 {
		t.Errorf("level after the refusal: %d, want %d", v, highest+1)
	}
}

// The refusal of a newer schema names "the pre-upgrade snapshot that the newer
// binary wrote before it migrated" as the way back. A subcommand migrates as
// well, so it writes that snapshot as the service does; otherwise a new binary
// that first touched the database through `user list` or `backup` would leave
// no snapshot to go back to (code review 32, IN-05).
func TestASubcommandSnapshotsBeforeItMigrates(t *testing.T) {
	dataDir := t.TempDir()
	t.Setenv("HOLZCLOUD_DATA_DIR", dataDir)

	if err := cmdCheck(nil); err != nil {
		t.Fatalf("check on an empty data directory: %v", err)
	}
	got := snapshotsIn(t, dataDir)
	if len(got) != 1 {
		t.Fatalf("after a subcommand migrated: %d pre-upgrade snapshots, want 1: %v", len(got), got)
	}
	if v := levelIn(t, dataDir); v == 0 {
		t.Fatal("the subcommand did not migrate")
	}

	// Nothing pending now: no further snapshot.
	if err := cmdCheck(nil); err != nil {
		t.Fatalf("second check: %v", err)
	}
	if again := snapshotsIn(t, dataDir); len(again) != 1 {
		t.Errorf("after a subcommand with nothing pending: %d pre-upgrade snapshots, want 1: %v", len(again), again)
	}
}

// restore.sh runs `check -copy -files` on a throwaway copy of the backup in a
// temporary directory. Migrating that copy needs no way back, and a full
// pre-upgrade snapshot beside it would double the temporary space the restore
// of an older backup needs, on a /tmp that may be a small tmpfs (code review
// 32, second pass, IN-02). Without -copy, check snapshots as before.
func TestCheckOnACopyMigratesWithoutASnapshot(t *testing.T) {
	dataDir := t.TempDir()
	t.Setenv("HOLZCLOUD_DATA_DIR", dataDir)

	if err := cmdCheck([]string{"-copy", "-files", t.TempDir()}); err != nil {
		t.Fatalf("check -copy -files on an empty data directory: %v", err)
	}
	if got := snapshotsIn(t, dataDir); len(got) != 0 {
		t.Errorf("check -copy wrote a pre-upgrade snapshot: %v", got)
	}
	if v := levelIn(t, dataDir); v == 0 {
		t.Fatal("check -copy did not migrate the copy")
	}
}

// A database one level below, with rows in it, is covered by
// TestTheGuardLetsAnOlderDatabasePassAndItMigrates in internal/db: package main
// cannot reach the unexported migration provider that migrates part-way, and no
// export is added for a test.
func TestTheServerStartsOnADatabaseAtOrBelowItsLevel(t *testing.T) {
	dataDir := t.TempDir()
	t.Setenv("HOLZCLOUD_DATA_DIR", dataDir)
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}

	// (a) An empty data directory is level 0, below the binary: the first
	// start writes its small snapshot and migrates, as before.
	database, err := openDatabase(cfg)
	if err != nil {
		t.Fatalf("first start: %v", err)
	}
	highest, err := db.CurrentVersion(database.Write)
	database.Close()
	if err != nil {
		t.Fatalf("CurrentVersion: %v", err)
	}
	if highest == 0 {
		t.Fatal("the first start did not migrate")
	}
	if got := snapshotsIn(t, dataDir); len(got) != 1 {
		t.Fatalf("after the first start: %d pre-upgrade snapshots, want 1: %v", len(got), got)
	}

	// (b) The same directory, now at the binary's own level: no refusal and
	// no further snapshot.
	database, err = openDatabase(cfg)
	if err != nil {
		t.Fatalf("second start at the binary's own level: %v", err)
	}
	database.Close()
	if got := snapshotsIn(t, dataDir); len(got) != 1 {
		t.Fatalf("after the second start: %d pre-upgrade snapshots, want 1: %v", len(got), got)
	}
	if v := levelIn(t, dataDir); v != highest {
		t.Fatalf("level after the second start: %d, want %d", v, highest)
	}
}

// The path snapshotBeforeMigrating returns is the one the way back cites —
// the migration-failure log line, openForCLI's error, the refusal of a newer
// schema. It has to exist when the pruning is done, also for the first
// two-digit release after five earlier snapshots, whose version string sorts
// before theirs (code review 32, second pass, CR-01).
func TestTheSnapshotJustWrittenSurvivesThePruning(t *testing.T) {
	dataDir := t.TempDir()
	t.Setenv("HOLZCLOUD_DATA_DIR", dataDir)
	for _, n := range []string{
		"pre-upgrade-v0.1.5-20200101-000000.sqlite",
		"pre-upgrade-v0.1.6-20200201-000000.sqlite",
		"pre-upgrade-v0.1.7-20200301-000000.sqlite",
		"pre-upgrade-v0.1.8-20200401-000000.sqlite",
		"pre-upgrade-v0.1.9-20200501-000000.sqlite",
	} {
		if err := os.WriteFile(filepath.Join(dataDir, n), []byte("x"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	saved := Version
	Version = "v0.1.10"
	t.Cleanup(func() { Version = saved })

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	database, err := db.Open(cfg.DBPath)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	defer database.Close()
	snapshot, err := snapshotBeforeMigrating(cfg, database)
	if err != nil {
		t.Fatalf("snapshotBeforeMigrating: %v", err)
	}
	if snapshot == "" {
		t.Fatal("no snapshot written for an empty database")
	}
	if _, err := os.Stat(snapshot); err != nil {
		t.Errorf("the snapshot the way back names is gone: %v", err)
	}
	if got := snapshotsIn(t, dataDir); len(got) != 5 {
		t.Errorf("%d snapshots kept, want 5: %v", len(got), got)
	}
	if _, err := os.Stat(filepath.Join(dataDir, "pre-upgrade-v0.1.5-20200101-000000.sqlite")); err == nil {
		t.Error("the oldest snapshot was kept instead of the one just written")
	}
}

// highestLevel is the head of this binary's migrations, read from a scratch
// database after RunMigrations, so no level is written here as a number.
func highestLevel(t *testing.T) int64 {
	t.Helper()
	database, err := db.Open(filepath.Join(t.TempDir(), "scratch.sqlite"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	defer database.Close()
	if err := db.RunMigrations(database.Write); err != nil {
		t.Fatalf("RunMigrations: %v", err)
	}
	highest, err := db.CurrentVersion(database.Write)
	if err != nil {
		t.Fatalf("CurrentVersion: %v", err)
	}
	return highest
}

// migrateUpTo applies this binary's own migration files up to level, through
// a goose provider over the files the binary embeds, read from the source
// tree (the test runs in cmd/holzcloud). It is how part of an update is
// applied without an export from internal/db: the half-migrated database a
// failed update leaves behind.
func migrateUpTo(t *testing.T, database *db.DB, level int64) {
	t.Helper()
	provider, err := goose.NewProvider(goose.DialectSQLite3, database.Write,
		os.DirFS(filepath.Join("..", "..", "internal", "db", "migrations")))
	if err != nil {
		t.Fatalf("goose provider over the migration files: %v", err)
	}
	if _, err := provider.UpTo(context.Background(), level); err != nil {
		t.Fatalf("migrate up to %d: %v", level, err)
	}
	if got, err := db.CurrentVersion(database.Write); err != nil || got != level {
		t.Fatalf("level after migrating up to %d: %d (%v)", level, got, err)
	}
}

// plantSnapshots writes files that look like the snapshots of earlier
// updates, older than anything written during the test.
func plantSnapshots(t *testing.T, dataDir string, names ...string) {
	t.Helper()
	for _, n := range names {
		if err := os.WriteFile(filepath.Join(dataDir, n), []byte("x"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

// markerIn reads the marker that names the snapshot of an unfinished
// migration: its base name and the level it holds. ok is false when there is
// none.
func markerIn(t *testing.T, dataDir string) (name, level string, ok bool) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dataDir, "pre-upgrade.pending"))
	if errors.Is(err, os.ErrNotExist) {
		return "", "", false
	}
	if err != nil {
		t.Fatalf("read marker: %v", err)
	}
	fields := strings.Fields(string(data))
	if len(fields) != 2 || strings.Contains(strings.TrimSpace(string(data)), "\n") {
		t.Fatalf("the marker is not one line of a name and a level: %q", data)
	}
	return fields[0], fields[1], true
}

// afterTheStampOf waits until the clock has left the second in the name of
// snapshot. Snapshot names carry the second, so two snapshots written within
// one second of a test would collide, which the service's restart loop
// (RestartSec=5) and an operator's rollback never do.
func afterTheStampOf(t *testing.T, snapshot string) {
	t.Helper()
	const layout = "20060102-150405"
	base := strings.TrimSuffix(filepath.Base(snapshot), ".sqlite")
	if len(base) < len(layout) {
		return
	}
	stamp := base[len(base)-len(layout):]
	for i := 0; i < 40 && time.Now().UTC().Format(layout) <= stamp; i++ {
		time.Sleep(50 * time.Millisecond)
	}
}

// restorable copies snapshot into a directory of its own, opens it, checks its
// integrity and returns its level: what an operator who follows "Rolling back"
// puts in place.
func restorable(t *testing.T, snapshot string) int64 {
	t.Helper()
	data, err := os.ReadFile(snapshot)
	if err != nil {
		t.Fatalf("read snapshot: %v", err)
	}
	copyPath := filepath.Join(t.TempDir(), "holzcloud.sqlite")
	if err := os.WriteFile(copyPath, data, 0o600); err != nil {
		t.Fatal(err)
	}
	database, err := db.Open(copyPath)
	if err != nil {
		t.Fatalf("open the snapshot's copy: %v", err)
	}
	defer database.Close()
	if result, err := db.QuickCheck(context.Background(), database.Read); err != nil || result != "ok" {
		t.Fatalf("the snapshot fails its integrity check: %q (%v)", result, err)
	}
	level, err := db.CurrentVersion(database.Write)
	if err != nil {
		t.Fatalf("CurrentVersion of the snapshot: %v", err)
	}
	return level
}

func sameFiles(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	seen := map[string]bool{}
	for _, p := range a {
		seen[filepath.Clean(p)] = true
	}
	for _, p := range b {
		if !seen[filepath.Clean(p)] {
			return false
		}
	}
	return true
}

// failedAttemptIn leaves dataDir as the first failed start of an update leaves
// it: three snapshots of earlier updates, the snapshot A this update wrote at
// level highest-2 before its first migration, and the database at highest-1,
// because the migration before the failing one was committed in its own
// transaction while the last one (highest) is still pending. It returns the
// config, the open database, A and highest.
func failedAttemptIn(t *testing.T, dataDir string) (config.Config, *db.DB, string, int64) {
	t.Helper()
	t.Setenv("HOLZCLOUD_DATA_DIR", dataDir)
	highest := highestLevel(t)
	plantSnapshots(t, dataDir,
		"pre-upgrade-v0.1.2-20200101-000000.sqlite",
		"pre-upgrade-v0.1.3-20200201-000000.sqlite",
		"pre-upgrade-v0.1.4-20200301-000000.sqlite",
	)
	saved := Version
	Version = "v9.9.9-test"
	t.Cleanup(func() { Version = saved })

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	database, err := db.Open(cfg.DBPath)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(database.Close)
	migrateUpTo(t, database, highest-2)

	snapshotA, err := snapshotBeforeMigrating(cfg, database)
	if err != nil || snapshotA == "" {
		t.Fatalf("first attempt: snapshot %q, err %v", snapshotA, err)
	}
	name, level, ok := markerIn(t, dataDir)
	if !ok {
		t.Fatal("the first attempt left no marker naming its snapshot")
	}
	if name != filepath.Base(snapshotA) || level != fmt.Sprint(highest-2) {
		t.Fatalf("marker: %s at level %s, want %s at level %d", name, level, filepath.Base(snapshotA), highest-2)
	}

	// The migration before the failing one is committed; the last one fails
	// and stays pending.
	migrateUpTo(t, database, highest-1)
	return cfg, database, snapshotA, highest
}

// F-1 of phase 34: under Restart=on-failure every restart after a failed
// migration wrote another snapshot, of the half-migrated database, and pruned
// to five, so the newest snapshot and every later "restore with:" line named a
// copy the previous binary refuses, and a second start loop pruned the one
// good snapshot away. Five failed starts must leave the snapshot written
// before the first attempt, restorable, and cited by every restart.
func TestRestartsAfterAFailedMigrationKeepTheSnapshotBeforeIt(t *testing.T) {
	dataDir := t.TempDir()
	cfg, database, snapshotA, highest := failedAttemptIn(t, dataDir)
	before := snapshotsIn(t, dataDir)
	if len(before) != 4 {
		t.Fatalf("after the first attempt: %d snapshots, want 4: %v", len(before), before)
	}

	for start := 2; start <= 5; start++ {
		got, err := snapshotBeforeMigrating(cfg, database)
		if err != nil {
			t.Fatalf("start %d: %v", start, err)
		}
		if got != snapshotA {
			t.Errorf("start %d cites %s, want the snapshot before the update %s", start, got, snapshotA)
		}
		if now := snapshotsIn(t, dataDir); !sameFiles(now, before) {
			t.Fatalf("start %d changed the snapshots: %v, want %v", start, now, before)
		}
	}
	if level := restorable(t, snapshotA); level != highest-2 {
		t.Errorf("the snapshot before the update holds level %d, want %d", level, highest-2)
	}
	if name, level, ok := markerIn(t, dataDir); !ok || name != filepath.Base(snapshotA) || level != fmt.Sprint(highest-2) {
		t.Errorf("marker after five starts: %q %q %v", name, level, ok)
	}
}

// A database put back to the level of the kept snapshot may hold writes the
// snapshot lacks (the operator rolled back and worked on), so the next
// migrating start snapshots it again instead of citing the old file.
func TestARolledBackDatabaseGetsAFreshSnapshot(t *testing.T) {
	dataDir := t.TempDir()
	cfg, database, snapshotA, _ := failedAttemptIn(t, dataDir)
	database.Close()

	data, err := os.ReadFile(snapshotA)
	if err != nil {
		t.Fatal(err)
	}
	for _, suffix := range []string{"-wal", "-shm"} {
		if err := os.Remove(cfg.DBPath + suffix); err != nil && !errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(cfg.DBPath, data, 0o600); err != nil {
		t.Fatal(err)
	}
	rolledBack, err := db.Open(cfg.DBPath)
	if err != nil {
		t.Fatalf("open the rolled-back database: %v", err)
	}
	defer rolledBack.Close()

	afterTheStampOf(t, snapshotA)
	got, err := snapshotBeforeMigrating(cfg, rolledBack)
	if err != nil {
		t.Fatalf("start after the rollback: %v", err)
	}
	if got == "" || got == snapshotA {
		t.Fatalf("the start after the rollback cites %q, want a fresh snapshot", got)
	}
	if _, err := os.Stat(got); err != nil {
		t.Errorf("the fresh snapshot is missing: %v", err)
	}
	if name, _, ok := markerIn(t, dataDir); !ok || name != filepath.Base(got) {
		t.Errorf("marker names %q (%v), want %s", name, ok, filepath.Base(got))
	}
}

// A migration that ran through ends the attempt the marker belongs to: the
// next update, even by a build with the same version string, writes a
// snapshot of its own instead of citing one that lacks everything since.
func TestAFinishedMigrationEndsTheKeptSnapshot(t *testing.T) {
	dataDir := t.TempDir()
	t.Setenv("HOLZCLOUD_DATA_DIR", dataDir)
	highest := highestLevel(t)
	saved := Version
	Version = "dev"
	t.Cleanup(func() { Version = saved })
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	database, err := db.Open(cfg.DBPath)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	defer database.Close()
	migrateUpTo(t, database, highest-2)

	s1, err := snapshotBeforeMigrating(cfg, database)
	if err != nil || s1 == "" {
		t.Fatalf("first update: snapshot %q, err %v", s1, err)
	}
	if _, _, ok := markerIn(t, dataDir); !ok {
		t.Fatal("no marker after the first update's snapshot")
	}

	// A binary whose head is highest-1 migrates through and cleans up.
	migrateUpTo(t, database, highest-1)
	forgetPreUpgradeSnapshot(cfg)
	if name, level, ok := markerIn(t, dataDir); ok {
		t.Fatalf("the marker survived a finished migration: %s %s", name, level)
	}

	afterTheStampOf(t, s1)
	s2, err := snapshotBeforeMigrating(cfg, database)
	if err != nil {
		t.Fatalf("next update: %v", err)
	}
	if s2 == "" || s2 == s1 {
		t.Fatalf("the next update cites %q, want a new snapshot beside %s", s2, s1)
	}
	if level := restorable(t, s2); level != highest-1 {
		t.Errorf("the new snapshot holds level %d, want %d", level, highest-1)
	}
}

// markerCase is a marker file's content and the files planted beside it.
type markerCase struct {
	marker string
	plant  func(t *testing.T, dataDir string)
}

// plantedAs is a marker naming a regular file that exists in the data
// directory, so only the name itself can be the reason to ignore it.
func plantedAs(name string) markerCase {
	return markerCase{
		marker: name + " 0\n",
		plant:  func(t *testing.T, dataDir string) { plantSnapshots(t, dataDir, name) },
	}
}

// The marker accepts exactly the shape of name the binary writes, whatever
// version string the build was stamped with.
func TestTheSnapshotNameIsOneItsMarkerAccepts(t *testing.T) {
	at := time.Date(2026, 11, 21, 11, 37, 3, 0, time.UTC)
	for _, version := range []string{"dev", "v0.1.4", "v0.1.4-247-g6986c32", "v0.1.5-beta.1", "v0.1.4-dirty", "v1.0.0+build.7"} {
		name := preUpgradeSnapshotName(version, at)
		if !preUpgradeSnapshotNameRE.MatchString(name) {
			t.Errorf("version %q: the name %q is not one the marker accepts", version, name)
		}
	}
	if got, want := preUpgradeSnapshotName("v0.1.4", at), "pre-upgrade-v0.1.4-20261121-113703.sqlite"; got != want {
		t.Errorf("name %q, want %q", got, want)
	}
}

// The "restore with:" line and the CLI's error are pasted as root. A path of
// plain characters stays as it was (the drivers and DEPLOY.md read that line);
// anything else is quoted so the shell hands cp exactly that path.
func TestTheRestoreLineQuotesWhatItCannotPassBare(t *testing.T) {
	const snap, live = "/opt/holzcloud/data/pre-upgrade-v0.1.4-247-g6986c32-20261121-113703.sqlite", "/opt/holzcloud/data/holzcloud.sqlite"
	if got, want := restoreCommand(snap, live),
		"systemctl stop holzcloud && cp "+snap+" "+live+" && systemctl start holzcloud"; got != want {
		t.Errorf("restoreCommand of plain paths:\n got %s\nwant %s", got, want)
	}
	for _, p := range []string{
		"/d/pre-upgrade-x;touch pwned-20200301-000000.sqlite",
		"/d/pre-upgrade-$(id)-20200301-000000.sqlite",
		"/d/pre-upgrade-`id`-20200301-000000.sqlite",
		"/d/pre-upgrade-a|sh-20200301-000000.sqlite",
		"/d d/pre-upgrade-a'b\"c\\d&e*f?g~h-20200301-000000.sqlite",
		"/d/${IFS}x",
		"",
	} {
		word := shellWord(p)
		out, err := exec.Command("sh", "-c", "printf '%s' "+word).CombinedOutput()
		if err != nil {
			t.Fatalf("sh on %s: %v (%s)", word, err, out)
		}
		if string(out) != p {
			t.Errorf("the shell reads %s as %q, want %q", word, out, p)
		}
		if line := restoreCommand(p, live); !strings.Contains(line, " cp "+word+" ") {
			t.Errorf("restoreCommand does not quote %q: %s", p, line)
		}
	}
}

// The marker is read back at every start and its name ends up in the
// "restore with:" line an operator pastes as root, so a name that is not a
// snapshot in the data directory, or a level that is not a number, is
// ignored: a new snapshot is written and the marker rewritten (T-34-01).
func TestAMarkerThatNamesNoSnapshotIsIgnored(t *testing.T) {
	highest := highestLevel(t)
	const planted = "pre-upgrade-v0.1.4-20200301-000000.sqlite"
	for name, c := range map[string]markerCase{
		"missing file": {marker: planted + " 0\n"},
		"slash": {
			marker: "sub/" + planted + " 0\n",
			plant: func(t *testing.T, dataDir string) {
				if err := os.MkdirAll(filepath.Join(dataDir, "sub"), 0o700); err != nil {
					t.Fatal(err)
				}
				plantSnapshots(t, filepath.Join(dataDir, "sub"), planted)
			},
		},
		"dot dot": {
			marker: "../" + planted + " 0\n",
			plant:  func(t *testing.T, dataDir string) { plantSnapshots(t, filepath.Dir(dataDir), planted) },
		},
		"dot dot in a base name": {
			marker: "pre-upgrade-a..b-20200301-000000.sqlite 0\n",
			plant: func(t *testing.T, dataDir string) {
				plantSnapshots(t, dataDir, "pre-upgrade-a..b-20200301-000000.sqlite")
			},
		},
		"not a snapshot name": {marker: "holzcloud.sqlite 0\n"},
		"level not a number": {
			marker: planted + " seventy\n",
			plant:  func(t *testing.T, dataDir string) { plantSnapshots(t, dataDir, planted) },
		},
		"no level": {
			marker: planted + "\n",
			plant:  func(t *testing.T, dataDir string) { plantSnapshots(t, dataDir, planted) },
		},
		// The name is pasted as root in "restore with:", so a file whose
		// name carries shell syntax is not trusted even when it exists and
		// matches pre-upgrade-*.sqlite (CR-01 of the phase 34 review).
		"semicolon":            plantedAs("pre-upgrade-x;touch${IFS}pwned-20200301-000000.sqlite"),
		"command substitution": plantedAs("pre-upgrade-$(id)-20200301-000000.sqlite"),
		"backtick":             plantedAs("pre-upgrade-`id`-20200301-000000.sqlite"),
		"pipe":                 plantedAs("pre-upgrade-a|sh-20200301-000000.sqlite"),
		"ampersand":            plantedAs("pre-upgrade-a&b-20200301-000000.sqlite"),
		"single quote":         plantedAs("pre-upgrade-a'b-20200301-000000.sqlite"),
		"space":                plantedAs("pre-upgrade-a b-20200301-000000.sqlite"),
		"no time stamp":        plantedAs("pre-upgrade-v0.1.4.sqlite"),
		"short time stamp":     plantedAs("pre-upgrade-v0.1.4-2020031-000000.sqlite"),
	} {
		t.Run(name, func(t *testing.T) {
			dataDir := filepath.Join(t.TempDir(), "data")
			if err := os.MkdirAll(dataDir, 0o700); err != nil {
				t.Fatal(err)
			}
			t.Setenv("HOLZCLOUD_DATA_DIR", dataDir)
			cfg, err := config.Load()
			if err != nil {
				t.Fatalf("config.Load: %v", err)
			}
			database, err := db.Open(cfg.DBPath)
			if err != nil {
				t.Fatalf("db.Open: %v", err)
			}
			defer database.Close()
			migrateUpTo(t, database, highest-1)
			if c.plant != nil {
				c.plant(t, dataDir)
			}
			if err := os.WriteFile(filepath.Join(dataDir, "pre-upgrade.pending"), []byte(c.marker), 0o600); err != nil {
				t.Fatal(err)
			}

			got, err := snapshotBeforeMigrating(cfg, database)
			if err != nil {
				t.Fatalf("snapshotBeforeMigrating: %v", err)
			}
			if got == "" || filepath.Dir(got) != filepath.Clean(dataDir) ||
				strings.HasPrefix(c.marker, filepath.Base(got)+" ") {
				t.Fatalf("cites %q for the marker %q, want a new snapshot in %s", got, c.marker, dataDir)
			}
			if _, err := os.Stat(got); err != nil {
				t.Fatalf("the new snapshot is missing: %v", err)
			}
			n, level, ok := markerIn(t, dataDir)
			if !ok || n != filepath.Base(got) || level != fmt.Sprint(highest-1) {
				t.Errorf("marker after the new snapshot: %q %q %v, want %s %d", n, level, ok, filepath.Base(got), highest-1)
			}
		})
	}
}

// A fresh snapshot that gets no marker of its own must not leave an older
// marker standing: the older marker survives a rollback to a binary before
// 0.1.5, and if this update then fails part-way every restart would cite that
// older snapshot, which lacks everything written since (WR-01 of the phase 34
// review). Both ways a snapshot is written without a marker are played: the
// level cannot be read, and the marker cannot be written.
func TestASnapshotWithoutAMarkerOfItsOwnEndsTheOlderMarker(t *testing.T) {
	for name, fail := range map[string]func(t *testing.T){
		"level unreadable": func(t *testing.T) {
			saved := schemaLevel
			schemaLevel = func(*sql.DB) (int64, error) { return 0, errors.New("database is locked") }
			t.Cleanup(func() { schemaLevel = saved })
		},
		"marker not written": func(t *testing.T) {
			saved := writeMarker
			writeMarker = func(string, string, int64) error { return errors.New("no space left on device") }
			t.Cleanup(func() { writeMarker = saved })
		},
	} {
		t.Run(name, func(t *testing.T) {
			dataDir := t.TempDir()
			t.Setenv("HOLZCLOUD_DATA_DIR", dataDir)
			highest := highestLevel(t)
			saved := Version
			Version = "v9.9.9-test"
			t.Cleanup(func() { Version = saved })
			cfg, err := config.Load()
			if err != nil {
				t.Fatalf("config.Load: %v", err)
			}
			database, err := db.Open(cfg.DBPath)
			if err != nil {
				t.Fatalf("db.Open: %v", err)
			}
			defer database.Close()
			migrateUpTo(t, database, highest-2)

			// A marker an earlier update left behind at this same level,
			// as a rollback to an older binary leaves it.
			const stale = "pre-upgrade-v0.1.4-20200301-000000.sqlite"
			plantSnapshots(t, dataDir, stale)
			if err := os.WriteFile(filepath.Join(dataDir, "pre-upgrade.pending"),
				[]byte(fmt.Sprintf("%s %d\n", stale, highest-2)), 0o600); err != nil {
				t.Fatal(err)
			}

			sub := t.Run("the update's first start", func(t *testing.T) {
				fail(t)
				fresh, err := snapshotBeforeMigrating(cfg, database)
				if err != nil || fresh == "" || filepath.Base(fresh) == stale {
					t.Fatalf("first start: snapshot %q, err %v, want a fresh one", fresh, err)
				}
				if n, l, ok := markerIn(t, dataDir); ok {
					t.Errorf("a marker outlived the fresh snapshot %s: %s %s", filepath.Base(fresh), n, l)
				}
				afterTheStampOf(t, fresh)
			})
			if !sub {
				return
			}

			// The update fails part-way; the restart must not cite the
			// stale snapshot.
			migrateUpTo(t, database, highest-1)
			got, err := snapshotBeforeMigrating(cfg, database)
			if err != nil {
				t.Fatalf("restart: %v", err)
			}
			if got == "" || filepath.Base(got) == stale {
				t.Fatalf("the restart cites %q, want anything but the stale %s", got, stale)
			}
		})
	}
}

// The fixed release meets the database a failed release left half-migrated:
// it finishes the migrations without a snapshot of the half-migrated state and
// removes the marker, through the server's start and through a subcommand
// alike.
func TestAHalfMigratedDatabaseIsFinishedWithoutANewSnapshot(t *testing.T) {
	t.Run("server", func(t *testing.T) {
		dataDir := t.TempDir()
		cfg, database, _, highest := failedAttemptIn(t, dataDir)
		database.Close()
		before := snapshotsIn(t, dataDir)

		started, err := openDatabase(cfg)
		if err != nil {
			t.Fatalf("openDatabase on the half-migrated database: %v", err)
		}
		started.Close()
		if v := levelIn(t, dataDir); v != highest {
			t.Errorf("level after the start: %d, want %d", v, highest)
		}
		if now := snapshotsIn(t, dataDir); !sameFiles(now, before) {
			t.Errorf("the start changed the snapshots: %v, want %v", now, before)
		}
		if name, level, ok := markerIn(t, dataDir); ok {
			t.Errorf("the marker survived the finished migration: %s %s", name, level)
		}
	})
	t.Run("subcommand", func(t *testing.T) {
		dataDir := t.TempDir()
		_, database, _, highest := failedAttemptIn(t, dataDir)
		database.Close()
		before := snapshotsIn(t, dataDir)

		if err := cmdCheck(nil); err != nil {
			t.Fatalf("check on the half-migrated database: %v", err)
		}
		if v := levelIn(t, dataDir); v != highest {
			t.Errorf("level after check: %d, want %d", v, highest)
		}
		if now := snapshotsIn(t, dataDir); !sameFiles(now, before) {
			t.Errorf("check changed the snapshots: %v, want %v", now, before)
		}
		if name, level, ok := markerIn(t, dataDir); ok {
			t.Errorf("the marker survived the finished migration: %s %s", name, level)
		}
	})
}
