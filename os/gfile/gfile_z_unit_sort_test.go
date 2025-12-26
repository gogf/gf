// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile_test

import (
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func Test_SortFiles(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			dirName   = "/testdir_sort_" + gconv.String(gtime.TimestampNano())
			fileName1 = dirName + "/b.txt"
			fileName2 = dirName + "/a.txt"
			subDir1   = dirName + "/subdir_b"
			subDir2   = dirName + "/subdir_a"
		)
		createDir(dirName)
		createDir(subDir1)
		createDir(subDir2)
		createTestFile(fileName1, "")
		createTestFile(fileName2, "")
		defer delTestFiles(dirName)

		// Test sorting: directories should come before files, then sorted alphabetically
		files := []string{
			testpath() + fileName1,
			testpath() + fileName2,
			testpath() + subDir1,
			testpath() + subDir2,
		}
		sorted := gfile.SortFiles(files)

		// Directories should come first, sorted alphabetically
		t.Assert(sorted[0], testpath()+subDir2) // subdir_a
		t.Assert(sorted[1], testpath()+subDir1) // subdir_b
		// Files should come after, sorted alphabetically
		t.Assert(sorted[2], testpath()+fileName2) // a.txt
		t.Assert(sorted[3], testpath()+fileName1) // b.txt
	})

	// Test with only files
	gtest.C(t, func(t *gtest.T) {
		var (
			dirName   = "/testdir_sort_files_" + gconv.String(gtime.TimestampNano())
			fileName1 = dirName + "/c.txt"
			fileName2 = dirName + "/a.txt"
			fileName3 = dirName + "/b.txt"
		)
		createDir(dirName)
		createTestFile(fileName1, "")
		createTestFile(fileName2, "")
		createTestFile(fileName3, "")
		defer delTestFiles(dirName)

		files := []string{
			testpath() + fileName1,
			testpath() + fileName2,
			testpath() + fileName3,
		}
		sorted := gfile.SortFiles(files)

		t.Assert(sorted[0], testpath()+fileName2) // a.txt
		t.Assert(sorted[1], testpath()+fileName3) // b.txt
		t.Assert(sorted[2], testpath()+fileName1) // c.txt
	})

	// Test with only directories
	gtest.C(t, func(t *gtest.T) {
		var (
			dirName = "/testdir_sort_dirs_" + gconv.String(gtime.TimestampNano())
			subDir1 = dirName + "/c_dir"
			subDir2 = dirName + "/a_dir"
			subDir3 = dirName + "/b_dir"
		)
		createDir(dirName)
		createDir(subDir1)
		createDir(subDir2)
		createDir(subDir3)
		defer delTestFiles(dirName)

		files := []string{
			testpath() + subDir1,
			testpath() + subDir2,
			testpath() + subDir3,
		}
		sorted := gfile.SortFiles(files)

		t.Assert(sorted[0], testpath()+subDir2) // a_dir
		t.Assert(sorted[1], testpath()+subDir3) // b_dir
		t.Assert(sorted[2], testpath()+subDir1) // c_dir
	})

	// Test with empty slice
	gtest.C(t, func(t *gtest.T) {
		files := []string{}
		sorted := gfile.SortFiles(files)
		t.Assert(len(sorted), 0)
	})

	// Test with single element
	gtest.C(t, func(t *gtest.T) {
		var (
			dirName  = "/testdir_sort_single_" + gconv.String(gtime.TimestampNano())
			fileName = dirName + "/single.txt"
		)
		createDir(dirName)
		createTestFile(fileName, "")
		defer delTestFiles(dirName)

		files := []string{testpath() + fileName}
		sorted := gfile.SortFiles(files)

		t.Assert(len(sorted), 1)
		t.Assert(sorted[0], testpath()+fileName)
	})

	// Test with mixed paths (some may not exist - testing sort behavior)
	gtest.C(t, func(t *gtest.T) {
		var (
			dirName  = "/testdir_sort_mixed_" + gconv.String(gtime.TimestampNano())
			fileName = dirName + "/existing.txt"
			subDir   = dirName + "/existing_dir"
		)
		createDir(dirName)
		createDir(subDir)
		createTestFile(fileName, "")
		defer delTestFiles(dirName)

		// Mix of existing dir, existing file
		files := []string{
			testpath() + fileName,
			testpath() + subDir,
		}
		sorted := gfile.SortFiles(files)

		// Directory should come first
		t.Assert(sorted[0], testpath()+subDir)
		t.Assert(sorted[1], testpath()+fileName)
	})
}
