package snapshot

import (
	"os"
	"path/filepath"

	"begit/internal/repo"
)

func RemoveAll(r *repo.Repo, relPath string) (int, error) {
	relPath = filepath.Clean(relPath)
	affected := 0
	for _, rec := range r.Info.Snapshots {
		snapRoot := r.SnapshotDir(rec.ID)
		target := filepath.Join(snapRoot, relPath)
		if _, err := os.Lstat(target); err != nil {
			continue
		}
		if err := os.RemoveAll(target); err != nil {
			return affected, err
		}
		affected++
		cleanupEmptyParents(filepath.Dir(target), snapRoot)
	}
	return affected, nil
}

func cleanupEmptyParents(dir, stop string) {
	for {
		if dir == stop || dir == "." || dir == "/" || dir == "" {
			return
		}
		entries, err := os.ReadDir(dir)
		if err != nil || len(entries) > 0 {
			return
		}
		if err := os.Remove(dir); err != nil {
			return
		}
		dir = filepath.Dir(dir)
	}
}
