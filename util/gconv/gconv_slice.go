// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gogf/gf/internal/utilstr"
)

// SliceInt is alias of Ints.
func SliceInt(i interface{}) []int {
	return Ints(i)
}

// SliceUint is alias of Uints.
func SliceUint(i interface{}) []uint {
	return Uints(i)
}

// SliceStr is alias of Strings.
func SliceStr(i interface{}) []string {
	return Strings(i)
}

// SliceAny is alias of Interfaces.
func SliceAny(i interface{}) []interface{} {
	return Interfaces(i)
}

// SliceFloat is alias of Floats.
func SliceFloat(i interface{}) []float64 {
	return Floats(i)
}

// SliceMap is alias of Maps.
func SliceMap(i interface{}) []map[string]interface{} {
	return Maps(i)
}

// SliceMapDeep is alias of MapsDeep.
func SliceMapDeep(i interface{}) []map[string]interface{} {
	return MapsDeep(i)
}

// SliceStruct is alias of Structs.
func SliceStruct(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	return Structs(params, pointer, mapping...)
}

// SliceStructDeep is alias of StructsDeep.
func SliceStructDeep(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	return StructsDeep(params, pointer, mapping...)
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

// Uints converts <i> to []uint.
func Uints(i interface{}) []uint {
	if i == nil {
		return nil
	}
	if r, ok := i.([]uint); ok {
		return r
	} else {
		var array []uint
		switch value := i.(type) {
		case []string:
			array = make([]uint, len(value))
			for k, v := range value {
				array[k] = Uint(v)
			}
		case []int8:
			array = make([]uint, len(value))
			for k, v := range value {
				array[k] = uint(v)
			}
		case []int16:
			array = make([]uint, len(value))
			for k, v := range value {
				array[k] = uint(v)
			}
		case []int32:
			array = make([]uint, len(value))
			for k, v := range value {
				array[k] = uint(v)
			}
		case []int64:
			array = make([]uint, len(value))
			for k, v := range value {
				array[k] = uint(v)
			}
		case []uint8:
			array = make([]uint, len(value))
			for k, v := range value {
				array[k] = uint(v)
			}
		case []uint16:
			array = make([]uint, len(value))
			for k, v := range value {
				array[k] = uint(v)
			}
		case []uint32:
			array = make([]uint, len(value))
			for k, v := range value {
				array[k] = uint(v)
			}
		case []uint64:
			array = make([]uint, len(value))
			for k, v := range value {
				array[k] = uint(v)
			}
		case []bool:
			array = make([]uint, len(value))
			for k, v := range value {
				if v {
					array[k] = 1
				} else {
					array[k] = 0
				}
			}
		case []float32:
			array = make([]uint, len(value))
			for k, v := range value {
				array[k] = Uint(v)
			}
		case []float64:
			array = make([]uint, len(value))
			for k, v := range value {
				array[k] = Uint(v)
			}
		case []interface{}:
			array = make([]uint, len(value))
			for k, v := range value {
				array[k] = Uint(v)
			}
		case [][]byte:
			array = make([]uint, len(value))
			for k, v := range value {
				array[k] = Uint(v)
			}
		default:
			return []uint{Uint(i)}
		}
		return array
	}
}

// Strings converts <i> to []string.
func Strings(i interface{}) []string {
	if i == nil {
		return nil
	}
	if r, ok := i.([]string); ok {
		return r
	} else {
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
		case [][]byte:
			array = make([]string, len(value))
			for k, v := range value {
				array[k] = String(v)
			}
		default:
			return []string{String(i)}
		}
		return array
	}
}

// Strings converts <i> to []float64.
func Floats(i interface{}) []float64 {
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

// Type assert api for Interfaces.
type apiInterfaces interface {
	Interfaces() []interface{}
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
			rv := reflect.ValueOf(i)
			kind := rv.Kind()
			for kind == reflect.Ptr {
				rv = rv.Elem()
				kind = rv.Kind()
			}
			switch kind {
			case reflect.Slice, reflect.Array:
				array = make([]interface{}, rv.Len())
				for i := 0; i < rv.Len(); i++ {
					array[i] = rv.Index(i).Interface()
				}
			case reflect.Struct:
				rt := rv.Type()
				array = make([]interface{}, rv.NumField())
				for i := 0; i < rv.NumField(); i++ {
					// Only public attributes.
					if !utilstr.IsLetterUpper(rt.Field(i).Name[0]) {
						continue
					}
					array[i] = rv.Field(i).Interface()
				}
			default:
				return []interface{}{i}
			}
		}
		return array
	}
}

