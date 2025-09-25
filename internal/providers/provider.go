package providers

import "flexplane/internal/models"

// DataProvider defines the interface for calendar and email data sources
// This is designed for unified providers like Outlook/Gmail that handle both
// calendar events and emails with shared OAuth authentication
type DataProvider interface {
	GetCalendarEvents() ([]models.Event, error)
	GetEmails() ([]models.Email, error)
}

// This interface allows easy swapping between different email/calendar providers
// For example:
// - MockProvider (development/testing with fake data)
// - OutlookProvider (Microsoft Graph API integration)  
// - GmailProvider (Gmail API integration)
// - NullProvider (empty data when no integration is configured)