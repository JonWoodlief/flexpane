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

func TestNullProvider(t *testing.T) {
	provider := NewNullProvider()

	// Test calendar events (should be empty)
	events, err := provider.GetCalendarEvents()
	if err != nil {
		t.Errorf("GetCalendarEvents failed: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("Expected 0 events from null provider, got %d", len(events))
	}

	// Test emails (should be empty)
	emails, err := provider.GetEmails()
	if err != nil {
		t.Errorf("GetEmails failed: %v", err)
	}
	if len(emails) != 0 {
		t.Errorf("Expected 0 emails from null provider, got %d", len(emails))
	}
}

func TestProviderFactory(t *testing.T) {
	factory := NewProviderFactory()

	// Test that production providers are registered (no mock)
	available := factory.GetAvailableProviders()
	expectedTypes := []string{"null"}
	
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

	// Test that mock is NOT available in production factory
	for _, available := range available {
		if available == "mock" {
			t.Error("Mock provider should not be available in production factory")
		}
	}

	// Test creating null provider
	nullProvider, err := factory.CreateProvider(ProviderConfig{Type: "null"})
	if err != nil {
		t.Errorf("Failed to create null provider: %v", err)
	}
	if nullProvider == nil {
		t.Error("Expected non-nil null provider")
	}

	// Test creating unknown provider
	_, err = factory.CreateProvider(ProviderConfig{Type: "unknown"})
	if err == nil {
		t.Error("Expected error when creating unknown provider type")
	}

	// Test creating mock provider should fail in production factory
	_, err = factory.CreateProvider(ProviderConfig{Type: "mock"})
	if err == nil {
		t.Error("Expected error when creating mock provider in production factory")
	}

	// Test custom provider registration
	factory.RegisterProvider("custom", func(args map[string]interface{}) (DataProvider, error) {
		return NewNullProvider(), nil
	})

	customProvider, err := factory.CreateProvider(ProviderConfig{Type: "custom"})
	if err != nil {
		t.Errorf("Failed to create custom provider: %v", err)
	}
	if customProvider == nil {
		t.Error("Expected non-nil custom provider")
	}
}

func TestProviderFactoryWithMocks(t *testing.T) {
	factory := NewProviderFactoryWithMocks()

	// Test that both production and mock providers are available
	available := factory.GetAvailableProviders()
	expectedTypes := []string{"null", "mock"}
	
	for _, expected := range expectedTypes {
		found := false
		for _, available := range available {
			if available == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected provider type '%s' to be available in development factory", expected)
		}
	}

	// Test creating mock provider in development factory
	mockProvider, err := factory.CreateProvider(ProviderConfig{Type: "mock"})
	if err != nil {
		t.Errorf("Failed to create mock provider in development factory: %v", err)
	}
	if mockProvider == nil {
		t.Error("Expected non-nil mock provider")
	}
}