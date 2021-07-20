// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
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
	var keyToAttributeNameMapping map[string]string
	if len(mapping) > 0 {
		keyToAttributeNameMapping = mapping[0]
	}
	return doStructs(params, pointer, keyToAttributeNameMapping, "")
}

// StructsTag acts as Structs but also with support for priority tag feature, which retrieves the
// specified tags for `params` key-value items to struct attribute names mapping.
// The parameter `priorityTag` supports multiple tags that can be joined with char ','.
func StructsTag(params interface{}, pointer interface{}, priorityTag string) (err error) {
	return doStructs(params, pointer, nil, priorityTag)
}

// StructsDeep converts any slice to given struct slice recursively.
// Deprecated, use Structs instead.
func StructsDeep(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	var keyToAttributeNameMapping map[string]string
	if len(mapping) > 0 {
		keyToAttributeNameMapping = mapping[0]
	}
	return doStructs(params, pointer, keyToAttributeNameMapping, "")
}

// doStructs converts any slice to given struct slice.
//
// It automatically checks and converts json string to []map if `params` is string/[]byte.
//
// The parameter `pointer` should be type of pointer to slice of struct.
// Note that if `pointer` is a pointer to another pointer of type of slice of struct,
// it will create the struct/pointer internally.
func doStructs(params interface{}, pointer interface{}, mapping map[string]string, priorityTag string) (err error) {
	if params == nil {
		// If `params` is nil, no conversion.
		return nil
	}
	if pointer == nil {
		return gerror.NewCode(gerror.CodeInvalidParameter, "object pointer cannot be nil")
	}

	if doStructsByDirectReflectSet(params, pointer) {
		return nil
	}

	defer func() {
		// Catch the panic, especially the reflect operation panics.
		if exception := recover(); exception != nil {
			if e, ok := exception.(errorStack); ok {
				err = e
			} else {
				err = gerror.NewCodeSkipf(gerror.CodeInternalError, 1, "%v", exception)
			}
		}
	}()
	// If given `params` is JSON, it then uses json.Unmarshal doing the converting.
	switch r := params.(type) {
	case []byte:
		if json.Valid(r) {
			if rv, ok := pointer.(reflect.Value); ok {
				if rv.Kind() == reflect.Ptr {
					return json.UnmarshalUseNumber(r, rv.Interface())
				}
			} else {
				return json.UnmarshalUseNumber(r, pointer)
			}
		}
	case string:
		if paramsBytes := []byte(r); json.Valid(paramsBytes) {
			if rv, ok := pointer.(reflect.Value); ok {
				if rv.Kind() == reflect.Ptr {
					return json.UnmarshalUseNumber(paramsBytes, rv.Interface())
				}
			} else {
				return json.UnmarshalUseNumber(paramsBytes, pointer)
			}
		}
	}
	// Pointer type check.
	pointerRv, ok := pointer.(reflect.Value)
	if !ok {
		pointerRv = reflect.ValueOf(pointer)
		if kind := pointerRv.Kind(); kind != reflect.Ptr {
			return gerror.NewCodef(gerror.CodeInvalidParameter, "pointer should be type of pointer, but got: %v", kind)
		}
	}
	// Converting `params` to map slice.
	var (
		paramsList []interface{}
		paramsRv   = reflect.ValueOf(params)
		paramsKind = paramsRv.Kind()
	)
	for paramsKind == reflect.Ptr {
		paramsRv = paramsRv.Elem()
		paramsKind = paramsRv.Kind()
	}
	switch paramsKind {
	case reflect.Slice, reflect.Array:
		paramsList = make([]interface{}, paramsRv.Len())
		for i := 0; i < paramsRv.Len(); i++ {
			paramsList[i] = paramsRv.Index(i)
		}
	default:
		var paramsMaps = Maps(params)
		paramsList = make([]interface{}, len(paramsMaps))
		for i := 0; i < len(paramsMaps); i++ {
			paramsList[i] = paramsMaps[i]
		}
	}
	// If `params` is an empty slice, no conversion.
	if len(paramsList) == 0 {
		return nil
	}
	var (
		reflectElemArray = reflect.MakeSlice(pointerRv.Type().Elem(), len(paramsList), len(paramsList))
		itemType         = reflectElemArray.Index(0).Type()
		itemTypeKind     = itemType.Kind()
		pointerRvElem    = pointerRv.Elem()
		pointerRvLength  = pointerRvElem.Len()
	)
	if itemTypeKind == reflect.Ptr {
		// Pointer element.
		for i := 0; i < len(paramsList); i++ {
			var tempReflectValue reflect.Value
			if i < pointerRvLength {
				// Might be nil.
				tempReflectValue = pointerRvElem.Index(i).Elem()
			}
			if !tempReflectValue.IsValid() {
				tempReflectValue = reflect.New(itemType.Elem()).Elem()
			}
			if err = doStruct(paramsList[i], tempReflectValue, mapping, priorityTag); err != nil {
				return err
			}
			reflectElemArray.Index(i).Set(tempReflectValue.Addr())
		}
	} else {
		// Struct element.
		for i := 0; i < len(paramsList); i++ {
			var tempReflectValue reflect.Value
			if i < pointerRvLength {
				tempReflectValue = pointerRvElem.Index(i)
			} else {
				tempReflectValue = reflect.New(itemType).Elem()
			}
			if err = doStruct(paramsList[i], tempReflectValue, mapping, priorityTag); err != nil {
				return err
			}
			reflectElemArray.Index(i).Set(tempReflectValue)
		}
	}
	pointerRv.Elem().Set(reflectElemArray)
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
