// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file implements resource-scoped link/unlink/plan operations against
// the unified agents registry. Callers pass a ResourceKind (skills /
// prompts / md); projections and the common engine do the rest. The
// former skills / prompts / md packages are intentionally gone so product
// maintenance stays in agents.go only.

package registry

import (
	"errors"
	"fmt"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/agents/common"
)

// LinkRequest captures one agents.<resource>.link invocation.
type LinkRequest struct {
	// Selectors is the list of agent names. Empty means "no selection".
	// "all" expands to every link-class agent for the resource.
	Selectors []string
	// Force rebuilds mismatched links and enables rootCollision creation
	// when the resource exposes that category (skills only today).
	Force bool
}

// UnlinkRequest captures one agents.<resource>.unlink invocation.
type UnlinkRequest struct {
	// Selectors mirrors LinkRequest.Selectors for unlink.
	Selectors []string
}

// Inspect returns the current Status for one resource projection.
func Inspect(repoRoot string, spec ResourceSpec) common.Result {
	return common.Inspect(repoRoot, spec)
}

// PlanList returns inspection results for every agent registered on kind.
func PlanList(repoRoot string, kind ResourceKind) []common.Result {
	specs := Specs(kind)
	out := make([]common.Result, 0, len(specs))
	for _, spec := range specs {
		out = append(out, common.Inspect(repoRoot, spec))
	}
	return out
}

// ApplyLink executes a link request for the given resource kind.
func ApplyLink(repoRoot string, kind ResourceKind, request LinkRequest) ([]common.Result, error) {
	if len(request.Selectors) == 0 {
		return nil, errors.New("no agent selected; pass agent=<name|all|csv>")
	}
	specs := Specs(kind)
	if len(specs) == 0 {
		return nil, fmt.Errorf("no agents registered for resource %q", kind)
	}
	policy := common.TargetPolicy{
		IncludeNative:        true,
		IncludeRootCollision: request.Force,
	}
	targets, err := common.ResolveTargets(request.Selectors, specs, policy)
	if err != nil {
		return nil, err
	}
	if len(targets) == 0 {
		return nil, errors.New("no agent selected")
	}
	results := make([]common.Result, 0, len(targets))
	for _, spec := range targets {
		results = append(results, common.ApplyOneLink(repoRoot, spec, request.Force))
	}
	return results, nil
}

// ApplyUnlink executes an unlink request for the given resource kind.
func ApplyUnlink(repoRoot string, kind ResourceKind, request UnlinkRequest) ([]common.Result, error) {
	if len(request.Selectors) == 0 {
		return nil, errors.New("no agent selected; pass agent=<name|all|csv>")
	}
	specs := Specs(kind)
	if len(specs) == 0 {
		return nil, fmt.Errorf("no agents registered for resource %q", kind)
	}
	policy := common.TargetPolicy{
		IncludeNative:        false,
		IncludeRootCollision: false,
	}
	targets, err := common.ResolveTargets(request.Selectors, specs, policy)
	if err != nil {
		return nil, err
	}
	if len(targets) == 0 {
		return nil, errors.New("no agent selected")
	}
	results := make([]common.Result, 0, len(targets))
	for _, spec := range targets {
		results = append(results, common.ApplyOneUnlink(repoRoot, spec))
	}
	return results, nil
}

// LinkCandidates returns interactive link candidates (link-class only).
func LinkCandidates(repoRoot string, kind ResourceKind) []common.SelectableEntry {
	out := make([]common.SelectableEntry, 0)
	for _, spec := range Specs(kind) {
		if spec.Category != common.CategoryLink {
			continue
		}
		result := common.Inspect(repoRoot, spec)
		out = append(out, common.SelectableEntry{
			Spec:          spec,
			CurrentStatus: result.Status,
			Detail:        result.Detail,
		})
	}
	return out
}

// UnlinkCandidates returns interactive unlink candidates (managed links only).
func UnlinkCandidates(repoRoot string, kind ResourceKind) []common.SelectableEntry {
	out := make([]common.SelectableEntry, 0)
	for _, spec := range Specs(kind) {
		if spec.Category == common.CategoryNative {
			continue
		}
		result := common.Inspect(repoRoot, spec)
		if result.Status != common.StatusOK {
			continue
		}
		out = append(out, common.SelectableEntry{
			Spec:          spec,
			CurrentStatus: result.Status,
			Detail:        result.Detail,
		})
	}
	return out
}

// Specs returns resource projections for kind.
func Specs(kind ResourceKind) []ResourceSpec {
	switch kind {
	case ResourceSkills:
		return SkillsSpecs()
	case ResourcePrompts:
		return PromptsSpecs()
	case ResourceMD:
		return MDSpecs()
	default:
		return nil
	}
}
