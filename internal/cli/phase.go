package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rice0649/fabric-lite/internal/core"
	"github.com/spf13/cobra"
)

func newPhaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "phase",
		Short: "Manage development phases",
		Long: `Manage the development phases of your project.

Phases represent structured stages in the development lifecycle:
  1. discovery     - Research and requirements gathering
  2. planning      - Architecture and component design
  3. design        - API and data model definition
  4. implementation - Code development
  5. testing       - Test creation and execution
  6. deployment    - Documentation and release`,
	}

	cmd.AddCommand(newPhaseListCmd())
	cmd.AddCommand(newPhaseStartCmd())
	cmd.AddCommand(newPhaseCompleteCmd())
	cmd.AddCommand(newPhaseInfoCmd())

	return cmd
}

func newPhaseListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all phases with status",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := core.LoadProjectState(".forge/state.yaml")
			if err != nil {
				// Show phases without status if not initialized
				fmt.Println("Development Phases:")
				fmt.Println()
				for _, p := range core.AllPhases {
					fmt.Printf("  ○ %s - %s\n", p.Name, p.Description)
				}
				return nil
			}

			fmt.Println("Development Phases:")
			fmt.Println()
			for _, p := range core.AllPhases {
				status := state.GetPhaseStatus(p.Name)
				icon := getStatusIcon(status, p.Name == state.CurrentPhase)
				fmt.Printf("  %s %s - %s\n", icon, p.Name, p.Description)
				if p.Name == state.CurrentPhase {
					fmt.Printf("      └─ Current phase (started %s)\n",
						formatDuration(time.Since(state.PhaseStartedAt)))
				}
			}

			return nil
		},
	}
}

func newPhaseStartCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "start <phase>",
		Short: "Start a development phase",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			var completions []string
			for _, p := range core.AllPhases {
				if strings.HasPrefix(p.Name, toComplete) {
					completions = append(completions, p.Name)
				}
			}
			return completions, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			phaseName := args[0]

			// Validate phase name
			if !core.IsValidPhase(phaseName) {
				return fmt.Errorf("invalid phase: %s. Valid phases: %s",
					phaseName, strings.Join(core.PhaseNames(), ", "))
			}

			state, err := core.LoadProjectState(".forge/state.yaml")
			if err != nil {
				return fmt.Errorf("not a forge project (run 'forge init' first)")
			}

			// Check if already in a phase
			if state.CurrentPhase != "" && !force {
				return fmt.Errorf("already in phase '%s'. Use 'forge phase complete' first or --force to switch",
					state.CurrentPhase)
			}

			// Check phase order (unless forcing)
			if !force {
				if err := validatePhaseOrder(state, phaseName); err != nil {
					return err
				}
			}

			// Start the phase
			state.CurrentPhase = phaseName
			state.PhaseStartedAt = time.Now()
			state.SetPhaseStatus(phaseName, "in_progress")
			state.AddActivity(fmt.Sprintf("Started phase: %s", phaseName))

			if err := state.Save(".forge/state.yaml"); err != nil {
				return fmt.Errorf("failed to save state: %w", err)
			}

			phase := core.GetPhase(phaseName)
			fmt.Printf("Started phase: %s\n", phaseName)
			fmt.Printf("Description: %s\n", phase.Description)
			fmt.Printf("Primary tool: %s\n", phase.PrimaryTool)
			fmt.Println("\nCheckpoint criteria:")
			for _, c := range phase.Checkpoint.Criteria {
				fmt.Printf("  • %s\n", c)
			}
			fmt.Println("\nRun 'forge run' to start working with the AI assistant.")

			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "force start even if another phase is active")

	return cmd
}

