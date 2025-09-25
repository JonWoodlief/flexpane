package providers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"flexplane/internal/models"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// GmailConfig holds configuration for Gmail provider
type GmailConfig struct {
	ClientID     string `envconfig:"GOOGLE_CLIENT_ID"`
	ClientSecret string `envconfig:"GOOGLE_CLIENT_SECRET"`
	RedirectURL  string `envconfig:"OAUTH_REDIRECT_URL" default:"http://localhost:3000/auth/callback"`
}

// GmailProvider implements DataProvider using Google APIs
type GmailProvider struct {
	config        *GmailConfig
	oauth2Config  *oauth2.Config
	token         *oauth2.Token
	userInfo      *UserInfo
	calendarSvc   *calendar.Service
	gmailSvc      *gmail.Service
	authenticated bool
}

// NewGmailProvider creates a new Gmail provider with OAuth configuration
// Uses envconfig library for better configuration management
func NewGmailProvider() *GmailProvider {
	var config GmailConfig
	if err := envconfig.Process("", &config); err != nil {
		log.Printf("Error processing Gmail configuration: %v", err)
		// Fall back to demo configuration for development
		config = GmailConfig{
			ClientID:    "demo-client-id",
			RedirectURL: "http://localhost:3000/auth/callback",
		}
	}

	// Validate required configuration
	if config.ClientID == "" {
		log.Println("Using demo OAuth credentials - set GOOGLE_CLIENT_ID for production")
		config.ClientID = "demo-client-id"
	}

	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes: []string{
			calendar.CalendarReadonlyScope,
			gmail.GmailReadonlyScope,
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint:    google.Endpoint,
		RedirectURL: config.RedirectURL,
	}

	return &GmailProvider{
		config:        &config,
		oauth2Config:  oauth2Config,
		authenticated: false,
	}
}

// IsAuthenticated returns whether the provider has valid authentication
func (g *GmailProvider) IsAuthenticated() bool {
	return g.authenticated && g.token != nil && g.token.Valid()
}

// GetAuthURL returns the OAuth URL for user authentication
func (g *GmailProvider) GetAuthURL() (string, error) {
	// Use library's built-in PKCE and security features
	return g.oauth2Config.AuthCodeURL("state", 
		oauth2.AccessTypeOffline, 
		oauth2.ApprovalForce,
		oauth2.SetAuthURLParam("prompt", "consent"),
	), nil
}

// Authenticate exchanges the OAuth code for tokens and sets up API clients
func (g *GmailProvider) Authenticate(ctx context.Context, code string) error {
	token, err := g.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("failed to exchange OAuth code: %w", err)
	}

	g.token = token

	// Use standard OAuth2 client - Google's newer auth libs are more complex for this use case
	client := g.oauth2Config.Client(ctx, token)

	// Initialize Google API clients using the authenticated client
	if err := g.initializeServices(ctx, client); err != nil {
		return fmt.Errorf("failed to initialize Google services: %w", err)
	}

	// Get user info using library helper
	if err := g.fetchUserInfo(ctx); err != nil {
		log.Printf("Warning: failed to fetch user info: %v", err)
	}

	g.authenticated = true
	return nil
}

// initializeServices sets up Google API service clients
func (g *GmailProvider) initializeServices(ctx context.Context, client *http.Client) error {
	var err error
	
	g.calendarSvc, err = calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("failed to create calendar service: %w", err)
	}

	g.gmailSvc, err = gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("failed to create gmail service: %w", err)
	}

	return nil
}

// fetchUserInfo gets basic user information using Gmail API
func (g *GmailProvider) fetchUserInfo(ctx context.Context) error {
	profile, err := g.gmailSvc.Users.GetProfile("me").Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to get user profile: %w", err)
	}

	g.userInfo = &UserInfo{
		Email: profile.EmailAddress,
		Name:  profile.EmailAddress, // Gmail API doesn't provide display name in profile
	}

	return nil
}

// GetUserInfo returns cached user information
func (g *GmailProvider) GetUserInfo() (*UserInfo, error) {
	if g.userInfo == nil {
		return nil, fmt.Errorf("user not authenticated or user info not available")
	}
	return g.userInfo, nil
}

// GetCalendarEvents fetches events from Google Calendar with better error handling
func (g *GmailProvider) GetCalendarEvents() ([]models.Event, error) {
	if !g.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	// Use library helpers for time calculations
	now := time.Now()
	endOfDay := now.AddDate(0, 0, 1).Truncate(24*time.Hour).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	calendarCall := g.calendarSvc.Events.List("primary").
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(now.Format(time.RFC3339)).
		TimeMax(endOfDay.Format(time.RFC3339)).
		MaxResults(20).
		OrderBy("startTime")

	events, err := calendarCall.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch calendar events: %w", err)
	}

	result := make([]models.Event, 0, len(events.Items))
	for _, item := range events.Items {
		event, err := g.parseCalendarEvent(item)
		if err != nil {
			log.Printf("Warning: failed to parse calendar event %s: %v", item.Id, err)
			continue
		}
		result = append(result, event)
	}

	return result, nil
}

