package cli

import (
	"fmt"
	"os"
)

func Run(args []string) int {
	if len(args) == 0 {
		printOverview()
		return 1
	}
	switch args[0] {
	case "help", "-h", "--help":
		return runHelp(args[1:])
	}
	c, ok := commands[args[0]]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[0])
		printOverview()
		return 1
	}
	return c.Run(args[1:])
}
