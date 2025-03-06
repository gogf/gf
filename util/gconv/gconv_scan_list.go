// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"reflect"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gstructs"
)

// ScanList converts `structSlice` to struct slice which contains other complex struct attributes.
// Note that the parameter `structSlicePointer` should be type of *[]struct/*[]*struct.
//
// Usage example 1: Normal attribute struct relation:
//
//	type EntityUser struct {
//	    Uid  int
//	    Name string
//	}
//
//	type EntityUserDetail struct {
//	    Uid     int
//	    Address string
//	}
//
//	type EntityUserScores struct {
//	    Id     int
//	    Uid    int
//	    Score  int
//	    Course string
//	}
//
//	type Entity struct {
//	    User       *EntityUser
//	    UserDetail *EntityUserDetail
//	    UserScores []*EntityUserScores
//	}
//
// var users []*Entity
// var userRecords   = EntityUser{Uid: 1, Name:"john"}
// var detailRecords = EntityUser{Uid: 1, Address: "chengdu"}
// var scoresRecords = EntityUser{Id: 1, Uid: 1, Score: 100, Course: "math"}
// ScanList(userRecords, &users, "User")
// ScanList(userRecords, &users, "User", "uid")
// ScanList(detailRecords, &users, "UserDetail", "User", "uid:Uid")
// ScanList(scoresRecords, &users, "UserScores", "User", "uid:Uid")
// ScanList(scoresRecords, &users, "UserScores", "User", "uid")
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
// var userRecords   = EntityUser{Uid: 1, Name:"john"}
// var detailRecords = EntityUser{Uid: 1, Address: "chengdu"}
// var scoresRecords = EntityUser{Id: 1, Uid: 1, Score: 100, Course: "math"}
// ScanList(userRecords, &users)
// ScanList(detailRecords, &users, "UserDetail", "uid")
// ScanList(scoresRecords, &users, "UserScores", "uid")
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
func ScanList(
	structSlice any, structSlicePointer any, bindToAttrName string, relationAttrNameAndFields ...string,
) (err error) {
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
	return doScanList(structSlice, structSlicePointer, bindToAttrName, relationAttrName, relationFields)
}

