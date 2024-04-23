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
)

func (m *Model) one(pointer any, where ...interface{}) error {
	var ctx = m.GetCtx()
	//if len(where) > 0 {
	//	return m.Where(where[0], where[1:]...).one(pointer)
	//}
	err := m.doScanToPointer(ctx, pointer, true)
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
	err := m.doScanToPointer(ctx, pointer, false)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) doScanToPointer(ctx context.Context, pointer any, limit1 bool) error {
	//if len(where) > 0 {
	//	return m.Where(where[0], where[1:]...).all(pointer)
	//}
	sqlWithHolder, holderArgs := m.getFormattedSqlAndArgs(ctx, queryTypeNormal, limit1)
	return m.doGetAllBySqlAndScanToPointer(ctx, pointer, queryTypeNormal, sqlWithHolder, holderArgs...)
}

func (m *Model) doGetAllBySqlAndScanToPointer(
	ctx context.Context,
	pointer any,
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
		Model:   m,
		Table:   m.tables,
		Sql:     sql,
		Args:    m.mergeArguments(args),
		Pointer: pointer,
	}
	if err = in.Next2(ctx); err != nil {
		return
	}
	// todo 缓存的api是否也需要更改？
	var result Result
	err = m.saveSelectResultToCache(ctx, queryType, result, sql, args...)
	return
}

// Next calls the next hook handler.
func (h *HookSelectInput) Next2(ctx context.Context) (err error) {
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
	err = h.Model.db.DoSelectAndScanToPointer(ctx, h.link, h.Pointer, toBeCommittedSql, h.Args...)
	return err
}

func (c *Core) DoSelectAndScanToPointer(ctx context.Context, link Link, pointer any, sql string, args ...interface{}) (err error) {
	return c.db.DoQueryAndScanToPointer(ctx, link, pointer, sql, args...)
}

// DoQueryScanStruct commits the sql string and its arguments to underlying driver
// through given link object and returns the execution result.
func (c *Core) DoQueryAndScanToPointer(ctx context.Context, link Link, pointer any, sql string, args ...interface{}) (err error) {
	// Transaction checks.
	if link == nil {
		if tx := TXFromCtx(ctx, c.db.GetGroup()); tx != nil {
			// Firstly, check and retrieve transaction link from context.
			link = &txLink{tx.GetSqlTX()}
		} else if link, err = c.SlaveLink(); err != nil {
			// Or else it creates one from master node.
			return err
		}
	} else if !link.IsTransaction() {
		// If current link is not transaction link, it checks and retrieves transaction from context.
		if tx := TXFromCtx(ctx, c.db.GetGroup()); tx != nil {
			link = &txLink{tx.GetSqlTX()}
		}
	}

	if c.db.GetConfig().QueryTimeout > 0 {
		ctx, _ = context.WithTimeout(ctx, c.db.GetConfig().QueryTimeout)
	}

	// Sql filtering.
	sql, args = c.FormatSqlBeforeExecuting(sql, args)
	sql, args, err = c.db.DoFilter(ctx, link, sql, args)
	if err != nil {
		return err
	}
	// SQL format and retrieve.
	if v := ctx.Value(ctxKeyCatchSQL); v != nil {
		var (
			manager      = v.(*CatchSQLManager)
			formattedSql = FormatSqlWithArgs(sql, args)
		)
		manager.SQLArray.Append(formattedSql)
		if !manager.DoCommit && ctx.Value(ctxKeyInternalProducedSQL) == nil {
			return nil
		}
	}
	// Link execution.
	_, err = c.db.DoCommit(ctx, DoCommitInput{
		Link:          link,
		Sql:           sql,
		Args:          args,
		Stmt:          nil,
		Type:          SqlTypeQueryContext,
		IsTransaction: link.IsTransaction(),
		Pointer:       pointer,
	})
	return err
}

type scanArgMappingToStructField struct {
	structFieldIndex int
	scanArgIndex     int
	columnType       *sql.ColumnType
}

