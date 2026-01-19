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

func newStatusCmd() *cobra.Command {
	var detailed bool

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show project status dashboard",
		Long: `Display a dashboard view of the project status.

Shows current phase, progress, checkpoint status, and suggested next steps.
Use --detailed for additional information including artifact list and full activity log.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return showStatus(detailed)
		},
	}

	cmd.Flags().BoolVarP(&detailed, "detailed", "d", false, "show detailed status with checkpoints and suggestions")

	return cmd
}

func newHistoryCmd() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "history",
		Short: "View activity history",
		RunE: func(cmd *cobra.Command, args []string) error {
			return showHistory(limit)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "n", 10, "number of entries to show")

	return cmd
}

func showStatus(detailed bool) error {
	cfg, err := core.LoadProjectConfig(".forge/config.yaml")
	if err != nil {
		return fmt.Errorf("not a forge project (run 'forge init' first)")
	}

	state, err := core.LoadProjectState(".forge/state.yaml")
	if err != nil {
		return fmt.Errorf("failed to load project state: %w", err)
	}

	// Header
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Printf("║  AI Project Forge - %s\n", padRight(cfg.Name, 38)+"║")
	fmt.Println("╠════════════════════════════════════════════════════════════╣")

	// Current Phase
	if state.CurrentPhase != "" {
		phase := core.GetPhase(state.CurrentPhase)
		fmt.Printf("║  Current Phase: %-42s ║\n", state.CurrentPhase)
		fmt.Printf("║  Primary Tool:  %-42s ║\n", phase.PrimaryTool)
		fmt.Printf("║  Duration:      %-42s ║\n", formatDuration(time.Since(state.PhaseStartedAt)))
	} else {
		fmt.Printf("║  %-58s ║\n", "No active phase")
	}

	fmt.Println("╠════════════════════════════════════════════════════════════╣")

	// Phase Progress
	fmt.Println("║  Progress:                                                 ║")
	completed := 0
	for _, p := range core.AllPhases {
		status := state.GetPhaseStatus(p.Name)
		icon := getStatusIcon(status, p.Name == state.CurrentPhase)
		if status == "completed" {
			completed++
		}
		fmt.Printf("║    %s %-54s ║\n", icon, p.Name)
	}

	fmt.Println("╠════════════════════════════════════════════════════════════╣")

	// Progress bar
	progress := float64(completed) / float64(len(core.AllPhases)) * 100
	barLen := 50
	filled := int(float64(barLen) * float64(completed) / float64(len(core.AllPhases)))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barLen-filled)
	fmt.Printf("║  [%s] %3.0f%% ║\n", bar, progress)

	fmt.Println("╠════════════════════════════════════════════════════════════╣")

	// Artifacts count
	artifactCount := countArtifacts()
	fmt.Printf("║  Artifacts Generated: %-36d ║\n", artifactCount)

	// Checkpoint status (if in a phase)
	if state.CurrentPhase != "" {
		fmt.Println("╠════════════════════════════════════════════════════════════╣")
		fmt.Println("║  Checkpoint Status:                                        ║")
		result := core.ValidateCheckpoint(state.CurrentPhase)
		passed := 0
		for _, check := range result.Checks {
			icon := "✗"
			if check.Passed {
				icon = "✓"
				passed++
			}
			name := check.Name
			if len(name) > 52 {
				name = name[:49] + "..."
			}
			fmt.Printf("║    %s %-54s ║\n", icon, name)
		}
		fmt.Printf("║    Progress: %d/%d criteria met%-28s║\n", passed, len(result.Checks), "")
	}

	// Recent activity
	if len(state.Activities) > 0 {
		fmt.Println("╠════════════════════════════════════════════════════════════╣")
		fmt.Println("║  Recent Activity:                                          ║")
		limit := 3
		if len(state.Activities) < limit {
			limit = len(state.Activities)
		}
		for i := len(state.Activities) - 1; i >= len(state.Activities)-limit; i-- {
			a := state.Activities[i]
			timeStr := a.Timestamp.Format("15:04")
			line := fmt.Sprintf("%s %s", timeStr, a.Message)
			if len(line) > 56 {
				line = line[:53] + "..."
			}
			fmt.Printf("║    %-56s ║\n", line)
		}
	}

	fmt.Println("╠════════════════════════════════════════════════════════════╣")

	// Next Steps
	fmt.Println("║  Suggested Next Steps:                                     ║")
	nextSteps := getNextSteps(state)
	for _, step := range nextSteps {
		if len(step) > 56 {
			step = step[:53] + "..."
		}
		fmt.Printf("║    → %-54s║\n", step)
	}

	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	// Detailed view
	if detailed {
		showDetailedStatus(state)
	}

	return nil
}

func getNextSteps(state *core.ProjectState) []string {
	steps := []string{}

	if state.CurrentPhase == "" {
		// No active phase - find next uncompleted phase
		for _, p := range core.AllPhases {
			status := state.GetPhaseStatus(p.Name)
			if status != "completed" {
				steps = append(steps, fmt.Sprintf("forge phase start %s", p.Name))
				break
			}
		}
		if len(steps) == 0 {
			steps = append(steps, "All phases complete!")
		}
	} else {
		steps = append(steps, "forge run")
		result := core.ValidateCheckpoint(state.CurrentPhase)
		if result.Passed {
			steps = append(steps, "forge phase complete")
		}
	}

	return steps
}

func showDetailedStatus(state *core.ProjectState) {
	fmt.Println()
	fmt.Println("═══ Detailed Information ═══")
	fmt.Println()

	// List artifacts
	fmt.Println("Artifacts:")
	hasArtifacts := false
	for _, p := range core.AllPhases {
		artifactDir := filepath.Join(".forge", "artifacts", p.Name)
		entries, err := os.ReadDir(artifactDir)
		if err != nil || len(entries) == 0 {
			continue
		}
		hasArtifacts = true
		fmt.Printf("  %s/\n", p.Name)
		for _, e := range entries {
			if !e.IsDir() {
				fmt.Printf("    - %s\n", e.Name())
			}
		}
	}
	if !hasArtifacts {
		fmt.Println("  (none)")
	}

	fmt.Println()
	fmt.Println("Tip: Run 'forge session save' to save full state for resuming later.")
}

func showHistory(limit int) error {
	state, err := core.LoadProjectState(".forge/state.yaml")
	if err != nil {
		return fmt.Errorf("not a forge project (run 'forge init' first)")
	}

	if len(state.Activities) == 0 {
		fmt.Println("No activity history yet.")
		return nil
	}

	fmt.Println("Activity History:")
	fmt.Println()

	start := 0
	if len(state.Activities) > limit {
		start = len(state.Activities) - limit
	}

	for i := len(state.Activities) - 1; i >= start; i-- {
		a := state.Activities[i]
		fmt.Printf("  %s  %s\n", a.Timestamp.Format("2006-01-02 15:04:05"), a.Message)
	}

	return nil
}

func countArtifacts() int {
	count := 0
	artifactDir := ".forge/artifacts"
	filepath.Walk(artifactDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	return count
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s[:n]
	}
	return s + strings.Repeat(" ", n-len(s))
}
