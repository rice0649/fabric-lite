package tools

import (
	"context"
	"os"
	"strings"
	"testing"
	"path/filepath"
	"gopkg.in/yaml.v3"
	"time"

	"github.com/rice0649/fabric-lite/internal/core"
	"github.com/rice0649/fabric-lite/internal/providers"
)

func TestExecutionContext(t *testing.T) {
	ctx := ExecutionContext{
		Phase:   "implementation",
		Pattern: "summarize",
		Prompt:  "Test prompt",
		Args:    []string{"--verbose"},
		Env:     map[string]string{"TEST_VAR": "test_value"},
		WorkDir: "/tmp",
	}

	if ctx.Phase != "implementation" {
		t.Errorf("Expected Phase to be 'implementation', got %s", ctx.Phase)
	}
	if ctx.Pattern != "summarize" {
		t.Errorf("Expected Pattern to be 'summarize', got %s", ctx.Pattern)
	}
	if ctx.Prompt != "Test prompt" {
		t.Errorf("Expected Prompt to be 'Test prompt', got %s", ctx.Prompt)
	}
	if len(ctx.Args) != 1 || ctx.Args[0] != "--verbose" {
		t.Errorf("Expected Args to be ['--verbose'], got %v", ctx.Args)
	}
	if ctx.Env["TEST_VAR"] != "test_value" {
		t.Errorf("Expected Env['TEST_VAR'] to be 'test_value', got %s", ctx.Env["TEST_VAR"])
	}
	if ctx.WorkDir != "/tmp" {
		t.Errorf("Expected WorkDir to be '/tmp', got %s", ctx.WorkDir)
	}
}

func TestExecutionResult(t *testing.T) {
	result := ExecutionResult{
		Output:   "test output",
		Error:    "test error",
		ExitCode: 1,
		Success:  false,
	}

	if result.Output != "test output" {
		t.Errorf("Expected Output to be 'test output', got %s", result.Output)
	}
	if result.Error != "test error" {
		t.Errorf("Expected Error to be 'test error', got %s", result.Error)
	}
	if result.ExitCode != 1 {
		t.Errorf("Expected ExitCode to be 1, got %d", result.ExitCode)
	}
	if result.Success {
		t.Errorf("Expected Success to be false, got %v", result.Success)
	}
}

func TestBaseTool(t *testing.T) {
	tool := BaseTool{
		name:        "test-tool",
		description: "Test tool description",
		command:     "test-command",
	}

	if tool.Name() != "test-tool" {
		t.Errorf("Expected Name() to return 'test-tool', got %s", tool.Name())
	}
	if tool.Description() != "Test tool description" {
		t.Errorf("Expected Description() to return 'Test tool description', got %s", tool.Description())
	}
	if tool.GetCommand() != "test-command" {
		t.Errorf("Expected GetCommand() to return 'test-command', got %s", tool.GetCommand())
	}
}

func TestCheckCommand(t *testing.T) {
	// Test with a command that should exist
	if !checkCommand("go") {
		t.Error("Expected 'go' command to be found")
	}

	// Test with a command that should not exist
	if checkCommand("nonexistent-command-12345") {
		t.Error("Expected 'nonexistent-command-12345' to not be found")
	}
}

// MockTool implements Tool interface for testing
type MockTool struct {
	BaseTool
	available bool
}

func (m *MockTool) IsAvailable() bool {
	return m.available
}

func (m *MockTool) Execute(ctx ExecutionContext) (*ExecutionResult, error) {
	return &ExecutionResult{
		Output:   "mock output",
		ExitCode: 0,
		Success:  true,
	}, nil
}

// mockOllamaProvider implements providers.Provider for testing purposes
type mockOllamaProvider struct {
	name string
}

func (m *mockOllamaProvider) Name() string {
	return m.name
}

