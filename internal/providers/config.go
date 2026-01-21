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
	// Patterns and Sessions config are now part of core.ProjectConfig
}


// DefaultConfig returns configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		DefaultProvider: "openai",
		Providers:       []ProviderConfig{},
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
