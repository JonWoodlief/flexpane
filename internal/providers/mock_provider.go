package providers

import (
	"flexplane/internal/models"
	"time"
)

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (m *MockProvider) GetCalendarEvents() ([]models.Event, error) {
	now := time.Now()
	return []models.Event{
		{ID: "1", Title: "Team Standup", Start: now.Add(time.Hour), End: now.Add(time.Hour + 30*time.Minute), Location: "Conference Room A"},
		{ID: "2", Title: "Product Review", Start: now.Add(2 * time.Hour), End: now.Add(3 * time.Hour), Location: "Zoom"},
		{ID: "3", Title: "Client Call", Start: now.Add(4 * time.Hour), End: now.Add(4*time.Hour + 45*time.Minute), Location: "Phone"},
	}, nil
}

func (m *MockProvider) GetEmails() ([]models.Email, error) {
	now := time.Now()
	return []models.Email{
		{ID: "1", Subject: "Budget Meeting", From: "sarah@company.com", Preview: "Q4 planning...", Time: now.Add(-2 * time.Hour), Read: false},
		{ID: "2", Subject: "Project Update", From: "mike@company.com", Preview: "Latest build ready...", Time: now.Add(-4 * time.Hour), Read: true},
		{ID: "3", Subject: "Newsletter", From: "news@tech.com", Preview: "AI developments...", Time: now.Add(-30 * time.Minute), Read: false},
	}, nil
}