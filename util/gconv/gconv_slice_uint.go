// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

// SliceUint is alias of Uints.
func SliceUint(anyInput any) []uint {
	return Uints(anyInput)
}

// SliceUint32 is alias of Uint32s.
func SliceUint32(anyInput any) []uint32 {
	return Uint32s(anyInput)
}

// SliceUint64 is alias of Uint64s.
func SliceUint64(anyInput any) []uint64 {
	return Uint64s(anyInput)
}

// Uints converts `any` to []uint.
func Uints(anyInput any) []uint {
	result, _ := defaultConverter.SliceUint(anyInput, SliceOption{
		ContinueOnError: true,
	})
	return result
}

// Uint32s converts `any` to []uint32.
func Uint32s(anyInput any) []uint32 {
	result, _ := defaultConverter.SliceUint32(anyInput, SliceOption{
		ContinueOnError: true,
	})
	return result
}

// Uint64s converts `any` to []uint64.
func Uint64s(anyInput any) []uint64 {
	result, _ := defaultConverter.SliceUint64(anyInput, SliceOption{
		ContinueOnError: true,
	})
	return result
}
