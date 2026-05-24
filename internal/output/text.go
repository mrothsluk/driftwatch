package output

import (
	"fmt"
	"io"

	"github.com/example/driftwatch/internal/drift"
)

// TextFormatter formats a drift report as plain text.
type TextFormatter struct{}

// Format writes the report as human-readable text to w.
func (f *TextFormatter) Format(report *drift.Report, w io.Writer) error {
	if !report.HasDrift() {
		_, err := fmt.Fprintln(w, "✓ No drift detected.")
		return err
	}

	fmt.Fprintf(w, "⚠ Drift detected in %d resource(s):\n\n", len(report.Drifts))
	for _, d := range report.Drifts {
		fmt.Fprintf(w, "  Resource: %s (%s)\n", d.ResourceID, d.ResourceType)
		for _, attr := range d.Attributes {
			fmt.Fprintf(w, "    - %s: want=%q, got=%q\n", attr.Key, attr.Expected, attr.Actual)
		}
		fmt.Fprintln(w)
	}
	return nil
}
