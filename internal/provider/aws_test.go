package provider_test

import (
	"context"
	"errors"
	"testing"

	"github.com/example/driftwatch/internal/provider"
)

func TestMockProvider_GetAttributes_Found(t *testing.T) {
	mp := &provider.MockProvider{
		Resources: map[string]map[string]string{
			"aws_instance/i-abc123": {
				"instance_type":  "t3.micro",
				"instance_state": "running",
			},
		},
	}

	attrs, err := mp.GetAttributes(context.Background(), "aws_instance", "i-abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["instance_type"] != "t3.micro" {
		t.Errorf("expected t3.micro, got %s", attrs["instance_type"])
	}
	if attrs["instance_state"] != "running" {
		t.Errorf("expected running, got %s", attrs["instance_state"])
	}
}

func TestMockProvider_GetAttributes_NotFound(t *testing.T) {
	mp := &provider.MockProvider{
		Resources: map[string]map[string]string{},
	}

	attrs, err := mp.GetAttributes(context.Background(), "aws_instance", "i-missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs != nil {
		t.Errorf("expected nil attrs for missing resource, got %v", attrs)
	}
}

func TestMockProvider_GetAttributes_Error(t *testing.T) {
	sentinel := errors.New("api unavailable")
	mp := &provider.MockProvider{Err: sentinel}

	_, err := mp.GetAttributes(context.Background(), "aws_instance", "i-abc123")
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestMockProvider_GetAttributes_MultipleResources(t *testing.T) {
	mp := &provider.MockProvider{
		Resources: map[string]map[string]string{
			"aws_instance/i-1": {"instance_type": "t3.small"},
			"aws_s3_bucket/my-bucket": {"bucket": "my-bucket", "exists": "true"},
		},
	}

	a1, err := mp.GetAttributes(context.Background(), "aws_instance", "i-1")
	if err != nil || a1["instance_type"] != "t3.small" {
		t.Errorf("unexpected result for aws_instance: %v %v", a1, err)
	}

	a2, err := mp.GetAttributes(context.Background(), "aws_s3_bucket", "my-bucket")
	if err != nil || a2["exists"] != "true" {
		t.Errorf("unexpected result for aws_s3_bucket: %v %v", a2, err)
	}
}
