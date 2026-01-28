See /home/oak38/projects/AGENTS.md for auto-resume and security protocol.

# CLAUDE.md

This file provides specific context and technical details for Claude, primarily focusing on the `openwebui-stack` project.

For a comprehensive overview of the workspace, all projects, and the general development workflow, please refer to the shared context file:

*   **Shared Context:** `/home/oak38/projects/SHARED_CONTEXT.md`

Please prioritize the information in the shared context file for all general inquiries and use this file for deep technical details about the `openwebui-stack`.

---

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Projects Overview

This directory contains multiple projects:

- **openwebui-stack/open-webui**: Main project - Self-hosted AI platform (Open WebUI)
- **coffeeproject**: Small markdown-based coffee blog planning project
- **tools**: Python utility scripts (e.g., yt_analyzer.py)
- **agents**: Agent prompts and templates
- **mine**: Personal persona documents

## Open WebUI - Main Project

### Build & Development Commands

**Frontend (SvelteKit + Tailwind):**
```bash
cd openwebui-stack/open-webui
npm install                    # Install frontend dependencies
npm run dev                    # Start frontend dev server (default port 5173)
npm run dev:5050               # Start frontend dev server on port 5050
npm run build                  # Production build
npm run check                  # TypeScript/Svelte type checking
npm run lint                   # Run all linting (frontend + types + backend)
npm run lint:frontend          # ESLint frontend only
npm run format                 # Prettier formatting
```

**Backend (FastAPI + Python 3.11+):**
```bash
cd openwebui-stack/open-webui/backend
./dev.sh                       # Start backend with hot reload (port 8080)
./start.sh                     # Production start script

# Or using the CLI:
pip install -e .               # Install in editable mode
open-webui serve               # Start production server
open-webui dev                 # Start development server with reload
```

**Testing:**
```bash
# Frontend tests
npm run test:frontend          # Vitest unit tests

# E2E tests
npm run cy:open                # Open Cypress test runner

# Backend tests
cd backend && pytest           # Run pytest (requires test dependencies)
```

**Docker:**
```bash
cd openwebui-stack/open-webui
make install                   # docker-compose up -d
make start                     # Start containers
make stop                      # Stop containers
make startAndBuild             # Rebuild and start
make update                    # Pull changes and rebuild
```

### Architecture

**Monorepo Structure:**
- `src/` - SvelteKit frontend application
- `backend/` - FastAPI Python backend

**Frontend (src/):**
- `src/lib/apis/` - API client modules matching backend routers (auths, chats, models, etc.)
- `src/lib/components/` - Svelte components organized by feature (admin, chat, workspace, etc.)
- `src/lib/stores/` - Svelte stores for state management
- `src/lib/i18n/` - Internationalization resources
- `src/routes/` - SvelteKit file-based routing

**Backend (backend/open_webui/):**
- `main.py` - FastAPI app entry point, middleware configuration, router mounting
- `routers/` - API endpoints (auths, chats, models, ollama, openai, retrieval, etc.)
- `models/` - SQLAlchemy database models
- `retrieval/` - RAG functionality (loaders, vector stores, web search)
- `utils/` - Shared utilities (auth, chat processing, embeddings, MCP integration)
- `socket/` - WebSocket handling for real-time features
- `migrations/` - Alembic database migrations

**Key Integrations:**
- Ollama and OpenAI-compatible API backends
- ChromaDB for vector storage (RAG)
- Multiple embedding providers
- Web search providers (SearXNG, Google, Brave, etc.)
- Image generation (AUTOMATIC1111, ComfyUI, DALL-E)

**Configuration:**
- Environment variables defined in `backend/open_webui/env.py` and `config.py`
- `.env.example` provides template for local configuration
- Key vars: `OLLAMA_BASE_URL`, `OPENAI_API_KEY`, `CORS_ALLOW_ORIGIN`

### Running Full Development Stack

1. Start backend: `cd backend && ./dev.sh` (port 8080)
2. Start frontend: `npm run dev` (port 5173)
3. Frontend proxies API calls to backend via vite config

Production serves both frontend and backend from the same FastAPI server on port 8080.

---

## Auto-Resume Protocol

If the user says “resume”, “good morning”, “where were we”, or similar:
1. Read `/home/oak38/projects/SHARED_CONTEXT.md`
2. Read the latest session log in the active project (see its `docs/session_log.md`)
3. Summarize last completed work, current blockers, and the next 1–3 actions
4. If multiple projects are active, ask which one to continue

## Repository Security Protocol

Follow the security rules in `/home/oak38/projects/AGENTS.md` before any commit or push.
