# Multi-Agent Orchestration System

## Overview

This system runs multiple specialized agents in separate terminal sessions, communicating via shared files.

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

| Agent | Input | Output | Terminal |
|-------|-------|--------|----------|
| Pattern Discovery | Original Fabric patterns | patterns_discovered.md, queue | 1 |
| Planning | Fabric codebase | 01_planning.md | 2 |
| Review | Planning output | 02_review.md | 2 |
| Code Analysis | Fabric source files | 03_code_analysis.md | 2 |
| Final Review | All outputs | 04_final_spec.md, FINAL_SUMMARY.md | 2 |
| Implementation | Final spec + patterns | Code changes | 3 |
