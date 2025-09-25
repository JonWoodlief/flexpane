package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"flexpane/internal/handlers"
	"flexpane/internal/panes"
	"flexpane/internal/providers"
	"flexpane/internal/services"
)

// Integration tests - test the full application flow
func TestFullApplication_HomePage(t *testing.T) {
	// Setup full application like main.go
	todoService := services.NewTodoService("test_integration_todos.json")

	// Create data provider
	dataProvider, err := providers.CreateProvider("mock")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// Create registry
	registry := services.NewPaneRegistry()
	registry.RegisterPane(panes.NewCalendarPane(dataProvider))
	registry.RegisterPane(panes.NewTodoPane(todoService))
	registry.RegisterPane(panes.NewEmailPane(dataProvider))
	registry.SetEnabledPanes([]string{"calendar", "todos", "email"})

	// Use real templates (this will help catch template errors)
	templates, err := loadTemplates()
	if err != nil {
		t.Skipf("Skipping integration test - templates not available: %v", err)
	}

	handler := handlers.NewHandler(registry, templates)

	// Test home page
	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	handler.Home(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
		t.Logf("Response body: %s", recorder.Body.String())
	}

	body := recorder.Body.String()

	// Check all panes are present
	expectedContent := []string{
		"Flexpane",     // Page title
		"Calendar",      // Calendar pane
		"Todos",         // Todo pane
		"Email Preview", // Email pane
	}

	for _, expected := range expectedContent {
		if !strings.Contains(body, expected) {
			t.Errorf("Expected response to contain '%s'", expected)
		}
	}
}

func TestFullApplication_TodosAPI(t *testing.T) {
	// Setup
	todoService := services.NewTodoService("test_integration_todos_api.json")
	registry := services.NewPaneRegistry()
	registry.RegisterPane(panes.NewTodoPane(todoService))
	registry.SetEnabledPanes([]string{"todos"})

	templates, err := loadTemplates()
	if err != nil {
		t.Skipf("Skipping integration test - templates not available: %v", err)
	}

	handler := handlers.NewHandler(registry, templates)

	// Test GET /api/todos
	req := httptest.NewRequest("GET", "/api/todos", nil)
	recorder := httptest.NewRecorder()

	handler.TodosAPI(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

	// Parse response
	var response map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	// Check structure
	if _, exists := response["Todos"]; !exists {
		t.Error("Expected 'Todos' field in response")
	}

	if _, exists := response["Count"]; !exists {
		t.Error("Expected 'Count' field in response")
	}
}

func TestFullApplication_StaticFiles(t *testing.T) {
	// Test that static file serving works
	req := httptest.NewRequest("GET", "/static/css/style.css", nil)
	recorder := httptest.NewRecorder()

	// Use the same static file setup as main.go
	fs := http.FileServer(http.Dir("web/static/"))
	handler := http.StripPrefix("/static/", fs)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200 for CSS file, got %d", recorder.Code)
	}

	contentType := recorder.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/css") {
		t.Errorf("Expected CSS content type, got %s", contentType)
	}
}

// Removed redundant configuration tests - covered by unit tests

// Helper function to load templates like main.go
func loadTemplates() (*template.Template, error) {
	tmpl, err := template.ParseGlob("web/templates/*.html")
	if err != nil {
		return nil, err
	}
	tmpl, err = tmpl.ParseGlob("web/templates/components/*.html")
	if err != nil {
		return nil, err
	}
	tmpl, err = tmpl.ParseGlob("web/templates/panes/*.html")
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func TestMain(m *testing.M) {
	// Setup: Clean up test files before running tests
	testFiles := []string{
		"test_integration_todos.json",
		"test_integration_todos_api.json",
		"test_config_todos.json",
		"test_todos.json",
	}

	for _, file := range testFiles {
		os.Remove(file)
	}

	// Run tests
	code := m.Run()

	// Cleanup: Remove test files after running tests
	for _, file := range testFiles {
		os.Remove(file)
	}

	os.Exit(code)
}