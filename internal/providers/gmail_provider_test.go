package providers

import (
	"testing"
)

func TestNewGmailProvider_RequiresConfiguration(t *testing.T) {
	// Test that Gmail provider requires configuration
	config := GmailConfig{}
	
	_, err := NewGmailProvider(config)
	if err == nil {
		t.Error("Expected error when creating Gmail provider without configuration")
	}
	
	expectedMsg := "gmail provider requires either credentials_path or access_token/refresh_token"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestGmailProvider_ImplementsInterface(t *testing.T) {
	// Test that GmailProvider implements DataProvider interface
	// This is a compile-time check
	var _ DataProvider = (*GmailProvider)(nil)
}

func TestCreateGmailProvider_FromFactory(t *testing.T) {
	factory := &ProviderFactory{
		config: DataProviderConfig{
			Providers: map[string]ProviderConfig{
				"test-gmail": {
					Type: ProviderTypeGmail,
					Config: map[string]interface{}{
						"client_id":     "test-client-id",
						"client_secret": "test-client-secret",
						"access_token":  "test-access-token",
						"refresh_token": "test-refresh-token",
					},
				},
			},
			Default: "test-gmail",
		},
	}
	
	// The provider should be created successfully with fake credentials
	provider, err := factory.CreateProvider("test-gmail")
	if err != nil {
		t.Errorf("Unexpected error creating Gmail provider: %v", err)
		return
	}
	
	if provider == nil {
		t.Error("Expected provider to be created")
		return
	}
	
	// Verify it implements DataProvider interface
	_, ok := provider.(DataProvider)
	if !ok {
		t.Error("Provider does not implement DataProvider interface")
	}
	
	// The actual API calls should fail with fake credentials
	// Test GetEmails (this will fail due to invalid credentials)
	_, err = provider.GetEmails()
	if err == nil {
		t.Error("Expected GetEmails to fail with fake credentials")
	} else {
		t.Logf("GetEmails failed as expected: %v", err)
	}
	
	// Test GetCalendarEvents (this should return empty slice for Gmail provider)
	events, err := provider.GetCalendarEvents()
	if err != nil {
		t.Errorf("GetCalendarEvents should not fail: %v", err)
	}
	if len(events) != 0 {
		t.Error("Gmail provider should return empty calendar events")
	}
}

func TestProviderTypeConstants(t *testing.T) {
	// Test that our constants are set correctly
	if ProviderTypeGmail != "gmail" {
		t.Errorf("Expected ProviderTypeGmail to be 'gmail', got '%s'", ProviderTypeGmail)
	}
	
	if ProviderTypeMock != "mock" {
		t.Errorf("Expected ProviderTypeMock to be 'mock', got '%s'", ProviderTypeMock)
	}
}