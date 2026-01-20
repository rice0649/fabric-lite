package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rice0649/fabric-lite/internal/providers"
	"gopkg.in/yaml.v3"
)

// ConfigManager handles unified configuration loading and management
type ConfigManager struct {
	config     *providers.Config
	configPath string
	loaded     bool
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(configPath string) *ConfigManager {
	if configPath == "" {
		configPath = filepath.Join(os.Getenv("HOME"), ".config", "fabric-lite", "config.yaml")
	}
	return &ConfigManager{
		configPath: configPath,
	}
}

// Load loads configuration from file, creating defaults if needed
func (cm *ConfigManager) Load() (*providers.Config, error) {
	if cm.loaded && cm.config != nil {
		return cm.config, nil
	}

	// Try to load main config file
	config, err := cm.loadMainConfig()
	if err != nil {
		// Fall back to providers.yaml for backward compatibility
		config, err = cm.loadLegacyProviderConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to load any configuration: %w", err)
		}
	}

	// Apply defaults for missing fields
	cm.applyDefaults(config)

	cm.config = config
	cm.loaded = true
	return config, nil
}

// loadMainConfig loads from the unified config.yaml file
func (cm *ConfigManager) loadMainConfig() (*providers.Config, error) {
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found: %s", cm.configPath)
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config providers.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// loadLegacyProviderConfig loads from old providers.yaml for compatibility
func (cm *ConfigManager) loadLegacyProviderConfig() (*providers.Config, error) {
	providersPath := filepath.Join(filepath.Dir(cm.configPath), "providers.yaml")

	providerConfigs, err := providers.LoadProviderConfigs(providersPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load legacy providers config: %w", err)
	}

	return &providers.Config{
		DefaultProvider: "openai",
		Providers:       providerConfigs,
		Patterns: providers.PatternsConfig{
			Directories: []string{"~/.config/fabric-lite/patterns", "./patterns"},
		},
		Sessions: providers.SessionsConfig{
			Directory:    "~/.config/fabric-lite/sessions",
			MaxHistory:   100,
			PersistState: true,
		},
	}, nil
}

// applyDefaults sets reasonable defaults for missing configuration
func (cm *ConfigManager) applyDefaults(config *providers.Config) {
	if config.DefaultProvider == "" {
		config.DefaultProvider = "openai"
	}

	if len(config.Patterns.Directories) == 0 {
		config.Patterns.Directories = []string{
			filepath.Join(os.Getenv("HOME"), ".config", "fabric-lite", "patterns"),
			"./patterns",
		}
	}

	if config.Sessions.Directory == "" {
		config.Sessions.Directory = filepath.Join(os.Getenv("HOME"), ".config", "fabric-lite", "sessions")
	}

	if config.Sessions.MaxHistory == 0 {
		config.Sessions.MaxHistory = 100
	}
}

// GetProvider returns configuration for a specific provider
func (cm *ConfigManager) GetProvider(name string) (*providers.ProviderConfig, error) {
	config, err := cm.Load()
	if err != nil {
		return nil, err
	}

	for _, provider := range config.Providers {
		if provider.Name == name {
			return &provider, nil
		}
	}

	return nil, fmt.Errorf("provider not found in configuration: %s", name)
}

// GetDefaultProvider returns the configured default provider
func (cm *ConfigManager) GetDefaultProvider() (string, error) {
	config, err := cm.Load()
	if err != nil {
		return "", err
	}
	return config.DefaultProvider, nil
}

// GetPatternPaths returns configured pattern directories
func (cm *ConfigManager) GetPatternPaths() ([]string, error) {
	config, err := cm.Load()
	if err != nil {
		return nil, err
	}
	return config.Patterns.Directories, nil
}

// GetSessionConfig returns session configuration
func (cm *ConfigManager) GetSessionConfig() (*providers.SessionsConfig, error) {
	config, err := cm.Load()
	if err != nil {
		return nil, err
	}
	return &config.Sessions, nil
}

// Save persists the current configuration to disk
func (cm *ConfigManager) Save() error {
	if cm.config == nil {
		return fmt.Errorf("no configuration to save")
	}

	// Ensure directory exists
	dir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(cm.config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// CreateDefaultConfig creates a default configuration file
func (cm *ConfigManager) CreateDefaultConfig() error {
	defaultConfig := &providers.Config{
		DefaultProvider: "openai",
		Providers: []providers.ProviderConfig{
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
		},
		Patterns: providers.PatternsConfig{
			Directories: []string{
				filepath.Join(os.Getenv("HOME"), ".config", "fabric-lite", "patterns"),
				"./patterns",
			},
		},
		Sessions: providers.SessionsConfig{
			Directory:    filepath.Join(os.Getenv("HOME"), ".config", "fabric-lite", "sessions"),
			MaxHistory:   100,
			PersistState: true,
		},
	}

	cm.config = defaultConfig
	cm.loaded = true
	return cm.Save()
}
