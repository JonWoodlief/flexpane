# Agent Guide - Flexplane

**Primary Documentation**: See `plan.md` for full architecture, phases, and design decisions.

## Quick Start
```bash
go run main.go    # Starts server on :8080
./scripts/screenshot.sh    # Take screenshot for Claude feedback
```

## Key Conventions

### Issue Tracking
- Use TODO comments: `// TODO: DESCRIPTION - explain impact and solution`

### Code Standards  
- **Delete all unused code.** Do not leave deprecated code or stubs
- **No versioned files** (`handler_v2.go`, etc.)
- **Prefer simplicity** over complex solutions
- **Tests must add value** (remove interface-only tests)
- **Configuration over hardcoding**

### Architecture - EXTENSIBLE Design
- **Panes can be added** - design interfaces accordingly
- **Provider pattern** for pluggable data sources
- **Server-side rendering** with Go templates

## Gmail Provider
- Set `GOOGLE_CLIENT_ID` environment variable to enable real Google data
- Falls back to mock data if not configured
- Supports OAuth 2.0 flow for distributed applications