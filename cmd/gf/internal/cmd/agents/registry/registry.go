// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package registry is the single source of truth for AI coding agent
// product bindings used by make agents.
//
// Each Agent describes one product once. Skills, Prompts and MD bindings
// declare how that product consumes GoFrame's canonical resources:
//   - Skills  → .agents/skills
//   - Prompts → .agents/prompts
//   - MD      → AGENTS.md
//
// Resource packages (skills, prompts, md) project the matching Binding
// into common.SpecLike for the shared link/unlink engine. The aggregate
// agents command reads Agents() directly for the interactive picker.
package registry

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/agents/common"
)

// Canonical source paths for the three managed resources.
const (
	// SkillsSourceDir is the repo-relative skills root.
	SkillsSourceDir = ".agents/skills"
	// PromptsSourceDir is the repo-relative prompts root.
	PromptsSourceDir = ".agents/prompts"
	// MDSourceFile is the repo-root project rules file.
	MDSourceFile = "AGENTS.md"
)

// Binding describes one resource binding for an agent. A zero Binding
// (Category empty) means the agent is not registered for that resource.
type Binding struct {
	// Category controls link/unlink behavior. Empty means unregistered.
	Category common.Category
	// ProjectPath is the agent-private path where a symlink may live.
	// Native skills typically use SkillsSourceDir; native md leaves this
	// empty because the agent reads AGENTS.md directly.
	ProjectPath string
	// SourcePath optionally overrides the resource default source.
	// Empty means use SkillsSourceDir / PromptsSourceDir / MDSourceFile.
	SourcePath string
}

// Registered reports whether this binding participates in a resource.
func (b Binding) Registered() bool {
	return b.Category != ""
}

// Agent is one AI coding product with optional skills / prompts / md
// bindings. Name is the canonical CLI id; Aliases are alternate selectors
// resolved by common.NormalizeAgentName.
type Agent struct {
	// Name is the canonical kebab-case identifier (e.g. claude).
	Name string
	// DisplayName is the human-readable picker label (e.g. Claude Code).
	DisplayName string
	// Aliases are alternate selector ids that normalize to Name.
	Aliases []string
	// Skills binding for .agents/skills.
	Skills Binding
	// Prompts binding for .agents/prompts (optional slash commands).
	Prompts Binding
	// MD binding for AGENTS.md / private guide files.
	MD Binding
}

// ResourceKind identifies one of the three managed resource projections.
type ResourceKind string

const (
	// ResourceSkills projects the Skills binding.
	ResourceSkills ResourceKind = "skills"
	// ResourcePrompts projects the Prompts binding.
	ResourcePrompts ResourceKind = "prompts"
	// ResourceMD projects the MD binding.
	ResourceMD ResourceKind = "md"
)

// ResourceSpec is a resource-scoped projection of Agent that implements
// common.SpecLike for the shared link/unlink engine.
type ResourceSpec struct {
	Name        string
	DisplayName string
	SourcePath  string
	ProjectPath string
	Category    common.Category
	Kind        common.Kind
}

// SpecName implements common.SpecLike.
func (s ResourceSpec) SpecName() string { return s.Name }

// SpecDisplayName implements common.SpecLike.
func (s ResourceSpec) SpecDisplayName() string { return s.DisplayName }

// SpecCategory implements common.SpecLike.
func (s ResourceSpec) SpecCategory() common.Category { return s.Category }

// SpecSourcePath implements common.SpecLike.
func (s ResourceSpec) SpecSourcePath() string { return s.SourcePath }

// SpecProjectPath implements common.SpecLike.
func (s ResourceSpec) SpecProjectPath() string { return s.ProjectPath }

// SpecKind implements common.SpecLike.
func (s ResourceSpec) SpecKind() common.Kind { return s.Kind }

// init normalizes paths, builds the alias map, sorts agents by Name and
// installs the alias resolver into common.NormalizeAgentName.
func init() {
	for index := range agents {
		normalizeAgent(&agents[index])
	}
	sort.Slice(agents, func(left, right int) bool {
		return agents[left].Name < agents[right].Name
	})
	rebuildAliasMap()
	common.SetAgentAliasResolver(ResolveAlias)
}

