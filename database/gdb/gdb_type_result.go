// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
	"math"
	"reflect"

	"github.com/gogf/gf/encoding/gparser"
)

// IsEmpty checks and returns whether <r> is empty.
func (r Result) IsEmpty() bool {
	return r.Len() == 0
}

// Len returns the length of result list.
func (r Result) Len() int {
	return len(r)
}

// Size is alias of function Len.
func (r Result) Size() int {
	return r.Len()
}

// Chunk splits an Result into multiple Results,
// the size of each array is determined by <size>.
// The last chunk may contain less than size elements.
func (r Result) Chunk(size int) []Result {
	if size < 1 {
		return nil
	}
	length := len(r)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n []Result
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}
		n = append(n, r[i*size:end])
		i++
	}
	return n
}

// Json converts <r> to JSON format content.
func (r Result) Json() string {
	content, _ := gparser.VarToJson(r.List())
	return string(content)
}

// Xml converts <r> to XML format content.
func (r Result) Xml(rootTag ...string) string {
	content, _ := gparser.VarToXml(r.List(), rootTag...)
	return string(content)
}

// List converts <r> to a List.
func (r Result) List() List {
	list := make(List, len(r))
	for k, v := range r {
		list[k] = v.Map()
	}
	return list
}

// Array retrieves and returns specified column values as slice.
// The parameter <field> is optional is the column field is only one.
func (r Result) Array(field ...string) []Value {
	array := make([]Value, len(r))
	if len(r) == 0 {
		return array
	}
	key := ""
	if len(field) > 0 && field[0] != "" {
		key = field[0]
	} else {
		for k, _ := range r[0] {
			key = k
			break
		}
	}
	for k, v := range r {
		array[k] = v[key]
	}
	return array
}

// MapKeyValue converts <r> to a map[string]Value of which key is specified by <key>.
// Note that the item value may be type of slice.
func (r Result) MapKeyValue(key string) map[string]Value {
	var (
		s              = ""
		m              = make(map[string]Value)
		tempMap        = make(map[string][]interface{})
		hasMultiValues bool
	)
	for _, item := range r {
		if k, ok := item[key]; ok {
			s = k.String()
			tempMap[s] = append(tempMap[s], item)
			if len(tempMap[s]) > 1 {
				hasMultiValues = true
			}
		}
	}
	for k, v := range tempMap {
		if hasMultiValues {
			m[k] = gvar.New(v)
		} else {
			m[k] = gvar.New(v[0])
		}
	}
	return m
}

// MapKeyStr converts <r> to a map[string]Map of which key is specified by <key>.
func (r Result) MapKeyStr(key string) map[string]Map {
	m := make(map[string]Map)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.String()] = item.Map()
		}
	}
	return m
}

// MapKeyInt converts <r> to a map[int]Map of which key is specified by <key>.
func (r Result) MapKeyInt(key string) map[int]Map {
	m := make(map[int]Map)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.Int()] = item.Map()
		}
	}
	return m
}

// MapKeyUint converts <r> to a map[uint]Map of which key is specified by <key>.
func (r Result) MapKeyUint(key string) map[uint]Map {
	m := make(map[uint]Map)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.Uint()] = item.Map()
		}
	}
	return m
}

// RecordKeyInt converts <r> to a map[int]Record of which key is specified by <key>.
func (r Result) RecordKeyStr(key string) map[string]Record {
	m := make(map[string]Record)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.String()] = item
		}
	}
	return m
}

// RecordKeyInt converts <r> to a map[int]Record of which key is specified by <key>.
func (r Result) RecordKeyInt(key string) map[int]Record {
	m := make(map[int]Record)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.Int()] = item
		}
	}
	return m
}

// RecordKeyUint converts <r> to a map[uint]Record of which key is specified by <key>.
func (r Result) RecordKeyUint(key string) map[uint]Record {
	m := make(map[uint]Record)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.Uint()] = item
		}
	}
	return m
}

// Structs converts <r> to struct slice.
// Note that the parameter <pointer> should be type of *[]struct/*[]*struct.
func (r Result) Structs(pointer interface{}) (err error) {
	var (
		reflectValue = reflect.ValueOf(pointer)
		reflectKind  = reflectValue.Kind()
	)
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
		reflectType = reflect.TypeOf(pointer)
		array       = reflect.MakeSlice(reflectType.Elem(), length, length)
		itemType    = array.Index(0).Type()
		itemKind    = itemType.Kind()
	)
	for i := 0; i < length; i++ {
		if itemKind == reflect.Ptr {
			e := reflect.New(itemType.Elem()).Elem()
			if err = r[i].Struct(e); err != nil {
				return fmt.Errorf(`slice element conversion failed: %s`, err.Error())
			}
			array.Index(i).Set(e.Addr())
		} else {
			e := reflect.New(itemType).Elem()
			if err = r[i].Struct(e); err != nil {
				return fmt.Errorf(`slice element conversion failed: %s`, err.Error())
			}
			array.Index(i).Set(e)
		}
	}
	reflect.ValueOf(pointer).Elem().Set(array)
	return nil
}

// ScanList converts <r> to struct slice which contains other complex struct attributes.
// Note that the parameter <listPointer> should be type of *[]struct/*[]*struct.
func (r Result) ScanList(listPointer interface{}, attributeName string, relation ...string) (err error) {
	// Necessary checks for parameters.
	if attributeName == "" {
		return errors.New(`attributeName should not be empty`)
	}
	if len(relation) > 0 {
		if len(relation) < 2 {
			return errors.New(`relation name and key should are both necessary`)
		}
		if relation[0] == "" || relation[1] == "" {
			return errors.New(`relation name and key should not be empty`)
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
			return fmt.Errorf(`invalid relation: %s, %s`, relation[0], relation[1])
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
					return fmt.Errorf(`invalid relation: %s, %s`, relation[0], relation[1])
				}
			} else {
				return fmt.Errorf(`relationKey should not be empty as field "%s" is slice`, attributeName)
			}

		case reflect.Ptr:
			e := reflect.New(attrType.Elem()).Elem()
			if len(relationDataMap) > 0 {
				relationField = relationValue.FieldByName(relationAttrName)
				if relationField.IsValid() {
					if err = gconv.Struct(relationDataMap[gconv.String(relationField.Interface())], e); err != nil {
						return err
					}
				} else {
					// May be the attribute does not exist yet.
					return fmt.Errorf(`invalid relation: %s, %s`, relation[0], relation[1])
				}
			} else {
				if err = gconv.Struct(r[i], e); err != nil {
					return err
				}
			}
			attrValue.Set(e.Addr())

		case reflect.Struct:
			e := reflect.New(attrType).Elem()
			if len(relationDataMap) > 0 {
				relationField = relationValue.FieldByName(relationAttrName)
				if relationField.IsValid() {
					if err = gconv.Struct(relationDataMap[gconv.String(relationField.Interface())], e); err != nil {
						return err
					}
				} else {
					// May be the attribute does not exist yet.
					return fmt.Errorf(`invalid relation: %s, %s`, relation[0], relation[1])
				}
			} else {
				if err = gconv.Struct(r[i], e); err != nil {
					return err
				}
			}
			attrValue.Set(e)

		case reflect.Map:
			// TODO
		}
	}
	reflect.ValueOf(listPointer).Elem().Set(arrayValue)
	return nil
}
