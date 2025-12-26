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

func Test_ReplaceFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			fileName = "/testfile_replace_" + gconv.String(gtime.TimestampNano()) + ".txt"
			content  = "hello world"
		)
		createTestFile(fileName, content)
		defer delTestFiles(fileName)

		// Test basic replacement
		err := gfile.ReplaceFile("world", "gf", testpath()+fileName)
		t.AssertNil(err)
		t.Assert(gfile.GetContents(testpath()+fileName), "hello gf")

		// Test replacement with non-existent search string
		err = gfile.ReplaceFile("notexist", "replaced", testpath()+fileName)
		t.AssertNil(err)
		t.Assert(gfile.GetContents(testpath()+fileName), "hello gf")

		// Test multiple occurrences replacement
		err = gfile.PutContents(testpath()+fileName, "hello hello hello")
		t.AssertNil(err)
		err = gfile.ReplaceFile("hello", "hi", testpath()+fileName)
		t.AssertNil(err)
		t.Assert(gfile.GetContents(testpath()+fileName), "hi hi hi")
	})
}

func Test_ReplaceFileFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			fileName = "/testfile_replacefunc_" + gconv.String(gtime.TimestampNano()) + ".txt"
			content  = "hello world"
		)
		createTestFile(fileName, content)
		defer delTestFiles(fileName)

		// Test replacement with callback function
		err := gfile.ReplaceFileFunc(func(path, content string) string {
			t.Assert(gfile.Basename(path), gfile.Basename(fileName))
			return content + " - modified"
		}, testpath()+fileName)
		t.AssertNil(err)
		t.Assert(gfile.GetContents(testpath()+fileName), "hello world - modified")
	})

	// Test when callback returns same content (no write should happen)
	gtest.C(t, func(t *gtest.T) {
		var (
			fileName = "/testfile_replacefunc2_" + gconv.String(gtime.TimestampNano()) + ".txt"
			content  = "unchanged content"
		)
		createTestFile(fileName, content)
		defer delTestFiles(fileName)

		err := gfile.ReplaceFileFunc(func(path, content string) string {
			return content // Return same content
		}, testpath()+fileName)
		t.AssertNil(err)
		t.Assert(gfile.GetContents(testpath()+fileName), "unchanged content")
	})

	// Test callback with path parameter
	gtest.C(t, func(t *gtest.T) {
		var (
			fileName = "/testfile_replacefunc3_" + gconv.String(gtime.TimestampNano()) + ".txt"
			content  = "test content"
		)
		createTestFile(fileName, content)
		defer delTestFiles(fileName)

		var receivedPath string
		err := gfile.ReplaceFileFunc(func(path, content string) string {
			receivedPath = path
			return "new content"
		}, testpath()+fileName)
		t.AssertNil(err)
		t.Assert(receivedPath, testpath()+fileName)
		t.Assert(gfile.GetContents(testpath()+fileName), "new content")
	})
}

