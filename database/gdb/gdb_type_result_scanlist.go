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
func (r Result) ScanList(structSlicePointer interface{}, bindToAttrName string, relationAttrNameAndFields ...string) (err error) {
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
	})
}

type checkGetSliceElementInfoForScanListOutput struct {
	SliceReflectValue reflect.Value
	BindToAttrType    reflect.Type
}

func checkGetSliceElementInfoForScanList(structSlicePointer interface{}, bindToAttrName string) (out *checkGetSliceElementInfoForScanListOutput, err error) {
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
	if reflectKind != reflect.Ptr {
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
	for reflectKind == reflect.Ptr {
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
	for reflectKind == reflect.Ptr {
		reflectType = reflectType.Elem()
		reflectKind = reflectType.Kind()
	}
	if reflectKind == reflect.Slice || reflectKind == reflect.Array {
		reflectType = reflectType.Elem()
		reflectKind = reflectType.Kind()
	}
	out.BindToAttrType = reflectType
	return
}

type doScanListInput struct {
	Model              *Model
	Result             Result
	StructSlicePointer interface{}
	StructSliceValue   reflect.Value
	BindToAttrName     string
	RelationAttrName   string
	RelationFields     string
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

	length := len(in.Result)
	if length == 0 {
		// The pointed slice is not empty.
		if in.StructSliceValue.Len() > 0 {
			// It here checks if it has struct item, which is already initialized.
			// It then returns error to warn the developer its empty and no conversion.
			if v := in.StructSliceValue.Index(0); v.Kind() != reflect.Ptr {
				return sql.ErrNoRows
			}
		}
		// Do nothing for empty struct slice.
		return nil
	}
	var (
		arrayValue    reflect.Value // Like: []*Entity
		arrayItemType reflect.Type  // Like: *Entity
		reflectType   = reflect.TypeOf(in.StructSlicePointer)
	)
	if in.StructSliceValue.Len() > 0 {
		arrayValue = in.StructSliceValue
	} else {
		arrayValue = reflect.MakeSlice(reflectType.Elem(), length, length)
	}

	// Slice element item.
	arrayItemType = arrayValue.Index(0).Type()

	// Relation variables.
	var (
		relationDataMap         map[string]Value
		relationFromFieldName   string // Eg: relationKV: id:uid  -> id
		relationBindToFieldName string // Eg: relationKV: id:uid  -> uid
	)
	if len(in.RelationFields) > 0 {
		// The relation key string of table field name and attribute name
		// can be joined with char '=' or ':'.
		array := gstr.SplitAndTrim(in.RelationFields, "=")
		if len(array) == 1 {
			// Compatible with old splitting char ':'.
			array = gstr.SplitAndTrim(in.RelationFields, ":")
		}
		if len(array) == 1 {
			// The relation names are the same.
			array = []string{in.RelationFields, in.RelationFields}
		}
		if len(array) == 2 {
			// Defined table field to relation attribute name.
			// Like:
			// uid:Uid
			// uid:UserId
			relationFromFieldName = array[0]
			relationBindToFieldName = array[1]
			if key, _ := gutil.MapPossibleItemByKey(in.Result[0].Map(), relationFromFieldName); key == "" {
				return gerror.NewCodef(
					gcode.CodeInvalidParameter,
					`cannot find possible related table field name "%s" from given relation fields "%s"`,
					relationFromFieldName,
					in.RelationFields,
				)
			} else {
				relationFromFieldName = key
			}
		} else {
			return gerror.NewCode(
				gcode.CodeInvalidParameter,
				`parameter relationKV should be format of "ResultFieldName:BindToAttrName"`,
			)
		}
		if relationFromFieldName != "" {
			// Note that the value might be type of slice.
			relationDataMap = in.Result.MapKeyValue(relationFromFieldName)
		}
		if len(relationDataMap) == 0 {
			return gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`cannot find the relation data map, maybe invalid relation fields given "%v"`,
				in.RelationFields,
			)
		}
	}
	// Bind to target attribute.
	var (
		ok              bool
		bindToAttrValue reflect.Value
		bindToAttrKind  reflect.Kind
		bindToAttrType  reflect.Type
		bindToAttrField reflect.StructField
	)
	if arrayItemType.Kind() == reflect.Ptr {
		if bindToAttrField, ok = arrayItemType.Elem().FieldByName(in.BindToAttrName); !ok {
			return gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid parameter bindToAttrName: cannot find attribute with name "%s" from slice element`,
				in.BindToAttrName,
			)
		}
	} else {
		if bindToAttrField, ok = arrayItemType.FieldByName(in.BindToAttrName); !ok {
			return gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid parameter bindToAttrName: cannot find attribute with name "%s" from slice element`,
				in.BindToAttrName,
			)
		}
	}
	bindToAttrType = bindToAttrField.Type
	bindToAttrKind = bindToAttrType.Kind()

	// Bind to relation conditions.
	var (
		relationFromAttrValue          reflect.Value
		relationFromAttrField          reflect.Value
		relationBindToFieldNameChecked bool
	)
	for i := 0; i < arrayValue.Len(); i++ {
		arrayElemValue := arrayValue.Index(i)
		// The FieldByName should be called on non-pointer reflect.Value.
		if arrayElemValue.Kind() == reflect.Ptr {
			// Like: []*Entity
			arrayElemValue = arrayElemValue.Elem()
			if !arrayElemValue.IsValid() {
				// The element is nil, then create one and set it to the slice.
				// The "reflect.New(itemType.Elem())" creates a new element and returns the address of it.
				// For example:
				// reflect.New(itemType.Elem())        => *Entity
				// reflect.New(itemType.Elem()).Elem() => Entity
				arrayElemValue = reflect.New(arrayItemType.Elem()).Elem()
				arrayValue.Index(i).Set(arrayElemValue.Addr())
			}
		} else {
			// Like: []Entity
		}
		bindToAttrValue = arrayElemValue.FieldByName(in.BindToAttrName)
		if in.RelationAttrName != "" {
			// Attribute value of current slice element.
			relationFromAttrValue = arrayElemValue.FieldByName(in.RelationAttrName)
			if relationFromAttrValue.Kind() == reflect.Ptr {
				relationFromAttrValue = relationFromAttrValue.Elem()
			}
		} else {
			// Current slice element.
			relationFromAttrValue = arrayElemValue
		}
		if len(relationDataMap) > 0 && !relationFromAttrValue.IsValid() {
			return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid relation fields specified: "%v"`, in.RelationFields)
		}
		// Check and find possible bind to attribute name.
		if in.RelationFields != "" && !relationBindToFieldNameChecked {
			relationFromAttrField = relationFromAttrValue.FieldByName(relationBindToFieldName)
			if !relationFromAttrField.IsValid() {
				fieldMap, _ := gstructs.FieldMap(gstructs.FieldMapInput{
					Pointer:         relationFromAttrValue,
					RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
				})
				if key, _ := gutil.MapPossibleItemByKey(gconv.Map(fieldMap), relationBindToFieldName); key == "" {
					return gerror.NewCodef(
						gcode.CodeInvalidParameter,
						`cannot find possible related attribute name "%s" from given relation fields "%s"`,
						relationBindToFieldName,
						in.RelationFields,
					)
				} else {
					relationBindToFieldName = key
				}
			}
			relationBindToFieldNameChecked = true
		}
		switch bindToAttrKind {
		case reflect.Array, reflect.Slice:
			if len(relationDataMap) > 0 {
				relationFromAttrField = relationFromAttrValue.FieldByName(relationBindToFieldName)
				if relationFromAttrField.IsValid() {
					results := make(Result, 0)
					for _, v := range relationDataMap[gconv.String(relationFromAttrField.Interface())].Slice() {
						results = append(results, v.(Record))
					}
					if err = results.Structs(bindToAttrValue.Addr()); err != nil {
						return err
					}
					// Recursively Scan.
					if in.Model != nil {
						if err = in.Model.doWithScanStructs(bindToAttrValue.Addr()); err != nil {
							return nil
						}
					}
				} else {
					// Maybe the attribute does not exist yet.
					return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid relation fields specified: "%v"`, in.RelationFields)
				}
			} else {
				return gerror.NewCodef(
					gcode.CodeInvalidParameter,
					`relationKey should not be empty as field "%s" is slice`,
					in.BindToAttrName,
				)
			}

		case reflect.Ptr:
			var element reflect.Value
			if bindToAttrValue.IsNil() {
				element = reflect.New(bindToAttrType.Elem()).Elem()
			} else {
				element = bindToAttrValue.Elem()
			}
			if len(relationDataMap) > 0 {
				relationFromAttrField = relationFromAttrValue.FieldByName(relationBindToFieldName)
				if relationFromAttrField.IsValid() {
					v := relationDataMap[gconv.String(relationFromAttrField.Interface())]
					if v == nil {
						// There's no relational data.
						continue
					}
					if v.IsSlice() {
						if err = v.Slice()[0].(Record).Struct(element); err != nil {
							return err
						}
					} else {
						if err = v.Val().(Record).Struct(element); err != nil {
							return err
						}
					}
				} else {
					// Maybe the attribute does not exist yet.
					return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid relation fields specified: "%v"`, in.RelationFields)
				}
			} else {
				if i >= len(in.Result) {
					// There's no relational data.
					continue
				}
				v := in.Result[i]
				if v == nil {
					// There's no relational data.
					continue
				}
				if err = v.Struct(element); err != nil {
					return err
				}
			}
			// Recursively Scan.
			if in.Model != nil {
				if err = in.Model.doWithScanStruct(element); err != nil {
					return err
				}
			}
			bindToAttrValue.Set(element.Addr())

		case reflect.Struct:
			if len(relationDataMap) > 0 {
				relationFromAttrField = relationFromAttrValue.FieldByName(relationBindToFieldName)
				if relationFromAttrField.IsValid() {
					relationDataItem := relationDataMap[gconv.String(relationFromAttrField.Interface())]
					if relationDataItem == nil {
						// There's no relational data.
						continue
					}
					if relationDataItem.IsSlice() {
						if err = relationDataItem.Slice()[0].(Record).Struct(bindToAttrValue); err != nil {
							return err
						}
					} else {
						if err = relationDataItem.Val().(Record).Struct(bindToAttrValue); err != nil {
							return err
						}
					}
				} else {
					// Maybe the attribute does not exist yet.
					return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid relation fields specified: "%v"`, in.RelationFields)
				}
			} else {
				if i >= len(in.Result) {
					// There's no relational data.
					continue
				}
				relationDataItem := in.Result[i]
				if relationDataItem == nil {
					// There's no relational data.
					continue
				}
				if err = relationDataItem.Struct(bindToAttrValue); err != nil {
					return err
				}
			}
			// Recursively Scan.
			if in.Model != nil {
				if err = in.Model.doWithScanStruct(bindToAttrValue); err != nil {
					return err
				}
			}

		default:
			return gerror.NewCodef(gcode.CodeInvalidParameter, `unsupported attribute type: %s`, bindToAttrKind.String())
		}
	}
	reflect.ValueOf(in.StructSlicePointer).Elem().Set(arrayValue)
	return nil
}
