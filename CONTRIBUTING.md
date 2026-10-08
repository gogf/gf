# Contributing

Thanks for taking the time to join our community and start contributing!

## With issues

- Use the search tool before opening a new issue.
- Please provide source code and commit sha if you found a bug.
- Review existing issues and provide feedback or react to them.

## With pull requests

- Open your pull request against `master`
- Your pull request should have no more than two commits, if not you should squash them.
- It should pass all tests in the available continuous integrations systems such as GitHub CI.
- You should add/modify tests to cover your proposed code changes.
- If your pull request contains a new feature, please document it on the README.

## Agent setup

Canonical AI coding resources are `.agents/` and `AGENTS.md`. Per-tool files such as `CLAUDE.md` are local symlinks and are gitignored. After cloning, run:

```bash
make agents agent=claude
gf agents --agent=claude
```

Use `make agents` or `gf agents` on a TTY to pick an agent interactively.
