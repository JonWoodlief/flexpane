package panes

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"flexpane/internal/models"
	"flexpane/internal/services"
)

// TodoPane implements both Pane and TypedPane interfaces for todo items
// The generic TypedPane provides compile-time type safety for the TodoPaneData
type TodoPane struct {
	todoService *services.TodoService
}

func NewTodoPane(todoService *services.TodoService) *TodoPane {
	return &TodoPane{
		todoService: todoService,
	}
}

func (tp *TodoPane) ID() string {
	return "todos"
}

func (tp *TodoPane) Title() string {
	return "Todos"
}

func (tp *TodoPane) Template() string {
	return "panes/todos.html"
}

// GetData maintains backward compatibility by returning interface{}
func (tp *TodoPane) GetData(ctx context.Context) (interface{}, error) {
	return tp.GetTypedData(ctx)
}

// GetTypedData provides type-safe access to todo data
// This eliminates the need for type assertions in calling code
func (tp *TodoPane) GetTypedData(ctx context.Context) (models.TodoPaneData, error) {
	todos := tp.todoService.GetTodos()

	return models.TodoPaneData{
		Todos: todos,
		Count: len(todos),
	}, nil
}

// HandleAPI implements the APIHandler interface for todo-specific operations
func (tp *TodoPane) HandleAPI(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		data, err := tp.GetData(r.Context())
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		return json.NewEncoder(w).Encode(data)

	case "POST":
		return tp.handleAddTodo(w, r)

	case "PATCH":
		return tp.handleToggleTodo(w, r)

	default:
		http.Error(w, "Method Not Allowed", 405)
		return nil
	}
}

func (tp *TodoPane) handleAddTodo(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Message string `json:"message"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return nil
	}
	
	if req.Message == "" {
		http.Error(w, "Message required", 400)
		return nil
	}
	
	if err := tp.todoService.AddTodo(req.Message); err != nil {
		return err
	}
	
	w.WriteHeader(201)
	return json.NewEncoder(w).Encode(map[string]string{"status": "created"})
}

func (tp *TodoPane) handleToggleTodo(w http.ResponseWriter, r *http.Request) error {
	indexStr := r.URL.Query().Get("index")
	if indexStr == "" {
		http.Error(w, "Index required", 400)
		return nil
	}
	
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 {
		http.Error(w, "Invalid index", 400)
		return nil
	}
	
	if err := tp.todoService.ToggleTodo(index); err != nil {
		return err
	}
	
	return json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}