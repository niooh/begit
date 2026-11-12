package main

import (
	"os"

	"begit/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
