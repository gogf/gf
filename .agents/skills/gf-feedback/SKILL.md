---
name: gf-feedback
description: >-
  Track, fix, verify, and test any bugs, improvements, or gaps reported against an OpenSpec change.
  MUST use this skill whenever user reports problems, defects, issues, bugs, or gaps related to
  existing implementations, even if they don't explicitly say "feedback" or mention OpenSpec.
compatibility: Requires openspec CLI, Go toolchain, and gf-review skill.
---

# Feedback: Structured Fix, Verification & Test Coverage Loop

When users discover bugs or improvement points after implementation, this skill captures those issues, organizes them into a traceable task list in `tasks.md`, systematically fixes and verifies each one, and ensures every behavior-changing fix is covered by focused unit tests.

**Core principles:**
1. **Spec is the source of truth** — Spec-level changes require spec update before task recording
2. **Write it down first, then fix it** — Every issue gets recorded before any code change
3. **Every fix deserves a test** — Behavior-changing code changes require unit test coverage in the affected package

---

## Workflow

### 1. Identify Target Change

**CRITICAL:** 
1. Always append to existing active changes. Only create new change when none exist.
2. An **active change** is any change directory that still exists directly under `openspec/changes/` and has **not** been moved into `openspec/changes/archive/`. Do **not** treat `status: complete`, all tasks checked off, or similar completion signals as "inactive" until archive actually happens.
3. Regardless of whether the feedback content is related to the main functionality of the current active iteration, it MUST be appended to the current active iteration. This ensures all changes are tracked in a single change record for unified management and archiving.

```bash
openspec list --json
# Or: ls openspec/changes/ | grep -v archive
```

When the two signals disagree, prefer the filesystem rule:

- If a change directory still exists under `openspec/changes/` and is not inside `archive/`, it is active.
- `openspec list --json` may still report such a change as `status: complete`; that only means implementation tasks are done, **not** that the change is inactive.
- Only archived changes under `openspec/changes/archive/` are inactive.

| Active Changes | Action |
|----------------|--------|
| None | Create new change (see below) |
| One | Auto-select it, announce and proceed |
| Multiple | Ask user to select from list |

**When multiple active changes exist:**
```
Multiple active changes detected. Which change should this feedback be appended to?

1. config-management — System config CRUD management
2. user-auth — User authentication enhancement

Please select 1 or 2:
```

**When no active change exists:**
1. Derive kebab-case name from feedback (e.g., "fix-menu-circular-ref")
2. If name exists, append suffix ("-2")
3. Create: `openspec new change "<name>"`
4. Generate minimal `proposal.md` (one paragraph summarizing context)
5. Skip `design.md` for pure bug fixes unless architectural changes needed

Announce: "Applying feedback fixes to change: **<name>**"

---

### 2. Read Current Context

| File | Purpose |
|------|---------|
| `tasks.md` | Task structure, naming conventions, numbering |
| `design.md` | Architectural context |
| `proposal.md` | Feature scope and intent |
| `specs/` | Delta spec definitions |

```bash
# Locate existing unit tests in the impacted package before adding a new one
rg --files <pkg-dir> | rg '_test\\.go$'
```

---

### 3. Analyze and Organize Issues

For each reported issue:

**Classify by type:**
- **bug** — Incorrect behavior, code doesn't match spec
- **missing** — Feature incomplete, gaps in implementation
- **ux** — UX improvement, no spec change needed
- **test-gap** — Missing test coverage only

**Classify by spec impact:**

| Level | Definition | Action |
|-------|------------|--------|
| **implementation** | Spec is correct, code is wrong | Fix code only |
| **spec-level** | Requirement missing/incomplete/changed | Update spec first, then fix |
| **internal** | No user-observable change | Fix code, test optional |

**Group related issues** — Same root cause → single task with multiple verification points.

---

### 4. Update Delta Specs (for Spec-Level Issues Only)

For spec-level issues, update specs **before** recording tasks:

1. Identify affected capability: `specs/<capability>/spec.md`
2. Apply delta operation:

```markdown
<!-- ADDED: New requirement -->
### Requirement: Parent Selector Circular Prevention
The system SHALL disable the current menu and all its descendants in the parent selector
to prevent circular references.

#### Scenario: Edit menu with children
WHEN user edits a menu that has child menus
THEN the parent selector SHALL disable the current menu and all descendant menus

<!-- MODIFIED: Changed requirement (include full original block) -->
### Requirement: Import Error Handling
The system SHALL display error messages when import fails.
**MODIFIED:** Error messages SHALL include row number, field name, and validation failure reason.

<!-- REMOVED: Deprecated requirement -->
### Requirement: Legacy Import Format
The system SHALL support legacy CSV format.
**REMOVED:** This format is no longer supported.
**Migration:** Use the new CSV format with header row.
```

