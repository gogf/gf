// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package converter

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

// StructOption is the option for Struct converting.
type StructOption struct {
	// ParamKeyToAttrMap is the map for custom parameter key to attribute name mapping.
	ParamKeyToAttrMap map[string]string

	// PriorityTag is the priority tag for struct converting.
	PriorityTag string

	// ContinueOnError specifies whether to continue converting the next element
	// if one element converting fails.
	ContinueOnError bool
}

func (c *Converter) getStructOption(option ...StructOption) StructOption {
	if len(option) > 0 {
		return option[0]
	}
	return StructOption{}
}

// Struct is the core internal converting function for any data to struct.
func (c *Converter) Struct(params, pointer any, option ...StructOption) (err error) {
	if params == nil {
		// If `params` is nil, no conversion.
		return nil
	}
	if pointer == nil {
		return gerror.NewCode(gcode.CodeInvalidParameter, "object pointer cannot be nil")
	}

	// JSON content converting.
	ok, err := c.doConvertWithJsonCheck(params, pointer)
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
		structOption            = c.getStructOption(option...)
		paramsReflectValue      reflect.Value
		paramsInterface         any // DO NOT use `params` directly as it might be type `reflect.Value`
		pointerReflectValue     reflect.Value
		pointerReflectKind      reflect.Kind
		pointerElemReflectValue reflect.Value // The reflection value to struct element.
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
	if ok = c.doConvertWithTypeCheck(paramsReflectValue, pointerElemReflectValue); ok {
		return nil
	}

	// custom convert.
	ok, err = c.callCustomConverter(paramsReflectValue, pointerReflectValue)
	if err != nil && !structOption.ContinueOnError {
		return err
	}
	if ok {
		return nil
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
	paramsMap, ok := paramsInterface.(map[string]any)
	if !ok {
		// paramsMap is the map[string]any type variable for params.
		// DO NOT use MapDeep here.
		paramsMap, err = c.doMapConvert(paramsInterface, RecursiveTypeAuto, true, MapOption{
			ContinueOnError: structOption.ContinueOnError,
		})
		if err != nil {
			return err
		}
		if paramsMap == nil {
			return gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`convert params from "%#v" to "map[string]any" failed`,
				params,
			)
		}
	}
	// Nothing to be done as the parameters are empty.
	if len(paramsMap) == 0 {
		return nil
	}
	// Get struct info from cache or parse struct and cache the struct info.
	cachedStructInfo := c.internalConverter.GetCachedStructInfo(
		pointerElemReflectValue.Type(), structOption.PriorityTag,
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
		cachedFieldInfo           *structcache.CachedFieldInfo
		paramsValue               any
	)
	defer structcache.PutUsedParamsKeyOrTagNameMapToPool(usedParamsKeyOrTagNameMap)

	// Firstly, search according to custom mapping rules.
	// If a possible direct assignment is found, reduce the number of subsequent map searches.
	for paramKey, fieldName := range structOption.ParamKeyToAttrMap {
		paramsValue, ok = paramsMap[paramKey]
		if !ok {
			continue
		}
		cachedFieldInfo = cachedStructInfo.GetFieldInfo(fieldName)
		if cachedFieldInfo != nil {
			fieldValue := cachedFieldInfo.GetFieldReflectValueFrom(pointerElemReflectValue)
			if err = c.bindVarToStructField(
				cachedFieldInfo,
				fieldValue,
				paramsValue,
				structOption,
			); err != nil {
				return err
			}
			if len(cachedFieldInfo.OtherSameNameField) > 0 {
				if err = c.setOtherSameNameField(
					cachedFieldInfo, paramsValue, pointerReflectValue, structOption,
				); err != nil {
					return err
				}
			}
			usedParamsKeyOrTagNameMap[paramKey] = struct{}{}
		}
	}
	// Already done converting for given `paramsMap`.
	if len(usedParamsKeyOrTagNameMap) == len(paramsMap) {
		return nil
	}
	return c.bindStructWithLoopFieldInfos(
		paramsMap, pointerElemReflectValue,
		usedParamsKeyOrTagNameMap, cachedStructInfo,
		structOption,
	)
}

func (c *Converter) setOtherSameNameField(
	cachedFieldInfo *structcache.CachedFieldInfo,
	srcValue any,
	structValue reflect.Value,
	option StructOption,
) (err error) {
	// loop the same field name of all sub attributes.
	for _, otherFieldInfo := range cachedFieldInfo.OtherSameNameField {
		fieldValue := cachedFieldInfo.GetOtherFieldReflectValueFrom(structValue, otherFieldInfo.FieldIndexes)
		if err = c.bindVarToStructField(otherFieldInfo, fieldValue, srcValue, option); err != nil {
			return err
		}
	}
	return nil
}

