// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
	"reflect"
)

// ScanList converts `r` to struct slice which contains other complex struct attributes.
// Note that the parameter `listPointer` should be type of *[]struct/*[]*struct.
// Usage example:
//
// type Entity struct {
// 	   User       *EntityUser
// 	   UserDetail *EntityUserDetail
//	   UserScores []*EntityUserScores
// }
// var users []*Entity
// or
// var users []Entity
//
// ScanList(&users, "User")
// ScanList(&users, "UserDetail", "User", "uid:Uid")
// ScanList(&users, "UserScores", "User", "uid:Uid")
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
func (r Result) ScanList(listPointer interface{}, bindToAttrName string, relationKV ...string) (err error) {
	return doScanList(nil, r, listPointer, bindToAttrName, relationKV...)
}

// doScanList converts `result` to struct slice which contains other complex struct attributes recursively.
// The parameter `model` is used for recursively scanning purpose, which means, it can scans the attribute struct/structs recursively but
// it needs the Model for database accessing.
// Note that the parameter `listPointer` should be type of *[]struct/*[]*struct.
func doScanList(model *Model, result Result, listPointer interface{}, bindToAttrName string, relationKV ...string) (err error) {
	if result.IsEmpty() {
		return nil
	}
	// Necessary checks for parameters.
	if bindToAttrName == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, `bindToAttrName should not be empty`)
	}

	var (
		reflectValue = reflect.ValueOf(listPointer)
		reflectKind  = reflectValue.Kind()
	)
	if reflectKind == reflect.Interface {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	if reflectKind != reflect.Ptr {
		return gerror.NewCodef(gcode.CodeInvalidParameter, "listPointer should be type of *[]struct/*[]*struct, but got: %v", reflectKind)
	}
	reflectValue = reflectValue.Elem()
	reflectKind = reflectValue.Kind()
	if reflectKind != reflect.Slice && reflectKind != reflect.Array {
		return gerror.NewCodef(gcode.CodeInvalidParameter, "listPointer should be type of *[]struct/*[]*struct, but got: %v", reflectKind)
	}
	length := len(result)
	if length == 0 {
		// The pointed slice is not empty.
		if reflectValue.Len() > 0 {
			// It here checks if it has struct item, which is already initialized.
			// It then returns error to warn the developer its empty and no conversion.
			if v := reflectValue.Index(0); v.Kind() != reflect.Ptr {
				return sql.ErrNoRows
			}
		}
		// Do nothing for empty struct slice.
		return nil
	}
	var (
		arrayValue    reflect.Value // Like: []*Entity
		arrayItemType reflect.Type  // Like: *Entity
		reflectType   = reflect.TypeOf(listPointer)
	)
	if reflectValue.Len() > 0 {
		arrayValue = reflectValue
	} else {
		arrayValue = reflect.MakeSlice(reflectType.Elem(), length, length)
	}

	// Slice element item.
	arrayItemType = arrayValue.Index(0).Type()

	// Relation variables.
	var (
		relationKVStr             string
		relationDataMap           map[string]Value
		relationFromAttrName      string // Eg: relationKV: User, uid:Uid -> User
		relationResultFieldName   string // Eg: relationKV: uid:Uid       -> uid
		relationBindToSubAttrName string // Eg: relationKV: uid:Uid       -> Uid
	)
	if len(relationKV) > 0 {
		if len(relationKV) == 1 {
			relationKVStr = relationKV[0]
		} else {
			relationFromAttrName = relationKV[0]
			relationKVStr = relationKV[1]
		}
		// The relation key string of table filed name and attribute name
		// can be joined with char '=' or ':'.
		array := gstr.SplitAndTrim(relationKVStr, "=")
		if len(array) == 1 {
			// Compatible with old splitting char ':'.
			array = gstr.SplitAndTrim(relationKVStr, ":")
		}
		if len(array) == 1 {
			// The relation names are the same.
			array = []string{relationKVStr, relationKVStr}
		}
		if len(array) == 2 {
			// Defined table field to relation attribute name.
			// Like:
			// uid:Uid
			// uid:UserId
			relationResultFieldName = array[0]
			relationBindToSubAttrName = array[1]
			if key, _ := gutil.MapPossibleItemByKey(result[0].Map(), relationResultFieldName); key == "" {
				return gerror.NewCodef(
					gcode.CodeInvalidParameter,
					`cannot find possible related table field name "%s" from given relation key "%s"`,
					relationResultFieldName,
					relationKVStr,
				)
			} else {
				relationResultFieldName = key
			}
		} else {
			return gerror.NewCode(gcode.CodeInvalidParameter, `parameter relationKV should be format of "ResultFieldName:BindToAttrName"`)
		}
		if relationResultFieldName != "" {
			// Note that the value might be type of slice.
			relationDataMap = result.MapKeyValue(relationResultFieldName)
		}
		if len(relationDataMap) == 0 {
			return gerror.NewCodef(gcode.CodeInvalidParameter, `cannot find the relation data map, maybe invalid relation given "%v"`, relationKV)
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
		relationFromAttrValue            reflect.Value
		relationFromAttrField            reflect.Value
		relationBindToSubAttrNameChecked bool
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
		if relationFromAttrName != "" {
			// Attribute value of current slice element.
			relationFromAttrValue = arrayElemValue.FieldByName(relationFromAttrName)
			if relationFromAttrValue.Kind() == reflect.Ptr {
				relationFromAttrValue = relationFromAttrValue.Elem()
			}
		} else {
			// Current slice element.
			relationFromAttrValue = arrayElemValue
		}
		if len(relationDataMap) > 0 && !relationFromAttrValue.IsValid() {
			return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid relation specified: "%v"`, relationKV)
		}
		// Check and find possible bind to attribute name.
		if relationKVStr != "" && !relationBindToSubAttrNameChecked {
			relationFromAttrField = relationFromAttrValue.FieldByName(relationBindToSubAttrName)
			if !relationFromAttrField.IsValid() {
				var (
					relationFromAttrType = relationFromAttrValue.Type()
					filedMap             = make(map[string]interface{})
				)
				for i := 0; i < relationFromAttrType.NumField(); i++ {
					filedMap[relationFromAttrType.Field(i).Name] = struct{}{}
				}
				if key, _ := gutil.MapPossibleItemByKey(filedMap, relationBindToSubAttrName); key == "" {
					return gerror.NewCodef(
						gcode.CodeInvalidParameter,
						`cannot find possible related attribute name "%s" from given relation key "%s"`,
						relationBindToSubAttrName,
						relationKVStr,
					)
				} else {
					relationBindToSubAttrName = key
				}
			}
			relationBindToSubAttrNameChecked = true
		}
		switch bindToAttrKind {
		case reflect.Array, reflect.Slice:
			if len(relationDataMap) > 0 {
				relationFromAttrField = relationFromAttrValue.FieldByName(relationBindToSubAttrName)
				if relationFromAttrField.IsValid() {
					results := make(Result, 0)
					for _, v := range relationDataMap[gconv.String(relationFromAttrField.Interface())].Slice() {
						results = append(results, v.(Record))
					}
					if err = results.Structs(bindToAttrValue.Addr()); err != nil {
						return err
					}
					// Recursively Scan.
					if model != nil {
						if err = model.doWithScanStructs(bindToAttrValue.Addr()); err != nil {
							return nil
						}
					}
				} else {
					// May be the attribute does not exist yet.
					return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid relation specified: "%v"`, relationKV)
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
				relationFromAttrField = relationFromAttrValue.FieldByName(relationBindToSubAttrName)
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
					// May be the attribute does not exist yet.
					return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid relation specified: "%v"`, relationKV)
				}
			} else {
				if i >= len(result) {
					// There's no relational data.
					continue
				}
				v := result[i]
				if v == nil {
					// There's no relational data.
					continue
				}
				if err = v.Struct(element); err != nil {
					return err
				}
			}
			// Recursively Scan.
			if model != nil {
				if err = model.doWithScanStruct(element); err != nil {
					return err
				}
			}
			bindToAttrValue.Set(element.Addr())

		case reflect.Struct:
			if len(relationDataMap) > 0 {
				relationFromAttrField = relationFromAttrValue.FieldByName(relationBindToSubAttrName)
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
					// May be the attribute does not exist yet.
					return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid relation specified: "%v"`, relationKV)
				}
			} else {
				if i >= len(result) {
					// There's no relational data.
					continue
				}
				relationDataItem := result[i]
				if relationDataItem == nil {
					// There's no relational data.
					continue
				}
				if err = relationDataItem.Struct(bindToAttrValue); err != nil {
					return err
				}
			}
			// Recursively Scan.
			if model != nil {
				if err = model.doWithScanStruct(bindToAttrValue); err != nil {
					return err
				}
			}

		default:
			return gerror.NewCodef(gcode.CodeInvalidParameter, `unsupported attribute type: %s`, bindToAttrKind.String())
		}
	}
	reflect.ValueOf(listPointer).Elem().Set(arrayValue)
	return nil
}
