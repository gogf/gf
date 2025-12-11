// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

// SoftTimeType custom defines the soft time field type.
type SoftTimeType int

const (
	SoftTimeTypeAuto           SoftTimeType = 0 // (Default)Auto detect the field type by table field type.
	SoftTimeTypeTime           SoftTimeType = 1 // Using datetime as the field value.
	SoftTimeTypeTimestamp      SoftTimeType = 2 // In unix seconds.
	SoftTimeTypeTimestampMilli SoftTimeType = 3 // In unix milliseconds.
	SoftTimeTypeTimestampMicro SoftTimeType = 4 // In unix microseconds.
	SoftTimeTypeTimestampNano  SoftTimeType = 5 // In unix nanoseconds.
)

// SoftTimeOption is the option to customize soft time feature for Model.
type SoftTimeOption struct {
	SoftTimeType SoftTimeType // The value type for soft time field.
}

type softTimeMaintainer struct {
	*Model
}

// SoftTimeFieldType represents different soft time field purposes
type SoftTimeFieldType int

const (
	SoftTimeFieldCreate SoftTimeFieldType = iota
	SoftTimeFieldUpdate
	SoftTimeFieldDelete
)

type iSoftTimeMaintainer interface {
	// GetFieldInfo returns field name and type for specified field purpose
	GetFieldInfo(ctx context.Context, schema, table string, fieldType SoftTimeFieldType) (fieldName string, localType LocalType)

	// GetFieldValue generates value for create/update/delete operations
	GetFieldValue(ctx context.Context, localType LocalType, isDeleted bool) any

	// GetDeleteCondition returns WHERE condition for soft delete query
	GetDeleteCondition(ctx context.Context) string

	// GetDeleteData returns UPDATE statement data for soft delete
	GetDeleteData(ctx context.Context, prefix, fieldName string, localType LocalType) (holder string, value any)
}

// getSoftFieldNameAndTypeCacheItem is the internal struct for storing create/update/delete fields.
type getSoftFieldNameAndTypeCacheItem struct {
	FieldName string
	FieldType LocalType
}

var (
	// Default field names of table for automatic-filled for record creating.
	createdFieldNames = []string{"created_at", "create_at"}
	// Default field names of table for automatic-filled for record updating.
	updatedFieldNames = []string{"updated_at", "update_at"}
	// Default field names of table for automatic-filled for record deleting.
	deletedFieldNames = []string{"deleted_at", "delete_at"}
)

// SoftTime sets the SoftTimeOption to customize soft time feature for Model.
func (m *Model) SoftTime(option SoftTimeOption) *Model {
	model := m.getModel()
	model.softTimeOption = option
	return model
}

// Unscoped disables the soft time feature for insert, update and delete operations.
func (m *Model) Unscoped() *Model {
	model := m.getModel()
	model.unscoped = true
	return model
}

func (m *Model) softTimeMaintainer() iSoftTimeMaintainer {
	return &softTimeMaintainer{
		m,
	}
}

// GetFieldInfo returns field name and type for specified field purpose.
// It checks the key with or without cases or chars '-'/'_'/'.'/' '.
func (m *softTimeMaintainer) GetFieldInfo(
	ctx context.Context, schema, table string, fieldType SoftTimeFieldType,
) (fieldName string, localType LocalType) {
	// Check if feature is disabled
	if m.db.GetConfig().TimeMaintainDisabled {
		return "", LocalTypeUndefined
	}

	// Determine table name
	tableName := table
	if tableName == "" {
		tableName = m.tablesInit
	}

	// Get config and field candidates
	config := m.db.GetConfig()
	var (
		configField   string
		defaultFields []string
	)

	switch fieldType {
	case SoftTimeFieldCreate:
		configField = config.CreatedAt
		defaultFields = createdFieldNames
	case SoftTimeFieldUpdate:
		configField = config.UpdatedAt
		defaultFields = updatedFieldNames
	case SoftTimeFieldDelete:
		configField = config.DeletedAt
		defaultFields = deletedFieldNames
	}

	// Use config field if specified, otherwise use defaults
	if configField != "" {
		return m.getSoftFieldNameAndType(ctx, schema, tableName, []string{configField})
	}
	return m.getSoftFieldNameAndType(ctx, schema, tableName, defaultFields)
}

