// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file contains unit tests for the resource-agnostic interactive
// primitives that survived the migration to charmbracelet/huh.
//
// We test only the parts that are testable without a real terminal:
//   - IsInteractiveTerminal.
//   - The non-TTY degrade behaviour of PromptSelection / PromptSingleSelection
//     / PromptYesNo. When stdin is not a *os.File pointing at a character
//     device the helpers must return safe zero values rather than spawn a
//     huh form. This protects test harnesses that wire bytes.Buffer streams
//     and CI invocations that should never reach the interactive path.
//
// The actual huh-driven TUI rendering and key handling is exercised by
// charmbracelet/huh's own test suite and verified manually; we do not
// attempt to assert on its output bytes here because they include
// terminal control sequences whose layout is intentionally not part of
// our public contract.

package common

import (
	"bytes"
	"strings"
	"testing"
)

// fakeSelectables produces a tiny SelectableEntry slice for tests in
// this file. It mirrors makeFakeRegistry/fakeSpec from common_test.go
// but composes the SelectableEntry rows directly.
func fakeSelectables() []SelectableEntry {
	return []SelectableEntry{
		{Spec: fakeSpec{name: "claude-code", category: CategoryLink}, CurrentStatus: StatusAbsent},
		{Spec: fakeSpec{name: "codebuddy", category: CategoryLink}, CurrentStatus: StatusOK},
	}
}

func TestIsInteractiveTerminalNilFile(t *testing.T) {
	if IsInteractiveTerminal(nil) {
		t.Fatalf("nil file must not be a terminal")
	}
}

// TestPromptSelectionNonTTYReturnsEmpty verifies that PromptSelection
// gracefully degrades to (nil, nil) when stdin is not a *os.File. This
// keeps mock-based command tests from accidentally spawning a huh form.
func TestPromptSelectionNonTTYReturnsEmpty(t *testing.T) {
	in := bytes.NewBufferString("")
	out := &bytes.Buffer{}
	got, err := PromptSelection(in, out, "Select:", fakeSelectables())
	if err != nil {
		t.Fatalf("PromptSelection: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty selection on non-TTY, got %v", got)
	}
}

// TestPromptSelectionEmptyCandidates verifies the explicit empty-list
// message that runs even on non-TTY callers, since it does not require
// terminal IO.
func TestPromptSelectionEmptyCandidates(t *testing.T) {
	out := &bytes.Buffer{}
	got, err := PromptSelection(bytes.NewBufferString(""), out, "Select:", nil)
	if err != nil {
		t.Fatalf("PromptSelection: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty selection")
	}
	if !strings.Contains(out.String(), "no candidates") {
		t.Fatalf("expected no-candidates message: %q", out.String())
	}
}

// TestPromptSingleSelectionNonTTYReturnsEmpty mirrors the multi-select
// degrade path for the single-select helper.
func TestPromptSingleSelectionNonTTYReturnsEmpty(t *testing.T) {
	got, err := PromptSingleSelection(bytes.NewBufferString(""), &bytes.Buffer{}, "Pick:", []SingleOption{
		{Value: "a", Label: "A"},
		{Value: "b", Label: "B"},
	})
	if err != nil {
		t.Fatalf("PromptSingleSelection: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty value on non-TTY, got %q", got)
	}
}

// TestPromptSingleSelectionEmptyOptionsErrors guards against passing an
// empty options slice, which is a programmer error.
func TestPromptSingleSelectionEmptyOptionsErrors(t *testing.T) {
	if _, err := PromptSingleSelection(bytes.NewBufferString(""), &bytes.Buffer{}, "Pick:", nil); err == nil {
		t.Fatalf("expected error for empty options")
	}
}

