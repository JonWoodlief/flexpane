# Flexplane Extensibility Guide

This document demonstrates how the new architecture enables easy extension of the Flexplane application with new pane types and data providers.

## Architecture Overview

The refactored architecture separates concerns into:

- **Providers**: Data sources that implement the `Provider` interface
- **Panes**: UI components that implement the `Pane` interface  
- **Factories**: Registration and creation systems for providers and panes
- **Configuration**: JSON-based configuration for enabling/configuring components

## Understanding Providers

A **provider** is a data source that implements our `DataProvider` interface. It fetches data from various sources:

- **Production Providers**: Connect to real services (APIs, databases, files)
  - `null`: Returns empty data (useful when integrations aren't configured)
  - Custom providers like `outlook`, `gmail` for real integrations
  
- **Development Providers**: For testing and development
  - `mock`: Returns fake demo data for all data types

### Provider Types

All provider types are available through the standard factory - no special handling needed. Simply configure the provider type in your JSON configuration.

#### Available Providers
- **`null`**: Returns empty data for all types - used when integrations aren't set up
- **`mock`**: Returns realistic fake data for demos and development
- **Custom providers**: Add your own by registering with the factory

Providers are configured purely through JSON - no code changes required to switch between production and development modes.

## Adding a New Provider

### 1. Implement the Provider Interface

```go
// internal/providers/weather_provider.go
package providers

import "flexplane/internal/models"

type WeatherProvider struct {
    apiKey string
    city   string
}

func NewWeatherProvider(apiKey, city string) *WeatherProvider {
    return &WeatherProvider{apiKey: apiKey, city: city}
}

// Implement required Provider interface methods
func (wp *WeatherProvider) GetCalendarEvents() ([]models.Event, error) { 
    return []models.Event{}, nil // Not supported
}

func (wp *WeatherProvider) GetEmails() ([]models.Email, error) { 
    return []models.Email{}, nil // Not supported
}

func (wp *WeatherProvider) GetTodos() []models.Todo { 
    return []models.Todo{} // Not supported
}

func (wp *WeatherProvider) AddTodo(message string) error { 
    return nil // Not supported
}

func (wp *WeatherProvider) ToggleTodo(index int) error { 
    return nil // Not supported
}

// Add weather-specific methods as needed
func (wp *WeatherProvider) GetWeather() (WeatherData, error) {
    // Implementation would call weather API
    return WeatherData{
        Temperature: "72°F",
        Condition:   "Sunny", 
        Location:    wp.city,
    }, nil
}
```

### 2. Register the Provider in the Factory

```go
// In main.go or initialization code
providerFactory := providers.NewProviderFactory()

providerFactory.RegisterProvider("weather", func(args map[string]interface{}) (providers.Provider, error) {
    apiKey := args["api_key"].(string)
    city := args["city"].(string)
    return NewWeatherProvider(apiKey, city), nil
})
```

### 3. Configure the Provider

#### Production Configuration (config/app.json)
```json
{
  "providers": {
    "weather": {
      "type": "weather",
      "args": {
        "api_key": "your_api_key_here",
        "city": "New York"
      }
    }
  }
}
```

#### Development Configuration (config/app-dev.json)
```json
{
  "providers": {
    "mock": {
      "type": "mock",
      "args": {}
    }
  },
  "panes": {
    "calendar": {
      "provider": "mock"
    }
  }
}
```

**Note**: Mock providers are automatically available when referenced in configuration. The factory detects this and loads the development factory with mock support.

## Adding a New Pane Type

### 1. Create the Pane Implementation

```go
// internal/panes/weather_pane.go
package panes

import (
    "context"
    "flexplane/internal/models"
    "flexplane/internal/providers"
)

type WeatherPane struct {
    provider providers.Provider
}

func NewWeatherPane(provider providers.Provider) *WeatherPane {
    return &WeatherPane{provider: provider}
}

func (wp *WeatherPane) ID() string { 
    return "weather" 
}

func (wp *WeatherPane) Title() string { 
    return "Weather" 
}

func (wp *WeatherPane) Template() string { 
    return "panes/weather.html" 
}

func (wp *WeatherPane) GetData(ctx context.Context) (interface{}, error) {
    // Type assert to get weather-specific methods
    if weatherProvider, ok := wp.provider.(*WeatherProvider); ok {
        weather, err := weatherProvider.GetWeather()
        if err != nil {
            return nil, err
        }
        
        return map[string]interface{}{
            "Weather": weather,
            "Location": weather.Location,
        }, nil
    }
    
    return map[string]interface{}{}, nil
}
```

### 2. Register the Pane Type

```go
// In main.go or initialization code
paneFactory := services.NewPaneFactory()

paneFactory.RegisterPaneType("weather", func(provider providers.Provider, args map[string]interface{}) models.Pane {
    return panes.NewWeatherPane(provider)
})
```

### 3. Configure the Pane

```json
// config/app.json
{
  "panes": {
    "weather": {
      "type": "weather",
      "enabled": true,
      "provider": "weather",
      "layout": {
        "grid_area": {
          "row": "1",
          "column": "span 2"
        }
      }
    }
  }
}
```

### 4. Create the Template

```html
<!-- web/templates/panes/weather.html -->
<div class="weather-pane">
    <h3>{{.Weather.Location}}</h3>
    <div class="temperature">{{.Weather.Temperature}}</div>
    <div class="condition">{{.Weather.Condition}}</div>
</div>
```

## Key Benefits of the New Architecture

### ✅ No Core Code Modification Required
- New panes and providers are added through configuration and new files
- Main application logic remains unchanged
- Factory pattern handles registration and instantiation

### ✅ Provider Flexibility  
- Providers can implement only the methods they support
- Composite providers can combine different data sources
- Runtime provider selection through configuration

### ✅ Comprehensive Interface Design
- Clear separation between data providers and UI panes
- Standardized error handling patterns
- Type-safe factory registration system

### ✅ Configuration-Driven Architecture
- Complete customization without code changes
- Provider and pane configurations externalized
- Layout and enablement controlled through JSON

### ✅ Backwards Compatibility
- Existing functionality preserved
- Legacy API routes maintained alongside new generic routes
- Gradual migration path for existing code

## API Access

All panes automatically get API access through the generic route:

```
GET /api/{pane_id}    # Get pane data
POST /api/{pane_id}   # Create/update (if pane supports APIHandler)
PATCH /api/{pane_id}  # Update (if pane supports APIHandler)
```

For example, the weather pane would be accessible at `/api/weather`.

## Testing New Components

The factory system makes testing straightforward:

```go
func TestWeatherPane(t *testing.T) {
    mockProvider := &WeatherProvider{city: "Test City"}
    pane := NewWeatherPane(mockProvider)
    
    data, err := pane.GetData(context.Background())
    if err != nil {
        t.Errorf("GetData failed: %v", err)
    }
    
    // Assert expected data structure
}
```

This architecture enables rapid development of new features while maintaining code quality and separation of concerns.