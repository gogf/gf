---
name: git-worktree
description: Create and actively use an isolated git worktree for the user's task, then continue the task inside that new directory. Use this whenever the user asks for a separate worktree, isolated checkout, clean branch directory, safer parallel changes, or a fresh workspace to avoid unrelated local edits.
---

# Git Worktree

Create a dedicated `git worktree` for the current task, then keep working inside that new directory instead of the original checkout.

This skill is about execution, not just advice. When it triggers, actually create the worktree unless the repository state makes that impossible.

Do not introduce helper scripts for this skill. Use direct `git` and shell commands inline.

## When To Use

- The user explicitly asks for a new `git worktree`, independent branch directory, or isolated workspace
- The current checkout contains unrelated local changes and isolation is the safest way forward
- The user wants parallel work on multiple tasks without stashing or disturbing the original worktree
- The user says things like "create a separate branch folder", "open a fresh worktree", "use a clean checkout", or "work in an isolated workspace"

## Core Rule

After creating the worktree, treat the new path as the active working directory for the rest of the task.

In any agent environment, "enter the directory" means:

- Run subsequent commands against the new worktree path
- Apply all edits under that worktree path
- Do not keep using the original checkout by accident
- Confirm the handoff by running at least one follow-up command in the new worktree

Never claim you "switched" unless your subsequent actions actually target the new `worktree_path`.

If your environment supports a per-command working directory, use it for every later command. If it does not, prefix later commands with an explicit `cd <worktree_path> && ...`.

## Name Derivation

- Derive a short ASCII kebab-case task slug from the user's real task, such as `login-timeout-fix` or `user-export`
- Do not use generic names like `git-worktree`, `new-worktree`, or `task` unless the request is too vague
- If the request is mostly non-ASCII or no good slug is obvious, fall back to `task-$(date +%Y%m%d-%H%M%S)`
- Default branch prefix is `worktree/`
- Default worktree directory is a sibling of the repository root, named `<repo-name>-<slug>`

## Default Workflow

1. Inspect the repository context from the current checkout:
   - `git rev-parse --show-toplevel`
   - `git branch --show-current`
   - `git status --short`
   - `git worktree list --porcelain`
2. Decide a task slug yourself using the rules above
3. Build branch and path names inline, then create the worktree with direct shell commands like:

   ```bash
   repo_root=$(git rev-parse --show-toplevel)
   repo_name=$(basename "$repo_root")
   parent_dir=$(dirname "$repo_root")
   source_branch=$(git -C "$repo_root" branch --show-current)
   if [ -n "$source_branch" ]; then
     source_ref="$source_branch"
   else
     source_ref="HEAD@$(git -C "$repo_root" rev-parse --short HEAD)"
   fi

   slug="<task-slug>"
   base_branch="worktree/$slug"
   branch_name="$base_branch"
   base_path="$parent_dir/$repo_name-$slug"
   worktree_path="$base_path"
   index=2

   while git -C "$repo_root" show-ref --verify --quiet "refs/heads/$branch_name" || [ -e "$worktree_path" ]; do
     branch_name="${base_branch}-$index"
     worktree_path="${base_path}-$index"
     index=$((index + 1))
   done

   git -C "$repo_root" worktree add -b "$branch_name" "$worktree_path" HEAD
   ```

4. Immediately verify the handoff inside the new worktree, for example:

   ```bash
   pwd
   git status --short --branch
   ```

   These verification commands must run against `worktree_path`.

5. Announce the new active path briefly, then continue the main task there
6. For the remainder of the task, use `worktree_path` as the working directory for every relevant command or edit operation

## Behavior Rules

- Default base ref is `HEAD` from the current checkout so uncommitted local changes are not dragged into the new worktree
- If a branch name or path already exists, auto-increment it instead of failing
- If you are already inside a non-default worktree and the user still wants another isolated workspace, create a new one from the current `HEAD`
- If the directory is not a Git repository, explain that clearly and do not pretend a worktree was created
- If worktree creation succeeds, continue the user's actual task instead of stopping at setup
- If worktree creation fails because of filesystem permissions, request the minimal approval needed and retry

## Uncommitted Change Policy

The safe default is isolation from uncommitted changes.

- If the source checkout is dirty, still create the new worktree from `HEAD` unless the user explicitly asks to carry local edits over
- Do not silently stash, reset, or move the user's existing changes
- If the user wants local edits copied into the new worktree, use an explicit flow such as a temporary commit, patch, or cherry-pick, and say what you are doing

## Output Contract

When you use this skill:

- Tell the user which branch and directory were created
- Make it clear that subsequent work is now happening inside that path
- Mention the source ref and whether the original checkout was dirty when that context matters
- Do not stop after setup if the user asked for additional work; continue the task in the new worktree

## Example

User request:

```text
Create a separate worktree for this task and then start implementing it.
```

Expected behavior:

- Inspect current repo status
- Create a new `worktree/...` branch and sibling directory with direct `git worktree` commands
- Switch all following commands to that directory
- Continue the requested implementation there

## Cleanup

Only remove a worktree when the user asks or when cleanup is clearly part of the task.

Before cleanup:

- Check status in the worktree you created
- Make sure you are removing the correct path
- Never remove the user's original checkout
