// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

//// MapStrStr converts `value` to map[string]string.
//// Note that there might be data copy for this map type converting.
//func (c *impConverter) MapStrStr(value any, option ...MapOption) (map[string]string, error) {
//	if r, ok := value.(map[string]string); ok {
//		return r, nil
//	}
//	m := Map(value, option...)
//	if len(m) > 0 {
//		vMap := make(map[string]string, len(m))
//		for k, v := range m {
//			s, err := c.String(v)
//			if err != nil {
//				return nil, err
//			}
//			vMap[k] = s
//		}
//		return vMap, nil
//	}
//	return nil, nil
//}
//
//// Map implements the map converting.
//// It automatically checks and converts json string to map if `value` is string/[]byte.
////
//// TODO completely implement the recursive converting for all types, especially the map.
//func (c *impConverter) Map(value any, option MapOption) (map[string]any, error) {
//	if value == nil {
//		return nil, nil
//	}
//	// It redirects to its underlying value if it has implemented interface iVal.
//	if v, ok := value.(localinterface.IVal); ok {
//		value = v.Val()
//	}
//	var (
//		newTags = gtag.StructTagPriority
//	)
//	switch len(option.Tags) {
//	case 0:
//		// No need handling.
//	case 1:
//		newTags = append(strings.Split(option.Tags[0], ","), gtag.StructTagPriority...)
//	default:
//		newTags = append(option.Tags, gtag.StructTagPriority...)
//	}
//	// Assert the common combination of types, and finally it uses reflection.
//	dataMap := make(map[string]any)
//	switch r := value.(type) {
//	case string:
//		// If it is a JSON string, automatically unmarshal it!
//		if len(r) > 0 && r[0] == '{' && r[len(r)-1] == '}' {
//			if err := json.UnmarshalUseNumber([]byte(r), &dataMap); err != nil {
//				return nil, err
//			}
//		} else {
//			return nil, nil
//		}
//	case []byte:
//		// If it is a JSON string, automatically unmarshal it!
//		if len(r) > 0 && r[0] == '{' && r[len(r)-1] == '}' {
//			if err := json.UnmarshalUseNumber(r, &dataMap); err != nil {
//				return nil, err
//			}
//		} else {
//			return nil, nil
//		}
//	case map[any]any:
//		recursiveOption := option
//		recursiveOption.Tags = newTags
//		for k, v := range r {
//			dataMap[String(k)] = c.doMapConvertForMapOrStructValue(
//				doMapConvertForMapOrStructValueInput{
//					IsRoot:    false,
//					Value:     v,
//					Recursive: option.Deep,
//					Option:    recursiveOption,
//				},
//			)
//		}
//	case map[any]string:
//		for k, v := range r {
//			dataMap[String(k)] = v
//		}
//	case map[any]int:
//		for k, v := range r {
//			dataMap[String(k)] = v
//		}
//	case map[any]uint:
//		for k, v := range r {
//			dataMap[String(k)] = v
//		}
//	case map[any]float32:
//		for k, v := range r {
//			dataMap[String(k)] = v
//		}
//	case map[any]float64:
//		for k, v := range r {
//			dataMap[String(k)] = v
//		}
//	case map[string]bool:
//		for k, v := range r {
//			dataMap[k] = v
//		}
//	case map[string]int:
//		for k, v := range r {
//			dataMap[k] = v
//		}
//	case map[string]uint:
//		for k, v := range r {
//			dataMap[k] = v
//		}
//	case map[string]float32:
//		for k, v := range r {
//			dataMap[k] = v
//		}
//	case map[string]float64:
//		for k, v := range r {
//			dataMap[k] = v
//		}
//	case map[string]string:
//		for k, v := range r {
//			dataMap[k] = v
//		}
//	case map[string]any:
//		if option.Deep {
//			recursiveOption := option
//			recursiveOption.Tags = newTags
//			// A copy of current map.
//			for k, v := range r {
//				dataMap[k] = c.doMapConvertForMapOrStructValue(
//					doMapConvertForMapOrStructValueInput{
//						IsRoot:    false,
//						Value:     v,
//						Recursive: option.Deep,
//						Option:    recursiveOption,
//					},
//				)
//			}
//		} else {
//			// It returns the map directly without any changing.
//			return r, nil
//		}
//	case map[int]any:
//		recursiveOption := option
//		recursiveOption.Tags = newTags
//		for k, v := range r {
//			dataMap[String(k)] = c.doMapConvertForMapOrStructValue(
//				doMapConvertForMapOrStructValueInput{
//					IsRoot:    false,
//					Value:     v,
//					Recursive: option.Deep,
//					Option:    recursiveOption,
//				},
//			)
//		}
//	case map[int]string:
//		for k, v := range r {
//			dataMap[String(k)] = v
//		}
//	case map[uint]string:
//		for k, v := range r {
//			dataMap[String(k)] = v
//		}
//
//	default:
//		// Not a common type, it then uses reflection for conversion.
//		var reflectValue reflect.Value
//		if v, ok := value.(reflect.Value); ok {
//			reflectValue = v
//		} else {
//			reflectValue = reflect.ValueOf(value)
//		}
//		reflectKind := reflectValue.Kind()
//		// If it is a pointer, we should find its real data type.
//		for reflectKind == reflect.Ptr {
//			reflectValue = reflectValue.Elem()
//			reflectKind = reflectValue.Kind()
//		}
//		switch reflectKind {
//		// If `value` is type of array, it converts the value of even number index as its key and
//		// the value of odd number index as its corresponding value, for example:
//		// []string{"k1","v1","k2","v2"} => map[string]any{"k1":"v1", "k2":"v2"}
//		// []string{"k1","v1","k2"}      => map[string]any{"k1":"v1", "k2":nil}
//		case reflect.Slice, reflect.Array:
//			length := reflectValue.Len()
//			for i := 0; i < length; i += 2 {
//				if i+1 < length {
//					dataMap[String(reflectValue.Index(i).Interface())] = reflectValue.Index(i + 1).Interface()
//				} else {
//					dataMap[String(reflectValue.Index(i).Interface())] = nil
//				}
//			}
//		case reflect.Map, reflect.Struct, reflect.Interface:
//			recursiveOption := option
//			recursiveOption.Tags = newTags
//			convertedValue := c.doMapConvertForMapOrStructValue(
//				doMapConvertForMapOrStructValueInput{
//					IsRoot:        true,
//					Value:         value,
//					Recursive:     option.Deep,
//					Option:        recursiveOption,
//					MustMapReturn: option.EmptyEvenNil,
//				},
//			)
//			if m, ok := convertedValue.(map[string]any); ok {
//				return m, nil
//			}
//			return nil, nil
//		default:
//			return nil, nil
//		}
//	}
//	return dataMap, nil
//}
//
//func getUsedMapOption(option ...MapOption) MapOption {
//	var usedOption MapOption
//	if len(option) > 0 {
//		usedOption = option[0]
//	}
//	return usedOption
//}
//
//type doMapConvertForMapOrStructValueInput struct {
//	IsRoot        bool      // It returns directly if it is not root and with no recursive converting.
//	Value         any       // Current operation value.
//	Recursive     bool      // Whether convert recursively for `current` operation.
//	Option        MapOption // Map converting option.
//	MustMapReturn bool      // Must return map instead of Value when empty.
//}
//
//func (c *impConverter) doMapConvertForMapOrStructValue(in doMapConvertForMapOrStructValueInput) any {
//	if !in.IsRoot && !in.Recursive {
//		return in.Value
//	}
//
//	var reflectValue reflect.Value
//	if v, ok := in.Value.(reflect.Value); ok {
//		reflectValue = v
//		in.Value = v.Interface()
//	} else {
//		reflectValue = reflect.ValueOf(in.Value)
//	}
//	reflectKind := reflectValue.Kind()
//	// If it is a pointer, we should find its real data type.
//	for reflectKind == reflect.Ptr {
//		reflectValue = reflectValue.Elem()
//		reflectKind = reflectValue.Kind()
//	}
//	switch reflectKind {
//	case reflect.Map:
//		var (
//			mapIter = reflectValue.MapRange()
//			dataMap = make(map[string]any)
//		)
//		for mapIter.Next() {
//			var (
//				mapKeyValue = mapIter.Value()
//				mapValue    any
//			)
//			switch {
//			case mapKeyValue.IsZero():
//				if utils.CanCallIsNil(mapKeyValue) && mapKeyValue.IsNil() {
//					// quick check for nil value.
//					mapValue = nil
//				} else {
//					// in case of:
//					// exception recovered: reflect: call of reflect.Value.Interface on zero Value
//					mapValue = reflect.New(mapKeyValue.Type()).Elem().Interface()
//				}
//			default:
//				mapValue = mapKeyValue.Interface()
//			}
//			dataMap[String(mapIter.Key().Interface())] = c.doMapConvertForMapOrStructValue(
//				doMapConvertForMapOrStructValueInput{
//					IsRoot:    false,
//					Value:     mapValue,
//					Recursive: in.Recursive,
//					Option:    in.Option,
//				},
//			)
//		}
//		return dataMap
//
//	case reflect.Struct:
//		var dataMap = make(map[string]any)
//		// Map converting interface check.
//		if v, ok := in.Value.(localinterface.IMapStrAny); ok {
//			// Value copy, in case of concurrent safety.
//			for mapK, mapV := range v.MapStrAny() {
//				if in.Recursive {
//					dataMap[mapK] = c.doMapConvertForMapOrStructValue(
//						doMapConvertForMapOrStructValueInput{
//							IsRoot:    false,
//							Value:     mapV,
//							Recursive: in.Recursive,
//							Option:    in.Option,
//						},
//					)
//				} else {
//					dataMap[mapK] = mapV
//				}
//			}
//			if len(dataMap) > 0 {
//				return dataMap
//			}
//		}
//		// Using reflect for converting.
//		var (
//			rtField     reflect.StructField
//			rvField     reflect.Value
//			reflectType = reflectValue.Type() // attribute value type.
//			mapKey      = ""                  // mapKey may be the tag name or the struct attribute name.
//		)
//		for i := 0; i < reflectValue.NumField(); i++ {
//			rtField = reflectType.Field(i)
//			rvField = reflectValue.Field(i)
//			// Only convert the public attributes.
//			fieldName := rtField.Name
//			if !utils.IsLetterUpper(fieldName[0]) {
//				continue
//			}
//			mapKey = ""
//			fieldTag := rtField.Tag
//			for _, tag := range in.Option.Tags {
//				if mapKey = fieldTag.Get(tag); mapKey != "" {
//					break
//				}
//			}
//			if mapKey == "" {
//				mapKey = fieldName
//			} else {
//				// Support json tag feature: -, omitempty
//				mapKey = strings.TrimSpace(mapKey)
//				if mapKey == "-" {
//					continue
//				}
//				array := strings.Split(mapKey, ",")
//				if len(array) > 1 {
//					switch strings.TrimSpace(array[1]) {
//					case "omitempty":
//						if in.Option.OmitEmpty && empty.IsEmpty(rvField.Interface()) {
//							continue
//						} else {
//							mapKey = strings.TrimSpace(array[0])
//						}
//					default:
//						mapKey = strings.TrimSpace(array[0])
//					}
//				}
//				if mapKey == "" {
//					mapKey = fieldName
//				}
//			}
//			if in.Recursive || rtField.Anonymous {
//				// Do map converting recursively.
//				var (
//					rvAttrField = rvField
//					rvAttrKind  = rvField.Kind()
//				)
//				if rvAttrKind == reflect.Ptr {
//					rvAttrField = rvField.Elem()
//					rvAttrKind = rvAttrField.Kind()
//				}
//				switch rvAttrKind {
//				case reflect.Struct:
//					// Embedded struct and has no fields, just ignores it.
//					// Eg: gmeta.Meta
//					if rvAttrField.Type().NumField() == 0 {
//						continue
//					}
//					var (
//						hasNoTag = mapKey == fieldName
//						// DO NOT use rvAttrField.Interface() here,
//						// as it might be changed from pointer to struct.
//						rvInterface = rvField.Interface()
//					)
//					switch {
//					case hasNoTag && rtField.Anonymous:
//						// It means this attribute field has no tag.
//						// Overwrite the attribute with sub-struct attribute fields.
//						anonymousValue := c.doMapConvertForMapOrStructValue(doMapConvertForMapOrStructValueInput{
//							IsRoot:    false,
//							Value:     rvInterface,
//							Recursive: true,
//							Option:    in.Option,
//						})
//						if m, ok := anonymousValue.(map[string]any); ok {
//							for k, v := range m {
//								dataMap[k] = v
//							}
//						} else {
//							dataMap[mapKey] = rvInterface
//						}
//
//					// It means this attribute field has desired tag.
//					case !hasNoTag && rtField.Anonymous:
//						dataMap[mapKey] = c.doMapConvertForMapOrStructValue(doMapConvertForMapOrStructValueInput{
//							IsRoot:    false,
//							Value:     rvInterface,
//							Recursive: true,
//							Option:    in.Option,
//						})
//
//					default:
//						dataMap[mapKey] = c.doMapConvertForMapOrStructValue(doMapConvertForMapOrStructValueInput{
//							IsRoot:    false,
//							Value:     rvInterface,
//							Recursive: in.Recursive,
//							Option:    in.Option,
//						})
//					}
//
//				// The struct attribute is type of slice.
//				case reflect.Array, reflect.Slice:
//					length := rvAttrField.Len()
//					if length == 0 {
//						dataMap[mapKey] = rvAttrField.Interface()
//						break
//					}
//					array := make([]any, length)
//					for arrayIndex := 0; arrayIndex < length; arrayIndex++ {
//						array[arrayIndex] = c.doMapConvertForMapOrStructValue(
//							doMapConvertForMapOrStructValueInput{
//								IsRoot:    false,
//								Value:     rvAttrField.Index(arrayIndex).Interface(),
//								Recursive: in.Recursive,
//								Option:    in.Option,
//							},
//						)
//					}
//					dataMap[mapKey] = array
//				case reflect.Map:
//					var (
//						mapIter   = rvAttrField.MapRange()
//						nestedMap = make(map[string]any)
//					)
//					for mapIter.Next() {
//						nestedMap[String(mapIter.Key().Interface())] = c.doMapConvertForMapOrStructValue(
//							doMapConvertForMapOrStructValueInput{
//								IsRoot:    false,
//								Value:     mapIter.Value().Interface(),
//								Recursive: in.Recursive,
//								Option:    in.Option,
//							},
//						)
//					}
//					dataMap[mapKey] = nestedMap
//				default:
//					if rvField.IsValid() {
//						dataMap[mapKey] = reflectValue.Field(i).Interface()
//					} else {
//						dataMap[mapKey] = nil
//					}
//				}
//			} else {
//				// No recursive map value converting
//				if rvField.IsValid() {
//					dataMap[mapKey] = reflectValue.Field(i).Interface()
//				} else {
//					dataMap[mapKey] = nil
//				}
//			}
//		}
//		if !in.MustMapReturn && len(dataMap) == 0 {
//			return in.Value
//		}
//		return dataMap
//
//	// The given value is type of slice.
//	case reflect.Array, reflect.Slice:
//		length := reflectValue.Len()
//		if length == 0 {
//			break
//		}
//		array := make([]any, reflectValue.Len())
//		for i := 0; i < length; i++ {
//			array[i] = c.doMapConvertForMapOrStructValue(doMapConvertForMapOrStructValueInput{
//				IsRoot:    false,
//				Value:     reflectValue.Index(i).Interface(),
//				Recursive: in.Recursive,
//				Option:    in.Option,
//			})
//		}
//		return array
//
//	default:
//	}
//	return in.Value
//}
