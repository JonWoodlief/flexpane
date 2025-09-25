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
	dataProviderConstructors map[string]func(providers.DataProvider, map[string]interface{}) models.Pane
	serviceConstructors      map[string]func(map[string]interface{}) models.Pane
	dataProviderRegistry     map[string]providers.DataProvider
	todoService              models.TodoService
}

func NewPaneFactory(todoService models.TodoService) *PaneFactory {
	factory := &PaneFactory{
		dataProviderConstructors: make(map[string]func(providers.DataProvider, map[string]interface{}) models.Pane),
		serviceConstructors:      make(map[string]func(map[string]interface{}) models.Pane),
		dataProviderRegistry:     make(map[string]providers.DataProvider),
		todoService:              todoService,
	}

	// Register built-in data provider pane types (calendar, email)
	factory.RegisterDataProviderPaneType("calendar", factory.createCalendarPane)
	factory.RegisterDataProviderPaneType("email", factory.createEmailPane)
	
	// Register built-in service pane types (todos)
	factory.RegisterServicePaneType("todos", factory.createTodoPane)

	return factory
}

// RegisterDataProviderPaneType registers a pane constructor that uses DataProvider
func (pf *PaneFactory) RegisterDataProviderPaneType(paneType string, constructor func(providers.DataProvider, map[string]interface{}) models.Pane) {
	pf.dataProviderConstructors[paneType] = constructor
}

// RegisterServicePaneType registers a pane constructor that uses services
func (pf *PaneFactory) RegisterServicePaneType(paneType string, constructor func(map[string]interface{}) models.Pane) {
	pf.serviceConstructors[paneType] = constructor
}

// RegisterDataProvider registers a data provider instance
func (pf *PaneFactory) RegisterDataProvider(name string, provider providers.DataProvider) {
	pf.dataProviderRegistry[name] = provider
}

// CreatePane creates a pane based on configuration
func (pf *PaneFactory) CreatePane(config PaneConfig) (models.Pane, error) {
	// Check if it's a data provider pane type
	if constructor, exists := pf.dataProviderConstructors[config.Type]; exists {
		// Get data provider for this pane
		var provider providers.DataProvider
		if config.Provider != "" {
			var ok bool
			provider, ok = pf.dataProviderRegistry[config.Provider]
			if !ok {
				return nil, fmt.Errorf("unknown data provider: %s", config.Provider)
			}
		} else {
			// Use default provider (first available)
			for _, p := range pf.dataProviderRegistry {
				provider = p
				break
			}
			if provider == nil {
				return nil, fmt.Errorf("no data providers available")
			}
		}
		
		return constructor(provider, config.Args), nil
	}
	
	// Check if it's a service pane type
	if constructor, exists := pf.serviceConstructors[config.Type]; exists {
		return constructor(config.Args), nil
	}
	
	return nil, fmt.Errorf("unknown pane type: %s", config.Type)
}

// GetAvailablePaneTypes returns list of available pane types
func (pf *PaneFactory) GetAvailablePaneTypes() []string {
	var types []string
	for paneType := range pf.dataProviderConstructors {
		types = append(types, paneType)
	}
	for paneType := range pf.serviceConstructors {
		types = append(types, paneType)
	}
	return types
}

// Built-in data provider pane constructors
func (pf *PaneFactory) createCalendarPane(provider providers.DataProvider, args map[string]interface{}) models.Pane {
	return panes.NewCalendarPane(provider)
}

func (pf *PaneFactory) createEmailPane(provider providers.DataProvider, args map[string]interface{}) models.Pane {
	return panes.NewEmailPane(provider)
}

// Built-in service pane constructors
func (pf *PaneFactory) createTodoPane(args map[string]interface{}) models.Pane {
	return panes.NewTodoPane(pf.todoService)
}