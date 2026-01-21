package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewProjectConfig(t *testing.T) {
	config := NewProjectConfig("test-project", "default")

	if config.Name != "test-project" {
		t.Errorf("Expected name 'test-project', got '%s'", config.Name)
	}
	if config.Template != "default" {
		t.Errorf("Expected template 'default', got '%s'", config.Template)
	}
	if config.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", config.Version)
	}

	// Check tools configuration
	if !config.Tools.Gemini.Enabled {
		t.Error("Expected Gemini to be enabled")
	}
	if config.Tools.Gemini.Model != "gemini-2.0-flash-exp" {
		t.Errorf("Expected Gemini model 'gemini-2.0-flash-exp', got '%s'", config.Tools.Gemini.Model)
	}

	if !config.Tools.Codex.Enabled {
		t.Error("Expected Codex to be enabled")
	}
	if config.Tools.Codex.Model != "o3-mini" {
		t.Errorf("Expected Codex model 'o3-mini', got '%s'", config.Tools.Codex.Model)
	}

	if !config.Tools.OpenCode.Enabled {
		t.Error("Expected OpenCode to be enabled")
	}
	if config.Tools.OpenCode.Provider != "anthropic" {
		t.Errorf("Expected OpenCode provider 'anthropic', got '%s'", config.Tools.OpenCode.Provider)
	}

	if !config.Tools.Fabric.Enabled {
		t.Error("Expected Fabric to be enabled")
	}
	if config.Tools.Fabric.Model != "gpt-4o-mini" {
		t.Errorf("Expected Fabric model 'gpt-4o-mini', got '%s'", config.Tools.Fabric.Model)
	}

	if config.Tools.Claude.Enabled {
		t.Error("Expected Claude to be disabled by default")
	}

	if config.Tools.Ollama.Enabled {
		t.Error("Expected Ollama to be disabled by default")
	}
}

func TestProjectConfigSaveAndLoad(t *testing.T) {
	// Create test config
	originalConfig := NewProjectConfig("test-project", "default")
	originalConfig.Description = "Test project description"
	originalConfig.Tools.Claude.Enabled = true
	originalConfig.Tools.Ollama.Enabled = true
	originalConfig.Phases = map[string]string{
		"discovery": "gemini",
		"design":    "opencode",
	}

	// Save to temporary file
	tempFile := filepath.Join(t.TempDir(), "test_config.yaml")
	err := originalConfig.Save(tempFile)
	if err != nil {
		t.Errorf("Expected no error saving config, got %v", err)
	}

	// Load from file
	loadedConfig, err := LoadProjectConfig(tempFile)
	if err != nil {
		t.Errorf("Expected no error loading config, got %v", err)
	}

	// Verify loaded config
	if loadedConfig.Name != originalConfig.Name {
		t.Errorf("Expected name '%s', got '%s'", originalConfig.Name, loadedConfig.Name)
	}
	if loadedConfig.Description != originalConfig.Description {
		t.Errorf("Expected description '%s', got '%s'", originalConfig.Description, loadedConfig.Description)
	}
	if loadedConfig.Version != originalConfig.Version {
		t.Errorf("Expected version '%s', got '%s'", originalConfig.Version, loadedConfig.Version)
	}

	if !loadedConfig.Tools.Claude.Enabled {
		t.Error("Expected Claude to be enabled in loaded config")
	}
	if !loadedConfig.Tools.Ollama.Enabled {
		t.Error("Expected Ollama to be enabled in loaded config")
	}

	if len(loadedConfig.Phases) != 2 {
		t.Errorf("Expected 2 phase overrides, got %d", len(loadedConfig.Phases))
	}
	if loadedConfig.Phases["discovery"] != "gemini" {
		t.Errorf("Expected discovery phase override to be 'gemini', got '%s'", loadedConfig.Phases["discovery"])
	}
}

func TestLoadProjectConfigNonExistent(t *testing.T) {
	_, err := LoadProjectConfig("/nonexistent/config.yaml")
	if err == nil {
		t.Error("Expected error for non-existent config file")
	}
}

func TestLoadProjectConfigInvalidYAML(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "invalid_config.yaml")
	invalidContent := `
name: test
invalid yaml: [unclosed array
description: test project
`
	err := os.WriteFile(tempFile, []byte(invalidContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid config: %v", err)
	}

	_, err = LoadProjectConfig(tempFile)
	if err == nil {
		t.Error("Expected error for invalid YAML")
	}
}

