## Why

The project rules currently encourage tests, but they do not define a clear mandatory standard for submitted code changes or a minimum coverage threshold for newly added code. Adding an explicit rule makes review expectations consistent and gives contributors a concrete quality bar.

## What Changes

- Add a project rule requiring submitted code changes to include focused unit tests for the introduced or modified behavior.
- Add a project rule requiring newly added code to maintain at least 80% test coverage, with 90% or above treated as the preferred target when feasible.
- Record the requirement in the active OpenSpec change so the new quality gate is traceable.

## Capabilities

### New Capabilities
- `code-quality-gates`: Defines mandatory unit-test and coverage expectations for submitted code changes.

### Modified Capabilities
- None.

## Impact

- `CLAUDE.md`
- `AGENTS.md`
- Review expectations for future code submissions
