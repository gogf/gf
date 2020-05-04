// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"encoding/json"
	"errors"
	"github.com/gogf/gf/errors/gerror"
	"reflect"
	"strings"

	"github.com/gogf/gf/internal/empty"
	"github.com/gogf/gf/internal/utils"
)

// apiMapStrAny is the interface support for converting struct parameter to map.
type apiMapStrAny interface {
	MapStrAny() map[string]interface{}
}

// Map converts any variable <value> to map[string]interface{}. If the parameter <value> is not a
// map/struct/*struct type, then the conversion will fail and returns nil.
//
// If <value> is a struct/*struct object, the second parameter <tags> specifies the most priority
// tags that will be detected, otherwise it detects the tags in order of:
// gconv, json, field name.
func Map(value interface{}, tags ...string) map[string]interface{} {
	return doMapConvert(value, false, tags...)
}

// MapDeep does Map function recursively, which means if the attribute of <value>
// is also a struct/*struct, calls Map function on this attribute converting it to
// a map[string]interface{} type variable.
// Also see Map.
func MapDeep(value interface{}, tags ...string) map[string]interface{} {
	return doMapConvert(value, true, tags...)
}

// doMapConvert implements the map converting.
// It automatically checks and converts json string to map if <value> is string/[]byte.
func doMapConvert(value interface{}, recursive bool, tags ...string) map[string]interface{} {
	if value == nil {
		return nil
	}

	// Assert the common combination of types, and finally it uses reflection.
	m := make(map[string]interface{})
	switch r := value.(type) {
	case string:
		if len(r) > 0 && r[0] == '{' && r[len(r)-1] == '}' {
			if err := json.Unmarshal([]byte(r), &m); err != nil {
				return nil
			}
		} else {
			return nil
		}
	case []byte:
		if len(r) > 0 && r[0] == '{' && r[len(r)-1] == '}' {
			if err := json.Unmarshal(r, &m); err != nil {
				return nil
			}
		} else {
			return nil
		}
	case map[interface{}]interface{}:
		for k, v := range r {
			m[String(k)] = v
		}
	case map[interface{}]string:
		for k, v := range r {
			m[String(k)] = v
		}
	case map[interface{}]int:
		for k, v := range r {
			m[String(k)] = v
		}
	case map[interface{}]uint:
		for k, v := range r {
			m[String(k)] = v
		}
	case map[interface{}]float32:
		for k, v := range r {
			m[String(k)] = v
		}
	case map[interface{}]float64:
		for k, v := range r {
			m[String(k)] = v
		}
	case map[string]bool:
		for k, v := range r {
			m[k] = v
		}
	case map[string]int:
		for k, v := range r {
			m[k] = v
		}
	case map[string]uint:
		for k, v := range r {
			m[k] = v
		}
	case map[string]float32:
		for k, v := range r {
			m[k] = v
		}
	case map[string]float64:
		for k, v := range r {
			m[k] = v
		}
	case map[string]interface{}:
		return r
	case map[int]interface{}:
		for k, v := range r {
			m[String(k)] = v
		}
	case map[int]string:
		for k, v := range r {
			m[String(k)] = v
		}
	case map[uint]string:
		for k, v := range r {
			m[String(k)] = v
		}
	// Not a common type, it then uses reflection for conversion.
	default:
		var rv reflect.Value
		if v, ok := value.(reflect.Value); ok {
			rv = v
		} else {
			rv = reflect.ValueOf(value)
		}
		kind := rv.Kind()
		// If it is a pointer, we should find its real data type.
		if kind == reflect.Ptr {
			rv = rv.Elem()
			kind = rv.Kind()
		}
		switch kind {
		// If <value> is type of array, it converts the value of even number index as its key and
		// the value of odd number index as its corresponding value.
		// Eg:
		// []string{"k1","v1","k2","v2"} => map[string]interface{}{"k1":"v1", "k2":"v2"}
		// []string{"k1","v1","k2"} => map[string]interface{}{"k1":"v1", "k2":nil}
		case reflect.Slice, reflect.Array:
			length := rv.Len()
			for i := 0; i < length; i += 2 {
				if i+1 < length {
					m[String(rv.Index(i).Interface())] = rv.Index(i + 1).Interface()
				} else {
					m[String(rv.Index(i).Interface())] = nil
				}
			}
		case reflect.Map:
			ks := rv.MapKeys()
			for _, k := range ks {
				m[String(k.Interface())] = rv.MapIndex(k).Interface()
			}
		case reflect.Struct:
			// Map converting interface check.
			if v, ok := value.(apiMapStrAny); ok {
				return v.MapStrAny()
			}
			// Using reflect for converting.
			var (
				rtField  reflect.StructField
				rvField  reflect.Value
				rvKind   reflect.Kind
				rt       = rv.Type()
				name     = ""
				tagArray = structTagPriority
			)
			switch len(tags) {
			case 0:
				// No need handle.
			case 1:
				tagArray = append(strings.Split(tags[0], ","), structTagPriority...)
			default:
				tagArray = append(tags, structTagPriority...)
			}
			for i := 0; i < rv.NumField(); i++ {
				rtField = rt.Field(i)
				rvField = rv.Field(i)
				// Only convert the public attributes.
				fieldName := rtField.Name
				if !utils.IsLetterUpper(fieldName[0]) {
					continue
				}
				name = ""
				fieldTag := rtField.Tag
				for _, tag := range tagArray {
					if name = fieldTag.Get(tag); name != "" {
						break
					}
				}
				if name == "" {
					name = fieldName
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
							if empty.IsEmpty(rvField.Interface()) {
								continue
							} else {
								name = strings.TrimSpace(array[0])
							}
						default:
							name = strings.TrimSpace(array[0])
						}
					}
				}
				if recursive {
					rvKind = rvField.Kind()
					if rvKind == reflect.Ptr {
						rvField = rvField.Elem()
						rvKind = rvField.Kind()
					}
					if rvKind == reflect.Struct {
						hasNoTag := name == fieldName
						if hasNoTag && rtField.Anonymous {
							// It means this attribute field has no tag.
							// Overwrite the attribute with sub-struct attribute fields.
							for k, v := range doMapConvert(rvField.Interface(), recursive, tags...) {
								m[k] = v
							}
						} else {
							// It means this attribute field has desired tag.
							m[name] = doMapConvert(rvField.Interface(), recursive, tags...)
						}

					} else {
						if rvField.IsValid() {
							m[name] = rvField.Interface()
						} else {
							m[name] = nil
						}
					}
				} else {
					if rvField.IsValid() {
						m[name] = rvField.Interface()
					} else {
						m[name] = nil
					}
				}
			}
		default:
			return nil
		}
	}
	return m
}

