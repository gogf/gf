// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/structs"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

// With creates and returns an ORM model based on metadata of given object.
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
			m.tablesInit = m.db.GetCore().QuotePrefixTableName(
				getTableNameFromOrmTag(object),
			)
			m.tables = m.tablesInit
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
	currentStructFieldMap, err := structs.FieldMap(structs.FieldMapInput{
		Pointer:          pointer,
		PriorityTagArray: nil,
		RecursiveOption:  structs.RecursiveOptionEmbeddedNoTag,
	})
	if err != nil {
		return err
	}
	// It checks the with array and automatically calls the ScanList to complete association querying.
	if !m.withAll {
		for _, field := range currentStructFieldMap {
			for _, withItem := range m.withArray {
				withItemReflectValueType, err := structs.StructType(withItem)
				if err != nil {
					return err
				}
				var (
					fieldTypeStr                = gstr.TrimAll(field.Type().String(), "*[]")
					withItemReflectValueTypeStr = gstr.TrimAll(withItemReflectValueType.String(), "*[]")
				)
				// It does select operation if the field type is in the specified "with" type array.
				if gstr.Compare(fieldTypeStr, withItemReflectValueTypeStr) == 0 {
					allowedTypeStrArray = append(allowedTypeStrArray, fieldTypeStr)
				}
			}
		}
	}
	for _, field := range currentStructFieldMap {
		var (
			fieldTypeStr    = gstr.TrimAll(field.Type().String(), "*[]")
			parsedTagOutput = m.parseWithTagInFieldStruct(field)
		)
		if parsedTagOutput.With == "" {
			continue
		}
		// It just handlers "with" type attribute struct, so it ignores other struct types.
		if !m.withAll && !gstr.InArray(allowedTypeStrArray, fieldTypeStr) {
			continue
		}
		array := gstr.SplitAndTrim(parsedTagOutput.With, "=")
		if len(array) == 1 {
			// It also supports using only one column name
			// if both tables associates using the same column name.
			array = append(array, parsedTagOutput.With)
		}
		var (
			model              *Model
			fieldKeys          []string
			relatedSourceName  = array[0]
			relatedTargetName  = array[1]
			relatedTargetValue interface{}
		)
		// Find the value of related attribute from `pointer`.
		for attributeName, attributeValue := range currentStructFieldMap {
			if utils.EqualFoldWithoutChars(attributeName, relatedTargetName) {
				relatedTargetValue = attributeValue.Value.Interface()
				break
			}
		}
		if relatedTargetValue == nil {
			return gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`cannot find the target related value of name "%s" in with tag "%s" for attribute "%s.%s"`,
				relatedTargetName, parsedTagOutput.With, reflect.TypeOf(pointer).Elem(), field.Name(),
			)
		}
		bindToReflectValue := field.Value
		if bindToReflectValue.Kind() != reflect.Ptr && bindToReflectValue.CanAddr() {
			bindToReflectValue = bindToReflectValue.Addr()
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
		if parsedTagOutput.Where != "" {
			model = model.Where(parsedTagOutput.Where)
		}
		if parsedTagOutput.Order != "" {
			model = model.Order(parsedTagOutput.Order)
		}
		err = model.Fields(fieldKeys).Where(relatedSourceName, relatedTargetValue).Scan(bindToReflectValue)
		// It ignores sql.ErrNoRows in with feature.
		if err != nil && err != sql.ErrNoRows {
			return err
		}
	}
	return nil
}

// doWithScanStructs handles model association operations feature for struct slice.
// Also see doWithScanStruct.
func (m *Model) doWithScanStructs(pointer interface{}) error {
	if v, ok := pointer.(reflect.Value); ok {
		pointer = v.Interface()
	}

	var (
		err                 error
		allowedTypeStrArray = make([]string, 0)
	)
	currentStructFieldMap, err := structs.FieldMap(structs.FieldMapInput{
		Pointer:          pointer,
		PriorityTagArray: nil,
		RecursiveOption:  structs.RecursiveOptionEmbeddedNoTag,
	})
	if err != nil {
		return err
	}
	// It checks the with array and automatically calls the ScanList to complete association querying.
	if !m.withAll {
		for _, field := range currentStructFieldMap {
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

	for fieldName, field := range currentStructFieldMap {
		var (
			fieldTypeStr    = gstr.TrimAll(field.Type().String(), "*[]")
			parsedTagOutput = m.parseWithTagInFieldStruct(field)
		)
		if parsedTagOutput.With == "" {
			continue
		}
		if !m.withAll && !gstr.InArray(allowedTypeStrArray, fieldTypeStr) {
			continue
		}
		array := gstr.SplitAndTrim(parsedTagOutput.With, "=")
		if len(array) == 1 {
			// It supports using only one column name
			// if both tables associates using the same column name.
			array = append(array, parsedTagOutput.With)
		}
		var (
			model              *Model
			fieldKeys          []string
			relatedSourceName  = array[0]
			relatedTargetName  = array[1]
			relatedTargetValue interface{}
		)
		// Find the value slice of related attribute from `pointer`.
		for attributeName, _ := range currentStructFieldMap {
			if utils.EqualFoldWithoutChars(attributeName, relatedTargetName) {
				relatedTargetValue = ListItemValuesUnique(pointer, attributeName)
				break
			}
		}
		if relatedTargetValue == nil {
			return gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`cannot find the related value for attribute name "%s" of with tag "%s"`,
				relatedTargetName, parsedTagOutput.With,
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
		if parsedTagOutput.Where != "" {
			model = model.Where(parsedTagOutput.Where)
		}
		if parsedTagOutput.Order != "" {
			model = model.Order(parsedTagOutput.Order)
		}

		err = model.Fields(fieldKeys).
			Where(relatedSourceName, relatedTargetValue).
			ScanList(pointer, fieldName, parsedTagOutput.With)
		// It ignores sql.ErrNoRows in with feature.
		if err != nil && err != sql.ErrNoRows {
			return err
		}
	}
	return nil
}

type parseWithTagInFieldStructOutput struct {
	With  string
	Where string
	Order string
}

func (m *Model) parseWithTagInFieldStruct(field *structs.Field) (output parseWithTagInFieldStructOutput) {
	var (
		match  []string
		ormTag = field.Tag(OrmTagForStruct)
	)
	// with tag.
	match, _ = gregex.MatchString(
		fmt.Sprintf(`%s\s*:\s*([^,]+),{0,1}`, OrmTagForWith),
		ormTag,
	)
	if len(match) > 1 {
		output.With = match[1]
	}
	if len(match) > 2 {
		output.Where = gstr.Trim(match[2])
	}
	// where string.
	match, _ = gregex.MatchString(
		fmt.Sprintf(`%s\s*:.+,\s*%s:\s*([^,]+),{0,1}`, OrmTagForWith, OrmTagForWithWhere),
		ormTag,
	)
	if len(match) > 1 {
		output.Where = gstr.Trim(match[1])
	}
	// order string.
	match, _ = gregex.MatchString(
		fmt.Sprintf(`%s\s*:.+,\s*%s:\s*([^,]+),{0,1}`, OrmTagForWith, OrmTagForWithOrder),
		ormTag,
	)
	if len(match) > 1 {
		output.Order = gstr.Trim(match[1])
	}
	return
}
