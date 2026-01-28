See /home/oak38/projects/AGENTS.md for auto-resume and security protocol.

# GEMINI.md

This file provides context for the Gemini CLI agent to interact with the `fabric-lite` repository.

## Project Overview

`fabric-lite` is a Go-based command-line tool designed as a lightweight, extensible AI augmentation framework. It allows users to execute "patterns" (reusable AI prompt templates) against text input using various AI providers like OpenAI, Anthropic, or a local Ollama instance.

The project also includes a "forge" workflow, which is a structured development process divided into phases (e.g., discovery, planning, implementation), each with a recommended AI tool. An `agents` directory suggests a multi-agent orchestration system for more complex, automated workflows.

## Key Concepts

*   **Patterns:** Reusable prompt templates for common tasks like summarizing text or extracting ideas. They are stored in the `patterns/` directory and can be extended by users.
*   **Providers:** The AI backends that execute the prompts. The tool is configured to support `openai`, `anthropic`, and `ollama`.
*   **Tools:** In addition to patterns, the CLI can directly invoke specialized tools like `codex`, `gemini`, and `claude`.
*   **Phases:** A structured development workflow (`discovery`, `planning`, `design`, etc.) that guides a project from idea to deployment, with specific AI tools recommended for each phase.
*   **Orchestration:** A multi-agent system where different AI agents collaborate on tasks, communicating via shared files.

## Building and Running

The project uses a `Makefile` for common development tasks.

*   **Build the CLI:**
    ```bash
    make build
    ```
    This creates the `fabric-lite` binary in the `./bin/` directory.

*   **Run tests:**
    ```bash
    make test
    ```

*   **Run the CLI:**
    ```bash
    # Using a pattern
    echo "This is some text to summarize." | ./bin/fabric-lite run --pattern summarize --provider ollama

    # Listing available patterns
    ./bin/fabric-lite list
    ```

## Development Conventions

*   **Language:** The project is written in Go.
*   **CLI:** The command-line interface is built using the `cobra` library.
*   **Configuration:** Configuration is managed with `viper`, supporting a `config.yaml` file and environment variables.
*   **Dependencies:** Go modules are used for dependency management (`go.mod`, `go.sum`).
*   **Formatting:** Standard `go fmt` is used for code formatting (`make fmt`).
*   **Linting:** `golangci-lint` is used for linting (`make lint`).
*   **Installation:** The `forgiving-setup.sh` script provides an easy way for users to set up the tool and its dependencies, particularly Ollama for local AI.

## Project Structure

-   `cmd/`: Main application entry points for `fabric-lite` and `forge`.
-   `internal/`: Core application logic, separated by concern (cli, core, executor, providers, tools).
-   `patterns/`: Default AI prompt patterns.
-   `config/`: Example configuration files.
-   `docs/`: Project documentation.
-   `agents/`: Files related to the multi-agent orchestration system.
-   `scripts/`: Helper scripts for installation and other tasks.
-   `Makefile`: Defines build, test, and other development commands.

---

## Auto-Resume Protocol
If the user says “resume”, “good morning”, “where were we”, or similar:
1. Read `/home/oak38/projects/SHARED_CONTEXT.md`
2. Read the latest session log in this project (if present)
3. Summarize last completed work, current blockers, and the next 1–3 actions

## Repository Security Protocol
Follow the security rules in `/home/oak38/projects/AGENTS.md`.