func normalizeAgent(agent *Agent) {
	agent.Skills.ProjectPath = filepath.ToSlash(agent.Skills.ProjectPath)
	agent.Skills.SourcePath = filepath.ToSlash(agent.Skills.SourcePath)
	agent.Prompts.ProjectPath = filepath.ToSlash(agent.Prompts.ProjectPath)
	agent.Prompts.SourcePath = filepath.ToSlash(agent.Prompts.SourcePath)
	agent.MD.ProjectPath = filepath.ToSlash(agent.MD.ProjectPath)
	agent.MD.SourcePath = filepath.ToSlash(agent.MD.SourcePath)
	if agent.Skills.Registered() && agent.Skills.SourcePath == "" {
		agent.Skills.SourcePath = SkillsSourceDir
	}
	if agent.Skills.Registered() && agent.Skills.Category == common.CategoryNative && agent.Skills.ProjectPath == "" {
		agent.Skills.ProjectPath = SkillsSourceDir
	}
	if agent.Prompts.Registered() && agent.Prompts.SourcePath == "" {
		agent.Prompts.SourcePath = PromptsSourceDir
	}
	if agent.MD.Registered() && agent.MD.SourcePath == "" {
		agent.MD.SourcePath = MDSourceFile
	}
}

// aliasToCanonical maps alternate selector ids onto Agent.Name. Populated
// from Agent.Aliases during init.
var aliasToCanonical map[string]string

func rebuildAliasMap() {
	aliasToCanonical = make(map[string]string)
	for _, agent := range agents {
		for _, alias := range agent.Aliases {
			alias = strings.TrimSpace(alias)
			if alias == "" {
				continue
			}
			aliasToCanonical[alias] = agent.Name
		}
	}
}

// ResolveAlias returns the canonical agent Name for a normalized selector
// token. When the token is not an alias, it is returned unchanged.
func ResolveAlias(normalizedName string) string {
	if canonical, ok := aliasToCanonical[normalizedName]; ok {
		return canonical
	}
	return normalizedName
}

// AliasMap returns a copy of the alias → canonical name map for tests and
// common.NormalizeAgentName wiring.
func AliasMap() map[string]string {
	out := make(map[string]string, len(aliasToCanonical))
	for key, value := range aliasToCanonical {
		out[key] = value
	}
	return out
}

// Agents returns a defensive copy of the unified registry sorted by Name.
func Agents() []Agent {
	out := make([]Agent, len(agents))
	copy(out, agents)
	return out
}

// Find returns the Agent for the given canonical name, or false if missing.
func Find(name string) (Agent, bool) {
	name = ResolveAlias(name)
	for _, agent := range agents {
		if agent.Name == name {
			return agent, true
		}
	}
	return Agent{}, false
}

// SkillsSpecs returns SpecLike projections for every agent with a Skills binding.
func SkillsSpecs() []ResourceSpec {
	return project(ResourceSkills)
}

// PromptsSpecs returns SpecLike projections for every agent with a Prompts binding.
func PromptsSpecs() []ResourceSpec {
	return project(ResourcePrompts)
}

// MDSpecs returns SpecLike projections for every agent with an MD binding.
func MDSpecs() []ResourceSpec {
	return project(ResourceMD)
}

// FindSpec returns the projection for name on kind.
func FindSpec(name string, kind ResourceKind) (ResourceSpec, bool) {
	return findProjection(name, kind)
}

// FindSkillsSpec returns the skills projection for one agent name.
func FindSkillsSpec(name string) (ResourceSpec, bool) {
	return findProjection(name, ResourceSkills)
}

// FindPromptsSpec returns the prompts projection for one agent name.
func FindPromptsSpec(name string) (ResourceSpec, bool) {
	return findProjection(name, ResourcePrompts)
}

// FindMDSpec returns the md projection for one agent name.
func FindMDSpec(name string) (ResourceSpec, bool) {
	return findProjection(name, ResourceMD)
}

func project(kind ResourceKind) []ResourceSpec {
	out := make([]ResourceSpec, 0, len(agents))
	for _, agent := range agents {
		if spec, ok := agent.asResourceSpec(kind); ok {
			out = append(out, spec)
		}
	}
	return out
}

func findProjection(name string, kind ResourceKind) (ResourceSpec, bool) {
	agent, ok := Find(name)
	if !ok {
		return ResourceSpec{}, false
	}
	return agent.asResourceSpec(kind)
}

func (a Agent) asResourceSpec(kind ResourceKind) (ResourceSpec, bool) {
	var binding Binding
	var resourceKind common.Kind
	switch kind {
	case ResourceSkills:
		binding = a.Skills
		resourceKind = common.KindDir
	case ResourcePrompts:
		binding = a.Prompts
		resourceKind = common.KindDir
	case ResourceMD:
		binding = a.MD
		resourceKind = common.KindFile
	default:
		return ResourceSpec{}, false
	}
	if !binding.Registered() {
		return ResourceSpec{}, false
	}
	return ResourceSpec{
		Name:        a.Name,
		DisplayName: a.DisplayName,
		SourcePath:  binding.SourcePath,
		ProjectPath: binding.ProjectPath,
		Category:    binding.Category,
		Kind:        resourceKind,
	}, true
}
