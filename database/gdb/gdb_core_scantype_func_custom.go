// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"fmt"
	"reflect"
)

// iUnmarshalValue is the interface for custom defined types customizing value assignment.
// Note that only pointer can implement interface iUnmarshalValue.
type iUnmarshalValue interface {
	UnmarshalValue(val interface{}) error
}

// 自定义类型转换函数的参数全部都是[]byte， 从sql.RawBytes转换而来

func checkFieldImplConvertInterface(structField reflect.StructField) (fn fieldScanFunc, arg any) {
	var impl = reflect.Value{}
	fieldType := structField.Type
	isptr := fieldType.Kind() == reflect.Ptr
	if fieldType.Kind() != reflect.Ptr {
		impl = reflect.New(fieldType)
	} else {
		impl = reflect.New(fieldType.Elem())
	}
	arg = impl.Interface()
	// 可能会导致顺序差异
	switch impl.Interface().(type) {
	case iUnmarshalValue:
		fn = getUnmarshalValueConvertFunc(isptr, impl.Elem().Type())
	case sql.Scanner:
		fn = getSqlScannerConvertFunc(isptr, impl.Elem().Type())
	}
	return
}

func getUnmarshalValueConvertFunc(isptr bool, typ reflect.Type) fieldScanFunc {

	if isptr == false {
		return func(src any, dst reflect.Value) error {

			fn, ok := dst.Addr().Interface().(iUnmarshalValue)
			if !ok {
				return fmt.Errorf("自定义类型:%v 转换到接口类型:%v 失败", dst.Type(), "iUnmarshalValue")
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
		if !ok {
			return fmt.Errorf("自定义类型:%v 转换到接口类型:%v 失败", dst.Type(), "iUnmarshalValue")
		}
		v := *src.(*sql.RawBytes)
		return fn.UnmarshalValue([]byte(v))
	}

}
func getSqlScannerConvertFunc(isptr bool, typ reflect.Type) fieldScanFunc {
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
		// todo 一定转换成功，可以取消检查
		if !ok {
			return fmt.Errorf("自定义类型:%v 转换到接口类型:%v 失败", dst.Type(), "sql.Scanner")
		}

		v := *src.(*sql.RawBytes)
		return fn.Scan([]byte(v))
	}

}
