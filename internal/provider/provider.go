package provider

import "context"

// Provider is the interface that cloud providers must implement
// to supply live resource attributes for drift detection.
type Provider interface {
	// GetAttributes returns the current live attributes for the resource
	// identified by resourceType (e.g. "aws_instance") and resourceID.
	GetAttributes(ctx context.Context, resourceType, resourceID string) (map[string]string, error)
}

// MockProvider is an in-memory Provider used in tests.
type MockProvider struct {
	// Resources maps "resourceType/resourceID" to attribute maps.
	Resources map[string]map[string]string
	// Err is returned for every call when non-nil.
	Err error
}

// GetAttributes implements Provider for MockProvider.
func (m *MockProvider) GetAttributes(_ context.Context, resourceType, resourceID string) (map[string]string, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	key := resourceType + "/" + resourceID
	attrs, ok := m.Resources[key]
	if !ok {
		return nil, nil
	}
	return attrs, nil
}