// getSoftFieldNameAndType retrieves and returns the field name of the table for possible key.
func (m *softTimeMaintainer) getSoftFieldNameAndType(
	ctx context.Context, schema, table string, candidates []string,
) (fieldName string, fieldType LocalType) {
	// Build cache key
	cacheKey := fmt.Sprintf(`soft_field:%s:%s:%s`,
		schema, table, strings.Join(candidates, ","))

	// Try to get from cache
	cache := m.db.GetCore().GetInnerMemCache()
	result, err := cache.GetOrSetFunc(ctx, cacheKey, func(ctx context.Context) (any, error) {
		// Get table fields
		fieldsMap, err := m.TableFields(table, schema)
		if err != nil || len(fieldsMap) == 0 {
			return nil, err
		}

		// Search for matching field
		for _, candidate := range candidates {
			if name := searchFieldNameFromMap(fieldsMap, candidate); name != "" {
				fType, _ := m.db.CheckLocalTypeForField(ctx, fieldsMap[name].Type, nil)
				return getSoftFieldNameAndTypeCacheItem{
					FieldName: name,
					FieldType: fType,
				}, nil
			}
		}
		return nil, nil
	}, gcache.DurationNoExpire)

	if err != nil || result == nil {
		return "", LocalTypeUndefined
	}

	item := result.Val().(getSoftFieldNameAndTypeCacheItem)
	return item.FieldName, item.FieldType
}

func searchFieldNameFromMap(fieldsMap map[string]*TableField, key string) string {
	if len(fieldsMap) == 0 {
		return ""
	}
	_, ok := fieldsMap[key]
	if ok {
		return key
	}
	key = utils.RemoveSymbols(key)
	for k := range fieldsMap {
		if strings.EqualFold(utils.RemoveSymbols(k), key) {
			return k
		}
	}
	return ""
}

// GetDeleteCondition returns WHERE condition for soft delete query.
// It supports multiple tables string like:
// "user u, user_detail ud"
// "user u LEFT JOIN user_detail ud ON(ud.uid=u.uid)"
// "user LEFT JOIN user_detail ON(user_detail.uid=user.uid)"
// "user u LEFT JOIN user_detail ud ON(ud.uid=u.uid) LEFT JOIN user_stats us ON(us.uid=u.uid)".
func (m *softTimeMaintainer) GetDeleteCondition(ctx context.Context) string {
	if m.unscoped {
		return ""
	}
	conditionArray := garray.NewStrArray()
	if gstr.Contains(m.tables, " JOIN ") {
		// Base table.
		tableMatch, _ := gregex.MatchString(`(.+?) [A-Z]+ JOIN`, m.tables)
		conditionArray.Append(m.getConditionOfTableStringForSoftDeleting(ctx, tableMatch[1]))
		// Multiple joined tables, exclude the sub query sql which contains char '(' and ')'.
		tableMatches, _ := gregex.MatchAllString(`JOIN ([^()]+?) ON`, m.tables)
		for _, match := range tableMatches {
			conditionArray.Append(m.getConditionOfTableStringForSoftDeleting(ctx, match[1]))
		}
	}
	if conditionArray.Len() == 0 && gstr.Contains(m.tables, ",") {
		// Multiple base tables.
		for _, s := range gstr.SplitAndTrim(m.tables, ",") {
			conditionArray.Append(m.getConditionOfTableStringForSoftDeleting(ctx, s))
		}
	}
	conditionArray.FilterEmpty()
	if conditionArray.Len() > 0 {
		return conditionArray.Join(" AND ")
	}
	// Only one table.
	fieldName, fieldType := m.GetFieldInfo(ctx, "", m.tablesInit, SoftTimeFieldDelete)
	if fieldName != "" {
		return m.buildDeleteCondition(ctx, "", fieldName, fieldType)
	}
	return ""
}

// getConditionOfTableStringForSoftDeleting does something as its name describes.
// Examples for `s`:
// - `test`.`demo` as b
// - `test`.`demo` b
// - `demo`
// - demo
func (m *softTimeMaintainer) getConditionOfTableStringForSoftDeleting(ctx context.Context, s string) string {
	var (
		table  string
		schema string
		array1 = gstr.SplitAndTrim(s, " ")
		array2 = gstr.SplitAndTrim(array1[0], ".")
	)
	if len(array2) >= 2 {
		table = array2[1]
		schema = array2[0]
	} else {
		table = array2[0]
	}
	fieldName, fieldType := m.GetFieldInfo(ctx, schema, table, SoftTimeFieldDelete)
	if fieldName == "" {
		return ""
	}
	if len(array1) >= 3 {
		return m.buildDeleteCondition(ctx, array1[2], fieldName, fieldType)
	}
	if len(array1) >= 2 {
		return m.buildDeleteCondition(ctx, array1[1], fieldName, fieldType)
	}
	return m.buildDeleteCondition(ctx, table, fieldName, fieldType)
}

