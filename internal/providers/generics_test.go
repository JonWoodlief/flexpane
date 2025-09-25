package providers

import (
	"testing"
	"flexpane/internal/models"
)

// TestTypedProviders validates the generic provider functionality
func TestTypedProviders(t *testing.T) {
	// Create a mock provider
	mockProvider := NewMockProvider()
	
	// Test CalendarProvider
	calendarProvider := NewCalendarProvider(mockProvider)
	events, err := calendarProvider.GetData()
	if err != nil {
		t.Errorf("Unexpected error from calendar provider: %v", err)
	}
	
	if len(events) != 3 {
		t.Errorf("Expected 3 events, got %d", len(events))
	}
	
	if events[0].Title != "Team Standup" {
		t.Errorf("Expected first event title 'Team Standup', got %s", events[0].Title)
	}

	// Test EmailProvider
	emailProvider := NewEmailProvider(mockProvider)
	emails, err := emailProvider.GetData()
	if err != nil {
		t.Errorf("Unexpected error from email provider: %v", err)
	}
	
	if len(emails) != 3 {
		t.Errorf("Expected 3 emails, got %d", len(emails))
	}
	
	if emails[0].Subject != "Budget Meeting" {
		t.Errorf("Expected first email subject 'Budget Meeting', got %s", emails[0].Subject)
	}
}

// TestCreateTypedProvider validates the generic factory pattern
func TestCreateTypedProvider(t *testing.T) {
	// Test creating a calendar provider
	calendarProvider, err := CreateTypedProvider[[]models.Event]("mock")
	if err != nil {
		t.Errorf("Unexpected error creating calendar provider: %v", err)
	}
	
	events, err := calendarProvider.GetData()
	if err != nil {
		t.Errorf("Unexpected error getting calendar data: %v", err)
	}
	
	if len(events) != 3 {
		t.Errorf("Expected 3 events, got %d", len(events))
	}

	// Test creating an email provider
	emailProvider, err := CreateTypedProvider[[]models.Email]("mock")
	if err != nil {
		t.Errorf("Unexpected error creating email provider: %v", err)
	}
	
	emails, err := emailProvider.GetData()
	if err != nil {
		t.Errorf("Unexpected error getting email data: %v", err)
	}
	
	if len(emails) != 3 {
		t.Errorf("Expected 3 emails, got %d", len(emails))
	}
}

// TestCreateTypedProviderUnsupportedType validates error handling for unsupported types
func TestCreateTypedProviderUnsupportedType(t *testing.T) {
	// Test with an unsupported type
	_, err := CreateTypedProvider[string]("mock")
	if err == nil {
		t.Error("Expected error for unsupported type, got nil")
	}
	
	expectedError := "unsupported typed provider for type string"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}