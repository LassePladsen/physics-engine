# AGENTS.md

## Project overview

This is a work-in-progress physics engine written in Go. Prefer small, focused
changes that improve correctness, clarity, or test coverage without assuming the
current structure is final.

## Working conventions

- Read the relevant code and existing tests before changing behavior.
- Preserve public APIs unless the task explicitly calls for an API change.
- Keep packages and responsibilities clear; avoid introducing unnecessary
  abstractions or dependencies.
- Use idiomatic Go and run `gofmt` on modified Go files.
- Add or update tests when changing observable behavior or fixing a bug.
- Run the narrowest relevant tests first, then `go test ./...` when practical.

## Physics and numerical code

- Favor correctness and predictable behavior over premature optimization.
- Be explicit about units, coordinate conventions, assumptions, and tolerance
  values when they matter.
- Treat floating-point comparisons carefully: use a documented epsilon where an
  exact comparison is not appropriate.
- Guard against invalid inputs and unstable states (for example zero mass,
  zero-length vectors, division by zero, or non-finite values) in a way that
  matches the surrounding API style.
- Keep simulation behavior deterministic when feasible, especially in tests.

## Collaboration

- Do not overwrite or revert unrelated working-tree changes.
- Keep commits and patches narrowly scoped and explain behavioral tradeoffs in
  code comments or the change description when they are non-obvious.
- Update documentation when a user-facing workflow or API changes.