// Maps converts <i> to []map[string]interface{}.
func Maps(value interface{}, tags ...string) []map[string]interface{} {
	if value == nil {
		return nil
	}
	if r, ok := value.([]map[string]interface{}); ok {
		return r
	} else {
		array := Interfaces(value)
		if len(array) == 0 {
			return nil
		}
		list := make([]map[string]interface{}, len(array))
		for k, v := range array {
			list[k] = Map(v, tags...)
		}
		return list
	}
}

// MapsDeep converts <i> to []map[string]interface{} recursively.
func MapsDeep(value interface{}, tags ...string) []map[string]interface{} {
	if value == nil {
		return nil
	}
	if r, ok := value.([]map[string]interface{}); ok {
		return r
	} else {
		array := Interfaces(value)
		if len(array) == 0 {
			return nil
		}
		list := make([]map[string]interface{}, len(array))
		for k, v := range array {
			list[k] = MapDeep(v, tags...)
		}
		return list
	}
}

// Structs converts any slice to given struct slice.
func Structs(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	return doStructs(params, pointer, false, mapping...)
}

// StructsDeep converts any slice to given struct slice recursively.
func StructsDeep(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	return doStructs(params, pointer, true, mapping...)
}

// doStructs converts any slice to given struct slice.
//
// The parameter <params> should be type of slice.
//
// The parameter <pointer> should be type of pointer to slice of struct.
// Note that if <pointer> is a pointer to another pointer of type of slice of struct,
// it will create the struct/pointer internally.
func doStructs(params interface{}, pointer interface{}, deep bool, mapping ...map[string]string) (err error) {
	if params == nil {
		return errors.New("params cannot be nil")
	}
	if pointer == nil {
		return errors.New("object pointer cannot be nil")
	}
	pointerRv, ok := pointer.(reflect.Value)
	if !ok {
		pointerRv = reflect.ValueOf(pointer)
		if kind := pointerRv.Kind(); kind != reflect.Ptr {
			return fmt.Errorf("pointer should be type of pointer, but got: %v", kind)
		}
	}
	rv := reflect.ValueOf(params)
	kind := rv.Kind()
	for kind == reflect.Ptr {
		rv = rv.Elem()
		kind = rv.Kind()
	}
	switch kind {
	case reflect.Slice, reflect.Array:
		// If <params> is an empty slice, no conversion.
		if rv.Len() == 0 {
			return nil
		}
		array := reflect.MakeSlice(pointerRv.Type().Elem(), rv.Len(), rv.Len())
		itemType := array.Index(0).Type()
		for i := 0; i < rv.Len(); i++ {
			if itemType.Kind() == reflect.Ptr {
				// Slice element is type pointer.
				e := reflect.New(itemType.Elem()).Elem()
				if deep {
					if err = StructDeep(rv.Index(i).Interface(), e, mapping...); err != nil {
						return err
					}
				} else {
					if err = Struct(rv.Index(i).Interface(), e, mapping...); err != nil {
						return err
					}
				}
				array.Index(i).Set(e.Addr())
			} else {
				// Slice element is not type of pointer.
				e := reflect.New(itemType).Elem()
				if deep {
					if err = StructDeep(rv.Index(i).Interface(), e, mapping...); err != nil {
						return err
					}
				} else {
					if err = Struct(rv.Index(i).Interface(), e, mapping...); err != nil {
						return err
					}
				}
				array.Index(i).Set(e)
			}
		}
		pointerRv.Elem().Set(array)
		return nil
	default:
		return fmt.Errorf("params should be type of slice, but got: %v", kind)
	}
}
