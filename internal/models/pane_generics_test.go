package models

import (
	"context"
	"testing"
	"time"
)

// TestTypedPaneData validates the generic TypedPaneData functionality
func TestTypedPaneData(t *testing.T) {
	// Test TodoPaneData
	todoData := TodoPaneData{
		Todos: []Todo{
			{Done: false, Message: "Test task"},
			{Done: true, Message: "Completed task"},
		},
		Count: 2,
	}

	typedTodoPaneData := TypedPaneData[TodoPaneData]{
		ID:       "todos",
		Title:    "Test Todos",
		GridArea: PaneGridArea{Row: "1", Column: "span 2"},
		Data:     todoData,
		Template: "todos.html",
	}

	// Convert to generic PaneData for backward compatibility
	paneData := typedTodoPaneData.ToPaneData()
	
	if paneData.ID != "todos" {
		t.Errorf("Expected ID 'todos', got %s", paneData.ID)
	}
	
	if paneData.Title != "Test Todos" {
		t.Errorf("Expected title 'Test Todos', got %s", paneData.Title)
	}

	// Test CalendarPaneData
	now := time.Now()
	calendarData := CalendarPaneData{
		Events: []Event{
			{ID: "1", Title: "Meeting", Start: now, End: now.Add(time.Hour)},
		},
		Count: 1,
	}

	typedCalendarPaneData := TypedPaneData[CalendarPaneData]{
		ID:       "calendar",
		Title:    "Test Calendar",
		GridArea: PaneGridArea{Row: "2", Column: "span 3"},
		Data:     calendarData,
		Template: "calendar.html",
	}

	calendarPaneData := typedCalendarPaneData.ToPaneData()
	
	if calendarPaneData.ID != "calendar" {
		t.Errorf("Expected ID 'calendar', got %s", calendarPaneData.ID)
	}
}

// MockTypedPane for testing
type MockTypedPane struct {
	id       string
	title    string
	template string
	data     TodoPaneData
}

func (m *MockTypedPane) ID() string {
	return m.id
}

func (m *MockTypedPane) Title() string {
	return m.title
}

func (m *MockTypedPane) Template() string {
	return m.template
}

func (m *MockTypedPane) GetData(ctx context.Context) (interface{}, error) {
	return m.data, nil
}

func (m *MockTypedPane) GetTypedData(ctx context.Context) (TodoPaneData, error) {
	return m.data, nil
}

// TestTypedPaneInterface validates the generic TypedPane interface
func TestTypedPaneInterface(t *testing.T) {
	todoData := TodoPaneData{
		Todos: []Todo{{Done: false, Message: "Test"}},
		Count: 1,
	}

	mockPane := &MockTypedPane{
		id:       "test-todos",
		title:    "Test Todos Pane",
		template: "test-todos.html",
		data:     todoData,
	}

	// Test that our mock implements both interfaces
	var pane Pane = mockPane
	var typedPane TypedPane[TodoPaneData] = mockPane

	// Test basic Pane interface methods
	if pane.ID() != "test-todos" {
		t.Errorf("Expected ID 'test-todos', got %s", pane.ID())
	}

	if pane.Title() != "Test Todos Pane" {
		t.Errorf("Expected title 'Test Todos Pane', got %s", pane.Title())
	}

	// Test typed data access
	typedData, err := typedPane.GetTypedData(context.Background())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(typedData.Todos) != 1 {
		t.Errorf("Expected 1 todo, got %d", len(typedData.Todos))
	}

	if typedData.Todos[0].Message != "Test" {
		t.Errorf("Expected message 'Test', got %s", typedData.Todos[0].Message)
	}
}