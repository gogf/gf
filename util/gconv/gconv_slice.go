// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogf/gf/errors/gerror"
	"reflect"
)

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

// Maps converts <i> to []map[string]interface{}.
func Maps(value interface{}, tags ...string) []map[string]interface{} {
	if value == nil {
		return nil
	}
	switch r := value.(type) {
	case string:
		list := make([]map[string]interface{}, 0)
		if len(r) > 0 && r[0] == '[' && r[len(r)-1] == ']' {
			if err := json.Unmarshal([]byte(r), &list); err != nil {
				return nil
			}
			return list
		} else {
			return nil
		}

	case []byte:
		list := make([]map[string]interface{}, 0)
		if len(r) > 0 && r[0] == '[' && r[len(r)-1] == ']' {
			if err := json.Unmarshal(r, &list); err != nil {
				return nil
			}
			return list
		} else {
			return nil
		}

	case []map[string]interface{}:
		return r

	default:
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
//
// TODO completely implement the recursive converting for all types.
func MapsDeep(value interface{}, tags ...string) []map[string]interface{} {
	if value == nil {
		return nil
	}
	switch r := value.(type) {
	case string:
		list := make([]map[string]interface{}, 0)
		if len(r) > 0 && r[0] == '[' && r[len(r)-1] == ']' {
			if err := json.Unmarshal([]byte(r), &list); err != nil {
				return nil
			}
			return list
		} else {
			return nil
		}

	case []byte:
		list := make([]map[string]interface{}, 0)
		if len(r) > 0 && r[0] == '[' && r[len(r)-1] == ']' {
			if err := json.Unmarshal(r, &list); err != nil {
				return nil
			}
			return list
		} else {
			return nil
		}

	case []map[string]interface{}:
		return r

	default:
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
	defer func() {
		// Catch the panic, especially the reflect operation panics.
		if e := recover(); e != nil {
			err = gerror.NewfSkip(1, "%v", e)
		}
	}()
	pointerRv, ok := pointer.(reflect.Value)
	if !ok {
		pointerRv = reflect.ValueOf(pointer)
		if kind := pointerRv.Kind(); kind != reflect.Ptr {
			return fmt.Errorf("pointer should be type of pointer, but got: %v", kind)
		}
	}
	params = Maps(params)
	var (
		reflectValue = reflect.ValueOf(params)
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Slice, reflect.Array:
		// If <params> is an empty slice, no conversion.
		if reflectValue.Len() == 0 {
			return nil
		}
		var (
			array    = reflect.MakeSlice(pointerRv.Type().Elem(), reflectValue.Len(), reflectValue.Len())
			itemType = array.Index(0).Type()
		)
		for i := 0; i < reflectValue.Len(); i++ {
			if itemType.Kind() == reflect.Ptr {
				// Slice element is type pointer.
				e := reflect.New(itemType.Elem()).Elem()
				if deep {
					if err = StructDeep(reflectValue.Index(i).Interface(), e, mapping...); err != nil {
						return err
					}
				} else {
					if err = Struct(reflectValue.Index(i).Interface(), e, mapping...); err != nil {
						return err
					}
				}
				array.Index(i).Set(e.Addr())
			} else {
				// Slice element is not type of pointer.
				e := reflect.New(itemType).Elem()
				if deep {
					if err = StructDeep(reflectValue.Index(i).Interface(), e, mapping...); err != nil {
						return err
					}
				} else {
					if err = Struct(reflectValue.Index(i).Interface(), e, mapping...); err != nil {
						return err
					}
				}
				array.Index(i).Set(e)
			}
		}
		pointerRv.Elem().Set(array)
		return nil
	default:
		return fmt.Errorf("params should be type of slice, but got: %v", reflectKind)
	}
}
