package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"flexplane/internal/handlers"
	"flexplane/internal/models"
	"flexplane/internal/providers"
	"flexplane/internal/services"
)

type AppConfig struct {
	Providers map[string]providers.ProviderConfig `json:"providers"`
	Panes     map[string]services.PaneConfig      `json:"panes"`
}

func main() {
	// Initialize factories
	providerFactory := providers.NewProviderFactory()
	paneFactory := services.NewPaneFactory()

	// Parse templates - include all template files
	tmpl := template.Must(template.ParseGlob("web/templates/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("web/templates/components/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("web/templates/panes/*.html"))

	// Load application configuration
	config := loadAppConfig()

	// Create providers based on configuration
	for name, providerConfig := range config.Providers {
		provider, err := providerFactory.CreateProvider(providerConfig)
		if err != nil {
			log.Fatalf("Failed to create provider %s: %v", name, err)
		}
		paneFactory.RegisterProvider(name, provider)
	}

	// Create pane registry
	registry := services.NewPaneRegistry()

	// Create and register panes based on configuration
	var enabledPanes []string
	var layoutConfig = make(map[string]services.PaneLayoutConfig)
	
	for paneID, paneConfig := range config.Panes {
		if paneConfig.Enabled {
			pane, err := paneFactory.CreatePane(paneConfig)
			if err != nil {
				log.Fatalf("Failed to create pane %s: %v", paneID, err)
			}
			registry.RegisterPane(pane)
			enabledPanes = append(enabledPanes, pane.ID())
			layoutConfig[pane.ID()] = paneConfig.Layout
		}
	}

	registry.SetEnabledPanes(enabledPanes)
	registry.SetLayoutConfig(layoutConfig)

	// Initialize handlers
	handler := handlers.NewHandler(registry, tmpl)

	// Routes
	http.HandleFunc("/", handler.Home)
	http.HandleFunc("/api/todos", handler.TodosAPI) // Legacy route for backward compatibility
	http.HandleFunc("/api/", handler.PaneAPI)       // Generic API route for all panes

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

	log.Println("Flexplane (extensible panes) server starting on :3000")
	log.Fatal(server.ListenAndServe())
}

func loadAppConfig() AppConfig {
	// Default configuration
	defaultConfig := AppConfig{
		Providers: map[string]providers.ProviderConfig{
			"default": {
				Type: "file",
				Args: map[string]interface{}{
					"todo_file": "data/todos.json",
				},
			},
		},
		Panes: map[string]services.PaneConfig{
			"calendar": {
				Type:     "calendar",
				Enabled:  true,
				Provider: "default",
				Layout: services.PaneLayoutConfig{
					GridArea: models.PaneGridArea{Row: "1", Column: "span 3"},
				},
			},
			"todos": {
				Type:     "todos",
				Enabled:  true,
				Provider: "default",
				Layout: services.PaneLayoutConfig{
					GridArea: models.PaneGridArea{Row: "2", Column: "span 6"},
				},
			},
			"email": {
				Type:     "email",
				Enabled:  true,
				Provider: "default",
				Layout: services.PaneLayoutConfig{
					GridArea: models.PaneGridArea{Row: "3", Column: "span 3"},
				},
			},
		},
	}

	// Try to load configuration from file
	if configData, err := os.ReadFile("config/app.json"); err == nil {
		var fileConfig AppConfig
		if err := json.Unmarshal(configData, &fileConfig); err == nil {
			// Merge configurations - file config overrides defaults
			for name, providerConfig := range fileConfig.Providers {
				defaultConfig.Providers[name] = providerConfig
			}
			for paneID, paneConfig := range fileConfig.Panes {
				defaultConfig.Panes[paneID] = paneConfig
			}
		}
	}

	return defaultConfig
}