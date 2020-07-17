// Copyright 2019 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package gconv

// SliceStr is alias of Strings.
func SliceStr(i interface{}) []string {
	return Strings(i)
}

// Strings converts <i> to []string.
func Strings(i interface{}) []string {
	if i == nil {
		return nil
	}
	var array []string
	switch value := i.(type) {
	case []int:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []int8:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []int16:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []int32:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []int64:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []uint:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []uint8:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []uint16:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []uint32:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []uint64:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []bool:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []float32:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []float64:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []interface{}:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []string:
		array = value
	case [][]byte:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	default:
		if v, ok := i.(apiStrings); ok {
			return v.Strings()
		}
		if v, ok := i.(apiInterfaces); ok {
			return Strings(v.Interfaces())
		}
		return []string{String(i)}
	}
	return array
}
