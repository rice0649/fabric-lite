package tools

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// CodexTool wraps the Codex CLI for implementation tasks
type CodexTool struct {
	BaseTool
}

// NewCodexTool creates a new Codex CLI tool wrapper
func NewCodexTool() *CodexTool {
	return &CodexTool{
		BaseTool: BaseTool{
			name:        "codex",
			description: "OpenAI Codex CLI - Advanced reasoning for implementation",
			command:     "codex",
		},
	}
}

// Execute runs the Codex CLI with the given context
func (t *CodexTool) Execute(ctx ExecutionContext) (*ExecutionResult, error) {
	// Build the prompt based on phase
	prompt := ctx.Prompt
	if prompt == "" {
		prompt = t.getPhasePrompt(ctx.Phase)
	}

	// Build command arguments
	args := []string{}

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

// ExecuteNonInteractive runs Codex and captures output
func (t *CodexTool) ExecuteNonInteractive(ctx ExecutionContext) (*ExecutionResult, error) {
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

func (t *CodexTool) getPhasePrompt(phase string) string {
	prompts := map[string]string{
		"implementation": `You are helping with the implementation phase of a software project.

Tasks:
1. Implement features based on the design specifications
2. Write clean, maintainable code following best practices
3. Add appropriate error handling and logging
4. Create unit tests for new functionality

Review the project context and help implement the planned features.`,
	}

	if p, ok := prompts[phase]; ok {
		return p
	}
	return ""
}
