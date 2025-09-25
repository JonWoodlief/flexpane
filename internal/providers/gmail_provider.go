package providers

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"flexplane/internal/models"
	
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// GmailProvider implements the DataProvider interface using Gmail API
type GmailProvider struct {
	service *gmail.Service
	ctx     context.Context
}

// GmailConfig holds configuration for Gmail provider
type GmailConfig struct {
	CredentialsPath string `json:"credentials_path,omitempty"`
	TokenPath       string `json:"token_path,omitempty"`
	ClientID        string `json:"client_id,omitempty"`
	ClientSecret    string `json:"client_secret,omitempty"`
	AccessToken     string `json:"access_token,omitempty"`
	RefreshToken    string `json:"refresh_token,omitempty"`
}

// NewGmailProvider creates a new Gmail provider with the given configuration
func NewGmailProvider(config GmailConfig) (*GmailProvider, error) {
	ctx := context.Background()
	
	var client *http.Client
	var err error
	
	if config.AccessToken != "" && config.RefreshToken != "" {
		// Use provided tokens directly
		token := &oauth2.Token{
			AccessToken:  config.AccessToken,
			RefreshToken: config.RefreshToken,
			TokenType:    "Bearer",
		}
		
		oauthConfig := &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			Scopes:       []string{gmail.GmailReadonlyScope},
			Endpoint:     google.Endpoint,
		}
		
		client = oauthConfig.Client(ctx, token)
	} else if config.CredentialsPath != "" {
		// TODO: SECURITY - Implement OAuth2 flow for production use
		// For now, this is a placeholder for credential-based auth
		return nil, fmt.Errorf("credentials file authentication not implemented yet - use access_token/refresh_token for now")
	} else {
		return nil, fmt.Errorf("gmail provider requires either credentials_path or access_token/refresh_token")
	}
	
	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gmail service: %w", err)
	}
	
	return &GmailProvider{
		service: service,
		ctx:     ctx,
	}, nil
}

// GetCalendarEvents implements DataProvider interface
// Note: Gmail API doesn't provide calendar events, so we return empty slice
func (g *GmailProvider) GetCalendarEvents() ([]models.Event, error) {
	// Gmail API doesn't provide calendar data - would need Google Calendar API
	// For now, return empty slice to satisfy interface
	// TODO: Consider adding Google Calendar API integration separately
	return []models.Event{}, nil
}

// GetEmails implements DataProvider interface by fetching emails from Gmail
func (g *GmailProvider) GetEmails() ([]models.Email, error) {
	// Get list of messages from inbox
	req := g.service.Users.Messages.List("me").Q("in:inbox").MaxResults(10)
	
	resp, err := req.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve messages: %w", err)
	}
	
	var emails []models.Email
	
	for _, msg := range resp.Messages {
		// Get full message details
		fullMsg, err := g.service.Users.Messages.Get("me", msg.Id).Do()
		if err != nil {
			log.Printf("Failed to get message %s: %v", msg.Id, err)
			continue
		}
		
		email := g.convertToEmail(fullMsg)
		emails = append(emails, email)
	}
	
	return emails, nil
}

// convertToEmail converts Gmail API message to our Email model
func (g *GmailProvider) convertToEmail(msg *gmail.Message) models.Email {
	email := models.Email{
		ID:   msg.Id,
		Time: time.Unix(msg.InternalDate/1000, 0), // Gmail uses milliseconds
		Read: !g.isUnread(msg),
	}
	
	// Extract headers
	for _, header := range msg.Payload.Headers {
		switch strings.ToLower(header.Name) {
		case "subject":
			email.Subject = header.Value
		case "from":
			email.From = header.Value
		}
	}
	
	// Extract preview from body
	email.Preview = g.extractPreview(msg.Payload)
	
	return email
}

// isUnread checks if message has UNREAD label
func (g *GmailProvider) isUnread(msg *gmail.Message) bool {
	for _, labelId := range msg.LabelIds {
		if labelId == "UNREAD" {
			return true
		}
	}
	return false
}

// extractPreview extracts a text preview from message payload
func (g *GmailProvider) extractPreview(payload *gmail.MessagePart) string {
	// Look for text/plain part first
	if payload.MimeType == "text/plain" && payload.Body != nil && payload.Body.Data != "" {
		decoded, err := base64URLDecode(payload.Body.Data)
		if err == nil {
			preview := string(decoded)
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}
			return preview
		}
	}
	
	// Recursively search parts
	for _, part := range payload.Parts {
		if preview := g.extractPreview(part); preview != "" {
			return preview
		}
	}
	
	return "No preview available"
}

// base64URLDecode decodes Gmail's base64URL encoding
func base64URLDecode(data string) ([]byte, error) {
	// Gmail uses base64url encoding without padding
	data = strings.ReplaceAll(data, "-", "+")
	data = strings.ReplaceAll(data, "_", "/")
	
	// Add padding if needed
	switch len(data) % 4 {
	case 2:
		data += "=="
	case 3:
		data += "="
	}
	
	return base64.StdEncoding.DecodeString(data)
}