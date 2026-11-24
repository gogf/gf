// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil_test

import (
	"math"
	"slices"
	"testing"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gutil"
)

func Test_ComparatorString(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorString(1, 1), 0)
		t.Assert(gutil.ComparatorString(1, 2), -1)
		t.Assert(gutil.ComparatorString(2, 1), 1)
	})
}

func Test_ComparatorInt(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorInt(1, 1), 0)
		t.Assert(gutil.ComparatorInt(1, 2), -1)
		t.Assert(gutil.ComparatorInt(2, 1), 1)
	})
}

func Test_ComparatorInt8(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorInt8(1, 1), 0)
		t.Assert(gutil.ComparatorInt8(1, 2), -1)
		t.Assert(gutil.ComparatorInt8(2, 1), 1)
	})
}

func Test_ComparatorInt16(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorInt16(1, 1), 0)
		t.Assert(gutil.ComparatorInt16(1, 2), -1)
		t.Assert(gutil.ComparatorInt16(2, 1), 1)
	})
}

func Test_ComparatorInt32(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorInt32(1, 1), 0)
		t.Assert(gutil.ComparatorInt32(1, 2), -1)
		t.Assert(gutil.ComparatorInt32(2, 1), 1)
	})
}

func Test_ComparatorInt64(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorInt64(1, 1), 0)
		t.Assert(gutil.ComparatorInt64(1, 2), -1)
		t.Assert(gutil.ComparatorInt64(2, 1), 1)
	})
}

func Test_ComparatorUint(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorUint(1, 1), 0)
		t.Assert(gutil.ComparatorUint(1, 2), -1)
		t.Assert(gutil.ComparatorUint(2, 1), 1)
	})
}

func Test_ComparatorUint8(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorUint8(1, 1), 0)
		t.Assert(gutil.ComparatorUint8(2, 6), -1)
		t.Assert(gutil.ComparatorUint8(2, 1), 1)
	})
}

func Test_ComparatorUint16(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorUint16(1, 1), 0)
		t.Assert(gutil.ComparatorUint16(1, 2), -1)
		t.Assert(gutil.ComparatorUint16(2, 1), 1)
	})
}

func Test_ComparatorUint32(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorUint32(1, 1), 0)
		t.Assert(gutil.ComparatorUint32(1, 2), -1)
		t.Assert(gutil.ComparatorUint32(-1000, 2147483640), 1)
		t.Assert(gutil.ComparatorUint32(2, 1), 1)
	})
}

func Test_ComparatorUint64(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorUint64(1, 1), 0)
		t.Assert(gutil.ComparatorUint64(1, 2), -1)
		t.Assert(gutil.ComparatorUint64(2, 1), 1)
	})
}

func Test_ComparatorFloat32(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorFloat32(1, 1), 0)
		t.Assert(gutil.ComparatorFloat32(1, 2), -1)
		t.Assert(gutil.ComparatorFloat32(2, 1), 1)
	})
}

func Test_ComparatorFloat64(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorFloat64(1, 1), 0)
		t.Assert(gutil.ComparatorFloat64(1, 2), -1)
		t.Assert(gutil.ComparatorFloat64(2, 1), 1)
	})
}

func Test_ComparatorByte(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorByte(1, 1), 0)
		t.Assert(gutil.ComparatorByte(1, 2), -1)
		t.Assert(gutil.ComparatorByte(2, 1), 1)
	})
}

func Test_ComparatorRune(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorRune(1, 1), 0)
		t.Assert(gutil.ComparatorRune(1, 2), -1)
		t.Assert(gutil.ComparatorRune(2, 1), 1)
	})
}