---

### 5. Write Task List to tasks.md

Append a **Feedback section** to `tasks.md`:

```markdown
## Feedback

- [ ] **FB-1**: Parent selector allows circular references in menu edit
- [ ] **FB-2**: Import error messages lack row and field details
- [ ] **FB-3**: No test coverage for reset password feature
```

**Numbering:** Sequential `FB-1`, `FB-2`, etc. Continue from last number if section exists.

**One line per task** — No sub-fields. Analysis happens during fix phase.

**Confirm with user** before writing to file.

**Test coverage planning (internal):**
- Behavior-changing code change → Unit test required
- Internal-only optimization → Unit test optional unless logic risk increased
- Prefer extending the nearest `*_z_unit*_test.go` or `*_test.go` in the same package

---

### 6. Execute Fixes (Loop)

For each task:

**a. Announce:** `## Fixing FB-X: <issue title>`

**b. Investigate** — Read source files, confirm root cause

**c. Implement** — Minimal, focused fix following existing patterns

**d. Write/update unit tests** — Prefer the affected package's existing `*_z_unit*_test.go` or `*_test.go` files and keep assertions focused on the changed logic

**e. Assess Impact Scope (MANDATORY)**

After implementing, identify regression risk:

| Change Type | Map To Tests |
|-------------|--------------|
| Package-level logic | Targeted test for changed function/method + package regression tests |
| Shared utility | Utility package unit tests + highest-value dependent package tests already covering reuse |
| DB/DAO logic | DAO/model package unit tests with focused fixtures, mocks, or test helpers |
| Public API validation | Handler/service package unit tests that assert the changed validation path |
| Refactor without behavior change | Existing package tests that prove behavior parity |

```bash
# Example: Find unit tests related to a changed symbol or package
rg -l "GenDao|gdao" . -g '*_test.go'
```

Announce:
```
### Impact Analysis for FB-X
- Modified: cmd/gf/internal/cmd/gendao/gendao.go
- Affected package: cmd/gf/internal/cmd/gendao
- Unit tests: cmd/gf/internal/cmd/gendao/gendao_test.go
- Regression command: go test ./cmd/gf/internal/cmd/gendao -run 'TestGenDao'
```

**f. Verify (MANDATORY before marking complete)**

1. Run new/updated unit tests for this task → **must pass**
2. Run ALL identified package-level regression tests → **must pass**
3. Only then: mark task `[x]` in tasks.md

If regression fails:
- Fix inline if related to current change
- Add as new FB task if separate issue

**g. Run review** — Invoke `gf-review` skill after completion

---

### 7. Comprehensive Verification

After all fixes:

1. Aggregate regression tests from all tasks
2. Run full set in single pass
3. Report:
```
### Comprehensive Verification Results
- Total tests: N
- Passed: N
- Failed: N (list with details)
- Regression tests: all passed ✓ / X failures
```

If failures → add new FB tasks, loop to Step 6.

---

### 8. Report Completion

```markdown
## Feedback Complete

**Change:** <name>
**Issues reported:** X
**Issues fixed:** Y/X
**Tests added:** Z unit tests / focused assertions
**Regression tests run:** R tests across N packages
**Verification:** all passed / N issues remaining

### Fixed This Session
- [x] FB-1: <title> ✓ (unit test: TestGenDao_FiltersInvalidTables | package: ./cmd/gf/internal/cmd/gendao ✓)
- [x] FB-2: <title> ✓ (unit test: existing package coverage extended | package: ./cmd/gf/internal/cmd ✓)

### Remaining (if any)
- [ ] FB-3: <title> — blocked by <reason>
```

---

## Edge Cases

| Situation | Handling |
|-----------|----------|
| Single issue | Still follow full workflow |
| Missing test cases only | Classify as test-gap, implement tests |
| Fix reveals more issues | Add as new FB tasks |
| "Bug" is actually feature request | Re-classify as spec-level, update specs first |
| Unit test not feasible (docs/spec only) | Note reason explicitly and skip only when no runtime code changes exist |
| Multiple feedback rounds | All tasks in single Feedback section, sequential numbering |

---

## Guardrails

- **Append to active change if exists** — Never create new change when active ones exist
- **Specs before tasks for spec-level issues** — Update delta specs first
- **Write tasks before fixing** — Never code without recording
- **Confirm task list with user** — User validates analysis
- **Minimal fixes** — No refactoring beyond issue scope
- **Behavior-changing fix needs unit test** — No exceptions unless the change is docs/spec only
- **No green check without green unit tests** — Mark `[x]` only after tests pass
- **Impact analysis mandatory** — Every fix requires package-level regression test identification
- **Regression failures block completion** — Must resolve before marking done
- **Update tasks.md in real time** — Mark complete immediately after verification
- **Match file language** — Use same language as existing content in target file
