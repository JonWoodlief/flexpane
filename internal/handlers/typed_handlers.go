package handlers

import (
	"encoding/json"
	"net/http"

	"flexpane/internal/models"
	"flexpane/internal/services"
)

// HandleTypedPaneAPI provides type-safe API handling for typed panes
func HandleTypedPaneAPI[T any](
	registry *services.PaneRegistry,
	paneID string,
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	// Use the generic function to get typed pane data
	typedPaneData, err := services.GetTypedPaneData[T](r.Context(), registry, paneID)
	if err != nil {
		http.Error(w, "Pane not found", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(typedPaneData); err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

// TypedTodosAPI is a type-safe version of the todos API
func (h *Handler) TypedTodosAPI(w http.ResponseWriter, r *http.Request) {
	HandleTypedPaneAPI[models.TodoPaneData](h.registry, "todos", w, r)
}

// TypedCalendarAPI is a type-safe version of the calendar API
func (h *Handler) TypedCalendarAPI(w http.ResponseWriter, r *http.Request) {
	HandleTypedPaneAPI[models.CalendarPaneData](h.registry, "calendar", w, r)
}

// TypedEmailAPI is a type-safe version of the email API  
func (h *Handler) TypedEmailAPI(w http.ResponseWriter, r *http.Request) {
	HandleTypedPaneAPI[models.EmailPaneData](h.registry, "email", w, r)
}