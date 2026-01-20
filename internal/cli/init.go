package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rice0649/fabric-lite/internal/core"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	var (
		name        string
		template    string
		interactive bool
	)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new forge project",
		Long: `Initialize a new project with forge configuration.

This creates a .forge directory with project configuration and state tracking.
Optionally use --template to scaffold from a project template.`,
		Example: `  # Initialize in current directory
  forge init --name myapp

  # Initialize with a template
  forge init --name myapi --template api

  # Interactive mode
  forge init --interactive`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if interactive {
				return runInteractiveInit()
			}
			return runInit(name, template)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "project name")
	cmd.Flags().StringVarP(&template, "template", "t", "", "project template (webapp, cli, api, library)")
	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "interactive initialization")

	return cmd
}

func runInit(name, template string) error {
	return runInitWithOptions(name, template, "", false)
}

func runInitWithOptions(name, template, description string, skipTemplate bool) error {
	if name == "" {
		// Use current directory name
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		name = filepath.Base(cwd)
	}

	fmt.Printf("Initializing forge project: %s\n", name)

	// Create .forge directory structure
	forgeDir := ".forge"
	dirs := []string{
		forgeDir,
		filepath.Join(forgeDir, "history"),
		filepath.Join(forgeDir, "artifacts"),
		filepath.Join(forgeDir, "artifacts", "discovery"),
		filepath.Join(forgeDir, "artifacts", "planning"),
		filepath.Join(forgeDir, "artifacts", "design"),
		filepath.Join(forgeDir, "artifacts", "implementation"),
		filepath.Join(forgeDir, "artifacts", "testing"),
		filepath.Join(forgeDir, "artifacts", "deployment"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create config.yaml
	cfg := core.NewProjectConfig(name, template)
	if description != "" {
		cfg.Description = description
	}
	if err := cfg.Save(filepath.Join(forgeDir, "config.yaml")); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Create state.yaml
	state := core.NewProjectState()
	if err := state.Save(filepath.Join(forgeDir, "state.yaml")); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	// Apply template if specified and not already done via AI
	if template != "" && !skipTemplate {
		if err := applyTemplate(template); err != nil {
			return fmt.Errorf("failed to apply template: %w", err)
		}
	}

	fmt.Println("\nProject initialized successfully!")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. forge phase start discovery")
	fmt.Println("  2. forge run")
	fmt.Println("  3. forge phase complete")

	return nil
}

func runInteractiveInit() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== AI Project Forge - Interactive Setup ===")
	fmt.Println()

	// Check for existing project
	if _, err := os.Stat(".forge"); err == nil {
		resume, err := offerResumeOrNew(reader)
		if err != nil {
			return err
		}
		if resume {
			return handleResume()
		}
		// Confirm overwrite
		fmt.Print("This will overwrite existing project. Continue? [y/N]: ")
		confirm, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
			return fmt.Errorf("initialization cancelled")
		}
	}

	// Get project name
	fmt.Print("Project name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		cwd, _ := os.Getwd()
		name = filepath.Base(cwd)
	}

	// Get template
	fmt.Println("\nAvailable templates:")
	fmt.Println("  1. webapp  - Web application (frontend + backend)")
	fmt.Println("  2. cli     - Command-line tool")
	fmt.Println("  3. api     - REST API service")
	fmt.Println("  4. library - Reusable library/package")
	fmt.Println("  5. none    - Empty project")
	fmt.Print("\nSelect template [1-5]: ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	templates := map[string]string{
		"1": "webapp",
		"2": "cli",
		"3": "api",
		"4": "library",
		"5": "",
	}
	template := templates[choice]

	// Ask template-specific questions
	var templateOpts *TemplateOptions
	if template != "" {
		var err error
		templateOpts, err = askTemplateQuestions(reader, template)
		if err != nil {
			return err
		}
	}

	// Get description
	fmt.Print("\nProject description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	fmt.Println()

	// Generate scaffold with AI (or fallback)
	aiScaffoldDone := false
	if template != "" && templateOpts != nil {
		ctx := ScaffoldContext{
			Name:            name,
			Description:     description,
			Template:        template,
			TemplateOptions: templateOpts.ToMap(template),
		}
		if err := scaffoldWithFallback(ctx); err != nil {
			// scaffoldWithFallback handles its own fallback, so this is a real error
			return err
		}
		aiScaffoldDone = true
	}

	// Initialize .forge directory (skip static template if AI scaffold succeeded)
	return runInitWithOptions(name, template, description, aiScaffoldDone)
}

// offerResumeOrNew detects existing project and offers options
func offerResumeOrNew(reader *bufio.Reader) (bool, error) {
	// Load existing config
	cfg, err := core.LoadProjectConfig(".forge/config.yaml")
	if err != nil {
		return false, nil // No valid config, proceed with new
	}

	state, _ := core.LoadProjectState(".forge/state.yaml")

	fmt.Println("Existing project detected!")
	fmt.Printf("  Name: %s\n", cfg.Name)
	if cfg.Description != "" {
		fmt.Printf("  Description: %s\n", cfg.Description)
	}
	if cfg.Template != "" {
		fmt.Printf("  Template: %s\n", cfg.Template)
	}
	if state != nil && state.CurrentPhase != "" {
		fmt.Printf("  Current phase: %s\n", state.CurrentPhase)
	}
	fmt.Println()

	fmt.Println("Options:")
	fmt.Println("  1. Resume existing project")
	fmt.Println("  2. Create new project (overwrite)")
	fmt.Println("  3. Cancel")
	fmt.Print("\nSelect [1-3]: ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		return true, nil
	case "2":
		return false, nil
	default:
		return false, fmt.Errorf("initialization cancelled")
	}
}

// handleResume loads existing project and suggests next steps
func handleResume() error {
	cfg, err := core.LoadProjectConfig(".forge/config.yaml")
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}

	state, err := core.LoadProjectState(".forge/state.yaml")
	if err != nil {
		return fmt.Errorf("failed to load project state: %w", err)
	}

	fmt.Println("\n=== Resuming Project ===")
	fmt.Printf("Project: %s\n", cfg.Name)
	if cfg.Description != "" {
		fmt.Printf("Description: %s\n", cfg.Description)
	}

	// Show phase status
	fmt.Println("\nPhase Status:")
	phases := []string{"discovery", "planning", "design", "implementation", "testing", "deployment"}
	for _, phase := range phases {
		status := state.GetPhaseStatus(phase)
		marker := "  "
		if phase == state.CurrentPhase {
			marker = "> "
		}
		fmt.Printf("%s%-15s [%s]\n", marker, phase, status)
	}

	// Show recent activity
	if len(state.Activities) > 0 {
		fmt.Println("\nRecent Activity:")
		start := len(state.Activities) - 5
		if start < 0 {
			start = 0
		}
		for _, activity := range state.Activities[start:] {
			fmt.Printf("  - %s: %s\n", activity.Timestamp.Format(time.RFC822), activity.Message)
		}
	}

	// Suggest next steps
	fmt.Println("\nSuggested next steps:")
	if state.CurrentPhase == "" {
		fmt.Println("  1. forge phase start discovery")
		fmt.Println("  2. forge run")
	} else {
		fmt.Printf("  1. forge run (continue %s phase)\n", state.CurrentPhase)
		fmt.Println("  2. forge phase complete (if ready)")
	}

	return nil
}

func applyTemplate(template string) error {
	// Template application will create standard directories
	switch template {
	case "webapp":
		dirs := []string{"src", "src/components", "src/pages", "public", "tests"}
		for _, d := range dirs {
			os.MkdirAll(d, 0755)
		}
	case "cli":
		dirs := []string{"cmd", "internal", "pkg", "tests"}
		for _, d := range dirs {
			os.MkdirAll(d, 0755)
		}
	case "api":
		dirs := []string{"cmd", "internal/handlers", "internal/models", "internal/services", "tests"}
		for _, d := range dirs {
			os.MkdirAll(d, 0755)
		}
	case "library":
		dirs := []string{"src", "tests", "examples"}
		for _, d := range dirs {
			os.MkdirAll(d, 0755)
		}
	default:
		return fmt.Errorf("unknown template: %s", template)
	}

	fmt.Printf("Applied template: %s\n", template)
	return nil
}
