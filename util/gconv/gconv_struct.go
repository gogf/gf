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
	paramsMap, ok := paramsInterface.(map[string]interface{})
	if !ok {
		// paramsMap is the map[string]interface{} type variable for params.
		// DO NOT use MapDeep here.
		paramsMap = doMapConvert(paramsInterface, recursiveTypeAuto, true)
		if paramsMap == nil {
			return gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`convert params from "%#v" to "map[string]interface{}" failed`,
				params,
			)
		}
	}
	// Nothing to be done as the parameters are empty.
	if len(paramsMap) == 0 {
		return nil
	}
	// parse struct
	convStructInfo := getConvStructInfo(pointerElemReflectValue.Type(), priorityTag)
	// Nothing to be converted.
	if convStructInfo == nil {
		return nil
	}
	// For the structure types of 0 fieldsMap,
	// they also need to be cached to prevent invalid logic
	if convStructInfo.NoFields() {
		return nil
	}
	var (
		// Indicates that those values have been used and cannot be reused.
		usedParamsKeyOrTagNameMap = poolGetUsedParamsKeyOrTagNameMap()
	)
	defer poolPutUsedParamsKeyOrTagNameMap(usedParamsKeyOrTagNameMap)

	setOtherField := func(fieldInfo *convertFieldInfo, srcValue any) error {
		for i := 0; i < len(fieldInfo.otherFieldIndex); i++ {
			fieldValue := fieldInfo.getOtherFieldReflectValue(pointerElemReflectValue, i)
			if err = bindVarToStructField(fieldValue, srcValue, fieldInfo, paramKeyToAttrMap); err != nil {
				return err
			}
		}
		return nil
	}

	// Firstly, search according to custom mapping rules.
	// If a possible direct assignment is found, reduce the number of subsequent map searches.
	for paramKey, fieldName := range paramKeyToAttrMap {
		fieldInfo := convStructInfo.GetFieldInfo(fieldName)
		if fieldInfo != nil {
			if paramsValue, ok := paramsMap[paramKey]; ok {
				fieldValue := fieldInfo.getFieldReflectValue(pointerElemReflectValue)
				if err = bindVarToStructField(
					fieldValue, paramsValue,
					fieldInfo, paramKeyToAttrMap); err != nil {
					return err
				}
				if len(fieldInfo.otherFieldIndex) > 0 {
					if err = setOtherField(fieldInfo, paramsValue); err != nil {
						return err
					}
				}
				usedParamsKeyOrTagNameMap[paramKey] = struct{}{}
			}
		}
	}
	if len(usedParamsKeyOrTagNameMap) == len(paramsMap) {
		return nil
	}
	var (
		fieldName   string
		fieldInfo   *convertFieldInfo
		fuzzLastKey string
		fieldValue  reflect.Value
		paramKey    string
		paramValue  any
	)
	for paramKey, paramValue = range paramsMap {
		// // We cannot make such a judgment to exit early
		// Assuming the following situation
		// All other fields have been matched.
		// At this point, we can exit by matching another one.
		// If we are not lucky, we will match [Limit], but what we want is [per_page]
		//  Limit int `json:"per_page"`
		// 	paramsMap = map{
		//		Limit:1,
		//  	per_page:2.
		//	}
		// if len(usedParamsKeyOrTagNameMap) >= convStructInfo.fieldsCount {
		// 	break
		// }
		fieldInfo = convStructInfo.GetFieldInfo(paramKey)
		if fieldInfo != nil {
			fieldName = fieldInfo.FieldName()
			_, ok = usedParamsKeyOrTagNameMap[fieldName]
			if ok && fieldInfo.isField {
				//  Limit int `json:"per_page"`
				//  Map has both [Limit] and [per_page] simultaneously
				//  [per_page] is required as the standard
				// 	paramsMap = map{
				//		Limit:1,
				//  	per_page:2.
				//	}
				continue
			}
			fieldValue = fieldInfo.getFieldReflectValue(pointerElemReflectValue)
			if err = bindVarToStructField(fieldValue, paramValue,
				fieldInfo, paramKeyToAttrMap); err != nil {
				return err
			}
			if len(fieldInfo.otherFieldIndex) > 0 {
				if err = setOtherField(fieldInfo, paramValue); err != nil {
					return err
				}
			}
			usedParamsKeyOrTagNameMap[fieldName] = struct{}{}
			continue
		}

		for fieldName, fieldInfo = range convStructInfo.fieldsMap {
			// Fuzzy matching only matches fields
			if fieldInfo.isField == false {
				continue
			}
			_, ok = usedParamsKeyOrTagNameMap[fieldName]
			if ok {
				continue
			}
			fuzzLastKey = fieldInfo.lastFuzzKey.Load().(string)
			paramValue, ok = paramsMap[fuzzLastKey]
			if !ok {
				if strings.EqualFold(fieldInfo.removeSymbolsFieldName, utils.RemoveSymbols(paramKey)) {
					paramValue, ok = paramsMap[paramKey]
					// If it is found this time, update it based on what was not found last time
					fieldInfo.lastFuzzKey.Store(paramKey)
				}
			}
			if ok {
				// If paramValue is nil, you don't need to set any value
				// Because getFieldReflectValue is already initialized
				fieldValue = fieldInfo.getFieldReflectValue(pointerElemReflectValue)
				if paramValue != nil {
					if err = bindVarToStructField(fieldValue, paramValue,
						fieldInfo, paramKeyToAttrMap); err != nil {
						return err
					}
					if len(fieldInfo.otherFieldIndex) > 0 {
						if err = setOtherField(fieldInfo, paramValue); err != nil {
							return err
						}
					}
				}
				usedParamsKeyOrTagNameMap[fieldInfo.FieldName()] = struct{}{}
				break
			}
		}
	}
	return nil
}

