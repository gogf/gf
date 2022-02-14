// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"reflect"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gstructs"
)

// Struct maps the params key-value pairs to the corresponding struct object's attributes.
// The third parameter `mapping` is unnecessary, indicating the mapping rules between the
// custom key name and the attribute name(case-sensitive).
//
// Note:
// 1. The `params` can be any type of map/struct, usually a map.
// 2. The `pointer` should be type of *struct/**struct, which is a pointer to struct object
//    or struct pointer.
// 3. Only the public attributes of struct object can be mapped.
// 4. If `params` is a map, the key of the map `params` can be lowercase.
//    It will automatically convert the first letter of the key to uppercase
//    in mapping procedure to do the matching.
//    It ignores the map key, if it does not match.
func Struct(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	return Scan(params, pointer, mapping...)
}

// StructTag acts as Struct but also with support for priority tag feature, which retrieves the
// specified tags for `params` key-value items to struct attribute names mapping.
// The parameter `priorityTag` supports multiple tags that can be joined with char ','.
func StructTag(params interface{}, pointer interface{}, priorityTag string) (err error) {
	return doStruct(params, pointer, nil, priorityTag)
}

// doStructWithJsonCheck checks if given `params` is JSON, it then uses json.Unmarshal doing the converting.
func doStructWithJsonCheck(params interface{}, pointer interface{}) (err error, ok bool) {
	switch r := params.(type) {
	case []byte:
		if json.Valid(r) {
			if rv, ok := pointer.(reflect.Value); ok {
				if rv.Kind() == reflect.Ptr {
					if rv.IsNil() {
						return nil, false
					}
					return json.UnmarshalUseNumber(r, rv.Interface()), true
				} else if rv.CanAddr() {
					return json.UnmarshalUseNumber(r, rv.Addr().Interface()), true
				}
			} else {
				return json.UnmarshalUseNumber(r, pointer), true
			}
		}
	case string:
		if paramsBytes := []byte(r); json.Valid(paramsBytes) {
			if rv, ok := pointer.(reflect.Value); ok {
				if rv.Kind() == reflect.Ptr {
					if rv.IsNil() {
						return nil, false
					}
					return json.UnmarshalUseNumber(paramsBytes, rv.Interface()), true
				} else if rv.CanAddr() {
					return json.UnmarshalUseNumber(paramsBytes, rv.Addr().Interface()), true
				}
			} else {
				return json.UnmarshalUseNumber(paramsBytes, pointer), true
			}
		}
	default:
		// The `params` might be struct that implements interface function Interface, eg: gvar.Var.
		if v, ok := params.(iInterface); ok {
			return doStructWithJsonCheck(v.Interface(), pointer)
		}
	}
	return nil, false
}

