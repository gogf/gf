// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file implements the resource-agnostic link state machine. Inspect,
// ApplyOneLink and ApplyOneUnlink operate on any SpecLike value and produce
// the same Status outcomes regardless of whether the binding manages a
// directory or a single-file symlink. Callers (resource subpackages)
// supply the concrete AgentSpec values via the SpecLike interface and
// receive Result values they can pass directly to Render.
//
// When force is true, the engine additionally detects "degraded symlink
// stubs" — plain files created by Git when core.symlinks=false whose
// content is the relative link target path. These stubs are automatically
// removed and replaced with a real symlink. Real directories and
// non-stub files are still never automatically removed.

package common

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Inspect returns the current Status and Detail for an agent without any
// filesystem mutation. It is used by the default no-selector listing flow
// shared across every resource subpackage.
func Inspect(repoRoot string, spec SpecLike) Result {
	if spec.SpecCategory() == CategoryNative {
		return Result{Spec: spec, Status: StatusNative}
	}
	target := filepath.Join(repoRoot, filepath.FromSlash(spec.SpecProjectPath()))
	info, lstatErr := os.Lstat(target)
	if errors.Is(lstatErr, os.ErrNotExist) {
		if spec.SpecCategory() == CategoryRootCollision {
			return Result{Spec: spec, Status: StatusSkippedRootCollision, Detail: "use force=1 to create"}
		}
		return Result{Spec: spec, Status: StatusAbsent}
	}
	if lstatErr != nil {
		return Result{Spec: spec, Status: StatusError, Detail: lstatErr.Error()}
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return Result{Spec: spec, Status: StatusConflict, Detail: conflictDetail(spec.SpecKind(), false)}
	}
	source := absoluteSource(repoRoot, spec)
	matches, currentTarget, err := LinkMatchesSource(target, source)
	if err != nil {
		return Result{Spec: spec, Status: StatusError, Detail: err.Error()}
	}
	if matches {
		return Result{Spec: spec, Status: StatusOK}
	}
	return Result{Spec: spec, Status: StatusMismatch, Detail: "-> " + currentTarget}
}

// ApplyOneLink performs the link action for a single agent. When force is
// true, degraded symlink stubs (plain files whose content matches the
// expected relative link target, typically created by Git with
// core.symlinks=false) are automatically replaced with real symlinks.
func ApplyOneLink(repoRoot string, spec SpecLike, force bool) Result {
	if spec.SpecCategory() == CategoryNative {
		return Result{Spec: spec, Status: StatusNative}
	}
	if spec.SpecCategory() == CategoryRootCollision && !force {
		return Result{Spec: spec, Status: StatusSkippedRootCollision, Detail: "use force=1 to create"}
	}
	target := filepath.Join(repoRoot, filepath.FromSlash(spec.SpecProjectPath()))
	info, lstatErr := os.Lstat(target)
	switch {
	case errors.Is(lstatErr, os.ErrNotExist):
		return createLink(repoRoot, spec, target)
	case lstatErr != nil:
		return Result{Spec: spec, Status: StatusError, Detail: lstatErr.Error()}
	}
	if info.Mode()&os.ModeSymlink == 0 {
		if force && isDegradedSymlinkStub(repoRoot, spec, target, info) {
			if removeErr := os.Remove(target); removeErr != nil {
				return Result{Spec: spec, Status: StatusError, Detail: "remove degraded stub: " + removeErr.Error()}
			}
			result := createLink(repoRoot, spec, target)
			if result.Status == StatusCreated {
				result.Status = StatusRebuilt
				result.Detail = "replaced degraded git stub; " + result.Detail
			}
			return result
		}
		return Result{Spec: spec, Status: StatusConflict, Detail: conflictDetail(spec.SpecKind(), force)}
	}
	source := absoluteSource(repoRoot, spec)
	matches, currentTarget, err := LinkMatchesSource(target, source)
	if err != nil {
		return Result{Spec: spec, Status: StatusError, Detail: err.Error()}
	}
	if matches {
		return Result{Spec: spec, Status: StatusOK}
	}
	if !force {
		return Result{Spec: spec, Status: StatusMismatch, Detail: "-> " + currentTarget + "; use force=1 to rebuild"}
	}
	if removeErr := os.Remove(target); removeErr != nil {
		return Result{Spec: spec, Status: StatusError, Detail: "remove existing link: " + removeErr.Error()}
	}
	result := createLink(repoRoot, spec, target)
	if result.Status == StatusCreated {
		result.Status = StatusRebuilt
		result.Detail = "previous: -> " + currentTarget
	}
	return result
}

