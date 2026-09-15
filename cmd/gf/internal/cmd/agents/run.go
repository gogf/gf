// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file implements the single gf agents command. Users either pass
// --agent=NAME for a one-shot setup, or run from a TTY to pick an agent
// from a two-group list (built-in support first, needs-setup second).
//
// Grouping is based only on skills (.agents) and md (AGENTS.md):
//   - built-in: skills and md are native (or unregistered) — no setup
//   - needs-setup: skills and/or md still need a managed symlink
// prompts never moves an agent into needs-setup. When a caller names an
// agent explicitly, optional prompt links are still applied.

package agents

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/agents/common"
	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/agents/registry"
)

// Request is one gf agents invocation.
type Request struct {
	// Root is the project directory that contains .agents and AGENTS.md.
	// Empty means discover it by walking upward from the working directory.
	Root string
	// Agent is a single agent selector such as "claude".
	Agent string
	// Action is "link" or "unlink". Empty defaults to link.
	Action string
	// Force rebuilds mismatched links and git symlink stubs.
	Force bool
	// Stdin is the interactive input stream.
	Stdin io.Reader
	// Stdout is the status/output stream.
	Stdout io.Writer
}

// runtime stores process streams and the resolved project root.
type runtime struct {
	stdout io.Writer
	stdin  io.Reader
	root   string
}

// Run executes the aggregate agents command.
func Run(req Request) error {
	if req.Stdout == nil {
		req.Stdout = os.Stdout
	}
	if req.Stdin == nil {
		req.Stdin = os.Stdin
	}
	if req.Root == "" {
		root, found, err := DiscoverProjectRoot()
		if err != nil {
			return err
		}
		if !found {
			return skipInvalidProject(req.Stdout)
		}
		req.Root = root
	} else if !isAgentsProjectRoot(req.Root) {
		return skipInvalidProject(req.Stdout)
	}
	rt := &runtime{
		stdout: req.Stdout,
		stdin:  req.Stdin,
		root:   req.Root,
	}
	return runAgents(rt, req.Agent, req.Action, req.Force)
}

// Picker section marker values for the aggregate TTY agent list. They are
// never valid agent names; selecting one re-prompts the user.
const (
	agentPickerSectionBuiltin = "__section_builtin__"
	agentPickerSectionSetup   = "__section_setup__"
)

// Picker styles: built-in (native-ready) agents render in green; agents that
// still need symlink setup render in yellow. Section headers are plain text
// and styled by the shared list picker so cursor chrome stays consistent.
var (
	agentBuiltinStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	agentSetupStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
)

// agentSetupAction enumerates the actions the aggregate command can
// dispatch. Only link and unlink are supported today; the type exists
// so the dispatcher can be extended without leaking string literals
// into business logic.
type agentSetupAction string

const (
	actionLink   agentSetupAction = "link"
	actionUnlink agentSetupAction = "unlink"
)

// resourceKind tags a resource entry inside the cross-resource agent
// universe with a stable identifier matching the registry it came from.
type resourceKind string

const (
	resourceSkills  resourceKind = "skills"
	resourcePrompts resourceKind = "prompts"
	resourceMd      resourceKind = "md"
)

// selectableAgent describes one agent's role across the three resource
// registries (skills / prompts / md). It is built once at command entry
// by collectAgentUniverse and consumed by both the interactive picker
// (huh option labels) and the dispatch loop (which resources to act on).
type selectableAgent struct {
	// Name is the canonical agent identifier shared across registries.
	Name string
	// DisplayName is the human-readable label aggregated from whichever
	// registry first carries it.
	DisplayName string
	// Roles records the agent's category in each resource it appears
	// in. A resource not present in the map means the agent is not
	// registered there.
	Roles map[resourceKind]common.Category
}

// optionLabel returns the compact human-readable label for status output
// and the interactive picker. The picker intentionally hides internal
// IDs, resource paths and status summaries; those remain available in
// the per-resource commands and the built-in readiness summary.
func (s selectableAgent) optionLabel() string {
	if label := strings.TrimSpace(s.DisplayName); label != "" {
		return label
	}
	return s.Name
}

