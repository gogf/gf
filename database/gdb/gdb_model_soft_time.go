// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

type softTimeMaintainer struct {
	*Model
}

type iSoftTimeMaintainer interface {
	GetSoftFieldNameAndTypeCreated(
		ctx context.Context, schema string, table string,
	) (fieldName string, fieldType LocalType)

	GetSoftFieldNameAndTypeUpdated(
		ctx context.Context, schema string, table string,
	) (fieldName string, fieldType LocalType)

	GetSoftFieldNameAndTypeDeleted(
		ctx context.Context, schema string, table string,
	) (fieldName string, fieldType LocalType)

	GetValueByFieldTypeForCreateOrUpdate(
		ctx context.Context, fieldType LocalType, isDeletedField bool,
	) (dataValue any)

	GetDataByFieldNameAndTypeForSoftDeleting(
		ctx context.Context, fieldPrefix, fieldName string, fieldType LocalType,
	) (dataHolder string, dataValue any)

	GetConditionForSoftDeleting(ctx context.Context) string
}

var (
	createdFieldNames = []string{"created_at", "create_at"} // Default field names of table for automatic-filled created datetime.
	updatedFieldNames = []string{"updated_at", "update_at"} // Default field names of table for automatic-filled updated datetime.
	deletedFieldNames = []string{"deleted_at", "delete_at"} // Default field names of table for automatic-filled deleted datetime.
)

// Unscoped disables the auto-update time feature for insert, update and delete options.
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

// GetSoftFieldNameAndTypeCreated checks and returns the field name for record creating time.
// If there's no field name for storing creating time, it returns an empty string.
// It checks the key with or without cases or chars '-'/'_'/'.'/' '.
func (m *softTimeMaintainer) GetSoftFieldNameAndTypeCreated(
	ctx context.Context, schema string, table string,
) (fieldName string, fieldType LocalType) {
	// It checks whether this feature disabled.
	if m.Model.db.GetConfig().TimeMaintainDisabled {
		return "", LocalTypeUndefined
	}
	tableName := ""
	if table != "" {
		tableName = table
	} else {
		tableName = m.tablesInit
	}
	config := m.db.GetConfig()
	if config.CreatedAt != "" {
		return m.getSoftFieldNameAndType(
			ctx, schema, tableName, []string{config.CreatedAt},
		)
	}
	return m.getSoftFieldNameAndType(
		ctx, schema, tableName, createdFieldNames,
	)
}

// GetSoftFieldNameAndTypeUpdated checks and returns the field name for record updating time.
// If there's no field name for storing updating time, it returns an empty string.
// It checks the key with or without cases or chars '-'/'_'/'.'/' '.
func (m *softTimeMaintainer) GetSoftFieldNameAndTypeUpdated(
	ctx context.Context, schema string, table string,
) (fieldName string, fieldType LocalType) {
	// It checks whether this feature disabled.
	if m.db.GetConfig().TimeMaintainDisabled {
		return "", LocalTypeUndefined
	}
	tableName := ""
	if table != "" {
		tableName = table
	} else {
		tableName = m.tablesInit
	}
	config := m.db.GetConfig()
	if config.UpdatedAt != "" {
		return m.getSoftFieldNameAndType(
			ctx, schema, tableName, []string{config.UpdatedAt},
		)
	}
	return m.getSoftFieldNameAndType(
		ctx, schema, tableName, updatedFieldNames,
	)
}

// GetSoftFieldNameAndTypeDeleted checks and returns the field name for record deleting time.
// If there's no field name for storing deleting time, it returns an empty string.
// It checks the key with or without cases or chars '-'/'_'/'.'/' '.
func (m *softTimeMaintainer) GetSoftFieldNameAndTypeDeleted(
	ctx context.Context, schema string, table string,
) (fieldName string, fieldType LocalType) {
	// It checks whether this feature disabled.
	if m.db.GetConfig().TimeMaintainDisabled {
		return "", LocalTypeUndefined
	}
	tableName := ""
	if table != "" {
		tableName = table
	} else {
		tableName = m.tablesInit
	}
	config := m.db.GetConfig()
	if config.DeletedAt != "" {
		return m.getSoftFieldNameAndType(
			ctx, schema, tableName, []string{config.DeletedAt},
		)
	}
	return m.getSoftFieldNameAndType(
		ctx, schema, tableName, deletedFieldNames,
	)
}

