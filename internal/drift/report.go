package drift

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

// Report holds the results of a drift detection run.
type Report struct {
	Drifts []ResourceDrift
}

// ResourceDrift describes the drift found for a single resource.
type ResourceDrift struct {
	ResourceType string
	ResourceName string
	Diffs        []AttributeDiff
}

// AttributeDiff describes a single attribute mismatch.
type AttributeDiff struct {
	Attribute string
	Expected  string
	Actual    string
}

// HasDrift returns true if any drift was detected.
func (r *Report) HasDrift() bool {
	return len(r.Drifts) > 0
}

// reportText writes a human-readable drift report to w.
func reportText(r *Report, w io.Writer) {
	if !r.HasDrift() {
		fmt.Fprintln(w, "✓ No drift detected. Infrastructure matches Terraform state.")
		return
	}

	fmt.Fprintf(w, "⚠ Drift detected in %d resource(s):\n\n", len(r.Drifts))

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, rd := range r.Drifts {
		fmt.Fprintf(tw, "Resource:\t%s.%s\n", rd.ResourceType, rd.ResourceName)
		fmt.Fprintf(tw, "%-30s\t%-30s\t%-30s\n",
			"ATTRIBUTE", "EXPECTED (state)", "ACTUAL (live)")
		fmt.Fprintf(tw, "%s\t%s\t%s\n",
			strings.Repeat("-", 30),
			strings.Repeat("-", 30),
			strings.Repeat("-", 30))
		for _, diff := range rd.Diffs {
			fmt.Fprintf(tw, "%-30s\t%-30s\t%-30s\n",
				diff.Attribute, diff.Expected, diff.Actual)
		}
		_ = tw.Flush()
		fmt.Fprintln(w)
	}
}

// reportSummary writes a one-line summary to w.
func reportSummary(r *Report, w io.Writer) {
	if !r.HasDrift() {
		fmt.Fprintln(w, "No drift detected.")
		return
	}
	total := 0
	for _, rd := range r.Drifts {
		total += len(rd.Diffs)
	}
	fmt.Fprintf(w, "Drift detected: %d resource(s) with %d attribute diff(s).\n",
		len(r.Drifts), total)
}

// Print writes the full drift report to stdout.
func (r *Report) Print() {
	reportText(r, os.Stdout)
}

// Summary writes a one-line summary to stdout.
func (r *Report) Summary() {
	reportSummary(r, os.Stdout)
}
