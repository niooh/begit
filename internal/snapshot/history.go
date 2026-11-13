package snapshot

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"begit/internal/repo"
)

type HistoryEntry struct {
	ID      string
	Message string
	Kind    string // "added", "modified", "deleted"
}

// PathHash returns a content hash of the file/folder at p, and whether it exists.
// For a folder, the hash covers every relative path plus content of each entry.
func PathHash(p string) (string, bool) {
	fi, err := os.Lstat(p)
	if err != nil {
		return "", false
	}
	h := sha256.New()

	if fi.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(p)
		if err != nil {
			return "", false
		}
		h.Write([]byte("l:" + target))
		return hex.EncodeToString(h.Sum(nil)), true
	}

	if !fi.IsDir() {
		f, err := os.Open(p)
		if err != nil {
			return "", false
		}
		defer f.Close()
		if _, err := io.Copy(h, f); err != nil {
			return "", false
		}
		return hex.EncodeToString(h.Sum(nil)), true
	}

	err = filepath.Walk(p, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(p, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if info.Mode()&os.ModeSymlink != 0 {
			h.Write([]byte("l:" + rel + ":"))
			target, err := os.Readlink(path)
			if err != nil {
				return err
			}
			h.Write([]byte(target + "\n"))
			return nil
		}
		if info.IsDir() {
			h.Write([]byte("d:" + rel + "\n"))
			return nil
		}
		h.Write([]byte("f:" + rel + ":"))
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := io.Copy(h, f); err != nil {
			return err
		}
		h.Write([]byte{'\n'})
		return nil
	})
	if err != nil {
		return "", false
	}
	return hex.EncodeToString(h.Sum(nil)), true
}

// History walks snapshots oldest→newest and returns only entries where
// the content at relPath actually changed. Consecutive snapshots with
// identical content are collapsed into a single row.
func History(r *repo.Repo, relPath string) []HistoryEntry {
	var out []HistoryEntry
	var prevHash string
	prevExists := false
	first := true

	for _, rec := range r.Info.Snapshots {
		h, exists := PathHash(filepath.Join(r.SnapshotDir(rec.ID), relPath))
		kind := ""
		if first {
			if exists {
				kind = "added"
			}
			first = false
		} else {
			switch {
			case !prevExists && exists:
				kind = "added"
			case prevExists && !exists:
				kind = "deleted"
			case exists && h != prevHash:
				kind = "modified"
			}
		}
		if kind != "" {
			out = append(out, HistoryEntry{ID: rec.ID, Message: rec.Message, Kind: kind})
		}
		prevHash = h
		prevExists = exists
	}
	return out
}
