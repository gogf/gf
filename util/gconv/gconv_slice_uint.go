// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

// SliceUint is alias of Uints.
func SliceUint(any interface{}) []uint {
	return Uints(any)
}

// SliceUint32 is alias of Uint32s.
func SliceUint32(any interface{}) []uint32 {
	return Uint32s(any)
}

// SliceUint64 is alias of Uint64s.
func SliceUint64(any interface{}) []uint64 {
	return Uint64s(any)
}

// Uints converts `any` to []uint.
func Uints(any interface{}) []uint {
	result, _ := defaultConverter.SliceUint(any, SliceOption{
		ContinueOnError: true,
	})
	return result
}

// Uint32s converts `any` to []uint32.
func Uint32s(any interface{}) []uint32 {
	result, _ := defaultConverter.SliceUint32(any, SliceOption{
		ContinueOnError: true,
	})
	return result
}

// Uint64s converts `any` to []uint64.
func Uint64s(any interface{}) []uint64 {
	result, _ := defaultConverter.SliceUint64(any, SliceOption{
		ContinueOnError: true,
	})
	return result
}
