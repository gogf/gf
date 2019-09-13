// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/util/gconv"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"

	"github.com/gogf/gf/test/gtest"
)

func Test_ConfigFromMap(t *testing.T) {
	gtest.Case(t, func() {
		m := g.Map{
			"addr":            ":8199",
			"readTimeout":     "60s",
			"indexFiles":      g.Slice{"index.php", "main.php"},
			"errorLogEnabled": true,
			"cookieMaxAge":    "1y",
		}
		config := ghttp.ConfigFromMap(m)
		d1, _ := time.ParseDuration(gconv.String(m["readTimeout"]))
		d2, _ := time.ParseDuration(gconv.String(m["cookieMaxAge"]))
		gtest.Assert(config.Addr, m["addr"])
		gtest.Assert(config.ReadTimeout, d1)
		gtest.Assert(config.CookieMaxAge, d2)
		gtest.Assert(config.IndexFiles, m["indexFiles"])
		gtest.Assert(config.ErrorLogEnabled, m["errorLogEnabled"])
	})
}
