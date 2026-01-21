package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rice0649/fabric-lite/internal/providers"
	"gopkg.in/yaml.v3"
)

var (
	defaultConfig *ProjectConfig
	configOnce    sync.Once
)

// SetDefaultConfig sets the global Config instance.
// This should only be called once during application initialization.
func SetDefaultConfig(cfg *ProjectConfig) {
	configOnce.Do(func() {
		defaultConfig = cfg
	})
}

// GetDefaultConfig returns the global Config instance.
// Panics if the config has not been set.
func GetDefaultConfig() *ProjectConfig {
	if defaultConfig == nil {
		panic("Default Config has not been initialized. Call SetDefaultConfig first.")
	}
	return defaultConfig
}

// ConfigManager handles unified configuration loading and management
type ConfigManager struct {
	config     *ProjectConfig
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
func (cm *ConfigManager) Load() (*ProjectConfig, error) {
	if cm.loaded && cm.config != nil {
		SetDefaultConfig(cm.config) // Ensure global config is set even on subsequent loads
		return cm.config, nil
	}

	// Try to load main config file
	config, err := cm.loadMainConfig()
	if err != nil {
		// If main config not found, create a default one
		if os.IsNotExist(err) || strings.Contains(err.Error(), "config file not found") {
			config = NewProjectConfig("", "") // Create a default ProjectConfig
		} else {
			return nil, fmt.Errorf("failed to load main configuration: %w", err)
		}
	}

	// Apply defaults for missing fields
	cm.applyDefaults(config)

	cm.config = config
	cm.loaded = true

	// Initialize and set the global ProviderManager
	// We need to construct a providers.Config from ProjectConfig for ProviderManager
	providerCfg := &providers.Config{
		DefaultProvider: "ollama", // Default to ollama if not specified
		Providers:       []providers.ProviderConfig{},
	}

	// Populate providers from ProjectConfig's tool configs
	if config.Tools.Ollama.Enabled {
		providerCfg.Providers = append(providerCfg.Providers, providers.ProviderConfig{
			Name: "ollama",
			Type: "ollama",
			Config: map[string]any{
				"endpoint": config.Tools.Ollama.Endpoint,
				"model":    config.Tools.Ollama.Model,
			},
		})
	}
	if config.Tools.Gemini.Enabled {
		providerCfg.Providers = append(providerCfg.Providers, providers.ProviderConfig{
			Name: "gemini",
			Type: "gemini", // Assuming a gemini provider type exists
			Config: map[string]any{
				"api_key": config.Tools.Gemini.APIKey,
				"model":   config.Tools.Gemini.Model,
			},
		})
	}
	if config.Tools.Claude.Enabled {
		providerCfg.Providers = append(providerCfg.Providers, providers.ProviderConfig{
			Name: "anthropic", // Claude tool uses anthropic provider type
			Type: "anthropic",
			Config: map[string]any{
				"api_key":     config.Tools.Claude.APIKey,
				"api_key_env": config.Tools.Claude.APIKeyEnv,
				"model":       config.Tools.Claude.Model,
				"max_tokens":  config.Tools.Claude.MaxTokens,
			},
		})
	}

	pm := NewProviderManager(providerCfg)
	if err := pm.InitializeAll(); err != nil {
		return nil, fmt.Errorf("failed to initialize providers: %w", err)
	}
	SetDefaultProviderManager(pm)
	SetDefaultConfig(config) // Set global config after successful load and PM init

	return config, nil
}

// loadMainConfig loads from the unified config.yaml file
func (cm *ConfigManager) loadMainConfig() (*ProjectConfig, error) {
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found: %s", cm.configPath)
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config ProjectConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// applyDefaults sets reasonable defaults for missing configuration
func (cm *ConfigManager) applyDefaults(config *ProjectConfig) {
	if config.Name == "" {
		config.Name = "default-project"
	}
	if config.Version == "" {
		config.Version = "1.0.0"
	}

	// Apply defaults for PatternsConfig
	if len(config.Patterns.Directories) == 0 {
		homeDir, _ := os.UserHomeDir()
		configDir := filepath.Join(homeDir, ".config", "fabric-lite")
		config.Patterns.Directories = []string{
			filepath.Join(configDir, "patterns"),
			"./patterns",
		}
	}

	// Apply defaults for SessionsConfig
	if config.Sessions.Directory == "" {
		homeDir, _ := os.UserHomeDir()
		configDir := filepath.Join(homeDir, ".config", "fabric-lite")
		config.Sessions.Directory = filepath.Join(configDir, "sessions")
	}
	if config.Sessions.MaxHistory == 0 {
		config.Sessions.MaxHistory = 100
	}
	// PersistState default is already true in NewProjectConfig

	// Apply defaults for ToolsConfig (if not set by NewProjectConfig)
	if config.Tools.Gemini.Model == "" {
		config.Tools.Gemini.Model = "gemini-2.0-flash-exp"
	}
	if config.Tools.Gemini.APIKey == "" {
		// Attempt to load from env var GEMINI_API_KEY if not explicitly set
		if os.Getenv("GEMINI_API_KEY") != "" {
			config.Tools.Gemini.APIKey = os.Getenv("GEMINI_API_KEY")
		}
	}
	// Gemini.Enabled is handled by NewProjectConfig

	if config.Tools.Codex.Model == "" {
		config.Tools.Codex.Model = "o3-mini"
	}
	if config.Tools.Codex.Provider == "" {
		config.Tools.Codex.Provider = "ollama"
	}
	// Codex.Enabled is handled by NewProjectConfig

	if config.Tools.OpenCode.Provider == "" {
		config.Tools.OpenCode.Provider = "anthropic"
	}
	if config.Tools.OpenCode.Model == "" {
		config.Tools.OpenCode.Model = "claude-sonnet-4-20250514"
	}
	// OpenCode.Enabled is handled by NewProjectConfig

	if config.Tools.Fabric.Model == "" {
		config.Tools.Fabric.Model = "gpt-4o-mini"
	}
	// Fabric.Enabled is handled by NewProjectConfig

	if config.Tools.Claude.Model == "" {
		config.Tools.Claude.Model = "claude-sonnet-4-20250514"
	}
	if config.Tools.Claude.MaxTokens == 0 {
		config.Tools.Claude.MaxTokens = 4096
	}
	// Claude.Enabled is handled by NewProjectConfig, APIKey and APIKeyEnv as well

	if config.Tools.Ollama.Model == "" {
		config.Tools.Ollama.Model = "llama3.2"
	}
	if config.Tools.Ollama.Endpoint == "" {
		config.Tools.Ollama.Endpoint = "http://localhost:11434"
	}
	// Ollama.Enabled is handled by NewProjectConfig
}

// Save persists the current configuration to disk
func (cm *ConfigManager) Save(config *ProjectConfig) error {
	if config == nil {
		return fmt.Errorf("no configuration to save")
	}

	// Ensure directory exists
	dir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
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
	defaultConfig := NewProjectConfig("default-project", "") // Using NewProjectConfig from internal/core/config.go

	// Apply any additional defaults or overrides if necessary
	cm.applyDefaults(defaultConfig)

	cm.config = defaultConfig
	cm.loaded = true
	return cm.Save(defaultConfig) // Pass the defaultConfig to the Save method
}
