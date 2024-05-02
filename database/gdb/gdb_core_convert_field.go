// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"database/sql"
	"reflect"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/internal/utils"
)

var (
	tablesMap sync.Map
)

const (
	scanPointerCtxKey = "gf.orm.scan.ctx.key"
)

func getTableName(pointerType reflect.Type) string {
	return pointerType.PkgPath() + "." + pointerType.Name()
}

type scanPointer struct {
	// True only when Scan is called
	scan    bool
	pointer any
}

type fieldConvertFunc = func(src any, dst reflect.Value) error

type fieldConvertInfo struct {
	// table field
	ColumnField      string
	ColumnFieldIndex int
	ColumnFieldType  *sql.ColumnType
	// struct field
	StructField      reflect.StructField
	StructFieldType  reflect.Type
	StructFieldIndex []int
	convertFunc      fieldConvertFunc
	// flag
	isCustomConvert bool
}

// GetReflectValue  reflect.Value.FieldByIndex
func (c *fieldConvertInfo) GetReflectValue(structValue reflect.Value) reflect.Value {
	if len(c.StructFieldIndex) == 1 {
		return structValue.Field(c.StructFieldIndex[0])
	}
	v := structValue
	for i, x := range c.StructFieldIndex {
		if i > 0 {
			if v.Kind() == reflect.Pointer {
				if v.IsNil() {
					v.Set(reflect.New(v.Type().Elem()))
				}
				v = v.Elem()
			}
		}
		v = v.Field(x)
	}
	return v
}

type Table struct {
	// tableField
	fields map[string]*fieldConvertInfo
	// Temporary parameters used to receive the driver
	scanArgs []any
}

func parseStruct(ctx context.Context, db DB, columnTypes []*sql.ColumnType) *Table {
	ctxKey := ctx.Value(scanPointerCtxKey)
	if ctxKey == nil {
		return nil
	}
	val := ctxKey.(*scanPointer)
	if val.scan == false {
		return nil
	}
	pointer := val.pointer

	pointerType := reflect.TypeOf(pointer).Elem()
	switch pointerType.Kind() {
	case reflect.Array, reflect.Slice:
		// 1.[]*struct => *struct
		// 2.[]struct => struct
		pointerType = pointerType.Elem()
		if pointerType.Kind() == reflect.Ptr {
			pointerType = pointerType.Elem()
		}
	case reflect.Ptr: // **struct
		pointerType = pointerType.Elem()
	}

	tableName := getTableName(pointerType)
	tableValue, ok := tablesMap.Load(tableName)
	if ok {
		return tableValue.(*Table)
	}

	var (
		scanArgs   = make([]any, len(columnTypes))
		fieldsInfo = make(map[string]*fieldConvertInfo)
	)
	for i := 0; i < len(columnTypes); i++ {
		column := columnTypes[i]
		fieldsInfo[column.Name()] = &fieldConvertInfo{
			ColumnFieldIndex: i,
			ColumnFieldType:  column,
		}
	}

	table := &Table{
		fields:   fieldsInfo,
		scanArgs: scanArgs,
	}

	var (
		existsColumn = make(map[string]struct{})
		scanCount    = table.getStructFields(ctx, db, pointerType, []int{}, existsColumn)
	)

	if scanCount == 0 {
		return nil
	}
	// Mainly mssql needs to be compared, and the others are not
	// When mssql uses a query with limit order, one more field will come out.
	for colName, field := range table.fields {
		_, ok := existsColumn[colName]
		if !ok {
			delete(table.fields, colName)
			// scanArgs do not need to be deleted
			table.scanArgs[field.ColumnFieldIndex] = new(any)
		}
	}

	for i := 0; i < len(table.scanArgs); i++ {
		if scanArgs[i] == nil {
			table.scanArgs[i] = &sql.RawBytes{}
		}
	}

	tablesMap.Store(tableName, table)
	return table
}

func (t *Table) getStructFields(ctx context.Context, db DB, structType reflect.Type, parentIndex []int, existsColumn map[string]struct{}) (scanCount int) {

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.IsExported() == false {
			continue
		}
		if field.Type.Kind() == reflect.Interface {
			// empty interface
			if field.Type.NumMethod() != 0 {
				continue
			}
		}
		// g.Meta
		if field.Type.String() == "gmeta.Meta" {
			continue
		}
		tag := field.Tag.Get(OrmTagForStruct)
		if field.Anonymous && tag == "" {
			if field.Type.Kind() == reflect.Ptr {
				field.Type = field.Type.Elem()
			}
			scanCount += t.getStructFields(ctx, db, field.Type, append(parentIndex, i), existsColumn)
			continue
		}

		// orm:"with:id1=id2" json:"name"
		if strings.Contains(tag, "with:") {
			tag = ""
		}
		// json
		if tag == "" {
			tag = field.Tag.Get("json")
		}
		if tag != "" {
			// json:"name,omitempty"
			// json:"-"
			// json:",omitempty"
			// orm:"id,omitempty"
			tag = strings.Split(tag, ",")[0]
			tag = strings.TrimSpace(tag)
			if tag == "-" {
				continue
			}
		}

		var (
			fieldInfo *fieldConvertInfo
			ok        bool
		)

		if tag != "" {
			fieldInfo, ok = t.fields[tag]
			if !ok {
				// There may not be a match to the tag
				fieldInfo, ok = t.fields[field.Name]
				if ok {
					tag = field.Name
				}
			}
		}

		// Neither the tag nor the field name matched
		if !ok {
			removeSymbolsFieldName := utils.RemoveSymbols(field.Name)
			for columnName, structField := range t.fields {
				if _, exists := existsColumn[columnName]; exists {
					continue
				}
				removeSymbolsColumnName := utils.RemoveSymbols(columnName)
				if strings.EqualFold(removeSymbolsFieldName, removeSymbolsColumnName) {
					tag = columnName
					ok = true
					fieldInfo = structField
					existsColumn[columnName] = struct{}{}
					break
				}
			}
		}

		if ok {
			fieldInfo.StructFieldIndex = append(parentIndex, i)
			fieldInfo.StructFieldType = field.Type
			fieldInfo.StructField = field
			var (
				convertFn       fieldConvertFunc
				tempArg         = any(&sql.RawBytes{})
				isCustomConvert = true
			)
			// Check the custom translation interface implementation
			convertFn, tempArg = checkFieldImplConvertInterface(field)
			if convertFn == nil {
				isCustomConvert = false
				convertFn, tempArg = RegisterFieldConverterFunc(ctx, db, fieldInfo.ColumnFieldType, fieldInfo.StructField)
			}
			fieldInfo.isCustomConvert = isCustomConvert
			fieldInfo.convertFunc = convertFn
			t.scanArgs[fieldInfo.ColumnFieldIndex] = tempArg
			scanCount++
			existsColumn[tag] = struct{}{}
		}
	}
	return
}
