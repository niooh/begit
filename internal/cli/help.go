package cli

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func runHelp(args []string) int {
	if len(args) == 0 {
		printOverview()
		return 0
	}
	name := args[0]
	if name == "help" {
		fmt.Println(Bold("Usage:") + " help [command]")
		fmt.Println(Bold("Desc:") + "  Show help for a command, or overview when no command is given")
		return 0
	}
	c, ok := commands[name]
	if !ok {
		fmt.Fprintf(os.Stderr, "%s unknown command: %s\n", Red("error:"), name)
		return 1
	}
	fmt.Println(Bold("Usage:") + " " + c.Usage)
	fmt.Println(Bold("Desc:") + "  " + c.Desc)
	if len(c.VerboseHelp) > 0 {
		fmt.Println()
		fmt.Println(Bold("Verbose:"))
		keys := make([]string, 0, len(c.VerboseHelp))
		for k := range c.VerboseHelp {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("  %s%s%s\n", Blue(k), spacesAfter(len(k), 12), c.VerboseHelp[k])
		}
	}
	return 0
}

func printOverview() {
	fmt.Println(Bold("begit") + ", a simple snapshot-based file manager")
	fmt.Println()
	fmt.Println(Bold("Usage:") + " begit <command> [args]")
	fmt.Println(Bold("Help:") + "  begit help <command>")
	fmt.Println()
	fmt.Println(Bold("Commands:"))
	names := make([]string, 0, len(commands))
	for n := range commands {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		c := commands[n]
		fmt.Printf("  %s%s%s\n", Blue(c.Name), spacesAfter(len(c.Name), 12), c.Desc)
	}
}

// spacesAfter returns the padding needed to reach width visible columns,
func spacesAfter(visibleLen, width int) string {
	if visibleLen >= width {
		return " "
	}
	return strings.Repeat(" ", width-visibleLen+1)
}
