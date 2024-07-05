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

const (
	scanPointerCtxKey = "gf.orm.scan.ctx.key"
)

var (
	convTableInfo = &convertTableInfo{}
	//
	useCacheTableExperiment = true
	// Mainly used to call the [convTableInfo.Delete] function
	// Due to the large number of declared structures within functions in the testing environment,
	// For example:
	// 		Declare an A structure within the A function, regardless of any fields
	// 		Then declare a structure inside the B function, regardless of any fields
	// 		So the go language believes that the names of the structures within these two functions are the same
	// Will result in duplicate names when storing structures in [convTableInfo]
	// So during testing, after testing a function, it is necessary to delete the registered ones
	// During formal development, there is no need to delete it
	isTextEnvironment = false
)

func EnableCacheTableExperiment(b bool) {
	useCacheTableExperiment = b
}

type structTypeName = string
type structFieldName = string

type convertTableInfo struct {
	// key   = string
	// value = *Table
	tablesMap sync.Map
	// Mainly used by [RegisterStructFieldConvertFunc] to register custom conversions
	customStructFieldConvertFunc map[structTypeName]map[structFieldName]fieldConvertFunc
}

func (c *convertTableInfo) getStructFieldConvertFunc(structType reflect.Type, fieldName string) fieldConvertFunc {
	tableConv, ok := c.customStructFieldConvertFunc[getTableName(structType)]
	if !ok {
		return nil
	}
	fn := tableConv[fieldName]
	return fn
}

func (c *convertTableInfo) Get(structType reflect.Type) *Table {
	var (
		tableName = getTableName(structType)
	)
	tableValue, ok := c.tablesMap.Load(tableName)
	if ok {
		return tableValue.(*Table)
	}
	return nil
}

func (c *convertTableInfo) Add(structType reflect.Type, table *Table) {
	var (
		tableName = getTableName(structType)
	)
	c.tablesMap.Store(tableName, table)
}

func (c *convertTableInfo) Delete(structType reflect.Type) {
	var (
		tableName = getTableName(structType)
	)
	c.tablesMap.Delete(tableName)
}

func getTableName(pointerType reflect.Type) string {
	if pointerType.Kind() == reflect.Ptr {
		pointerType = pointerType.Elem()
	}
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

func (t *Table) GetFieldInfo(fieldName string) *fieldConvertInfo {
	for _, field := range t.fields {
		if field.StructField.Name == fieldName {
			return field
		}
	}
	return nil
}

func parseStruct(ctx context.Context, db DB, columnTypes []*sql.ColumnType) *Table {
	if useCacheTableExperiment == false {
		return nil
	}
	ctxKey := ctx.Value(scanPointerCtxKey)
	if ctxKey == nil {
		return nil
	}
	val := ctxKey.(*scanPointer)
	if val.scan == false {
		return nil
	}
	var (
		// 1.[]*struct  => *struct
		// 2.[]struct   => struct
		// 3.*[]*struct => []*struct
		// 4.*[]struct  => []struct
		// 5.*struct    => struct
		// 6.**struct   => *struct
		pointerType reflect.Type
	)
	pointer, ok := val.pointer.(reflect.Value)
	if ok {
		pointerType = pointer.Type().Elem()
	} else {
		pointerType = reflect.TypeOf(val.pointer).Elem()
	}
	switch pointerType.Kind() {
	case reflect.Array, reflect.Slice:
		// 1.[]*struct => *struct
		// 2.[]struct  => struct
		pointerType = pointerType.Elem()
		if pointerType.Kind() == reflect.Ptr {
			pointerType = pointerType.Elem()
		}
	case reflect.Ptr:
		// *struct => struct
		pointerType = pointerType.Elem()
	}

	tableValue := convTableInfo.Get(pointerType)
	if tableValue != nil {
		return tableValue
	}
	var (
		fieldsConvertInfoMap = make(map[string]*fieldConvertInfo)
	)
	for i := 0; i < len(columnTypes); i++ {
		column := columnTypes[i]
		fieldsConvertInfoMap[column.Name()] = &fieldConvertInfo{
			ColumnFieldIndex: i,
			ColumnFieldType:  column,
			ColumnFieldName:  column.Name(),
		}
	}

	var (
		table         = &Table{}
		matchedColumn = make(map[string]struct{})
		matchedCount  = table.getStructFields(ctx, db, fieldsConvertInfoMap, pointerType, []int{}, matchedColumn)
	)

	if matchedCount == 0 {
		return nil
	}

	table.fields = make([]*fieldConvertInfo, len(columnTypes))
	for _, v := range fieldsConvertInfoMap {
		table.fields[v.ColumnFieldIndex] = v
	}

	convTableInfo.Add(pointerType, table)
	return table
}

// parentIndex is the index of anonymous structures, or in other words, Nested index in Hello
// For example, Hello.A parentIndex is []int{0}, Hello.B parentIndex is also []int{0},
//
//	type Nested struct {
//		A int
//		B int
//	}
//	type Hello struct{
//		Nested
//		ID int
//	}
func (t *Table) getStructFields(ctx context.Context, db DB, fieldsConvertInfoMap map[string]*fieldConvertInfo,
	structType reflect.Type, parentIndex []int, matchedColumn map[string]struct{}) (matchedCount int) {
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
		if field.Anonymous && field.Tag == "" {
			if field.Type.Kind() == reflect.Ptr {
				field.Type = field.Type.Elem()
			}
			matchedCount += t.getStructFields(ctx, db, fieldsConvertInfoMap, field.Type, append(parentIndex, i), matchedColumn)
			continue
		}

		fieldInfo := t.parseTagAndMatchColumn(field.Tag, field.Name, fieldsConvertInfoMap, matchedColumn)

		if fieldInfo != nil {
			fieldInfo.StructFieldIndex = append(parentIndex, i)
			fieldInfo.StructFieldType = field.Type
			fieldInfo.StructField = field
			convertFn := registerFieldConvertFunc(ctx, db, fieldInfo.ColumnFieldType, fieldInfo.StructField, structType)
			fieldInfo.convertFunc = convertFn
			matchedCount++
		}
	}
	return
}

func (t *Table) parseTagAndMatchColumn(fieldTag reflect.StructTag, fieldName string,
	fieldsConvertInfoMap map[string]*fieldConvertInfo, matchedColumn map[string]struct{}) *fieldConvertInfo {
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
			tag = fieldName
		}
	}

	if tag != "" {
		fieldInfo, ok := fieldsConvertInfoMap[tag]
		if ok {
			matchedColumn[tag] = struct{}{}
			return fieldInfo
		}
	}
	// There may not be a match to the tag
	fieldInfo, ok := fieldsConvertInfoMap[fieldName]
	if ok {
		matchedColumn[fieldName] = struct{}{}
		return fieldInfo
	}

	// Neither the tag nor the field name matched
	removeSymbolsFieldName := utils.RemoveSymbols(fieldName)
	for columnName, fieldInfo := range fieldsConvertInfoMap {
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
