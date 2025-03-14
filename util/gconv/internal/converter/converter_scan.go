// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package converter

import (
	"reflect"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

// ScanOption is the option for the Scan function.
type ScanOption struct {
	// ParamKeyToAttrMap specifies the mapping between parameter keys and struct attribute names.
	ParamKeyToAttrMap map[string]string

	// ContinueOnError specifies whether to continue converting the next element
	// if one element converting fails.
	ContinueOnError bool
}

func (c *Converter) getScanOption(option ...ScanOption) ScanOption {
	if len(option) > 0 {
		return option[0]
	}
	return ScanOption{}
}

// Scan automatically checks the type of `pointer` and converts `params` to `pointer`.
func (c *Converter) Scan(srcValue any, dstPointer any, option ...ScanOption) (err error) {
	// Check if srcValue is nil, in which case no conversion is needed
	if srcValue == nil {
		return nil
	}
	// Check if dstPointer is nil, which is an invalid parameter
	if dstPointer == nil {
		return gerror.NewCode(
			gcode.CodeInvalidParameter,
			`destination pointer should not be nil`,
		)
	}

	// Get the reflection type and value of dstPointer
	var (
		dstPointerReflectType  reflect.Type
		dstPointerReflectValue reflect.Value
	)
	if v, ok := dstPointer.(reflect.Value); ok {
		dstPointerReflectValue = v
		dstPointerReflectType = v.Type()
	} else {
		dstPointerReflectValue = reflect.ValueOf(dstPointer)
		// Do not use dstPointerReflectValue.Type() as dstPointerReflectValue might be zero
		dstPointerReflectType = reflect.TypeOf(dstPointer)
	}

	// Validate the kind of dstPointer
	var dstPointerReflectKind = dstPointerReflectType.Kind()
	if dstPointerReflectKind != reflect.Ptr {
		// If dstPointer is not a pointer, try to get its address
		if dstPointerReflectValue.CanAddr() {
			dstPointerReflectValue = dstPointerReflectValue.Addr()
			dstPointerReflectType = dstPointerReflectValue.Type()
			dstPointerReflectKind = dstPointerReflectType.Kind()
		} else {
			// If dstPointer is not a pointer and cannot be addressed, return an error
			return gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`destination pointer should be type of pointer, but got type: %v`,
				dstPointerReflectType,
			)
		}
	}

	// Get the reflection value of srcValue
	var srcValueReflectValue reflect.Value
	if v, ok := srcValue.(reflect.Value); ok {
		srcValueReflectValue = v
	} else {
		srcValueReflectValue = reflect.ValueOf(srcValue)
	}

	// Get the element type and kind of dstPointer
	var (
		dstPointerReflectValueElem     = dstPointerReflectValue.Elem()
		dstPointerReflectValueElemKind = dstPointerReflectValueElem.Kind()
	)
	// Handle multiple level pointers
	if dstPointerReflectValueElemKind == reflect.Ptr {
		if dstPointerReflectValueElem.IsNil() {
			// Create a new value for the pointer dereference
			nextLevelPtr := reflect.New(dstPointerReflectValueElem.Type().Elem())
			// Recursively scan into the dereferenced pointer
			if err = c.Scan(srcValueReflectValue, nextLevelPtr, option...); err == nil {
				dstPointerReflectValueElem.Set(nextLevelPtr)
			}
			return
		}
		return c.Scan(srcValueReflectValue, dstPointerReflectValueElem, option...)
	}

	// Check if srcValue and dstPointer are the same type, in which case direct assignment can be performed
	if ok := c.doConvertWithTypeCheck(srcValueReflectValue, dstPointerReflectValueElem); ok {
		return nil
	}

	scanOption := c.getScanOption(option...)
	// Handle different destination types
	switch dstPointerReflectValueElemKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := c.Int64(srcValue)
		if err != nil && !scanOption.ContinueOnError {
			return err
		}
		dstPointerReflectValueElem.SetInt(v)
		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := c.Uint64(srcValue)
		if err != nil && !scanOption.ContinueOnError {
			return err
		}
		dstPointerReflectValueElem.SetUint(v)
		return nil

	case reflect.Float32, reflect.Float64:
		v, err := c.Float64(srcValue)
		if err != nil && !scanOption.ContinueOnError {
			return err
		}
		dstPointerReflectValueElem.SetFloat(v)
		return nil

	case reflect.String:
		v, err := c.String(srcValue)
		if err != nil && !scanOption.ContinueOnError {
			return err
		}
		dstPointerReflectValueElem.SetString(v)
		return nil

	case reflect.Bool:
		v, err := c.Bool(srcValue)
		if err != nil && !scanOption.ContinueOnError {
			return err
		}
		dstPointerReflectValueElem.SetBool(v)
		return nil

	case reflect.Slice:
		// Handle slice type conversion
		var (
			dstElemType = dstPointerReflectValueElem.Type().Elem()
			dstElemKind = dstElemType.Kind()
		)
		// The slice element might be a pointer type
		if dstElemKind == reflect.Ptr {
			dstElemType = dstElemType.Elem()
			dstElemKind = dstElemType.Kind()
		}
		// Special handling for struct or map slice elements
		if dstElemKind == reflect.Struct || dstElemKind == reflect.Map {
			return c.doScanForComplicatedTypes(srcValue, dstPointer, dstPointerReflectType, scanOption)
		}
		// Handle basic type slice conversions
		var srcValueReflectValueKind = srcValueReflectValue.Kind()
		if srcValueReflectValueKind == reflect.Slice || srcValueReflectValueKind == reflect.Array {
			var (
				srcLen   = srcValueReflectValue.Len()
				newSlice = reflect.MakeSlice(dstPointerReflectValueElem.Type(), srcLen, srcLen)
			)
			for i := 0; i < srcLen; i++ {
				srcElem := srcValueReflectValue.Index(i).Interface()
				switch dstElemType.Kind() {
				case reflect.String:
					v, err := c.String(srcElem)
					if err != nil && !scanOption.ContinueOnError {
						return err
					}
					newSlice.Index(i).SetString(v)
				case reflect.Int:
					v, err := c.Int64(srcElem)
					if err != nil && !scanOption.ContinueOnError {
						return err
					}
					newSlice.Index(i).SetInt(v)
				case reflect.Int64:
					v, err := c.Int64(srcElem)
					if err != nil && !scanOption.ContinueOnError {
						return err
					}
					newSlice.Index(i).SetInt(v)
				case reflect.Float64:
					v, err := c.Float64(srcElem)
					if err != nil && !scanOption.ContinueOnError {
						return err
					}
					newSlice.Index(i).SetFloat(v)
				case reflect.Bool:
					v, err := c.Bool(srcElem)
					if err != nil && !scanOption.ContinueOnError {
						return err
					}
					newSlice.Index(i).SetBool(v)
				default:
					return c.Scan(
						srcElem, newSlice.Index(i).Addr().Interface(), option...,
					)
				}
			}
			dstPointerReflectValueElem.Set(newSlice)
			return nil
		}
		return c.doScanForComplicatedTypes(srcValue, dstPointer, dstPointerReflectType, scanOption)

	default:
		// Handle complex types (structs, maps, etc.)
		return c.doScanForComplicatedTypes(srcValue, dstPointer, dstPointerReflectType, scanOption)
	}
}