// GetDeleteData returns UPDATE statement data for soft delete.
func (m *softTimeMaintainer) GetDeleteData(
	ctx context.Context, prefix, fieldName string, fieldType LocalType,
) (holder string, value any) {
	core := m.db.GetCore()
	quotedName := core.QuoteWord(fieldName)

	if prefix != "" {
		quotedName = fmt.Sprintf(`%s.%s`, core.QuoteWord(prefix), quotedName)
	}

	holder = fmt.Sprintf(`%s=?`, quotedName)
	value = m.GetFieldValue(ctx, fieldType, false)
	return
}

// buildDeleteCondition builds WHERE condition for soft delete filtering.
func (m *softTimeMaintainer) buildDeleteCondition(
	ctx context.Context, prefix, fieldName string, fieldType LocalType,
) string {
	core := m.db.GetCore()
	quotedName := core.QuoteWord(fieldName)

	if prefix != "" {
		quotedName = fmt.Sprintf(`%s.%s`, core.QuoteWord(prefix), quotedName)
	}
	switch m.softTimeOption.SoftTimeType {
	case SoftTimeTypeAuto:
		switch fieldType {
		case LocalTypeDate, LocalTypeTime, LocalTypeDatetime:
			return fmt.Sprintf(`%s IS NULL`, quotedName)
		case LocalTypeInt, LocalTypeUint, LocalTypeInt64, LocalTypeUint64, LocalTypeBool:
			return fmt.Sprintf(`%s=0`, quotedName)
		default:
			intlog.Errorf(ctx, `invalid field type "%s" for soft delete condition: %s`, fieldType, fieldName)
			return ""
		}

	case SoftTimeTypeTime:
		return fmt.Sprintf(`%s IS NULL`, quotedName)

	default:
		return fmt.Sprintf(`%s=0`, quotedName)
	}
}

// GetFieldValue generates value for create/update/delete operations.
func (m *softTimeMaintainer) GetFieldValue(
	ctx context.Context, fieldType LocalType, isDeleted bool,
) any {
	// For deleted field, return "empty" value
	if isDeleted {
		return m.getEmptyValue(fieldType)
	}

	// For create/update/delete, return current time value
	switch m.softTimeOption.SoftTimeType {
	case SoftTimeTypeAuto:
		return m.getAutoValue(ctx, fieldType)
	default:
		switch fieldType {
		case LocalTypeBool:
			return 1
		default:
			return m.getTimestampValue()
		}
	}
}

// getTimestampValue returns timestamp value for soft time.
func (m *softTimeMaintainer) getTimestampValue() any {
	switch m.softTimeOption.SoftTimeType {
	case SoftTimeTypeTime:
		return gtime.Now()
	case SoftTimeTypeTimestamp:
		return gtime.Timestamp()
	case SoftTimeTypeTimestampMilli:
		return gtime.TimestampMilli()
	case SoftTimeTypeTimestampMicro:
		return gtime.TimestampMicro()
	case SoftTimeTypeTimestampNano:
		return gtime.TimestampNano()
	default:
		panic(gerror.NewCodef(
			gcode.CodeInternalPanic,
			`unrecognized SoftTimeType "%d"`, m.softTimeOption.SoftTimeType,
		))
	}
}

// getEmptyValue returns "empty" value for deleted field.
func (m *softTimeMaintainer) getEmptyValue(fieldType LocalType) any {
	switch fieldType {
	case LocalTypeDate, LocalTypeTime, LocalTypeDatetime:
		return nil
	default:
		return 0
	}
}

// getAutoValue returns auto-detected value based on field type.
func (m *softTimeMaintainer) getAutoValue(ctx context.Context, fieldType LocalType) any {
	switch fieldType {
	case LocalTypeDate, LocalTypeTime, LocalTypeDatetime:
		return gtime.Now()
	case LocalTypeInt, LocalTypeUint, LocalTypeInt64, LocalTypeUint64:
		return gtime.Timestamp()
	case LocalTypeBool:
		return 1
	default:
		intlog.Errorf(ctx, `invalid field type "%s" for soft time auto value`, fieldType)
		return nil
	}
}

// getTimeValue returns time value for datetime field type.
// For non-datetime types, returns 1 as a fallback marker.
func (m *softTimeMaintainer) getTimeValue(fieldType LocalType) any {
	switch fieldType {
	case LocalTypeDate, LocalTypeTime, LocalTypeDatetime:
		return gtime.Now()
	default:
		// For Bool or other types, use 1 as delete marker
		return 1
	}
}
