// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file defines the Status enum and Result type shared by every
// resource subpackage. Statuses describe the outcome of one agent action
// emitted by Inspect/ApplyOneLink/ApplyOneUnlink and rendered by the
// command output table.

package common

// Status is the outcome of one agent action emitted by the command output.
type Status string

const (
	// StatusNative indicates the agent project path is already the canonical
	// source path (no symlink needed).
	StatusNative Status = "native"
	// StatusOK indicates the link already exists and points at the canonical
	// source path.
	StatusOK Status = "ok"
	// StatusCreated indicates a new link was just created.
	StatusCreated Status = "created"
	// StatusRebuilt indicates a mismatched link was removed and recreated.
	StatusRebuilt Status = "rebuilt"
	// StatusMismatch indicates a link exists but points at a different
	// target than the canonical source path, and force was not provided.
	StatusMismatch Status = "mismatch"
	// StatusConflict indicates the project path exists as a real directory
	// or file. Auto-resolution is never attempted.
	StatusConflict Status = "conflict"
	// StatusSkippedRootCollision indicates a rootCollision agent was
	// skipped because force was not provided.
	StatusSkippedRootCollision Status = "skipped-root-collision"
	// StatusRemoved indicates an unlink call removed a managed link.
	StatusRemoved Status = "removed"
	// StatusSkippedForeignTarget indicates an unlink target is a symlink
	// pointing at a non-managed location and was preserved.
	StatusSkippedForeignTarget Status = "skipped-foreign"
	// StatusSkippedNotManaged indicates an unlink target is a real
	// directory or file and was preserved.
	StatusSkippedNotManaged Status = "skipped-not-managed"
	// StatusAbsent indicates the unlink target does not exist.
	StatusAbsent Status = "absent"
	// StatusError indicates an unrecoverable error processing the agent.
	// The caller renders Detail to show the underlying error message.
	StatusError Status = "error"
)

// Result describes the outcome of one agent action across every resource
// type. The Spec field is held as a SpecLike interface so the same Result
// works for skills, prompts and md resources without copy-pasting struct
// definitions.
type Result struct {
	// Spec is the agent under inspection (any concrete AgentSpec value
	// implementing SpecLike).
	Spec SpecLike
	// Status is the action outcome.
	Status Status
	// Detail provides additional context for non-trivial statuses such as
	// the actual link target during a mismatch or the underlying error
	// message when Status is StatusError.
	Detail string
}

// IsError reports whether the result represents an unrecoverable error.
func (r Result) IsError() bool {
	return r.Status == StatusError
}

// HasError reports whether any result in the slice represents an error.
func HasError(results []Result) bool {
	for _, result := range results {
		if result.IsError() {
			return true
		}
	}
	return false
}
