package core

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// ProjectState tracks the current state of a forge project
type ProjectState struct {
	CurrentPhase   string            `yaml:"current_phase"`
	PhaseStartedAt time.Time         `yaml:"phase_started_at,omitempty"`
	PhaseStatuses  map[string]string `yaml:"phase_statuses"` // phase -> status (pending, in_progress, completed)
	Activities     []Activity        `yaml:"activities"`
	CreatedAt      time.Time         `yaml:"created_at"`
	UpdatedAt      time.Time         `yaml:"updated_at"`
}

// Activity represents a logged activity in the project
type Activity struct {
	Timestamp time.Time `yaml:"timestamp"`
	Message   string    `yaml:"message"`
	Phase     string    `yaml:"phase,omitempty"`
}

// PhaseHistory stores completion data for a phase
type PhaseHistory struct {
	Phase       string        `yaml:"phase"`
	StartedAt   time.Time     `yaml:"started_at"`
	CompletedAt time.Time     `yaml:"completed_at"`
	Duration    time.Duration `yaml:"duration"`
	Notes       string        `yaml:"notes,omitempty"`
}

// NewProjectState creates a new project state
func NewProjectState() *ProjectState {
	now := time.Now()
	return &ProjectState{
		PhaseStatuses: map[string]string{
			"discovery":      "pending",
			"planning":       "pending",
			"design":         "pending",
			"implementation": "pending",
			"testing":        "pending",
			"deployment":     "pending",
		},
		Activities: []Activity{
			{
				Timestamp: now,
				Message:   "Project initialized",
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// GetPhaseStatus returns the status of a phase
func (s *ProjectState) GetPhaseStatus(phase string) string {
	if status, ok := s.PhaseStatuses[phase]; ok {
		return status
	}
	return "pending"
}

// SetPhaseStatus updates the status of a phase
func (s *ProjectState) SetPhaseStatus(phase, status string) {
	if s.PhaseStatuses == nil {
		s.PhaseStatuses = make(map[string]string)
	}
	s.PhaseStatuses[phase] = status
	s.UpdatedAt = time.Now()
}

// AddActivity logs a new activity
func (s *ProjectState) AddActivity(message string) {
	s.Activities = append(s.Activities, Activity{
		Timestamp: time.Now(),
		Message:   message,
		Phase:     s.CurrentPhase,
	})
	s.UpdatedAt = time.Now()
}

// Save writes the state to a YAML file
func (s *ProjectState) Save(path string) error {
	s.UpdatedAt = time.Now()
	data, err := yaml.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadProjectState loads state from a YAML file
func LoadProjectState(path string) (*ProjectState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var state ProjectState
	if err := yaml.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	if state.PhaseStatuses == nil {
		state.PhaseStatuses = make(map[string]string)
	}

	return &state, nil
}

// Save writes the phase history to a YAML file
func (h *PhaseHistory) Save(path string) error {
	data, err := yaml.Marshal(h)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
