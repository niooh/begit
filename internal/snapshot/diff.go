package snapshot

import (
	"path/filepath"
	"sort"
)

type Change struct {
	Kind byte // 'A' Add, 'M' Modify, 'D' Delete
	Path string
}

func Diff(prevDir, curDir string) ([]Change, error) {
	cur, err := walkMeta(curDir)
	if err != nil {
		return nil, err
	}
	prev := map[string]meta{}
	if prevDir != "" {
		prev, err = walkMeta(prevDir)
		if err != nil {
			return nil, err
		}
	}

	var changes []Change
	for p, cm := range cur {
		pm, ok := prev[p]
		if !ok {
			changes = append(changes, Change{Kind: 'A', Path: p})
			continue
		}
		if cm.kind != pm.kind || cm.size != pm.size {
			changes = append(changes, Change{Kind: 'M', Path: p})
			continue
		}
		if cm.kind == 'f' {
			eq, err := filesEqual(filepath.Join(prevDir, p), filepath.Join(curDir, p))
			if err != nil {
				return nil, err
			}
			if !eq {
				changes = append(changes, Change{Kind: 'M', Path: p})
			}
		}
	}
	for p := range prev {
		if _, ok := cur[p]; !ok {
			changes = append(changes, Change{Kind: 'D', Path: p})
		}
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].Path < changes[j].Path })
	return changes, nil
}
