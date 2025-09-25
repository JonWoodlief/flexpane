# Flexplane

A lightweight personal productivity app combining calendar, todos, and email in a text-focused, extensible pane interface.

## Features

- **📅 Calendar** - View today's events and upcoming schedule
- **📝 Todos** - Simple task management with JSON persistence  
- **📧 Email** - Inbox preview with read/unread status
- **🔧 Extensible** - Modular pane architecture for easy customization
- **⚡ Fast** - Server-side rendered with minimal JavaScript (~6-19ms page loads)

## Tech Stack

- **Backend:** Go with `html/template` and JSON file storage
- **Frontend:** Responsive HTML/CSS with CSS Grid layout
- **Architecture:** Extensible pane registry with provider pattern
- **Data:** JSON files for persistence, mock providers for development

## Quick Start

```bash
# Clone the repository
git clone https://github.com/JonWoodlief/flexplane.git
cd flexplane

# Run the server
go run main.go

# Open in browser
open http://localhost:3000
```

The app will start with mock data for calendar and email, while todos are persisted to `data/todos.json`.

## Configuration

Customize enabled panes and layout in `config/panes.json`:

```json
{
  "enabled": ["calendar", "todos", "email"],
  "layout": {
    "calendar": {"grid_area": {"row": "1", "column": "span 3"}},
    "todos": {"grid_area": {"row": "2", "column": "span 6"}},
    "email": {"grid_area": {"row": "1", "column": "span 3"}}
  }
}
```

## Architecture

```
flexplane/
├── main.go              # Server entry point
├── internal/
│   ├── handlers/        # HTTP request handlers
│   ├── models/          # Data structures  
│   ├── panes/           # Pane implementations
│   ├── providers/       # Data providers (mock → real APIs)
│   └── services/        # Business logic
├── web/
│   ├── templates/       # HTML templates
│   └── static/          # CSS and JavaScript
├── data/               # JSON data files
└── config/             # Pane configuration
```

## Development

### Provider Pattern
Easily swap mock data providers for real API integrations:

```go
// Mock provider (current)
mockProvider := services.NewMockProvider()

// Future: Real API provider  
realProvider := services.NewGraphProvider(credentials)
```

### Adding New Panes
1. Implement the `Pane` interface in `internal/panes/`
2. Register in `main.go`: `registry.RegisterPane(yourPane)`
3. Add template in `web/templates/panes/`
4. Configure in `config/panes.json`

### Performance Target
- Page load: ~6-19ms (localhost + template rendering + minimal JS)
- Architecture: Server-side rendering with progressive enhancement

## Roadmap

- **Phase 1:** ✅ Foundation with extensible pane architecture
- **Phase 2:** Real API integration (Microsoft Graph)
- **Phase 3:** Native app wrapper (Tauri/Wails)

## Documentation

- [`plan.md`](plan.md) - Complete development plan and architecture decisions
- [`agents.md`](agents.md) - Development conventions and workflow guide

---

*Flexplane: Extensible panes for personal productivity*