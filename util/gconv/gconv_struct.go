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
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
	"github.com/gogf/gf/v2/util/gconv/internal/structcache"
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
// specified priorityTagAndFieldName for `params` key-value items to struct attribute names mapping.
// The parameter `priorityTag` supports multiple priorityTagAndFieldName that can be joined with char ','.
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
		// if v, ok := pointerElemReflectValue.Interface().(localinterface.IUnmarshalValue); ok {
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
	// Get struct info from cache or parse struct and cache the struct info.
	cachedStructInfo := structcache.GetCachedStructInfo(
		pointerElemReflectValue.Type(), priorityTag,
	)
	// Nothing to be converted.
	if cachedStructInfo == nil {
		return nil
	}
	// For the structure types of 0 tagOrFiledNameToFieldInfoMap,
	// they also need to be cached to prevent invalid logic
	if cachedStructInfo.HasNoFields() {
		return nil
	}
	var (
		// Indicates that those values have been used and cannot be reused.
		usedParamsKeyOrTagNameMap = structcache.GetUsedParamsKeyOrTagNameMapFromPool()
	)
	defer structcache.PutUsedParamsKeyOrTagNameMapToPool(usedParamsKeyOrTagNameMap)

	// Firstly, search according to custom mapping rules.
	// If a possible direct assignment is found, reduce the number of subsequent map searches.
	for paramKey, fieldName := range paramKeyToAttrMap {
		fieldInfo := cachedStructInfo.GetFieldInfo(fieldName)
		if fieldInfo != nil {
			if paramsValue, ok := paramsMap[paramKey]; ok {
				fieldValue := fieldInfo.GetFieldReflectValue(pointerElemReflectValue)
				if err = bindVarToStructField(
					fieldValue,
					paramsValue,
					fieldInfo,
					paramKeyToAttrMap,
				); err != nil {
					return err
				}
				if len(fieldInfo.OtherSameNameFieldIndex) > 0 {
					if err = setOtherSameNameField(
						fieldInfo, paramsValue, pointerReflectValue, paramKeyToAttrMap,
					); err != nil {
						return err
					}
				}
				usedParamsKeyOrTagNameMap[paramKey] = struct{}{}
			}
		}
	}
	// Already done converting for given `paramsMap`.
	if len(usedParamsKeyOrTagNameMap) == len(paramsMap) {
		return nil
	}
	// If the length of `paramsMap` is less than the number of fields, then loop based on `paramsMap`
	if len(paramsMap) < len(cachedStructInfo.FieldConvertInfos) {
		return bindStructWithLoopParamsMap(
			paramsMap, pointerElemReflectValue, paramKeyToAttrMap, usedParamsKeyOrTagNameMap, cachedStructInfo,
		)
	}
	return bindStructWithLoopFieldInfos(
		paramsMap, pointerElemReflectValue, paramKeyToAttrMap, usedParamsKeyOrTagNameMap, cachedStructInfo,
	)
}

func setOtherSameNameField(
	fieldInfo *structcache.CachedFieldInfo,
	srcValue any,
	structValue reflect.Value,
	paramKeyToAttrMap map[string]string,
) (err error) {
	// loop the same field name of all sub attributes.
	for i := range fieldInfo.OtherSameNameFieldIndex {
		fieldValue := fieldInfo.GetOtherFieldReflectValue(structValue, i)
		if err = bindVarToStructField(fieldValue, srcValue, fieldInfo, paramKeyToAttrMap); err != nil {
			return err
		}
	}
	return nil
}

