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
	// todo 直接存储索引和字段信息，不再需要fieldsIndex
	fieldsMap map[string]*fieldConvertInfo
	// columnIndex => fieldsMapKey
	fieldsIndex map[int]string
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
		fieldsInfo = make(map[string]*fieldConvertInfo)
	)
	for i := 0; i < len(columnTypes); i++ {
		column := columnTypes[i]
		fieldsInfo[column.Name()] = &fieldConvertInfo{
			ColumnFieldIndex: i,
			ColumnFieldType:  column,
			ColumnField:      column.Name(),
		}
	}

	table := &Table{
		fieldsMap: fieldsInfo,
	}

	var (
		existsColumn = make(map[string]struct{})
		scanCount    = table.getStructFields(ctx, db, pointerType, []int{}, existsColumn)
	)

	if scanCount == 0 {
		return nil
	}

	table.fieldsIndex = make(map[int]string)
	for k, v := range table.fieldsMap {
		table.fieldsIndex[v.ColumnFieldIndex] = k
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
			continue
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
			fieldInfo, ok = t.fieldsMap[tag]
			if !ok {
				// There may not be a match to the tag
				fieldInfo, ok = t.fieldsMap[field.Name]
				if ok {
					tag = field.Name
				}
			}
		}

		// Neither the tag nor the field name matched
		if !ok {
			removeSymbolsFieldName := utils.RemoveSymbols(field.Name)
			for columnName, structField := range t.fieldsMap {
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
			convertFn := registerFieldConvertFunc(ctx, db, fieldInfo.ColumnFieldType, fieldInfo.StructField)
			fieldInfo.convertFunc = convertFn
			scanCount++
			existsColumn[tag] = struct{}{}
		}
	}
	return
}