// doStruct is the core internal converting function for any data to struct.
func doStruct(params interface{}, pointer interface{}, mapping map[string]string, priorityTag string) (err error) {
	if params == nil {
		// If `params` is nil, no conversion.
		return nil
	}
	if pointer == nil {
		return gerror.NewCode(gcode.CodeInvalidParameter, "object pointer cannot be nil")
	}

	defer func() {
		// Catch the panic, especially the reflection operation panics.
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodeSkipf(gcode.CodeInternalError, 1, "%+v", exception)
			}
		}
	}()

	// JSON content converting.
	err, ok := doStructWithJsonCheck(params, pointer)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	var (
		paramsReflectValue      reflect.Value
		paramsInterface         interface{} // DO NOT use `params` directly as it might be type `reflect.Value`
		pointerReflectValue     reflect.Value
		pointerReflectKind      reflect.Kind
		pointerElemReflectValue reflect.Value // The pointed element.
	)
	if v, ok := params.(reflect.Value); ok {
		paramsReflectValue = v
	} else {
		paramsReflectValue = reflect.ValueOf(params)
	}
	paramsInterface = paramsReflectValue.Interface()
	if v, ok := pointer.(reflect.Value); ok {
		pointerReflectValue = v
		pointerElemReflectValue = v
	} else {
		pointerReflectValue = reflect.ValueOf(pointer)
		pointerReflectKind = pointerReflectValue.Kind()
		if pointerReflectKind != reflect.Ptr {
			return gerror.NewCodef(gcode.CodeInvalidParameter, "object pointer should be type of '*struct', but got '%v'", pointerReflectKind)
		}
		// Using IsNil on reflect.Ptr variable is OK.
		if !pointerReflectValue.IsValid() || pointerReflectValue.IsNil() {
			return gerror.NewCode(gcode.CodeInvalidParameter, "object pointer cannot be nil")
		}
		pointerElemReflectValue = pointerReflectValue.Elem()
	}

	// If `params` and `pointer` are the same type, the do directly assignment.
	// For performance enhancement purpose.
	if pointerElemReflectValue.IsValid() && pointerElemReflectValue.Type() == paramsReflectValue.Type() {
		pointerElemReflectValue.Set(paramsReflectValue)
		return nil
	}

	// Normal unmarshalling interfaces checks.
	if err, ok = bindVarToReflectValueWithInterfaceCheck(pointerReflectValue, paramsInterface); ok {
		return err
	}

	// It automatically creates struct object if necessary.
	// For example, if `pointer` is **User, then `elem` is *User, which is a pointer to User.
	if pointerElemReflectValue.Kind() == reflect.Ptr {
		if !pointerElemReflectValue.IsValid() || pointerElemReflectValue.IsNil() {
			e := reflect.New(pointerElemReflectValue.Type().Elem()).Elem()
			pointerElemReflectValue.Set(e.Addr())
		}
		// if v, ok := pointerElemReflectValue.Interface().(iUnmarshalValue); ok {
		//	return v.UnmarshalValue(params)
		// }
		// Note that it's `pointerElemReflectValue` here not `pointerReflectValue`.
		if err, ok = bindVarToReflectValueWithInterfaceCheck(pointerElemReflectValue, paramsInterface); ok {
			return err
		}
		// Retrieve its element, may be struct at last.
		pointerElemReflectValue = pointerElemReflectValue.Elem()
	}

	// paramsMap is the map[string]interface{} type variable for params.
	// DO NOT use MapDeep here.
	paramsMap := Map(paramsInterface)
	if paramsMap == nil {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`convert params from "%#v" to "map[string]interface{}" failed`,
			params,
		)
	}

	// It only performs one converting to the same attribute.
	// doneMap is used to check repeated converting, its key is the real attribute name
	// of the struct.
	doneMap := make(map[string]struct{})

	// The key of the attrMap is the attribute name of the struct,
	// and the value is its replaced name for later comparison to improve performance.
	var (
		tempName       string
		elemFieldType  reflect.StructField
		elemFieldValue reflect.Value
		elemType       = pointerElemReflectValue.Type()
		attrMap        = make(map[string]string)
	)
	for i := 0; i < pointerElemReflectValue.NumField(); i++ {
		elemFieldType = elemType.Field(i)
		// Only do converting to public attributes.
		if !utils.IsLetterUpper(elemFieldType.Name[0]) {
			continue
		}
		// Maybe it's struct/*struct embedded.
		if elemFieldType.Anonymous {
			elemFieldValue = pointerElemReflectValue.Field(i)
			// Ignore the interface attribute if it's nil.
			if elemFieldValue.Kind() == reflect.Interface {
				elemFieldValue = elemFieldValue.Elem()
				if !elemFieldValue.IsValid() {
					continue
				}
			}
			if err = doStruct(paramsMap, elemFieldValue, mapping, priorityTag); err != nil {
				return err
			}
		} else {
			tempName = elemFieldType.Name
			attrMap[tempName] = utils.RemoveSymbols(tempName)
		}
	}
	if len(attrMap) == 0 {
		return nil
	}

	// The key of the tagMap is the attribute name of the struct,
	// and the value is its replaced tag name for later comparison to improve performance.
	var (
		tagMap           = make(map[string]string)
		priorityTagArray []string
	)
	if priorityTag != "" {
		priorityTagArray = append(utils.SplitAndTrim(priorityTag, ","), StructTagPriority...)
	} else {
		priorityTagArray = StructTagPriority
	}
	tagToNameMap, err := gstructs.TagMapName(pointerElemReflectValue, priorityTagArray)
	if err != nil {
		return err
	}
	for tagName, attributeName := range tagToNameMap {
		// If there's something else in the tag string,
		// it uses the first part which is split using char ','.
		// Eg:
		// orm:"id, priority"
		// orm:"name, with:uid=id"
		tagMap[attributeName] = utils.RemoveSymbols(strings.Split(tagName, ",")[0])

		// If tag and attribute values both exist in `paramsMap`,
		// it then uses the tag value overwriting the attribute value in `paramsMap`.
		if paramsMap[tagName] != nil && paramsMap[attributeName] != nil {
			paramsMap[attributeName] = paramsMap[tagName]
		}
	}

	var (
		attrName  string
		checkName string
	)
	for mapK, mapV := range paramsMap {
		attrName = ""
		// It firstly checks the passed mapping rules.
		if len(mapping) > 0 {
			if passedAttrKey, ok := mapping[mapK]; ok {
				attrName = passedAttrKey
			}
		}
		// It secondly checks the predefined tags and matching rules.
		if attrName == "" {
			checkName = utils.RemoveSymbols(mapK)
			// Loop to find the matched attribute name with or without
			// string cases and chars like '-'/'_'/'.'/' '.

			// Matching the parameters to struct tag names.
			// The `attrKey` is the attribute name of the struct.
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
		}

		// No matching, it gives up this attribute converting.
		if attrName == "" {
			continue
		}
		// If the attribute name is already checked converting, then skip it.
		if _, ok = doneMap[attrName]; ok {
			continue
		}
		// Mark it done.
		doneMap[attrName] = struct{}{}
		if err = bindVarToStructAttr(pointerElemReflectValue, attrName, mapV, mapping); err != nil {
			return err
		}
	}
	return nil
}

