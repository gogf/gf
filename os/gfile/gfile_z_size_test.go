// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile_test

import (
	"testing"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/test/gtest"
)

func Test_Size(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			paths1 string = "/testfile_t1.txt"
			sizes  int64
		)

		createTestFile(paths1, "abcdefghijklmn")
		defer delTestFiles(paths1)

		sizes = gfile.Size(testpath() + paths1)
		t.Assert(sizes, 14)

		sizes = gfile.Size("")
		t.Assert(sizes, 0)

	})
}

func Test_FormatSize(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gfile.FormatSize(0), "0.00B")
		t.Assert(gfile.FormatSize(16), "16.00B")

		t.Assert(gfile.FormatSize(1024), "1.00K")

		t.Assert(gfile.FormatSize(16000000), "15.26M")

		t.Assert(gfile.FormatSize(1600000000), "1.49G")

		t.Assert(gfile.FormatSize(9600000000000), "8.73T")
		t.Assert(gfile.FormatSize(9600000000000000), "8.53P")
	})
}

func Test_ReadableSize(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {

		var (
			paths1 string = "/testfile_t1.txt"
		)
		createTestFile(paths1, "abcdefghijklmn")
		defer delTestFiles(paths1)
		t.Assert(gfile.ReadableSize(testpath()+paths1), "14.00B")
		t.Assert(gfile.ReadableSize(""), "0.00B")

	})
}
