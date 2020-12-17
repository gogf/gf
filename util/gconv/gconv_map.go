// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/json"
	"reflect"
	"strings"

	"github.com/gogf/gf/internal/empty"
	"github.com/gogf/gf/internal/utils"
)

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
//
// TODO completely implement the recursive converting for all types, especially the map.
func doMapConvert(value interface{}, recursive bool, tags ...string) map[string]interface{} {
	if value == nil {
		return nil
	}
	newTags := StructTagPriority
	switch len(tags) {
	case 0:
		// No need handle.
	case 1:
		newTags = append(strings.Split(tags[0], ","), StructTagPriority...)
	default:
		newTags = append(tags, StructTagPriority...)
	}
	// Assert the common combination of types, and finally it uses reflection.
	dataMap := make(map[string]interface{})
	switch r := value.(type) {
	case string:
		// If it is a JSON string, automatically unmarshal it!
		if len(r) > 0 && r[0] == '{' && r[len(r)-1] == '}' {
			if err := json.Unmarshal([]byte(r), &dataMap); err != nil {
				return nil
			}
		} else {
			return nil
		}
	case []byte:
		// If it is a JSON string, automatically unmarshal it!
		if len(r) > 0 && r[0] == '{' && r[len(r)-1] == '}' {
			if err := json.Unmarshal(r, &dataMap); err != nil {
				return nil
			}
		} else {
			return nil
		}
	case map[interface{}]interface{}:
		for k, v := range r {
			dataMap[String(k)] = doMapConvertForMapOrStructValue(false, v, recursive, newTags...)
		}
	case map[interface{}]string:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	case map[interface{}]int:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	case map[interface{}]uint:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	case map[interface{}]float32:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	case map[interface{}]float64:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	case map[string]bool:
		for k, v := range r {
			dataMap[k] = v
		}
	case map[string]int:
		for k, v := range r {
			dataMap[k] = v
		}
	case map[string]uint:
		for k, v := range r {
			dataMap[k] = v
		}
	case map[string]float32:
		for k, v := range r {
			dataMap[k] = v
		}
	case map[string]float64:
		for k, v := range r {
			dataMap[k] = v
		}
	case map[string]interface{}:
		if recursive {
			// A copy of current map.
			for k, v := range r {
				dataMap[k] = doMapConvertForMapOrStructValue(false, v, recursive, newTags...)
			}
		} else {
			// It returns the map directly without any changing.
			return r
		}
	case map[int]interface{}:
		for k, v := range r {
			dataMap[String(k)] = doMapConvertForMapOrStructValue(false, v, recursive, newTags...)
		}
	case map[int]string:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	case map[uint]string:
		for k, v := range r {
			dataMap[String(k)] = v
		}

	default:
		// Not a common type, it then uses reflection for conversion.
		var reflectValue reflect.Value
		if v, ok := value.(reflect.Value); ok {
			reflectValue = v
		} else {
			reflectValue = reflect.ValueOf(value)
		}
		reflectKind := reflectValue.Kind()
		// If it is a pointer, we should find its real data type.
		for reflectKind == reflect.Ptr {
			reflectValue = reflectValue.Elem()
			reflectKind = reflectValue.Kind()
		}
		switch reflectKind {
		// If <value> is type of array, it converts the value of even number index as its key and
		// the value of odd number index as its corresponding value, for example:
		// []string{"k1","v1","k2","v2"} => map[string]interface{}{"k1":"v1", "k2":"v2"}
		// []string{"k1","v1","k2"}      => map[string]interface{}{"k1":"v1", "k2":nil}
		case reflect.Slice, reflect.Array:
			length := reflectValue.Len()
			for i := 0; i < length; i += 2 {
				if i+1 < length {
					dataMap[String(reflectValue.Index(i).Interface())] = reflectValue.Index(i + 1).Interface()
				} else {
					dataMap[String(reflectValue.Index(i).Interface())] = nil
				}
			}
		case reflect.Map, reflect.Struct:
			convertedValue := doMapConvertForMapOrStructValue(true, value, recursive, newTags...)
			if m, ok := convertedValue.(map[string]interface{}); ok {
				return m
			}
			return nil
		default:
			return nil
		}
	}
	return dataMap
}

