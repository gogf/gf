// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package common provides resource-agnostic primitives shared by every
// resource subpackage of the multi-resource agents framework (skills,
// prompts, md, ...).
//
// It owns:
//   - The Status enum and Result type used by every resource.
//   - The SpecLike interface every resource subpackage's AgentSpec implements.
//   - The link state-machine that handles directory and single-file targets
//     with identical semantics (Inspect/ApplyOneLink/ApplyOneUnlink).
//   - Generic selector parsing and target resolution helpers.
//   - Tabular result rendering with optional extra columns.
//   - Interactive selection helpers (terminal detection, numbered grid,
//     yes/no prompts) used in TTY mode.
//   - Filesystem helpers for symlink target comparison and platform-specific
//     error guidance.
//
// All filesystem mutations use Go standard library primitives (os.Symlink,
// os.Readlink, os.Lstat, os.Remove, os.MkdirAll) only. Real directories and
// files are never automatically removed, even with Force=true.
package common