// getSoftFieldName retrieves and returns the field name of the table for possible key.
func (m *softTimeMaintainer) getSoftFieldNameAndType(
	ctx context.Context,
	schema string, table string, checkFiledNames []string,
) (fieldName string, fieldType LocalType) {
	// Ignore the error from TableFields.
	fieldsMap, _ := m.TableFields(table, schema)
	if len(fieldsMap) > 0 {
		for _, checkFiledName := range checkFiledNames {
			fieldName, _ = gutil.MapPossibleItemByKey(
				gconv.Map(fieldsMap), checkFiledName,
			)
			if fieldName != "" {
				fieldType, _ = m.db.CheckLocalTypeForField(
					ctx, fieldsMap[fieldName].Type, nil,
				)
				return
			}
		}
	}
	return
}

// GetConditionForSoftDeleting retrieves and returns the condition string for soft deleting.
// It supports multiple tables string like:
// "user u, user_detail ud"
// "user u LEFT JOIN user_detail ud ON(ud.uid=u.uid)"
// "user LEFT JOIN user_detail ON(user_detail.uid=user.uid)"
// "user u LEFT JOIN user_detail ud ON(ud.uid=u.uid) LEFT JOIN user_stats us ON(us.uid=u.uid)".
func (m *softTimeMaintainer) GetConditionForSoftDeleting(ctx context.Context) string {
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
	fieldName, fieldType := m.GetSoftFieldNameAndTypeDeleted(ctx, "", m.tablesInit)
	if fieldName != "" {
		return m.getConditionByFieldNameAndTypeForSoftDeleting(ctx, "", fieldName, fieldType)
	}
	return ""
}

func (m *softTimeMaintainer) getConditionByFieldNameAndTypeForSoftDeleting(
	ctx context.Context, fieldPrefix, fieldName string, fieldType LocalType,
) string {
	var (
		quotedFieldPrefix = m.db.GetCore().QuoteWord(fieldPrefix)
		quotedFieldName   = m.db.GetCore().QuoteWord(fieldName)
	)
	if quotedFieldPrefix != "" {
		quotedFieldName = fmt.Sprintf(`%s.%s`, quotedFieldPrefix, quotedFieldName)
	}
	switch fieldType {
	case LocalTypeDate, LocalTypeDatetime:
		return fmt.Sprintf(`%s IS NULL`, quotedFieldName)
	case LocalTypeInt, LocalTypeUint, LocalTypeInt64:
		return fmt.Sprintf(`%s=0`, quotedFieldName)
	case LocalTypeBool:
		return fmt.Sprintf(`%s=0`, quotedFieldName)
	default:
		intlog.Errorf(
			ctx,
			`invalid field type "%s" of field name "%s" for soft deleting condition`,
			fieldType,
		)
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
	fieldName, fieldType := m.GetSoftFieldNameAndTypeDeleted(ctx, schema, table)
	if fieldName == "" {
		return ""
	}
	if len(array1) >= 3 {
		return m.getConditionByFieldNameAndTypeForSoftDeleting(ctx, array1[2], fieldName, fieldType)
	}
	if len(array1) >= 2 {
		return m.getConditionByFieldNameAndTypeForSoftDeleting(ctx, array1[1], fieldName, fieldType)
	}
	return m.getConditionByFieldNameAndTypeForSoftDeleting(ctx, table, fieldName, fieldType)
}

// GetDataByFieldNameAndTypeForSoftDeleting creates and returns the placeholder and value for
// specified field name and type in soft-deleting scenario.
func (m *softTimeMaintainer) GetDataByFieldNameAndTypeForSoftDeleting(
	ctx context.Context, fieldPrefix, fieldName string, fieldType LocalType,
) (dataHolder string, dataValue any) {
	var (
		quotedFieldPrefix = m.db.GetCore().QuoteWord(fieldPrefix)
		quotedFieldName   = m.db.GetCore().QuoteWord(fieldName)
	)
	if quotedFieldPrefix != "" {
		quotedFieldName = fmt.Sprintf(`%s.%s`, quotedFieldPrefix, quotedFieldName)
	}
	dataHolder = fmt.Sprintf(`%s=?`, quotedFieldName)
	dataValue = m.GetValueByFieldTypeForCreateOrUpdate(ctx, fieldType, false)
	return
}

// GetValueByFieldTypeForCreateOrUpdate creates and returns the value for specified field type,
// usually for creating or updating operations.
func (m *softTimeMaintainer) GetValueByFieldTypeForCreateOrUpdate(
	ctx context.Context, fieldType LocalType, isDeletedField bool,
) any {
	switch fieldType {
	case LocalTypeDate, LocalTypeDatetime:
		if isDeletedField {
			return nil
		}
		return gtime.Now()
	case LocalTypeInt, LocalTypeUint, LocalTypeInt64:
		if isDeletedField {
			return 0
		}
		return gtime.Timestamp()
	case LocalTypeBool:
		if isDeletedField {
			return 0
		}
		return 1
	default:
		intlog.Errorf(
			ctx,
			`invalid field type "%s" of field name "%s" for soft deleting data`,
			fieldType,
		)
	}
	return nil
}
