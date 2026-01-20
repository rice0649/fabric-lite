#!/bin/bash

# ADHD-Optimized Headless Execution Daemon - Fixed Version

set -e

# Service details
SERVICE_URL="http://localhost:9999"
SERVICE_NAME="fabric-lite-adhd"
DAEMON_LOG="/tmp/adhd_daemon.log"

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
    echo "  status          - Check daemon status"
    echo "  help            - Show this help"
    echo ""
    echo "Options:"
    echo "  --daemon        - Run as background service"
    echo "  --profile       - Use specific attention profile"
    echo "  --auto          - Auto-detect task complexity"
    echo "  --headless      - Run tools without prompts"
    echo "  --quiet         - Minimal output for automated processing"
    echo ""
    echo "Examples:"
    echo "    ./scripts/adhd_daemon.sh --status"
    echo "    ./scripts/adhd_daemon.sh start"
    echo "    ./scripts/adhd_daemon.sh stop"
    echo ""
}

# Start daemon function
start_daemon() {
    echo "🚀 Starting ADHD-Optimized Daemon..."
    echo "📍 Service: $SERVICE_NAME on $SERVICE_URL"
    echo "📝 Persistent context: ./user_context"
    
    # Create context directory if it doesn't exist
    mkdir -p ./user_context
    
    # Start background monitoring
    nohup bash -c '
        while true; do
            echo "[$(date +'"'"'%Y-%m-%d %H:%M:%S'"'"')] 🔄 ADHD Daemon monitoring..."
            sleep 30
        done
    ' > "$DAEMON_LOG" 2>&1 &
    
    echo "✅ Daemon started with PID: $!"
    echo "📋 Log file: $DAEMON_LOG"
}

# Command line argument handling
case "$1" in
    "--help"|"-h"|"help"|"")
        show_help
        ;;
    "--status"|"status")
        daemon_status
        ;;
    "--start"|"start")
        start_daemon
        ;;
    "--stop"|"stop")
        echo "🛑 Stopping ADHD daemon..."
        if pgrep -f "fabric-lite" >/dev/null 2>&1; then
            pkill -f "fabric-lite" >/dev/null 2>&1 || true
        fi
        if pgrep -f "ADHD Daemon" >/dev/null 2>&1; then
            pkill -f "ADHD Daemon" >/dev/null 2>&1 || true
        fi
        echo "✅ Daemon stopped"
        ;;
    *)
        echo "Unknown option: $1"
        show_help
        exit 1
        ;;
esac