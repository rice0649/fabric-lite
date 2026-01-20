package providers

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the main configuration structure
type Config struct {
	DefaultProvider string           `yaml:"default_provider"`
	Providers       []ProviderConfig `yaml:"providers"`
	Patterns        PatternsConfig   `yaml:"patterns"`
	Sessions        SessionsConfig   `yaml:"sessions"`
}

// PatternsConfig configures pattern directories
type PatternsConfig struct {
	Directories []string `yaml:"directories"`
}

// SessionsConfig configures session storage
type SessionsConfig struct {
	Directory    string `yaml:"directory"`
	MaxHistory   int    `yaml:"max_history"`
	PersistState bool   `yaml:"persist_state"`
}

// DefaultConfig returns configuration with sensible defaults
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config", "fabric-lite")

	return &Config{
		DefaultProvider: "openai",
		Providers:       []ProviderConfig{},
		Patterns: PatternsConfig{
			Directories: []string{
				filepath.Join(configDir, "patterns"),
				"./patterns",
			},
		},
		Sessions: SessionsConfig{
			Directory:    filepath.Join(configDir, "sessions"),
			MaxHistory:   100,
			PersistState: true,
		},
	}
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = DefaultConfigPath()
	}

	// Start with defaults
	config := DefaultConfig()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config if file doesn't exist
			return config, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// Substitute environment variables before parsing
	expandedData := expandEnvVars(string(data))

	err = yaml.Unmarshal([]byte(expandedData), config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return config, nil
}

// SaveConfig saves configuration to a YAML file
func SaveConfig(config *Config, configPath string) error {
	if configPath == "" {
		configPath = DefaultConfigPath()
	}

	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(configPath, data, 0600) // Restrictive permissions for security
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// DefaultConfigPath returns the default configuration file path
func DefaultConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "fabric-lite", "config.yaml")
}

// DefaultProvidersPath returns the default providers configuration path
func DefaultProvidersPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "fabric-lite", "providers.yaml")
}

// LoadProviderConfigs loads provider configurations from providers.yaml
func LoadProviderConfigs(configPath string) ([]ProviderConfig, error) {
	if configPath == "" {
		configPath = DefaultProvidersPath()
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default providers if file doesn't exist
			return defaultProviderConfigs(), nil
		}
		return nil, fmt.Errorf("failed to read providers config: %w", err)
	}

	// Substitute environment variables before parsing
	expandedData := expandEnvVars(string(data))

	var providers struct {
		Providers []ProviderConfig `yaml:"providers"`
	}

	err = yaml.Unmarshal([]byte(expandedData), &providers)
	if err != nil {
		return nil, fmt.Errorf("failed to parse providers config: %w", err)
	}

	return providers.Providers, nil
}

// defaultProviderConfigs returns default provider configurations
func defaultProviderConfigs() []ProviderConfig {
	return []ProviderConfig{
		{
			Name: "openai",
			Type: "http",
			Config: map[string]any{
				"endpoint":    "https://api.openai.com/v1/chat/completions",
				"api_key_env": "OPENAI_API_KEY",
				"model":       "gpt-4o-mini",
				"max_tokens":  4096,
			},
		},
		{
			Name: "anthropic",
			Type: "anthropic",
			Config: map[string]any{
				"endpoint":    "https://api.anthropic.com/v1/messages",
				"api_key_env": "ANTHROPIC_API_KEY",
				"model":       "claude-sonnet-4-20250514",
				"max_tokens":  4096,
			},
		},
		{
			Name: "ollama",
			Type: "ollama",
			Config: map[string]any{
				"endpoint": "http://localhost:11434",
				"model":    "llama3.2",
			},
		},
	}
}

// expandEnvVars expands environment variables in the format ${VAR} or $VAR
func expandEnvVars(input string) string {
	// Match ${VAR} or ${VAR:-default} patterns
	re := regexp.MustCompile(`\$\{([a-zA-Z_][a-zA-Z0-9_]*)(:-[^}]*)?\}`)
	result := re.ReplaceAllStringFunc(input, func(match string) string {
		// Extract variable name and optional default
		parts := re.FindStringSubmatch(match)
		varName := parts[1]
		defaultVal := ""
		if len(parts) > 2 && parts[2] != "" {
			defaultVal = strings.TrimPrefix(parts[2], ":-")
		}

		value := os.Getenv(varName)
		if value == "" {
			return defaultVal
		}
		return value
	})

	// Also support simple $VAR format (without braces)
	re2 := regexp.MustCompile(`\$([a-zA-Z_][a-zA-Z0-9_]*)`)
	result = re2.ReplaceAllStringFunc(result, func(match string) string {
		varName := strings.TrimPrefix(match, "$")
		return os.Getenv(varName)
	})

	return result
}

// EnsureConfigDir ensures the configuration directory exists
func EnsureConfigDir() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".config", "fabric-lite")
	subdirs := []string{
		configDir,
		filepath.Join(configDir, "patterns"),
		filepath.Join(configDir, "sessions"),
	}

	for _, dir := range subdirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// WriteDefaultProviders writes default provider configuration if it doesn't exist
func WriteDefaultProviders() error {
	configPath := DefaultProvidersPath()

	// Check if already exists
	if _, err := os.Stat(configPath); err == nil {
		return nil // Already exists
	}

	// Create default content
	content := `# Fabric-Lite Provider Configuration
# Environment variables can be used with ${VAR} or ${VAR:-default} syntax

providers:
  - name: openai
    type: http
    config:
      endpoint: "https://api.openai.com/v1/chat/completions"
      api_key_env: "OPENAI_API_KEY"
      model: "gpt-4o-mini"
      max_tokens: 4096

  - name: anthropic
    type: anthropic
    config:
      endpoint: "https://api.anthropic.com/v1/messages"
      api_key_env: "ANTHROPIC_API_KEY"
      model: "claude-sonnet-4-20250514"
      max_tokens: 4096

  - name: ollama
    type: ollama
    config:
      endpoint: "http://localhost:11434"
      model: "llama3.2"

# Example custom executable provider:
#  - name: custom_script
#    type: executable
#    config:
#      executable: "/path/to/script.sh"
#      args: ["--flag", "value"]
#      timeout_seconds: 60
`

	return os.WriteFile(configPath, []byte(content), 0600)
}

// InitializeProviders loads and registers all providers from configuration
func InitializeProviders(configPath string) error {
	configs, err := LoadProviderConfigs(configPath)
	if err != nil {
		return err
	}

	return LoadProviders(configs)
}
