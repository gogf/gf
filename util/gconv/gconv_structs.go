// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/json"
	"reflect"
)

// Structs converts any slice to given struct slice.
func Structs(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	return doStructs(params, pointer, mapping...)
}

// StructsDeep converts any slice to given struct slice recursively.
// Deprecated, use Structs instead.
func StructsDeep(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	return doStructs(params, pointer, mapping...)
}

// doStructs converts any slice to given struct slice.
//
// It automatically checks and converts json string to []map if <params> is string/[]byte.
//
// The parameter <pointer> should be type of pointer to slice of struct.
// Note that if <pointer> is a pointer to another pointer of type of slice of struct,
// it will create the struct/pointer internally.
func doStructs(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	if params == nil {
		// If <params> is nil, no conversion.
		return nil
	}
	if pointer == nil {
		return gerror.New("object pointer cannot be nil")
	}

	if doStructsByDirectReflectSet(params, pointer) {
		return nil
	}

	defer func() {
		// Catch the panic, especially the reflect operation panics.
		if e := recover(); e != nil {
			err = gerror.NewSkipf(1, "%v", e)
		}
	}()
	// If given <params> is JSON, it then uses json.Unmarshal doing the converting.
	switch r := params.(type) {
	case []byte:
		if json.Valid(r) {
			if rv, ok := pointer.(reflect.Value); ok {
				if rv.Kind() == reflect.Ptr {
					return json.Unmarshal(r, rv.Interface())
				}
			} else {
				return json.Unmarshal(r, pointer)
			}
		}
	case string:
		if paramsBytes := []byte(r); json.Valid(paramsBytes) {
			if rv, ok := pointer.(reflect.Value); ok {
				if rv.Kind() == reflect.Ptr {
					return json.Unmarshal(paramsBytes, rv.Interface())
				}
			} else {
				return json.Unmarshal(paramsBytes, pointer)
			}
		}
	}
	// Pointer type check.
	pointerRv, ok := pointer.(reflect.Value)
	if !ok {
		pointerRv = reflect.ValueOf(pointer)
		if kind := pointerRv.Kind(); kind != reflect.Ptr {
			return gerror.Newf("pointer should be type of pointer, but got: %v", kind)
		}
	}
	// Converting <params> to map slice.
	paramsMaps := Maps(params)
	// If <params> is an empty slice, no conversion.
	if len(paramsMaps) == 0 {
		return nil
	}
	var (
		array    = reflect.MakeSlice(pointerRv.Type().Elem(), len(paramsMaps), len(paramsMaps))
		itemType = array.Index(0).Type()
	)
	for i := 0; i < len(paramsMaps); i++ {
		if itemType.Kind() == reflect.Ptr {
			// Slice element is type pointer.
			e := reflect.New(itemType.Elem()).Elem()
			if err = Struct(paramsMaps[i], e, mapping...); err != nil {
				return err
			}
			array.Index(i).Set(e.Addr())
		} else {
			// Slice element is not type of pointer.
			e := reflect.New(itemType).Elem()
			if err = Struct(paramsMaps[i], e, mapping...); err != nil {
				return err
			}
			array.Index(i).Set(e)
		}
	}
	pointerRv.Elem().Set(array)
	return nil
}

// doStructsByDirectReflectSet do the converting directly using reflect Set.
// It returns true if success, or else false.
func doStructsByDirectReflectSet(params interface{}, pointer interface{}) (ok bool) {
	v1 := reflect.ValueOf(pointer)
	v2 := reflect.ValueOf(params)
	if v1.Kind() == reflect.Ptr {
		if elem := v1.Elem(); elem.IsValid() && elem.Type() == v2.Type() {
			elem.Set(v2)
			ok = true
		}
	}
	return ok
}
