package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"flexplane/internal/models"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// GmailProvider implements DataProvider using Google APIs
type GmailProvider struct {
	config       *oauth2.Config
	token        *oauth2.Token
	userInfo     *UserInfo
	calendarSvc  *calendar.Service
	gmailSvc     *gmail.Service
	authenticated bool
}

// NewGmailProvider creates a new Gmail provider with OAuth configuration
// For demo purposes, we use a development OAuth client
// In production, users should provide their own OAuth credentials
func NewGmailProvider() *GmailProvider {
	// Try to get OAuth credentials from environment
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	
	// Fall back to demo credentials (these would be for a demo app)
	// Note: In a real distributed app, you'd need proper OAuth setup
	if clientID == "" {
		log.Println("Using demo OAuth credentials - set GOOGLE_CLIENT_ID for production")
		clientID = "your-demo-client-id"  // This needs to be replaced with actual demo credentials
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret, // Can be empty for some OAuth flows
		Scopes: []string{
			calendar.CalendarReadonlyScope,
			gmail.GmailReadonlyScope,
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint:    google.Endpoint,
		RedirectURL: "http://localhost:3000/auth/callback",
	}

	return &GmailProvider{
		config:        config,
		authenticated: false,
	}
}

// IsAuthenticated returns whether the provider has valid authentication
func (g *GmailProvider) IsAuthenticated() bool {
	return g.authenticated && g.token != nil && g.token.Valid()
}

// GetAuthURL returns the OAuth URL for user authentication
func (g *GmailProvider) GetAuthURL() (string, error) {
	// Use PKCE for additional security without client secret
	return g.config.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce), nil
}

// Authenticate exchanges the OAuth code for tokens and sets up API clients
func (g *GmailProvider) Authenticate(ctx context.Context, code string) error {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("failed to exchange OAuth code: %v", err)
	}

	g.token = token
	client := g.config.Client(ctx, token)

	// Initialize Google API clients
	g.calendarSvc, err = calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("failed to create calendar service: %v", err)
	}

	g.gmailSvc, err = gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("failed to create gmail service: %v", err)
	}

	// Get user info
	if err := g.fetchUserInfo(ctx); err != nil {
		log.Printf("Warning: failed to fetch user info: %v", err)
	}

	g.authenticated = true
	return nil
}

// fetchUserInfo gets basic user information
func (g *GmailProvider) fetchUserInfo(ctx context.Context) error {
	profile, err := g.gmailSvc.Users.GetProfile("me").Context(ctx).Do()
	if err != nil {
		return err
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

// GetCalendarEvents fetches events from Google Calendar
func (g *GmailProvider) GetCalendarEvents() ([]models.Event, error) {
	if !g.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	now := time.Now()
	endOfDay := time.Date(now.Year(), now.Month(), now.Day()+1, 23, 59, 59, 0, now.Location())

	calendarCall := g.calendarSvc.Events.List("primary").
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(now.Format(time.RFC3339)).
		TimeMax(endOfDay.Format(time.RFC3339)).
		MaxResults(20).
		OrderBy("startTime")

	events, err := calendarCall.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch calendar events: %v", err)
	}

	var result []models.Event
	for _, item := range events.Items {
		start, _ := time.Parse(time.RFC3339, item.Start.DateTime)
		end, _ := time.Parse(time.RFC3339, item.End.DateTime)

		// Handle all-day events
		if item.Start.DateTime == "" {
			start, _ = time.Parse("2006-01-02", item.Start.Date)
			end, _ = time.Parse("2006-01-02", item.End.Date)
		}

		location := ""
		if item.Location != "" {
			location = item.Location
		}

		result = append(result, models.Event{
			ID:       item.Id,
			Title:    item.Summary,
			Start:    start,
			End:      end,
			Location: location,
		})
	}

	return result, nil
}

// GetEmails fetches recent emails from Gmail
func (g *GmailProvider) GetEmails() ([]models.Email, error) {
	if !g.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	// Get recent messages from inbox
	messagesCall := g.gmailSvc.Users.Messages.List("me").
		Q("in:inbox").
		MaxResults(10)

	messages, err := messagesCall.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list emails: %v", err)
	}

	var result []models.Email
	for _, message := range messages.Messages {
		msg, err := g.gmailSvc.Users.Messages.Get("me", message.Id).Do()
		if err != nil {
			log.Printf("Failed to get message %s: %v", message.Id, err)
			continue
		}

		email := g.parseGmailMessage(msg)
		result = append(result, email)
	}

	return result, nil
}

// parseGmailMessage converts a Gmail message to our Email model
func (g *GmailProvider) parseGmailMessage(msg *gmail.Message) models.Email {
	email := models.Email{
		ID:   msg.Id,
		Read: true, // Default to read
	}

	// Parse message timestamp
	timestamp, err := parseInternalDate(msg.InternalDate)
	if err == nil {
		email.Time = timestamp
	}

	// Check if message is unread
	for _, labelId := range msg.LabelIds {
		if labelId == "UNREAD" {
			email.Read = false
			break
		}
	}

	// Parse headers for subject and from
	for _, header := range msg.Payload.Headers {
		switch header.Name {
		case "Subject":
			email.Subject = header.Value
		case "From":
			email.From = header.Value
		}
	}

	// Get message snippet as preview
	email.Preview = msg.Snippet

	return email
}

// parseInternalDate converts Gmail's internal date to time.Time
func parseInternalDate(internalDate int64) (time.Time, error) {
	return time.Unix(internalDate/1000, (internalDate%1000)*1000000), nil
}

// SaveToken saves the OAuth token (for demo purposes, we'll just log it)
// In a real application, you'd want to save this securely
func (g *GmailProvider) SaveToken() error {
	if g.token == nil {
		return fmt.Errorf("no token to save")
	}

	tokenJSON, err := json.Marshal(g.token)
	if err != nil {
		return err
	}

	log.Printf("Token saved (in production, save this securely): %s", string(tokenJSON))
	return nil
}

// LoadToken loads a saved OAuth token (for demo purposes)
// In a real application, you'd load this from secure storage
func (g *GmailProvider) LoadToken(tokenJSON []byte) error {
	var token oauth2.Token
	if err := json.Unmarshal(tokenJSON, &token); err != nil {
		return err
	}

	if !token.Valid() {
		return fmt.Errorf("token is expired")
	}

	g.token = &token
	
	// Re-initialize services with the loaded token
	ctx := context.Background()
	client := g.config.Client(ctx, &token)

	var err error
	g.calendarSvc, err = calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("failed to create calendar service: %v", err)
	}

	g.gmailSvc, err = gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("failed to create gmail service: %v", err)
	}

	if err := g.fetchUserInfo(ctx); err != nil {
		log.Printf("Warning: failed to fetch user info: %v", err)
	}

	g.authenticated = true
	return nil
}