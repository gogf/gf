// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"fmt"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/empty"
	"github.com/gogf/gf/internal/json"
	"reflect"
	"regexp"
	"strings"

	"github.com/gogf/gf/internal/structs"
	"github.com/gogf/gf/internal/utils"
)

var (
	// replaceCharReg is the regular expression object for replacing chars
	// in map keys and attribute names.
	replaceCharReg, _ = regexp.Compile(`[\-\.\_\s]+`)
)

// Struct maps the params key-value pairs to the corresponding struct object's attributes.
// The third parameter <mapping> is unnecessary, indicating the mapping rules between the
// custom key name and the attribute name(case sensitive).
//
// Note:
// 1. The <params> can be any type of map/struct, usually a map.
// 2. The <pointer> should be type of *struct/**struct, which is a pointer to struct object
//    or struct pointer.
// 3. Only the public attributes of struct object can be mapped.
// 4. If <params> is a map, the key of the map <params> can be lowercase.
//    It will automatically convert the first letter of the key to uppercase
//    in mapping procedure to do the matching.
//    It ignores the map key, if it does not match.
func Struct(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	return doStruct(params, pointer, false, mapping...)
}

// StructDeep do Struct function recursively.
// See Struct.
func StructDeep(params interface{}, pointer interface{}, mapping ...map[string]string) error {
	return doStruct(params, pointer, true, mapping...)
}

