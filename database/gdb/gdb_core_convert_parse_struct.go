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
)

func EnableCacheTableExperiment(b bool) {
	useCacheTableExperiment = b
}

type convertTableInfo struct {
	// key   = go type
	// value = *Table
	tablesMap sync.Map
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
	c.tablesMap.Delete(getTableName(structType))
}

func getTableName(pointerType reflect.Type) reflect.Type {
	if pointerType.Kind() == reflect.Ptr {
		pointerType = pointerType.Elem()
	}
	return pointerType
}

type scanPointer struct {
	// True only when Scan is called
	scan    bool
	pointer any
}

type Table struct {
	// tableFields
	fields []*fieldConvertInfo
	// Check if the [iUnmanshalValue] interface is implemented.
	// If it is, call it directly
	unmarshalValue fieldConvertFunc
}

func (t *Table) ScanToSlice(elemType, sliceType reflect.Type, records []Record, structIsPtr bool) (sliceValue reflect.Value, err error) {
	sliceValue = reflect.MakeSlice(sliceType, 0, len(records))
	for _, record := range records {
		var structValue reflect.Value
		structValue, err = t.ScanToStruct(elemType, record, structIsPtr)
		if err != nil {
			return
		}
		sliceValue = reflect.Append(sliceValue, structValue)
	}
	return
}

func (t *Table) ScanToStruct(elemType reflect.Type, record Record, secondPtr bool) (structValue reflect.Value, err error) {
	structValue = reflect.New(elemType)
	if t.unmarshalValue != nil {
		err = t.unmarshalValue(structValue, record)
		if err != nil {
			return
		}
		if secondPtr == false {
			structValue = structValue.Elem()
		}
		return
	}
	structValue = structValue.Elem()
	for _, field := range t.fields {
		// This field does not exist in the structure, for example,
		// when querying some databases, an additional column may be added, such as MSSQL or Oracle
		if field.convertFunc == nil {
			continue
		}
		fieldValue := field.GetReflectValue(structValue)
		value := record[field.ColumnFieldName]
		if value == nil {
			continue
		}
		fieldValue.Set(reflect.ValueOf(value.Val()))
	}
	// If the user passes a second level pointer dst=[**struct],
	// then when creating the struct,
	// a layer of reference has already been resolved,
	// which is a value type and needs to be converted to [*struct] again.
	// Then dst resolves a layer of reference to become a first level pointer
	// that matches the current struct exactly
	if secondPtr {
		structValue = structValue.Addr()
	}
	return
}

func (t *Table) GetFieldInfo(fieldName string) *fieldConvertInfo {
	for _, field := range t.fields {
		if field.StructField.Name == fieldName {
			return field
		}
	}
	return nil
}

func checkTypeImplUnmarshalValue(typ reflect.Type) fieldConvertFunc {
	// 1.[]*struct
	// 2.[]struct
	// 3.*struct
	switch typ.Kind() {
	case reflect.Slice, reflect.Array:
		// 1.[]*struct  => *struct
		// 2.[]struct   => struct
		typ = typ.Elem()
	case reflect.Ptr:
		// 3.*struct    => struct
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PointerTo(typ)
	}
	if typ.Implements(unmarshalValueType) {
		return func(dest reflect.Value, src any) error {
			fn, ok := dest.Interface().(iUnmarshalValue)
			if !ok {
				return fmt.Errorf("not implements iUnmarshalValue")
			}
			return fn.UnmarshalValue(src)
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
	// Used to check whether the type implements the interface.
	// If it does, the interface can be directly called during subsequent assignment
	unmarshalValueFunc := checkTypeImplUnmarshalValue(pointerType)

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
		table = &Table{
			unmarshalValue: unmarshalValueFunc,
		}
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
			// todo  empty interface ?
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
			convertFn := registerFieldConvertFunc(ctx, db, fieldInfo.ColumnFieldType, fieldInfo.StructField)
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
