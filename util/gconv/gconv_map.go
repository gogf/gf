// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gogf/gf/internal/empty"
	"github.com/gogf/gf/internal/utilstr"
)

// Interface support for package gmap.
type apiMapStrAny interface {
	MapStrAny() map[string]interface{}
}

// Map converts any variable <value> to map[string]interface{}.
//
// If the parameter <value> is not a map/struct/*struct type, then the conversion will fail and returns nil.
//
// If <value> is a struct/*struct object, the second parameter <tags> specifies the most priority
// tags that will be detected, otherwise it detects the tags in order of: gconv, json, and then the field name.
func Map(value interface{}, tags ...string) map[string]interface{} {
	if value == nil {
		return nil
	}
	if r, ok := value.(map[string]interface{}); ok {
		return r
	} else {
		// Only assert the common combination of types, and finally it uses reflection.
		m := make(map[string]interface{})
		switch value.(type) {
		case map[interface{}]interface{}:
			for k, v := range value.(map[interface{}]interface{}) {
				m[String(k)] = v
			}
		case map[interface{}]string:
			for k, v := range value.(map[interface{}]string) {
				m[String(k)] = v
			}
		case map[interface{}]int:
			for k, v := range value.(map[interface{}]int) {
				m[String(k)] = v
			}
		case map[interface{}]uint:
			for k, v := range value.(map[interface{}]uint) {
				m[String(k)] = v
			}
		case map[interface{}]float32:
			for k, v := range value.(map[interface{}]float32) {
				m[String(k)] = v
			}
		case map[interface{}]float64:
			for k, v := range value.(map[interface{}]float64) {
				m[String(k)] = v
			}
		case map[string]bool:
			for k, v := range value.(map[string]bool) {
				m[k] = v
			}
		case map[string]int:
			for k, v := range value.(map[string]int) {
				m[k] = v
			}
		case map[string]uint:
			for k, v := range value.(map[string]uint) {
				m[k] = v
			}
		case map[string]float32:
			for k, v := range value.(map[string]float32) {
				m[k] = v
			}
		case map[string]float64:
			for k, v := range value.(map[string]float64) {
				m[k] = v
			}
		case map[int]interface{}:
			for k, v := range value.(map[int]interface{}) {
				m[String(k)] = v
			}
		case map[int]string:
			for k, v := range value.(map[int]string) {
				m[String(k)] = v
			}
		case map[uint]string:
			for k, v := range value.(map[uint]string) {
				m[String(k)] = v
			}
		// Not a common type, use reflection
		default:
			rv := reflect.ValueOf(value)
			kind := rv.Kind()
			// If it is a pointer, we should find its real data type.
			if kind == reflect.Ptr {
				rv = rv.Elem()
				kind = rv.Kind()
			}
			switch kind {
			case reflect.Map:
				ks := rv.MapKeys()
				for _, k := range ks {
					m[String(k.Interface())] = rv.MapIndex(k).Interface()
				}
			case reflect.Struct:
				if v, ok := value.(apiMapStrAny); ok {
					return v.MapStrAny()
				}
				rt := rv.Type()
				name := ""
				tagArray := structTagPriority
				switch len(tags) {
				case 0:
					// No need handle.
				case 1:
					tagArray = append(strings.Split(tags[0], ","), structTagPriority...)
				default:
					tagArray = append(tags, structTagPriority...)
				}
				for i := 0; i < rv.NumField(); i++ {
					// Only convert the public attributes.
					fieldName := rt.Field(i).Name
					if !utilstr.IsLetterUpper(fieldName[0]) {
						continue
					}
					name = ""
					fieldTag := rt.Field(i).Tag
					for _, tag := range tagArray {
						if name = fieldTag.Get(tag); name != "" {
							break
						}
					}
					if name == "" {
						name = strings.TrimSpace(fieldName)
					} else {
						// Support json tag feature: -, omitempty
						name = strings.TrimSpace(name)
						if name == "-" {
							continue
						}
						array := strings.Split(name, ",")
						if len(array) > 1 {
							switch strings.TrimSpace(array[1]) {
							case "omitempty":
								if empty.IsEmpty(rv.Field(i).Interface()) {
									continue
								} else {
									name = strings.TrimSpace(array[0])
								}
							default:
								name = strings.TrimSpace(array[0])
							}
						}
					}
					m[name] = rv.Field(i).Interface()
				}
			default:
				return nil
			}
		}
		return m
	}
}

// MapDeep do Map function recursively.
// See Map.
func MapDeep(value interface{}, tags ...string) map[string]interface{} {
	data := Map(value, tags...)
	for key, value := range data {
		rv := reflect.ValueOf(value)
		kind := rv.Kind()
		if kind == reflect.Ptr {
			rv = rv.Elem()
			kind = rv.Kind()
		}
		switch kind {
		case reflect.Struct:
			delete(data, key)
			for k, v := range MapDeep(value, tags...) {
				data[k] = v
			}
		}
	}
	return data
}

