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
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

// ScanList converts `r` to struct slice which contains other complex struct attributes.
// Note that the parameter `structSlicePointer` should be type of *[]struct/*[]*struct.
//
// Usage example 1: Normal attribute struct relation:
//
//	type EntityUser struct {
//		   Uid  int
//		   Name string
//	}
//
//	type EntityUserDetail struct {
//		   Uid     int
//		   Address string
//	}
//
//	type EntityUserScores struct {
//		   Id     int
//		   Uid    int
//		   Score  int
//		   Course string
//	}
//
//	type Entity struct {
//	    User       *EntityUser
//		   UserDetail *EntityUserDetail
//		   UserScores []*EntityUserScores
//	}
//
// var users []*Entity
// ScanList(&users, "User")
// ScanList(&users, "User", "uid")
// ScanList(&users, "UserDetail", "User", "uid:Uid")
// ScanList(&users, "UserScores", "User", "uid:Uid")
// ScanList(&users, "UserScores", "User", "uid")
//
// Usage example 2: Embedded attribute struct relation:
//
//	type EntityUser struct {
//		   Uid  int
//		   Name string
//	}
//
//	type EntityUserDetail struct {
//		   Uid     int
//		   Address string
//	}
//
//	type EntityUserScores struct {
//		   Id    int
//		   Uid   int
//		   Score int
//	}
//
//	type Entity struct {
//		   EntityUser
//		   UserDetail EntityUserDetail
//		   UserScores []EntityUserScores
//	}
//
// var users []*Entity
// ScanList(&users)
// ScanList(&users, "UserDetail", "uid")
// ScanList(&users, "UserScores", "uid")
//
// The parameters "User/UserDetail/UserScores" in the example codes specify the target attribute struct
// that current result will be bound to.
//
// The "uid" in the example codes is the table field name of the result, and the "Uid" is the relational
// struct attribute name - not the attribute name of the bound to target. In the example codes, it's attribute
// name "Uid" of "User" of entity "Entity". It automatically calculates the HasOne/HasMany relationship with
// given `relation` parameter.
//
// See the example or unit testing cases for clear understanding for this function.
func (r Result) ScanList(structSlicePointer any, bindToAttrName string, relationAttrNameAndFields ...string) (err error) {
	out, err := checkGetSliceElementInfoForScanList(structSlicePointer, bindToAttrName)
	if err != nil {
		return err
	}

	var (
		relationAttrName string
		relationFields   string
	)
	switch len(relationAttrNameAndFields) {
	case 2:
		relationAttrName = relationAttrNameAndFields[0]
		relationFields = relationAttrNameAndFields[1]
	case 1:
		relationFields = relationAttrNameAndFields[0]
	}
	return doScanList(doScanListInput{
		Model:              nil,
		Result:             r,
		StructSlicePointer: structSlicePointer,
		StructSliceValue:   out.SliceReflectValue,
		BindToAttrName:     bindToAttrName,
		RelationAttrName:   relationAttrName,
		RelationFields:     relationFields,
		BatchEnabled:       false,
		BatchOptions:       nil,
		Depth:              1,
	})
}

type checkGetSliceElementInfoForScanListOutput struct {
	SliceReflectValue reflect.Value
	BindToAttrType    reflect.Type
}

