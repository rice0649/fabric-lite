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

func newSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session",
		Short: "Manage session state for project continuity",
		Long: `Save and resume session state for seamless project continuity.

The session command helps you pick up where you left off by saving
a comprehensive snapshot of your project state, progress, and next steps.`,
	}

	cmd.AddCommand(newSessionSaveCmd())
	cmd.AddCommand(newSessionResumeCmd())
	cmd.AddCommand(newSessionShowCmd())

	return cmd
}

func newSessionSaveCmd() *cobra.Command {
	var outputPath string
	var includeContext bool

	cmd := &cobra.Command{
		Use:   "save",
		Short: "Save current session state",
		Long: `Save a comprehensive snapshot of the current project state.

This creates a markdown file with:
- Project overview and purpose
- Current phase and progress
- Generated artifacts
- Recent activity log
- Next steps and suggestions
- Key file locations for quick reference`,
		Example: `  # Save to default location (.forge/session.md)
  forge session save

  # Save to custom location
  forge session save -o PROJECT_STATE.md

  # Include extended context for AI assistants
  forge session save --context`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return saveSession(outputPath, includeContext)
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", ".forge/session.md", "output file path")
	cmd.Flags().BoolVar(&includeContext, "context", false, "include extended context for AI assistants")

	return cmd
}

func newSessionResumeCmd() *cobra.Command {
	var inputPath string

	cmd := &cobra.Command{
		Use:   "resume",
		Short: "Display session state for resuming work",
		Long: `Display the saved session state to help resume work.

This command reads the session file and presents a formatted
overview suitable for starting a new coding session.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return resumeSession(inputPath)
		},
	}

	cmd.Flags().StringVarP(&inputPath, "input", "i", ".forge/session.md", "session file path")

	return cmd
}

func newSessionShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current session state without saving",
		Long:  `Display current session state to stdout without saving to a file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			content, err := generateSessionContent(true)
			if err != nil {
				return err
			}
			fmt.Println(content)
			return nil
		},
	}
}

func saveSession(outputPath string, includeContext bool) error {
	content, err := generateSessionContent(includeContext)
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	fmt.Printf("Session saved to: %s\n", outputPath)
	fmt.Println("\nTo resume in a new session, run:")
	fmt.Printf("  forge session resume -i %s\n", outputPath)
	fmt.Println("\nOr copy the file contents as context for your AI assistant.")

	return nil
}

func resumeSession(inputPath string) error {
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read session file: %w", err)
	}

	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("                    SESSION RESUME                              ")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println(string(content))

	return nil
}

