package snapshot

import "path/filepath"

func isExcluded(relPath string, isDir bool) bool {
	base := filepath.Base(relPath)
	if isDir && (base == ".__begit__" || base == ".git") {
		return true
	}
	if base == ".DS_Store" || base == "Thumbs.db" {
		return true
	}
	if !isDir && filepath.Ext(base) == ".swp" {
		return true
	}
	return false
}