// needsSetup reports whether the aggregate `make agents` picker should
// treat this agent as "Needs setup" rather than "Built-in support".
//
// Classification intentionally considers only skills (.agents/skills)
// and md (AGENTS.md). prompts is excluded: slash-command / prompt
// directory symlinks are optional enhancements managed by
// agents.prompts.* and must not force Codex, Cursor, Grok, and similar
// AGENTS.md + .agents-native tools out of the built-in group.
//
// skills/md categories that still need work:
//   - CategoryLink: create a relative symlink
//   - CategoryRootCollision: opt-in via force=1 (e.g. openclaw skills/)
//
// Missing skills or md registration is not treated as needs-setup; the
// agent remains built-in when every registered skills/md role is native.
func (s selectableAgent) needsSetup() bool {
	for _, kind := range []resourceKind{resourceSkills, resourceMd} {
		category, present := s.Roles[kind]
		if !present {
			continue
		}
		if category == common.CategoryLink || category == common.CategoryRootCollision {
			return true
		}
	}
	return false
}

// writeLine prints a single line to the writer, wrapping any write error
// so the agents CLI never silently drops stdout failures (e.g. a broken pipe).
// The helper is shared by every agents.* command.
func writeLine(out io.Writer, line string) error {
	if _, err := fmt.Fprintln(out, line); err != nil {
		return fmt.Errorf("write output line: %w", err)
	}
	return nil
}

// writeLines prints multiple lines, returning the first write error.
func writeLines(out io.Writer, lines ...string) error {
	for _, line := range lines {
		if err := writeLine(out, line); err != nil {
			return err
		}
	}
	return nil
}

// stdinAsFile returns the *os.File backing the app's stdin when available,
// or nil when the test harness wired an in-memory reader. It is shared by
// every agents.* command's TTY-detection branch.
func stdinAsFile(rt *runtime) *os.File {
	if file, ok := rt.stdin.(*os.File); ok {
		return file
	}
	return nil
}

// runAgents dispatches the aggregate agents command. Behavior:
//
//   - agent=<name> [action=link|unlink] [force=1] : one-shot setup; the
//     selected action runs against every resource type the agent is
//     link-class in.
//   - no agent + TTY                              : two-step menu (pick
//     agent -> pick action) backed by huh.
//   - no agent + non-TTY                          : print usage and
//     return successfully.
//
// Unknown agents, the literal string "all", and comma-separated lists
// are explicitly rejected to keep the one-shot path safe; multi-agent
// batch flows remain available through the per-resource
// `agents.<resource>.<action>` subcommands.
func runAgents(rt *runtime, rawAgent string, rawAction string, force bool) error {
	rawAgent = strings.TrimSpace(rawAgent)
	rawAction = strings.TrimSpace(strings.ToLower(rawAction))

	universe := collectAgentUniverse(rt.root)

	if rawAgent == "" {
		if !common.IsInteractiveTerminal(stdinAsFile(rt)) {
			return printAgentsUsage(rt.stdout)
		}
		return runAgentInteractiveMenu(rt, universe, force)
	}

	agentName, err := validateSingleAgentName(rawAgent, universe)
	if err != nil {
		return err
	}
	action, err := parseAgentSetupAction(rawAction, actionLink)
	if err != nil {
		return err
	}
	return dispatchAgentSetup(rt, agentName, action, force, universe)
}

// printAgentsUsage emits the non-interactive usage hint pointing
// callers at the one-shot mode and the advanced subcommands.
func printAgentsUsage(out io.Writer) error {
	return writeLines(out,
		"Usage: gf agents [--agent=NAME] [--action=link|unlink] [--force]",
		"",
		"One-shot mode (works in any environment):",
		"  gf agents --agent=claude",
		"  make agents agent=claude",
		"  - agent must name a single supported agent (no 'all', no csv).",
		"  - action defaults to 'link'.",
		"  - The selected action runs against every resource type the agent supports.",
		"  - Built-in agents natively read .agents/skills + AGENTS.md;",
		"    optional prompt links are still created on an explicit --agent.",
		"",
		"Interactive mode (TTY only):",
		"  gf agents",
		"  make agents",
		"  - Step 1: arrow-key pick the agent from a two-group list",
		"    (built-in = skills/md native; needs-setup = skills/md need link;",
		"    both groups A-Z; prompts is ignored for grouping).",
		"  - Step 2: for needs-setup agents, pick link or unlink;",
		"    built-in agents only show a readiness summary.",
	)
}

