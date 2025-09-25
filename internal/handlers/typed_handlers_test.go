package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"flexpane/internal/models"
	"flexpane/internal/panes"
	"flexpane/internal/providers"
	"flexpane/internal/services"
)

// TestTypedHandlers validates the type-safe API handlers
func TestTypedHandlers(t *testing.T) {
	// Setup test environment
	todoService := services.NewTodoService("test_typed_todos.json")
	dataProvider, err := providers.CreateProvider("mock")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	registry := services.NewPaneRegistry()
	registry.RegisterPane(panes.NewTodoPane(todoService))
	registry.RegisterPane(panes.NewCalendarPane(dataProvider))
	registry.RegisterPane(panes.NewEmailPane(dataProvider))
	
	registry.SetEnabledPanes([]string{"todos", "calendar", "email"})
	registry.SetLayoutConfig(map[string]services.PaneLayoutConfig{
		"todos": {
			GridArea: models.PaneGridArea{Row: "1", Column: "span 2"},
		},
		"calendar": {
			GridArea: models.PaneGridArea{Row: "2", Column: "span 3"},
		},
		"email": {
			GridArea: models.PaneGridArea{Row: "3", Column: "span 2"},
		},
	})

	handler := &Handler{registry: registry}

	// Test TypedTodosAPI
	req := httptest.NewRequest("GET", "/api/typed/todos", nil)
	recorder := httptest.NewRecorder()
	
	handler.TypedTodosAPI(recorder, req)
	
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200 for todos, got %d", recorder.Code)
	}

	// Test TypedCalendarAPI
	req = httptest.NewRequest("GET", "/api/typed/calendar", nil)
	recorder = httptest.NewRecorder()
	
	handler.TypedCalendarAPI(recorder, req)
	
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200 for calendar, got %d", recorder.Code)
	}

	// Test TypedEmailAPI
	req = httptest.NewRequest("GET", "/api/typed/email", nil)
	recorder = httptest.NewRecorder()
	
	handler.TypedEmailAPI(recorder, req)
	
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200 for email, got %d", recorder.Code)
	}
}

// TestTypedHandlerMethodNotAllowed validates method restrictions
func TestTypedHandlerMethodNotAllowed(t *testing.T) {
	registry := services.NewPaneRegistry()
	handler := &Handler{registry: registry}

	req := httptest.NewRequest("POST", "/api/typed/todos", nil)
	recorder := httptest.NewRecorder()
	
	handler.TypedTodosAPI(recorder, req)
	
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405 for POST, got %d", recorder.Code)
	}
}

// TestTypedHandlerPaneNotFound validates error handling
func TestTypedHandlerPaneNotFound(t *testing.T) {
	registry := services.NewPaneRegistry()
	handler := &Handler{registry: registry}

	req := httptest.NewRequest("GET", "/api/typed/todos", nil)
	recorder := httptest.NewRecorder()
	
	handler.TypedTodosAPI(recorder, req)
	
	if recorder.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for missing pane, got %d", recorder.Code)
	}
}