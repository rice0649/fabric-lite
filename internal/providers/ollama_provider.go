package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultOllamaEndpoint = "http://localhost:11434"

// OllamaProvider implements Provider for local Ollama instance
type OllamaProvider struct {
	name     string
	endpoint string
	model    string
	client   *http.Client
}

type ollamaChatRequest struct {
	Model    string              `json:"model"`
	Messages []ollamaChatMessage `json:"messages"`
	Stream   bool                `json:"stream"`
}

type ollamaChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
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

// NewOllamaProvider creates a new Ollama provider
func NewOllamaProvider(name string, config map[string]any) (*OllamaProvider, error) {
	endpoint := getConfigString(config, "endpoint", defaultOllamaEndpoint)
	model := getConfigString(config, "model", "llama3.2")

	return &OllamaProvider{
		name:     name,
		endpoint: endpoint,
		model:    model,
		client: &http.Client{
			Timeout: 300 * time.Second, // Ollama can be slow
		},
	}, nil
}

func (p *OllamaProvider) Name() string {
	return p.name
}

func (p *OllamaProvider) IsAvailable() bool {
	// Check if Ollama is running
	resp, err := p.client.Get(p.endpoint + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (p *OllamaProvider) GetModels() []string {
	models, err := p.listModels()
	if err != nil {
		return []string{"llama3.2", "mistral", "codellama"}
	}
	return models
}

func (p *OllamaProvider) listModels() ([]string, error) {
	resp, err := p.client.Get(p.endpoint + "/api/tags")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	models := make([]string, len(result.Models))
	for i, m := range result.Models {
		models[i] = m.Name
	}
	return models, nil
}

func (p *OllamaProvider) Execute(ctx context.Context, request CompletionRequest) (*CompletionResponse, error) {
	start := time.Now()

	model := request.Model
	if model == "" {
		model = p.model
	}

	messages := []ollamaChatMessage{}
	if request.System != "" {
		messages = append(messages, ollamaChatMessage{Role: "system", Content: request.System})
	}
	messages = append(messages, ollamaChatMessage{Role: "user", Content: request.Prompt})

	ollamaReq := ollamaChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
	}

	jsonData, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := p.endpoint + "/api/chat"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed (is Ollama running?): %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	var ollamaResp ollamaChatResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if ollamaResp.Error != "" {
		return nil, fmt.Errorf("Ollama error: %s", ollamaResp.Error)
	}

	return &CompletionResponse{
		Content:  ollamaResp.Message.Content,
		Model:    ollamaResp.Model,
		Duration: time.Since(start),
	}, nil
}

func (p *OllamaProvider) ExecuteStream(ctx context.Context, request CompletionRequest) (<-chan StreamChunk, error) {
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
