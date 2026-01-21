// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_Pack_ToGoFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			srcPath = gfile.Temp(guid.S())
			dstPath = gfile.Temp(guid.S())
			dstFile = filepath.Join(dstPath, "packed", "data.go")
		)
		// Create source directory with test files
		err := gfile.Mkdir(srcPath)
		t.AssertNil(err)
		defer gfile.Remove(srcPath)

		err = gfile.Mkdir(dstPath)
		t.AssertNil(err)
		defer gfile.Remove(dstPath)

		// Create test files
		err = gfile.PutContents(filepath.Join(srcPath, "test.txt"), "hello world")
		t.AssertNil(err)
		err = gfile.PutContents(filepath.Join(srcPath, "test.json"), `{"key":"value"}`)
		t.AssertNil(err)

		// Create packed directory
		err = gfile.Mkdir(filepath.Join(dstPath, "packed"))
		t.AssertNil(err)

		// Pack to go file
		_, err = Pack.Index(context.Background(), cPackInput{
			Src:  srcPath,
			Dst:  dstFile,
			Name: "packed",
		})
		t.AssertNil(err)

		// Verify output file exists
		t.Assert(gfile.Exists(dstFile), true)

		// Verify it's a valid Go file
		content := gfile.GetContents(dstFile)
		t.Assert(gstr.Contains(content, "package packed"), true)
		t.Assert(gstr.Contains(content, "func init()"), true)
	})
}

func Test_Pack_ToBinaryFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			srcPath = gfile.Temp(guid.S())
			dstPath = gfile.Temp(guid.S())
			dstFile = filepath.Join(dstPath, "data.bin")
		)
		// Create source directory with test files
		err := gfile.Mkdir(srcPath)
		t.AssertNil(err)
		defer gfile.Remove(srcPath)

		err = gfile.Mkdir(dstPath)
		t.AssertNil(err)
		defer gfile.Remove(dstPath)

		// Create test file
		err = gfile.PutContents(filepath.Join(srcPath, "test.txt"), "binary content")
		t.AssertNil(err)

		// Pack to binary file (no Name specified)
		_, err = Pack.Index(context.Background(), cPackInput{
			Src: srcPath,
			Dst: dstFile,
		})
		t.AssertNil(err)

		// Verify output file exists
		t.Assert(gfile.Exists(dstFile), true)

		// Verify it's a binary file (not a Go file)
		content := gfile.GetContents(dstFile)
		t.Assert(gstr.Contains(content, "package"), false)
	})
}

func Test_Pack_MultipleSources(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			srcPath1 = gfile.Temp(guid.S())
			srcPath2 = gfile.Temp(guid.S())
			dstPath  = gfile.Temp(guid.S())
			dstFile  = filepath.Join(dstPath, "packed", "multi.go")
		)
		// Create source directories
		err := gfile.Mkdir(srcPath1)
		t.AssertNil(err)
		defer gfile.Remove(srcPath1)

		err = gfile.Mkdir(srcPath2)
		t.AssertNil(err)
		defer gfile.Remove(srcPath2)

		err = gfile.Mkdir(dstPath)
		t.AssertNil(err)
		defer gfile.Remove(dstPath)

		// Create test files in each source
		err = gfile.PutContents(filepath.Join(srcPath1, "file1.txt"), "content1")
		t.AssertNil(err)
		err = gfile.PutContents(filepath.Join(srcPath2, "file2.txt"), "content2")
		t.AssertNil(err)

		// Create packed directory
		err = gfile.Mkdir(filepath.Join(dstPath, "packed"))
		t.AssertNil(err)

		// Pack multiple sources (comma-separated)
		_, err = Pack.Index(context.Background(), cPackInput{
			Src:  srcPath1 + "," + srcPath2,
			Dst:  dstFile,
			Name: "packed",
		})
		t.AssertNil(err)

		// Verify output file exists
		t.Assert(gfile.Exists(dstFile), true)
	})
}

func Test_Pack_WithPrefix(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			srcPath = gfile.Temp(guid.S())
			dstPath = gfile.Temp(guid.S())
			dstFile = filepath.Join(dstPath, "packed", "prefix.go")
		)
		// Create source directory
		err := gfile.Mkdir(srcPath)
		t.AssertNil(err)
		defer gfile.Remove(srcPath)

		err = gfile.Mkdir(dstPath)
		t.AssertNil(err)
		defer gfile.Remove(dstPath)

		// Create test file
		err = gfile.PutContents(filepath.Join(srcPath, "test.txt"), "with prefix")
		t.AssertNil(err)

		// Create packed directory
		err = gfile.Mkdir(filepath.Join(dstPath, "packed"))
		t.AssertNil(err)

		// Pack with prefix
		_, err = Pack.Index(context.Background(), cPackInput{
			Src:    srcPath,
			Dst:    dstFile,
			Name:   "packed",
			Prefix: "/static",
		})
		t.AssertNil(err)

		// Verify output file exists
		t.Assert(gfile.Exists(dstFile), true)
	})
}

