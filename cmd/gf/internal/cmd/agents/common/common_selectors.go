// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file implements selector parsing and target resolution against a
// generic agent registry. ResolveTargets is parameterized over the
// concrete AgentSpec type so callers receive their own type back without
// runtime assertions; a SpecLike constraint lets the resolver consult
// Name and Category fields uniformly.

package common

import (
	"fmt"
	"slices"
	"sort"
	"strings"
	"unicode"
)

// SelectorAll is the special selector value that targets every link-class
// agent. native and rootCollision agents are skipped by default with this
// selector; rootCollision agents only execute when force=true is also set.
const SelectorAll = "all"

// ParseSelectors parses a comma-separated agent selector value. Empty
// values and whitespace-only tokens are dropped. An empty input yields a
// nil slice.
func ParseSelectors(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		token := NormalizeAgentName(part)
		if token == "" {
			continue
		}
		out = append(out, token)
	}
	return out
}

// agentAliasResolver optionally maps a kebab-case selector onto a
// canonical agent Name (e.g. kimi-code-cli → kimi-cli). The unified
// agents registry installs the real resolver during its package init so
// common does not import registry (avoids an import cycle).
var agentAliasResolver = func(normalized string) string { return normalized }

// SetAgentAliasResolver installs the alias lookup used by
// NormalizeAgentName. Passing nil restores the identity resolver.
// Intended for registry package init and tests.
func SetAgentAliasResolver(resolve func(normalized string) string) {
	if resolve == nil {
		agentAliasResolver = func(normalized string) string { return normalized }
		return
	}
	agentAliasResolver = resolve
}

// NormalizeAgentName converts user-facing agent selector input to the
// canonical kebab-case identifier used by registries. It accepts common
// variants such as "ClaudeCode", "Claude Code" and "claude_code", and
// then applies the installed alias resolver (when configured).
func NormalizeAgentName(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	normalized := caseKebab(value)
	return agentAliasResolver(normalized)
}

// caseKebab converts user-facing identifiers to kebab-case without pulling
// the GoFrame root module into this standalone tool.
func caseKebab(value string) string {
	runes := []rune(strings.TrimSpace(value))
	if len(runes) == 0 {
		return ""
	}
	var builder strings.Builder
	prevIsDash := false
	prevIsLower := false
	prevIsUpper := false
	for index, r := range runes {
		switch {
		case r == '_' || r == ' ' || r == '-':
			if builder.Len() > 0 && !prevIsDash {
				builder.WriteByte('-')
				prevIsDash = true
			}
			prevIsLower = false
			prevIsUpper = false
		case unicode.IsUpper(r):
			nextIsLower := index+1 < len(runes) && unicode.IsLower(runes[index+1])
			if builder.Len() > 0 && !prevIsDash && (prevIsLower || (prevIsUpper && nextIsLower)) {
				builder.WriteByte('-')
			}
			builder.WriteRune(unicode.ToLower(r))
			prevIsDash = false
			prevIsLower = false
			prevIsUpper = true
		default:
			builder.WriteRune(unicode.ToLower(r))
			prevIsDash = false
			prevIsLower = unicode.IsLetter(r)
			prevIsUpper = false
		}
	}
	return builder.String()
}

// TargetPolicy controls which agent categories an "all" selector expands
// to. Specific agent names always match regardless of policy; the policy
// only affects the special "all" expansion.
type TargetPolicy struct {
	// IncludeNative includes native-class agents in expansion. They are
	// always reported in status output regardless of this flag.
	IncludeNative bool
	// IncludeRootCollision includes rootCollision-class agents. Should be
	// set to true only when force is also true.
	IncludeRootCollision bool
}

// ResolveTargets returns the agents in the registry matched by selectors.
// When selectors contains SelectorAll the policy filters apply; otherwise
// specific agent names are looked up and missing names are returned as a
// single error listing every unknown name in input order.
//
// The result is sorted by SpecName for stable rendering.
func ResolveTargets[S SpecLike](selectors []string, registry []S, policy TargetPolicy) ([]S, error) {
	if len(selectors) == 0 {
		return nil, nil
	}
	normalizedSelectors := normalizeSelectors(selectors)
	if hasAll(normalizedSelectors) {
		out := make([]S, 0, len(registry))
		for _, spec := range registry {
			switch spec.SpecCategory() {
			case CategoryNative:
				if policy.IncludeNative {
					out = append(out, spec)
				}
			case CategoryLink:
				out = append(out, spec)
			case CategoryRootCollision:
				if policy.IncludeRootCollision {
					out = append(out, spec)
				}
			}
		}
		return out, nil
	}
	byName := make(map[string]S, len(registry))
	for _, spec := range registry {
		byName[spec.SpecName()] = spec
	}
	seen := make(map[string]struct{}, len(normalizedSelectors))
	out := make([]S, 0, len(normalizedSelectors))
	var unknown []string
	for _, name := range normalizedSelectors {
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}
		spec, ok := byName[name]
		if !ok {
			unknown = append(unknown, name)
			continue
		}
		out = append(out, spec)
	}
	if len(unknown) > 0 {
		return nil, fmt.Errorf("unknown agent(s): %s", strings.Join(unknown, ", "))
	}
	sort.Slice(out, func(left, right int) bool {
		return out[left].SpecName() < out[right].SpecName()
	})
	return out, nil
}

// hasAll reports whether a selector list contains SelectorAll.
func hasAll(selectors []string) bool {
	return slices.Contains(selectors, SelectorAll)
}

// normalizeSelectors applies NormalizeAgentName to a selector list while
// preserving input order and dropping empty tokens.
func normalizeSelectors(selectors []string) []string {
	out := make([]string, 0, len(selectors))
	for _, selector := range selectors {
		normalized := NormalizeAgentName(selector)
		if normalized == "" {
			continue
		}
		out = append(out, normalized)
	}
	return out
}
