// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"context"
	"reflect"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv/internal/structcache"
)

type (
	converterInType  = reflect.Type
	converterOutType = reflect.Type
	converterFunc    = reflect.Value
)

// Converter is the manager for type converting.
type Converter struct {
	internalConvertConfig *structcache.ConvertConfig
	typeConverterFuncMap  map[converterInType]map[converterOutType]converterFunc
}

var (
	intType   = reflect.TypeOf(0)
	int8Type  = reflect.TypeOf(int8(0))
	int16Type = reflect.TypeOf(int16(0))
	int32Type = reflect.TypeOf(int32(0))
	int64Type = reflect.TypeOf(int64(0))

	uintType   = reflect.TypeOf(uint(0))
	uint8Type  = reflect.TypeOf(uint8(0))
	uint16Type = reflect.TypeOf(uint16(0))
	uint32Type = reflect.TypeOf(uint32(0))
	uint64Type = reflect.TypeOf(uint64(0))

	float32Type = reflect.TypeOf(float32(0))
	float64Type = reflect.TypeOf(float64(0))

	stringType = reflect.TypeOf("")
	bytesType  = reflect.TypeOf([]byte{})

	boolType = reflect.TypeOf(false)

	timeType  = reflect.TypeOf((*time.Time)(nil)).Elem()
	gtimeType = reflect.TypeOf((*gtime.Time)(nil)).Elem()
)

// NewConverter creates and returns management object for type converting.
func NewConverter() *Converter {
	cf := &Converter{
		internalConvertConfig: structcache.NewConvertConfig(),
		typeConverterFuncMap:  make(map[converterInType]map[converterOutType]converterFunc),
	}
	cf.registerBuiltInConverter()
	return cf
}

func (c *Converter) registerBuiltInConverter() {
	c.registerAnyConvertFuncForTypes(
		c.builtInAnyConvertFuncForInt64, intType, int8Type, int16Type, int32Type, int64Type,
	)
	c.registerAnyConvertFuncForTypes(
		c.builtInAnyConvertFuncForUint64, uintType, uint8Type, uint16Type, uint32Type, uint64Type,
	)
	c.registerAnyConvertFuncForTypes(
		c.builtInAnyConvertFuncForString, stringType,
	)
	c.registerAnyConvertFuncForTypes(
		c.builtInAnyConvertFuncForFloat64, float32Type, float64Type,
	)
	c.registerAnyConvertFuncForTypes(
		c.builtInAnyConvertFuncForBool, boolType,
	)
	c.registerAnyConvertFuncForTypes(
		c.builtInAnyConvertFuncForBytes, bytesType,
	)
	c.registerAnyConvertFuncForTypes(
		c.builtInAnyConvertFuncForTime, timeType,
	)
	c.registerAnyConvertFuncForTypes(
		c.builtInAnyConvertFuncForGTime, gtimeType,
	)
}

func (c *Converter) registerAnyConvertFuncForTypes(convertFunc AnyConvertFunc, types ...reflect.Type) {
	for _, t := range types {
		c.internalConvertConfig.RegisterAnyConvertFunc(t, convertFunc)
	}
}

// RegisterTypeConverterFunc registers custom converter.
// It must be registered before you use this custom converting feature.
// It is suggested to do it in boot procedure of the process.
//
// Note:
//  1. The parameter `fn` must be defined as pattern `func(T1) (T2, error)`.
//     It will convert type `T1` to type `T2`.
//  2. The `T1` should not be type of pointer, but the `T2` should be type of pointer.
func (c *Converter) RegisterTypeConverterFunc(fn any) (err error) {
	var (
		fnReflectType = reflect.TypeOf(fn)
		errType       = reflect.TypeOf((*error)(nil)).Elem()
	)
	if fnReflectType.Kind() != reflect.Func ||
		fnReflectType.NumIn() != 1 || fnReflectType.NumOut() != 2 ||
		!fnReflectType.Out(1).Implements(errType) {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"parameter must be type of converter function and defined as pattern `func(T1) (T2, error)`, "+
				"but defined as `%s`",
			fnReflectType.String(),
		)
		return
	}

	// The Key and Value of the converter map should not be pointer.
	var (
		inType  = fnReflectType.In(0)
		outType = fnReflectType.Out(0)
	)
	if inType.Kind() == reflect.Pointer {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"invalid converter function `%s`: invalid input parameter type `%s`, should not be type of pointer",
			fnReflectType.String(), inType.String(),
		)
		return
	}
	if outType.Kind() != reflect.Pointer {
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"invalid converter function `%s`: invalid output parameter type `%s` should be type of pointer",
			fnReflectType.String(), outType.String(),
		)
		return
	}

	registeredOutTypeMap, ok := c.typeConverterFuncMap[inType]
	if !ok {
		registeredOutTypeMap = make(map[converterOutType]converterFunc)
		c.typeConverterFuncMap[inType] = registeredOutTypeMap
	}
	if _, ok = registeredOutTypeMap[outType]; ok {
		err = gerror.NewCodef(
			gcode.CodeInvalidOperation,
			"the converter parameter type `%s` to type `%s` has already been registered",
			inType.String(), outType.String(),
		)
		return
	}
	registeredOutTypeMap[outType] = reflect.ValueOf(fn)
	c.internalConvertConfig.RegisterTypeConvertFunc(outType)
	return
}