// MapStrStr converts <value> to map[string]string.
// Note that there might be data copy for this map type converting.
func MapStrStr(value interface{}, tags ...string) map[string]string {
	if r, ok := value.(map[string]string); ok {
		return r
	}
	m := Map(value, tags...)
	if len(m) > 0 {
		vMap := make(map[string]string)
		for k, v := range m {
			vMap[k] = String(v)
		}
		return vMap
	}
	return nil
}

// MapStrStrDeep converts <value> to map[string]string recursively.
// Note that there might be data copy for this map type converting.
func MapStrStrDeep(value interface{}, tags ...string) map[string]string {
	if r, ok := value.(map[string]string); ok {
		return r
	}
	m := MapDeep(value, tags...)
	if len(m) > 0 {
		vMap := make(map[string]string)
		for k, v := range m {
			vMap[k] = String(v)
		}
		return vMap
	}
	return nil
}

// MapToMap converts map type variable <params> to another map type variable <pointer> using
// reflect.
// See doMapToMap.
func MapToMap(params interface{}, pointer interface{}, mapping ...map[string]string) error {
	return doMapToMap(params, pointer, false, mapping...)
}

// MapToMapDeep converts map type variable <params> to another map type variable <pointer> using
// reflect recursively.
// See doMapToMap.
func MapToMapDeep(params interface{}, pointer interface{}, mapping ...map[string]string) error {
	return doMapToMap(params, pointer, true, mapping...)
}

