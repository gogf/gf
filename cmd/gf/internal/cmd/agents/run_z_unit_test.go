// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package agents

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Run_ClaudeLinkAndUnlink(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		root := t.TempDir()
		t.Assert(os.MkdirAll(filepath.Join(root, ".agents", "skills"), 0o755), nil)
		t.Assert(os.MkdirAll(filepath.Join(root, ".agents", "prompts"), 0o755), nil)
		t.Assert(os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("# AGENTS.md\n"), 0o644), nil)

		stdout := &bytes.Buffer{}
		err := Run(Request{
			Root:   root,
			Agent:  "claude",
			Stdin:  strings.NewReader(""),
			Stdout: stdout,
		})
		t.AssertNil(err)
		info, statErr := os.Lstat(filepath.Join(root, "CLAUDE.md"))
		t.AssertNil(statErr)
		t.Assert(info.Mode()&os.ModeSymlink != 0, true)
		info, statErr = os.Lstat(filepath.Join(root, ".claude", "skills"))
		t.AssertNil(statErr)
		t.Assert(info.Mode()&os.ModeSymlink != 0, true)

		err = Run(Request{
			Root:   root,
			Agent:  "claude",
			Action: "unlink",
			Stdin:  strings.NewReader(""),
			Stdout: &bytes.Buffer{},
		})
		t.AssertNil(err)
		_, statErr = os.Lstat(filepath.Join(root, "CLAUDE.md"))
		t.Assert(os.IsNotExist(statErr), true)
		_, statErr = os.Lstat(filepath.Join(root, "AGENTS.md"))
		t.AssertNil(statErr)
	})
}

func Test_Run_RejectsAgentAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		root := t.TempDir()
		t.Assert(os.MkdirAll(filepath.Join(root, ".agents"), 0o755), nil)
		t.Assert(os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("# AGENTS.md\n"), 0o644), nil)
		err := Run(Request{
			Root:   root,
			Agent:  "all",
			Stdin:  strings.NewReader(""),
			Stdout: &bytes.Buffer{},
		})
		t.AssertNE(err, nil)
		t.Assert(strings.Contains(err.Error(), "all"), true)
	})
}

func Test_Run_SkipsInvalidProjectLayout(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		root := t.TempDir()
		stdout := &bytes.Buffer{}
		err := Run(Request{
			Root:   root,
			Agent:  "claude",
			Stdin:  strings.NewReader(""),
			Stdout: stdout,
		})
		t.AssertNil(err)
		t.Assert(strings.Contains(stdout.String(), "skip:"), true)
		t.Assert(strings.Contains(stdout.String(), ".agents/"), true)
		_, statErr := os.Lstat(filepath.Join(root, "CLAUDE.md"))
		t.Assert(os.IsNotExist(statErr), true)
		_, statErr = os.Lstat(filepath.Join(root, ".claude"))
		t.Assert(os.IsNotExist(statErr), true)
	})
}

func Test_Run_SkipsWhenWorkingDirectoryHasNoLayout(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		dir := t.TempDir()
		cwd, err := os.Getwd()
		t.AssertNil(err)
		chdirErr := os.Chdir(dir)
		t.AssertNil(chdirErr)
		t.Cleanup(func() {
			if restoreErr := os.Chdir(cwd); restoreErr != nil {
				t.Errorf("restore cwd: %v", restoreErr)
			}
		})
		stdout := &bytes.Buffer{}
		err = Run(Request{
			Agent:  "claude",
			Stdin:  strings.NewReader(""),
			Stdout: stdout,
		})
		t.AssertNil(err)
		t.Assert(strings.Contains(stdout.String(), "skip:"), true)
		_, statErr := os.Lstat(filepath.Join(dir, "CLAUDE.md"))
		t.Assert(os.IsNotExist(statErr), true)
	})
}
