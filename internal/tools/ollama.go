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
	defaultOllamaEndpoint = "http://localhost:11434"
	defaultOllamaModel    = "llama3.2"
)

// OllamaTool provides access to local Ollama LLM instance
type OllamaTool struct {
	BaseTool
	endpoint string
	model    string
	client   *http.Client
}

// OllamaGenerateRequest represents a request to the Ollama generate API
type OllamaGenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	System string `json:"system,omitempty"`
	Stream bool   `json:"stream"`
}

// OllamaGenerateResponse represents a response from Ollama generate API
type OllamaGenerateResponse struct {
	Model     string `json:"model"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"created_at"`
	Context   []int  `json:"context,omitempty"`
	Error     string `json:"error,omitempty"`
}

// OllamaChatRequest represents a request to the Ollama chat API
type OllamaChatRequest struct {
	Model    string              `json:"model"`
	Messages []OllamaChatMessage `json:"messages"`
	Stream   bool                `json:"stream"`
}

// OllamaChatMessage represents a message in the Ollama chat
type OllamaChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OllamaChatResponse represents a response from Ollama chat API
type OllamaChatResponse struct {
	Model   string `json:"model"`
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done  bool   `json:"done"`
	Error string `json:"error,omitempty"`
}

// NewOllamaTool creates a new Ollama tool with default settings
func NewOllamaTool() *OllamaTool {
	return &OllamaTool{
		BaseTool: BaseTool{
			name:        "ollama",
			description: "Ollama - Local LLM inference (Llama, Mistral, CodeLlama, etc.)",
			command:     "ollama",
		},
		endpoint: defaultOllamaEndpoint,
		model:    defaultOllamaModel,
		client: &http.Client{
			Timeout: 300 * time.Second, // Local inference can be slow
		},
	}
}

// NewOllamaToolWithConfig creates an Ollama tool with specific configuration
func NewOllamaToolWithConfig(model string, endpoint string) *OllamaTool {
	if endpoint == "" {
		endpoint = defaultOllamaEndpoint
	}
	if model == "" {
		model = defaultOllamaModel
	}

	return &OllamaTool{
		BaseTool: BaseTool{
			name:        "ollama",
			description: "Ollama - Local LLM inference (Llama, Mistral, CodeLlama, etc.)",
			command:     "ollama",
		},
		endpoint: endpoint,
		model:    model,
		client: &http.Client{
			Timeout: 300 * time.Second,
		},
	}
}

// IsAvailable checks if Ollama is running and accessible
func (t *OllamaTool) IsAvailable() bool {
	// Check if Ollama API is accessible
	resp, err := t.client.Get(t.endpoint + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// Execute runs a prompt through the local Ollama instance
func (t *OllamaTool) Execute(ctx ExecutionContext) (*ExecutionResult, error) {
	prompt := ctx.Prompt
	if prompt == "" {
		prompt = t.getPhasePrompt(ctx.Phase)
	}

	// Use chat API for better results
	request := OllamaChatRequest{
		Model: t.model,
		Messages: []OllamaChatMessage{
			{Role: "system", Content: t.getSystemPrompt(ctx.Phase)},
			{Role: "user", Content: prompt},
		},
		Stream: false,
	}

	response, err := t.callChatAPI(request)
	if err != nil {
		return &ExecutionResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	output := response.Message.Content

	// Print output for interactive use
	fmt.Println(output)

	return &ExecutionResult{
		Output:   output,
		Success:  true,
		ExitCode: 0,
	}, nil
}

// ExecuteNonInteractive runs Ollama and returns output without printing
func (t *OllamaTool) ExecuteNonInteractive(ctx ExecutionContext) (*ExecutionResult, error) {
	prompt := ctx.Prompt
	if prompt == "" {
		prompt = t.getPhasePrompt(ctx.Phase)
	}

	request := OllamaChatRequest{
		Model: t.model,
		Messages: []OllamaChatMessage{
			{Role: "system", Content: t.getSystemPrompt(ctx.Phase)},
			{Role: "user", Content: prompt},
		},
		Stream: false,
	}

	response, err := t.callChatAPI(request)
	if err != nil {
		return &ExecutionResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &ExecutionResult{
		Output:   response.Message.Content,
		Success:  true,
		ExitCode: 0,
	}, nil
}

// callChatAPI makes the HTTP request to Ollama chat API
func (t *OllamaTool) callChatAPI(request OllamaChatRequest) (*OllamaChatResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := t.endpoint + "/api/chat"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Ollama request failed (is Ollama running?): %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	var response OllamaChatResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if response.Error != "" {
		return nil, fmt.Errorf("Ollama error: %s", response.Error)
	}

	return &response, nil
}

// callGenerateAPI makes the HTTP request to Ollama generate API (alternative)
func (t *OllamaTool) callGenerateAPI(request OllamaGenerateRequest) (*OllamaGenerateResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := t.endpoint + "/api/generate"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Ollama request failed (is Ollama running?): %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	var response OllamaGenerateResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if response.Error != "" {
		return nil, fmt.Errorf("Ollama error: %s", response.Error)
	}

	return &response, nil
}

func (t *OllamaTool) getSystemPrompt(phase string) string {
	prompts := map[string]string{
		"discovery": `You are an expert software architect helping with project discovery.
Analyze requirements, research best practices, and identify technical constraints.
Be thorough but concise. Focus on actionable insights.`,

		"planning": `You are an expert software architect helping with project planning.
Create clear architecture documents, component breakdowns, and technical decisions.
Consider scalability, maintainability, and best practices.`,

		"design": `You are an expert software designer helping with API and data model design.
Create clean, well-documented interfaces and data structures.
Follow industry standards and conventions.`,

		"implementation": `You are an expert software developer helping with code implementation.
Write clean, maintainable, well-tested code.
Follow best practices and coding standards.`,

		"testing": `You are an expert QA engineer helping with test creation.
Design comprehensive test suites covering edge cases.
Focus on both unit and integration tests.`,

		"deployment": `You are a DevOps expert helping with deployment preparation.
Create clear documentation, changelogs, and deployment guides.
Consider security, monitoring, and rollback procedures.`,
	}

	if p, ok := prompts[phase]; ok {
		return p
	}
	return "You are an expert software engineer. Be helpful, accurate, and concise."
}

func (t *OllamaTool) getPhasePrompt(phase string) string {
	prompts := map[string]string{
		"discovery": "Analyze this project and help gather requirements. What are the key technical decisions needed?",
		"planning":  "Help create an architecture plan for this project. What components are needed?",
		"design":    "Help design the API and data models for this project.",
		"implementation": "Help implement the core functionality for this project.",
		"testing":   "Help create a comprehensive test plan for this project.",
		"deployment": "Help prepare deployment documentation for this project.",
	}

	if p, ok := prompts[phase]; ok {
		return p
	}
	return ""
}

// ListModels returns available models from the Ollama instance
func (t *OllamaTool) ListModels() ([]string, error) {
	resp, err := t.client.Get(t.endpoint + "/api/tags")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse models: %w", err)
	}

	models := make([]string, len(result.Models))
	for i, m := range result.Models {
		models[i] = m.Name
	}
	return models, nil
}
