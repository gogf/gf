// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

// Uint converts `any` to uint.
func Uint(any any) uint {
	v, _ := defaultConverter.Uint(any)
	return v
}

// Uint8 converts `any` to uint8.
func Uint8(any any) uint8 {
	v, _ := defaultConverter.Uint8(any)
	return v
}

// Uint16 converts `any` to uint16.
func Uint16(any any) uint16 {
	v, _ := defaultConverter.Uint16(any)
	return v
}

// Uint32 converts `any` to uint32.
func Uint32(any any) uint32 {
	v, _ := defaultConverter.Uint32(any)
	return v
}

// Uint64 converts `any` to uint64.
func Uint64(any any) uint64 {
	v, _ := defaultConverter.Uint64(any)
	return v
}
