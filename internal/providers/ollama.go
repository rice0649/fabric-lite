package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	ollamaDefaultURL = "http://localhost:11434/api/generate"
)

// OllamaRequest represents the request structure for the Ollama API.
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// OllamaResponse represents the response structure for the Ollama API.
type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// SendOllamaRequest sends a prompt to a local Ollama server and gets a response.
func SendOllamaRequest(prompt, model string) (string, error) {
	requestBody, err := json.Marshal(OllamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal ollama request: %w", err)
	}

	resp, err := http.Post(ollamaDefaultURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to send request to ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK status from ollama: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read ollama response body: %w", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal ollama response: %w", err)
	}

	return ollamaResp.Response, nil
}