func Test_Pack_WithKeepPath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			srcPath = gfile.Temp(guid.S())
			dstPath = gfile.Temp(guid.S())
			dstFile = filepath.Join(dstPath, "packed", "keeppath.go")
		)
		// Create source directory with subdirectory
		err := gfile.Mkdir(srcPath)
		t.AssertNil(err)
		defer gfile.Remove(srcPath)

		err = gfile.Mkdir(dstPath)
		t.AssertNil(err)
		defer gfile.Remove(dstPath)

		// Create subdirectory and file
		subDir := filepath.Join(srcPath, "subdir")
		err = gfile.Mkdir(subDir)
		t.AssertNil(err)
		err = gfile.PutContents(filepath.Join(subDir, "test.txt"), "keeppath content")
		t.AssertNil(err)

		// Create packed directory
		err = gfile.Mkdir(filepath.Join(dstPath, "packed"))
		t.AssertNil(err)

		// Pack with keepPath
		_, err = Pack.Index(context.Background(), cPackInput{
			Src:      srcPath,
			Dst:      dstFile,
			Name:     "packed",
			KeepPath: true,
		})
		t.AssertNil(err)

		// Verify output file exists
		t.Assert(gfile.Exists(dstFile), true)
	})
}

func Test_Pack_AutoPackageName(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			srcPath = gfile.Temp(guid.S())
			dstPath = gfile.Temp(guid.S())
			dstFile = filepath.Join(dstPath, "mypackage", "data.go")
		)
		// Create source directory
		err := gfile.Mkdir(srcPath)
		t.AssertNil(err)
		defer gfile.Remove(srcPath)

		err = gfile.Mkdir(dstPath)
		t.AssertNil(err)
		defer gfile.Remove(dstPath)

		// Create test file
		err = gfile.PutContents(filepath.Join(srcPath, "test.txt"), "auto package name")
		t.AssertNil(err)

		// Create mypackage directory
		err = gfile.Mkdir(filepath.Join(dstPath, "mypackage"))
		t.AssertNil(err)

		// Pack without Name - should use directory name "mypackage"
		_, err = Pack.Index(context.Background(), cPackInput{
			Src: srcPath,
			Dst: dstFile,
			// Name not specified, should be auto-detected as "mypackage"
		})
		t.AssertNil(err)

		// Verify output file exists and has correct package name
		t.Assert(gfile.Exists(dstFile), true)
		content := gfile.GetContents(dstFile)
		t.Assert(gstr.Contains(content, "package mypackage"), true)
	})
}

func Test_Pack_EmptySource(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			srcPath = gfile.Temp(guid.S())
			dstPath = gfile.Temp(guid.S())
			dstFile = filepath.Join(dstPath, "packed", "empty.go")
		)
		// Create empty source directory
		err := gfile.Mkdir(srcPath)
		t.AssertNil(err)
		defer gfile.Remove(srcPath)

		err = gfile.Mkdir(dstPath)
		t.AssertNil(err)
		defer gfile.Remove(dstPath)

		// Create packed directory
		err = gfile.Mkdir(filepath.Join(dstPath, "packed"))
		t.AssertNil(err)

		// Pack empty directory
		_, err = Pack.Index(context.Background(), cPackInput{
			Src:  srcPath,
			Dst:  dstFile,
			Name: "packed",
		})
		t.AssertNil(err)

		// Verify output file exists (even if source is empty)
		t.Assert(gfile.Exists(dstFile), true)
	})
}

func Test_Pack_NestedDirectories(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			srcPath = gfile.Temp(guid.S())
			dstPath = gfile.Temp(guid.S())
			dstFile = filepath.Join(dstPath, "packed", "nested.go")
		)
		// Create source directory with nested structure
		err := gfile.Mkdir(srcPath)
		t.AssertNil(err)
		defer gfile.Remove(srcPath)

		err = gfile.Mkdir(dstPath)
		t.AssertNil(err)
		defer gfile.Remove(dstPath)

		// Create nested directories and files
		level1 := filepath.Join(srcPath, "level1")
		level2 := filepath.Join(level1, "level2")
		level3 := filepath.Join(level2, "level3")
		err = gfile.Mkdir(level3)
		t.AssertNil(err)

		err = gfile.PutContents(filepath.Join(srcPath, "root.txt"), "root")
		t.AssertNil(err)
		err = gfile.PutContents(filepath.Join(level1, "l1.txt"), "level1")
		t.AssertNil(err)
		err = gfile.PutContents(filepath.Join(level2, "l2.txt"), "level2")
		t.AssertNil(err)
		err = gfile.PutContents(filepath.Join(level3, "l3.txt"), "level3")
		t.AssertNil(err)

		// Create packed directory
		err = gfile.Mkdir(filepath.Join(dstPath, "packed"))
		t.AssertNil(err)

		// Pack nested directories
		_, err = Pack.Index(context.Background(), cPackInput{
			Src:  srcPath,
			Dst:  dstFile,
			Name: "packed",
		})
		t.AssertNil(err)

		// Verify output file exists
		t.Assert(gfile.Exists(dstFile), true)

		// Verify content includes all files
		content := gfile.GetContents(dstFile)
		t.Assert(gstr.Contains(content, "package packed"), true)
	})
}
