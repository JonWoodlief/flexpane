package providers

import (
	"fmt"
	"flexpane/internal/models"
)

// CreateProvider creates a data provider based on the specified provider type
// This replaces the verbose factory pattern with a simple switch statement
func CreateProvider(providerType string) (DataProvider, error) {
	switch providerType {
	case "mock":
		return NewMockProvider(), nil
	case "": // Default to mock if empty
		return NewMockProvider(), nil
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
}

// CreateTypedProvider creates a typed provider based on the specified provider type
func CreateTypedProvider[T any](providerType string) (TypedProvider[T], error) {
	baseProvider, err := CreateProvider(providerType)
	if err != nil {
		return nil, err
	}

	// This is a bit of a hack using type assertion, but provides type safety at compile time
	// In a real implementation, you'd have specific factory methods for each type
	var provider interface{}
	
	// Use type inference to determine the correct provider
	switch any(*new(T)).(type) {
	case []models.Event:
		provider = NewCalendarProvider(baseProvider)
	case []models.Email:
		provider = NewEmailProvider(baseProvider)
	default:
		return nil, fmt.Errorf("unsupported typed provider for type %T", *new(T))
	}

	typedProvider, ok := provider.(TypedProvider[T])
	if !ok {
		return nil, fmt.Errorf("failed to create typed provider for type %T", *new(T))
	}

	return typedProvider, nil
}
