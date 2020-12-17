// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"reflect"
)

// SliceAny is alias of Interfaces.
func SliceAny(i interface{}) []interface{} {
	return Interfaces(i)
}

// Interfaces converts <i> to []interface{}.
func Interfaces(i interface{}) []interface{} {
	if i == nil {
		return nil
	}
	if r, ok := i.([]interface{}); ok {
		return r
	} else if r, ok := i.(apiInterfaces); ok {
		return r.Interfaces()
	} else {
		var array []interface{}
		switch value := i.(type) {
		case []string:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []int:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []int8:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []int16:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []int32:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []int64:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []uint:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []uint8:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []uint16:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []uint32:
			for _, v := range value {
				array = append(array, v)
			}
		case []uint64:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []bool:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []float32:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []float64:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		default:
			// Finally we use reflection.
			var (
				reflectValue = reflect.ValueOf(i)
				reflectKind  = reflectValue.Kind()
			)
			for reflectKind == reflect.Ptr {
				reflectValue = reflectValue.Elem()
				reflectKind = reflectValue.Kind()
			}
			switch reflectKind {
			case reflect.Slice, reflect.Array:
				array = make([]interface{}, reflectValue.Len())
				for i := 0; i < reflectValue.Len(); i++ {
					array[i] = reflectValue.Index(i).Interface()
				}
			// Deprecated.
			//// Eg: {"K1": "v1", "K2": "v2"} => ["K1", "v1", "K2", "v2"]
			//case reflect.Map:
			//	array = make([]interface{}, 0)
			//	for _, key := range reflectValue.MapKeys() {
			//		array = append(array, key.Interface())
			//		array = append(array, reflectValue.MapIndex(key).Interface())
			//	}
			//// Eg: {"K1": "v1", "K2": "v2"} => ["K1", "v1", "K2", "v2"]
			//case reflect.Struct:
			//	array = make([]interface{}, 0)
			//	// Note that, it uses the gconv tag name instead of the attribute name if
			//	// the gconv tag is fined in the struct attributes.
			//	for k, v := range Map(reflectValue) {
			//		array = append(array, k)
			//		array = append(array, v)
			//	}
			default:
				return []interface{}{i}
			}
		}
		return array
	}
}
