package models

import (
	"context"
	"net/http"
	"time"
)

// PaneSize represents the relative size of a pane in the grid
type PaneSize int

const (
	PaneSmall  PaneSize = 1
	PaneMedium PaneSize = 2
	PaneLarge  PaneSize = 3
)

// PaneGridArea represents where a pane should be placed in the grid
type PaneGridArea struct {
	Row    string `json:"row"`    // CSS grid-row value (e.g., "1", "2", "1 / -1")
	Column string `json:"column"` // CSS grid-column value (e.g., "span 3", "1 / 4")
}

// Pane interface defines the contract for all panes
type Pane interface {
	ID() string
	Title() string
	GetData(ctx context.Context) (interface{}, error)
	Template() string
}

// APIHandler interface for panes that need API endpoints
type APIHandler interface {
	HandleAPI(w http.ResponseWriter, r *http.Request) error
}

// PaneData holds the rendered data for a pane
type PaneData struct {
	ID       string       `json:"id"`
	Title    string       `json:"title"`
	GridArea PaneGridArea `json:"grid_area"`
	Data     interface{}  `json:"data"`
	Template string       `json:"template"`
}

// Simple data models
type Todo struct {
	Done    bool   `json:"done"`
	Message string `json:"message"`
}

type Event struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Location string    `json:"location,omitempty"`
}

type Email struct {
	ID      string    `json:"id"`
	Subject string    `json:"subject"`
	From    string    `json:"from"`
	Preview string    `json:"preview"`
	Time    time.Time `json:"time"`
	Read    bool      `json:"read"`
}

// PageData contains all data for the main page
type PageData struct {
	Panes []PaneData `json:"panes"`
}