func (m *mockOllamaProvider) Execute(ctx context.Context, request providers.CompletionRequest) (*providers.CompletionResponse, error) {
	return &providers.CompletionResponse{
		Content:  "mock Ollama response",
		Model:    request.Model,
		Tokens:   10,
		Duration: 1 * time.Millisecond,
	}, nil
}

func (m *mockOllamaProvider) ExecuteStream(ctx context.Context, request providers.CompletionRequest) (<-chan providers.StreamChunk, error) {
	chunks := make(chan providers.StreamChunk, 1)
	go func() {
		defer close(chunks)
		chunks <- providers.StreamChunk{Content: "mock Ollama stream chunk", Done: true}
	}()
	return chunks, nil
}

func (m *mockOllamaProvider) IsAvailable() bool {
	return true // Always available for testing
}

func (m *mockOllamaProvider) GetModels() []string {
	return []string{"llama3:latest", "llama2"}
}

func TestToolRegistry(t *testing.T) {
	// Save original registry
	originalRegistry := toolRegistry
	defer func() { toolRegistry = originalRegistry }()

	// Clear registry for testing
	toolRegistry = make(map[string]Tool)

	mockTool := &MockTool{
		BaseTool: BaseTool{
			name:        "mock-tool",
			description: "Mock tool for testing",
			command:     "echo",
		},
		available: true,
	}

	// Test RegisterTool
	RegisterTool(mockTool)

	if len(toolRegistry) != 1 {
		t.Errorf("Expected registry to have 1 tool, got %d", len(toolRegistry))
	}

	// Test GetTool
	tool, err := GetTool("mock-tool")
	if err != nil {
		t.Errorf("Expected to find mock-tool, got error: %v", err)
	}
	if tool.Name() != "mock-tool" {
		t.Errorf("Expected tool name to be 'mock-tool', got %s", tool.Name())
	}

	_, err = GetTool("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent tool, got nil")
	}

	// Test ListTools
	names := ListTools()
	if len(names) != 1 {
		t.Errorf("Expected ListTools to return 1 name, got %d", len(names))
	}
	if names[0] != "mock-tool" {
		t.Errorf("Expected tool name to be 'mock-tool', got %s", names[0])
	}

	// Test ListAvailableTools with available command
	available := ListAvailableTools()
	if len(available) != 1 {
		t.Errorf("Expected 1 available tool, got %d", len(available))
	}

	// Test with unavailable command
	unavailableTool := &MockTool{
		BaseTool: BaseTool{
			name:        "unavailable-tool",
			description: "Unavailable tool",
			command:     "nonexistent-command-12345",
		},
		available: false,
	}
	RegisterTool(unavailableTool)

	available = ListAvailableTools()
	if len(available) != 1 { // Only mock-tool should be available
		t.Errorf("Expected 1 available tool, got %d", len(available))
	}
}

func TestNewClaudeTool(t *testing.T) {
	tool := NewClaudeTool()

	if tool.Name() != "claude" {
		t.Errorf("Expected tool name to be 'claude', got %s", tool.Name())
	}
	expectedDesc := "The Large-Scale Architect. Specialized in handling large codebases, complex refactoring, and architectural analysis."
	if tool.Description() != expectedDesc {
		t.Errorf("Expected description to be '%s', got '%s'", expectedDesc, tool.Description())
	}
	if tool.GetCommand() != "claude" {
		t.Errorf("Expected command to be 'claude', got %s", tool.GetCommand())
	}
}

