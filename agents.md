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
- Examples in: `internal/handlers/handler.go`

### Code Standards  
- **Delete unused code** (don't comment out)
- **No versioned files** (`handler_v2.go`, etc.)
- **Prefer simplicity** over complex solutions
- **Tests must add value** (remove interface-only tests)
- **Configuration over hardcoding**

### Architecture
- Server-side rendering with Go templates
- Provider pattern for mock → real API swapping
- Minimal JS (progressive enhancement only)
- Target: ~6-19ms page loads

## Data Models
- **Todos**: `[{done: bool, message: string}]` 
- **Mock providers**: `internal/services/mock_provider.go`