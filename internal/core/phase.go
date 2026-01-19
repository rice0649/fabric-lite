package core

// Phase represents a development phase in the forge workflow
type Phase struct {
	Name        string
	Description string
	PrimaryTool string
	ToolReason  string
	Checkpoint  Checkpoint
	Artifacts   []string
}

// Checkpoint defines validation criteria for completing a phase
type Checkpoint struct {
	Criteria []string
}

// AllPhases defines the ordered list of development phases
var AllPhases = []Phase{
	{
		Name:        "discovery",
		Description: "Research and requirements gathering",
		PrimaryTool: "gemini",
		ToolReason:  "Free tier, 1M context window, Google Search integration for research",
		Checkpoint: Checkpoint{
			Criteria: []string{
				"Requirements document exists",
				"User stories or use cases defined",
				"Technical constraints identified",
				"Research notes compiled",
			},
		},
		Artifacts: []string{
			"requirements.md",
			"user_stories.md",
			"research_notes.md",
		},
	},
	{
		Name:        "planning",
		Description: "Architecture and component design",
		PrimaryTool: "opencode",
		ToolReason:  "Read-only exploration mode, provider-agnostic planning",
		Checkpoint: Checkpoint{
			Criteria: []string{
				"Architecture document exists",
				"Component breakdown defined",
				"Technology decisions documented",
				"Dependencies identified",
			},
		},
		Artifacts: []string{
			"architecture.md",
			"components.md",
			"tech_decisions.md",
		},
	},
	{
		Name:        "design",
		Description: "API and data model definition",
		PrimaryTool: "opencode",
		ToolReason:  "Continue planning context, ideal for API/schema design",
		Checkpoint: Checkpoint{
			Criteria: []string{
				"API specification defined",
				"Data models documented",
				"Interface contracts specified",
				"Error handling strategy defined",
			},
		},
		Artifacts: []string{
			"api_spec.md",
			"data_models.md",
			"interfaces.md",
		},
	},
	{
		Name:        "implementation",
		Description: "Code development and feature building",
		PrimaryTool: "codex",
		ToolReason:  "Advanced reasoning, code review capabilities, multimodal support",
		Checkpoint: Checkpoint{
			Criteria: []string{
				"Code builds successfully",
				"Core features implemented",
				"Basic tests exist",
				"No critical errors in linting",
			},
		},
		Artifacts: []string{
			"implementation_notes.md",
			"code_review.md",
		},
	},
	{
		Name:        "testing",
		Description: "Test creation and quality assurance",
		PrimaryTool: "gemini",
		ToolReason:  "Large context for analyzing coverage, good at test generation",
		Checkpoint: Checkpoint{
			Criteria: []string{
				"Tests pass",
				"Coverage threshold met (if configured)",
				"Edge cases covered",
				"Integration tests exist",
			},
		},
		Artifacts: []string{
			"test_plan.md",
			"coverage_report.md",
		},
	},
	{
		Name:        "deployment",
		Description: "Documentation and release preparation",
		PrimaryTool: "fabric",
		ToolReason:  "Pattern-based generation for docs, changelogs, release notes",
		Checkpoint: Checkpoint{
			Criteria: []string{
				"README is complete",
				"Changelog updated",
				"Deployment docs exist",
				"Release notes prepared",
			},
		},
		Artifacts: []string{
			"changelog.md",
			"release_notes.md",
			"deployment_guide.md",
		},
	},
}

// phaseMap provides O(1) lookup by name
var phaseMap = func() map[string]*Phase {
	m := make(map[string]*Phase)
	for i := range AllPhases {
		m[AllPhases[i].Name] = &AllPhases[i]
	}
	return m
}()

// GetPhase returns a phase by name
func GetPhase(name string) *Phase {
	if p, ok := phaseMap[name]; ok {
		return p
	}
	return nil
}

// IsValidPhase checks if a phase name is valid
func IsValidPhase(name string) bool {
	_, ok := phaseMap[name]
	return ok
}

// PhaseNames returns all phase names in order
func PhaseNames() []string {
	names := make([]string, len(AllPhases))
	for i, p := range AllPhases {
		names[i] = p.Name
	}
	return names
}

// GetDefaultTool returns the primary tool for a phase
func GetDefaultTool(phase string) string {
	if p := GetPhase(phase); p != nil {
		return p.PrimaryTool
	}
	return "gemini" // default fallback
}

// NextPhase returns the next phase after the given one
func NextPhase(current string) string {
	for i, p := range AllPhases {
		if p.Name == current && i < len(AllPhases)-1 {
			return AllPhases[i+1].Name
		}
	}
	return ""
}

// PreviousPhase returns the phase before the given one
func PreviousPhase(current string) string {
	for i, p := range AllPhases {
		if p.Name == current && i > 0 {
			return AllPhases[i-1].Name
		}
	}
	return ""
}