// bindVarToStructAttr sets value to struct object attribute by name.
func bindVarToStructAttr(elem reflect.Value, name string, value interface{}, mapping map[string]string) (err error) {
	structFieldValue := elem.FieldByName(name)
	if !structFieldValue.IsValid() {
		return nil
	}
	// CanSet checks whether attribute is public accessible.
	if !structFieldValue.CanSet() {
		return nil
	}
	defer func() {
		if exception := recover(); exception != nil {
			if err = bindVarToReflectValue(structFieldValue, value, mapping); err != nil {
				err = gerror.Wrapf(err, `error binding value to attribute "%s"`, name)
			}
		}
	}()
	// Directly converting.
	if empty.IsNil(value) {
		structFieldValue.Set(reflect.Zero(structFieldValue.Type()))
	} else {
		structFieldValue.Set(reflect.ValueOf(doConvert(
			doConvertInput{
				FromValue:  value,
				ToTypeName: structFieldValue.Type().String(),
				ReferValue: structFieldValue,
			},
		)))
	}
	return nil
}

// bindVarToReflectValueWithInterfaceCheck does bind using common interfaces checks.
func bindVarToReflectValueWithInterfaceCheck(reflectValue reflect.Value, value interface{}) (err error, ok bool) {
	var pointer interface{}
	if reflectValue.Kind() != reflect.Ptr && reflectValue.CanAddr() {
		reflectValueAddr := reflectValue.Addr()
		if reflectValueAddr.IsNil() || !reflectValueAddr.IsValid() {
			return nil, false
		}
		// Not a pointer, but can token address, that makes it can be unmarshalled.
		pointer = reflectValue.Addr().Interface()
	} else {
		if reflectValue.IsNil() || !reflectValue.IsValid() {
			return nil, false
		}
		pointer = reflectValue.Interface()
	}
	// UnmarshalValue.
	if v, ok := pointer.(iUnmarshalValue); ok {
		return v.UnmarshalValue(value), ok
	}
	// UnmarshalText.
	if v, ok := pointer.(iUnmarshalText); ok {
		var valueBytes []byte
		if b, ok := value.([]byte); ok {
			valueBytes = b
		} else if s, ok := value.(string); ok {
			valueBytes = []byte(s)
		}
		if len(valueBytes) > 0 {
			return v.UnmarshalText(valueBytes), ok
		}
	}
	// UnmarshalJSON.
	if v, ok := pointer.(iUnmarshalJSON); ok {
		var valueBytes []byte
		if b, ok := value.([]byte); ok {
			valueBytes = b
		} else if s, ok := value.(string); ok {
			valueBytes = []byte(s)
		}

		if len(valueBytes) > 0 {
			// If it is not a valid JSON string, it then adds char `"` on its both sides to make it is.
			if !json.Valid(valueBytes) {
				newValueBytes := make([]byte, len(valueBytes)+2)
				newValueBytes[0] = '"'
				newValueBytes[len(newValueBytes)-1] = '"'
				copy(newValueBytes[1:], valueBytes)
				valueBytes = newValueBytes
			}
			return v.UnmarshalJSON(valueBytes), ok
		}
	}
	if v, ok := pointer.(iSet); ok {
		v.Set(value)
		return nil, ok
	}
	return nil, false
}

