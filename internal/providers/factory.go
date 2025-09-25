package providers

import (
	"log"

	"flexplane/internal/services"
	"github.com/kelseyhightower/envconfig"
)

// ProviderConfig holds configuration for provider selection
type ProviderConfig struct {
	GoogleClientID string `envconfig:"GOOGLE_CLIENT_ID"`
}

// ProviderFactory creates the appropriate provider based on configuration
type ProviderFactory struct {
	config ProviderConfig
}

// NewProviderFactory creates a new provider factory with configuration loaded from environment
func NewProviderFactory() *ProviderFactory {
	var config ProviderConfig
	if err := envconfig.Process("", &config); err != nil {
		log.Printf("Error loading provider configuration: %v", err)
	}

	return &ProviderFactory{
		config: config,
	}
}

// CreateProvider returns either a Gmail provider (if configured) or Mock provider (default)
func (f *ProviderFactory) CreateProvider() DataProvider {
	if f.config.GoogleClientID != "" {
		log.Println("Google OAuth configured - using Gmail provider")
		return NewGmailProvider()
	}

	// For distributed apps without setup, we'll use mock data
	// This allows the app to work out of the box
	log.Println("No Google OAuth configured - using mock provider for demo")
	log.Println("Set GOOGLE_CLIENT_ID environment variable to enable Gmail integration")
	return services.NewMockProvider()
}

// CreateAuthenticatedProvider returns a Gmail provider that supports authentication
// Returns nil if Gmail is not configured
func (f *ProviderFactory) CreateAuthenticatedProvider() AuthenticatedProvider {
	if f.config.GoogleClientID != "" {
		log.Println("Google OAuth configured - using Gmail provider")
		return NewGmailProvider()
	}

	return nil
}