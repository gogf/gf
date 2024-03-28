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
	"github.com/gogf/gf/v2/util/gtag"
)

// Struct maps the params key-value pairs to the corresponding struct object's attributes.
// The third parameter `mapping` is unnecessary, indicating the mapping rules between the
// custom key name and the attribute name(case-sensitive).
//
// Note:
//  1. The `params` can be any type of map/struct, usually a map.
//  2. The `pointer` should be type of *struct/**struct, which is a pointer to struct object
//     or struct pointer.
//  3. Only the public attributes of struct object can be mapped.
//  4. If `params` is a map, the key of the map `params` can be lowercase.
//     It will automatically convert the first letter of the key to uppercase
//     in mapping procedure to do the matching.
//     It ignores the map key, if it does not match.
func Struct(params interface{}, pointer interface{}, paramKeyToAttrMap ...map[string]string) (err error) {
	return Scan(params, pointer, paramKeyToAttrMap...)
}

// StructTag acts as Struct but also with support for priority tag feature, which retrieves the
// specified tags for `params` key-value items to struct attribute names mapping.
// The parameter `priorityTag` supports multiple tags that can be joined with char ','.
func StructTag(params interface{}, pointer interface{}, priorityTag string) (err error) {
	return doStruct(params, pointer, nil, priorityTag)
}

