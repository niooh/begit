package snapshot

import (
	"io"
	"os"
	"path/filepath"
)

// copy src to dst
func copyTree(src, dst string, skip func(rel string, isDir bool) bool) (int, error) {
	count := 0
	err := filepath.Walk(src, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		isDir := info.IsDir()
		if skip != nil && skip(rel, isDir) {
			if isDir {
				return filepath.SkipDir
			}
			return nil
		}
		target := filepath.Join(dst, rel)
		switch {
		case info.Mode()&os.ModeSymlink != 0:
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			_ = os.Remove(target)
			return os.Symlink(link, target)
		case isDir:
			return os.MkdirAll(target, info.Mode().Perm())
		case info.Mode().IsRegular():
			if err := copyFile(path, target, info.Mode().Perm()); err != nil {
				return err
			}
			count++
			return nil
		default:
			return nil
		}
	})
	return count, err
}

func copyFile(src, dst string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
