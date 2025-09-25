package services

import (
	"context"
	"testing"

	"flexplane/internal/models"
	"flexplane/internal/providers"
)

func TestPaneFactory_CreatePane(t *testing.T) {
	factory := NewPaneFactory()
	mockProvider := providers.NewMockProvider()
	factory.RegisterProvider("test", mockProvider)

	// Test creating calendar pane
	calendarPane, err := factory.CreatePane(PaneConfig{
		Type:     "calendar",
		Provider: "test",
	})
	if err != nil {
		t.Errorf("Failed to create calendar pane: %v", err)
	}
	if calendarPane == nil {
		t.Error("Expected non-nil calendar pane")
	}
	if calendarPane.ID() != "calendar" {
		t.Errorf("Expected calendar ID, got %s", calendarPane.ID())
	}

	// Test creating todos pane
	todosPane, err := factory.CreatePane(PaneConfig{
		Type:     "todos",
		Provider: "test",
	})
	if err != nil {
		t.Errorf("Failed to create todos pane: %v", err)
	}
	if todosPane == nil {
		t.Error("Expected non-nil todos pane")
	}
	if todosPane.ID() != "todos" {
		t.Errorf("Expected todos ID, got %s", todosPane.ID())
	}

	// Test creating email pane
	emailPane, err := factory.CreatePane(PaneConfig{
		Type:     "email",
		Provider: "test",
	})
	if err != nil {
		t.Errorf("Failed to create email pane: %v", err)
	}
	if emailPane == nil {
		t.Error("Expected non-nil email pane")
	}
	if emailPane.ID() != "email" {
		t.Errorf("Expected email ID, got %s", emailPane.ID())
	}
}

func TestPaneFactory_UnknownPaneType(t *testing.T) {
	factory := NewPaneFactory()
	mockProvider := providers.NewMockProvider()
	factory.RegisterProvider("test", mockProvider)

	_, err := factory.CreatePane(PaneConfig{
		Type:     "unknown",
		Provider: "test",
	})
	if err == nil {
		t.Error("Expected error when creating unknown pane type")
	}
}

func TestPaneFactory_UnknownProvider(t *testing.T) {
	factory := NewPaneFactory()

	_, err := factory.CreatePane(PaneConfig{
		Type:     "calendar",
		Provider: "unknown",
	})
	if err == nil {
		t.Error("Expected error when using unknown provider")
	}
}

func TestPaneFactory_DefaultProvider(t *testing.T) {
	factory := NewPaneFactory()
	mockProvider := providers.NewMockProvider()
	factory.RegisterProvider("default", mockProvider)

	// Create pane without specifying provider - should use default
	pane, err := factory.CreatePane(PaneConfig{
		Type: "calendar",
		// No provider specified
	})
	if err != nil {
		t.Errorf("Failed to create pane with default provider: %v", err)
	}
	if pane == nil {
		t.Error("Expected non-nil pane with default provider")
	}
}

func TestPaneFactory_NoProviders(t *testing.T) {
	factory := NewPaneFactory()

	_, err := factory.CreatePane(PaneConfig{
		Type: "calendar",
	})
	if err == nil {
		t.Error("Expected error when no providers are available")
	}
}

func TestPaneFactory_GetAvailablePaneTypes(t *testing.T) {
	factory := NewPaneFactory()

	availableTypes := factory.GetAvailablePaneTypes()
	expectedTypes := []string{"calendar", "todos", "email"}

	if len(availableTypes) != len(expectedTypes) {
		t.Errorf("Expected %d pane types, got %d", len(expectedTypes), len(availableTypes))
	}

	for _, expected := range expectedTypes {
		found := false
		for _, available := range availableTypes {
			if available == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected pane type '%s' to be available", expected)
		}
	}
}

func TestPaneFactory_CustomPaneType(t *testing.T) {
	factory := NewPaneFactory()
	mockProvider := providers.NewMockProvider()
	factory.RegisterProvider("test", mockProvider)

	// Register a custom pane type
	factory.RegisterPaneType("custom", func(provider providers.Provider, args map[string]interface{}) models.Pane {
		return &mockPane{id: "custom", title: "Custom Pane"}
	})

	// Create custom pane
	pane, err := factory.CreatePane(PaneConfig{
		Type:     "custom",
		Provider: "test",
	})
	if err != nil {
		t.Errorf("Failed to create custom pane: %v", err)
	}
	if pane == nil {
		t.Error("Expected non-nil custom pane")
	}
	if pane.ID() != "custom" {
		t.Errorf("Expected custom ID, got %s", pane.ID())
	}
	if pane.Title() != "Custom Pane" {
		t.Errorf("Expected 'Custom Pane' title, got %s", pane.Title())
	}
}

// Mock pane for testing
type mockPane struct {
	id       string
	title    string
	template string
	data     interface{}
	err      error
}

func (m *mockPane) ID() string                                         { return m.id }
func (m *mockPane) Title() string                                      { return m.title }
func (m *mockPane) Template() string                                   { return m.template }
func (m *mockPane) GetData(ctx context.Context) (interface{}, error) { return m.data, m.err }