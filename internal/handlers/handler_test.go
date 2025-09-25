package handlers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"flexplane/internal/models"
	"flexplane/internal/panes"
	"flexplane/internal/services"
)

// MockDataProvider for testing
type MockDataProvider struct {
	events []models.Event
	emails []models.Email
	err    error
}

func (m *MockDataProvider) GetCalendarEvents() ([]models.Event, error) {
	return m.events, m.err
}

func (m *MockDataProvider) GetEmails() ([]models.Email, error) {
	return m.emails, m.err
}

func setupTestHandler(t *testing.T) *Handler {
	// Create test templates
	tmpl := template.New("test")

	// Add minimal layout template
	layoutTemplate := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
{{range .Panes}}
<div class="pane" data-pane-id="{{.ID}}">
    <h2>{{.Title}}</h2>
    <div>{{.Data.Count}} items</div>
</div>
{{end}}
</body>
</html>`

	tmpl = template.Must(tmpl.New("layout.html").Parse(layoutTemplate))

	// Create registry with mock panes
	registry := services.NewPaneRegistry()

	// Use test provider
	mockProvider := &MockDataProvider{
		events: []models.Event{
			{ID: "1", Title: "Test Event"},
		},
	}

	// Register test panes
	registry.RegisterPane(panes.NewCalendarPane(mockProvider))
	registry.SetEnabledPanes([]string{"calendar"})

	// Set layout config
	registry.SetLayoutConfig(map[string]services.PaneLayoutConfig{
		"calendar": {
			GridArea: models.PaneGridArea{Row: "1", Column: "span 3"},
		},
	})

	return NewHandler(registry, tmpl)
}

func TestHandler_Home_Success(t *testing.T) {
	handler := setupTestHandler(t)

	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	handler.Home(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

	contentType := recorder.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected HTML content type, got %s", contentType)
	}

	body := recorder.Body.String()
	if body == "" {
		t.Error("Expected non-empty response body")
	}

	// Check that the pane is rendered
	if !containsString(body, "Calendar") {
		t.Error("Expected response to contain 'Calendar'")
	}
}

func TestHandler_TodosAPI_PaneNotFound(t *testing.T) {
	handler := setupTestHandler(t)

	req := httptest.NewRequest("GET", "/api/todos", nil)
	recorder := httptest.NewRecorder()

	handler.TodosAPI(recorder, req)

	// Should return 404 since we didn't register todos pane
	if recorder.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", recorder.Code)
	}

	body := recorder.Body.String()
	if body != "Todos pane not found\n" {
		t.Errorf("Expected 'Todos pane not found', got %s", body)
	}
}

func TestHandler_TodosAPI_WithTodosPane(t *testing.T) {
	// Setup handler with todos pane
	tmpl := template.New("test")
	layoutTemplate := `<div>{{range .Panes}}{{.Title}}{{end}}</div>`
	tmpl = template.Must(tmpl.New("layout.html").Parse(layoutTemplate))

	registry := services.NewPaneRegistry()
	todoService := services.NewTodoService("test_todos.json")
	registry.RegisterPane(panes.NewTodoPane(todoService))
	registry.SetEnabledPanes([]string{"todos"})

	// Set layout config
	registry.SetLayoutConfig(map[string]services.PaneLayoutConfig{
		"todos": {
			GridArea: models.PaneGridArea{Row: "2", Column: "span 6"},
		},
	})

	handler := NewHandler(registry, tmpl)

	req := httptest.NewRequest("GET", "/api/todos", nil)
	recorder := httptest.NewRecorder()

	handler.TodosAPI(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

	contentType := recorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected JSON content type, got %s", contentType)
	}
}

func TestHandler_TodosAPI_MethodNotAllowed(t *testing.T) {
	// This test needs a handler with todos pane registered
	tmpl := template.New("test")
	layoutTemplate := `<div>test</div>`
	tmpl = template.Must(tmpl.New("layout.html").Parse(layoutTemplate))

	registry := services.NewPaneRegistry()
	todoService := services.NewTodoService("test_method_not_allowed.json")
	registry.RegisterPane(panes.NewTodoPane(todoService))
	registry.SetEnabledPanes([]string{"todos"})

	// Set layout config
	registry.SetLayoutConfig(map[string]services.PaneLayoutConfig{
		"todos": {
			GridArea: models.PaneGridArea{Row: "2", Column: "span 6"},
		},
	})

	handler := NewHandler(registry, tmpl)

	req := httptest.NewRequest("DELETE", "/api/todos", nil)
	recorder := httptest.NewRecorder()

	handler.TodosAPI(recorder, req)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", recorder.Code)
	}
}

// Helper function for string contains check
func containsString(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}