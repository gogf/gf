// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file implements the resource-agnostic terminal primitives shared
// by every interactive helper in the common subpackage.
//
// It owns:
//   - IsInteractiveTerminal: detects whether a stdin file points at a
//     real character device, used by every command_agents.* entry point
//     to decide whether to enter the huh-based interactive flow or fall
//     back to a non-interactive usage hint.
//   - SelectableEntry: the row description type consumed by huh-based
//     PromptSelection / PromptSingleSelection helpers in
//     common_interactive_huh.go.
//   - StatusGlyph: the single-character status indicator embedded into
//     huh option labels and the resource Render output.
//
// All terminal interaction (arrow-key navigation, single/multi select,
// yes/no confirmation) lives in common_interactive_huh.go and goes
// through charmbracelet/huh. The legacy numbered-grid prompts and
// renderCandidateGrid helper that previously lived here have been
// removed: they are no longer reachable from any command path.

package common

import (
	"os"
)

// IsInteractiveTerminal reports whether the provided file looks like an
// interactive terminal. Returning false also covers nil files, pipes and
// regular files used by tests.
func IsInteractiveTerminal(file *os.File) bool {
	if file == nil {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

// SelectableEntry describes one candidate row for interactive selection.
// The Spec field is held as SpecLike so the helper works for any resource
// subpackage's AgentSpec type.
type SelectableEntry struct {
	// Spec is the agent backing this row.
	Spec SpecLike
	// CurrentStatus is the current Inspect status, used to render context.
	CurrentStatus Status
	// Detail mirrors Inspect's detail field for the same row.
	Detail string
}

// StatusGlyph maps a Status to a single-character indicator embedded in
// huh option labels and in resource Render output. Unknown statuses fall
// back to "?" so the UI never silently drops a state.
//
// Glyphs:
//
//	[+] linked         ok / created / rebuilt / removed
//	[~] mismatch       linked but pointing at a foreign target
//	[.] absent         not linked yet (or native, no action)
//	[!] conflict       a real directory or file blocks linking
//	[*] root collision agent uses a colliding repo-root path (openclaw)
//	[?] error          inspection failed; see status table for details
func StatusGlyph(status Status) string {
	switch status {
	case StatusOK, StatusCreated, StatusRebuilt, StatusRemoved:
		return "[+]"
	case StatusMismatch:
		return "[~]"
	case StatusAbsent, StatusNative:
		return "[.]"
	case StatusConflict:
		return "[!]"
	case StatusSkippedRootCollision:
		return "[*]"
	case StatusSkippedForeignTarget, StatusSkippedNotManaged:
		return "[~]"
	case StatusError:
		return "[?]"
	default:
		return "[?]"
	}
}
