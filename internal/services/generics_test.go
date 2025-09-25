package services

import (
	"context"
	"testing"

	"flexpane/internal/models"
)

// MockTypedPane for testing
type MockTypedTodoPane struct {
	id       string
	title    string
	template string
	data     models.TodoPaneData
}

func (m *MockTypedTodoPane) ID() string {
	return m.id
}

func (m *MockTypedTodoPane) Title() string {
	return m.title
}

func (m *MockTypedTodoPane) Template() string {
	return m.template
}

func (m *MockTypedTodoPane) GetData(ctx context.Context) (interface{}, error) {
	return m.data, nil
}

func (m *MockTypedTodoPane) GetTypedData(ctx context.Context) (models.TodoPaneData, error) {
	return m.data, nil
}

// TestGetTypedPane validates the generic pane retrieval functionality
func TestGetTypedPane(t *testing.T) {
	registry := NewPaneRegistry()
	
	todoData := models.TodoPaneData{
		Todos: []models.Todo{{Done: false, Message: "Test"}},
		Count: 1,
	}

	mockPane := &MockTypedTodoPane{
		id:       "test-todos",
		title:    "Test Todos",
		template: "todos.html",
		data:     todoData,
	}

	// Register the pane
	registry.RegisterPane(mockPane)
	
	// Test getting typed pane
	typedPane, exists := GetTypedPane[models.TodoPaneData](registry, "test-todos")
	if !exists {
		t.Error("Expected to find typed pane, but it was not found")
	}
	
	if typedPane == nil {
		t.Fatal("Expected typed pane to be non-nil")
	}
	
	if typedPane.ID() != "test-todos" {
		t.Errorf("Expected pane ID 'test-todos', got %s", typedPane.ID())
	}

	// Test getting typed pane data
	ctx := context.Background()
	paneData, err := GetTypedPaneData[models.TodoPaneData](ctx, registry, "test-todos")
	if err != nil {
		t.Errorf("Unexpected error getting typed pane data: %v", err)
	}
	
	if paneData.ID != "test-todos" {
		t.Errorf("Expected pane data ID 'test-todos', got %s", paneData.ID)
	}
	
	if len(paneData.Data.Todos) != 1 {
		t.Errorf("Expected 1 todo, got %d", len(paneData.Data.Todos))
	}
}

// TestGetTypedPaneNotFound validates error handling for missing panes
func TestGetTypedPaneNotFound(t *testing.T) {
	registry := NewPaneRegistry()
	
	// Test getting non-existent pane
	_, exists := GetTypedPane[models.TodoPaneData](registry, "non-existent")
	if exists {
		t.Error("Expected pane not to exist, but it was found")
	}
	
	// Test getting typed pane data for non-existent pane
	ctx := context.Background()
	_, err := GetTypedPaneData[models.TodoPaneData](ctx, registry, "non-existent")
	if err == nil {
		t.Error("Expected error for non-existent pane, got nil")
	}
	
	expectedError := "pane not found: non-existent"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

// TestGetTypedPaneWrongType validates type safety
func TestGetTypedPaneWrongType(t *testing.T) {
	registry := NewPaneRegistry()
	
	todoData := models.TodoPaneData{
		Todos: []models.Todo{{Done: false, Message: "Test"}},
		Count: 1,
	}

	mockPane := &MockTypedTodoPane{
		id:       "test-todos",
		title:    "Test Todos",
		template: "todos.html",
		data:     todoData,
	}

	// Register the pane
	registry.RegisterPane(mockPane)
	
	// Try to get it as a different type (should fail)
	_, exists := GetTypedPane[models.CalendarPaneData](registry, "test-todos")
	if exists {
		t.Error("Expected type mismatch to result in pane not found, but it was found")
	}
}