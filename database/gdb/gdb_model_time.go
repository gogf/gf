// Copyright 2020 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package gdb

import (
	"fmt"
	"github.com/jin502437344/gf/container/garray"
	"github.com/jin502437344/gf/text/gregex"
	"github.com/jin502437344/gf/text/gstr"
	"github.com/jin502437344/gf/util/gconv"
	"github.com/jin502437344/gf/util/gutil"
)

const (
	gSOFT_FIELD_NAME_CREATE = "create_at"
	gSOFT_FIELD_NAME_UPDATE = "update_at"
	gSOFT_FIELD_NAME_DELETE = "delete_at"
)

// getSoftFieldNameCreate checks and returns the field name for record creating time.
// If there's no field name for storing creating time, it returns an empty string.
// It checks the key with or without cases or chars '-'/'_'/'.'/' '.
func (m *Model) getSoftFieldNameCreate(table ...string) string {
	tableName := ""
	if len(table) > 0 {
		tableName = table[0]
	} else {
		tableName = m.getPrimaryTableName()
	}
	return m.getSoftFieldName(tableName, gSOFT_FIELD_NAME_CREATE)
}

// getSoftFieldNameUpdate checks and returns the field name for record updating time.
// If there's no field name for storing updating time, it returns an empty string.
// It checks the key with or without cases or chars '-'/'_'/'.'/' '.
func (m *Model) getSoftFieldNameUpdate(table ...string) (field string) {
	tableName := ""
	if len(table) > 0 {
		tableName = table[0]
	} else {
		tableName = m.getPrimaryTableName()
	}
	return m.getSoftFieldName(tableName, gSOFT_FIELD_NAME_UPDATE)
}

// getSoftFieldNameDelete checks and returns the field name for record deleting time.
// If there's no field name for storing deleting time, it returns an empty string.
// It checks the key with or without cases or chars '-'/'_'/'.'/' '.
func (m *Model) getSoftFieldNameDelete(table ...string) (field string) {
	tableName := ""
	if len(table) > 0 {
		tableName = table[0]
	} else {
		tableName = m.getPrimaryTableName()
	}
	return m.getSoftFieldName(tableName, gSOFT_FIELD_NAME_DELETE)
}

// getSoftFieldName retrieves and returns the field name of the table for possible key.
func (m *Model) getSoftFieldName(table string, key string) (field string) {
	fieldsMap, _ := m.db.TableFields(table)
	if len(fieldsMap) > 0 {
		field, _ = gutil.MapPossibleItemByKey(
			gconv.Map(fieldsMap), key,
		)
	}
	return
}

// getConditionForSoftDeleting retrieves and returns the condition string for soft deleting.
// It supports multiple tables string like:
// "user u, user_detail ud"
// "user u LEFT JOIN user_detail ud ON(ud.uid=u.uid)"
// "user LEFT JOIN user_detail ON(user_detail.uid=user.uid)"
// "user u LEFT JOIN user_detail ud ON(ud.uid=u.uid) LEFT JOIN user_stats us ON(us.uid=u.uid)"
func (m *Model) getConditionForSoftDeleting() string {
	if m.unscoped {
		return ""
	}
	conditionArray := garray.NewStrArray()
	if gstr.Contains(m.tables, " JOIN ") {
		// Base table.
		match, _ := gregex.MatchString(`(.+?) [A-Z]+ JOIN`, m.tables)
		conditionArray.Append(m.getConditionOfTableStringForSoftDeleting(match[1]))
		// Multiple joined tables.
		matches, _ := gregex.MatchAllString(`JOIN (.+?) ON`, m.tables)
		for _, match := range matches {
			conditionArray.Append(m.getConditionOfTableStringForSoftDeleting(match[1]))
		}
	}
	if conditionArray.Len() == 0 && gstr.Contains(m.tables, ",") {
		// Multiple base tables.
		for _, s := range gstr.SplitAndTrim(m.tables, ",") {
			conditionArray.Append(m.getConditionOfTableStringForSoftDeleting(s))
		}
	}
	conditionArray.FilterEmpty()
	if conditionArray.Len() > 0 {
		return conditionArray.Join(" AND ")
	}
	// Only one table.
	if fieldName := m.getSoftFieldNameDelete(); fieldName != "" {
		return fmt.Sprintf(`%s IS NULL`, m.db.QuoteWord(fieldName))
	}
	return ""
}

// getConditionOfTableStringForSoftDeleting does something as its name describes.
func (m *Model) getConditionOfTableStringForSoftDeleting(s string) string {
	var (
		field  = ""
		table  = ""
		array1 = gstr.SplitAndTrim(s, " ")
		array2 = gstr.SplitAndTrim(array1[0], ".")
	)
	if len(array2) >= 2 {
		table = array2[1]
	} else {
		table = array2[0]
	}
	field = m.getSoftFieldNameDelete(table)
	if field == "" {
		return ""
	}
	if len(array1) >= 3 {
		return fmt.Sprintf(`%s.%s IS NULL`, m.db.QuoteWord(array1[2]), m.db.QuoteWord(field))
	}
	if len(array1) >= 2 {
		return fmt.Sprintf(`%s.%s IS NULL`, m.db.QuoteWord(array1[1]), m.db.QuoteWord(field))
	}
	return fmt.Sprintf(`%s.%s IS NULL`, m.db.QuoteWord(table), m.db.QuoteWord(field))
}

// getPrimaryTableName parses and returns the primary table name.
func (m *Model) getPrimaryTableName() string {
	array1 := gstr.SplitAndTrim(m.tables, ",")
	array2 := gstr.SplitAndTrim(array1[0], " ")
	array3 := gstr.SplitAndTrim(array2[0], ".")
	if len(array3) >= 2 {
		return array3[1]
	}
	return array3[0]
}
