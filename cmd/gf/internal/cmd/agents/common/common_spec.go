// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file defines the resource-agnostic agent specification interface
// (SpecLike), the shared category enum, and the link kind enum (directory
// or single file). Every resource subpackage defines its own concrete
// AgentSpec struct that implements SpecLike, allowing common engine code
// to operate on any resource type uniformly.

package common

// Category classifies an agent's project skill path. The same values are
// reused across every resource type (skills, prompts, md, ...).
type Category string

const (
	// CategoryNative marks agents whose project path is already the
	// canonical source path. No symlink is required for these agents.
	CategoryNative Category = "native"
	// CategoryLink marks agents whose project path differs from the
	// canonical source path; a relative symlink is created at the project
	// path pointing back to the source path.
	CategoryLink Category = "link"
	// CategoryRootCollision marks agents whose project path is exactly
	// "skills" (or another bare repo-root path). Creating that link would
	// shadow any real directory at that path, so it is skipped by default
	// and only enabled when the caller passes Force=true.
	CategoryRootCollision Category = "rootCollision"
)

// Kind describes whether a resource binding manages a directory symlink or
// a single-file symlink. Both kinds share the same Status state machine,
// only the Lstat-after-success and conflict detail wording differ.
type Kind int

const (
	// KindDir indicates the binding manages a directory symlink (e.g.
	// .claude/skills -> .agents/skills).
	KindDir Kind = iota
	// KindFile indicates the binding manages a single-file symlink (e.g.
	// CLAUDE.md -> AGENTS.md).
	KindFile
)

// SpecLike is the minimal contract every resource subpackage's AgentSpec
// must satisfy. It exposes the fields the common engine needs to inspect,
// link, unlink, render and select agents without knowing which resource
// type the spec belongs to.
type SpecLike interface {
	// SpecName returns the CLI identifier used on the command line (e.g.
	// "claude"). It must be unique within a registry.
	SpecName() string
	// SpecDisplayName returns the human-readable label rendered in status
	// output and interactive selection grids.
	SpecDisplayName() string
	// SpecCategory returns the agent category controlling link/unlink
	// behavior.
	SpecCategory() Category
	// SpecSourcePath returns the canonical source path that managed
	// symlinks must point at. The path is repo-relative and uses forward
	// slashes; the engine converts it to OS-native separators.
	SpecSourcePath() string
	// SpecProjectPath returns the project-relative target path where the
	// symlink should live. The path uses forward slashes.
	SpecProjectPath() string
	// SpecKind reports whether this binding manages a directory or a
	// single file.
	SpecKind() Kind
}
