// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package converter

import (
	"reflect"
	"strings"

	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
	"github.com/gogf/gf/v2/util/gtag"
)

// MapOption specifies the option for map converting.
type MapOption struct {
	// Deep marks doing Map function recursively, which means if the attribute of given converting value
	// is also a struct/*struct, it automatically calls Map function on this attribute converting it to
	// a map[string]any type variable.
	Deep bool

	// OmitEmpty ignores the attributes that has json `omitempty` tag.
	OmitEmpty bool

	// Tags specifies the converted map key name by struct tag name.
	Tags []string

	// ContinueOnError specifies whether to continue converting the next element
	// if one element converting fails.
	ContinueOnError bool
}

func (c *Converter) getMapOption(option ...MapOption) MapOption {
	if len(option) > 0 {
		return option[0]
	}
	return MapOption{}
}

// Map converts any variable `value` to map[string]any. If the parameter `value` is not a
// map/struct/*struct type, then the conversion will fail and returns nil.
//
// If `value` is a struct/*struct object, the second parameter `priorityTagAndFieldName` specifies the most priority
// priorityTagAndFieldName that will be detected, otherwise it detects the priorityTagAndFieldName in order of:
// gconv, json, field name.
func (c *Converter) Map(value any, option ...MapOption) (map[string]any, error) {
	return c.doMapConvert(value, RecursiveTypeAuto, false, c.getMapOption(option...))
}

// MapStrStr converts `value` to map[string]string.
// Note that there might be data copy for this map type converting.
func (c *Converter) MapStrStr(value any, option ...MapOption) (map[string]string, error) {
	if r, ok := value.(map[string]string); ok {
		return r, nil
	}
	m, err := c.Map(value, option...)
	if err != nil {
		return nil, err
	}
	if len(m) > 0 {
		var (
			s         string
			vMap      = make(map[string]string, len(m))
			mapOption = c.getMapOption(option...)
		)
		for k, v := range m {
			s, err = c.String(v)
			if err != nil && !mapOption.ContinueOnError {
				return nil, err
			}
			vMap[k] = s
		}
		return vMap, nil
	}
	return nil, nil
}

