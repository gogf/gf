// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file contains unit tests for the agents aggregate command's
// non-interactive surface and the cross-resource agent universe builder.
// The huh-driven interactive menu (runAgentInteractiveMenu) is exercised
// manually because automating arrow-key TUI input is brittle and
// charmbracelet/huh ships its own test coverage for form mechanics.
//
// Coverage in this file:
//   - collectAgentUniverse merges three registries, includes native-only
//     agents and is sorted alphabetically by agent name.
//   - partitionAgentUniverse / buildAgentPickerOptions split built-in and
//     needs-setup groups with alphabetical order inside each group.
//   - validateSingleAgentName rejects empty / "all" / comma-separated /
//     unknown values and accepts a single supported agent.
//   - parseAgentSetupAction normalizes link/unlink and rejects others.
//   - runAgents prints the usage hint when invoked without an agent on
//     a non-TTY stdin (the standard CI path).
//   - runAgents accepts a known agent on a non-TTY stdin and dispatches
//     the per-resource execute helpers without prompting.

package agents

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// newTestRequest builds a Request wired with buffers so interactive
// helpers degrade to their non-TTY safe paths.
func newTestRequest(t *testing.T) (Request, *bytes.Buffer) {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".agents"), 0o755); err != nil {
		t.Fatalf("mkdir .agents: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("# AGENTS.md\n"), 0o644); err != nil {
		t.Fatalf("write AGENTS.md: %v", err)
	}
	stdout := &bytes.Buffer{}
	req := Request{
		Root:   root,
		Stdin:  strings.NewReader(""),
		Stdout: stdout,
	}
	return req, stdout
}

func TestCollectAgentUniverseIncludesNativeOnly(t *testing.T) {
	universe := collectAgentUniverse(t.TempDir())
	if len(universe) == 0 {
		t.Fatalf("expected non-empty agent universe")
	}
	var (
		hasBuiltin bool
		hasSetup   bool
	)
	for _, agent := range universe {
		if agent.needsSetup() {
			hasSetup = true
			continue
		}
		hasBuiltin = true
	}
	if !hasBuiltin {
		t.Fatalf("expected at least one built-in (native-only) agent in universe")
	}
	if !hasSetup {
		t.Fatalf("expected at least one needs-setup agent in universe")
	}
}

func TestCollectAgentUniverseSortedAlphabetically(t *testing.T) {
	universe := collectAgentUniverse(t.TempDir())
	if len(universe) < 2 {
		t.Fatalf("expected at least two agents for alphabetical sort check")
	}
	for index := 1; index < len(universe); index++ {
		if universe[index-1].Name >= universe[index].Name {
			t.Fatalf("universe not alphabetically sorted: %q before %q", universe[index-1].Name, universe[index].Name)
		}
	}
}

func TestPartitionAgentUniverseGroupsAndOrder(t *testing.T) {
	universe := collectAgentUniverse(t.TempDir())
	builtin, setup := partitionAgentUniverse(universe)
	if len(builtin) == 0 || len(setup) == 0 {
		t.Fatalf("expected both groups non-empty; builtin=%d setup=%d", len(builtin), len(setup))
	}
	for _, agent := range builtin {
		if agent.needsSetup() {
			t.Fatalf("built-in group contains needs-setup agent %q", agent.Name)
		}
	}
	for _, agent := range setup {
		if !agent.needsSetup() {
			t.Fatalf("setup group contains built-in agent %q", agent.Name)
		}
	}
	for index := 1; index < len(builtin); index++ {
		if builtin[index-1].Name >= builtin[index].Name {
			t.Fatalf("built-in group not alphabetical: %q before %q", builtin[index-1].Name, builtin[index].Name)
		}
	}
	for index := 1; index < len(setup); index++ {
		if setup[index-1].Name >= setup[index].Name {
			t.Fatalf("setup group not alphabetical: %q before %q", setup[index-1].Name, setup[index].Name)
		}
	}
}

