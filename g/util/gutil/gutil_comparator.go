package gutil

import (
	"github.com/gogf/gf/g/util/gconv"
	"strings"
)

// Comparator is a function that compare a and b, and returns the result as int.
//
// Should return a number:
//    negative , if a < b
//    zero     , if a == b
//    positive , if a > b
type Comparator func(a, b interface{}) int

// StringComparator provides a fast comparison on strings.
func StringComparator(a, b interface{}) int {
	return strings.Compare(gconv.String(a), gconv.String(b))
}

// IntComparator provides a basic comparison on int.
func IntComparator(a, b interface{}) int {
	return gconv.Int(a) - gconv.Int(b)
}

// Int8Comparator provides a basic comparison on int8.
func Int8Comparator(a, b interface{}) int {
	return int(gconv.Int8(a) - gconv.Int8(b))
}

// Int16Comparator provides a basic comparison on int16.
func Int16Comparator(a, b interface{}) int {
	return int(gconv.Int16(a) - gconv.Int16(b))
}

// Int32Comparator provides a basic comparison on int32.
func Int32Comparator(a, b interface{}) int {
	return int(gconv.Int32(a) - gconv.Int32(b))
}

// Int64Comparator provides a basic comparison on int64.
func Int64Comparator(a, b interface{}) int {
	return int(gconv.Int64(a) - gconv.Int64(b))
}

// UintComparator provides a basic comparison on uint.
func UintComparator(a, b interface{}) int {
	return int(gconv.Uint(a) - gconv.Uint(b))
}

// Uint8Comparator provides a basic comparison on uint8.
func Uint8Comparator(a, b interface{}) int {
	return int(gconv.Uint8(a) - gconv.Uint8(b))
}

// Uint16Comparator provides a basic comparison on uint16.
func Uint16Comparator(a, b interface{}) int {
	return int(gconv.Uint16(a) - gconv.Uint16(b))
}

// Uint32Comparator provides a basic comparison on uint32.
func Uint32Comparator(a, b interface{}) int {
	return int(gconv.Uint32(a) - gconv.Uint32(b))
}

// Uint64Comparator provides a basic comparison on uint64.
func Uint64Comparator(a, b interface{}) int {
	return int(gconv.Uint64(a) - gconv.Uint64(b))
}

// Float32Comparator provides a basic comparison on float32.
func Float32Comparator(a, b interface{}) int {
	return int(gconv.Float32(a) - gconv.Float32(b))
}

// Float64Comparator provides a basic comparison on float64.
func Float64Comparator(a, b interface{}) int {
	return int(gconv.Float64(a) - gconv.Float64(b))
}

// ByteComparator provides a basic comparison on byte.
func ByteComparator(a, b interface{}) int {
	return int(gconv.Byte(a) - gconv.Byte(b))
}

// RuneComparator provides a basic comparison on rune.
func RuneComparator(a, b interface{}) int {
	return int(gconv.Rune(a) - gconv.Rune(b))
}

// TimeComparator provides a basic comparison on time.Time.
func TimeComparator(a, b interface{}) int {
	aTime := gconv.Time(a)
	bTime := gconv.Time(b)
	switch {
		case aTime.After(bTime):
			return 1
		case aTime.Before(bTime):
			return -1
		default:
			return 0
	}
}