package cli

import (
	"fmt"
	"os"

	"begit/internal/repo"
)

func cmdInit(args []string) int {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	r, err := repo.Init(cwd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	fmt.Printf("Initialized empty begit repository in %s\n", r.BegitDir)
	return 0
}
