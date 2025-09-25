package providers

import (
	"os"
	"testing"
)

func TestMockProvider(t *testing.T) {
	provider := NewMockProvider()

	// Test calendar events
	events, err := provider.GetCalendarEvents()
	if err != nil {
		t.Errorf("GetCalendarEvents failed: %v", err)
	}
	if len(events) == 0 {
		t.Error("Expected mock calendar events, got none")
	}

	// Test emails
	emails, err := provider.GetEmails()
	if err != nil {
		t.Errorf("GetEmails failed: %v", err)
	}
	if len(emails) == 0 {
		t.Error("Expected mock emails, got none")
	}

	// Test todos (should be empty for mock)
	todos := provider.GetTodos()
	if len(todos) != 0 {
		t.Errorf("Expected empty todos from mock provider, got %d", len(todos))
	}

	// Test add todo (should not error but does nothing)
	err = provider.AddTodo("test todo")
	if err != nil {
		t.Errorf("AddTodo failed: %v", err)
	}

	// Test toggle todo (should not error but does nothing)
	err = provider.ToggleTodo(0)
	if err != nil {
		t.Errorf("ToggleTodo failed: %v", err)
	}
}

func TestTodoFileProvider(t *testing.T) {
	testFile := "test_provider_todos.json"
	defer os.Remove(testFile) // cleanup

	provider := NewTodoFileProvider(testFile)

	// Initially should have no todos
	todos := provider.GetTodos()
	if len(todos) != 0 {
		t.Errorf("Expected 0 todos initially, got %d", len(todos))
	}

	// Add a todo
	err := provider.AddTodo("Test todo item")
	if err != nil {
		t.Errorf("AddTodo failed: %v", err)
	}

	// Should now have one todo
	todos = provider.GetTodos()
	if len(todos) != 1 {
		t.Errorf("Expected 1 todo after adding, got %d", len(todos))
	}

	if todos[0].Message != "Test todo item" {
		t.Errorf("Expected 'Test todo item', got '%s'", todos[0].Message)
	}

	if todos[0].Done {
		t.Error("Expected new todo to be not done")
	}

	// Toggle the todo
	err = provider.ToggleTodo(0)
	if err != nil {
		t.Errorf("ToggleTodo failed: %v", err)
	}

	todos = provider.GetTodos()
	if !todos[0].Done {
		t.Error("Expected todo to be done after toggle")
	}

	// Toggle again
	err = provider.ToggleTodo(0)
	if err != nil {
		t.Errorf("ToggleTodo failed: %v", err)
	}

	todos = provider.GetTodos()
	if todos[0].Done {
		t.Error("Expected todo to be not done after second toggle")
	}

	// Test invalid index
	err = provider.ToggleTodo(999)
	if err != nil {
		t.Errorf("ToggleTodo with invalid index should not error: %v", err)
	}
}

func TestCompositeProvider(t *testing.T) {
	mockProvider := NewMockProvider()
	todoProvider := NewTodoFileProvider("test_composite_todos.json")
	defer os.Remove("test_composite_todos.json") // cleanup

	composite := NewCompositeProvider(mockProvider, todoProvider)

	// Test that it correctly delegates to data provider
	events, err := composite.GetCalendarEvents()
	if err != nil {
		t.Errorf("GetCalendarEvents failed: %v", err)
	}
	if len(events) == 0 {
		t.Error("Expected calendar events from composite provider")
	}

	emails, err := composite.GetEmails()
	if err != nil {
		t.Errorf("GetEmails failed: %v", err)
	}
	if len(emails) == 0 {
		t.Error("Expected emails from composite provider")
	}

	// Test that it correctly delegates to todo provider
	todos := composite.GetTodos()
	if len(todos) != 0 {
		t.Errorf("Expected 0 todos initially from composite, got %d", len(todos))
	}

	err = composite.AddTodo("Composite test todo")
	if err != nil {
		t.Errorf("AddTodo failed on composite: %v", err)
	}

	todos = composite.GetTodos()
	if len(todos) != 1 {
		t.Errorf("Expected 1 todo after adding to composite, got %d", len(todos))
	}

	err = composite.ToggleTodo(0)
	if err != nil {
		t.Errorf("ToggleTodo failed on composite: %v", err)
	}

	todos = composite.GetTodos()
	if !todos[0].Done {
		t.Error("Expected todo to be done after toggle on composite")
	}
}

func TestProviderFactory(t *testing.T) {
	factory := NewProviderFactory()

	// Test that built-in providers are registered
	available := factory.GetAvailableProviders()
	expectedTypes := []string{"mock", "file"}
	
	for _, expected := range expectedTypes {
		found := false
		for _, available := range available {
			if available == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected provider type '%s' to be available", expected)
		}
	}

	// Test creating mock provider
	mockProvider, err := factory.CreateProvider(ProviderConfig{Type: "mock"})
	if err != nil {
		t.Errorf("Failed to create mock provider: %v", err)
	}
	if mockProvider == nil {
		t.Error("Expected non-nil mock provider")
	}

	// Test creating file provider
	fileProvider, err := factory.CreateProvider(ProviderConfig{
		Type: "file",
		Args: map[string]interface{}{
			"todo_file": "test_factory_todos.json",
		},
	})
	defer os.Remove("test_factory_todos.json") // cleanup
	
	if err != nil {
		t.Errorf("Failed to create file provider: %v", err)
	}
	if fileProvider == nil {
		t.Error("Expected non-nil file provider")
	}

	// Test creating unknown provider
	_, err = factory.CreateProvider(ProviderConfig{Type: "unknown"})
	if err == nil {
		t.Error("Expected error when creating unknown provider type")
	}

	// Test custom provider registration
	factory.RegisterProvider("custom", func(args map[string]interface{}) (Provider, error) {
		return NewMockProvider(), nil
	})

	customProvider, err := factory.CreateProvider(ProviderConfig{Type: "custom"})
	if err != nil {
		t.Errorf("Failed to create custom provider: %v", err)
	}
	if customProvider == nil {
		t.Error("Expected non-nil custom provider")
	}
}