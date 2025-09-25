package providers

import "flexplane/internal/models"

// DataProvider defines the interface for calendar and email data sources
type DataProvider interface {
	GetCalendarEvents() ([]models.Event, error)
	GetEmails() ([]models.Email, error)
}

// TodoProvider defines the interface for todo data sources
type TodoProvider interface {
	GetTodos() []models.Todo
	AddTodo(message string) error
	ToggleTodo(index int) error
}

// Provider defines a unified interface for all data providers
// This allows for comprehensive provider implementations that can handle multiple data types
type Provider interface {
	DataProvider
	TodoProvider
}

// This interface allows easy swapping between mock and real providers
// For example:
// - MockProvider (current implementation)
// - OutlookProvider (future Microsoft Graph API integration)
// - GmailProvider (future Gmail API integration)
// - TodoFileProvider (file-based todo storage)
// - TodoDatabaseProvider (database-based todo storage)