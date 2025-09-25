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
	h.handlePaneAPI("todos", w, r)
}

// handlePaneAPI provides a generic API handler for panes that implement APIHandler
func (h *Handler) handlePaneAPI(paneID string, w http.ResponseWriter, r *http.Request) {
	pane, exists := h.registry.GetPane(paneID)
	if !exists {
		http.Error(w, "Pane not found", 404)
		return
	}

	// Check if pane supports API operations
	if apiHandler, ok := pane.(models.APIHandler); ok {
		if err := apiHandler.HandleAPI(w, r); err != nil {
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}

	// Fallback for panes without API support
	if r.Method == "GET" {
		data, err := pane.GetData(r.Context())
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
		return
	}
	
	http.Error(w, "Method Not Allowed", 405)
}