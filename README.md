# Flexplane

Personal productivity app with extensible panes (calendar, todos, email).

## Quick Start

```bash
go run main.go
open http://localhost:3000
```

## Configuration

The app uses JSON configuration files:

- **`config/app.json`**: Production configuration (no mock data)
- **`config/app-dev.json`**: Development configuration (with demo data)

### Provider Types

- **`file`**: File-based storage (todos from JSON, calendar/email disabled)
- **`null`**: Empty data provider (for when integrations aren't configured)
- **`mock`**: Demo data provider (automatically enabled in development)

The app automatically detects if mock providers are needed and switches to development mode.

## Documentation

- [`EXTENSIBILITY.md`](EXTENSIBILITY.md) - How to add new panes and providers
- [`plan.md`](plan.md) - Full development plan and architecture
- [`agents.md`](agents.md) - Development conventions