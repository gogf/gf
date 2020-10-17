// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gconv"
)

func Test_Time(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t1 := "2011-10-10 01:02:03.456"
		t.AssertEQ(gconv.GTime(t1), gtime.NewFromStr(t1))
		t.AssertEQ(gconv.Time(t1), gtime.NewFromStr(t1).Time)
		t.AssertEQ(gconv.Duration(100), 100*time.Nanosecond)
	})
}
