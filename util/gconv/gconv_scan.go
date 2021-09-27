// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"reflect"
)

// Scan automatically checks the type of `pointer` and converts `params` to `pointer`. It supports `pointer`
// with type of `*map/*[]map/*[]*map/*struct/**struct/*[]struct/*[]*struct` for converting.
//
// It calls function `doMapToMap`  internally if `pointer` is type of *map                 for converting.
// It calls function `doMapToMaps` internally if `pointer` is type of *[]map/*[]*map       for converting.
// It calls function `doStruct`    internally if `pointer` is type of *struct/**struct     for converting.
// It calls function `doStructs`   internally if `pointer` is type of *[]struct/*[]*struct for converting.
func Scan(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	var (
		pointerType reflect.Type
		pointerKind reflect.Kind
	)
	if v, ok := pointer.(reflect.Value); ok {
		pointerType = v.Type()
	} else {
		pointerType = reflect.TypeOf(pointer)
	}
	if pointerType == nil {
		return gerror.NewCode(gcode.CodeInvalidParameter, "parameter pointer should not be nil")
	}
	pointerKind = pointerType.Kind()
	if pointerKind != reflect.Ptr {
		return gerror.NewCodef(gcode.CodeInvalidParameter, "params should be type of pointer, but got type: %v", pointerKind)
	}
	var (
		pointerElem               = pointerType.Elem()
		pointerElemKind           = pointerElem.Kind()
		keyToAttributeNameMapping map[string]string
	)
	if len(mapping) > 0 {
		keyToAttributeNameMapping = mapping[0]
	}
	switch pointerElemKind {
	case reflect.Map:
		return doMapToMap(params, pointer, mapping...)

	case reflect.Array, reflect.Slice:
		var (
			sliceElem     = pointerElem.Elem()
			sliceElemKind = sliceElem.Kind()
		)
		for sliceElemKind == reflect.Ptr {
			sliceElem = sliceElem.Elem()
			sliceElemKind = sliceElem.Kind()
		}
		if sliceElemKind == reflect.Map {
			return doMapToMaps(params, pointer, mapping...)
		}
		return doStructs(params, pointer, keyToAttributeNameMapping, "")

	default:

		return doStruct(params, pointer, keyToAttributeNameMapping, "")
	}
}
