// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcompress_test

import (
	"testing"

	"github.com/gogf/gf/v2/encoding/gcompress"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Gzip_UnGzip(t *testing.T) {
	src := "Hello World!!"

	gtest.C(t, func(t *gtest.T) {
		arr := []byte(src)
		data, err := gcompress.Gzip(arr)
		t.AssertNil(err)
		t.AssertGT(len(data), 0)

		result, err := gcompress.UnGzip(data)
		t.AssertNil(err)
		t.Assert(result, arr)

		result, _ = gcompress.UnGzip(data[1:])
		t.Assert(result, nil)
	})
}

func Test_Gzip_UnGzip_File(t *testing.T) {
	var (
		srcPath  = gtest.DataPath("gzip", "file.txt")
		dstPath1 = gfile.Temp(gtime.TimestampNanoStr(), "gzip.zip")
		dstPath2 = gfile.Temp(gtime.TimestampNanoStr(), "file.txt")
	)

	// Compress.
	gtest.C(t, func(t *gtest.T) {
		err := gcompress.GzipFile(srcPath, dstPath1, 9)
		t.AssertNil(err)
		defer gfile.Remove(dstPath1)
		t.Assert(gfile.Exists(dstPath1), true)

		// Decompress.
		err = gcompress.UnGzipFile(dstPath1, dstPath2)
		t.AssertNil(err)
		defer gfile.Remove(dstPath2)
		t.Assert(gfile.Exists(dstPath2), true)

		t.Assert(gfile.GetContents(srcPath), gfile.GetContents(dstPath2))
	})
}