// doStruct is the core internal converting function for any data to struct.
func doStruct(
	params interface{}, pointer interface{}, paramKeyToAttrMap map[string]string, priorityTag string,
) (err error) {
	if params == nil {
		// If `params` is nil, no conversion.
		return nil
	}
	if pointer == nil {
		return gerror.NewCode(gcode.CodeInvalidParameter, "object pointer cannot be nil")
	}

	// JSON content converting.
	ok, err := doConvertWithJsonCheck(params, pointer)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	defer func() {
		// Catch the panic, especially the reflection operation panics.
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodeSkipf(gcode.CodeInternalPanic, 1, "%+v", exception)
			}
		}
	}()

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
			return gerror.NewCodef(
				gcode.CodeInvalidParameter,
				"destination pointer should be type of '*struct', but got '%v'",
				pointerReflectKind,
			)
		}
		// Using IsNil on reflect.Ptr variable is OK.
		if !pointerReflectValue.IsValid() || pointerReflectValue.IsNil() {
			return gerror.NewCode(
				gcode.CodeInvalidParameter,
				"destination pointer cannot be nil",
			)
		}
		pointerElemReflectValue = pointerReflectValue.Elem()
	}

	// If `params` and `pointer` are the same type, the do directly assignment.
	// For performance enhancement purpose.
	if ok = doConvertWithTypeCheck(paramsReflectValue, pointerElemReflectValue); ok {
		return nil
	}

	// custom convert.
	if ok, err = callCustomConverter(paramsReflectValue, pointerReflectValue); ok {
		return err
	}

	// Normal unmarshalling interfaces checks.
	if ok, err = bindVarToReflectValueWithInterfaceCheck(pointerReflectValue, paramsInterface); ok {
		return err
	}

	// It automatically creates struct object if necessary.
	// For example, if `pointer` is **User, then `elem` is *User, which is a pointer to User.
	if pointerElemReflectValue.Kind() == reflect.Ptr {
		if !pointerElemReflectValue.IsValid() || pointerElemReflectValue.IsNil() {
			e := reflect.New(pointerElemReflectValue.Type().Elem())
			pointerElemReflectValue.Set(e)
			defer func() {
				if err != nil {
					// If it is converted failed, it reset the `pointer` to nil.
					pointerReflectValue.Elem().Set(reflect.Zero(pointerReflectValue.Type().Elem()))
				}
			}()
		}
		// if v, ok := pointerElemReflectValue.Interface().(iUnmarshalValue); ok {
		//	return v.UnmarshalValue(params)
		// }
		// Note that it's `pointerElemReflectValue` here not `pointerReflectValue`.
		if ok, err = bindVarToReflectValueWithInterfaceCheck(pointerElemReflectValue, paramsInterface); ok {
			return err
		}
		// Retrieve its element, may be struct at last.
		pointerElemReflectValue = pointerElemReflectValue.Elem()
	}

	// paramsMap is the map[string]interface{} type variable for params.
	// DO NOT use MapDeep here.
	paramsMap := doMapConvert(paramsInterface, recursiveTypeAuto, true)
	if paramsMap == nil {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`convert params from "%#v" to "map[string]interface{}" failed`,
			params,
		)
	}

	// Nothing to be done as the parameters are empty.
	if len(paramsMap) == 0 {
		return nil
	}

	// It only performs one converting to the same attribute.
	// doneMap is used to check repeated converting, its key is the real attribute name
	// of the struct.
	doneAttrMap := make(map[string]struct{})

	// The key of the attrMap is the attribute name of the struct,
	// and the value is its replaced name for later comparison to improve performance.
	var (
		tempName       string
		elemFieldType  reflect.StructField
		elemFieldValue reflect.Value
		elemType       = pointerElemReflectValue.Type()
		// Attribute name to its symbols-removed name,
		// in order to quick index and comparison in following logic.
		attrToCheckNameMap = make(map[string]string)
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
			if err = doStruct(paramsMap, elemFieldValue, paramKeyToAttrMap, priorityTag); err != nil {
				return err
			}
		} else {
			tempName = elemFieldType.Name
			attrToCheckNameMap[tempName] = utils.RemoveSymbols(tempName)
		}
	}
	if len(attrToCheckNameMap) == 0 {
		return nil
	}

	// The key of the `attrToTagCheckNameMap` is the attribute name of the struct,
	// and the value is its replaced tag name for later comparison to improve performance.
	var priorityTagArray []string
	if priorityTag != "" {
		priorityTagArray = append(utils.SplitAndTrim(priorityTag, ","), gtag.StructTagPriority...)
	} else {
		priorityTagArray = gtag.StructTagPriority
	}
	tagToAttrNameMap, err := gstructs.TagMapName(pointerElemReflectValue, priorityTagArray)
	if err != nil {
		return err
	}
	var toBeDeletedTagNames = make([]string, 0)
	for tagName, attributeName := range tagToAttrNameMap {
		// If there's something else in the tag string,
		// it uses the first part which is split using char ','.
		// Example:
		// orm:"id, priority"
		// orm:"name, with:uid=id"
		if array := strings.Split(tagName, ","); len(array) > 1 {
			toBeDeletedTagNames = append(toBeDeletedTagNames, tagName)
			tagName = array[0]
			tagToAttrNameMap[tagName] = attributeName
		}
		// If tag and attribute values both exist in `paramsMap`,
		// it then uses the tag value overwriting the attribute value in `paramsMap`.
		if paramsMap[tagName] != nil && paramsMap[attributeName] != nil {
			paramsMap[attributeName] = paramsMap[tagName]
		}
	}
	for _, tagName := range toBeDeletedTagNames {
		delete(tagToAttrNameMap, tagName)
	}

	// To convert value base on custom parameter key to attribute name map.
	err = doStructBaseOnParamKeyToAttrMap(
		pointerElemReflectValue,
		paramsMap,
		paramKeyToAttrMap,
		doneAttrMap,
	)
	if err != nil {
		return err
	}

	err = doStructBaseOnTagToAttrNameMap(
		pointerElemReflectValue,
		paramsMap,
		paramKeyToAttrMap,
		doneAttrMap,
		tagToAttrNameMap,
	)
	if err != nil {
		return err
	}

	// To convert value base on precise attribute name.
	err = doStructBaseOnAttrToCheckNameMap(
		pointerElemReflectValue,
		paramsMap,
		paramKeyToAttrMap,
		doneAttrMap,
		attrToCheckNameMap,
	)
	if err != nil {
		return err
	}

	// To convert value base on parameter map.
	err = doStructBaseOnParamMap(
		pointerElemReflectValue,
		paramsMap,
		paramKeyToAttrMap,
		doneAttrMap,
		attrToCheckNameMap,
		tagToAttrNameMap,
	)
	if err != nil {
		return err
	}

	return nil
}

func doStructBaseOnParamKeyToAttrMap(
	pointerElemReflectValue reflect.Value,
	paramsMap map[string]interface{},
	paramKeyToAttrMap map[string]string,
	doneAttrMap map[string]struct{},
) error {
	if len(paramKeyToAttrMap) == 0 {
		return nil
	}
	for paramKey, attrName := range paramKeyToAttrMap {
		paramValue, ok := paramsMap[paramKey]
		if !ok {
			continue
		}
		// If the attribute name is already checked converting, then skip it.
		if _, ok = doneAttrMap[attrName]; ok {
			continue
		}
		// Mark it done.
		doneAttrMap[attrName] = struct{}{}
		if err := bindVarToStructAttr(
			pointerElemReflectValue, attrName, paramValue, paramKeyToAttrMap,
		); err != nil {
			return err
		}
	}
	return nil
}

