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
	
	// Register built-in providers
	factory.RegisterProvider("mock", factory.createMockProvider)
	factory.RegisterProvider("file", factory.createFileProvider)
	
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

func (pf *ProviderFactory) createFileProvider(args map[string]interface{}) (Provider, error) {
	// For file provider, we create a composite provider that uses MockProvider for calendar/email
	// and TodoFileProvider for todos
	todoFilename, ok := args["todo_file"].(string)
	if !ok {
		todoFilename = "data/todos.json" // default
	}
	
	return NewCompositeProvider(
		NewMockProvider(),           // For calendar and email
		NewTodoFileProvider(todoFilename), // For todos
	), nil
}