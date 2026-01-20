# Persona Creation Workflow for Attention-Considerate Users

## 🎯 Purpose
Helps users with ADHD or other attention challenges make informed decisions about project setup by understanding their personal needs, preferences, and working styles before they begin development.

## 📋 Workflow Overview

1. **Personal Assessment** - User preferences, work style, attention profile
2. **Project Scoping** - Project complexity, team structure, timeline
3. **Tool Selection** - Optimal AI provider configuration
4. **Environment Setup** - Development environment optimization
5. **Workflow Creation** - Customized development plan with ADHD-optimized settings

## 🔧 Interactive Persona Creator

```bash
create_attention_persona() {
    echo "🎯 Creating personalized fabric-lite setup..."
    echo ""
    
    # Assess working style
    echo "🔍 How do you prefer to work?"
    echo "1) 🔥 Highly structured with clear routines and minimal distractions"
    echo "2) 🎨 Moderately structured with some flexibility"
    echo "3) 🌊 Flexible with creative freedom"
    echo "4) 🧘 Minimal structure, high autonomy"
    echo ""
    echo "Enter your choice (1-4): "
    
    read -p working_style
    echo ""
    
    case "$working_style" in
        1)
            echo "✅ Selected: Highly structured approach"
            echo "Creating ADHD-optimized configuration..."
            ;;
        2)
            echo "✅ Selected: Moderately structured approach"
            echo "Creating balanced configuration..."
            ;;
        3)
            echo "✅ Selected: Flexible with creative freedom"
            echo "Creating flexible configuration..."
            ;;
        4)
            echo "✅ Selected: Minimal structure with high autonomy"
            echo "Creating autonomous configuration..."
            ;;
        *)
            echo "❌ Invalid choice. Please select 1-4."
            return 1
    esac
    
    echo "🎍 Please provide details about your work:"
    echo ""
    echo "1. Project scope and complexity:"
    echo "   - What type of project are you working on?"
    echo "   - How complex is the codebase?"
    echo "   - Are you working alone or in a team?"
    echo "   - What are your main challenges in development?"
    
    echo "2. Attention and focus needs:"
    echo "   - How easily distracted do you get?"
    echo "   - What time of day are you most productive?"
    echo "   - Do you prefer working in focused blocks or with background music?"
    echo "   - Are there any specific sensory needs we should accommodate?"
    
    echo "3. Working environment preferences:"
    echo "   - What environment setup helps you focus best?"
    echo "   - Do you prefer visual, auditory, or kinesthetic learning?"
    
    echo "4. AI tool experience:"
    echo "   - Which AI tools have you found most helpful?"
    echo "   - How do you prefer to receive AI assistance?"
    echo "   - Have you tried tools specifically designed for ADHD?"
    
    echo "5. Persona preferences:"
    echo "   - Do you prefer formal or casual communication style?"
    echo "   - Should the AI have personality or be neutral?"
    
    echo ""
    read -p "Enter your preferences (comma-separated, or 'done' when finished): "
    
    # Parse preferences
    local structured=true
    working_style=""
    creative_freedom=false
    formal=false
    ai_personality=false
    
    IFS="," # Read comma-separated list
    
    while [ -n "$structured" ]; do
        case "$preference" in
            1)  ;; 2) ;; 3) ;; 4) ;; *) structured=false ;; esac
        read -r "$preference"
        
        if [ "$structured" = "true" ]; then
            structured=true
            echo "   ✅ Structured preferences detected"
        fi
        
        echo "$preference" | while IFS=',' read -ra preference; do
            preference="${preference%% *}"
            echo "Processing preference: $preference"
            
            case "$preference" in
                "working_style") working_style="$1" ;;
                "creative_freedom") creative_freedom="$1" ;;
                "formal") formal="$2" ;;
                *) structured=false ;;
            esac
        done
    done
    
    # Show summary
    echo ""
    echo "📋 Personal Assessment Summary:"
    echo "   Working Style: Option $working_style"
    if [ "$creative_freedom" = "1" ]; then
        echo "   🎨 Creative freedom enabled"
    elif [ "$formal" = "2" ]; then
        echo "   🏛️ Formal approach enabled"
    else
        echo "   🌊 Balanced approach selected"
    fi
    
    echo ""
    echo "   AI Personaility: Option $ai_personality"
    if [ "$ai_personality" = "1" ]; then
        echo "   🤖 Friendly persona with warmth and encouragement"
    else
        echo "   🤖 Neutral persona with professional assistance"
    fi
    
    echo ""
    echo "   Creative Freedom: Option $creative_freedom"
    if [ "$creative_freedom" = "1" ]; then
        echo "   🎨 Maximum creativity with flexibility"
    else
        echo "   🌊 Structured approach selected"
    fi
    
    echo ""
    echo "   📝 ADHD Profile Analysis:"
    echo "   Based on preferences: $working_style, $structured"
    
    # Generate personalized profile
    case "$working_style" in
        "$1")
            echo "   🎯 ADHD-Optimized Profile: High focus, structured routines, 25-minute blocks"
            ;;
        "$2") 
            echo "   🎨 ADHD-Optimized Profile: Balanced approach, moderate complexity"
            ;;
        "$3") 
            echo "   🌊 Flexible Profile: Creative freedom with some structure"
            ;;
        "$4") 
            echo "   🧘 Autonomous Profile: High autonomy, minimal structure"
            ;;
    esac
    
    echo ""
    echo "   📋 Profile Features:"
    echo "   • Context awareness and adaptation"
    echo "   • Automatic tool optimization for task complexity"
    echo "   • Focus area prioritization and workflow selection"
    echo "   • Multi-tool coordination capabilities"
    echo "   • Background processing with progress tracking"
    
    echo ""
    echo "📋 Personalization Options:"
    echo "   • Visual organization (color-coded focus areas)"
    echo "   • Workflow templates for common project types"
    echo "   • Sensory accommodations if needed"
    echo "   • Notification preferences and timing management"
    
    # Create profile configuration
    echo "🔧 Generating personalized configuration..."
    
    python3 tools/create_attention_profile.py \
        --working-style "$working_style" \
        --structured "$structured" \
        --creative_freedom "$creative_freedom" \
        --formal "$formal" \
        --ai_personality "$ai_personality" \
        --project_scope "general" \
        --complexity "moderate" \
        --team_size "small"
        --adhd_focused true
    
    echo "✅ Personalized configuration created successfully!"
    echo "💾 Profile saved to: ./user_context/attention_profile.json"
    
    echo "🔧 ADHD-Optimized setup complete!"
}

# Start monitoring with the new profile
echo "🚀 Starting daemon with personalized ADHD-optimized settings..."
    
    # Apply new profile to existing context (if any)
    if [ -f "./user_context/attention_profile.json" ]; then
        export FABRIC_LITE_PROFILE="$user_context/attention_profile.json"
        echo "📋 Applying ADHD-optimized profile..."
    else
        echo "📋 Creating new default profile..."
        create_attention_profile --default
        export FABRIC_LITE_PROFILE="$user_context/attention_profile_default.json"
        echo "✅ Default profile created"
    fi
    
    # Start daemon with optimized settings
    ./scripts/adhd_daemon.sh --daemon --profile adhd_focused
    
    echo "🎉 System ready for ADHD-optimized workflows!"
    echo ""
    echo "🔧 Commands: start | stop | status | task | analyze | help"
    echo "🔧 Profile: $(cat $FABRIC_LITE_PROFILE) | jq -r '.name')"
    echo "🎯 Daemon running with: $(pgrep -f "fabric-lite" >/dev/null 2>&1 && echo $! || echo "Process ID: $!")"
}

echo ""
echo "🎯 ADHD-Optimized Fabric-Lite v2.0 is ready for deployment!"
echo ""
echo "💡 Enter a command to begin (or 'help' for available options): "
