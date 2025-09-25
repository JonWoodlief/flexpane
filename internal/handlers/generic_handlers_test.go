package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"flexpane/internal/models"
	"flexpane/internal/panes"
	"flexpane/internal/providers"
)

// TestGenericAPIHandler demonstrates the advantages of type-safe API handlers
func TestGenericAPIHandler(t *testing.T) {
	// Create a simple typed handler
	handler := NewGenericAPIHandler(func(ctx context.Context, req AddTodoRequest) (AddTodoResponse, error) {
		if req.Message == "" {
			return AddTodoResponse{Status: "error"}, nil
		}
		return AddTodoResponse{Status: "created", ID: "123"}, nil
	})
	
	// Test valid request
	reqData := AddTodoRequest{Message: "Test todo"}
	body, _ := json.Marshal(reqData)
	
	req := httptest.NewRequest("POST", "/api/test", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	
	err := handler.HandleHTTP(recorder, req)
	if err != nil {
		t.Fatalf("Handler failed: %v", err)
	}
	
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}
	
	var resp AddTodoResponse
	if err := json.NewDecoder(recorder.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	if resp.Status != "created" {
		t.Errorf("Expected status 'created', got '%s'", resp.Status)
	}
	
	if resp.ID != "123" {
		t.Errorf("Expected ID '123', got '%s'", resp.ID)
	}
}

// TestGenericAPIHandler_InvalidJSON tests error handling
func TestGenericAPIHandler_InvalidJSON(t *testing.T) {
	handler := NewGenericAPIHandler(func(ctx context.Context, req AddTodoRequest) (AddTodoResponse, error) {
		return AddTodoResponse{Status: "created"}, nil
	})
	
	// Send invalid JSON
	req := httptest.NewRequest("POST", "/api/test", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	
	err := handler.HandleHTTP(recorder, req)
	if err != nil {
		t.Errorf("Expected no error for bad JSON handling, got: %v", err)
	}
	
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", recorder.Code)
	}
	
	body := recorder.Body.String()
	if body != "Invalid JSON\n" {
		t.Errorf("Expected 'Invalid JSON' error, got '%s'", body)
	}
}

// TestTypedPaneAPIHandler demonstrates type-safe pane API handling
func TestTypedPaneAPIHandler(t *testing.T) {
	mockProvider := providers.NewMockProvider()
	calendarPane := panes.NewCalendarPane(mockProvider)
	
	// Create a typed API handler for the calendar pane
	handler := TypedPaneAPIHandler[models.CalendarPaneData](calendarPane)
	
	req := httptest.NewRequest("GET", "/api/calendar", nil)
	recorder := httptest.NewRecorder()
	
	handler.ServeHTTP(recorder, req)
	
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}
	
	contentType := recorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected JSON content type, got %s", contentType)
	}
	
	var data models.CalendarPaneData
	if err := json.NewDecoder(recorder.Body).Decode(&data); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	// We get compile-time type safety here
	if len(data.Events) == 0 {
		t.Error("Expected events in response")
	}
	
	if data.Count != len(data.Events) {
		t.Error("Count should match events length")
	}
	
	// Test that we can access typed fields without assertions
	for _, event := range data.Events {
		if event.Title == "" {
			t.Error("Event should have title")
		}
		// event.Start is guaranteed to be time.Time - no type assertion needed
		if event.Start.IsZero() {
			t.Error("Event should have start time")
		}
	}
}

// TestTypedPaneAPIHandler_MethodNotAllowed tests method restriction
func TestTypedPaneAPIHandler_MethodNotAllowed(t *testing.T) {
	mockProvider := providers.NewMockProvider()
	emailPane := panes.NewEmailPane(mockProvider)
	
	handler := TypedPaneAPIHandler[models.EmailPaneData](emailPane)
	
	req := httptest.NewRequest("POST", "/api/email", nil)
	recorder := httptest.NewRecorder()
	
	handler.ServeHTTP(recorder, req)
	
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", recorder.Code)
	}
}

// BenchmarkTypedVsUntyped demonstrates performance characteristics
func BenchmarkTypedVsUntyped(b *testing.B) {
	mockProvider := providers.NewMockProvider()
	calendarPane := panes.NewCalendarPane(mockProvider)
	ctx := context.Background()
	
	b.Run("Typed", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data, err := calendarPane.GetTypedData(ctx)
			if err != nil {
				b.Fatal(err)
			}
			// Access data directly - no type assertion needed
			_ = len(data.Events)
		}
	})
	
	b.Run("Untyped", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data, err := calendarPane.GetData(ctx)
			if err != nil {
				b.Fatal(err)
			}
			// Would need type assertion in real code
			typedData := data.(models.CalendarPaneData)
			_ = len(typedData.Events)
		}
	})
}