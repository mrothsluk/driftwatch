package drift

import (
	"fmt"
	"io"
	"strings"
)

// ReportFormat defines the output format for drift reports.
type ReportFormat string

const (
	FormatText ReportFormat = "text"
	FormatSummary ReportFormat = "summary"
)

// Report writes a drift report for the given diffs to the provided writer.
func Report(w io.Writer, diffs []ResourceDiff, format ReportFormat) error {
	switch format {
	case FormatText:
		return reportText(w, diffs)
	case FormatSummary:
		return reportSummary(w, diffs)
	default:
		return fmt.Errorf("unknown report format: %q", format)
	}
}

func reportText(w io.Writer, diffs []ResourceDiff) error {
	if len(diffs) == 0 {
		_, err := fmt.Fprintln(w, "No drift detected. Infrastructure matches Terraform state.")
		return err
	}
	_, err := fmt.Fprintf(w, "Detected %d drift(s):\n", len(diffs))
	if err != nil {
		return err
	}
	for _, d := range diffs {
		_, err = fmt.Fprintf(w, "  - %s\n", d.String())
		if err != nil {
			return err
		}
	}
	return nil
}

func reportSummary(w io.Writer, diffs []ResourceDiff) error {
	if len(diffs) == 0 {
		_, err := fmt.Fprintln(w, "Status: OK (0 drifts)")
		return err
	}
	resources := make(map[string]int)
	for _, d := range diffs {
		key := d.ResourceType + "." + d.ResourceName
		resources[key]++
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Status: DRIFT DETECTED (%d attribute(s) across %d resource(s))\n", len(diffs), len(resources)))
	for res, count := range resources {
		sb.WriteString(fmt.Sprintf("  %s: %d attribute(s) drifted\n", res, count))
	}
	_, err := fmt.Fprint(w, sb.String())
	return err
}
