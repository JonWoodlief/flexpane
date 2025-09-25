package models

// TodoService defines the interface for todo operations
type TodoService interface {
	GetTodos() []Todo
	AddTodo(message string) error
	ToggleTodo(index int) error
}