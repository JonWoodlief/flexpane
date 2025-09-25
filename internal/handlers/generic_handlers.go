package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"flexpane/internal/models"
)

// GenericAPIHandler provides type-safe API handling with compile-time guarantees
// This eliminates the need for manual type assertions and improves error handling
type GenericAPIHandler[TReq, TResp any] struct {
	handler func(ctx context.Context, req TReq) (TResp, error)
}

// NewGenericAPIHandler creates a new type-safe API handler
func NewGenericAPIHandler[TReq, TResp any](handler func(ctx context.Context, req TReq) (TResp, error)) *GenericAPIHandler[TReq, TResp] {
	return &GenericAPIHandler[TReq, TResp]{
		handler: handler,
	}
}

// HandleHTTP provides HTTP handling with automatic JSON marshaling/unmarshaling
// This eliminates boilerplate code and provides type safety
func (h *GenericAPIHandler[TReq, TResp]) HandleHTTP(w http.ResponseWriter, r *http.Request) error {
	var req TReq
	
	// Only decode JSON for requests with body
	if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return nil
		}
	}
	
	// Call the type-safe handler
	resp, err := h.handler(r.Context(), req)
	if err != nil {
		return fmt.Errorf("handler error: %w", err)
	}
	
	// Automatically encode the response
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resp)
}

// TypedPaneAPIHandler provides a bridge between typed panes and HTTP APIs
// This shows how generics can eliminate boilerplate in API handling
func TypedPaneAPIHandler[T any](pane models.TypedPane[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			data, err := pane.GetTypedData(r.Context())
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(data); err != nil {
				http.Error(w, "Encoding Error", http.StatusInternalServerError)
				return
			}
			
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

// Example of how generic handlers could be used for specific pane types
type AddTodoRequest struct {
	Message string `json:"message"`
}

type AddTodoResponse struct {
	Status string `json:"status"`
	ID     string `json:"id,omitempty"`
}

type ToggleTodoRequest struct {
	Index int `json:"index"`
}

type ToggleTodoResponse struct {
	Status string `json:"status"`
	Done   bool   `json:"done"`
}

// These would be methods on TodoPane to demonstrate typed API handlers
// func (tp *TodoPane) HandleAddTodo() *GenericAPIHandler[AddTodoRequest, AddTodoResponse] {
//     return NewGenericAPIHandler(func(ctx context.Context, req AddTodoRequest) (AddTodoResponse, error) {
//         if req.Message == "" {
//             return AddTodoResponse{}, fmt.Errorf("message required")
//         }
//         
//         if err := tp.todoService.AddTodo(req.Message); err != nil {
//             return AddTodoResponse{}, err
//         }
//         
//         return AddTodoResponse{Status: "created"}, nil
//     })
// }