// parseAgentSetupAction normalizes an action string (link/unlink). An
// empty value falls back to fallback. Any other value yields an error
// so typos are caught at the CLI boundary.
func parseAgentSetupAction(raw string, fallback agentSetupAction) (agentSetupAction, error) {
	switch raw {
	case "":
		return fallback, nil
	case string(actionLink):
		return actionLink, nil
	case string(actionUnlink):
		return actionUnlink, nil
	default:
		return "", fmt.Errorf("invalid action %q: expected 'link' or 'unlink'", raw)
	}
}

// validateSingleAgentName enforces the one-shot mode contract: exactly
// one supported agent name, no "all", no comma list. The candidate
// listing in error messages covers the full cross-resource universe
// (built-in and needs-setup) so users can discover native-ready tools.
func validateSingleAgentName(raw string, universe []selectableAgent) (string, error) {
	normalized := common.NormalizeAgentName(raw)
	if normalized == "" {
		return "", fmt.Errorf("agent= must be set; pass a single supported agent name")
	}
	if strings.Contains(raw, ",") {
		return "", fmt.Errorf("agent=%q: comma-separated lists are not supported by `gf agents`; pass a single agent name", raw)
	}
	if normalized == common.SelectorAll {
		return "", fmt.Errorf("agent=all is not supported by `gf agents` (safety guard); pass a specific agent name")
	}
	for _, candidate := range universe {
		if candidate.Name == normalized {
			return candidate.Name, nil
		}
	}
	return "", fmt.Errorf("unknown agent %q; supported agents: %s", raw, joinAgentNames(universe))
}

// joinAgentNames flattens the universe slice into a comma-separated
// candidate listing for error messages.
func joinAgentNames(universe []selectableAgent) string {
	names := make([]string, 0, len(universe))
	for _, agent := range universe {
		names = append(names, agent.Name)
	}
	return strings.Join(names, ", ")
}

// collectAgentUniverse builds the aggregate picker universe from the
// unified agents registry. Every product appears once; Roles records the
// category of each registered resource binding (skills / prompts / md).
// The returned slice is sorted alphabetically by canonical agent name.
//
// The repoRoot parameter is intentionally unused. It remains in the
// signature because callers already pass it and future registries may need
// repository context, but the aggregate picker does not inspect runtime
// link state while building its compact labels.
func collectAgentUniverse(_ string) []selectableAgent {
	all := registry.Agents()
	out := make([]selectableAgent, 0, len(all))
	for _, agent := range all {
		entry := selectableAgent{
			Name:        agent.Name,
			DisplayName: agent.DisplayName,
			Roles:       map[resourceKind]common.Category{},
		}
		if agent.Skills.Registered() {
			entry.Roles[resourceSkills] = agent.Skills.Category
		}
		if agent.Prompts.Registered() {
			entry.Roles[resourcePrompts] = agent.Prompts.Category
		}
		if agent.MD.Registered() {
			entry.Roles[resourceMd] = agent.MD.Category
		}
		out = append(out, entry)
	}
	sort.Slice(out, func(left, right int) bool {
		return out[left].Name < out[right].Name
	})
	return out
}

// partitionAgentUniverse splits agents into built-in ready (no symlink
// work) and needs-setup groups. Both output slices keep alphabetical
// order by agent Name.
func partitionAgentUniverse(universe []selectableAgent) (builtin []selectableAgent, setup []selectableAgent) {
	for _, agent := range universe {
		if agent.needsSetup() {
			setup = append(setup, agent)
			continue
		}
		builtin = append(builtin, agent)
	}
	return builtin, setup
}

