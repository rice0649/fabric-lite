package core

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ProjectConfig represents the forge project configuration
type ProjectConfig struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description,omitempty"`
	Template    string            `yaml:"template,omitempty"`
	Version     string            `yaml:"version"`
	Tools       ToolsConfig       `yaml:"tools"`
	Phases      map[string]string `yaml:"phases,omitempty"` // phase -> custom tool override
}

// ToolsConfig holds configuration for AI tools
type ToolsConfig struct {
	Gemini   GeminiConfig   `yaml:"gemini"`
	Codex    CodexConfig    `yaml:"codex"`
	OpenCode OpenCodeConfig `yaml:"opencode"`
	Fabric   FabricConfig   `yaml:"fabric"`
	Claude   ClaudeConfig   `yaml:"claude"`
	Ollama   OllamaConfig   `yaml:"ollama"`
}

type GeminiConfig struct {
	Model   string `yaml:"model,omitempty"`
	APIKey  string `yaml:"api_key,omitempty"` // Usually from env
	Enabled bool   `yaml:"enabled"`
}

type CodexConfig struct {
	Model    string `yaml:"model,omitempty"`
	Provider string `yaml:"provider,omitempty"` // openai, azure, etc.
	Enabled  bool   `yaml:"enabled"`
}

type OpenCodeConfig struct {
	Provider string `yaml:"provider,omitempty"` // anthropic, openai, etc.
	Model    string `yaml:"model,omitempty"`
	Enabled  bool   `yaml:"enabled"`
}

type FabricConfig struct {
	PatternsDir string `yaml:"patterns_dir,omitempty"`
	Model       string `yaml:"model,omitempty"`
	Enabled     bool   `yaml:"enabled"`
}

type ClaudeConfig struct {
	Model     string `yaml:"model,omitempty"`
	APIKey    string `yaml:"api_key,omitempty"`     // Usually from ANTHROPIC_API_KEY env
	APIKeyEnv string `yaml:"api_key_env,omitempty"` // Env var name for API key
	MaxTokens int    `yaml:"max_tokens,omitempty"`
	Enabled   bool   `yaml:"enabled"`
}

type OllamaConfig struct {
	Model    string `yaml:"model,omitempty"`
	Endpoint string `yaml:"endpoint,omitempty"` // Default: http://localhost:11434
	Enabled  bool   `yaml:"enabled"`
}

// NewProjectConfig creates a new project configuration
func NewProjectConfig(name, template string) *ProjectConfig {
	return &ProjectConfig{
		Name:     name,
		Template: template,
		Version:  "1.0.0",
		Tools: ToolsConfig{
			Gemini: GeminiConfig{
				Model:   "gemini-2.0-flash-exp",
				Enabled: true,
			},
			Codex: CodexConfig{
				Model:   "o3-mini",
				Enabled: true,
			},
			OpenCode: OpenCodeConfig{
				Provider: "anthropic",
				Model:    "claude-sonnet-4-20250514",
				Enabled:  true,
			},
			Fabric: FabricConfig{
				Model:   "gpt-4o-mini",
				Enabled: true,
			},
			Claude: ClaudeConfig{
				Model:     "claude-sonnet-4-20250514",
				APIKeyEnv: "ANTHROPIC_API_KEY",
				MaxTokens: 4096,
				Enabled:   false, // Disabled by default, user must enable
			},
			Ollama: OllamaConfig{
				Model:    "llama3.2",
				Endpoint: "http://localhost:11434",
				Enabled:  false, // Disabled by default, user must enable
			},
		},
	}
}

// Save writes the config to a YAML file
func (c *ProjectConfig) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadProjectConfig loads configuration from a YAML file
func LoadProjectConfig(path string) (*ProjectConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg ProjectConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
