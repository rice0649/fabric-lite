package cli

import (
	"fmt"
	"os"

	"github.com/rice0649/fabric-lite/internal/core"
	"github.com/rice0649/fabric-lite/internal/tools"
	"github.com/spf13/cobra"
)

func newRunCmd() *cobra.Command {
	var (
		toolName    string
		patternName string
		prompt      string
		dryRun      bool
	)

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Execute AI tool for current phase",
		Long: `Run the appropriate AI tool for the current development phase.

By default, forge automatically selects the best tool for the current phase:
  - Discovery: Gemini CLI (research, context gathering)
  - Planning: OpenCode (architecture exploration)
  - Design: OpenCode (API/data modeling)
  - Implementation: Codex CLI (code generation)
  - Testing: Gemini CLI (coverage analysis)
  - Deployment: fabric-lite (documentation)

Use --tool to override the automatic selection.`,
		Example: `  # Auto-select tool for current phase
  forge run

  # Use a specific tool
  forge run --tool gemini
  forge run --tool codex

  # Use a fabric-lite pattern
  forge run --pattern summarize

  # Provide a custom prompt
  forge run --prompt "Analyze the authentication flow"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTool(toolName, patternName, prompt, dryRun)
		},
	}

	cmd.Flags().StringVarP(&toolName, "tool", "t", "", "tool to use (gemini, codex, opencode, fabric)")
	cmd.Flags().StringVarP(&patternName, "pattern", "p", "", "fabric-lite pattern to use")
	cmd.Flags().StringVar(&prompt, "prompt", "", "custom prompt for the tool")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be executed")

	return cmd
}

func runTool(toolName, patternName, prompt string, dryRun bool) error {
	// Load project state
	state, err := core.LoadProjectState(".forge/state.yaml")
	if err != nil {
		return fmt.Errorf("not a forge project (run 'forge init' first): %w", err)
	}

	if state.CurrentPhase == "" {
		return fmt.Errorf("no active phase. Run 'forge phase start <phase>' first")
	}

	// Determine which tool to use
	if patternName != "" {
		toolName = "fabric"
	}
	if toolName == "" {
		toolName = core.GetDefaultTool(state.CurrentPhase)
	}

	// Get the tool
	tool, err := tools.GetTool(toolName)
	if err != nil {
		return fmt.Errorf("failed to get tool %s: %w", toolName, err)
	}

	// Check if tool is available
	if !tool.IsAvailable() {
		return fmt.Errorf("tool %s is not installed or not in PATH", toolName)
	}

	// Build execution context
	ctx := tools.ExecutionContext{
		Phase:   state.CurrentPhase,
		Pattern: patternName,
		Prompt:  prompt,
	}

	if dryRun {
		fmt.Printf("Would execute: %s\n", tool.Name())
		fmt.Printf("Phase: %s\n", state.CurrentPhase)
		if patternName != "" {
			fmt.Printf("Pattern: %s\n", patternName)
		}
		if prompt != "" {
			fmt.Printf("Prompt: %s\n", prompt)
		}
		return nil
	}

	fmt.Printf("Running %s for phase: %s\n", tool.Name(), state.CurrentPhase)

	// Execute the tool
	result, err := tool.Execute(ctx)
	if err != nil {
		return fmt.Errorf("tool execution failed: %w", err)
	}

	// Save output as artifact
	if result.Output != "" {
		artifactPath := fmt.Sprintf(".forge/artifacts/%s/%s_output.md", state.CurrentPhase, toolName)
		if err := os.WriteFile(artifactPath, []byte(result.Output), 0644); err != nil {
			fmt.Printf("Warning: failed to save artifact: %v\n", err)
		}
	}

	// Update state
	state.AddActivity(fmt.Sprintf("Ran %s", tool.Name()))
	if err := state.Save(".forge/state.yaml"); err != nil {
		fmt.Printf("Warning: failed to update state: %v\n", err)
	}

	return nil
}
