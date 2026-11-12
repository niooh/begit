package cli

import (
	"fmt"
	"os"

	"begit/internal/exporter"
)

func cmdExportMD(args []string) int {
	return runExport("pwd_report.md", exporter.ExportMarkdown)
}

func cmdExportTyp(args []string) int {
	return runExport("pwd_report.typ", exporter.ExportTypst)
}

func runExport(out string, fn func(string) error) int {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	if err := fn(cwd); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	fmt.Println("report generated: " + out)
	return 0
}
