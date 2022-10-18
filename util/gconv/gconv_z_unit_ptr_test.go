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

func Test_Ptr_Functions(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var v interface{} = 1
		t.AssertEQ(gconv.PtrAny(v), &v)
	})
	gtest.C(t, func(t *gtest.T) {
		var v string = "1"
		t.AssertEQ(gconv.PtrString(v), &v)
	})
	gtest.C(t, func(t *gtest.T) {
		var v bool = true
		t.AssertEQ(gconv.PtrBool(v), &v)
	})
	gtest.C(t, func(t *gtest.T) {
		var v int = 1
		t.AssertEQ(gconv.PtrInt(v), &v)
	})
	gtest.C(t, func(t *gtest.T) {
		var v int8 = 1
		t.AssertEQ(gconv.PtrInt8(v), &v)
	})
	gtest.C(t, func(t *gtest.T) {
		var v int16 = 1
		t.AssertEQ(gconv.PtrInt16(v), &v)
	})
	gtest.C(t, func(t *gtest.T) {
		var v int32 = 1
		t.AssertEQ(gconv.PtrInt32(v), &v)
	})
	gtest.C(t, func(t *gtest.T) {
		var v int64 = 1
		t.AssertEQ(gconv.PtrInt64(v), &v)
	})
	gtest.C(t, func(t *gtest.T) {
		var v uint = 1
		t.AssertEQ(gconv.PtrUint(v), &v)
	})
	gtest.C(t, func(t *gtest.T) {
		var v uint8 = 1
		t.AssertEQ(gconv.PtrUint8(v), &v)
	})
	gtest.C(t, func(t *gtest.T) {
		var v uint16 = 1
		t.AssertEQ(gconv.PtrUint16(v), &v)
	})
	gtest.C(t, func(t *gtest.T) {
		var v uint32 = 1
		t.AssertEQ(gconv.PtrUint32(v), &v)
	})
	gtest.C(t, func(t *gtest.T) {
		var v uint64 = 1
		t.AssertEQ(gconv.PtrUint64(v), &v)
	})
	gtest.C(t, func(t *gtest.T) {
		var v float32 = 1.01
		t.AssertEQ(gconv.PtrFloat32(v), &v)
	})
	gtest.C(t, func(t *gtest.T) {
		var v float64 = 1.01
		t.AssertEQ(gconv.PtrFloat64(v), &v)
	})
}
