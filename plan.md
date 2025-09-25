# Flexplane Development Plan

A lightweight personal productivity app combining calendar, todos, and email in a text-focused interface.

## Project Overview

**Goal**: Build a Go-served HTML app with pane-based layout for calendar (top), todos (bottom), and future email integration.

**Tech Stack**:
- Backend: Go with html/template, JSON file storage
- Frontend: Lightweight HTML/CSS with responsive design
- Development: Claude feedback loop via Firefox screenshots
- Future: Native app wrapper (Tauri/Wails)

## Architecture

```
flexplane/
├── plan.md                    # This file
├── main.go                    # Server entry point
├── go.mod                     # Go module definition
├── internal/
│   ├── handlers/              # HTTP request handlers
│   ├── models/                # Data structures
│   ├── providers/             # Pluggable data providers (mock → real)
│   └── services/              # Business logic
├── web/
│   ├── templates/
│   │   └── index.html         # Main layout with panes
│   └── static/
│       ├── css/style.css      # Text-focused styling
│       └── js/app.js          # Minimal JavaScript
├── data/
│   └── todos.json             # Persistent todo storage
├── screenshots/               # UI feedback screenshots
└── scripts/
    └── screenshot.sh          # Automated screenshot tool
```

## Development Phases

### Phase 1: Foundation (MVP)
1. **Project Setup**
   - [x] Create project structure
   - [ ] Initialize Go module (`flexplane`)
   - [ ] Basic HTTP server with static file serving
   - [ ] HTML template rendering setup

2. **Claude Feedback Loop**
   - [ ] Create screenshot automation script for Firefox
   - [ ] Test screenshot → Claude analysis workflow
   - [ ] Establish iterative UI improvement process

3. **Core Layout**
   - [ ] HTML pane-based layout (calendar top, todos bottom)
   - [ ] CSS Grid/Flexbox for responsive panes
   - [ ] Text-focused, minimal styling
   - [ ] Easy pane resizing via CSS adjustments

4. **Mock Data & Providers**
   - [ ] Pluggable provider interface
   - [ ] Mock calendar provider (realistic events)
   - [ ] Mock email provider (inbox simulation)
   - [ ] Todo JSON persistence in `data/todos.json`

### Phase 2: Core Functionality
5. **Todo Operations**
   - [ ] CRUD operations for todos
   - [ ] Simple JSON structure: `{done: bool, message: string}`
   - [ ] Priority via array order
   - [ ] Real-time updates without page refresh

6. **Calendar Integration**
   - [ ] Display today's events prominently
   - [ ] Week view in compact format
   - [ ] Mock → real provider swap preparation

7. **Email Preview**
   - [ ] Inbox preview pane
   - [ ] Read/unread status
   - [ ] Expandable email content

### Phase 3: Polish & Integration
8. **UI Refinement**
   - [ ] Responsive design (mobile stacking)
   - [ ] Keyboard navigation
   - [ ] Accessibility improvements
   - [ ] Performance optimization

9. **Real Provider Integration**
   - [ ] Microsoft Graph API setup
   - [ ] OAuth 2.0 authentication flow
   - [ ] Real calendar/email data
   - [ ] Error handling and offline fallback

10. **Native App Preparation**
    - [ ] Evaluate Tauri vs Wails
    - [ ] **Research**: Does native webview need network calls or support direct function calls?
    - [ ] Package web app for native wrapper
    - [ ] macOS app bundle creation

## Data Models

### Todo Structure
```json
[
  {"done": false, "message": "Review quarterly budget"},
  {"done": true, "message": "Update team calendar"}
]
```

### Calendar Event (Mock)
```go
type Event struct {
    ID       string    `json:"id"`
    Title    string    `json:"title"`
    Start    time.Time `json:"start"`
    End      time.Time `json:"end"`
    Location string    `json:"location,omitempty"`
}
```

### Email Message (Mock)
```go
type Email struct {
    ID      string    `json:"id"`
    Subject string    `json:"subject"`
    From    string    `json:"from"`
    Preview string    `json:"preview"`
    Time    time.Time `json:"time"`
    Read    bool      `json:"read"`
}
```

