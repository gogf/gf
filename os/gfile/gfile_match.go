// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile

import (
	"path"
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
//   - '*'         matches any sequence of non-separator characters
//   - '**'        matches any sequence of characters including separators (globstar)
//   - '?'         matches any single non-separator character
//   - '[abc]'     matches any character in the bracket
//   - '[a-z]'     matches any character in the range
//   - '[^abc]'    matches any character not in the bracket (negation)
//   - '[^a-z]'    matches any character not in the range (negation)
//
// Globstar rules:
//   - "**" only has globstar semantics when it appears as a complete path component
//     (e.g., "a/**/b", "**/a", "a/**", "**").
//   - Patterns like "a**b" or "**a" treat "**" as two regular "*" wildcards,
//     matching only within a single path component.
//   - Both "/" and "\" are treated as path separators (cross-platform support).
//
// Error handling:
//   - Returns an error for malformed patterns (e.g., unclosed brackets "[abc").
//   - Errors from filepath.Match are propagated.
//
// Example:
//
//	MatchGlob("src/**/*.go", "src/foo/bar/main.go")  => true, nil
//	MatchGlob("*.go", "main.go")                     => true, nil
//	MatchGlob("**", "any/path/file.go")              => true, nil
//	MatchGlob("a**b", "axxb")                        => true, nil  (** as two *)
//	MatchGlob("a**b", "a/b")                         => false, nil (no separator match)
//	MatchGlob("[abc]", "a")                          => true, nil
//	MatchGlob("[", "a")                              => false, error (malformed)
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

	// Clean up paths (handles multiple slashes, . and ..)
	// Using path.Clean for consistent cross-platform behavior with forward slashes
	pattern = path.Clean(pattern)
	name = path.Clean(name)

	// Check if "**" appears as a valid globstar (complete path component).
	// If not, treat "**" as two regular "*" wildcards.
	if !hasValidGlobstar(pattern) {
		// Replace "**" with a placeholder, then use filepath.Match
		// Since filepath.Match treats "*" as matching non-separator chars,
		// "**" is equivalent to "*" in terms of matching (both match any
		// sequence of non-separator characters).
		normalizedPattern := strings.ReplaceAll(pattern, "**", "*")
		return filepath.Match(normalizedPattern, name)
	}

	return doMatchGlobstar(pattern, name)
}

// hasValidGlobstar checks if the pattern contains "**" as a valid globstar
// (i.e., as a complete path component). Valid globstar patterns:
//   - "**" (the entire pattern)
//   - "**/" (at the start)
//   - "/**" (at the end)
//   - "/**/" (in the middle)
func hasValidGlobstar(pattern string) bool {
	// Check each occurrence of "**"
	idx := 0
	for {
		pos := strings.Index(pattern[idx:], "**")
		if pos == -1 {
			return false
		}
		pos += idx

		// Check if this "**" is a valid globstar
		if isValidGlobstarAt(pattern, pos) {
			return true
		}

		idx = pos + 2
		if idx >= len(pattern) {
			break
		}
	}
	return false
}

// isValidGlobstarAt checks if the "**" at position pos is a valid globstar.
// A valid globstar must be a complete path component:
//   - At start: "**" or "**/"
//   - At end: "/**"
//   - In middle: "/**/"
func isValidGlobstarAt(pattern string, pos int) bool {
	// Check character before "**"
	if pos > 0 && pattern[pos-1] != '/' {
		return false
	}

	// Check character after "**"
	endPos := pos + 2
	if endPos < len(pattern) && pattern[endPos] != '/' {
		return false
	}

	return true
}

// findValidGlobstar finds the first valid globstar in the pattern.
// Returns the position or -1 if not found.
func findValidGlobstar(pattern string) int {
	idx := 0
	for {
		pos := strings.Index(pattern[idx:], "**")
		if pos == -1 {
			return -1
		}
		pos += idx

		if isValidGlobstarAt(pattern, pos) {
			return pos
		}

		idx = pos + 2
		if idx >= len(pattern) {
			break
		}
	}
	return -1
}

// doMatchGlobstar recursively matches pattern with globstar support.
// Uses memoization to avoid exponential time complexity with multiple "**" operators.
func doMatchGlobstar(pattern, name string) (bool, error) {
	memo := make(map[string]bool)
	return doMatchGlobstarMemo(pattern, name, memo)
}

// doMatchGlobstarMemo is the memoized implementation of globstar matching.
func doMatchGlobstarMemo(pattern, name string, memo map[string]bool) (bool, error) {
	// Create cache key
	cacheKey := pattern + "\x00" + name
	if cached, ok := memo[cacheKey]; ok {
		return cached, nil
	}

	result, err := doMatchGlobstarCore(pattern, name, memo)
	if err != nil {
		return false, err
	}

	memo[cacheKey] = result
	return result, nil
}

// doMatchGlobstarCore contains the core matching logic.
func doMatchGlobstarCore(pattern, name string, memo map[string]bool) (bool, error) {
	// Find the first valid globstar
	pos := findValidGlobstar(pattern)
	if pos == -1 {
		// No valid globstar, use standard match
		// Replace any "**" with "*" since they're not valid globstars
		normalizedPattern := strings.ReplaceAll(pattern, "**", "*")
		return filepath.Match(normalizedPattern, name)
	}

	// Split pattern at the valid globstar position
	prefix := pattern[:pos]
	suffix := pattern[pos+2:]

	// Remove trailing slash from prefix
	prefix = strings.TrimSuffix(prefix, "/")
	// Remove leading slash from suffix
	suffix = strings.TrimPrefix(suffix, "/")

	// Match prefix
	if prefix != "" {
		// Check if name starts with prefix pattern
		if !strings.Contains(prefix, "*") && !strings.Contains(prefix, "?") && !strings.Contains(prefix, "[") {
			// Prefix is literal, check directly against full path component
			if !strings.HasPrefix(name, prefix) {
				return false, nil
			}
			if len(name) == len(prefix) {
				// Name is exactly the prefix
				name = ""
			} else {
				// Ensure the prefix ends at a path separator boundary
				if name[len(prefix)] != '/' {
					return false, nil
				}
				// Skip the separator as well
				name = name[len(prefix)+1:]
			}
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
		return doMatchGlobstarMemo(suffix, "", memo)
	}

	nameParts := strings.Split(name, "/")

	// Try "**" matching 0, 1, 2, ... N segments
	for i := 0; i <= len(nameParts); i++ {
		remaining := strings.Join(nameParts[i:], "/")
		matched, err := doMatchGlobstarMemo(suffix, remaining, memo)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}

	return false, nil
}
