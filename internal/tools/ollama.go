package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

// OllamaTool represents the Ollama CLI tool
type OllamaTool struct {
	endpoint string
}

func NewOllamaTool(endpoint string) *OllamaTool {
	return &OllamaTool{endpoint: endpoint}
}

func (t *OllamaTool) Name() string {
	return "ollama"
}

func (t *OllamaTool) Description() string {
	return "Interact with the Ollama CLI to pull and list models. Requires Ollama to be installed and running."
}

func (t *OllamaTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": map[string]any{
				"type":        "string",
				"description": "The ollama command to execute (e.g., 'pull llama3', 'list').",
				"enum":        []string{"pull", "list"},
			},
			"model": map[string]any{
				"type":        "string",
				"description": "The model name, required for 'pull' command.",
			},
		},
		"required": []string{"command"},
	}
}

func (t *OllamaTool) GetCommand() string {
	return "ollama"
}

func (t *OllamaTool) Execute(ctx ExecutionContext) (*ExecutionResult, error) {
	// The input for OllamaTool comes from ctx.Args (command and model)
	if len(ctx.Args) == 0 {
		return nil, fmt.Errorf("missing command for ollama tool")
	}

	command := ctx.Args[0]
	input := make(map[string]any)
	input["command"] = command

	if command == "pull" {
		if len(ctx.Args) < 2 {
			return nil, fmt.Errorf("'model' is required for 'pull' command")
		}
		input["model"] = ctx.Args[1]
	}

	var output string
	var err error

	switch command {
	case "pull":
		model := input["model"].(string)
		output, err = t.pullModel(context.Background(), model)
	case "list":
		output, err = t.listModels(context.Background())
	default:
		return nil, fmt.Errorf("unsupported ollama command: %s", command)
	}

	if err != nil {
		return &ExecutionResult{
			Output:   "",
			Error:    err.Error(),
			ExitCode: 1,
			Success:  false,
		}, nil
	}

	return &ExecutionResult{
		Output:   output,
		Error:    "",
		ExitCode: 0,
		Success:  true,
	}, nil
}

func (t *OllamaTool) pullModel(ctx context.Context, model string) (string, error) {
	cmd := exec.CommandContext(ctx, "ollama", "pull", model)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to pull model '%s': %v, output: %s", model, err, string(output))
	}
	return string(output), nil
}

func (t *OllamaTool) listModels(ctx context.Context) (string, error) {
	// Directly query the Ollama API for models as it's more reliable than parsing CLI output
	resp, err := http.Get(t.endpoint + "/api/tags")
	if err != nil {
		return "", fmt.Errorf("failed to connect to Ollama API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode Ollama API response: %w", err)
	}

	var models []string
	for _, m := range result.Models {
		models = append(models, m.Name)
	}

	return strings.Join(models, "\n"), nil
}

func (t *OllamaTool) IsAvailable() bool {
	// Check if ollama is in PATH
	_, err := exec.LookPath("ollama")
	return err == nil
}