package cli

import "os"

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	red    = "\033[31m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	cyan   = "\033[36m"
)

// colorEnabled is false when NO_COLOR is set or stdout is not a terminal.
var colorEnabled = detectColor()

func detectColor() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func paint(code, s string) string {
	if !colorEnabled {
		return s
	}
	return code + s + reset
}

func Bold(s string) string   { return paint(bold, s) }
func Red(s string) string    { return paint(red, s) }
func Yellow(s string) string { return paint(yellow, s) }
func Blue(s string) string   { return paint(blue, s) }
func Cyan(s string) string   { return paint(cyan, s) }
