package cli

import (
	"errors"
	"fmt"
	"os"

	"begit/internal/repo"
	"begit/internal/snapshot"
)

func cmdSend(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, `usage: begit send "message"`)
		return 1
	}
	msg := args[0]
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
	s, err := snapshot.Create(r, msg)
	if err != nil {
		if errors.Is(err, snapshot.ErrNoChanges) {
			fmt.Println("no changes since last snapshot")
			return 0
		}
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	fmt.Printf("[%s] %s\n", s.ID, s.Message)
	return 0
}
