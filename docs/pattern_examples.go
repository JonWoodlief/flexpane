package main

import (
	"fmt"
	"flexplane/internal/models"
	"flexplane/internal/providers"
)

// EXAMPLE 1: Direct Case Statement Approach (NOT RECOMMENDED)
// This is what a simple case statement implementation would look like

func createProviderDirectCase(providerType string) (providers.DataProvider, error) {
	switch providerType {
	case "mock":
		return providers.NewMockProvider(), nil
	case "gmail":
		// return providers.NewGmailProvider(), nil  // Would need implementation
		return nil, fmt.Errorf("gmail provider not implemented")
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
}

// PROBLEMS with direct case approach:
// 1. Adding new provider requires code change here
// 2. No external configuration support
// 3. No fallback logic
// 4. No validation of provider settings
// 5. Violates Open/Closed Principle

// EXAMPLE 2: Enhanced Factory Pattern with Plugin Registry
// This shows how the factory pattern could be extended for even more flexibility

type ProviderCreator func(config map[string]interface{}) (providers.DataProvider, error)

type EnhancedProviderFactory struct {
	creators map[providers.ProviderType]ProviderCreator
	config   providers.DataProviderConfig
}

func NewEnhancedProviderFactory(configPath string) (*EnhancedProviderFactory, error) {
	// ... config loading logic similar to current factory ...
	
	factory := &EnhancedProviderFactory{
		creators: make(map[providers.ProviderType]ProviderCreator),
	}
	
	// Register built-in providers
	factory.RegisterProviderType(providers.ProviderTypeMock, func(config map[string]interface{}) (providers.DataProvider, error) {
		return providers.NewMockProvider(), nil
	})
	
	// Could register more providers dynamically
	// factory.RegisterProviderType("gmail", func(config map[string]interface{}) (providers.DataProvider, error) {
	//     return providers.NewGmailProvider(config), nil
	// })
	
	return factory, nil
}

func (f *EnhancedProviderFactory) RegisterProviderType(providerType providers.ProviderType, creator ProviderCreator) {
	f.creators[providerType] = creator
}

func (f *EnhancedProviderFactory) CreateProvider(name string) (providers.DataProvider, error) {
	providerConfig, exists := f.config.Providers[name]
	if !exists {
		// Fallback logic...
		return nil, fmt.Errorf("provider not found")
	}
	
	creator, exists := f.creators[providerConfig.Type]
	if !exists {
		return nil, fmt.Errorf("unsupported provider type: %s", providerConfig.Type)
	}
	
	return creator(providerConfig.Config)
}

// ADVANTAGES of enhanced factory:
// 1. Plugin-style registration
// 2. No switch/case statement needed
// 3. Completely extensible without code changes
// 4. Still maintains configuration benefits

// EXAMPLE 3: Registry Pattern for Panes (Current Approach)
// This shows why the registry pattern works well for panes

type PaneCreator func(dependencies interface{}) models.Pane

type EnhancedPaneRegistry struct {
	creators map[string]PaneCreator
	panes    map[string]models.Pane
}

func NewEnhancedPaneRegistry() *EnhancedPaneRegistry {
	registry := &EnhancedPaneRegistry{
		creators: make(map[string]PaneCreator),
		panes:    make(map[string]models.Pane),
	}
	
	// Register pane creators (still compile-time registration)
	// registry.RegisterPaneType("calendar", func(deps interface{}) models.Pane {
	//     return panes.NewCalendarPane(deps.(providers.DataProvider))
	// })
	
	return registry
}

// This is more complex than needed for panes because:
// 1. Panes are known at compile time
// 2. Dependencies are simple and available at startup
// 3. No runtime configuration needed
// 4. Performance matters more than flexibility

// EXAMPLE 4: Configuration Comparison

// Provider configuration (complex, runtime-flexible):
const exampleProviderConfig = `{
  "providers": {
    "gmail": {
      "type": "gmail",
      "config": {
        "api_key": "secret",
        "client_id": "client123",
        "scopes": ["readonly"]
      }
    },
    "mock": {
      "type": "mock",
      "config": {}
    }
  },
  "default": "mock"
}`

// Pane configuration (simple, layout-focused):
const examplePaneConfig = `{
  "enabled": ["calendar", "todos", "email"],
  "layout": {
    "calendar": {"grid_area": {"row": "1", "column": "span 2"}},
    "todos": {"grid_area": {"row": "2", "column": "1"}},
    "email": {"grid_area": {"row": "2", "column": "2"}}
  }
}`

// The configuration differences show why different patterns are appropriate:
// - Providers need complex runtime configuration
// - Panes need simple layout configuration