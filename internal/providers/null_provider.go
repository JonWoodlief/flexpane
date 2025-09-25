package providers

import "flexplane/internal/models"

// NullProvider implements Provider interface but returns empty data
// This is used for production when real integrations aren't configured yet
type NullProvider struct{}

func NewNullProvider() *NullProvider {
	return &NullProvider{}
}

func (np *NullProvider) GetCalendarEvents() ([]models.Event, error) {
	// Return empty events - no calendar integration configured
	return []models.Event{}, nil
}

func (np *NullProvider) GetEmails() ([]models.Email, error) {
	// Return empty emails - no email integration configured
	return []models.Email{}, nil
}

func (np *NullProvider) GetTodos() []models.Todo {
	// Return empty todos - this provider doesn't handle todos
	return []models.Todo{}
}

func (np *NullProvider) AddTodo(message string) error {
	// No-op - this provider doesn't handle todos
	return nil
}

func (np *NullProvider) ToggleTodo(index int) error {
	// No-op - this provider doesn't handle todos
	return nil
}