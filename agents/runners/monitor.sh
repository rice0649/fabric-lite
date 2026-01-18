#!/bin/bash
# Monitor Agent Outputs
# Run this in your main terminal to watch progress

PROJECT_DIR="/home/oak38/projects/fabric-lite"
OUTPUT_DIR="$PROJECT_DIR/agents/outputs"

echo "============================================"
echo "  Agent Output Monitor"
echo "============================================"
echo ""
echo "Watching: $OUTPUT_DIR"
echo "Press Ctrl+C to stop"
echo ""
echo "============================================"
echo ""

# Function to show status of all outputs
show_status() {
    echo "=== Agent Outputs Status ==="
    echo ""
    for f in "$OUTPUT_DIR"/*.md; do
        if [ -f "$f" ]; then
            filename=$(basename "$f")
            status=$(grep -m1 "^status:" "$f" 2>/dev/null | cut -d: -f2 | tr -d ' ' || echo "unknown")
            lines=$(wc -l < "$f")
            echo "  $filename: $status ($lines lines)"
        fi
    done
    echo ""
    echo "=== Queue Status ==="
    if [ -f "$PROJECT_DIR/agents/queue/patterns_to_implement.md" ]; then
        patterns=$(grep -c "^[0-9]" "$PROJECT_DIR/agents/queue/patterns_to_implement.md" 2>/dev/null || echo "0")
        echo "  Patterns queued: $patterns"
    else
        echo "  No patterns queued yet"
    fi
    echo ""
}

# Watch for changes
while true; do
    clear
    show_status

    # Show final summary if it exists
    if [ -f "$OUTPUT_DIR/FINAL_SUMMARY.md" ]; then
        echo "=== FINAL SUMMARY AVAILABLE ==="
        echo ""
        head -50 "$OUTPUT_DIR/FINAL_SUMMARY.md"
        echo ""
        echo "[... see full file for more ...]"
    fi

    sleep 5
done
