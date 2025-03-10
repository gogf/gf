// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

// SliceFloat is alias of Floats.
func SliceFloat(any interface{}) []float64 {
	return Floats(any)
}

// SliceFloat32 is alias of Float32s.
func SliceFloat32(any interface{}) []float32 {
	return Float32s(any)
}

// SliceFloat64 is alias of Float64s.
func SliceFloat64(any interface{}) []float64 {
	return Floats(any)
}

// Floats converts `any` to []float64.
func Floats(any interface{}) []float64 {
	return Float64s(any)
}

// Float32s converts `any` to []float32.
func Float32s(any interface{}) []float32 {
	result, _ := defaultConverter.SliceFloat32(any, SliceOption{
		ContinueOnError: true,
	})
	return result
}

// Float64s converts `any` to []float64.
func Float64s(any interface{}) []float64 {
	result, _ := defaultConverter.SliceFloat64(any, SliceOption{
		ContinueOnError: true,
	})
	return result
}
