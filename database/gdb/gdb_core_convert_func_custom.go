// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"reflect"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

// iUnmarshalValue is the interface for custom defined types customizing value assignment.
// Note that only pointer can implement interface iUnmarshalValue.
type iUnmarshalValue interface {
	UnmarshalValue(val interface{}) error
}

func checkFieldImplConvertInterface(structField reflect.StructField) (fn fieldConvertFunc, scanArg any) {
	var (
		impl      = reflect.Value{}
		fieldType = structField.Type
		isptr     = false
		// gtime.Time
		isGtime = false
	)

	if fieldType.Kind() != reflect.Ptr {
		impl = reflect.New(fieldType)
		isGtime = fieldType == reflect.TypeOf(gtime.Time{})
	} else {
		impl = reflect.New(fieldType.Elem())
		isptr = true
		isGtime = fieldType == reflect.TypeOf(&gtime.Time{})
	}

	if isGtime {
		return nil, nil
	}
	// Differences in order may occur
	switch impl.Interface().(type) {
	case iUnmarshalValue:
		fn = getUnmarshalValueConvertFunc(isptr, impl.Elem().Type())
	case sql.Scanner:
		fn = getSqlScannerConvertFunc(isptr, impl.Elem().Type())
	}
	scanArg = &sql.RawBytes{}
	return
}

func getUnmarshalValueConvertFunc(isptr bool, typ reflect.Type) fieldConvertFunc {
	// The arguments of the custom type conversion function are all []byte, from sql.RawBytes
	if isptr == false {
		return func(src any, dst reflect.Value) error {
			fn, ok := dst.Addr().Interface().(iUnmarshalValue)
			// todo: If the conversion is successful, you can cancel the check
			if !ok {
				return gerror.Newf("custom Type: %v Conversion to Interface Type: %v failed", dst.Type(), "iUnmarshalValue")
			}
			v := *src.(*sql.RawBytes)
			return fn.UnmarshalValue([]byte(v))
		}
	}
	return func(src any, dst reflect.Value) error {
		if dst.IsNil() {
			dst.Set(reflect.New(typ))
		}
		fn, ok := dst.Interface().(iUnmarshalValue)
		// todo: If the conversion is successful, you can cancel the check
		if !ok {
			return gerror.Newf("custom Type: %v Conversion to Interface Type: %v failed", dst.Type(), "iUnmarshalValue")
		}
		v := *src.(*sql.RawBytes)
		return fn.UnmarshalValue([]byte(v))
	}
}

func getSqlScannerConvertFunc(isptr bool, typ reflect.Type) fieldConvertFunc {
	// The arguments of the custom type conversion function are all []byte, from sql.RawBytes
	return func(src any, dst reflect.Value) error {
		var (
			fn sql.Scanner
			ok bool
		)
		if isptr {
			if dst.IsNil() {
				dst.Set(reflect.New(typ))
			}
			fn, ok = dst.Interface().(sql.Scanner)
		} else {
			fn, ok = dst.Addr().Interface().(sql.Scanner)
		}
		// todo: If the conversion is successful, you can cancel the check
		if !ok {
			return gerror.Newf("custom Type: %v Conversion to Interface Type: %v failed", dst.Type(), "sql.Scanner")
		}
		v := *src.(*sql.RawBytes)
		return fn.Scan([]byte(v))
	}
}
