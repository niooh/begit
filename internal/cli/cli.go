package cli

import (
	"fmt"
	"os"
)

func Run(args []string) int {
	if len(args) == 0 {
		usage()
		return 1
	}
	cmd, rest := args[0], args[1:]
	switch cmd {
	case "init":
		return cmdInit(rest)
	case "send":
		return cmdSend(rest)
	case "select":
		return cmdSelect(rest)
	case "remove-all":
		return cmdRemoveAll(rest)
	case "md":
		return cmdExportMD(rest)
	case "typ":
		return cmdExportTyp(rest)
	case "-h", "--help", "help":
		usage()
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		usage()
		return 1
	}
}

func usage() {
	fmt.Print(`begit - a simple snapshot-based file manager

Usage:
  begit init                  initialize a repository in the current directory
  begit send "message"        create a full snapshot of the current directory
  begit select                browse snapshots interactively
  begit remove-all <path>     remove a relative path from all snapshots
  begit md                    export as markdown (not implemented yet)
  begit typ                   export as typst    (not implemented yet)
`)
}