func TestBuildAgentPickerOptionsTwoGroups(t *testing.T) {
	universe := collectAgentUniverse(t.TempDir())
	options := buildAgentPickerOptions(universe)
	if len(options) < 4 {
		t.Fatalf("expected section headers plus agents; got %d options", len(options))
	}
	if options[0].Value != agentPickerSectionBuiltin || !options[0].Section {
		t.Fatalf("first option should be built-in section header, got value=%q section=%v", options[0].Value, options[0].Section)
	}

	var (
		sawSetupSection bool
		builtinNames    []string
		setupNames      []string
	)
	for _, option := range options[1:] {
		switch {
		case option.Value == agentPickerSectionSetup:
			sawSetupSection = true
		case isAgentPickerSection(option.Value):
			t.Fatalf("unexpected section value %q", option.Value)
		case !sawSetupSection:
			builtinNames = append(builtinNames, option.Value)
		default:
			setupNames = append(setupNames, option.Value)
		}
	}
	if !sawSetupSection {
		t.Fatalf("expected needs-setup section header")
	}
	if len(builtinNames) == 0 || len(setupNames) == 0 {
		t.Fatalf("expected agents in both groups; builtin=%v setup=%v", builtinNames, setupNames)
	}
	for index := 1; index < len(builtinNames); index++ {
		if builtinNames[index-1] >= builtinNames[index] {
			t.Fatalf("picker built-in names not alphabetical: %q before %q", builtinNames[index-1], builtinNames[index])
		}
	}
	for index := 1; index < len(setupNames); index++ {
		if setupNames[index-1] >= setupNames[index] {
			t.Fatalf("picker setup names not alphabetical: %q before %q", setupNames[index-1], setupNames[index])
		}
	}
}

func TestIsAgentPickerSection(t *testing.T) {
	if !isAgentPickerSection(agentPickerSectionBuiltin) || !isAgentPickerSection(agentPickerSectionSetup) {
		t.Fatalf("section markers should be detected")
	}
	if isAgentPickerSection("claude") {
		t.Fatalf("real agent name must not be treated as section")
	}
}

func TestCollectAgentUniverseMergesAcrossRegistries(t *testing.T) {
	universe := collectAgentUniverse(t.TempDir())
	// claude-code is link-class in both skills and md; its roles map
	// must contain at least those two resources.
	for _, agent := range universe {
		if agent.Name != "claude" {
			continue
		}
		if _, ok := agent.Roles[resourceSkills]; !ok {
			t.Fatalf("claude missing skills role: %v", agent.Roles)
		}
		if _, ok := agent.Roles[resourceMd]; !ok {
			t.Fatalf("claude missing md role: %v", agent.Roles)
		}
		return
	}
	t.Fatalf("claude not found in universe")
}

func TestSelectableAgentOptionLabelUsesDisplayNameOnly(t *testing.T) {
	universe := collectAgentUniverse(t.TempDir())
	cases := map[string]string{
		"claude": "Claude Code",
		"codex":  "Codex",
		"cursor": "Cursor",
	}
	for name, want := range cases {
		t.Run(name, func(t *testing.T) {
			agent, ok := lookupAgent(universe, name)
			if !ok {
				t.Fatalf("agent %q not found in universe", name)
			}
			got := agent.optionLabel()
			if got != want {
				t.Fatalf("optionLabel(%q) got=%q want=%q", name, got, want)
			}
			for _, forbidden := range []string{name, "skills:", "prompts:", "md:", "link", "native", "["} {
				if strings.Contains(got, forbidden) {
					t.Fatalf("option label %q should not include internal fragment %q", got, forbidden)
				}
			}
		})
	}
}

