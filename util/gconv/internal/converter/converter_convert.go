// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package converter

import (
	"reflect"
	"time"

	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/os/gtime"
)

// ConvertOption is the option for converting.
type ConvertOption struct {
	// ExtraParams are extra values for implementing the converting.
	ExtraParams  []any
	SliceOption  SliceOption
	MapOption    MapOption
	StructOption StructOption
}

func (c *Converter) getConvertOption(option ...ConvertOption) ConvertOption {
	if len(option) > 0 {
		return option[0]
	}
	return ConvertOption{}
}

// ConvertWithTypeName converts the variable `fromValue` to the type `toTypeName`, the type `toTypeName` is specified by string.
func (c *Converter) ConvertWithTypeName(fromValue any, toTypeName string, option ...ConvertOption) (any, error) {
	return c.doConvert(
		doConvertInput{
			FromValue:  fromValue,
			ToTypeName: toTypeName,
			ReferValue: nil,
		},
		c.getConvertOption(option...),
	)
}

// ConvertWithRefer converts the variable `fromValue` to the type referred by value `referValue`.
func (c *Converter) ConvertWithRefer(fromValue, referValue any, option ...ConvertOption) (any, error) {
	var referValueRf reflect.Value
	if v, ok := referValue.(reflect.Value); ok {
		referValueRf = v
	} else {
		referValueRf = reflect.ValueOf(referValue)
	}
	return c.doConvert(
		doConvertInput{
			FromValue:  fromValue,
			ToTypeName: referValueRf.Type().String(),
			ReferValue: referValue,
		},
		c.getConvertOption(option...),
	)
}

type doConvertInput struct {
	FromValue  any    // Value that is converted from.
	ToTypeName string // Target value type name in string.
	ReferValue any    // Referred value, a value in type `ToTypeName`. Note that its type might be reflect.Value.

	// Marks that the value is already converted and set to `ReferValue`. Caller can ignore the returned result.
	// It is an attribute for internal usage purpose.
	alreadySetToReferValue bool
}

