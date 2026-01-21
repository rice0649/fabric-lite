package providers

import (
	"context"
	"testing"
	"time"
)

func TestCompletionRequest(t *testing.T) {
	req := CompletionRequest{
		System:    "system prompt",
		Prompt:    "user prompt",
		Model:     "test-model",
		MaxTokens: 1000,
		Stream:    true,
		Options:   map[string]any{"temperature": 0.7},
	}

	if req.System != "system prompt" {
		t.Errorf("Expected System to be 'system prompt', got %s", req.System)
	}
	if req.Prompt != "user prompt" {
		t.Errorf("Expected Prompt to be 'user prompt', got %s", req.Prompt)
	}
	if req.Model != "test-model" {
		t.Errorf("Expected Model to be 'test-model', got %s", req.Model)
	}
	if req.MaxTokens != 1000 {
		t.Errorf("Expected MaxTokens to be 1000, got %d", req.MaxTokens)
	}
	if !req.Stream {
		t.Errorf("Expected Stream to be true, got %v", req.Stream)
	}
	if req.Options["temperature"] != 0.7 {
		t.Errorf("Expected Options[temperature] to be 0.7, got %v", req.Options["temperature"])
	}
}

func TestCompletionResponse(t *testing.T) {
	duration := 100 * time.Millisecond
	resp := CompletionResponse{
		Content:  "test response",
		Model:    "test-model",
		Tokens:   50,
		Duration: duration,
		Error:    nil,
	}

	if resp.Content != "test response" {
		t.Errorf("Expected Content to be 'test response', got %s", resp.Content)
	}
	if resp.Model != "test-model" {
		t.Errorf("Expected Model to be 'test-model', got %s", resp.Model)
	}
	if resp.Tokens != 50 {
		t.Errorf("Expected Tokens to be 50, got %d", resp.Tokens)
	}
	if resp.Duration != duration {
		t.Errorf("Expected Duration to be %v, got %v", duration, resp.Duration)
	}
	if resp.Error != nil {
		t.Errorf("Expected Error to be nil, got %v", resp.Error)
	}
}

func TestStreamChunk(t *testing.T) {
	chunk := StreamChunk{
		Content: "test chunk",
		Done:    true,
		Error:   nil,
	}

	if chunk.Content != "test chunk" {
		t.Errorf("Expected Content to be 'test chunk', got %s", chunk.Content)
	}
	if !chunk.Done {
		t.Errorf("Expected Done to be true, got %v", chunk.Done)
	}
	if chunk.Error != nil {
		t.Errorf("Expected Error to be nil, got %v", chunk.Error)
	}

	errorChunk := StreamChunk{
		Content: "",
		Done:    false,
		Error:   context.Canceled,
	}

	if errorChunk.Error == nil {
		t.Errorf("Expected Error to be non-nil, got nil")
	}
}

func TestProviderConfig(t *testing.T) {
	config := ProviderConfig{
		Name: "test-provider",
		Type: "http",
		Config: map[string]any{
			"endpoint": "https://api.example.com",
			"api_key":  "test-key",
		},
	}

	if config.Name != "test-provider" {
		t.Errorf("Expected Name to be 'test-provider', got %s", config.Name)
	}
	if config.Type != "http" {
		t.Errorf("Expected Type to be 'http', got %s", config.Type)
	}
	if config.Config["endpoint"] != "https://api.example.com" {
		t.Errorf("Expected endpoint to be 'https://api.example.com', got %v", config.Config["endpoint"])
	}
	if config.Config["api_key"] != "test-key" {
		t.Errorf("Expected api_key to be 'test-key', got %v", config.Config["api_key"])
	}
}

// Mock provider for testing
type mockProvider struct {
	name      string
	available bool
	models    []string
}

func (m *mockProvider) Name() string {
	return m.name
}

func (m *mockProvider) Execute(ctx context.Context, request CompletionRequest) (*CompletionResponse, error) {
	return &CompletionResponse{
		Content:  "mock response",
		Model:    m.models[0],
		Tokens:   10,
		Duration: 10 * time.Millisecond,
	}, nil
}

func (m *mockProvider) ExecuteStream(ctx context.Context, request CompletionRequest) (<-chan StreamChunk, error) {
	chunks := make(chan StreamChunk, 1)
	go func() {
		defer close(chunks)
		chunks <- StreamChunk{Content: "mock response", Done: true}
	}()
	return chunks, nil
}

func (m *mockProvider) IsAvailable() bool {
	return m.available
}

func (m *mockProvider) GetModels() []string {
	return m.models
}

