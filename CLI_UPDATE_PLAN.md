See /home/oak38/projects/AGENTS.md for auto-resume and security protocol.

# CLI Documentation Update Plan - fabric-lite & Forge

## 1. CURRENT STATE ASSESSMENT

### fabric-lite Capabilities
- ✅ **Pattern-based AI execution** - Core pattern framework working
- ✅ **Multi-provider support** - OpenAI, Anthropic, Ollama integration
- ✅ **CLI-tools bridge** - Direct tool invocation via `fabric-lite run <tool>`
- ✅ **Unified run command** - Single interface for patterns AND tools
- ✅ **6 AI tools available** - codex, gemini, claude, opencode, fabric, ollama

### Forge Orchestration Status
- ✅ **Multi-agent orchestration** - Agent system with phase-based workflow
- ✅ **Tool coordination** - Can use any of the 6 available tools
- ✅ **Phase management** - Discovery → Planning → Design → Implementation → Testing → Deployment
- ✅ **State persistence** - Checkpoint system with resumable workflows

## 2. TOOL CAPABILITY MAPPING

| Tool | Strength | Best Use Cases | Integration Point |
|------|----------|----------------|-------------------|
| **codex** | Meta-coding tool | Delegates coding tasks, code generation | Implementation phase, complex coding tasks |
| **gemini** | Research & analysis | Discovery, testing, documentation | Discovery phase, test planning |
| **claude** | Advanced reasoning | Architecture design, complex analysis | Planning & Design phases |
| **opencode** | Interactive coding | Real-time code assistance | Implementation, debugging |
| **ollama** | Local processing | Fast, private tasks | Quick local analysis, prototyping |
| **fabric** | Pattern execution | Documentation, structured outputs | Deployment, documentation tasks |

## 3. FILES REQUIRING UPDATES

### Primary Documentation Files
1. **README.md** - Missing tool integration capabilities
2. **docs/getting-started.md** - Needs tool execution section
3. **docs/forge-getting-started.md** - Outdated tool references
4. **AGENTS.md** - May need tool integration updates

### Secondary Files (Review Needed)
5. **agents/QUICKSTART.md** - Tool coordination workflows
6. **docs/phases.md** - Phase-to-tool mapping updates
7. **agents/ORCHESTRATOR.md** - Tool orchestration details

## 4. EXECUTION STRATEGY

### Phase 1: Core Documentation Updates
1. Update README.md with new CLI-tools bridge section
2. Enhance getting-started.md with tool execution examples
3. Update forge-getting-started.md with current tool ecosystem

### Phase 2: Advanced Integration Documentation
4. Update AGENTS.md with tool coordination capabilities
5. Enhance phase documentation with tool mappings
6. Create tool integration examples

### Phase 3: Specialized Content
7. Update quickstart guides with multi-tool workflows
8. Create advanced orchestration examples

## 5. HEADLESS EXECUTION PROTOCOL

### Tool Execution Commands
```bash
# Direct tool execution
fabric-lite run codex -p "implement REST API"
fabric-lite run gemini -p "research microservices patterns"
fabric-lite run claude -p "design system architecture"

# Pattern execution with tool override
fabric-lite run --pattern implement_feature --provider claude
```

### Background Processing Strategy
1. **Batch Processing** - Use tool combinations for complex tasks
2. **Progress Tracking** - Leverage forge checkpoint system
3. **State Persistence** - Save intermediate results for resumption
4. **Error Recovery** - Tool-specific fallback mechanisms

## 6. PROGRESS TRACKING SYSTEM

### Checkpoint Structure
```yaml
state:
  phase: documentation_updates
  current_file: README.md
  completed_sections:
    - installation
    - basic_usage
  pending_sections:
    - tool_integration
    - examples
  tools_used:
    - claude: architecture_analysis
    - gemini: research_completion
```

### Resumable Workflows
- **File-level checkpoints** - Track completion per .md file
- **Section-level progress** - Track completion within files
- **Tool usage logs** - Document which tools contributed to each section
- **Validation checkpoints** - Verify accuracy before proceeding

## 7. INTEGRATION WORKFLOWS

### Multi-Tool Coordination Examples

#### Workflow 1: Documentation Generation
1. **gemini** - Research current capabilities and features
2. **claude** - Structure documentation outline
3. **codex** - Generate code examples
4. **fabric** - Apply formatting patterns
5. **opencode** - Review and refine content

#### Workflow 2: Feature Analysis
1. **ollama** - Quick local code analysis
2. **claude** - Deep architectural review
3. **gemini** - Research best practices
4. **codex** - Generate implementation examples

## 8. EXECUTION TIMELINE

### Session 1: Core Updates (Priority 1)
- Update README.md with tool integration section
- Enhance getting-started.md with tool execution examples
- Update forge-getting-started.md tool references

### Session 2: Integration Documentation (Priority 2)
- Update AGENTS.md with tool coordination
- Enhance phase documentation
- Create tool mapping tables

### Session 3: Advanced Content (Priority 3)
- Update quickstart guides
- Create orchestration examples
- Add troubleshooting sections

## 9. SUCCESS METRICS

### Completion Criteria
- [ ] All 7 target .md files updated with current capabilities
- [ ] Tool integration examples working for all 6 tools
- [ ] Forge orchestration documentation matches implementation
- [ ] All examples tested and validated
- [ ] Consistent formatting and structure across files

### Quality Assurance
- [ ] Tool commands tested and working
- [ ] Examples produce expected results
- [ ] Cross-references between files are accurate
- [ ] No outdated information remains

This plan provides a structured approach to updating all CLI documentation while maximizing the combined capabilities of the available tools and maintaining resumable progress tracking.