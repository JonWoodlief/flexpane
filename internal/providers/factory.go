package providers

import (
	"fmt"
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
