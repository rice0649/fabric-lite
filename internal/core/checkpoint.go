package core

import (
	"os"
	"path/filepath"
)

// CheckResult represents the result of a single check
type CheckResult struct {
	Name    string
	Passed  bool
	Message string
}

// ValidationResult represents the overall checkpoint validation result
type ValidationResult struct {
	Phase  string
	Passed bool
	Checks []CheckResult
}

// ValidateCheckpoint runs checkpoint validation for a phase
func ValidateCheckpoint(phaseName string) ValidationResult {
	result := ValidationResult{
		Phase:  phaseName,
		Passed: true,
	}

	phase := GetPhase(phaseName)
	if phase == nil {
		result.Passed = false
		result.Checks = append(result.Checks, CheckResult{
			Name:    "Phase validation",
			Passed:  false,
			Message: "Unknown phase",
		})
		return result
	}

	// Run phase-specific checks
	switch phaseName {
	case "discovery":
		result.Checks = validateDiscovery()
	case "planning":
		result.Checks = validatePlanning()
	case "design":
		result.Checks = validateDesign()
	case "implementation":
		result.Checks = validateImplementation()
	case "testing":
		result.Checks = validateTesting()
	case "deployment":
		result.Checks = validateDeployment()
	default:
		// Generic artifact check for any phase
		result.Checks = validateArtifacts(phaseName, phase.Artifacts)
	}

	// Check if all passed
	for _, check := range result.Checks {
		if !check.Passed {
			result.Passed = false
			break
		}
	}

	return result
}

func validateDiscovery() []CheckResult {
	checks := []CheckResult{}

	// Check for artifacts
	artifactDir := ".forge/artifacts/discovery"

	// Check requirements
	checks = append(checks, checkFileExists(
		"Requirements document",
		filepath.Join(artifactDir, "requirements.md"),
		"docs/requirements.md",
		"REQUIREMENTS.md",
	))

	// Check user stories
	checks = append(checks, checkFileExists(
		"User stories or use cases",
		filepath.Join(artifactDir, "user_stories.md"),
		"docs/user_stories.md",
	))

	// Check research notes
	checks = append(checks, checkFileExists(
		"Research notes",
		filepath.Join(artifactDir, "research_notes.md"),
		filepath.Join(artifactDir, "notes.md"),
	))

	return checks
}

func validatePlanning() []CheckResult {
	checks := []CheckResult{}

	artifactDir := ".forge/artifacts/planning"

	// Check architecture
	checks = append(checks, checkFileExists(
		"Architecture document",
		filepath.Join(artifactDir, "architecture.md"),
		"docs/architecture.md",
		"ARCHITECTURE.md",
	))

	// Check components
	checks = append(checks, checkFileExists(
		"Component breakdown",
		filepath.Join(artifactDir, "components.md"),
		"docs/components.md",
	))

	// Check tech decisions
	checks = append(checks, checkFileExists(
		"Technology decisions",
		filepath.Join(artifactDir, "tech_decisions.md"),
		"docs/tech_stack.md",
		"docs/adr/",
	))

	return checks
}

func validateDesign() []CheckResult {
	checks := []CheckResult{}

	artifactDir := ".forge/artifacts/design"

	// Check API spec
	checks = append(checks, checkFileExists(
		"API specification",
		filepath.Join(artifactDir, "api_spec.md"),
		"docs/api.md",
		"openapi.yaml",
		"swagger.yaml",
	))

	// Check data models
	checks = append(checks, checkFileExists(
		"Data models",
		filepath.Join(artifactDir, "data_models.md"),
		"docs/models.md",
		"docs/schema.md",
	))

	return checks
}

func validateImplementation() []CheckResult {
	checks := []CheckResult{}

	// Check if code exists
	codeExists := false
	for _, dir := range []string{"src", "cmd", "internal", "pkg", "lib", "app"} {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			codeExists = true
			break
		}
	}
	checks = append(checks, CheckResult{
		Name:   "Source code exists",
		Passed: codeExists,
		Message: func() string {
			if codeExists {
				return ""
			}
			return "No source code directory found"
		}(),
	})

	// Check if builds (by looking for common build artifacts or scripts)
	buildable := false
	buildFiles := []string{
		"Makefile", "go.mod", "package.json", "Cargo.toml",
		"setup.py", "pyproject.toml", "build.gradle", "pom.xml",
	}
	for _, f := range buildFiles {
		if _, err := os.Stat(f); err == nil {
			buildable = true
			break
		}
	}
	checks = append(checks, CheckResult{
		Name:   "Build configuration exists",
		Passed: buildable,
		Message: func() string {
			if buildable {
				return ""
			}
			return "No build file found (Makefile, go.mod, package.json, etc.)"
		}(),
	})

	return checks
}

func validateTesting() []CheckResult {
	checks := []CheckResult{}

	// Check if tests exist
	testExists := false
	testDirs := []string{"tests", "test", "__tests__", "spec"}
	for _, dir := range testDirs {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			testExists = true
			break
		}
	}

	// Also check for test files in common patterns
	testPatterns := []string{"*_test.go", "*_test.py", "*.test.js", "*.test.ts", "*.spec.js", "*.spec.ts"}
	for _, pattern := range testPatterns {
		matches, _ := filepath.Glob(pattern)
		if len(matches) > 0 {
			testExists = true
			break
		}
		// Also check in src directory
		matches, _ = filepath.Glob(filepath.Join("src", "**", pattern))
		if len(matches) > 0 {
			testExists = true
			break
		}
	}

	checks = append(checks, CheckResult{
		Name:   "Test files exist",
		Passed: testExists,
		Message: func() string {
			if testExists {
				return ""
			}
			return "No test directory or test files found"
		}(),
	})

	return checks
}

func validateDeployment() []CheckResult {
	checks := []CheckResult{}

	// Check README
	checks = append(checks, checkFileExists(
		"README exists",
		"README.md",
		"README",
		"readme.md",
	))

	// Check changelog
	checks = append(checks, checkFileExists(
		"Changelog exists",
		"CHANGELOG.md",
		"CHANGELOG",
		"HISTORY.md",
	))

	return checks
}

func validateArtifacts(phase string, artifacts []string) []CheckResult {
	checks := []CheckResult{}
	artifactDir := filepath.Join(".forge", "artifacts", phase)

	for _, artifact := range artifacts {
		path := filepath.Join(artifactDir, artifact)
		checks = append(checks, checkFileExists(artifact, path))
	}

	return checks
}

func checkFileExists(name string, paths ...string) CheckResult {
	for _, path := range paths {
		if info, err := os.Stat(path); err == nil {
			if info.IsDir() {
				// Check if directory has any files
				entries, _ := os.ReadDir(path)
				if len(entries) > 0 {
					return CheckResult{Name: name, Passed: true}
				}
			} else {
				return CheckResult{Name: name, Passed: true}
			}
		}
	}

	return CheckResult{
		Name:    name,
		Passed:  false,
		Message: "File not found (tried: " + paths[0] + ")",
	}
}
