# Multi-Agent Orchestration System - Your Augmented "Sleeve"

## Overview

This system operates as your augmented "sleeve," a high-performance shell designed to execute your intent. It functions by first "sleeving" your core "consciousness" (your mission, goals, and workflow preferences - your "stack") from the `SHARED_CONTEXT.md` file. It then runs multiple specialized AI agents ("constructs") in separate terminal sessions, communicating via shared files to augment your capabilities.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        AGENT ORCHESTRATION FLOW                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  TERMINAL 1                TERMINAL 2              TERMINAL 3                │
│  ┌──────────────┐         ┌──────────────┐        ┌──────────────┐          │
│  │ Pattern      │         │ Planning     │        │ Main Session │          │
│  │ Discovery    │────────▶│ Pipeline     │───────▶│ (Receives    │          │
│  │ Agent        │         │ (Review)     │        │  Final Only) │          │
│  └──────────────┘         └──────────────┘        └──────────────┘          │
│         │                        │                       │                   │
│         │                        ▼                       │                   │
│         │                 ┌──────────────┐               │                   │
│         │                 │ Final Review │               │                   │
│         │                 │ Agent        │───────────────┤                   │
│         │                 └──────────────┘               │                   │
│         │                        │                       │                   │
│         │                        ▼                       ▼                   │
│         │                 ┌──────────────┐        ┌──────────────┐          │
│         └────────────────▶│ Implementation│◀──────│ User Reviews │          │
│                           │ Agent        │        │ & Approves   │          │
│                           └──────────────┘        └──────────────┘          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Directory Structure

```
agents/
├── ORCHESTRATOR.md              # This file
├── outputs/                     # Agent outputs (shared)
│   ├── 01_planning.md
│   ├── 02_review.md
│   ├── 03_code_analysis.md
│   ├── 04_final_spec.md
│   ├── 05_patterns_discovered.md
│   └── FINAL_SUMMARY.md         # Main session reads this
├── queue/                       # Task queue
│   ├── patterns_to_review.md
│   └── patterns_to_implement.md
├── planning_agent.md
├── review_agent.md
├── code_analysis_agent.md
├── final_review_agent.md
├── pattern_discovery_agent.md
├── implementation_agent.md
└── runners/
    ├── run_planning.sh
    ├── run_pattern_discovery.sh
    └── run_implementation.sh
```

## Communication Protocol

### File-Based Message Passing

Agents communicate via markdown files in `outputs/` and `queue/`:

1. **Output Files**: Each agent writes its output to a numbered file
2. **Queue Files**: Discovery agent queues patterns for review/implementation
3. **Final Summary**: Last agent writes `FINAL_SUMMARY.md` for main session

### Status Markers

Each output file includes a status header:

```markdown
---
status: complete | in_progress | error
agent: agent_name
timestamp: ISO-8601
next: next_agent_name (or "none")
---
```

## Tool Specializations

Within this multi-agent system, specific AI tools ("constructs") are delegated tasks based on their unique strengths, ensuring efficient and high-quality outcomes.

*   **`OpenCode` (Master Planner):** Focuses on breaking down high-level objectives into clear, actionable, multi-phase plans. It orchestrates the overall development strategy.
*   **`Codex` (Code Generation Specialist):** Dedicated to translating detailed plans into functional code. It handles code writing, refactoring, and in-depth code analysis.
*   **`Claude` (Systems Architect):** Excels at complex, large-scale code analysis, architectural design, and ensuring structural integrity across large codebases.
*   **`Gemini` (Researcher & Reviewer):** Specializes in discovery, external research (grounded by Google Search), and critical review of proposals, plans, and implementations against best practices.
*   **`Ollama` (Quick Task Automator):** Designed for rapid execution of simple, well-defined tasks such as boilerplate generation, formatting, or quick checks, offloading simpler work from more powerful constructs.


## Running the System

### Terminal 1: Pattern Discovery
```bash
cd /home/oak38/projects/fabric-lite
./agents/runners/run_pattern_discovery.sh
```

### Terminal 2: Planning Pipeline
```bash
cd /home/oak38/projects/fabric-lite
./agents/runners/run_planning.sh
```

### Terminal 3: Implementation (after planning completes)
```bash
cd /home/oak38/projects/fabric-lite
./agents/runners/run_implementation.sh
```

### Main Session: Monitor & Review
```bash
# Watch for final output
tail -f agents/outputs/FINAL_SUMMARY.md

# Or check status
cat agents/outputs/*.md | grep "^status:"
```

## Agent Responsibilities

| Agent | Input | Output | Primary Tools Utilized | Terminal |
|-------|-------|--------|------------------------|----------|
| Pattern Discovery | Original Fabric patterns | patterns_discovered.md, queue | Gemini (for research and categorization) | 1 |
| Planning | User goals, SHARED_CONTEXT.md | 01_planning.md | OpenCode (master planning), Gemini (initial research) | 2 |
| Review | Planning output, Code Analysis | 02_review.md | Gemini (validation, best practices), Claude (architectural consistency) | 2 |
| Code Analysis | Fabric source files | 03_code_analysis.md | Codex (deep code review), Claude (large-scale analysis) | 2 |
| Final Review | All outputs | 04_final_spec.md, FINAL_SUMMARY.md | Gemini (final synthesis), Claude (comprehensive check) | 2 |
| Implementation | Final spec + patterns | Code changes | Codex (code generation), Claude (complex refactoring), Ollama (boilerplate) | 3 |
