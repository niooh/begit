package cli

import (
	"fmt"
	"os"

	"begit/internal/repo"
)

func cmdLS(args []string) int {
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
	printList(r)
	return 0
}
