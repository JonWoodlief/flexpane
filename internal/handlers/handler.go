package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"flexplane/internal/models"
	"flexplane/internal/panes"
	"flexplane/internal/services"
)

type Handler struct {
	registry  *services.PaneRegistry
	templates *template.Template
}

func NewHandler(registry *services.PaneRegistry, templates *template.Template) *Handler {
	return &Handler{
		registry:  registry,
		templates: templates,
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get all enabled panes with their data
	panes, err := h.registry.GetEnabledPanes(ctx)
	if err != nil {
		log.Printf("Error getting enabled panes: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prepare template data
	data := models.PageData{
		Panes: panes,
	}

	// Render template
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "layout.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// TODO: CONCURRENCY BUG - Index-based operations are unsafe with concurrent
// reordering. Need to add unique IDs or implement proper locking before
// multi-user support or background sync.
func (h *Handler) TodosAPI(w http.ResponseWriter, r *http.Request) {
	// Get the todos pane from registry
	todosPane, exists := h.registry.GetPane("todos")
	if !exists {
		log.Printf("Todos pane not found in registry")
		http.Error(w, "Todos pane not found", http.StatusNotFound)
		return
	}

	// Cast to TodoPane to access service
	// TODO: Better way to handle pane-specific APIs
	switch r.Method {
	case "GET":
		data, err := todosPane.GetData(r.Context())
		if err != nil {
			log.Printf("Error getting todos data: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("Error encoding JSON response: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

	case "POST":
		h.handleAddTodo(w, r)

	case "PATCH":
		h.handleToggleTodo(w, r)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleAddTodo(w http.ResponseWriter, r *http.Request) {
	// Limit request size to prevent DoS
	r.Body = http.MaxBytesReader(w, r.Body, 1024) // 1KB max
	
	var req struct {
		Message string `json:"message"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding add todo request: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	// Validate message
	if req.Message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}
	
	if len(req.Message) > 200 {
		http.Error(w, "Message too long (max 200 characters)", http.StatusBadRequest)
		return
	}
	
	// Get the TodoPane and service
	todosPane, exists := h.registry.GetPane("todos")
	if !exists {
		http.Error(w, "Todos pane not found", http.StatusNotFound)
		return
	}
	
	// This is a type assertion - in production, would use a better interface
	if todoPane, ok := todosPane.(*panes.TodoPane); ok {
		if err := todoPane.AddTodo(req.Message); err != nil {
			log.Printf("Error adding todo: %v", err)
			http.Error(w, "Failed to add todo", http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "created"})
	} else {
		http.Error(w, "Invalid pane type", http.StatusInternalServerError)
	}
}

func (h *Handler) handleToggleTodo(w http.ResponseWriter, r *http.Request) {
	// Get index from query parameters
	indexStr := r.URL.Query().Get("index")
	if indexStr == "" {
		http.Error(w, "Index parameter is required", http.StatusBadRequest)
		return
	}
	
	index := 0
	if _, err := fmt.Sscanf(indexStr, "%d", &index); err != nil {
		http.Error(w, "Invalid index parameter", http.StatusBadRequest)
		return
	}
	
	if index < 0 {
		http.Error(w, "Index must be non-negative", http.StatusBadRequest)
		return
	}
	
	// Get the TodoPane and service
	todosPane, exists := h.registry.GetPane("todos")
	if !exists {
		http.Error(w, "Todos pane not found", http.StatusNotFound)
		return
	}
	
	// This is a type assertion - in production, would use a better interface
	if todoPane, ok := todosPane.(*panes.TodoPane); ok {
		if err := todoPane.ToggleTodo(index); err != nil {
			log.Printf("Error toggling todo: %v", err)
			http.Error(w, "Failed to toggle todo", http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
	} else {
		http.Error(w, "Invalid pane type", http.StatusInternalServerError)
	}
}