func doMapConvertForMapOrStructValue(isRoot bool, value interface{}, recursive bool, tags ...string) interface{} {
	if isRoot == false && recursive == false {
		return value
	}
	var reflectValue reflect.Value
	if v, ok := value.(reflect.Value); ok {
		reflectValue = v
		value = v.Interface()
	} else {
		reflectValue = reflect.ValueOf(value)
	}
	reflectKind := reflectValue.Kind()
	// If it is a pointer, we should find its real data type.
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Map:
		var (
			mapKeys = reflectValue.MapKeys()
			dataMap = make(map[string]interface{})
		)
		for _, k := range mapKeys {
			dataMap[String(k.Interface())] = doMapConvertForMapOrStructValue(
				false,
				reflectValue.MapIndex(k).Interface(),
				recursive,
				tags...,
			)
		}
		if len(dataMap) == 0 {
			return value
		}
		return dataMap
	case reflect.Struct:
		// Map converting interface check.
		if v, ok := value.(apiMapStrAny); ok {
			m := v.MapStrAny()
			if recursive {
				for k, v := range m {
					m[k] = doMapConvertForMapOrStructValue(false, v, recursive, tags...)
				}
			}
			return m
		}
		// Using reflect for converting.
		var (
			rtField     reflect.StructField
			rvField     reflect.Value
			dataMap     = make(map[string]interface{}) // result map.
			reflectType = reflectValue.Type()          // attribute value type.
			name        = ""                           // name may be the tag name or the struct attribute name.
		)
		for i := 0; i < reflectValue.NumField(); i++ {
			rtField = reflectType.Field(i)
			rvField = reflectValue.Field(i)
			// Only convert the public attributes.
			fieldName := rtField.Name
			if !utils.IsLetterUpper(fieldName[0]) {
				continue
			}
			name = ""
			fieldTag := rtField.Tag
			for _, tag := range tags {
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
			if recursive || rtField.Anonymous {
				// Do map converting recursively.
				var (
					rvAttrField = rvField
					rvAttrKind  = rvField.Kind()
				)
				if rvAttrKind == reflect.Ptr {
					rvAttrField = rvField.Elem()
					rvAttrKind = rvAttrField.Kind()
				}
				switch rvAttrKind {
				case reflect.Struct:
					var (
						hasNoTag        = name == fieldName
						rvAttrInterface = rvAttrField.Interface()
					)
					if hasNoTag && rtField.Anonymous {
						// It means this attribute field has no tag.
						// Overwrite the attribute with sub-struct attribute fields.
						anonymousValue := doMapConvertForMapOrStructValue(false, rvAttrInterface, true, tags...)
						if m, ok := anonymousValue.(map[string]interface{}); ok {
							for k, v := range m {
								dataMap[k] = v
							}
						} else {
							dataMap[name] = rvAttrInterface
						}
					} else if !hasNoTag && rtField.Anonymous {
						// It means this attribute field has desired tag.
						dataMap[name] = doMapConvertForMapOrStructValue(false, rvAttrInterface, true, tags...)
					} else {
						dataMap[name] = doMapConvertForMapOrStructValue(false, rvAttrInterface, false, tags...)
					}

				// The struct attribute is type of slice.
				case reflect.Array, reflect.Slice:
					length := rvField.Len()
					if length == 0 {
						dataMap[name] = rvField.Interface()
						break
					}
					array := make([]interface{}, length)
					for i := 0; i < length; i++ {
						array[i] = doMapConvertForMapOrStructValue(false, rvField.Index(i), recursive, tags...)
					}
					dataMap[name] = array

				default:
					if rvField.IsValid() {
						dataMap[name] = reflectValue.Field(i).Interface()
					} else {
						dataMap[name] = nil
					}
				}
			} else {
				// No recursive map value converting
				if rvField.IsValid() {
					dataMap[name] = reflectValue.Field(i).Interface()
				} else {
					dataMap[name] = nil
				}
			}
		}
		if len(dataMap) == 0 {
			return value
		}
		return dataMap

	// The given value is type of slice.
	case reflect.Array, reflect.Slice:
		length := reflectValue.Len()
		if length == 0 {
			break
		}
		array := make([]interface{}, reflectValue.Len())
		for i := 0; i < length; i++ {
			array[i] = doMapConvertForMapOrStructValue(false, reflectValue.Index(i), recursive, tags...)
		}
		return array
	}
	return value
}

