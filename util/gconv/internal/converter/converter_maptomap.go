// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package converter

import (
	"reflect"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// MapToMap converts any map type variable `params` to another map type variable `pointer`.
//
// The parameter `params` can be any type of map, like:
// map[string]string, map[string]struct, map[string]*struct, reflect.Value, etc.
//
// The parameter `pointer` should be type of *map, like:
// map[int]string, map[string]struct, map[string]*struct, reflect.Value, etc.
//
// The optional parameter `mapping` is used for struct attribute to map key mapping, which makes
// sense only if the items of original map `params` is type struct.
func (c *Converter) MapToMap(
	params, pointer any, mapping map[string]string, option ...MapOption,
) (err error) {
	var (
		paramsRv   reflect.Value
		paramsKind reflect.Kind
	)
	if v, ok := params.(reflect.Value); ok {
		paramsRv = v
	} else {
		paramsRv = reflect.ValueOf(params)
	}
	paramsKind = paramsRv.Kind()
	if paramsKind == reflect.Ptr {
		paramsRv = paramsRv.Elem()
		paramsKind = paramsRv.Kind()
	}
	if paramsKind != reflect.Map {
		m, err := c.Map(params, option...)
		if err != nil {
			return err
		}
		return c.MapToMap(m, pointer, mapping, option...)
	}
	// Empty params map, no need continue.
	if paramsRv.Len() == 0 {
		return nil
	}
	var pointerRv reflect.Value
	if v, ok := pointer.(reflect.Value); ok {
		pointerRv = v
	} else {
		pointerRv = reflect.ValueOf(pointer)
	}
	pointerKind := pointerRv.Kind()
	for pointerKind == reflect.Ptr {
		pointerRv = pointerRv.Elem()
		pointerKind = pointerRv.Kind()
	}
	if pointerKind != reflect.Map {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`destination pointer should be type of *map, but got: %s`,
			pointerKind,
		)
	}
	defer func() {
		// Catch the panic, especially the reflection operation panics.
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodeSkipf(gcode.CodeInternalPanic, 1, "%+v", exception)
			}
		}
	}()
	var (
		paramsKeys       = paramsRv.MapKeys()
		pointerKeyType   = pointerRv.Type().Key()
		pointerValueType = pointerRv.Type().Elem()
		pointerValueKind = pointerValueType.Kind()
		dataMap          = reflect.MakeMapWithSize(pointerRv.Type(), len(paramsKeys))
		mapOption        = c.getMapOption(option...)
		convertOption    = ConvertOption{
			StructOption: StructOption{ContinueOnError: mapOption.ContinueOnError},
			SliceOption:  SliceOption{ContinueOnError: mapOption.ContinueOnError},
			MapOption:    mapOption,
		}
	)
	// Retrieve the true element type of target map.
	if pointerValueKind == reflect.Ptr {
		pointerValueKind = pointerValueType.Elem().Kind()
	}
	for _, key := range paramsKeys {
		mapValue := reflect.New(pointerValueType).Elem()
		switch pointerValueKind {
		case reflect.Map, reflect.Struct:
			structOption := StructOption{
				ParamKeyToAttrMap: mapping,
				PriorityTag:       "",
				ContinueOnError:   mapOption.ContinueOnError,
			}
			if err = c.Struct(paramsRv.MapIndex(key).Interface(), mapValue, structOption); err != nil {
				return err
			}
		default:
			convertResult, err := c.doConvert(
				doConvertInput{
					FromValue:  paramsRv.MapIndex(key).Interface(),
					ToTypeName: pointerValueType.String(),
					ReferValue: mapValue,
				},
				convertOption,
			)
			if err != nil {
				return err
			}
			mapValue.Set(reflect.ValueOf(convertResult))
		}
		convertResult, err := c.doConvert(
			doConvertInput{
				FromValue:  key.Interface(),
				ToTypeName: pointerKeyType.Name(),
				ReferValue: reflect.New(pointerKeyType).Elem().Interface(),
			},
			convertOption,
		)
		if err != nil {
			return err
		}
		var mapKey = reflect.ValueOf(convertResult)
		dataMap.SetMapIndex(mapKey, mapValue)
	}
	pointerRv.Set(dataMap)
	return nil
}
