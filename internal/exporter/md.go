package exporter

import (
	"fmt"
	"strings"
)

type markdownRenderer struct{}

func (markdownRenderer) Header(m Meta) string {
	return fmt.Sprintf("# Report\n\n- Generated: `%s`\n\n", m.Generated.Format(rfc3339))
}

func (markdownRenderer) Heading(level int, text string) string {
	return strings.Repeat("#", level) + " " + text + "\n\n"
}
