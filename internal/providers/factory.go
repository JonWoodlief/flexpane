package providers

import (
	"log"
	"os"

	"flexplane/internal/services"
)

// ProviderFactory creates the appropriate provider based on configuration
type ProviderFactory struct{}

// NewProviderFactory creates a new provider factory
func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{}
}

// CreateProvider returns either a Gmail provider (if configured) or Mock provider (default)
func (f *ProviderFactory) CreateProvider() DataProvider {
	// Check if Google OAuth is configured
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	
	if clientID != "" {
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
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	
	if clientID != "" {
		log.Println("Google OAuth configured - using Gmail provider")
		return NewGmailProvider()
	}

	return nil
}