package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"flexpane/internal/handlers"
	"flexpane/internal/panes"
	"flexpane/internal/providers"
	"flexpane/internal/services"
)

type PaneConfig struct {
	Enabled []string                                   `json:"enabled"`
	Layout  map[string]services.PaneLayoutConfig `json:"layout"`
}

func main() {
	// Initialize services
	todoService := services.NewTodoService("data/todos.json")

	// Create data provider
	dataProvider, err := providers.CreateProvider("mock")
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}

	// Parse templates - include all template files
	tmpl := template.Must(template.ParseGlob("web/templates/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("web/templates/components/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("web/templates/panes/*.html"))

	// Create pane registry
	registry := services.NewPaneRegistry()

	// Register available panes
	registry.RegisterPane(panes.NewCalendarPane(dataProvider))
	registry.RegisterPane(panes.NewTodoPane(todoService))
	registry.RegisterPane(panes.NewEmailPane(dataProvider))

	// Load pane configuration
	defaultPanes := []string{"calendar", "todos", "email"}
	registry.SetEnabledPanes(defaultPanes) // Set defaults first
	
	if configData, err := os.ReadFile("config/panes.json"); err == nil {
		var config PaneConfig
		if err := json.Unmarshal(configData, &config); err == nil {
			registry.SetEnabledPanes(config.Enabled)
			registry.SetLayoutConfig(config.Layout)
		}
	}

	// Initialize handlers
	handler := handlers.NewHandler(registry, tmpl)

	// Routes
	http.HandleFunc("/", handler.Home)
	http.HandleFunc("/api/todos", handler.TodosAPI) // Legacy route for backward compatibility
	// TODO: Add generic /api/{pane} route pattern for extensibility

	// Static files  
	// TODO: SECURITY - Static file serving vulnerable to directory traversal attacks (../../../etc/passwd)
	// Consider implementing path validation or using a more secure static file handler
	fs := http.FileServer(http.Dir("web/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start server
	server := &http.Server{
		Addr:         ":3000",
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Flexpane (extensible panes) server starting on :3000")
	log.Fatal(server.ListenAndServe())
}