// doConvert does commonly use types converting.
func (c *Converter) doConvert(in doConvertInput, option ConvertOption) (convertedValue any, err error) {
	switch in.ToTypeName {
	case "int":
		return c.Int(in.FromValue)
	case "*int":
		if _, ok := in.FromValue.(*int); ok {
			return in.FromValue, nil
		}
		v, err := c.Int(in.FromValue)
		if err != nil {
			return nil, err
		}
		return &v, nil

	case "int8":
		return c.Int8(in.FromValue)
	case "*int8":
		if _, ok := in.FromValue.(*int8); ok {
			return in.FromValue, nil
		}
		v, err := c.Int8(in.FromValue)
		if err != nil {
			return nil, err
		}
		return &v, nil

	case "int16":
		return c.Int16(in.FromValue)
	case "*int16":
		if _, ok := in.FromValue.(*int16); ok {
			return in.FromValue, nil
		}
		v, err := c.Int16(in.FromValue)
		if err != nil {
			return nil, err
		}
		return &v, nil

	case "int32":
		return c.Int32(in.FromValue)
	case "*int32":
		if _, ok := in.FromValue.(*int32); ok {
			return in.FromValue, nil
		}
		v, err := c.Int32(in.FromValue)
		if err != nil {
			return nil, err
		}
		return &v, nil

	case "int64":
		return c.Int64(in.FromValue)
	case "*int64":
		if _, ok := in.FromValue.(*int64); ok {
			return in.FromValue, nil
		}
		v, err := c.Int64(in.FromValue)
		if err != nil {
			return nil, err
		}
		return &v, nil

	case "uint":
		return c.Uint(in.FromValue)
	case "*uint":
		if _, ok := in.FromValue.(*uint); ok {
			return in.FromValue, nil
		}
		v, err := c.Uint(in.FromValue)
		if err != nil {
			return nil, err
		}
		return &v, nil

	case "uint8":
		return c.Uint8(in.FromValue)
	case "*uint8":
		if _, ok := in.FromValue.(*uint8); ok {
			return in.FromValue, nil
		}
		v, err := c.Uint8(in.FromValue)
		if err != nil {
			return nil, err
		}
		return &v, nil

	case "uint16":
		return c.Uint16(in.FromValue)
	case "*uint16":
		if _, ok := in.FromValue.(*uint16); ok {
			return in.FromValue, nil
		}
		v, err := c.Uint16(in.FromValue)
		if err != nil {
			return nil, err
		}
		return &v, nil

	case "uint32":
		return c.Uint32(in.FromValue)
	case "*uint32":
		if _, ok := in.FromValue.(*uint32); ok {
			return in.FromValue, nil
		}
		v, err := c.Uint32(in.FromValue)
		if err != nil {
			return nil, err
		}
		return &v, nil

	case "uint64":
		return c.Uint64(in.FromValue)
	case "*uint64":
		if _, ok := in.FromValue.(*uint64); ok {
			return in.FromValue, nil
		}
		v, err := c.Uint64(in.FromValue)
		if err != nil {
			return nil, err
		}
		return &v, nil

	case "float32":
		return c.Float32(in.FromValue)
	case "*float32":
		if _, ok := in.FromValue.(*float32); ok {
			return in.FromValue, nil
		}
		v, err := c.Float32(in.FromValue)
		if err != nil {
			return nil, err
		}
		return &v, nil

	case "float64":
		return c.Float64(in.FromValue)
	case "*float64":
		if _, ok := in.FromValue.(*float64); ok {
			return in.FromValue, nil
		}
		v, err := c.Float64(in.FromValue)
		if err != nil {
			return nil, err
		}
		return &v, nil

	case "bool":
		return c.Bool(in.FromValue)
	case "*bool":
		if _, ok := in.FromValue.(*bool); ok {
			return in.FromValue, nil
		}
		v, err := c.Bool(in.FromValue)
		if err != nil {
			return nil, err
		}
		return &v, nil

	case "string":
		return c.String(in.FromValue)
	case "*string":
		if _, ok := in.FromValue.(*string); ok {
			return in.FromValue, nil
		}
		v, err := c.String(in.FromValue)
		if err != nil {
			return nil, err
		}
		return &v, nil

	case "[]byte":
		return c.Bytes(in.FromValue)
	case "[]int":
		return c.SliceInt(in.FromValue, option.SliceOption)
	case "[]int32":
		return c.SliceInt32(in.FromValue, option.SliceOption)
	case "[]int64":
		return c.SliceInt64(in.FromValue, option.SliceOption)
	case "[]uint":
		return c.SliceUint(in.FromValue, option.SliceOption)
	case "[]uint8":
		return c.Bytes(in.FromValue)
	case "[]uint32":
		return c.SliceUint32(in.FromValue, option.SliceOption)
	case "[]uint64":
		return c.SliceUint64(in.FromValue, option.SliceOption)
	case "[]float32":
		return c.SliceFloat32(in.FromValue, option.SliceOption)
	case "[]float64":
		return c.SliceFloat64(in.FromValue, option.SliceOption)
	case "[]string":
		return c.SliceStr(in.FromValue, option.SliceOption)

	case "Time", "time.Time":
		if len(option.ExtraParams) > 0 {
			s, err := c.String(option.ExtraParams[0])
			if err != nil {
				return nil, err
			}
			return c.Time(in.FromValue, s)
		}
		return c.Time(in.FromValue)
	case "*time.Time":
		var v time.Time
		if len(option.ExtraParams) > 0 {
			s, err := c.String(option.ExtraParams[0])
			if err != nil {
				return time.Time{}, err
			}
			v, err = c.Time(in.FromValue, s)
			if err != nil {
				return time.Time{}, err
			}
		} else {
			if _, ok := in.FromValue.(*time.Time); ok {
				return in.FromValue, nil
			}
			v, err = c.Time(in.FromValue)
			if err != nil {
				return time.Time{}, err
			}
		}
		return &v, nil

	case "GTime", "gtime.Time":
		if len(option.ExtraParams) > 0 {
			s, err := c.String(option.ExtraParams[0])
			if err != nil {
				return *gtime.New(), err
			}
			v, err := c.GTime(in.FromValue, s)
			if err != nil {
				return *gtime.New(), err
			}
			if v != nil {
				return *v, nil
			}
			return *gtime.New(), nil
		}
		v, err := c.GTime(in.FromValue)
		if err != nil {
			return *gtime.New(), err
		}
		if v != nil {
			return *v, nil
		}
		return *gtime.New(), nil
	case "*gtime.Time":
		if len(option.ExtraParams) > 0 {
			s, err := c.String(option.ExtraParams[0])
			if err != nil {
				return gtime.New(), err
			}
			v, err := c.GTime(in.FromValue, s)
			if err != nil {
				return gtime.New(), err
			}
			if v != nil {
				return v, nil
			}
			return gtime.New(), nil
		}
		v, err := c.GTime(in.FromValue)
		if err != nil {
			return gtime.New(), err
		}
		if v != nil {
			return v, nil
		}
		return gtime.New(), nil

	case "Duration", "time.Duration":
		return c.Duration(in.FromValue)
	case "*time.Duration":
		if _, ok := in.FromValue.(*time.Duration); ok {
			return in.FromValue, nil
		}
		v, err := c.Duration(in.FromValue)
		if err != nil {
			return nil, err
		}
		return &v, nil

	case "map[string]string":
		return c.MapStrStr(in.FromValue, option.MapOption)

	case "map[string]interface {}":
		return c.Map(in.FromValue, option.MapOption)

	case "[]map[string]interface {}":
		return c.SliceMap(in.FromValue, SliceMapOption{
			SliceOption: option.SliceOption,
			MapOption:   option.MapOption,
		})

	case "RawMessage", "json.RawMessage":
		// issue 3449
		bytes, err := json.Marshal(in.FromValue)
		if err != nil {
			return nil, err
		}
		return bytes, nil

	default:
		return c.doConvertForDefault(in, option)
	}
}

