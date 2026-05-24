package drift

import (
	"fmt"

	"github.com/driftwatch/internal/tfstate"
)

// ResourceDiff represents a detected difference between Terraform state and live cloud state.
type ResourceDiff struct {
	ResourceType string
	ResourceName string
	Attribute    string
	Expected     interface{}
	Actual       interface{}
}

// String returns a human-readable representation of the diff.
func (d ResourceDiff) String() string {
	return fmt.Sprintf(
		"[%s.%s] attribute %q: expected=%v, actual=%v",
		d.ResourceType, d.ResourceName, d.Attribute, d.Expected, d.Actual,
	)
}

// LiveStateProvider is an interface for fetching live resource attributes from a cloud provider.
type LiveStateProvider interface {
	GetAttributes(resourceType, resourceName string) (map[string]interface{}, error)
}

// Detector compares Terraform state against live cloud state.
type Detector struct {
	Provider LiveStateProvider
}

// NewDetector creates a new Detector with the given LiveStateProvider.
func NewDetector(provider LiveStateProvider) *Detector {
	return &Detector{Provider: provider}
}

// Detect compares all resources in the given TerraformState against live state
// and returns a list of detected diffs.
func (d *Detector) Detect(state *tfstate.TerraformState) ([]ResourceDiff, error) {
	var diffs []ResourceDiff

	for _, resource := range state.Resources {
		liveAttrs, err := d.Provider.GetAttributes(resource.Type, resource.Name)
		if err != nil {
			return nil, fmt.Errorf("fetching live state for %s.%s: %w", resource.Type, resource.Name, err)
		}

		for _, instance := range resource.Instances {
			for key, expectedVal := range instance.Attributes {
				actualVal, exists := liveAttrs[key]
				if !exists {
					diffs = append(diffs, ResourceDiff{
						ResourceType: resource.Type,
						ResourceName: resource.Name,
						Attribute:    key,
						Expected:     expectedVal,
						Actual:       nil,
					})
					continue
				}
				if fmt.Sprintf("%v", expectedVal) != fmt.Sprintf("%v", actualVal) {
					diffs = append(diffs, ResourceDiff{
						ResourceType: resource.Type,
						ResourceName: resource.Name,
						Attribute:    key,
						Expected:     expectedVal,
						Actual:       actualVal,
					})
				}
			}
		}
	}

	return diffs, nil
}