// buildAgentPickerOptions builds the aggregate TTY select options as two
// colored groups: built-in support first, needs-setup second. Section
// headers use reserved values that are not real agent names.
func buildAgentPickerOptions(universe []selectableAgent) []common.SingleOption {
	builtin, setup := partitionAgentUniverse(universe)
	options := make([]common.SingleOption, 0, len(universe)+2)
	if len(builtin) > 0 {
		options = append(options, common.SingleOption{
			Value:   agentPickerSectionBuiltin,
			Label:   "── Built-in support (no setup needed) ──",
			Section: true,
		})
		for _, agent := range builtin {
			options = append(options, common.SingleOption{
				Value: agent.Name,
				Label: agentBuiltinStyle.Render(agent.optionLabel()),
			})
		}
	}
	if len(setup) > 0 {
		options = append(options, common.SingleOption{
			Value:   agentPickerSectionSetup,
			Label:   "── Needs setup (symlink) ──",
			Section: true,
		})
		for _, agent := range setup {
			options = append(options, common.SingleOption{
				Value: agent.Name,
				Label: agentSetupStyle.Render(agent.optionLabel()),
			})
		}
	}
	return options
}

// isAgentPickerSection reports whether value is a visual section header
// rather than a real agent name.
func isAgentPickerSection(value string) bool {
	return value == agentPickerSectionBuiltin || value == agentPickerSectionSetup
}

// runAgentInteractiveMenu drives the TTY flow:
//  1. select one agent from the two-group list (built-in / needs-setup);
//  2. for needs-setup agents, select link or unlink; for built-in agents,
//     print a readiness summary and exit.
//
// Cancellation at either step returns nil after printing "Cancelled.".
func runAgentInteractiveMenu(rt *runtime, universe []selectableAgent, force bool) error {
	if len(universe) == 0 {
		return writeLine(rt.stdout, "No agents are registered; nothing to configure.")
	}

	options := buildAgentPickerOptions(universe)
	if len(options) == 0 {
		return writeLine(rt.stdout, "No agents are registered; nothing to configure.")
	}

	agentName, err := common.PromptSingleSelection(rt.stdin, rt.stdout, "Select an agent to configure:", options)
	if err != nil {
		return err
	}
	if agentName == "" {
		return writeLine(rt.stdout, "Cancelled.")
	}
	// Defensive: section rows are not selectable in the list picker, but
	// reject reserved values if a future caller constructs them as normal
	// options without Section=true.
	if isAgentPickerSection(agentName) {
		return writeLine(rt.stdout, "Please select a tool name (section headers are not selectable).")
	}

	target, ok := lookupAgent(universe, agentName)
	if !ok {
		return fmt.Errorf("agent %q not found in the cross-resource registry", agentName)
	}

	// Built-in agents already read canonical paths; only show readiness.
	if !target.needsSetup() {
		return printBuiltinAgentReady(rt.stdout, target)
	}

	actionChoice, err := common.PromptSingleSelection(rt.stdin, rt.stdout,
		fmt.Sprintf("What should we do for %s?", agentName), []common.SingleOption{
			{Value: string(actionLink), Label: "link    Create or rebuild symlinks for this agent"},
			{Value: string(actionUnlink), Label: "unlink  Remove managed symlinks for this agent"},
		})
	if err != nil {
		return err
	}
	if actionChoice == "" {
		return writeLine(rt.stdout, "Cancelled.")
	}

	action, err := parseAgentSetupAction(actionChoice, actionLink)
	if err != nil {
		return err
	}
	return dispatchAgentSetup(rt, agentName, action, force, universe)
}

