package exporter

import (
	"path/filepath"
	"strings"
)

var defaultIgnoreDirs = map[string]struct{}{
	".__begit__": {}, ".git": {}, ".svn": {}, ".hg": {},
	"node_modules": {}, "dist": {}, "build": {}, "out": {},
	"coverage": {}, "bin": {}, ".next": {}, ".nuxt": {},
	".cache": {}, "__pycache__": {}, ".venv": {}, "venv": {},
	".vscode": {}, ".idea": {}, ".history": {},
}

var defaultIgnoreFiles = map[string]struct{}{
	"README.md": {}, "README": {}, "LICENSE": {},
	"pwd_report.md": {}, "package-lock.json": {}, "pnpm-lock.yaml": {},
}

var defaultIgnoreExts = map[string]struct{}{
	".exe": {}, ".dll": {}, ".so": {}, ".dylib": {},
	".bin": {}, ".class": {}, ".o": {}, ".a": {},
	".pyc": {}, ".pyo": {}, ".jar": {}, ".war": {},
	".zip": {}, ".tar": {}, ".gz": {}, ".7z": {},
	".png": {}, ".jpg": {}, ".jpeg": {}, ".gif": {},
	".pdf": {}, ".doc": {}, ".docx": {}, ".xls": {},
	".lock": {}, ".swp": {}, ".swo": {}, ".db": {},
}

func defaultIgnoreSet() map[string]struct{} {
	m := make(map[string]struct{}, len(defaultIgnoreDirs)+len(defaultIgnoreFiles))
	for k := range defaultIgnoreDirs {
		m[k] = struct{}{}
	}
	for k := range defaultIgnoreFiles {
		m[k] = struct{}{}
	}
	return m
}

func isIgnored(name, rel string, isDir bool, ignore map[string]struct{}) bool {
	if _, ok := ignore[name]; ok {
		return true
	}
	if _, ok := ignore[rel]; ok {
		return true
	}
	if isDir {
		_, ok := defaultIgnoreDirs[name]
		return ok
	}
	if _, ok := defaultIgnoreFiles[name]; ok {
		return true
	}
	_, ok := defaultIgnoreExts[strings.ToLower(filepath.Ext(name))]
	return ok
}
