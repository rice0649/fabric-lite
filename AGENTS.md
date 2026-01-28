See /home/oak38/projects/AGENTS.md for auto-resume and security protocol.

# Repository Guidelines for Agentic Coding

## Project Structure & Module Organization

- `cmd/fabric-lite/` and `cmd/forge/` hold CLI entry points
- `internal/` contains core logic, CLI handling, and provider integrations:
  - `core/` - Business logic, config management, state, phases, checkpoints
  - `cli/` - Command-line interface commands and initialization
  - `providers/` - AI provider integrations (OpenAI, Anthropic, Ollama, HTTP, Executable)
  - `executor/` - Pattern execution engine
  - `tools/` - Direct tool integrations (codex, gemini, claude, opencode, fabric, ollama)
- `patterns/` stores built-in prompt patterns using `snake_case` naming
- `config/` provides configuration templates; `agents/` contains orchestration assets

## Build, Test, and Development Commands

### Build Commands
- `make build` - Build fabric-lite CLI to `bin/fabric-lite`
- `make build-forge` - Build forge CLI to `bin/forge`
- `make build-all` - Build both binaries
- `make install` - Install fabric-lite to `~/go/bin`
- `make install-forge` - Install forge to `~/go/bin`
- `make install-all` - Full install (binaries + patterns)

### Test Commands
- `make test` - Run all tests with `go test -v ./...`
- `make test-coverage` - Generate coverage report in `coverage.html`
- **Single test**: `go test -v ./internal/core -run TestNewProjectConfig`
- **Verbose single test**: `go test -v ./internal/core -run TestNewProjectConfig -v`

### Development Commands
- `make run ARGS='--help'` - Run CLI directly via `go run`
- `make lint` - Run `golangci-lint` (install if missing)
- `make fmt` - Format Go code with `go fmt ./...`
- `make clean` - Remove build artifacts

## Code Style Guidelines

### Formatting & Imports
- Use **tabs for indentation** in Go files (per `.editorconfig`)
- Run `make fmt` before committing - uses standard `go fmt`
- Import groups: standard library, third-party, internal (separated by blank lines)
- Sort imports alphabetically within groups

### Naming Conventions
- **Exported identifiers**: `CamelCase` (e.g., `ProjectConfig`, `NewProvider`)
- **Unexported identifiers**: `camelCase` (e.g., `configManager`, `executeRequest`)
- **Constants**: `UPPER_SNAKE_CASE` for exported, `camelCase` for unexported
- **Interfaces**: End with `-er` suffix (e.g., `Provider`, `Executor`)
- **Files**: `snake_case.go` for packages, `camelCase_test.go` for tests

### Types & Structures
- Use struct tags for YAML/JSON: `yaml:"name,omitempty"`
- Group related fields in structs; use embedding carefully
- Prefer composition over inheritance
- Define error types as var errors or custom types with `Error() string`

### Error Handling
- Always handle errors explicitly; never ignore with `_`
- Use wrapped errors: `fmt.Errorf("operation failed: %w", err)`
- Return errors as last return value
- Use descriptive error messages that explain context

### Testing Patterns
- Test files: `filename_test.go` beside source code
- Test functions: `TestFunctionName_Condition`
- Subtests: `t.Run("subcase name", func(t *testing.T) { ... })`
- Table-driven tests for multiple scenarios
- Mock external dependencies using interfaces

### Configuration & Constants
- Configuration structs in `internal/core/config.go`
- Use `yaml:"field,omitempty"` for optional fields
- Default values in constructors like `NewProjectConfig()`
- Environment variable hints in config comments

## Architecture Guidelines

### Provider Pattern
- All AI providers implement `Provider` interface in `internal/providers/provider.go`
- Use provider manager for abstraction and fallback logic
- HTTP and Executable providers for generic integrations

### CLI Design
- Use Cobra for CLI structure (standard for Go CLI tools)
- Commands in `internal/cli/` with separate files per command
- Persistent flags for global options, local flags for command-specific

### State Management
- Forge projects use `.forge/` directory with `config.yaml` and `state.yaml`
- Checkpoints for resumable workflows in `.forge/checkpoints/`
- User context in `SHARED_CONTEXT.md` for personalized workflows

## Security & Configuration

- API keys from environment variables, not hardcoded
- Use example config templates in `config/`
- Never commit personal configurations or secrets
- Local LLM support via Ollama for privacy-conscious users

## Testing Requirements

Before committing, ensure:
1. `make fmt` passes (code is formatted)
2. `make lint` passes (if golangci-lint is installed)
3. `make test` passes with all tests
4. New features include appropriate test coverage
5. Documentation is updated for user-facing changes

---

## Auto-Resume Protocol
If the user says “resume”, “good morning”, “where were we”, or similar:
1. Read `/home/oak38/projects/SHARED_CONTEXT.md`
2. Read the latest session log in this project (if present)
3. Summarize last completed work, current blockers, and the next 1–3 actions

## Repository Security Protocol
- Never commit secrets (API keys, tokens, passwords, private keys, `.env`, SSNs, PII/CUI)
- Scan before pushing (use `rg` and review hits)
- Use `.env.example` for variable names only
- Ensure `.gitignore` includes `.env`, `*.key`, `*.pem`, `*.p12`, `*.crt`, `*.sqlite`, `*.db`, `secrets.*`, `.claude/`
- Default GitHub repos to private unless explicitly approved for public
