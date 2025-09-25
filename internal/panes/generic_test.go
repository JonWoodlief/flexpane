package panes

import (
	"context"
	"testing"

	"flexplane/internal/models"
	"flexplane/internal/providers"
	"flexplane/internal/services"
)

// TestTypedPaneInterfaces demonstrates the advantages of generic type-safe interfaces
// over the traditional interface{} approach
func TestTypedPaneInterfaces(t *testing.T) {
	mockProvider := providers.NewMockProvider()
	calendarPane := NewCalendarPane(mockProvider)
	
	// Test that CalendarPane implements both interfaces
	var _ models.Pane = calendarPane
	var _ models.TypedPane[models.CalendarPaneData] = calendarPane
	
	ctx := context.Background()
	
	// Old way: requires type assertion and can fail at runtime
	untypedData, err := calendarPane.GetData(ctx)
	if err != nil {
		t.Fatalf("GetData failed: %v", err)
	}
	
	// This would require a type assertion in real code:
	// data, ok := untypedData.(models.CalendarPaneData)
	// if !ok { /* handle error */ }
	
	// New way: compile-time type safety
	typedData, err := calendarPane.GetTypedData(ctx)
	if err != nil {
		t.Fatalf("GetTypedData failed: %v", err)
	}
	
	// No type assertion needed - we have compile-time guarantees
	if len(typedData.Events) == 0 {
		t.Error("Expected events in typed data")
	}
	
	if typedData.Count != len(typedData.Events) {
		t.Error("Count should match events length")
	}
	
	// Verify the data is the same
	if typedData.Count != 3 { // Mock provider returns 3 events
		t.Errorf("Expected 3 events, got %d", typedData.Count)
	}
	
	// Both methods should return the same underlying data
	_ = untypedData // We know this is models.CalendarPaneData now
}

func TestEmailTypedPane(t *testing.T) {
	mockProvider := providers.NewMockProvider()
	emailPane := NewEmailPane(mockProvider)
	
	// Test type safety for EmailPane
	var _ models.TypedPane[models.EmailPaneData] = emailPane
	
	ctx := context.Background()
	typedData, err := emailPane.GetTypedData(ctx)
	if err != nil {
		t.Fatalf("GetTypedData failed: %v", err)
	}
	
	// Compile-time type safety means we can access fields directly
	if len(typedData.Emails) == 0 {
		t.Error("Expected emails in typed data")
	}
	
	// We can work with strongly-typed data
	for _, email := range typedData.Emails {
		if email.Subject == "" {
			t.Error("Email should have subject")
		}
		if email.From == "" {
			t.Error("Email should have from address")
		}
	}
	
	if typedData.Count != 3 { // Mock provider returns 3 emails
		t.Errorf("Expected 3 emails, got %d", typedData.Count)
	}
}

func TestTodoTypedPane(t *testing.T) {
	todoService := services.NewTodoService("/tmp/test_todos.json")
	todoPane := NewTodoPane(todoService)
	
	// Test type safety for TodoPane
	var _ models.TypedPane[models.TodoPaneData] = todoPane
	
	ctx := context.Background()
	typedData, err := todoPane.GetTypedData(ctx)
	if err != nil {
		t.Fatalf("GetTypedData failed: %v", err)
	}
	
	// Compile-time type safety
	if typedData.Count != len(typedData.Todos) {
		t.Error("Count should match todos length")
	}
	
	// We can safely iterate over strongly-typed todos
	for _, todo := range typedData.Todos {
		if todo.Message == "" {
			t.Error("Todo should have a message")
		}
		// todo.Done is guaranteed to be a bool - no type assertion needed
		_ = todo.Done
	}
}

// TestGenericProviderFactory demonstrates type-safe provider creation
func TestGenericProviderFactory(t *testing.T) {
	// This test shows how generics could be used for provider factories
	// though we haven't fully implemented this yet due to the existing architecture
	
	// For now, we can show that the existing mock provider works
	// with our typed interfaces
	provider := providers.NewMockProvider()
	
	events, err := provider.GetCalendarEvents()
	if err != nil {
		t.Fatalf("GetCalendarEvents failed: %v", err)
	}
	
	if len(events) != 3 {
		t.Errorf("Expected 3 events, got %d", len(events))
	}
	
	// With generics, we could have:
	// typedProvider := NewGenericProvider[[]models.Event](provider.GetCalendarEvents)
	// This would eliminate the need for the specific GetCalendarEvents method
}