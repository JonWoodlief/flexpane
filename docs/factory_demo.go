package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// This demo shows the advantages of factory pattern vs case statement
// Run with: go run docs/factory_demo.go

// Simple interface for demo
type DemoProvider interface {
	GetData() string
	GetType() string
}

// Mock implementation
type MockDemoProvider struct{}

func (m *MockDemoProvider) GetData() string { return "Mock data" }
func (m *MockDemoProvider) GetType() string { return "mock" }

// Email implementation (simulated)
type EmailDemoProvider struct {
	apiKey string
}

func (e *EmailDemoProvider) GetData() string {
	// Never log sensitive data like API keys
	return "Email data (configured with API key)"
}
func (e *EmailDemoProvider) GetType() string { return "email" }

// Configuration structures
type DemoProviderConfig struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config,omitempty"`
}

type DemoConfig struct {
	Providers map[string]DemoProviderConfig `json:"providers"`
	Default   string                        `json:"default"`
}

// APPROACH 1: Direct Case Statement (Limited)
func createProviderCaseStatement(providerType string, config map[string]interface{}) (DemoProvider, error) {
	switch providerType {
	case "mock":
		return &MockDemoProvider{}, nil
	case "email":
		apiKey, ok := config["api_key"].(string)
		if !ok {
			return nil, fmt.Errorf("email provider requires api_key")
		}
		return &EmailDemoProvider{apiKey: apiKey}, nil
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
}

// APPROACH 2: Factory Pattern (Flexible)
type DemoProviderFactory struct {
	config DemoConfig
}

func NewDemoProviderFactory(configData []byte) (*DemoProviderFactory, error) {
	var config DemoConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return &DemoProviderFactory{config: config}, nil
}

func (f *DemoProviderFactory) CreateProvider(name string) (DemoProvider, error) {
	providerConfig, exists := f.config.Providers[name]
	if !exists {
		if f.config.Default == "" {
			return nil, fmt.Errorf("provider '%s' not found and no default", name)
		}
		// Fallback to default
		providerConfig, exists = f.config.Providers[f.config.Default]
		if !exists {
			return nil, fmt.Errorf("default provider '%s' not found", f.config.Default)
		}
	}

	return createProviderCaseStatement(providerConfig.Type, providerConfig.Config)
}

func (f *DemoProviderFactory) GetAvailableProviders() []string {
	providers := make([]string, 0, len(f.config.Providers))
	for name := range f.config.Providers {
		providers = append(providers, name)
	}
	return providers
}

func main() {
	// Create demo configuration
	config1 := DemoConfig{
		Providers: map[string]DemoProviderConfig{
			"test-mock": {Type: "mock"},
			"gmail": {
				Type: "email",
				Config: map[string]interface{}{
					"api_key": "secret_gmail_key_12345",
				},
			},
		},
		Default: "test-mock",
	}

	config2 := DemoConfig{
		Providers: map[string]DemoProviderConfig{
			"production-email": {
				Type: "email",
				Config: map[string]interface{}{
					"api_key": "production_key_67890",
				},
			},
		},
		Default: "production-email",
	}

	fmt.Println("=== Factory Pattern Advantages Demo ===\n")

	// Demo 1: Configuration-driven flexibility
	fmt.Println("1. Configuration-driven flexibility:")
	for i, config := range []DemoConfig{config1, config2} {
		configBytes, _ := json.MarshalIndent(config, "", "  ")
		
		factory, err := NewDemoProviderFactory(configBytes)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\nConfiguration %d:\n%s\n", i+1, string(configBytes))
		
		provider, err := factory.GetDefaultProvider()
		if err != nil {
			log.Fatal(err)
		}
		
		fmt.Printf("Default provider type: %s\n", provider.GetType())
		fmt.Printf("Default provider data: %s\n", provider.GetData())
	}

	// Demo 2: Fallback behavior
	fmt.Println("\n2. Fallback behavior:")
	configBytes, _ := json.Marshal(config1)
	factory, _ := NewDemoProviderFactory(configBytes)
	
	// Try non-existent provider - should fall back to default
	provider, err := factory.CreateProvider("non-existent")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Fallback provider type: %s\n", provider.GetType())
	}

	// Demo 3: Multiple configured instances
	fmt.Println("\n3. Multiple configured instances of same type:")
	multiConfig := DemoConfig{
		Providers: map[string]DemoProviderConfig{
			"gmail": {
				Type: "email",
				Config: map[string]interface{}{"api_key": "gmail_key_123"},
			},
			"outlook": {
				Type: "email", // Same type, different config
				Config: map[string]interface{}{"api_key": "outlook_key_456"},
			},
		},
	}
	
	configBytes, _ = json.Marshal(multiConfig)
	factory, _ = NewDemoProviderFactory(configBytes)
	
	gmailProvider, _ := factory.CreateProvider("gmail")
	outlookProvider, _ := factory.CreateProvider("outlook")
	
	fmt.Printf("Gmail provider: %s\n", gmailProvider.GetData())
	fmt.Printf("Outlook provider: %s\n", outlookProvider.GetData())

	// Demo 4: Runtime discovery
	fmt.Println("\n4. Runtime discovery:")
	fmt.Printf("Available providers: %v\n", factory.GetAvailableProviders())

	fmt.Println("\n=== Case Statement Limitations ===")
	fmt.Println("- Requires code changes to add new provider configurations")
	fmt.Println("- No fallback logic built-in")
	fmt.Println("- Cannot have multiple instances of same type with different configs") 
	fmt.Println("- No runtime discovery of available options")
	fmt.Println("- Configuration hardcoded in switch statement")

	// Save example configurations for reference
	saveExampleConfigs(config1, config2)
}

func saveExampleConfigs(configs ...DemoConfig) {
	os.MkdirAll("docs/examples", 0755)
	
	for i, config := range configs {
		configBytes, _ := json.MarshalIndent(config, "", "  ")
		filename := filepath.Join("docs", "examples", fmt.Sprintf("provider_config_%d.json", i+1))
		os.WriteFile(filename, configBytes, 0644)
	}
	fmt.Printf("\n✓ Example configurations saved to docs/examples/\n")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (f *DemoProviderFactory) GetDefaultProvider() (DemoProvider, error) {
	if f.config.Default == "" {
		return nil, fmt.Errorf("no default provider configured")
	}
	return f.CreateProvider(f.config.Default)
}