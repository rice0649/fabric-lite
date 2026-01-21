package providers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestNewAnthropicProvider(t *testing.T) {
	tests := []struct {
		name        string
		config      map[string]any
		expectError bool
		expectAvail bool
	}{
		{
			name: "Valid config with API key",
			config: map[string]any{
				"endpoint":   "https://api.anthropic.com/v1/messages",
				"api_key":    "test-key",
				"model":      "claude-sonnet-4-20250514",
				"max_tokens": 4096,
			},
			expectError: false,
			expectAvail: true,
		},
		{
			name: "Valid config with API key env var",
			config: map[string]any{
				"endpoint":    "https://api.anthropic.com/v1/messages",
				"api_key_env": "NONEXISTENT_API_KEY", // Use non-existent env var
				"model":       "claude-sonnet-4-20250514",
				"max_tokens":  4096,
			},
			expectError: false,
			expectAvail: false, // No env var set
		},
		{
			name: "Missing API key",
			config: map[string]any{
				"endpoint": "https://api.anthropic.com/v1/messages",
				"model":    "claude-sonnet-4-20250514",
			},
			expectError: false,
			expectAvail: false,
		},
		{
			name: "Default values",
			config: map[string]any{
				"api_key": "test-key",
			},
			expectError: false,
			expectAvail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable if needed
			if apiKeyEnv, ok := tt.config["api_key_env"].(string); ok && apiKeyEnv == "ANTHROPIC_API_KEY" {
				os.Setenv("ANTHROPIC_API_KEY", "env-test-key")
				defer os.Unsetenv("ANTHROPIC_API_KEY")
			}

			provider, err := NewAnthropicProvider("test", tt.config)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, got nil", tt.name)
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error for %s, got %v", tt.name, err)
				return
			}

			if provider.Name() != "test" {
				t.Errorf("Expected provider name to be 'test', got %s", provider.Name())
			}

			if provider.IsAvailable() != tt.expectAvail {
				t.Errorf("Expected IsAvailable to be %v, got %v", tt.expectAvail, provider.IsAvailable())
			}

			models := provider.GetModels()
			if len(models) == 0 {
				t.Error("Expected at least one model")
			}
		})
	}
}

