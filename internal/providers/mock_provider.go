package providers

import (
	"time"

	"flexplane/internal/models"
)

// MockProvider implements the Provider interface with mock data
type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (m *MockProvider) GetCalendarEvents() ([]models.Event, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	events := []models.Event{
		{
			ID:       "1",
			Title:    "Team Standup",
			Start:    today.Add(9 * time.Hour),
			End:      today.Add(9*time.Hour + 30*time.Minute),
			Location: "Conference Room A",
		},
		{
			ID:       "2",
			Title:    "Product Review",
			Start:    today.Add(14 * time.Hour),
			End:      today.Add(16 * time.Hour),
			Location: "Conference Room B",
		},
		{
			ID:       "3",
			Title:    "Client Meeting",
			Start:    today.Add(17 * time.Hour),
			End:      today.Add(18 * time.Hour),
			Location: "Zoom",
		},
		{
			ID:       "4",
			Title:    "Sprint Planning",
			Start:    today.Add(24*time.Hour + 10*time.Hour),
			End:      today.Add(24*time.Hour + 12*time.Hour),
			Location: "Conference Room B",
		},
		{
			ID:       "5",
			Title:    "1:1 with Manager",
			Start:    today.Add(24*time.Hour + 15*time.Hour),
			End:      today.Add(24*time.Hour + 15*time.Hour + 30*time.Minute),
			Location: "Office",
		},
	}

	return events, nil
}

func (m *MockProvider) GetEmails() ([]models.Email, error) {
	now := time.Now()

	emails := []models.Email{
		{
			ID:      "1",
			Subject: "Q4 Budget Planning Meeting",
			From:    "sarah.johnson@company.com",
			Preview: "Hi team, I'd like to schedule our Q4 budget planning session. Please review the attached documents...",
			Time:    now.Add(-2 * time.Hour),
			Read:    false,
		},
		{
			ID:      "2",
			Subject: "Project Alpha Update",
			From:    "mike.chen@company.com",
			Preview: "The latest build is ready for testing. We've addressed the performance issues mentioned in the last review...",
			Time:    now.Add(-4 * time.Hour),
			Read:    true,
		},
		{
			ID:      "3",
			Subject: "Weekly Newsletter - Tech Insights",
			From:    "newsletter@techinsights.com",
			Preview: "This week: New developments in AI, cloud security best practices, and the future of remote work...",
			Time:    now.Add(-6 * time.Hour),
			Read:    false,
		},
		{
			ID:      "4",
			Subject: "Re: Client Requirements Discussion",
			From:    "alex.rodriguez@clientcorp.com",
			Preview: "Thanks for the detailed proposal. I've shared it with the stakeholders and we have a few questions...",
			Time:    now.Add(-8 * time.Hour),
			Read:    true,
		},
		{
			ID:      "5",
			Subject: "System Maintenance Window - Sunday",
			From:    "ops-team@company.com",
			Preview: "Scheduled maintenance will occur this Sunday from 2 AM to 6 AM EST. Services may be intermittently unavailable...",
			Time:    now.Add(-12 * time.Hour),
			Read:    false,
		},
	}

	return emails, nil
}