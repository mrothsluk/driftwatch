// Package output provides formatting utilities for drift reports.
package output

import (
	"io"

	"github.com/example/driftwatch/internal/drift"
)

// Format is the output format type.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
	FormatMarkdown Format = "markdown"
)

// Formatter writes a drift report in a specific format.
type Formatter interface {
	Format(report *drift.Report, w io.Writer) error
}

// New returns a Formatter for the given format string.
// It returns a TextFormatter for unknown formats.
func New(format Format) Formatter {
	switch format {
	case FormatJSON:
		return &JSONFormatter{}
	case FormatMarkdown:
		return &MarkdownFormatter{}
	default:
		return &TextFormatter{}
	}
}
