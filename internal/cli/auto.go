package cli

import (
	"fmt"

	"github.com/rice0649/fabric-lite/internal/core"
	"github.com/rice0649/fabric-lite/internal/tools"
	"github.com/spf13/cobra"
)

func newAutoCmd() *cobra.Command {
	var (
		fromPhase      string
		untilPhase     string
		skipValidation bool
		dryRun         bool
	)

	cmd := &cobra.Command{
		Use:   "auto",
		Short: "Run phases automatically with validation checkpoints",
		Long: `Run development phases sequentially with optional AI-powered validation.

By default, runs all phases from the current state to deployment.
Use --from and --until to specify a range of phases.

Between each phase, an AI validator checks that the phase output
meets quality criteria before proceeding. Use --skip-validation
to disable this.

If interrupted, forge auto will resume from where it left off.`,
		Example: `  # Run all phases from current state
  forge auto

  # Run phases from planning to design
  forge auto --from planning --until design

  # Run without validation checkpoints
  forge auto --skip-validation

  # Preview what would run
  forge auto --dry-run`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAuto(fromPhase, untilPhase, skipValidation, dryRun)
		},
	}

	cmd.Flags().StringVar(&fromPhase, "from", "", "start from this phase")
	cmd.Flags().StringVar(&untilPhase, "until", "", "stop after this phase")
	cmd.Flags().BoolVar(&skipValidation, "skip-validation", false, "skip AI validation between phases")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be executed without running")

	return cmd
}

func runAuto(fromPhase, untilPhase string, skipValidation, dryRun bool) error {
	// Load project config
	config, err := core.LoadProjectConfig(".forge/config.yaml")
	if err != nil {
		return fmt.Errorf("not a forge project (run 'forge init' first): %w", err)
	}

	// Load project state
	statePath := ".forge/state.yaml"
	state, err := core.LoadProjectState(statePath)
	if err != nil {
		return fmt.Errorf("failed to load project state: %w", err)
	}

	// Create auto runner
	runner := core.NewAutoRunner(config, state, statePath)

	// Check for resumable state
	canResume, lastPhase, nextPhase := runner.GetResumeInfo()
	if canResume && fromPhase == "" {
		fmt.Printf("Resuming from phase: %s (last completed: %s)\n", nextPhase, lastPhase)
	}

	// Set up validator if not skipping
	if !skipValidation {
		runner.Validator = createValidator()
	}

	// Set up executor
	runner.Executor = &defaultPhaseExecutor{
		config: config,
		state:  state,
		dryRun: dryRun,
	}

	// Dry run mode
	if dryRun {
		return showDryRun(runner, fromPhase, untilPhase, skipValidation)
	}

	// Run phases
	return runner.Run(fromPhase, untilPhase, skipValidation)
}

// showDryRun displays what would be executed
func showDryRun(runner *core.AutoRunner, from, until string, skipValidation bool) error {
	phases, err := getPhaseRangeForDisplay(from, until, runner)
	if err != nil {
		return err
	}

	fmt.Println("Dry run - would execute:")
	fmt.Println()

	for i, phase := range phases {
		phaseInfo := core.GetPhase(phase)
		fmt.Printf("%d. Phase: %s\n", i+1, phase)
		fmt.Printf("   Tool: %s\n", phaseInfo.PrimaryTool)
		fmt.Printf("   Artifacts: %v\n", phaseInfo.Artifacts)
		if !skipValidation {
			fmt.Printf("   Validation: enabled\n")
		}
		fmt.Println()
	}

	if skipValidation {
		fmt.Println("Note: Validation checkpoints disabled")
	}

	return nil
}

// getPhaseRangeForDisplay returns phases for dry-run display
func getPhaseRangeForDisplay(from, until string, runner *core.AutoRunner) ([]string, error) {
	allPhases := core.PhaseNames()

	startIdx := 0
	if from != "" {
		for i, p := range allPhases {
			if p == from {
				startIdx = i
				break
			}
		}
	} else if runner.State.Auto != nil && runner.State.Auto.LastCompletedPhase != "" {
		for i, p := range allPhases {
			if p == runner.State.Auto.LastCompletedPhase {
				startIdx = i + 1
				break
			}
		}
	}

	endIdx := len(allPhases) - 1
	if until != "" {
		for i, p := range allPhases {
			if p == until {
				endIdx = i
				break
			}
		}
	}

	if startIdx > endIdx || startIdx >= len(allPhases) {
		return nil, fmt.Errorf("invalid phase range")
	}

	return allPhases[startIdx : endIdx+1], nil
}

// createValidator creates a phase validator using fabric-lite
func createValidator() *core.PhaseValidator {
	fabricTool := tools.NewFabricTool()
	if !fabricTool.IsAvailable() {
		fmt.Println("Warning: fabric-lite not available, validation disabled")
		return nil
	}

	return &core.PhaseValidator{
		ExecuteFunc: func(pattern, input string) (string, error) {
			ctx := tools.ExecutionContext{
				Pattern: pattern,
				Prompt:  input,
			}
			result, err := fabricTool.ExecuteNonInteractive(ctx)
			if err != nil {
				return "", err
			}
			if !result.Success {
				return "", fmt.Errorf("validation failed: %s", result.Error)
			}
			return result.Output, nil
		},
	}
}

// defaultPhaseExecutor executes phases using the configured tools
type defaultPhaseExecutor struct {
	config *core.ProjectConfig
	state  *core.ProjectState
	dryRun bool
}

func (e *defaultPhaseExecutor) Execute(phase string) error {
	if e.dryRun {
		return nil
	}

	phaseInfo := core.GetPhase(phase)
	if phaseInfo == nil {
		return fmt.Errorf("unknown phase: %s", phase)
	}

	toolName := phaseInfo.PrimaryTool

	// Get the tool
	tool, err := tools.GetTool(toolName)
	if err != nil {
		return fmt.Errorf("get tool %s: %w", toolName, err)
	}

	if !tool.IsAvailable() {
		return fmt.Errorf("tool %s is not installed or not in PATH", toolName)
	}

	// Build execution context
	ctx := tools.ExecutionContext{
		Phase: phase,
	}

	fmt.Printf("  → Executing with %s...\n", toolName)

	// Execute the tool
	result, err := tool.Execute(ctx)
	if err != nil {
		return fmt.Errorf("tool execution failed: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("tool returned error: %s", result.Error)
	}

	return nil
}
