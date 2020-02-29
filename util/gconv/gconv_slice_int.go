// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

// SliceInt is alias of Ints.
func SliceInt(i interface{}) []int {
	return Ints(i)
}

// SliceInt32 is alias of Int32s.
func SliceInt32(i interface{}) []int32 {
	return Int32s(i)
}

// SliceInt is alias of Int64s.
func SliceInt64(i interface{}) []int64 {
	return Int64s(i)
}

// Ints converts <i> to []int.
func Ints(i interface{}) []int {
	if i == nil {
		return nil
	}
	if r, ok := i.([]int); ok {
		return r
	} else {
		var array []int
		switch value := i.(type) {
		case []string:
			array = make([]int, len(value))
			for k, v := range value {
				array[k] = Int(v)
			}
		case []int8:
			array = make([]int, len(value))
			for k, v := range value {
				array[k] = int(v)
			}
		case []int16:
			array = make([]int, len(value))
			for k, v := range value {
				array[k] = int(v)
			}
		case []int32:
			array = make([]int, len(value))
			for k, v := range value {
				array[k] = int(v)
			}
		case []int64:
			array = make([]int, len(value))
			for k, v := range value {
				array[k] = int(v)
			}
		case []uint:
			array = make([]int, len(value))
			for k, v := range value {
				array[k] = int(v)
			}
		case []uint8:
			array = make([]int, len(value))
			for k, v := range value {
				array[k] = int(v)
			}
		case []uint16:
			array = make([]int, len(value))
			for k, v := range value {
				array[k] = int(v)
			}
		case []uint32:
			array = make([]int, len(value))
			for k, v := range value {
				array[k] = int(v)
			}
		case []uint64:
			array = make([]int, len(value))
			for k, v := range value {
				array[k] = int(v)
			}
		case []bool:
			array = make([]int, len(value))
			for k, v := range value {
				if v {
					array[k] = 1
				} else {
					array[k] = 0
				}
			}
		case []float32:
			array = make([]int, len(value))
			for k, v := range value {
				array[k] = Int(v)
			}
		case []float64:
			array = make([]int, len(value))
			for k, v := range value {
				array[k] = Int(v)
			}
		case []interface{}:
			array = make([]int, len(value))
			for k, v := range value {
				array[k] = Int(v)
			}
		case [][]byte:
			array = make([]int, len(value))
			for k, v := range value {
				array[k] = Int(v)
			}
		default:
			return []int{Int(i)}
		}
		return array
	}
}

// Int32s converts <i> to []int32.
func Int32s(i interface{}) []int32 {
	if i == nil {
		return nil
	}
	if r, ok := i.([]int32); ok {
		return r
	} else {
		var array []int32
		switch value := i.(type) {
		case []string:
			array = make([]int32, len(value))
			for k, v := range value {
				array[k] = Int32(v)
			}
		case []int:
			array = make([]int32, len(value))
			for k, v := range value {
				array[k] = int32(v)
			}
		case []int8:
			array = make([]int32, len(value))
			for k, v := range value {
				array[k] = int32(v)
			}
		case []int16:
			array = make([]int32, len(value))
			for k, v := range value {
				array[k] = int32(v)
			}
		case []int64:
			array = make([]int32, len(value))
			for k, v := range value {
				array[k] = int32(v)
			}
		case []uint:
			array = make([]int32, len(value))
			for k, v := range value {
				array[k] = int32(v)
			}
		case []uint8:
			array = make([]int32, len(value))
			for k, v := range value {
				array[k] = int32(v)
			}
		case []uint16:
			array = make([]int32, len(value))
			for k, v := range value {
				array[k] = int32(v)
			}
		case []uint32:
			array = make([]int32, len(value))
			for k, v := range value {
				array[k] = int32(v)
			}
		case []uint64:
			array = make([]int32, len(value))
			for k, v := range value {
				array[k] = int32(v)
			}
		case []bool:
			array = make([]int32, len(value))
			for k, v := range value {
				if v {
					array[k] = 1
				} else {
					array[k] = 0
				}
			}
		case []float32:
			array = make([]int32, len(value))
			for k, v := range value {
				array[k] = Int32(v)
			}
		case []float64:
			array = make([]int32, len(value))
			for k, v := range value {
				array[k] = Int32(v)
			}
		case []interface{}:
			array = make([]int32, len(value))
			for k, v := range value {
				array[k] = Int32(v)
			}
		case [][]byte:
			array = make([]int32, len(value))
			for k, v := range value {
				array[k] = Int32(v)
			}
		default:
			return []int32{Int32(i)}
		}
		return array
	}
}

// Int64s converts <i> to []int64.
func Int64s(i interface{}) []int64 {
	if i == nil {
		return nil
	}
	if r, ok := i.([]int64); ok {
		return r
	} else {
		var array []int64
		switch value := i.(type) {
		case []string:
			array = make([]int64, len(value))
			for k, v := range value {
				array[k] = Int64(v)
			}
		case []int:
			array = make([]int64, len(value))
			for k, v := range value {
				array[k] = int64(v)
			}
		case []int8:
			array = make([]int64, len(value))
			for k, v := range value {
				array[k] = int64(v)
			}
		case []int16:
			array = make([]int64, len(value))
			for k, v := range value {
				array[k] = int64(v)
			}
		case []int32:
			array = make([]int64, len(value))
			for k, v := range value {
				array[k] = int64(v)
			}
		case []uint:
			array = make([]int64, len(value))
			for k, v := range value {
				array[k] = int64(v)
			}
		case []uint8:
			array = make([]int64, len(value))
			for k, v := range value {
				array[k] = int64(v)
			}
		case []uint16:
			array = make([]int64, len(value))
			for k, v := range value {
				array[k] = int64(v)
			}
		case []uint32:
			array = make([]int64, len(value))
			for k, v := range value {
				array[k] = int64(v)
			}
		case []uint64:
			array = make([]int64, len(value))
			for k, v := range value {
				array[k] = int64(v)
			}
		case []bool:
			array = make([]int64, len(value))
			for k, v := range value {
				if v {
					array[k] = 1
				} else {
					array[k] = 0
				}
			}
		case []float32:
			array = make([]int64, len(value))
			for k, v := range value {
				array[k] = Int64(v)
			}
		case []float64:
			array = make([]int64, len(value))
			for k, v := range value {
				array[k] = Int64(v)
			}
		case []interface{}:
			array = make([]int64, len(value))
			for k, v := range value {
				array[k] = Int64(v)
			}
		case [][]byte:
			array = make([]int64, len(value))
			for k, v := range value {
				array[k] = Int64(v)
			}
		default:
			return []int64{Int64(i)}
		}
		return array
	}
}
