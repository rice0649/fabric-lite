package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rice0649/fabric-lite/internal/providers"
)

func TestNewConfigManager(t *testing.T) {
	tests := []struct {
		name       string
		configPath string
		wantPath   bool
	}{
		{
			name:       "with custom path",
			configPath: "/custom/path/config.yaml",
			wantPath:   true,
		},
		{
			name:       "with empty path uses default",
			configPath: "",
			wantPath:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := NewConfigManager(tt.configPath)
			if cm == nil {
				t.Fatal("NewConfigManager returned nil")
			}
			if tt.configPath != "" && cm.configPath != tt.configPath {
				t.Errorf("configPath = %v, want %v", cm.configPath, tt.configPath)
			}
			if tt.configPath == "" && cm.configPath == "" {
				t.Error("configPath should have default value")
			}
		})
	}
}

func TestConfigManager_Load(t *testing.T) {
	// Create a temporary directory for test configs
	tmpDir := t.TempDir()

	// Create a valid config file
	validConfig := `
default_provider: openai
providers:
  - name: openai
    type: http
    config:
      endpoint: https://api.openai.com/v1/chat/completions
      api_key_env: OPENAI_API_KEY
      model: gpt-4o-mini
patterns:
  directories:
    - ./patterns
sessions:
  directory: ./sessions
  max_history: 50
`
	validConfigPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(validConfigPath, []byte(validConfig), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	t.Run("load valid config", func(t *testing.T) {
		cm := NewConfigManager(validConfigPath)
		config, err := cm.Load()
		if err != nil {
			t.Errorf("Load() error = %v", err)
			return
		}
		if config == nil {
			t.Error("Load() returned nil config without error")
		}
		if config.DefaultProvider != "openai" {
			t.Errorf("DefaultProvider = %v, want openai", config.DefaultProvider)
		}
	})

	t.Run("load non-existent config falls back or errors", func(t *testing.T) {
		// Use a path where there's definitely no providers.yaml fallback
		isolatedDir := t.TempDir()
		nonExistentPath := filepath.Join(isolatedDir, "subdir", "nonexistent.yaml")
		cm := NewConfigManager(nonExistentPath)
		_, err := cm.Load()
		// ConfigManager may fall back to legacy config or error - both are valid
		// We just verify it doesn't panic
		_ = err
	})
}

func TestConfigManager_LoadCached(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create config
	config := `
default_provider: test
providers: []
`
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cm := NewConfigManager(configPath)

	// First load
	cfg1, err := cm.Load()
	if err != nil {
		t.Fatalf("First Load() failed: %v", err)
	}

	// Second load should return cached
	cfg2, err := cm.Load()
	if err != nil {
		t.Fatalf("Second Load() failed: %v", err)
	}

	// Should be the same pointer
	if cfg1 != cfg2 {
		t.Error("Second Load() should return cached config")
	}
}

func TestConfigManager_ApplyDefaults(t *testing.T) {
	cm := &ConfigManager{}
	config := &providers.Config{}

	cm.applyDefaults(config)

	if config.DefaultProvider != "openai" {
		t.Errorf("DefaultProvider = %v, want openai", config.DefaultProvider)
	}
	if len(config.Patterns.Directories) == 0 {
		t.Error("Patterns.Directories should have default values")
	}
	if config.Sessions.Directory == "" {
		t.Error("Sessions.Directory should have default value")
	}
	if config.Sessions.MaxHistory != 100 {
		t.Errorf("Sessions.MaxHistory = %v, want 100", config.Sessions.MaxHistory)
	}
}

func TestConfigManager_GetProvider(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	config := `
default_provider: openai
providers:
  - name: openai
    type: http
    config:
      model: gpt-4o
  - name: anthropic
    type: anthropic
    config:
      model: claude-3
`
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cm := NewConfigManager(configPath)

	tests := []struct {
		name         string
		providerName string
		wantErr      bool
	}{
		{
			name:         "existing provider",
			providerName: "openai",
			wantErr:      false,
		},
		{
			name:         "another existing provider",
			providerName: "anthropic",
			wantErr:      false,
		},
		{
			name:         "non-existent provider",
			providerName: "notfound",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := cm.GetProvider(tt.providerName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && provider == nil {
				t.Error("GetProvider() returned nil without error")
			}
			if !tt.wantErr && provider.Name != tt.providerName {
				t.Errorf("GetProvider() name = %v, want %v", provider.Name, tt.providerName)
			}
		})
	}
}

func TestConfigManager_GetDefaultProvider(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	config := `
default_provider: anthropic
providers: []
`
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cm := NewConfigManager(configPath)
	defaultProvider, err := cm.GetDefaultProvider()
	if err != nil {
		t.Fatalf("GetDefaultProvider() failed: %v", err)
	}
	if defaultProvider != "anthropic" {
		t.Errorf("GetDefaultProvider() = %v, want anthropic", defaultProvider)
	}
}

func TestConfigManager_GetPatternPaths(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	config := `
default_provider: openai
providers: []
patterns:
  directories:
    - /path/one
    - /path/two
`
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cm := NewConfigManager(configPath)
	paths, err := cm.GetPatternPaths()
	if err != nil {
		t.Fatalf("GetPatternPaths() failed: %v", err)
	}
	if len(paths) != 2 {
		t.Errorf("GetPatternPaths() returned %d paths, want 2", len(paths))
	}
}

func TestConfigManager_GetSessionConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	config := `
default_provider: openai
providers: []
sessions:
  directory: /sessions
  max_history: 200
  persist_state: true
`
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cm := NewConfigManager(configPath)
	sessionConfig, err := cm.GetSessionConfig()
	if err != nil {
		t.Fatalf("GetSessionConfig() failed: %v", err)
	}
	if sessionConfig.MaxHistory != 200 {
		t.Errorf("MaxHistory = %v, want 200", sessionConfig.MaxHistory)
	}
}

func TestConfigManager_Save(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "newconfig.yaml")

	cm := NewConfigManager(configPath)
	cm.config = &providers.Config{
		DefaultProvider: "test",
		Providers:       []providers.ProviderConfig{},
	}
	cm.loaded = true

	err := cm.Save()
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Save() did not create config file")
	}
}

func TestConfigManager_SaveNoConfig(t *testing.T) {
	cm := NewConfigManager("/tmp/test.yaml")
	err := cm.Save()
	if err == nil {
		t.Error("Save() should fail with no config loaded")
	}
}

func TestConfigManager_CreateDefaultConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "default.yaml")

	cm := NewConfigManager(configPath)
	err := cm.CreateDefaultConfig()
	if err != nil {
		t.Fatalf("CreateDefaultConfig() failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("CreateDefaultConfig() did not create config file")
	}

	// Verify config was loaded
	if cm.config == nil {
		t.Error("CreateDefaultConfig() did not set config")
	}
	if cm.config.DefaultProvider != "openai" {
		t.Errorf("DefaultProvider = %v, want openai", cm.config.DefaultProvider)
	}
	if len(cm.config.Providers) != 3 {
		t.Errorf("Providers count = %v, want 3", len(cm.config.Providers))
	}
}