// doScanForComplicatedTypes handles the scanning of complex data types.
// It supports converting between maps, structs, and slices of these types.
// The function first attempts JSON conversion, then falls back to specific type handling.
//
// It supports `pointer` in type of `*map/*[]map/*[]*map/*struct/**struct/*[]struct/*[]*struct` for converting.
//
// Parameters:
// - srcValue: The source value to convert from
// - dstPointer: The destination pointer to convert to
// - dstPointerReflectType: The reflection type of the destination pointer
// - paramKeyToAttrMap: Optional mapping between parameter keys and struct attribute names
func (c *Converter) doScanForComplicatedTypes(
	srcValue, dstPointer any,
	dstPointerReflectType reflect.Type,
	option ScanOption,
) error {
	// Try JSON conversion first
	ok, err := c.doConvertWithJsonCheck(srcValue, dstPointer)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	// Handle specific type conversions
	var (
		dstPointerReflectTypeElem     = dstPointerReflectType.Elem()
		dstPointerReflectTypeElemKind = dstPointerReflectTypeElem.Kind()
		keyToAttributeNameMapping     = option.ParamKeyToAttrMap
	)
	// Handle different destination types
	switch dstPointerReflectTypeElemKind {
	case reflect.Map:
		// Convert map to map
		return c.MapToMap(srcValue, dstPointer, keyToAttributeNameMapping, MapOption{
			ContinueOnError: option.ContinueOnError,
		})

	case reflect.Array, reflect.Slice:
		var (
			sliceElem     = dstPointerReflectTypeElem.Elem()
			sliceElemKind = sliceElem.Kind()
		)
		// Handle pointer elements
		for sliceElemKind == reflect.Ptr {
			sliceElem = sliceElem.Elem()
			sliceElemKind = sliceElem.Kind()
		}
		if sliceElemKind == reflect.Map {
			// Convert to slice of maps
			return c.MapToMaps(srcValue, dstPointer, keyToAttributeNameMapping, MapOption{
				ContinueOnError: option.ContinueOnError,
			})
		}
		// Convert to slice of structs
		var (
			sliceOption = SliceOption{
				ContinueOnError: option.ContinueOnError,
			}
			mapOption = StructOption{
				ParamKeyToAttrMap: keyToAttributeNameMapping,
				ContinueOnError:   option.ContinueOnError,
			}
		)
		return c.Structs(srcValue, dstPointer, StructsOption{
			SliceOption:  sliceOption,
			StructOption: mapOption,
		})

	default:
		structOption := StructOption{
			ParamKeyToAttrMap: keyToAttributeNameMapping,
			PriorityTag:       "",
			ContinueOnError:   option.ContinueOnError,
		}
		return c.Struct(srcValue, dstPointer, structOption)
	}
}

