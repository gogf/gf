// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file contains unit tests for the unified agents registry:
// integrity of Names/DisplayNames, alias resolution, and resource
// projections used by skills / prompts / md packages.

package registry

import (
	"strings"
	"testing"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/agents/common"
)

func TestAgentsNonEmptyAndSorted(t *testing.T) {
	specs := Agents()
	if len(specs) == 0 {
		t.Fatalf("expected non-empty unified registry")
	}
	for index := 1; index < len(specs); index++ {
		if specs[index-1].Name >= specs[index].Name {
			t.Fatalf("registry not sorted: %q before %q", specs[index-1].Name, specs[index].Name)
		}
	}
}

func TestUniqueNamesAndDisplayNames(t *testing.T) {
	names := make(map[string]struct{})
	displays := make(map[string]string)
	for _, agent := range Agents() {
		if agent.Name == "" {
			t.Fatalf("agent missing Name: %+v", agent)
		}
		if _, dup := names[agent.Name]; dup {
			t.Fatalf("duplicate Name %q", agent.Name)
		}
		names[agent.Name] = struct{}{}
		display := strings.TrimSpace(agent.DisplayName)
		if display == "" {
			t.Fatalf("agent %s missing DisplayName", agent.Name)
		}
		if previous, dup := displays[display]; dup {
			t.Fatalf("duplicate DisplayName %q on %q and %q", display, previous, agent.Name)
		}
		displays[display] = agent.Name
	}
}

func TestUniqueLinkProjectPathsPerResource(t *testing.T) {
	check := func(kind ResourceKind, specs []ResourceSpec) {
		t.Helper()
		seen := make(map[string]string)
		for _, spec := range specs {
			if spec.Category == common.CategoryNative {
				continue
			}
			path := spec.ProjectPath
			if path == "" {
				t.Fatalf("%s agent %s missing ProjectPath for category %s", kind, spec.Name, spec.Category)
			}
			if previous, dup := seen[path]; dup {
				t.Fatalf("%s duplicate ProjectPath %q on %q and %q", kind, path, previous, spec.Name)
			}
			seen[path] = spec.Name
		}
	}
	check(ResourceSkills, SkillsSpecs())
	check(ResourcePrompts, PromptsSpecs())
	check(ResourceMD, MDSpecs())
}

func TestRequiredMainstreamAgents(t *testing.T) {
	for _, name := range []string{"claude", "claude-code", "codex", "cursor", "gemini-cli", "grok"} {
		agent, ok := Find(name)
		if !ok {
			t.Fatalf("expected mainstream agent %q", name)
		}
		if !agent.Skills.Registered() {
			t.Fatalf("%s must register Skills", name)
		}
	}
	// Codex / Grok: skills+md native; prompts optional link.
	for _, name := range []string{"codex", "grok"} {
		agent, _ := Find(name)
		if agent.Skills.Category != common.CategoryNative {
			t.Fatalf("%s skills want native, got %s", name, agent.Skills.Category)
		}
		if agent.MD.Category != common.CategoryNative {
			t.Fatalf("%s md want native, got %s", name, agent.MD.Category)
		}
		if agent.Prompts.Category != common.CategoryLink {
			t.Fatalf("%s prompts want link, got %s", name, agent.Prompts.Category)
		}
	}
}

func TestAliasesResolve(t *testing.T) {
	cases := map[string]string{
		"kimi-code-cli":   "kimi-cli",
		"antigravity-cli": "antigravity",
		"qoder-cn":        "qoder",
		"trae-cn":         "trae",
		"zenflow":         "zencoder",
		"claude-code":     "claude",
		"codex":           "codex",
	}
	for input, want := range cases {
		if got := ResolveAlias(input); got != want {
			t.Fatalf("ResolveAlias(%q) got=%q want=%q", input, got, want)
		}
		agent, ok := Find(input)
		if !ok {
			t.Fatalf("Find(%q) missing", input)
		}
		if agent.Name != want {
			t.Fatalf("Find(%q).Name got=%q want=%q", input, agent.Name, want)
		}
	}
}

