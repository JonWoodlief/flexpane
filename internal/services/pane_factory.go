package services

import (
	"fmt"

	"flexplane/internal/models"
	"flexplane/internal/panes"
	"flexplane/internal/providers"
)

// PaneConfig holds configuration for pane initialization
type PaneConfig struct {
	Type     string                 `json:"type"`
	Enabled  bool                   `json:"enabled"`
	Layout   PaneLayoutConfig       `json:"layout"`
	Provider string                 `json:"provider,omitempty"` // Provider name to use
	Args     map[string]interface{} `json:"args,omitempty"`     // Additional pane-specific arguments
}

// PaneFactory creates panes based on configuration
type PaneFactory struct {
	constructors     map[string]func(providers.Provider, map[string]interface{}) models.Pane
	providerRegistry map[string]providers.Provider
}

func NewPaneFactory() *PaneFactory {
	factory := &PaneFactory{
		constructors:     make(map[string]func(providers.Provider, map[string]interface{}) models.Pane),
		providerRegistry: make(map[string]providers.Provider),
	}

	// Register built-in pane types
	factory.RegisterPaneType("calendar", factory.createCalendarPane)
	factory.RegisterPaneType("todos", factory.createTodoPane)
	factory.RegisterPaneType("email", factory.createEmailPane)

	return factory
}

// RegisterPaneType registers a pane constructor
func (pf *PaneFactory) RegisterPaneType(paneType string, constructor func(providers.Provider, map[string]interface{}) models.Pane) {
	pf.constructors[paneType] = constructor
}

// RegisterProvider registers a provider instance
func (pf *PaneFactory) RegisterProvider(name string, provider providers.Provider) {
	pf.providerRegistry[name] = provider
}

// CreatePane creates a pane based on configuration
func (pf *PaneFactory) CreatePane(config PaneConfig) (models.Pane, error) {
	constructor, exists := pf.constructors[config.Type]
	if !exists {
		return nil, fmt.Errorf("unknown pane type: %s", config.Type)
	}

	// Get provider for this pane
	var provider providers.Provider
	if config.Provider != "" {
		var ok bool
		provider, ok = pf.providerRegistry[config.Provider]
		if !ok {
			return nil, fmt.Errorf("unknown provider: %s", config.Provider)
		}
	} else {
		// Use default provider (first available)
		for _, p := range pf.providerRegistry {
			provider = p
			break
		}
		if provider == nil {
			return nil, fmt.Errorf("no providers available")
		}
	}

	return constructor(provider, config.Args), nil
}

// GetAvailablePaneTypes returns list of available pane types
func (pf *PaneFactory) GetAvailablePaneTypes() []string {
	var types []string
	for paneType := range pf.constructors {
		types = append(types, paneType)
	}
	return types
}

// Built-in pane constructors
func (pf *PaneFactory) createCalendarPane(provider providers.Provider, args map[string]interface{}) models.Pane {
	return panes.NewCalendarPane(provider)
}

func (pf *PaneFactory) createTodoPane(provider providers.Provider, args map[string]interface{}) models.Pane {
	return panes.NewTodoPane(provider)
}

func (pf *PaneFactory) createEmailPane(provider providers.Provider, args map[string]interface{}) models.Pane {
	return panes.NewEmailPane(provider)
}