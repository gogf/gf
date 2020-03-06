// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcompress_test

import (
	"testing"

	"github.com/gogf/gf/encoding/gcompress"
	"github.com/gogf/gf/test/gtest"
)

func Test_Gzip_UnGzip(t *testing.T) {
	src := "Hello World!!"

	gzip := []byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0xff,
		0xf2, 0x48, 0xcd, 0xc9, 0xc9,
		0x57, 0x08, 0xcf, 0x2f, 0xca,
		0x49, 0x51, 0x54, 0x04, 0x04,
		0x00, 0x00, 0xff, 0xff, 0x9d,
		0x24, 0xa8, 0xd1, 0x0d, 0x00,
		0x00, 0x00,
	}

	arr := []byte(src)
	data, _ := gcompress.Gzip(arr)
	gtest.Assert(data, gzip)

	data, _ = gcompress.UnGzip(gzip)
	gtest.Assert(data, arr)

	data, _ = gcompress.UnGzip(gzip[1:])
	gtest.Assert(data, nil)
}
