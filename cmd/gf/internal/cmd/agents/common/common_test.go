// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file contains tests for the resource-agnostic common subpackage.
// Tests cover the generic ResolveTargets selector resolver and the
// shared StatusGlyph mapping; per-resource Inspect / ApplyOneLink /
// ApplyOneUnlink behavior is exercised through each resource subpackage's
// own test suite.

package common

import (
	"testing"
)

// fakeSpec is a minimal SpecLike implementation used by ResolveTargets
// tests. It carries just enough information to exercise the all-expansion
// policy filters and missing-name lookup paths.
type fakeSpec struct {
	name     string
	category Category
}

func (s fakeSpec) SpecName() string        { return s.name }
func (s fakeSpec) SpecDisplayName() string { return s.name }
func (s fakeSpec) SpecCategory() Category  { return s.category }
func (s fakeSpec) SpecSourcePath() string  { return ".source" }
func (s fakeSpec) SpecProjectPath() string { return ".project" }
func (s fakeSpec) SpecKind() Kind          { return KindDir }

func makeFakeRegistry() []fakeSpec {
	return []fakeSpec{
		{name: "a-native", category: CategoryNative},
		{name: "b-link-1", category: CategoryLink},
		{name: "c-link-2", category: CategoryLink},
		{name: "d-root", category: CategoryRootCollision},
	}
}

func TestResolveTargetsAllExpansionExcludesNativeAndRootByDefault(t *testing.T) {
	registry := makeFakeRegistry()
	got, err := ResolveTargets([]string{SelectorAll}, registry, TargetPolicy{})
	if err != nil {
		t.Fatalf("ResolveTargets(all): %v", err)
	}
	if len(got) == 0 {
		t.Fatalf("expected link-class agents in default expansion")
	}
	for _, spec := range got {
		if spec.SpecCategory() != CategoryLink {
			t.Fatalf("agent %s leaked into default all expansion (category=%s)",
				spec.SpecName(), spec.SpecCategory())
		}
	}
}

func TestResolveTargetsAllExpansionIncludesPolicyOptIns(t *testing.T) {
	registry := makeFakeRegistry()
	got, err := ResolveTargets([]string{SelectorAll}, registry, TargetPolicy{
		IncludeNative:        true,
		IncludeRootCollision: true,
	})
	if err != nil {
		t.Fatalf("ResolveTargets(all,full): %v", err)
	}
	if len(got) != len(registry) {
		t.Fatalf("expected full expansion to cover %d agents, got %d", len(registry), len(got))
	}
}

func TestResolveTargetsUnknownAgentReturnsError(t *testing.T) {
	registry := makeFakeRegistry()
	if _, err := ResolveTargets([]string{"no-such-agent"}, registry, TargetPolicy{}); err == nil {
		t.Fatalf("expected unknown agent error")
	}
}

func TestNormalizeAgentName(t *testing.T) {
	// Install the same alias map the unified registry uses so this package
	// test does not depend on importing registry (would force common→registry
	// for every common unit test).
	SetAgentAliasResolver(func(normalized string) string {
		aliases := map[string]string{
			"kimi-code-cli":   "kimi-cli",
			"antigravity-cli": "antigravity",
			"qoder-cn":        "qoder",
			"trae-cn":         "trae",
			"zenflow":         "zencoder",
		}
		if canonical, ok := aliases[normalized]; ok {
			return canonical
		}
		return normalized
	})
	t.Cleanup(func() { SetAgentAliasResolver(nil) })

	cases := []struct {
		input string
		want  string
	}{
		{input: "", want: ""},
		{input: "  ", want: ""},
		{input: "claude-code", want: "claude-code"},
		{input: "ClaudeCode", want: "claude-code"},
		{input: "Claude Code", want: "claude-code"},
		{input: "claude_code", want: "claude-code"},
		{input: "QwenCode", want: "qwen-code"},
		// Product aliases: vendor alternate ids collapse to one registry name.
		{input: "kimi-code-cli", want: "kimi-cli"},
		{input: "KimiCodeCLI", want: "kimi-cli"},
		{input: "antigravity-cli", want: "antigravity"},
		{input: "qoder-cn", want: "qoder"},
		{input: "trae-cn", want: "trae"},
		{input: "zenflow", want: "zencoder"},
	}
	for _, testCase := range cases {
		if got := NormalizeAgentName(testCase.input); got != testCase.want {
			t.Fatalf("NormalizeAgentName(%q) got=%q want=%q", testCase.input, got, testCase.want)
		}
	}
}

