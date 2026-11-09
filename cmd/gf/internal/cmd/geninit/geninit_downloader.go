// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package geninit

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

// downloadTemplate fetches the remote repository using go get
func downloadTemplate(ctx context.Context, repo string) (string, error) {
	// 1. Create a temporary directory workspace
	tempDir := gfile.Temp("gf-init-cli")
	if tempDir == "" {
		return "", fmt.Errorf("failed to create temporary directory")
	}
	if err := gfile.Mkdir(tempDir); err != nil {
		return "", err
	}
	defer func() {
		if err := gfile.Remove(tempDir); err != nil {
			mlog.Debugf("Failed to remove temp directory %s: %v", tempDir, err)
		}
	}() // Clean up the temp workspace

	mlog.Debugf("Using temp workspace: %s", tempDir)

	// 2. Initialize a temp go module to perform go get
	// We run commands inside the temp directory
	if err := runCmd(ctx, tempDir, "go", "mod", "init", "temp"); err != nil {
		return "", err
	}

	// 3. Run go get <repo>
	// Try different version strategies: original -> @latest -> @master
	moduleName := repo
	if gstr.Contains(repo, "@") {
		moduleName = gstr.Split(repo, "@")[0]
	}

	var downloadErrs []string
	versionsToTry := []string{repo}
	if !gstr.Contains(repo, "@") {
		versionsToTry = append(versionsToTry, repo+"@latest", repo+"@master")
	}

	var successRepo string
	for _, tryRepo := range versionsToTry {
		mlog.Printf("Downloading template %s...", tryRepo)
		if err := runCmd(ctx, tempDir, "go", "get", tryRepo); err == nil {
			successRepo = tryRepo
			break
		} else {
			downloadErrs = append(downloadErrs, fmt.Sprintf("%s: %v", tryRepo, err))
			mlog.Debugf("Failed to download %s, trying next...", tryRepo)
		}
	}

	if successRepo == "" {
		errMsg := "all download attempts failed"
		if len(downloadErrs) > 0 {
			errMsg = strings.Join(downloadErrs, "; ")
		}
		return "", fmt.Errorf("failed to download repo %s: %s", repo, errMsg)
	}

	// 4. Find the local path using go list -m -json <repo>
	listCmd := exec.CommandContext(ctx, "go", "list", "-m", "-json", moduleName)
	listCmd.Dir = tempDir
	output, err := listCmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("go list failed: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to locate module path: %w", err)
	}

	var modInfo struct {
		Dir string `json:"Dir"`
	}
	if err := json.Unmarshal(output, &modInfo); err != nil {
		return "", fmt.Errorf("failed to parse go list output: %w", err)
	}

	if modInfo.Dir == "" {
		return "", fmt.Errorf("module directory not found for %s", repo)
	}

	return modInfo.Dir, nil
}

func runCmd(ctx context.Context, dir string, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