// MapToMap converts map type variable <params> to another map type variable <pointer>.
// The elements of <pointer> should be type of *map.
func MapToMap(params interface{}, pointer interface{}, mapping ...map[string]string) error {
	return doMapToMap(params, pointer, false, mapping...)
}

// MapToMapDeep recursively converts map type variable <params> to another map type variable <pointer>.
// The elements of <pointer> should be type of *map.
func MapToMapDeep(params interface{}, pointer interface{}, mapping ...map[string]string) error {
	return doMapToMap(params, pointer, true, mapping...)
}

// doMapStruct converts map type variable <params> to another map type variable <pointer>.
// The elements of <pointer> should be type of *map.
func doMapToMap(params interface{}, pointer interface{}, deep bool, mapping ...map[string]string) error {
	paramsRv := reflect.ValueOf(params)
	paramsKind := paramsRv.Kind()
	if paramsKind == reflect.Ptr {
		paramsRv = paramsRv.Elem()
		paramsKind = paramsRv.Kind()
	}
	if paramsKind != reflect.Map {
		return errors.New("params should be type of map")
	}

	pointerRv := reflect.ValueOf(pointer)
	pointerKind := pointerRv.Kind()
	for pointerKind == reflect.Ptr {
		pointerRv = pointerRv.Elem()
		pointerKind = pointerRv.Kind()
	}
	if pointerKind != reflect.Map {
		return errors.New("pointer should be type of map")
	}
	err := (error)(nil)
	paramsKeys := paramsRv.MapKeys()
	pointerKeyType := pointerRv.Type().Key()
	pointerValueType := pointerRv.Type().Elem()
	dataMap := reflect.MakeMapWithSize(pointerRv.Type(), len(paramsKeys))
	for _, key := range paramsKeys {
		e := reflect.New(pointerValueType).Elem()
		if deep {
			if err = StructDeep(paramsRv.MapIndex(key).Interface(), e, mapping...); err != nil {
				return err
			}
		} else {
			if err = Struct(paramsRv.MapIndex(key).Interface(), e, mapping...); err != nil {
				return err
			}
		}
		dataMap.SetMapIndex(
			reflect.ValueOf(Convert(key.Interface(), pointerKeyType.Name())),
			e,
		)
	}
	pointerRv.Set(dataMap)
	return nil
}

// MapToMaps converts map type variable <params> to another map type variable <pointer>.
// The elements of <pointer> should be type of []map/*map.
func MapToMaps(params interface{}, pointer interface{}, mapping ...map[string]string) error {
	return doMapToMaps(params, pointer, false, mapping...)
}

// MapToMapsDeep recursively converts map type variable <params> to another map type variable <pointer>.
// The elements of <pointer> should be type of []map/*map.
func MapToMapsDeep(params interface{}, pointer interface{}, mapping ...map[string]string) error {
	return doMapToMaps(params, pointer, true, mapping...)
}

// doMapToMaps converts map type variable <params> to another map type variable <pointer>.
// The elements of <pointer> should be type of []map/*map.
func doMapToMaps(params interface{}, pointer interface{}, deep bool, mapping ...map[string]string) error {
	paramsRv := reflect.ValueOf(params)
	paramsKind := paramsRv.Kind()
	if paramsKind == reflect.Ptr {
		paramsRv = paramsRv.Elem()
		paramsKind = paramsRv.Kind()
	}
	if paramsKind != reflect.Map {
		return errors.New("params should be type of map")
	}

	pointerRv := reflect.ValueOf(pointer)
	pointerKind := pointerRv.Kind()
	for pointerKind == reflect.Ptr {
		pointerRv = pointerRv.Elem()
		pointerKind = pointerRv.Kind()
	}
	if pointerKind != reflect.Map {
		return errors.New("pointer should be type of map")
	}
	err := (error)(nil)
	paramsKeys := paramsRv.MapKeys()
	pointerKeyType := pointerRv.Type().Key()
	pointerValueType := pointerRv.Type().Elem()
	dataMap := reflect.MakeMapWithSize(pointerRv.Type(), len(paramsKeys))
	for _, key := range paramsKeys {
		e := reflect.New(pointerValueType).Elem().Addr()
		if deep {
			if err = StructsDeep(paramsRv.MapIndex(key).Interface(), e, mapping...); err != nil {
				return err
			}
		} else {
			if err = Structs(paramsRv.MapIndex(key).Interface(), e, mapping...); err != nil {
				return err
			}
		}
		dataMap.SetMapIndex(
			reflect.ValueOf(Convert(key.Interface(), pointerKeyType.Name())),
			e.Elem(),
		)
	}
	pointerRv.Set(dataMap)
	return nil
}
