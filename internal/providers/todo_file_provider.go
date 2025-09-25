package providers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"flexplane/internal/models"
)

// TodoFileProvider implements TodoProvider interface with file-based storage
type TodoFileProvider struct {
	filename string
	todos    []models.Todo
	mutex    sync.RWMutex
}

func NewTodoFileProvider(filename string) *TodoFileProvider {
	provider := &TodoFileProvider{
		filename: filename,
		todos:    []models.Todo{},
	}
	provider.load()
	return provider
}

func (p *TodoFileProvider) GetTodos() []models.Todo {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.todos
}

func (p *TodoFileProvider) AddTodo(message string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.todos = append(p.todos, models.Todo{
		Done:    false,
		Message: message,
	})

	return p.save()
}

func (p *TodoFileProvider) ToggleTodo(index int) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if index < 0 || index >= len(p.todos) {
		return nil // Invalid index, ignore
	}

	p.todos[index].Done = !p.todos[index].Done
	return p.save()
}

func (p *TodoFileProvider) load() error {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(p.filename), 0755); err != nil {
		return err
	}

	data, err := os.ReadFile(p.filename)
	if os.IsNotExist(err) {
		// File doesn't exist, start with empty todos
		p.todos = []models.Todo{}
		return p.save() // Create the file
	}
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &p.todos)
}

func (p *TodoFileProvider) save() error {
	data, err := json.MarshalIndent(p.todos, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(p.filename, data, 0644)
}