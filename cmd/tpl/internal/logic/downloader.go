package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

// downloadTemplate fetches the remote repository using go get
func downloadTemplate(ctx context.Context, repo string) (string, error) {
	// 1. Create a temporary directory workspace
	tempDir := gfile.Temp("tpl-cli")
	if err := gfile.Mkdir(tempDir); err != nil {
		return "", err
	}
	defer gfile.Remove(tempDir) // Clean up the temp workspace

	g.Log().Debug(ctx, "Using temp workspace:", tempDir)

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

	var downloadErr error
	versionsToTry := []string{repo}
	if !gstr.Contains(repo, "@") {
		versionsToTry = append(versionsToTry, repo+"@latest", repo+"@master")
	}

	var successRepo string
	for _, tryRepo := range versionsToTry {
		g.Log().Infof(ctx, "Downloading template %s...", tryRepo)
		if err := runCmd(ctx, tempDir, "go", "get", tryRepo); err == nil {
			successRepo = tryRepo
			break
		} else {
			downloadErr = err
			g.Log().Debugf(ctx, "Failed to download %s, trying next...", tryRepo)
		}
	}

	if successRepo == "" {
		return "", fmt.Errorf("failed to download repo %s: %w", repo, downloadErr)
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
