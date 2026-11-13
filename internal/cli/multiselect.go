package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"begit/internal/repo"
)

func cmdMultiSelect(args []string) int {
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
	if len(r.Info.Snapshots) < 2 {
		fmt.Println("need at least two snapshots for multi-select")
		return 0
	}

	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 4096), 1024*1024)

	for {
		if len(r.Info.Snapshots) < 2 {
			fmt.Println("need at least two snapshots for multi-select")
			return 0
		}
		printList(r)
		fmt.Print("\nSelect new-old (e.g. 1-3, q to quit): ")
		if !sc.Scan() {
			return 0
		}
		input := strings.TrimSpace(sc.Text())
		if input == "" || input == "q" {
			return 0
		}

		parts := strings.Split(input, "-")
		if len(parts) != 2 {
			fmt.Println("please enter in the form new-old (e.g. 1-3)")
			continue
		}
		newIdx, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		oldIdx, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err1 != nil || err2 != nil {
			fmt.Println("invalid number(s)")
			continue
		}

		n := len(r.Info.Snapshots)
		if newIdx < 1 || newIdx > n || oldIdx < 1 || oldIdx > n {
			fmt.Println("invalid number(s)")
			continue
		}
		if newIdx == oldIdx {
			fmt.Println("please pick two different snapshots")
			continue
		}
		// 1 = newest, so new should have a smaller index than old
		if newIdx > oldIdx {
			fmt.Println("new must be more recent than old (lower number = newer)")
			continue
		}

		// display index -> position (1 = newest)
		newPos := n - newIdx
		oldPos := n - oldIdx

		handleMultiSelect(sc, r, oldPos, newPos)

		// removal may have changed info.json
		if r2, err := repo.Open(r.Root); err == nil {
			r.Info = r2.Info
		}
	}
}

func handleMultiSelect(sc *bufio.Scanner, r *repo.Repo, oldPos, newPos int) {
	idOld := r.Info.Snapshots[oldPos].ID
	idNew := r.Info.Snapshots[newPos].ID
	fmt.Printf("Selected new=%s old=%s\n", Blue(idNew), Blue(idOld))
	for {
		fmt.Println("\n  " + Yellow("(d)") + " diff  " + Yellow("(x)") +
			" remove range  " + Yellow("(b)") + " back")
		fmt.Print("> ")
		if !sc.Scan() {
			return
		}
		switch strings.TrimSpace(sc.Text()) {
		case "d":
			// doDiff(r, cur, prev): show changes prev -> cur
			doDiff(r, idNew, idOld)
		case "x":
			selected := r.Info.Snapshots[oldPos : newPos+1]
			ids := make([]string, len(selected))
			for i, rec := range selected {
				ids[i] = rec.ID
			}
			if confirm(sc, fmt.Sprintf(
				"Remove %d snapshot(s) from %s to %s? [y/N]: ",
				len(ids), Blue(idOld), Blue(idNew))) {
				if err := removeSnapshots(r, ids...); err != nil {
					fmt.Fprintln(os.Stderr, "error:", err)
				} else {
					fmt.Printf("removed %d snapshot(s)\n", len(ids))
				}
				return
			}
		case "b", "q":
			return
		}
	}
}

func removeSnapshots(r *repo.Repo, ids ...string) error {
	drop := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		if err := os.RemoveAll(r.SnapshotDir(id)); err != nil {
			return err
		}
		drop[id] = struct{}{}
	}
	kept := r.Info.Snapshots[:0]
	for _, rec := range r.Info.Snapshots {
		if _, ok := drop[rec.ID]; !ok {
			kept = append(kept, rec)
		}
	}
	r.Info.Snapshots = kept
	return r.SaveInfo()
}
