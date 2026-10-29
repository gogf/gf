// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcompress_test

import (
	"testing"

	"github.com/gogf/gf/v2/encoding/gcompress"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Zlib_UnZlib(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		src := "hello, world\n"
		compressed, err := gcompress.Zlib([]byte(src))
		t.AssertNil(err)
		t.AssertGT(len(compressed), 0)

		result, err := gcompress.UnZlib(compressed)
		t.AssertNil(err)
		t.Assert(result, []byte(src))

		data, _ := gcompress.Zlib(nil)
		t.Assert(data, nil)
		data, _ = gcompress.UnZlib(nil)
		t.Assert(data, nil)

		data, _ = gcompress.UnZlib(compressed[1:])
		t.Assert(data, nil)
	})
}