func (c *Converter) doConvertForDefault(in doConvertInput, option ConvertOption) (convertedValue any, err error) {
	if in.ReferValue != nil {
		var referReflectValue reflect.Value
		if v, ok := in.ReferValue.(reflect.Value); ok {
			referReflectValue = v
		} else {
			referReflectValue = reflect.ValueOf(in.ReferValue)
		}
		var fromReflectValue reflect.Value
		if v, ok := in.FromValue.(reflect.Value); ok {
			fromReflectValue = v
		} else {
			fromReflectValue = reflect.ValueOf(in.FromValue)
		}

		// custom converter.
		dstReflectValue, ok, err := c.callCustomConverterWithRefer(fromReflectValue, referReflectValue)
		if err != nil {
			return nil, err
		}
		if ok {
			return dstReflectValue.Interface(), nil
		}

		defer func() {
			if recover() != nil {
				in.alreadySetToReferValue = false
				if err = c.bindVarToReflectValue(referReflectValue, in.FromValue, option.StructOption); err == nil {
					in.alreadySetToReferValue = true
					convertedValue = referReflectValue.Interface()
				}
			}
		}()
		switch referReflectValue.Kind() {
		case reflect.Ptr:
			// Type converting for custom type pointers.
			// Eg:
			// type PayMode int
			// type Req struct{
			//     Mode *PayMode
			// }
			//
			// Struct(`{"Mode": 1000}`, &req)
			originType := referReflectValue.Type().Elem()
			switch originType.Kind() {
			case reflect.Struct:
				// Not support some kinds.
			default:
				in.ToTypeName = originType.Kind().String()
				in.ReferValue = nil
				result, err := c.doConvert(in, option)
				if err != nil {
					return nil, err
				}
				refElementValue := reflect.ValueOf(result)
				originTypeValue := reflect.New(refElementValue.Type()).Elem()
				originTypeValue.Set(refElementValue)
				in.alreadySetToReferValue = true
				return originTypeValue.Addr().Convert(referReflectValue.Type()).Interface(), nil
			}

		case reflect.Map:
			var targetValue = reflect.New(referReflectValue.Type()).Elem()
			if err = c.MapToMap(in.FromValue, targetValue, nil, option.MapOption); err == nil {
				in.alreadySetToReferValue = true
			}
			return targetValue.Interface(), nil

		default:
		}
		in.ToTypeName = referReflectValue.Kind().String()
		in.ReferValue = nil
		in.alreadySetToReferValue = true
		result, err := c.doConvert(in, option)
		if err != nil {
			return nil, err
		}
		convertedValue = reflect.ValueOf(result).Convert(referReflectValue.Type()).Interface()
		return convertedValue, nil
	}
	return in.FromValue, nil
}

func (c *Converter) doConvertWithReflectValueSet(reflectValue reflect.Value, in doConvertInput, option ConvertOption) error {
	convertedValue, err := c.doConvert(in, option)
	if err != nil {
		return err
	}
	if !in.alreadySetToReferValue {
		reflectValue.Set(reflect.ValueOf(convertedValue))
	}
	return err
}

