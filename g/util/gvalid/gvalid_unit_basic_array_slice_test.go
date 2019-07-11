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

func Test_Array(t *testing.T) {
	gtest.Case(t, func() {
		arrayData := [2]int{7, 8}
		msgs := map[string]string{
			"between": "not between",
			"in":      "not in",
		}

		err := gvalid.Check(arrayData, "between:6,10|in:7,8", msgs)
		gtest.Assert(err, nil)

		if e := gvalid.Check(arrayData, "between:9,10|in:7,8", msgs); e != nil {
			gtest.Assert(e.String(), "not between")
		}

		if e := gvalid.Check(arrayData, "between:6,10|in:7", msgs); e != nil {
			gtest.Assert(e.String(), "not in")
		}

		err1 := gvalid.Check(&arrayData, "between:6,10|in:7,8", msgs)
		gtest.Assert(err1, nil)

		if e := gvalid.Check(&arrayData, "between:9,10|in:7,8", msgs); e != nil {
			gtest.Assert(e.String(), "not between")
		}

		if e := gvalid.Check(&arrayData, "between:6,10|in:7", msgs); e != nil {
			gtest.Assert(e.String(), "not in")
		}
	})
}

func Test_Slice(t *testing.T) {
	gtest.Case(t, func() {
		sliceData := [][]string{[]string{"12345678", "12345678"}, []string{"12345678", "12345678"}}

		msgs := map[string]string{
			"length": "length err",
		}

		err := gvalid.Check(sliceData, "length:3,16", msgs)
		gtest.Assert(err, nil)

		if e := gvalid.Check(sliceData, "length:9,16", msgs); e != nil {
			gtest.Assert(e.String(), "length err")
		}

	})
}
