package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	BegitDirName     = ".__begit__"
	InfoFileName     = "info.json"
	SnapshotsDirName = "snapshots"
)

var ErrNotARepo = errors.New("not a begit repository (or any parent up to root): .__begit__ not found")

type SnapshotRecord struct {
	ID      string    `json:"id"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

type Info struct {
	Version   int              `json:"version"`
	Snapshots []SnapshotRecord `json:"snapshots"`
}

type Repo struct {
	Root     string
	BegitDir string
	Info     *Info
}

func Init(root string) (*Repo, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	begitDir := filepath.Join(root, BegitDirName)
	if _, err := os.Stat(begitDir); err == nil {
		return nil, fmt.Errorf("begit repository already exists at %s", begitDir)
	}
	if err := os.MkdirAll(filepath.Join(begitDir, SnapshotsDirName), 0o755); err != nil {
		return nil, err
	}
	r := &Repo{
		Root:     root,
		BegitDir: begitDir,
		Info:     &Info{Version: 1, Snapshots: []SnapshotRecord{}},
	}
	if err := r.SaveInfo(); err != nil {
		return nil, err
	}
	return r, nil
}

func Open(start string) (*Repo, error) {
	start, err := filepath.Abs(start)
	if err != nil {
		return nil, err
	}
	dir := start
	for {
		begitDir := filepath.Join(dir, BegitDirName)
		if st, err := os.Stat(begitDir); err == nil && st.IsDir() {
			r := &Repo{Root: dir, BegitDir: begitDir}
			if err := r.loadInfo(); err != nil {
				return nil, err
			}
			return r, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return nil, ErrNotARepo
		}
		dir = parent
	}
}

func (r *Repo) loadInfo() error {
	data, err := os.ReadFile(filepath.Join(r.BegitDir, InfoFileName))
	if err != nil {
		return err
	}
	var info Info
	if err := json.Unmarshal(data, &info); err != nil {
		return err
	}
	if info.Version != 1 {
		return fmt.Errorf("unsupported info.json version: %d", info.Version)
	}
	if info.Snapshots == nil {
		info.Snapshots = []SnapshotRecord{}
	}
	r.Info = &info
	return nil
}

func (r *Repo) SaveInfo() error {
	p := filepath.Join(r.BegitDir, InfoFileName)
	data, err := json.MarshalIndent(r.Info, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, p)
}

func (r *Repo) SnapshotsDir() string {
	return filepath.Join(r.BegitDir, SnapshotsDirName)
}

func (r *Repo) SnapshotDir(id string) string {
	return filepath.Join(r.SnapshotsDir(), id)
}