## Development Workflow

### Claude Feedback Loop
1. **Development Cycle**:
   ```bash
   # Start server
   go run main.go

   # Make changes to templates/CSS
   # Auto-screenshot with script
   ./scripts/screenshot.sh

   # Share screenshot with Claude for feedback
   # Implement suggestions, repeat
   ```

2. **Screenshot Strategy**:
   - Firefox automation for consistency
   - Single `current-ui.png` file (overwrite)
   - Capture full window at 1200x800
   - Include before/after comparisons for major changes

### Code Organization Principles
- **Separation of Concerns**: Handlers, services, models clearly divided
- **Provider Pattern**: Easy swapping of mock → real data sources
- **Template-Driven**: Server-side rendering with minimal client JS
- **Configuration**: Environment-based settings for dev/prod

## Caching & Performance Strategy

### Passthrough Cache + Progressive Enhancement
**Architecture:**
- **Go server**: Acts as smart proxy, caches API responses, serves pre-rendered HTML
- **Progressive JS**: Adds interactivity to server-rendered content (no data fetching)
- **Background sync**: Optional goroutine keeps cache fresh

**Performance Profile:**
- Initial page load: `~5-16ms` (localhost Go server + template rendering)
- JS enhancement: `+1-3ms` (lightweight interactivity, no DOM building)
- **Total perceived load**: `~6-19ms` with full functionality

**Why This Wins:**
- Users see complete content immediately (no loading states)
- Interactions feel instant (optimistic UI updates)
- Simple architecture (server owns all data)
- Easy debugging (all rendering in one place)

**Implementation:**
```
┌─ API Cache ─┐    ┌─ Template Render ─┐    ┌─ Progressive JS ─┐
│ Background  │ →  │ Pre-populated     │ →  │ Add click       │
│ sync or     │    │ HTML served       │    │ handlers,       │
│ on-demand   │    │ instantly         │    │ optimistic UI   │
└─────────────┘    └───────────────────┘    └─────────────────┘
```

## Key Design Decisions

### Why These Choices?
- **Go + Templates**: Fast, simple, easy deployment
- **Passthrough Caching**: Best performance with simple architecture
- **JSON File Storage**: Lightweight, no database overhead initially
- **Extensible Pane Architecture**: Registry pattern for modular UI components
- **Mock-First**: Rapid development without external dependencies
- **Text-Focused UI**: Fast loading, accessible, distraction-free

### Extensible Pane Architecture (Added)
**Decision**: Implement registry-based pane system with Go template composition
**Rationale**:
- Easy to add/remove panes without touching core layout
- Each pane is self-contained (data fetching + template)
- CSS Grid handles responsive sizing automatically
- Configuration-driven (JSON controls enabled panes)

**Implementation**:
- `Pane` interface: ID, Title, Size, Order, GetData, Template
- `PaneRegistry`: Manages enabled panes and data fetching
- Template structure: `panes/calendar.html`, `panes/todos.html`, etc.
- CSS classes: `.pane-small`, `.pane-medium`, `.pane-large`

### Future Considerations
- **Database Migration**: SQLite → PostgreSQL for multi-user
- **Real-time Updates**: WebSocket or SSE for live data
- **Plugin System**: Easy addition of new panes/providers
- **Sync**: Cross-device data synchronization
- **Native App Integration**: Research if webview can bypass network and call Go functions directly

## Getting Started

```bash
# Create project
mkdir flexplane && cd flexplane

# Initialize Go module
go mod init flexplane

# Run development server
go run main.go

# Take screenshot for Claude feedback
./scripts/screenshot.sh
```

## Success Criteria

**Phase 1 Complete When**:
- [x] Basic server serves HTML with calendar/todo panes
- [ ] Screenshot feedback loop working with Claude
- [ ] Mock data displays realistically in both panes
- [ ] Todos persist to JSON file
- [ ] CSS allows easy pane size adjustments

**Phase 2 Complete When**:
- [ ] Full CRUD operations for todos working
- [ ] Responsive design works on mobile
- [ ] Provider interface ready for real API swap
- [ ] Performance is snappy (<100ms page loads)

Ready to start implementation!