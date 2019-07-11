// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid_test

import (
	"testing"

	"github.com/gogf/gf/g/test/gtest"
	"github.com/gogf/gf/g/util/gvalid"
)

func Test_Map(t *testing.T) {
	type Test struct {
		Id int
	}
	gtest.Case(t, func() {
		mapData := map[string]interface{}{
			"123": map[string]int{
				"aaa": 6,
				"bbb": 7,
				"ccc": 8,
			},
			"456": &Test{
				Id: 9,
			},
		}

		err := gvalid.Check(mapData, "between:6,10|in:6,7,8,9", nil)
		gtest.Assert(err, nil)

		msgs := map[string]string{
			"between": "not between",
			"in":      "not in",
		}
		if e := gvalid.Check(mapData, "between:10,10|in:6,7,8,9", msgs); e != nil {
			gtest.Assert(e.String(), "not between")
		}
		if e := gvalid.Check(mapData, "between:6,10|in:10", msgs); e != nil {
			gtest.Assert(e.String(), "not in")
		}
	})
}
