package services

import (
	"context"

	"flexpane/internal/models"
)

// PaneRegistry manages all available panes
type PaneRegistry struct {
	panes   map[string]models.Pane
	enabled []string
	layout  map[string]PaneLayoutConfig
}

// PaneLayoutConfig holds layout configuration for a pane
type PaneLayoutConfig struct {
	GridArea models.PaneGridArea `json:"grid_area"`
}

func NewPaneRegistry() *PaneRegistry {
	return &PaneRegistry{
		panes:   make(map[string]models.Pane),
		enabled: []string{},
		layout:  make(map[string]PaneLayoutConfig),
	}
}

// RegisterPane adds a pane to the registry
func (pr *PaneRegistry) RegisterPane(pane models.Pane) {
	pr.panes[pane.ID()] = pane
}

// SetEnabledPanes sets which panes should be displayed
func (pr *PaneRegistry) SetEnabledPanes(paneIDs []string) {
	pr.enabled = paneIDs
}

// SetLayoutConfig sets layout configuration for panes
func (pr *PaneRegistry) SetLayoutConfig(layout map[string]PaneLayoutConfig) {
	pr.layout = layout
}

// GetEnabledPanes returns all enabled panes with their data
func (pr *PaneRegistry) GetEnabledPanes(ctx context.Context) ([]models.PaneData, error) {
	var paneData []models.PaneData

	for _, paneID := range pr.enabled {
		pane, exists := pr.panes[paneID]
		if !exists {
			continue // Skip missing panes gracefully
		}

		data, err := pane.GetData(ctx)
		if err != nil {
			// TODO: Add logging, for now continue with nil data
			data = nil
		}

		// Get layout config for this pane (required)
		layoutConfig := pr.layout[paneID]

		paneData = append(paneData, models.PaneData{
			ID:       pane.ID(),
			Title:    pane.Title(),
			GridArea: layoutConfig.GridArea,
			Data:     data,
			Template: pane.Template(),
		})
	}

	// Panes are returned in configuration order

	return paneData, nil
}

// GetPane returns a specific pane by ID
func (pr *PaneRegistry) GetPane(paneID string) (models.Pane, bool) {
	pane, exists := pr.panes[paneID]
	return pane, exists
}