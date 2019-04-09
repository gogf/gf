// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
package gcompress_test

import (
	"github.com/gogf/gf/g/util/grand"
	"github.com/gogf/gf/g/encoding/gcompress"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

var times = 10
var length = 2000

func TestCompress(t *testing.T) {
	for i := 0; i < times; i++ {
		src := grand.RandStr(length + i * 10)
		zlibVal := gcompress.Zlib([]byte(src))
		dst := gcompress.UnZlib(zlibVal)
		gtest.Assert(dst, []byte(src))

		dst1 := gcompress.UnGzip(gcompress.Gzip([]byte(src)))
		gtest.Assert(dst1, []byte(src))
	}
}
