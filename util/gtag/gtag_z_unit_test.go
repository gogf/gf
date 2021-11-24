// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtag_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gtag"
)

func Test_Set_Get(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		k := gtime.TimestampNanoStr()
		v := gtime.TimestampNanoStr()
		gtag.Set(k, v)
		t.Assert(gtag.Get(k), v)
	})
}

func Test_Sets_Get(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		k1 := gtime.TimestampNanoStr()
		k2 := gtime.TimestampNanoStr()
		v1 := gtime.TimestampNanoStr()
		v2 := gtime.TimestampNanoStr()
		gtag.Sets(g.MapStrStr{
			k1: v1,
			k2: v2,
		})
		t.Assert(gtag.Get(k1), v1)
		t.Assert(gtag.Get(k2), v2)
	})
}

func Test_Parse(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			k1      = gtime.TimestampNanoStr()
			k2      = gtime.TimestampNanoStr()
			v1      = gtime.TimestampNanoStr()
			v2      = gtime.TimestampNanoStr()
			content = fmt.Sprintf(`this is {%s} and {%s}`, k1, k2)
			expect  = fmt.Sprintf(`this is %s and %s`, v1, v2)
		)
		gtag.Sets(g.MapStrStr{
			k1: v1,
			k2: v2,
		})
		t.Assert(gtag.Parse(content), expect)
	})
}
