package core

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewProjectState(t *testing.T) {
	state := NewProjectState()

	if state.CurrentPhase != "" {
		t.Errorf("Expected CurrentPhase to be empty, got %s", state.CurrentPhase)
	}

	if state.PhaseStatuses == nil {
		t.Error("Expected PhaseStatuses to be initialized")
	}

	expectedPhases := []string{"discovery", "planning", "design", "implementation", "testing", "deployment"}
	for _, phase := range expectedPhases {
		if status, ok := state.PhaseStatuses[phase]; !ok {
			t.Errorf("Expected phase %s to be initialized", phase)
		} else if status != "pending" {
			t.Errorf("Expected phase %s status to be 'pending', got %s", phase, status)
		}
	}

	if len(state.Activities) != 1 {
		t.Errorf("Expected 1 initial activity, got %d", len(state.Activities))
	}

	if state.Activities[0].Message != "Project initialized" {
		t.Errorf("Expected initial activity message to be 'Project initialized', got %s", state.Activities[0].Message)
	}

	if state.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if state.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

func TestProjectStateGetPhaseStatus(t *testing.T) {
	state := NewProjectState()

	// Test existing phase
	status := state.GetPhaseStatus("discovery")
	if status != "pending" {
		t.Errorf("Expected 'pending', got %s", status)
	}

	// Test non-existent phase
	status = state.GetPhaseStatus("nonexistent")
	if status != "pending" {
		t.Errorf("Expected 'pending' for non-existent phase, got %s", status)
	}

	// Test after setting status
	state.SetPhaseStatus("discovery", "completed")
	status = state.GetPhaseStatus("discovery")
	if status != "completed" {
		t.Errorf("Expected 'completed', got %s", status)
	}
}

func TestProjectStateSetPhaseStatus(t *testing.T) {
	state := &ProjectState{}

	// Test with nil PhaseStatuses
	state.SetPhaseStatus("test", "in_progress")
	if state.PhaseStatuses == nil {
		t.Error("Expected PhaseStatuses to be initialized")
	}
	if state.PhaseStatuses["test"] != "in_progress" {
		t.Errorf("Expected 'in_progress', got %s", state.PhaseStatuses["test"])
	}

	// Test updating existing status
	state.SetPhaseStatus("test", "completed")
	if state.PhaseStatuses["test"] != "completed" {
		t.Errorf("Expected 'completed', got %s", state.PhaseStatuses["test"])
	}

	// Test UpdatedAt is set
	oldUpdatedAt := state.UpdatedAt
	time.Sleep(1 * time.Millisecond)
	state.SetPhaseStatus("test2", "pending")
	if !state.UpdatedAt.After(oldUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
}

func TestProjectStateAddActivity(t *testing.T) {
	state := NewProjectState()
	initialCount := len(state.Activities)

	state.CurrentPhase = "planning"
	state.AddActivity("Test activity")

	if len(state.Activities) != initialCount+1 {
		t.Errorf("Expected %d activities, got %d", initialCount+1, len(state.Activities))
	}

	newActivity := state.Activities[len(state.Activities)-1]
	if newActivity.Message != "Test activity" {
		t.Errorf("Expected message 'Test activity', got %s", newActivity.Message)
	}
	if newActivity.Phase != "planning" {
		t.Errorf("Expected phase 'planning', got %s", newActivity.Phase)
	}
	if newActivity.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}
}

func TestProjectStateSaveAndLoad(t *testing.T) {
	// Create test state
	state := NewProjectState()
	state.CurrentPhase = "planning"
	state.SetPhaseStatus("discovery", "completed")
	state.AddActivity("Test activity")

	// Save to temporary file
	tempFile := filepath.Join(t.TempDir(), "test_state.yaml")
	err := state.Save(tempFile)
	if err != nil {
		t.Errorf("Expected no error saving state, got %v", err)
	}

	// Load from file
	loadedState, err := LoadProjectState(tempFile)
	if err != nil {
		t.Errorf("Expected no error loading state, got %v", err)
	}

	// Verify loaded state
	if loadedState.CurrentPhase != state.CurrentPhase {
		t.Errorf("Expected CurrentPhase %s, got %s", state.CurrentPhase, loadedState.CurrentPhase)
	}

	if loadedState.GetPhaseStatus("discovery") != "completed" {
		t.Errorf("Expected discovery status to be 'completed', got %s", loadedState.GetPhaseStatus("discovery"))
	}

	if len(loadedState.Activities) != len(state.Activities) {
		t.Errorf("Expected %d activities, got %d", len(state.Activities), len(loadedState.Activities))
	}
}

func TestLoadProjectStateNonExistent(t *testing.T) {
	_, err := LoadProjectState("/nonexistent/file.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestLoadProjectStateWithNilPhaseStatuses(t *testing.T) {
	// Create a file with no PhaseStatuses
	tempFile := filepath.Join(t.TempDir(), "test_state.yaml")
	content := `
current_phase: "planning"
activities: []
created_at: "2023-01-01T00:00:00Z"
updated_at: "2023-01-01T00:00:00Z"
`
	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	state, err := LoadProjectState(tempFile)
	if err != nil {
		t.Errorf("Expected no error loading state, got %v", err)
	}

	if state.PhaseStatuses == nil {
		t.Error("Expected PhaseStatuses to be initialized")
	}
}

func TestPhaseHistorySave(t *testing.T) {
	history := &PhaseHistory{
		Phase:       "planning",
		StartedAt:   time.Now(),
		CompletedAt: time.Now().Add(1 * time.Hour),
		Duration:    1 * time.Hour,
		Notes:       "Test notes",
	}

	tempFile := filepath.Join(t.TempDir(), "test_history.yaml")
	err := history.Save(tempFile)
	if err != nil {
		t.Errorf("Expected no error saving history, got %v", err)
	}

	// Verify file exists and has content
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Error("Expected history file to exist")
	}
}

func TestAutoState(t *testing.T) {
	autoState := &AutoState{
		LastCompletedPhase: "planning",
		CurrentPhaseStatus: "running",
		FromPhase:          "discovery",
		UntilPhase:         "implementation",
		SkipValidation:     true,
		Feedback:           "Test feedback",
		StartedAt:          time.Now(),
	}

	// Test that all fields are set correctly
	if autoState.LastCompletedPhase != "planning" {
		t.Errorf("Expected LastCompletedPhase to be 'planning', got %s", autoState.LastCompletedPhase)
	}
	if autoState.CurrentPhaseStatus != "running" {
		t.Errorf("Expected CurrentPhaseStatus to be 'running', got %s", autoState.CurrentPhaseStatus)
	}
	if !autoState.SkipValidation {
		t.Error("Expected SkipValidation to be true")
	}
	if autoState.Feedback != "Test feedback" {
		t.Errorf("Expected Feedback to be 'Test feedback', got %s", autoState.Feedback)
	}
}

func TestActivity(t *testing.T) {
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
		t.Errorf("Expected message to be 'Test activity', got %s", activity.Message)
	}
	if activity.Phase != "planning" {
		t.Errorf("Expected phase to be 'planning', got %s", activity.Phase)
	}
}
