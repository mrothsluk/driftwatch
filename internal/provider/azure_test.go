package provider

import (
	"context"
	"errors"
	"testing"
)

// newTestAzureProvider returns an AzureProvider whose fetchFn is controlled
// by the caller, allowing deterministic unit tests without real Azure creds.
func newTestAzureProvider(fn func(ctx context.Context, resourceType, resourceID string) (map[string]string, error)) *AzureProvider {
	return &AzureProvider{fetchFn: fn}
}

func TestAzureProvider_GetAttributes_Found(t *testing.T) {
	want := map[string]string{"location": "eastus", "size": "Standard_D2s_v3"}

	p := newTestAzureProvider(func(_ context.Context, _, _ string) (map[string]string, error) {
		return want, nil
	})

	got, err := p.GetAttributes(context.Background(), "azurerm_virtual_machine", "vm-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("attr %q: got %q, want %q", k, got[k], v)
		}
	}
}

func TestAzureProvider_GetAttributes_EmptyID(t *testing.T) {
	p := newTestAzureProvider(func(_ context.Context, _, _ string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	_, err := p.GetAttributes(context.Background(), "azurerm_storage_account", "")
	if err == nil {
		t.Fatal("expected error for empty resource ID, got nil")
	}
}

func TestAzureProvider_GetAttributes_FetchError(t *testing.T) {
	sentinel := errors.New("azure api unavailable")

	p := newTestAzureProvider(func(_ context.Context, _, _ string) (map[string]string, error) {
		return nil, sentinel
	})

	_, err := p.GetAttributes(context.Background(), "azurerm_resource_group", "rg-prod")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, sentinel) {
		t.Errorf("error chain should contain sentinel; got: %v", err)
	}
}

func TestAzureProvider_TypePrefixStripped(t *testing.T) {
	var capturedType string

	p := newTestAzureProvider(func(_ context.Context, resourceType, _ string) (map[string]string, error) {
		capturedType = resourceType
		return map[string]string{}, nil
	})

	_, _ = p.GetAttributes(context.Background(), "azurerm_virtual_network", "vnet-1")

	if capturedType != "virtual_network" {
		t.Errorf("expected stripped type %q, got %q", "virtual_network", capturedType)
	}
}

func TestAzureProvider_ImplementsProvider(t *testing.T) {
	var _ Provider = NewAzureProvider()
}
