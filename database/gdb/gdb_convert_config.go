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

func init() {
	convertConfig.RegisterDefaultConvertFuncs()

	convertConfig.RegisterTypeConvertFunc(reflect.TypeFor[[]string](), convertToSliceFunc)
	convertConfig.RegisterTypeConvertFunc(reflect.TypeFor[[]float32](), convertToSliceFunc)
	convertConfig.RegisterTypeConvertFunc(reflect.TypeFor[[]float32](), convertToSliceFunc)
	convertConfig.RegisterTypeConvertFunc(reflect.TypeFor[[]int64](), convertToSliceFunc)
	convertConfig.RegisterTypeConvertFunc(reflect.TypeFor[map[string]any](), convertToSliceFunc)

	convertConfig.RegisterInterfaceTypeConvertFunc(reflect.TypeFor[sql.Scanner](), sqlScanner)
}

func convertToSliceFunc(from any, to reflect.Value) error {
	dst := to.Addr().Interface()
	sv := from.(*gvar.Var).Bytes()
	err := json.Unmarshal(sv, dst)
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
