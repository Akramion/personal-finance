# AGENTS.md

## Project

This is a personal finance management project written in Go.

The goal is to help a user:
- track accounts
- track debts and credit cards
- track income and expenses
- calculate monthly budget
- calculate debt payoff scenarios
- estimate free money until next salary

This is not a banking system and should not handle real payments.

## Tech stack

- Go
- Standard library first
- No framework unless explicitly requested
- Tests are required for domain and calculator logic

## Architecture rules

Use simple layered architecture:

- `internal/domain` — pure domain types and business rules
- `internal/calculator` — financial calculations
- `internal/usecase` — application scenarios
- `internal/storage` — persistence interfaces and implementations
- `cmd/app` — entry point

Do not put business logic in `main.go`.

Prefer simple code over abstractions.

## Coding style

- Keep functions small and readable
- Avoid premature interfaces
- Avoid global mutable state
- Use explicit names
- Prefer table-driven tests
- All money values should be represented carefully

## Money rules

Do not use float64 for money calculations.

Use one of:
- integer minor units, for example kopecks/cents
- or a dedicated Money type

Currency must be explicit.

## Before changing code

Before implementing a feature:
1. Read related domain types
2. Check existing tests
3. Explain the intended change briefly
4. Make the smallest reasonable change
5. Add or update tests

## Commands

Run tests:
```bash
go test ./...
```

Run app:
```bash
go run ./cmd/app
```

Format code:
```bash
gofmt -w .
```

## Done criteria
A task is done only if:
- code compiles
- tests pass
- important business logic has tests
- no unrelated refactoring was made