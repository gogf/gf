// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"reflect"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
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

// WithBatch enables or disables the batch recursive scanning feature for association operations.
// The batch recursive scanning feature is used to solve the N+1 problem by batching multiple
// association queries into one or fewer queries.
// It is disabled by default.
// 开启或关闭关联查询的批量递归扫描功能（解决N+1问题）。
// 默认关闭，开启后可大幅提升存在大量关联数据时的查询性能。
func (m *Model) WithBatch(enabled ...bool) *Model {
	model := m.getModel()
	model.withBatchEnabled = true
	if len(enabled) > 0 {
		model.withBatchEnabled = enabled[0]
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
			fieldTypeStr    = gstr.TrimAll(field.Type().String(), "*[]")
			parsedTagOutput = parseWithTagInField(field.Field)
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
				relatedTargetName, parsedTagOutput.With, reflect.TypeOf(pointer).Elem(), field.Name(),
			)
		}
		bindToReflectValue := field.Value
		if bindToReflectValue.Kind() != reflect.Pointer && bindToReflectValue.CanAddr() {
			bindToReflectValue = bindToReflectValue.Addr()
		}

		// Recursively with feature checks.
		model = m.db.With(field.Value).Hook(m.hookHandler)
		model.withBatchEnabled = m.withBatchEnabled
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
		if parsedTagOutput.Unscoped == "true" {
			model = model.Unscoped()
		}
		// Apply cache option if enabled (for query result caching, not field metadata).
		if m.cacheEnabled && m.cacheOption.Name == "" {
			model = model.Cache(m.cacheOption)
		}
		// Fields will be automatically determined from the struct type
		err = model.Fields(field.Value).
			Where(relatedSourceName, relatedTargetValue).
			Scan(bindToReflectValue)
		// It ignores sql.ErrNoRows in with feature.
		if err != nil && err != sql.ErrNoRows {
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
	if v, ok := pointer.(reflect.Value); ok {
		pointer = v.Interface()
	}
	var (
		err                 error
		allowedTypeStrArray = make([]string, 0)
		reflectValue        = reflect.ValueOf(pointer)
		reflectKind         = reflectValue.Kind()
	)
	if reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	if reflectKind != reflect.Slice && reflectKind != reflect.Array {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`the parameter "pointer" for doWithScanStructs should be type of slice, invalid type: %v`,
			reflect.TypeOf(pointer),
		)
	}

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
				// It does select operation if the field type is in the specified with type array.
				if gstr.Compare(fieldTypeStr, withItemReflectValueTypeStr) == 0 {
					allowedTypeStrArray = append(allowedTypeStrArray, fieldTypeStr)
				}
			}
		}
	}

	for fieldName, field := range currentStructFieldMap {
		var (
			fieldTypeStr = gstr.TrimAll(field.Type().String(), "*[]")
		)
		// Parse withTag directly from field instead of using cache to avoid cache pollution
		// when multiple tests define struct with same name but different tags
		parsedTagOutput := parseWithTagInField(field.Field)
		if parsedTagOutput.With == "" {
			continue
		}
		if !m.withAll && !gstr.InArray(allowedTypeStrArray, fieldTypeStr) {
			continue
		}
		array := gstr.SplitAndTrim(parsedTagOutput.With, "=")
		if len(array) == 1 {
			array = append(array, parsedTagOutput.With)
		}
		var (
			model              *Model
			relatedSourceName  = array[0]
			relatedTargetName  = array[1]
			relatedTargetValue any
		)
		// Find the value slice of related attribute from `pointer`.
		for attributeName := range currentStructFieldMap {
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
		// If related value is empty, it does nothing but just returns.
		if gutil.IsEmpty(relatedTargetValue) {
			continue
		}
		// Recursively with feature checks.
		model = m.db.With(field.Value).Hook(m.hookHandler)
		model.withBatchEnabled = m.withBatchEnabled
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
		if parsedTagOutput.Unscoped == "true" {
			model = model.Unscoped()
		}
		// Apply cache option if enabled (for query result caching, not field metadata).
		if m.cacheEnabled && m.cacheOption.Name == "" {
			model = model.Cache(m.cacheOption)
		}

		var (
			batchSize      int
			batchThreshold int
			results        Result
		)

		if m.withBatchEnabled {
			batchSize = parsedTagOutput.BatchSize
			batchThreshold = parsedTagOutput.BatchThreshold
		}

		if m.withBatchEnabled && batchSize > 0 && len(gconv.SliceAny(relatedTargetValue)) >= batchThreshold {
			var ids = gconv.SliceAny(relatedTargetValue)
			for i := 0; i < len(ids); i += batchSize {
				end := i + batchSize
				if end > len(ids) {
					end = len(ids)
				}
				// 使用 Clone() 避免条件累加
				// Fields will be automatically determined from the struct type
				result, err := model.Clone().Fields(field.Value).
					Where(relatedSourceName, ids[i:end]).
					All()
				if err != nil {
					return err
				}
				results = append(results, result...)
			}
		} else {
			// Fields will be automatically determined from the struct type
			results, err = model.Clone().Fields(field.Value).
				Where(relatedSourceName, relatedTargetValue).
				All()
			if err != nil && err != sql.ErrNoRows {
				return err
			}
		}

		if results.IsEmpty() {
			continue
		}

		err = doScanList(doScanListInput{
			Model:              model,
			Result:             results,
			StructSlicePointer: pointer,
			StructSliceValue:   reflect.ValueOf(pointer).Elem(),
			BindToAttrName:     fieldName,
			RelationAttrName:   "",
			RelationFields:     parsedTagOutput.With,
			BatchEnabled:       m.withBatchEnabled,
			BatchSize:          batchSize,
			BatchThreshold:     batchThreshold,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

type withTagOutput struct {
	With           string
	Where          string
	Order          string
	Unscoped       string
	BatchSize      int
	BatchThreshold int
}

func parseWithTagInField(field reflect.StructField) (output withTagOutput) {
	var (
		ormTag = field.Tag.Get(OrmTagForStruct)
		data   = make(map[string]string)
	)
	// Parse tags, support key:value and nested batch:threshold=1000,batchSize=100
	for _, v := range gstr.SplitAndTrim(ormTag, ",") {
		v = gstr.Trim(v)
		if v == "" {
			continue
		}

		// 处理 batch: 开头的特殊配置
		if gstr.HasPrefix(v, "batch:") {
			// 提取 batch: 后面的内容
			batchConfig := gstr.TrimLeft(v, "batch:")
			// 解析 batch 内部的配置项（如 threshold=1000,batchSize=100）
			for _, batchItem := range gstr.SplitAndTrim(batchConfig, ",") {
				parts := gstr.Split(batchItem, "=")
				if len(parts) == 2 {
					data[gstr.Trim(parts[0])] = gstr.Trim(parts[1])
				}
			}
			continue
		}

		// Process normal key:value or key=value
		var (
			key   string
			value string
			parts = gstr.Split(v, ":")
		)
		if len(parts) == 2 {
			key = gstr.Trim(parts[0])
			value = gstr.Trim(parts[1])
		} else {
			parts = gstr.Split(v, "=")
			if len(parts) == 2 {
				key = gstr.Trim(parts[0])
				value = gstr.Trim(parts[1])
			}
		}

		if key != "" {
			data[key] = value
		}
	}
	output.With = data[OrmTagForWith]
	output.Where = data[OrmTagForWithWhere]
	output.Order = data[OrmTagForWithOrder]
	output.Unscoped = data[OrmTagForWithUnscoped]
	output.BatchSize = gconv.Int(data["batchSize"])
	output.BatchThreshold = gconv.Int(data["threshold"])
	return
}
