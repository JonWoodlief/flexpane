package main

import (
	"context"
	"testing"

	"flexpane/internal/models"
	"flexpane/internal/panes"
	"flexpane/internal/providers"
	"flexpane/internal/services"
)

// TestTypedPaneRegistry demonstrates type-safe pane registration and retrieval
func TestTypedPaneRegistry(t *testing.T) {
	baseRegistry := services.NewPaneRegistry()
	calendarRegistry := services.NewTypedPaneRegistry[models.CalendarPaneData](baseRegistry)
	
	mockProvider := providers.NewMockProvider()
	calendarPane := panes.NewCalendarPane(mockProvider)
	
	// Register with compile-time type safety
	calendarRegistry.RegisterTypedPane(calendarPane)
	
	// Retrieve with compile-time type safety
	retrievedPane, exists := calendarRegistry.GetTypedPane("calendar")
	if !exists {
		t.Fatal("Calendar pane should exist")
	}
	
	if retrievedPane.ID() != "calendar" {
		t.Errorf("Expected ID 'calendar', got '%s'", retrievedPane.ID())
	}
	
	// Get typed data without type assertions
	ctx := context.Background()
	data, err := calendarRegistry.GetTypedData(ctx, "calendar")
	if err != nil {
		t.Fatalf("GetTypedData failed: %v", err)
	}
	
	// Compile-time type safety - no type assertion needed
	if len(data.Events) == 0 {
		t.Error("Expected events in calendar data")
	}
	
	if data.Count != len(data.Events) {
		t.Error("Count should match events length")
	}
}

// TestGenericPaneManager demonstrates a comprehensive type-safe pane management system
func TestGenericPaneManager(t *testing.T) {
	baseRegistry := services.NewPaneRegistry()
	manager := services.NewGenericPaneManager(baseRegistry)
	
	// Set up test data
	mockProvider := providers.NewMockProvider()
	todoService := services.NewTodoService("/tmp/test_generic_manager.json")
	
	calendarPane := panes.NewCalendarPane(mockProvider)
	emailPane := panes.NewEmailPane(mockProvider)
	todoPane := panes.NewTodoPane(todoService)
	
	// Register panes with compile-time type safety
	manager.RegisterCalendarPane(calendarPane)
	manager.RegisterEmailPane(emailPane)
	manager.RegisterTodoPane(todoPane)
	
	ctx := context.Background()
	
	// Get calendar data with type safety
	calendarData, err := manager.GetCalendarData(ctx, "calendar")
	if err != nil {
		t.Fatalf("GetCalendarData failed: %v", err)
	}
	
	// No type assertion needed - we have compile-time guarantees
	if len(calendarData.Events) != 3 {
		t.Errorf("Expected 3 events, got %d", len(calendarData.Events))
	}
	
	// Test that we can work with the typed data directly
	for _, event := range calendarData.Events {
		if event.Title == "" {
			t.Error("Event should have title")
		}
		// event.Start is guaranteed to be time.Time
		if event.Start.IsZero() {
			t.Error("Event should have start time")
		}
	}
	
	// Get email data with type safety
	emailData, err := manager.GetEmailData(ctx, "email")
	if err != nil {
		t.Fatalf("GetEmailData failed: %v", err)
	}
	
	if len(emailData.Emails) != 3 {
		t.Errorf("Expected 3 emails, got %d", len(emailData.Emails))
	}
	
	// Get todo data with type safety
	todoData, err := manager.GetTodoData(ctx, "todos")
	if err != nil {
		t.Fatalf("GetTodoData failed: %v", err)
	}
	
	// todoData.Todos is guaranteed to be []models.Todo
	if todoData.Count != len(todoData.Todos) {
		t.Error("Count should match todos length")
	}
}

// TestGenericPaneManager_NonexistentPane tests error handling
func TestGenericPaneManager_NonexistentPane(t *testing.T) {
	baseRegistry := services.NewPaneRegistry()
	manager := services.NewGenericPaneManager(baseRegistry)
	
	ctx := context.Background()
	
	// Try to get data from nonexistent pane
	_, err := manager.GetCalendarData(ctx, "nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent pane")
	}
	
	if err.Error() != "pane not found: nonexistent" {
		t.Errorf("Expected 'pane not found' error, got: %v", err)
	}
}

// TestCompileTimeTypeSafety demonstrates that the generic approach prevents
// common runtime errors at compile time
func TestCompileTimeTypeSafety(t *testing.T) {
	mockProvider := providers.NewMockProvider()
	calendarPane := panes.NewCalendarPane(mockProvider)
	
	// These would be compile errors if types don't match:
	
	// 1. Cannot register wrong type in typed registry
	baseRegistry := services.NewPaneRegistry()
	calendarRegistry := services.NewTypedPaneRegistry[models.CalendarPaneData](baseRegistry)
	calendarRegistry.RegisterTypedPane(calendarPane) // ✓ Correct type
	
	// This would be a compile error:
	// emailRegistry := services.NewTypedPaneRegistry[models.EmailPaneData](baseRegistry)
	// emailRegistry.RegisterTypedPane(calendarPane) // ✗ Wrong type
	
	// 2. Cannot assign wrong return type
	ctx := context.Background()
	data, err := calendarPane.GetTypedData(ctx)
	if err != nil {
		t.Fatal(err)
	}
	
	// This is compile-time safe:
	var calendarData models.CalendarPaneData = data // ✓ Correct type
	_ = calendarData.Events                         // ✓ Type-safe field access
	
	// This would be a compile error:
	// var emailData models.EmailPaneData = data // ✗ Wrong type
	
	t.Log("Compile-time type safety prevents runtime type assertion errors")
}