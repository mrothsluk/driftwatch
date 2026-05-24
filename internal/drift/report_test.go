package drift

import (
	"bytes"
	"strings"
	"testing"
)

func sampleReport(drifts []ResourceDrift) *Report {
	return &Report{Drifts: drifts}
}

func TestReport_HasDrift_Empty(t *testing.T) {
	r := sampleReport(nil)
	if r.HasDrift() {
		t.Error("expected HasDrift() == false for empty report")
	}
}

func TestReport_HasDrift_WithDrift(t *testing.T) {
	r := sampleReport([]ResourceDrift{
		{ResourceType: "aws_instance", ResourceName: "web",
			Diffs: []AttributeDiff{{Attribute: "instance_type", Expected: "t2.micro", Actual: "t3.small"}}},
	})
	if !r.HasDrift() {
		t.Error("expected HasDrift() == true")
	}
}

func TestReportText_NoDrift(t *testing.T) {
	r := sampleReport(nil)
	var buf bytes.Buffer
	reportText(r, &buf)
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestReportText_WithDrift(t *testing.T) {
	r := sampleReport([]ResourceDrift{
		{
			ResourceType: "aws_instance",
			ResourceName: "web",
			Diffs: []AttributeDiff{
				{Attribute: "instance_type", Expected: "t2.micro", Actual: "t3.small"},
				{Attribute: "ami", Expected: "ami-abc123", Actual: "ami-xyz789"},
			},
		},
	})
	var buf bytes.Buffer
	reportText(r, &buf)
	out := buf.String()

	for _, want := range []string{"aws_instance.web", "instance_type", "t2.micro", "t3.small", "ami"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestReportSummary_NoDrift(t *testing.T) {
	r := sampleReport(nil)
	var buf bytes.Buffer
	reportSummary(r, &buf)
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Errorf("unexpected summary: %s", buf.String())
	}
}

func TestReportSummary_WithDrift(t *testing.T) {
	r := sampleReport([]ResourceDrift{
		{
			ResourceType: "aws_s3_bucket",
			ResourceName: "assets",
			Diffs: []AttributeDiff{
				{Attribute: "versioning", Expected: "true", Actual: "false"},
			},
		},
		{
			ResourceType: "aws_instance",
			ResourceName: "api",
			Diffs: []AttributeDiff{
				{Attribute: "instance_type", Expected: "t2.micro", Actual: "t3.medium"},
			},
		},
	})
	var buf bytes.Buffer
	reportSummary(r, &buf)
	out := buf.String()
	if !strings.Contains(out, "2 resource(s)") || !strings.Contains(out, "2 attribute diff(s)") {
		t.Errorf("unexpected summary output: %s", out)
	}
}
