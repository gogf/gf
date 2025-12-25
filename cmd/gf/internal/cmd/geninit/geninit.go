// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package geninit

import (
	"context"
	"path/filepath"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

// ProcessOptions contains options for the Process function
type ProcessOptions struct {
	SelectVersion bool   // Enable interactive version selection
	ModulePath    string // Custom go.mod module path (e.g., github.com/xxx/xxx)
	UpgradeDeps   bool   // Upgrade dependencies to latest (go get -u ./...)
}

// Process handles the template generation flow from remote repository
func Process(ctx context.Context, repo, name string, opts *ProcessOptions) error {
	if opts == nil {
		opts = &ProcessOptions{}
	}

	// 0. Check Go environment first
	mlog.Print("Checking Go environment...")
	goEnv, err := CheckGoEnv(ctx)
	if err != nil {
		mlog.Printf("Go environment check failed: %v", err)
		return err
	}
	mlog.Printf("Go environment OK (version: %s)", goEnv.GOVERSION)

	// Check if this is a git subdirectory URL
	if IsSubdirRepo(repo) {
		return processGitSubdir(ctx, repo, name, opts)
	}

	return processGoModule(ctx, repo, name, opts)
}

// processGoModule handles standard Go module download via go get
func processGoModule(ctx context.Context, repo, name string, opts *ProcessOptions) error {
	// Extract module path (without version)
	modulePath := repo
	specifiedVersion := ""
	if gstr.Contains(repo, "@") {
		parts := gstr.Split(repo, "@")
		modulePath = parts[0]
		specifiedVersion = parts[1]
	}

	// Default name to repo basename if empty
	if name == "" {
		name = filepath.Base(modulePath)
	}

	// Determine the target module path for go.mod
	targetModulePath := name
	if opts.ModulePath != "" {
		targetModulePath = opts.ModulePath
	}

	// 1. Determine version to use
	var targetVersion string
	if specifiedVersion != "" {
		// User specified version
		targetVersion = specifiedVersion
		mlog.Printf("Using specified version: %s", targetVersion)
	} else if opts.SelectVersion {
		// Interactive version selection
		mlog.Print("Fetching available versions...")
		versionInfo, err := GetModuleVersions(ctx, modulePath)
		if err != nil {
			mlog.Printf("Failed to get versions: %v", err)
			return err
		}

		targetVersion, err = SelectVersion(ctx, versionInfo.Versions, modulePath)
		if err != nil {
			mlog.Printf("Version selection failed: %v", err)
			return err
		}
	} else {
		// Default: use latest version
		mlog.Print("Fetching latest version...")
		latest, err := GetLatestVersion(ctx, modulePath)
		if err != nil {
			mlog.Printf("Failed to get latest version, will try @latest tag: %v", err)
			targetVersion = "latest"
		} else {
			targetVersion = latest
			mlog.Printf("Latest version: %s", targetVersion)
		}
	}

	// 2. Download Template with determined version
	repoWithVersion := modulePath + "@" + targetVersion
	srcDir, err := downloadTemplate(ctx, repoWithVersion)
	if err != nil {
		mlog.Printf("Download failed: %v", err)
		return err
	}

	mlog.Debugf("Template located at: %s", srcDir)

	// 3. Generate Project
	if err := generateProject(ctx, srcDir, name, modulePath, targetModulePath); err != nil {
		mlog.Printf("Generation failed: %v", err)
		return err
	}

	// 4. Handle dependencies
	var projectDir string
	if name == "." {
		projectDir = gfile.Pwd()
	} else {
		projectDir = filepath.Join(gfile.Pwd(), name)
	}
	if opts.UpgradeDeps {
		// Upgrade all dependencies to latest
		if err := upgradeDependencies(ctx, projectDir); err != nil {
			mlog.Printf("Failed to upgrade dependencies: %v", err)
		}
	} else {
		// Default: just tidy dependencies
		if err := tidyDependencies(ctx, projectDir); err != nil {
			mlog.Printf("Failed to tidy dependencies: %v", err)
		}
	}

	return nil
}

// processGitSubdir handles git subdirectory download via sparse checkout
func processGitSubdir(ctx context.Context, repo, name string, opts *ProcessOptions) error {
	mlog.Print("Detected subdirectory URL, using git sparse checkout...")

	// Check if git is available
	gitVersion, err := CheckGitEnv(ctx)
	if err != nil {
		mlog.Printf("Git is required for subdirectory templates: %v", err)
		return err
	}
	mlog.Printf("Git available (%s)", gitVersion)

	// Download via git sparse checkout
	srcDir, gitInfo, err := downloadGitSubdir(ctx, repo)
	if err != nil {
		mlog.Printf("Git download failed: %v", err)
		return err
	}

	// Clean up temp directory after generation
	// The temp dir is parent of parent of srcDir (tempDir/repo/subpath)
	tempDir := filepath.Dir(filepath.Dir(srcDir))
	if gstr.Contains(tempDir, "gf-init-git") {
		defer gfile.Remove(tempDir)
	}

	// Default name to subpath basename if empty
	if name == "" {
		name = filepath.Base(gitInfo.SubPath)
	}

	// Get original module name from go.mod (might be "main" or something else)
	oldModule := GetModuleNameFromGoMod(srcDir)
	if oldModule == "" {
		// Fallback: construct from git info
		oldModule = gitInfo.Host + "/" + gitInfo.Owner + "/" + gitInfo.Repo + "/" + gitInfo.SubPath
	}

	// Determine the target module path for go.mod
	targetModulePath := name
	if opts.ModulePath != "" {
		targetModulePath = opts.ModulePath
	}

	mlog.Debugf("Template located at: %s", srcDir)
	mlog.Debugf("Original module: %s", oldModule)

	// Generate Project
	if err := generateProject(ctx, srcDir, name, oldModule, targetModulePath); err != nil {
		mlog.Printf("Generation failed: %v", err)
		return err
	}

	// Handle dependencies
	var projectDir string
	if name == "." {
		projectDir = gfile.Pwd()
	} else {
		projectDir = filepath.Join(gfile.Pwd(), name)
	}
	if opts.UpgradeDeps {
		// Upgrade all dependencies to latest
		if err := upgradeDependencies(ctx, projectDir); err != nil {
			mlog.Printf("Failed to upgrade dependencies: %v", err)
		}
	} else {
		// Default: just tidy dependencies
		if err := tidyDependencies(ctx, projectDir); err != nil {
			mlog.Printf("Failed to tidy dependencies: %v", err)
		}
	}

	return nil
}
