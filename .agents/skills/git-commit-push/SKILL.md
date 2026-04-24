---
name: git-commit-push
description: Review the current git working tree, generate a commit message from the actual diff using the repository's commit or PR-title convention, commit all current changes on the active branch, and push that branch to `origin`. Use this whenever the user asks to "commit", "push", "commit and push", "generate a commit message", "commit the current changes", or wants the current branch changes sent upstream without hand-writing the git commands.
---

# Git Commit Push

Inspect the current repository changes, derive a concise commit subject that matches the repository convention, commit every current modification on the active branch, and push that branch to `origin`.

This skill is for execution, not just advice. When it triggers, actually run the git workflow unless the repository state makes that unsafe or impossible.

## When To Use

- The user asks you to commit the current changes, with or without asking for push
- The user wants you to write the commit message from the diff instead of inventing one up front
- The user mentions the repo's PR or commit naming convention and wants you to follow it
- The user says things like "commit the current branch", "help me commit", "commit and push", "generate a commit message and push", or "send these changes to origin"

## Core Behavior

1. Confirm you are inside a Git repository and detect the active branch with `git branch --show-current`.
2. Inspect the working tree before committing:
   - `git status --short --branch`
   - `git diff --stat`
   - `git diff --cached --stat`
   - `git diff -- . ':(exclude)package-lock.json'` or narrower path filters only when needed for readability
3. If the repository contains `.github/PULL_REQUEST_TEMPLATE.MD`, read it and treat its PR-title rules as the default commit-subject convention.
4. Generate a commit subject from the actual changed files and diff content, not from the user prompt alone.
5. Stage every current modification on the branch with `git add -A`.
6. Commit once with the generated message.
7. Push the current branch to `origin` with `git push origin <current-branch>`.

## Commit Message Rules

Prefer a single-line subject unless the user explicitly asks for a body.

When `.github/PULL_REQUEST_TEMPLATE.MD` exists, follow its title rules:

- Use `<type>[optional scope]: <description>`
- Choose `type` from the allowed list in the template
- Pick a scope from the dominant changed package, module, directory, or feature area when it is clear
- Use a lowercase verb after the colon
- Do not end the subject with a period
- Keep it short, ideally 76 characters or fewer

Map the diff to the narrowest honest type:

- `feat` for a new user-facing capability
- `fix` for a bug fix
- `docs` for documentation-only changes
- `test` for test-only changes
- `refactor` for structural cleanup without behavior change
- `chore` for tooling, repo maintenance, or housekeeping changes
- `build`, `ci`, `style`, or `perf` only when those clearly fit better

Do not invent issue numbers. If the template mentions `Fixes #1234` or `Updates #1234`, treat that as PR-comment guidance unless the user explicitly asks for a multi-line commit message that includes it.

## Execution Rules

- Commit all current tracked and untracked changes in the working tree, because this skill is for "commit the current state" requests
- If there are no changes, say so clearly and stop before commit or push
- If `git branch --show-current` is empty, explain that `HEAD` is detached and stop unless the user explicitly asks you to commit from detached `HEAD`
- Never use `--force`, `--force-with-lease`, or history-rewriting commands unless the user explicitly asks
- If push fails because the remote branch moved, report the exact failure and stop instead of auto-rebasing or auto-merging
- Do not silently drop files from the commit unless the user asked to exclude them

## Suggested Command Flow

```bash
git status --short --branch
git diff --stat
git diff --cached --stat
test -f .github/PULL_REQUEST_TEMPLATE.MD && sed -n '1,220p' .github/PULL_REQUEST_TEMPLATE.MD
branch_name=$(git branch --show-current)
git add -A
git commit -m "<generated-subject>"
git push origin "$branch_name"
```

Inspect `git diff --cached` again after staging if the pre-stage diff was noisy or if untracked files materially change the scope.

## Output Contract

When you use this skill:

- Tell the user which branch you committed
- Provide the final commit subject you used
- Mention that you staged all current changes
- Report the push target as `origin/<branch>`
- If commit or push did not happen, explain exactly why

## Example

User request:

```text
Generate a commit message that follows this repository's convention, then commit and push the current branch
```

Expected behavior:

- Inspect the repo status and diff
- Read `.github/PULL_REQUEST_TEMPLATE.MD` if present
- Generate a conventional subject from the real changes
- Run one commit for the whole current working tree
- Push the active branch to `origin`