func TestValidateSingleAgentName(t *testing.T) {
	universe := collectAgentUniverse(t.TempDir())
	cases := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "valid claude-code alias", input: "claude-code", want: "claude", wantErr: false},
		{name: "valid ClaudeCode", input: "ClaudeCode", want: "claude", wantErr: false},
		{name: "valid Claude Code", input: "Claude Code", want: "claude", wantErr: false},
		{name: "valid claude_code", input: "claude_code", want: "claude", wantErr: false},
		{name: "valid native amp", input: "amp", want: "amp", wantErr: false},
		{name: "empty", input: "", wantErr: true},
		{name: "literal all", input: "all", wantErr: true},
		{name: "case-insensitive all", input: "ALL", wantErr: true},
		{name: "csv", input: "claude,codex", wantErr: true},
		{name: "unknown", input: "no-such-agent", wantErr: true},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := validateSingleAgentName(testCase.input, universe)
			if testCase.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q, got value %q", testCase.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", testCase.input, err)
			}
			if got != testCase.want {
				t.Fatalf("validate(%q) got=%q want=%q", testCase.input, got, testCase.want)
			}
		})
	}
}

func TestRunAgentsClaudeAliasCreatesLinks(t *testing.T) {
	req, stdout := newTestRequest(t)
	if err := os.MkdirAll(filepath.Join(req.Root, ".agents", "skills"), 0o755); err != nil {
		t.Fatalf("mkdir skills: %v", err)
	}
	if err := os.WriteFile(filepath.Join(req.Root, "AGENTS.md"), []byte("# AGENTS.md\n"), 0o644); err != nil {
		t.Fatalf("write AGENTS.md: %v", err)
	}
	req.Agent = "claude-code"
	if err := Run(req); err != nil {
		t.Fatalf("Run claude-code alias: %v", err)
	}
	if !strings.Contains(stdout.String(), "Agent: Claude Code") {
		t.Fatalf("alias output: %q", stdout.String())
	}
	for _, rel := range []string{
		filepath.Join(".claude", "skills"),
		filepath.Join(".claude", "commands"),
		"CLAUDE.md",
	} {
		info, err := os.Lstat(filepath.Join(req.Root, rel))
		if err != nil || info.Mode()&os.ModeSymlink == 0 {
			t.Fatalf("expected symlink %s: %v", rel, err)
		}
	}
}

func TestRunAgentsNormalizesOneShotAgentName(t *testing.T) {
	req, stdout := newTestRequest(t)
	if err := func() error { req.Agent = "ClaudeCode"; return Run(req) }(); err != nil {
		t.Fatalf("runAgents normalized: %v", err)
	}
	output := stdout.String()
	for _, fragment := range []string{
		"Agent: Claude Code",
		"Action: link",
		"RESOURCE",
		"skills",
		"prompts",
		"md",
	} {
		if !strings.Contains(output, fragment) {
			t.Fatalf("normalized one-shot output missing %q; got %q", fragment, output)
		}
	}
}

func TestRunAgentsBuiltinAgentPrintsReadySummary(t *testing.T) {
	req, stdout := newTestRequest(t)
	// amp is skills-native only. An explicit --agent still prints the
	// compact resource table and skips native / unregistered resources.
	if err := func() error { req.Agent = "amp"; return Run(req) }(); err != nil {
		t.Fatalf("runAgents builtin: %v", err)
	}
	output := stdout.String()
	for _, fragment := range []string{
		"Agent: Amp",
		"Action: link",
		"skills",
		"skipped",
		"native",
	} {
		if !strings.Contains(output, fragment) {
			t.Fatalf("builtin summary missing %q; got %q", fragment, output)
		}
	}
}

func TestCollectAgentUniverseIncludesMainstreamTools(t *testing.T) {
	universe := collectAgentUniverse(t.TempDir())
	required := map[string]bool{
		"claude":     false,
		"codex":      false,
		"cursor":     false,
		"gemini-cli": false,
		"grok":       false,
		"zed":        false,
		"lingma":     false,
	}
	for _, agent := range universe {
		if _, ok := required[agent.Name]; ok {
			required[agent.Name] = true
		}
	}
	for name, found := range required {
		if !found {
			t.Fatalf("expected mainstream agent %q in make agents universe", name)
		}
	}
}

