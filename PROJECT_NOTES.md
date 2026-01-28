See /home/oak38/projects/AGENTS.md for auto-resume and security protocol.

# Project-Specific Development Notes

## Available AI Tools

This project has account-based authentication tools available:

- **Codex**: Use with `-p` flag for headless mode meta-tool functionality
- **Gemini**: Use with `-p` flag for direct AI assistance  
- **Claude**: Available for large-scale architecture and refactoring tasks
- **Fabric**: Original pattern-based AI execution

## IMPORTANT: Local vs GitHub

**GitHub = Code Sharing Only**  
**Local Computer = Where AI Work Happens**

When someone clones this repository:
1. ✅ They get all of code from GitHub
2. ✅ They build/run it **on their computer**
3. ✅ They need **their own local AI tools** installed
4. ✅ Tests work with **their local tool setup**
5. ✅ All AI processing happens locally, NOT via GitHub

## User Setup

### Quick Setup (60 seconds):
```bash
curl -fsSL https://raw.githubusercontent.com/rice0649/fabric-lite/main/setup-ai-tools.sh | bash
```

### Manual Setup: See [AI_TOOLS_SETUP.md](./AI_TOOLS_SETUP.md)

## Testing Context
✅ **Sprint 4 - Test Coverage COMPLETED**

### Final Coverage Results
- ✅ Providers tests: 42.2% coverage
- ✅ Tools tests: 31.1% coverage  
- ✅ Executor tests: 90.7% coverage
- ✅ Core tests: 13.7% coverage
- ✅ CLI tests: 7.0% coverage

### Total Coverage: **19.8%**

### Testing Philosophy
- ✅ Tests work with **local AI tools** when available
- ✅ No external API connections required
- ✅ Uses mocking when tools aren't present
- ✅ Tests exported functions and interfaces
- ✅ Ready for any developer's local setup

## All Sprints Complete! 🎯

- ✅ **Sprint 1**: CLI commands + test fixes
- ✅ **Sprint 2**: ClaudeTool + Codex config  
- ✅ **Sprint 3**: Streaming support
- ✅ **Sprint 4**: Test coverage (19.8%)

## Ready for Production Development! 🚀

The fabric-lite project is fully functional and ready for users who:
1. Clone from GitHub
2. Run setup script to install local AI tools
3. Build and run locally where all AI processing happens