func bindStructWithLoopParamsMap(
	paramsMap map[string]any,
	structValue reflect.Value,
	paramKeyToAttrMap map[string]string,
	usedParamsKeyOrTagNameMap map[string]struct{},
	cachedStructInfo *structcache.CachedStructInfo,
) (err error) {
	var (
		fieldName   string
		fieldInfo   *structcache.CachedFieldInfo
		fuzzLastKey string
		fieldValue  reflect.Value
		paramKey    string
		paramValue  any
		ok          bool
	)
	for paramKey, paramValue = range paramsMap {
		if _, ok = usedParamsKeyOrTagNameMap[paramKey]; ok {
			continue
		}
		fieldInfo = cachedStructInfo.GetFieldInfo(paramKey)
		if fieldInfo != nil {
			fieldName = fieldInfo.FieldName()
			// already converted using its field name?
			// the field name has the more priority than tag name.
			_, ok = usedParamsKeyOrTagNameMap[fieldName]
			if ok && fieldInfo.IsField {
				continue
			}
			fieldValue = fieldInfo.GetFieldReflectValue(structValue)
			if err = bindVarToStructField(
				fieldValue, paramValue, fieldInfo, paramKeyToAttrMap,
			); err != nil {
				return err
			}
			// handle same field name in nested struct.
			if len(fieldInfo.OtherSameNameFieldIndex) > 0 {
				if err = setOtherSameNameField(fieldInfo, paramValue, structValue, paramKeyToAttrMap); err != nil {
					return err
				}
			}
			usedParamsKeyOrTagNameMap[fieldName] = struct{}{}
			continue
		}

		// fuzzy matching.
		for _, fieldInfo = range cachedStructInfo.FieldConvertInfos {
			fieldName = fieldInfo.FieldName()
			if _, ok = usedParamsKeyOrTagNameMap[fieldName]; ok {
				continue
			}
			fuzzLastKey = fieldInfo.LastFuzzyKey.Load().(string)
			paramValue, ok = paramsMap[fuzzLastKey]
			if !ok {
				if strings.EqualFold(
					fieldInfo.RemoveSymbolsFieldName, utils.RemoveSymbols(paramKey),
				) {
					paramValue, ok = paramsMap[paramKey]
					// If it is found this time, update it based on what was not found last time.
					fieldInfo.LastFuzzyKey.Store(paramKey)
				}
			}
			if ok {
				fieldValue = fieldInfo.GetFieldReflectValue(structValue)
				if paramValue != nil {
					if err = bindVarToStructField(
						fieldValue, paramValue, fieldInfo, paramKeyToAttrMap,
					); err != nil {
						return err
					}
					// handle same field name in nested struct.
					if len(fieldInfo.OtherSameNameFieldIndex) > 0 {
						if err = setOtherSameNameField(
							fieldInfo, paramValue, structValue, paramKeyToAttrMap,
						); err != nil {
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

func bindStructWithLoopFieldInfos(
	paramsMap map[string]any,
	structValue reflect.Value,
	paramKeyToAttrMap map[string]string,
	usedParamsKeyOrTagNameMap map[string]struct{},
	cachedStructInfo *structcache.CachedStructInfo,
) (err error) {
	var (
		fieldInfo   *structcache.CachedFieldInfo
		fuzzLastKey string
		fieldValue  reflect.Value
		paramKey    string
		paramValue  any
		matched     bool
		ok          bool
	)
	for _, fieldInfo = range cachedStructInfo.FieldConvertInfos {
		for _, fieldTag := range fieldInfo.PriorityTagAndFieldName {
			if paramValue, ok = paramsMap[fieldTag]; !ok {
				continue
			}
			if _, ok = usedParamsKeyOrTagNameMap[fieldTag]; ok {
				matched = true
				break
			}
			fieldValue = fieldInfo.GetFieldReflectValue(structValue)
			if err = bindVarToStructField(
				fieldValue, paramValue, fieldInfo, paramKeyToAttrMap,
			); err != nil {
				return err
			}
			// handle same field name in nested struct.
			if len(fieldInfo.OtherSameNameFieldIndex) > 0 {
				if err = setOtherSameNameField(
					fieldInfo, paramValue, structValue, paramKeyToAttrMap,
				); err != nil {
					return err
				}
			}
			usedParamsKeyOrTagNameMap[fieldTag] = struct{}{}
			matched = true
			break
		}
		if matched {
			matched = false
			continue
		}

		fuzzLastKey = fieldInfo.LastFuzzyKey.Load().(string)
		if paramValue, ok = paramsMap[fuzzLastKey]; !ok {
			paramKey, paramValue = fuzzyMatchingFieldName(
				fieldInfo.RemoveSymbolsFieldName, paramsMap, usedParamsKeyOrTagNameMap,
			)
			ok = paramKey != ""
			fieldInfo.LastFuzzyKey.Store(paramKey)
		}
		if ok {
			fieldValue = fieldInfo.GetFieldReflectValue(structValue)
			if paramValue != nil {
				if err = bindVarToStructField(
					fieldValue, paramValue, fieldInfo, paramKeyToAttrMap,
				); err != nil {
					return err
				}
				// handle same field name in nested struct.
				if len(fieldInfo.OtherSameNameFieldIndex) > 0 {
					if err = setOtherSameNameField(
						fieldInfo, paramValue, structValue, paramKeyToAttrMap,
					); err != nil {
						return err
					}
				}
			}
			usedParamsKeyOrTagNameMap[paramKey] = struct{}{}
		}
	}
	return nil
}

// fuzzy matching rule:
// to match field name and param key in case-insensitive and without symbols.
func fuzzyMatchingFieldName(
	fieldName string,
	paramsMap map[string]any,
	usedParamsKeyMap map[string]struct{},
) (string, any) {
	for paramKey, paramValue := range paramsMap {
		if _, ok := usedParamsKeyMap[paramKey]; ok {
			continue
		}
		removeParamKeyUnderline := utils.RemoveSymbols(paramKey)
		if strings.EqualFold(fieldName, removeParamKeyUnderline) {
			return paramKey, paramValue
		}
	}
	return "", nil
}

// bindVarToStructField sets value to struct object attribute by name.
func bindVarToStructField(
	fieldValue reflect.Value,
	srcValue interface{},
	fieldInfo *structcache.CachedFieldInfo,
	paramKeyToAttrMap map[string]string,
) (err error) {
	if !fieldValue.IsValid() {
		return nil
	}
	// CanSet checks whether attribute is public accessible.
	if !fieldValue.CanSet() {
		return nil
	}
	defer func() {
		if exception := recover(); exception != nil {
			if err = bindVarToReflectValue(fieldValue, srcValue, paramKeyToAttrMap); err != nil {
				err = gerror.Wrapf(err, `error binding srcValue to attribute "%s"`, fieldInfo.FieldName())
			}
		}
	}()
	// Directly converting.
	if empty.IsNil(srcValue) {
		fieldValue.Set(reflect.Zero(fieldValue.Type()))
		return nil
	}
	// Try to call custom converter.
	// Issue: https://github.com/gogf/gf/issues/3099
	var (
		customConverterInput reflect.Value
		ok                   bool
	)
	if fieldInfo.IsCustomConvert {
		if customConverterInput, ok = srcValue.(reflect.Value); !ok {
			customConverterInput = reflect.ValueOf(srcValue)
		}
		if ok, err = callCustomConverter(customConverterInput, fieldValue); ok || err != nil {
			return
		}
	}
	if fieldInfo.IsCommonInterface {
		if ok, err = bindVarToReflectValueWithInterfaceCheck(fieldValue, srcValue); ok || err != nil {
			return
		}
	}
	// Common types use fast assignment logic
	if fieldInfo.ConvertFunc != nil {
		fieldInfo.ConvertFunc(srcValue, fieldValue)
		return nil
	}
	doConvertWithReflectValueSet(fieldValue, doConvertInput{
		FromValue:  srcValue,
		ToTypeName: fieldInfo.StructField.Type.String(),
		ReferValue: fieldValue,
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
	if v, ok := pointer.(localinterface.IUnmarshalValue); ok {
		return ok, v.UnmarshalValue(value)
	}
	// UnmarshalText.
	if v, ok := pointer.(localinterface.IUnmarshalText); ok {
		var valueBytes []byte
		if b, ok := value.([]byte); ok {
			valueBytes = b
		} else if s, ok := value.(string); ok {
			valueBytes = []byte(s)
		} else if f, ok := value.(localinterface.IString); ok {
			valueBytes = []byte(f.String())
		}
		if len(valueBytes) > 0 {
			return ok, v.UnmarshalText(valueBytes)
		}
	}
	// UnmarshalJSON.
	if v, ok := pointer.(localinterface.IUnmarshalJSON); ok {
		var valueBytes []byte
		if b, ok := value.([]byte); ok {
			valueBytes = b
		} else if s, ok := value.(string); ok {
			valueBytes = []byte(s)
		} else if f, ok := value.(localinterface.IString); ok {
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
	if v, ok := pointer.(localinterface.ISet); ok {
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
			if v, ok := structFieldValue.Interface().(localinterface.ISet); ok {
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
