# YouTube Analyzer for Attention Deficit Disorder Management

This tool analyzes YouTube video transcripts and provides structured summaries specifically designed for ADHD workflows and attention management.

## 🎯 Core Features

### 1. Multi-Modal Transcript Analysis
- Extract and structure transcript data (speakers, topics, key phrases)
- Identify attention-demanding vs attention-friendly segments
- Provide time-stamped content summaries

### 2. ADHD-Friendly Summaries
- Bullet-point key information (reduces cognitive load)
- Executive summary for quick scanning
- Actionable suggestions tailored to attention profiles

### 3. Attention Management Tools
- Focus mode suggestions based on content analysis
- Break recommendations into manageable steps
- Distraction-free workflow recommendations

## 📋 Implementation

```python
import re
import json
import sys
from pathlib import Path

class YouTubeAnalyzer:
    def __init__(self):
        self.transcript_dir = Path("./transcripts")
        self.output_dir = Path("./summaries")
        self.focus_keywords = [
            "attention", "focus", "key point", "summary", "action item",
            "distraction", "overwhelming", "complex", "confusing"
        ]
        
    def analyze_transcript(self, video_id: str, transcript_path: str):
        """Analyze transcript for ADHD-friendly structured output."""
        
        # Read transcript
        with open(transcript_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # Basic structure
        structure = self._parse_transcript_structure(content)
        
        # Analysis
        analysis = self._analyze_for_adhd(structure)
        
        # Generate outputs
        self._save_outputs(video_id, structure, analysis)
        
        return analysis
    
    def _parse_transcript_structure(self, content: str) -> dict:
        """Parse transcript into structured format."""
        segments = []
        current_speaker = None
        current_time = 0
        
        lines = content.strip().split('\n')
        for line in lines:
            if line.strip():
                if '[' in line and ']' in line:
                    # Timestamp
                    current_time = self._extract_time(line)
                elif ':' in line and not line.startswith('['):
                    # Content with speaker
                    speaker = line.split(':')[0].strip()
                    content = ':'.join(line.split(':')[1:]).strip()
                    
                    segments.append({
                        'speaker': speaker,
                        'content': content.strip(),
                        'timestamp': current_time,
                        'type': 'speech',
                        'attention_level': self._assess_attention_demand(content)
                    })
                elif line.strip():
                    # Continuation of previous speaker
                    segments.append({
                        'speaker': current_speaker,
                        'content': line.strip(),
                        'timestamp': current_time,
                        'type': 'speech',
                        'attention_level': self._assess_attention_demand(content)
                    })
        
        return {
            'video_id': video_id,
            'segments': segments,
            'total_duration': current_time,
            'speakers': list(set(seg['speaker'] for seg in segments)),
            'attention_profile': self._create_attention_profile(segments)
        }
    
    def _analyze_for_adhd(self, structure: dict) -> dict:
        """Analyze content for ADHD-friendly output."""
        segments = structure['segments']
        high_attention_segments = []
        low_attention_segments = []
        total_words = 0
        
        for segment in segments:
            attention_level = segment.get('attention_level', 1)
            word_count = len(segment['content'].split())
            total_words += word_count
            
            if attention_level >= 3:  # High attention demand
                high_attention_segments.append(segment)
            else:  # Low attention demand
                low_attention_segments.append(segment)
        
        # Focus suggestions
        focus_areas = self._identify_focus_areas(high_attention_segments)
        
        return {
            'attention_profile': self._create_attention_profile(segments),
            'focus_areas': focus_areas,
            'high_attention_segments': len(high_attention_segments),
            'low_attention_segments': len(low_attention_segments),
            'total_words': total_words,
            'recommendations': self._generate_recommendations(focus_areas)
        }
    
    def _assess_attention_demand(self, text: str) -> int:
        """Assess cognitive load of content (1-5 scale)."""
        complexity_indicators = ['complex', 'technical', 'abstract', 'rapid']
        indicators_found = sum(1 for indicator in complexity_indicators if indicator.lower() in text.lower())
        
        # Base score on length and complexity
        base_score = min(len(text.split()) / 50, 2) + indicators_found
        
        return min(base_score, 5)
    
    def _create_attention_profile(self, segments: list) -> dict:
        """Create attention profile for ADHD management."""
        attention_levels = {
            1: 'Very Low',
            2: 'Low', 
            3: 'Moderate',
            4: 'High',
            5: 'Very High'
        }
        
        if not segments:
            return {'overall_level': 1, 'distribution': {}}
        
        level_counts = {}
        for segment in segments:
            level = segment.get('attention_level', 1)
            level_counts[level] = level_counts.get(level, 0) + 1
        
        return {
            'overall_level': max(level_counts.keys(), key=level_counts.get),
            'distribution': level_counts,
            'total_segments': len(segments)
        }
    
    def _identify_focus_areas(self, high_attention_segments: list) -> list:
        """Identify key focus areas from high-attention segments."""
        if not high_attention_segments:
            return []
        
        focus_words = []
        for segment in high_attention_segments:
            words = re.findall(r'\b(word)s?\b', segment['content'], re.IGNORECASE)
            focus_words.extend([word.lower() for word in words if len(word) > 4])
        
        # Remove duplicates and count frequency
        unique_focus_words = list(set(focus_words))
        word_frequency = {}
        for word in unique_focus_words:
            word_frequency[word] = focus_words.count(word)
        
        # Sort by frequency
        sorted_words = sorted(unique_focus_words, key=lambda x: word_frequency[x], reverse=True)
        
        return sorted_words[:10]  # Top 10 focus areas
    
    def _generate_recommendations(self, focus_areas: list) -> list:
        """Generate ADHD-friendly recommendations."""
        recommendations = []
        
        for area in focus_areas:
            if 'technical' in area.lower():
                recommendations.append({
                    'area': area,
                    'strategy': 'Use pomodoro technique (25min focus, 5min break)',
                    'tools': 'Screen recording app, website blockers',
                    'environment': 'Quiet workspace, noise-cancelling headphones'
                })
            elif 'planning' in area.lower():
                recommendations.append({
                    'area': area,
                    'strategy': 'Break into 15-minute focused work sessions',
                    'tools': 'Visual planning board, mind mapping software',
                    'environment': 'Meeting-free zones, structured task lists'
                })
            elif 'learning' in area.lower():
                recommendations.append({
                    'area': area,
                    'strategy': 'Use multi-modal learning (video + audio)',
                    'tools': 'Interactive note-taking, spaced repetition apps',
                    'environment': 'Adaptive learning platforms with speed control'
                })
        
        return recommendations
    
    def _save_outputs(self, video_id: str, structure: dict, analysis: dict):
        """Save analysis outputs in structured format."""
        import json
        
        # Create output directory
        self.output_dir.mkdir(parents=True, exist_ok=True)
        
        # Save structured analysis
        with open(self.output_dir / f"{video_id}_structure.json", 'w', encoding='utf-8') as f:
            json.dump(structure, f, indent=2)
        
        # Save ADHD-friendly summary
        summary = self._create_adhd_summary(analysis)
        with open(self.output_dir / f"{video_id}_summary.md", 'w', encoding='utf-8') as f:
            f.write(summary)
        
        # Save focus areas
        with open(self.output_dir / f"{video_id}_focus_areas.json", 'w', encoding='utf-8') as f:
            json.dump(analysis['focus_areas'], f, indent=2)
        
        print(f"Analysis saved for {video_id}")
    
    def _create_adhd_summary(self, analysis: dict) -> str:
        """Create ADHD-optimized content summary."""
        summary_parts = []
        
        # Executive summary
        summary_parts.append("## 🎯 QUICK OVERVIEW\n")
        summary_parts.append(f"Video ID: {analysis['video_id']}\n")
        summary_parts.append(f"Attention Profile: Level {analysis['attention_profile']['overall_level']}/5\n")
        
        # Key focus areas
        if analysis['focus_areas']:
            summary_parts.append("\n## 🎯 KEY FOCUS AREAS\n")
            for i, area in enumerate(analysis['focus_areas'][:5], 1):
                summary_parts.append(f"{i+1}. **{area.upper()}**\n")
                summary_parts.append(f"   Strategy: {area['strategy']}\n")
                summary_parts.append(f"   Tools: {', '.join(area['tools'])}\n")
        
        # Attention recommendations
        if analysis['recommendations']:
            summary_parts.append("\n## 📋 ADHD-SPECIFIC RECOMMENDATIONS\n")
            for rec in analysis['recommendations'][:3]:
                summary_parts.append(f"• {rec['area']}: {rec['strategy']}\n")
        
        # ADHD-friendly action items
        summary_parts.append("\n## 🚀 IMMEDIATE ACTIONS\n")
        if analysis['high_attention_segments'] > 3:
            summary_parts.append("⚠️  High cognitive load detected - consider breaking into smaller sessions")
        else:
            summary_parts.append("✅  Manageable cognitive load - proceed with focused work")
        
        return '\n'.join(summary_parts)

def main():
    if len(sys.argv) < 3:
        print("Usage: python yt_analyzer.py <video_id> <transcript_file>")
        sys.exit(1)
    
    video_id = sys.argv[1]
    transcript_path = sys.argv[2]
    
    analyzer = YouTubeAnalyzer()
    analysis = analyzer.analyze_transcript(video_id, transcript_path)
    
    print(f"\n🎯 ANALYSIS COMPLETE: {video_id}")
    print(f"📊 Files saved to: {analyzer.output_dir}")
    print(f"🧠 ADHD Profile: Level {analysis['attention_profile']['overall_level']}/5")

if __name__ == "__main__":
    main()
```