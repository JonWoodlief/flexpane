package providers

import (
	"fmt"
)

// ProviderConfig holds configuration for provider initialization
type ProviderConfig struct {
	Type string                 `json:"type"`
	Args map[string]interface{} `json:"args"`
}

// ProviderFactory creates email/calendar providers based on configuration  
type ProviderFactory struct {
	constructors map[string]func(map[string]interface{}) (DataProvider, error)
}

func NewProviderFactory() *ProviderFactory {
	factory := &ProviderFactory{
		constructors: make(map[string]func(map[string]interface{}) (DataProvider, error)),
	}
	
	// Register all available provider types
	factory.RegisterProvider("null", factory.createNullProvider)
	factory.RegisterProvider("mock", factory.createMockProvider)
	
	return factory
}

// RegisterProvider registers a provider constructor
func (pf *ProviderFactory) RegisterProvider(providerType string, constructor func(map[string]interface{}) (DataProvider, error)) {
	pf.constructors[providerType] = constructor
}

// CreateProvider creates a provider based on configuration
func (pf *ProviderFactory) CreateProvider(config ProviderConfig) (DataProvider, error) {
	constructor, exists := pf.constructors[config.Type]
	if !exists {
		return nil, fmt.Errorf("unknown provider type: %s", config.Type)
	}
	
	return constructor(config.Args)
}

// GetAvailableProviders returns list of available provider types
func (pf *ProviderFactory) GetAvailableProviders() []string {
	var types []string
	for providerType := range pf.constructors {
		types = append(types, providerType)
	}
	return types
}

// Built-in provider constructors
func (pf *ProviderFactory) createMockProvider(args map[string]interface{}) (DataProvider, error) {
	return NewMockProvider(), nil
}

func (pf *ProviderFactory) createNullProvider(args map[string]interface{}) (DataProvider, error) {
	return NewNullProvider(), nil
}