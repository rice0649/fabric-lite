package executor

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rice0649/fabric-lite/internal/providers"
)

// MockProvider implements providers.Provider for testing
type MockProvider struct {
	ProviderName string
	Available    bool
	Models       []string
}

func (m *MockProvider) Name() string {
	return m.ProviderName
}

func (m *MockProvider) Execute(ctx context.Context, request providers.CompletionRequest) (*providers.CompletionResponse, error) {
	return &providers.CompletionResponse{
		Content:  "mock response",
		Model:    m.Models[0],
		Tokens:   10,
		Duration: 10 * time.Millisecond,
	}, nil
}

func (m *MockProvider) ExecuteStream(ctx context.Context, request providers.CompletionRequest) (<-chan providers.StreamChunk, error) {
	chunks := make(chan providers.StreamChunk, 1)
	go func() {
		defer close(chunks)
		chunks <- providers.StreamChunk{Content: "mock response", Done: true}
	}()
	return chunks, nil
}

func (m *MockProvider) IsAvailable() bool {
	return m.Available
}

func (m *MockProvider) GetModels() []string {
	return m.Models
}

func TestNewPatternExecutor(t *testing.T) {
	executor := NewPatternExecutor()

	if executor == nil {
		t.Error("Expected executor to be non-nil")
	}

	if executor.providers == nil {
		t.Error("Expected providers map to be initialized")
	}

	if executor.patternsDir == "" {
		t.Error("Expected patternsDir to be set")
	}
}

func TestGetPatternsDir(t *testing.T) {
	// Test with local patterns directory
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	os.Chdir(tempDir)

	// Create patterns directory
	err := os.Mkdir("patterns", 0755)
	if err != nil {
		t.Fatalf("Failed to create patterns directory: %v", err)
	}

	executor := NewPatternExecutor()
	patternsDir := executor.GetPatternsDir()

	expected := "patterns"
	if patternsDir != expected {
		t.Errorf("Expected patterns directory to be '%s', got '%s'", expected, patternsDir)
	}
}

func TestGetPatternsDirFallback(t *testing.T) {
	// Test fallback to home directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	os.Chdir(t.TempDir()) // Directory without patterns

	executor := NewPatternExecutor()
	patternsDir := executor.GetPatternsDir()

	homeDir := os.Getenv("HOME")
	expected := filepath.Join(homeDir, ".config", "fabric-lite", "patterns")

	if patternsDir != expected {
		t.Errorf("Expected patterns directory to be '%s', got '%s'", expected, patternsDir)
	}
}

func TestLoadProvider(t *testing.T) {
	executor := NewPatternExecutor()

	// Test with nil config
	err := executor.LoadProvider("test", nil)
	if err == nil {
		t.Error("Expected error for nil config")
	}

	// Test with empty config
	config := &providers.Config{
		Providers: []providers.ProviderConfig{},
	}

	err = executor.LoadProvider("nonexistent", config)
	if err == nil {
		t.Error("Expected error for nonexistent provider")
	}

	// Test with valid config
	config.Providers = []providers.ProviderConfig{
		{
			Name: "test-provider",
			Type: "http",
			Config: map[string]any{
				"endpoint": "https://api.test.com",
				"api_key":  "test-key",
			},
		},
	}

	err = executor.LoadProvider("test-provider", config)
	if err != nil {
		t.Errorf("Expected no error for valid provider config, got %v", err)
	}

	// Check that provider was loaded
	if _, ok := executor.providers["test-provider"]; !ok {
		t.Error("Expected provider to be loaded")
	}
}

func TestLoadProviderDirect(t *testing.T) {
	executor := NewPatternExecutor()

	mockProvider := &MockProvider{
		ProviderName: "mock-provider",
		Available:    true,
		Models:       []string{"test-model"},
	}

	executor.LoadProviderDirect("mock-provider", mockProvider)

	// Check that provider was loaded
	if _, ok := executor.providers["mock-provider"]; !ok {
		t.Error("Expected provider to be loaded")
	}
}