// printBuiltinAgentReady explains that the selected agent already works
// against GoFrame skills (.agents) and AGENTS.md without creating managed
// symlinks. prompts is listed for transparency but is not part of the
// built-in / needs-setup classification.
func printBuiltinAgentReady(out io.Writer, target selectableAgent) error {
	lines := []string{
		fmt.Sprintf("Agent: %s", target.optionLabel()),
		"Status: built-in support (no setup needed)",
		"",
		"Built-in means this tool natively reads .agents/skills and AGENTS.md.",
		"No managed symlink is required for skills / project rules:",
		"",
	}
	for _, kind := range []resourceKind{resourceSkills, resourceMd} {
		category, present := target.Roles[kind]
		if !present {
			lines = append(lines, fmt.Sprintf("  %-7s  not registered", kind))
			continue
		}
		lines = append(lines, fmt.Sprintf("  %-7s  %s", kind, category))
	}
	if category, present := target.Roles[resourcePrompts]; present {
		lines = append(lines,
			"",
			fmt.Sprintf("  %-7s  %s (optional; not required for built-in)", resourcePrompts, category),
		)
		if category == common.CategoryLink {
			lines = append(lines,
				"  Tip: gf agents --agent="+target.Name,
				"       also creates the optional prompts symlink on an explicit --agent.",
			)
		}
	} else {
		lines = append(lines,
			"",
			fmt.Sprintf("  %-7s  not registered (optional)", resourcePrompts),
		)
	}
	lines = append(lines,
		"",
		"Tip: clone the repo and start coding. Optional prompt links are created",
		"when you pass an explicit --agent to gf agents.",
	)
	return writeLines(out, lines...)
}