func TestClaudeToolGetPhasePrompt(t *testing.T) {
	tool := NewClaudeTool()

	tests := []struct {
		phase    string
		expected string
	}{
		{
			phase:    "planning",
			expected: "You are helping with the planning phase",
		},
		{
			phase:    "design",
			expected: "You are helping with the design phase",
		},
		{
			phase:    "implementation",
			expected: "You are helping with the implementation phase",
		},
		{
			phase:    "deployment",
			expected: "You are helping with the deployment phase",
		},
		{
			phase:    "nonexistent",
			expected: "",
		},
		{
			phase:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.phase, func(t *testing.T) {
			prompt := tool.getPhasePrompt(tt.phase)
			if tt.expected == "" {
				if prompt != "" {
					t.Errorf("Expected empty prompt for phase '%s', got '%s'", tt.phase, prompt)
				}
			} else {
				if prompt == "" {
					t.Errorf("Expected non-empty prompt for phase '%s', got empty string", tt.phase)
				}
				// Check if prompt contains expected text
				if len(prompt) < len(tt.expected) || prompt[:len(tt.expected)] != tt.expected {
					t.Errorf("Expected prompt to contain '%s', got '%s'", tt.expected, prompt[:min(len(prompt), len(tt.expected))])
				}
			}
		})
	}
}

func TestClaudeToolExecuteNonInteractive(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	tool := NewClaudeTool()

	// Create execution context
	ctx := ExecutionContext{
		Phase:   "planning",
		Prompt:  "Test prompt",
		WorkDir: tempDir,
		Env:     map[string]string{"TEST_VAR": "test_value"},
		Args:    []string{"--version"},
	}

	// Execute() tool (this will fail if claude is not installed, but we can test the structure)
	result, err := tool.ExecuteNonInteractive(ctx)

	// If claude is not available, we expect a failure
	if !tool.IsAvailable() {
		if result == nil {
			t.Error("Expected result to be returned even if command fails")
		}
		return
	}

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected result to be non-nil")
	}
}

func TestNewCodexTool(t *testing.T) {
	// Mock config for CodexTool
	mockCodexConfig := core.CodexConfig{
		Provider: "ollama",
		Model:    "llama3:latest",
		Enabled:  true,
	}

	// Mock ProviderManager (no actual functionality needed for this test)
	mockProviderManager := core.NewProviderManager(&providers.Config{})

	tool := NewCodexTool(mockCodexConfig, mockProviderManager)

	if tool.Name() != "codex" {
		t.Errorf("Expected tool name to be 'codex', got %s", tool.Name())
	}
	expectedDesc := "The Code Generation Specialist. A meta-tool for writing, refactoring, and explaining code based on a plan."
	if tool.Description() != expectedDesc {
		t.Errorf("Expected description to be '%s', got '%s'", expectedDesc, tool.Description())
	}
}

func TestCodexToolExecute(t *testing.T) {
	// Setup: Create a temporary config file for testing
	tempDir := t.TempDir()
	tempConfigFile := filepath.Join(tempDir, "config.yaml")

	mockConfig := core.NewProjectConfig("test-project", "")
	mockConfig.Tools.Codex.Enabled = true
	mockConfig.Tools.Codex.Provider = "ollama"
	mockConfig.Tools.Ollama.Enabled = true // Ensure Ollama is enabled for ProviderManager
	mockConfig.Tools.Ollama.Model = "llama3:latest"
	mockConfig.Tools.Ollama.Endpoint = "http://localhost:11434"

	configData, err := yaml.Marshal(mockConfig)
	if err != nil {
		t.Fatalf("Failed to marshal mock config: %v", err)
	}
	if err := os.WriteFile(tempConfigFile, configData, 0644); err != nil {
		t.Fatalf("Failed to write mock config file: %v", err)
	}

	cm := core.NewConfigManager(tempConfigFile)
	cfg, err := cm.Load() // cm.Load now also sets the default ProviderManager
	if err != nil {
		t.Fatalf("Failed to load config for test: %v", err)
	}

	// Register mock Ollama provider
	mockOllama := &mockOllamaProvider{name: "ollama"}
	core.GetDefaultProviderManager().AddProvider("ollama", mockOllama)

	// Register tools after config and provider manager are loaded
	RegisterConfiguredTools(cfg.Tools.Codex, core.GetDefaultProviderManager())

	tool, err := GetTool("codex")
	if err != nil {
		t.Fatalf("Failed to get codex tool: %v", err)
	}

	// Create execution context
	ctx := ExecutionContext{
		Prompt:  "Write a hello world function",
		WorkDir: t.TempDir(),
	}

	// Execute without config (should load defaults)
	result, err := tool.Execute(ctx)

	// We expect this to potentially fail if provider is not available
	// but we can test the structure
	if err != nil {
		// Expected if provider is not available
		if strings.Contains(err.Error(), "provider not available") {
			return // Expected failure
		}
		t.Errorf("Unexpected error: %v", err)
	}

	if result != nil {
		// Should have exit code
		if result.ExitCode == 0 && !result.Success {
			t.Error("Inconsistent result: exit code 0 but success false")
		}
	}
}

