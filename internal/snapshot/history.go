package snapshot

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"begit/internal/repo"
)

// EndID empty means the path is still present in the latest snapshot.
type Interval struct {
	StartID  string
	StartMsg string
	EndID    string
	EndMsg   string
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


// History returns the contiguous existence intervals of relPath across
func History(r *repo.Repo, relPath string) []Interval {
	var out []Interval
	var open *Interval
	prevExists := false

	for _, rec := range r.Info.Snapshots {
		_, exists := PathHash(filepath.Join(r.SnapshotDir(rec.ID), relPath))
		switch {
		case !prevExists && exists:
			open = &Interval{StartID: rec.ID, StartMsg: rec.Message}
		case prevExists && !exists:
			if open != nil {
				open.EndID = rec.ID
				open.EndMsg = rec.Message
				out = append(out, *open)
				open = nil
			}
		}
		prevExists = exists
	}
	if open != nil {
		out = append(out, *open) // still present at the end
	}
	return out
}
