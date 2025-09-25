package providers

import "flexplane/internal/models"

// NullProvider implements DataProvider interface but returns empty data
// This is used for production when real email/calendar integrations aren't configured yet
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