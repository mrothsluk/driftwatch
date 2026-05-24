package provider

import (
	"context"
	"fmt"
	"strings"
)

// AzureProvider fetches live resource attributes from Azure.
// In production this would use the Azure SDK; here we use a pluggable
// fetch function so tests can inject a mock.
type AzureProvider struct {
	// fetchFn is called with (resourceType, resourceID) and returns
	// a map of attribute name → value or an error.
	fetchFn func(ctx context.Context, resourceType, resourceID string) (map[string]string, error)
}

// NewAzureProvider returns an AzureProvider wired to real Azure SDK calls.
// Currently it returns a provider whose fetchFn delegates to the stub below;
// swap in real ARM client calls here when the SDK dependency is added.
func NewAzureProvider() *AzureProvider {
	return &AzureProvider{
		fetchFn: realAzureFetch,
	}
}

// realAzureFetch is the production implementation stub.
// Replace the body with actual Azure ARM / SDK calls.
func realAzureFetch(_ context.Context, resourceType, resourceID string) (map[string]string, error) {
	if resourceID == "" {
		return nil, fmt.Errorf("azure: resource ID must not be empty")
	}
	// TODO: integrate github.com/Azure/azure-sdk-for-go
	return nil, fmt.Errorf("azure: real fetch not yet implemented for type %s id %s", resourceType, resourceID)
}

// GetAttributes satisfies the Provider interface.
// It maps Terraform resource types (e.g. "azurerm_virtual_machine") to the
// Azure resource kind understood by fetchFn.
func (p *AzureProvider) GetAttributes(ctx context.Context, resourceType, resourceID string) (map[string]string, error) {
	if resourceID == "" {
		return nil, fmt.Errorf("azure: resource ID must not be empty")
	}

	// Normalise the Terraform type prefix so callers don't have to.
	azureType := strings.TrimPrefix(resourceType, "azurerm_")

	attrs, err := p.fetchFn(ctx, azureType, resourceID)
	if err != nil {
		return nil, fmt.Errorf("azure: GetAttributes(%s, %s): %w", resourceType, resourceID, err)
	}
	return attrs, nil
}
