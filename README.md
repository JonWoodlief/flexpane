# Flexplane

Extensible personal productivity app with configurable panes (calendar, todos, email).

## Quick Start

```bash
go run main.go
# Open http://localhost:3000
```

## Configuration

Configure providers and panes through JSON files:

- **`config/app.json`**: Production configuration 
- **`config/app-dev.json`**: Development configuration (with demo data)

Switch between configurations by changing provider types in JSON - no code changes needed.

## Documentation

- [`EXTENSIBILITY.md`](EXTENSIBILITY.md) - How to add new panes and providers