func checkGetSliceElementInfoForScanList(structSlicePointer any, bindToAttrName string) (out *checkGetSliceElementInfoForScanListOutput, err error) {
	// Necessary checks for parameters.
	if structSlicePointer == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, `structSlicePointer cannot be nil`)
	}
	if bindToAttrName == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, `bindToAttrName should not be empty`)
	}
	var (
		reflectType  reflect.Type
		reflectValue = reflect.ValueOf(structSlicePointer)
		reflectKind  = reflectValue.Kind()
	)
	if reflectKind == reflect.Interface {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	if reflectKind != reflect.Pointer {
		return nil, gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"structSlicePointer should be type of *[]struct/*[]*struct, but got: %s",
			reflect.TypeOf(structSlicePointer).String(),
		)
	}
	out = &checkGetSliceElementInfoForScanListOutput{
		SliceReflectValue: reflectValue.Elem(),
	}
	// Find the element struct type of the slice.
	reflectType = reflectValue.Type().Elem().Elem()
	reflectKind = reflectType.Kind()
	for reflectKind == reflect.Pointer {
		reflectType = reflectType.Elem()
		reflectKind = reflectType.Kind()
	}
	if reflectKind != reflect.Struct {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"structSlicePointer should be type of *[]struct/*[]*struct, but got: %s",
			reflect.TypeOf(structSlicePointer).String(),
		)
		return
	}
	// Find the target field by given name.
	structField, ok := reflectType.FieldByName(bindToAttrName)
	if !ok {
		return nil, gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`field "%s" not found in element of "%s"`,
			bindToAttrName,
			reflect.TypeOf(structSlicePointer).String(),
		)
	}
	// Find the attribute struct type for ORM fields filtering.
	reflectType = structField.Type
	reflectKind = reflectType.Kind()
	for reflectKind == reflect.Pointer {
		reflectType = reflectType.Elem()
		reflectKind = reflectType.Kind()
	}
	if reflectKind == reflect.Slice || reflectKind == reflect.Array {
		reflectType = reflectType.Elem()
		// reflectKind = reflectType.Kind()
	}
	out.BindToAttrType = reflectType
	return
}

type doScanListInput struct {
	Model              *Model
	Result             Result
	StructSlicePointer any
	StructSliceValue   reflect.Value
	BindToAttrName     string
	RelationAttrName   string
	RelationFields     string
	BatchEnabled       bool
	BatchOptions       []WithBatchOption
	Depth              int
}

// doScanListRelation is the relation metadata for doScanList.
type doScanListRelation struct {
	DataMap         map[string]Value // Relation data map, which is Map[RelationValue]Record/Result.
	FromFieldName   string           // The field name of the result that is used for relation.
	BindToFieldName string           // The attribute name of the struct that is used for relation.
}

// doScanListBindAttr is the binding attribute information for doScanList.
type doScanListBindAttr struct {
	Field reflect.StructField // The struct field of the attribute.
	Kind  reflect.Kind        // The kind of the attribute.
	Type  reflect.Type        // The type of the attribute.
}

// doScanList converts `result` to struct slice which contains other complex struct attributes recursively.
// The parameter `model` is used for recursively scanning purpose, which means, it can scan the attribute struct/structs recursively,
// but it needs the Model for database accessing.
// Note that the parameter `structSlicePointer` should be type of *[]struct/*[]*struct.
func doScanList(in doScanListInput) (err error) {
	if in.Result.IsEmpty() {
		return nil
	}
	if in.BindToAttrName == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, `bindToAttrName should not be empty`)
	}

	var (
		length        = len(in.Result)
		arrayValue    reflect.Value
		arrayItemType reflect.Type
		reflectType   = reflect.TypeOf(in.StructSlicePointer)
	)
	if length == 0 {
		if in.StructSliceValue.Len() > 0 {
			if v := in.StructSliceValue.Index(0); v.Kind() != reflect.Pointer {
				return sql.ErrNoRows
			}
		}
		return nil
	}

	if in.StructSliceValue.Len() > 0 {
		arrayValue = in.StructSliceValue
	} else {
		arrayValue = reflect.MakeSlice(reflectType.Elem(), length, length)
	}
	arrayItemType = arrayValue.Index(0).Type()

	// 1. Parse relation metadata.
	relation, err := doScanListParseRelation(in)
	if err != nil {
		return err
	}

	// 2. Get target attribute info.
	attr, err := doScanListGetBindAttrInfo(arrayItemType, in.BindToAttrName)
	if err != nil {
		return err
	}

	// 3. Batch recursive scanning optimization.
	// 批量递归扫描优化：通过预先提取所有子数据的ID并执行批量查询，解决传统的N+1性能问题。
	structsMap, err := doScanListGetBatchRecursiveMap(in, attr, relation)
	if err != nil {
		return err
	}

	// 4. Final assignment loop.
	// 最终赋值循环：将查询到的数据分发到各个结构体属性中。
	if err = doScanListAssignmentLoop(in, arrayValue, attr, &relation, structsMap); err != nil {
		return err
	}

	reflect.ValueOf(in.StructSlicePointer).Elem().Set(arrayValue)
	return nil
}

