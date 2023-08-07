// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"reflect"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

var customConverters map[reflect.Type]map[reflect.Type]reflect.Value

func init() {
	customConverters = make(map[reflect.Type]map[reflect.Type]reflect.Value)
}

// RegisterConverter to register custom converter.
// It must be register before you use gconv. So suggest to do it in boot.
// Note:
//  1. The fn must be func(T1)(T2,error). It will convert T1 to T2.
//  2. The T1 and T2 must be pointer.
func RegisterConverter(fn interface{}) (err error) {
	fnReflectValue := reflect.ValueOf(fn)
	fnReflectType := fnReflectValue.Type()
	errType := reflect.TypeOf((*error)(nil)).Elem()

	if fnReflectType.Kind() != reflect.Func ||
		fnReflectType.NumIn() != 1 || fnReflectType.NumOut() != 2 ||
		!fnReflectType.Out(1).Implements(errType) {
		err = gerror.NewCode(gcode.CodeInvalidParameter, "The gconv.RegisterConverter's parameter must be a function as func(T1)(T2,error).")
		return
	}

	inType := fnReflectType.In(0)
	outType := fnReflectType.Out(0)

	subMap, ok := customConverters[inType]
	if !ok {
		subMap = make(map[reflect.Type]reflect.Value)
		customConverters[inType] = subMap
	}

	if _, ok := subMap[outType]; ok {
		err = gerror.NewCode(gcode.CodeOperationFailed, "The converter has been registered.")
		return
	}

	subMap[outType] = reflect.ValueOf(fn)
	return
}

// callCustomConverter call the custom converter. It will try some possible type.
func callCustomConverter(reflectValue reflect.Value, pointerReflectValue reflect.Value) (ok bool, err error) {
	if reflectValue.Kind() != reflect.Pointer && reflectValue.CanAddr() {
		reflectValue = reflectValue.Addr()
	}
	if pointerReflectValue.Kind() != reflect.Pointer && pointerReflectValue.CanAddr() {
		pointerReflectValue = pointerReflectValue.Addr()
	}
	for {
		if !reflectValue.IsValid() {
			break
		}
		subMap, ok := customConverters[reflectValue.Type()]
		if ok {
			pointerTmpReflectValue := pointerReflectValue
			for {
				if !pointerTmpReflectValue.IsValid() {
					break
				}
				if converter, ok := subMap[pointerTmpReflectValue.Type()]; ok {
					ret := converter.Call([]reflect.Value{reflectValue})
					if pointerTmpReflectValue.CanSet() {
						pointerTmpReflectValue.Set(ret[0])
					} else if pointerTmpReflectValue.Elem().CanSet() {
						pointerTmpReflectValue.Elem().Set(ret[0].Elem())
					}
					if ret[1].IsNil() {
						err = nil
					} else {
						err = ret[1].Interface().(error)
					}
					return true, err
				}
				if pointerTmpReflectValue.Kind() == reflect.Pointer {
					pointerTmpReflectValue = pointerTmpReflectValue.Elem()
				} else {
					break
				}
			}
		}
		if reflectValue.Kind() == reflect.Pointer {
			reflectValue = reflectValue.Elem()
		} else {
			break
		}
	}
	return false, nil
}