func TestAnthropicProviderExecute(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("x-api-key") != "test-key" {
			t.Errorf("Expected x-api-key to be 'test-key', got %s", r.Header.Get("x-api-key"))
		}
		if r.Header.Get("anthropic-version") != "2023-06-01" {
			t.Errorf("Expected anthropic-version to be '2023-06-01', got %s", r.Header.Get("anthropic-version"))
		}

		// Send mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "msg_123",
			"type": "message",
			"role": "assistant",
			"content": [
				{
					"type": "text",
					"text": "Test response"
				}
			],
			"model": "claude-sonnet-4-20250514",
			"stop_reason": "end_turn",
			"usage": {
				"input_tokens": 10,
				"output_tokens": 5
			}
		}`))
	}))
	defer server.Close()

	config := map[string]any{
		"endpoint":   server.URL,
		"api_key":    "test-key",
		"model":      "claude-sonnet-4-20250514",
		"max_tokens": 4096,
	}

	provider, err := NewAnthropicProvider("test", config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	ctx := context.Background()
	request := CompletionRequest{
		System:    "You are a helpful assistant",
		Prompt:    "Hello, world!",
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 1000,
	}

	response, err := provider.Execute(ctx, request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Content != "Test response" {
		t.Errorf("Expected content 'Test response', got %s", response.Content)
	}

	if response.Model != "claude-sonnet-4-20250514" {
		t.Errorf("Expected model 'claude-sonnet-4-20250514', got %s", response.Model)
	}

	if response.Tokens != 15 { // 10 input + 5 output
		t.Errorf("Expected tokens 15, got %d", response.Tokens)
	}

	if response.Duration <= 0 {
		t.Errorf("Expected positive duration, got %v", response.Duration)
	}
}

func TestAnthropicProviderExecuteError(t *testing.T) {
	// Create mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": {
				"type": "invalid_request_error",
				"message": "Invalid request"
			}
		}`))
	}))
	defer server.Close()

	config := map[string]any{
		"endpoint": server.URL,
		"api_key":  "test-key",
	}

	provider, err := NewAnthropicProvider("test", config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	ctx := context.Background()
	request := CompletionRequest{
		Prompt: "Hello, world!",
	}

	_, err = provider.Execute(ctx, request)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestAnthropicProviderExecuteUnavailable(t *testing.T) {
	config := map[string]any{
		"endpoint": "https://api.anthropic.com/v1/messages",
		// No API key
	}

	provider, err := NewAnthropicProvider("test", config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	ctx := context.Background()
	request := CompletionRequest{
		Prompt: "Hello, world!",
	}

	_, err = provider.Execute(ctx, request)
	if err == nil {
		t.Error("Expected error for unavailable provider, got nil")
	}
}

func TestAnthropicProviderExecuteStream(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "msg_123",
			"type": "message",
			"role": "assistant",
			"content": [
				{
					"type": "text",
					"text": "Test response"
				}
			],
			"model": "claude-sonnet-4-20250514",
			"stop_reason": "end_turn",
			"usage": {
				"input_tokens": 10,
				"output_tokens": 5
			}
		}`))
	}))
	defer server.Close()

	config := map[string]any{
		"endpoint": server.URL,
		"api_key":  "test-key",
	}

	provider, err := NewAnthropicProvider("test", config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	ctx := context.Background()
	request := CompletionRequest{
		Prompt: "Hello, world!",
		Stream: true,
	}

	chunks, err := provider.ExecuteStream(ctx, request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Collect chunks
	var collectedChunks []StreamChunk
	for chunk := range chunks {
		collectedChunks = append(collectedChunks, chunk)
	}

	if len(collectedChunks) != 1 {
		t.Errorf("Expected 1 chunk, got %d", len(collectedChunks))
	}

	if !collectedChunks[0].Done {
		t.Errorf("Expected chunk to be done, got %v", collectedChunks[0].Done)
	}

	if collectedChunks[0].Content != "Test response" {
		t.Errorf("Expected content 'Test response', got %s", collectedChunks[0].Content)
	}

	if collectedChunks[0].Error != nil {
		t.Errorf("Expected no error, got %v", collectedChunks[0].Error)
	}
}

func TestAnthropicProviderExecuteStreamError(t *testing.T) {
	config := map[string]any{
		"endpoint": "https://api.anthropic.com/v1/messages",
		// No API key - provider unavailable
	}

	provider, err := NewAnthropicProvider("test", config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	ctx := context.Background()
	request := CompletionRequest{
		Prompt: "Hello, world!",
		Stream: true,
	}

	chunks, err := provider.ExecuteStream(ctx, request)
	if err != nil {
		t.Errorf("Expected no error creating stream, got %v", err)
	}

	// Should get error chunk
	chunk := <-chunks
	if !chunk.Done {
		t.Errorf("Expected chunk to be done, got %v", chunk.Done)
	}

	if chunk.Error == nil {
		t.Error("Expected error in chunk, got nil")
	}
}

func TestAnthropicProviderRequestDefaults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read request body to verify defaults
		body := make([]byte, 1024)
		n, _ := r.Body.Read(body)
		bodyStr := string(body[:n])

		if !contains(bodyStr, "claude-sonnet-4-20250514") {
			t.Error("Expected default model to be used")
		}
		if !contains(bodyStr, "4096") {
			t.Error("Expected default max_tokens to be used")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "msg_123",
			"type": "message",
			"role": "assistant",
			"content": [{"type": "text", "text": "Test"}],
			"model": "claude-sonnet-4-20250514",
			"stop_reason": "end_turn",
			"usage": {"input_tokens": 1, "output_tokens": 1}
		}`))
	}))
	defer server.Close()

	config := map[string]any{
		"endpoint": server.URL,
		"api_key":  "test-key",
		// Using defaults for model and max_tokens
	}

	provider, err := NewAnthropicProvider("test", config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	ctx := context.Background()
	request := CompletionRequest{
		Prompt: "Hello, world!",
		// No model or max_tokens specified
	}

	_, err = provider.Execute(ctx, request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestAnthropicProviderContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"content": [{"type": "text", "text": "response"}]}`))
	}))
	defer server.Close()

	config := map[string]any{
		"endpoint": server.URL,
		"api_key":  "test-key",
	}

	provider, err := NewAnthropicProvider("test", config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// Cancel context immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	request := CompletionRequest{
		Prompt: "Hello, world!",
	}

	_, err = provider.Execute(ctx, request)
	if err == nil {
		t.Error("Expected context cancellation error, got nil")
	}
}
