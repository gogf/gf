// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

func registerFieldConvertFunc(ctx context.Context, db DB,
	tableField *sql.ColumnType, structField reflect.StructField) (convertFn fieldConvertFunc) {
	convertFn = getGoTypeConvertFunc(structField.Type)
	if convertFn != nil {
		return convertFn
	}
	convertFn = fromDriverGetFieldConvFunc(db.GetConfig().Type, tableField.DatabaseTypeName())
	if convertFn != nil {
		return convertFn
	}
	localType, _ := db.CheckLocalTypeForField(ctx, tableField.DatabaseTypeName(), nil)
	// There are several special types that require special handling
	switch localType {
	case LocalTypeUint64Bytes:
		// mysql bit
		convertFn = getBitConvertFunc(structField.Type, 0)
	case LocalTypeInt64Bytes:
		// mysql bit
		convertFn = getBitConvertFunc(structField.Type, 0)
	case LocalTypeDecimal:
		// decimal numeric money
		convertFn = getDecimalConvertFunc(structField.Type, 0)
	default:
		convertFn, _ = getConverter(structField.Type, 0)
	}
	if convertFn == nil {
		panic(&typeConvertError{
			driverName:  db.GetConfig().Type,
			columnName:  tableField.Name(),
			columnType:  tableField.DatabaseTypeName(),
			structField: structField,
		})
	}
	return
}

type typeConvertError struct {
	driverName  string
	columnName  string
	columnType  string
	structField reflect.StructField
}

func (t *typeConvertError) Error() string {
	err := `Driver: %s does not support conversion from (%s: %s) to (%s: %s)`
	return fmt.Sprintf(err, t.driverName, t.columnName, t.columnType, t.structField.Name, t.structField.Type)
}

type driverConvertFunc struct {
	driverName string
	// key = {bigint, text, char, json, ...}
	columnTypeConvertFunc map[string]fieldConvertFunc
}

func (d *driverConvertFunc) GetOrSetColumnTypeConvertFunc(columnType string, fn fieldConvertFunc) fieldConvertFunc {
	if d.columnTypeConvertFunc == nil {
		d.columnTypeConvertFunc = make(map[string]fieldConvertFunc)
	}
	columnType = strings.ToLower(columnType)
	convFunc, ok := d.columnTypeConvertFunc[columnType]
	if !ok {
		if fn != nil {
			d.columnTypeConvertFunc[columnType] = fn
		}
		return fn
	}
	return convFunc
}

var (
	// key = {mysql,mssql,oracle,pgsql, ...}
	customDriverFieldTypeConvertFuncs = map[string]*driverConvertFunc{}
	// key = {int,int8, others...}
	customGoTypeConvertFuncs = map[reflect.Type]fieldConvertFunc{}
)

func getDriverConvertFunc(driverName string, init bool) *driverConvertFunc {
	convertFunc, ok := customDriverFieldTypeConvertFuncs[driverName]
	if !ok && init {
		convertFunc = &driverConvertFunc{
			driverName: driverName,
		}
		customDriverFieldTypeConvertFuncs[driverName] = convertFunc
	}
	return convertFunc
}

func fromDriverGetFieldConvFunc(driverName, columnType string) fieldConvertFunc {
	if driverName == "" {
		return nil
	}
	driverConv := getDriverConvertFunc(driverName, false)
	if driverConv == nil {
		return nil
	}
	return driverConv.GetOrSetColumnTypeConvertFunc(columnType, nil)
}

// RegisterDatabaseConvertFunc
// Provide user-defined field conversion functions for smooth transitions
// driverName = {mysql,mssql,oracle,pgsql, ...}
// columnType = {bigint,datetime,char,text, ...}
func RegisterDatabaseConvertFunc(driverName, columnType string, fn fieldConvertFunc) {
	if useCacheTableExperiment == false {
		return
	}
	if fn == nil || checkStringIsEmpty(driverName, columnType) {
		panic(fmt.Errorf("parameter cannot be empty"))
	}
	databaseConv := getDriverConvertFunc(driverName, true)
	databaseConv.GetOrSetColumnTypeConvertFunc(columnType, fn)
}

// RegisterGoTypeConvertFunc
// Registering Field Conversion Functions for Go Language Types
// For example, if a field in the structure is of a third-party library type,
// it is not possible to implement [sql.Scanner]
// You can use this function to register a field conversion function,
// which has a higher priority than [RegisterBaseConvertFunc]
func RegisterGoTypeConvertFunc(goType reflect.Type, fn fieldConvertFunc) {
	if useCacheTableExperiment == false {
		return
	}
	if goType == nil || fn == nil {
		panic(fmt.Errorf("parameter cannot be empty"))
	}
	addGoTypeConvertFunc(goType, fn)
}

func addGoTypeConvertFunc(goType reflect.Type, fn fieldConvertFunc) {
	elemType, deref := getElemType(goType)
	_, ok := customGoTypeConvertFuncs[elemType]
	if ok {
		panic(fmt.Errorf("repeatedly registering conversion functions for go type [%v]", goType))
	}
	customGoTypeConvertFuncs[elemType] = getPtrConvFunc(deref, fn)
}

func getGoTypeConvertFunc(goType reflect.Type) fieldConvertFunc {
	elemType, _ := getElemType(goType)
	conv := customGoTypeConvertFuncs[elemType]
	return conv
}

func getElemType(typ reflect.Type) (reflect.Type, int) {
	deref := 0
	for {
		if typ.Kind() != reflect.Ptr {
			break
		}
		deref++
		typ = typ.Elem()
	}
	return typ, deref
}

// If you are registering * int, you need to wrap it with a pointer conversion function
// Users ensure that the type is non pointer type before operation
// The same goes for multi-level pointers
func getPtrConvFunc(deref int, fn fieldConvertFunc) fieldConvertFunc {
	if deref > 0 {
		return getPtrConvFunc(deref-1, ptrConverter(fn))
	}
	return fn
}

func checkStringIsEmpty(strs ...string) bool {
	for _, s := range strs {
		if strings.TrimSpace(s) == "" {
			return true
		}
	}
	return false
}
