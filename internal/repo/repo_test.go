package repo

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFormatParseSnapshotID(t *testing.T) {
	tm := time.Date(2026, 11, 14, 23, 32, 9, 123_000_000, time.UTC)
	id := FormatSnapshotID(tm)
	if len(id) != idLen {
		t.Fatalf("id length = %d, want %d (%q)", len(id), idLen, id)
	}
	if id != "20261114_233209123" {
		t.Fatalf("id = %q, want 20261114_233209123", id)
	}
	got, err := ParseSnapshotID(id)
	if err != nil {
		t.Fatalf("ParseSnapshotID: %v", err)
	}
	if !got.Equal(tm) {
		t.Fatalf("roundtrip = %v, want %v", got, tm)
	}
}

func TestParseSnapshotIDInvalid(t *testing.T) {
	for _, c := range []string{
		"",
		"short",
		"20261114_233209",     // missing ms
		"20261114_2332099999", // too long
		"20261114_233209abc",  // non-numeric ms
		"20261114_XXXXXX123",  // bad timestamp
	} {
		if _, err := ParseSnapshotID(c); err == nil {
			t.Errorf("ParseSnapshotID(%q) = nil error, want error", c)
		}
	}
}

func TestInitAndOpen(t *testing.T) {
	root := t.TempDir()
	r, err := Init(root)
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	if r.Root != root {
		t.Errorf("root = %q, want %q", r.Root, root)
	}
	if _, err := os.Stat(filepath.Join(root, BegitDirName, InfoFileName)); err != nil {
		t.Errorf("info.json missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, BegitDirName, SnapshotsDirName)); err != nil {
		t.Errorf("snapshots dir missing: %v", err)
	}
	if _, err := Init(root); err == nil {
		t.Error("Init on existing repo should fail")
	}

	sub := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	got, err := Open(sub)
	if err != nil {
		t.Fatalf("Open from subdir: %v", err)
	}
	if got.Root != root {
		t.Errorf("Open root = %q, want %q", got.Root, root)
	}
}

func TestOpenNoRepo(t *testing.T) {
	_, err := Open(t.TempDir())
	if !errors.Is(err, ErrNotARepo) {
		t.Fatalf("err = %v, want ErrNotARepo", err)
	}
}

func TestSaveInfoRoundtrip(t *testing.T) {
	root := t.TempDir()
	r, err := Init(root)
	if err != nil {
		t.Fatal(err)
	}
	r.Info.Snapshots = append(r.Info.Snapshots, SnapshotRecord{ID: "20261114_233209123", Message: "hi"})
	if err := r.SaveInfo(); err != nil {
		t.Fatal(err)
	}
	r2, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(r2.Info.Snapshots) != 1 || r2.Info.Snapshots[0].Message != "hi" {
		t.Fatalf("roundtrip mismatch: %+v", r2.Info)
	}
}

func TestDirs(t *testing.T) {
	r, err := Init(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(r.BegitDir, SnapshotsDirName); r.SnapshotsDir() != want {
		t.Errorf("SnapshotsDir = %q, want %q", r.SnapshotsDir(), want)
	}
	if want := filepath.Join(r.SnapshotsDir(), "abc"); r.SnapshotDir("abc") != want {
		t.Errorf("SnapshotDir = %q, want %q", r.SnapshotDir("abc"), want)
	}
}
