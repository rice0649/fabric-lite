# Repository Guidelines

## Project Structure & Module Organization

- `cmd/fabric-lite/` holds the CLI entry point and wiring.
- `internal/` contains core logic, CLI handling, and provider integrations.
- `patterns/` stores built-in prompt patterns (e.g., `summarize/`, `extract_ideas/`).
- `config/` provides configuration templates; `docs/` is project documentation.
- `agents/` contains orchestration assets; `scripts/` has build/install helpers.

## Build, Test, and Development Commands

- `make build` compiles the CLI to `bin/fabric-lite`.
- `make install` copies the binary into your `GOPATH/bin` (or `~/go/bin`).
- `make run ARGS='--help'` runs the CLI directly via `go run`.
- `make test` runs `go test -v ./...` (no tests currently detected).
- `make test-coverage` produces `coverage.out` via `go test -coverprofile=...`.
- `make lint` runs `golangci-lint` if installed.
- `make fmt` runs `go fmt ./...` to format source files.

## Coding Style & Naming Conventions

- Go code should be gofmt-formatted (`make fmt`) and keep standard Go layout.
- Prefer Go idioms: short, clear names; exported identifiers in `CamelCase`.
- Pattern directories use `snake_case` (see `patterns/explain_code`).
- Keep CLI flags and config keys consistent with `config/` templates.

## Testing Guidelines

- Use Go’s standard `testing` package; place tests beside code as `*_test.go`.
- Name tests as `TestXxx` and subtests with `t.Run("case")`.
- Add coverage for new logic where practical; `make test-coverage` is available.

## Commit & Pull Request Guidelines

- Existing history uses short, direct subjects (e.g., "Add multi-agent orchestration system").
- Keep commit subjects imperative and under ~70 characters when possible.
- PRs should describe intent, list key changes, and include test commands run.
- Include examples for user-facing changes (CLI output, new patterns, config).

## Security & Configuration Tips

- Store API keys in environment variables referenced by `~/.config/fabric-lite/config.yaml`.
- Avoid committing secrets or personal config files; use templates in `config/`.
