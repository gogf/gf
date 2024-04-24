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

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	gfOrmScanToStructPointerCtxKey = "gf.orm.scan_to_struct_pointer.ctx.key"
)

func (m *Model) one(pointer any, where ...interface{}) error {
	var ctx = m.GetCtx()
	//if len(where) > 0 {
	//	return m.Where(where[0], where[1:]...).one(pointer)
	//}
	ctx = context.WithValue(ctx, gfOrmScanToStructPointerCtxKey, pointer)
	err := m.doScanToPointer(ctx, true)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) all(pointer any, where ...interface{}) error {
	var ctx = m.GetCtx()
	//if len(where) > 0 {
	//	return m.Where(where[0], where[1:]...).all(pointer)
	//}
	ctx = context.WithValue(ctx, gfOrmScanToStructPointerCtxKey, pointer)
	err := m.doScanToPointer(ctx, false)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) doScanToPointer(ctx context.Context, limit1 bool) error {
	//if len(where) > 0 {
	//	return m.Where(where[0], where[1:]...).all(pointer)
	//}
	sqlWithHolder, holderArgs := m.getFormattedSqlAndArgs(ctx, queryTypeNormal, limit1)
	return m.doGetAllBySqlAndScanToPointer(ctx, queryTypeNormal, sqlWithHolder, holderArgs...)
}

func (m *Model) doGetAllBySqlAndScanToPointer(
	ctx context.Context,
	queryType queryType, sql string, args ...interface{}) (err error) {
	// todo 缓存的api是否也需要更改？
	if _, err = m.getSelectResultFromCache(ctx, sql, args...); err != nil {
		return
	}

	in := &HookSelectInput{
		internalParamHookSelect: internalParamHookSelect{
			internalParamHook: internalParamHook{
				link: m.getLink(false),
			},
			handler: m.hookHandler.Select,
		},
		Model: m,
		Table: m.tables,
		Sql:   sql,
		Args:  m.mergeArguments(args),
	}
	if err = in.Next2(ctx); err != nil {
		return
	}
	// todo 缓存的api是否也需要更改？
	var result Result
	err = m.saveSelectResultToCache(ctx, queryType, result, sql, args...)
	return
}

func (h *HookSelectInput) Next2(ctx context.Context) (err error) {
	structPointer := ctx.Value(gfOrmScanToStructPointerCtxKey)
	if structPointer == nil {
		panic("请考虑调用Scan方法来查询数据库")
	}

	if h.originalTableName.IsNil() {
		h.originalTableName = gvar.New(h.Table)
	}
	if h.originalSchemaName.IsNil() {
		h.originalSchemaName = gvar.New(h.Schema)
	}
	// Custom hook handler call.
	if h.handler != nil && !h.handlerCalled {
		h.handlerCalled = true
		_, err = h.handler(ctx, h)
		return err
	}
	var toBeCommittedSql = h.Sql
	// Table change.
	if h.Table != h.originalTableName.String() {
		toBeCommittedSql, err = gregex.ReplaceStringFuncMatch(
			`(?i) FROM ([\S]+)`,
			toBeCommittedSql,
			func(match []string) string {
				charL, charR := h.Model.db.GetChars()
				return fmt.Sprintf(` FROM %s%s%s`, charL, h.Table, charR)
			},
		)
	}
	// Schema change.
	if h.Schema != "" && h.Schema != h.originalSchemaName.String() {
		h.link, err = h.Model.db.GetCore().SlaveLink(h.Schema)
		if err != nil {
			return
		}
	}
	_, err = h.Model.db.DoSelect(ctx, h.link, toBeCommittedSql, h.Args...)
	return err
}

type columnToStructField struct {
	// 可能是嵌套的
	structFieldIndex []int
	structFieldType  string
	//=====
	columnIndex int
	columnType  *sql.ColumnType
}

