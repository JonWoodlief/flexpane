# Flexplane

Personal productivity app with extensible panes (calendar, todos, email).

## Quick Start

```bash
go run main.go
open http://localhost:3000
```

## Data Providers

Flexplane supports multiple data sources:

- **Mock Provider** (default): Generates sample data for testing
- **Gmail Provider**: Connect to your Gmail inbox via Google APIs

### Gmail Setup

To use your real Gmail data:

1. Follow the [Gmail setup guide](docs/gmail-setup.md)
2. Update `config/providers.json` with your credentials  
3. Set `"default": "gmail"` in the providers config

## Documentation

- [`plan.md`](plan.md) - Full development plan and architecture
- [`AGENTS.md`](AGENTS.md) - Development conventions
- [`docs/gmail-setup.md`](docs/gmail-setup.md) - Gmail provider setup