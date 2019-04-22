// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
package gcompress_test

import (
	"github.com/gogf/gf/g/encoding/gcompress"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func TestZlib(t *testing.T) {
	gtest.Case(t, func() {
		src := "hello, world\n"
		dst := []byte{120, 156, 202, 72, 205, 201, 201, 215, 81, 40, 207, 47, 202, 73, 225, 2, 4, 0, 0, 255, 255, 33, 231, 4, 147}
		gtest.Assert(gcompress.Zlib([]byte(src)), dst)

		gtest.Assert(gcompress.UnZlib(dst), []byte(src))
	})

}

func TestGzip(t *testing.T) {
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

	gtest.Assert(gcompress.Gzip([]byte(src)), gzip)

	gtest.Assert(gcompress.UnGzip(gzip), []byte(src))
}
