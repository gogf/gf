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

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	convertConfig = gconv.NewConvertConfig("gf.orm")
)

func ConvertConfig() *gconv.ConvertConfig {
	return convertConfig
}

func reflectTypeFor[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func init() {
	convertConfig.RegisterDefaultConvertFuncs()

	convertConfig.RegisterTypeConvertFunc(reflectTypeFor[[]string](), convertToSliceFunc)
	convertConfig.RegisterTypeConvertFunc(reflectTypeFor[[]float32](), convertToSliceFunc)
	convertConfig.RegisterTypeConvertFunc(reflectTypeFor[[]float32](), convertToSliceFunc)
	convertConfig.RegisterTypeConvertFunc(reflectTypeFor[[]int64](), convertToSliceFunc)
	convertConfig.RegisterTypeConvertFunc(reflectTypeFor[map[string]any](), convertToSliceFunc)

	convertConfig.RegisterInterfaceTypeConvertFunc(reflectTypeFor[sql.Scanner](), sqlScanner)
}

func convertToSliceFunc(from any, to reflect.Value) (err error) {
	fromVal := from.(*gvar.Var).Val()
	switch x := fromVal.(type) {
	case []byte:
		dst := to.Addr().Interface()
		err = json.Unmarshal(x, dst)
	case string:
		dst := to.Addr().Interface()
		err = json.Unmarshal([]byte(x), dst)
	default:
		fromType := reflect.TypeOf(fromVal)
		switch fromType.Kind() {
		case reflect.Slice:
			dv := gconv.Convert(fromVal, to.Type().String())
			to.Set(reflect.ValueOf(dv))
		default:
			err = fmt.Errorf("conversion from `%v(%T)` to `%v(%T)` is not supported", fromVal, fromVal, to, to)
		}
	}
	return err
}

func sqlScanner(from any, to reflect.Value) error {
	dv := to.Addr().Interface()
	scanner, ok := dv.(sql.Scanner)
	if ok {
		return scanner.Scan(from)
	}
	return fmt.Errorf("type [%T] does not implement [sql.Scanner]", dv)
}