func doStructBaseOnAttrToCheckNameMap(
	pointerElemReflectValue reflect.Value,
	paramsMap map[string]interface{},
	paramKeyToAttrMap map[string]string,
	doneAttrMap map[string]struct{},
	attrToCheckNameMap map[string]string,
) error {
	for attrName := range attrToCheckNameMap {
		// The value by precise attribute name.
		paramValue, ok := paramsMap[attrName]
		if !ok {
			continue
		}
		// If the attribute name is already converted, it then skips it.
		if _, ok = doneAttrMap[attrName]; ok {
			continue
		}
		// Mark it done.
		doneAttrMap[attrName] = struct{}{}
		if err := bindVarToStructAttr(
			pointerElemReflectValue, attrName, paramValue, paramKeyToAttrMap,
		); err != nil {
			return err
		}
	}
	return nil
}

func doStructBaseOnTagToAttrNameMap(
	pointerElemReflectValue reflect.Value,
	paramsMap map[string]interface{},
	paramKeyToAttrMap map[string]string,
	doneAttrMap map[string]struct{},
	tagToAttrNameMap map[string]string,
) error {
	var (
		paramValue interface{}
		ok         bool
	)
	for tagName, attrName := range tagToAttrNameMap {
		// It firstly considers `paramName` as accurate tag name,
		// and retrieve attribute name from `tagToAttrNameMap` .
		paramValue, ok = paramsMap[tagName]
		if !ok {
			continue
		}
		// If the attribute name is already converted, it then skips it.
		if _, ok = doneAttrMap[attrName]; ok {
			continue
		}
		// Mark it done.
		doneAttrMap[attrName] = struct{}{}
		if err := bindVarToStructAttr(
			pointerElemReflectValue, attrName, paramValue, paramKeyToAttrMap,
		); err != nil {
			return err
		}
	}
	return nil
}

func doStructBaseOnParamMap(
	pointerElemReflectValue reflect.Value,
	paramsMap map[string]interface{},
	paramKeyToAttrMap map[string]string,
	doneAttrMap map[string]struct{},
	attrToCheckNameMap map[string]string,
	tagToAttrNameMap map[string]string,
) error {
	var (
		attrName                string
		paramNameWithoutSymbols string
		ok                      bool
	)
	for paramName, paramValue := range paramsMap {
		// It was already converted in previous procedure.
		if _, ok = tagToAttrNameMap[paramName]; ok {
			continue
		}
		// It was already converted in previous procedure.
		if _, ok = attrToCheckNameMap[paramName]; ok {
			continue
		}

		paramNameWithoutSymbols = utils.RemoveSymbols(paramName)
		// Matching the parameters to struct attributes.
		for attrKey, cmpKey := range attrToCheckNameMap {
			// Example:
			// UserName  eq user_name
			// User-Name eq username
			// username  eq userName
			// etc.
			if strings.EqualFold(paramNameWithoutSymbols, cmpKey) {
				attrName = attrKey
				break
			}
		}

		// No matching, it gives up this attribute converting.
		if attrName == "" {
			continue
		}
		// If the attribute name is already converted, it then skips it.
		if _, ok = doneAttrMap[attrName]; ok {
			continue
		}
		// Mark it done.
		doneAttrMap[attrName] = struct{}{}
		if err := bindVarToStructAttr(
			pointerElemReflectValue, attrName, paramValue, paramKeyToAttrMap,
		); err != nil {
			return err
		}
	}
	return nil
}

