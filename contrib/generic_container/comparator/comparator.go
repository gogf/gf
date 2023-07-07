// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package comparator

import (
	"strings"
	"time"

	"github.com/gogf/gf/v2/util/gconv"
)

// Comparator is a function that compare a and b, and returns the result as int.
//
// Should return a number:
//
//	negative , if a < b
//	zero     , if a == b
//	positive , if a > b
type Comparator[T comparable] func(a, b T) int

type IComparable[T any] interface {
	Compare(other T) int
}

// ComparatorAny provides a comparison on any types
func ComparatorAny[T comparable](a, b T) int {
	if any(a) == nil && any(b) == nil {
		return 0
	}
	if any(a) == nil && any(b) != nil {
		return -1
	}
	if any(a) != nil && any(b) == nil {
		return 1
	}
	switch va := any(a).(type) {
	case string:
		return ComparatorString(va, any(b).(string))
	case int:
		return ComparatorInt(va, any(b).(int))
	case int8:
		return ComparatorInt8(va, any(b).(int8))
	case int16:
		return ComparatorInt16(va, any(b).(int16))
	case int32:
		return ComparatorInt32(va, any(b).(int32))
	case int64:
		return ComparatorInt64(va, any(b).(int64))
	case uint:
		return ComparatorUint(va, any(b).(uint))
	case uint8:
		return ComparatorUint8(va, any(b).(uint8))
	case uint16:
		return ComparatorUint16(va, any(b).(uint16))
	case uint32:
		return ComparatorUint32(va, any(b).(uint32))
	case uint64:
		return ComparatorUint64(va, any(b).(uint64))
	case float32:
		return ComparatorFloat32(va, any(b).(float32))
	case float64:
		return ComparatorFloat64(va, any(b).(float64))
	case time.Time:
		return ComparatorTime(va, any(b).(time.Time))
	default:
		if aComp, ok := any(a).(IComparable[T]); ok {
			return aComp.Compare(b)
		}
		return strings.Compare(gconv.String(a), gconv.String(b))
	}
}

// ComparatorString provides a fast comparison on strings.
func ComparatorString(a, b string) int {
	return strings.Compare(a, b)
}

// ComparatorInt provides a basic comparison on int.
func ComparatorInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// ComparatorInt8 provides a basic comparison on int8.
func ComparatorInt8(a, b int8) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// ComparatorInt16 provides a basic comparison on int16.
func ComparatorInt16(a, b int16) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// ComparatorInt32 provides a basic comparison on int32.
func ComparatorInt32(a, b int32) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// ComparatorInt64 provides a basic comparison on int64.
func ComparatorInt64(a, b int64) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// ComparatorUint provides a basic comparison on uint.
func ComparatorUint(a, b uint) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// ComparatorUint8 provides a basic comparison on uint8.
func ComparatorUint8(a, b uint8) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// ComparatorUint16 provides a basic comparison on uint16.
func ComparatorUint16(a, b uint16) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// ComparatorUint32 provides a basic comparison on uint32.
func ComparatorUint32(a, b uint32) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// ComparatorUint64 provides a basic comparison on uint64.
func ComparatorUint64(a, b uint64) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// ComparatorFloat32 provides a basic comparison on float32.
func ComparatorFloat32(a, b float32) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// ComparatorFloat64 provides a basic comparison on float64.
func ComparatorFloat64(a, b float64) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// ComparatorByte provides a basic comparison on byte.
func ComparatorByte(a, b byte) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// ComparatorRune provides a basic comparison on rune.
func ComparatorRune(a, b rune) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// ComparatorTime provides a basic comparison on time.Time.
func ComparatorTime(a, b time.Time) int {
	switch {
	case a.After(b):
		return 1
	case a.Before(b):
		return -1
	default:
		return 0
	}
}