// bindVarToReflectValue sets `value` to reflect value object `structFieldValue`.
func bindVarToReflectValue(structFieldValue reflect.Value, value interface{}, mapping map[string]string) (err error) {
	// JSON content converting.
	err, ok := doStructWithJsonCheck(value, structFieldValue)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	// Common interface check.
	if err, ok := bindVarToReflectValueWithInterfaceCheck(structFieldValue, value); ok {
		return err
	}

	kind := structFieldValue.Kind()
	// Converting using interface, for some kinds.
	switch kind {
	case reflect.Slice, reflect.Array, reflect.Ptr, reflect.Interface:
		if !structFieldValue.IsNil() {
			if v, ok := structFieldValue.Interface().(iSet); ok {
				v.Set(value)
				return nil
			}
		}
	}

	// Converting by kind.
	switch kind {
	case reflect.Map:
		return doMapToMap(value, structFieldValue, mapping)

	case reflect.Struct:
		// Recursively converting for struct attribute.
		if err = doStruct(value, structFieldValue, nil, ""); err != nil {
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
						if err = doStruct(v.Index(i).Interface(), e, nil, ""); err != nil {
							// Note there's reflect conversion mechanism here.
							e.Set(reflect.ValueOf(v.Index(i).Interface()).Convert(t))
						}
						a.Index(i).Set(e.Addr())
					} else {
						e := reflect.New(t).Elem()
						if err = doStruct(v.Index(i).Interface(), e, nil, ""); err != nil {
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
				// Pointer element.
				e := reflect.New(t.Elem()).Elem()
				if err = doStruct(value, e, nil, ""); err != nil {
					// Note there's reflect conversion mechanism here.
					e.Set(reflect.ValueOf(value).Convert(t))
				}
				a.Index(0).Set(e.Addr())
			} else {
				// Just consider it as struct element. (Although it might be other types but not basic types, eg: map)
				e := reflect.New(t).Elem()
				if err = doStruct(value, e, nil, ""); err != nil {
					return err
				}
				a.Index(0).Set(e)
			}
		}
		structFieldValue.Set(a)

	case reflect.Ptr:
		item := reflect.New(structFieldValue.Type().Elem())
		if err, ok := bindVarToReflectValueWithInterfaceCheck(item, value); ok {
			structFieldValue.Set(item)
			return err
		}
		elem := item.Elem()
		if err = bindVarToReflectValue(elem, value, mapping); err == nil {
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
			if exception := recover(); exception != nil {
				err = gerror.NewCodef(
					gcode.CodeInternalError,
					`cannot convert value "%+v" to type "%s":%+v`,
					value,
					structFieldValue.Type().String(),
					exception,
				)
			}
		}()
		// It here uses reflect converting `value` to type of the attribute and assigns
		// the result value to the attribute. It might fail and panic if the usual Go
		// conversion rules do not allow conversion.
		structFieldValue.Set(reflect.ValueOf(value).Convert(structFieldValue.Type()))
	}
	return nil
}
