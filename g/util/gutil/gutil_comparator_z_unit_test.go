package gutil_test

import (
	"testing"

	"github.com/gogf/gf/g/test/gtest"
	"github.com/gogf/gf/g/util/gutil"
)

func Test_ComparatorString(t *testing.T) {
	j := gutil.ComparatorString(1, 1)
	gtest.Assert(j, 0)
}

func Test_ComparatorInt(t *testing.T) {
	j := gutil.ComparatorInt(1, 1)
	gtest.Assert(j, 0)
}

func Test_ComparatorInt8(t *testing.T) {
	j := gutil.ComparatorInt8(1, 1)
	gtest.Assert(j, 0)
}

func Test_ComparatorInt16(t *testing.T) {
	j := gutil.ComparatorInt16(1, 1)
	gtest.Assert(j, 0)
}

func Test_ComparatorInt32(t *testing.T) {
	j := gutil.ComparatorInt32(1, 1)
	gtest.Assert(j, 0)
}

func Test_ComparatorInt64(t *testing.T) {
	j := gutil.ComparatorInt64(1, 1)
	gtest.Assert(j, 0)
}

func Test_ComparatorUint(t *testing.T) {
	j := gutil.ComparatorUint(1, 1)
	gtest.Assert(j, 0)
}

func Test_ComparatorUint8(t *testing.T) {
	j := gutil.ComparatorUint8(1, 1)
	gtest.Assert(j, 0)
}

func Test_ComparatorUint16(t *testing.T) {
	j := gutil.ComparatorUint16(1, 1)
	gtest.Assert(j, 0)
}

func Test_ComparatorUint32(t *testing.T) {
	j := gutil.ComparatorUint32(1, 1)
	gtest.Assert(j, 0)
}

func Test_ComparatorUint64(t *testing.T) {
	j := gutil.ComparatorUint64(1, 1)
	gtest.Assert(j, 0)
}

func Test_ComparatorFloat32(t *testing.T) {
	j := gutil.ComparatorFloat32(1, 1)
	gtest.Assert(j, 0)
}

func Test_ComparatorFloat64(t *testing.T) {
	j := gutil.ComparatorFloat64(1, 1)
	gtest.Assert(j, 0)
}

func Test_ComparatorByte(t *testing.T) {
	j := gutil.ComparatorByte(1, 1)
	gtest.Assert(j, 0)
}

func Test_ComparatorRune(t *testing.T) {
	j := gutil.ComparatorRune(1, 1)
	gtest.Assert(j, 0)
}

func Test_ComparatorTime(t *testing.T) {
	j := gutil.ComparatorTime("2019-06-14", "2019-06-14")
	gtest.Assert(j, 0)
	k := gutil.ComparatorTime("2019-06-15", "2019-06-14")
	gtest.Assert(k, 1)
	l := gutil.ComparatorTime("2019-06-13", "2019-06-14")
	gtest.Assert(l, -1)
}