// Test_Comparator_Overflow checks that the integer comparators keep the sign of the comparison
// when the difference of both values does not fit into the compared type.
func Test_Comparator_Overflow(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorInt(math.MaxInt, -1), 1)
		t.Assert(gutil.ComparatorInt8(100, -100), 1)
		t.Assert(gutil.ComparatorInt8(-100, 100), -1)
		t.Assert(gutil.ComparatorInt16(30000, -30000), 1)
		t.Assert(gutil.ComparatorInt32(math.MaxInt32, -1), 1)
		t.Assert(gutil.ComparatorInt64(int64(math.MaxInt64), -1), 1)
		t.Assert(gutil.ComparatorUint(0, 1), -1)
		t.Assert(gutil.ComparatorUint32(0, uint32(math.MaxUint32)), -1)
		t.Assert(gutil.ComparatorUint64(uint64(math.MaxUint64), 0), 1)
		t.Assert(gutil.ComparatorByte(200, 100), 1)
		t.Assert(gutil.ComparatorRune(math.MinInt32, 1), -1)

		// A sorted array relies on the sign of the comparator result.
		array := garray.NewSortedArray(gutil.ComparatorUint8)
		array.Add(2, 6, 1, 200)
		t.Assert(array.Slice(), []any{1, 2, 6, 200})
	})
}

func Test_ComparatorTime(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		j := gutil.ComparatorTime("2019-06-14", "2019-06-14")
		t.Assert(j, 0)

		k := gutil.ComparatorTime("2019-06-15", "2019-06-14")
		t.Assert(k, 1)

		l := gutil.ComparatorTime("2019-06-13", "2019-06-14")
		t.Assert(l, -1)
	})
}

func Test_ComparatorFloat32OfFixed(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorFloat32(0.1, 0.1), 0)
		t.Assert(gutil.ComparatorFloat32(1.1, 2.1), -1)
		t.Assert(gutil.ComparatorFloat32(2.1, 1.1), 1)
	})
}

func Test_ComparatorFloat64OfFixed(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorFloat64(0.1, 0.1), 0)
		t.Assert(gutil.ComparatorFloat64(1.1, 2.1), -1)
		t.Assert(gutil.ComparatorFloat64(2.1, 1.1), 1)
	})
}

func Test_ComparatorIntOverflow(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorInt(math.MaxInt, math.MinInt), 1)
		t.Assert(gutil.ComparatorInt(math.MinInt, math.MaxInt), -1)
		t.Assert(gutil.ComparatorInt8(math.MinInt8, math.MaxInt8), -1)
		t.Assert(gutil.ComparatorInt16(math.MinInt16, math.MaxInt16), -1)
		t.Assert(gutil.ComparatorInt32(math.MinInt32, math.MaxInt32), -1)
		t.Assert(gutil.ComparatorInt64(int64(math.MinInt64), int64(math.MaxInt64)), -1)
		t.Assert(gutil.ComparatorInt64(int64(math.MaxInt64), int64(math.MinInt64)), 1)
	})
}

func Test_ComparatorUint64Overflow(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorUint64(uint64(math.MaxUint64), 1), 1)
		t.Assert(gutil.ComparatorUint64(1, uint64(math.MaxUint64)), -1)
	})
}

func Test_ComparatorFloatNaN(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.ComparatorFloat32(float32(math.NaN()), 1), -1)
		t.Assert(gutil.ComparatorFloat32(1, float32(math.NaN())), 1)
		t.Assert(gutil.ComparatorFloat32(float32(math.NaN()), float32(math.NaN())), 0)

		t.Assert(gutil.ComparatorFloat64(math.NaN(), 1), -1)
		t.Assert(gutil.ComparatorFloat64(1, math.NaN()), 1)
		t.Assert(gutil.ComparatorFloat64(math.NaN(), math.NaN()), 0)
	})
}

func Test_ComparatorInt64SortExtremes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		values := []int64{math.MaxInt64, 0, math.MinInt64, 1, -1}
		slices.SortFunc(values, func(a, b int64) int {
			return gutil.ComparatorInt64(a, b)
		})
		t.Assert(values, []int64{math.MinInt64, -1, 0, 1, math.MaxInt64})
	})
}
