package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/rice0649/fabric-lite/internal/tools"
)

// ScaffoldContext contains all information needed for AI scaffold generation
type ScaffoldContext struct {
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	Template        string                 `json:"template"`
	TemplateOptions map[string]interface{} `json:"template_options"`
}

// ScaffoldOutput represents the AI-generated scaffold structure
type ScaffoldOutput struct {
	Directories []string       `json:"directories"`
	Files       []ScaffoldFile `json:"files"`
	Commands    []string       `json:"commands"`
}

// ScaffoldFile represents a file to be created with its content
type ScaffoldFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// generateDynamicScaffold uses AI to generate project scaffolding
func generateDynamicScaffold(ctx ScaffoldContext) (*ScaffoldOutput, error) {
	fabricTool := tools.NewFabricTool()

	if !fabricTool.IsAvailable() {
		return nil, fmt.Errorf("fabric-lite not available")
	}

	// Serialize context to JSON for the AI
	contextJSON, err := json.MarshalIndent(ctx, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to serialize context: %w", err)
	}

	// Create execution context
	execCtx := tools.ExecutionContext{
		Pattern: "init/scaffold_project",
		Prompt:  string(contextJSON),
	}

	// Use context with timeout
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Execute in a goroutine with timeout
	resultChan := make(chan *tools.ExecutionResult, 1)
	errChan := make(chan error, 1)

	go func() {
		result, err := fabricTool.ExecuteNonInteractive(execCtx)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- result
	}()

	var result *tools.ExecutionResult
	select {
	case <-timeoutCtx.Done():
		return nil, fmt.Errorf("scaffold generation timed out")
	case err := <-errChan:
		return nil, err
	case result = <-resultChan:
	}

	if !result.Success {
		return nil, fmt.Errorf("scaffold generation failed: %s", result.Error)
	}

	// Extract JSON from output (may be wrapped in markdown code blocks)
	jsonStr := extractJSON(result.Output)
	if jsonStr == "" {
		return nil, fmt.Errorf("no valid JSON found in output")
	}

	// Parse the output
	var output ScaffoldOutput
	if err := json.Unmarshal([]byte(jsonStr), &output); err != nil {
		return nil, fmt.Errorf("failed to parse scaffold output: %w", err)
	}

	return &output, nil
}

// extractJSON extracts JSON from a string that may contain markdown code blocks
func extractJSON(s string) string {
	// Try to find JSON in code blocks first
	codeBlockRegex := regexp.MustCompile("```(?:json)?\\s*([\\s\\S]*?)```")
	matches := codeBlockRegex.FindStringSubmatch(s)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// Try to find raw JSON (starts with {)
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "{") {
		// Find the matching closing brace
		depth := 0
		for i, c := range s {
			if c == '{' {
				depth++
			} else if c == '}' {
				depth--
				if depth == 0 {
					return s[:i+1]
				}
			}
		}
	}

	return ""
}

// createFilesFromScaffold creates directories and files from the scaffold output
func createFilesFromScaffold(output *ScaffoldOutput) error {
	createdDirs := 0
	createdFiles := 0

	// Create directories
	for _, dir := range output.Directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		createdDirs++
	}

	// Create files
	for _, file := range output.Files {
		// Ensure parent directory exists
		dir := filepath.Dir(file.Path)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create parent directory for %s: %w", file.Path, err)
			}
		}

		if err := os.WriteFile(file.Path, []byte(file.Content), 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %w", file.Path, err)
		}
		createdFiles++
	}

	fmt.Printf("Created %d directories and %d files\n", createdDirs, createdFiles)

	// Print created file summary
	if len(output.Files) > 0 {
		fmt.Println("\nGenerated files:")
		for _, file := range output.Files {
			fmt.Printf("  - %s\n", file.Path)
		}
	}

	// Print suggested commands
	if len(output.Commands) > 0 {
		fmt.Println("\nSuggested setup commands:")
		for _, cmd := range output.Commands {
			fmt.Printf("  $ %s\n", cmd)
		}
	}

	return nil
}

// scaffoldWithFallback attempts AI scaffold generation with fallback to static templates
func scaffoldWithFallback(ctx ScaffoldContext) error {
	fmt.Println("\nGenerating project scaffold with AI...")

	output, err := generateDynamicScaffold(ctx)
	if err != nil {
		fmt.Printf("\nWarning: AI scaffold generation failed (%v)\n", err)
		fmt.Println("Falling back to static template...")

		if ctx.Template != "" {
			return applyTemplate(ctx.Template)
		}
		return nil
	}

	return createFilesFromScaffold(output)
}
