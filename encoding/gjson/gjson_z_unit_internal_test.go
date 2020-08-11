// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson

import (
	"github.com/gogf/gf/test/gtest"
	"testing"
)

func Test_checkDataType(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		data := []byte(`
bb           = """
                   dig := dig;                         END;"""
`)
		t.Assert(checkDataType(data), "toml")
	})

	gtest.C(t, func(t *gtest.T) {
		data := []byte(`
# 模板引擎目录
viewpath = "/home/www/templates/"
# MySQL数据库配置
[redis]
dd = 11
[redis]
    disk  = "127.0.0.1:6379,0"
    cache = "127.0.0.1:6379,1"
`)
		t.Assert(checkDataType(data), "toml")
	})
}
