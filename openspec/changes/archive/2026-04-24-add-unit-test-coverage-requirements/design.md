## Context

`CLAUDE.md` and `AGENTS.md` currently act as the canonical project rule documents used during implementation and review. The requested change is documentation-only, but the wording must clearly separate the hard minimum quality bar from the preferred target so contributors and reviewers apply the same standard.

## Goals / Non-Goals

**Goals:**
- Add an explicit rule that submitted code changes must include unit tests.
- Add an explicit rule that newly added code must reach at least 80% coverage.
- Preserve a stronger preferred target of 90% or above when feasible.
- Keep the change limited to the canonical rule documents and OpenSpec artifacts.

**Non-Goals:**
- Adding CI automation or coverage tooling in this change.
- Defining repository-wide historical coverage requirements.
- Changing unrelated contribution or coding rules.

## Decisions

- Update both `CLAUDE.md` and `AGENTS.md` so the duplicated project rule documents remain aligned.
- Place the new requirement under `Code Development Rules`, where existing mandatory engineering requirements already live.
- Use 80% as the hard minimum for newly added code and describe 90% or above as the preferred target when feasible.
- Keep the wording focused on unit tests and coverage expectations rather than implementation details about how coverage is measured.

## Risks / Trade-offs

- [Documentation-only enforcement] -> The rule depends on reviewer and contributor discipline until tooling is added; placing it in the canonical rule documents keeps the expectation visible during development and review.
- [Coverage interpretation] -> Different packages may measure coverage differently; using a clear minimum plus a preferred target reduces ambiguity without over-specifying tooling.
