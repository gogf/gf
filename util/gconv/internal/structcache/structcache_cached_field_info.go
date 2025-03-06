// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package structcache

import (
	"reflect"
	"sync/atomic"
)

// CachedFieldInfo holds the cached info for struct field.
type CachedFieldInfo struct {
	// WARN:
	//  The [CachedFieldInfoBase] structure cannot be merged with the following [IsField] field into one structure.
	// 	The [IsField] field should be used separately in the [bindStructWithLoopParamsMap] method
	*CachedFieldInfoBase

	// This field is mainly used in the [bindStructWithLoopParamsMap] method.
	// This field is needed when both `fieldName` and `tag` of a field exist in the map.
	// For example:
	// field string `json:"name"`
	// map = {
	//     "field" : "f1",
	//     "name" : "n1",
	// }
	// The `name` should be used here.
	// In the bindStructWithLoopParamsMap method, due to the disorder of `map`, `field` may be traversed first.
	// This field is more about priority, that is, the priority of `name` is higher than that of `field`,
	// even if it has been set before.
	IsField bool
}

// CachedFieldInfoBase holds the cached info for struct field.
type CachedFieldInfoBase struct {
	// FieldIndexes holds the global index number from struct info.
	// The field may belong to an embedded structure, so it is defined here as []int.
	FieldIndexes []int

	// PriorityTagAndFieldName holds the tag value(conv, param, p, c, json) and the field name.
	// PriorityTagAndFieldName contains the field name, which is the last item of slice.
	PriorityTagAndFieldName []string

	// IsCommonInterface marks this field implements common interfaces as:
	// - iUnmarshalValue
	// - iUnmarshalText
	// - iUnmarshalJSON
	// Purpose: reduce the interface asserting cost in runtime.
	IsCommonInterface bool

	// HasCustomConvert marks there custom converting function for this field type.
	// A custom converting function is a function that user defined for converting specified type
	// to another type.
	HasCustomConvert bool

	// StructField is the type info of this field.
	StructField reflect.StructField

	// OtherSameNameField stores fields with the same name and type or different types of nested structures.
	//
	// For example:
	// type ID struct{
	//     ID1  string
	//     ID2 int
	// }
	// type Card struct{
	//     ID
	//     ID1  uint64
	//     ID2 int64
	// }
	//
	// We will cache each ID1 and ID2 separately,
	// even if their types are different and their indexes are different
	OtherSameNameField []*CachedFieldInfo

	// ConvertFunc is the converting function for this field.
	ConvertFunc AnyConvertFunc

	// The last fuzzy matching key for this field.
	// The fuzzy matching occurs only if there are no direct tag and field name matching in the params map.
	// TODO If different paramsMaps contain paramKeys in different formats and all hit the same fieldName,
	//      the cached value may be continuously updated.
	// LastFuzzyKey string.
	LastFuzzyKey atomic.Value

	// removeSymbolsFieldName is used for quick fuzzy match for parameter key.
	// removeSymbolsFieldName = utils.RemoveSymbols(fieldName)
	RemoveSymbolsFieldName string
}

// FieldName returns the field name of current field info.
func (cfi *CachedFieldInfo) FieldName() string {
	return cfi.PriorityTagAndFieldName[len(cfi.PriorityTagAndFieldName)-1]
}

// GetFieldReflectValueFrom retrieves and returns the `reflect.Value` of given struct field,
// which is used for directly value assignment.
//
// Note that, the input parameter `structValue` might be initialized internally.
func (cfi *CachedFieldInfo) GetFieldReflectValueFrom(structValue reflect.Value) reflect.Value {
	if len(cfi.FieldIndexes) == 1 {
		// no nested struct.
		return structValue.Field(cfi.FieldIndexes[0])
	}
	return cfi.fieldReflectValue(structValue, cfi.FieldIndexes)
}

// GetOtherFieldReflectValueFrom retrieves and returns the `reflect.Value` of given struct field with nested index
// by `fieldLevel`, which is used for directly value assignment.
//
// Note that, the input parameter `structValue` might be initialized internally.
func (cfi *CachedFieldInfo) GetOtherFieldReflectValueFrom(structValue reflect.Value, fieldIndex []int) reflect.Value {
	if len(fieldIndex) == 1 {
		// no nested struct.
		return structValue.Field(fieldIndex[0])
	}
	return cfi.fieldReflectValue(structValue, fieldIndex)
}

func (cfi *CachedFieldInfo) fieldReflectValue(v reflect.Value, fieldIndexes []int) reflect.Value {
	for i, x := range fieldIndexes {
		if i > 0 {
			// it means nested struct.
			switch v.Kind() {
			case reflect.Pointer:
				if v.IsNil() {
					// Initialization.
					v.Set(reflect.New(v.Type().Elem()))
				}
				v = v.Elem()

			case reflect.Interface:
				// Compatible with previous code
				// Interface => struct
				v = v.Elem()
				if v.Kind() == reflect.Ptr {
					// maybe *struct or other types
					v = v.Elem()
				}
			default:
			}
		}
		v = v.Field(x)
	}
	return v
}
