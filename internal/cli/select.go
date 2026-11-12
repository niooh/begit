package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"begit/internal/repo"
	"begit/internal/snapshot"
)

func cmdSelect(args []string) int {
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

	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 4096), 1024*1024)

	for {
		if len(r.Info.Snapshots) == 0 {
			fmt.Println("no snapshots left")
			return 0
		}
		printList(r)
		fmt.Print("Select number (q to quit): ")
		if !sc.Scan() {
			return 0
		}
		input := strings.TrimSpace(sc.Text())
		if input == "" || input == "q" {
			return 0
		}
		n, err := strconv.Atoi(input)
		if err != nil || n < 1 || n > len(r.Info.Snapshots) {
			fmt.Println("invalid number")
			continue
		}
		pos := len(r.Info.Snapshots) - n
		rec := r.Info.Snapshots[pos]
		var prevID string
		if pos > 0 {
			prevID = r.Info.Snapshots[pos-1].ID
		}
		handleSnapshot(sc, r, rec.ID, prevID)

		// `remove` operation may have modified info.json, reload it
		if r2, err := repo.Open(r.Root); err == nil {
			r.Info = r2.Info
		}
	}
}

func printList(r *repo.Repo) {
	fmt.Println("Snapshots (newest first):")
	for i := len(r.Info.Snapshots) - 1; i >= 0; i-- {
		rec := r.Info.Snapshots[i]
		idx := len(r.Info.Snapshots) - i
		fmt.Printf("  %s %s  %s  %s\n", Bold(fmt.Sprintf("%d", idx)),
		Blue(rec.ID), rec.Time.Format("2006-01-02 15:04:05"), rec.Message)
	}
}

func handleSnapshot(sc *bufio.Scanner, r *repo.Repo, id, prevID string) {
	for {
		fmt.Printf("Selected [%s]\n", id)
		fmt.Println("  (d) diff    (r) reset    (x) remove    (b) back")
		fmt.Print("> ")
		if !sc.Scan() {
			return
		}
		switch strings.TrimSpace(sc.Text()) {
		case "d":
			doDiff(r, id, prevID)
		case "r":
			if confirm(sc, "Reset working directory to this snapshot? [y/N]: ") {
				if err := snapshot.ResetTo(r.SnapshotDir(id), r.Root); err != nil {
					fmt.Fprintln(os.Stderr, "error:", err)
				} else {
					fmt.Println("reset done")
				}
			}
		case "x":
			if confirm(sc, "Remove this snapshot permanently? [y/N]: ") {
				if err := removeSnapshot(r, id); err != nil {
					fmt.Fprintln(os.Stderr, "error:", err)
				} else {
					fmt.Printf("removed [%s]\n", id)
				}
				return
			}
		case "b", "q":
			return
		}
	}
}

func doDiff(r *repo.Repo, id, prevID string) {
	var prevDir string
	if prevID != "" {
		prevDir = r.SnapshotDir(prevID)
	}
	changes, err := snapshot.Diff(prevDir, r.SnapshotDir(id))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return
	}
	if len(changes) == 0 {
		fmt.Println("(no changes)")
		return
	}
	for _, c := range changes {
		fmt.Printf("  %c %s\n", c.Kind, c.Path)
	}
}

func removeSnapshot(r *repo.Repo, id string) error {
	if err := os.RemoveAll(r.SnapshotDir(id)); err != nil {
		return err
	}
	kept := r.Info.Snapshots[:0]
	for _, rec := range r.Info.Snapshots {
		if rec.ID != id {
			kept = append(kept, rec)
		}
	}
	r.Info.Snapshots = kept
	return r.SaveInfo()
}

func confirm(sc *bufio.Scanner, prompt string) bool {
	fmt.Print(prompt)
	if !sc.Scan() {
		return false
	}
	ans := strings.ToLower(strings.TrimSpace(sc.Text()))
	return ans == "y" || ans == "yes"
}
