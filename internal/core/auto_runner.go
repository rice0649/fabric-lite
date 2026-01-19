package core

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// AutoRunner orchestrates automated phase execution
type AutoRunner struct {
	Config    *ProjectConfig
	State     *ProjectState
	StatePath string
	Validator *PhaseValidator
	Executor  PhaseExecutor
	Verbose   bool
}

// PhaseExecutor interface for executing phases (allows mocking in tests)
type PhaseExecutor interface {
	Execute(phase string) error
}

// PhaseValidator validates phase outputs using AI
type PhaseValidator struct {
	ExecuteFunc func(pattern, input string) (string, error)
}

// AIValidationResult holds the result of AI validation
type AIValidationResult struct {
	Valid    bool   `json:"valid"`
	Feedback string `json:"feedback"`
}

// NewAutoRunner creates a new AutoRunner with validation
func NewAutoRunner(config *ProjectConfig, state *ProjectState, statePath string) *AutoRunner {
	return &AutoRunner{
		Config:    config,
		State:     state,
		StatePath: statePath,
	}
}

// Run executes phases automatically from `from` to `until`
func (r *AutoRunner) Run(from, until string, skipValidation bool) error {
	// Validate runner state
	if r == nil || r.State == nil {
		return fmt.Errorf("auto runner not initialized")
	}

	// Initialize auto state if needed
	if r.State.Auto == nil {
		r.State.Auto = &AutoState{
			StartedAt: time.Now(),
		}
	}

	// Get phase range
	phases, err := r.getPhaseRange(from, until)
	if err != nil {
		return fmt.Errorf("get phases: %w", err)
	}
	if len(phases) == 0 {
		return fmt.Errorf("no phases to run")
	}

	// Store configuration in state for resume
	r.State.Auto.FromPhase = from
	r.State.Auto.UntilPhase = until
	r.State.Auto.SkipValidation = skipValidation
	r.State.Auto.StartedAt = time.Now()

	fmt.Printf("Running phases: %s\n", strings.Join(phases, " → "))

	for _, phase := range phases {
		if err := r.runPhase(phase, skipValidation); err != nil {
			return err
		}
	}

	fmt.Println("\n✓ All phases completed successfully!")
	return nil
}

// runPhase executes a single phase with checkpointing and optional validation
func (r *AutoRunner) runPhase(phase string, skipValidation bool) error {
	fmt.Printf("\n─── Phase: %s ───\n", phase)

	// Save state BEFORE execution (crash recovery)
	r.State.CurrentPhase = phase
	r.State.PhaseStartedAt = time.Now()
	r.State.Auto.CurrentPhaseStatus = "running"
	r.State.SetPhaseStatus(phase, "in_progress")
	if err := r.State.Save(r.StatePath); err != nil {
		return fmt.Errorf("save state before %s: %w", phase, err)
	}

	// Execute the phase
	if r.Executor != nil {
		if err := r.Executor.Execute(phase); err != nil {
			r.State.Auto.CurrentPhaseStatus = "failed"
			r.State.Auto.Feedback = err.Error()
			r.State.Save(r.StatePath) // Best effort save
			return fmt.Errorf("phase %s failed: %w", phase, err)
		}
	} else {
		// Default execution: just mark as needing manual run
		fmt.Printf("  → Run: forge run (tool: %s)\n", GetDefaultTool(phase))
		fmt.Printf("  → Artifacts: %v\n", GetPhase(phase).Artifacts)
	}

	// Save state AFTER execution
	r.State.Auto.LastCompletedPhase = phase
	r.State.Auto.CurrentPhaseStatus = "completed"
	r.State.SetPhaseStatus(phase, "completed")
	r.State.AddActivity(fmt.Sprintf("Auto: completed phase %s", phase))
	if err := r.State.Save(r.StatePath); err != nil {
		return fmt.Errorf("save state after %s: %w", phase, err)
	}

	// Validate if enabled
	if !skipValidation && r.Validator != nil {
		fmt.Printf("  → Validating phase output...\n")
		result, err := r.validatePhase(phase)
		if err != nil {
			r.State.Auto.CurrentPhaseStatus = "validation_error"
			r.State.Auto.Feedback = err.Error()
			r.State.Save(r.StatePath)
			return fmt.Errorf("validate %s: %w", phase, err)
		}
		if !result.Valid {
			r.State.Auto.CurrentPhaseStatus = "validation_failed"
			r.State.Auto.Feedback = result.Feedback
			r.State.SetPhaseStatus(phase, "validation_failed")
			if err := r.State.Save(r.StatePath); err != nil {
				return fmt.Errorf("save validation feedback: %w", err)
			}
			return fmt.Errorf("validation failed for %s: %s", phase, result.Feedback)
		}
		fmt.Printf("  ✓ Validation passed\n")
	}

	fmt.Printf("  ✓ Phase %s completed\n", phase)
	return nil
}