// TestPromptYesNoNonTTYReturnsDefault verifies the helper returns the
// caller-supplied default when stdin is not a TTY, instead of spawning
// a form or erroring.
func TestPromptYesNoNonTTYReturnsDefault(t *testing.T) {
	got, err := PromptYesNo(bytes.NewBufferString(""), &bytes.Buffer{}, "OK?", true)
	if err != nil || !got {
		t.Fatalf("default yes path: got=%v err=%v", got, err)
	}
	got, err = PromptYesNo(bytes.NewBufferString(""), &bytes.Buffer{}, "OK?", false)
	if err != nil || got {
		t.Fatalf("default no path: got=%v err=%v", got, err)
	}
}

// TestFormatCandidateLabelEmbedsGlyphAndStatus exercises the label
// builder used to compose huh option text. We verify glyph + name +
// status descriptor are all present.
func TestFormatCandidateLabelEmbedsGlyphAndStatus(t *testing.T) {
	entry := SelectableEntry{
		Spec:          fakeSpec{name: "claude-code", category: CategoryLink},
		CurrentStatus: StatusMismatch,
		Detail:        "-> /elsewhere",
	}
	got := formatCandidateLabel(entry)
	for _, want := range []string{"[~]", "claude-code", "-> /elsewhere"} {
		if !strings.Contains(got, want) {
			t.Fatalf("formatCandidateLabel missing %q in %q", want, got)
		}
	}
}

// TestClampListWindowKeepsCursorVisibleWithContext verifies the sliding
// window does not pin the cursor to the top (the huh.Select failure mode).
func TestClampListWindowKeepsCursorVisibleWithContext(t *testing.T) {
	// Cursor moves from 0 -> 1 inside a tall enough window: start stays 0
	// so the previous row remains visible.
	if got := clampListWindow(1, 0, 10, 50); got != 0 {
		t.Fatalf("expected windowStart to stay 0, got %d", got)
	}
	// Cursor walks past the bottom edge: window scrolls just enough.
	if got := clampListWindow(10, 0, 10, 50); got != 1 {
		t.Fatalf("expected windowStart 1 when cursor leaves bottom, got %d", got)
	}
	// Cursor walks above the top edge: window follows upward.
	if got := clampListWindow(2, 5, 10, 50); got != 2 {
		t.Fatalf("expected windowStart 2 when cursor leaves top, got %d", got)
	}
}

// TestMoveSelectableFilteredSkipsSections ensures group headers are never
// the resting place of the cursor.
func TestMoveSelectableFilteredSkipsSections(t *testing.T) {
	items := []singleListItem{
		{value: "sec-a", label: "A", section: true},
		{value: "a1", label: "A1"},
		{value: "a2", label: "A2"},
		{value: "sec-b", label: "B", section: true},
		{value: "b1", label: "B1"},
	}
	filtered := []int{0, 1, 2, 3, 4}
	// From first selectable (index 1), move down once -> a2 (index 2).
	if got := moveSelectableFiltered(items, filtered, 1, 1); got != 2 {
		t.Fatalf("move down from a1: got %d want 2", got)
	}
	// From a2, move down once -> skip section B, land on b1 (index 4).
	if got := moveSelectableFiltered(items, filtered, 2, 1); got != 4 {
		t.Fatalf("move down across section: got %d want 4", got)
	}
	// From b1, move up once -> a2 (index 2), not the section.
	if got := moveSelectableFiltered(items, filtered, 4, -1); got != 2 {
		t.Fatalf("move up across section: got %d want 2", got)
	}
}

// TestFirstSelectableFilteredSkipsLeadingSection covers initial cursor
// placement on a two-group agent list.
func TestFirstSelectableFilteredSkipsLeadingSection(t *testing.T) {
	items := []singleListItem{
		{value: "sec", label: "Built-in", section: true},
		{value: "amp", label: "Amp"},
	}
	filtered := []int{0, 1}
	if got := firstSelectableFiltered(items, filtered); got != 1 {
		t.Fatalf("first selectable: got %d want 1", got)
	}
}

