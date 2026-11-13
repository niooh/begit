package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"begit/internal/repo"
	"begit/internal/snapshot"
)

func cmdView(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: begit view <relative-path>")
		return 1
	}
	relPath := filepath.Clean(args[0])
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

	history := snapshot.History(r, relPath)
	if len(history) == 0 {
		fmt.Printf("no history for %q\n", relPath)
		return 0
	}

	fmt.Printf("History for %q:\n", relPath)
	for _, h := range history {
		ts := ""
		if t, err := repo.ParseSnapshotID(h.ID); err == nil {
			ts = t.Format("2006-01-02 15:04:05")
		}
		var marker string
		switch h.Kind {
		case "added":
			marker = Green("+")
		case "modified":
			marker = Yellow("*")
		case "deleted":
			marker = Red("-")
		default:
			marker = " "
		}
		fmt.Printf("  %s %s  %s  %s\n", marker, Blue(h.ID), ts, h.Message)
	}
	return 0
}
