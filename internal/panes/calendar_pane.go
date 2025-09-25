package panes

import (
	"context"

	"flexplane/internal/models"
	"flexplane/internal/providers"
)

// CalendarPane implements the Pane interface for calendar events
type CalendarPane struct {
	provider providers.DataProvider
}

func NewCalendarPane(provider providers.DataProvider) *CalendarPane {
	return &CalendarPane{
		provider: provider,
	}
}

func (cp *CalendarPane) ID() string {
	return "calendar"
}

func (cp *CalendarPane) Title() string {
	return "Calendar"
}


func (cp *CalendarPane) Template() string {
	return "panes/calendar.html"
}

func (cp *CalendarPane) GetData(ctx context.Context) (interface{}, error) {
	events, err := cp.provider.GetCalendarEvents()
	if err != nil {
		return map[string]interface{}{
			"Events": []models.Event{},
			"Count":  0,
		}, err
	}

	return map[string]interface{}{
		"Events": events,
		"Count":  len(events),
	}, nil
}