func TestListPatterns(t *testing.T) {
	// Create temporary patterns directory
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	os.Chdir(tempDir)
	err := os.Mkdir("patterns", 0755)
	if err != nil {
		t.Fatalf("Failed to create patterns directory: %v", err)
	}

	// Create test patterns
	createTestPattern(t, "patterns/test1", "Test System 1", "Test User 1")
	createTestPattern(t, "patterns/test2", "Test System 2", "") // No user prompt

	// Create a file (should be ignored)
	err = os.WriteFile("patterns/not-a-pattern.txt", []byte("not a pattern"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	executor := NewPatternExecutor()
	executor.patternsDir = "patterns"

	patterns, err := executor.ListPatterns()
	if err != nil {
		t.Errorf("Expected no error listing patterns, got %v", err)
	}

	if len(patterns) != 2 {
		t.Errorf("Expected 2 patterns, got %d", len(patterns))
	}

	// Check pattern names
	patternNames := make(map[string]bool)
	for _, pattern := range patterns {
		patternNames[pattern.Name] = true
	}

	if !patternNames["test1"] {
		t.Error("Expected to find test1 pattern")
	}
	if !patternNames["test2"] {
		t.Error("Expected to find test2 pattern")
	}
}

func TestLoadPattern(t *testing.T) {
	// Create temporary patterns directory
	tempDir := t.TempDir()
	executor := NewPatternExecutor()
	executor.patternsDir = tempDir

	// Create test pattern
	createTestPattern(t, tempDir+"/test", "Test System Content", "Test User Content")

	pattern, err := executor.loadPattern("test")
	if err != nil {
		t.Errorf("Expected no error loading pattern, got %v", err)
	}

	if pattern.Name != "test" {
		t.Errorf("Expected pattern name to be 'test', got %s", pattern.Name)
	}
	if pattern.System != "Test System Content\n" {
		t.Errorf("Expected system content to be 'Test System Content\\n', got %s", pattern.System)
	}
	if pattern.User != "Test User Content\n" {
		t.Errorf("Expected user content to be 'Test User Content\\n', got %s", pattern.User)
	}
}

func TestLoadPatternMissingSystem(t *testing.T) {
	tempDir := t.TempDir()
	executor := NewPatternExecutor()
	executor.patternsDir = tempDir

	// Create pattern directory without system.md
	err := os.MkdirAll(tempDir+"/test", 0755)
	if err != nil {
		t.Fatalf("Failed to create pattern directory: %v", err)
	}

	_, err = executor.loadPattern("test")
	if err == nil {
		t.Error("Expected error for missing system.md")
	}
}

func TestLoadPatternOnlySystem(t *testing.T) {
	tempDir := t.TempDir()
	executor := NewPatternExecutor()
	executor.patternsDir = tempDir

	// Create pattern with only system.md
	createTestPattern(t, tempDir+"/test", "Test System Content", "")

	pattern, err := executor.loadPattern("test")
	if err != nil {
		t.Errorf("Expected no error loading pattern with only system, got %v", err)
	}

	if pattern.User != "" {
		t.Errorf("Expected user content to be empty, got %s", pattern.User)
	}
}

func TestExtractDescription(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "Description line",
			content:  "# DESCRIPTION This is a test pattern\nMore content",
			expected: "This is a test pattern",
		},
		{
			name:     "Identity and purpose",
			content:  "# IDENTITY and PURPOSE\nYou are a helpful assistant\nMore content",
			expected: "You are a helpful assistant",
		},
		{
			name:     "No description",
			content:  "Just some content\nWithout special markers",
			expected: "Pattern execution",
		},
		{
			name:     "Mixed content",
			content:  "Some content\n# DESCRIPTION Test description\nMore content",
			expected: "Test description",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractDescription(tt.content)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestBuildRequest(t *testing.T) {
	executor := NewPatternExecutor()

	pattern := &PatternInfo{
		Name:   "test",
		System: "Test system",
		User:   "Test user",
	}

	// Test with user prompt
	request := executor.buildRequest(pattern, "test input", "test-model", false)
	if request.System != "Test system" {
		t.Errorf("Expected system to be 'Test system', got %s", request.System)
	}
	expectedPrompt := "Test user\n\nInput:\ntest input"
	if request.Prompt != expectedPrompt {
		t.Errorf("Expected prompt to be '%s', got '%s'", expectedPrompt, request.Prompt)
	}
	if request.Model != "test-model" {
		t.Errorf("Expected model to be 'test-model', got %s", request.Model)
	}
	if request.Stream {
		t.Error("Expected stream to be false")
	}
	if request.MaxTokens != 4096 {
		t.Errorf("Expected max tokens to be 4096, got %d", request.MaxTokens)
	}

	// Test without user prompt
	pattern.User = ""
	request = executor.buildRequest(pattern, "test input", "", true)
	if request.Prompt != "test input" {
		t.Errorf("Expected prompt to be 'test input', got %s", request.Prompt)
	}
	if request.Model != "gpt-4o-mini" { // Default model
		t.Errorf("Expected model to be 'gpt-4o-mini', got %s", request.Model)
	}
	if !request.Stream {
		t.Error("Expected stream to be true")
	}
}

func TestExecute(t *testing.T) {
	executor := NewPatternExecutor()
	ctx := context.Background()

	// Test with no provider loaded
	_, err := executor.Execute(ctx, "test", "input", "nonexistent")
	if err == nil {
		t.Error("Expected error for unloaded provider")
	}

	// Load a mock provider
	mockProvider := &MockProvider{
		ProviderName: "mock-provider",
		Available:    true,
		Models:       []string{"test-model"},
	}
	executor.LoadProviderDirect("mock-provider", mockProvider)

	// Test with non-existent pattern
	_, err = executor.Execute(ctx, "nonexistent", "input", "mock-provider")
	if err == nil {
		t.Error("Expected error for non-existent pattern")
	}
}

func TestExecuteWithUnavailableProvider(t *testing.T) {
	executor := NewPatternExecutor()
	ctx := context.Background()

	// Load unavailable mock provider
	mockProvider := &MockProvider{
		ProviderName: "unavailable-provider",
		Available:    false,
		Models:       []string{"test-model"},
	}
	executor.LoadProviderDirect("unavailable-provider", mockProvider)

	_, err := executor.Execute(ctx, "test", "input", "unavailable-provider")
	if err == nil {
		t.Error("Expected error for unavailable provider")
	}
}

func TestExecuteWithOptions(t *testing.T) {
	executor := NewPatternExecutor()
	ctx := context.Background()

	// Load a mock provider
	mockProvider := &MockProvider{
		ProviderName: "mock-provider",
		Available:    true,
		Models:       []string{"test-model"},
	}
	executor.LoadProviderDirect("mock-provider", mockProvider)

	// Test with non-existent pattern
	_, err := executor.ExecuteWithOptions(ctx, "nonexistent", "input", "mock-provider", "test-model", false)
	if err == nil {
		t.Error("Expected error for non-existent pattern")
	}
}

func TestExecuteStream(t *testing.T) {
	executor := NewPatternExecutor()
	ctx := context.Background()

	// Test with no provider loaded
	_, err := executor.ExecuteStream(ctx, "test", "input", "nonexistent", "test-model")
	if err == nil {
		t.Error("Expected error for unloaded provider")
	}

	// Load a mock provider
	mockProvider := &MockProvider{
		ProviderName: "mock-provider",
		Available:    true,
		Models:       []string{"test-model"},
	}
	executor.LoadProviderDirect("mock-provider", mockProvider)

	// Test with non-existent pattern
	_, err = executor.ExecuteStream(ctx, "nonexistent", "input", "mock-provider", "test-model")
	if err == nil {
		t.Error("Expected error for non-existent pattern")
	}
}

func TestExecuteStreamWithUnavailableProvider(t *testing.T) {
	executor := NewPatternExecutor()
	ctx := context.Background()

	// Load unavailable mock provider
	mockProvider := &MockProvider{
		ProviderName: "unavailable-provider",
		Available:    false,
		Models:       []string{"test-model"},
	}
	executor.LoadProviderDirect("unavailable-provider", mockProvider)

	_, err := executor.ExecuteStream(ctx, "test", "input", "unavailable-provider", "test-model")
	if err == nil {
		t.Error("Expected error for unavailable provider")
	}
}

// Helper function to create test pattern files
func createTestPattern(t *testing.T, path, systemContent, userContent string) {
	t.Helper()

	err := os.MkdirAll(path, 0755)
	if err != nil {
		t.Fatalf("Failed to create pattern directory: %v", err)
	}

	if systemContent != "" {
		err = os.WriteFile(filepath.Join(path, "system.md"), []byte(systemContent+"\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to create system.md: %v", err)
		}
	}

	if userContent != "" {
		err = os.WriteFile(filepath.Join(path, "user.md"), []byte(userContent+"\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to create user.md: %v", err)
		}
	}
}