func Test_ReplaceDir(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			dirName  = "/testdir_replace_" + gconv.String(gtime.TimestampNano())
			fileName = dirName + "/test.txt"
			content  = "hello world"
		)
		createDir(dirName)
		createTestFile(fileName, content)
		defer delTestFiles(dirName)

		// Test directory replacement with pattern
		err := gfile.ReplaceDir("world", "gf", testpath()+dirName, "*.txt")
		t.AssertNil(err)
		t.Assert(gfile.GetContents(testpath()+fileName), "hello gf")
	})

	// Test recursive replacement
	gtest.C(t, func(t *gtest.T) {
		var (
			dirName     = "/testdir_replace_recursive_" + gconv.String(gtime.TimestampNano())
			subDirName  = dirName + "/subdir"
			fileName1   = dirName + "/test1.txt"
			fileName2   = subDirName + "/test2.txt"
			content     = "hello world"
		)
		createDir(dirName)
		createDir(subDirName)
		createTestFile(fileName1, content)
		createTestFile(fileName2, content)
		defer delTestFiles(dirName)

		// Non-recursive replacement
		err := gfile.ReplaceDir("world", "gf", testpath()+dirName, "*.txt", false)
		t.AssertNil(err)
		t.Assert(gfile.GetContents(testpath()+fileName1), "hello gf")
		t.Assert(gfile.GetContents(testpath()+fileName2), "hello world") // Should not be changed

		// Reset content
		err = gfile.PutContents(testpath()+fileName1, content)
		t.AssertNil(err)

		// Recursive replacement
		err = gfile.ReplaceDir("world", "gf", testpath()+dirName, "*.txt", true)
		t.AssertNil(err)
		t.Assert(gfile.GetContents(testpath()+fileName1), "hello gf")
		t.Assert(gfile.GetContents(testpath()+fileName2), "hello gf")
	})

	// Test with pattern matching
	gtest.C(t, func(t *gtest.T) {
		var (
			dirName   = "/testdir_replace_pattern_" + gconv.String(gtime.TimestampNano())
			fileName1 = dirName + "/test.txt"
			fileName2 = dirName + "/test.log"
			content   = "hello world"
		)
		createDir(dirName)
		createTestFile(fileName1, content)
		createTestFile(fileName2, content)
		defer delTestFiles(dirName)

		// Only replace in .txt files
		err := gfile.ReplaceDir("world", "gf", testpath()+dirName, "*.txt")
		t.AssertNil(err)
		t.Assert(gfile.GetContents(testpath()+fileName1), "hello gf")
		t.Assert(gfile.GetContents(testpath()+fileName2), "hello world") // .log should not be changed
	})

	// Test with non-existent directory
	gtest.C(t, func(t *gtest.T) {
		err := gfile.ReplaceDir("search", "replace", "/nonexistent_dir_12345", "*.txt")
		t.AssertNE(err, nil)
	})
}

func Test_ReplaceDirFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			dirName   = "/testdir_replacefunc_" + gconv.String(gtime.TimestampNano())
			fileName1 = dirName + "/test1.txt"
			fileName2 = dirName + "/test2.txt"
			content1  = "content1"
			content2  = "content2"
		)
		createDir(dirName)
		createTestFile(fileName1, content1)
		createTestFile(fileName2, content2)
		defer delTestFiles(dirName)

		// Test directory replacement with callback function
		processedFiles := make(map[string]bool)
		err := gfile.ReplaceDirFunc(func(path, content string) string {
			processedFiles[gfile.Basename(path)] = true
			return content + " - modified"
		}, testpath()+dirName, "*.txt")
		t.AssertNil(err)
		t.Assert(gfile.GetContents(testpath()+fileName1), "content1 - modified")
		t.Assert(gfile.GetContents(testpath()+fileName2), "content2 - modified")
		t.Assert(processedFiles["test1.txt"], true)
		t.Assert(processedFiles["test2.txt"], true)
	})

	// Test recursive replacement with callback
	gtest.C(t, func(t *gtest.T) {
		var (
			dirName    = "/testdir_replacefunc_recursive_" + gconv.String(gtime.TimestampNano())
			subDirName = dirName + "/subdir"
			fileName1  = dirName + "/test1.txt"
			fileName2  = subDirName + "/test2.txt"
			content    = "original"
		)
		createDir(dirName)
		createDir(subDirName)
		createTestFile(fileName1, content)
		createTestFile(fileName2, content)
		defer delTestFiles(dirName)

		// Non-recursive
		err := gfile.ReplaceDirFunc(func(path, content string) string {
			return "changed"
		}, testpath()+dirName, "*.txt", false)
		t.AssertNil(err)
		t.Assert(gfile.GetContents(testpath()+fileName1), "changed")
		t.Assert(gfile.GetContents(testpath()+fileName2), "original") // Should not be changed

		// Reset
		err = gfile.PutContents(testpath()+fileName1, content)
		t.AssertNil(err)

		// Recursive
		err = gfile.ReplaceDirFunc(func(path, content string) string {
			return "changed"
		}, testpath()+dirName, "*.txt", true)
		t.AssertNil(err)
		t.Assert(gfile.GetContents(testpath()+fileName1), "changed")
		t.Assert(gfile.GetContents(testpath()+fileName2), "changed")
	})

	// Test with non-existent directory
	gtest.C(t, func(t *gtest.T) {
		err := gfile.ReplaceDirFunc(func(path, content string) string {
			return content
		}, "/nonexistent_dir_12345", "*.txt")
		t.AssertNE(err, nil)
	})
}
