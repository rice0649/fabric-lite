package tools

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// ClaudeTool wraps the Claude CLI for large-scale architecture and refactoring tasks
type ClaudeTool struct {
	BaseTool
}

// NewClaudeTool creates a new Claude CLI tool wrapper
func NewClaudeTool() *ClaudeTool {
	return &ClaudeTool{
		BaseTool: BaseTool{
			name:        "claude",
			description: "The Large-Scale Architect. Specialized in handling large codebases, complex refactoring, and architectural analysis.",
			command:     "claude",
		},
	}
}

// Execute runs the Claude CLI with the given context
func (t *ClaudeTool) Execute(ctx ExecutionContext) (*ExecutionResult, error) {
	// Build the prompt based on phase
	prompt := ctx.Prompt
	if prompt == "" {
		prompt = t.getPhasePrompt(ctx.Phase)
	}

	// Build command arguments
	args := []string{}

	// Add prompt if provided (claude uses -p for print/prompt mode)
	if prompt != "" {
		args = append(args, "-p", prompt)
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

// ExecuteNonInteractive runs Claude and captures output (for automation)
func (t *ClaudeTool) ExecuteNonInteractive(ctx ExecutionContext) (*ExecutionResult, error) {
	prompt := ctx.Prompt
	if prompt == "" {
		prompt = t.getPhasePrompt(ctx.Phase)
	}

	args := []string{}
	if prompt != "" {
		args = append(args, "-p", prompt)
	}
	args = append(args, ctx.Args...)

	cmd := exec.Command(t.command, args...)
	cmd.Dir = ctx.WorkDir
	if cmd.Dir == "" {
		cmd.Dir, _ = os.Getwd()
	}

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

func (t *ClaudeTool) getPhasePrompt(phase string) string {
	prompts := map[string]string{
		"planning": `You are helping with the planning phase of a software project.

Tasks:
1. Analyze the project requirements and constraints
2. Design the high-level architecture
3. Break down work into manageable components
4. Identify potential risks and mitigations

Please help create a comprehensive project plan.`,

		"design": `You are helping with the design phase of a software project.

Tasks:
1. Create detailed technical designs
2. Define interfaces and data structures
3. Document design decisions and trade-offs
4. Review designs for completeness and correctness

Please help refine the technical design.`,

		"implementation": `You are helping with the implementation phase of a software project.

Tasks:
1. Write clean, maintainable code
2. Follow established patterns and conventions
3. Handle edge cases and error conditions
4. Ensure code is well-documented

Please help implement the required functionality.`,

		"deployment": `You are helping with the deployment phase of a software project.

Tasks:
1. Prepare deployment configurations
2. Document deployment procedures
3. Set up monitoring and logging
4. Create rollback procedures

Please help prepare for production deployment.`,
	}

	if p, ok := prompts[phase]; ok {
		return p
	}
	return ""
}