func TestCodexToolExecuteWithMissingProvider(t *testing.T) {
	// Setup: Create a temporary config file with an intentionally missing provider for this test
	tempDir := t.TempDir()
	tempConfigFile := filepath.Join(tempDir, "config.yaml")

	mockConfig := core.NewProjectConfig("test-project-missing", "")
	mockConfig.Tools.Codex.Enabled = true
	mockConfig.Tools.Codex.Provider = "nonexistent-provider" // This will cause the failure
	// Ensure at least one provider is enabled in ProjectConfig so ProviderManager initializes
	mockConfig.Tools.Ollama.Enabled = true 
	mockConfig.Tools.Ollama.Model = "llama3:latest"
	mockConfig.Tools.Ollama.Endpoint = "http://localhost:11434"

	configData, err := yaml.Marshal(mockConfig)
	if err != nil {
		t.Fatalf("Failed to marshal mock config: %v", err)
	}
	if err := os.WriteFile(tempConfigFile, configData, 0644); err != nil {
		t.Fatalf("Failed to write mock config file: %v", err)
	}

	cm := core.NewConfigManager(tempConfigFile)
	cfg, err := cm.Load() // cm.Load now also sets the default ProviderManager
	if err != nil {
		t.Fatalf("Failed to load config for test: %v", err)
	}

	// Register mock Ollama provider (even though this test uses "nonexistent-provider" for Codex)
	// This ensures that GetDefaultProviderManager doesn't panic if other parts of the system
	// (e.g., IsAvailable for other tools) try to access it.
	mockOllama := &mockOllamaProvider{name: "ollama"}
	core.GetDefaultProviderManager().AddProvider("ollama", mockOllama)

	// Register tools after config and provider manager are loaded
	RegisterConfiguredTools(cfg.Tools.Codex, core.GetDefaultProviderManager())

	tool, err := GetTool("codex")
	if err != nil {
		t.Fatalf("Failed to get codex tool: %v", err)
	}

	ctx := ExecutionContext{
		Prompt:  "Test prompt",
		WorkDir: tempDir,
	}

	_, err = tool.Execute(ctx)

	if err == nil {
		t.Error("Expected an error for missing provider, got nil")
	} else if !strings.Contains(err.Error(), "codex tool is not available (either not enabled or provider 'nonexistent-provider' is not ready)") {
		t.Errorf("Unexpected error type: %v", err)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestNewGeminiTool(t *testing.T) {
	tool := NewGeminiTool()

	if tool.Name() != "gemini" {
		t.Errorf("Expected tool name to be 'gemini', got %s", tool.Name())
	}
	expectedDesc := "The Researcher/Reviewer. Expert in discovery, validation, and proposal review, grounded by Google Search."
	if tool.Description() != expectedDesc {
		t.Errorf("Expected description to be '%s', got '%s'", expectedDesc, tool.Description())
	}
	if tool.GetCommand() != "gemini" {
		t.Errorf("Expected command to be 'gemini', got %s", tool.GetCommand())
	}
}

func TestGeminiToolGetPhasePrompt(t *testing.T) {
	tool := NewGeminiTool()

	tests := []struct {
		phase    string
		expected string
	}{
		{
			phase:    "discovery",
			expected: "You are helping with the discovery phase",
		},
		{
			phase:    "testing",
			expected: "You are helping with the testing phase",
		},
		{
			phase:    "nonexistent",
			expected: "",
		},
		{
			phase:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.phase, func(t *testing.T) {
			prompt := tool.getPhasePrompt(tt.phase)
			if tt.expected == "" {
				if prompt != "" {
					t.Errorf("Expected empty prompt for phase '%s', got '%s'", tt.phase, prompt)
				}
			} else {
				if prompt == "" {
					t.Errorf("Expected non-empty prompt for phase '%s', got empty string", tt.phase)
				}
				if !strings.Contains(prompt, tt.expected) {
					t.Errorf("Expected prompt to contain '%s', got '%s'", tt.expected, prompt)
				}
			}
		})
	}
}

func TestNewOllamaTool(t *testing.T) {
	tool := NewOllamaTool("http://localhost:11434")

	if tool.Name() != "ollama" {
		t.Errorf("Expected tool name to be 'ollama', got %s", tool.Name())
	}
	expectedDesc := "Interact with the Ollama CLI to pull and list models. Requires Ollama to be installed and running."
	if tool.Description() != expectedDesc {
		t.Errorf("Expected description to be '%s', got '%s'", expectedDesc, tool.Description())
	}
	if tool.GetCommand() != "ollama" {
		t.Errorf("Expected command to be 'ollama', got %s", tool.GetCommand())
	}
}

func TestOllamaToolIsAvailable(t *testing.T) {
	tool := NewOllamaTool("http://localhost:11434")

	// IsAvailable should return a boolean without panicking
	// The actual result depends on whether Ollama is running
	available := tool.IsAvailable()

	// Just verify it returns without error (bool type check is implicit)
	_ = available
}

func TestNewOpenCodeTool(t *testing.T) {
	tool := NewOpenCodeTool()

	if tool.Name() != "opencode" {
		t.Errorf("Expected tool name to be 'opencode', got %s", tool.Name())
	}
	expectedDesc := "The Master Planner and Orchestrator. A meta-tool for creating comprehensive, multi-phase development plans."
	if tool.Description() != expectedDesc {
		t.Errorf("Expected description to be '%s', got '%s'", expectedDesc, tool.Description())
	}
	if tool.GetCommand() != "opencode" {
		t.Errorf("Expected command to be 'opencode', got %s", tool.GetCommand())
	}
}

func TestOpenCodeToolGetPhasePrompt(t *testing.T) {
	tool := NewOpenCodeTool()

	tests := []struct {
		phase    string
		expected string
	}{
		{
			phase:    "planning",
			expected: "You are helping with the planning phase",
		},
		{
			phase:    "design",
			expected: "You are helping with the design phase",
		},
		{
			phase:    "nonexistent",
			expected: "",
		},
		{
			phase:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.phase, func(t *testing.T) {
			prompt := tool.getPhasePrompt(tt.phase)
			if tt.expected == "" {
				if prompt != "" {
					t.Errorf("Expected empty prompt for phase '%s', got '%s'", tt.phase, prompt)
				}
			} else {
				if prompt == "" {
					t.Errorf("Expected non-empty prompt for phase '%s', got empty string", tt.phase)
				}
				if !strings.Contains(prompt, tt.expected) {
					t.Errorf("Expected prompt to contain '%s', got '%s'", tt.expected, prompt)
				}
			}
		})
	}
}

func TestNewFabricTool(t *testing.T) {
	tool := NewFabricTool()

	if tool.Name() != "fabric" {
		t.Errorf("Expected tool name to be 'fabric', got %s", tool.Name())
	}
	expectedDesc := "fabric-lite - Pattern-based document generation"
	if tool.Description() != expectedDesc {
		t.Errorf("Expected description to be '%s', got '%s'", expectedDesc, tool.Description())
	}
	if tool.GetCommand() != "fabric-lite" {
		t.Errorf("Expected command to be 'fabric-lite', got %s", tool.GetCommand())
	}
}
