// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file implements the arrow-key driven interactive helpers.
//
// Three helpers are exposed:
//   - PromptSelection:       multi-select form for picking 0..N agents
//                            from a candidate list (space toggles, enter
//                            confirms, Esc/Ctrl+C cancels) via huh.
//   - PromptSingleSelection: single-select list for picking exactly one
//                            value. Implemented with a custom Bubble Tea
//                            picker (see common_interactive_list.go) so
//                            long lists keep rows above the cursor
//                            visible; huh.Select pins the cursor to the
//                            top of its viewport once the list exceeds
//                            the terminal height.
//   - PromptYesNo:           Yes/No confirmation form with a default via huh.
//
// Design notes:
//   - Interactive helpers require a real terminal. Callers MUST guard
//     each call with IsInteractiveTerminal(*os.File) so the form never
//     runs in CI or piped contexts. When stdin is not a *os.File the
//     helpers conservatively treat the call as cancelled (returning
//     empty/zero values) to keep test harnesses that wire bytes.Buffer
//     streams working without spawning a TUI.
//   - Output is routed through the caller-supplied io.Writer and stdin
//     through WithInput / tea.WithInput so process IO stays
//     consistent with the rest of the tool.
//   - User abort (Esc / Ctrl+C) is translated to "cancelled" rather than
//     an error, matching the historical numbered-grid behavior where
//     blank/q cancelled.

package common

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

// SingleOption describes one choice in a single-select prompt. Value is
// returned to the caller; Label is the human-readable text rendered in
// the selection list.
type SingleOption struct {
	// Value is the underlying identifier returned when this option is
	// chosen. It must be unique within an options slice. Section headers
	// may reuse reserved values; they are never returned to callers
	// because Section rows are not selectable.
	Value string
	// Label is the human-readable text shown to the user. When empty
	// the helper falls back to Value.
	Label string
	// Section marks a non-selectable group header row (for example
	// "Built-in support"). Cursor movement skips these rows.
	Section bool
}

// promptIOReady reports whether the supplied reader/writer pair is
// suitable for driving a huh form. We require stdin to be a *os.File
// pointing at a character device; non-file readers (bytes.Buffer,
// strings.Reader) cannot drive bubbletea's input loop.
func promptIOReady(in io.Reader) (*os.File, bool) {
	file, ok := in.(*os.File)
	if !ok {
		return nil, false
	}
	if !IsInteractiveTerminal(file) {
		return nil, false
	}
	return file, true
}

// abortKeyMap returns a huh KeyMap whose Quit binding accepts both
// `ctrl+c` (huh's default) AND `esc`, so users can dismiss a prompt
// with either key. huh's default keymap only binds `ctrl+c` to Quit;
// `esc` is reserved for sub-features like filter clear and is mostly
// disabled at the form level, which means `esc` does nothing in our
// single-group prompts. Adding `esc` to Quit keeps the documented
// "Esc / Ctrl+C cancels" behavior consistent with reality.
func abortKeyMap() *huh.KeyMap {
	keymap := huh.NewDefaultKeyMap()
	keymap.Quit = key.NewBinding(key.WithKeys("ctrl+c", "esc"), key.WithHelp("esc", "cancel"))
	return keymap
}

// PromptSelection renders a multi-select huh form and returns the agent
// names selected by the user (in candidate order, deduplicated).
//
// Returns (nil, nil) when:
//   - candidates is empty;
//   - stdin is not an interactive terminal (non-TTY callers must not
//     reach this path; reaching it here means the caller forgot to
//     guard with IsInteractiveTerminal — we degrade silently rather
//     than spawn a broken TUI);
//   - the user aborts the form (Esc / Ctrl+C / "no" choice).
//
// Other errors from huh are wrapped and returned to the caller.
func PromptSelection(in io.Reader, out io.Writer, title string, candidates []SelectableEntry) ([]string, error) {
	if len(candidates) == 0 {
		if _, err := fmt.Fprintln(out, title+": no candidates available."); err != nil {
			return nil, fmt.Errorf("write empty-candidates message: %w", err)
		}
		return nil, nil
	}
	stdin, ok := promptIOReady(in)
	if !ok {
		return nil, nil
	}

	options := make([]huh.Option[string], 0, len(candidates))
	for _, entry := range candidates {
		options = append(options, huh.NewOption(formatCandidateLabel(entry), entry.Spec.SpecName()))
	}

	var selected []string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title(title).
				Description("Use arrow keys to navigate, space to toggle, enter to confirm, Esc to cancel.").
				Options(options...).
				Filterable(true).
				Value(&selected),
		),
	).WithInput(stdin).WithOutput(out).WithKeyMap(abortKeyMap())

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil, nil
		}
		return nil, fmt.Errorf("run multi-select form: %w", err)
	}

	// Preserve candidate order (huh returns values in selection order
	// rather than option order). Deduplicate defensively even though
	// huh.MultiSelect already enforces uniqueness across options.
	picked := make(map[string]struct{}, len(selected))
	ordered := make([]string, 0, len(selected))
	for _, entry := range candidates {
		name := entry.Spec.SpecName()
		for _, value := range selected {
			if value != name {
				continue
			}
			if _, dup := picked[name]; dup {
				continue
			}
			picked[name] = struct{}{}
			ordered = append(ordered, name)
			break
		}
	}
	return ordered, nil
}

// PromptSingleSelection renders a single-select list and returns the
// chosen value. Returns ("", nil) when the user aborts or stdin is not
// an interactive terminal. options must be non-empty.
//
// Implementation note: this uses the custom Bubble Tea list in
// common_interactive_list.go rather than huh.Select. huh pins the
// selected option to the top of its viewport once the list is taller
// than the terminal, which hides every row above the cursor while the
// user navigates — unacceptable for the agents aggregate picker with
// 100+ registered tools.
func PromptSingleSelection(in io.Reader, out io.Writer, title string, options []SingleOption) (string, error) {
	if len(options) == 0 {
		return "", fmt.Errorf("PromptSingleSelection: options must not be empty")
	}
	if _, ok := promptIOReady(in); !ok {
		return "", nil
	}
	return runSingleListSelection(in, out, title, options)
}

// PromptYesNo asks a yes/no question with a default answer used when
// the user presses enter without changing the focus. Esc/Ctrl+C aborts
// the form and returns the supplied default to match the historical
// behavior where blank input meant "use the default".
func PromptYesNo(in io.Reader, out io.Writer, question string, defaultYes bool) (bool, error) {
	stdin, ok := promptIOReady(in)
	if !ok {
		return defaultYes, nil
	}

	answer := defaultYes
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(question).
				Affirmative("Yes").
				Negative("No").
				Value(&answer),
		),
	).WithInput(stdin).WithOutput(out).WithKeyMap(abortKeyMap())

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return defaultYes, nil
		}
		return false, fmt.Errorf("run confirm form: %w", err)
	}
	return answer, nil
}

// formatCandidateLabel composes the visible label for a SelectableEntry
// inside huh option lists. The label embeds the status glyph, agent
// name and an optional detail/status descriptor so users can see each
// candidate's current binding state without leaving the prompt.
func formatCandidateLabel(entry SelectableEntry) string {
	var builder strings.Builder
	builder.WriteString(StatusGlyph(entry.CurrentStatus))
	builder.WriteString(" ")
	builder.WriteString(entry.Spec.SpecName())
	suffix := strings.TrimSpace(entry.Detail)
	if suffix == "" {
		suffix = string(entry.CurrentStatus)
	}
	if suffix != "" {
		builder.WriteString("  (")
		builder.WriteString(suffix)
		builder.WriteString(")")
	}
	return builder.String()
}
