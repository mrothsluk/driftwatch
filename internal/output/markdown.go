package output

import (
	"fmt"
	"io"

	"github.com/example/driftwatch/internal/drift"
)

// MarkdownFormatter formats a drift report as Markdown.
type MarkdownFormatter struct{}

// Format writes the report as a Markdown document to w.
func (f *MarkdownFormatter) Format(report *drift.Report, w io.Writer) error {
	fmt.Fprintln(w, "# Driftwatch Report")
	fmt.Fprintln(w)

	if !report.HasDrift() {
		fmt.Fprintln(w, "**Status:** ✅ No drift detected.")
		return nil
	}

	fmt.Fprintf(w, "**Status:** ⚠️ Drift detected in %d resource(s).\n\n", len(report.Drifts))

	for _, d := range report.Drifts {
		fmt.Fprintf(w, "## `%s` (%s)\n\n", d.ResourceID, d.ResourceType)
		fmt.Fprintln(w, "| Attribute | Expected | Actual |")
		fmt.Fprintln(w, "|-----------|----------|--------|")
		for _, attr := range d.Attributes {
			fmt.Fprintf(w, "| `%s` | `%s` | `%s` |\n", attr.Key, attr.Expected, attr.Actual)
		}
		fmt.Fprintln(w)
	}
	return nil
}
