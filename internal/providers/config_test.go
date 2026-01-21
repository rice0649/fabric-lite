package providers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.DefaultProvider != "openai" {
		t.Errorf("Expected default provider to be 'openai', got %s", config.DefaultProvider)
	}

	if len(config.Providers) != 0 {
		t.Errorf("Expected no providers in default config, got %d", len(config.Providers))
	}

	if len(config.Patterns.Directories) != 2 {
		t.Errorf("Expected 2 pattern directories, got %d", len(config.Patterns.Directories))
	}

	if config.Sessions.MaxHistory != 100 {
		t.Errorf("Expected max history 100, got %d", config.Sessions.MaxHistory)
	}

	if !config.Sessions.PersistState {
		t.Errorf("Expected persist state to be true, got %v", config.Sessions.PersistState)
	}
}

func TestDefaultConfigPath(t *testing.T) {
	path := DefaultConfigPath()
	if path == "" {
		t.Error("Expected config path to not be empty")
	}

	expectedSuffix := string(os.PathSeparator) + ".config" + string(os.PathSeparator) + "fabric-lite" + string(os.PathSeparator) + "config.yaml"
	if len(path) < len(expectedSuffix) || path[len(path)-len(expectedSuffix):] != expectedSuffix {
		t.Errorf("Expected path to end with %s, got %s", expectedSuffix, path)
	}
}

func TestDefaultProvidersPath(t *testing.T) {
	path := DefaultProvidersPath()
	if path == "" {
		t.Error("Expected providers path to not be empty")
	}

	expectedSuffix := string(os.PathSeparator) + ".config" + string(os.PathSeparator) + "fabric-lite" + string(os.PathSeparator) + "providers.yaml"
	if len(path) < len(expectedSuffix) || path[len(path)-len(expectedSuffix):] != expectedSuffix {
		t.Errorf("Expected path to end with %s, got %s", expectedSuffix, path)
	}
}

