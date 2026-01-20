package tools

import (
	"fmt"
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
// A more robust check would verify if its configured provider is available.
func (t *CodexTool) IsAvailable() bool {
	// For now, we assume it's available if we can load its config.
	// This will be improved when config loading is implemented.
	return true
}

// Execute runs the codex meta-tool. It's headless by nature.
func (t *CodexTool) Execute(ctx ExecutionContext) (*ExecutionResult, error) {
	// TODO: Replace this with actual config loading once implemented.
	// For now, we'll simulate the config to use Ollama.
	t.Config.Provider = "ollama"
	t.Config.Model = "llama3:latest"

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
		// Pass the model to the provider via Args, a common pattern.
		Args: []string{t.Config.Model},
	}

	// 4. Delegate execution to the provider tool.
	result, err := providerTool.Execute(providerCtx)
	if err != nil {
		return nil, fmt.Errorf("codex execution via provider '%s' failed: %w", t.Config.Provider, err)
	}

	return result, nil
}
