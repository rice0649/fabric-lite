# Fabric-Lite User Manual - Your First AI Assistant

**Welcome to the world of local AI tools!** 🎯

This manual is written for **complete beginners** - people who have never used command line tools or AI before. No technical knowledge required!

---

## 📚 Chapter 1: What You Need to Know First

### What is Fabric-Lite?
Think of fabric-lite as a **translator** between your thoughts and powerful AI assistants. You write or speak normally, and fabric-lite makes AI tools understand and respond helpfully.

**Simple Example:**
- You: "I need to summarize this meeting transcript"
- Fabric-Lite + AI: Gives you a clean, organized summary

### Why Local AI?
🔒 **Your Privacy Matters**: Everything happens on YOUR computer
- Your conversations never go to the internet
- Your ideas stay private and secure
- You're in complete control

### What You'll Accomplish
✅ **Summarize** long documents into key points  
✅ **Extract Ideas** from brainstorming sessions  
✅ **Explain Code** you don't understand  
✅ **Generate Content** like emails or reports  
✅ **Organize Thoughts** into structured formats

---

## 🚀 Chapter 2: Getting Started in 5 Minutes

### The Magic Setup Command
Open your terminal (that black box with text) and type:

```bash
curl -fsSL https://raw.githubusercontent.com/rice0649/fabric-lite/main/forgiving-setup.sh | bash
```

**What This Does:**
- 🔧 Installs Ollama (your local AI brain)
- 🏗 Builds fabric-lite (your command translator)
- 🧪 Tests everything works
- ✅ Sets you up for success

### What Happens During Setup:
```
Step 1: Checking your computer...
✅ Found: curl, terminal
✅ Found: space for AI tools

Step 2: Setting up AI assistant...
📦 Installing Ollama (local AI brain)
✅ Downloaded Ollama for your computer
✅ Started Ollama at http://localhost:11434

Step 3: Building fabric-lite...
🔨 Compiling your command translator
✅ Built successfully at ./bin/fabric-lite

Step 4: Testing everything...
🧪 Running quick tests
✅ All systems working!

🎉 YOU'RE READY TO USE AI TOOLS!
```

### First Test - Your First AI Conversation:
Type this immediately after setup:

```bash
./bin/fabric-lite run "Hello AI, can you help me?" --provider ollama
```

**Expected Response:**
```
Hello! I'd be happy to help you. What would you like to work on today? I can assist with:
- Answering questions
- Summarizing text
- Explaining concepts
- Generating ideas
- And much more!

What specific task can I help you with?
```

**🎯 CONGRATULATIONS!** You just had your first AI conversation entirely on your computer!

---

## 🎮 Chapter 3: Your First Day with AI Tools

### The Basic Recipe
All fabric-lite commands follow this pattern:

```bash
./bin/fabric-lite [what-you-want-to-do] --with-options
```

### Task 1: Summarize a Webpage
```bash
echo "https://example.com/interesting-article" | ./bin/fabric-lite run --pattern summarize --provider ollama
```

### Task 2: Extract Ideas from Notes
```bash
./bin/fabric-lite run --pattern extract_ideas --provider ollama << 'EOF
Meeting notes:
- Discussed new product features
- Customer wants faster delivery
- Budget constraints mentioned
- Team suggested subscription model
EOF
```

### Task 3: Explain Something Confusing
```bash
./bin/fabric-lite run --pattern explain_code --provider ollama << 'EOF'
Can someone explain why this code isn't working?
function mysteriousFunction() {
  // complicated code here
  return something unexpected
}
EOF
```

---

## 📖 Chapter 4: Understanding Patterns

Patterns are like **templates for AI conversations**. Each pattern is designed for a specific type of task.

### Available Patterns:
```bash
./bin/fabric-lite list
```

**You'll see:**
- `explain_code` - Ask AI to explain programming concepts
- `summarize` - Turn long text into short summaries
- `extract_ideas` - Find interesting insights in content
- `test` - Test your understanding of concepts

### Using Patterns in Real Life:

**For Students:**
```bash
# Summarize lecture notes
./bin/fabric-lite run --pattern summarize --provider ollama < lecture_notes.txt

# Extract study ideas
./bin/fabric-lite run --pattern extract_ideas --provider ollama < research_paper.pdf
```

**For Professionals:**
```bash
# Explain complex topics
./bin/fabric-lite run --pattern explain_code --provider ollama << 'EOF'
Explain quantum computing for business presentation
EOF

# Generate professional summaries
./bin/fabric-lite run --pattern summarize --provider ollama < meeting_transcript.txt
```

---

## 🛠️ Chapter 5: Working with Files

### Reading from Files
```bash
# Process a document
./bin/fabric-lite run --pattern summarize --provider ollama < my_document.txt

# Process multiple files
cat *.txt | ./bin/fabric-lite run --pattern summarize --provider ollama
```

### Saving AI Responses
```bash
# Save to a file
./bin/fabric-lite run --pattern summarize --provider ollama < document.txt > summary.txt

# Both input and output
./bin/fabric-lite run --pattern summarize --provider ollama < document.txt | tee summary.txt
```

