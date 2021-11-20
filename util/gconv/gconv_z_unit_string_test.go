// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

type stringStruct1 struct {
	Name string
}

type stringStruct2 struct {
	Name string
}

func (s *stringStruct1) String() string {
	return s.Name
}

func Test_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.String(int(123)), "123")
		t.AssertEQ(gconv.String(int(-123)), "-123")
		t.AssertEQ(gconv.String(int8(123)), "123")
		t.AssertEQ(gconv.String(int8(-123)), "-123")
		t.AssertEQ(gconv.String(int16(123)), "123")
		t.AssertEQ(gconv.String(int16(-123)), "-123")
		t.AssertEQ(gconv.String(int32(123)), "123")
		t.AssertEQ(gconv.String(int32(-123)), "-123")
		t.AssertEQ(gconv.String(int64(123)), "123")
		t.AssertEQ(gconv.String(int64(-123)), "-123")
		t.AssertEQ(gconv.String(int64(1552578474888)), "1552578474888")
		t.AssertEQ(gconv.String(int64(-1552578474888)), "-1552578474888")

		t.AssertEQ(gconv.String(uint(123)), "123")
		t.AssertEQ(gconv.String(uint8(123)), "123")
		t.AssertEQ(gconv.String(uint16(123)), "123")
		t.AssertEQ(gconv.String(uint32(123)), "123")
		t.AssertEQ(gconv.String(uint64(155257847488898765)), "155257847488898765")

		t.AssertEQ(gconv.String(float32(123.456)), "123.456")
		t.AssertEQ(gconv.String(float32(-123.456)), "-123.456")
		t.AssertEQ(gconv.String(float64(1552578474888.456)), "1552578474888.456")
		t.AssertEQ(gconv.String(float64(-1552578474888.456)), "-1552578474888.456")

		t.AssertEQ(gconv.String(true), "true")
		t.AssertEQ(gconv.String(false), "false")

		t.AssertEQ(gconv.String([]byte("bytes")), "bytes")

		t.AssertEQ(gconv.String(stringStruct1{"john"}), `{"Name":"john"}`)
		t.AssertEQ(gconv.String(&stringStruct1{"john"}), "john")

		t.AssertEQ(gconv.String(stringStruct2{"john"}), `{"Name":"john"}`)
		t.AssertEQ(gconv.String(&stringStruct2{"john"}), `{"Name":"john"}`)
	})
}
