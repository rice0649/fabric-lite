#!/bin/bash

# ADHD-Optimized Headless Execution Daemon - Simplified Version

set -e

# Service details
SERVICE_URL="http://localhost:9999"
SERVICE_NAME="fabric-lite-adhd"
DAEMON_LOG="/tmp/adhd_daemon.log"

echo "🚀 Starting ADHD-optimized daemon..."
echo "📡 Service: $SERVICE_NAME on $SERVICE_URL"
echo "📝 Persistent context directory: ./user_context"

# Tool execution function - simplified
run_tool_headless() {
    local tool="$1"
    local prompt="$2"
    
    # Basic logging
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] 🔄 Running $tool with ADHD optimization"
    
    # Execute tool (simplified)
    ./bin/fabric-lite run "$tool" -P "$prompt" >/dev/null 2>&1 &
    
    # Capture PID
    local pid=$!
    
    # Wait for completion (simplified)
    wait $!
    
    # Save context
    wait $!
    if [ -f "/tmp/tool_output_${tool}_${pid}" ]; then
        user_context=$(cat "/tmp/tool_output_${tool}_${pid}")
        timestamp=$(date +"%Y-%m-%d %H:%M:%S")
        echo "📝 Capturing user context at $timestamp"
        echo "Tool: $tool" >> "./user_context/session_${timestamp}.md"
        echo "Prompt: $prompt" >> "./user_context/session_${timestamp}.md"
        echo "---" >> "./user_context/session_${timestamp}.md"
        cat "/tmp/tool_output_${tool}_${pid}" >> "./user_context/session_${timestamp}.md"
        
        echo "✅ $tool execution completed and context saved"
    else
        echo "❌ Tool execution failed"
    fi
    
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] 🎯 $tool execution result: $?"
}

