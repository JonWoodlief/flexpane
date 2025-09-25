package services

import (
	"context"
	"fmt"

	"flexplane/internal/models"
)

// TypedPaneRegistry provides type-safe pane registration and retrieval
// This eliminates the need for type assertions when working with specific pane types
type TypedPaneRegistry[T any] struct {
	registry *PaneRegistry
	panes    map[string]models.TypedPane[T]
}

// NewTypedPaneRegistry creates a new type-safe pane registry
func NewTypedPaneRegistry[T any](registry *PaneRegistry) *TypedPaneRegistry[T] {
	return &TypedPaneRegistry[T]{
		registry: registry,
		panes:    make(map[string]models.TypedPane[T]),
	}
}

// RegisterTypedPane registers a typed pane with compile-time type safety
func (tr *TypedPaneRegistry[T]) RegisterTypedPane(pane models.TypedPane[T]) {
	tr.registry.RegisterPane(pane)
	tr.panes[pane.ID()] = pane
}

// GetTypedPane returns a typed pane with compile-time guarantees
func (tr *TypedPaneRegistry[T]) GetTypedPane(paneID string) (models.TypedPane[T], bool) {
	pane, exists := tr.panes[paneID]
	return pane, exists
}

// GetTypedData gets typed data from a specific pane
func (tr *TypedPaneRegistry[T]) GetTypedData(ctx context.Context, paneID string) (T, error) {
	var zero T
	pane, exists := tr.GetTypedPane(paneID)
	if !exists {
		return zero, fmt.Errorf("pane not found: %s", paneID)
	}
	return pane.GetTypedData(ctx)
}

// GenericPaneManager demonstrates how generics could be used for cross-cutting concerns
// This provides type-safe operations across all pane types
type GenericPaneManager struct {
	calendarRegistry *TypedPaneRegistry[models.CalendarPaneData]
	emailRegistry    *TypedPaneRegistry[models.EmailPaneData]
	todoRegistry     *TypedPaneRegistry[models.TodoPaneData]
}

// NewGenericPaneManager creates a new generic pane manager with type-safe registries
func NewGenericPaneManager(baseRegistry *PaneRegistry) *GenericPaneManager {
	return &GenericPaneManager{
		calendarRegistry: NewTypedPaneRegistry[models.CalendarPaneData](baseRegistry),
		emailRegistry:    NewTypedPaneRegistry[models.EmailPaneData](baseRegistry),
		todoRegistry:     NewTypedPaneRegistry[models.TodoPaneData](baseRegistry),
	}
}

// RegisterCalendarPane provides type-safe calendar pane registration
func (gpm *GenericPaneManager) RegisterCalendarPane(pane models.TypedPane[models.CalendarPaneData]) {
	gpm.calendarRegistry.RegisterTypedPane(pane)
}

// RegisterEmailPane provides type-safe email pane registration
func (gpm *GenericPaneManager) RegisterEmailPane(pane models.TypedPane[models.EmailPaneData]) {
	gpm.emailRegistry.RegisterTypedPane(pane)
}

// RegisterTodoPane provides type-safe todo pane registration
func (gpm *GenericPaneManager) RegisterTodoPane(pane models.TypedPane[models.TodoPaneData]) {
	gpm.todoRegistry.RegisterTypedPane(pane)
}

// GetCalendarData provides compile-time type safety for calendar data
func (gpm *GenericPaneManager) GetCalendarData(ctx context.Context, paneID string) (models.CalendarPaneData, error) {
	return gpm.calendarRegistry.GetTypedData(ctx, paneID)
}

// GetEmailData provides compile-time type safety for email data
func (gpm *GenericPaneManager) GetEmailData(ctx context.Context, paneID string) (models.EmailPaneData, error) {
	return gpm.emailRegistry.GetTypedData(ctx, paneID)
}

// GetTodoData provides compile-time type safety for todo data
func (gpm *GenericPaneManager) GetTodoData(ctx context.Context, paneID string) (models.TodoPaneData, error) {
	return gpm.todoRegistry.GetTypedData(ctx, paneID)
}