func TestProviderRegistry(t *testing.T) {
	// Clear registry for testing
	providerRegistry = make(map[string]Provider)

	mockProv1 := &mockProvider{
		name:      "mock1",
		available: true,
		models:    []string{"model1", "model2"},
	}

	mockProv2 := &mockProvider{
		name:      "mock2",
		available: false,
		models:    []string{"model3"},
	}

	// Test RegisterProvider
	RegisterProvider(mockProv1)
	RegisterProvider(mockProv2)

	if len(providerRegistry) != 2 {
		t.Errorf("Expected registry to have 2 providers, got %d", len(providerRegistry))
	}

	// Test GetProvider
	prov, err := GetProvider("mock1")
	if err != nil {
		t.Errorf("Expected to find mock1, got error: %v", err)
	}
	if prov.Name() != "mock1" {
		t.Errorf("Expected provider name to be 'mock1', got %s", prov.Name())
	}

	_, err = GetProvider("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent provider, got nil")
	}

	// Test ListProviders
	names := ListProviders()
	if len(names) != 2 {
		t.Errorf("Expected ListProviders to return 2 names, got %d", len(names))
	}

	// Test ListAvailableProviders
	available := ListAvailableProviders()
	if len(available) != 1 {
		t.Errorf("Expected 1 available provider, got %d", len(available))
	}
	if available[0].Name() != "mock1" {
		t.Errorf("Expected available provider to be 'mock1', got %s", available[0].Name())
	}
}

func TestNewProvider(t *testing.T) {
	tests := []struct {
		name        string
		config      ProviderConfig
		expectError bool
	}{
		{
			name: "Valid HTTP provider",
			config: ProviderConfig{
				Name:   "test-http",
				Type:   "http",
				Config: map[string]any{"endpoint": "https://api.test.com"},
			},
			expectError: false,
		},
		{
			name: "Valid OpenAI provider",
			config: ProviderConfig{
				Name:   "test-openai",
				Type:   "openai",
				Config: map[string]any{"endpoint": "https://api.openai.com/v1/chat/completions"},
			},
			expectError: false,
		},
		{
			name: "Valid Ollama provider",
			config: ProviderConfig{
				Name:   "test-ollama",
				Type:   "ollama",
				Config: map[string]any{"endpoint": "http://localhost:11434"},
			},
			expectError: false,
		},
		{
			name: "Valid Anthropic provider",
			config: ProviderConfig{
				Name:   "test-anthropic",
				Type:   "anthropic",
				Config: map[string]any{"endpoint": "https://api.anthropic.com/v1/messages"},
			},
			expectError: false,
		},
		{
			name: "Valid Executable provider",
			config: ProviderConfig{
				Name:   "test-exec",
				Type:   "executable",
				Config: map[string]any{"executable": "/bin/echo"},
			},
			expectError: false,
		},
		{
			name: "Invalid provider type",
			config: ProviderConfig{
				Name:   "test-invalid",
				Type:   "invalid",
				Config: map[string]any{},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := NewProvider(tt.config)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, got nil", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for %s, got %v", tt.name, err)
				}
				if provider == nil {
					t.Errorf("Expected provider for %s, got nil", tt.name)
				}
				if provider != nil && provider.Name() != tt.config.Name {
					t.Errorf("Expected provider name %s, got %s", tt.config.Name, provider.Name())
				}
			}
		})
	}
}

func TestLoadProviders(t *testing.T) {
	// Clear registry for testing
	providerRegistry = make(map[string]Provider)

	configs := []ProviderConfig{
		{
			Name:   "test1",
			Type:   "http",
			Config: map[string]any{"endpoint": "https://api.test1.com"},
		},
		{
			Name:   "test2",
			Type:   "ollama",
			Config: map[string]any{"endpoint": "http://localhost:11434"},
		},
	}

	err := LoadProviders(configs)
	if err != nil {
		t.Errorf("Expected no error loading providers, got %v", err)
	}

	if len(providerRegistry) != 2 {
		t.Errorf("Expected 2 providers in registry, got %d", len(providerRegistry))
	}

	// Verify providers are registered
	_, err = GetProvider("test1")
	if err != nil {
		t.Errorf("Expected to find test1 provider, got error: %v", err)
	}

	_, err = GetProvider("test2")
	if err != nil {
		t.Errorf("Expected to find test2 provider, got error: %v", err)
	}
}

func TestLoadProvidersWithError(t *testing.T) {
	configs := []ProviderConfig{
		{
			Name:   "invalid",
			Type:   "nonexistent",
			Config: map[string]any{},
		},
	}

	err := LoadProviders(configs)
	if err == nil {
		t.Error("Expected error loading invalid provider, got nil")
	}
}

