package panes

import (
	"context"

	"flexpane/internal/models"
	"flexpane/internal/providers"
)

// CalendarPane implements both Pane and TypedPane interfaces for calendar events
// The generic TypedPane provides compile-time type safety for the CalendarPaneData
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

// GetData maintains backward compatibility by returning interface{}
func (cp *CalendarPane) GetData(ctx context.Context) (interface{}, error) {
	return cp.GetTypedData(ctx)
}

// GetTypedData provides type-safe access to calendar data
// This eliminates the need for type assertions in calling code
func (cp *CalendarPane) GetTypedData(ctx context.Context) (models.CalendarPaneData, error) {
	events, err := cp.provider.GetCalendarEvents()
	if err != nil {
		return models.CalendarPaneData{
			Events: []models.Event{},
			Count:  0,
		}, err
	}

	return models.CalendarPaneData{
		Events: events,
		Count:  len(events),
	}, nil
}