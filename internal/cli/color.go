package cli

import "os"

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	green  = "\033[32m"
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

// add
func Green(s string) string { return paint(green, s) }

// error, delete
func Red(s string) string { return paint(red, s) }

// emphasize
func Yellow(s string) string { return paint(yellow, s) }

// emphasize, list
func Blue(s string) string { return paint(blue, s) }

// alternative
func Cyan(s string) string { return paint(cyan, s) }

// slight emphasis
func Bold(s string) string { return paint(bold, s) }
