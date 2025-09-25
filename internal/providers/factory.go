package providers

import (
	"encoding/json"
	"fmt"
	"os"
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

// DataProviderConfig represents the configuration for all data providers
type DataProviderConfig struct {
	Providers map[string]ProviderConfig `json:"providers"`
	Default   string                    `json:"default"`
}

// ProviderFactory creates data providers based on configuration
type ProviderFactory struct {
	config DataProviderConfig
}

// GenericProviderFactory provides type-safe provider creation
// T represents the data type the provider will return
type GenericProviderFactory[T any] struct {
	*ProviderFactory
	createFunc func(ProviderConfig) (T, error)
}

// NewGenericProviderFactory creates a type-safe provider factory
// This eliminates the need for type assertions when creating providers
func NewGenericProviderFactory[T any](configPath string, createFunc func(ProviderConfig) (T, error)) (*GenericProviderFactory[T], error) {
	factory, err := NewProviderFactory(configPath)
	if err != nil {
		return nil, err
	}
	
	return &GenericProviderFactory[T]{
		ProviderFactory: factory,
		createFunc:      createFunc,
	}, nil
}

// CreateTypedProvider creates a provider with compile-time type safety
func (gf *GenericProviderFactory[T]) CreateTypedProvider(name string) (T, error) {
	var zero T
	providerConfig, exists := gf.config.Providers[name]
	if !exists {
		if gf.config.Default == "" {
			return zero, fmt.Errorf("provider '%s' not found and no default configured", name)
		}
		providerConfig, exists = gf.config.Providers[gf.config.Default]
		if !exists {
			return zero, fmt.Errorf("default provider '%s' not found", gf.config.Default)
		}
	}
	
	return gf.createFunc(providerConfig)
}

// NewProviderFactory creates a new provider factory from configuration file
func NewProviderFactory(configPath string) (*ProviderFactory, error) {
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read provider config: %w", err)
	}

	var config DataProviderConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("failed to parse provider config: %w", err)
	}

	return &ProviderFactory{config: config}, nil
}

// CreateProvider creates a data provider based on the specified name
// Falls back to default provider if name is not found
func (f *ProviderFactory) CreateProvider(name string) (DataProvider, error) {
	providerConfig, exists := f.config.Providers[name]
	if !exists {
		if f.config.Default == "" {
			return nil, fmt.Errorf("provider '%s' not found and no default configured", name)
		}
		providerConfig, exists = f.config.Providers[f.config.Default]
		if !exists {
			return nil, fmt.Errorf("default provider '%s' not found", f.config.Default)
		}
	}

	switch providerConfig.Type {
	case ProviderTypeMock:
		return f.createMockProvider(providerConfig.Config)
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerConfig.Type)
	}
}

// GetDefaultProvider returns the default data provider
func (f *ProviderFactory) GetDefaultProvider() (DataProvider, error) {
	if f.config.Default == "" {
		return nil, fmt.Errorf("no default provider configured")
	}
	return f.CreateProvider(f.config.Default)
}

// createMockProvider creates a mock provider (no configuration needed)
func (f *ProviderFactory) createMockProvider(_ map[string]interface{}) (DataProvider, error) {
	return NewMockProvider(), nil
}