// TestNeedsSetupIgnoresPrompts locks the product rule: built-in vs
// needs-setup is decided only by skills (.agents) and md (AGENTS.md).
// prompts may still be link-class without pushing the agent out of built-in.
func TestNeedsSetupIgnoresPrompts(t *testing.T) {
	universe := collectAgentUniverse(t.TempDir())

	// skills + md native (prompts may still be link) → built-in
	for _, name := range []string{"codex", "cursor", "grok"} {
		agent, ok := lookupAgent(universe, name)
		if !ok {
			t.Fatalf("expected %q in universe", name)
		}
		if agent.needsSetup() {
			t.Fatalf("%s must be built-in (skills/md native); roles=%v", name, agent.Roles)
		}
	}

	// skills and/or md still need a symlink → needs-setup
	for _, name := range []string{"claude", "lingma"} {
		agent, ok := lookupAgent(universe, name)
		if !ok {
			t.Fatalf("expected %q in universe", name)
		}
		if !agent.needsSetup() {
			t.Fatalf("%s must be needs-setup (skills/md link); roles=%v", name, agent.Roles)
		}
	}
}

func TestRunAgentsCodexPrintsBuiltinReady(t *testing.T) {
	req, stdout := newTestRequest(t)
	if err := func() error { req.Agent = "codex"; return Run(req) }(); err != nil {
		t.Fatalf("runAgents codex: %v", err)
	}
	output := stdout.String()
	for _, fragment := range []string{
		"Agent: Codex",
		"Action: link",
		"skills",
		"md",
		"prompts",
		"skipped",
		"applied",
	} {
		if !strings.Contains(output, fragment) {
			t.Fatalf("codex one-shot output missing %q; got %q", fragment, output)
		}
	}
}

func TestCollectAgentUniverseUniqueDisplayNames(t *testing.T) {
	universe := collectAgentUniverse(t.TempDir())
	seen := make(map[string]string, len(universe))
	for _, agent := range universe {
		label := agent.optionLabel()
		if previous, dup := seen[label]; dup {
			t.Fatalf("duplicate picker label %q for agents %q and %q", label, previous, agent.Name)
		}
		seen[label] = agent.Name
	}
}

func TestAgentAliasesResolveInAggregateCommand(t *testing.T) {
	req, stdout := newTestRequest(t)
	// kimi-code-cli is an alias of kimi-cli; both must resolve without
	// listing two "Kimi Code CLI" rows in the picker.
	if err := func() error { req.Agent = "kimi-code-cli"; return Run(req) }(); err != nil {
		t.Fatalf("runAgents kimi-code-cli alias: %v", err)
	}
	if !strings.Contains(stdout.String(), "Agent: Kimi Code CLI") {
		t.Fatalf("alias resolution missing Kimi Code CLI summary; got %q", stdout.String())
	}
}

func TestParseAgentSetupAction(t *testing.T) {
	cases := []struct {
		input    string
		fallback agentSetupAction
		want     agentSetupAction
		wantErr  bool
	}{
		{input: "", fallback: actionLink, want: actionLink},
		{input: "link", fallback: actionLink, want: actionLink},
		{input: "unlink", fallback: actionLink, want: actionUnlink},
		{input: "wat", fallback: actionLink, wantErr: true},
	}
	for _, testCase := range cases {
		got, err := parseAgentSetupAction(testCase.input, testCase.fallback)
		if testCase.wantErr {
			if err == nil {
				t.Fatalf("expected error for %q", testCase.input)
			}
			continue
		}
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", testCase.input, err)
		}
		if got != testCase.want {
			t.Fatalf("parse(%q) got=%q want=%q", testCase.input, got, testCase.want)
		}
	}
}

// TestRunAgentsNonTTYPrintsUsage verifies the no-agent, non-TTY path
// emits the usage hint and returns successfully (no error, no
// dispatch). This is the standard CI invocation: `gf agents`
// in a piped context should never block on input.
func TestRunAgentsNonTTYPrintsUsage(t *testing.T) {
	req, stdout := newTestRequest(t)
	if err := func() error { req.Agent = ""; return Run(req) }(); err != nil {
		t.Fatalf("runAgents: %v", err)
	}
	output := stdout.String()
	for _, fragment := range []string{
		"Usage: gf agents",
		"One-shot mode",
		"Interactive mode",
		"two-group list",
	} {
		if !strings.Contains(output, fragment) {
			t.Fatalf("usage hint missing %q; got %q", fragment, output)
		}
	}
}

