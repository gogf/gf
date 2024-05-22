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

func TestPtrAny(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var v interface{} = 1
		t.AssertEQ(gconv.PtrAny(v), &v)
	})
}

func TestPtrString(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var v string = "Pluto"
		t.AssertEQ(gconv.PtrString(v), &v)
	})
}

func TestPtrBool(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var v bool = true
		t.AssertEQ(gconv.PtrBool(v), &v)
	})
}

func TestPtrInt(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var v int = 123
		t.AssertEQ(gconv.PtrInt(v), &v)
	})
}

func TestPtrInt8(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var v int8 = 123
		t.AssertEQ(gconv.PtrInt8(v), &v)
	})
}

func TestPtrInt16(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var v int16 = 123
		t.AssertEQ(gconv.PtrInt16(v), &v)
	})
}

func TestPtrInt32(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var v int32 = 123
		t.AssertEQ(gconv.PtrInt32(v), &v)
	})
}

func TestPtrInt64(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var v int64 = 123
		t.AssertEQ(gconv.PtrInt64(v), &v)
	})
}

func TestPtrUint(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var v uint = 123
		t.AssertEQ(gconv.PtrUint(v), &v)
	})
}

func TestPtrUint8(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var v uint8 = 123
		t.AssertEQ(gconv.PtrUint8(v), &v)
	})
}

func TestPtrUint16(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var v uint16 = 123
		t.AssertEQ(gconv.PtrUint16(v), &v)
	})
}

func TestPtrUint32(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var v uint32 = 123
		t.AssertEQ(gconv.PtrUint32(v), &v)
	})
}

func TestPtrUint64(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var v uint64 = 123
		t.AssertEQ(gconv.PtrUint64(v), &v)
	})
}

func TestPtrFloat32(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var v float32 = 123.456
		t.AssertEQ(gconv.PtrFloat32(v), &v)
	})
}

func TestPtrFloat64(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var v float64 = 123.456
		t.AssertEQ(gconv.PtrFloat64(v), &v)
	})
}