// MapConvert implements the map converting.
// It automatically checks and converts json string to map if `value` is string/[]byte.
//
// TODO completely implement the recursive converting for all types, especially the map.
func (c *Converter) doMapConvert(
	value any, recursive RecursiveType, mustMapReturn bool, option MapOption,
) (map[string]any, error) {
	if value == nil {
		return nil, nil
	}
	// It redirects to its underlying value if it has implemented interface iVal.
	if v, ok := value.(localinterface.IVal); ok {
		value = v.Val()
	}
	var (
		err     error
		newTags = gtag.StructTagPriority
	)
	if option.Deep {
		recursive = RecursiveTypeTrue
	}
	switch len(option.Tags) {
	case 0:
		// No need handling.
	case 1:
		newTags = append(strings.Split(option.Tags[0], ","), gtag.StructTagPriority...)
	default:
		newTags = append(option.Tags, gtag.StructTagPriority...)
	}
	// Assert the common combination of types, and finally it uses reflection.
	dataMap := make(map[string]interface{})
	switch r := value.(type) {
	case string:
		// If it is a JSON string, automatically unmarshal it!
		if len(r) > 0 && r[0] == '{' && r[len(r)-1] == '}' {
			if err = json.UnmarshalUseNumber([]byte(r), &dataMap); err != nil {
				return nil, err
			}
		} else {
			return nil, nil
		}
	case []byte:
		// If it is a JSON string, automatically unmarshal it!
		if len(r) > 0 && r[0] == '{' && r[len(r)-1] == '}' {
			if err = json.UnmarshalUseNumber(r, &dataMap); err != nil {
				return nil, err
			}
		} else {
			return nil, nil
		}
	case map[interface{}]interface{}:
		recursiveOption := option
		recursiveOption.Tags = newTags
		for k, v := range r {
			s, err := c.String(k)
			if err != nil && !option.ContinueOnError {
				return nil, err
			}
			dataMap[s], err = c.doMapConvertForMapOrStructValue(
				doMapConvertForMapOrStructValueInput{
					IsRoot:          false,
					Value:           v,
					RecursiveType:   recursive,
					RecursiveOption: recursive == RecursiveTypeTrue,
					Option:          recursiveOption,
				},
			)
			if err != nil && !option.ContinueOnError {
				return nil, err
			}
		}
	case map[interface{}]string:
		for k, v := range r {
			s, err := c.String(k)
			if err != nil && !option.ContinueOnError {
				return nil, err
			}
			dataMap[s] = v
		}
	case map[interface{}]int:
		for k, v := range r {
			s, err := c.String(k)
			if err != nil && !option.ContinueOnError {
				return nil, err
			}
			dataMap[s] = v
		}
	case map[interface{}]uint:
		for k, v := range r {
			s, err := c.String(k)
			if err != nil && !option.ContinueOnError {
				return nil, err
			}
			dataMap[s] = v
		}
	case map[interface{}]float32:
		for k, v := range r {
			s, err := c.String(k)
			if err != nil && !option.ContinueOnError {
				return nil, err
			}
			dataMap[s] = v
		}
	case map[interface{}]float64:
		for k, v := range r {
			s, err := c.String(k)
			if err != nil && !option.ContinueOnError {
				return nil, err
			}
			dataMap[s] = v
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
		if recursive == RecursiveTypeTrue {
			recursiveOption := option
			recursiveOption.Tags = newTags
			// A copy of current map.
			for k, v := range r {
				dataMap[k], err = c.doMapConvertForMapOrStructValue(
					doMapConvertForMapOrStructValueInput{
						IsRoot:          false,
						Value:           v,
						RecursiveType:   recursive,
						RecursiveOption: recursive == RecursiveTypeTrue,
						Option:          recursiveOption,
					},
				)
				if err != nil && !option.ContinueOnError {
					return nil, err
				}
			}
		} else {
			// It returns the map directly without any changing.
			return r, nil
		}
	case map[int]interface{}:
		recursiveOption := option
		recursiveOption.Tags = newTags
		for k, v := range r {
			s, err := c.String(k)
			if err != nil && !option.ContinueOnError {
				return nil, err
			}
			dataMap[s], err = c.doMapConvertForMapOrStructValue(
				doMapConvertForMapOrStructValueInput{
					IsRoot:          false,
					Value:           v,
					RecursiveType:   recursive,
					RecursiveOption: recursive == RecursiveTypeTrue,
					Option:          recursiveOption,
				},
			)
			if err != nil && !option.ContinueOnError {
				return nil, err
			}
		}
	case map[int]string:
		for k, v := range r {
			s, err := c.String(k)
			if err != nil && !option.ContinueOnError {
				return nil, err
			}
			dataMap[s] = v
		}
	case map[uint]string:
		for k, v := range r {
			s, err := c.String(k)
			if err != nil && !option.ContinueOnError {
				return nil, err
			}
			dataMap[s] = v
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
				s, err := c.String(reflectValue.Index(i).Interface())
				if err != nil && !option.ContinueOnError {
					return nil, err
				}
				if i+1 < length {
					dataMap[s] = reflectValue.Index(i + 1).Interface()
				} else {
					dataMap[s] = nil
				}
			}
		case reflect.Map, reflect.Struct, reflect.Interface:
			recursiveOption := option
			recursiveOption.Tags = newTags
			convertedValue, err := c.doMapConvertForMapOrStructValue(
				doMapConvertForMapOrStructValueInput{
					IsRoot:          true,
					Value:           value,
					RecursiveType:   recursive,
					RecursiveOption: recursive == RecursiveTypeTrue,
					Option:          recursiveOption,
					MustMapReturn:   mustMapReturn,
				},
			)
			if err != nil && !option.ContinueOnError {
				return nil, err
			}
			if m, ok := convertedValue.(map[string]interface{}); ok {
				return m, nil
			}
			return nil, nil
		default:
			return nil, nil
		}
	}
	return dataMap, nil
}

type doMapConvertForMapOrStructValueInput struct {
	IsRoot          bool          // It returns directly if it is not root and with no recursive converting.
	Value           interface{}   // Current operation value.
	RecursiveType   RecursiveType // The type from top function entry.
	RecursiveOption bool          // Whether convert recursively for `current` operation.
	Option          MapOption     // Map converting option.
	MustMapReturn   bool          // Must return map instead of Value when empty.
}

