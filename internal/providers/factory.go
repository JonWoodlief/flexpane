package providers

import (
	"fmt"
)

// ProviderConfig holds configuration for provider initialization
type ProviderConfig struct {
	Type string                 `json:"type"`
	Args map[string]interface{} `json:"args"`
}

// ProviderFactory creates providers based on configuration
type ProviderFactory struct {
	constructors map[string]func(map[string]interface{}) (Provider, error)
}

func NewProviderFactory() *ProviderFactory {
	factory := &ProviderFactory{
		constructors: make(map[string]func(map[string]interface{}) (Provider, error)),
	}
	
	// Register production providers only
	factory.RegisterProvider("file", factory.createFileProvider)
	factory.RegisterProvider("null", factory.createNullProvider)
	
	return factory
}

// NewProviderFactoryWithMocks creates a factory with mock providers for development/testing
func NewProviderFactoryWithMocks() *ProviderFactory {
	factory := NewProviderFactory()
	
	// Add mock providers for development/testing
	factory.RegisterProvider("mock", factory.createMockProvider)
	
	return factory
}

// RegisterProvider registers a provider constructor
func (pf *ProviderFactory) RegisterProvider(providerType string, constructor func(map[string]interface{}) (Provider, error)) {
	pf.constructors[providerType] = constructor
}

// CreateProvider creates a provider based on configuration
func (pf *ProviderFactory) CreateProvider(config ProviderConfig) (Provider, error) {
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
func (pf *ProviderFactory) createMockProvider(args map[string]interface{}) (Provider, error) {
	return NewMockProvider(), nil
}

func (pf *ProviderFactory) createNullProvider(args map[string]interface{}) (Provider, error) {
	return NewNullProvider(), nil
}

func (pf *ProviderFactory) createFileProvider(args map[string]interface{}) (Provider, error) {
	// For file provider, we create a composite provider that uses NullProvider for calendar/email
	// and TodoFileProvider for todos. This is appropriate for production when real calendar/email
	// integrations aren't configured yet.
	todoFilename, ok := args["todo_file"].(string)
	if !ok {
		todoFilename = "data/todos.json" // default
	}
	
	return NewCompositeProvider(
		NewNullProvider(),              // For calendar and email - no fake data in production
		NewTodoFileProvider(todoFilename), // For todos
	), nil
}