func TestParseSelectorsNormalizesAgentNames(t *testing.T) {
	got := ParseSelectors(" ClaudeCode , qwen_code , Claude Code ")
	want := []string{"claude-code", "qwen-code", "claude-code"}
	if len(got) != len(want) {
		t.Fatalf("ParseSelectors length got=%v want=%v", got, want)
	}
	for index := range got {
		if got[index] != want[index] {
			t.Fatalf("ParseSelectors[%d] got=%q want=%q", index, got[index], want[index])
		}
	}
}

func TestResolveTargetsNormalizesExplicitNames(t *testing.T) {
	registry := []fakeSpec{
		{name: "claude-code", category: CategoryLink},
		{name: "qwen-code", category: CategoryLink},
	}
	got, err := ResolveTargets([]string{"ClaudeCode", "qwen_code", "Claude Code"}, registry, TargetPolicy{})
	if err != nil {
		t.Fatalf("ResolveTargets normalized explicit: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected de-duplicated 2 specs, got %d: %v", len(got), got)
	}
	if got[0].SpecName() != "claude-code" || got[1].SpecName() != "qwen-code" {
		t.Fatalf("expected sorted normalized output, got %v", got)
	}
}

func TestResolveTargetsExplicitNamesReturnsRequestedSpecs(t *testing.T) {
	registry := makeFakeRegistry()
	got, err := ResolveTargets([]string{"b-link-1", "a-native"}, registry, TargetPolicy{})
	if err != nil {
		t.Fatalf("ResolveTargets explicit: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 specs, got %d", len(got))
	}
	// Output is sorted by SpecName, so a-native should precede b-link-1.
	if got[0].SpecName() != "a-native" || got[1].SpecName() != "b-link-1" {
		t.Fatalf("expected sorted output, got %v", got)
	}
}

func TestResolveTargetsEmptySelectorsReturnsNil(t *testing.T) {
	registry := makeFakeRegistry()
	got, err := ResolveTargets([]string{}, registry, TargetPolicy{})
	if err != nil {
		t.Fatalf("ResolveTargets empty: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil for empty selectors, got %v", got)
	}
}

func TestStatusGlyphCoversAllStatuses(t *testing.T) {
	cases := []struct {
		status Status
		want   string
	}{
		{StatusOK, "[+]"},
		{StatusCreated, "[+]"},
		{StatusRebuilt, "[+]"},
		{StatusRemoved, "[+]"},
		{StatusMismatch, "[~]"},
		{StatusSkippedForeignTarget, "[~]"},
		{StatusSkippedNotManaged, "[~]"},
		{StatusAbsent, "[.]"},
		{StatusNative, "[.]"},
		{StatusConflict, "[!]"},
		{StatusSkippedRootCollision, "[*]"},
		{StatusError, "[?]"},
		{Status("unknown"), "[?]"},
	}
	for _, testCase := range cases {
		if got := StatusGlyph(testCase.status); got != testCase.want {
			t.Fatalf("StatusGlyph(%s) got=%s want=%s", testCase.status, got, testCase.want)
		}
	}
}

func TestHasErrorReports(t *testing.T) {
	if HasError([]Result{{Status: StatusOK}}) {
		t.Fatalf("expected false")
	}
	if !HasError([]Result{{Status: StatusError}}) {
		t.Fatalf("expected true")
	}
}
