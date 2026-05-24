package provider_test

import (
	"context"
	"testing"

	"github.com/your-org/driftwatch/internal/provider"
)

func TestGCPProvider_GetAttributes_ComputeInstance(t *testing.T) {
	p := provider.NewGCPProvider("my-project")
	attrs, err := p.GetAttributes(context.Background(), "google_compute_instance", "my-instance")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["name"] != "my-instance" {
		t.Errorf("expected name=my-instance, got %q", attrs["name"])
	}
	if attrs["project"] != "my-project" {
		t.Errorf("expected project=my-project, got %q", attrs["project"])
	}
	if attrs["status"] != "RUNNING" {
		t.Errorf("expected status=RUNNING, got %q", attrs["status"])
	}
}

func TestGCPProvider_GetAttributes_StorageBucket(t *testing.T) {
	p := provider.NewGCPProvider("my-project")
	attrs, err := p.GetAttributes(context.Background(), "google_storage_bucket", "my-bucket")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["name"] != "my-bucket" {
		t.Errorf("expected name=my-bucket, got %q", attrs["name"])
	}
	if attrs["storage_class"] != "STANDARD" {
		t.Errorf("expected storage_class=STANDARD, got %q", attrs["storage_class"])
	}
}

func TestGCPProvider_GetAttributes_UnsupportedType(t *testing.T) {
	p := provider.NewGCPProvider("my-project")
	_, err := p.GetAttributes(context.Background(), "google_unknown_resource", "some-id")
	if err == nil {
		t.Fatal("expected error for unsupported resource type, got nil")
	}
}

func TestGCPProvider_GetAttributes_EmptyID(t *testing.T) {
	p := provider.NewGCPProvider("my-project")
	_, err := p.GetAttributes(context.Background(), "google_compute_instance", "")
	if err == nil {
		t.Fatal("expected error for empty resource ID, got nil")
	}
}

func TestGCPProvider_ImplementsProvider(t *testing.T) {
	// Ensure GCPProvider satisfies the Provider interface at compile time.
	var _ provider.Provider = provider.NewGCPProvider("proj")
}
