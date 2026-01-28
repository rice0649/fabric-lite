See /home/oak38/projects/AGENTS.md for auto-resume and security protocol.

# Development Phases Guide

AI Project Forge uses a structured phase-based approach to software development. Each phase has a specific purpose, recommended AI tool, and checkpoint criteria.

## Phase Overview

```
Discovery → Planning → Design → Implementation → Testing → Deployment
```

## 1. Discovery Phase

**Purpose**: Research and requirements gathering

**Primary Tool**: Gemini CLI
- Free tier available
- 1M token context window
- Google Search integration for research

**Checkpoint Criteria**:
- Requirements document exists
- User stories or use cases defined
- Technical constraints identified
- Research notes compiled

**Artifacts**:
- `requirements.md`
- `user_stories.md`
- `research_notes.md`

**Commands**:
```bash
forge phase start discovery
forge run                    # Uses Gemini CLI
forge run --pattern research_topic  # Use fabric-lite pattern
forge phase complete
```

## 2. Planning Phase

**Purpose**: Architecture and component design

**Primary Tool**: OpenCode
- Provider-agnostic (Anthropic, OpenAI, etc.)
- Plan mode for read-only exploration
- Great for architecture discussions

**Checkpoint Criteria**:
- Architecture document exists
- Component breakdown defined
- Technology decisions documented
- Dependencies identified

**Artifacts**:
- `architecture.md`
- `components.md`
- `tech_decisions.md`

**Commands**:
```bash
forge phase start planning
forge run                    # Uses OpenCode
forge run --pattern create_architecture
forge phase complete
```

## 3. Design Phase

**Purpose**: API and data model definition

**Primary Tool**: OpenCode
- Continues context from planning
- Ideal for API design discussions
- Schema and interface modeling

**Checkpoint Criteria**:
- API specification defined
- Data models documented
- Interface contracts specified
- Error handling strategy defined

**Artifacts**:
- `api_spec.md`
- `data_models.md`
- `interfaces.md`

**Commands**:
```bash
forge phase start design
forge run                    # Uses OpenCode
forge run --pattern create_api_spec
forge phase complete
```

## 4. Implementation Phase

**Purpose**: Code development and feature building

**Primary Tool**: Codex CLI
- Advanced reasoning (o3-mini)
- Code generation capabilities
- Multimodal support for diagrams

**Checkpoint Criteria**:
- Code builds successfully
- Core features implemented
- Basic tests exist
- No critical linting errors

**Artifacts**:
- `implementation_notes.md`
- `code_review.md`

**Commands**:
```bash
forge phase start implementation
forge run                    # Uses Codex CLI
forge run --prompt "Implement the user authentication feature"
forge phase complete
```

## 5. Testing Phase

**Purpose**: Test creation and quality assurance

**Primary Tool**: Gemini CLI
- Large context for analyzing codebase
- Good at generating test cases
- Coverage analysis

**Checkpoint Criteria**:
- Tests pass
- Coverage threshold met
- Edge cases covered
- Integration tests exist

**Artifacts**:
- `test_plan.md`
- `coverage_report.md`

**Commands**:
```bash
forge phase start testing
forge run                    # Uses Gemini CLI
forge run --pattern create_test_plan
forge phase complete
```

## 6. Deployment Phase

**Purpose**: Documentation and release preparation

**Primary Tool**: fabric-lite
- Pattern-based document generation
- Consistent output format
- Release notes and changelog generation

**Checkpoint Criteria**:
- README is complete
- Changelog updated
- Deployment docs exist
- Release notes prepared

**Artifacts**:
- `changelog.md`
- `release_notes.md`
- `deployment_guide.md`

**Commands**:
```bash
forge phase start deployment
forge run --pattern create_release_notes
forge run --pattern create_changelog
forge phase complete
```

## Skipping Phases

You can skip phases or go back to previous phases:

```bash
# Skip to implementation (not recommended)
forge phase start implementation --force

# Go back to planning
forge phase start planning --force
```

## Customizing Tool Selection

Override the default tool for a phase:

```bash
# Use Codex for discovery instead of Gemini
forge run --tool codex

# Use OpenCode for implementation
forge run --tool opencode
```

Or configure in `.forge/config.yaml`:

```yaml
phases:
  discovery: codex
  planning: gemini
```

## Checkpoint Validation

When completing a phase, forge validates checkpoint criteria:

```bash
$ forge phase complete
Running checkpoint validation for phase: discovery

  ✓ Requirements document
  ✓ User stories
  ✗ Research notes
      File not found (tried: .forge/artifacts/discovery/research_notes.md)

Checkpoint validation failed. Address the issues above or use --skip-check.
```

Skip validation (not recommended):

```bash
forge phase complete --skip-check
```

## Best Practices

1. **Don't skip phases**: Each phase builds on the previous
2. **Save artifacts**: Store AI outputs in the artifacts directory
3. **Review checkpoints**: They ensure quality before moving on
4. **Iterate within phases**: Run the tool multiple times as needed
5. **Use patterns**: fabric-lite patterns provide consistent outputs
