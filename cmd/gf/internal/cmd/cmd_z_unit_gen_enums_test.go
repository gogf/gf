// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"path/filepath"
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gogf/gf/v2/util/gutil"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/genenums"
)

// https://github.com/gogf/gf/issues/4387
// Test that the output path is relative to the original working directory,
// not the source directory after Chdir.
func Test_Gen_Enums_Issue4387_RelativePath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			// Create temp directory to simulate user's project
			tempPath = gfile.Temp(guid.S())
			// Copy testdata to temp directory
			srcTestData = gtest.DataPath("issue", "4387")
		)

		// Setup: create temp project structure
		err := gfile.CopyDir(srcTestData, tempPath)
		t.AssertNil(err)
		defer gfile.Remove(tempPath)

		// Save original working directory
		originalWd := gfile.Pwd()

		// Change to temp directory (simulate user being in project root)
		err = gfile.Chdir(tempPath)
		t.AssertNil(err)
		defer gfile.Chdir(originalWd) // Restore original working directory

		// Run gen enums with relative paths
		var (
			srcFolder  = "api"
			outputPath = filepath.FromSlash("internal/packed/packed_enums.go")
			in         = genenums.CGenEnumsInput{
				Src:  srcFolder,
				Path: outputPath,
			}
		)
		err = gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		_, err = genenums.CGenEnums{}.Enums(ctx, in)
		t.AssertNil(err)

		// Expected: file should be created at tempPath/internal/packed/packed_enums.go
		expectedPath := filepath.Join(tempPath, "internal", "packed", "packed_enums.go")
		// Bug: file is created at tempPath/api/internal/packed/packed_enums.go
		wrongPath := filepath.Join(tempPath, "api", "internal", "packed", "packed_enums.go")

		// Assert the file is at the expected location
		t.Assert(gfile.Exists(expectedPath), true)
		// Assert the file is NOT at the wrong location
		t.Assert(gfile.Exists(wrongPath), false)
	})
}

// Test gen enums with absolute output path (should work correctly)
func Test_Gen_Enums_AbsolutePath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			tempPath    = gfile.Temp(guid.S())
			srcTestData = gtest.DataPath("issue", "4387")
		)

		err := gfile.CopyDir(srcTestData, tempPath)
		t.AssertNil(err)
		defer gfile.Remove(tempPath)

		originalWd := gfile.Pwd()
		err = gfile.Chdir(tempPath)
		t.AssertNil(err)
		defer gfile.Chdir(originalWd)

		// Use absolute path for output
		var (
			srcFolder  = "api"
			outputPath = filepath.Join(tempPath, "internal", "packed", "packed_enums.go")
			in         = genenums.CGenEnumsInput{
				Src:  srcFolder,
				Path: outputPath,
			}
		)
		err = gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		_, err = genenums.CGenEnums{}.Enums(ctx, in)
		t.AssertNil(err)

		// Assert the file exists at absolute path
		t.Assert(gfile.Exists(outputPath), true)
	})
}

// Test gen enums in monorepo mode (cd app/xxx/ then run command)
func Test_Gen_Enums_Issue4387_Monorepo(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			// Simulate monorepo structure
			tempPath    = gfile.Temp(guid.S())
			srcTestData = gtest.DataPath("issue", "4387")
			// app/myapp is the subdirectory in monorepo
			appPath = filepath.Join(tempPath, "app", "myapp")
		)

		// Create monorepo structure: tempPath/app/myapp/api/...
		err := gfile.Mkdir(appPath)
		t.AssertNil(err)
		// Copy testdata into app/myapp
		err = gfile.CopyDir(srcTestData, appPath)
		t.AssertNil(err)
		defer gfile.Remove(tempPath)

		originalWd := gfile.Pwd()

		// cd app/myapp (simulate user in monorepo subdirectory)
		err = gfile.Chdir(appPath)
		t.AssertNil(err)
		defer gfile.Chdir(originalWd)

		var (
			srcFolder  = "api"
			outputPath = filepath.FromSlash("internal/packed/packed_enums.go")
			in         = genenums.CGenEnumsInput{
				Src:  srcFolder,
				Path: outputPath,
			}
		)
		err = gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		_, err = genenums.CGenEnums{}.Enums(ctx, in)
		t.AssertNil(err)

		// Expected: file at app/myapp/internal/packed/packed_enums.go
		expectedPath := filepath.Join(appPath, "internal", "packed", "packed_enums.go")
		// Bug: file at app/myapp/api/internal/packed/packed_enums.go
		wrongPath := filepath.Join(appPath, "api", "internal", "packed", "packed_enums.go")

		t.Assert(gfile.Exists(expectedPath), true)
		t.Assert(gfile.Exists(wrongPath), false)
	})
}
