See /home/oak38/projects/AGENTS.md for auto-resume and security protocol.

# Quick Start: Multi-Agent Workflow

## Overview

Run specialized agents in separate terminals to keep each session focused and lightweight.

## Terminal Setup

You need **3-4 terminal windows**:

```
┌─────────────────┬─────────────────┬─────────────────┐
│   TERMINAL 1    │   TERMINAL 2    │   TERMINAL 3    │
│                 │                 │                 │
│   Pattern       │   Planning      │   Monitor/      │
│   Discovery     │   Pipeline      │   Main Session  │
│                 │                 │                 │
└─────────────────┴─────────────────┴─────────────────┘
                          │
                          ▼
                  ┌─────────────────┐
                  │   TERMINAL 4    │
                  │                 │
                  │ Implementation  │
                  │ (after planning)│
                  └─────────────────┘
```

## Step-by-Step

### Step 1: Open Terminal 3 (Monitor)
```bash
cd /home/oak38/projects/fabric-lite
./agents/runners/monitor.sh
```
This shows real-time status of all agent outputs.

### Step 2: Open Terminal 1 (Pattern Discovery)
```bash
cd /home/oak38/projects/fabric-lite
./agents/runners/run_pattern_discovery.sh
```
When Gemini starts, say:
> "Read agents/pattern_discovery_agent.md and begin scanning patterns"

### Step 3: Open Terminal 2 (Planning Pipeline)
```bash
cd /home/oak38/projects/fabric-lite
./agents/runners/run_planning.sh
```
When Gemini starts, say:
> "Read agents/planning_agent.md and begin the planning pipeline"

### Step 4: Wait for Completion

Watch Terminal 3 (monitor) for:
- `05_patterns_discovered.md: complete`
- `04_final_spec.md: complete`
- `FINAL_SUMMARY.md` appears

### Step 5: Open Terminal 4 (Implementation)
```bash
cd /home/oak38/projects/fabric-lite
./agents/runners/run_implementation.sh
```
When Gemini starts, say:
> "Read agents/implementation_agent.md and begin implementation"

## Output Files

All outputs go to `agents/outputs/`:

| File | Agent | Contains |
|------|-------|----------|
| `01_planning.md` | Planning | Architecture plan |
| `02_review.md` | Review | Feasibility assessment |
| `03_code_analysis.md` | Code Analysis | Source code analysis |
| `04_final_spec.md` | Final Review | Implementation spec |
| `05_patterns_discovered.md` | Pattern Discovery | Top patterns list |
| `FINAL_SUMMARY.md` | Final Review | Executive summary for main session |
| `implementation_progress.md` | Implementation | Progress tracking |

## Tips

1. **Don't rush**: Let each agent complete before moving to the next pipeline stage
2. **Check monitor**: Terminal 3 shows real-time status
3. **Read summaries**: `FINAL_SUMMARY.md` has everything you need to review
4. **Implementation is last**: Only start after planning + patterns are done

## Simplified Flow (If Short on Time)

If you want a faster flow, run just:

**Terminal 1**: Pattern Discovery (can run in background)
**Terminal 2**: Planning → Implementation (sequential in same terminal)

```bash
# Terminal 1
cd /home/oak38/projects/fabric-lite
gemini "Read agents/pattern_discovery_agent.md and scan patterns"

# Terminal 2
cd /home/oak38/projects/fabric-lite
gemini "Run the full pipeline: planning_agent.md → review_agent.md → final_review_agent.md → implementation_agent.md"
```

## Troubleshooting

**Agent not writing output?**
- Make sure it knows the output path: `agents/outputs/`
- Tell it: "Write your output to agents/outputs/[filename].md"

**Implementation agent can't find spec?**
- Run planning pipeline first
- Check that `04_final_spec.md` exists

**Monitor not updating?**
- Outputs are only created when agents finish writing
- Check that agents are actually running
