package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

	// Parse templates - include all template files with error handling
	tmpl, err := template.ParseGlob("web/templates/*.html")
	if err != nil {
		log.Fatalf("Failed to parse main templates: %v", err)
	}
	
	tmpl, err = tmpl.ParseGlob("web/templates/components/*.html")
	if err != nil {
		log.Fatalf("Failed to parse component templates: %v", err)
	}
	
	tmpl, err = tmpl.ParseGlob("web/templates/panes/*.html")
	if err != nil {
		log.Fatalf("Failed to parse pane templates: %v", err)
	}

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

	// Static files - secure file serving to prevent directory traversal
	staticDir := http.Dir("web/static/")
	fs := http.FileServer(staticDir)
	http.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent directory traversal by cleaning the path
		if strings.Contains(r.URL.Path, "..") {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		fs.ServeHTTP(w, r)
	})))

	// Start server with graceful shutdown
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Flexplane (extensible panes) server starting on :%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup failed: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}