package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/output"
)

func sampleReport(hasDrift bool) *drift.Report {
	if !hasDrift {
		return &drift.Report{Drifts: []drift.ResourceDrift{}}
	}
	return &drift.Report{
		Drifts: []drift.ResourceDrift{
			{
				ResourceID:   "i-abc123",
				ResourceType: "aws_instance",
				Attributes: []drift.AttributeDiff{
					{Key: "instance_type", Expected: "t3.micro", Actual: "t3.small"},
				},
			},
		},
	}
}

func TestNew_ReturnsTextByDefault(t *testing.T) {
	f := output.New("unknown")
	if _, ok := f.(*output.TextFormatter); !ok {
		t.Errorf("expected TextFormatter, got %T", f)
	}
}

func TestTextFormatter_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := &output.TextFormatter{}
	if err := f.Format(sampleReport(false), &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No drift") {
		t.Errorf("expected 'No drift' in output, got: %s", buf.String())
	}
}

func TestTextFormatter_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	f := &output.TextFormatter{}
	if err := f.Format(sampleReport(true), &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "i-abc123") {
		t.Errorf("expected resource id in output")
	}
	if !strings.Contains(out, "instance_type") {
		t.Errorf("expected attribute key in output")
	}
}

func TestJSONFormatter_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	f := &output.JSONFormatter{}
	if err := f.Format(sampleReport(true), &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, `"has_drift": true`) {
		t.Errorf("expected has_drift true in JSON output, got: %s", out)
	}
}

func TestMarkdownFormatter_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := &output.MarkdownFormatter{}
	if err := f.Format(sampleReport(false), &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No drift") {
		t.Errorf("expected 'No drift' in markdown output")
	}
}

func TestMarkdownFormatter_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	f := &output.MarkdownFormatter{}
	if err := f.Format(sampleReport(true), &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "| Attribute |") {
		t.Errorf("expected markdown table in output, got: %s", out)
	}
}
