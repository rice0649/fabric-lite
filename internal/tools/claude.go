package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	claudeAPIEndpoint   = "https://api.anthropic.com/v1/messages"
	claudeAPIVersion    = "2023-06-01"
	defaultClaudeModel  = "claude-sonnet-4-20250514"
	defaultMaxTokens    = 4096
)

// ClaudeTool provides direct access to Claude/Anthropic API
type ClaudeTool struct {
	BaseTool
	apiKey    string
	model     string
	maxTokens int
	client    *http.Client
}

// ClaudeRequest represents a request to the Anthropic API
type ClaudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []ClaudeMessage `json:"messages"`
	System    string          `json:"system,omitempty"`
}

// ClaudeMessage represents a message in the conversation
type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ClaudeResponse represents a response from the Anthropic API
type ClaudeResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model        string `json:"model"`
	StopReason   string `json:"stop_reason"`
	StopSequence string `json:"stop_sequence,omitempty"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// NewClaudeTool creates a new Claude API tool
func NewClaudeTool() *ClaudeTool {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")

	return &ClaudeTool{
		BaseTool: BaseTool{
			name:        "claude",
			description: "Anthropic Claude API - Advanced reasoning and code generation",
			command:     "claude-api", // Not a CLI command, but satisfies interface
		},
		apiKey:    apiKey,
		model:     defaultClaudeModel,
		maxTokens: defaultMaxTokens,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// NewClaudeToolWithConfig creates a Claude tool with specific configuration
func NewClaudeToolWithConfig(model string, apiKey string, apiKeyEnv string, maxTokens int) *ClaudeTool {
	// Resolve API key
	key := apiKey
	if key == "" && apiKeyEnv != "" {
		key = os.Getenv(apiKeyEnv)
	}
	if key == "" {
		key = os.Getenv("ANTHROPIC_API_KEY")
	}

	if model == "" {
		model = defaultClaudeModel
	}
	if maxTokens == 0 {
		maxTokens = defaultMaxTokens
	}

	return &ClaudeTool{
		BaseTool: BaseTool{
			name:        "claude",
			description: "Anthropic Claude API - Advanced reasoning and code generation",
			command:     "claude-api",
		},
		apiKey:    key,
		model:     model,
		maxTokens: maxTokens,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// IsAvailable checks if the Claude API is accessible (has API key)
func (t *ClaudeTool) IsAvailable() bool {
	return t.apiKey != ""
}

// Execute runs a prompt through the Claude API
func (t *ClaudeTool) Execute(ctx ExecutionContext) (*ExecutionResult, error) {
	if t.apiKey == "" {
		return &ExecutionResult{
			Success: false,
			Error:   "ANTHROPIC_API_KEY not set",
		}, fmt.Errorf("ANTHROPIC_API_KEY environment variable not set")
	}

	prompt := ctx.Prompt
	if prompt == "" {
		prompt = t.getPhasePrompt(ctx.Phase)
	}

	// Build request
	request := ClaudeRequest{
		Model:     t.model,
		MaxTokens: t.maxTokens,
		Messages: []ClaudeMessage{
			{Role: "user", Content: prompt},
		},
	}

	// Add system prompt based on phase
	if systemPrompt := t.getSystemPrompt(ctx.Phase); systemPrompt != "" {
		request.System = systemPrompt
	}

	// Execute API call
	response, err := t.callAPI(request)
	if err != nil {
		return &ExecutionResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	// Extract text from response
	var output string
	for _, content := range response.Content {
		if content.Type == "text" {
			output += content.Text
		}
	}

	// Print output to stdout for interactive use
	fmt.Println(output)

	return &ExecutionResult{
		Output:   output,
		Success:  true,
		ExitCode: 0,
	}, nil
}

// ExecuteNonInteractive runs Claude API and returns output without printing
func (t *ClaudeTool) ExecuteNonInteractive(ctx ExecutionContext) (*ExecutionResult, error) {
	if t.apiKey == "" {
		return &ExecutionResult{
			Success: false,
			Error:   "ANTHROPIC_API_KEY not set",
		}, fmt.Errorf("ANTHROPIC_API_KEY environment variable not set")
	}

	prompt := ctx.Prompt
	if prompt == "" {
		prompt = t.getPhasePrompt(ctx.Phase)
	}

	request := ClaudeRequest{
		Model:     t.model,
		MaxTokens: t.maxTokens,
		Messages: []ClaudeMessage{
			{Role: "user", Content: prompt},
		},
	}

	if systemPrompt := t.getSystemPrompt(ctx.Phase); systemPrompt != "" {
		request.System = systemPrompt
	}

	response, err := t.callAPI(request)
	if err != nil {
		return &ExecutionResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	var output string
	for _, content := range response.Content {
		if content.Type == "text" {
			output += content.Text
		}
	}

	return &ExecutionResult{
		Output:   output,
		Success:  true,
		ExitCode: 0,
	}, nil
}

// callAPI makes the actual HTTP request to the Anthropic API
func (t *ClaudeTool) callAPI(request ClaudeRequest) (*ClaudeResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", claudeAPIEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", t.apiKey)
	req.Header.Set("anthropic-version", claudeAPIVersion)

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response ClaudeResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("API error: %s - %s", response.Error.Type, response.Error.Message)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return &response, nil
}

func (t *ClaudeTool) getSystemPrompt(phase string) string {
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

func (t *ClaudeTool) getPhasePrompt(phase string) string {
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