func TestGeminiConfig(t *testing.T) {
	config := GeminiConfig{
		Model:   "gemini-pro",
		APIKey:  "test-key",
		Enabled: true,
	}

	if config.Model != "gemini-pro" {
		t.Errorf("Expected model 'gemini-pro', got '%s'", config.Model)
	}
	if config.APIKey != "test-key" {
		t.Errorf("Expected API key 'test-key', got '%s'", config.APIKey)
	}
	if !config.Enabled {
		t.Error("Expected enabled to be true")
	}
}

func TestCodexConfig(t *testing.T) {
	config := CodexConfig{
		Model:    "gpt-4",
		Provider: "openai",
		Enabled:  true,
	}

	if config.Model != "gpt-4" {
		t.Errorf("Expected model 'gpt-4', got '%s'", config.Model)
	}
	if config.Provider != "openai" {
		t.Errorf("Expected provider 'openai', got '%s'", config.Provider)
	}
	if !config.Enabled {
		t.Error("Expected enabled to be true")
	}
}

func TestOpenCodeConfig(t *testing.T) {
	config := OpenCodeConfig{
		Provider: "anthropic",
		Model:    "claude-sonnet-4-20250514",
		Enabled:  true,
	}

	if config.Provider != "anthropic" {
		t.Errorf("Expected provider 'anthropic', got '%s'", config.Provider)
	}
	if config.Model != "claude-sonnet-4-20250514" {
		t.Errorf("Expected model 'claude-sonnet-4-20250514', got '%s'", config.Model)
	}
	if !config.Enabled {
		t.Error("Expected enabled to be true")
	}
}

func TestFabricConfig(t *testing.T) {
	config := FabricConfig{
		PatternsDir: "/custom/patterns",
		Model:       "gpt-4o",
		Enabled:     true,
	}

	if config.PatternsDir != "/custom/patterns" {
		t.Errorf("Expected patterns dir '/custom/patterns', got '%s'", config.PatternsDir)
	}
	if config.Model != "gpt-4o" {
		t.Errorf("Expected model 'gpt-4o', got '%s'", config.Model)
	}
	if !config.Enabled {
		t.Error("Expected enabled to be true")
	}
}

func TestClaudeConfig(t *testing.T) {
	config := ClaudeConfig{
		Model:     "claude-sonnet-4-20250514",
		APIKey:    "test-key",
		APIKeyEnv: "ANTHROPIC_API_KEY",
		MaxTokens: 8192,
		Enabled:   true,
	}

	if config.Model != "claude-sonnet-4-20250514" {
		t.Errorf("Expected model 'claude-sonnet-4-20250514', got '%s'", config.Model)
	}
	if config.APIKey != "test-key" {
		t.Errorf("Expected API key 'test-key', got '%s'", config.APIKey)
	}
	if config.APIKeyEnv != "ANTHROPIC_API_KEY" {
		t.Errorf("Expected API key env 'ANTHROPIC_API_KEY', got '%s'", config.APIKeyEnv)
	}
	if config.MaxTokens != 8192 {
		t.Errorf("Expected max tokens 8192, got %d", config.MaxTokens)
	}
	if !config.Enabled {
		t.Error("Expected enabled to be true")
	}
}

func TestOllamaConfig(t *testing.T) {
	config := OllamaConfig{
		Model:    "llama3.2",
		Endpoint: "http://localhost:11434",
		Enabled:  true,
	}

	if config.Model != "llama3.2" {
		t.Errorf("Expected model 'llama3.2', got '%s'", config.Model)
	}
	if config.Endpoint != "http://localhost:11434" {
		t.Errorf("Expected endpoint 'http://localhost:11434', got '%s'", config.Endpoint)
	}
	if !config.Enabled {
		t.Error("Expected enabled to be true")
	}
}

func TestProjectConfigWithDefaults(t *testing.T) {
	config := &ProjectConfig{
		Name:    "minimal-project",
		Version: "1.0.0",
		Tools:   ToolsConfig{
			// All tools should have zero values by default
		},
		Phases: make(map[string]string), // Explicitly initialize
	}

	if config.Tools.Gemini.Enabled {
		t.Error("Expected Gemini to be false by default in minimal config")
	}
	if config.Tools.Codex.Enabled {
		t.Error("Expected Codex to be false by default in minimal config")
	}
	if config.Phases == nil {
		t.Error("Expected phases map to be initialized, got nil")
	}
}
