package providers

import "flexplane/internal/models"

// CompositeProvider combines different providers to create a unified provider
type CompositeProvider struct {
	dataProvider DataProvider
	todoProvider TodoProvider
}

func NewCompositeProvider(dataProvider DataProvider, todoProvider TodoProvider) *CompositeProvider {
	return &CompositeProvider{
		dataProvider: dataProvider,
		todoProvider: todoProvider,
	}
}

// DataProvider methods
func (cp *CompositeProvider) GetCalendarEvents() ([]models.Event, error) {
	return cp.dataProvider.GetCalendarEvents()
}

func (cp *CompositeProvider) GetEmails() ([]models.Email, error) {
	return cp.dataProvider.GetEmails()
}

// TodoProvider methods
func (cp *CompositeProvider) GetTodos() []models.Todo {
	return cp.todoProvider.GetTodos()
}

func (cp *CompositeProvider) AddTodo(message string) error {
	return cp.todoProvider.AddTodo(message)
}

func (cp *CompositeProvider) ToggleTodo(index int) error {
	return cp.todoProvider.ToggleTodo(index)
}