// validatePhase runs AI validation on phase output
func (r *AutoRunner) validatePhase(phase string) (*AIValidationResult, error) {
	if r.Validator == nil || r.Validator.ExecuteFunc == nil {
		return &AIValidationResult{Valid: true, Feedback: "Validation skipped (no validator)"}, nil
	}

	phaseInfo := GetPhase(phase)
	if phaseInfo == nil {
		return nil, fmt.Errorf("unknown phase: %s", phase)
	}

	// Build validation input
	input := fmt.Sprintf(`Phase: %s
Description: %s
Expected Artifacts: %v
Checkpoint Criteria: %v

Please validate that this phase has been completed successfully.`,
		phase, phaseInfo.Description, phaseInfo.Artifacts, phaseInfo.Checkpoint.Criteria)

	// Execute validation pattern
	output, err := r.Validator.ExecuteFunc("validation/validate_phase_output", input)
	if err != nil {
		return nil, fmt.Errorf("validation execution failed: %w", err)
	}

	// Parse JSON response
	result, err := parseValidationResult(output)
	if err != nil {
		return nil, fmt.Errorf("parse validation result: %w", err)
	}

	return result, nil
}

// getPhaseRange returns phases between from and until (inclusive)
func (r *AutoRunner) getPhaseRange(from, until string) ([]string, error) {
	allPhases := PhaseNames()

	// Determine start phase
	startIdx := 0
	if from != "" {
		if !IsValidPhase(from) {
			return nil, fmt.Errorf("invalid start phase: %s", from)
		}
		for i, p := range allPhases {
			if p == from {
				startIdx = i
				break
			}
		}
	} else if r.State.Auto != nil && r.State.Auto.LastCompletedPhase != "" {
		// Resume from last completed phase
		for i, p := range allPhases {
			if p == r.State.Auto.LastCompletedPhase {
				startIdx = i + 1 // Start from next phase
				break
			}
		}
		if startIdx >= len(allPhases) {
			return nil, fmt.Errorf("all phases already completed")
		}
	}

	// Determine end phase
	endIdx := len(allPhases) - 1
	if until != "" {
		if !IsValidPhase(until) {
			return nil, fmt.Errorf("invalid end phase: %s", until)
		}
		for i, p := range allPhases {
			if p == until {
				endIdx = i
				break
			}
		}
	}

	// Validate range
	if startIdx > endIdx {
		return nil, fmt.Errorf("start phase '%s' comes after end phase '%s'",
			allPhases[startIdx], allPhases[endIdx])
	}

	return allPhases[startIdx : endIdx+1], nil
}

// GetResumeInfo returns information about resumable auto run
func (r *AutoRunner) GetResumeInfo() (canResume bool, lastPhase string, nextPhase string) {
	if r.State.Auto == nil || r.State.Auto.LastCompletedPhase == "" {
		return false, "", PhaseNames()[0]
	}

	lastPhase = r.State.Auto.LastCompletedPhase
	nextPhase = NextPhase(lastPhase)
	canResume = nextPhase != ""

	return canResume, lastPhase, nextPhase
}

// parseValidationResult extracts AIValidationResult from AI output
func parseValidationResult(output string) (*AIValidationResult, error) {
	// Try to find JSON in the output (may be wrapped in markdown)
	jsonStr := extractJSONFromOutput(output)
	if jsonStr == "" {
		// If no JSON found, treat as invalid with the output as feedback
		return &AIValidationResult{
			Valid:    false,
			Feedback: "Could not parse validation response: " + output,
		}, nil
	}

	var result AIValidationResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("invalid JSON in validation response: %w", err)
	}

	return &result, nil
}

// extractJSONFromOutput extracts JSON from potentially markdown-wrapped output
func extractJSONFromOutput(s string) string {
	// Try to find JSON in code blocks first
	codeBlockRegex := regexp.MustCompile("```(?:json)?\\s*([\\s\\S]*?)```")
	matches := codeBlockRegex.FindStringSubmatch(s)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// Try to find raw JSON object
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "{") {
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