// doScanListParseRelation parses the relation metadata from input.
func doScanListParseRelation(in doScanListInput) (relation doScanListRelation, err error) {
	if len(in.RelationFields) > 0 {
		array := gstr.SplitAndTrim(in.RelationFields, "=")
		if len(array) == 1 {
			array = gstr.SplitAndTrim(in.RelationFields, ":")
		}
		if len(array) == 1 {
			array = []string{in.RelationFields, in.RelationFields}
		}
		if len(array) == 2 {
			relation.FromFieldName = array[0]
			relation.BindToFieldName = array[1]
			if key, _ := gutil.MapPossibleItemByKey(in.Result[0].Map(), relation.FromFieldName); key == "" {
				return relation, gerror.NewCodef(
					gcode.CodeInvalidParameter,
					`cannot find possible related table field name "%s" from given relation fields "%s"`,
					relation.FromFieldName,
					in.RelationFields,
				)
			} else {
				relation.FromFieldName = key
			}
		} else {
			return relation, gerror.NewCode(
				gcode.CodeInvalidParameter,
				`parameter relationKV should be format of "ResultFieldName:BindToAttrName"`,
			)
		}
		if relation.FromFieldName != "" {
			relation.DataMap = in.Result.MapKeyValue(relation.FromFieldName)
		}
		if len(relation.DataMap) == 0 {
			return relation, gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`cannot find the relation data map, maybe invalid relation fields given "%v"`,
				in.RelationFields,
			)
		}
	}
	return relation, nil
}

// doScanListGetBindAttrInfo gets the binding attribute information from given array item type and name.
func doScanListGetBindAttrInfo(arrayItemType reflect.Type, bindToAttrName string) (attr doScanListBindAttr, err error) {
	var ok bool
	if arrayItemType.Kind() == reflect.Pointer {
		if attr.Field, ok = arrayItemType.Elem().FieldByName(bindToAttrName); !ok {
			return attr, gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid parameter bindToAttrName: cannot find attribute with name "%s" from slice element`,
				bindToAttrName,
			)
		}
	} else {
		if attr.Field, ok = arrayItemType.FieldByName(bindToAttrName); !ok {
			return attr, gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid parameter bindToAttrName: cannot find attribute with name "%s" from slice element`,
				bindToAttrName,
			)
		}
	}
	attr.Type = attr.Field.Type
	attr.Kind = attr.Type.Kind()
	return attr, nil
}

