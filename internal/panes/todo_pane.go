package panes

import (
	"context"

	"flexplane/internal/services"
)

// TodoPane implements the Pane interface for todo items
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

func (tp *TodoPane) GetData(ctx context.Context) (interface{}, error) {
	todos := tp.todoService.GetTodos()

	return map[string]interface{}{
		"Todos": todos,
		"Count": len(todos),
	}, nil
}

// AddTodo adds a new todo item
func (tp *TodoPane) AddTodo(message string) error {
	return tp.todoService.AddTodo(message)
}

// ToggleTodo toggles the done status of a todo item by index
func (tp *TodoPane) ToggleTodo(index int) error {
	return tp.todoService.ToggleTodo(index)
}