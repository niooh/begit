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
		fmt.Print("\nSelect two numbers separated by space (q to quit): ")
		if !sc.Scan() {
			return 0
		}
		input := strings.TrimSpace(sc.Text())
		if input == "" || input == "q" {
			return 0
		}
		fields := strings.Fields(input)
		if len(fields) != 2 {
			fmt.Println("please enter exactly two numbers")
			continue
		}
		a, err1 := strconv.Atoi(fields[0])
		b, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			fmt.Println("invalid number(s)")
			continue
		}
		n := len(r.Info.Snapshots)
		if a < 1 || a > n || b < 1 || b > n {
			fmt.Println("invalid number(s)")
			continue
		}
		if a == b {
			fmt.Println("please pick two different snapshots")
			continue
		}

		// translate display index -> position (1 = newest)
		posA, posB := n-a, n-b
		oldPos, newPos := posA, posB
		if oldPos > newPos {
			oldPos, newPos = newPos, oldPos
		}
		handleMultiSelect(sc, r, r.Info.Snapshots[oldPos].ID, r.Info.Snapshots[newPos].ID)

		// removal may have changed info.json
		if r2, err := repo.Open(r.Root); err == nil {
			r.Info = r2.Info
		}
	}
}

func handleMultiSelect(sc *bufio.Scanner, r *repo.Repo, idOld, idNew string) {
	fmt.Printf("Selected old=%s new=%s\n", Blue(idOld), Blue(idNew))
	for {
		fmt.Println("\n  " + Yellow("(d)") + " diff  " + Yellow("(x)") +
			" remove both  " + Yellow("(b)") + " back")
		fmt.Print("> ")
		if !sc.Scan() {
			return
		}
		switch strings.TrimSpace(sc.Text()) {
		case "d":
			// doDiff(r, cur, prev): show changes prev -> cur
			doDiff(r, idNew, idOld)
		case "x":
			if confirm(sc, "Remove both selected snapshots permanently? [y/N]: ") {
				if err := removeTwoSnapshots(r, idOld, idNew); err != nil {
					fmt.Fprintln(os.Stderr, "error:", err)
				} else {
					fmt.Printf("removed %s and %s\n", Blue(idOld), Blue(idNew))
				}
				return
			}
		case "b", "q":
			return
		}
	}
}

func removeTwoSnapshots(r *repo.Repo, ids ...string) error {
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
