package notify

import (
	"fmt"
	"log"

	"github.com/example/driftwatch/internal/drift"
)

// logNotifier writes drift summaries to the standard logger.
type logNotifier struct{}

// Notify logs a human-readable summary of the drift report.
func (l *logNotifier) Notify(report *drift.Report) error {
	if !report.HasDrift() {
		log.Println("[driftwatch] no drift detected")
		return nil
	}
	log.Printf("[driftwatch] drift detected in %d resource(s):", len(report.Drifts))
	for _, d := range report.Drifts {
		for _, diff := range d.Diffs {
			log.Println(fmt.Sprintf("  %s %s: want=%q got=%q",
				d.ResourceID, diff.Attribute, diff.Expected, diff.Actual))
		}
	}
	return nil
}
