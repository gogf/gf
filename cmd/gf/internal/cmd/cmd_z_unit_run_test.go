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
		for _, wp := range watchPaths {
			t.Log("watch path:", wp)
			t.Assert(strings.HasSuffix(wp.Path, "2572"), false)
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

		// Should watch tempDir non-recursively (to catch top-level files) and src recursively
		t.Assert(len(watchPaths), 2)
		// First path is tempDir (non-recursive)
		t.Assert(watchPaths[0].Path, tempDir)
		t.Assert(watchPaths[0].Recursive, false)
		// Second path is src (recursive, since it has no ignored descendants)
		t.Assert(watchPaths[1].Path, filepath.Join(tempDir, "src"))
		t.Assert(watchPaths[1].Recursive, true)
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

		// Should watch the root directory recursively since no ignored directories exist
		t.Assert(len(watchPaths), 1)
		t.Assert(watchPaths[0].Path, tempDir)
		t.Assert(watchPaths[0].Recursive, true)
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

		// Should watch tempDir non-recursively and src recursively
		t.Assert(len(watchPaths), 2)
		t.Assert(watchPaths[0].Path, tempDir)
		t.Assert(watchPaths[0].Recursive, false)
		t.Assert(watchPaths[1].Path, filepath.Join(tempDir, "src"))
		t.Assert(watchPaths[1].Recursive, true)
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
		for _, wp := range watchPaths {
			t.Assert(strings.Contains(wp.Path, "vendor"), false)
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

		// Should watch both root directories recursively
		t.Assert(len(watchPaths), 2)

		// Both directories should be in the watch list
		foundDir1, foundDir2 := false, false
		for _, wp := range watchPaths {
			if wp.Path == tempDir1 {
				foundDir1 = true
				t.Assert(wp.Recursive, true)
			}
			if wp.Path == tempDir2 {
				foundDir2 = true
				t.Assert(wp.Recursive, true)
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
		for _, wp := range watchPaths {
			if wp.Path == currentDir {
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

func Test_hasIgnoredDescendant(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create a temporary directory structure
		tempDir := gfile.Temp("gf_run_test_has_ignored")
		defer gfile.Remove(tempDir)

		// Create directory structure:
		// tempDir/
		//   ├── a/
		//   │   ├── b/
		//   │   │   └── vendor/  <-- deeply nested ignored
		//   │   └── c/
		//   └── d/
		gfile.Mkdir(filepath.Join(tempDir, "a", "b", "vendor"))
		gfile.Mkdir(filepath.Join(tempDir, "a", "c"))
		gfile.Mkdir(filepath.Join(tempDir, "d"))

		// Test: tempDir should have ignored descendant (vendor is 3 levels deep)
		t.Assert(hasIgnoredDescendant(tempDir, defaultIgnorePatterns), true)

		// Test: d/ should NOT have ignored descendant
		t.Assert(hasIgnoredDescendant(filepath.Join(tempDir, "d"), defaultIgnorePatterns), false)

		// Test: a/c/ should NOT have ignored descendant
		t.Assert(hasIgnoredDescendant(filepath.Join(tempDir, "a", "c"), defaultIgnorePatterns), false)

		// Test: a/ should have ignored descendant (vendor in a/b/)
		t.Assert(hasIgnoredDescendant(filepath.Join(tempDir, "a"), defaultIgnorePatterns), true)

		// Test: a/b/ should have ignored descendant (vendor directly inside)
		t.Assert(hasIgnoredDescendant(filepath.Join(tempDir, "a", "b"), defaultIgnorePatterns), true)
	})
}

func Test_cRunApp_getWatchPaths_DeeplyNestedIgnore(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create a temporary directory structure with deeply nested ignored directory
		tempDir := gfile.Temp("gf_run_test_deeply_nested")
		defer gfile.Remove(tempDir)

		// Create directory structure:
		// tempDir/
		//   ├── a/
		//   │   ├── b/
		//   │   │   ├── c/
		//   │   │   │   └── vendor/  <-- deeply nested ignored (4 levels)
		//   │   │   └── d/
		//   │   └── e/
		//   └── f/
		gfile.Mkdir(filepath.Join(tempDir, "a", "b", "c", "vendor"))
		gfile.Mkdir(filepath.Join(tempDir, "a", "b", "d"))
		gfile.Mkdir(filepath.Join(tempDir, "a", "e"))
		gfile.Mkdir(filepath.Join(tempDir, "f"))

		app := &cRunApp{
			WatchPaths: []string{tempDir},
		}
		watchPaths := app.getWatchPaths()

		// Expected watch paths:
		// 1. tempDir (non-recursive) - has ignored descendant
		// 2. a (non-recursive) - has ignored descendant in b/c/vendor
		// 3. b (non-recursive) - has ignored descendant in c/vendor
		// 4. c (non-recursive) - has ignored child vendor
		// 5. d (recursive) - no ignored descendants
		// 6. e (recursive) - no ignored descendants
		// 7. f (recursive) - no ignored descendants

		t.AssertGT(len(watchPaths), 0)

		// Verify vendor is not in watch paths
		for _, wp := range watchPaths {
			t.Assert(strings.Contains(wp.Path, "vendor"), false)
		}

		// Find specific paths and verify their recursive flags
		foundF := false
		for _, wp := range watchPaths {
			if wp.Path == filepath.Join(tempDir, "f") {
				foundF = true
				t.Assert(wp.Recursive, true) // f should be recursive (no ignored descendants)
			}
		}
		t.Assert(foundF, true)
	})
}

func Test_cRunApp_getWatchPaths_EmptyDirectory(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create an empty temporary directory
		tempDir := gfile.Temp("gf_run_test_empty")
		defer gfile.Remove(tempDir)

		gfile.Mkdir(tempDir)

		app := &cRunApp{
			WatchPaths: []string{tempDir},
		}
		watchPaths := app.getWatchPaths()

		// Empty directory should be watched recursively (no ignored descendants)
		t.Assert(len(watchPaths), 1)
		t.Assert(watchPaths[0].Path, tempDir)
		t.Assert(watchPaths[0].Recursive, true)
	})
}