func generateSessionContent(includeContext bool) (string, error) {
	var sb strings.Builder

	// Load project config and state
	cfg, err := core.LoadProjectConfig(".forge/config.yaml")
	if err != nil {
		return "", fmt.Errorf("not a forge project (run 'forge init' first)")
	}

	state, err := core.LoadProjectState(".forge/state.yaml")
	if err != nil {
		return "", fmt.Errorf("failed to load project state: %w", err)
	}

	// Header
	sb.WriteString("# Forge Session State\n\n")
	sb.WriteString(fmt.Sprintf("**Generated:** %s\n", time.Now().Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("**Project:** %s\n", cfg.Name))
	if cfg.Description != "" {
		sb.WriteString(fmt.Sprintf("**Description:** %s\n", cfg.Description))
	}
	sb.WriteString("\n---\n\n")

	// Project Overview
	sb.WriteString("## Project Overview\n\n")
	if cfg.Description != "" {
		sb.WriteString(cfg.Description + "\n\n")
	}
	if cfg.Template != "" {
		sb.WriteString(fmt.Sprintf("**Template:** %s\n", cfg.Template))
	}
	sb.WriteString(fmt.Sprintf("**Version:** %s\n\n", cfg.Version))

	// Current State
	sb.WriteString("## Current State\n\n")
	if state.CurrentPhase != "" {
		phase := core.GetPhase(state.CurrentPhase)
		sb.WriteString(fmt.Sprintf("**Active Phase:** %s\n", state.CurrentPhase))
		sb.WriteString(fmt.Sprintf("**Phase Description:** %s\n", phase.Description))
		sb.WriteString(fmt.Sprintf("**Primary Tool:** %s\n", phase.PrimaryTool))
		sb.WriteString(fmt.Sprintf("**Started:** %s\n", state.PhaseStartedAt.Format("2006-01-02 15:04")))
		sb.WriteString(fmt.Sprintf("**Duration:** %s\n\n", formatDurationLong(time.Since(state.PhaseStartedAt))))
	} else {
		sb.WriteString("**Active Phase:** None\n\n")
	}

	// Progress Overview
	sb.WriteString("## Progress\n\n")
	completed := 0
	for _, p := range core.AllPhases {
		status := state.GetPhaseStatus(p.Name)
		icon := "⬜"
		statusText := "pending"
		if status == "completed" {
			icon = "✅"
			statusText = "completed"
			completed++
		} else if status == "in_progress" {
			icon = "🔄"
			statusText = "in progress"
		}
		sb.WriteString(fmt.Sprintf("- %s **%s** - %s\n", icon, p.Name, statusText))
	}
	progress := float64(completed) / float64(len(core.AllPhases)) * 100
	sb.WriteString(fmt.Sprintf("\n**Overall Progress:** %.0f%% (%d/%d phases)\n\n", progress, completed, len(core.AllPhases)))

	// Generated Artifacts
	sb.WriteString("## Generated Artifacts\n\n")
	artifactCount := 0
	for _, p := range core.AllPhases {
		artifactDir := filepath.Join(".forge", "artifacts", p.Name)
		entries, err := os.ReadDir(artifactDir)
		if err != nil || len(entries) == 0 {
			continue
		}
		sb.WriteString(fmt.Sprintf("### %s\n", p.Name))
		for _, e := range entries {
			if !e.IsDir() {
				sb.WriteString(fmt.Sprintf("- `%s`\n", filepath.Join(artifactDir, e.Name())))
				artifactCount++
			}
		}
		sb.WriteString("\n")
	}
	if artifactCount == 0 {
		sb.WriteString("*No artifacts generated yet.*\n\n")
	}

	// Checkpoint Status (if in a phase)
	if state.CurrentPhase != "" {
		sb.WriteString("## Checkpoint Status\n\n")
		result := core.ValidateCheckpoint(state.CurrentPhase)
		passed := 0
		for _, check := range result.Checks {
			icon := "❌"
			if check.Passed {
				icon = "✅"
				passed++
			}
			sb.WriteString(fmt.Sprintf("- %s %s\n", icon, check.Name))
		}
		sb.WriteString(fmt.Sprintf("\n**Checkpoint Progress:** %d/%d criteria met\n\n", passed, len(result.Checks)))
	}

	// Recent Activity
	sb.WriteString("## Recent Activity\n\n")
	limit := 10
	if len(state.Activities) < limit {
		limit = len(state.Activities)
	}
	if limit == 0 {
		sb.WriteString("*No activity recorded yet.*\n\n")
	} else {
		for i := len(state.Activities) - 1; i >= len(state.Activities)-limit; i-- {
			a := state.Activities[i]
			sb.WriteString(fmt.Sprintf("- `%s` %s", a.Timestamp.Format("2006-01-02 15:04"), a.Message))
			if a.Phase != "" {
				sb.WriteString(fmt.Sprintf(" (%s)", a.Phase))
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	// Next Steps
	sb.WriteString("## Suggested Next Steps\n\n")
	nextSteps := generateNextSteps(state)
	for i, step := range nextSteps {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, step))
	}
	sb.WriteString("\n")

	// Key Commands
	sb.WriteString("## Quick Commands\n\n")
	sb.WriteString("```bash\n")
	sb.WriteString("# Check current status\n")
	sb.WriteString("forge status\n\n")
	if state.CurrentPhase == "" {
		sb.WriteString("# Start a phase\n")
		sb.WriteString("forge phase start discovery\n\n")
	} else {
		sb.WriteString("# Run AI tool for current phase\n")
		sb.WriteString("forge run\n\n")
		sb.WriteString("# Complete current phase\n")
		sb.WriteString("forge phase complete\n\n")
	}
	sb.WriteString("# View phase details\n")
	sb.WriteString("forge phase info\n")
	sb.WriteString("```\n\n")

	// Extended context for AI assistants
	if includeContext {
		sb.WriteString("---\n\n")
		sb.WriteString("## AI Assistant Context\n\n")
		sb.WriteString("### Key File Locations\n\n")
		sb.WriteString("```\n")
		sb.WriteString("Project Root: " + mustGetwd() + "\n")
		sb.WriteString("Config:       .forge/config.yaml\n")
		sb.WriteString("State:        .forge/state.yaml\n")
		sb.WriteString("Artifacts:    .forge/artifacts/\n")
		sb.WriteString("History:      .forge/history/\n")
		sb.WriteString("```\n\n")

		sb.WriteString("### Project Structure\n\n")
		sb.WriteString("```\n")
		sb.WriteString(generateProjectTree())
		sb.WriteString("```\n\n")

		sb.WriteString("### Resume Prompt\n\n")
		sb.WriteString("Use this prompt to resume work with an AI assistant:\n\n")
		sb.WriteString("```\n")
		sb.WriteString(generateResumePrompt(cfg, state))
		sb.WriteString("```\n")
	}

	return sb.String(), nil
}

func generateNextSteps(state *core.ProjectState) []string {
	steps := []string{}

	if state.CurrentPhase == "" {
		// No active phase
		for _, p := range core.AllPhases {
			status := state.GetPhaseStatus(p.Name)
			if status != "completed" {
				steps = append(steps, fmt.Sprintf("Start the **%s** phase: `forge phase start %s`", p.Name, p.Name))
				break
			}
		}
		if len(steps) == 0 {
			steps = append(steps, "All phases complete! Project is ready for release.")
		}
	} else {
		// Active phase
		phase := core.GetPhase(state.CurrentPhase)
		steps = append(steps, fmt.Sprintf("Run AI tool for %s: `forge run`", state.CurrentPhase))

		// Check what artifacts are missing
		result := core.ValidateCheckpoint(state.CurrentPhase)
		for _, check := range result.Checks {
			if !check.Passed {
				steps = append(steps, fmt.Sprintf("Create: %s", check.Name))
			}
		}

		if result.Passed {
			steps = append(steps, fmt.Sprintf("Complete the phase: `forge phase complete`"))
			if next := core.NextPhase(state.CurrentPhase); next != "" {
				steps = append(steps, fmt.Sprintf("Then start **%s**: `forge phase start %s`", next, next))
			}
		} else {
			steps = append(steps, fmt.Sprintf("Use %s to generate required artifacts", phase.PrimaryTool))
		}
	}

	return steps
}

func generateProjectTree() string {
	var sb strings.Builder

	// Walk common directories
	dirs := []string{"src", "cmd", "internal", "pkg", "lib", "tests", "docs"}
	for _, dir := range dirs {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			sb.WriteString(dir + "/\n")
			entries, _ := os.ReadDir(dir)
			for i, e := range entries {
				prefix := "├── "
				if i == len(entries)-1 {
					prefix = "└── "
				}
				if e.IsDir() {
					sb.WriteString(prefix + e.Name() + "/\n")
				} else {
					sb.WriteString(prefix + e.Name() + "\n")
				}
			}
		}
	}

	// Always show .forge structure
	sb.WriteString(".forge/\n")
	sb.WriteString("├── config.yaml\n")
	sb.WriteString("├── state.yaml\n")
	sb.WriteString("├── artifacts/\n")
	sb.WriteString("└── history/\n")

	return sb.String()
}

func generateResumePrompt(cfg *core.ProjectConfig, state *core.ProjectState) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("I'm working on '%s'", cfg.Name))
	if cfg.Description != "" {
		sb.WriteString(fmt.Sprintf(" - %s", cfg.Description))
	}
	sb.WriteString(".\n\n")

	if state.CurrentPhase != "" {
		phase := core.GetPhase(state.CurrentPhase)
		sb.WriteString(fmt.Sprintf("Currently in the **%s** phase (%s).\n", state.CurrentPhase, phase.Description))
		sb.WriteString(fmt.Sprintf("Primary tool: %s\n\n", phase.PrimaryTool))
	}

	sb.WriteString("Run `forge status` to see current progress.\n")
	sb.WriteString("Run `forge session show` for full context.\n")

	return sb.String()
}

func formatDurationLong(d time.Duration) string {
	if d < time.Minute {
		return "just started"
	}
	if d < time.Hour {
		mins := int(d.Minutes())
		if mins == 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", mins)
	}
	if d < 24*time.Hour {
		hours := int(d.Hours())
		if hours == 1 {
			return "1 hour"
		}
		return fmt.Sprintf("%d hours", hours)
	}
	days := int(d.Hours() / 24)
	if days == 1 {
		return "1 day"
	}
	return fmt.Sprintf("%d days", days)
}

func mustGetwd() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}
