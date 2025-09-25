package providers

import "flexplane/internal/models"

// DataProvider defines the interface for calendar and email data sources
type DataProvider interface {
	GetCalendarEvents() ([]models.Event, error)
	GetEmails() ([]models.Email, error)
}

// This interface allows easy swapping between mock and real providers
// For example:
// - MockProvider (current implementation)
// - OutlookProvider (future Microsoft Graph API integration)
// - GmailProvider (future Gmail API integration)