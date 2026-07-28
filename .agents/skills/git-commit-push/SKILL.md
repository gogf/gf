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

The commit message is formatted as follows: `<type>[optional scope]: <description>` For example, `fix(os/gtime): fix time zone issue`
  + `<type>` is mandatory and can be one of `fix`, `feat`, `build`, `ci`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`
    + `fix`: Used when a bug has been fixed.
    + `feat`: Used when a new feature has been added.
    + `build`: Used for modifications to the project build system, such as changes to dependencies, external interfaces, or upgrading Node version.
    + `ci`: Used for modifications to continuous integration processes, such as changes to Travis, Jenkins workflow configurations.
    + `docs`: Used for modifications to documentation, such as changes to README files, API documentation, etc.
    + `style`: Used for changes to code style, such as adjustments to indentation, spaces, blank lines, etc.
    + `refactor`: Used for code refactoring, such as changes to code structure, variable names, function names, without altering functionality.
    + `perf`: Used for performance optimization, such as improving code performance, reducing memory usage, etc.
    + `test`: Used for modifications to test cases, such as adding, deleting, or modifying test cases for code.
    + `chore`: Used for modifications to non-business-related code, such as changes to build processes or tool configurations.
  + After `<type>`, specify the affected package name or scope in parentheses, for example, `(os/gtime)`.
  + The part after the colon uses the verb tense + phrase that completes the blank in
  + Lowercase verb after the colon
  + No trailing period
  + Keep the title as short as possible. ideally under 76 characters or shorter
+ If there is a corresponding issue, add either `fixes #1234` (the latter if this is not a complete fix) to this comment

### Examples
#### Commit message with description and breaking change footer
```
feat: allow provided config object to extend other configs
BREAKING CHANGE: `extends` key in config file is now used for extending other config files
```

#### Commit message with ! to draw attention to breaking change
```
feat!: send an email to the customer when a product is shipped
```

#### Commit message with scope and ! to draw attention to breaking change
```
feat(api)!: send an email to the customer when a product is shipped
```

#### Commit message with both ! and BREAKING CHANGE footer
```
feat!: drop support for Node 6
BREAKING CHANGE: use JavaScript features not available in Node 6.
```

#### Commit message with no body
```
docs: correct spelling of CHANGELOG
```

#### Commit message with scope
```
feat(lang): add Polish language
```

#### Commit message with multi-paragraph body and multiple footers
```
fix: prevent racing of requests

Introduce a request id and a reference to latest request. Dismiss
incoming responses other than from latest request.

Remove timeouts which were used to mitigate the racing issue but are
obsolete now.

Reviewed-by: Z
Refs: #123
```

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
- Generate a conventional subject from the real changes
- Run one commit for the whole current working tree
- Push the active branch to `origin`
