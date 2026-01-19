package tools

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// GeminiTool wraps the Gemini CLI for research and discovery tasks
type GeminiTool struct {
	BaseTool
}

// NewGeminiTool creates a new Gemini CLI tool wrapper
func NewGeminiTool() *GeminiTool {
	return &GeminiTool{
		BaseTool: BaseTool{
			name:        "gemini",
			description: "Google Gemini CLI - Research and discovery with 1M context window",
			command:     "gemini",
		},
	}
}

// Execute runs the Gemini CLI with the given context
func (t *GeminiTool) Execute(ctx ExecutionContext) (*ExecutionResult, error) {
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

// ExecuteNonInteractive runs Gemini and captures output (for automation)
func (t *GeminiTool) ExecuteNonInteractive(ctx ExecutionContext) (*ExecutionResult, error) {
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

func (t *GeminiTool) getPhasePrompt(phase string) string {
	prompts := map[string]string{
		"discovery": `You are helping with the discovery phase of a software project.

Tasks:
1. Gather requirements and understand the problem space
2. Research similar solutions and best practices
3. Identify technical constraints and dependencies
4. Document findings for the planning phase

Please analyze the project context and help gather requirements.`,

		"testing": `You are helping with the testing phase of a software project.

Tasks:
1. Analyze the codebase for test coverage gaps
2. Suggest test cases for edge conditions
3. Help write unit and integration tests
4. Review test quality and completeness

Please examine the code and help improve test coverage.`,
	}

	if p, ok := prompts[phase]; ok {
		return p
	}
	return ""
}
