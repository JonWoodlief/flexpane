package providers

import "flexpane/internal/models"

// DataProvider defines the interface for calendar and email data sources
type DataProvider interface {
	GetCalendarEvents() ([]models.Event, error)
	GetEmails() ([]models.Email, error)
}

// This interface allows easy swapping between mock and real providers