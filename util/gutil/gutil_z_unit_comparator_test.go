// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil_test

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gutil"
)

func TestComparatorString(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorString(1, 1), 0)
		t.Assert(gutil.ComparatorString(1, 2), -1)
		t.Assert(gutil.ComparatorString(2, 1), 1)
	})
}

func TestComparatorInt(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorInt(1, 1), 0)
		t.Assert(gutil.ComparatorInt(1, 2), -1)
		t.Assert(gutil.ComparatorInt(2, 1), 1)
	})
}

func TestComparatorInt8(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorInt8(1, 1), 0)
		t.Assert(gutil.ComparatorInt8(1, 2), -1)
		t.Assert(gutil.ComparatorInt8(2, 1), 1)
	})
}

func TestComparatorInt16(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorInt16(1, 1), 0)
		t.Assert(gutil.ComparatorInt16(1, 2), -1)
		t.Assert(gutil.ComparatorInt16(2, 1), 1)
	})
}

func TestComparatorInt32(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorInt32(1, 1), 0)
		t.Assert(gutil.ComparatorInt32(1, 2), -1)
		t.Assert(gutil.ComparatorInt32(2, 1), 1)
	})
}

func TestComparatorInt64(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorInt64(1, 1), 0)
		t.Assert(gutil.ComparatorInt64(1, 2), -1)
		t.Assert(gutil.ComparatorInt64(2, 1), 1)
	})
}

func TestComparatorUint(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorUint(1, 1), 0)
		t.Assert(gutil.ComparatorUint(1, 2), -1)
		t.Assert(gutil.ComparatorUint(2, 1), 1)
	})
}

func TestComparatorUint8(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorUint8(1, 1), 0)
		t.Assert(gutil.ComparatorUint8(2, 6), 252)
		t.Assert(gutil.ComparatorUint8(2, 1), 1)
	})
}

func TestComparatorUint16(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorUint16(1, 1), 0)
		t.Assert(gutil.ComparatorUint16(1, 2), 65535)
		t.Assert(gutil.ComparatorUint16(2, 1), 1)
	})
}

func TestComparatorUint32(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorUint32(1, 1), 0)
		t.Assert(gutil.ComparatorUint32(-1000, 2147483640), 2147482656)
		t.Assert(gutil.ComparatorUint32(2, 1), 1)
	})
}

func TestComparatorUint64(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorUint64(1, 1), 0)
		t.Assert(gutil.ComparatorUint64(1, 2), -1)
		t.Assert(gutil.ComparatorUint64(2, 1), 1)
	})
}

func TestComparatorFloat32(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorFloat32(1, 1), 0)
		t.Assert(gutil.ComparatorFloat32(1, 2), -1)
		t.Assert(gutil.ComparatorFloat32(2, 1), 1)
	})
}

func TestComparatorFloat64(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorFloat64(1, 1), 0)
		t.Assert(gutil.ComparatorFloat64(1, 2), -1)
		t.Assert(gutil.ComparatorFloat64(2, 1), 1)
	})
}

func TestComparatorByte(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorByte(1, 1), 0)
		t.Assert(gutil.ComparatorByte(1, 2), 255)
		t.Assert(gutil.ComparatorByte(2, 1), 1)
	})
}

func TestComparatorRune(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorRune(1, 1), 0)
		t.Assert(gutil.ComparatorRune(1, 2), -1)
		t.Assert(gutil.ComparatorRune(2, 1), 1)
	})
}

func TestComparatorTime(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		j := gutil.ComparatorTime("2019-06-14", "2019-06-14")
		t.Assert(j, 0)

		k := gutil.ComparatorTime("2019-06-15", "2019-06-14")
		t.Assert(k, 1)

		l := gutil.ComparatorTime("2019-06-13", "2019-06-14")
		t.Assert(l, -1)
	})
}

func TestComparatorFloat32OfFixed(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorFloat32(0.1, 0.1), 0)
		t.Assert(gutil.ComparatorFloat32(1.1, 2.1), -1)
		t.Assert(gutil.ComparatorFloat32(2.1, 1.1), 1)
	})
}

func TestComparatorFloat64OfFixed(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorFloat64(0.1, 0.1), 0)
		t.Assert(gutil.ComparatorFloat64(1.1, 2.1), -1)
		t.Assert(gutil.ComparatorFloat64(2.1, 1.1), 1)
	})
}
