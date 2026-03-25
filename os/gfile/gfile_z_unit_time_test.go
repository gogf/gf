// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile_test

import (
	"os"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_MTime(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {

		var (
			file1   = "/testfile_t1.txt"
			err     error
			fileobj os.FileInfo
		)

		createTestFile(file1, "")
		defer delTestFiles(file1)
		fileobj, err = os.Stat(testpath() + file1)
		t.AssertNil(err)

		t.Assert(gfile.MTime(testpath()+file1), fileobj.ModTime())
		t.Assert(gfile.MTime(""), "")
	})
}

func Test_MTimeMillisecond(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			file1   = "/testfile_t1.txt"
			err     error
			fileobj os.FileInfo
		)

		createTestFile(file1, "")
		defer delTestFiles(file1)
		fileobj, err = os.Stat(testpath() + file1)
		t.AssertNil(err)

		time.Sleep(time.Millisecond * 100)
		t.AssertGE(
			gfile.MTimestampMilli(testpath()+file1),
			fileobj.ModTime().UnixNano()/1000000,
		)
		t.Assert(gfile.MTimestampMilli(""), -1)
	})
}

func Test_MTimestamp(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			file1   = "/testfile_mtimestamp.txt"
			err     error
			fileobj os.FileInfo
		)

		createTestFile(file1, "")
		defer delTestFiles(file1)
		fileobj, err = os.Stat(testpath() + file1)
		t.AssertNil(err)

		// Test MTimestamp returns correct unix timestamp
		timestamp := gfile.MTimestamp(testpath() + file1)
		t.Assert(timestamp, fileobj.ModTime().Unix())
		t.Assert(timestamp > 0, true)

		// Test with non-existent file
		t.Assert(gfile.MTimestamp("/nonexistent_file_12345.txt"), -1)

		// Test with empty path
		t.Assert(gfile.MTimestamp(""), -1)
	})

	// Test MTimestamp with directory
	gtest.C(t, func(t *gtest.T) {
		tempDir := gfile.Temp()
		timestamp := gfile.MTimestamp(tempDir)
		t.Assert(timestamp > 0, true)
	})
}
