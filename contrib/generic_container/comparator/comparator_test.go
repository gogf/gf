package comparator

import (
	"sort"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func TestComparators(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a := []string{"2", "3", "1", "255", "127"}
		sort.Slice(a, func(i, j int) bool {
			return ComparatorString(a[i], a[j]) < 0
		})
		t.Assert(a, []string{"1", "127", "2", "255", "3"})
	})
	gtest.C(t, func(t *gtest.T) {
		a := []int{2, 3, 1, -2147483648, 255, 2147483647, 127}
		sort.Slice(a, func(i, j int) bool {
			return ComparatorInt(a[i], a[j]) < 0
		})
		t.Assert(a, []int{-2147483648, 1, 2, 3, 127, 255, 2147483647})
	})
	gtest.C(t, func(t *gtest.T) {
		a := []int8{2, 3, 1, -128, 127}
		sort.Slice(a, func(i, j int) bool {
			return ComparatorInt8(a[i], a[j]) < 0
		})
		t.Assert(a, []int8{-128, 1, 2, 3, 127})
	})
	gtest.C(t, func(t *gtest.T) {
		a := []int16{2, 3, 1, 32767, -32768}
		sort.Slice(a, func(i, j int) bool {
			return ComparatorInt16(a[i], a[j]) < 0
		})
		t.Assert(a, []int16{-32768, 1, 2, 3, 32767})
	})
	gtest.C(t, func(t *gtest.T) {
		a := []int32{2, 3, 1, -2147483648, 255, 2147483647, 127}
		sort.Slice(a, func(i, j int) bool {
			return ComparatorInt32(a[i], a[j]) < 0
		})
		t.Assert(a, []int32{-2147483648, 1, 2, 3, 127, 255, 2147483647})
	})
	gtest.C(t, func(t *gtest.T) {
		a := []int64{2, 3, 1, -9223372036854775808, 9223372036854775807}
		sort.Slice(a, func(i, j int) bool {
			return ComparatorInt64(a[i], a[j]) < 0
		})
		t.Assert(a, []int64{-9223372036854775808, 1, 2, 3, 9223372036854775807})
	})
	gtest.C(t, func(t *gtest.T) {
		a := []uint{2, 3, 1, 0, 255, 4294967296, 127}
		sort.Slice(a, func(i, j int) bool {
			return ComparatorUint(a[i], a[j]) < 0
		})
		t.Assert(a, []uint{0, 1, 2, 3, 127, 255, 4294967296})
	})
	gtest.C(t, func(t *gtest.T) {
		a := []uint8{2, 3, 1, 255, 127}
		sort.Slice(a, func(i, j int) bool {
			return ComparatorUint8(a[i], a[j]) < 0
		})
		t.Assert(a, []uint8{1, 2, 3, 127, 255})
	})
	gtest.C(t, func(t *gtest.T) {
		a := []uint16{2, 3, 1, 0, 255, 65535, 127}
		sort.Slice(a, func(i, j int) bool {
			return ComparatorUint16(a[i], a[j]) < 0
		})
		t.Assert(a, []uint16{0, 1, 2, 3, 127, 255, 65535})
	})
	gtest.C(t, func(t *gtest.T) {
		a := []uint32{2, 3, 1, 0, 255, 4294967295, 127}
		sort.Slice(a, func(i, j int) bool {
			return ComparatorUint32(a[i], a[j]) < 0
		})
		t.Assert(a, []uint32{0, 1, 2, 3, 127, 255, 4294967295})
	})
	gtest.C(t, func(t *gtest.T) {
		a := []uint64{2, 3, 1, 0, 255, 18446744073709551615, 127}
		sort.Slice(a, func(i, j int) bool {
			return ComparatorUint64(a[i], a[j]) < 0
		})
		t.Assert(a, []uint64{0, 1, 2, 3, 127, 255, 18446744073709551615})
	})
	gtest.C(t, func(t *gtest.T) {
		a := []float32{2, 3, 1, -2.45534534, 1.000000001}
		sort.Slice(a, func(i, j int) bool {
			return ComparatorFloat32(a[i], a[j]) < 0
		})
		t.Assert(a, []float32{-2.45534534, 1, 1.000000001, 2, 3})
	})
	gtest.C(t, func(t *gtest.T) {
		a := []float64{2, 3, 1, -2.45534534, 1.000000001}
		sort.Slice(a, func(i, j int) bool {
			return ComparatorFloat64(a[i], a[j]) < 0
		})
		t.Assert(a, []float64{-2.45534534, 1, 1.000000001, 2, 3})
	})
}