// doScanList converts `structSlice` to struct slice which contains other complex struct attributes recursively.
// Note that the parameter `structSlicePointer` should be type of *[]struct/*[]*struct.
func doScanList(
	structSlice any, structSlicePointer any, bindToAttrName, relationAttrName, relationFields string,
) (err error) {
	var (
		maps    = Maps(structSlice)
		lenMaps = len(maps)
	)
	if lenMaps == 0 {
		return nil
	}
	// Necessary checks for parameters.
	if bindToAttrName == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, `bindToAttrName should not be empty`)
	}

	if relationAttrName == "." {
		relationAttrName = ""
	}

	var (
		reflectValue = reflect.ValueOf(structSlicePointer)
		reflectKind  = reflectValue.Kind()
	)
	if reflectKind == reflect.Interface {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	if reflectKind != reflect.Ptr {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"structSlicePointer should be type of *[]struct/*[]*struct, but got: %v",
			reflectKind,
		)
	}
	reflectValue = reflectValue.Elem()
	reflectKind = reflectValue.Kind()
	if reflectKind != reflect.Slice && reflectKind != reflect.Array {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"structSlicePointer should be type of *[]struct/*[]*struct, but got: %v",
			reflectKind,
		)
	}
	var (
		arrayValue    reflect.Value // Like: []*Entity
		arrayItemType reflect.Type  // Like: *Entity
		reflectType   = reflect.TypeOf(structSlicePointer)
	)
	if reflectValue.Len() > 0 {
		arrayValue = reflectValue
	} else {
		arrayValue = reflect.MakeSlice(reflectType.Elem(), lenMaps, lenMaps)
	}

	// Slice element item.
	arrayItemType = arrayValue.Index(0).Type()

	// Relation variables.
	var (
		relationDataMap         map[string]any
		relationFromFieldName   string // Eg: relationKV: id:uid  -> id
		relationBindToFieldName string // Eg: relationKV: id:uid  -> uid
	)
	if len(relationFields) > 0 {
		// The relation key string of table field name and attribute name
		// can be joined with char '=' or ':'.
		array := utils.SplitAndTrim(relationFields, "=")
		if len(array) == 1 {
			// Compatible with old splitting char ':'.
			array = utils.SplitAndTrim(relationFields, ":")
		}
		if len(array) == 1 {
			// The relation names are the same.
			array = []string{relationFields, relationFields}
		}
		if len(array) == 2 {
			// Defined table field to relation attribute name.
			// Like:
			// uid:Uid
			// uid:UserId
			relationFromFieldName = array[0]
			relationBindToFieldName = array[1]
			if key, _ := utils.MapPossibleItemByKey(maps[0], relationFromFieldName); key == "" {
				return gerror.NewCodef(
					gcode.CodeInvalidParameter,
					`cannot find possible related table field name "%s" from given relation fields "%s"`,
					relationFromFieldName,
					relationFields,
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
			relationDataMap = utils.ListToMapByKey(maps, relationFromFieldName)
		}
		if len(relationDataMap) == 0 {
			return gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`cannot find the relation data map, maybe invalid relation fields given "%v"`,
				relationFields,
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
		if bindToAttrField, ok = arrayItemType.Elem().FieldByName(bindToAttrName); !ok {
			return gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid parameter bindToAttrName: cannot find attribute with name "%s" from slice element`,
				bindToAttrName,
			)
		}
	} else {
		if bindToAttrField, ok = arrayItemType.FieldByName(bindToAttrName); !ok {
			return gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid parameter bindToAttrName: cannot find attribute with name "%s" from slice element`,
				bindToAttrName,
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
		bindToAttrValue = arrayElemValue.FieldByName(bindToAttrName)
		if relationAttrName != "" {
			// Attribute value of current slice element.
			relationFromAttrValue = arrayElemValue.FieldByName(relationAttrName)
			if relationFromAttrValue.Kind() == reflect.Ptr {
				relationFromAttrValue = relationFromAttrValue.Elem()
			}
		} else {
			// Current slice element.
			relationFromAttrValue = arrayElemValue
		}
		if len(relationDataMap) > 0 && !relationFromAttrValue.IsValid() {
			return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid relation fields specified: "%v"`, relationFields)
		}
		// Check and find possible bind to attribute name.
		if relationFields != "" && !relationBindToFieldNameChecked {
			relationFromAttrField = relationFromAttrValue.FieldByName(relationBindToFieldName)
			if !relationFromAttrField.IsValid() {
				var (
					fieldMap, _ = gstructs.FieldMap(gstructs.FieldMapInput{
						Pointer:         relationFromAttrValue,
						RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
					})
				)
				if key, _ := utils.MapPossibleItemByKey(Map(fieldMap), relationBindToFieldName); key == "" {
					return gerror.NewCodef(
						gcode.CodeInvalidParameter,
						`cannot find possible related attribute name "%s" from given relation fields "%s"`,
						relationBindToFieldName,
						relationFields,
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
					// results := make(Result, 0)
					results := make([]any, 0)
					for _, v := range SliceAny(relationDataMap[String(relationFromAttrField.Interface())]) {
						item := v
						results = append(results, item)
					}
					if err = Structs(results, bindToAttrValue.Addr()); err != nil {
						return err
					}
				} else {
					// Maybe the attribute does not exist yet.
					return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid relation fields specified: "%v"`, relationFields)
				}
			} else {
				return gerror.NewCodef(
					gcode.CodeInvalidParameter,
					`relationKey should not be empty as field "%s" is slice`,
					bindToAttrName,
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
					v := relationDataMap[String(relationFromAttrField.Interface())]
					if v == nil {
						// There's no relational data.
						continue
					}
					if utils.IsSlice(v) {
						if err = Struct(SliceAny(v)[0], element); err != nil {
							return err
						}
					} else {
						if err = Struct(v, element); err != nil {
							return err
						}
					}
				} else {
					// Maybe the attribute does not exist yet.
					return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid relation fields specified: "%v"`, relationFields)
				}
			} else {
				if i >= len(maps) {
					// There's no relational data.
					continue
				}
				v := maps[i]
				if v == nil {
					// There's no relational data.
					continue
				}
				if err = Struct(v, element); err != nil {
					return err
				}
			}
			bindToAttrValue.Set(element.Addr())

		case reflect.Struct:
			if len(relationDataMap) > 0 {
				relationFromAttrField = relationFromAttrValue.FieldByName(relationBindToFieldName)
				if relationFromAttrField.IsValid() {
					relationDataItem := relationDataMap[String(relationFromAttrField.Interface())]
					if relationDataItem == nil {
						// There's no relational data.
						continue
					}
					if utils.IsSlice(relationDataItem) {
						if err = Struct(SliceAny(relationDataItem)[0], bindToAttrValue); err != nil {
							return err
						}
					} else {
						if err = Struct(relationDataItem, bindToAttrValue); err != nil {
							return err
						}
					}
				} else {
					// Maybe the attribute does not exist yet.
					return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid relation fields specified: "%v"`, relationFields)
				}
			} else {
				if i >= len(maps) {
					// There's no relational data.
					continue
				}
				relationDataItem := maps[i]
				if relationDataItem == nil {
					// There's no relational data.
					continue
				}
				if err = Struct(relationDataItem, bindToAttrValue); err != nil {
					return err
				}
			}

		default:
			return gerror.NewCodef(gcode.CodeInvalidParameter, `unsupported attribute type: %s`, bindToAttrKind.String())
		}
	}
	reflect.ValueOf(structSlicePointer).Elem().Set(arrayValue)
	return nil
}