func (c *Converter) doMapConvertForMapOrStructValue(in doMapConvertForMapOrStructValueInput) (any, error) {
	if !in.IsRoot && !in.RecursiveOption {
		return in.Value, nil
	}

	var (
		err          error
		reflectValue reflect.Value
	)
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
			mapIter = reflectValue.MapRange()
			dataMap = make(map[string]interface{})
		)
		for mapIter.Next() {
			var (
				mapKeyValue = mapIter.Value()
				mapValue    interface{}
			)
			switch {
			case mapKeyValue.IsZero():
				if utils.CanCallIsNil(mapKeyValue) && mapKeyValue.IsNil() {
					// quick check for nil value.
					mapValue = nil
				} else {
					// in case of:
					// exception recovered: reflect: call of reflect.Value.Interface on zero Value
					mapValue = reflect.New(mapKeyValue.Type()).Elem().Interface()
				}
			default:
				mapValue = mapKeyValue.Interface()
			}
			s, err := c.String(mapIter.Key().Interface())
			if err != nil && !in.Option.ContinueOnError {
				return nil, err
			}
			dataMap[s], err = c.doMapConvertForMapOrStructValue(
				doMapConvertForMapOrStructValueInput{
					IsRoot:          false,
					Value:           mapValue,
					RecursiveType:   in.RecursiveType,
					RecursiveOption: in.RecursiveType == RecursiveTypeTrue,
					Option:          in.Option,
				},
			)
			if err != nil && !in.Option.ContinueOnError {
				return nil, err
			}
		}
		return dataMap, nil

	case reflect.Struct:
		var dataMap = make(map[string]interface{})
		// Map converting interface check.
		if v, ok := in.Value.(localinterface.IMapStrAny); ok {
			// Value copy, in case of concurrent safety.
			for mapK, mapV := range v.MapStrAny() {
				if in.RecursiveOption {
					dataMap[mapK], err = c.doMapConvertForMapOrStructValue(
						doMapConvertForMapOrStructValueInput{
							IsRoot:          false,
							Value:           mapV,
							RecursiveType:   in.RecursiveType,
							RecursiveOption: in.RecursiveType == RecursiveTypeTrue,
							Option:          in.Option,
						},
					)
					if err != nil && !in.Option.ContinueOnError {
						return nil, err
					}
				} else {
					dataMap[mapK] = mapV
				}
			}
			if len(dataMap) > 0 {
				return dataMap, nil
			}
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
			for _, tag := range in.Option.Tags {
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
						if in.Option.OmitEmpty && empty.IsEmpty(rvField.Interface()) {
							continue
						} else {
							mapKey = strings.TrimSpace(array[0])
						}
					default:
						mapKey = strings.TrimSpace(array[0])
					}
				}
				if mapKey == "" {
					mapKey = fieldName
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
						anonymousValue, err := c.doMapConvertForMapOrStructValue(
							doMapConvertForMapOrStructValueInput{
								IsRoot:          false,
								Value:           rvInterface,
								RecursiveType:   in.RecursiveType,
								RecursiveOption: true,
								Option:          in.Option,
							},
						)
						if err != nil && !in.Option.ContinueOnError {
							return nil, err
						}
						if m, ok := anonymousValue.(map[string]interface{}); ok {
							for k, v := range m {
								dataMap[k] = v
							}
						} else {
							dataMap[mapKey] = rvInterface
						}

					// It means this attribute field has desired tag.
					case !hasNoTag && rtField.Anonymous:
						dataMap[mapKey], err = c.doMapConvertForMapOrStructValue(
							doMapConvertForMapOrStructValueInput{
								IsRoot:          false,
								Value:           rvInterface,
								RecursiveType:   in.RecursiveType,
								RecursiveOption: true,
								Option:          in.Option,
							},
						)
						if err != nil && !in.Option.ContinueOnError {
							return nil, err
						}

					default:
						dataMap[mapKey], err = c.doMapConvertForMapOrStructValue(
							doMapConvertForMapOrStructValueInput{
								IsRoot:          false,
								Value:           rvInterface,
								RecursiveType:   in.RecursiveType,
								RecursiveOption: in.RecursiveType == RecursiveTypeTrue,
								Option:          in.Option,
							},
						)
						if err != nil && !in.Option.ContinueOnError {
							return nil, err
						}
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
						array[arrayIndex], err = c.doMapConvertForMapOrStructValue(
							doMapConvertForMapOrStructValueInput{
								IsRoot:          false,
								Value:           rvAttrField.Index(arrayIndex).Interface(),
								RecursiveType:   in.RecursiveType,
								RecursiveOption: in.RecursiveType == RecursiveTypeTrue,
								Option:          in.Option,
							},
						)
						if err != nil && !in.Option.ContinueOnError {
							return nil, err
						}
					}
					dataMap[mapKey] = array
				case reflect.Map:
					var (
						mapIter   = rvAttrField.MapRange()
						nestedMap = make(map[string]interface{})
					)
					for mapIter.Next() {
						s, err := c.String(mapIter.Key().Interface())
						if err != nil && !in.Option.ContinueOnError {
							return nil, err
						}
						nestedMap[s], err = c.doMapConvertForMapOrStructValue(
							doMapConvertForMapOrStructValueInput{
								IsRoot:          false,
								Value:           mapIter.Value().Interface(),
								RecursiveType:   in.RecursiveType,
								RecursiveOption: in.RecursiveType == RecursiveTypeTrue,
								Option:          in.Option,
							},
						)
						if err != nil && !in.Option.ContinueOnError {
							return nil, err
						}
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
		if !in.MustMapReturn && len(dataMap) == 0 {
			return in.Value, nil
		}
		return dataMap, nil

	// The given value is type of slice.
	case reflect.Array, reflect.Slice:
		length := reflectValue.Len()
		if length == 0 {
			break
		}
		array := make([]interface{}, reflectValue.Len())
		for i := 0; i < length; i++ {
			array[i], err = c.doMapConvertForMapOrStructValue(doMapConvertForMapOrStructValueInput{
				IsRoot:          false,
				Value:           reflectValue.Index(i).Interface(),
				RecursiveType:   in.RecursiveType,
				RecursiveOption: in.RecursiveType == RecursiveTypeTrue,
				Option:          in.Option,
			})
			if err != nil && !in.Option.ContinueOnError {
				return nil, err
			}
		}
		return array, nil

	default:
	}
	return in.Value, nil
}
