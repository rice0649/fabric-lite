package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	ollamaDefaultURL     = "http://localhost:11434/api/chat"
	ollamaRequestTimeout = 120 * time.Second
	ollamaHealthURL      = "http://localhost:11434/api/tags"
	ollamaDefaultModel   = "llama3:latest"
)

// OllamaTool is a tool for interacting with a local Ollama server.
type OllamaTool struct {
	BaseTool
}

type ollamaChatResponse struct {
	Model   string `json:"model"`
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done  bool   `json:"done"`
	Error string `json:"error,omitempty"`
}

// NewOllamaTool creates a new OllamaTool.
func NewOllamaTool() *OllamaTool {
	return &OllamaTool{
		BaseTool: BaseTool{
			name:        "ollama",
			description: "The Quick Task Automator. Runs local models like Llama3 for fast, simple tasks like boilerplate generation and formatting.",
			command:     "", // Not an external command
		},
	}
}

// IsAvailable checks if the Ollama server is running.
func (t *OllamaTool) IsAvailable() bool {
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get(ollamaHealthURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// Execute sends the prompt to the Ollama server.
func (t *OllamaTool) Execute(ctx ExecutionContext) (*ExecutionResult, error) {
	// Ollama is non-interactive, so we capture the output.
	// Default to llama3 if no other model is specified in args
	model := ollamaDefaultModel
	if len(ctx.Args) > 0 {
		model = ctx.Args[0]
	}

	// Create request body matching Ollama API format
	requestBody, err := json.Marshal(map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": ctx.Prompt},
		},
		"stream": false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ollama request: %w", err)
	}

	client := http.Client{
		Timeout: ollamaRequestTimeout,
	}

	resp, err := client.Post(ollamaDefaultURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return &ExecutionResult{Success: false, Error: err.Error()}, fmt.Errorf("failed to send request to ollama: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ExecutionResult{Success: false, Error: err.Error()}, fmt.Errorf("failed to read ollama response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		errorMsg := fmt.Sprintf("received non-OK status from ollama: %s. Body: %s", resp.Status, string(body))
		return &ExecutionResult{Success: false, Error: errorMsg, ExitCode: resp.StatusCode}, nil
	}

	var ollamaResp ollamaChatResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return &ExecutionResult{Success: false, Error: err.Error()}, fmt.Errorf("failed to unmarshal ollama response: %w", err)
	}

	output := ollamaResp.Message.Content
	return &ExecutionResult{
		Output:   output,
		Success:  true,
		ExitCode: 0,
	}, nil
}
