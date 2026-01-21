package core

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCheckpoint(t *testing.T) {
	checkpoint := Checkpoint{
		Criteria: []string{
			"Requirements document exists",
			"User stories defined",
			"Technical constraints identified",
		},
	}

	if len(checkpoint.Criteria) != 3 {
		t.Errorf("Expected 3 criteria, got %d", len(checkpoint.Criteria))
	}

	if checkpoint.Criteria[0] != "Requirements document exists" {
		t.Errorf("Expected first criterion to be 'Requirements document exists', got '%s'", checkpoint.Criteria[0])
	}
}

func TestEmptyCheckpoint(t *testing.T) {
	checkpoint := Checkpoint{}
	if len(checkpoint.Criteria) != 0 {
		t.Errorf("Expected empty criteria, got %d", len(checkpoint.Criteria))
	}
}

func TestActivityValidation(t *testing.T) {
	// Test activity with all fields
	now := time.Now()
	activity := Activity{
		Timestamp: now,
		Message:   "Test activity",
		Phase:     "planning",
	}

	if activity.Timestamp != now {
		t.Error("Expected timestamp to match")
	}
	if activity.Message != "Test activity" {
		t.Errorf("Expected message 'Test activity', got '%s'", activity.Message)
	}
	if activity.Phase != "planning" {
		t.Errorf("Expected phase 'planning', got '%s'", activity.Phase)
	}
}

func TestPhaseHistoryValidation(t *testing.T) {
	startTime := time.Now()
	endTime := startTime.Add(2 * time.Hour)

	history := &PhaseHistory{
		Phase:       "planning",
		StartedAt:   startTime,
		CompletedAt: endTime,
		Duration:    2 * time.Hour,
		Notes:       "Planning completed successfully",
	}

	if history.Phase != "planning" {
		t.Errorf("Expected phase 'planning', got '%s'", history.Phase)
	}
	if history.StartedAt != startTime {
		t.Error("Expected start time to match")
	}
	if history.CompletedAt != endTime {
		t.Error("Expected completion time to match")
	}
	if history.Duration != 2*time.Hour {
		t.Errorf("Expected duration 2h, got %v", history.Duration)
	}
	if history.Notes != "Planning completed successfully" {
		t.Errorf("Expected notes 'Planning completed successfully', got '%s'", history.Notes)
	}
}

func TestPhaseHistorySaveDetailed(t *testing.T) {
	history := &PhaseHistory{
		Phase:       "testing",
		StartedAt:   time.Now(),
		CompletedAt: time.Now().Add(1 * time.Hour),
		Duration:    1 * time.Hour,
		Notes:       "All tests passed",
	}

	tempFile := filepath.Join(t.TempDir(), "test_history_detailed.yaml")
	err := history.Save(tempFile)
	if err != nil {
		t.Errorf("Expected no error saving phase history, got %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Error("Expected history file to exist")
	}
}

func TestAutoStateValidation(t *testing.T) {
	now := time.Now()
	autoState := &AutoState{
		LastCompletedPhase: "planning",
		CurrentPhaseStatus: "completed",
		FromPhase:          "discovery",
		UntilPhase:         "implementation",
		SkipValidation:     true,
		Feedback:           "All requirements met",
		StartedAt:          now,
	}

	if autoState.LastCompletedPhase != "planning" {
		t.Errorf("Expected LastCompletedPhase 'planning', got '%s'", autoState.LastCompletedPhase)
	}
	if autoState.CurrentPhaseStatus != "completed" {
		t.Errorf("Expected CurrentPhaseStatus 'completed', got '%s'", autoState.CurrentPhaseStatus)
	}
	if autoState.FromPhase != "discovery" {
		t.Errorf("Expected FromPhase 'discovery', got '%s'", autoState.FromPhase)
	}
	if autoState.UntilPhase != "implementation" {
		t.Errorf("Expected UntilPhase 'implementation', got '%s'", autoState.UntilPhase)
	}
	if !autoState.SkipValidation {
		t.Error("Expected SkipValidation to be true")
	}
	if autoState.Feedback != "All requirements met" {
		t.Errorf("Expected Feedback 'All requirements met', got '%s'", autoState.Feedback)
	}
	if autoState.StartedAt != now {
		t.Error("Expected StartedAt to match")
	}
}

func TestAutoStateDefaults(t *testing.T) {
	autoState := &AutoState{}

	// Test zero values
	if autoState.LastCompletedPhase != "" {
		t.Errorf("Expected empty LastCompletedPhase, got '%s'", autoState.LastCompletedPhase)
	}
	if autoState.CurrentPhaseStatus != "" {
		t.Errorf("Expected empty CurrentPhaseStatus, got '%s'", autoState.CurrentPhaseStatus)
	}
	if autoState.SkipValidation {
		t.Error("Expected SkipValidation to be false by default")
	}
	if autoState.Feedback != "" {
		t.Errorf("Expected empty Feedback, got '%s'", autoState.Feedback)
	}
	if !autoState.StartedAt.IsZero() {
		t.Error("Expected StartedAt to be zero by default")
	}
}

func TestPhaseValidation(t *testing.T) {
	phase := Phase{
		Name:        "test-phase",
		Description: "Test phase description",
		PrimaryTool: "test-tool",
		ToolReason:  "Testing purposes",
		Checkpoint: Checkpoint{
			Criteria: []string{"Test criterion 1", "Test criterion 2"},
		},
		Artifacts: []string{"artifact1.md", "artifact2.md"},
	}

	if phase.Name != "test-phase" {
		t.Errorf("Expected name 'test-phase', got '%s'", phase.Name)
	}
	if phase.Description != "Test phase description" {
		t.Errorf("Expected description 'Test phase description', got '%s'", phase.Description)
	}
	if phase.PrimaryTool != "test-tool" {
		t.Errorf("Expected primary tool 'test-tool', got '%s'", phase.PrimaryTool)
	}
	if phase.ToolReason != "Testing purposes" {
		t.Errorf("Expected tool reason 'Testing purposes', got '%s'", phase.ToolReason)
	}
	if len(phase.Checkpoint.Criteria) != 2 {
		t.Errorf("Expected 2 checkpoint criteria, got %d", len(phase.Checkpoint.Criteria))
	}
	if len(phase.Artifacts) != 2 {
		t.Errorf("Expected 2 artifacts, got %d", len(phase.Artifacts))
	}
}

func TestPhaseDefaults(t *testing.T) {
	phase := Phase{}

	if phase.Name != "" {
		t.Errorf("Expected empty name, got '%s'", phase.Name)
	}
	if phase.Description != "" {
		t.Errorf("Expected empty description, got '%s'", phase.Description)
	}
	if phase.PrimaryTool != "" {
		t.Errorf("Expected empty primary tool, got '%s'", phase.PrimaryTool)
	}
	if phase.ToolReason != "" {
		t.Errorf("Expected empty tool reason, got '%s'", phase.ToolReason)
	}
	if len(phase.Checkpoint.Criteria) != 0 {
		t.Errorf("Expected no checkpoint criteria, got %d", len(phase.Checkpoint.Criteria))
	}
	if len(phase.Artifacts) != 0 {
		t.Errorf("Expected no artifacts, got %d", len(phase.Artifacts))
	}
}
