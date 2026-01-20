package tools

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// OpenCodeTool wraps OpenCode for planning and design tasks
type OpenCodeTool struct {
	BaseTool
}

// NewOpenCodeTool creates a new OpenCode tool wrapper
func NewOpenCodeTool() *OpenCodeTool {
	return &OpenCodeTool{
		BaseTool: BaseTool{
			name:        "opencode",
			description: "The Master Planner and Orchestrator. A meta-tool for creating comprehensive, multi-phase development plans.",
			command:     "opencode",
		},
	}
}

// Execute runs OpenCode with the given context
func (t *OpenCodeTool) Execute(ctx ExecutionContext) (*ExecutionResult, error) {
	// Build the prompt based on phase
	prompt := ctx.Prompt
	if prompt == "" {
		prompt = t.getPhasePrompt(ctx.Phase)
	}

	// Build command arguments
	args := []string{}

	// For planning phases, start in plan mode
	if ctx.Phase == "planning" || ctx.Phase == "design" {
		// OpenCode may support plan mode or similar
		// This depends on the actual CLI interface
	}

	// Add prompt if provided
	if prompt != "" {
		args = append(args, prompt)
	}

	// Add any additional arguments
	args = append(args, ctx.Args...)

	// Create command
	cmd := exec.Command(t.command, args...)
	cmd.Dir = ctx.WorkDir
	if cmd.Dir == "" {
		cmd.Dir, _ = os.Getwd()
	}

	// Set environment
	cmd.Env = os.Environ()
	for k, v := range ctx.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	// Connect to terminal for interactive use
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()

	result := &ExecutionResult{
		ExitCode: cmd.ProcessState.ExitCode(),
		Success:  err == nil,
	}

	if err != nil {
		result.Error = err.Error()
	}

	return result, nil
}

// ExecuteNonInteractive runs OpenCode and captures output
func (t *OpenCodeTool) ExecuteNonInteractive(ctx ExecutionContext) (*ExecutionResult, error) {
	prompt := ctx.Prompt
	if prompt == "" {
		prompt = t.getPhasePrompt(ctx.Phase)
	}

	args := []string{}
	if prompt != "" {
		args = append(args, prompt)
	}
	args = append(args, ctx.Args...)

	cmd := exec.Command(t.command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := &ExecutionResult{
		Output:   stdout.String(),
		Error:    stderr.String(),
		ExitCode: cmd.ProcessState.ExitCode(),
		Success:  err == nil,
	}

	return result, nil
}

func (t *OpenCodeTool) getPhasePrompt(phase string) string {
	prompts := map[string]string{
		"planning": `You are helping with the planning phase of a software project.

Tasks:
1. Design the high-level architecture
2. Break down the system into components
3. Define interfaces between components
4. Document technology decisions and rationale
5. Identify risks and mitigation strategies

Explore the codebase (if any) and help plan the architecture.`,

		"design": `You are helping with the design phase of a software project.

Tasks:
1. Define API endpoints and contracts
2. Design data models and schemas
3. Specify interface contracts
4. Document error handling strategies
5. Create sequence diagrams for complex flows

Build on the architecture and create detailed designs.`,
	}

	if p, ok := prompts[phase]; ok {
		return p
	}
	return ""
}
