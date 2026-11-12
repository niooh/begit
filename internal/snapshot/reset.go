package snapshot

import (
	"os"
	"path/filepath"
)

func ResetTo(snapDir, workdir string) error {
	if err := clearTree(workdir); err != nil {
		return err
	}
	_, err := copyTree(snapDir, workdir, nil)
	return err
}

func clearTree(root string) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if isExcluded(e.Name(), e.IsDir()) {
			continue
		}
		if err := os.RemoveAll(filepath.Join(root, e.Name())); err != nil {
			return err
		}
	}
	return nil
}
