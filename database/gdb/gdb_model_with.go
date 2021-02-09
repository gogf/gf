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
	"reflect"
)

// With creates and returns an ORM model based on meta data of given object.
// It also enables model association operations feature on given `object`.
// It can be called multiple times to add one or more objects to model and enable
// their mode association operations feature.
// For example, if given struct definition:
// type User struct {
//	 gmeta.Meta `orm:"table:user"`
// 	 Id         int           `json:"id"`
//	 Name       string        `json:"name"`
//	 UserDetail *UserDetail   `orm:"with:uid=id"`
//	 UserScores []*UserScores `orm:"with:uid=id"`
// }
// We can enable model association operations on attribute `UserDetail` and `UserScores` by:
//     db.With(User{}.UserDetail).With(User{}.UserDetail).Scan(xxx)
// Or:
//     db.With(UserDetail{}).With(UserDetail{}).Scan(xxx)
func (m *Model) With(object interface{}) *Model {
	model := m.getModel()
	if m.tables == "" {
		m.tables = m.db.QuotePrefixTableName(getTableNameFromOrmTag(object))
		return model
	}
	model.withArray = append(model.withArray, object)
	return model
}

// WithAll enables model association operations on all objects that have "with" tag in the struct.
func (m *Model) WithAll() *Model {
	model := m.getModel()
	model.withAll = true
	return model
}

// getWithTagObjectArrayFrom retrieves and returns object array that have "with" tag in the struct.
func (m *Model) getWithTagObjectArrayFrom(pointer interface{}) ([]interface{}, error) {
	fieldMap, err := structs.FieldMap(pointer, nil)
	if err != nil {
		return nil, err
	}
	withTagObjectArray := make([]interface{}, 0)
	for _, fieldValue := range fieldMap {
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
		withTagObjectArray = append(withTagObjectArray, fieldValue.Value.Interface())
	}
	return withTagObjectArray, nil
}

// doWithScanStruct handles model association operations feature for single struct.
func (m *Model) doWithScanStruct(pointer interface{}) error {
	var (
		err       error
		withArray = m.withArray
	)
	if m.withAll {
		withArray, err = m.getWithTagObjectArrayFrom(pointer)
		if err != nil {
			return err
		}
	}
	if len(withArray) == 0 {
		return nil
	}
	fieldMap, err := structs.FieldMap(pointer, nil)
	if err != nil {
		return err
	}
	for withIndex, withItem := range withArray {
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
				bindToReflectValue := fieldValue.Value
				switch bindToReflectValue.Kind() {
				case reflect.Array, reflect.Slice:
					if bindToReflectValue.CanAddr() {
						bindToReflectValue = bindToReflectValue.Addr()
					}
				}
				model := m.db.With(fieldValue.Value)
				for i, v := range withArray {
					if i == withIndex {
						continue
					}
					model = model.With(v)
				}
				err = model.Fields(withItemReflectValueType.FieldKeys()).
					Where(relatedFieldName, relatedFieldValue).
					Scan(bindToReflectValue)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// doWithScanStructs handles model association operations feature for struct slice.
func (m *Model) doWithScanStructs(pointer interface{}) error {
	var (
		err       error
		withArray = m.withArray
	)
	if m.withAll {
		withArray, err = m.getWithTagObjectArrayFrom(pointer)
		if err != nil {
			return err
		}
	}
	if len(withArray) == 0 {
		return nil
	}
	fieldMap, err := structs.FieldMap(pointer, nil)
	if err != nil {
		return err
	}
	for withIndex, withItem := range withArray {
		withItemReflectValueType, err := structs.StructType(withItem)
		if err != nil {
			return err
		}
		withItemReflectValueTypeStr := gstr.TrimAll(withItemReflectValueType.String(), "*[]")
		for fieldName, fieldValue := range fieldMap {
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
				// Find the value slice of related attribute from `pointer`.
				for attributeName, _ := range fieldMap {
					if utils.EqualFoldWithoutChars(attributeName, relatedAttrName) {
						relatedFieldValue = ListItemValuesUnique(pointer, attributeName)
						break
					}
				}
				if relatedFieldValue == nil {
					return gerror.Newf(
						`cannot find the related value for attribute name "%s" of with tag "%s"`,
						relatedAttrName, withTag,
					)
				}
				model := m.db.With(fieldValue.Value)
				for i, v := range withArray {
					if i == withIndex {
						continue
					}
					model = model.With(v)
				}
				err = model.Fields(withItemReflectValueType.FieldKeys()).
					Where(relatedFieldName, relatedFieldValue).
					ScanList(pointer, fieldName, withTag)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
