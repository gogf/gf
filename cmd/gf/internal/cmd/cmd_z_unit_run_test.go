// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_cRunApp_getWatchPaths_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		app := &cRunApp{
			WatchPaths: []string{"."},
		}
		watchPaths := app.getWatchPaths()

		t.AssertGT(len(watchPaths), 0)
		for _, v := range watchPaths {
			t.Log(v)
		}
	})
}

func Test_cRunApp_getWatchPaths_EmptyWatchPaths(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		app := &cRunApp{
			WatchPaths: []string{},
		}
		watchPaths := app.getWatchPaths()

		// Should default to current directory "."
		t.AssertGT(len(watchPaths), 0)
	})
}

func Test_cRunApp_getWatchPaths_CustomIgnorePattern(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		app := &cRunApp{
			WatchPaths:     []string{"testdata"},
			IgnorePatterns: []string{"2572"},
		}
		watchPaths := app.getWatchPaths()

		// Ensure the "2572" directory is not watched directly.
		for _, path := range watchPaths {
			t.Log("watch path:", path)
			t.Assert(strings.HasSuffix(path, "2572"), false)
		}
		t.AssertGT(len(watchPaths), 0)
	})
}

func Test_cRunApp_getWatchPaths_WithIgnoredDirectories(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create a temporary directory structure for testing
		tempDir := gfile.Temp("gf_run_test")
		defer gfile.Remove(tempDir)

		// Create directory structure:
		// tempDir/
		//   ├── src/
		//   │   ├── api/
		//   │   └── internal/
		//   ├── vendor/  <-- ignored
		//   └── node_modules/  <-- ignored
		gfile.Mkdir(filepath.Join(tempDir, "src", "api"))
		gfile.Mkdir(filepath.Join(tempDir, "src", "internal"))
		gfile.Mkdir(filepath.Join(tempDir, "vendor"))
		gfile.Mkdir(filepath.Join(tempDir, "node_modules"))

		app := &cRunApp{
			WatchPaths: []string{tempDir},
		}
		watchPaths := app.getWatchPaths()

		// Should watch the src directory since parent has ignored children
		t.Assert(len(watchPaths), 1)
		t.Assert(watchPaths[0], filepath.Join(tempDir, "src"))
	})
}

func Test_cRunApp_getWatchPaths_NoIgnoredDirectories(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create a temporary directory structure without ignored directories
		tempDir := gfile.Temp("gf_run_test_no_ignore")
		defer gfile.Remove(tempDir)

		// Create directory structure without ignored patterns:
		// tempDir/
		//   ├── src/
		//   │   ├── api/
		//   │   └── internal/
		gfile.Mkdir(filepath.Join(tempDir, "src", "api"))
		gfile.Mkdir(filepath.Join(tempDir, "src", "internal"))

		app := &cRunApp{
			WatchPaths: []string{tempDir},
		}
		watchPaths := app.getWatchPaths()

		// Should watch the root directory since no ignored directories exist
		t.Assert(len(watchPaths), 1)
		t.Assert(watchPaths[0], tempDir)
	})
}

func Test_cRunApp_getWatchPaths_CustomIgnorePatterns(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create a temporary directory structure
		tempDir := gfile.Temp("gf_run_test_custom_ignore")
		defer gfile.Remove(tempDir)

		// Create directory structure:
		// tempDir/
		//   ├── src/
		//   │   ├── api/
		//   │   └── internal/
		//   ├── build/  <-- ignored
		//   └── dist/  <-- ignored
		gfile.Mkdir(filepath.Join(tempDir, "src", "api"))
		gfile.Mkdir(filepath.Join(tempDir, "src", "internal"))
		gfile.Mkdir(filepath.Join(tempDir, "build"))
		gfile.Mkdir(filepath.Join(tempDir, "dist"))

		app := &cRunApp{
			WatchPaths:     []string{tempDir},
			IgnorePatterns: []string{"build", "dist"},
		}
		watchPaths := app.getWatchPaths()

		// Should watch the src directory since parent has ignored children (build, dist)
		t.Assert(len(watchPaths), 1)
		t.Assert(watchPaths[0], filepath.Join(tempDir, "src"))
	})
}

