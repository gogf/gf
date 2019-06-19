// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"github.com/gogf/gf/g/test/gtest"
	"github.com/gogf/gf/g/util/gconv"
	"testing"
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
	gtest.Case(t, func() {
		gtest.AssertEQ(gconv.String(int(123)), "123")
		gtest.AssertEQ(gconv.String(int(-123)), "-123")
		gtest.AssertEQ(gconv.String(int8(123)), "123")
		gtest.AssertEQ(gconv.String(int8(-123)), "-123")
		gtest.AssertEQ(gconv.String(int16(123)), "123")
		gtest.AssertEQ(gconv.String(int16(-123)), "-123")
		gtest.AssertEQ(gconv.String(int32(123)), "123")
		gtest.AssertEQ(gconv.String(int32(-123)), "-123")
		gtest.AssertEQ(gconv.String(int64(123)), "123")
		gtest.AssertEQ(gconv.String(int64(-123)), "-123")
		gtest.AssertEQ(gconv.String(int64(1552578474888)), "1552578474888")
		gtest.AssertEQ(gconv.String(int64(-1552578474888)), "-1552578474888")

		gtest.AssertEQ(gconv.String(uint(123)), "123")
		gtest.AssertEQ(gconv.String(uint8(123)), "123")
		gtest.AssertEQ(gconv.String(uint16(123)), "123")
		gtest.AssertEQ(gconv.String(uint32(123)), "123")
		gtest.AssertEQ(gconv.String(uint64(155257847488898765)), "155257847488898765")

		gtest.AssertEQ(gconv.String(float32(123.456)), "123.456")
		gtest.AssertEQ(gconv.String(float32(-123.456)), "-123.456")
		gtest.AssertEQ(gconv.String(float64(1552578474888.456)), "1552578474888.456")
		gtest.AssertEQ(gconv.String(float64(-1552578474888.456)), "-1552578474888.456")

		gtest.AssertEQ(gconv.String(true), "true")
		gtest.AssertEQ(gconv.String(false), "false")

		gtest.AssertEQ(gconv.String([]byte("bytes")), "bytes")

		gtest.AssertEQ(gconv.String(stringStruct1{"john"}), `{"Name":"john"}`)
		gtest.AssertEQ(gconv.String(&stringStruct1{"john"}), "john")

		gtest.AssertEQ(gconv.String(stringStruct2{"john"}), `{"Name":"john"}`)
		gtest.AssertEQ(gconv.String(&stringStruct2{"john"}), `{"Name":"john"}`)
	})
}