// dispatchAgentSetup executes the chosen action across every resource
// type the agent participates in and renders a compact resource-level
// summary. Per-resource commands still own the verbose path/category
// tables; the aggregate command stays focused on the high-level outcome
// for the selected agent.
//
// Built-in agents (skills + AGENTS.md already native) short-circuit to a
// readiness summary. prompts-only link work is not triggered here so
// Codex/Grok/etc. stay zero-config; use agents.prompts.link when the
// optional prompt directory bridge is desired.
func dispatchAgentSetup(rt *runtime, agentName string, action agentSetupAction, force bool, universe []selectableAgent) error {
	target, ok := lookupAgent(universe, agentName)
	if !ok {
		// Should not happen: validateSingleAgentName already rejected
		// unknown names. Defensive guard so callers that bypass the
		// validator (future internal callers) still get a clear error.
		return fmt.Errorf("agent %q not found in the cross-resource registry", agentName)
	}

	if err := writeLines(rt.stdout,
		fmt.Sprintf("Agent: %s", target.optionLabel()),
		fmt.Sprintf("Action: %s", action),
		"",
	); err != nil {
		return err
	}

	outcomes := make([]aggregateResourceOutcome, 0, 3)
	var firstErr error

	for _, kind := range []resourceKind{resourceSkills, resourcePrompts, resourceMd} {
		category, present := target.Roles[kind]
		if !present {
			outcomes = append(outcomes, aggregateResourceOutcome{
				kind:   kind,
				status: aggregateStatusSkipped,
				detail: "not registered",
			})
			continue
		}
		if category != common.CategoryLink {
			outcomes = append(outcomes, aggregateResourceOutcome{
				kind:   kind,
				status: aggregateStatusSkipped,
				detail: fmt.Sprintf("%s (no symlink work)", category),
			})
			continue
		}
		results, err := runAggregateResourceAction(rt, kind, agentName, action, force)
		if err != nil {
			outcomes = append(outcomes, aggregateResourceOutcome{
				kind:   kind,
				status: aggregateStatusFailed,
				detail: err.Error(),
			})
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		outcome := aggregateOutcomeFromResults(kind, results)
		outcomes = append(outcomes, outcome)
		if outcome.status == aggregateStatusFailed && firstErr == nil {
			firstErr = fmt.Errorf("one or more agents failed; see result table")
		}
	}

	if err := renderAggregateOutcomes(rt.stdout, outcomes); err != nil {
		return err
	}

	return firstErr
}

// aggregateResourceOutcome is one row in the aggregate agents result
// table. status intentionally uses coarse values so the table stays
// scannable; the original per-resource Status is included in detail.
type aggregateResourceOutcome struct {
	kind   resourceKind
	status string
	detail string
}

const (
	aggregateStatusApplied = "applied"
	aggregateStatusSkipped = "skipped"
	aggregateStatusFailed  = "failed"
)

// runAggregateResourceAction executes one resource without rendering the
// verbose per-resource table. The aggregate command renders its own
// compact table after all resources have been processed.
func runAggregateResourceAction(rt *runtime, kind resourceKind, agentName string, action agentSetupAction, force bool) ([]common.Result, error) {
	selectors := []string{agentName}
	resource, err := registryResourceKind(kind)
	if err != nil {
		return nil, err
	}
	if action == actionLink {
		return registry.ApplyLink(rt.root, resource, registry.LinkRequest{Selectors: selectors, Force: force})
	}
	return registry.ApplyUnlink(rt.root, resource, registry.UnlinkRequest{Selectors: selectors})
}

// registryResourceKind maps aggregate resourceKind onto registry.ResourceKind.
func registryResourceKind(kind resourceKind) (registry.ResourceKind, error) {
	switch kind {
	case resourceSkills:
		return registry.ResourceSkills, nil
	case resourcePrompts:
		return registry.ResourcePrompts, nil
	case resourceMd:
		return registry.ResourceMD, nil
	default:
		return "", fmt.Errorf("unsupported resource %q", kind)
	}
}

// aggregateOutcomeFromResults converts a per-resource result list into
// one coarse summary row. The aggregate command always resolves a single
// agent, but the loop keeps the function robust if a future caller passes
// more than one result.
func aggregateOutcomeFromResults(kind resourceKind, results []common.Result) aggregateResourceOutcome {
	if len(results) == 0 {
		return aggregateResourceOutcome{kind: kind, status: aggregateStatusSkipped, detail: "no result"}
	}
	status := aggregateStatusApplied
	parts := make([]string, 0, len(results))
	for _, result := range results {
		if result.Status == common.StatusError {
			status = aggregateStatusFailed
		} else if status != aggregateStatusFailed && !aggregateResultApplied(result.Status) {
			status = aggregateStatusSkipped
		}
		detail := string(result.Status)
		if trimmed := strings.TrimSpace(result.Detail); trimmed != "" {
			detail = fmt.Sprintf("%s: %s", result.Status, trimmed)
		}
		parts = append(parts, detail)
	}
	return aggregateResourceOutcome{
		kind:   kind,
		status: status,
		detail: strings.Join(parts, "; "),
	}
}

// aggregateResultApplied reports whether a per-resource status represents
// successful application or an already-satisfied binding.
func aggregateResultApplied(status common.Status) bool {
	switch status {
	case common.StatusOK, common.StatusCreated, common.StatusRebuilt, common.StatusRemoved:
		return true
	default:
		return false
	}
}

// renderAggregateOutcomes writes the compact aggregate result table.
func renderAggregateOutcomes(out io.Writer, outcomes []aggregateResourceOutcome) error {
	const (
		columnResource = "RESOURCE"
		columnStatus   = "STATUS"
		columnDetail   = "DETAIL"
	)
	maxResource := len(columnResource)
	maxStatus := len(columnStatus)
	for _, outcome := range outcomes {
		if width := len(string(outcome.kind)); width > maxResource {
			maxResource = width
		}
		if width := len(outcome.status); width > maxStatus {
			maxStatus = width
		}
	}
	if _, err := fmt.Fprintf(out, "%-*s  %-*s  %s\n", maxResource, columnResource, maxStatus, columnStatus, columnDetail); err != nil {
		return fmt.Errorf("write aggregate header: %w", err)
	}
	for _, outcome := range outcomes {
		if _, err := fmt.Fprintf(out, "%-*s  %-*s  %s\n", maxResource, outcome.kind, maxStatus, outcome.status, outcome.detail); err != nil {
			return fmt.Errorf("write aggregate row: %w", err)
		}
	}
	return nil
}

// lookupAgent finds a selectableAgent by name in the universe slice.
// Used by dispatchAgentSetup to resolve role information once name
// validation has succeeded.
func lookupAgent(universe []selectableAgent, name string) (selectableAgent, bool) {
	for _, agent := range universe {
		if agent.Name == name {
			return agent, true
		}
	}
	return selectableAgent{}, false
}
