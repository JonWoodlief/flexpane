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
- **Delete unused code** (don't comment out)
- **No versioned files** (`handler_v2.go`, etc.)
- **Prefer simplicity** over complex solutions
- **Tests must add value** (remove interface-only tests)
- **Configuration over hardcoding**

### Architecture - EXTENSIBLE Design
- **Agents may be added or removed** - design interfaces accordingly
- **Provider pattern** for pluggable data sources
- **Pane registry system** for modular UI components
- **Server-side rendering** with Go templates
- **Minimal JS** (progressive enhancement only)
- Target: ~6-19ms page loads