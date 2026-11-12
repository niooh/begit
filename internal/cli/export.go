package cli

import (
	"fmt"
	"os"

	"begit/internal/exporter"
)

func cmdExportMD(args []string) int {
	if err := exporter.ExportMarkdown("."); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	return 0
}

func cmdExportTyp(args []string) int {
	if err := exporter.ExportTypst("."); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	return 0
}