func TestNormalizeAgentNameUsesRegistryAliases(t *testing.T) {
	// registry.init installs the resolver into common.
	if got := common.NormalizeAgentName("kimi-code-cli"); got != "kimi-cli" {
		t.Fatalf("NormalizeAgentName via registry resolver: got=%q want=kimi-cli", got)
	}
	if got := common.NormalizeAgentName("ClaudeCode"); got != "claude" {
		t.Fatalf("NormalizeAgentName kebab: got=%q want=claude", got)
	}
	if got := common.NormalizeAgentName("claude-code"); got != "claude" {
		t.Fatalf("NormalizeAgentName alias: got=%q want=claude", got)
	}
}

func TestFindSpecAndResourceLookups(t *testing.T) {
	if _, ok := Find("no-such-agent"); ok {
		t.Fatalf("unknown agent should be missing")
	}
	skills, ok := FindSkillsSpec("claude-code")
	if !ok || skills.ProjectPath != ".claude/skills" {
		t.Fatalf("FindSkillsSpec claude-code: %+v ok=%v", skills, ok)
	}
	prompts, ok := FindPromptsSpec("codex")
	if !ok || prompts.ProjectPath != ".codex/prompts" {
		t.Fatalf("FindPromptsSpec codex: %+v ok=%v", prompts, ok)
	}
	md, ok := FindMDSpec("claude-code")
	if !ok || md.ProjectPath != "CLAUDE.md" {
		t.Fatalf("FindMDSpec claude-code: %+v ok=%v", md, ok)
	}
	if _, ok := FindSpec("claude-code", ResourceKind("nope")); ok {
		t.Fatalf("unknown resource kind should be missing")
	}
	copied := AliasMap()
	if copied["kimi-code-cli"] != "kimi-cli" {
		t.Fatalf("AliasMap missing kimi-code-cli: %v", copied)
	}
	copied["sentinel"] = "x"
	if _, ok := AliasMap()["sentinel"]; ok {
		t.Fatalf("AliasMap should return a copy")
	}
}

func TestUnlinkCandidatesOnlyManagedLinks(t *testing.T) {
	root := newSkillsFixture(t)
	if _, err := ApplyLink(root, ResourceSkills, LinkRequest{Selectors: []string{"claude-code"}}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	found := false
	for _, entry := range UnlinkCandidates(root, ResourceSkills) {
		if entry.Spec.SpecName() == "claude" {
			found = true
		}
		if entry.CurrentStatus != common.StatusOK {
			t.Fatalf("unlink candidate %s status=%s", entry.Spec.SpecName(), entry.CurrentStatus)
		}
	}
	if !found {
		t.Fatalf("expected claude in unlink candidates")
	}
}

func TestInspectDelegatesToCommon(t *testing.T) {
	root := newSkillsFixture(t)
	spec, ok := FindSkillsSpec("claude-code")
	if !ok {
		t.Fatalf("missing claude-code skills spec")
	}
	result := Inspect(root, spec)
	if result.Status != common.StatusAbsent {
		t.Fatalf("expected absent, got %s", result.Status)
	}
}

func TestSpecsUnknownKind(t *testing.T) {
	if got := Specs(ResourceKind("nope")); got != nil {
		t.Fatalf("unknown kind should return nil, got %d", len(got))
	}
}

func TestProjectionsMatchBindings(t *testing.T) {
	skills := SkillsSpecs()
	if len(skills) == 0 {
		t.Fatalf("expected skills projections")
	}
	for _, spec := range skills {
		if spec.Kind != common.KindDir {
			t.Fatalf("skills %s kind=%v want dir", spec.Name, spec.Kind)
		}
		if spec.SourcePath != SkillsSourceDir {
			t.Fatalf("skills %s source=%q want %q", spec.Name, spec.SourcePath, SkillsSourceDir)
		}
	}
	for _, spec := range MDSpecs() {
		if spec.Kind != common.KindFile {
			t.Fatalf("md %s kind=%v want file", spec.Name, spec.Kind)
		}
		if spec.SourcePath != MDSourceFile {
			t.Fatalf("md %s source=%q want %q", spec.Name, spec.SourcePath, MDSourceFile)
		}
	}
	for _, spec := range PromptsSpecs() {
		if spec.Kind != common.KindDir {
			t.Fatalf("prompts %s kind=%v want dir", spec.Name, spec.Kind)
		}
		if spec.Category != common.CategoryLink {
			t.Fatalf("prompts %s category=%s want link", spec.Name, spec.Category)
		}
	}
}
