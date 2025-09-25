package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"flexplane/internal/models"
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
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Prepare template data
	data := models.PageData{
		Panes: panes,
	}

	// Render template
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, "Internal Server Error", 500)
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
		http.Error(w, "Todos pane not found", 404)
		return
	}

	// Cast to TodoPane to access service
	// TODO: Better way to handle pane-specific APIs
	switch r.Method {
	case "GET":
		data, err := todosPane.GetData(r.Context())
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)

	case "POST":
		// TODO: Implement add todo via pane interface
		http.Error(w, "Not implemented yet", 501)

	case "PATCH":
		// TODO: Implement toggle todo via pane interface
		http.Error(w, "Not implemented yet", 501)

	default:
		http.Error(w, "Method Not Allowed", 405)
	}
}