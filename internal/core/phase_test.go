package core

import (
	"testing"
)

func TestAllPhases(t *testing.T) {
	if len(AllPhases) != 6 {
		t.Errorf("Expected 6 phases, got %d", len(AllPhases))
	}

	expectedPhases := []string{"discovery", "planning", "design", "implementation", "testing", "deployment"}
	actualPhases := PhaseNames()

	if len(actualPhases) != len(expectedPhases) {
		t.Errorf("Expected %d phase names, got %d", len(expectedPhases), len(actualPhases))
	}

	for i, expected := range expectedPhases {
		if i >= len(actualPhases) || actualPhases[i] != expected {
			t.Errorf("Expected phase %d to be '%s', got '%s'", i, getPhaseAtTest(actualPhases, i), actualPhases[i])
		}
	}
}

func TestGetPhase(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		found    bool
	}{
		{"discovery", "discovery", true},
		{"planning", "planning", true},
		{"implementation", "implementation", true},
		{"nonexistent", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phase := GetPhase(tt.name)
			if tt.found {
				if phase == nil {
					t.Errorf("Expected to find phase '%s', got nil", tt.name)
				} else if phase.Name != tt.expected {
					t.Errorf("Expected phase name '%s', got '%s'", tt.expected, phase.Name)
				}
			} else {
				if phase != nil {
					t.Errorf("Expected not to find phase '%s', got %v", tt.name, phase)
				}
			}
		})
	}
}

func TestIsValidPhase(t *testing.T) {
	validPhases := []string{"discovery", "planning", "design", "implementation", "testing", "deployment"}
	invalidPhases := []string{"nonexistent", "", "invalid", "DISCOVERY"}

	for _, phase := range validPhases {
		if !IsValidPhase(phase) {
			t.Errorf("Expected phase '%s' to be valid", phase)
		}
	}

	for _, phase := range invalidPhases {
		if IsValidPhase(phase) {
			t.Errorf("Expected phase '%s' to be invalid", phase)
		}
	}
}

func TestPhaseNames(t *testing.T) {
	names := PhaseNames()
	if len(names) != len(AllPhases) {
		t.Errorf("Expected %d names, got %d", len(AllPhases), len(names))
	}

	// Check that names match phase order
	for i, name := range names {
		if i >= len(AllPhases) || AllPhases[i].Name != name {
			t.Errorf("Expected name %d to be '%s', got '%s'", i, getPhaseNameAtTest(i), name)
		}
	}
}

func TestGetDefaultTool(t *testing.T) {
	tests := []struct {
		phase    string
		expected string
	}{
		{"discovery", "gemini"},
		{"planning", "opencode"},
		{"design", "opencode"},
		{"implementation", "codex"},
		{"testing", "gemini"},
		{"deployment", "fabric"},
		{"nonexistent", "gemini"}, // default fallback
	}

	for _, tt := range tests {
		t.Run(tt.phase, func(t *testing.T) {
			tool := GetDefaultTool(tt.phase)
			if tool != tt.expected {
				t.Errorf("Expected tool '%s' for phase '%s', got '%s'", tt.expected, tt.phase, tool)
			}
		})
	}
}

func TestNextPhase(t *testing.T) {
	tests := []struct {
		current  string
		expected string
	}{
		{"discovery", "planning"},
		{"planning", "design"},
		{"design", "implementation"},
		{"implementation", "testing"},
		{"testing", "deployment"},
		{"deployment", ""},  // last phase
		{"nonexistent", ""}, // invalid phase
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.current, func(t *testing.T) {
			next := NextPhase(tt.current)
			if next != tt.expected {
				t.Errorf("Expected next phase '%s' after '%s', got '%s'", tt.expected, tt.current, next)
			}
		})
	}
}

func TestPreviousPhase(t *testing.T) {
	tests := []struct {
		current  string
		expected string
	}{
		{"discovery", ""}, // first phase
		{"planning", "discovery"},
		{"design", "planning"},
		{"implementation", "design"},
		{"testing", "implementation"},
		{"deployment", "testing"},
		{"nonexistent", ""}, // invalid phase
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.current, func(t *testing.T) {
			prev := PreviousPhase(tt.current)
			if prev != tt.expected {
				t.Errorf("Expected previous phase '%s' before '%s', got '%s'", tt.expected, tt.current, prev)
			}
		})
	}
}

func TestPhaseStructure(t *testing.T) {
	// Test that each phase has required fields
	for _, phase := range AllPhases {
		if phase.Name == "" {
			t.Error("Phase name should not be empty")
		}
		if phase.Description == "" {
			t.Errorf("Phase '%s' should have a description", phase.Name)
		}
		if phase.PrimaryTool == "" {
			t.Errorf("Phase '%s' should have a primary tool", phase.Name)
		}
		if phase.ToolReason == "" {
			t.Errorf("Phase '%s' should have a tool reason", phase.Name)
		}
		if len(phase.Checkpoint.Criteria) == 0 {
			t.Errorf("Phase '%s' should have checkpoint criteria", phase.Name)
		}
		if len(phase.Artifacts) == 0 {
			t.Errorf("Phase '%s' should have artifacts", phase.Name)
		}
	}
}

func TestPhaseSpecifics(t *testing.T) {
	discovery := GetPhase("discovery")
	if discovery == nil {
		t.Fatal("Expected to find discovery phase")
	}
	if discovery.PrimaryTool != "gemini" {
		t.Errorf("Expected discovery primary tool to be 'gemini', got '%s'", discovery.PrimaryTool)
	}
	if len(discovery.Checkpoint.Criteria) != 4 {
		t.Errorf("Expected discovery to have 4 criteria, got %d", len(discovery.Checkpoint.Criteria))
	}

	implementation := GetPhase("implementation")
	if implementation == nil {
		t.Fatal("Expected to find implementation phase")
	}
	if implementation.PrimaryTool != "codex" {
		t.Errorf("Expected implementation primary tool to be 'codex', got '%s'", implementation.PrimaryTool)
	}
}

// Helper functions
func getPhaseAtTest(phases []string, index int) string {
	if index >= 0 && index < len(phases) {
		return phases[index]
	}
	return ""
}

func getPhaseNameAtTest(index int) string {
	if index >= 0 && index < len(AllPhases) {
		return AllPhases[index].Name
	}
	return ""
}

func getPhaseNameAt(index int) string {
	if index >= 0 && index < len(AllPhases) {
		return AllPhases[index].Name
	}
	return ""
}
