package snapshot

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"begit/internal/repo"
)

const IDLayout = "20060102_150405000"

var ErrNoChanges = errors.New("no changes since last snapshot")

type Snapshot struct {
	ID      string
	Dir     string
	Message string
	Time    time.Time
}

func Create(r *repo.Repo, msg string) (*Snapshot, error) {
	id, dir, err := allocateID(r)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	if _, err := copyTree(r.Root, dir, isExcluded); err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	if prev := lastSnapshot(r); prev != nil {
		same, err := sameAs(prev.Dir, dir)
		if err != nil {
			os.RemoveAll(dir)
			return nil, err
		}
		if same {
			os.RemoveAll(dir)
			return nil, ErrNoChanges
		}
	}
	now := time.Now()
	rec := repo.SnapshotRecord{ID: id, Message: msg, Time: now}
	r.Info.Snapshots = append(r.Info.Snapshots, rec)
	if err := r.SaveInfo(); err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	return &Snapshot{ID: id, Dir: dir, Message: msg, Time: now}, nil
}

func allocateID(r *repo.Repo) (string, string, error) {
	for i := 0; i < 20; i++ {
		id := time.Now().Format(IDLayout)
		dir := r.SnapshotDir(id)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return id, dir, nil
		}
		time.Sleep(time.Millisecond)
	}
	return "", "", fmt.Errorf("failed to allocate unique snapshot id after retries")
}

func lastSnapshot(r *repo.Repo) *Snapshot {
	if len(r.Info.Snapshots) == 0 {
		return nil
	}
	rec := r.Info.Snapshots[len(r.Info.Snapshots)-1]
	return &Snapshot{ID: rec.ID, Dir: r.SnapshotDir(rec.ID), Message: rec.Message, Time: rec.Time}
}

type meta struct {
	kind byte // 'f' regular, 'd' dir, 'l' symlink
	size int64
}

func sameAs(a, b string) (bool, error) {
	am, err := walkMeta(a)
	if err != nil {
		return false, err
	}
	bm, err := walkMeta(b)
	if err != nil {
		return false, err
	}
	if len(am) != len(bm) {
		return false, nil
	}
	for k, v := range am {
		w, ok := bm[k]
		if !ok || v.kind != w.kind || v.size != w.size {
			return false, nil
		}
		if v.kind == 'f' {
			eq, err := filesEqual(filepath.Join(a, k), filepath.Join(b, k))
			if err != nil {
				return false, err
			}
			if !eq {
				return false, nil
			}
		}
	}
	return true, nil
}

func walkMeta(root string) (map[string]meta, error) {
	m := make(map[string]meta)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		var k byte
		switch {
		case info.Mode()&os.ModeSymlink != 0:
			k = 'l'
		case info.IsDir():
			k = 'd'
		default:
			k = 'f'
		}
		m[rel] = meta{kind: k, size: info.Size()}
		return nil
	})
	return m, err
}

func filesEqual(a, b string) (bool, error) {
	ha, err := hashFile(a)
	if err != nil {
		return false, err
	}
	hb, err := hashFile(b)
	if err != nil {
		return false, err
	}
	return bytes.Equal(ha, hb), nil
}

func hashFile(p string) ([]byte, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
