// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

// Uint converts `any` to uint.
func Uint(anyInput any) uint {
	v, _ := defaultConverter.Uint(anyInput)
	return v
}

// Uint8 converts `any` to uint8.
func Uint8(anyInput any) uint8 {
	v, _ := defaultConverter.Uint8(anyInput)
	return v
}

// Uint16 converts `any` to uint16.
func Uint16(anyInput any) uint16 {
	v, _ := defaultConverter.Uint16(anyInput)
	return v
}

// Uint32 converts `any` to uint32.
func Uint32(anyInput any) uint32 {
	v, _ := defaultConverter.Uint32(anyInput)
	return v
}

// Uint64 converts `any` to uint64.
func Uint64(anyInput any) uint64 {
	v, _ := defaultConverter.Uint64(anyInput)
	return v
}
