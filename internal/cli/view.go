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

	intervals := snapshot.History(r, relPath)
	if len(intervals) == 0 {
		fmt.Printf("no history for %q\n", relPath)
		return 0
	}

	fmt.Printf("History for %q:\n", relPath)
	for i, iv := range intervals {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("  %s %s  %s  %s\n",
			Green("+"), Blue(iv.StartID), formatID(iv.StartID), iv.StartMsg)

		if iv.EndID == "" {
			fmt.Printf("  %s still present\n", Cyan("·"))
			continue
		}
		fmt.Printf("  %s %s  %s  %s\n",
			Red("-"), Blue(iv.EndID), formatID(iv.EndID), iv.EndMsg)
	}
	return 0
}

func formatID(id string) string {
	if t, err := repo.ParseSnapshotID(id); err == nil {
		return t.Format("2006-01-02 15:04:05")
	}
	return ""
}
