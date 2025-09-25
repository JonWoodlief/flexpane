package providers

import (
	"flexplane/internal/models"
)

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (m *MockProvider) GetCalendarEvents() ([]models.Event, error) {
	return []models.Event{
		{ID: "1", Title: "Team Standup", Location: "Conference Room A"},
		{ID: "2", Title: "Product Review", Location: "Zoom"},
		{ID: "3", Title: "Client Call", Location: "Phone"},
	}, nil
}

func (m *MockProvider) GetEmails() ([]models.Email, error) {
	return []models.Email{
		{ID: "1", Subject: "Budget Meeting", From: "sarah@company.com", Preview: "Q4 planning...", Read: false},
		{ID: "2", Subject: "Project Update", From: "mike@company.com", Preview: "Latest build ready...", Read: true},
		{ID: "3", Subject: "Newsletter", From: "news@tech.com", Preview: "AI developments...", Read: false},
	}, nil
}