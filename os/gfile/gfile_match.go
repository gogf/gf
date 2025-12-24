// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile

import (
	"path/filepath"
	"strings"
)

// MatchGlob reports whether name matches the shell pattern.
// It extends filepath.Match (https://pkg.go.dev/path/filepath#Match)
// with support for "**" (globstar) pattern, similar to bash's globstar
// (https://www.gnu.org/software/bash/manual/html_node/The-Shopt-Builtin.html)
// and gitignore patterns (https://git-scm.com/docs/gitignore#_pattern_format).
//
// Pattern syntax:
//   - '*'      matches any sequence of non-separator characters
//   - '**'     matches any sequence of characters including separators (globstar)
//   - '?'      matches any single non-separator character
//   - '[abc]'  matches any character in the bracket
//   - '[a-z]'  matches any character in the range
//
// Example:
//
//	MatchGlob("src/**/*.go", "src/foo/bar/main.go")  => true
//	MatchGlob("*.go", "main.go")                     => true
//	MatchGlob("**", "any/path/file.go")              => true
func MatchGlob(pattern, name string) (bool, error) {
	// If no **, use standard filepath.Match
	if !strings.Contains(pattern, "**") {
		return filepath.Match(pattern, name)
	}
	return matchGlobstar(pattern, name)
}

// matchGlobstar handles patterns containing "**".
func matchGlobstar(pattern, name string) (bool, error) {
	// Normalize path separators to / (handle both Windows and Unix)
	pattern = strings.ReplaceAll(pattern, "\\", "/")
	name = strings.ReplaceAll(name, "\\", "/")

	// Clean up multiple slashes
	for strings.Contains(pattern, "//") {
		pattern = strings.ReplaceAll(pattern, "//", "/")
	}
	for strings.Contains(name, "//") {
		name = strings.ReplaceAll(name, "//", "/")
	}

	return doMatchGlobstar(pattern, name)
}

// doMatchGlobstar recursively matches pattern with globstar support.
func doMatchGlobstar(pattern, name string) (bool, error) {
	// Split pattern by "**"
	parts := strings.SplitN(pattern, "**", 2)
	if len(parts) == 1 {
		// No "**" found, use standard match
		return filepath.Match(pattern, name)
	}

	prefix := parts[0]
	suffix := parts[1]

	// Remove trailing slash from prefix
	prefix = strings.TrimSuffix(prefix, "/")
	// Remove leading slash from suffix
	suffix = strings.TrimPrefix(suffix, "/")

	// Match prefix
	if prefix != "" {
		// Check if name starts with prefix pattern
		if !strings.Contains(prefix, "*") && !strings.Contains(prefix, "?") && !strings.Contains(prefix, "[") {
			// Prefix is literal, check directly
			if !strings.HasPrefix(name, prefix) {
				return false, nil
			}
			name = strings.TrimPrefix(name, prefix)
			name = strings.TrimPrefix(name, "/")
		} else {
			// Prefix contains wildcards, need to match each segment
			prefixParts := strings.Split(prefix, "/")
			nameParts := strings.Split(name, "/")

			if len(nameParts) < len(prefixParts) {
				return false, nil
			}

			for i, pp := range prefixParts {
				matched, err := filepath.Match(pp, nameParts[i])
				if err != nil {
					return false, err
				}
				if !matched {
					return false, nil
				}
			}
			name = strings.Join(nameParts[len(prefixParts):], "/")
		}
	}

	// If suffix is empty, "**" matches everything remaining
	if suffix == "" {
		return true, nil
	}

	// Try matching "**" with 0 to N path segments
	if name == "" {
		// No remaining name, check if suffix can match empty
		return doMatchGlobstar(suffix, "")
	}

	nameParts := strings.Split(name, "/")

	// Try "**" matching 0, 1, 2, ... N segments
	for i := 0; i <= len(nameParts); i++ {
		remaining := strings.Join(nameParts[i:], "/")
		matched, err := doMatchGlobstar(suffix, remaining)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}

	return false, nil
}
