package panes

import (
	"context"

	"flexpane/internal/models"
	"flexpane/internal/providers"
)

// CalendarPane implements the Pane interface for calendar events
type CalendarPane struct {
	provider providers.DataProvider
	typedProvider providers.CalendarProvider
}

func NewCalendarPane(provider providers.DataProvider) *CalendarPane {
	return &CalendarPane{
		provider: provider,
		typedProvider: providers.NewCalendarProvider(provider),
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

// GetTypedData implements the TypedPane interface for type-safe data access
func (cp *CalendarPane) GetTypedData(ctx context.Context) (models.CalendarPaneData, error) {
	events, err := cp.typedProvider.GetData()
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