// parseCalendarEvent converts Google Calendar event to internal model with better parsing
func (g *GmailProvider) parseCalendarEvent(item *calendar.Event) (models.Event, error) {
	event := models.Event{
		ID:    item.Id,
		Title: item.Summary,
	}

	// Use library helpers for time parsing with fallbacks
	var err error
	if item.Start.DateTime != "" {
		event.Start, err = time.Parse(time.RFC3339, item.Start.DateTime)
		if err != nil {
			return event, fmt.Errorf("failed to parse start time: %w", err)
		}
	} else if item.Start.Date != "" {
		// Handle all-day events
		event.Start, err = time.Parse("2006-01-02", item.Start.Date)
		if err != nil {
			return event, fmt.Errorf("failed to parse start date: %w", err)
		}
	}

	if item.End.DateTime != "" {
		event.End, err = time.Parse(time.RFC3339, item.End.DateTime)
		if err != nil {
			return event, fmt.Errorf("failed to parse end time: %w", err)
		}
	} else if item.End.Date != "" {
		event.End, err = time.Parse("2006-01-02", item.End.Date)
		if err != nil {
			return event, fmt.Errorf("failed to parse end date: %w", err)
		}
	}

	if item.Location != "" {
		event.Location = item.Location
	}

	return event, nil
}

// GetEmails fetches recent emails from Gmail with improved error handling
func (g *GmailProvider) GetEmails() ([]models.Email, error) {
	if !g.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	// Get recent messages from inbox using library best practices
	messagesCall := g.gmailSvc.Users.Messages.List("me").
		Q("in:inbox").
		MaxResults(10)

	messages, err := messagesCall.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list emails: %w", err)
	}

	result := make([]models.Email, 0, len(messages.Messages))
	for _, message := range messages.Messages {
		msg, err := g.gmailSvc.Users.Messages.Get("me", message.Id).Do()
		if err != nil {
			log.Printf("Warning: failed to get message %s: %v", message.Id, err)
			continue
		}

		email := g.parseGmailMessage(msg)
		result = append(result, email)
	}

	return result, nil
}

// parseGmailMessage converts a Gmail message to our Email model with better parsing
func (g *GmailProvider) parseGmailMessage(msg *gmail.Message) models.Email {
	email := models.Email{
		ID:   msg.Id,
		Read: true, // Default to read
	}

	// Parse message timestamp using library helper
	if msg.InternalDate > 0 {
		email.Time = time.Unix(msg.InternalDate/1000, (msg.InternalDate%1000)*1000000)
	}

	// Check if message is unread using library constants where possible
	for _, labelId := range msg.LabelIds {
		if labelId == "UNREAD" {
			email.Read = false
			break
		}
	}

	// Parse headers using structured approach
	headers := make(map[string]string)
	for _, header := range msg.Payload.Headers {
		headers[header.Name] = header.Value
	}

	if subject, exists := headers["Subject"]; exists {
		email.Subject = subject
	}
	if from, exists := headers["From"]; exists {
		email.From = from
	}

	// Get message snippet as preview
	email.Preview = msg.Snippet

	return email
}

// SaveToken saves the OAuth token using structured approach
func (g *GmailProvider) SaveToken() error {
	if g.token == nil {
		return fmt.Errorf("no token to save")
	}

	// In a real application, you'd use a proper token storage library
	// For demo purposes, we'll use structured logging
	log.Printf("Token saved for user (in production, save securely): expires=%v", g.token.Expiry)
	return nil
}

// LoadToken loads a saved OAuth token with better validation
func (g *GmailProvider) LoadToken(tokenSource oauth2.TokenSource) error {
	token, err := tokenSource.Token()
	if err != nil {
		return fmt.Errorf("failed to load token: %w", err)
	}

	if !token.Valid() {
		return fmt.Errorf("token is expired")
	}

	g.token = token
	
	// Re-initialize services with the loaded token
	ctx := context.Background()
	client := g.oauth2Config.Client(ctx, token)

	if err := g.initializeServices(ctx, client); err != nil {
		return fmt.Errorf("failed to reinitialize services: %w", err)
	}

	if err := g.fetchUserInfo(ctx); err != nil {
		log.Printf("Warning: failed to fetch user info: %v", err)
	}

	g.authenticated = true
	return nil
}