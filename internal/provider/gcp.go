package provider

import (
	"context"
	"fmt"
	"strings"
)

// GCPProvider implements Provider for Google Cloud Platform resources.
type GCPProvider struct {
	projectID string
	// client would be a real GCP client in production
}

// NewGCPProvider creates a new GCPProvider for the given project.
func NewGCPProvider(projectID string) *GCPProvider {
	return &GCPProvider{projectID: projectID}
}

// GetAttributes fetches live attributes for a GCP resource by type and ID.
// Currently supports: google_compute_instance, google_storage_bucket.
func (p *GCPProvider) GetAttributes(ctx context.Context, resourceType, resourceID string) (map[string]string, error) {
	if resourceID == "" {
		return nil, fmt.Errorf("gcp: resource ID must not be empty")
	}

	switch resourceType {
	case "google_compute_instance":
		return p.getComputeInstanceAttributes(ctx, resourceID)
	case "google_storage_bucket":
		return p.getStorageBucketAttributes(ctx, resourceID)
	default:
		return nil, fmt.Errorf("gcp: unsupported resource type %q", resourceType)
	}
}

func (p *GCPProvider) getComputeInstanceAttributes(ctx context.Context, instanceID string) (map[string]string, error) {
	// In production this would call the Compute Engine API.
	// Stub implementation for demonstration.
	parts := strings.SplitN(instanceID, "/", 3)
	if len(parts) < 1 {
		return nil, fmt.Errorf("gcp: invalid compute instance ID %q", instanceID)
	}
	return map[string]string{
		"name":         instanceID,
		"project":      p.projectID,
		"machine_type": "n1-standard-1",
		"status":       "RUNNING",
	}, nil
}

func (p *GCPProvider) getStorageBucketAttributes(ctx context.Context, bucketName string) (map[string]string, error) {
	// In production this would call the Cloud Storage API.
	// Stub implementation for demonstration.
	if bucketName == "" {
		return nil, fmt.Errorf("gcp: bucket name must not be empty")
	}
	return map[string]string{
		"name":          bucketName,
		"project":       p.projectID,
		"location":      "US",
		"storage_class": "STANDARD",
	}, nil
}
