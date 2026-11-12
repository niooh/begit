package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	BegitDirName     = ".__begit__"
	InfoFileName     = "info.json"
	SnapshotsDirName = "snapshots"

	idTimeLayout = "20060102_150405" // seconds only
	idLen        = 18                // 15 + 3-digit milliseconds
)

var ErrNotARepo = errors.New("not a begit repository, as .__begit__ not found")

type SnapshotRecord struct {
	ID      string `json:"id"`
	Message string `json:"message"`
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

// FormatSnapshotID returns a 18-char ID in UTC: "20060102_150405" + 3-digit ms.
func FormatSnapshotID(t time.Time) string {
	return t.UTC().Format(idTimeLayout) + fmt.Sprintf("%03d", t.Nanosecond()/1_000_000)
}

// ParseSnapshotID parses an ID produced by FormatSnapshotID and returns UTC time.
func ParseSnapshotID(id string) (time.Time, error) {
	if len(id) != idLen {
		return time.Time{}, fmt.Errorf("bad snapshot id: %q", id)
	}
	t, err := time.ParseInLocation(idTimeLayout, id[:15], time.UTC)
	if err != nil {
		return time.Time{}, err
	}
	ms, err := strconv.Atoi(id[15:])
	if err != nil || ms > 999 {
		return time.Time{}, fmt.Errorf("bad snapshot id: %q", id)
	}
	return t.Add(time.Duration(ms) * time.Millisecond), nil
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