func (c *Converter) bindStructWithLoopFieldInfos(
	paramsMap map[string]any,
	structValue reflect.Value,
	usedParamsKeyOrTagNameMap map[string]struct{},
	cachedStructInfo *structcache.CachedStructInfo,
	option StructOption,
) (err error) {
	var (
		cachedFieldInfo *structcache.CachedFieldInfo
		fuzzLastKey     string
		fieldValue      reflect.Value
		paramKey        string
		paramValue      any
		matched         bool
		ok              bool
	)
	for _, cachedFieldInfo = range cachedStructInfo.GetFieldConvertInfos() {
		for _, fieldTag := range cachedFieldInfo.PriorityTagAndFieldName {
			if paramValue, ok = paramsMap[fieldTag]; !ok {
				continue
			}
			fieldValue = cachedFieldInfo.GetFieldReflectValueFrom(structValue)
			if err = c.bindVarToStructField(
				cachedFieldInfo, fieldValue, paramValue, option,
			); err != nil && !option.ContinueOnError {
				return err
			}
			// handle same field name in nested struct.
			if len(cachedFieldInfo.OtherSameNameField) > 0 {
				if err = c.setOtherSameNameField(
					cachedFieldInfo, paramValue, structValue, option,
				); err != nil && !option.ContinueOnError {
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

		fuzzLastKey = cachedFieldInfo.LastFuzzyKey.Load().(string)
		if paramValue, ok = paramsMap[fuzzLastKey]; !ok {
			paramKey, paramValue = fuzzyMatchingFieldName(
				cachedFieldInfo.RemoveSymbolsFieldName, paramsMap, usedParamsKeyOrTagNameMap,
			)
			ok = paramKey != ""
			cachedFieldInfo.LastFuzzyKey.Store(paramKey)
		}
		if ok {
			fieldValue = cachedFieldInfo.GetFieldReflectValueFrom(structValue)
			if paramValue != nil {
				if err = c.bindVarToStructField(
					cachedFieldInfo, fieldValue, paramValue, option,
				); err != nil && !option.ContinueOnError {
					return err
				}
				// handle same field name in nested struct.
				if len(cachedFieldInfo.OtherSameNameField) > 0 {
					if err = c.setOtherSameNameField(
						cachedFieldInfo, paramValue, structValue, option,
					); err != nil && !option.ContinueOnError {
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
// each value to attribute converting comes into in this function.
func (c *Converter) bindVarToStructField(
	cachedFieldInfo *structcache.CachedFieldInfo,
	fieldValue reflect.Value,
	srcValue any,
	option StructOption,
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
			if err = c.bindVarToReflectValue(fieldValue, srcValue, option); err != nil {
				err = gerror.Wrapf(err, `error binding srcValue to attribute "%s"`, cachedFieldInfo.FieldName())
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
	if cachedFieldInfo.HasCustomConvert {
		if customConverterInput, ok = srcValue.(reflect.Value); !ok {
			customConverterInput = reflect.ValueOf(srcValue)
		}
		if ok, err = c.callCustomConverter(customConverterInput, fieldValue); ok || err != nil {
			return
		}
	}
	if cachedFieldInfo.IsCommonInterface {
		if ok, err = bindVarToReflectValueWithInterfaceCheck(fieldValue, srcValue); ok || err != nil {
			return
		}
	}
	// Common types use fast assignment logic
	if cachedFieldInfo.ConvertFunc != nil {
		return cachedFieldInfo.ConvertFunc(srcValue, fieldValue)
	}
	convertOption := ConvertOption{
		StructOption: option,
		SliceOption:  SliceOption{ContinueOnError: option.ContinueOnError},
		MapOption:    MapOption{ContinueOnError: option.ContinueOnError},
	}
	err = c.doConvertWithReflectValueSet(
		fieldValue, doConvertInput{
			FromValue:  srcValue,
			ToTypeName: cachedFieldInfo.StructField.Type.String(),
			ReferValue: fieldValue,
		},
		convertOption,
	)
	return err
}

// bindVarToReflectValueWithInterfaceCheck does bind using common interfaces checks.
func bindVarToReflectValueWithInterfaceCheck(reflectValue reflect.Value, value any) (bool, error) {
	var pointer any
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
func (c *Converter) bindVarToReflectValue(structFieldValue reflect.Value, value any, option StructOption) (err error) {
	// JSON content converting.
	ok, err := c.doConvertWithJsonCheck(value, structFieldValue)
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
	default:
	}

	// Converting using reflection by kind.
	switch kind {
	case reflect.Map:
		return c.MapToMap(value, structFieldValue, option.ParamKeyToAttrMap, MapOption{
			ContinueOnError: option.ContinueOnError,
		})

	case reflect.Struct:
		// Recursively converting for struct attribute.
		if err = c.Struct(value, structFieldValue, option); err != nil {
			// Note there's reflect conversion mechanism here.
			structFieldValue.Set(reflect.ValueOf(value).Convert(structFieldValue.Type()))
		}

	// Note that the slice element might be type of struct,
	// so it uses Struct function doing the converting internally.
	case reflect.Slice, reflect.Array:
		var (
			reflectArray  reflect.Value
			reflectValue  = reflect.ValueOf(value)
			convertOption = ConvertOption{
				StructOption: option,
				SliceOption:  SliceOption{ContinueOnError: option.ContinueOnError},
				MapOption:    MapOption{ContinueOnError: option.ContinueOnError},
			}
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
						if err = c.Struct(reflectValue.Index(i).Interface(), elem, option); err == nil {
							converted = true
						}
					}
					if !converted {
						err = c.doConvertWithReflectValueSet(
							elem, doConvertInput{
								FromValue:  reflectValue.Index(i).Interface(),
								ToTypeName: elemTypeName,
								ReferValue: elem,
							},
							convertOption,
						)
						if err != nil {
							return err
						}
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
					default:
					}
				}
			default:
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
				if err = c.Struct(value, elem, option); err == nil {
					converted = true
				}
			}
			if !converted {
				err = c.doConvertWithReflectValueSet(
					elem, doConvertInput{
						FromValue:  value,
						ToTypeName: elemTypeName,
						ReferValue: elem,
					},
					convertOption,
				)
				if err != nil {
					return err
				}
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
			if err = c.bindVarToReflectValue(elem, value, option); err == nil {
				structFieldValue.Set(elem.Addr())
			}
		} else {
			// Not empty pointer, it assigns values to it.
			return c.bindVarToReflectValue(structFieldValue.Elem(), value, option)
		}

	// It mainly and specially handles the interface of nil value.
	case reflect.Interface:
		if value == nil {
			// Specially.
			structFieldValue.Set(reflect.ValueOf((*any)(nil)))
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
