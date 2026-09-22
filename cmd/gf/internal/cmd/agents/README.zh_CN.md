# gf agents

`gf agents`把本地 AI 编码助手路径软链到项目的规范源。

## 必需的项目目录结构

请在已经具备以下结构的 GoFrame 项目（或其子目录）中执行：

```text
.
├── AGENTS.md
└── .agents/
    ├── skills/
    └── prompts/
```

| 路径 | 作用 |
| --- | --- |
| `AGENTS.md` | 项目规范文件。`CLAUDE.md`等工具私有文件会软链到它。 |
| `.agents/` | 规范的 AI 资源目录，必须是目录。 |
| `.agents/skills/` | 共享 skills。会软链到`.claude/skills`等路径。 |
| `.agents/prompts/` | 可选的斜杠命令 / prompt。会软链到`.claude/commands`等路径。 |

如果缺少`.agents/`或`AGENTS.md`，`gf agents`什么都不做：不会创建、重建或删除文件，也不会把这次调用当成失败。

## 用法

```bash
gf agents
gf agents --agent=claude
gf agents --agent=claude --action=unlink
gf agents --agent=claude --force
```

在 GoFrame 仓库中，`make agents`是同一条命令的包装：

```bash
make agents agent=claude
```

`--agent`只接受一个受支持的工具名。`claude-code`是`claude`的别名。`--action`默认是`link`。`--force`会重建指向错误目标的软链。