---

## 🔧 Chapter 6: Customizing Your Experience

### Getting Better Responses
Add `--model llama3:8b` for more intelligent responses:
```bash
./bin/fabric-lite run --pattern summarize --model llama3:8b --provider ollama < text.txt
```

### Speed vs Quality Trade-off
- **llama3:2** = Fast, good for quick tasks
- **llama3:8b** = Slower, smarter, better for complex tasks

### Using Streaming (Real-time Responses)
Add `--stream` to see responses as they're generated:
```bash
./bin/fabric-lite run --pattern summarize --stream --provider ollama < long_document.txt
```

---

## 🚨 Chapter 7: Troubleshooting - When Things Go Wrong

### Problem: "Ollama model not found"
**Solution:**
```bash
# Download the model
ollama pull llama3.2

# Try again
./bin/fabric-lite run --pattern summarize --provider ollama < test.txt
```

### Problem: "fabric-lite command not found"
**Solution:**
```bash
# You're probably in wrong directory
cd /path/to/fabric-lite
./bin/fabric-lite --help

# Or add to your PATH (one time setup)
echo 'export PATH="$PATH:/path/to/fabric-lite"' >> ~/.bashrc
source ~/.bashrc
```

### Problem: "Response is very slow"
**Solution:**
```bash
# Use a smaller, faster model
./bin/fabric-lite run --pattern summarize --model llama3.2 --provider ollama < text.txt

# Or get more RAM (if you can)
ollama pull mistral  # Smaller, faster model
```

### Problem: "AI gives weird answers"
**Solution:**
```bash
# Be more specific in your request
# Bad: "summarize"
# Good: "summarize this business report into 3 key points"

./bin/fabric-lite run --pattern summarize --provider ollama << 'EOF'
Please summarize this business report into 3 key points for executives:
[Your business report content here]
EOF
```

---

## 💡 Chapter 8: Pro Tips and Tricks

### Tip 1: Create Shortcuts
Add to your `~/.bashrc`:
```bash
# Quick summarization shortcut
alias sum='./bin/fabric-lite run --pattern summarize --provider ollama'
alias explain='./bin/fabric-lite run --pattern explain_code --provider ollama'

# Now you can just type:
sum < document.txt
explain < some_code.py
```

### Tip 2: Use Environment Variables
```bash
# Set your preferred model
export FABRIC_LITE_MODEL=llama3:8b

# Use it everywhere
./bin/fabric-lite run --pattern summarize --provider ollama < text.txt
# Automatically uses your preferred model
```

### Tip 3: Batch Processing
```bash
# Process multiple files at once
for file in *.txt; do
  echo "Processing $file..."
  ./bin/fabric-lite run --pattern summarize --provider ollama < "$file" > "${file%.txt}_summary.txt"
done
```

### Tip 4: Save Good Conversations
```bash
# Save conversations you like
./bin/fabric-lite run --pattern extract_ideas --provider ollama < great_conversation.txt > saved_insights.txt
```

---

## 🎓 Chapter 9: Your AI Journey Continues

### What to Do Next
1. **Practice Daily** - Use fabric-lite for real tasks
2. **Explore Patterns** - Try all available patterns
3. **Combine Tools** - Use fabric-lite alongside your work
4. **Share Results** - Show others what AI can help you create

### Learning Resources
- **Project Website**: [GitHub repository link]
- **Community**: Join discussions and ask questions
- **Patterns**: Create your own custom patterns
- **Advanced**: Explore tool integration and automation

### Safety First
🔒 **Your Data Stays Local** - Nothing is sent to clouds
🧠 **You're in Control** - AI tools work only for you
💡 **Privacy Matters** - Your thoughts and ideas remain private

---

## 🆘 Chapter 10: Quick Reference

### Cheat Sheet
```bash
# Basic usage
./bin/fabric-lite run --pattern [name] --provider ollama < input

# List patterns
./bin/fabric-lite list

# Get help
./bin/fabric-lite --help

# Stream responses
./bin/fabric-lite run --pattern [name] --stream --provider ollama < input

# Save responses
./bin/fabric-lite run --pattern [name] --provider ollama < input > output.txt
```

### Pattern Quick Guide
- `summarize` - Make long text short
- `extract_ideas` - Find insights and opportunities  
- `explain_code` - Understand programming concepts
- `test` - Check your knowledge

### Common Flags
- `--provider ollama` - Use your local AI
- `--model llama3.2` - Use specific model
- `--stream` - See responses in real-time
- `--help` - Get help with any command

---

## 🎯 You're Ready!

**Congratulations!** 🎉 You now have everything you need to:
✅ Use AI tools privately and securely  
✅ Solve real problems with AI assistance  
✅ Learn and explore with AI guidance  
✅ Create better work with AI collaboration  

**Remember:**
- Start simple, practice daily
- Don't worry about "breaking" anything
- AI tools are here to help, not to judge
- Have fun exploring what's possible!

**Your AI journey starts now!** 🚀

---

*This manual covers everything a beginner needs to start using fabric-lite and local AI tools effectively and safely.*