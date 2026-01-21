package tools

import (
	"os"
	"strings"
	"testing"
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
	tool := NewCodexTool()

	if tool.Name() != "codex" {
		t.Errorf("Expected tool name to be 'codex', got %s", tool.Name())
	}
	expectedDesc := "The Code Generation Specialist. A meta-tool for writing, refactoring, and explaining code based on a plan."
	if tool.Description() != expectedDesc {
		t.Errorf("Expected description to be '%s', got '%s'", expectedDesc, tool.Description())
	}
}

func TestCodexToolLoadConfig(t *testing.T) {
	tool := NewCodexTool()

	// Test loading config (should use defaults if no config files exist)
	tool.loadConfig()

	if tool.Config.Provider == "" {
		t.Error("Expected provider to be set after loadConfig")
	}

	if tool.Config.Model == "" {
		t.Error("Expected model to be set after loadConfig")
	}
}

func TestCodexToolExecute(t *testing.T) {
	tool := NewCodexTool()

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
		if strings.Contains(err.Error(), "failed to get codex provider") {
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
	tool := NewCodexTool()

	// Create a temp dir with no config files to force defaults
	tempDir := t.TempDir()
	ctx := ExecutionContext{
		Prompt:  "Test prompt",
		WorkDir: tempDir,
	}

	// Save original HOME and set temp to avoid loading real config
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	// Execute should work with default provider (ollama), but if ollama is not available
	// it should return an error from the provider lookup
	_, err := tool.Execute(ctx)

	// Either success (if ollama is available) or error (if not) is acceptable
	// The important thing is no panic
	if err != nil && !strings.Contains(err.Error(), "failed to get codex provider") {
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
	tool := NewOllamaTool()

	if tool.Name() != "ollama" {
		t.Errorf("Expected tool name to be 'ollama', got %s", tool.Name())
	}
	expectedDesc := "The Quick Task Automator. Runs local models like Llama3 for fast, simple tasks like boilerplate generation and formatting."
	if tool.Description() != expectedDesc {
		t.Errorf("Expected description to be '%s', got '%s'", expectedDesc, tool.Description())
	}
	// Ollama doesn't use an external command
	if tool.GetCommand() != "" {
		t.Errorf("Expected command to be empty, got %s", tool.GetCommand())
	}
}

func TestOllamaToolIsAvailable(t *testing.T) {
	tool := NewOllamaTool()

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
