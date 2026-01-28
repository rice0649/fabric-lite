See /home/oak38/projects/AGENTS.md for auto-resume and security protocol.

### **Implementation Plan for `opencode`: `forge` User Context Onboarding**

**Objective:**
To implement a "User Context Onboarding" feature within the `forge init` command. This feature will interactively collect a user's persona, goals, and work style, then generate a `SHARED_CONTEXT.md` file. This file will be used by all `fabric-lite` AI agents to create a personalized, optimal workflow that adapts to the user and keeps the project on track.

---

#### **Phase 1: Modify `forge` CLI for Interactive Onboarding**

1.  **Target File:** `internal/cli/init.go` (or wherever the `forge init` command logic resides).

2.  **Logic to Implement:**
    *   After the existing project scaffolding logic, add a call to a new function, e.g., `runUserContextOnboarding()`.
    *   This function should first ask the user if they wish to proceed with the onboarding.
    *   If yes, use a Go library for interactive prompts (like `survey`) to ask a series of questions based on the `my_persona.md` template. The questions should be organized into sections (Mission, Skills, Workflow, etc.).
    *   Store the user's answers in a dedicated struct, e.g., `core.UserContextData`.

3.  **Example Implementation (`init.go`):**
    ```go
    // In the 'forge init' command execution logic...
    
    // ... after scaffolding project files ...
    
    fmt.Println("Project scaffolding complete.")
    if prompt.Confirm("To create an optimal, personalized workflow, would you like to proceed with User Context Onboarding?") {
        contextData, err := onboarding.RunInteractiveQuestions()
        if err != nil {
            return fmt.Errorf("failed during user context onboarding: %w", err)
        }
    
        // Pass contextData to the generator in the next phase
        // ...
    }
    ```

---

#### **Phase 2: Generate Context Artifacts (`SHARED_CONTEXT.md` and `ideas.md`)**

1.  **Target File:** Create a new file: `internal/core/context_generator.go`.

2.  **Logic to Implement:**
    *   Create a function `GenerateSharedContext(data core.UserContextData, projectPath string) error`.
    *   This function will use Go's `text/template` package to render two files in the `projectPath`:
        1.  **`SHARED_CONTEXT.md`**: The template will structure the user's answers from the `UserContextData` struct into a clean Markdown file. It will also append the pre-defined "AI Constructs - Specialized Roles" section, ensuring the user's context and the AI's capabilities are in one place.
        2.  **`ideas.md`**: A simple file containing a header like `# Ideas and Musings` and a brief explanation of its purpose for capturing "squirrel" moments.

3.  **Example Implementation (`context_generator.go`):**
    ```go
    const sharedContextTemplate = `
    # SHARED_CONTEXT.md - Project Stack
    
    ## 1. User Mission & Goals
    **Primary Mission:** {{.Mission}}
    
    ## 2. User Strengths
    **Top Skills:** {{.Skills}}
    
    // ... and so on for all collected data ...
    
    ---
    ## AI Constructs - Specialized Roles
    *   **OpenCode (Master Planner):** ...
    *   **Ollama (Quick Task Automator):** ...
    // ...
    `
    
    func GenerateSharedContext(...) error {
        // 1. Create and parse the template.
        // 2. Open the output file (e.g., filepath.Join(projectPath, "SHARED_CONTEXT.md")).
        // 3. Execute the template with the user's data.
        // 4. Create the ideas.md file.
        return nil
    }
    ```

---

#### **Phase 3: Standardize Agent Behavior (Documentation and Convention)**

This phase ensures the generated context is actually used by the system.

1.  **Target Files:**
    *   `agents/ORCHESTRATOR.md`
    *   The `system.md` file for all key patterns (e.g., `patterns/planning/create_architecture/system.md`, `patterns/implementation/...,` etc.).

2.  **Content to Add:**
    *   A strict convention must be documented in `agents/ORCHESTRATOR.md`:
        > **"Agent Protocol: Stack Loading. The first action of any agent or tool execution must be to load and parse the `SHARED_CONTEXT.md` file from the project root. All subsequent actions, plans, and outputs must be guided by the context defined within."**
    *   The system prompts for all relevant patterns must be updated to include this as their absolute first instruction. The update previously made to `create_architecture` can serve as the template.

---

**Acceptance Criteria:**
1.  Running `forge init` on a new project triggers the interactive questionnaire.
2.  Upon completion, `SHARED_CONTEXT.md` and `ideas.md` are successfully created in the new project's root directory.
3.  The `SHARED_CONTEXT.md` contains the user's answers and the AI tool personas.
4.  Core agent documentation and patterns are updated to enforce the context-loading protocol.