// bindVarToStructField sets value to struct object attribute by name.
func bindVarToStructField(
	structFieldValue reflect.Value, srcValue interface{},
	fieldInfo *convertFieldInfo, paramKeyToAttrMap map[string]string,
) (err error) {
	if !structFieldValue.IsValid() {
		return nil
	}
	// CanSet checks whether attribute is public accessible.
	if !structFieldValue.CanSet() {
		return nil
	}
	defer func() {
		if exception := recover(); exception != nil {
			if err = bindVarToReflectValue(structFieldValue, srcValue, paramKeyToAttrMap); err != nil {
				err = gerror.Wrapf(err, `error binding srcValue to attribute "%s"`, fieldInfo.FieldName())
			}
		}
	}()
	// Directly converting.
	if empty.IsNil(srcValue) {
		structFieldValue.Set(reflect.Zero(structFieldValue.Type()))
		return nil
	}
	// Try to call custom converter.
	// Issue: https://github.com/gogf/gf/issues/3099
	var (
		customConverterInput reflect.Value
		ok                   bool
	)
	if fieldInfo.isCustomConvert {
		if customConverterInput, ok = srcValue.(reflect.Value); !ok {
			customConverterInput = reflect.ValueOf(srcValue)
		}
		if ok, err = callCustomConverter(customConverterInput, structFieldValue); ok || err != nil {
			return
		}
	}
	if fieldInfo.isCommonInterface {
		if ok, err = bindVarToReflectValueWithInterfaceCheck(structFieldValue, srcValue); ok || err != nil {
			return
		}
	}
	// Common types use fast assignment logic
	if fieldInfo.convFunc != nil {
		fieldInfo.convFunc(srcValue, structFieldValue)
		return nil
	}
	doConvertWithReflectValueSet(structFieldValue, doConvertInput{
		FromValue:  srcValue,
		ToTypeName: fieldInfo.structField.Type.String(),
		ReferValue: structFieldValue,
	})
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
