package tools

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// FabricTool wraps fabric-lite for pattern-based generation
type FabricTool struct {
	BaseTool
	patternsDir string
}

// NewFabricTool creates a new fabric-lite tool wrapper
func NewFabricTool() *FabricTool {
	// Default patterns directory
	patternsDir := filepath.Join(os.Getenv("HOME"), ".config", "fabric-lite", "patterns")
	if dir := os.Getenv("FABRIC_PATTERNS_DIR"); dir != "" {
		patternsDir = dir
	}

	return &FabricTool{
		BaseTool: BaseTool{
			name:        "fabric",
			description: "fabric-lite - Pattern-based document generation",
			command:     "fabric-lite",
		},
		patternsDir: patternsDir,
	}
}

// Execute runs fabric-lite with a pattern
func (t *FabricTool) Execute(ctx ExecutionContext) (*ExecutionResult, error) {
	// Determine pattern to use
	pattern := ctx.Pattern
	if pattern == "" {
		pattern = t.getPhasePattern(ctx.Phase)
	}

	// Build command arguments
	args := []string{}

	if pattern != "" {
		args = append(args, "--pattern", pattern)
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

// ExecuteNonInteractive runs fabric-lite and captures output
func (t *FabricTool) ExecuteNonInteractive(ctx ExecutionContext) (*ExecutionResult, error) {
	pattern := ctx.Pattern
	if pattern == "" {
		pattern = t.getPhasePattern(ctx.Phase)
	}

	args := []string{}
	if pattern != "" {
		args = append(args, "--pattern", pattern)
	}
	args = append(args, ctx.Args...)

	cmd := exec.Command(t.command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Pipe input if provided
	if ctx.Prompt != "" {
		cmd.Stdin = bytes.NewBufferString(ctx.Prompt)
	}

	err := cmd.Run()

	result := &ExecutionResult{
		Output:   stdout.String(),
		Error:    stderr.String(),
		ExitCode: cmd.ProcessState.ExitCode(),
		Success:  err == nil,
	}

	return result, nil
}

// ListPatterns returns available patterns
func (t *FabricTool) ListPatterns() ([]string, error) {
	entries, err := os.ReadDir(t.patternsDir)
	if err != nil {
		// Try local patterns directory
		entries, err = os.ReadDir("patterns")
		if err != nil {
			return nil, err
		}
	}

	var patterns []string
	for _, entry := range entries {
		if entry.IsDir() {
			patterns = append(patterns, entry.Name())
		}
	}
	return patterns, nil
}

// GetPatternPath returns the path to a pattern
func (t *FabricTool) GetPatternPath(name string) string {
	// Check global patterns
	path := filepath.Join(t.patternsDir, name)
	if _, err := os.Stat(path); err == nil {
		return path
	}

	// Check local patterns
	path = filepath.Join("patterns", name)
	if _, err := os.Stat(path); err == nil {
		return path
	}

	return ""
}

func (t *FabricTool) getPhasePattern(phase string) string {
	// Map phases to default patterns
	patterns := map[string]string{
		"discovery":      "research_topic",
		"planning":       "create_architecture",
		"design":         "create_api_spec",
		"implementation": "explain_code",
		"testing":        "create_test_plan",
		"deployment":     "create_release_notes",
	}

	if p, ok := patterns[phase]; ok {
		return p
	}
	return "summarize"
}

// DeploymentPatterns returns patterns suitable for deployment phase
var DeploymentPatterns = []string{
	"create_release_notes",
	"create_changelog",
	"summarize",
	"extract_ideas",
}
