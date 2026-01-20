package providers

import (
	"context"
	"fmt"
	"time"
)

// CompletionRequest represents a request to an AI provider
type CompletionRequest struct {
	System    string         `json:"system"`
	Prompt    string         `json:"prompt"`
	Model     string         `json:"model"`
	MaxTokens int            `json:"max_tokens,omitempty"`
	Stream    bool           `json:"stream,omitempty"`
	Options   map[string]any `json:"options,omitempty"`
}

// CompletionResponse represents a response from an AI provider
type CompletionResponse struct {
	Content  string        `json:"content"`
	Model    string        `json:"model"`
	Tokens   int           `json:"tokens"`
	Duration time.Duration `json:"duration"`
	Error    error         `json:"-"`
}

// StreamChunk represents a chunk of streamed response
type StreamChunk struct {
	Content string
	Done    bool
	Error   error
}

// Provider defines the interface for AI providers
type Provider interface {
	// Name returns the provider's identifier
	Name() string

	// Execute sends a completion request and returns the response
	Execute(ctx context.Context, request CompletionRequest) (*CompletionResponse, error)

	// ExecuteStream sends a completion request and streams the response
	ExecuteStream(ctx context.Context, request CompletionRequest) (<-chan StreamChunk, error)

	// IsAvailable checks if the provider is configured and accessible
	IsAvailable() bool

	// GetModels returns available models for this provider
	GetModels() []string
}

// ProviderConfig represents configuration for a provider
type ProviderConfig struct {
	Name   string         `yaml:"name"`
	Type   string         `yaml:"type"` // http, ollama, executable
	Config map[string]any `yaml:"config"`
}

// providerRegistry holds all registered providers
var providerRegistry = make(map[string]Provider)

// RegisterProvider adds a provider to the registry
func RegisterProvider(provider Provider) {
	providerRegistry[provider.Name()] = provider
}

// GetProvider returns a provider by name
func GetProvider(name string) (Provider, error) {
	if provider, ok := providerRegistry[name]; ok {
		return provider, nil
	}
	return nil, fmt.Errorf("provider not found: %s", name)
}

// ListProviders returns all registered provider names
func ListProviders() []string {
	names := make([]string, 0, len(providerRegistry))
	for name := range providerRegistry {
		names = append(names, name)
	}
	return names
}

// ListAvailableProviders returns all providers that are currently available
func ListAvailableProviders() []Provider {
	available := make([]Provider, 0)
	for _, provider := range providerRegistry {
		if provider.IsAvailable() {
			available = append(available, provider)
		}
	}
	return available
}

// NewProvider creates a provider from configuration
func NewProvider(config ProviderConfig) (Provider, error) {
	switch config.Type {
	case "http", "openai":
		return NewHTTPProvider(config.Name, config.Config)
	case "ollama":
		return NewOllamaProvider(config.Name, config.Config)
	case "anthropic":
		return NewAnthropicProvider(config.Name, config.Config)
	case "executable":
		return NewExecutableProvider(config.Name, config.Config)
	default:
		return nil, fmt.Errorf("unknown provider type: %s", config.Type)
	}
}

// LoadProviders loads providers from configuration and registers them
func LoadProviders(configs []ProviderConfig) error {
	for _, cfg := range configs {
		provider, err := NewProvider(cfg)
		if err != nil {
			return fmt.Errorf("failed to create provider %s: %w", cfg.Name, err)
		}
		RegisterProvider(provider)
	}
	return nil
}
