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
	"os/exec"
	"strings"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

// GoEnv represents Go environment variables
type GoEnv struct {
	GOVERSION   string `json:"GOVERSION"`
	GOROOT      string `json:"GOROOT"`
	GOPATH      string `json:"GOPATH"`
	GOMODCACHE  string `json:"GOMODCACHE"`
	GOPROXY     string `json:"GOPROXY"`
	GO111MODULE string `json:"GO111MODULE"`
}

// CheckGoEnv verifies Go is installed and properly configured
func CheckGoEnv(ctx context.Context) (*GoEnv, error) {
	// 1. Check if go binary exists
	goPath, err := exec.LookPath("go")
	if err != nil {
		return nil, fmt.Errorf("go is not installed or not in PATH: %w", err)
	}
	mlog.Debugf("Found go binary at: %s", goPath)

	// 2. Get go env as JSON
	cmd := exec.CommandContext(ctx, "go", "env", "-json")
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("go env failed: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("failed to run go env: %w", err)
	}

	// 3. Parse JSON output
	var env GoEnv
	if err := json.Unmarshal(output, &env); err != nil {
		return nil, fmt.Errorf("failed to parse go env output: %w", err)
	}

	// 4. Validate critical environment variables
	if env.GOROOT == "" {
		return nil, fmt.Errorf("GOROOT is not set")
	}
	if env.GOMODCACHE == "" && env.GOPATH == "" {
		return nil, fmt.Errorf("neither GOMODCACHE nor GOPATH is set")
	}

	mlog.Debugf("Go Version: %s", env.GOVERSION)
	mlog.Debugf("GOROOT: %s", env.GOROOT)
	mlog.Debugf("GOMODCACHE: %s", env.GOMODCACHE)
	mlog.Debugf("GOPROXY: %s", env.GOPROXY)

	return &env, nil
}

// CheckGitEnv verifies Git is installed and returns its version
func CheckGitEnv(ctx context.Context) (string, error) {
	// 1. Check if git binary exists
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return "", fmt.Errorf("git is not installed or not in PATH: %w", err)
	}
	mlog.Debugf("Found git binary at: %s", gitPath)

	// 2. Get git version
	cmd := exec.CommandContext(ctx, "git", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git version: %w", err)
	}

	version := strings.TrimSpace(string(output))
	mlog.Debugf("Git version: %s", version)

	return version, nil
}
