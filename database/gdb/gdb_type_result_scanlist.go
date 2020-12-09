// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"fmt"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
	"reflect"
)

// ScanList converts <r> to struct slice which contains other complex struct attributes.
// Note that the parameter <listPointer> should be type of *[]struct/*[]*struct.
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
// The parameters "User"/"UserDetail"/"UserScores" in the example codes specify the target attribute struct
// that current result will be bound to.
//
// The "uid" in the example codes is the table field name of the result, and the "Uid" is the relational
// struct attribute name - not the attribute name of the bound to target. In the example codes, it's attribute
// name "Uid" of "User" of entity "Entity". It automatically calculates the HasOne/HasMany relationship with
// given <relation> parameter.
//
// See the example or unit testing cases for clear understanding for this function.
func (r Result) ScanList(listPointer interface{}, attributeName string, relation ...string) (err error) {
	// Necessary checks for parameters.
	if attributeName == "" {
		return gerror.New(`attributeName should not be empty`)
	}
	if len(relation) > 0 {
		if len(relation) < 2 {
			return gerror.New(`relation name and key should are both necessary`)
		}
		if relation[0] == "" || relation[1] == "" {
			return gerror.New(`relation name and key should not be empty`)
		}
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
		return fmt.Errorf("parameter should be type of *[]struct/*[]*struct, but got: %v", reflectKind)
	}
	reflectValue = reflectValue.Elem()
	reflectKind = reflectValue.Kind()
	if reflectKind != reflect.Slice && reflectKind != reflect.Array {
		return fmt.Errorf("parameter should be type of *[]struct/*[]*struct, but got: %v", reflectKind)
	}
	length := len(r)
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
		relationDataMap   map[string]Value
		relationFieldName string
		relationAttrName  string
	)
	if len(relation) > 0 {
		array := gstr.Split(relation[1], ":")
		if len(array) > 1 {
			// Defined table field to relation attribute name.
			// Like:
			// uid:Uid
			// uid:UserId
			relationFieldName = array[0]
			relationAttrName = array[1]
		} else {
			relationAttrName = relation[1]
			// Find the possible map key by given only struct attribute name.
			// Like:
			// Uid
			if k, _ := gutil.MapPossibleItemByKey(r[0].Map(), relation[1]); k != "" {
				relationFieldName = k
			}
		}
		if relationFieldName != "" {
			relationDataMap = r.MapKeyValue(relationFieldName)
		}
		if len(relationDataMap) == 0 {
			return fmt.Errorf(`cannot find the relation data map, maybe invalid relation key given: %s`, relation[1])
		}
	}
	// Bind to target attribute.
	var (
		ok        bool
		attrValue reflect.Value
		attrKind  reflect.Kind
		attrType  reflect.Type
		attrField reflect.StructField
	)
	if arrayItemType.Kind() == reflect.Ptr {
		if attrField, ok = arrayItemType.Elem().FieldByName(attributeName); !ok {
			return fmt.Errorf(`invalid field name: %s`, attributeName)
		}
	} else {
		if attrField, ok = arrayItemType.FieldByName(attributeName); !ok {
			return fmt.Errorf(`invalid field name: %s`, attributeName)
		}
	}
	attrType = attrField.Type
	attrKind = attrType.Kind()

	// Bind to relation conditions.
	var (
		relationValue reflect.Value
		relationField reflect.Value
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
		attrValue = arrayElemValue.FieldByName(attributeName)
		if len(relation) > 0 {
			relationValue = arrayElemValue.FieldByName(relation[0])
			if relationValue.Kind() == reflect.Ptr {
				relationValue = relationValue.Elem()
			}
		}
		if len(relationDataMap) > 0 && !relationValue.IsValid() {
			return fmt.Errorf(`invalid relation: "%s:%s"`, relation[0], relation[1])
		}
		switch attrKind {
		case reflect.Array, reflect.Slice:
			if len(relationDataMap) > 0 {
				relationField = relationValue.FieldByName(relationAttrName)
				if relationField.IsValid() {
					if err = gconv.Structs(
						relationDataMap[gconv.String(relationField.Interface())],
						attrValue.Addr(),
					); err != nil {
						return err
					}
				} else {
					// May be the attribute does not exist yet.
					return fmt.Errorf(`invalid relation: "%s:%s"`, relation[0], relation[1])
				}
			} else {
				return fmt.Errorf(`relationKey should not be empty as field "%s" is slice`, attributeName)
			}

		case reflect.Ptr:
			e := reflect.New(attrType.Elem()).Elem()
			if len(relationDataMap) > 0 {
				relationField = relationValue.FieldByName(relationAttrName)
				if relationField.IsValid() {
					v := relationDataMap[gconv.String(relationField.Interface())]
					if v == nil {
						// There's no relational data.
						continue
					}
					if err = gconv.Struct(v, e); err != nil {
						return err
					}
				} else {
					// May be the attribute does not exist yet.
					return fmt.Errorf(`invalid relation: "%s:%s"`, relation[0], relation[1])
				}
			} else {
				v := r[i]
				if v == nil {
					// There's no relational data.
					continue
				}
				if err = gconv.Struct(v, e); err != nil {
					return err
				}
			}
			attrValue.Set(e.Addr())

		case reflect.Struct:
			e := reflect.New(attrType).Elem()
			if len(relationDataMap) > 0 {
				relationField = relationValue.FieldByName(relationAttrName)
				if relationField.IsValid() {
					v := relationDataMap[gconv.String(relationField.Interface())]
					if v == nil {
						// There's no relational data.
						continue
					}
					if err = gconv.Struct(v, e); err != nil {
						return err
					}
				} else {
					// May be the attribute does not exist yet.
					return fmt.Errorf(`invalid relation: "%s:%s"`, relation[0], relation[1])
				}
			} else {
				v := r[i]
				if v == nil {
					// There's no relational data.
					continue
				}
				if err = gconv.Struct(v, e); err != nil {
					return err
				}
			}
			attrValue.Set(e)

		default:
			return fmt.Errorf(`unsupport attribute type: %s`, attrKind.String())
		}
	}
	reflect.ValueOf(listPointer).Elem().Set(arrayValue)
	return nil
}
