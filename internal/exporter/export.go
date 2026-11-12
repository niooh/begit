package exporter

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const rfc3339 = "2006-01-02T15:04:05Z07:00"
const maxFileSize = 1 << 20 // 1MB

type Options struct {
	OutFile   string
	PathsOnly bool
	Ignore    map[string]struct{}
}

type Meta struct {
	Generated time.Time
}

// Renderer abstracts only the parts that differ between output formats.
// Code blocks, tree blocks, binary detection and fence computation are shared.
type Renderer interface {
	Header(m Meta) string
	Heading(level int, text string) string
}

func ExportMarkdown(workdir string) error {
	return export(workdir, Options{OutFile: "pwd_report.md"}, markdownRenderer{})
}

func ExportTypst(workdir string) error {
	return export(workdir, Options{OutFile: "pwd_report.typ"}, typstRenderer{})
}

func export(workdir string, opts Options, r Renderer) error {
	if opts.OutFile == "" {
		return fmt.Errorf("output file is required")
	}
	if opts.Ignore == nil {
		opts.Ignore = defaultIgnoreSet()
	}
	outPath := filepath.Join(workdir, opts.OutFile)
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create report: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	if err := buildDocument(w, workdir, opts, r); err != nil {
		return err
	}
	return w.Flush()
}

func buildDocument(w *bufio.Writer, workdir string, opts Options, r Renderer) error {
	fmt.Fprint(w, r.Header(Meta{Generated: time.Now()}))

	entries, err := collectEntries(workdir, opts.Ignore)
	if err != nil {
		return err
	}

	fmt.Fprint(w, r.Heading(2, "Directory Tree"))
	writeRawBlock(w, "text", renderTree(entries))

	if opts.PathsOnly {
		return nil
	}

	fmt.Fprint(w, r.Heading(2, "Files"))
	for _, e := range entries {
		if e.isDir {
			continue
		}
		fmt.Fprint(w, r.Heading(3, "./"+e.rel))
		writeFileBlock(w, workdir, e.rel)
	}
	return nil
}

func writeFileBlock(w *bufio.Writer, workdir, rel string) {
	abs := filepath.Join(workdir, filepath.FromSlash(rel))
	fi, err := os.Stat(abs)
	if err != nil {
		writeRawBlock(w, "text", "[read failed: "+err.Error()+"]")
		return
	}
	if fi.Size() > maxFileSize {
		writeRawBlock(w, "text", "[skipped: file too large (>1MB)]")
		return
	}
	data, err := os.ReadFile(abs)
	if err != nil {
		writeRawBlock(w, "text", "[read failed: "+err.Error()+"]")
		return
	}
	if looksBinary(data) {
		writeRawBlock(w, "text", "[skipped: binary file]")
		return
	}
	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	ext := strings.ToLower(filepath.Ext(rel))
	writeRawBlock(w, langByExt(ext), content)
}

// writeRawBlock emits a fenced code block. The fence length is chosen so that
// no line inside content can close it early (matters when exporting .md/.mdx).
func writeRawBlock(w *bufio.Writer, lang, content string) {
	fence := computeFence(content)
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	if lang != "" {
		fmt.Fprintf(w, "%s%s\n%s%s\n\n", fence, lang, content, fence)
	} else {
		fmt.Fprintf(w, "%s\n%s%s\n\n", fence, content, fence)
	}
}

// computeFence returns a fence at least one backtick longer than the longest
// pure-backtick line in content.
func computeFence(content string) string {
	maxTicks := 0
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimRight(line, " \t\r")
		if len(line) < 3 {
			continue
		}
		all := true
		for _, ch := range line {
			if ch != '`' {
				all = false
				break
			}
		}
		if all && len(line) > maxTicks {
			maxTicks = len(line)
		}
	}
	if maxTicks == 0 {
		return "```"
	}
	return strings.Repeat("`", maxTicks+1)
}

func looksBinary(data []byte) bool {
	n := len(data)
	if n > 512 {
		n = 512
	}
	for i := 0; i < n; i++ {
		if data[i] == 0 {
			return true
		}
	}
	return false
}

type entry struct {
	rel   string // slash-separated, relative to workdir
	isDir bool
}

func collectEntries(workdir string, ignore map[string]struct{}) ([]entry, error) {
	var out []entry
	var walk func(dir, rel string) error
	walk = func(dir, rel string) error {
		des, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		sort.Slice(des, func(i, j int) bool {
			if des[i].IsDir() != des[j].IsDir() {
				return des[i].IsDir()
			}
			return strings.ToLower(des[i].Name()) < strings.ToLower(des[j].Name())
		})
		for _, de := range des {
			name := de.Name()
			childRel := name
			if rel != "" {
				childRel = rel + "/" + name
			}
			if isIgnored(name, childRel, de.IsDir(), ignore) {
				continue
			}
			if de.IsDir() {
				before := len(out)
				out = append(out, entry{rel: childRel, isDir: true})
				if err := walk(filepath.Join(dir, name), childRel); err != nil {
					return err
				}
				if len(out) == before+1 {
					out = out[:before] // prune empty dirs left by filtering
				}
				continue
			}
			out = append(out, entry{rel: childRel})
		}
		return nil
	}
	return out, walk(workdir, "")
}

func renderTree(entries []entry) string {
	var b strings.Builder
	for _, e := range entries {
		depth := strings.Count(e.rel, "/")
		b.WriteString(strings.Repeat("  ", depth))
		name := filepath.Base(e.rel)
		if e.isDir {
			b.WriteString(name + "/\n")
		} else {
			b.WriteString(name + "\n")
		}
	}
	return b.String()
}

func langByExt(ext string) string {
	m := map[string]string{
		".js": "javascript", ".jsx": "javascript", ".ts": "typescript",
		".tsx": "typescript", ".html": "html", ".py": "python",
		".pl": "perl", ".java": "java", ".go": "go", ".rs": "rust",
		".c": "c", ".h": "c", ".cpp": "cpp", ".hpp": "cpp",
		".cs": "csharp", ".php": "php", ".rb": "ruby", ".kt": "kotlin",
		".swift": "swift", ".sql": "sql", ".sh": "bash", ".bash": "bash",
		".zsh": "zsh", ".yml": "yaml", ".yaml": "yaml", ".json": "json",
		".toml": "toml", ".ini": "ini", ".gradle": "groovy",
	}
	return m[ext]
}