func TestExpandEnvVars(t *testing.T) {
	// Set test environment variables
	os.Setenv("TEST_VAR", "test_value")
	os.Setenv("TEST_VAR_WITH_DEFAULT", "")
	defer func() {
		os.Unsetenv("TEST_VAR")
		os.Unsetenv("TEST_VAR_WITH_DEFAULT")
	}()

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "No variables here",
			expected: "No variables here",
		},
		{
			input:    "${TEST_VAR}",
			expected: "test_value",
		},
		{
			input:    "${TEST_VAR:-default_value}",
			expected: "test_value",
		},
		{
			input:    "${TEST_VAR_WITH_DEFAULT:-default_value}",
			expected: "default_value",
		},
		{
			input:    "$TEST_VAR",
			expected: "test_value",
		},
		{
			input:    "prefix_${TEST_VAR}_suffix",
			expected: "prefix_test_value_suffix",
		},
		{
			input:    "prefix_${TEST_VAR}_suffix",
			expected: "prefix_test_value_suffix",
		},
		{
			input:    "${NONEXISTENT_VAR:-fallback}",
			expected: "fallback",
		},
		{
			input:    "${NONEXISTENT_VAR}",
			expected: "",
		},
		{
			input:    "$NONEXISTENT_VAR",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := expandEnvVars(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestDefaultProviderConfigs(t *testing.T) {
	configs := defaultProviderConfigs()

	if len(configs) != 3 {
		t.Errorf("Expected 3 default provider configs, got %d", len(configs))
	}

	// Check OpenAI config
	var openaiConfig *ProviderConfig
	for _, cfg := range configs {
		if cfg.Name == "openai" {
			openaiConfig = &cfg
			break
		}
	}
	if openaiConfig == nil {
		t.Error("Expected OpenAI provider config")
	} else {
		if openaiConfig.Type != "http" {
			t.Errorf("Expected OpenAI type to be 'http', got %s", openaiConfig.Type)
		}
		if openaiConfig.Config["model"] != "gpt-4o-mini" {
			t.Errorf("Expected OpenAI model to be 'gpt-4o-mini', got %v", openaiConfig.Config["model"])
		}
	}

	// Check Anthropic config
	var anthropicConfig *ProviderConfig
	for _, cfg := range configs {
		if cfg.Name == "anthropic" {
			anthropicConfig = &cfg
			break
		}
	}
	if anthropicConfig == nil {
		t.Error("Expected Anthropic provider config")
	} else {
		if anthropicConfig.Type != "anthropic" {
			t.Errorf("Expected Anthropic type to be 'anthropic', got %s", anthropicConfig.Type)
		}
		if anthropicConfig.Config["model"] != "claude-sonnet-4-20250514" {
			t.Errorf("Expected Anthropic model to be 'claude-sonnet-4-20250514', got %v", anthropicConfig.Config["model"])
		}
	}

	// Check Ollama config
	var ollamaConfig *ProviderConfig
	for _, cfg := range configs {
		if cfg.Name == "ollama" {
			ollamaConfig = &cfg
			break
		}
	}
	if ollamaConfig == nil {
		t.Error("Expected Ollama provider config")
	} else {
		if ollamaConfig.Type != "ollama" {
			t.Errorf("Expected Ollama type to be 'ollama', got %s", ollamaConfig.Type)
		}
		if ollamaConfig.Config["endpoint"] != "http://localhost:11434" {
			t.Errorf("Expected Ollama endpoint to be 'http://localhost:11434', got %v", ollamaConfig.Config["endpoint"])
		}
	}
}

func TestSaveConfig(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.yaml")

	config := &Config{
		DefaultProvider: "test",
		Providers: []ProviderConfig{
			{
				Name: "test-provider",
				Type: "http",
				Config: map[string]any{
					"endpoint": "https://api.test.com",
				},
			},
		},
		Patterns: PatternsConfig{
			Directories: []string{"/tmp/patterns"},
		},
		Sessions: SessionsConfig{
			Directory:    "/tmp/sessions",
			MaxHistory:   50,
			PersistState: false,
		},
	}

	err := SaveConfig(config, configPath)
	if err != nil {
		t.Errorf("Expected no error saving config, got %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Expected config file to exist")
	}

	// Load and verify content
	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Errorf("Expected no error loading config, got %v", err)
	}

	if loadedConfig.DefaultProvider != "test" {
		t.Errorf("Expected loaded default provider to be 'test', got %s", loadedConfig.DefaultProvider)
	}

	if len(loadedConfig.Providers) != 1 {
		t.Errorf("Expected 1 provider in loaded config, got %d", len(loadedConfig.Providers))
	}

	if loadedConfig.Sessions.MaxHistory != 50 {
		t.Errorf("Expected loaded max history to be 50, got %d", loadedConfig.Sessions.MaxHistory)
	}
}

func TestLoadConfigNonExistent(t *testing.T) {
	// Test loading non-existent config returns default
	nonExistentPath := filepath.Join(t.TempDir(), "nonexistent.yaml")
	config, err := LoadConfig(nonExistentPath)

	if err != nil {
		t.Errorf("Expected no error loading non-existent config, got %v", err)
	}

	if config == nil {
		t.Error("Expected config to be returned for non-existent file")
	}

	// Should return default config values
	if config.DefaultProvider != "openai" {
		t.Errorf("Expected default provider to be 'openai', got %s", config.DefaultProvider)
	}
}

func TestLoadProviderConfigs(t *testing.T) {
	// Create temporary providers config
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-providers.yaml")

	providersContent := `
providers:
  - name: test-http
    type: http
    config:
      endpoint: "https://api.test.com"
      api_key_env: "TEST_API_KEY"
      model: "test-model"

  - name: test-ollama
    type: ollama
    config:
      endpoint: "http://localhost:11434"
      model: "llama3.2"
`

	err := os.WriteFile(configPath, []byte(providersContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test providers config: %v", err)
	}

	configs, err := LoadProviderConfigs(configPath)
	if err != nil {
		t.Errorf("Expected no error loading provider configs, got %v", err)
	}

	if len(configs) != 2 {
		t.Errorf("Expected 2 provider configs, got %d", len(configs))
	}

	// Check first provider
	if configs[0].Name != "test-http" {
		t.Errorf("Expected first provider name to be 'test-http', got %s", configs[0].Name)
	}
	if configs[0].Type != "http" {
		t.Errorf("Expected first provider type to be 'http', got %s", configs[0].Type)
	}
	if configs[0].Config["endpoint"] != "https://api.test.com" {
		t.Errorf("Expected endpoint to be 'https://api.test.com', got %v", configs[0].Config["endpoint"])
	}

	// Check second provider
	if configs[1].Name != "test-ollama" {
		t.Errorf("Expected second provider name to be 'test-ollama', got %s", configs[1].Name)
	}
	if configs[1].Type != "ollama" {
		t.Errorf("Expected second provider type to be 'ollama', got %s", configs[1].Type)
	}
}

func TestLoadProviderConfigsWithEnvVars(t *testing.T) {
	// Set test environment variable
	os.Setenv("TEST_ENDPOINT", "https://api.from-env.com")
	defer os.Unsetenv("TEST_ENDPOINT")

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-providers.yaml")

	providersContent := `
providers:
  - name: test-http
    type: http
    config:
      endpoint: "${TEST_ENDPOINT}"
      api_key_env: "TEST_API_KEY"
`

	err := os.WriteFile(configPath, []byte(providersContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test providers config: %v", err)
	}

	configs, err := LoadProviderConfigs(configPath)
	if err != nil {
		t.Errorf("Expected no error loading provider configs, got %v", err)
	}

	if len(configs) != 1 {
		t.Errorf("Expected 1 provider config, got %d", len(configs))
	}

	if configs[0].Config["endpoint"] != "https://api.from-env.com" {
		t.Errorf("Expected endpoint to be 'https://api.from-env.com', got %v", configs[0].Config["endpoint"])
	}
}

func TestEnsureConfigDir(t *testing.T) {
	// Test with temporary home directory
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	os.Setenv("HOME", tempHome)

	err := EnsureConfigDir()
	if err != nil {
		t.Errorf("Expected no error ensuring config dir, got %v", err)
	}

	configDir := filepath.Join(tempHome, ".config", "fabric-lite")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Error("Expected config directory to exist")
	}

	patternsDir := filepath.Join(configDir, "patterns")
	if _, err := os.Stat(patternsDir); os.IsNotExist(err) {
		t.Error("Expected patterns directory to exist")
	}

	sessionsDir := filepath.Join(configDir, "sessions")
	if _, err := os.Stat(sessionsDir); os.IsNotExist(err) {
		t.Error("Expected sessions directory to exist")
	}
}

func TestWriteDefaultProviders(t *testing.T) {
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	os.Setenv("HOME", tempDir)

	// Ensure config directory exists first
	configDir := filepath.Join(tempDir, ".config", "fabric-lite")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	err = WriteDefaultProviders()
	if err != nil {
		t.Errorf("Expected no error writing default providers, got %v", err)
	}

	providersPath := DefaultProvidersPath()
	if _, err := os.Stat(providersPath); os.IsNotExist(err) {
		t.Error("Expected providers file to exist")
	}

	// Check file content
	content, err := os.ReadFile(providersPath)
	if err != nil {
		t.Errorf("Expected no error reading providers file, got %v", err)
	}

	contentStr := string(content)
	if !contains(contentStr, "openai") {
		t.Error("Expected content to contain 'openai'")
	}
	if !contains(contentStr, "anthropic") {
		t.Error("Expected content to contain 'anthropic'")
	}
	if !contains(contentStr, "ollama") {
		t.Error("Expected content to contain 'ollama'")
	}
}

func TestInitializeProviders(t *testing.T) {
	// Clear registry for testing
	providerRegistry = make(map[string]Provider)

	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	os.Setenv("HOME", tempDir)

	// Ensure config directory exists first
	configDir := filepath.Join(tempDir, ".config", "fabric-lite")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	// Write default providers first
	err = WriteDefaultProviders()
	if err != nil {
		t.Fatalf("Failed to write default providers: %v", err)
	}

	err = InitializeProviders("")
	if err != nil {
		t.Errorf("Expected no error initializing providers, got %v", err)
	}

	// Check that providers were loaded
	names := ListProviders()
	if len(names) == 0 {
		t.Error("Expected providers to be loaded")
	}

	// Should have at least the default providers
	expectedProviders := []string{"openai", "anthropic", "ollama"}
	for _, expected := range expectedProviders {
		found := false
		for _, name := range names {
			if name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected provider %s to be loaded", expected)
		}
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
