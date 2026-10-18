# gf agents

`gf agents` links local AI coding-agent paths onto the project's canonical sources.

## Required project layout

Run the command from a GoFrame project (or any subdirectory) that already has this layout:

```text
.
├── AGENTS.md
└── .agents/
    ├── skills/
    └── prompts/
```

| Path | Role |
| --- | --- |
| `AGENTS.md` | Project rules file. Per-tool files such as `CLAUDE.md` are linked to it. |
| `.agents/` | Canonical AI resource directory. Must exist as a directory. |
| `.agents/skills/` | Shared skills. Linked to paths such as `.claude/skills`. |
| `.agents/prompts/` | Optional slash-command / prompt files. Linked to paths such as `.claude/commands`. |

If `.agents/` or `AGENTS.md` is missing, `gf agents` does nothing: it does not create, rebuild, or remove files, and it does not fail the command.

## Usage

```bash
gf agents
gf agents --agent=claude
gf agents --agent=claude --action=unlink
gf agents --agent=claude --force
```

In the GoFrame repository, `make agents` is a wrapper around the same command:

```bash
make agents agent=claude
```

`--agent` names one supported tool. `claude-code` is accepted as an alias of `claude`. `--action` defaults to `link`. `--force` rebuilds a mismatched symlink.
