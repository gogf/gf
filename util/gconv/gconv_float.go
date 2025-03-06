// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

// Float32 converts `any` to float32.
func Float32(any any) float32 {
	v, _ := defaultConverter.Float32(any)
	return v
}

// Float64 converts `any` to float64.
func Float64(any any) float64 {
	v, _ := defaultConverter.Float64(any)
	return v
}
