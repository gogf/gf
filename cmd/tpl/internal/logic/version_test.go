package logic

import (
	"context"
	"testing"
)

func TestGetModuleVersions(t *testing.T) {
	ctx := context.Background()

	// Test with a well-known module that has multiple versions
	info, err := GetModuleVersions(ctx, "github.com/gogf/gf/v2")
	if err != nil {
		t.Skipf("Skipping test due to network issue: %v", err)
	}

	if info.Module == "" {
		t.Error("Module name should not be empty")
	}

	if len(info.Versions) == 0 {
		t.Error("Should have at least one version")
	}

	if info.Latest == "" {
		t.Error("Latest version should not be empty")
	}

	t.Logf("Module: %s, Latest: %s, Total versions: %d", info.Module, info.Latest, len(info.Versions))
}

func TestGetLatestVersion(t *testing.T) {
	ctx := context.Background()

	version, err := GetLatestVersion(ctx, "github.com/gogf/gf/v2")
	if err != nil {
		t.Skipf("Skipping test due to network issue: %v", err)
	}

	if version == "" {
		t.Error("Version should not be empty")
	}

	t.Logf("Latest version: %s", version)
}
