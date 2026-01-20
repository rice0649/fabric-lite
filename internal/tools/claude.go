package tools

// ClaudeTool is a placeholder for the Anthropic Claude tool.
type ClaudeTool struct {
	BaseTool
}

// NewClaudeTool creates a new, non-functional ClaudeTool.
func NewClaudeTool() *ClaudeTool {
	return &ClaudeTool{
		BaseTool: BaseTool{
			name:        "claude",
			description: "The Large-Scale Architect. Specialized in handling large codebases, complex refactoring, and architectural analysis.",
			command:     "claude", // Assumes a CLI tool might exist
		},
	}
}

// Execute is a placeholder and currently does nothing.
func (t *ClaudeTool) Execute(ctx ExecutionContext) (*ExecutionResult, error) {
	result := &ExecutionResult{
		Output:   "Claude tool is not yet implemented.",
		Success:  false,
		ExitCode: 1,
	}
	return result, nil
}