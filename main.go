package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"flexplane/internal/handlers"
	"flexplane/internal/panes"
	"flexplane/internal/services"
)

type PaneConfig struct {
	Enabled []string                                   `json:"enabled"`
	Layout  map[string]services.PaneLayoutConfig `json:"layout"`
}

func main() {
	// Initialize services
	todoService := services.NewTodoService("data/todos.json")
	mockProvider := services.NewMockProvider()

	// Parse templates - include all template files
	tmpl := template.Must(template.ParseGlob("web/templates/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("web/templates/components/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("web/templates/panes/*.html"))

	// Create pane registry
	registry := services.NewPaneRegistry()

	// Register available panes
	registry.RegisterPane(panes.NewCalendarPane(mockProvider))
	registry.RegisterPane(panes.NewTodoPane(todoService))
	registry.RegisterPane(panes.NewEmailPane(mockProvider))

	// Load pane configuration
	configData, err := os.ReadFile("config/panes.json")
	if err != nil {
		log.Printf("Could not read pane config, using defaults: %v", err)
		registry.SetEnabledPanes([]string{"calendar", "todos", "email"})
	} else {
		var config PaneConfig
		if err := json.Unmarshal(configData, &config); err != nil {
			log.Printf("Could not parse pane config, using defaults: %v", err)
			registry.SetEnabledPanes([]string{"calendar", "todos", "email"})
		} else {
			registry.SetEnabledPanes(config.Enabled)
			registry.SetLayoutConfig(config.Layout)
		}
	}

	// Initialize handlers
	handler := handlers.NewHandler(registry, tmpl)

	// Routes
	http.HandleFunc("/", handler.Home)
	http.HandleFunc("/api/todos", handler.TodosAPI)
	// TODO: Add pane-specific API endpoints

	// Static files
	fs := http.FileServer(http.Dir("web/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start server
	server := &http.Server{
		Addr:         ":3000",
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Flexplane (extensible panes) server starting on :3000")
	log.Fatal(server.ListenAndServe())
}