// doConvertWithTypeCheck supports `pointer` in type of `*map/*[]map/*[]*map/*struct/**struct/*[]struct/*[]*struct`
// for converting.
func (c *Converter) doConvertWithTypeCheck(srcValueReflectValue, dstPointerReflectValueElem reflect.Value) (ok bool) {
	if !dstPointerReflectValueElem.IsValid() || !srcValueReflectValue.IsValid() {
		return false
	}
	switch {
	// Examples:
	// UploadFile       => UploadFile
	// []UploadFile     => []UploadFile
	// *UploadFile      => *UploadFile
	// *[]UploadFile    => *[]UploadFile
	// map[int][int]    => map[int][int]
	// []map[int][int]  => []map[int][int]
	// *[]map[int][int] => *[]map[int][int]
	case dstPointerReflectValueElem.Type() == srcValueReflectValue.Type():
		dstPointerReflectValueElem.Set(srcValueReflectValue)
		return true

	// Examples:
	// UploadFile      => *UploadFile
	// []UploadFile    => *[]UploadFile
	// map[int][int]   => *map[int][int]
	// []map[int][int] => *[]map[int][int]
	case dstPointerReflectValueElem.Kind() == reflect.Ptr &&
		dstPointerReflectValueElem.Elem().IsValid() &&
		dstPointerReflectValueElem.Elem().Type() == srcValueReflectValue.Type():
		dstPointerReflectValueElem.Elem().Set(srcValueReflectValue)
		return true

	// Examples:
	// *UploadFile      => UploadFile
	// *[]UploadFile    => []UploadFile
	// *map[int][int]   => map[int][int]
	// *[]map[int][int] => []map[int][int]
	case srcValueReflectValue.Kind() == reflect.Ptr &&
		srcValueReflectValue.Elem().IsValid() &&
		dstPointerReflectValueElem.Type() == srcValueReflectValue.Elem().Type():
		dstPointerReflectValueElem.Set(srcValueReflectValue.Elem())
		return true

	default:
		return false
	}
}

// doConvertWithJsonCheck attempts to convert the source value to the destination
// using JSON marshaling and unmarshaling. This is particularly useful for complex
// types that can be represented as JSON.
//
// Parameters:
// - srcValue: The source value to convert from
// - dstPointer: The destination pointer to convert to
//
// Returns:
// - bool: true if JSON conversion was successful
// - error: any error that occurred during conversion
func (c *Converter) doConvertWithJsonCheck(srcValue any, dstPointer any) (ok bool, err error) {
	switch valueResult := srcValue.(type) {
	case []byte:
		if json.Valid(valueResult) {
			if dstPointerReflectType, ok := dstPointer.(reflect.Value); ok {
				if dstPointerReflectType.Kind() == reflect.Ptr {
					if dstPointerReflectType.IsNil() {
						return false, nil
					}
					return true, json.UnmarshalUseNumber(valueResult, dstPointerReflectType.Interface())
				} else if dstPointerReflectType.CanAddr() {
					return true, json.UnmarshalUseNumber(valueResult, dstPointerReflectType.Addr().Interface())
				}
			} else {
				return true, json.UnmarshalUseNumber(valueResult, dstPointer)
			}
		}

	case string:
		if valueBytes := []byte(valueResult); json.Valid(valueBytes) {
			if dstPointerReflectType, ok := dstPointer.(reflect.Value); ok {
				if dstPointerReflectType.Kind() == reflect.Ptr {
					if dstPointerReflectType.IsNil() {
						return false, nil
					}
					return true, json.UnmarshalUseNumber(valueBytes, dstPointerReflectType.Interface())
				} else if dstPointerReflectType.CanAddr() {
					return true, json.UnmarshalUseNumber(valueBytes, dstPointerReflectType.Addr().Interface())
				}
			} else {
				return true, json.UnmarshalUseNumber(valueBytes, dstPointer)
			}
		}

	default:
		// The `params` might be struct that implements interface function Interface, eg: gvar.Var.
		if v, ok := srcValue.(localinterface.IInterface); ok {
			return c.doConvertWithJsonCheck(v.Interface(), dstPointer)
		}
	}
	return false, nil
}