func (c *Converter) getRegisteredConverterFuncAndSrcType(
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

func (c *Converter) callCustomConverterWithRefer(
	srcReflectValue, referReflectValue reflect.Value,
) (dstReflectValue reflect.Value, converted bool, err error) {
	registeredConverterFunc, srcType, ok := c.getRegisteredConverterFuncAndSrcType(srcReflectValue, referReflectValue)
	if !ok {
		return reflect.Value{}, false, nil
	}
	dstReflectValue = reflect.New(referReflectValue.Type()).Elem()
	converted, err = c.doCallCustomConverter(srcReflectValue, dstReflectValue, registeredConverterFunc, srcType)
	return
}

// callCustomConverter call the custom converter. It will try some possible type.
func (c *Converter) callCustomConverter(srcReflectValue, dstReflectValue reflect.Value) (converted bool, err error) {
	registeredConverterFunc, srcType, ok := c.getRegisteredConverterFuncAndSrcType(srcReflectValue, dstReflectValue)
	if !ok {
		return false, nil
	}
	return c.doCallCustomConverter(srcReflectValue, dstReflectValue, registeredConverterFunc, srcType)
}

func (c *Converter) doCallCustomConverter(
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

type doConvertInput struct {
	FromValue  interface{}   // Value that is converted from.
	ToTypeName string        // Target value type name in string.
	ReferValue interface{}   // Referred value, a value in type `ToTypeName`. Note that its type might be reflect.Value.
	Extra      []interface{} // Extra values for implementing the converting.

	// Marks that the value is already converted and set to `ReferValue`. Caller can ignore the returned result.
	// It is an attribute for internal usage purpose.
	alreadySetToReferValue bool
}

// doConvert does commonly use types converting.
func (c *Converter) doConvert(in doConvertInput) (convertedValue interface{}) {
	switch in.ToTypeName {
	case "int":
		return Int(in.FromValue)
	case "*int":
		if _, ok := in.FromValue.(*int); ok {
			return in.FromValue
		}
		v := Int(in.FromValue)
		return &v

	case "int8":
		return Int8(in.FromValue)
	case "*int8":
		if _, ok := in.FromValue.(*int8); ok {
			return in.FromValue
		}
		v := Int8(in.FromValue)
		return &v

	case "int16":
		return Int16(in.FromValue)
	case "*int16":
		if _, ok := in.FromValue.(*int16); ok {
			return in.FromValue
		}
		v := Int16(in.FromValue)
		return &v

	case "int32":
		return Int32(in.FromValue)
	case "*int32":
		if _, ok := in.FromValue.(*int32); ok {
			return in.FromValue
		}
		v := Int32(in.FromValue)
		return &v

	case "int64":
		return Int64(in.FromValue)
	case "*int64":
		if _, ok := in.FromValue.(*int64); ok {
			return in.FromValue
		}
		v := Int64(in.FromValue)
		return &v

	case "uint":
		return Uint(in.FromValue)
	case "*uint":
		if _, ok := in.FromValue.(*uint); ok {
			return in.FromValue
		}
		v := Uint(in.FromValue)
		return &v

	case "uint8":
		return Uint8(in.FromValue)
	case "*uint8":
		if _, ok := in.FromValue.(*uint8); ok {
			return in.FromValue
		}
		v := Uint8(in.FromValue)
		return &v

	case "uint16":
		return Uint16(in.FromValue)
	case "*uint16":
		if _, ok := in.FromValue.(*uint16); ok {
			return in.FromValue
		}
		v := Uint16(in.FromValue)
		return &v

	case "uint32":
		return Uint32(in.FromValue)
	case "*uint32":
		if _, ok := in.FromValue.(*uint32); ok {
			return in.FromValue
		}
		v := Uint32(in.FromValue)
		return &v

	case "uint64":
		return Uint64(in.FromValue)
	case "*uint64":
		if _, ok := in.FromValue.(*uint64); ok {
			return in.FromValue
		}
		v := Uint64(in.FromValue)
		return &v

	case "float32":
		return Float32(in.FromValue)
	case "*float32":
		if _, ok := in.FromValue.(*float32); ok {
			return in.FromValue
		}
		v := Float32(in.FromValue)
		return &v

	case "float64":
		return Float64(in.FromValue)
	case "*float64":
		if _, ok := in.FromValue.(*float64); ok {
			return in.FromValue
		}
		v := Float64(in.FromValue)
		return &v

	case "bool":
		return Bool(in.FromValue)
	case "*bool":
		if _, ok := in.FromValue.(*bool); ok {
			return in.FromValue
		}
		v := Bool(in.FromValue)
		return &v

	case "string":
		return String(in.FromValue)
	case "*string":
		if _, ok := in.FromValue.(*string); ok {
			return in.FromValue
		}
		v := String(in.FromValue)
		return &v

	case "[]byte":
		return Bytes(in.FromValue)
	case "[]int":
		return Ints(in.FromValue)
	case "[]int32":
		return Int32s(in.FromValue)
	case "[]int64":
		return Int64s(in.FromValue)
	case "[]uint":
		return Uints(in.FromValue)
	case "[]uint8":
		return Bytes(in.FromValue)
	case "[]uint32":
		return Uint32s(in.FromValue)
	case "[]uint64":
		return Uint64s(in.FromValue)
	case "[]float32":
		return Float32s(in.FromValue)
	case "[]float64":
		return Float64s(in.FromValue)
	case "[]string":
		return Strings(in.FromValue)

	case "Time", "time.Time":
		if len(in.Extra) > 0 {
			return Time(in.FromValue, String(in.Extra[0]))
		}
		return Time(in.FromValue)
	case "*time.Time":
		var v time.Time
		if len(in.Extra) > 0 {
			v = Time(in.FromValue, String(in.Extra[0]))
		} else {
			if _, ok := in.FromValue.(*time.Time); ok {
				return in.FromValue
			}
			v = Time(in.FromValue)
		}
		return &v

	case "GTime", "gtime.Time":
		if len(in.Extra) > 0 {
			if v := GTime(in.FromValue, String(in.Extra[0])); v != nil {
				return *v
			} else {
				return *gtime.New()
			}
		}
		if v := GTime(in.FromValue); v != nil {
			return *v
		} else {
			return *gtime.New()
		}
	case "*gtime.Time":
		if len(in.Extra) > 0 {
			if v := GTime(in.FromValue, String(in.Extra[0])); v != nil {
				return v
			} else {
				return gtime.New()
			}
		}
		if v := GTime(in.FromValue); v != nil {
			return v
		} else {
			return gtime.New()
		}

	case "Duration", "time.Duration":
		return Duration(in.FromValue)
	case "*time.Duration":
		if _, ok := in.FromValue.(*time.Duration); ok {
			return in.FromValue
		}
		v := Duration(in.FromValue)
		return &v

	case "map[string]string":
		return MapStrStr(in.FromValue)

	case "map[string]interface {}":
		return Map(in.FromValue)

	case "[]map[string]interface {}":
		return Maps(in.FromValue)

	case "RawMessage", "json.RawMessage":
		// issue 3449
		bytes, err := json.Marshal(in.FromValue)
		if err != nil {
			intlog.Errorf(context.TODO(), `%+v`, err)
		}
		return bytes

	default:
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
			if dstReflectValue, ok, _ := c.callCustomConverterWithRefer(fromReflectValue, referReflectValue); ok {
				return dstReflectValue.Interface()
			}

			defer func() {
				if recover() != nil {
					in.alreadySetToReferValue = false
					if err := c.bindVarToReflectValue(referReflectValue, in.FromValue, nil); err == nil {
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
					refElementValue := reflect.ValueOf(c.doConvert(in))
					originTypeValue := reflect.New(refElementValue.Type()).Elem()
					originTypeValue.Set(refElementValue)
					in.alreadySetToReferValue = true
					return originTypeValue.Addr().Convert(referReflectValue.Type()).Interface()
				}

			case reflect.Map:
				var targetValue = reflect.New(referReflectValue.Type()).Elem()
				if err := c.MapToMap(in.FromValue, targetValue); err == nil {
					in.alreadySetToReferValue = true
				}
				return targetValue.Interface()

			default:

			}
			in.ToTypeName = referReflectValue.Kind().String()
			in.ReferValue = nil
			in.alreadySetToReferValue = true
			convertedValue = reflect.ValueOf(c.doConvert(in)).Convert(referReflectValue.Type()).Interface()
			return convertedValue
		}
		return in.FromValue
	}
}

func (c *Converter) doConvertWithReflectValueSet(reflectValue reflect.Value, in doConvertInput) {
	convertedValue := c.doConvert(in)
	if !in.alreadySetToReferValue {
		reflectValue.Set(reflect.ValueOf(convertedValue))
	}
}