// doMapToMap converts map type variable <params> to another map type variable <pointer>.
//
// The parameter <params> can be any type of map, like:
// map[string]string, map[string]struct, , map[string]*struct, etc.
//
// The parameter <pointer> should be type of *map, like:
// map[int]string, map[string]struct, , map[string]*struct, etc.
//
// The optional parameter <mapping> is used for struct attribute to map key mapping, which makes
// sense only if the items of original map <params> is type struct.
func doMapToMap(params interface{}, pointer interface{}, deep bool, mapping ...map[string]string) (err error) {
	var (
		paramsRv   = reflect.ValueOf(params)
		paramsKind = paramsRv.Kind()
	)
	if paramsKind == reflect.Ptr {
		paramsRv = paramsRv.Elem()
		paramsKind = paramsRv.Kind()
	}
	if paramsKind != reflect.Map {
		return errors.New("params should be type of map")
	}
	// Empty params map, no need continue.
	if paramsRv.Len() == 0 {
		return nil
	}
	var pointerRv reflect.Value
	if v, ok := pointer.(reflect.Value); ok {
		pointerRv = v
	} else {
		pointerRv = reflect.ValueOf(pointer)
	}
	pointerKind := pointerRv.Kind()
	for pointerKind == reflect.Ptr {
		pointerRv = pointerRv.Elem()
		pointerKind = pointerRv.Kind()
	}
	if pointerKind != reflect.Map {
		return errors.New("pointer should be type of *map")
	}
	defer func() {
		// Catch the panic, especially the reflect operation panics.
		if e := recover(); e != nil {
			err = gerror.NewfSkip(1, "%v", e)
		}
	}()
	var (
		paramsKeys       = paramsRv.MapKeys()
		pointerKeyType   = pointerRv.Type().Key()
		pointerValueType = pointerRv.Type().Elem()
		pointerValueKind = pointerValueType.Kind()
		dataMap          = reflect.MakeMapWithSize(pointerRv.Type(), len(paramsKeys))
	)
	// Retrieve the true element type of target map.
	if pointerValueKind == reflect.Ptr {
		pointerValueKind = pointerValueType.Elem().Kind()
	}
	for _, key := range paramsKeys {
		e := reflect.New(pointerValueType).Elem()
		switch pointerValueKind {
		case reflect.Map, reflect.Struct:
			if deep {
				if err = StructDeep(paramsRv.MapIndex(key).Interface(), e, mapping...); err != nil {
					return err
				}
			} else {
				if err = Struct(paramsRv.MapIndex(key).Interface(), e, mapping...); err != nil {
					return err
				}
			}
		default:
			e.Set(
				reflect.ValueOf(
					Convert(
						paramsRv.MapIndex(key).Interface(),
						pointerValueType.String(),
					),
				),
			)
		}
		dataMap.SetMapIndex(
			reflect.ValueOf(
				Convert(
					key.Interface(),
					pointerKeyType.Name(),
				),
			),
			e,
		)
	}
	pointerRv.Set(dataMap)
	return nil
}

// MapToMaps converts map type variable <params> to another map type variable <pointer>.
// See doMapToMaps.
func MapToMaps(params interface{}, pointer interface{}, mapping ...map[string]string) error {
	return doMapToMaps(params, pointer, false, mapping...)
}

// MapToMapsDeep converts map type variable <params> to another map type variable
// <pointer> recursively.
// See doMapToMaps.
func MapToMapsDeep(params interface{}, pointer interface{}, mapping ...map[string]string) error {
	return doMapToMaps(params, pointer, true, mapping...)
}

// doMapToMaps converts map type variable <params> to another map type variable <pointer>.
//
// The parameter <params> can be any type of map, of which the item type is slice map, like:
// map[int][]map, map[string][]map.
//
// The parameter <pointer> should be type of *map, of which the item type is slice map, like:
// map[string][]struct, map[string][]*struct.
//
// The optional parameter <mapping> is used for struct attribute to map key mapping, which makes
// sense only if the items of original map is type struct.
//
// TODO it's supposed supporting target type <pointer> like: map[int][]map, map[string][]map.
func doMapToMaps(params interface{}, pointer interface{}, deep bool, mapping ...map[string]string) (err error) {
	var (
		paramsRv   = reflect.ValueOf(params)
		paramsKind = paramsRv.Kind()
	)
	if paramsKind == reflect.Ptr {
		paramsRv = paramsRv.Elem()
		paramsKind = paramsRv.Kind()
	}
	if paramsKind != reflect.Map {
		return errors.New("params should be type of map")
	}
	// Empty params map, no need continue.
	if paramsRv.Len() == 0 {
		return nil
	}
	var (
		pointerRv   = reflect.ValueOf(pointer)
		pointerKind = pointerRv.Kind()
	)
	for pointerKind == reflect.Ptr {
		pointerRv = pointerRv.Elem()
		pointerKind = pointerRv.Kind()
	}
	if pointerKind != reflect.Map {
		return errors.New("pointer should be type of *map/**map")
	}
	defer func() {
		// Catch the panic, especially the reflect operation panics.
		if e := recover(); e != nil {
			err = gerror.NewfSkip(1, "%v", e)
		}
	}()
	var (
		paramsKeys       = paramsRv.MapKeys()
		pointerKeyType   = pointerRv.Type().Key()
		pointerValueType = pointerRv.Type().Elem()
		dataMap          = reflect.MakeMapWithSize(pointerRv.Type(), len(paramsKeys))
	)
	for _, key := range paramsKeys {
		e := reflect.New(pointerValueType).Elem()
		if deep {
			if err = StructsDeep(paramsRv.MapIndex(key).Interface(), e.Addr(), mapping...); err != nil {
				return err
			}
		} else {
			if err = Structs(paramsRv.MapIndex(key).Interface(), e.Addr(), mapping...); err != nil {
				return err
			}
		}
		dataMap.SetMapIndex(
			reflect.ValueOf(
				Convert(
					key.Interface(),
					pointerKeyType.Name(),
				),
			),
			e,
		)
	}
	pointerRv.Set(dataMap)
	return nil
}