// RowsToResult converts underlying data record type sql.Rows to Result type.
func (c *Core) RowsToResult(ctx context.Context, pointer any, rows *sql.Rows) (int64, error) {
	if pointer == nil {
		return 0, nil
	}
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

	var getScanArgMappingToStructFieldsMap = func(pointerType reflect.Type, scanArgMappingToStructFieldsMap map[string]scanArgMappingToStructField) {
		for i := 0; i < len(columnTypes); i++ {
			column := columnTypes[i]
			scanArgMappingToStructFieldsMap[column.Name()] = scanArgMappingToStructField{
				scanArgIndex: i,
				columnType:   column,
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
		for i := 0; i < pointerType.NumField(); i++ {

			fieldType := pointerType.Field(i)
			if fieldType.IsExported() == false {
				continue
			}
			// todo 1.如果为空，是否还需要根据json tag
			ormTag := fieldType.Tag.Get(OrmTagForStruct)
			if ormTag != "" {
				v, ok := scanArgMappingToStructFieldsMap[ormTag]
				if ok {
					// todo 需要判断当前字段是否是匿名字段，如果是匿名字段，索引是错误的
					v.structFieldIndex = fieldType.Index[0]
					scanArgMappingToStructFieldsMap[ormTag] = v
				}
			}
			// todo 2.是否需要循环遍历 columns ，去做模糊匹配 和每一个字段名
		}

	}

	// 判断pointer 是*struct 还是*[]*struct
	// 前面已经判断过指针了，这里直接解引用
	structType := reflect.TypeOf(pointer).Elem()

	var (
		values                          = make([]interface{}, len(columnTypes))
		scanArgs                        = make([]interface{}, len(values))
		scanArgMappingToStructFieldsMap = make(map[string]scanArgMappingToStructField)
		//=================
		ptr         = reflect.ValueOf(pointer).Elem()
		result_rows = int64(1)
	)

	getScanArgMappingToStructFieldsMap(structType, scanArgMappingToStructFieldsMap)

	for i := range values {
		scanArgs[i] = &values[i]
	}

	switch ptr.Kind() {
	case reflect.Slice, reflect.Array: // slice array
		// []struct []*struct
		sliceStructValue, err := c.rowsConvertToSliceStruct(ctx, rows, structType, scanArgs, scanArgMappingToStructFieldsMap)
		if err != nil {
			return 0, err
		}
		ptr.Set(sliceStructValue)
		result_rows = int64(sliceStructValue.Len())

		//  kind= []*map[string]any []map[string]any
		// kind= *map[string]any
	case reflect.Struct:
		structValue, err := c.rowsConvertToStruct(ctx, rows, structType, scanArgs, scanArgMappingToStructFieldsMap)
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
	scanArgMappingToStructFieldsMap map[string]scanArgMappingToStructField,
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
	// 假设to = string
	fieldConvertFunc := func(from any) (to any) {
		switch f := from.(type) {
		case []byte:
			return string(f)
		case string:
			return f
		case *[]byte:
			return string(*f)
		case *string:
			return *f
		}
		return ""
	}

	for {
		dest := reflect.New(structType).Elem()
		if err = rows.Scan(scanArgs...); err != nil {
			return sliceStructValue, err
		}

		for _, field := range scanArgMappingToStructFieldsMap {
			dstField := dest.Field(field.structFieldIndex)
			// var convertedValue any
			columnValue := *(scanArgs[field.scanArgIndex].(*any))
			if columnValue == nil {
				continue
			}
			//convertedValue, err := c.columnValueToLocalValue(ctx, columnValue, field.columnType)
			//if err != nil {
			//	return sliceStructValue, err
			//}
			convertedValue := fieldConvertFunc(columnValue)
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
	scanArgMappingToStructFieldsMap map[string]scanArgMappingToStructField,
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

		for _, field := range scanArgMappingToStructFieldsMap {
			dstField := dest.Field(field.structFieldIndex)
			// var convertedValue any
			columnValue := *(scanArgs[field.scanArgIndex].(*any))
			if columnValue == nil {
				continue
			}
			convertedValue, err := c.columnValueToLocalValue(ctx, columnValue, field.columnType)
			if err != nil {
				return structValue, err
			}
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
