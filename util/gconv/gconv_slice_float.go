// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

// SliceFloat is alias of Floats.
func SliceFloat(anyInput any) []float64 {
	return Floats(anyInput)
}

// SliceFloat32 is alias of Float32s.
func SliceFloat32(anyInput any) []float32 {
	return Float32s(anyInput)
}

// SliceFloat64 is alias of Float64s.
func SliceFloat64(anyInput any) []float64 {
	return Floats(anyInput)
}

// Floats converts `any` to []float64.
func Floats(anyInput any) []float64 {
	return Float64s(anyInput)
}

// Float32s converts `any` to []float32.
func Float32s(anyInput any) []float32 {
	result, _ := defaultConverter.SliceFloat32(anyInput, SliceOption{
		ContinueOnError: true,
	})
	return result
}

// Float64s converts `any` to []float64.
func Float64s(anyInput any) []float64 {
	result, _ := defaultConverter.SliceFloat64(anyInput, SliceOption{
		ContinueOnError: true,
	})
	return result
}
