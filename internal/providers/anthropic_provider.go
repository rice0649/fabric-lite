package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	anthropicAPIEndpoint = "https://api.anthropic.com/v1/messages"
	anthropicAPIVersion  = "2023-06-01"
)

// AnthropicProvider implements Provider for Anthropic Claude API
type AnthropicProvider struct {
	name      string
	endpoint  string
	apiKey    string
	model     string
	maxTokens int
	client    *http.Client
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system,omitempty"`
	Messages  []anthropicMessage `json:"messages"`
	Stream    bool               `json:"stream,omitempty"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model      string `json:"model"`
	StopReason string `json:"stop_reason"`
	Usage      struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// NewAnthropicProvider creates a new Anthropic Claude provider
func NewAnthropicProvider(name string, config map[string]any) (*AnthropicProvider, error) {
	endpoint := getConfigString(config, "endpoint", anthropicAPIEndpoint)
	apiKeyEnv := getConfigString(config, "api_key_env", "ANTHROPIC_API_KEY")
	apiKey := getConfigString(config, "api_key", "")
	model := getConfigString(config, "model", "claude-sonnet-4-20250514")
	maxTokens := getConfigInt(config, "max_tokens", 4096)

	if apiKey == "" {
		apiKey = os.Getenv(apiKeyEnv)
	}

	return &AnthropicProvider{
		name:      name,
		endpoint:  endpoint,
		apiKey:    apiKey,
		model:     model,
		maxTokens: maxTokens,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}, nil
}

func (p *AnthropicProvider) Name() string {
	return p.name
}

func (p *AnthropicProvider) IsAvailable() bool {
	return p.apiKey != ""
}

func (p *AnthropicProvider) GetModels() []string {
	return []string{"claude-sonnet-4-20250514", "claude-3-5-sonnet-20241022", "claude-3-opus-20240229"}
}

func (p *AnthropicProvider) Execute(ctx context.Context, request CompletionRequest) (*CompletionResponse, error) {
	if !p.IsAvailable() {
		return nil, fmt.Errorf("provider %s is not available (missing API key)", p.name)
	}

	start := time.Now()

	model := request.Model
	if model == "" {
		model = p.model
	}

	maxTokens := request.MaxTokens
	if maxTokens == 0 {
		maxTokens = p.maxTokens
	}

	anthropicReq := anthropicRequest{
		Model:     model,
		MaxTokens: maxTokens,
		System:    request.System,
		Messages: []anthropicMessage{
			{Role: "user", Content: request.Prompt},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(anthropicReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", anthropicAPIVersion)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var anthropicResp anthropicResponse
	if err := json.Unmarshal(body, &anthropicResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if anthropicResp.Error != nil {
		return nil, fmt.Errorf("API error: %s - %s", anthropicResp.Error.Type, anthropicResp.Error.Message)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var content string
	for _, c := range anthropicResp.Content {
		if c.Type == "text" {
			content += c.Text
		}
	}

	return &CompletionResponse{
		Content:  content,
		Model:    anthropicResp.Model,
		Tokens:   anthropicResp.Usage.InputTokens + anthropicResp.Usage.OutputTokens,
		Duration: time.Since(start),
	}, nil
}

func (p *AnthropicProvider) ExecuteStream(ctx context.Context, request CompletionRequest) (<-chan StreamChunk, error) {
	// For now, use non-streaming and return all at once
	// Full SSE streaming can be added later
	chunks := make(chan StreamChunk, 1)

	go func() {
		defer close(chunks)
		resp, err := p.Execute(ctx, request)
		if err != nil {
			chunks <- StreamChunk{Error: err, Done: true}
			return
		}
		chunks <- StreamChunk{Content: resp.Content, Done: true}
	}()

	return chunks, nil
}
