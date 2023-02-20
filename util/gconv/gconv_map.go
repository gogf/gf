// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"reflect"
	"strings"

	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/utils"
)

type recursiveType string

const (
	recursiveTypeAuto recursiveType = "auto"
	recursiveTypeTrue recursiveType = "true"
)

// Map converts any variable `value` to map[string]interface{}. If the parameter `value` is not a
// map/struct/*struct type, then the conversion will fail and returns nil.
//
// If `value` is a struct/*struct object, the second parameter `tags` specifies the most priority
// tags that will be detected, otherwise it detects the tags in order of:
// gconv, json, field name.
func Map(value interface{}, tags ...string) map[string]interface{} {
	return doMapConvert(value, recursiveTypeAuto, tags...)
}

// MapDeep does Map function recursively, which means if the attribute of `value`
// is also a struct/*struct, calls Map function on this attribute converting it to
// a map[string]interface{} type variable.
// Also see Map.
func MapDeep(value interface{}, tags ...string) map[string]interface{} {
	return doMapConvert(value, recursiveTypeTrue, tags...)
}

// doMapConvert implements the map converting.
// It automatically checks and converts json string to map if `value` is string/[]byte.
//
// TODO completely implement the recursive converting for all types, especially the map.
func doMapConvert(value interface{}, recursive recursiveType, tags ...string) map[string]interface{} {
	if value == nil {
		return nil
	}
	newTags := StructTagPriority
	switch len(tags) {
	case 0:
		// No need handling.
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
			if err := json.UnmarshalUseNumber([]byte(r), &dataMap); err != nil {
				return nil
			}
		} else {
			return nil
		}
	case []byte:
		// If it is a JSON string, automatically unmarshal it!
		if len(r) > 0 && r[0] == '{' && r[len(r)-1] == '}' {
			if err := json.UnmarshalUseNumber(r, &dataMap); err != nil {
				return nil
			}
		} else {
			return nil
		}
	case map[interface{}]interface{}:
		for k, v := range r {
			dataMap[String(k)] = doMapConvertForMapOrStructValue(
				doMapConvertForMapOrStructValueInput{
					IsRoot:          false,
					Value:           v,
					RecursiveType:   recursive,
					RecursiveOption: recursive == recursiveTypeTrue,
					Tags:            newTags,
				},
			)
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
	case map[string]string:
		for k, v := range r {
			dataMap[k] = v
		}
	case map[string]interface{}:
		if recursive == recursiveTypeTrue {
			// A copy of current map.
			for k, v := range r {
				dataMap[k] = doMapConvertForMapOrStructValue(
					doMapConvertForMapOrStructValueInput{
						IsRoot:          false,
						Value:           v,
						RecursiveType:   recursive,
						RecursiveOption: recursive == recursiveTypeTrue,
						Tags:            newTags,
					},
				)
			}
		} else {
			// It returns the map directly without any changing.
			return r
		}
	case map[int]interface{}:
		for k, v := range r {
			dataMap[String(k)] = doMapConvertForMapOrStructValue(
				doMapConvertForMapOrStructValueInput{
					IsRoot:          false,
					Value:           v,
					RecursiveType:   recursive,
					RecursiveOption: recursive == recursiveTypeTrue,
					Tags:            newTags,
				},
			)
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
		// If `value` is type of array, it converts the value of even number index as its key and
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
		case reflect.Map, reflect.Struct, reflect.Interface:
			convertedValue := doMapConvertForMapOrStructValue(
				doMapConvertForMapOrStructValueInput{
					IsRoot:          true,
					Value:           value,
					RecursiveType:   recursive,
					RecursiveOption: recursive == recursiveTypeTrue,
					Tags:            newTags,
				},
			)
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

type doMapConvertForMapOrStructValueInput struct {
	IsRoot          bool          // It returns directly if it is not root and with no recursive converting.
	Value           interface{}   // Current operation value.
	RecursiveType   recursiveType // The type from top function entry.
	RecursiveOption bool          // Whether convert recursively for `current` operation.
	Tags            []string      // Map key mapping.
}

func doMapConvertForMapOrStructValue(in doMapConvertForMapOrStructValueInput) interface{} {
	if in.IsRoot == false && in.RecursiveOption == false {
		return in.Value
	}

	var reflectValue reflect.Value
	if v, ok := in.Value.(reflect.Value); ok {
		reflectValue = v
		in.Value = v.Interface()
	} else {
		reflectValue = reflect.ValueOf(in.Value)
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
				doMapConvertForMapOrStructValueInput{
					IsRoot:          false,
					Value:           reflectValue.MapIndex(k).Interface(),
					RecursiveType:   in.RecursiveType,
					RecursiveOption: in.RecursiveType == recursiveTypeTrue,
					Tags:            in.Tags,
				},
			)
		}
		return dataMap

	case reflect.Struct:
		var dataMap = make(map[string]interface{})
		// Map converting interface check.
		if v, ok := in.Value.(iMapStrAny); ok {
			// Value copy, in case of concurrent safety.
			for mapK, mapV := range v.MapStrAny() {
				if in.RecursiveOption {
					dataMap[mapK] = doMapConvertForMapOrStructValue(
						doMapConvertForMapOrStructValueInput{
							IsRoot:          false,
							Value:           mapV,
							RecursiveType:   in.RecursiveType,
							RecursiveOption: in.RecursiveType == recursiveTypeTrue,
							Tags:            in.Tags,
						},
					)
				} else {
					dataMap[mapK] = mapV
				}
			}
			return dataMap
		}
		// Using reflect for converting.
		var (
			rtField     reflect.StructField
			rvField     reflect.Value
			reflectType = reflectValue.Type() // attribute value type.
			mapKey      = ""                  // mapKey may be the tag name or the struct attribute name.
		)
		for i := 0; i < reflectValue.NumField(); i++ {
			rtField = reflectType.Field(i)
			rvField = reflectValue.Field(i)
			// Only convert the public attributes.
			fieldName := rtField.Name
			if !utils.IsLetterUpper(fieldName[0]) {
				continue
			}
			mapKey = ""
			fieldTag := rtField.Tag
			for _, tag := range in.Tags {
				if mapKey = fieldTag.Get(tag); mapKey != "" {
					break
				}
			}
			if mapKey == "" {
				mapKey = fieldName
			} else {
				// Support json tag feature: -, omitempty
				mapKey = strings.TrimSpace(mapKey)
				if mapKey == "-" {
					continue
				}
				array := strings.Split(mapKey, ",")
				if len(array) > 1 {
					switch strings.TrimSpace(array[1]) {
					case "omitempty":
						if empty.IsEmpty(rvField.Interface()) {
							continue
						} else {
							mapKey = strings.TrimSpace(array[0])
						}
					default:
						mapKey = strings.TrimSpace(array[0])
					}
				}
			}
			if in.RecursiveOption || rtField.Anonymous {
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
					// Embedded struct and has no fields, just ignores it.
					// Eg: gmeta.Meta
					if rvAttrField.Type().NumField() == 0 {
						continue
					}
					var (
						hasNoTag = mapKey == fieldName
						// DO NOT use rvAttrField.Interface() here,
						// as it might be changed from pointer to struct.
						rvInterface = rvField.Interface()
					)
					switch {
					case hasNoTag && rtField.Anonymous:
						// It means this attribute field has no tag.
						// Overwrite the attribute with sub-struct attribute fields.
						anonymousValue := doMapConvertForMapOrStructValue(doMapConvertForMapOrStructValueInput{
							IsRoot:          false,
							Value:           rvInterface,
							RecursiveType:   in.RecursiveType,
							RecursiveOption: true,
							Tags:            in.Tags,
						})
						if m, ok := anonymousValue.(map[string]interface{}); ok {
							for k, v := range m {
								dataMap[k] = v
							}
						} else {
							dataMap[mapKey] = rvInterface
						}

					// It means this attribute field has desired tag.
					case !hasNoTag && rtField.Anonymous:
						dataMap[mapKey] = doMapConvertForMapOrStructValue(doMapConvertForMapOrStructValueInput{
							IsRoot:          false,
							Value:           rvInterface,
							RecursiveType:   in.RecursiveType,
							RecursiveOption: true,
							Tags:            in.Tags,
						})

					default:
						dataMap[mapKey] = doMapConvertForMapOrStructValue(doMapConvertForMapOrStructValueInput{
							IsRoot:          false,
							Value:           rvInterface,
							RecursiveType:   in.RecursiveType,
							RecursiveOption: in.RecursiveType == recursiveTypeTrue,
							Tags:            in.Tags,
						})
					}

				// The struct attribute is type of slice.
				case reflect.Array, reflect.Slice:
					length := rvAttrField.Len()
					if length == 0 {
						dataMap[mapKey] = rvAttrField.Interface()
						break
					}
					array := make([]interface{}, length)
					for arrayIndex := 0; arrayIndex < length; arrayIndex++ {
						array[arrayIndex] = doMapConvertForMapOrStructValue(
							doMapConvertForMapOrStructValueInput{
								IsRoot:          false,
								Value:           rvAttrField.Index(arrayIndex).Interface(),
								RecursiveType:   in.RecursiveType,
								RecursiveOption: in.RecursiveType == recursiveTypeTrue,
								Tags:            in.Tags,
							},
						)
					}
					dataMap[mapKey] = array
				case reflect.Map:
					var (
						mapKeys   = rvAttrField.MapKeys()
						nestedMap = make(map[string]interface{})
					)
					for _, k := range mapKeys {
						nestedMap[String(k.Interface())] = doMapConvertForMapOrStructValue(
							doMapConvertForMapOrStructValueInput{
								IsRoot:          false,
								Value:           rvAttrField.MapIndex(k).Interface(),
								RecursiveType:   in.RecursiveType,
								RecursiveOption: in.RecursiveType == recursiveTypeTrue,
								Tags:            in.Tags,
							},
						)
					}
					dataMap[mapKey] = nestedMap
				default:
					if rvField.IsValid() {
						dataMap[mapKey] = reflectValue.Field(i).Interface()
					} else {
						dataMap[mapKey] = nil
					}
				}
			} else {
				// No recursive map value converting
				if rvField.IsValid() {
					dataMap[mapKey] = reflectValue.Field(i).Interface()
				} else {
					dataMap[mapKey] = nil
				}
			}
		}
		if len(dataMap) == 0 {
			return in.Value
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
			array[i] = doMapConvertForMapOrStructValue(doMapConvertForMapOrStructValueInput{
				IsRoot:          false,
				Value:           reflectValue.Index(i).Interface(),
				RecursiveType:   in.RecursiveType,
				RecursiveOption: in.RecursiveType == recursiveTypeTrue,
				Tags:            in.Tags,
			})
		}
		return array
	}
	return in.Value
}

// MapStrStr converts `value` to map[string]string.
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

// MapStrStrDeep converts `value` to map[string]string recursively.
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
