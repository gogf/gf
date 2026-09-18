// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file is the single source of truth for AI coding agent bindings
// across skills (.agents/skills), prompts (.agents/prompts) and md
// (AGENTS.md). Add or edit agents here only; resource packages project
// the relevant Binding fields into common.SpecLike for the shared engine.
//
// Vendor alternate ids (e.g. kimi-code-cli) belong in Aliases, not as a
// second Agent row, so `make agents` never shows duplicate labels.

package registry

import "github.com/gogf/gf/cmd/gf/v2/internal/cmd/agents/common"

// agents is the unified product registry. Sorted alphabetically by Name
// in init() for stable iteration. Empty Binding fields mean the agent is
// not registered for that resource.
var agents = []Agent{
	{
		Name:        "adal",
		DisplayName: "AdaL",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".adal/skills"},
	},
	{
		Name:        "aider-desk",
		DisplayName: "AiderDesk",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".aider-desk/skills"},
		MD:          Binding{Category: common.CategoryLink, ProjectPath: "CONVENTIONS.md"},
	},
	{
		Name:        "amp",
		DisplayName: "Amp",
		Skills:      Binding{Category: common.CategoryNative, ProjectPath: ".agents/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "antigravity",
		DisplayName: "Antigravity",
		Aliases:     []string{"antigravity-cli"},
		Skills:      Binding{Category: common.CategoryNative, ProjectPath: ".agents/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "augment",
		DisplayName: "Augment",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".augment/skills"},
		MD:          Binding{Category: common.CategoryLink, ProjectPath: ".augment-guidelines"},
	},
	{
		Name:        "autohand-code",
		DisplayName: "Autohand Code CLI",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".autohand/skills"},
	},
	{
		Name:        "bob",
		DisplayName: "IBM Bob",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".bob/skills"},
	},
	{
		Name:        "claude",
		DisplayName: "Claude Code",
		Aliases:     []string{"claude-code"},
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".claude/skills"},
		Prompts:     Binding{Category: common.CategoryLink, ProjectPath: ".claude/commands", SourcePath: ".agents/prompts"},
		MD:          Binding{Category: common.CategoryLink, ProjectPath: "CLAUDE.md"},
	},
	{
		Name:        "cline",
		DisplayName: "Cline",
		Skills:      Binding{Category: common.CategoryNative, ProjectPath: ".agents/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "codearts-agent",
		DisplayName: "CodeArts Agent",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".codeartsdoer/skills"},
	},
	{
		Name:        "codebuddy",
		DisplayName: "CodeBuddy",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".codebuddy/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "codemaker",
		DisplayName: "Codemaker",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".codemaker/skills"},
	},
	{
		Name:        "codestudio",
		DisplayName: "Code Studio",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".codestudio/skills"},
	},
	{
		Name:        "codex",
		DisplayName: "Codex",
		Skills:      Binding{Category: common.CategoryNative, ProjectPath: ".agents/skills"},
		Prompts:     Binding{Category: common.CategoryLink, ProjectPath: ".codex/prompts", SourcePath: ".agents/prompts"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "command-code",
		DisplayName: "Command Code",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".commandcode/skills"},
	},
	{
		Name:        "continue",
		DisplayName: "Continue",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".continue/skills"},
		MD:          Binding{Category: common.CategoryLink, ProjectPath: ".continuerules"},
	},
	{
		Name:        "cortex",
		DisplayName: "Cortex Code",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".cortex/skills"},
	},
	{
		Name:        "crush",
		DisplayName: "Crush",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".crush/skills"},
		MD:          Binding{Category: common.CategoryLink, ProjectPath: "CRUSH.md"},
	},
	{
		Name:        "cursor",
		DisplayName: "Cursor",
		Skills:      Binding{Category: common.CategoryNative, ProjectPath: ".agents/skills"},
		Prompts:     Binding{Category: common.CategoryLink, ProjectPath: ".cursor/commands", SourcePath: ".agents/prompts"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "deepagents",
		DisplayName: "Deep Agents",
		Skills:      Binding{Category: common.CategoryNative, ProjectPath: ".agents/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "devin",
		DisplayName: "Devin for Terminal",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".devin/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "dexto",
		DisplayName: "Dexto",
		Skills:      Binding{Category: common.CategoryNative, ProjectPath: ".agents/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "droid",
		DisplayName: "Droid",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".factory/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "firebender",
		DisplayName: "Firebender",
		Skills:      Binding{Category: common.CategoryNative, ProjectPath: ".agents/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "forgecode",
		DisplayName: "ForgeCode",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".forge/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "gemini-cli",
		DisplayName: "Gemini CLI",
		Skills:      Binding{Category: common.CategoryNative, ProjectPath: ".agents/skills"},
		Prompts:     Binding{Category: common.CategoryLink, ProjectPath: ".gemini/commands", SourcePath: ".agents/prompts"},
		MD:          Binding{Category: common.CategoryLink, ProjectPath: "GEMINI.md"},
	},
	{
		Name:        "github-copilot",
		DisplayName: "GitHub Copilot",
		Skills:      Binding{Category: common.CategoryNative, ProjectPath: ".agents/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "goose",
		DisplayName: "Goose",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".goose/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "grok",
		DisplayName: "Grok",
		Skills:      Binding{Category: common.CategoryNative, ProjectPath: ".agents/skills"},
		Prompts:     Binding{Category: common.CategoryLink, ProjectPath: ".grok/commands", SourcePath: ".agents/prompts"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "hermes-agent",
		DisplayName: "Hermes Agent",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".hermes/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "iflow-cli",
		DisplayName: "iFlow CLI",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".iflow/skills"},
		MD:          Binding{Category: common.CategoryLink, ProjectPath: "IFLOW.md"},
	},
	{
		Name:        "inference-sh",
		DisplayName: "inference.sh",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".inferencesh/skills"},
	},
	{
		Name:        "jazz",
		DisplayName: "Jazz",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".jazz/skills"},
	},
	{
		Name:        "junie",
		DisplayName: "Junie",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".junie/skills"},
		MD:          Binding{Category: common.CategoryLink, ProjectPath: ".junie/guidelines.md"},
	},
	{
		Name:        "kilo",
		DisplayName: "Kilo Code",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".kilocode/skills"},
	},
	{
		Name:        "kimi-cli",
		DisplayName: "Kimi Code CLI",
		Aliases:     []string{"kimi-code-cli"},
		Skills:      Binding{Category: common.CategoryNative, ProjectPath: ".agents/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "kiro-cli",
		DisplayName: "Kiro CLI",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".kiro/skills"},
	},
	{
		Name:        "kode",
		DisplayName: "Kode",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".kode/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "lingma",
		DisplayName: "Lingma",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".lingma/skills"},
	},
	{
		Name:        "loaf",
		DisplayName: "Loaf",
		Skills:      Binding{Category: common.CategoryNative, ProjectPath: ".agents/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "mcpjam",
		DisplayName: "MCPJam",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".mcpjam/skills"},
	},
	{
		Name:        "mistral-vibe",
		DisplayName: "Mistral Vibe",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".vibe/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "moxby",
		DisplayName: "Moxby",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".moxby/skills"},
	},
	{
		Name:        "mux",
		DisplayName: "Mux",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".mux/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "neovate",
		DisplayName: "Neovate",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".neovate/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "ona",
		DisplayName: "Ona",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".ona/skills"},
	},
	{
		Name:        "openclaw",
		DisplayName: "OpenClaw",
		Skills:      Binding{Category: common.CategoryRootCollision, ProjectPath: "skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "opencode",
		DisplayName: "OpenCode",
		Skills:      Binding{Category: common.CategoryNative, ProjectPath: ".agents/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "openhands",
		DisplayName: "OpenHands",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".openhands/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "pi",
		DisplayName: "Pi",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".pi/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "pochi",
		DisplayName: "Pochi",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".pochi/skills"},
	},
	{
		Name:        "promptscript",
		DisplayName: "PromptScript",
		Skills:      Binding{Category: common.CategoryNative, ProjectPath: ".agents/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "qoder",
		DisplayName: "Qoder",
		Aliases:     []string{"qoder-cn"},
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".qoder/skills"},
	},
	{
		Name:        "qwen-code",
		DisplayName: "Qwen Code",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".qwen/skills"},
		MD:          Binding{Category: common.CategoryLink, ProjectPath: "QWEN.md"},
	},
	{
		Name:        "reasonix",
		DisplayName: "Reasonix",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".reasonix/skills"},
	},
	{
		Name:        "replit",
		DisplayName: "Replit",
		Skills:      Binding{Category: common.CategoryNative, ProjectPath: ".agents/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "roo",
		DisplayName: "Roo Code",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".roo/skills"},
		MD:          Binding{Category: common.CategoryLink, ProjectPath: ".roo/rules/AGENTS.md"},
	},
	{
		Name:        "rovodev",
		DisplayName: "Rovo Dev",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".rovodev/skills"},
	},
	{
		Name:        "tabnine-cli",
		DisplayName: "Tabnine CLI",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".tabnine/agent/skills"},
		MD:          Binding{Category: common.CategoryLink, ProjectPath: "TABNINE.md"},
	},
	{
		Name:        "terramind",
		DisplayName: "Terramind",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".terramind/skills"},
	},
	{
		Name:        "tinycloud",
		DisplayName: "Tinycloud",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".tinycloud/skills"},
	},
	{
		Name:        "trae",
		DisplayName: "Trae",
		Aliases:     []string{"trae-cn"},
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".trae/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "universal",
		DisplayName: "Universal",
		Skills:      Binding{Category: common.CategoryNative, ProjectPath: ".agents/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "warp",
		DisplayName: "Warp",
		Skills:      Binding{Category: common.CategoryNative, ProjectPath: ".agents/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "windsurf",
		DisplayName: "Windsurf",
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".windsurf/skills"},
		MD:          Binding{Category: common.CategoryLink, ProjectPath: ".windsurfrules"},
	},
	{
		Name:        "zed",
		DisplayName: "Zed",
		Skills:      Binding{Category: common.CategoryNative, ProjectPath: ".agents/skills"},
		MD:          Binding{Category: common.CategoryNative},
	},
	{
		Name:        "zencoder",
		DisplayName: "Zencoder",
		Aliases:     []string{"zenflow"},
		Skills:      Binding{Category: common.CategoryLink, ProjectPath: ".zencoder/skills"},
	},
}