// callCustomConverter call the custom converter. It will try some possible type.
func (c *Converter) callCustomConverter(srcReflectValue, dstReflectValue reflect.Value) (converted bool, err error) {
	// search type converter function.
	registeredConverterFunc, srcType, ok := c.getRegisteredTypeConverterFuncAndSrcType(srcReflectValue, dstReflectValue)
	if ok {
		return c.doCallCustomTypeConverter(srcReflectValue, dstReflectValue, registeredConverterFunc, srcType)
	}

	// search any converter function.
	anyConverterFunc := c.getRegisteredAnyConverterFunc(dstReflectValue)
	if anyConverterFunc == nil {
		return false, nil
	}
	err = anyConverterFunc(srcReflectValue.Interface(), dstReflectValue)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *Converter) callCustomConverterWithRefer(
	srcReflectValue, referReflectValue reflect.Value,
) (dstReflectValue reflect.Value, converted bool, err error) {
	// search type converter function.
	registeredConverterFunc, srcType, ok := c.getRegisteredTypeConverterFuncAndSrcType(
		srcReflectValue, referReflectValue,
	)
	if ok {
		dstReflectValue = reflect.New(referReflectValue.Type()).Elem()
		converted, err = c.doCallCustomTypeConverter(srcReflectValue, dstReflectValue, registeredConverterFunc, srcType)
		return
	}

	// search any converter function.
	anyConverterFunc := c.getRegisteredAnyConverterFunc(referReflectValue)
	if anyConverterFunc == nil {
		return reflect.Value{}, false, nil
	}
	dstReflectValue = reflect.New(referReflectValue.Type()).Elem()
	err = anyConverterFunc(srcReflectValue.Interface(), dstReflectValue)
	if err != nil {
		return reflect.Value{}, false, err
	}
	return dstReflectValue, true, nil
}

func (c *Converter) getRegisteredTypeConverterFuncAndSrcType(
	srcReflectValue, dstReflectValueForRefer reflect.Value,
) (f converterFunc, srcType reflect.Type, ok bool) {
	if len(c.typeConverterFuncMap) == 0 {
		return reflect.Value{}, nil, false
	}
	srcType = srcReflectValue.Type()
	for srcType.Kind() == reflect.Pointer {
		srcType = srcType.Elem()
	}
	var registeredOutTypeMap map[converterOutType]converterFunc
	// firstly, it searches the map by input parameter type.
	registeredOutTypeMap, ok = c.typeConverterFuncMap[srcType]
	if !ok {
		return reflect.Value{}, nil, false
	}
	var dstType = dstReflectValueForRefer.Type()
	if dstType.Kind() == reflect.Pointer {
		// Might be **struct, which is support as designed.
		if dstType.Elem().Kind() == reflect.Pointer {
			dstType = dstType.Elem()
		}
	} else if dstReflectValueForRefer.IsValid() && dstReflectValueForRefer.CanAddr() {
		dstType = dstReflectValueForRefer.Addr().Type()
	} else {
		dstType = reflect.PointerTo(dstType)
	}
	// secondly, it searches the input parameter type map
	// and finds the result converter function by the output parameter type.
	f, ok = registeredOutTypeMap[dstType]
	if !ok {
		return reflect.Value{}, nil, false
	}
	return
}

func (c *Converter) getRegisteredAnyConverterFunc(dstReflectValueForRefer reflect.Value) (f AnyConvertFunc) {
	if c.internalConverter.IsAnyConvertFuncEmpty() {
		return nil
	}
	if !dstReflectValueForRefer.IsValid() {
		return nil
	}
	var dstType = dstReflectValueForRefer.Type()
	if dstType.Kind() == reflect.Pointer {
		// Might be **struct, which is support as designed.
		if dstType.Elem().Kind() == reflect.Pointer {
			dstType = dstType.Elem()
		}
	} else if dstReflectValueForRefer.IsValid() && dstReflectValueForRefer.CanAddr() {
		dstType = dstReflectValueForRefer.Addr().Type()
	} else {
		dstType = reflect.PointerTo(dstType)
	}
	return c.internalConverter.GetAnyConvertFuncByType(dstType)
}

func (c *Converter) doCallCustomTypeConverter(
	srcReflectValue reflect.Value,
	dstReflectValue reflect.Value,
	registeredConverterFunc converterFunc,
	srcType reflect.Type,
) (converted bool, err error) {
	// Converter function calling.
	for srcReflectValue.Type() != srcType {
		srcReflectValue = srcReflectValue.Elem()
	}
	result := registeredConverterFunc.Call([]reflect.Value{srcReflectValue})
	if !result[1].IsNil() {
		return false, result[1].Interface().(error)
	}
	// The `result[0]` is a pointer.
	if result[0].IsNil() {
		return false, nil
	}
	var resultValue = result[0]
	for {
		if resultValue.Type() == dstReflectValue.Type() && dstReflectValue.CanSet() {
			dstReflectValue.Set(resultValue)
			converted = true
		} else if dstReflectValue.Kind() == reflect.Pointer {
			if resultValue.Type() == dstReflectValue.Elem().Type() && dstReflectValue.Elem().CanSet() {
				dstReflectValue.Elem().Set(resultValue)
				converted = true
			}
		}
		if converted {
			break
		}
		if resultValue.Kind() == reflect.Pointer {
			resultValue = resultValue.Elem()
		} else {
			break
		}
	}

	return converted, nil
}
