package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"begit/internal/repo"
	"begit/internal/snapshot"
)

func cmdRemoveAll(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: begit remove-all <relative-path>")
		return 1
	}
	relPath := args[0]
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	r, err := repo.Open(cwd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	if len(r.Info.Snapshots) == 0 {
		fmt.Println("no snapshots yet")
		return 0
	}

	clean := filepath.Clean(relPath)
	var hits []string
	for _, rec := range r.Info.Snapshots {
		if _, err := os.Lstat(filepath.Join(r.SnapshotDir(rec.ID), clean)); err == nil {
			hits = append(hits, rec.ID)
		}
	}
	if len(hits) == 0 {
		fmt.Printf("no snapshot contains %q\n", clean)
		return 0
	}
	fmt.Printf("Will remove %q from %d snapshot(s):\n", clean, len(hits))
	for _, id := range hits {
		fmt.Printf("  %s\n", Blue(id))
	}
	sc := bufio.NewScanner(os.Stdin)
	if !confirm(sc, "Proceed? [y/N]: ") {
		fmt.Println("aborted")
		return 0
	}

	n, err := snapshot.RemoveAll(r, clean)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	fmt.Printf("removed from %d snapshot(s)\n", n)
	return 0
}
