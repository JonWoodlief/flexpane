package panes

import (
	"context"

	"flexplane/internal/models"
	"flexplane/internal/providers"
)

// EmailPane implements both Pane and TypedPane interfaces for email messages
// The generic TypedPane provides compile-time type safety for the EmailPaneData
type EmailPane struct {
	provider providers.DataProvider
}

func NewEmailPane(provider providers.DataProvider) *EmailPane {
	return &EmailPane{
		provider: provider,
	}
}

func (ep *EmailPane) ID() string {
	return "email"
}

func (ep *EmailPane) Title() string {
	return "Email Preview"
}

func (ep *EmailPane) Template() string {
	return "panes/email.html"
}

// GetData maintains backward compatibility by returning interface{}
func (ep *EmailPane) GetData(ctx context.Context) (interface{}, error) {
	return ep.GetTypedData(ctx)
}

// GetTypedData provides type-safe access to email data
// This eliminates the need for type assertions in calling code
func (ep *EmailPane) GetTypedData(ctx context.Context) (models.EmailPaneData, error) {
	emails, err := ep.provider.GetEmails()
	if err != nil {
		return models.EmailPaneData{
			Emails: []models.Email{},
			Count:  0,
		}, err
	}

	return models.EmailPaneData{
		Emails: emails,
		Count:  len(emails),
	}, nil
}