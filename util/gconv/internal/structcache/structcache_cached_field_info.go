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
	//  The [cachedFieldInfoBase] structure cannot be merged with the following [IsField] field into one structure.
	// 	The [IsField] field should be used separately in the [bindStructWithLoopParamsMap] method
	*cachedFieldInfoBase

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

// cachedFieldInfoBase holds the cached info for struct field.
type cachedFieldInfoBase struct {
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

	// IsCustomConvert marks there custom converting function for this field type.
	IsCustomConvert bool

	// StructField is the type info of this field.
	StructField reflect.StructField

	// OtherSameNameFieldIndex holds the sub attributes of the same field name.
	// For example:
	// type Name struct{
	//     LastName  string
	//     FirstName string
	// }
	// type User struct{
	//     Name
	//     LastName  string
	//     FirstName string
	// }
	//
	// As the `LastName` in `User`, its internal attributes:
	//   FieldIndexes = []int{0,1}
	//   // item length 1, as there's only one repeat item with the same field name.
	//   OtherSameNameFieldIndex = [][]int{[]int{1}}
	//
	// In value assignment, the value will be assigned to index {0,1} and {1}.
	OtherSameNameFieldIndex [][]int

	// ConvertFunc is the converting function for this field.
	ConvertFunc func(from any, to reflect.Value)

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

// GetFieldReflectValue retrieves and returns the reflect.Value of given struct value,
// which is used for directly value assignment.
func (cfi *CachedFieldInfo) GetFieldReflectValue(structValue reflect.Value) reflect.Value {
	if len(cfi.FieldIndexes) == 1 {
		return structValue.Field(cfi.FieldIndexes[0])
	}
	return cfi.fieldReflectValue(structValue, cfi.FieldIndexes)
}

// GetOtherFieldReflectValue retrieves and returns the reflect.Value of given struct value with nested index
// by `fieldLevel`, which is used for directly value assignment.
func (cfi *CachedFieldInfo) GetOtherFieldReflectValue(structValue reflect.Value, fieldLevel int) reflect.Value {
	fieldIndex := cfi.OtherSameNameFieldIndex[fieldLevel]
	if len(fieldIndex) == 1 {
		return structValue.Field(fieldIndex[0])
	}
	return cfi.fieldReflectValue(structValue, fieldIndex)
}

func (cfi *CachedFieldInfo) fieldReflectValue(v reflect.Value, fieldIndexes []int) reflect.Value {
	for i, x := range fieldIndexes {
		if i > 0 {
			switch v.Kind() {
			case reflect.Pointer:
				if v.IsNil() {
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
			}
		}
		v = v.Field(x)
	}
	return v
}
