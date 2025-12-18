package logic

import (
	"context"
	"path/filepath"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

// ProcessOptions contains options for the Process function
type ProcessOptions struct {
	SelectVersion bool   // Enable interactive version selection
	ModulePath    string // Custom go.mod module path (e.g., github.com/xxx/xxx)
	UpgradeDeps   bool   // Upgrade dependencies to latest (go get -u ./...)
}

// Process handles the template generation flow
func Process(ctx context.Context, repo, name string, opts *ProcessOptions) error {
	if opts == nil {
		opts = &ProcessOptions{}
	}

	// 0. Check Go environment first
	g.Log().Info(ctx, "Checking Go environment...")
	goEnv, err := CheckGoEnv(ctx)
	if err != nil {
		g.Log().Error(ctx, "Go environment check failed:", err)
		return err
	}
	g.Log().Infof(ctx, "Go environment OK (version: %s)", goEnv.GOVERSION)

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
		g.Log().Infof(ctx, "Using specified version: %s", targetVersion)
	} else if opts.SelectVersion {
		// Interactive version selection
		g.Log().Info(ctx, "Fetching available versions...")
		versionInfo, err := GetModuleVersions(ctx, modulePath)
		if err != nil {
			g.Log().Error(ctx, "Failed to get versions:", err)
			return err
		}

		targetVersion, err = SelectVersion(ctx, versionInfo.Versions, modulePath)
		if err != nil {
			g.Log().Error(ctx, "Version selection failed:", err)
			return err
		}
	} else {
		// Default: use latest version
		g.Log().Info(ctx, "Fetching latest version...")
		latest, err := GetLatestVersion(ctx, modulePath)
		if err != nil {
			g.Log().Warning(ctx, "Failed to get latest version, will try @latest tag:", err)
			targetVersion = "latest"
		} else {
			targetVersion = latest
			g.Log().Infof(ctx, "Latest version: %s", targetVersion)
		}
	}

	// 2. Download Template with determined version
	repoWithVersion := modulePath + "@" + targetVersion
	srcDir, err := downloadTemplate(ctx, repoWithVersion)
	if err != nil {
		g.Log().Error(ctx, "Download failed:", err)
		return err
	}

	g.Log().Debug(ctx, "Template located at:", srcDir)

	// 3. Generate Project
	if err := generateProject(ctx, srcDir, name, modulePath, targetModulePath); err != nil {
		g.Log().Error(ctx, "Generation failed:", err)
		return err
	}

	// 4. Handle dependencies
	projectDir := filepath.Join(gfile.Pwd(), name)
	if opts.UpgradeDeps {
		// Upgrade all dependencies to latest
		if err := upgradeDependencies(ctx, projectDir); err != nil {
			g.Log().Warning(ctx, "Failed to upgrade dependencies:", err)
		}
	} else {
		// Default: just tidy dependencies
		if err := tidyDependencies(ctx, projectDir); err != nil {
			g.Log().Warning(ctx, "Failed to tidy dependencies:", err)
		}
	}

	return nil
}

// processGitSubdir handles git subdirectory download via sparse checkout
func processGitSubdir(ctx context.Context, repo, name string, opts *ProcessOptions) error {
	g.Log().Info(ctx, "Detected subdirectory URL, using git sparse checkout...")

	// Check if git is available
	gitVersion, err := CheckGitEnv(ctx)
	if err != nil {
		g.Log().Error(ctx, "Git is required for subdirectory templates:", err)
		return err
	}
	g.Log().Infof(ctx, "Git available (%s)", gitVersion)

	// Download via git sparse checkout
	srcDir, gitInfo, err := downloadGitSubdir(ctx, repo)
	if err != nil {
		g.Log().Error(ctx, "Git download failed:", err)
		return err
	}

	// Clean up temp directory after generation
	// The temp dir is parent of parent of srcDir (tempDir/repo/subpath)
	tempDir := filepath.Dir(filepath.Dir(srcDir))
	if gstr.Contains(tempDir, "tpl-git") {
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

	g.Log().Debug(ctx, "Template located at:", srcDir)
	g.Log().Debug(ctx, "Original module:", oldModule)

	// Generate Project
	if err := generateProject(ctx, srcDir, name, oldModule, targetModulePath); err != nil {
		g.Log().Error(ctx, "Generation failed:", err)
		return err
	}

	// Handle dependencies
	projectDir := filepath.Join(gfile.Pwd(), name)
	if opts.UpgradeDeps {
		// Upgrade all dependencies to latest
		if err := upgradeDependencies(ctx, projectDir); err != nil {
			g.Log().Warning(ctx, "Failed to upgrade dependencies:", err)
		}
	} else {
		// Default: just tidy dependencies
		if err := tidyDependencies(ctx, projectDir); err != nil {
			g.Log().Warning(ctx, "Failed to tidy dependencies:", err)
		}
	}

	return nil
}
