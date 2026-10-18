// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package agents

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverProjectRootFromFindsAgentsAndMarkdown(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "internal", "cmd")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".agents"), 0o755); err != nil {
		t.Fatalf("mkdir .agents: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("# rules\n"), 0o644); err != nil {
		t.Fatalf("write AGENTS.md: %v", err)
	}
	got, found := DiscoverProjectRootFrom(nested)
	if !found {
		t.Fatalf("expected to find project root")
	}
	if filepath.Clean(got) != filepath.Clean(root) {
		t.Fatalf("root got=%q want=%q", got, root)
	}
}

func TestDiscoverProjectRootFromRequiresBothResources(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".agents"), 0o755); err != nil {
		t.Fatalf("mkdir .agents: %v", err)
	}
	if _, found := DiscoverProjectRootFrom(root); found {
		t.Fatalf("expected not found when AGENTS.md is missing")
	}
}

func TestDiscoverProjectRootUsesWorkingDirectory(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "app")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".agents"), 0o755); err != nil {
		t.Fatalf("mkdir .agents: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("# rules\n"), 0o644); err != nil {
		t.Fatalf("write AGENTS.md: %v", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err = os.Chdir(nested); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		if restoreErr := os.Chdir(cwd); restoreErr != nil {
			t.Errorf("restore cwd: %v", restoreErr)
		}
	})
	got, found, err := DiscoverProjectRoot()
	if err != nil {
		t.Fatalf("DiscoverProjectRoot: %v", err)
	}
	if !found {
		t.Fatalf("expected to find project root")
	}
	want, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatalf("eval root: %v", err)
	}
	gotResolved, err := filepath.EvalSymlinks(got)
	if err != nil {
		t.Fatalf("eval got: %v", err)
	}
	if gotResolved != want {
		t.Fatalf("root got=%q want=%q", gotResolved, want)
	}
}
