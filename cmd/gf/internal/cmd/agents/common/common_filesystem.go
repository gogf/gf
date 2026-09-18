// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file implements the resource-agnostic filesystem helpers used by the
// link state machine: symlink target comparison, platform-aware path
// equality and Windows-specific symlink error formatting. The helpers use
// only Go standard library primitives (os.Readlink, filepath.Abs,
// filepath.Clean, filepath.IsAbs) so they behave consistently across
// Linux, macOS and Windows.

package common

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// LinkMatchesSource checks whether the symlink at link points at the
// supplied canonical source path. It accepts both absolute and relative
// link target values and returns the resolved current target string
// alongside the comparison result so callers can render diagnostic detail.
//
// sourceAbs must be the absolute, cleaned, OS-native path the link is
// expected to point at. Callers typically derive it via
// filepath.Join(repoRoot, filepath.FromSlash(spec.SpecSourcePath())).
func LinkMatchesSource(link string, sourceAbs string) (bool, string, error) {
	currentTarget, err := os.Readlink(link)
	if err != nil {
		return false, "", fmt.Errorf("readlink %s: %w", link, err)
	}
	resolved := currentTarget
	if !filepath.IsAbs(resolved) {
		resolved = filepath.Join(filepath.Dir(link), resolved)
	}
	resolvedClean, err := filepath.Abs(filepath.Clean(resolved))
	if err != nil {
		return false, currentTarget, fmt.Errorf("resolve link target: %w", err)
	}
	sourceClean, err := filepath.Abs(filepath.Clean(sourceAbs))
	if err != nil {
		return false, currentTarget, fmt.Errorf("resolve source path: %w", err)
	}
	return pathsEqual(resolvedClean, sourceClean), filepath.ToSlash(currentTarget), nil
}

// pathsEqual compares two cleaned paths using the platform's case-folding
// rules. On Windows path comparisons must be case-insensitive.
func pathsEqual(left string, right string) bool {
	if runtime.GOOS == "windows" {
		return equalFoldPath(left, right)
	}
	return left == right
}

// equalFoldPath performs an ASCII case-insensitive path comparison suitable
// for Windows file systems where path components compare case-insensitively.
func equalFoldPath(left string, right string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := 0; index < len(left); index++ {
		leftByte := left[index]
		rightByte := right[index]
		if leftByte >= 'A' && leftByte <= 'Z' {
			leftByte += 'a' - 'A'
		}
		if rightByte >= 'A' && rightByte <= 'Z' {
			rightByte += 'a' - 'A'
		}
		if leftByte != rightByte {
			return false
		}
	}
	return true
}

// SymlinkErrorDetail formats an os.Symlink error with platform-specific
// guidance when permission is denied on Windows. The returned string is
// safe to surface in command output and CI logs.
func SymlinkErrorDetail(err error) string {
	if err == nil {
		return ""
	}
	message := err.Error()
	if runtime.GOOS == "windows" {
		return message + "; Windows requires Developer Mode or Administrator to create symlinks"
	}
	return message
}