// TestRunAgentsRejectsAgentAll guards the safety rule: agent=all is
// explicitly rejected by the aggregate command.
func TestRunAgentsRejectsAgentAll(t *testing.T) {
	req, _ := newTestRequest(t)
	err := func() error { req.Agent = "all"; return Run(req) }()
	if err == nil {
		t.Fatalf("expected error for agent=all")
	}
	if !strings.Contains(err.Error(), "all") {
		t.Fatalf("expected error to mention 'all', got %q", err)
	}
}

// TestRunAgentsRejectsCSV guards the safety rule: comma-separated lists
// are explicitly rejected.
func TestRunAgentsRejectsCSV(t *testing.T) {
	req, _ := newTestRequest(t)
	err := func() error { req.Agent = "claude,codex"; return Run(req) }()
	if err == nil {
		t.Fatalf("expected error for csv input")
	}
	if !strings.Contains(err.Error(), "comma-separated") {
		t.Fatalf("expected error to mention comma-separated, got %q", err)
	}
}

// TestRunAgentsUnknownAgentReportsCandidates verifies the error message
// for an unknown agent includes the candidate listing so users can
// recover without consulting docs.
func TestRunAgentsUnknownAgentReportsCandidates(t *testing.T) {
	req, _ := newTestRequest(t)
	err := func() error { req.Agent = "no-such-agent"; return Run(req) }()
	if err == nil {
		t.Fatalf("expected error for unknown agent")
	}
	if !strings.Contains(err.Error(), "supported agents") {
		t.Fatalf("expected candidate listing; got %q", err)
	}
}

// TestRunAgentsRejectsBadAction verifies typos in action surface at the
// CLI boundary rather than silently falling back.
func TestRunAgentsRejectsBadAction(t *testing.T) {
	req, _ := newTestRequest(t)
	err := func() error { req.Agent = "claude"; req.Action = "wat"; return Run(req) }()
	if err == nil {
		t.Fatalf("expected error for bad action")
	}
	if !strings.Contains(err.Error(), "invalid action") {
		t.Fatalf("expected invalid action error, got %q", err)
	}
}

// TestDispatchAgentSetupSkipsUnregisteredResources verifies the aggregate
// dispatcher prints one compact resource-level table.
func TestDispatchAgentSetupSkipsUnregisteredResources(t *testing.T) {
	req, stdout := newTestRequest(t)
	universe := collectAgentUniverse(req.Root)
	// Pick the first agent that is link-class in skills only, or fall
	// back to a known link-only-in-skills entry. `codebuddy` is a
	// strong candidate (skills link, no md/prompts entry today).
	target := "codebuddy"
	if _, ok := lookupAgent(universe, target); !ok {
		t.Skipf("expected %q in universe; got %v", target, universe)
	}

	// We expect no error in non-TTY because the temp dir has no
	// pre-existing collisions, so a fresh symlink should be created.
	if err := dispatchAgentSetup(&runtime{stdout: stdout, stdin: strings.NewReader(""), root: req.Root}, target, actionLink, false, universe); err != nil {
		t.Fatalf("dispatchAgentSetup: %v", err)
	}
	output := stdout.String()
	for _, fragment := range []string{
		"Agent: CodeBuddy",
		"Action: link",
		"RESOURCE",
		"STATUS",
		"DETAIL",
		"skills",
		"applied",
		"skipped",
		"not registered",
	} {
		if !strings.Contains(output, fragment) {
			t.Fatalf("expected compact output to contain %q; got %q", fragment, output)
		}
	}
	for _, forbidden := range []string{"Summary:", "== skills ==", "PROJECT PATH", "CATEGORY"} {
		if strings.Contains(output, forbidden) {
			t.Fatalf("compact aggregate output should not contain %q; got %q", forbidden, output)
		}
	}
}
