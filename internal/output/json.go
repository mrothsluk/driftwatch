package output

import (
	"encoding/json"
	"io"

	"github.com/example/driftwatch/internal/drift"
)

// JSONFormatter formats a drift report as JSON.
type JSONFormatter struct{}

type jsonReport struct {
	HasDrift bool         `json:"has_drift"`
	Drifts   []jsonDrift  `json:"drifts"`
}

type jsonDrift struct {
	ResourceID   string          `json:"resource_id"`
	ResourceType string          `json:"resource_type"`
	Attributes   []jsonAttribute `json:"attributes"`
}

type jsonAttribute struct {
	Key      string `json:"key"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
}

// Format writes the report as a JSON object to w.
func (f *JSONFormatter) Format(report *drift.Report, w io.Writer) error {
	out := jsonReport{
		HasDrift: report.HasDrift(),
		Drifts:   make([]jsonDrift, 0, len(report.Drifts)),
	}
	for _, d := range report.Drifts {
		attrs := make([]jsonAttribute, 0, len(d.Attributes))
		for _, a := range d.Attributes {
			attrs = append(attrs, jsonAttribute{Key: a.Key, Expected: a.Expected, Actual: a.Actual})
		}
		out.Drifts = append(out.Drifts, jsonDrift{
			ResourceID:   d.ResourceID,
			ResourceType: d.ResourceType,
			Attributes:   attrs,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