func Test_cRunApp_getWatchPaths_DeepNestedStructure(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create a deep nested directory structure
		tempDir := gfile.Temp("gf_run_test_deep")
		defer gfile.Remove(tempDir)

		// Create deep directory structure:
		// tempDir/
		//   ├── a/
		//   │   ├── b/
		//   │   │   └── c/
		//   │   └── vendor/  <-- ignored
		//   └── d/
		gfile.Mkdir(filepath.Join(tempDir, "a", "b", "c"))
		gfile.Mkdir(filepath.Join(tempDir, "a", "vendor"))
		gfile.Mkdir(filepath.Join(tempDir, "d"))

		app := &cRunApp{
			WatchPaths: []string{tempDir},
		}
		watchPaths := app.getWatchPaths()

		// Should watch individual valid directories due to ignored vendor directory
		t.AssertGT(len(watchPaths), 0)

		// Verify that vendor directory is not in watch list
		for _, path := range watchPaths {
			t.Assert(strings.Contains(path, "vendor"), false)
		}
	})
}

func Test_cRunApp_getWatchPaths_MultipleRoots(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create multiple temporary directories
		tempDir1 := gfile.Temp("gf_run_test_multi1")
		tempDir2 := gfile.Temp("gf_run_test_multi2")
		defer gfile.Remove(tempDir1)
		defer gfile.Remove(tempDir2)

		gfile.Mkdir(filepath.Join(tempDir1, "src"))
		gfile.Mkdir(filepath.Join(tempDir2, "api"))

		app := &cRunApp{
			WatchPaths: []string{tempDir1, tempDir2},
		}
		watchPaths := app.getWatchPaths()

		// Should watch both root directories
		t.Assert(len(watchPaths), 2)

		// Both directories should be in the watch list
		foundDir1, foundDir2 := false, false
		for _, path := range watchPaths {
			if path == tempDir1 {
				foundDir1 = true
			}
			if path == tempDir2 {
				foundDir2 = true
			}
		}
		t.Assert(foundDir1, true)
		t.Assert(foundDir2, true)
	})
}

func Test_cRunApp_getWatchPaths_NonExistentDirectory(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		app := &cRunApp{
			WatchPaths: []string{"/non/existent/path"},
		}
		watchPaths := app.getWatchPaths()

		// Should fall back to current directory when no valid paths found
		t.AssertGT(len(watchPaths), 0)

		// Should contain current directory
		currentDir, _ := os.Getwd()
		foundCurrentDir := false
		for _, path := range watchPaths {
			if path == currentDir {
				foundCurrentDir = true
				break
			}
		}
		t.Assert(foundCurrentDir, true)
	})
}

func Test_isIgnoredDirName(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test default ignore patterns
		t.Assert(isIgnoredDirName("node_modules", defaultIgnorePatterns), true)
		t.Assert(isIgnoredDirName("vendor", defaultIgnorePatterns), true)
		t.Assert(isIgnoredDirName(".git", defaultIgnorePatterns), true)
		t.Assert(isIgnoredDirName("_private", defaultIgnorePatterns), true)
		t.Assert(isIgnoredDirName("src", defaultIgnorePatterns), false)
		t.Assert(isIgnoredDirName("api", defaultIgnorePatterns), false)

		// Test custom ignore patterns
		customPatterns := []string{"build", "dist", "*.tmp"}
		t.Assert(isIgnoredDirName("build", customPatterns), true)
		t.Assert(isIgnoredDirName("dist", customPatterns), true)
		t.Assert(isIgnoredDirName("test.tmp", customPatterns), true)
		t.Assert(isIgnoredDirName("src", customPatterns), false)
	})
}
