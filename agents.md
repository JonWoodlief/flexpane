# Agent Guide - Flexplane

## Overview
Flexplane is a lightweight personal productivity app (calendar + todos + email) built with Go server + HTML templates.

## Key References
- **Primary Documentation**: `plan.md` - Contains full architecture, phases, and design decisions
- **Project Structure**: See `plan.md` Architecture section for complete file layout

## Development Conventions

### Issue Tracking
- **Use TODO comments** in code instead of GitHub issues
- Format: `// TODO: DESCRIPTION - explain impact and solution`
- Examples in: `internal/handlers/handler.go`

### Code Cleanup
- **Clean up deprecated files** when refactoring architecture
- **No versioned files** (no `handler_v2.go`, `main_v2.go` etc.)
- **Use git** if concerned about data loss during major refactors
- Keep codebase minimal and focused

### Code Review Preferences
- **Eliminate unneeded code** ruthlessly - delete rather than comment out
- **Remove low-value code** that doesn't add meaningful functionality
- **Prefer simplicity** - complex solutions must justify their complexity
- **Value extensibility** over premature optimization
- **Tests must add value** - remove tests that just verify interface contracts
- **Configuration over hardcoding** for layout and behavior

### Architecture Principles
- **Server-side rendering**: Go templates with pre-populated data
- **Passthrough caching**: Go server caches API responses, serves rendered HTML
- **Progressive enhancement**: Minimal JS for interactivity only (no data fetching)
- **Provider pattern**: Easy swapping mock → real APIs

### Performance Target
- Page load: ~6-19ms (localhost server + template rendering + minimal JS)

## Quick Start
```bash
cd flexplane
go run main.go    # Starts server on :8080
```

## Feedback Loop
- Uses Firefox screenshots for visual feedback with Claude
- Script: `scripts/screenshot.sh` (auto-opens browser, captures UI)
- Screenshots saved to `screenshots/current-ui.png`

## Data Models
- **Todos**: Simple JSON array `[{done: bool, message: string}]`
- **Calendar/Email**: Mock providers (see `internal/services/mock_provider.go`)
- **Priority**: Array order (not explicit field)

## Current Phase
Phase 1: Foundation - building basic server structure and pane layout with mock data.

## Tech Stack
Go, html/template, JSON files, CSS Grid, minimal vanilla JS