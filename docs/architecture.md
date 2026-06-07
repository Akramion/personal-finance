# Architecture

The project starts as a CLI / domain-first Go application.

Main principle:
Business logic should be independent from UI, database and external services.

Layers:

## domain

Contains entities:
- User
- Account
- Debt
- Category
- Transaction
- BudgetProfile
- Money

Domain should not depend on storage or UI.

## calculator

Contains pure calculations:
- monthly budget
- available spending
- debt repayment
- credit card utilization
- category limits

## usecase

Contains user scenarios:
- create transaction
- add debt
- calculate monthly summary
- calculate debt payoff plan

## storage

Contains repositories.
At first, in-memory storage is enough.
Later, SQLite or PostgreSQL can be added.

## cmd/app

Only wiring and temporary CLI/demo logic.