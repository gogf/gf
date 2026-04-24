---
name: openspec-review
description: >-
  Code and specification review for OpenSpec workflow. Triggers automatically after /opsx:apply
  task completion, after /opsx:feedback task completion, and before /opsx:archive. Use when
  user requests code review, spec compliance check, or when explicitly invoked via /openspec-review.
compatibility: Requires OpenSpec CLI, GoFrame v2 skill, openspec-e2e skill.
---

# OpenSpec Review

Structured code and specification review for the OpenSpec development workflow.

**Spec Source**: `CLAUDE.md` is the single source of truth for all review criteria.

---

## When This Skill Activates

**Automatic triggers:**
- After completing each task in `/opsx:apply`
- After completing each task in `/opspec-feedback`
- Before executing `/opsx:archive`

**Manual trigger:**
- User explicitly requests: "review this code", "check spec compliance", "/openspec-review"

---

## Review Workflow

### 1. Identify Scope

Determine what needs to be reviewed:

1. **After task completion** — Review files modified/created by the completed task
2. **Before archive** — Review all changes in the current OpenSpec change
3. **Manual invocation** — Ask user to specify scope or use current change

**Mandatory scope collection rules:**

1. Start with repository status, not `git diff` alone:
   ```bash
   git status --short
   git ls-files --others --exclude-standard
   ```
2. Treat **all tracked and untracked changes** as review candidates, including:
   - staged files
   - unstaged files
   - untracked files shown as `??`
   - untracked directories shown as `?? path/`
3. When `git status --short` reports an untracked directory, expand it to concrete files before review:
   ```bash
   find <path> -type f
   # Or prefer:
   rg --files <path>
   ```
4. If the task ran generators such as `make ctrl`, `make dao`, codegen scripts, or produced new test files, explicitly include the generated untracked files in review scope even if they do not appear in `git diff`.
5. `git diff` may be used only as a secondary narrowing aid after status collection. It is **never sufficient by itself** for review scope definition.

Run `openspec status --change "<name>" --json` to understand the current change state.

### 2. Load Specifications

Read `CLAUDE.md` to load all specifications. This is the single source of truth.

### 3. Backend Code Review

**Trigger**: Changes to files under `apps/lina-core` directory

1. Invoke `goframe-v2` skill for GoFrame framework conventions
2. Check against `CLAUDE.md` backend code specifications

### 4. RESTful API Review

**Trigger**: Any API endpoint changes

Check against `CLAUDE.md` API design specifications.

### 5. Project Specification Review

**Trigger**: Any implementation changes

Check against `CLAUDE.md` architecture design specifications and code development specifications.

### 6. SQL Review

**Trigger**: New or modified files under `apps/lina-core/manifest/sql/`、`apps/lina-core/manifest/sql/mock-data/`、`apps/lina-plugins/**/manifest/sql/` or SQL snippets embedded in related delivery docs

Check against `CLAUDE.md` SQL file management specifications, at minimum covering:
1. File naming, versioning, and single-iteration single-file rules
2. Seed DML vs mock data separation
3. **Idempotent execution safety** — SQL must be safe to run multiple times without duplicate-object errors or duplicate seed data; verify use of `IF [NOT] EXISTS`, `IF EXISTS`, `INSERT IGNORE`, or equivalent safe re-entry patterns
4. **Seed write style compliance** — delivered SQL must reject `INSERT ... ON DUPLICATE KEY UPDATE` and reject explicit writes to `AUTO_INCREMENT` `id` columns in seed/mock/install data
5. Whether schema/data changes still match the current change scope and deployment path

### 7. E2E Test Review

**Trigger**: New or modified E2E test files in `hack/tests/e2e/` directory

1. Invoke `openspec-e2e` skill for test conventions
2. Check against `CLAUDE.md` E2E test specifications

### 8. Generate Review Report

```markdown
## OpenSpec Review Report

**Change:** <change-name>
**Scope:** <task-specific / full change>
**Files Reviewed:** <count>
**Scope Source:** `git status --short` + `git ls-files --others --exclude-standard` + task/change context

### Backend Code Review
✓ All checks passed / ⚠ N issues found

### RESTful API Review
✓ All endpoints compliant / ⚠ N violations found

### Project Spec Review
✓ Compliant with CLAUDE.md / ⚠ N violations found

### SQL Review
✓ No SQL changes / ✓ SQL changes compliant / ⚠ N SQL issues found

### E2E Test Review
✓ Tests follow conventions / ⚠ N issues found

### Summary
- **Critical:** N (must fix before archive)
- **Warnings:** N (recommended to fix)

### Recommended Actions
1. [Specific action with CLAUDE.md reference]
```

---

## Issue Severity

| Level | Behavior |
|-------|----------|
| **Critical** | Block archive, must fix |
| **Warning** | Show but allow proceed |

---

## Integration Points

| Workflow Step | Behavior |
|---------------|----------|
| `/opsx:apply` task done | Review, offer to fix issues before next task |
| `/opspec-feedback` task done | Review, fix before marking complete |
| `/opsx:archive` | Review all changes, block on critical issues |

---

## Guardrails

- **CLAUDE.md is the single source of truth** — All spec references point to it
- Only check categories relevant to changed files
- Scope identification MUST include untracked files and expanded untracked directories; never rely on `git diff` alone
- Don't block on warnings — only critical issues block archive
- Include file paths and line numbers in issue reports
- Offer to fix issues automatically when straightforward