// TestSectionStartForCursorFindsNearestHeader covers walking upward from
// a selectable row to its group title.
func TestSectionStartForCursorFindsNearestHeader(t *testing.T) {
	items := []singleListItem{
		{value: "sec-a", label: "Built-in", section: true},
		{value: "a1", label: "A1"},
		{value: "a2", label: "A2"},
		{value: "sec-b", label: "Needs setup", section: true},
		{value: "b1", label: "B1"},
	}
	filtered := []int{0, 1, 2, 3, 4}
	if got := sectionStartForCursor(items, filtered, 1); got != 0 {
		t.Fatalf("cursor on a1: sectionStart got=%d want=0", got)
	}
	if got := sectionStartForCursor(items, filtered, 2); got != 0 {
		t.Fatalf("cursor on a2: sectionStart got=%d want=0", got)
	}
	if got := sectionStartForCursor(items, filtered, 4); got != 3 {
		t.Fatalf("cursor on b1: sectionStart got=%d want=3", got)
	}
	if got := sectionStartForCursor(items, filtered, -1); got != -1 {
		t.Fatalf("invalid cursor: sectionStart got=%d want=-1", got)
	}
}

// TestClampListWindowWithSectionRestoresGroupHeader is the regression for
// "return to top of Built-in and still see the section title" instead of
// only "… 1 more above".
func TestClampListWindowWithSectionRestoresGroupHeader(t *testing.T) {
	// Layout: [0]=section, [1]=first agent. After scrolling, windowStart
	// was 1 so the title was hidden. Cursor returns to first agent (1).
	// Plain clamp leaves start=1; with section recovery start becomes 0.
	if got := clampListWindowWithSection(1, 1, 10, 50, 0); got != 0 {
		t.Fatalf("restore built-in header: got=%d want=0", got)
	}
	// Section already in view: do not jump the window.
	if got := clampListWindowWithSection(5, 0, 10, 50, 0); got != 0 {
		t.Fatalf("section already visible: got=%d want=0", got)
	}
	// Deep in a long group: section is far above and cannot fit with the
	// cursor in the same window — keep following the cursor.
	// cursor=25, section=0, windowSize=10 → 25-0 >= 10, no pull.
	if got := clampListWindowWithSection(25, 20, 10, 50, 0); got != 20 {
		// clampListWindow(25, 20, 10, 50): cursor in [20,30) → 20
		t.Fatalf("deep in group keep scroll: got=%d want=20", got)
	}
	// Cursor walks above current window onto first item; section pulls in.
	// clamp alone would set start=1; with section 0 and fit, start=0.
	if got := clampListWindowWithSection(1, 15, 10, 50, 0); got != 0 {
		t.Fatalf("scroll up to first item with header: got=%d want=0", got)
	}
	// Needs-setup group: section at 10, first item at 11, window scrolled
	// so start=11. Returning to item 11 should reveal section 10.
	if got := clampListWindowWithSection(11, 11, 10, 50, 10); got != 10 {
		t.Fatalf("restore needs-setup header: got=%d want=10", got)
	}
}

// TestEnsureVisibleRestoresBuiltInHeader exercises the model helper end
// to end with the two-group agent list layout.
func TestEnsureVisibleRestoresBuiltInHeader(t *testing.T) {
	model := singleListModel{
		items: []singleListItem{
			{value: "sec-a", label: "── Built-in support ──", section: true},
			{value: "amp", label: "Amp"},
			{value: "codex", label: "Codex"},
			{value: "sec-b", label: "── Needs setup ──", section: true},
			{value: "claude", label: "Claude Code"},
		},
		filtered:    []int{0, 1, 2, 3, 4},
		cursor:      1,
		windowStart: 1, // title already scrolled away
		windowSize:  10,
	}
	model.ensureVisible()
	if model.windowStart != 0 {
		t.Fatalf("ensureVisible should reveal built-in header: windowStart=%d want=0", model.windowStart)
	}
}
