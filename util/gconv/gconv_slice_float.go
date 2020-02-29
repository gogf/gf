// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

// SliceFloat is alias of Floats.
func SliceFloat(i interface{}) []float64 {
	return Floats(i)
}

// SliceFloat32 is alias of Float32s.
func SliceFloat32(i interface{}) []float32 {
	return Float32s(i)
}

// SliceFloat64 is alias of Float64s.
func SliceFloat64(i interface{}) []float64 {
	return Floats(i)
}

// Floats converts <i> to []float64.
func Floats(i interface{}) []float64 {
	return Float64s(i)
}

// Float32s converts <i> to []float32.
func Float32s(i interface{}) []float32 {
	if i == nil {
		return nil
	}
	if r, ok := i.([]float32); ok {
		return r
	} else {
		var array []float32
		switch value := i.(type) {
		case []string:
			array = make([]float32, len(value))
			for k, v := range value {
				array[k] = Float32(v)
			}
		case []int:
			array = make([]float32, len(value))
			for k, v := range value {
				array[k] = Float32(v)
			}
		case []int8:
			array = make([]float32, len(value))
			for k, v := range value {
				array[k] = Float32(v)
			}
		case []int16:
			array = make([]float32, len(value))
			for k, v := range value {
				array[k] = Float32(v)
			}
		case []int32:
			array = make([]float32, len(value))
			for k, v := range value {
				array[k] = Float32(v)
			}
		case []int64:
			array = make([]float32, len(value))
			for k, v := range value {
				array[k] = Float32(v)
			}
		case []uint:
			for _, v := range value {
				array = append(array, Float32(v))
			}
		case []uint8:
			array = make([]float32, len(value))
			for k, v := range value {
				array[k] = Float32(v)
			}
		case []uint16:
			array = make([]float32, len(value))
			for k, v := range value {
				array[k] = Float32(v)
			}
		case []uint32:
			array = make([]float32, len(value))
			for k, v := range value {
				array[k] = Float32(v)
			}
		case []uint64:
			array = make([]float32, len(value))
			for k, v := range value {
				array[k] = Float32(v)
			}
		case []bool:
			array = make([]float32, len(value))
			for k, v := range value {
				array[k] = Float32(v)
			}
		case []float64:
			array = make([]float32, len(value))
			for k, v := range value {
				array[k] = Float32(v)
			}
		case []interface{}:
			array = make([]float32, len(value))
			for k, v := range value {
				array[k] = Float32(v)
			}
		default:
			return []float32{Float32(i)}
		}
		return array
	}
}

// Float64s converts <i> to []float64.
func Float64s(i interface{}) []float64 {
	if i == nil {
		return nil
	}
	if r, ok := i.([]float64); ok {
		return r
	} else {
		var array []float64
		switch value := i.(type) {
		case []string:
			array = make([]float64, len(value))
			for k, v := range value {
				array[k] = Float64(v)
			}
		case []int:
			array = make([]float64, len(value))
			for k, v := range value {
				array[k] = Float64(v)
			}
		case []int8:
			array = make([]float64, len(value))
			for k, v := range value {
				array[k] = Float64(v)
			}
		case []int16:
			array = make([]float64, len(value))
			for k, v := range value {
				array[k] = Float64(v)
			}
		case []int32:
			array = make([]float64, len(value))
			for k, v := range value {
				array[k] = Float64(v)
			}
		case []int64:
			array = make([]float64, len(value))
			for k, v := range value {
				array[k] = Float64(v)
			}
		case []uint:
			for _, v := range value {
				array = append(array, Float64(v))
			}
		case []uint8:
			array = make([]float64, len(value))
			for k, v := range value {
				array[k] = Float64(v)
			}
		case []uint16:
			array = make([]float64, len(value))
			for k, v := range value {
				array[k] = Float64(v)
			}
		case []uint32:
			array = make([]float64, len(value))
			for k, v := range value {
				array[k] = Float64(v)
			}
		case []uint64:
			array = make([]float64, len(value))
			for k, v := range value {
				array[k] = Float64(v)
			}
		case []bool:
			array = make([]float64, len(value))
			for k, v := range value {
				array[k] = Float64(v)
			}
		case []float32:
			array = make([]float64, len(value))
			for k, v := range value {
				array[k] = Float64(v)
			}
		case []interface{}:
			array = make([]float64, len(value))
			for k, v := range value {
				array[k] = Float64(v)
			}
		default:
			return []float64{Float64(i)}
		}
		return array
	}
}
