package providers

import (
	"fmt"
)

// ProviderType represents the type of data provider
type ProviderType string

const (
	ProviderTypeMock ProviderType = "mock"
)

// ProviderConfig holds configuration for data providers
type ProviderConfig struct {
	Type   ProviderType           `json:"type"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// GenericProviderFactory provides type-safe provider creation
// T represents the data type the provider will return
type GenericProviderFactory[T any] struct {
	createFunc func(string) (T, error)
}

// NewGenericProviderFactory creates a type-safe provider factory
// This eliminates the need for type assertions when creating providers
func NewGenericProviderFactory[T any](createFunc func(string) (T, error)) *GenericProviderFactory[T] {
	return &GenericProviderFactory[T]{
		createFunc: createFunc,
	}
}

// CreateTypedProvider creates a provider with compile-time type safety
func (gf *GenericProviderFactory[T]) CreateTypedProvider(providerType string) (T, error) {
	return gf.createFunc(providerType)
}

// CreateProvider creates a data provider based on the specified provider type
// This replaces the verbose factory pattern with a simple switch statement
func CreateProvider(providerType string) (DataProvider, error) {
	switch providerType {
	case "mock":
		return NewMockProvider(), nil
	case "": // Default to mock if empty
		return NewMockProvider(), nil
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
}
