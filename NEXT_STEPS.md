# Fabric-Lite: Next Steps Implementation Plan

**Generated**: 2026-01-20
**Analysis Source**: Codex code review + manual verification
**Status**: Ready for implementation

---

## Executive Summary

The fabric-lite project has a solid foundation with working providers and pattern execution. However, several key gaps prevent full functionality:

| Category | Status | Priority |
|----------|--------|----------|
| Provider System | Complete | - |
| Pattern Executor | 90% (missing streaming) | HIGH |
| CLI Commands | 40% wired | CRITICAL |
| Tool Implementations | 5/6 working | HIGH |
| Test Coverage | 0% | MEDIUM |

---

## Critical Gaps Identified

### 1. CLI Commands Not Wired (CRITICAL)

**Location**: `internal/cli/root.go:39-43`

Currently registered:
- `run`, `list`, `config`, `version`

**Missing (exist but not registered)**:
| Command | File | Purpose |
|---------|------|---------|
| `auto` | `internal/cli/auto.go` | Automated phase execution |
| `phase` | `internal/cli/phase.go` | Phase management |
| `status` | `internal/cli/status.go` | Project status dashboard |
| `session` | `internal/cli/session.go` | Session state management |
| `init` | `internal/cli/init.go` | Project initialization |

**Fix**: Add to `internal/cli/root.go`:
```go
rootCmd.AddCommand(newAutoCmd())
rootCmd.AddCommand(newPhaseCmd())
rootCmd.AddCommand(newStatusCmd())
rootCmd.AddCommand(newSessionCmd())
rootCmd.AddCommand(newInitCmd())
```

---

### 2. Claude Tool is Stub (HIGH)

**Location**: `internal/tools/claude.go:20`

**Current state**: Returns "Claude tool is not yet implemented"

**Fix**: Implement Execute() method similar to GeminiTool:
- Wrap the `claude` CLI command
- Support interactive and non-interactive modes
- Add ExecuteNonInteractive() for automation

---

### 3. Streaming Not Implemented (HIGH)

**Location**: `internal/executor/pattern.go:183`

**Current state**: `Stream: false` hardcoded with TODO comment

**Fix**:
- Implement streaming path in pattern executor
- Use provider's streaming capabilities
- Add SSE output for real-time responses

---

### 4. Codex Config Loading (MEDIUM)

**Location**: `internal/tools/codex.go:36`

**Current state**: Hardcoded provider/model, ignores user config

**Fix**:
- Load config from core.LoadConfig()
- Use configured provider and model
- Respect user preferences

---

### 5. Test Failures (MEDIUM)

**Current state**:
- `TestRunListFlag` fails (unknown flag: --list)
- 0% coverage across all packages

**Fix**:
- Update test to match current CLI structure
- Add unit tests for providers, executor, tools
- Target: 60%+ coverage

---

### 6. Naming Inconsistency (LOW)

**Issue**: Some CLI files reference "forge" instead of "fabric-lite"

**Affected files**:
- `internal/cli/auto.go`
- `internal/cli/phase.go`
- `internal/cli/session.go`
- `internal/cli/init.go`

**Fix**: Update command names and help text to use "fabric-lite"

---

## Implementation Priority Order

### Sprint 1: CLI Completeness (Est. effort: Small)
1. Wire missing CLI commands in root.go
2. Fix TestRunListFlag test
3. Update naming from "forge" to "fabric-lite"

### Sprint 2: Tool Completion (Est. effort: Medium)
1. Implement ClaudeTool.Execute() properly
2. Fix Codex config loading
3. Add tests for all 6 tools

### Sprint 3: Streaming Support (Est. effort: Medium)
1. Implement streaming in pattern executor
2. Add streaming support to providers that support it
3. Add --stream flag handling in CLI

### Sprint 4: Test Coverage (Est. effort: Large)
1. Add unit tests for providers
2. Add unit tests for executor
3. Add integration tests for CLI commands
4. Target 60%+ coverage

---

## Quick Wins (Can Do Now)

1. **Wire CLI commands** - 5 lines of code in root.go
2. **Fix failing test** - Update test expectations
3. **Update help text** - Search/replace "forge" → "fabric-lite"

---

## Files to Modify

| File | Changes |
|------|---------|
| `internal/cli/root.go` | Add 5 missing command registrations |
| `internal/tools/claude.go` | Implement Execute() method |
| `internal/tools/codex.go` | Load config properly |
| `internal/executor/pattern.go` | Add streaming support |
| `cmd/fabric-lite/main_test.go` | Fix failing test |

---

## Success Criteria

- [ ] All CLI commands accessible via `fabric-lite --help`
- [ ] All 6 tools have working Execute() methods
- [ ] Streaming works with `--stream` flag
- [ ] Tests pass with 60%+ coverage
- [ ] No "forge" references in user-facing text
