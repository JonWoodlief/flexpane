package providers

import (
	"context"
	"flexplane/internal/models"
)

// DataProvider defines the interface for calendar and email data sources
type DataProvider interface {
	GetCalendarEvents() ([]models.Event, error)
	GetEmails() ([]models.Email, error)
}

// AuthenticatedProvider extends DataProvider with authentication capabilities
type AuthenticatedProvider interface {
	DataProvider
	IsAuthenticated() bool
	GetAuthURL() (string, error)
	Authenticate(ctx context.Context, code string) error
	GetUserInfo() (*UserInfo, error)
}

// UserInfo represents basic user information from the provider
type UserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// This interface allows easy swapping between mock and real providers
// For example:
// - MockProvider (current implementation)
// - GmailProvider (Google OAuth integration)
// - OutlookProvider (future Microsoft Graph API integration)