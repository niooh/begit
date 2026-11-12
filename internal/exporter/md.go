package exporter

import "errors"

var ErrNotImplemented = errors.New("not implemented yet")

func ExportMarkdown(workdir string) error { return ErrNotImplemented }