// bindVarToStructAttr sets value to struct object attribute by name.
func bindVarToStructAttr(
	structReflectValue reflect.Value,
	attrName string, value interface{}, paramKeyToAttrMap map[string]string,
) (err error) {
	structFieldValue := structReflectValue.FieldByName(attrName)
	if !structFieldValue.IsValid() {
		return nil
	}
	// CanSet checks whether attribute is public accessible.
	if !structFieldValue.CanSet() {
		return nil
	}
	defer func() {
		if exception := recover(); exception != nil {
			if err = bindVarToReflectValue(structFieldValue, value, paramKeyToAttrMap); err != nil {
				err = gerror.Wrapf(err, `error binding value to attribute "%s"`, attrName)
			}
		}
	}()
	// Directly converting.
	if empty.IsNil(value) {
		structFieldValue.Set(reflect.Zero(structFieldValue.Type()))
	} else {
		// Try to call custom converter.
		// Issue: https://github.com/gogf/gf/issues/3099
		var (
			customConverterInput reflect.Value
			ok                   bool
		)
		if customConverterInput, ok = value.(reflect.Value); !ok {
			customConverterInput = reflect.ValueOf(value)
		}

		if ok, err = callCustomConverter(customConverterInput, structFieldValue); ok || err != nil {
			return
		}

		// Special handling for certain types:
		// - Overwrite the default type converting logic of stdlib for time.Time/*time.Time.
		var structFieldTypeName = structFieldValue.Type().String()
		switch structFieldTypeName {
		case "time.Time", "*time.Time":
			doConvertWithReflectValueSet(structFieldValue, doConvertInput{
				FromValue:  value,
				ToTypeName: structFieldTypeName,
				ReferValue: structFieldValue,
			})
			return
		// Hold the time zone consistent in recursive
		// Issue: https://github.com/gogf/gf/issues/2980
		case "*gtime.Time", "gtime.Time":
			doConvertWithReflectValueSet(structFieldValue, doConvertInput{
				FromValue:  value,
				ToTypeName: structFieldTypeName,
				ReferValue: structFieldValue,
			})
			return
		}

		// Common interface check.
		if ok, err = bindVarToReflectValueWithInterfaceCheck(structFieldValue, value); ok {
			return err
		}

		// Default converting.
		doConvertWithReflectValueSet(structFieldValue, doConvertInput{
			FromValue:  value,
			ToTypeName: structFieldTypeName,
			ReferValue: structFieldValue,
		})
	}
	return nil
}

// bindVarToReflectValueWithInterfaceCheck does bind using common interfaces checks.
func bindVarToReflectValueWithInterfaceCheck(reflectValue reflect.Value, value interface{}) (bool, error) {
	var pointer interface{}
	if reflectValue.Kind() != reflect.Ptr && reflectValue.CanAddr() {
		reflectValueAddr := reflectValue.Addr()
		if reflectValueAddr.IsNil() || !reflectValueAddr.IsValid() {
			return false, nil
		}
		// Not a pointer, but can token address, that makes it can be unmarshalled.
		pointer = reflectValue.Addr().Interface()
	} else {
		if reflectValue.IsNil() || !reflectValue.IsValid() {
			return false, nil
		}
		pointer = reflectValue.Interface()
	}
	// UnmarshalValue.
	if v, ok := pointer.(iUnmarshalValue); ok {
		return ok, v.UnmarshalValue(value)
	}
	// UnmarshalText.
	if v, ok := pointer.(iUnmarshalText); ok {
		var valueBytes []byte
		if b, ok := value.([]byte); ok {
			valueBytes = b
		} else if s, ok := value.(string); ok {
			valueBytes = []byte(s)
		} else if f, ok := value.(iString); ok {
			valueBytes = []byte(f.String())
		}
		if len(valueBytes) > 0 {
			return ok, v.UnmarshalText(valueBytes)
		}
	}
	// UnmarshalJSON.
	if v, ok := pointer.(iUnmarshalJSON); ok {
		var valueBytes []byte
		if b, ok := value.([]byte); ok {
			valueBytes = b
		} else if s, ok := value.(string); ok {
			valueBytes = []byte(s)
		} else if f, ok := value.(iString); ok {
			valueBytes = []byte(f.String())
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
			return ok, v.UnmarshalJSON(valueBytes)
		}
	}
	if v, ok := pointer.(iSet); ok {
		v.Set(value)
		return ok, nil
	}
	return false, nil
}

