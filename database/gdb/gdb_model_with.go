// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/structs"
	"github.com/gogf/gf/internal/utils"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
)

func (m *Model) With(structAttrPointer interface{}) *Model {
	model := m.getModel()
	if m.tables == "" {
		m.tables = m.db.QuotePrefixTableName(getTableNameFromOrmTag(structAttrPointer))
		return model
	}
	model.withArray = append(model.withArray, structAttrPointer)
	return model
}

func (m *Model) doWithScan(pointer interface{}) error {
	if len(m.withArray) == 0 {
		return nil
	}
	fieldMap, err := structs.FieldMap(pointer, nil)
	if err != nil {
		return err
	}
	for _, withItem := range m.withArray {
		withItemReflectValueType, err := structs.StructType(withItem)
		if err != nil {
			return err
		}
		withItemReflectValueTypeStr := gstr.TrimAll(withItemReflectValueType.String(), "*[]")
		for _, fieldValue := range fieldMap {
			var (
				fieldType    = fieldValue.Type()
				fieldTypeStr = gstr.TrimAll(fieldType.String(), "*[]")
			)
			if gstr.Compare(fieldTypeStr, withItemReflectValueTypeStr) == 0 {
				var (
					withTag  string
					ormTag   = fieldValue.Tag(OrmTagForStruct)
					match, _ = gregex.MatchString(
						fmt.Sprintf(`%s\s*:\s*([^,]+)`, OrmTagForWith),
						ormTag,
					)
				)
				if len(match) > 1 {
					withTag = match[1]
				}
				if withTag == "" {
					continue
				}
				array := gstr.SplitAndTrim(withTag, "=")
				if len(array) != 2 {
					return gerror.Newf(`invalid with tag "%s"`, withTag)
				}
				var (
					relatedFieldName  = array[0]
					relatedAttrName   = array[1]
					relatedFieldValue interface{}
				)
				// Find the value of related attribute from `pointer`.
				for attributeName, attributeValue := range fieldMap {
					if utils.EqualFoldWithoutChars(attributeName, relatedAttrName) {
						relatedFieldValue = attributeValue.Value.Interface()
						break
					}
				}
				if relatedFieldValue == nil {
					return gerror.Newf(
						`cannot find the related value for attribute name "%s" of with tag "%s"`,
						relatedAttrName, withTag,
					)
				}
				err = m.db.With(fieldValue.Value).
					Fields(withItemReflectValueType.FieldKeys()).
					Where(relatedFieldName, relatedFieldValue).
					Scan(fieldValue.Value)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
