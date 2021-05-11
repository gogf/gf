// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid_test

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gvalid"
	"testing"
)

func TestValidator_CheckStructWithParamMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type UserApiSearch struct {
			Uid      int64  `json:"uid" v:"required"`
			Nickname string `json:"nickname" v:"required-with:Uid"`
		}
		data := UserApiSearch{}
		t.Assert(gvalid.CheckStructWithParamMap(data, g.Map{}, nil), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		type UserApiSearch struct {
			Uid       int64       `json:"uid"`
			Nickname  string      `json:"nickname" v:"required-with:Uid"`
			StartTime *gtime.Time `json:"start_time" v:"required-with:EndTime"`
			EndTime   *gtime.Time `json:"end_time" v:"required-with:StartTime"`
		}
		data := UserApiSearch{
			StartTime: nil,
			EndTime:   nil,
		}
		t.Assert(gvalid.CheckStructWithParamMap(data, g.Map{}, nil), nil)
	})
}
