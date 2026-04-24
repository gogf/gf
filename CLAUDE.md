# Karpathy Guidelines

**Tradeoff:** These guidelines bias toward caution over speed. For trivial tasks, use judgment.

## 1. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:
- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

## 2. Simplicity First

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:
- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:
- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

## 4. Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:
- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:
```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

# Development Workflow Rules

This project follows `SDD` and uses the `OpenSpec` tool to drive implementation. Change records are stored under `openspec/changes/`. Each change includes `proposal.md`, `design.md`, `specs/`, and `tasks.md`.

**Execution workflow**:
1. Use the `/opsx:explore` slash command at `.agents/prompts/opsx/explore.md` to conduct an exploratory discussion based on the requirement description, analyze the problem, design the solution, and assess risks.
2. Once the exploratory discussion finishes and the solution is clear, use the `/opsx:propose` slash command at `.agents/prompts/opsx/propose.md` to turn it into a formal `OpenSpec` change proposal. The command format is `/opsx:propose feature-name`, where `feature-name` is a descriptive name for the current change in `kebab-case`, such as `user-auth` or `data-export`. A new change directory will then be generated automatically under `openspec/changes`, containing incremental spec documents (`spec/`), the technical implementation plan (`design.md`), the proposal and rationale (`proposal.md`), and the implementation task list (`tasks.md`).
3. Then run the `/opsx:apply` slash command at `.agents/prompts/opsx/apply.md` to execute the items in `tasks.md` one by one, completing code changes, tests, and documentation updates. After the work is done, the `/gf-review` skill must be invoked for code and spec review.
4. When users report issues or improvement requests, the `/gf-feedback` skill must be used to fix and verify them, and the related `OpenSpec` documents must be updated. After the work is done, the `/gf-review` skill must be invoked for review.
5. After the user confirms that the current iteration is complete and has no remaining issues, run the `/opsx:archive` slash command at `.agents/prompts/opsx/archive.md` to archive the change. Before archiving, the `/gf-review` skill must be used for a full change review to ensure code quality and compliance with the spec.

**Key rules**:
- **An `OpenSpec` change is considered active until it is archived**: any change directory that still exists directly under `openspec/changes/` and has **not been moved to** `openspec/changes/archive/` is an active change. **Even if the change has completed all tasks and `openspec list --json` shows `status: complete`, it must still be treated as active until the archive step has been executed.**
- When a user reports a bug, defect, or improvement request in either Chinese or English, and there is an active `OpenSpec` change in the project, the `gf-feedback` skill must be used. **Unless the user explicitly asks to create a new change, the feedback must always be appended to the current active iteration, even if it is unrelated to the main feature of that iteration**, so that everything can be managed and archived together.
- The `/gf-review` review skill is triggered automatically after `/opsx:apply` completes, after `/opsx:feedback` completes, and before `/opsx:archive`.
- During development tasks executed with tools such as `Claude Code` or `Codex CLI`, if the work can be parallelized effectively with `SubAgent` and doing so would clearly improve efficiency, that option must be evaluated first and adopted whenever appropriate. Only skip `SubAgent` when the task is strongly dependent on serial context, the split cost is too high, or it introduces obvious collaboration risk.
- When creating new iteration documents, the content of `proposal.md`, `design.md`, `tasks.md`, and the incremental specs must be written in English.

# Code Development Rules
- All source code must include comments, such as package comments, file comments, method comments for both public and private methods, constant comments, variable comments, and comments for key logic.
- **All submitted code changes must include unit tests**: every submitted code change must add or update focused unit tests that directly cover the affected logic and expected behavior of the changed code path, and the coverage for newly added code must stay at or above 80%; 90% or above is the preferred target when feasible.
- **Do not hardcode string literals with enum semantics in backend implementation code**: values that represent statuses, types, stages, actions, execution modes, sort directions, filter operators, or any other enum-like semantics must be managed through Go named types and constants. Writing raw string literals directly in business branching, comparisons, assignments, or persistence logic is forbidden.
- **Do not ignore any `error` return value**: every call that may return an `error` must be handled explicitly. Do not use patterns such as `_ = someFunc()`, `_, _ = someFunc()`, or direct calls that discard returned errors. In business flows, errors should be returned explicitly or converted before returning; in initialization, startup, or other critical non-recoverable paths, they should `panic`; in tests and cleanup paths, they must still be asserted, logged, or otherwise handled explicitly rather than silently ignored.
- **Do not use stand-alone assignments like `_ = var` to mask unused parameters or local variables**: this placeholder pattern has no business meaning and creates misleading signals about whether the variable was supposed to participate in the logic. Prefer deleting unused variables. If a parameter must be kept to satisfy an interface signature or callback contract, use the blank identifier directly in the function signature, such as `func(ctx context.Context, _ gdb.TX) error`, or omit an unused receiver name instead of adding one-line statements like `_ = tx`, `_ = req`, or `_ = ctx` in the function body.
- **File header comment rules**:
  - Every `Go` source file must include a file-purpose comment at the top of the file. Component-level comments should appear in the component's main file, meaning the file with the same name as the component, such as `plugin.go`, `config.go`, or `file.go`.
  - In a main file, the component comment must be placed immediately before the `package xxx` declaration with no blank line in between. For example:
  ```go
  // Package plugin implements plugin manifest discovery, lifecycle orchestration,
  // governance metadata synchronization, and host integration for LinaPro plugins.
  package plugin
  ```
  - Other implementation files must keep only file comments that describe the responsibility of the current file, such as `plugin_xxx.go` or `config_xxx.go`. There must be one blank line between the file comment and `package xxx`, and non-main files must not duplicate component-level descriptions.
- **Variable Declarations**: When defining multiple variables, use a `var` block to group them for better alignment and readability:
  ```go
  // Good - aligned and clean
  var (
      authSvc       *auth.Service
      bizCtxSvc     *bizctx.Service
      k8sSvc        *svcK8s.Service
      notebookSvc   *notebook.Service
      middlewareSvc *middleware.Service
  )

  // Avoid - scattered declarations
  authSvc := auth.New()
  bizCtxSvc := bizctx.New()
  k8sSvc := svcK8s.New()
  ```
  Apply this pattern when you have 3 or more related variable declarations in the same scope.
