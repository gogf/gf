// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// With creates and returns an ORM model based on metadata of given object.
// It also enables model association operations feature on given `object`.
// It can be called multiple times to add one or more objects to model and enable
// their mode association operations feature.
// For example, if given struct definition:
//
//	type User struct {
//		 gmeta.Meta `orm:"table:user"`
//		 Id         int           `json:"id"`
//		 Name       string        `json:"name"`
//		 UserDetail *UserDetail   `orm:"with:uid=id"`
//		 UserScores []*UserScores `orm:"with:uid=id"`
//	}
//
// We can enable model association operations on attribute `UserDetail` and `UserScores` by:
//
//	db.With(User{}.UserDetail).With(User{}.UserScores).Scan(xxx)
//
// Or:
//
//	db.With(UserDetail{}).With(UserScores{}).Scan(xxx)
//
// Or:
//
//	db.With(UserDetail{}, UserScores{}).Scan(xxx)
func (m *Model) With(objects ...any) *Model {
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

// WithOptions sets the batch association configuration options.
// It matches fields by chunkName and allows runtime override of chunk settings.
// Multiple options can be provided to configure different chunkName groups.
func (m *Model) WithOptions(options ...WithOption) *Model {
	model := m.getModel()
	if model.withOptions == nil {
		model.withOptions = make(map[ChunkName]*WithOption)
	}

	for _, opt := range options {
		// Skip empty chunkName
		if opt.ChunkName == "" {
			continue
		}
		// Store a copy of the option
		optCopy := opt
		model.withOptions[opt.ChunkName] = &optCopy
	}

	return model
}

// doWithScanStruct handles model association operations feature for single struct.
func (m *Model) doWithScanStruct(pointer any) error {
	if len(m.withArray) == 0 && !m.withAll {
		return nil
	}
	var (
		err                 error
		allowedTypeStrArray = make([]string, 0)
	)
	currentStructFieldMap, err := gstructs.FieldMap(gstructs.FieldMapInput{
		Pointer:          pointer,
		PriorityTagArray: nil,
		RecursiveOption:  gstructs.RecursiveOptionEmbeddedNoTag,
	})
	if err != nil {
		return err
	}
	// It checks the with array and automatically calls the ScanList to complete association querying.
	if !m.withAll {
		for _, field := range currentStructFieldMap {
			for _, withItem := range m.withArray {
				withItemReflectValueType, err := gstructs.StructType(withItem)
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
			fieldTypeStr = gstr.TrimAll(field.Type().String(), "*[]")
			withTag      = parseWithTag(field)
		)
		if withTag.With == "" {
			continue
		}
		// It just handlers "with" type attribute struct, so it ignores other struct types.
		if !m.withAll && !gstr.InArray(allowedTypeStrArray, fieldTypeStr) {
			continue
		}
		array := gstr.SplitAndTrim(withTag.With, "=")
		if len(array) == 1 {
			// It also supports using only one column name
			// if both tables associates using the same column name.
			array = append(array, withTag.With)
		}
		var (
			model              *Model
			fieldKeys          []string
			relatedSourceName  = array[0]
			relatedTargetName  = array[1]
			relatedTargetValue any
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
				relatedTargetName, withTag.With, reflect.TypeOf(pointer).Elem(), field.Name(),
			)
		}
		bindToReflectValue := field.Value
		if bindToReflectValue.Kind() != reflect.Pointer && bindToReflectValue.CanAddr() {
			bindToReflectValue = bindToReflectValue.Addr()
		}

		if structFields, err := gstructs.Fields(gstructs.FieldsInput{
			Pointer:         field.Value,
			RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
		}); err != nil {
			return err
		} else {
			fieldKeys = make([]string, len(structFields))
			for i, field := range structFields {
				fieldKeys[i] = field.Name()
			}
		}
		// Recursively with feature checks.
		model = m.db.With(field.Value).Hook(m.hookHandler)
		if m.withAll {
			model = model.WithAll()
		} else {
			model = model.With(m.withArray...)
		}
		if withTag.Where != "" {
			model = model.Where(withTag.Where)
		}
		if withTag.Order != "" {
			model = model.Order(withTag.Order)
		}
		if withTag.Unscoped == "true" {
			model = model.Unscoped()
		}
		// With cache feature.
		if m.cacheEnabled && m.cacheOption.Name == "" {
			model = model.Cache(m.cacheOption)
		}
		err = model.Fields(fieldKeys).
			Where(relatedSourceName, relatedTargetValue).
			Scan(bindToReflectValue)
		// It ignores sql.ErrNoRows in with feature.
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}
	return nil
}

// doWithScanStructs handles model association operations feature for struct slice.
// Also see doWithScanStruct.
func (m *Model) doWithScanStructs(pointer any) error {
	if len(m.withArray) == 0 && !m.withAll {
		return nil
	}
	return m.doBatchWithScan(pointer)
}

// withTagConfig holds the basic ORM tag configuration for a relation field.
type withTagConfig struct {
	With     string
	Where    string
	Order    string
	Unscoped string
}

// chunkTagConfig holds the chunk-related ORM tag configuration for a relation field.
type chunkTagConfig struct {
	ChunkName    string
	ChunkSize    int
	ChunkMinRows int
	Chunked      bool
}

// parseWithTag parses the basic ORM tag configuration (with, where, order, unscoped) from a struct field.
func parseWithTag(field gstructs.Field) withTagConfig {
	ormTag := field.Tag(OrmTagForStruct)
	data := parseOrmTagData(ormTag)
	return withTagConfig{
		With:     data[OrmTagForWith],
		Where:    data[OrmTagForWithWhere],
		Order:    data[OrmTagForWithOrder],
		Unscoped: data[OrmTagForWithUnscoped],
	}
}

// parseChunkTag parses the chunk-related ORM tag configuration from a struct field.
func parseChunkTag(field gstructs.Field) chunkTagConfig {
	ormTag := field.Tag(OrmTagForStruct)
	data := parseOrmTagData(ormTag)

	chunkName := data[OrmTagForChunkName]
	_, ifChunkSize := data[OrmTagForChunkSize]
	_, ifChunkMinRows := data[OrmTagForChunkMinRows]
	chunkSize := gconv.Int(data[OrmTagForChunkSize])
	chunkMinRows := gconv.Int(data[OrmTagForChunkMinRows])

	return chunkTagConfig{
		ChunkName:    chunkName,
		ChunkSize:    chunkSize,
		ChunkMinRows: chunkMinRows,
		Chunked:      ifChunkSize && ifChunkMinRows && chunkSize > 0 && chunkMinRows > 0,
	}
}

// parseOrmTagData parses the raw ORM tag string into a key-value map.
func parseOrmTagData(ormTag string) map[string]string {
	var (
		data  = make(map[string]string)
		array []string
		key   string
	)
	for _, v := range gstr.SplitAndTrim(ormTag, ",") {
		v = gstr.Trim(v)
		if v == "" {
			continue
		}
		array = gstr.Split(v, ":")
		if len(array) == 2 {
			key = array[0]
			data[key] = gstr.Trim(array[1])
		} else {
			switch key {
			case OrmTagForWithOrder:
				data[key] += "," + gstr.Trim(v)
			default:
				data[key] += " " + gstr.Trim(v)
			}
		}
	}
	return data
}