# Background transcript monitoring
monitor_transcripts() {
    echo "👀 Monitoring for new YouTube transcripts..."
    
    # Check for new transcript files every 10 seconds
    while true; do
        new_files=0
        for file in ./transcripts/*.srt; do
            if [ ! -f "./user_context/session_$(date -d '1 day ago' +%_*)" ]; then
                echo "📥 New transcript detected: $(basename "$file")"
                ((new_files++))
                echo "🎥 Processing new transcript: $file..."
                python3 tools/youtube_analyzer.py "process_youtube_content" "$(basename "$file")" --context "./user_context/" 2>/dev/null &
                wait $!
                echo "📝 Analysis saved for $(basename "$file")"
            fi
        done
        
        if [ $new_files -gt 0 ]; then
            echo "📊 Found $new_files new transcript(s)"
        fi
        
        sleep 10
    done &
}

# Auto-detects task types
auto_detect_task_type() {
    task_type="$1"
    
    # Simple heuristics
    if [[ "$1" == *"implement"* ]] || [[ "$1" == *"code"* ]]; then
        task_type="complex"
    elif [[ "$1" == *"research"* ]] || [[ "$1" == *"plan"* ]] || [[ "$1" == *"test"* ]]; then
        task_type="testing"
    elif [[ "$1" == *"analyze"* ]] || [[ "$1" == *"review"* ]]; then
        task_type="validation"
    else
        task_type="general"
    fi
    
    echo "$task_type"
}

# Main daemon with simplified task execution
main() {
    echo "🚀 ADHD-Optimized Daemon v2.0"
    echo "📍 Service: $SERVICE_NAME on $SERVICE_URL"
    echo "📝 Persistent context: ./user_context"
    echo "🔄 Auto-detecting task types..."
    
    # Start transcript monitoring
    monitor_transcripts &
    
    # Main control loop
    while true; do
        # Display status
        daemon_status
        
        # Get user input
        echo ""
        echo "🤖 Available Commands:"
        echo "  start   - Start daemon"
        echo "  stop    - Stop daemon"
        echo "  status   - Check daemon status"
        echo "  task    - Get next task (auto-detects)"
        echo "  analyze  - Analyze transcript content"
        echo "  help    - Show this help"
        echo "  run     - Execute tool with ADHD-optimized settings"
        echo ""
        echo "💡 Options:"
        echo "    --daemon    Run as background service"
        echo "    --profile   - Use specific attention profile"
        echo "    --auto     - Auto-detect task complexity"
        echo "    --headless  - Run tools without prompts"
        echo ""
        
        read -p "Enter command: " next_task
        
        # Execute command
        case "$next_task" in
            "start")
                start_daemon
                ;;
            "stop")
                if pgrep -f "fabric-lite" >/dev/null 2>&1; then
                    pkill -f "fabric-lite" >/dev/null 2>&1 || true
                    echo "🛑 Daemon stopped"
                else
                    echo "✅ Daemon already stopped"
                fi
                ;;
            "status")
                daemon_status
                ;;
            "task")
                local task_type
                if [ -n "$1" ]; then
                    echo "🔄 Auto-detecting task complexity..."
                    auto_detect_task_type
                else
                    echo "📝 Task complexity: general"
                fi
                get_next_task
                ;;
            "analyze")
                echo "📊 Analyzing with fabric-lite tools..."
                # In real implementation, this would use the ADHD analysis tools
                ;;
            "run")
                # Simplified execution without auto-detection
                run_tool_headless "$2" "Implement ADHD-optimized workflow system based on transcript"
                ;;
            "help")
                echo "🎯 Available Commands:"
                echo "  start   - Start daemon"
                echo "  stop    - Stop daemon"
                echo "  status   - Check daemon status"
                echo "  task    - Get next task (auto-detects)"
                echo "  analyze  - Analyze content with fabric-lite tools"
                echo "  help    - Show this help"
                echo "  run     - Execute tool with ADHD-optimized settings"
                echo "  --daemon    Run as background service"
                echo "  --profile   - Use specific attention profile"
                echo "  --auto     - Auto-detect task complexity"
                echo "  --headless - Run tools without prompts"
                echo "    --quiet    - Minimal output for automated processing"
                echo ""
                echo "💡 Examples:"
                echo "    # Start with optimized daemon"
                echo "        ./scripts/adhd_daemon.sh --daemon --auto"
                echo "    # Check status"
                echo "        ./scripts/adhd_daemon.sh --status"
                echo "        # Process specific YouTube content"
                echo "        ./scripts/adhd_daemon.sh --task 'research' --file VIDEO_ID.srt'"
                echo ""
                ;;
        esac
        ;;
    esac
done &
}

# Service status function
daemon_status() {
    if pgrep -f "fabric-lite" >/dev/null 2>&1; then
        echo "✅ Fabric-Lite daemon is running"
    else
        echo "❌ Fabric-Lite daemon is not running"
    fi
}

# Help function
show_help() {
    echo "🎯 ADHD-Optimized Daemon v2.0 Help"
    echo ""
    echo "This service automatically:"
    echo "• Monitors YouTube transcripts for new content"
    echo "• Auto-detects optimal AI tools for each task"
    echo "• Provides ADHD-optimized workflows with appropriate tools"
    echo "• Maintains persistent user context across sessions"
    echo "• Uses fabric-lite's enhanced tool coordination system"
    echo ""
    echo "Available Commands:"
    echo "  start           - Start daemon in background"
    echo "  stop            - Stop background daemon"
    echo "  status           - Check daemon status"
    echo "  task            - Get next task (auto-detects)"
    echo "  analyze           - Analyze transcript with fabric-lite tools"
    echo "  run             - Execute tool with settings"
    echo "  help             - Show this help"
    echo ""
    echo "Options:"
    echo "  --daemon         - Run as background service"
    echo "  --profile        - Use specific attention profile"
    echo "  --auto          - Auto-detect task complexity"
    echo "  --headless       - Run tools without prompts"
    echo "  --quiet          - Minimal output for automated processing"
    echo ""
    echo "Examples:"
    echo "    # Start optimized daemon"
    echo "        # Check status"
    echo "        # Process YouTube content"
    echo "        ./scripts/adhd_daemon.sh --task 'research' --file VIDEO_ID.srt'"
    echo "        # Status check"
    echo "        # Complete task with codex"
    echo "        ./scripts/adhd_daemon.sh --task 'implement' --file TRANSCRIPT.srt'"
    echo ""
    echo ""
}

# Set executable
chmod +x scripts/adhd_daemon.sh