// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file discovers the project root used by gf agents: the nearest
// directory that contains both .agents/ and AGENTS.md.

package agents

import (
	"io"
	"os"
	"path/filepath"
)

// skipInvalidProjectMessage is printed when the working tree is not an
// agents workspace. The command then returns success and makes no changes.
const skipInvalidProjectMessage = "skip: current project is missing .agents/ or AGENTS.md; gf agents did nothing"

// DiscoverProjectRoot walks upward from the current directory until it
// finds a directory that contains both .agents/ and AGENTS.md.
func DiscoverProjectRoot() (string, bool, error) {
	current, err := os.Getwd()
	if err != nil {
		return "", false, err
	}
	root, found := DiscoverProjectRootFrom(current)
	return root, found, nil
}

// DiscoverProjectRootFrom walks upward from start looking for .agents/ and AGENTS.md.
func DiscoverProjectRootFrom(start string) (string, bool) {
	current := start
	for {
		if isAgentsProjectRoot(current) {
			return current, true
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return "", false
}

// isAgentsProjectRoot reports whether dir contains canonical agent resources.
func isAgentsProjectRoot(dir string) bool {
	agentsInfo, err := os.Stat(filepath.Join(dir, ".agents"))
	if err != nil || !agentsInfo.IsDir() {
		return false
	}
	if _, err = os.Lstat(filepath.Join(dir, "AGENTS.md")); err != nil {
		return false
	}
	return true
}

// skipInvalidProject writes the no-op skip message and returns success.
func skipInvalidProject(out io.Writer) error {
	return writeLine(out, skipInvalidProjectMessage)
}
