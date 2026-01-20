# Attention-Aware Context Structure

This directory structure allows fabric-lite to automatically adapt to users' attention profiles and maintain context across different cognitive states and workflows.

## 📂 Directory Organization

```
/user_context/
├── profiles/
│   ├── adhd_focused/
│   ├── standard/
│   ├── low_cognitive_load/
│   └── high_stimulation/
├── workflows/
│   ├── research/
│   ├── implementation/
│   ├── testing/
│   └── deployment/
├── templates/
│   ├── focused/
│   ├── detailed/
│   └── quick_actions/
└── personal/
    ├── ideas.md
    ├── journal_entries/
    └── custom_prompts.md
```

## 🧠 Profile-Based Configuration

Each profile automatically customizes:
- **Interface complexity** based on cognitive load
- **Tool selection strategies** optimal for attention level
- **Workflow patterns** sequential vs parallel based on focus capacity
- **Communication style** concise vs detailed based on preference
- **Progress tracking frequency** adaptive based on task complexity

## 📋 Context-Aware Features

### **Profile Detection**
```bash
# Automatic profile selection (simplified)
analyze_current_state() {
    local cognitive_load=$(detect_cognitive_load)
    current_focus_areas=$(identify_primary_focus_areas)
    
    case "$local_cognitive_load" in
        echo "Profile: Low cognitive load detected"
        echo "Strategy: Sequential tasks, minimal distractions"
        ;;
    "$high_stimulation" in
        echo "Profile: High stimulation detected"
        echo "Strategy: Structured workflows with clear boundaries"
        ;;
    *)
        echo "Profile: Standard cognitive load"
        echo "Strategy: Balanced approach with moderate complexity"
        ;;
    esac
}
```

### **Workflow Adaptation**
```bash
# Context-aware tool selection
adapt_workflow_for_profile() {
    profile="$1"
    local tools=("")
    
    case "$profile" in
        "adhd_focused") tools="codex gemini claude" ;;
        "standard") tools="codex claude" ;;
        "low_cognitive_load") tools="codex ollama" ;;
        "high_stimulation") tools="codex" ;;
    esac
    
    for tool in $tools; do
        echo "Initializing $tool with ADHD-optimized settings..."
        fabric-lite run "$tool" --profile "$profile"
    done
}
```

## 📝 Implementation Guidelines

### **For Attention Deficit Disorder**
1. **Chunking**: Break large tasks into 15-minute focused blocks
2. **Visual Timers**: Use integrated timers for work/break cycles
3. **Progress Visualization**: Clear progress bars with milestone markers
4. **Flexibility**: Allow mid-task tool switching without losing context

### **For Development Teams**
1. **Shared Context Files**: All tools write to the same context files
2. **Profile-Based Development**: Each developer works with their preferred ADHD profile
3. **Attention-Friendly Code Reviews**: Focus on structure, clarity, and executive function
4. **User Story Integration**: Include user workflows as test cases in development

This structure makes fabric-lite naturally adaptive to different cognitive states while maintaining professional development workflows.