package cli

import (
	"fmt"
	"os"
	"sort"
)

func runHelp(args []string) int {
	if len(args) == 0 {
		printOverview()
		return 0
	}
	name := args[0]
	if name == "help" {
		fmt.Println("Usage: help [command]")
		fmt.Println("Desc:  Show help for a command, or overview when no command is given")
		return 0
	}
	c, ok := commands[name]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", name)
		return 1
	}
	fmt.Printf("Usage: %s\n", c.Usage)
	fmt.Printf("Desc:  %s\n", c.Desc)
	if len(c.VerboseHelp) > 0 {
		fmt.Println("\nOptions:")
		keys := make([]string, 0, len(c.VerboseHelp))
		for k := range c.VerboseHelp {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("  %-18s %s\n", k, c.VerboseHelp[k])
		}
	}
	return 0
}

func printOverview() {
	fmt.Println("begit - a simple snapshot-based file manager")
	fmt.Println()
	fmt.Println("Usage: begit <command> [args]")
	fmt.Println("Help:  begit help <command>")
	fmt.Println()
	fmt.Println("Commands:")
	names := make([]string, 0, len(commands))
	for n := range commands {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		c := commands[n]
		fmt.Printf("  %-12s %s\n", c.Name, c.Desc)
	}
}