// bindVarToReflectValue sets `value` to reflect value object `structFieldValue`.
func bindVarToReflectValue(
	structFieldValue reflect.Value, value interface{}, paramKeyToAttrMap map[string]string,
) (err error) {
	// JSON content converting.
	ok, err := doConvertWithJsonCheck(value, structFieldValue)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	kind := structFieldValue.Kind()
	// Converting using `Set` interface implements, for some types.
	switch kind {
	case reflect.Slice, reflect.Array, reflect.Ptr, reflect.Interface:
		if !structFieldValue.IsNil() {
			if v, ok := structFieldValue.Interface().(iSet); ok {
				v.Set(value)
				return nil
			}
		}
	}

	// Converting using reflection by kind.
	switch kind {
	case reflect.Map:
		return doMapToMap(value, structFieldValue, paramKeyToAttrMap)

	case reflect.Struct:
		// Recursively converting for struct attribute.
		if err = doStruct(value, structFieldValue, nil, ""); err != nil {
			// Note there's reflect conversion mechanism here.
			structFieldValue.Set(reflect.ValueOf(value).Convert(structFieldValue.Type()))
		}

	// Note that the slice element might be type of struct,
	// so it uses Struct function doing the converting internally.
	case reflect.Slice, reflect.Array:
		var (
			reflectArray reflect.Value
			reflectValue = reflect.ValueOf(value)
		)
		if reflectValue.Kind() == reflect.Slice || reflectValue.Kind() == reflect.Array {
			reflectArray = reflect.MakeSlice(structFieldValue.Type(), reflectValue.Len(), reflectValue.Len())
			if reflectValue.Len() > 0 {
				var (
					elemType     = reflectArray.Index(0).Type()
					elemTypeName string
					converted    bool
				)
				for i := 0; i < reflectValue.Len(); i++ {
					converted = false
					elemTypeName = elemType.Name()
					if elemTypeName == "" {
						elemTypeName = elemType.String()
					}
					var elem reflect.Value
					if elemType.Kind() == reflect.Ptr {
						elem = reflect.New(elemType.Elem()).Elem()
					} else {
						elem = reflect.New(elemType).Elem()
					}
					if elem.Kind() == reflect.Struct {
						if err = doStruct(reflectValue.Index(i).Interface(), elem, nil, ""); err == nil {
							converted = true
						}
					}
					if !converted {
						doConvertWithReflectValueSet(elem, doConvertInput{
							FromValue:  reflectValue.Index(i).Interface(),
							ToTypeName: elemTypeName,
							ReferValue: elem,
						})
					}
					if elemType.Kind() == reflect.Ptr {
						// Before it sets the `elem` to array, do pointer converting if necessary.
						elem = elem.Addr()
					}
					reflectArray.Index(i).Set(elem)
				}
			}
		} else {
			var (
				elem         reflect.Value
				elemType     = structFieldValue.Type().Elem()
				elemTypeName = elemType.Name()
				converted    bool
			)
			switch reflectValue.Kind() {
			case reflect.String:
				// Value is empty string.
				if reflectValue.IsZero() {
					var elemKind = elemType.Kind()
					// Try to find the original type kind of the slice element.
					if elemKind == reflect.Ptr {
						elemKind = elemType.Elem().Kind()
					}
					switch elemKind {
					case reflect.String:
						// Empty string cannot be assigned to string slice.
						return nil
					}
				}
			}
			if elemTypeName == "" {
				elemTypeName = elemType.String()
			}
			if elemType.Kind() == reflect.Ptr {
				elem = reflect.New(elemType.Elem()).Elem()
			} else {
				elem = reflect.New(elemType).Elem()
			}
			if elem.Kind() == reflect.Struct {
				if err = doStruct(value, elem, nil, ""); err == nil {
					converted = true
				}
			}
			if !converted {
				doConvertWithReflectValueSet(elem, doConvertInput{
					FromValue:  value,
					ToTypeName: elemTypeName,
					ReferValue: elem,
				})
			}
			if elemType.Kind() == reflect.Ptr {
				// Before it sets the `elem` to array, do pointer converting if necessary.
				elem = elem.Addr()
			}
			reflectArray = reflect.MakeSlice(structFieldValue.Type(), 1, 1)
			reflectArray.Index(0).Set(elem)
		}
		structFieldValue.Set(reflectArray)

	case reflect.Ptr:
		if structFieldValue.IsNil() || structFieldValue.IsZero() {
			// Nil or empty pointer, it creates a new one.
			item := reflect.New(structFieldValue.Type().Elem())
			if ok, err = bindVarToReflectValueWithInterfaceCheck(item, value); ok {
				structFieldValue.Set(item)
				return err
			}
			elem := item.Elem()
			if err = bindVarToReflectValue(elem, value, paramKeyToAttrMap); err == nil {
				structFieldValue.Set(elem.Addr())
			}
		} else {
			// Not empty pointer, it assigns values to it.
			return bindVarToReflectValue(structFieldValue.Elem(), value, paramKeyToAttrMap)
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
					gcode.CodeInternalPanic,
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