func TestGetConfigString(t *testing.T) {
	tests := []struct {
		name       string
		config     map[string]any
		key        string
		defaultVal string
		expected   string
	}{
		{
			name:       "key exists",
			config:     map[string]any{"key": "value"},
			key:        "key",
			defaultVal: "default",
			expected:   "value",
		},
		{
			name:       "key missing",
			config:     map[string]any{},
			key:        "missing",
			defaultVal: "default",
			expected:   "default",
		},
		{
			name:       "key wrong type",
			config:     map[string]any{"key": 123},
			key:        "key",
			defaultVal: "default",
			expected:   "default",
		},
		{
			name:       "nil config",
			config:     nil,
			key:        "key",
			defaultVal: "default",
			expected:   "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getConfigString(tt.config, tt.key, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetConfigInt(t *testing.T) {
	tests := []struct {
		name       string
		config     map[string]any
		key        string
		defaultVal int
		expected   int
	}{
		{
			name:       "int key exists",
			config:     map[string]any{"key": 42},
			key:        "key",
			defaultVal: 0,
			expected:   42,
		},
		{
			name:       "float64 key exists",
			config:     map[string]any{"key": 42.0},
			key:        "key",
			defaultVal: 0,
			expected:   42,
		},
		{
			name:       "key missing",
			config:     map[string]any{},
			key:        "missing",
			defaultVal: 100,
			expected:   100,
		},
		{
			name:       "key wrong type",
			config:     map[string]any{"key": "not-an-int"},
			key:        "key",
			defaultVal: 50,
			expected:   50,
		},
		{
			name:       "nil config",
			config:     nil,
			key:        "key",
			defaultVal: 25,
			expected:   25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getConfigInt(tt.config, tt.key, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestNewHTTPProviderDefaults(t *testing.T) {
	provider, err := NewHTTPProvider("test", map[string]any{})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if provider == nil {
		t.Fatal("Expected provider to not be nil")
	}

	if provider.Name() != "test" {
		t.Errorf("Expected name to be 'test', got %s", provider.Name())
	}

	// Check that default endpoint is set
	if provider.endpoint != "https://api.openai.com/v1/chat/completions" {
		t.Errorf("Expected default endpoint, got %s", provider.endpoint)
	}

	// Check that default model is set
	if provider.model != "gpt-4o-mini" {
		t.Errorf("Expected default model gpt-4o-mini, got %s", provider.model)
	}

	// Check that default max tokens is set
	if provider.maxTokens != 4096 {
		t.Errorf("Expected default max tokens 4096, got %d", provider.maxTokens)
	}
}

func TestNewHTTPProviderCustomConfig(t *testing.T) {
	config := map[string]any{
		"endpoint":   "https://custom.api.com/v1/chat",
		"api_key":    "test-key",
		"model":      "custom-model",
		"max_tokens": 2048,
	}

	provider, err := NewHTTPProvider("custom", config)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if provider == nil {
		t.Fatal("Expected provider to not be nil")
	}

	if provider.endpoint != "https://custom.api.com/v1/chat" {
		t.Errorf("Expected custom endpoint, got %s", provider.endpoint)
	}
	if provider.model != "custom-model" {
		t.Errorf("Expected custom model, got %s", provider.model)
	}
	if provider.maxTokens != 2048 {
		t.Errorf("Expected custom max tokens, got %d", provider.maxTokens)
	}
	if provider.apiKey != "test-key" {
		t.Errorf("Expected api key to be set, got %s", provider.apiKey)
	}
}

func TestNewOllamaProvider(t *testing.T) {
	provider, err := NewOllamaProvider("test-ollama", map[string]any{
		"endpoint": "http://localhost:11434",
		"model":    "llama3",
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if provider == nil {
		t.Fatal("Expected provider to not be nil")
	}

	if provider.Name() != "test-ollama" {
		t.Errorf("Expected name to be 'test-ollama', got %s", provider.Name())
	}
}

func TestNewExecutableProvider(t *testing.T) {
	provider, err := NewExecutableProvider("test-exec", map[string]any{
		"executable": "/bin/echo",
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if provider == nil {
		t.Fatal("Expected provider to not be nil")
	}

	if provider.Name() != "test-exec" {
		t.Errorf("Expected name to be 'test-exec', got %s", provider.Name())
	}
}

func TestNewExecutableProviderMissingExecutable(t *testing.T) {
	provider, err := NewExecutableProvider("test-exec", map[string]any{})
	if err == nil {
		t.Error("Expected error for missing executable, got nil")
	}
	if provider != nil {
		t.Error("Expected provider to be nil on error")
	}
}
