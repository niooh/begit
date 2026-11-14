package main

import (
	"os"

	"begit/internal/cli"
	"begit/internal/commands"
)

func main() {
	commands.RegisterAll()
	os.Exit(cli.Run(os.Args[1:]))
}
