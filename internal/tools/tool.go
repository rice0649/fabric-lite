package tools

import (
	"fmt"
	"os/exec"

	"github.com/rice0649/fabric-lite/internal/core"
)

// Tool defines the interface for AI coding assistant tools
type Tool interface {
	// Name returns the tool's identifier
	Name() string

	// Description returns a human-readable description
	Description() string

	// IsAvailable checks if the tool is installed and accessible
	IsAvailable() bool

	// Execute runs the tool with the given context
	Execute(ctx ExecutionContext) (*ExecutionResult, error)

	// GetCommand returns the base command for the tool
	GetCommand() string
}

// ExecutionContext provides context for tool execution
type ExecutionContext struct {
	Phase   string            // Current development phase
	Pattern string            // Pattern name (for fabric-lite)
	Prompt  string            // Custom prompt
	Args    []string          // Additional arguments
	Env     map[string]string // Environment variables
	WorkDir string            // Working directory
}

// ExecutionResult contains the result of a tool execution
type ExecutionResult struct {
	Output   string // Standard output
	Error    string // Standard error
	ExitCode int    // Exit code
	Success  bool   // Whether execution was successful
}

// toolRegistry holds all registered tools
var toolRegistry = make(map[string]Tool)

// RegisterTool adds a tool to the registry
func RegisterTool(tool Tool) {
	toolRegistry[tool.Name()] = tool
}

// GetTool returns a tool by name
func GetTool(name string) (Tool, error) {
	if tool, ok := toolRegistry[name]; ok {
		return tool, nil
	}
	return nil, fmt.Errorf("unknown tool: %s", name)
}

// ListTools returns all registered tool names
func ListTools() []string {
	names := make([]string, 0, len(toolRegistry))
	for name := range toolRegistry {
		names = append(names, name)
	}
	return names
}

// ListAvailableTools returns all tools that are currently available
func ListAvailableTools() []Tool {
	available := make([]Tool, 0)
	for _, tool := range toolRegistry {
		if tool.IsAvailable() {
			available = append(available, tool)
		}
	}
	return available
}

// checkCommand checks if a command exists in PATH
func checkCommand(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// BaseTool provides common functionality for tools
type BaseTool struct {
	name        string
	description string
	command     string
}

func (t *BaseTool) Name() string        { return t.name }
func (t *BaseTool) Description() string { return t.description }
func (t *BaseTool) GetCommand() string  { return t.command }

func (t *BaseTool) IsAvailable() bool {
	return checkCommand(t.command)
}

func init() {
	// Register all tools
	RegisterTool(NewGeminiTool())
	// RegisterTool(NewCodexTool()) // CodexTool is now registered via RegisterConfiguredTools
	RegisterTool(NewOpenCodeTool())
	RegisterTool(NewFabricTool())
	RegisterTool(NewClaudeTool())
	RegisterTool(NewOllamaTool("http://localhost:11434"))
}

// RegisterConfiguredTools registers tools that require configuration from the main config.
func RegisterConfiguredTools(codexConfig core.CodexConfig, pm *core.ProviderManager) {
	RegisterTool(NewCodexTool(codexConfig, pm))
}
