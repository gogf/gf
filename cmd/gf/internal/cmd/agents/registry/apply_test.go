// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file covers resource-scoped ApplyLink / ApplyUnlink / PlanList
// behavior on the unified registry (formerly split across skills /
// prompts / md packages).

package registry

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/agents/common"
)

func newSkillsFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".agents", "skills"), 0o755); err != nil {
		t.Fatalf("create skills source: %v", err)
	}
	return root
}

func newMDFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, MDSourceFile), []byte("# AGENTS.md\n"), 0o644); err != nil {
		t.Fatalf("create AGENTS.md: %v", err)
	}
	return root
}

func newPromptsFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".agents", "prompts", "opsx"), 0o755); err != nil {
		t.Fatalf("create prompts source: %v", err)
	}
	return root
}

func trySymlink(t *testing.T, oldname, newname string) {
	t.Helper()
	if err := os.Symlink(oldname, newname); err != nil {
		if runtime.GOOS == "windows" {
			t.Skipf("skip: symlink unsupported on this Windows host: %v", err)
		}
		t.Fatalf("symlink: %v", err)
	}
}

func TestSkillsApplyLinkCreatesAndIsIdempotent(t *testing.T) {
	root := newSkillsFixture(t)
	results, err := ApplyLink(root, ResourceSkills, LinkRequest{Selectors: []string{"claude-code"}})
	if err != nil {
		t.Fatalf("first apply: %v", err)
	}
	if len(results) != 1 || results[0].Status != common.StatusCreated {
		t.Fatalf("expected created, got %+v", results)
	}
	link := filepath.Join(root, ".claude", "skills")
	info, err := os.Lstat(link)
	if err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("expected symlink at %s", link)
	}
	again, err := ApplyLink(root, ResourceSkills, LinkRequest{Selectors: []string{"claude-code"}})
	if err != nil {
		t.Fatalf("second apply: %v", err)
	}
	if again[0].Status != common.StatusOK {
		t.Fatalf("expected ok, got %s", again[0].Status)
	}
}

func TestSkillsApplyLinkNativeSkipped(t *testing.T) {
	root := newSkillsFixture(t)
	results, err := ApplyLink(root, ResourceSkills, LinkRequest{Selectors: []string{"cursor"}})
	if err != nil {
		t.Fatalf("apply native: %v", err)
	}
	if results[0].Status != common.StatusNative {
		t.Fatalf("expected native, got %s", results[0].Status)
	}
}

func TestSkillsApplyLinkMismatchRequiresForce(t *testing.T) {
	root := newSkillsFixture(t)
	link := filepath.Join(root, ".claude", "skills")
	if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	other := filepath.Join(root, "other-target")
	if err := os.MkdirAll(other, 0o755); err != nil {
		t.Fatalf("mkdir other: %v", err)
	}
	trySymlink(t, other, link)
	results, err := ApplyLink(root, ResourceSkills, LinkRequest{Selectors: []string{"claude-code"}})
	if err != nil {
		t.Fatalf("apply mismatch: %v", err)
	}
	if results[0].Status != common.StatusMismatch {
		t.Fatalf("expected mismatch, got %s", results[0].Status)
	}
	rebuilt, err := ApplyLink(root, ResourceSkills, LinkRequest{Selectors: []string{"claude-code"}, Force: true})
	if err != nil {
		t.Fatalf("rebuild: %v", err)
	}
	if rebuilt[0].Status != common.StatusRebuilt {
		t.Fatalf("expected rebuilt, got %s", rebuilt[0].Status)
	}
}

func TestSkillsApplyLinkRootCollisionRequiresForce(t *testing.T) {
	root := newSkillsFixture(t)
	defaultRun, err := ApplyLink(root, ResourceSkills, LinkRequest{Selectors: []string{"openclaw"}})
	if err != nil {
		t.Fatalf("default: %v", err)
	}
	if defaultRun[0].Status != common.StatusSkippedRootCollision {
		t.Fatalf("expected skipped-root-collision, got %s", defaultRun[0].Status)
	}
	forced, err := ApplyLink(root, ResourceSkills, LinkRequest{Selectors: []string{"openclaw"}, Force: true})
	if err != nil {
		t.Fatalf("force: %v", err)
	}
	if forced[0].Status != common.StatusCreated {
		t.Fatalf("expected created, got %s", forced[0].Status)
	}
}