// doStruct is the core internal converting function for any data to struct recursively or not.
func doStruct(params interface{}, pointer interface{}, recursive bool, mapping ...map[string]string) (err error) {
	if params == nil {
		// If <params> is nil, no conversion.
		return nil
	}
	if pointer == nil {
		return gerror.New("object pointer cannot be nil")
	}
	defer func() {
		// Catch the panic, especially the reflect operation panics.
		if e := recover(); e != nil {
			err = gerror.NewfSkip(1, "%v", e)
		}
	}()

	// If given <params> is JSON, it then uses json.Unmarshal doing the converting.
	switch r := params.(type) {
	case []byte:
		if json.Valid(r) {
			return json.Unmarshal(r, pointer)
		}
	case string:
		if paramsBytes := []byte(r); json.Valid(paramsBytes) {
			return json.Unmarshal(paramsBytes, pointer)
		}
	}

	// UnmarshalValue.
	// Assign value with interface UnmarshalValue.
	// Note that only pointer can implement interface UnmarshalValue.
	if v, ok := pointer.(apiUnmarshalValue); ok {
		return v.UnmarshalValue(params)
	}

	// paramsMap is the map[string]interface{} type variable for params.
	// DO NOT use MapDeep here.
	paramsMap := Map(params)
	if paramsMap == nil {
		return gerror.Newf("convert params to map failed: %v", params)
	}

	// Using reflect to do the converting,
	// it also supports type of reflect.Value for <pointer>(always in internal usage).
	elem, ok := pointer.(reflect.Value)
	if !ok {
		rv := reflect.ValueOf(pointer)
		if kind := rv.Kind(); kind != reflect.Ptr {
			return gerror.Newf("object pointer should be type of '*struct', but got '%v'", kind)
		}
		// Using IsNil on reflect.Ptr variable is OK.
		if !rv.IsValid() || rv.IsNil() {
			return gerror.New("object pointer cannot be nil")
		}
		elem = rv.Elem()
	}

	// Check if an invalid interface.
	if elem.Kind() == reflect.Interface {
		elem = elem.Elem()
		if !elem.IsValid() {
			return gerror.New("interface type converting is not supported")
		}
	}

	// It automatically creates struct object if necessary.
	// For example, if <pointer> is **User, then <elem> is *User, which is a pointer to User.
	if elem.Kind() == reflect.Ptr {
		if !elem.IsValid() || elem.IsNil() {
			e := reflect.New(elem.Type().Elem()).Elem()
			elem.Set(e.Addr())
			elem = e
		} else {
			elem = elem.Elem()
		}
	}

	// UnmarshalValue checks again.
	// Assign value with interface UnmarshalValue.
	// Note that only pointer can implement interface UnmarshalValue.
	if elem.Kind() == reflect.Struct && elem.CanAddr() {
		if v, ok := elem.Addr().Interface().(apiUnmarshalValue); ok {
			return v.UnmarshalValue(params)
		}
	}

	// It only performs one converting to the same attribute.
	// doneMap is used to check repeated converting, its key is the real attribute name
	// of the struct.
	doneMap := make(map[string]struct{})
	// It first checks the passed mapping rules.
	if len(mapping) > 0 && len(mapping[0]) > 0 {
		for mapK, mapV := range mapping[0] {
			// mapV is the the attribute name of the struct.
			if paramV, ok := paramsMap[mapK]; ok {
				doneMap[mapV] = struct{}{}
				if err := bindVarToStructAttr(elem, mapV, paramV, recursive, mapping...); err != nil {
					return err
				}
			}
		}
	}

	// The key of the attrMap is the attribute name of the struct,
	// and the value is its replaced name for later comparison to improve performance.
	var (
		tempName       string
		elemFieldType  reflect.StructField
		elemFieldValue reflect.Value
		elemType       = elem.Type()
		attrMap        = make(map[string]string)
	)
	for i := 0; i < elem.NumField(); i++ {
		elemFieldType = elemType.Field(i)
		// Only do converting to public attributes.
		if !utils.IsLetterUpper(elemFieldType.Name[0]) {
			continue
		}
		// Maybe it's struct/*struct.
		if recursive && elemFieldType.Anonymous {
			elemFieldValue = elem.Field(i)
			// Ignore the interface attribute if it's nil.
			if elemFieldValue.Kind() == reflect.Interface {
				elemFieldValue = elemFieldValue.Elem()
				if !elemFieldValue.IsValid() {
					continue
				}
			}
			if err = doStruct(paramsMap, elemFieldValue, recursive, mapping...); err != nil {
				return err
			}
		} else {
			tempName = elemFieldType.Name
			attrMap[tempName] = replaceCharReg.ReplaceAllString(tempName, "")
		}
	}
	if len(attrMap) == 0 {
		return nil
	}

	// The key of the tagMap is the attribute name of the struct,
	// and the value is its replaced tag name for later comparison to improve performance.
	tagMap := make(map[string]string)
	for k, v := range structs.TagMapName(pointer, StructTagPriority, true) {
		tagMap[v] = replaceCharReg.ReplaceAllString(k, "")
	}

	var (
		attrName  string
		checkName string
	)
	for mapK, mapV := range paramsMap {
		attrName = ""
		checkName = replaceCharReg.ReplaceAllString(mapK, "")
		// Loop to find the matched attribute name with or without
		// string cases and chars like '-'/'_'/'.'/' '.

		// Matching the parameters to struct tag names.
		// The <tagV> is the attribute name of the struct.
		for attrKey, cmpKey := range tagMap {
			if strings.EqualFold(checkName, cmpKey) {
				attrName = attrKey
				break
			}
		}

		// Matching the parameters to struct attributes.
		if attrName == "" {
			for attrKey, cmpKey := range attrMap {
				// Eg:
				// UserName  eq user_name
				// User-Name eq username
				// username  eq userName
				// etc.
				if strings.EqualFold(checkName, cmpKey) {
					attrName = attrKey
					break
				}
			}
		}

		// No matching, give up this attribute converting.
		if attrName == "" {
			continue
		}
		// If the attribute name is already checked converting, then skip it.
		if _, ok := doneMap[attrName]; ok {
			continue
		}
		// Mark it done.
		doneMap[attrName] = struct{}{}
		if err := bindVarToStructAttr(elem, attrName, mapV, recursive, mapping...); err != nil {
			return err
		}
	}
	// Recursively concerting for struct attributes with the same params map.
	if recursive && elem.Kind() == reflect.Struct {
		for i := 0; i < elemType.NumField(); i++ {
			// Only do converting to public attributes.
			if !utils.IsLetterUpper(elemType.Field(i).Name[0]) {
				continue
			}
			fieldValue := elem.Field(i)
			if fieldValue.Kind() == reflect.Struct {
				if err := doStruct(paramsMap, fieldValue, recursive, mapping...); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// bindVarToStructAttr sets value to struct object attribute by name.
func bindVarToStructAttr(elem reflect.Value, name string, value interface{}, recursive bool, mapping ...map[string]string) (err error) {
	structFieldValue := elem.FieldByName(name)
	if !structFieldValue.IsValid() {
		return nil
	}
	// CanSet checks whether attribute is public accessible.
	if !structFieldValue.CanSet() {
		return nil
	}
	defer func() {
		if e := recover(); e != nil {
			err = bindVarToReflectValue(structFieldValue, value, recursive, mapping...)
			if err != nil {
				err = gerror.Wrapf(err, `error binding value to attribute "%s"`, name)
			}
		}
	}()
	if empty.IsNil(value) {
		structFieldValue.Set(reflect.Zero(structFieldValue.Type()))
	} else {
		structFieldValue.Set(reflect.ValueOf(Convert(value, structFieldValue.Type().String())))
	}
	return nil
}

// bindVarToReflectValue sets <value> to reflect value object <structFieldValue>.
func bindVarToReflectValue(structFieldValue reflect.Value, value interface{}, recursive bool, mapping ...map[string]string) (err error) {
	kind := structFieldValue.Kind()

	// Converting using interface, for some kinds.
	switch kind {
	case reflect.Slice, reflect.Array, reflect.Ptr, reflect.Interface:
		if !structFieldValue.IsNil() {
			if v, ok := structFieldValue.Interface().(apiSet); ok {
				v.Set(value)
				return nil
			} else if v, ok := structFieldValue.Interface().(apiUnmarshalValue); ok {
				if err = v.UnmarshalValue(value); err == nil {
					return err
				}
			}
		}
	}

	// Converting by kind.
	switch kind {
	case reflect.Struct:
		// UnmarshalValue.
		if v, ok := structFieldValue.Addr().Interface().(apiUnmarshalValue); ok {
			return v.UnmarshalValue(value)
		}

		// Recursively converting for struct attribute.
		if err := doStruct(value, structFieldValue, recursive); err != nil {
			// Note there's reflect conversion mechanism here.
			structFieldValue.Set(reflect.ValueOf(value).Convert(structFieldValue.Type()))
		}

	// Note that the slice element might be type of struct,
	// so it uses Struct function doing the converting internally.
	case reflect.Slice, reflect.Array:
		a := reflect.Value{}
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
			a = reflect.MakeSlice(structFieldValue.Type(), v.Len(), v.Len())
			if v.Len() > 0 {
				t := a.Index(0).Type()
				for i := 0; i < v.Len(); i++ {
					if t.Kind() == reflect.Ptr {
						e := reflect.New(t.Elem()).Elem()
						if err := doStruct(v.Index(i).Interface(), e, recursive); err != nil {
							// Note there's reflect conversion mechanism here.
							e.Set(reflect.ValueOf(v.Index(i).Interface()).Convert(t))
						}
						a.Index(i).Set(e.Addr())
					} else {
						e := reflect.New(t).Elem()
						if err := doStruct(v.Index(i).Interface(), e, recursive); err != nil {
							// Note there's reflect conversion mechanism here.
							e.Set(reflect.ValueOf(v.Index(i).Interface()).Convert(t))
						}
						a.Index(i).Set(e)
					}
				}
			}
		} else {
			a = reflect.MakeSlice(structFieldValue.Type(), 1, 1)
			t := a.Index(0).Type()
			if t.Kind() == reflect.Ptr {
				e := reflect.New(t.Elem()).Elem()
				if err := doStruct(value, e, recursive); err != nil {
					// Note there's reflect conversion mechanism here.
					e.Set(reflect.ValueOf(value).Convert(t))
				}
				a.Index(0).Set(e.Addr())
			} else {
				e := reflect.New(t).Elem()
				if err := doStruct(value, e, recursive); err != nil {
					// Note there's reflect conversion mechanism here.
					e.Set(reflect.ValueOf(value).Convert(t))
				}
				a.Index(0).Set(e)
			}
		}
		structFieldValue.Set(a)

	case reflect.Ptr:
		item := reflect.New(structFieldValue.Type().Elem())
		// Assign value with interface Set.
		// Note that only pointer can implement interface Set.
		if v, ok := item.Interface().(apiUnmarshalValue); ok {
			err = v.UnmarshalValue(value)
			structFieldValue.Set(item)
			return err
		}
		elem := item.Elem()
		if err = bindVarToReflectValue(elem, value, recursive, mapping...); err == nil {
			structFieldValue.Set(elem.Addr())
		}

	// It mainly and specially handles the interface of nil value.
	case reflect.Interface:
		if value == nil {
			// Specially.
			structFieldValue.Set(reflect.ValueOf((*interface{})(nil)))
		} else {
			// Note there's reflect conversion mechanism here.
			structFieldValue.Set(reflect.ValueOf(value).Convert(structFieldValue.Type()))
		}

	default:
		defer func() {
			if e := recover(); e != nil {
				err = gerror.New(
					fmt.Sprintf(`cannot convert value "%+v" to type "%s"`,
						value,
						structFieldValue.Type().String(),
					),
				)
			}
		}()
		// It here uses reflect converting <value> to type of the attribute and assigns
		// the result value to the attribute. It might fail and panic if the usual Go
		// conversion rules do not allow conversion.
		structFieldValue.Set(reflect.ValueOf(value).Convert(structFieldValue.Type()))
	}
	return nil
}