func newPhaseCompleteCmd() *cobra.Command {
	var skipCheck bool

	cmd := &cobra.Command{
		Use:   "complete",
		Short: "Complete the current phase",
		Long: `Complete the current development phase.

This will run checkpoint validation to ensure phase criteria are met
before marking the phase as complete.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := core.LoadProjectState(".forge/state.yaml")
			if err != nil {
				return fmt.Errorf("not a forge project (run 'forge init' first)")
			}

			if state.CurrentPhase == "" {
				return fmt.Errorf("no active phase to complete")
			}

			_ = core.GetPhase(state.CurrentPhase) // validate phase exists

			// Run checkpoint validation
			if !skipCheck {
				fmt.Printf("Running checkpoint validation for phase: %s\n\n", state.CurrentPhase)

				result := core.ValidateCheckpoint(state.CurrentPhase)

				for _, check := range result.Checks {
					icon := "✓"
					if !check.Passed {
						icon = "✗"
					}
					fmt.Printf("  %s %s\n", icon, check.Name)
					if check.Message != "" {
						fmt.Printf("      %s\n", check.Message)
					}
				}

				if !result.Passed {
					fmt.Println("\nCheckpoint validation failed. Address the issues above or use --skip-check.")
					return fmt.Errorf("checkpoint validation failed")
				}

				fmt.Println("\nAll checkpoints passed!")
			}

			// Save phase history
			historyPath := filepath.Join(".forge", "history", fmt.Sprintf("%s_%s.yaml",
				state.CurrentPhase, time.Now().Format("20060102_150405")))
			history := core.PhaseHistory{
				Phase:       state.CurrentPhase,
				StartedAt:   state.PhaseStartedAt,
				CompletedAt: time.Now(),
				Duration:    time.Since(state.PhaseStartedAt),
			}
			if err := history.Save(historyPath); err != nil {
				fmt.Printf("Warning: failed to save history: %v\n", err)
			}

			// Update state
			state.SetPhaseStatus(state.CurrentPhase, "completed")
			state.AddActivity(fmt.Sprintf("Completed phase: %s", state.CurrentPhase))
			completedPhase := state.CurrentPhase
			state.CurrentPhase = ""
			state.PhaseStartedAt = time.Time{}

			if err := state.Save(".forge/state.yaml"); err != nil {
				return fmt.Errorf("failed to save state: %w", err)
			}

			fmt.Printf("\nPhase '%s' completed successfully!\n", completedPhase)

			// Suggest next phase
			if next := core.NextPhase(completedPhase); next != "" {
				fmt.Printf("\nNext: forge phase start %s\n", next)
			} else {
				fmt.Println("\nAll phases complete! Project ready for release.")
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&skipCheck, "skip-check", false, "skip checkpoint validation")

	return cmd
}

func newPhaseInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info [phase]",
		Short: "Show detailed phase information",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var phaseName string

			if len(args) > 0 {
				phaseName = args[0]
			} else {
				state, err := core.LoadProjectState(".forge/state.yaml")
				if err != nil || state.CurrentPhase == "" {
					return fmt.Errorf("specify a phase or start one with 'forge phase start'")
				}
				phaseName = state.CurrentPhase
			}

			if !core.IsValidPhase(phaseName) {
				return fmt.Errorf("invalid phase: %s", phaseName)
			}

			phase := core.GetPhase(phaseName)

			fmt.Printf("Phase: %s\n", phase.Name)
			fmt.Printf("Description: %s\n", phase.Description)
			fmt.Printf("Primary Tool: %s\n", phase.PrimaryTool)
			fmt.Printf("Tool Reason: %s\n", phase.ToolReason)
			fmt.Println("\nCheckpoint Criteria:")
			for _, c := range phase.Checkpoint.Criteria {
				fmt.Printf("  • %s\n", c)
			}
			fmt.Println("\nArtifacts:")
			for _, a := range phase.Artifacts {
				fmt.Printf("  • %s\n", a)
			}

			// Check artifact directory
			artifactDir := filepath.Join(".forge", "artifacts", phaseName)
			if entries, err := os.ReadDir(artifactDir); err == nil && len(entries) > 0 {
				fmt.Println("\nGenerated Artifacts:")
				for _, e := range entries {
					fmt.Printf("  • %s\n", e.Name())
				}
			}

			return nil
		},
	}
}

func validatePhaseOrder(state *core.ProjectState, targetPhase string) error {
	// Allow starting from discovery without prerequisites
	if targetPhase == "discovery" {
		return nil
	}

	// Check if previous phases are completed
	for _, p := range core.AllPhases {
		if p.Name == targetPhase {
			break
		}
		status := state.GetPhaseStatus(p.Name)
		if status != "completed" {
			return fmt.Errorf("phase '%s' must be completed before starting '%s'. Use --force to override",
				p.Name, targetPhase)
		}
	}

	return nil
}

func getStatusIcon(status string, current bool) string {
	if current {
		return "▶"
	}
	switch status {
	case "completed":
		return "✓"
	case "in_progress":
		return "◐"
	default:
		return "○"
	}
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "just now"
	}
	if d < time.Hour {
		return fmt.Sprintf("%d minutes ago", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%d hours ago", int(d.Hours()))
	}
	return fmt.Sprintf("%d days ago", int(d.Hours()/24))
}
