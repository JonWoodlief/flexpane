package providers

import "flexpane/internal/models"

// DataProvider defines the interface for calendar and email data sources
type DataProvider interface {
	GetCalendarEvents() ([]models.Event, error)
	GetEmails() ([]models.Email, error)
}

// TypedProvider is a generic interface for type-safe data providers
type TypedProvider[T any] interface {
	GetData() (T, error)
}

// CalendarProvider provides calendar-specific data
type CalendarProvider interface {
	TypedProvider[[]models.Event]
}

// EmailProvider provides email-specific data  
type EmailProvider interface {
	TypedProvider[[]models.Email]
}

// This interface allows easy swapping between mock and real providers