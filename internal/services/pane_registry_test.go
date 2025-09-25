package services

import (
	"context"
	"testing"

	"flexplane/internal/models"
)

// MockPane for testing
type MockPane struct {
	id       string
	title    string
	template string
	data     interface{}
	err      error
}

func (m *MockPane) ID() string                                        { return m.id }
func (m *MockPane) Title() string                                     { return m.title }
func (m *MockPane) Template() string                                  { return m.template }
func (m *MockPane) GetData(ctx context.Context) (interface{}, error) { return m.data, m.err }

func TestPaneRegistry_RegisterPane(t *testing.T) {
	registry := NewPaneRegistry()

	pane := &MockPane{
		id:    "test",
		title: "Test Pane",
	}

	registry.RegisterPane(pane)

	retrieved, exists := registry.GetPane("test")
	if !exists {
		t.Fatal("Expected pane to exist after registration")
	}

	if retrieved.ID() != "test" {
		t.Errorf("Expected ID 'test', got '%s'", retrieved.ID())
	}
}

func TestPaneRegistry_GetEnabledPanes(t *testing.T) {
	registry := NewPaneRegistry()

	pane := &MockPane{id: "test", title: "Test", data: "test-data"}
	registry.RegisterPane(pane)
	registry.SetEnabledPanes([]string{"test"})

	// Set required layout config
	registry.SetLayoutConfig(map[string]PaneLayoutConfig{
		"test": {
			GridArea: models.PaneGridArea{Row: "1", Column: "span 2"},
		},
	})

	paneData, err := registry.GetEnabledPanes(context.Background())
	if err != nil {
		t.Fatalf("GetEnabledPanes failed: %v", err)
	}

	if len(paneData) != 1 {
		t.Errorf("Expected 1 pane, got %d", len(paneData))
	}

	if paneData[0].ID != "test" {
		t.Errorf("Expected ID 'test', got '%s'", paneData[0].ID)
	}
}

// Removed ordering and error tests - not needed with simplified design