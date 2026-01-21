package tools

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rice0649/fabric-lite/internal/core"
)

const codexSystemPrompt = `You are an expert AI programming assistant, a meta-tool known as 'Codex'. Your purpose is to help with code planning, review, and execution. You will be given a prompt and context, and you must return a clear, concise, and actionable response, formatted in markdown unless otherwise specified. When generating code, provide complete, runnable snippets. When reviewing code, be specific and provide examples. When planning, break down the task into clear steps.`

// CodexTool is a meta-tool that uses other providers for code-specific tasks.
type CodexTool struct {
	BaseTool
	Config core.CodexConfig // Using the config struct from the core package
}

// NewCodexTool creates a new CodexTool.
func NewCodexTool() *CodexTool {
	return &CodexTool{
		BaseTool: BaseTool{
			name:        "codex",
			description: "The Code Generation Specialist. A meta-tool for writing, refactoring, and explaining code based on a plan.",
		},
	}
}

// IsAvailable checks if the tool is enabled in the config.
func (t *CodexTool) IsAvailable() bool {
	t.loadConfig()
	return t.Config.Enabled
}

// loadConfig loads the codex configuration from project config or uses defaults
func (t *CodexTool) loadConfig() {
	// Try to load from .forge/config.yaml first
	configPath := filepath.Join(".forge", "config.yaml")
	if cfg, err := core.LoadProjectConfig(configPath); err == nil {
		t.Config = cfg.Tools.Codex
		return
	}

	// Try home config directory
	homeDir, _ := os.UserHomeDir()
	if homeDir != "" {
		configPath = filepath.Join(homeDir, ".config", "fabric-lite", "config.yaml")
		if cfg, err := core.LoadProjectConfig(configPath); err == nil {
			t.Config = cfg.Tools.Codex
			return
		}
	}

	// Use sensible defaults if no config found
	t.Config = core.CodexConfig{
		Provider: "ollama",       // Default to local ollama
		Model:    "llama3:latest",
		Enabled:  true,
	}
}

// Execute runs the codex meta-tool. It's headless by nature.
func (t *CodexTool) Execute(ctx ExecutionContext) (*ExecutionResult, error) {
	// Load config if not already loaded
	if t.Config.Provider == "" {
		t.loadConfig()
	}

	if t.Config.Provider == "" {
		return nil, fmt.Errorf("codex tool requires a provider (e.g., ollama, gemini) to be configured")
	}

	// 1. Get the underlying provider tool from the registry.
	providerTool, err := GetTool(t.Config.Provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get codex provider '%s': %w", t.Config.Provider, err)
	}

	// 2. Combine the codex system prompt with the user's prompt.
	fullPrompt := fmt.Sprintf("%s\n\n--- Prompt ---\n\n%s", codexSystemPrompt, ctx.Prompt)

	// 3. Create a new execution context for the underlying tool.
	providerCtx := ExecutionContext{
		Prompt:  fullPrompt,
		WorkDir: ctx.WorkDir,
		Env:     ctx.Env,
		Phase:   ctx.Phase,
		// Pass the model to the provider via Args if specified
		Args: ctx.Args,
	}

	// Add model to args if configured
	if t.Config.Model != "" {
		providerCtx.Args = append([]string{"--model", t.Config.Model}, providerCtx.Args...)
	}

	// 4. Delegate execution to the provider tool.
	result, err := providerTool.Execute(providerCtx)
	if err != nil {
		return nil, fmt.Errorf("codex execution via provider '%s' failed: %w", t.Config.Provider, err)
	}

	return result, nil
}
