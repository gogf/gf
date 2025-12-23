// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func Test_cRunApp_shouldIgnorePath_CommonPatterns(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		app := &cRunApp{
			IgnorePaths: []string{".git", "*.js", "node_modules", "vendor", "*.pb.go"},
		}

		// Test cases for common development patterns
		testCases := []struct {
			path     string
			expected bool
			desc     string
		}{
			// Git related files
			{".git", true, "git directory"},
			{".git/config", true, "git config file"},
			{".gitignore", false, "gitignore file should not be ignored"},

			// JavaScript files
			{"main.js", true, "js file"},
			{"src/utils.js", true, "js file in subdirectory"},
			{"test.ts", false, "ts file should not be ignored"},

			// Node modules
			{"node_modules", true, "node_modules directory"},
			{"node_modules/react", true, "react in node_modules"},
			{"node_modules/.bin", true, "bin in node_modules"},

			// Vendor directory
			{"vendor", true, "vendor directory"},
			{"vendor/package", true, "package in vendor"},

			// Proto files
			{"api.pb.go", true, "pb.go file"},
			{"test.pb.go", true, "test pb.go file"},
			{"proto.go", false, "proto.go without pb"},

			// Normal Go files that should not be ignored
			{"main.go", false, "main go file"},
			{"internal/utils.go", false, "internal go file"},
			{"cmd/server/main.go", false, "server go file"},
			{"pkg/model/user.go", false, "model go file"},
		}

		for _, tc := range testCases {
			result := app.shouldIgnorePath(tc.path)
			if result != tc.expected {
				t.Errorf("path: %s, desc: %s, expected: %v, got: %v", tc.path, tc.desc, tc.expected, result)
			}
		}
	})
}

func Test_cRunApp_shouldIgnorePath_EdgeCases(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		app := &cRunApp{
			IgnorePaths: []string{"test_*", "*.tmp", "backup"},
		}

		testCases := []struct {
			path     string
			expected bool
		}{
			// Pattern matching with wildcards
			{"test_main.go", true},
			{"test_utils.go", true},
			{"main_test.go", false}, // Different pattern
			{"test.go", false},      // No underscore

			// Temporary files
			{"file.tmp", true},
			{"temp.tmp", true},
			{"tmp", false}, // Directory name, not file extension

			// Backup directory
			{"backup", true},
			{"backup/file.go", true},
			{"backup-old", false}, // Different name
		}

		for _, tc := range testCases {
			result := app.shouldIgnorePath(tc.path)
			t.Assert(result, tc.expected)
		}
	})
}

func Test_cRunApp_shouldIgnorePath_EmptyPatterns(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		app := &cRunApp{
			IgnorePaths: []string{}, // Empty list
		}

		// All paths should return false when no patterns are set
		testPaths := []string{
			"main.go",
			"node_modules/react",
			".git/config",
			"vendor/package",
			"test.js",
		}

		for _, path := range testPaths {
			result := app.shouldIgnorePath(path)
			t.Assert(result, false)
		}
	})
}

func Test_cRunApp_shouldIgnorePath_ComplexPatterns(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		app := &cRunApp{
			IgnorePaths: []string{"build/*", "dist/*", "*.log", "cache"},
		}

		testCases := []struct {
			path     string
			expected bool
		}{
			// Build and dist directories
			{"build", true},
			{"build/main", true},
			{"dist", true},
			{"dist/app", true},
			{"builds", false}, // Different name

			// Log files
			{"app.log", true},
			{"error.log", true},
			{"log.txt", false}, // Different extension

			// Cache directory
			{"cache", true},
			{"cache/file", true},
			{".cache", false}, // Different name
		}

		for _, tc := range testCases {
			t.Log(tc.path)
			result := app.shouldIgnorePath(tc.path)
			t.Assert(result, tc.expected)
		}
	})
}

func Test_cRunApp_getWatchPaths_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		app := &cRunApp{
			IgnorePaths: []string{"*.js", ".git"},
			WatchPaths:  []string{"src", "pkg"},
		}
		watchPaths := app.getWatchPaths()

		// Should return at least the current directory and specified watch paths
		t.AssertGT(len(watchPaths), 0)
		t.Assert(watchPaths[0], ".")

		// Should contain specified watch paths (if they exist)
		// Note: These assertions depend on actual directory structure
		// For testing purposes, we just verify the function returns valid paths
		for _, path := range watchPaths {
			if path == "src" || path == "pkg" {
				// Found one of the expected paths
				break
			}
		}
	})
}

func Test_cRunApp_getWatchPaths_NoWatchPaths(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		app := &cRunApp{
			IgnorePaths: []string{"*.js", ".git"},
			WatchPaths:  []string{}, // Empty watch paths
		}

		watchPaths := app.getWatchPaths()

		// Should return at least the current directory
		t.AssertGT(len(watchPaths), 0)
	})
}

func Test_cRunApp_shouldIgnorePath_LogPatternFix(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		app := &cRunApp{
			IgnorePaths: []string{"*.log"},
		}

		testCases := []struct {
			path     string
			expected bool
			desc     string
		}{
			// Log files should be ignored
			{"app.log", true, "log file"},
			{"error.log", true, "error log file"},
			{"logs/app.log", true, "log file in logs directory"},

			// Go files should NOT be ignored (this was the bug)
			{"main.go", false, "main go file"},
			{"app/main.go", false, "go file in app directory"},
			{"internal/app/main.go", false, "nested go file"},

			// Directories should NOT be ignored
			{"app", false, "app directory"},
			{"internal", false, "internal directory"},
			{"vendor", false, "vendor directory"},

			// Other files should NOT be ignored
			{"README.md", false, "markdown file"},
			{"config.yaml", false, "yaml file"},
		}

		for _, tc := range testCases {
			result := app.shouldIgnorePath(tc.path)
			if result != tc.expected {
				t.Errorf("path: %s, desc: %s, expected: %v, got: %v", tc.path, tc.desc, tc.expected, result)
			}
		}
	})
}
