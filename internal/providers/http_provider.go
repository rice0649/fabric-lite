package providers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// HTTPProvider implements Provider for OpenAI-compatible APIs
type HTTPProvider struct {
	name      string
	endpoint  string
	apiKey    string
	model     string
	headers   map[string]string
	client    *http.Client
	maxTokens int
}

// OpenAI API request/response structures
type openAIRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
	Stream      bool            `json:"stream,omitempty"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// NewHTTPProvider creates a new HTTP-based provider for OpenAI-compatible APIs
func NewHTTPProvider(name string, config map[string]any) (*HTTPProvider, error) {
	endpoint := getConfigString(config, "endpoint", "https://api.openai.com/v1/chat/completions")
	apiKeyEnv := getConfigString(config, "api_key_env", "OPENAI_API_KEY")
	apiKey := getConfigString(config, "api_key", "")
	model := getConfigString(config, "model", "gpt-4o-mini")
	maxTokens := getConfigInt(config, "max_tokens", 4096)

	// Resolve API key from environment if not set directly
	if apiKey == "" {
		apiKey = os.Getenv(apiKeyEnv)
	}

	return &HTTPProvider{
		name:      name,
		endpoint:  endpoint,
		apiKey:    apiKey,
		model:     model,
		maxTokens: maxTokens,
		headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + apiKey,
		},
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}, nil
}

func (p *HTTPProvider) Name() string {
	return p.name
}

func (p *HTTPProvider) IsAvailable() bool {
	return p.apiKey != "" && p.endpoint != ""
}

func (p *HTTPProvider) GetModels() []string {
	return []string{"gpt-4o", "gpt-4o-mini", "gpt-4-turbo", "gpt-3.5-turbo"}
}

func (p *HTTPProvider) Execute(ctx context.Context, request CompletionRequest) (*CompletionResponse, error) {
	if !p.IsAvailable() {
		return nil, fmt.Errorf("provider %s is not available (missing API key or endpoint)", p.name)
	}

	start := time.Now()

	// Build OpenAI request
	model := request.Model
	if model == "" {
		model = p.model
	}

	maxTokens := request.MaxTokens
	if maxTokens == 0 {
		maxTokens = p.maxTokens
	}

	messages := []openAIMessage{}
	if request.System != "" {
		messages = append(messages, openAIMessage{Role: "system", Content: request.System})
	}
	messages = append(messages, openAIMessage{Role: "user", Content: request.Prompt})

	oaiReq := openAIRequest{
		Model:     model,
		Messages:  messages,
		MaxTokens: maxTokens,
		Stream:    false,
	}

	// Make HTTP request
	jsonData, err := json.Marshal(oaiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for k, v := range p.headers {
		req.Header.Set(k, v)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var oaiResp openAIResponse
	if err := json.Unmarshal(body, &oaiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if oaiResp.Error != nil {
		return nil, fmt.Errorf("API error: %s", oaiResp.Error.Message)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	if len(oaiResp.Choices) == 0 {
		return nil, fmt.Errorf("no response choices returned")
	}

	return &CompletionResponse{
		Content:  oaiResp.Choices[0].Message.Content,
		Model:    oaiResp.Model,
		Tokens:   oaiResp.Usage.TotalTokens,
		Duration: time.Since(start),
	}, nil
}

func (p *HTTPProvider) ExecuteStream(ctx context.Context, request CompletionRequest) (<-chan StreamChunk, error) {
	if !p.IsAvailable() {
		return nil, fmt.Errorf("provider %s is not available", p.name)
	}

	chunks := make(chan StreamChunk, 100)

	go func() {
		defer close(chunks)

		model := request.Model
		if model == "" {
			model = p.model
		}

		maxTokens := request.MaxTokens
		if maxTokens == 0 {
			maxTokens = p.maxTokens
		}

		messages := []openAIMessage{}
		if request.System != "" {
			messages = append(messages, openAIMessage{Role: "system", Content: request.System})
		}
		messages = append(messages, openAIMessage{Role: "user", Content: request.Prompt})

		oaiReq := openAIRequest{
			Model:     model,
			Messages:  messages,
			MaxTokens: maxTokens,
			Stream:    true,
		}

		jsonData, err := json.Marshal(oaiReq)
		if err != nil {
			chunks <- StreamChunk{Error: err, Done: true}
			return
		}

		req, err := http.NewRequestWithContext(ctx, "POST", p.endpoint, bytes.NewBuffer(jsonData))
		if err != nil {
			chunks <- StreamChunk{Error: err, Done: true}
			return
		}

		for k, v := range p.headers {
			req.Header.Set(k, v)
		}

		resp, err := p.client.Do(req)
		if err != nil {
			chunks <- StreamChunk{Error: err, Done: true}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			chunks <- StreamChunk{Error: fmt.Errorf("API error: %s", string(body)), Done: true}
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				chunks <- StreamChunk{Done: true}
				return
			}

			var oaiResp openAIResponse
			if err := json.Unmarshal([]byte(data), &oaiResp); err != nil {
				continue
			}

			if len(oaiResp.Choices) > 0 && oaiResp.Choices[0].Delta.Content != "" {
				chunks <- StreamChunk{Content: oaiResp.Choices[0].Delta.Content}
			}
		}

		if err := scanner.Err(); err != nil {
			chunks <- StreamChunk{Error: err, Done: true}
		}
	}()

	return chunks, nil
}

// Helper functions
func getConfigString(config map[string]any, key, defaultVal string) string {
	if val, ok := config[key].(string); ok {
		return val
	}
	return defaultVal
}

func getConfigInt(config map[string]any, key string, defaultVal int) int {
	if val, ok := config[key].(int); ok {
		return val
	}
	if val, ok := config[key].(float64); ok {
		return int(val)
	}
	return defaultVal
}
