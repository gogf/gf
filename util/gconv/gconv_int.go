// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

// Int converts `any` to int.
func Int(anyInput any) int {
	v, _ := defaultConverter.Int(anyInput)
	return v
}

// Int8 converts `any` to int8.
func Int8(anyInput any) int8 {
	v, _ := defaultConverter.Int8(anyInput)
	return v
}

// Int16 converts `any` to int16.
func Int16(anyInput any) int16 {
	v, _ := defaultConverter.Int16(anyInput)
	return v
}

// Int32 converts `any` to int32.
func Int32(anyInput any) int32 {
	v, _ := defaultConverter.Int32(anyInput)
	return v
}

// Int64 converts `any` to int64.
func Int64(anyInput any) int64 {
	v, _ := defaultConverter.Int64(anyInput)
	return v
}
