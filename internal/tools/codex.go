package tools

import (
	"context"
	"fmt"

	"github.com/rice0649/fabric-lite/internal/core"
	"github.com/rice0649/fabric-lite/internal/providers"
)

const codexSystemPrompt = `You are an expert AI programming assistant, a meta-tool known as 'Codex'. Your purpose is to help with code planning, review, and execution. You will be given a prompt and context, and you must return a clear, concise, and actionable response, formatted in markdown unless otherwise specified. When generating code, provide complete, runnable snippets. When reviewing code, be specific and provide examples. When planning, break down the task into clear steps.`

// CodexTool is a meta-tool that uses other providers for code-specific tasks.
type CodexTool struct {
	BaseTool
	config          core.CodexConfig
	providerManager *core.ProviderManager
}

// NewCodexTool creates a new CodexTool.
func NewCodexTool(config core.CodexConfig, pm *core.ProviderManager) *CodexTool {
	return &CodexTool{
		BaseTool: BaseTool{
			name:        "codex",
			description: "The Code Generation Specialist. A meta-tool for writing, refactoring, and explaining code based on a plan.",
		},
		config:          config,
		providerManager: pm,
	}
}

// IsAvailable checks if the tool is enabled and its configured provider is available.
func (t *CodexTool) IsAvailable() bool {
	if !t.config.Enabled {
		return false
	}
	if t.config.Provider == "" {
		return false
	}
	return t.providerManager.CheckAvailability(t.config.Provider)
}

// Execute runs the codex meta-tool. It uses an LLM provider.
func (t *CodexTool) Execute(ctx ExecutionContext) (*ExecutionResult, error) {
	if !t.IsAvailable() {
		return nil, fmt.Errorf("codex tool is not available (either not enabled or provider '%s' is not ready)", t.config.Provider)
	}

	// Combine the codex system prompt with the user's prompt.
	fullPrompt := fmt.Sprintf("%s\n\n--- Prompt ---\n\n%s", codexSystemPrompt, ctx.Prompt)

	// Create a CompletionRequest for the LLM provider
	completionRequest := providers.CompletionRequest{
		System: fullPrompt, // Using System for the main prompt as Codex is the system
		Model:  t.config.Model,
		Prompt: "", // User's prompt is integrated into System
	}

	// Delegate execution to the LLM provider
	resp, err := t.providerManager.Execute(context.Background(), t.config.Provider, completionRequest)
	if err != nil {
		return nil, fmt.Errorf("codex execution via provider '%s' failed: %w", t.config.Provider, err)
	}

	return &ExecutionResult{
		Output:   resp.Content,
		Error:    "",
		ExitCode: 0,
		Success:  true,
	}, nil
}
