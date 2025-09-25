package providers

import (
	"testing"
)

func TestMockProvider(t *testing.T) {
	provider := NewMockProvider()

	// Test calendar events
	events, err := provider.GetCalendarEvents()
	if err != nil {
		t.Errorf("GetCalendarEvents failed: %v", err)
	}
	if len(events) == 0 {
		t.Error("Expected mock calendar events, got none")
	}

	// Test emails
	emails, err := provider.GetEmails()
	if err != nil {
		t.Errorf("GetEmails failed: %v", err)
	}
	if len(emails) == 0 {
		t.Error("Expected mock emails, got none")
	}
}

func TestProviderFactory(t *testing.T) {
	factory := NewProviderFactory()

	// Test that all providers are registered
	available := factory.GetAvailableProviders()
	expectedTypes := []string{"mock"}
	
	for _, expected := range expectedTypes {
		found := false
		for _, available := range available {
			if available == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected provider type '%s' to be available", expected)
		}
	}

	// Test creating mock provider
	mockProvider, err := factory.CreateProvider(ProviderConfig{Type: "mock"})
	if err != nil {
		t.Errorf("Failed to create mock provider: %v", err)
	}
	if mockProvider == nil {
		t.Error("Expected non-nil mock provider")
	}

	// Test creating unknown provider
	_, err = factory.CreateProvider(ProviderConfig{Type: "unknown"})
	if err == nil {
		t.Error("Expected error when creating unknown provider type")
	}

	// Test custom provider registration
	factory.RegisterProvider("custom", func(args map[string]interface{}) (DataProvider, error) {
		return NewMockProvider(), nil
	})

	customProvider, err := factory.CreateProvider(ProviderConfig{Type: "custom"})
	if err != nil {
		t.Errorf("Failed to create custom provider: %v", err)
	}
	if customProvider == nil {
		t.Error("Expected non-nil custom provider")
	}
}