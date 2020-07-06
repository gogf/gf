// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/empty"
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
// 2. The second parameter <pointer> should be a pointer to the struct object.
// 3. Only the public attributes of struct object can be mapped.
// 4. If <params> is a map, the key of the map <params> can be lowercase.
//    It will automatically convert the first letter of the key to uppercase
//    in mapping procedure to do the matching.
//    It ignores the map key, if it does not match.
func Struct(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	if params == nil {
		return errors.New("params cannot be nil")
	}
	if pointer == nil {
		return errors.New("object pointer cannot be nil")
	}
	defer func() {
		// Catch the panic, especially the reflect operation panics.
		if e := recover(); e != nil {
			err = gerror.NewfSkip(1, "%v", e)
		}
	}()

	// paramsMap is the map[string]interface{} type variable for params.
	paramsMap := Map(params)
	if paramsMap == nil {
		return fmt.Errorf("invalid params: %v", params)
	}

	// UnmarshalValue.
	// Assign value with interface UnmarshalValue.
	// Note that only pointer can implement interface UnmarshalValue.
	if v, ok := pointer.(apiUnmarshalValue); ok {
		return v.UnmarshalValue(params)
	}

	// Using reflect to do the converting,
	// it also supports type of reflect.Value for <pointer>(always in internal usage).
	elem, ok := pointer.(reflect.Value)
	if !ok {
		rv := reflect.ValueOf(pointer)
		if kind := rv.Kind(); kind != reflect.Ptr {
			return fmt.Errorf("object pointer should be type of '*struct', but got '%v'", kind)
		}
		// Using IsNil on reflect.Ptr variable is OK.
		if !rv.IsValid() || rv.IsNil() {
			return errors.New("object pointer cannot be nil")
		}
		elem = rv.Elem()
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

	// UnmarshalValue.
	// Assign value with interface UnmarshalValue.
	// Note that only pointer can implement interface UnmarshalValue.
	if elem.Kind() == reflect.Struct {
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
			if v, ok := paramsMap[mapK]; ok {
				doneMap[mapV] = struct{}{}
				if err := bindVarToStructAttr(elem, mapV, v); err != nil {
					return err
				}
			}
		}
	}

	// The key of the attrMap is the attribute name of the struct,
	// and the value is its replaced name for later comparison to improve performance.
	var (
		attrMap  = make(map[string]string)
		elemType = elem.Type()
		tempName = ""
	)
	for i := 0; i < elem.NumField(); i++ {
		// Only do converting to public attributes.
		if !utils.IsLetterUpper(elemType.Field(i).Name[0]) {
			continue
		}
		tempName = elemType.Field(i).Name
		attrMap[tempName] = replaceCharReg.ReplaceAllString(tempName, "")
	}
	if len(attrMap) == 0 {
		return nil
	}

	// The key of the tagMap is the attribute name of the struct,
	// and the value is its replaced tag name for later comparison to improve performance.
	tagMap := make(map[string]string)
	for k, v := range structs.TagMapName(pointer, structTagPriority, true) {
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
		if err := bindVarToStructAttr(elem, attrName, mapV); err != nil {
			return err
		}
	}
	return nil
}

// StructDeep do Struct function recursively.
// See Struct.
func StructDeep(params interface{}, pointer interface{}, mapping ...map[string]string) error {
	if params == nil {
		return nil
	}
	if err := Struct(params, pointer, mapping...); err != nil {
		return err
	} else {
		rv, ok := pointer.(reflect.Value)
		if !ok {
			rv = reflect.ValueOf(pointer)
		}
		kind := rv.Kind()
		for kind == reflect.Ptr {
			rv = rv.Elem()
			kind = rv.Kind()
		}
		switch kind {
		case reflect.Struct:
			rt := rv.Type()
			for i := 0; i < rv.NumField(); i++ {
				// Only do converting to public attributes.
				if !utils.IsLetterUpper(rt.Field(i).Name[0]) {
					continue
				}
				trv := rv.Field(i)
				switch trv.Kind() {
				case reflect.Struct:
					if err := StructDeep(params, trv, mapping...); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

// bindVarToStructAttr sets value to struct object attribute by name.
func bindVarToStructAttr(elem reflect.Value, name string, value interface{}) (err error) {
	structFieldValue := elem.FieldByName(name)
	if !structFieldValue.IsValid() {
		return nil
	}
	// CanSet checks whether attribute is public accessible.
	if !structFieldValue.CanSet() {
		return nil
	}
	defer func() {
		if recover() != nil {
			err = bindVarToReflectValue(structFieldValue, value)
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
func bindVarToReflectValue(structFieldValue reflect.Value, value interface{}) (err error) {
	kind := structFieldValue.Kind()

	// Converting using interface, for some kinds.
	switch kind {
	case reflect.Slice, reflect.Array, reflect.Ptr, reflect.Interface:
		if !structFieldValue.IsNil() {
			if v, ok := structFieldValue.Interface().(apiSet); ok {
				v.Set(value)
				return nil
			} else if v, ok := structFieldValue.Interface().(apiUnmarshalValue); ok {
				err = v.UnmarshalValue(value)
				if err == nil {
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

		if err := Struct(value, structFieldValue); err != nil {
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
						if err := Struct(v.Index(i).Interface(), e); err != nil {
							// Note there's reflect conversion mechanism here.
							e.Set(reflect.ValueOf(v.Index(i).Interface()).Convert(t))
						}
						a.Index(i).Set(e.Addr())
					} else {
						e := reflect.New(t).Elem()
						if err := Struct(v.Index(i).Interface(), e); err != nil {
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
				if err := Struct(value, e); err != nil {
					// Note there's reflect conversion mechanism here.
					e.Set(reflect.ValueOf(value).Convert(t))
				}
				a.Index(0).Set(e.Addr())
			} else {
				e := reflect.New(t).Elem()
				if err := Struct(value, e); err != nil {
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
		if err = bindVarToReflectValue(elem, value); err == nil {
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
				err = errors.New(
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