func TestSkillsApplyUnlinkOnlyManaged(t *testing.T) {
	root := newSkillsFixture(t)
	if _, err := ApplyLink(root, ResourceSkills, LinkRequest{Selectors: []string{"claude-code"}}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	results, err := ApplyUnlink(root, ResourceSkills, UnlinkRequest{Selectors: []string{"claude-code"}})
	if err != nil {
		t.Fatalf("unlink: %v", err)
	}
	if results[0].Status != common.StatusRemoved {
		t.Fatalf("expected removed, got %s", results[0].Status)
	}
}

func TestMDApplyLinkCreatesFileSymlink(t *testing.T) {
	root := newMDFixture(t)
	results, err := ApplyLink(root, ResourceMD, LinkRequest{Selectors: []string{"claude-code"}})
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if results[0].Status != common.StatusCreated {
		t.Fatalf("expected created, got %s", results[0].Status)
	}
	info, err := os.Lstat(filepath.Join(root, "CLAUDE.md"))
	if err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("expected CLAUDE.md symlink")
	}
}

func TestPromptsApplyLinkCreatesDirSymlink(t *testing.T) {
	root := newPromptsFixture(t)
	results, err := ApplyLink(root, ResourcePrompts, LinkRequest{Selectors: []string{"codex"}})
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if results[0].Status != common.StatusCreated {
		t.Fatalf("expected created, got %s", results[0].Status)
	}
	info, err := os.Lstat(filepath.Join(root, ".codex", "prompts"))
	if err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("expected .codex/prompts symlink")
	}
}

func TestPlanListCoversResource(t *testing.T) {
	root := newSkillsFixture(t)
	results := PlanList(root, ResourceSkills)
	if len(results) != len(SkillsSpecs()) {
		t.Fatalf("PlanList size got=%d want=%d", len(results), len(SkillsSpecs()))
	}
}

func TestMDUnlinkPreservesAuthoredFile(t *testing.T) {
	root := newMDFixture(t)
	if err := os.WriteFile(filepath.Join(root, "CLAUDE.md"), []byte("authored\n"), 0o644); err != nil {
		t.Fatalf("write authored CLAUDE.md: %v", err)
	}
	results, err := ApplyUnlink(root, ResourceMD, UnlinkRequest{Selectors: []string{"claude-code"}})
	if err != nil {
		t.Fatalf("unlink: %v", err)
	}
	if results[0].Status != common.StatusSkippedNotManaged {
		t.Fatalf("expected not-managed, got %s (%s)", results[0].Status, results[0].Detail)
	}
	content, err := os.ReadFile(filepath.Join(root, "CLAUDE.md"))
	if err != nil {
		t.Fatalf("read authored file: %v", err)
	}
	if string(content) != "authored\n" {
		t.Fatalf("authored file was changed: %q", content)
	}
}

func TestSkillsForceRebuildsDegradedGitStub(t *testing.T) {
	root := newSkillsFixture(t)
	link := filepath.Join(root, ".claude", "skills")
	if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(link, []byte("../.agents/skills\n"), 0o644); err != nil {
		t.Fatalf("write stub: %v", err)
	}
	withoutForce, err := ApplyLink(root, ResourceSkills, LinkRequest{Selectors: []string{"claude-code"}})
	if err != nil {
		t.Fatalf("apply stub: %v", err)
	}
	if withoutForce[0].Status != common.StatusConflict {
		t.Fatalf("expected conflict on stub without force, got %s", withoutForce[0].Status)
	}
	rebuilt, err := ApplyLink(root, ResourceSkills, LinkRequest{Selectors: []string{"claude-code"}, Force: true})
	if err != nil {
		t.Fatalf("force stub: %v", err)
	}
	if rebuilt[0].Status != common.StatusRebuilt && rebuilt[0].Status != common.StatusCreated {
		t.Fatalf("expected rebuilt/created, got %s (%s)", rebuilt[0].Status, rebuilt[0].Detail)
	}
	info, err := os.Lstat(link)
	if err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("expected symlink after force rebuild")
	}
}

func TestApplyLinkRequiresSelector(t *testing.T) {
	root := newSkillsFixture(t)
	if _, err := ApplyLink(root, ResourceSkills, LinkRequest{}); err == nil {
		t.Fatalf("expected error")
	}
	if _, err := ApplyUnlink(root, ResourceSkills, UnlinkRequest{}); err == nil {
		t.Fatalf("expected error")
	}
}

func TestLinkCandidatesOnlyLinkClass(t *testing.T) {
	root := newSkillsFixture(t)
	for _, entry := range LinkCandidates(root, ResourceSkills) {
		if entry.Spec.SpecCategory() != common.CategoryLink {
			t.Fatalf("candidate %s category=%s", entry.Spec.SpecName(), entry.Spec.SpecCategory())
		}
	}
}
