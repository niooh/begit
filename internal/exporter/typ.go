package exporter

import "strings"

type typstRenderer struct{}

func (typstRenderer) Header(m Meta) string {
	return "= Report\n\n- Generated: `" + m.Generated.Format(rfc3339) + "`\n\n"
}

func (typstRenderer) Heading(level int, text string) string {
	return strings.Repeat("=", level) + " " + text + "\n\n"
}