// MapStrStr converts <value> to map[string]string.
// Note that there might be data copy for this map type converting.
func MapStrStr(value interface{}, tags ...string) map[string]string {
	if r, ok := value.(map[string]string); ok {
		return r
	}
	m := Map(value, tags...)
	if len(m) > 0 {
		vMap := make(map[string]string, len(m))
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
		vMap := make(map[string]string, len(m))
		for k, v := range m {
			vMap[k] = String(v)
		}
		return vMap
	}
	return nil
}

// MapToMap converts any map type variable <params> to another map type variable <pointer>
// using reflect.
// See doMapToMap.
func MapToMap(params interface{}, pointer interface{}, mapping ...map[string]string) error {
	return doMapToMap(params, pointer, mapping...)
}

// MapToMapDeep converts any map type variable <params> to another map type variable <pointer>
// using reflect recursively.
// Deprecated, use MapToMap instead.
func MapToMapDeep(params interface{}, pointer interface{}, mapping ...map[string]string) error {
	return doMapToMap(params, pointer, mapping...)
}

// doMapToMap converts any map type variable <params> to another map type variable <pointer>.
//
// The parameter <params> can be any type of map, like:
// map[string]string, map[string]struct, , map[string]*struct, etc.
//
// The parameter <pointer> should be type of *map, like:
// map[int]string, map[string]struct, , map[string]*struct, etc.
//
// The optional parameter <mapping> is used for struct attribute to map key mapping, which makes
// sense only if the items of original map <params> is type struct.
func doMapToMap(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	var (
		paramsRv   = reflect.ValueOf(params)
		paramsKind = paramsRv.Kind()
	)
	if paramsKind == reflect.Ptr {
		paramsRv = paramsRv.Elem()
		paramsKind = paramsRv.Kind()
	}
	if paramsKind != reflect.Map {
		return gerror.New("params should be type of map")
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
		return gerror.New("pointer should be type of *map")
	}
	defer func() {
		// Catch the panic, especially the reflect operation panics.
		if e := recover(); e != nil {
			err = gerror.NewSkipf(1, "%v", e)
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
			if err = Struct(paramsRv.MapIndex(key).Interface(), e, mapping...); err != nil {
				return err
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

// MapToMaps converts any map type variable <params> to another map type variable <pointer>.
// See doMapToMaps.
func MapToMaps(params interface{}, pointer interface{}, mapping ...map[string]string) error {
	return doMapToMaps(params, pointer, mapping...)
}

// MapToMapsDeep converts any map type variable <params> to another map type variable
// <pointer> recursively.
// Deprecated, use MapToMaps instead.
func MapToMapsDeep(params interface{}, pointer interface{}, mapping ...map[string]string) error {
	return doMapToMaps(params, pointer, mapping...)
}

// doMapToMaps converts any map type variable <params> to another map type variable <pointer>.
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
func doMapToMaps(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	var (
		paramsRv   = reflect.ValueOf(params)
		paramsKind = paramsRv.Kind()
	)
	if paramsKind == reflect.Ptr {
		paramsRv = paramsRv.Elem()
		paramsKind = paramsRv.Kind()
	}
	if paramsKind != reflect.Map {
		return gerror.New("params should be type of map")
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
		return gerror.New("pointer should be type of *map/**map")
	}
	defer func() {
		// Catch the panic, especially the reflect operation panics.
		if e := recover(); e != nil {
			err = gerror.NewSkipf(1, "%v", e)
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
		if err = Structs(paramsRv.MapIndex(key).Interface(), e.Addr(), mapping...); err != nil {
			return err
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