// doScanListGetBatchRecursiveMap executes the batch recursive scanning optimization (Solving N+1 problem).
// It returns a map that contains the relational structs, which can be used for fast assignment in the loop.
// 批量递归扫描的核心逻辑：
// 1. 根据当前层级（Depth）获取对应的批量配置（BatchSize, BatchThreshold）。
// 2. 将 Result 中的所有记录通过 Chunking 分段扫描到一个临时切片中。
// 3. 对该临时切片执行递归关联查询（doWithScanStructs）。
// 4. 构建一个以关联字段为 Key 的 Map，供后续 O(1) 复杂度的赋值循环使用。
func doScanListGetBatchRecursiveMap(
	in doScanListInput, attr doScanListBindAttr, relation doScanListRelation,
) (relationStructsMap map[string]reflect.Value, err error) {
	var (
		opt                  = in.getOptionForCurrentLayer()
		allChildStructsValue reflect.Value
	)
	if !opt.Enabled || len(in.Result) < opt.BatchThreshold {
		return nil, nil
	}

	if in.Model != nil && len(in.Result) > 0 {
		var (
			allChildStructsSlice reflect.Value
		)
		// Step 1: Prepare the container for bulk scanning.
		if attr.Kind == reflect.Array || attr.Kind == reflect.Slice {
			allChildStructsSlice = reflect.MakeSlice(attr.Field.Type, 0, len(in.Result))
		} else {
			allChildStructsSlice = reflect.MakeSlice(reflect.SliceOf(attr.Field.Type), 0, len(in.Result))
		}

		// Step 2: Scan all result records into the container and execute batch queries.
		// It uses chunking both for memory processing and database queries to:
		// 1. Avoid placeholder limit in SQL (MySQL limit: 65535)
		// 2. Control memory usage for large datasets
		// 3. Enable fine-grained control over query batch size
		batchSize := opt.BatchSize
		if batchSize <= 0 {
			batchSize = 1000
		}
		for i := 0; i < len(in.Result); i += batchSize {
			end := i + batchSize
			if end > len(in.Result) {
				end = len(in.Result)
			}
			var (
				resultChunk        = in.Result[i:end]
				allChildStructsPtr any
			)
			if attr.Kind == reflect.Array || attr.Kind == reflect.Slice {
				allChildStructsPtr = reflect.New(attr.Field.Type).Interface()
			} else {
				allChildStructsPtr = reflect.New(reflect.SliceOf(attr.Field.Type)).Interface()
			}
			if err = resultChunk.Structs(allChildStructsPtr); err != nil {
				return nil, err
			}

			// Step 3: Execute recursive relation queries for this batch.
			// Each batch will generate its own WHERE IN(...) clause with at most batchSize IDs.
			if err = in.Model.doWithScanStructs(allChildStructsPtr); err != nil {
				return nil, err
			}

			allChildStructsSlice = reflect.AppendSlice(allChildStructsSlice, reflect.ValueOf(allChildStructsPtr).Elem())
		}

		// Step 4: Build a map for fast lookup in the main assignment loop.
		// Note: After the batched queries above, allChildStructsSlice contains all results.
		allChildStructsValue = allChildStructsSlice
		if relation.FromFieldName != "" {
			relationStructsMap = make(map[string]reflect.Value)
			for i := 0; i < allChildStructsValue.Len(); i++ {
				kv := gconv.String(in.Result[i][relation.FromFieldName])
				if attr.Kind == reflect.Array || attr.Kind == reflect.Slice {
					if _, ok := relationStructsMap[kv]; !ok {
						relationStructsMap[kv] = reflect.MakeSlice(attr.Field.Type, 0, 0)
					}
					relationStructsMap[kv] = reflect.Append(relationStructsMap[kv], allChildStructsValue.Index(i))
				} else {
					if _, ok := relationStructsMap[kv]; !ok {
						relationStructsMap[kv] = allChildStructsValue.Index(i)
					}
				}
			}
		}
	}
	return relationStructsMap, nil
}

func (in doScanListInput) getOptionForCurrentLayer() WithBatchOption {
	// Default value.
	finalOpt := WithBatchOption{
		Enabled:   in.BatchEnabled,
		BatchSize: 1000,
	}
	// 1. Find the global default option (Layer == 0).
	for _, opt := range in.BatchOptions {
		if opt.Layer == 0 {
			finalOpt = opt
			// Do not return here, as we might have more specific layer option.
		}
	}
	// 2. Find the specific layer option (Layer == in.Depth).
	for _, opt := range in.BatchOptions {
		if opt.Layer == in.Depth {
			finalOpt = opt
			return finalOpt
		}
	}
	return finalOpt
}

