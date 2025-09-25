# Flexplane Extensibility Guide

This document explains how to extend Flexplane with new pane types and data providers.

## Architecture Overview

The refactored architecture separates concerns into:

- **DataProviders**: Handle email/calendar data from external sources (Outlook, Gmail, etc.)
- **TodoService**: Independent todo management system (easily forkable/customizable)
- **Panes**: UI components that display data
- **Factories**: Registration systems for providers and panes
- **Configuration**: JSON-based configuration for all components

## Understanding Providers

**DataProviders** focus exclusively on email/calendar data - perfect for shared OAuth scenarios like Outlook/Gmail integration.

**Available Provider Types:**
- **`mock`**: Returns demo data for development/testing
- **Custom providers**: Implement DataProvider interface for real integrations

**TodoService** handles todos independently outside the provider system - can be forked/customized without affecting email/calendar providers.

## Adding a New Email/Calendar Provider

### 1. Implement DataProvider Interface
Create a new provider in `internal/providers/` that implements:
- `GetCalendarEvents() ([]models.Event, error)`
- `GetEmails() ([]models.Email, error)`

### 2. Register with Factory
Add registration in the provider factory initialization.

### 3. Configure via JSON
Add your provider configuration to `config/app.json`.

## Adding a New Pane Type  

### 1. Implement Pane Interface
Create a new pane in `internal/panes/` that implements:
- `ID() string`
- `Title() string` 
- `GetData(ctx context.Context) (interface{}, error)`
- `Template() string`

### 2. Register with PaneFactory
Add registration in the pane factory initialization.

### 3. Configure via JSON
Add pane configuration to enable and position it.

### 4. Create Template
Add HTML template in `web/templates/panes/`.

## Key Benefits

### ✅ Pure Configuration
- Switch between providers entirely through JSON files
- No code changes needed for different environments
- Enable/disable panes through configuration

### ✅ Clean Separation
- Email/calendar providers share OAuth authentication
- Todos completely independent (forkable without affecting providers)
- UI panes separated from data sources

### ✅ Extensibility
- Add new email/calendar providers without core changes
- Fork todo system independently
- Generic API routing supports any pane type

### ✅ Real-World Alignment
- Provider system designed for Outlook/Gmail integration
- Todo system designed for easy customization
- No overengineered abstractions

## API Access

All panes automatically get API endpoints:
- `GET /api/{pane_id}` - Get pane data  
- `POST /api/{pane_id}` - Create/update (if pane supports it)

The architecture enables rapid development while maintaining clean separation of concerns.