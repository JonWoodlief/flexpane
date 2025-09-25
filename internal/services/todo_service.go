package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"flexplane/internal/models"
)

// TodoService handles todo operations independently of the provider system
type TodoService struct {
	filename string
	todos    []models.Todo
	mutex    sync.RWMutex
}

// Ensure TodoService implements the interface
var _ models.TodoService = (*TodoService)(nil)

func NewTodoService(filename string) *TodoService {
	service := &TodoService{
		filename: filename,
		todos:    []models.Todo{},
	}
	service.load()
	return service
}

func (s *TodoService) GetTodos() []models.Todo {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.todos
}

func (s *TodoService) AddTodo(message string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.todos = append(s.todos, models.Todo{
		Done:    false,
		Message: message,
	})

	return s.save()
}

func (s *TodoService) ToggleTodo(index int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if index < 0 || index >= len(s.todos) {
		return nil // Invalid index, ignore
	}

	s.todos[index].Done = !s.todos[index].Done
	return s.save()
}

func (s *TodoService) load() error {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(s.filename), 0755); err != nil {
		return err
	}

	data, err := os.ReadFile(s.filename)
	if os.IsNotExist(err) {
		// File doesn't exist, start with empty todos
		s.todos = []models.Todo{}
		return s.save() // Create the file
	}
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &s.todos)
}

func (s *TodoService) save() error {
	data, err := json.MarshalIndent(s.todos, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, data, 0644)
}