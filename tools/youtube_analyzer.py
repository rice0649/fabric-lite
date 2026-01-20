#!/usr/bin/env python3
"""
YouTube Analyzer for Attention Management

Analyzes YouTube video transcripts and provides structured summaries 
specifically designed for ADHD workflows and attention management.
"""

import argparse
import json
import re
import sys
from datetime import datetime
from pathlib import Path

def parse_srt_file(srt_file):
    """Parse SRT subtitle file and extract transcript with timestamps."""
    try:
        with open(srt_file, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # Split into subtitle blocks
        blocks = re.split(r'\n\s*\n', content.strip())
        transcript = []
        
        for block in blocks:
            lines = block.strip().split('\n')
            if len(lines) >= 3:
                # Extract timestamp and text
                timestamp = lines[1]
                text = ' '.join(lines[2:])
                transcript.append({
                    'timestamp': timestamp,
                    'text': text
                })
        
        return transcript
    except Exception as e:
        print(f"Error parsing SRT file: {e}")
        return None

def analyze_for_adhd(transcript):
    """Analyze transcript for ADHD-friendly content structure."""
    if not transcript:
        return None
    
    # Extract key information
    total_duration = len(transcript)
    key_topics = []
    actionable_items = []
    
    # Simple keyword-based analysis
    for entry in transcript:
        text = entry['text'].lower()
        
        # Look for action items
        if any(word in text for word in ['should', 'must', 'need to', 'important', 'remember']):
            actionable_items.append({
                'timestamp': entry['timestamp'],
                'content': entry['text']
            })
        
        # Look for topic indicators
        if any(word in text for word in ['first', 'second', 'finally', 'conclusion', 'summary']):
            key_topics.append({
                'timestamp': entry['timestamp'],
                'content': entry['text']
            })
    
    return {
        'total_segments': total_duration,
        'key_topics': key_topics[:5],  # Limit to top 5
        'actionable_items': actionable_items[:10],  # Limit to top 10
        'analysis_timestamp': datetime.now().isoformat()
    }

def generate_adhd_summary(analysis):
    """Generate ADHD-friendly summary."""
    if not analysis:
        return "No analysis available"
    
    summary = []
    summary.append("# 🎯 Executive Summary")
    summary.append(f"Video analyzed with {analysis['total_segments']} transcript segments")
    summary.append(f"Found {len(analysis['key_topics'])} key topics and {len(analysis['actionable_items'])} actionable items")
    summary.append("")
    
    if analysis['key_topics']:
        summary.append("## 🔑 Key Topics")
        for i, topic in enumerate(analysis['key_topics'], 1):
            summary.append(f"{i}. [{topic['timestamp']}] {topic['content']}")
        summary.append("")
    
    if analysis['actionable_items']:
        summary.append("## ✅ Action Items")
        for i, item in enumerate(analysis['actionable_items'], 1):
            summary.append(f"{i}. [{item['timestamp']}] {item['content']}")
        summary.append("")
    
    summary.append("## 💡 ADHD Recommendations")
    summary.append("• Review key topics first for context")
    summary.append("• Focus on actionable items for implementation")
    summary.append("• Use timestamps to jump to important sections")
    
    return '\n'.join(summary)

def main():
    parser = argparse.ArgumentParser(
        description='Analyze YouTube transcripts for ADHD-friendly summaries'
    )
    parser.add_argument(
        'srt_file',
        help='Path to SRT transcript file'
    )
    parser.add_argument(
        '--output', '-o',
        help='Output file for analysis results'
    )
    parser.add_argument(
        '--json',
        action='store_true',
        help='Output results in JSON format'
    )
    
    args = parser.parse_args()
    
    # Parse transcript
    transcript = parse_srt_file(args.srt_file)
    if not transcript:
        sys.exit(1)
    
    # Analyze for ADHD
    analysis = analyze_for_adhd(transcript)
    if not analysis:
        sys.exit(1)
    
    # Generate output
    if args.json:
        output = json.dumps(analysis, indent=2)
    else:
        output = generate_adhd_summary(analysis)
    
    # Save or print
    if args.output:
        with open(args.output, 'w', encoding='utf-8') as f:
            f.write(output)
        print(f"Analysis saved to {args.output}")
    else:
        print(output)

if __name__ == '__main__':
    main()