// doScanListAssignmentLoop executes the final assignment loop for ScanList.
func doScanListAssignmentLoop(
	in doScanListInput,
	arrayValue reflect.Value,
	attr doScanListBindAttr,
	relation *doScanListRelation,
	structsMap map[string]reflect.Value,
) (err error) {
	var (
		arrayItemType                  = arrayValue.Index(0).Type()
		relationFromAttrValue          reflect.Value
		relationFromAttrField          reflect.Value
		relationBindToFieldNameChecked bool
	)

	// Phase 1 优化：使用缓存管理器获取字段索引缓存
	// 这里缓存了确定性的字段访问信息，避免循环内重复反射
	cache, err := globalFieldCacheManager.GetOrBuild(
		arrayItemType,
		in.BindToAttrName,
		in.RelationAttrName,
	)
	if err != nil {
		return err
	}

	for i := 0; i < arrayValue.Len(); i++ {
		arrayElemValue := arrayValue.Index(i)

		// 使用缓存的类型判断结果
		if cache.isPointerElem {
			arrayElemValue = arrayElemValue.Elem()
			if !arrayElemValue.IsValid() {
				arrayElemValue = reflect.New(arrayItemType.Elem()).Elem()
				arrayValue.Index(i).Set(arrayElemValue.Addr())
			}
		}

		// 使用缓存的字段索引直接访问（避免 FieldByName）
		bindToAttrValue := arrayElemValue.Field(cache.bindToAttrIndex)

		// 获取关系属性值
		if cache.relationAttrIndex >= 0 {
			relationFromAttrValue = arrayElemValue.Field(cache.relationAttrIndex)
			if relationFromAttrValue.Kind() == reflect.Pointer {
				relationFromAttrValue = relationFromAttrValue.Elem()
			}
		} else {
			relationFromAttrValue = arrayElemValue
		}

		if len(relation.DataMap) > 0 && !relationFromAttrValue.IsValid() {
			return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid relation fields specified: "%v"`, in.RelationFields)
		}

		// 嵌入字段查找保持原有的动态查找逻辑（不缓存）
		if in.RelationFields != "" && !relationBindToFieldNameChecked {
			relationFromAttrField = relationFromAttrValue.FieldByName(relation.BindToFieldName)
			if !relationFromAttrField.IsValid() {
				fieldMap, _ := gstructs.FieldMap(gstructs.FieldMapInput{
					Pointer:         relationFromAttrValue,
					RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
				})
				if key, _ := gutil.MapPossibleItemByKey(gconv.Map(fieldMap), relation.BindToFieldName); key == "" {
					return gerror.NewCodef(
						gcode.CodeInvalidParameter,
						`cannot find possible related attribute name "%s" from given relation fields "%s"`,
						relation.BindToFieldName,
						in.RelationFields,
					)
				} else {
					relation.BindToFieldName = key
				}
			}
			relationBindToFieldNameChecked = true
		}

		// 使用缓存的类型信息进行分发
		switch attr.Kind {
		case reflect.Array, reflect.Slice:
			if err = doScanListHandleAssignmentSlice(in, bindToAttrValue, relationFromAttrValue, *relation, structsMap); err != nil {
				return err
			}

		case reflect.Pointer:
			if err = doScanListHandleAssignmentPointer(in, bindToAttrValue, relationFromAttrValue, *relation, structsMap, attr, i); err != nil {
				return err
			}

		case reflect.Struct:
			if err = doScanListHandleAssignmentStruct(in, bindToAttrValue, relationFromAttrValue, *relation, structsMap, i); err != nil {
				return err
			}

		default:
			return gerror.NewCodef(gcode.CodeInvalidParameter, `unsupported attribute type: %s`, attr.Kind.String())
		}
	}
	return nil
}

// doScanListHandleAssignmentSlice handles the assignment for slice attribute.
func doScanListHandleAssignmentSlice(
	in doScanListInput,
	bindToAttrValue reflect.Value,
	relationFromAttrValue reflect.Value,
	relation doScanListRelation,
	structsMap map[string]reflect.Value,
) error {
	if len(structsMap) > 0 {
		relationFromAttrField := relationFromAttrValue.FieldByName(relation.BindToFieldName)
		if relationFromAttrField.IsValid() {
			key := gconv.String(relationFromAttrField.Interface())
			if structs, ok := structsMap[key]; ok {
				bindToAttrValue.Set(structs)
			}
		} else {
			return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid relation fields specified: "%v"`, in.RelationFields)
		}
	} else if len(relation.DataMap) > 0 {
		relationFromAttrField := relationFromAttrValue.FieldByName(relation.BindToFieldName)
		if relationFromAttrField.IsValid() {
			results := make(Result, 0)
			for _, v := range relation.DataMap[gconv.String(relationFromAttrField.Interface())].Slice() {
				results = append(results, v.(Record))
			}
			if err := results.Structs(bindToAttrValue.Addr()); err != nil {
				return err
			}
			if in.Model != nil {
				if err := in.Model.doWithScanStructs(bindToAttrValue.Addr()); err != nil {
					return err
				}
			}
		} else {
			return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid relation fields specified: "%v"`, in.RelationFields)
		}
	} else {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`relationKey should not be empty as field "%s" is slice`,
			in.BindToAttrName,
		)
	}
	return nil
}

// doScanListHandleAssignmentPointer handles the assignment for pointer attribute.
func doScanListHandleAssignmentPointer(
	in doScanListInput,
	bindToAttrValue reflect.Value,
	relationFromAttrValue reflect.Value,
	relation doScanListRelation,
	structsMap map[string]reflect.Value,
	attr doScanListBindAttr,
	index int,
) error {
	var element reflect.Value
	if bindToAttrValue.Kind() == reflect.Pointer && bindToAttrValue.IsNil() {
		element = reflect.New(attr.Type.Elem()).Elem()
	} else {
		element = bindToAttrValue.Elem()
	}
	if len(structsMap) > 0 {
		relationFromAttrField := relationFromAttrValue.FieldByName(relation.BindToFieldName)
		if relationFromAttrField.IsValid() {
			key := gconv.String(relationFromAttrField.Interface())
			if structs, ok := structsMap[key]; ok {
				bindToAttrValue.Set(structs)
			}
		} else {
			return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid relation fields specified: "%v"`, in.RelationFields)
		}
	} else if len(relation.DataMap) > 0 {
		relationFromAttrField := relationFromAttrValue.FieldByName(relation.BindToFieldName)
		if relationFromAttrField.IsValid() {
			v := relation.DataMap[gconv.String(relationFromAttrField.Interface())]
			if v == nil {
				return nil
			}
			if v.IsSlice() {
				if err := v.Slice()[0].(Record).Struct(element); err != nil {
					return err
				}
			} else {
				if err := v.Val().(Record).Struct(element); err != nil {
					return err
				}
			}
		} else {
			return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid relation fields specified: "%v"`, in.RelationFields)
		}
	} else {
		if index >= len(in.Result) {
			return nil
		}
		v := in.Result[index]
		if v == nil {
			return nil
		}
		if err := v.Struct(element); err != nil {
			return err
		}
	}
	if in.Model != nil && len(structsMap) == 0 {
		if err := in.Model.doWithScanStruct(element); err != nil {
			return err
		}
	}
	if len(structsMap) == 0 {
		bindToAttrValue.Set(element.Addr())
	}
	return nil
}

// doScanListHandleAssignmentStruct handles the assignment for struct attribute.
func doScanListHandleAssignmentStruct(
	in doScanListInput,
	bindToAttrValue reflect.Value,
	relationFromAttrValue reflect.Value,
	relation doScanListRelation,
	structsMap map[string]reflect.Value,
	index int,
) error {
	if len(structsMap) > 0 {
		relationFromAttrField := relationFromAttrValue.FieldByName(relation.BindToFieldName)
		if relationFromAttrField.IsValid() {
			key := gconv.String(relationFromAttrField.Interface())
			if structs, ok := structsMap[key]; ok {
				bindToAttrValue.Set(structs)
			}
		} else {
			return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid relation fields specified: "%v"`, in.RelationFields)
		}
	} else if len(relation.DataMap) > 0 {
		relationFromAttrField := relationFromAttrValue.FieldByName(relation.BindToFieldName)
		if relationFromAttrField.IsValid() {
			relationDataItem := relation.DataMap[gconv.String(relationFromAttrField.Interface())]
			if relationDataItem == nil {
				return nil
			}
			if relationDataItem.IsSlice() {
				if err := relationDataItem.Slice()[0].(Record).Struct(bindToAttrValue); err != nil {
					return err
				}
			} else {
				if err := relationDataItem.Val().(Record).Struct(bindToAttrValue); err != nil {
					return err
				}
			}
		} else {
			return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid relation fields specified: "%v"`, in.RelationFields)
		}
	} else {
		if index >= len(in.Result) {
			return nil
		}
		relationDataItem := in.Result[index]
		if relationDataItem == nil {
			return nil
		}
		if err := relationDataItem.Struct(bindToAttrValue); err != nil {
			return err
		}
	}
	if in.Model != nil && len(structsMap) == 0 {
		if err := in.Model.doWithScanStruct(bindToAttrValue); err != nil {
			return err
		}
	}
	return nil
}
