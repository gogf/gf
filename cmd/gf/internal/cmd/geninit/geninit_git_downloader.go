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
	"strings"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

// GitRepoInfo holds parsed git repository information
type GitRepoInfo struct {
	Host     string // e.g., github.com
	Owner    string // e.g., gogf
	Repo     string // e.g., examples
	Branch   string // e.g., main (default: main)
	SubPath  string // e.g., httpserver/jwt
	CloneURL string // e.g., https://github.com/gogf/examples.git
}

// ParseGitURL parses a git URL and extracts repository info
// Supports formats:
//   - github.com/owner/repo
//   - github.com/owner/repo/subdir/path
//   - github.com/owner/repo/tree/branch/subdir/path (from GitHub web URL)
func ParseGitURL(url string) (*GitRepoInfo, error) {
	// Remove protocol prefix if present
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimSuffix(url, ".git")

	// Remove version suffix like @v1.0.0
	if idx := strings.Index(url, "@"); idx != -1 {
		url = url[:idx]
	}

	parts := strings.Split(url, "/")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid git URL: %s", url)
	}

	info := &GitRepoInfo{
		Host:   parts[0],
		Owner:  parts[1],
		Repo:   parts[2],
		Branch: "main", // default branch
	}

	// Check for /tree/branch/ pattern (GitHub web URL)
	if len(parts) > 4 && parts[3] == "tree" {
		info.Branch = parts[4]
		if len(parts) > 5 {
			info.SubPath = strings.Join(parts[5:], "/")
		}
	} else if len(parts) > 3 {
		// Direct subpath: github.com/owner/repo/subdir/path
		info.SubPath = strings.Join(parts[3:], "/")
	}

	info.CloneURL = fmt.Sprintf("https://%s/%s/%s.git", info.Host, info.Owner, info.Repo)

	return info, nil
}

// IsSubdirRepo checks if the URL points to a subdirectory of a repository
// Returns false for Go module paths (which may have /vN suffix or nested module paths)
func IsSubdirRepo(url string) bool {
	info, err := ParseGitURL(url)
	if err != nil {
		return false
	}
	if info.SubPath == "" {
		return false
	}

	// Check if this looks like a Go module path rather than a git subdirectory
	// Go modules can have nested paths like github.com/owner/repo/cmd/tool/v2
	// We should try to resolve it as a Go module first

	// If the URL can be resolved as a Go module, it's not a subdir repo
	// We use a heuristic: check if the full path looks like a valid Go module
	// by checking if it ends with /vN (major version) or contains common module patterns

	// Remove version suffix for checking
	cleanURL := url
	if idx := strings.Index(url, "@"); idx != -1 {
		cleanURL = url[:idx]
	}

	// Check if the path ends with /vN (Go module major version)
	parts := strings.Split(cleanURL, "/")
	if len(parts) > 0 {
		lastPart := parts[len(parts)-1]
		if len(lastPart) >= 2 && lastPart[0] == 'v' {
			// Check if it's v2, v3, etc.
			if _, err := fmt.Sscanf(lastPart, "v%d", new(int)); err == nil {
				// This looks like a Go module with major version suffix
				// It could be either a versioned module or a subdir ending in vN
				// We'll treat it as a Go module and let go get handle it
				return false
			}
		}
	}

	// For GitHub URLs, check if the subpath could be a nested Go module
	// Common patterns: cmd/*, internal/*, pkg/*, contrib/*
	subPathParts := strings.Split(info.SubPath, "/")
	if len(subPathParts) > 0 {
		firstPart := subPathParts[0]
		// These are common Go module nesting patterns
		if firstPart == "cmd" || firstPart == "contrib" || firstPart == "tools" {
			// This might be a nested Go module, not a simple subdirectory
			// Let go get try first
			return false
		}
	}

	return true
}

// downloadGitSubdir downloads a subdirectory from a git repository using sparse checkout
func downloadGitSubdir(ctx context.Context, repoURL string) (string, *GitRepoInfo, error) {
	info, err := ParseGitURL(repoURL)
	if err != nil {
		return "", nil, err
	}

	if info.SubPath == "" {
		return "", nil, fmt.Errorf("not a subdirectory URL: %s", repoURL)
	}

	// Create temp directory for clone
	tempDir := gfile.Temp("gf-init-git")
	if err := gfile.Mkdir(tempDir); err != nil {
		return "", nil, err
	}

	cloneDir := filepath.Join(tempDir, info.Repo)
	mlog.Debugf("Using git temp workspace: %s", tempDir)
	mlog.Printf("Cloning %s (sparse checkout: %s)...", info.CloneURL, info.SubPath)

	// 1. Clone with no checkout, filter, and sparse
	if err := runCmd(ctx, tempDir, "git", "clone", "--filter=blob:none", "--no-checkout", "--sparse", info.CloneURL); err != nil {
		// Fallback: try without filter for older git versions
		mlog.Debugf("Sparse clone failed, trying full clone...")
		gfile.Remove(cloneDir)
		if err := runCmd(ctx, tempDir, "git", "clone", "--no-checkout", info.CloneURL); err != nil {
			gfile.Remove(tempDir)
			return "", nil, fmt.Errorf("git clone failed: %w", err)
		}
	}

	// 2. Set sparse-checkout to the subpath
	if err := runCmd(ctx, cloneDir, "git", "sparse-checkout", "set", info.SubPath); err != nil {
		// Fallback for older git: use sparse-checkout init + echo
		mlog.Debugf("sparse-checkout set failed, trying legacy method...")
		runCmd(ctx, cloneDir, "git", "sparse-checkout", "init", "--cone")
		runCmd(ctx, cloneDir, "git", "sparse-checkout", "set", info.SubPath)
	}

	// 3. Checkout the branch
	if err := runCmd(ctx, cloneDir, "git", "checkout", info.Branch); err != nil {
		// Try master if main fails
		if info.Branch == "main" {
			mlog.Debugf("Branch 'main' not found, trying 'master'...")
			info.Branch = "master"
			if err := runCmd(ctx, cloneDir, "git", "checkout", "master"); err != nil {
				gfile.Remove(tempDir)
				return "", nil, fmt.Errorf("git checkout failed: %w", err)
			}
		} else {
			gfile.Remove(tempDir)
			return "", nil, fmt.Errorf("git checkout failed: %w", err)
		}
	}

	// Return the path to the subdirectory
	subDirPath := filepath.Join(cloneDir, info.SubPath)
	if !gfile.Exists(subDirPath) {
		gfile.Remove(tempDir)
		return "", nil, fmt.Errorf("subdirectory not found: %s", info.SubPath)
	}

	mlog.Debugf("Subdirectory located at: %s", subDirPath)
	return subDirPath, info, nil
}

// GetModuleNameFromGoMod reads module name from go.mod file
func GetModuleNameFromGoMod(dir string) string {
	goModPath := filepath.Join(dir, "go.mod")
	if !gfile.Exists(goModPath) {
		return ""
	}

	content := gfile.GetContents(goModPath)
	lines := gstr.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return ""
}
