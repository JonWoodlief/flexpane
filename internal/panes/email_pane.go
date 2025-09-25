package panes

import (
	"context"

	"flexpane/internal/models"
	"flexpane/internal/providers"
)

// EmailPane implements the Pane interface for email messages
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

func (ep *EmailPane) GetData(ctx context.Context) (interface{}, error) {
	emails, err := ep.provider.GetEmails()
	if err != nil {
		return map[string]interface{}{
			"Emails": []models.Email{},
			"Count":  0,
		}, err
	}

	return map[string]interface{}{
		"Emails": emails,
		"Count":  len(emails),
	}, nil
}