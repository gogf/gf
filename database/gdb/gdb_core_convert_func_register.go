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
	tableField *sql.ColumnType, structField reflect.StructField, structType reflect.Type) (convertFn fieldConvertFunc) {
	convertFn = convTableInfo.getStructFieldConvertFunc(structType, structField.Name)
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
	// key = columnType
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
	driverConvertFuncs = map[string]*driverConvertFunc{}
)

func getDriverConvertFunc(driverName string, init bool) *driverConvertFunc {
	convertFunc, ok := driverConvertFuncs[driverName]
	if !ok && init {
		convertFunc = &driverConvertFunc{
			driverName: driverName,
		}
		driverConvertFuncs[driverName] = convertFunc
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

// RegisterStructFieldConvertFunc
// Registering Field Conversion Functions for Structures
// For example, if a field in the structure is of a third-party library type,
// it is not possible to implement [sql.Scanner]
// You can use this function to register a field conversion function,
// which has a higher priority than [RegisterBaseConvertFunc]
func RegisterStructFieldConvertFunc(structType reflect.Type, fieldName string, fn fieldConvertFunc) {
	if useCacheTableExperiment == false {
		return
	}
	tableConv, ok := convTableInfo.customStructFieldConvertFunc[getTableName(structType)]
	if !ok {
		tableConv = make(map[structFieldName]fieldConvertFunc)
		convTableInfo.customStructFieldConvertFunc[getTableName(structType)] = tableConv
	}
	tableConv[fieldName] = fn
}

func checkStringIsEmpty(strs ...string) bool {
	for _, s := range strs {
		if strings.TrimSpace(s) == "" {
			return true
		}
	}
	return false
}