// ApplyOneUnlink performs the unlink action for a single agent. It only
// removes symlinks pointing at the canonical source path; real files and
// directories, and symlinks pointing elsewhere, are preserved.
func ApplyOneUnlink(repoRoot string, spec SpecLike) Result {
	target := filepath.Join(repoRoot, filepath.FromSlash(spec.SpecProjectPath()))
	info, lstatErr := os.Lstat(target)
	if errors.Is(lstatErr, os.ErrNotExist) {
		return Result{Spec: spec, Status: StatusAbsent}
	}
	if lstatErr != nil {
		return Result{Spec: spec, Status: StatusError, Detail: lstatErr.Error()}
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return Result{Spec: spec, Status: StatusSkippedNotManaged, Detail: notManagedDetail(spec.SpecKind())}
	}
	source := absoluteSource(repoRoot, spec)
	matches, currentTarget, err := LinkMatchesSource(target, source)
	if err != nil {
		return Result{Spec: spec, Status: StatusError, Detail: err.Error()}
	}
	if !matches {
		return Result{Spec: spec, Status: StatusSkippedForeignTarget, Detail: "-> " + currentTarget}
	}
	if removeErr := os.Remove(target); removeErr != nil {
		return Result{Spec: spec, Status: StatusError, Detail: removeErr.Error()}
	}
	return Result{Spec: spec, Status: StatusRemoved}
}

// createLink resolves the relative source path and creates a symlink at
// the project path. Parent directories are created on demand for both
// directory-kind and file-kind bindings.
func createLink(repoRoot string, spec SpecLike, target string) Result {
	if mkErr := os.MkdirAll(filepath.Dir(target), 0o755); mkErr != nil {
		return Result{Spec: spec, Status: StatusError, Detail: "create parent directory: " + mkErr.Error()}
	}
	source := absoluteSource(repoRoot, spec)
	relativeSource, relErr := filepath.Rel(filepath.Dir(target), source)
	if relErr != nil {
		return Result{Spec: spec, Status: StatusError, Detail: "compute relative source: " + relErr.Error()}
	}
	if symErr := os.Symlink(relativeSource, target); symErr != nil {
		return Result{Spec: spec, Status: StatusError, Detail: SymlinkErrorDetail(symErr)}
	}
	return Result{Spec: spec, Status: StatusCreated, Detail: "-> " + filepath.ToSlash(relativeSource)}
}

// absoluteSource returns the OS-native absolute path of the canonical
// source for a spec.
func absoluteSource(repoRoot string, spec SpecLike) string {
	return filepath.Join(repoRoot, filepath.FromSlash(spec.SpecSourcePath()))
}

// conflictDetail returns the detail message used when a non-symlink path
// blocks linking. Directory- and file-kind bindings phrase the conflict
// slightly differently so users can spot which kind of object survived.
// When force is already true, it omits the "use force" hint since the
// user has already passed the flag.
func conflictDetail(kind Kind, force bool) string {
	suffix := "; use force=1 to auto-replace git stubs"
	if force {
		suffix = "; not a degraded git stub, resolve manually"
	}
	if kind == KindFile {
		return "real file exists" + suffix
	}
	return "real path exists" + suffix
}

// isDegradedSymlinkStub detects whether the non-symlink at target is a
// "degraded symlink stub" that Git creates when core.symlinks=false. Git
// writes the intended link target path as the file content. This function
// checks:
//  1. The path is a regular file (not a directory).
//  2. The file size is small (≤512 bytes, symlink targets are short paths).
//  3. The trimmed file content matches the expected relative source path.
//
// When all conditions hold, it is safe for force mode to replace the stub
// with a real symlink.
func isDegradedSymlinkStub(repoRoot string, spec SpecLike, target string, info os.FileInfo) bool {
	// Only regular files can be degraded stubs; directories are never stubs.
	if info.IsDir() {
		return false
	}
	// Git stub files are tiny (just a relative path string). Reject large files.
	if info.Size() > 512 {
		return false
	}
	content, err := os.ReadFile(target)
	if err != nil {
		return false
	}
	trimmed := strings.TrimSpace(string(content))
	if trimmed == "" {
		return false
	}
	// Compute the expected relative source path as Git would write it.
	source := absoluteSource(repoRoot, spec)
	expectedRel, relErr := filepath.Rel(filepath.Dir(target), source)
	if relErr != nil {
		return false
	}
	// Compare using forward slashes (Git normalizes to forward slashes).
	expectedSlash := filepath.ToSlash(expectedRel)
	trimmedSlash := filepath.ToSlash(trimmed)
	return trimmedSlash == expectedSlash
}

// notManagedDetail returns the detail message used by ApplyOneUnlink when
// the target is a real directory or file. The wording mirrors
// conflictDetail so the table column reads consistently across resources.
func notManagedDetail(kind Kind) string {
	if kind == KindFile {
		return "real file; not removed"
	}
	return "real path; not removed"
}
