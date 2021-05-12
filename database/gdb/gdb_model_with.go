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
// Or:
//     db.With(UserDetail{}, UserDetail{}).Scan(xxx)
func (m *Model) With(objects ...interface{}) *Model {
	model := m.getModel()
	for _, object := range objects {
		if m.tables == "" {
			m.tables = m.db.QuotePrefixTableName(getTableNameFromOrmTag(object))
			return model
		}
		model.withArray = append(model.withArray, object)
	}
	return model
}

// WithAll enables model association operations on all objects that have "with" tag in the struct.
func (m *Model) WithAll() *Model {
	model := m.getModel()
	model.withAll = true
	return model
}

// doWithScanStruct handles model association operations feature for single struct.
func (m *Model) doWithScanStruct(pointer interface{}) error {
	var (
		err                 error
		allowedTypeStrArray = make([]string, 0)
	)
	fieldMap, err := structs.FieldMap(pointer, nil, false)
	if err != nil {
		return err
	}
	// It checks the with array and automatically calls the ScanList to complete association querying.
	if !m.withAll {
		for _, field := range fieldMap {
			for _, withItem := range m.withArray {
				withItemReflectValueType, err := structs.StructType(withItem)
				if err != nil {
					return err
				}
				var (
					fieldTypeStr                = gstr.TrimAll(field.Type().String(), "*[]")
					withItemReflectValueTypeStr = gstr.TrimAll(withItemReflectValueType.String(), "*[]")
				)
				// It does select operation if the field type is in the specified with type array.
				if gstr.Compare(fieldTypeStr, withItemReflectValueTypeStr) == 0 {
					allowedTypeStrArray = append(allowedTypeStrArray, fieldTypeStr)
				}
			}
		}
	}
	for _, field := range fieldMap {
		var (
			withTag      string
			ormTag       = field.Tag(OrmTagForStruct)
			fieldTypeStr = gstr.TrimAll(field.Type().String(), "*[]")
			match, _     = gregex.MatchString(
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
		if !m.withAll && !gstr.InArray(allowedTypeStrArray, fieldTypeStr) {
			continue
		}
		array := gstr.SplitAndTrim(withTag, "=")
		if len(array) == 1 {
			// It supports using only one column name
			// if both tables associates using the same column name.
			array = append(array, withTag)
		}
		var (
			model             *Model
			fieldKeys         []string
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
		bindToReflectValue := field.Value
		switch bindToReflectValue.Kind() {
		case reflect.Array, reflect.Slice:
			if bindToReflectValue.CanAddr() {
				bindToReflectValue = bindToReflectValue.Addr()
			}
		}

		// It automatically retrieves struct field names from current attribute struct/slice.
		if structType, err := structs.StructType(field.Value); err != nil {
			return err
		} else {
			fieldKeys = structType.FieldKeys()
		}

		// Recursively with feature checks.
		model = m.db.With(field.Value)
		if m.withAll {
			model = model.WithAll()
		} else {
			model = model.With(m.withArray...)
		}

		err = model.Fields(fieldKeys).Where(relatedFieldName, relatedFieldValue).Scan(bindToReflectValue)
		if err != nil {
			return err
		}

	}
	return nil
}

// doWithScanStructs handles model association operations feature for struct slice.
// Also see doWithScanStruct.
func (m *Model) doWithScanStructs(pointer interface{}) error {
	var (
		err                 error
		allowedTypeStrArray = make([]string, 0)
	)
	fieldMap, err := structs.FieldMap(pointer, nil, false)
	if err != nil {
		return err
	}
	// It checks the with array and automatically calls the ScanList to complete association querying.
	if !m.withAll {
		for _, field := range fieldMap {
			for _, withItem := range m.withArray {
				withItemReflectValueType, err := structs.StructType(withItem)
				if err != nil {
					return err
				}
				var (
					fieldTypeStr                = gstr.TrimAll(field.Type().String(), "*[]")
					withItemReflectValueTypeStr = gstr.TrimAll(withItemReflectValueType.String(), "*[]")
				)
				// It does select operation if the field type is in the specified with type array.
				if gstr.Compare(fieldTypeStr, withItemReflectValueTypeStr) == 0 {
					allowedTypeStrArray = append(allowedTypeStrArray, fieldTypeStr)
				}
			}
		}
	}

	for fieldName, field := range fieldMap {
		var (
			withTag      string
			ormTag       = field.Tag(OrmTagForStruct)
			fieldTypeStr = gstr.TrimAll(field.Type().String(), "*[]")
			match, _     = gregex.MatchString(
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
		if !m.withAll && !gstr.InArray(allowedTypeStrArray, fieldTypeStr) {
			continue
		}
		array := gstr.SplitAndTrim(withTag, "=")
		if len(array) == 1 {
			// It supports using only one column name
			// if both tables associates using the same column name.
			array = append(array, withTag)
		}
		var (
			model             *Model
			fieldKeys         []string
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

		// It automatically retrieves struct field names from current attribute struct/slice.
		if structType, err := structs.StructType(field.Value); err != nil {
			return err
		} else {
			fieldKeys = structType.FieldKeys()
		}

		// Recursively with feature checks.
		model = m.db.With(field.Value)
		if m.withAll {
			model = model.WithAll()
		} else {
			model = model.With(m.withArray...)
		}

		err = model.Fields(fieldKeys).Where(relatedFieldName, relatedFieldValue).ScanList(pointer, fieldName, withTag)
		if err != nil {
			return err
		}
	}
	return nil
}