// RowsScanToPointer converts underlying data record type sql.Rows to Result type.
func (c *Core) RowsScanToPointer(ctx context.Context, pointer any, rows *sql.Rows) (int64, error) {

	if rows == nil {
		return 0, nil
	}
	defer func() {
		if err := rows.Close(); err != nil {
			intlog.Errorf(ctx, `%+v`, err)
		}
	}()
	if !rows.Next() {
		return 0, nil
	}
	// Column names and types.
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return 0, err
	}

	if len(columnTypes) > 0 {
		if internalData := c.getInternalColumnFromCtx(ctx); internalData != nil {
			internalData.FirstResultColumn = columnTypes[0].Name()
		}
	}

	var getStructFields func(structType reflect.Type, parentIndex []int, columnNameToStructFieldMap map[string]columnToStructField)

	getStructFields = func(structType reflect.Type, parentIndex []int, columnNameToStructFieldMap map[string]columnToStructField) {
		for i := 0; i < structType.NumField(); i++ {

			field := structType.Field(i)
			if field.IsExported() == false || field.Type.Kind() == reflect.Interface {
				continue
			}
			if field.Anonymous {
				getStructFields(field.Type, append(parentIndex, i), columnNameToStructFieldMap)
				continue
			}

			// todo 1.如果为空，是否还需要根据其他tag，如json，conv之类的
			ormTag := field.Tag.Get(OrmTagForStruct)
			if ormTag != "" {
				v, ok := columnNameToStructFieldMap[ormTag]
				if ok {
					v.structFieldIndex = append(parentIndex, i)
					v.structFieldType = field.Type.Name()
					columnNameToStructFieldMap[ormTag] = v
				}
			}
			// todo 2.是否需要循环遍历 columns ，去做模糊匹配 和每一个字段名
		}
	}

	var getColumnNameToStructFieldMap = func(pointerType reflect.Type, columnNameToStructFieldMap map[string]columnToStructField) {
		for i := 0; i < len(columnTypes); i++ {
			column := columnTypes[i]
			columnNameToStructFieldMap[column.Name()] = columnToStructField{
				columnIndex: i,
				columnType:  column,
			}
		}
		// *struct -> struct
		// *[]struct -> []struct
		if pointerType.Kind() == reflect.Ptr {
			pointerType = pointerType.Elem()
		}
		switch pointerType.Kind() {
		case reflect.Array, reflect.Slice:
			// 1.[]*struct => *struct
			// 2.[]struct => struct
			pointerType = pointerType.Elem()
			if pointerType.Kind() == reflect.Ptr {
				pointerType = pointerType.Elem()
			}
		case reflect.Struct:

		}
		getStructFields(pointerType, []int{}, columnNameToStructFieldMap)

	}

	// 判断pointer 是*struct 还是*[]*struct
	// 前面已经判断过指针了，这里直接解引用
	structType := reflect.TypeOf(pointer).Elem()

	var (
		values                     = make([]interface{}, len(columnTypes))
		scanArgs                   = make([]interface{}, len(values))
		columnNameToStructFieldMap = make(map[string]columnToStructField)
		//=================
		ptr         = reflect.ValueOf(pointer).Elem()
		result_rows = int64(1)
	)

	getColumnNameToStructFieldMap(structType, columnNameToStructFieldMap)

	for i := range values {
		scanArgs[i] = &values[i]
	}

	switch ptr.Kind() {
	case reflect.Slice, reflect.Array: // slice array
		// []struct []*struct
		sliceStructValue, err := c.rowsConvertToSliceStruct(ctx, rows, structType, scanArgs, columnNameToStructFieldMap)
		if err != nil {
			return 0, err
		}
		ptr.Set(sliceStructValue)
		result_rows = int64(sliceStructValue.Len())

		// kind = []*map[string]any []map[string]any
		// kind = *map[string]any
	case reflect.Struct:
		structValue, err := c.rowsConvertToStruct(ctx, rows, structType, scanArgs, columnNameToStructFieldMap)
		if err != nil {
			return 0, err
		}
		ptr.Set(structValue)
	}
	return result_rows, nil
}

func (c *Core) rowsConvertToSliceStruct(
	ctx context.Context, rows *sql.Rows,
	sliceType reflect.Type,
	scanArgs []any,
	columnNameToStructFieldMap map[string]columnToStructField,
) (sliceStructValue reflect.Value, err error) {

	sliceStruct := reflect.MakeSlice(sliceType, 0, 4)

	// []struct -> struct
	// []*struct -> *struct
	structType := sliceType.Elem()
	deref := false
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
		deref = true
	}
	// todo 直接提前缓存一份所有字段的从数据库字段类型到go语言类型的映射函数
	for {
		dest := reflect.New(structType).Elem()
		if err = rows.Scan(scanArgs...); err != nil {
			return sliceStructValue, err
		}

		for _, structField := range columnNameToStructFieldMap {
			dstField := dest.FieldByIndex(structField.structFieldIndex)
			// var convertedValue any
			columnValue := *(scanArgs[structField.columnIndex].(*any))
			if columnValue == nil {
				continue
			}

			convertedValue := gconv.Convert(columnValue, structField.structFieldType)
			dstField.Set(reflect.ValueOf(convertedValue))
		}
		if deref {
			dest = dest.Addr()
		}
		sliceStruct = reflect.Append(sliceStruct, dest)
		if !rows.Next() {
			break
		}
	}

	return sliceStruct, nil
}

func (c *Core) rowsConvertToStruct(
	ctx context.Context, rows *sql.Rows,
	structType reflect.Type,
	scanArgs []any,
	columnNameToStructFieldMap map[string]columnToStructField,
) (structValue reflect.Value, err error) {
	deref := false
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
		deref = true
	}

	dest := reflect.New(structType).Elem()

	for {

		if err = rows.Scan(scanArgs...); err != nil {
			return structValue, err
		}

		for _, structField := range columnNameToStructFieldMap {
			dstField := dest.FieldByIndex(structField.structFieldIndex)
			// var convertedValue any
			columnValue := *(scanArgs[structField.columnIndex].(*any))
			if columnValue == nil {
				continue
			}
			convertedValue := gconv.Convert(columnValue, structField.structFieldType)
			dstField.Set(reflect.ValueOf(convertedValue))
		}

		if !rows.Next() {
			break
		}
	}
	if deref {
		dest = dest.Addr()
	}

	return dest, nil
}
