package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	if err := cfg.Save(filepath.Join(forgeDir, "config.yaml")); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Create state.yaml
	state := core.NewProjectState()
	if err := state.Save(filepath.Join(forgeDir, "state.yaml")); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	// Apply template if specified
	if template != "" {
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

	// Get description
	fmt.Print("\nProject description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	fmt.Println()
	return runInitWithDescription(name, template, description)
}

func runInitWithDescription(name, template, description string) error {
	if err := runInit(name, template); err != nil {
		return err
	}

	if description != "" {
		// Update config with description
		cfg, err := core.LoadProjectConfig(".forge/config.yaml")
		if err != nil {
			return err
		}
		cfg.Description = description
		return cfg.Save(".forge/config.yaml")
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
