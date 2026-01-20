# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-01-19

### Added

#### AI Project Forge CLI (`forge`)
- **Multi-tool orchestration**: Coordinate AI coding assistants through structured development phases
- **Supported tools**:
  - Gemini CLI - Research and discovery (free tier, 1M context window)
  - OpenCode - Planning and design (read-only exploration)
  - Codex CLI - Implementation (advanced reasoning)
  - fabric-lite - Pattern-based document generation

#### Development Phase System
- Six structured phases: discovery, planning, design, implementation, testing, deployment
- Phase-specific prompts and tool assignments
- Checkpoint system to validate phase completion
- Automatic progress tracking

#### CLI Commands
- `forge init` - Initialize new forge project with configuration
- `forge run` - Execute AI tool for current phase
- `forge phase` - Manage development phases (info, complete, skip)
- `forge status` - Dashboard showing project progress
- `forge history` - View activity history
- `forge session` - Session management (save, show, resume)

#### Project Configuration
- YAML-based configuration (`.forge/config.yaml`)
- Persistent state tracking (`.forge/state.yaml`)
- Artifact storage (`.forge/artifacts/`)
- Activity history (`.forge/history/`)

#### Session Management
- Save session state for project continuity
- Resume prompts for AI assistants
- Full context export for handoffs

### fabric-lite Core
- Pattern-based AI prompt system
- Multiple provider support (OpenAI, Ollama, Anthropic)
- Customizable pattern templates
- Configuration via `~/.config/fabric-lite/`

## [Unreleased]

### Planned
- Testing phase automation
- Additional AI tool integrations
- Plugin system for custom tools
- Web dashboard for project monitoring
