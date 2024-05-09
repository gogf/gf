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

type Table struct {
	// tableFields
	fields []*fieldConvertInfo
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

	var (
		pointer     = val.pointer
		pointerType = reflect.TypeOf(pointer).Elem()
	)

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
		fieldsMap = make(map[string]*fieldConvertInfo)
	)
	for i := 0; i < len(columnTypes); i++ {
		column := columnTypes[i]
		fieldsMap[column.Name()] = &fieldConvertInfo{
			ColumnFieldIndex: i,
			ColumnFieldType:  column,
			ColumnFieldName:  column.Name(),
		}
	}

	var (
		table         = &Table{}
		matchedColumn = make(map[string]struct{})
		matchedCount  = table.getStructFields(ctx, db, fieldsMap, pointerType, []int{}, matchedColumn)
	)

	if matchedCount == 0 {
		return nil
	}

	table.fields = make([]*fieldConvertInfo, len(columnTypes))
	for _, v := range fieldsMap {
		table.fields[v.ColumnFieldIndex] = v
	}

	tablesMap.Store(tableName, table)
	return table
}

func (t *Table) getStructFields(ctx context.Context, db DB, fieldsMap map[string]*fieldConvertInfo, structType reflect.Type, parentIndex []int, matchedColumn map[string]struct{}) (matchedCount int) {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.IsExported() == false {
			continue
		}
		if field.Type.Kind() == reflect.Interface {
			continue
		}
		// gmeta.Meta
		if field.Type.String() == "gmeta.Meta" {
			continue
		}

		tag := field.Tag.Get(OrmTagForStruct)
		if field.Anonymous && tag == "" {
			if field.Type.Kind() == reflect.Ptr {
				field.Type = field.Type.Elem()
			}
			matchedCount += t.getStructFields(ctx, db, fieldsMap, field.Type, append(parentIndex, i), matchedColumn)
			continue
		}

		fieldInfo := t.parseTagAndMatchColumn(field.Tag, field.Name, fieldsMap, matchedColumn)

		if fieldInfo != nil {
			fieldInfo.StructFieldIndex = append(parentIndex, i)
			fieldInfo.StructFieldType = field.Type
			fieldInfo.StructField = field
			convertFn := registerFieldConvertFunc(ctx, db, fieldInfo.ColumnFieldType, fieldInfo.StructField)
			fieldInfo.convertFunc = convertFn
			matchedCount++
		}
	}
	return
}

func (t *Table) parseTagAndMatchColumn(fieldTag reflect.StructTag, fieldName string, fieldsMap map[string]*fieldConvertInfo, matchedColumn map[string]struct{}) *fieldConvertInfo {
	tag := fieldTag.Get("orm")
	// If there is with, skip it directly
	// type User struct {
	//	 gmeta.Meta `orm:"table:user"`
	//	 Id         int           `json:"id"`
	//	 Name       string        `json:"name"`
	//	 UserDetail *UserDetail   `orm:"with:uid=id"`
	//	 UserScores []*UserScores `orm:"with:uid=id"`
	// }
	if strings.Contains(tag, "with:") {
		// tag = ""
		return nil
	}
	// json
	if tag == "" {
		tag = fieldTag.Get("json")
	}
	if tag != "" {
		// json:"name,omitempty"
		// json:"-"
		// json:",omitempty"
		// orm:"id,omitempty"
		tag = strings.Split(tag, ",")[0]
		tag = strings.TrimSpace(tag)
		if tag == "-" {
			return nil
		}
	}

	if tag != "" {
		fieldInfo, ok := fieldsMap[tag]
		if ok {
			matchedColumn[tag] = struct{}{}
			return fieldInfo
		}
	}
	// There may not be a match to the tag
	fieldInfo, ok := fieldsMap[fieldName]
	if ok {
		matchedColumn[fieldName] = struct{}{}
		return fieldInfo
	}

	// Neither the tag nor the field name matched
	removeSymbolsFieldName := utils.RemoveSymbols(fieldName)
	for columnName, fieldInfo := range fieldsMap {
		if _, matched := matchedColumn[columnName]; matched {
			continue
		}
		removeSymbolsColumnName := utils.RemoveSymbols(columnName)
		if strings.EqualFold(removeSymbolsFieldName, removeSymbolsColumnName) {
			matchedColumn[columnName] = struct{}{}
			return fieldInfo
		}
	}
	return nil
}
