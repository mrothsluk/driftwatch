package drift_test

import (
	"errors"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/tfstate"
)

// mockProvider is a test implementation of LiveStateProvider.
type mockProvider struct {
	attrs map[string]map[string]interface{}
	err   error
}

func (m *mockProvider) GetAttributes(resourceType, resourceName string) (map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	key := resourceType + "." + resourceName
	return m.attrs[key], nil
}

func sampleState() *tfstate.TerraformState {
	return &tfstate.TerraformState{
		Version: 4,
		Resources: []tfstate.Resource{
			{
				Type: "aws_instance",
				Name: "web",
				Instances: []tfstate.Instance{
					{Attributes: map[string]interface{}{"instance_type": "t2.micro", "ami": "ami-123"}},
				},
			},
		},
	}
}

func TestDetect_NoDrift(t *testing.T) {
	provider := &mockProvider{
		attrs: map[string]map[string]interface{}{
			"aws_instance.web": {"instance_type": "t2.micro", "ami": "ami-123"},
		},
	}
	detector := drift.NewDetector(provider)
	diffs, err := detector.Detect(sampleState())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diffs) != 0 {
		t.Errorf("expected no diffs, got %d: %v", len(diffs), diffs)
	}
}

func TestDetect_AttributeMismatch(t *testing.T) {
	provider := &mockProvider{
		attrs: map[string]map[string]interface{}{
			"aws_instance.web": {"instance_type": "t3.medium", "ami": "ami-123"},
		},
	}
	detector := drift.NewDetector(provider)
	diffs, err := detector.Detect(sampleState())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Attribute != "instance_type" {
		t.Errorf("expected diff on instance_type, got %q", diffs[0].Attribute)
	}
}

func TestDetect_ProviderError(t *testing.T) {
	provider := &mockProvider{err: errors.New("api error")}
	detector := drift.NewDetector(provider)
	_, err := detector.Detect(sampleState())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
