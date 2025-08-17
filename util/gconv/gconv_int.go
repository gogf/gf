// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

// Int converts `any` to int.
func Int(any any) int {
	v, _ := defaultConverter.Int(any)
	return v
}

// Int8 converts `any` to int8.
func Int8(any any) int8 {
	v, _ := defaultConverter.Int8(any)
	return v
}

// Int16 converts `any` to int16.
func Int16(any any) int16 {
	v, _ := defaultConverter.Int16(any)
	return v
}

// Int32 converts `any` to int32.
func Int32(any any) int32 {
	v, _ := defaultConverter.Int32(any)
	return v
}

// Int64 converts `any` to int64.
func Int64(any any) int64 {
	v, _ := defaultConverter.Int64(any)
	return v
}
