package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"golang.org/x/mod/semver"
)

// VersionInfo contains module version information
type VersionInfo struct {
	Module   string   `json:"module"`
	Versions []string `json:"versions"`
	Latest   string   `json:"latest"`
}

// GetModuleVersions fetches available versions for a Go module
func GetModuleVersions(ctx context.Context, modulePath string) (*VersionInfo, error) {
	// Create a temporary directory for go list
	tempDir := gfile.Temp("tpl-version")
	if err := gfile.Mkdir(tempDir); err != nil {
		return nil, err
	}
	defer gfile.Remove(tempDir)

	// Initialize a temp go module
	if err := runCmd(ctx, tempDir, "go", "mod", "init", "temp"); err != nil {
		return nil, fmt.Errorf("failed to init temp module: %w", err)
	}

	// Get versions using go list -m -versions
	cmd := exec.CommandContext(ctx, "go", "list", "-m", "-versions", modulePath)
	cmd.Dir = tempDir
	output, err := cmd.Output()
	if err != nil {
		// Try with @latest to see if module exists
		g.Log().Debugf(ctx, "go list -versions failed, trying @latest: %v", err)
		return getLatestOnly(ctx, tempDir, modulePath)
	}

	// Parse output: "module/path v1.0.0 v1.1.0 v2.0.0"
	parts := strings.Fields(strings.TrimSpace(string(output)))
	if len(parts) < 1 {
		return nil, fmt.Errorf("no version information found for %s", modulePath)
	}

	info := &VersionInfo{
		Module:   parts[0],
		Versions: []string{},
	}

	if len(parts) > 1 {
		info.Versions = parts[1:]
		// Sort versions in descending order (newest first)
		sort.Slice(info.Versions, func(i, j int) bool {
			return semver.Compare(info.Versions[i], info.Versions[j]) > 0
		})
		info.Latest = info.Versions[0]
	}

	// If no tagged versions, try to get latest
	if len(info.Versions) == 0 {
		latestInfo, err := getLatestOnly(ctx, tempDir, modulePath)
		if err != nil {
			return nil, err
		}
		info.Latest = latestInfo.Latest
		if latestInfo.Latest != "" {
			info.Versions = []string{latestInfo.Latest}
		}
	}

	return info, nil
}

// getLatestOnly gets only the latest version when go list -versions fails
func getLatestOnly(ctx context.Context, tempDir, modulePath string) (*VersionInfo, error) {
	// Try go list -m modulePath@latest
	cmd := exec.CommandContext(ctx, "go", "list", "-m", "-json", modulePath+"@latest")
	cmd.Dir = tempDir
	output, err := cmd.Output()
	if err != nil {
		// Try without @latest
		cmd = exec.CommandContext(ctx, "go", "list", "-m", "-json", modulePath)
		cmd.Dir = tempDir
		output, err = cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("failed to get module info for %s: %w", modulePath, err)
		}
	}

	var modInfo struct {
		Path    string `json:"Path"`
		Version string `json:"Version"`
	}
	if err := json.Unmarshal(output, &modInfo); err != nil {
		return nil, fmt.Errorf("failed to parse module info: %w", err)
	}

	return &VersionInfo{
		Module:   modInfo.Path,
		Versions: []string{modInfo.Version},
		Latest:   modInfo.Version,
	}, nil
}

// GetLatestVersion returns the latest version of a module
func GetLatestVersion(ctx context.Context, modulePath string) (string, error) {
	info, err := GetModuleVersions(ctx, modulePath)
	if err != nil {
		return "", err
	}
	if info.Latest == "" {
		return "", fmt.Errorf("no version found for %s", modulePath)
	}
	return info.Latest, nil
}
