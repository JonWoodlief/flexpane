package panes

import (
	"context"

	"flexpane/internal/models"
	"flexpane/internal/providers"
)

// EmailPane implements the Pane interface for email messages
type EmailPane struct {
	provider providers.DataProvider
	typedProvider providers.EmailProvider
}

func NewEmailPane(provider providers.DataProvider) *EmailPane {
	return &EmailPane{
		provider: provider,
		typedProvider: providers.NewEmailProvider(provider),
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

// GetTypedData implements the TypedPane interface for type-safe data access
func (ep *EmailPane) GetTypedData(ctx context.Context) (models.EmailPaneData, error) {
	emails, err := ep.typedProvider.GetData()
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