package executor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rice0649/fabric-lite/internal/providers"
)

type PatternInfo struct {
	Name        string
	Description string
	System      string
	User        string
}

type PatternExecutor struct {
	providers   map[string]providers.Provider
	patternsDir string
}

func NewPatternExecutor() *PatternExecutor {
	return &PatternExecutor{
		patternsDir: getPatternsDir(),
		providers:   make(map[string]providers.Provider),
	}
}

func (e *PatternExecutor) GetPatternsDir() string {
	return getPatternsDir()
}

func getPatternsDir() string {
	// Check local patterns first
	if _, err := os.Stat("patterns"); err == nil {
		return "patterns"
	}
	// Fall back to user config dir
	return filepath.Join(os.Getenv("HOME"), ".config", "fabric-lite", "patterns")
}

// LoadProvider loads and initializes a provider from configuration
func (e *PatternExecutor) LoadProvider(name string, config *providers.Config) error {
	if config == nil {
		return fmt.Errorf("no config provided")
	}

	// Find provider config
	var targetConfig *providers.ProviderConfig
	for _, provider := range config.Providers {
		if provider.Name == name {
			targetConfig = &provider
			break
		}
	}

	if targetConfig == nil {
		return fmt.Errorf("provider not found in config: %s", name)
	}

	provider, err := providers.NewProvider(*targetConfig)
	if err != nil {
		return fmt.Errorf("failed to create provider %s: %w", name, err)
	}

	e.providers[name] = provider
	return nil
}

// LoadProviderDirect directly loads a pre-created provider
func (e *PatternExecutor) LoadProviderDirect(name string, provider providers.Provider) {
	e.providers[name] = provider
}

func (e *PatternExecutor) ListPatterns() ([]PatternInfo, error) {
	patternsDir := e.patternsDir
	entries, err := os.ReadDir(patternsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read patterns directory: %w", err)
	}

	var patterns []PatternInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		patternName := entry.Name()
		pattern, err := e.loadPattern(patternName)
		if err != nil {
			continue // Skip invalid patterns
		}

		patterns = append(patterns, *pattern)
	}

	return patterns, nil
}

func (e *PatternExecutor) loadPattern(name string) (*PatternInfo, error) {
	patternDir := filepath.Join(e.patternsDir, name)

	// Load system prompt
	systemFile := filepath.Join(patternDir, "system.md")
	systemContent, err := os.ReadFile(systemFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read system prompt: %w", err)
	}

	// Load user prompt (optional)
	userFile := filepath.Join(patternDir, "user.md")
	userData, err := os.ReadFile(userFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read user prompt: %w", err)
	}

	var userContent string
	if err == nil {
		userContent = string(userData)
	}

	return &PatternInfo{
		Name:        name,
		Description: extractDescription(string(systemContent)),
		System:      string(systemContent),
		User:        userContent,
	}, nil
}

func extractDescription(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, "# ") && strings.Contains(line, "DESCRIPTION") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# DESCRIPTION"))
		}
		if strings.Contains(line, "# IDENTITY and PURPOSE") {
			// Take the first line after this
			parts := strings.Split(content, "\n")
			for i := 0; i < len(parts); i++ {
				if strings.Contains(parts[i], "You are") {
					// Extract the role description
					roleLine := strings.TrimSpace(parts[i])
					return roleLine
				}
			}
		}
	}
	return "Pattern execution"
}

func (e *PatternExecutor) Execute(ctx context.Context, patternName, input, providerName string) (*providers.CompletionResponse, error) {
	// Check if provider is loaded
	provider, ok := e.providers[providerName]
	if !ok {
		return nil, fmt.Errorf("provider not loaded: %s", providerName)
	}
	if !provider.IsAvailable() {
		return nil, fmt.Errorf("provider not available: %s", providerName)
	}

	// Load pattern
	pattern, err := e.loadPattern(patternName)
	if err != nil {
		return nil, fmt.Errorf("failed to load pattern %s: %w", patternName, err)
	}

	// Build full prompt
	fullPrompt := pattern.User
	if fullPrompt != "" {
		fullPrompt = pattern.User + "\n\nInput:\n" + input
	} else {
		fullPrompt = input
	}

	// Execute with provider
	request := providers.CompletionRequest{
		System:    pattern.System,
		Prompt:    fullPrompt,
		Model:     "gpt-4o-mini", // Default model for now
		Stream:    false,         // TODO: Add streaming support
		MaxTokens: 4096,
	}

	return provider.Execute(ctx, request)
}
