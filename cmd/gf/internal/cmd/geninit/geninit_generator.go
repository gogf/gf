// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package geninit

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

// generateProject copies the template to the destination and performs cleanup
// oldModule: original module path from template
// newModule: target module path for go.mod (can be different from project name)
func generateProject(ctx context.Context, srcPath, name, oldModule, newModule string) error {
	pwd := gfile.Pwd()

	dstPath := filepath.Join(pwd, name)
	if name == "." {
		dstPath = pwd
	}

	if gfile.Exists(dstPath) && !gfile.IsEmpty(dstPath) {
		return fmt.Errorf("target directory %s is not empty", dstPath)
	}

	mlog.Printf("Generating project in %s...", dstPath)

	// 1. Copy files
	if err := gfile.Copy(srcPath, dstPath); err != nil {
		return err
	}

	// 2. Clean up .git directory
	gitDir := filepath.Join(dstPath, ".git")
	if gfile.Exists(gitDir) {
		if err := gfile.Remove(gitDir); err != nil {
			mlog.Debugf("Failed to remove .git directory: %v", err)
		}
	}

	// 3. Clean up go.work and go.work.sum (workspace files should not be in generated project)
	for _, workFile := range []string{"go.work", "go.work.sum"} {
		workPath := filepath.Join(dstPath, workFile)
		if gfile.Exists(workPath) {
			if err := gfile.Remove(workPath); err != nil {
				mlog.Printf("Failed to remove %s: %v", workFile, err)
			} else {
				mlog.Debugf("Removed %s", workFile)
			}
		}
	}

	// 4. Update go.mod module name
	goModPath := filepath.Join(dstPath, "go.mod")
	if gfile.Exists(goModPath) {
		content := gfile.GetContents(goModPath)
		lines := gstr.Split(content, "\n")
		if len(lines) > 0 && gstr.HasPrefix(lines[0], "module ") {
			lines[0] = "module " + newModule
			newContent := gstr.Join(lines, "\n")
			if err := gfile.PutContents(goModPath, newContent); err != nil {
				mlog.Printf("Failed to update go.mod: %v", err)
			}
		}
	}

	// 5. Use AST to replace import paths in all Go files
	if oldModule != "" && oldModule != newModule {
		replacer := NewASTReplacer(oldModule, newModule)
		if err := replacer.ReplaceInDir(ctx, dstPath); err != nil {
			return fmt.Errorf("failed to replace imports: %w", err)
		}
	}

	mlog.Print("Project generated successfully!")
	return nil
}

// tidyDependencies runs go mod tidy in the project directory
func tidyDependencies(ctx context.Context, projectDir string) error {
	mlog.Print("Tidying dependencies (go mod tidy)...")
	if err := runCmd(ctx, projectDir, "go", "mod", "tidy"); err != nil {
		return fmt.Errorf("go mod tidy failed: %w", err)
	}
	mlog.Print("Dependencies tidied successfully!")
	return nil
}

// upgradeDependencies runs go get -u ./... to upgrade all dependencies to latest
func upgradeDependencies(ctx context.Context, projectDir string) error {
	mlog.Print("Upgrading dependencies to latest (go get -u ./...)...")
	if err := runCmd(ctx, projectDir, "go", "get", "-u", "./..."); err != nil {
		return fmt.Errorf("go get -u failed: %w", err)
	}
	// Run tidy again after upgrade
	if err := runCmd(ctx, projectDir, "go", "mod", "tidy"); err != nil {
		return fmt.Errorf("go mod tidy after upgrade failed: %w", err)
	}
	mlog.Print("Dependencies upgraded successfully!")
	return nil
}
