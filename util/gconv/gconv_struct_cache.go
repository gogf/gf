// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gtag"
)

var (
	poolUsedParamsKeyOrTagNameMap = &sync.Pool{
		New: func() any {
			return make(map[string]struct{})
		},
	}
	// Allow users to choose whether to enable this feature
	convCacheExperiment = true
)

func UseConvCacheExperiment(b bool) {
	convCacheExperiment = b
}

func poolGetUsedParamsKeyOrTagNameMap() map[string]struct{} {
	return poolUsedParamsKeyOrTagNameMap.Get().(map[string]struct{})
}

func poolPutUsedParamsKeyOrTagNameMap(m map[string]struct{}) {
	// need to be cleared, otherwise there will be a bug
	for k := range m {
		delete(m, k)
	}
	poolUsedParamsKeyOrTagNameMap.Put(m)
}

type convertFieldInfo struct {
	// The index of a field may be a nested structure, so [] int is required
	fieldIndex []int
	// All tags in the field include the field name, and the field name is the last one
	tags []string
	// 1.iUnmarshalValue
	// 2.iUnmarshalText
	// 3.iUnmarshalJSON
	// Implemented one of the three types of interfaces mentioned above,
	// except for [time.Time] and [gtime.Time]
	isCommonInterface bool
	fieldTypeName     string
	// Field index used to store duplicate names
	// Generally used for nested structures
	otherFieldIndex [][]int
	// Cache type conversion function
	// For example:
	// if the type of the field is int, then directly cache the [gconv.Int] function
	convFunc func(from any, to reflect.Value)
}

func (c *convertFieldInfo) FieldName() string {
	return c.tags[len(c.tags)-1]
}

type convertStructInfo struct {
	// key = field's name
	fields map[string]*convertFieldInfo
}

func (structInfo *convertStructInfo) AddField(field reflect.StructField, fieldIndex []int, priorityTags []string) *convertFieldInfo {
	convFieldInfo, ok := structInfo.fields[field.Name]
	if !ok {
		fieldInfo := &convertFieldInfo{
			isCommonInterface: checkTypeIsImplCommonInterface(field),
			fieldTypeName:     field.Type.String(),
			fieldIndex:        fieldIndex,
			convFunc:          getFieldConvFunc(field.Type.String()),
		}

		structInfo.fields[field.Name] = fieldInfo
		fieldInfo.tags = getFieldTags(field, priorityTags)
		return fieldInfo
	}
	if convFieldInfo.otherFieldIndex == nil {
		convFieldInfo.otherFieldIndex = make([][]int, 0, 2)
	}

	convFieldInfo.otherFieldIndex = append(convFieldInfo.otherFieldIndex, fieldIndex)
	return convFieldInfo
}

func ptrConvFunc(convFunc func(from any, to reflect.Value)) func(from any, to reflect.Value) {
	return func(from any, to reflect.Value) {
		if to.IsNil() {
			to.Set(reflect.New(to.Type().Elem()))
		}
		convFunc(from, to.Elem())
	}
}

func getFieldConvFunc(fieldType string) (convFunc func(from any, to reflect.Value)) {
	if fieldType[0] == '*' {
		convFunc = getFieldConvFunc(fieldType[1:])
		if convFunc == nil {
			return nil
		}
		return ptrConvFunc(convFunc)
	}
	switch fieldType {
	case "int":
		convFunc = func(from any, to reflect.Value) {
			to.SetInt(int64(Int(from)))
		}
	case "int64":
		convFunc = func(from any, to reflect.Value) {
			to.SetInt(Int64(from))
		}
	case "uint":
		convFunc = func(from any, to reflect.Value) {
			to.SetUint(uint64(Uint(from)))
		}
	case "uint64":
		convFunc = func(from any, to reflect.Value) {
			to.SetUint(Uint64(from))
		}
	case "string":
		convFunc = func(from any, to reflect.Value) {
			to.SetString(String(from))
		}
	case "float32":
		convFunc = func(from any, to reflect.Value) {
			to.SetFloat(float64(Float32(from)))
		}
	case "float64":
		convFunc = func(from any, to reflect.Value) {
			to.SetFloat(Float64(from))
		}
	case "Time", "time.Time":
		convFunc = func(from any, to reflect.Value) {
			*to.Addr().Interface().(*time.Time) = Time(from)
		}
	case "GTime", "gtime.Time":
		convFunc = func(from any, to reflect.Value) {
			v := GTime(from)
			if v == nil {
				v = gtime.New()
			}
			*to.Addr().Interface().(*gtime.Time) = *v
		}
	default:
		return nil
	}

	return convFunc
}

type toBeConvertedStructInfo struct {
	// key = field's name
	fields map[string]toBeConvertedFieldInfo
}

func (t *toBeConvertedStructInfo) AddField(fieldName string, fieldInfo *convertFieldInfo) {
	convertInfo, ok := t.fields[fieldName]
	if !ok {
		t.fields[fieldName] = toBeConvertedFieldInfo{
			Value:            nil,
			convertFieldInfo: fieldInfo,
		}
		return
	}
	if convertInfo.otherFieldIndex == nil {
		convertInfo.otherFieldIndex = make([][]int, 0, 2)
	}
	convertInfo.otherFieldIndex = append(convertInfo.otherFieldIndex, fieldInfo.fieldIndex)
}

var (
	cacheConvStructsInfo = make(map[reflect.Type]*convertStructInfo)
)

func setCacheConvStructInfo(structType reflect.Type, info *convertStructInfo) {
	if convCacheExperiment {
		cacheConvStructsInfo[structType] = info
	}
}

func getCacheConvStructInfo(structType reflect.Type) (*convertStructInfo, bool) {
	if convCacheExperiment {
		structInfo, ok := cacheConvStructsInfo[structType]
		return structInfo, ok
	}
	return nil, false
}

func getConvStructInfo(structType reflect.Type, priorityTag string) *toBeConvertedStructInfo {
	if structType.Kind() != reflect.Struct {
		return nil
	}
	// key=field's name
	toBeConvertedFieldNameToInfo := &toBeConvertedStructInfo{
		fields: make(map[string]toBeConvertedFieldInfo),
	}

	// Check if it has been cached
	structInfo, ok := getCacheConvStructInfo(structType)
	if ok {
		for k, v := range structInfo.fields {
			toBeConvertedFieldNameToInfo.AddField(k, v)
		}
		return toBeConvertedFieldNameToInfo
	}

	structInfo = &convertStructInfo{
		fields: make(map[string]*convertFieldInfo),
	}
	var (
		priorityTagArray []string
		parentIndex      = make([]int, 0)
	)
	if priorityTag != "" {
		priorityTagArray = append(utils.SplitAndTrim(priorityTag, ","), gtag.StructTagPriority...)
	} else {
		priorityTagArray = gtag.StructTagPriority
	}
	parseStruct(structType, parentIndex, structInfo, priorityTagArray)
	setCacheConvStructInfo(structType, structInfo)
	for k, v := range structInfo.fields {
		toBeConvertedFieldNameToInfo.AddField(k, v)
	}
	return toBeConvertedFieldNameToInfo
}

func parseStruct(structType reflect.Type, parentIndex []int, structInfo *convertStructInfo, priorityTagArray []string) {
	var (
		fieldName   string
		structField reflect.StructField
		fieldType   reflect.Type
	)
	// TODO:
	//  Check if the structure has already been cached in the cache.
	//  If it has been cached, some information can be reused,
	//  but the [FieldIndex] needs to be reset.
	//  We will not implement it temporarily because it is somewhat complex
	for i := 0; i < structType.NumField(); i++ {
		structField = structType.Field(i)
		fieldType = structField.Type
		fieldName = structField.Name
		// Only do converting to public attributes.
		if !utils.IsLetterUpper(fieldName[0]) {
			continue
		}
		if structField.Anonymous == false {
			structInfo.AddField(structField, append(parentIndex, i), priorityTagArray)
			continue
		}
		// Maybe it's struct/*struct embedded.
		if fieldType.Kind() == reflect.Interface {
		} else {
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}
		}
		if fieldType.Kind() != reflect.Struct {
			continue
		}
		// TODO: If it's an anonymous field with a tag, doesn't it need to be recursive?
		structInfo.AddField(structField, append(parentIndex, i), priorityTagArray)
		parseStruct(fieldType, append(parentIndex, i), structInfo, priorityTagArray)
	}
}

// Holds the info for subsequent converting.
type toBeConvertedFieldInfo struct {
	Value any // Found value by tag name or field name from input.
	*convertFieldInfo
}

func fieldReflectValue(v reflect.Value, fieldIndex []int) reflect.Value {
	for i, x := range fieldIndex {
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

// Only serving value
func (f *toBeConvertedFieldInfo) getFieldReflectValue(structValue reflect.Value) reflect.Value {
	if len(f.fieldIndex) == 1 {
		return structValue.Field(f.fieldIndex[0])
	}
	return fieldReflectValue(structValue, f.fieldIndex)
}

func (f *toBeConvertedFieldInfo) getOtherFieldReflectValue(structValue reflect.Value, fieldLevel int) reflect.Value {
	fieldIndex := f.otherFieldIndex[fieldLevel]
	if len(fieldIndex) == 1 {
		return structValue.Field(fieldIndex[0])
	}
	return fieldReflectValue(structValue, fieldIndex)
}

func getFieldTags(field reflect.StructField, priorityTags []string) (tags []string) {
	for _, tag := range priorityTags {
		value, ok := field.Tag.Lookup(tag)
		if ok {
			// If there's something else in the tag string,
			// it uses the first part which is split using char ','.
			// Example:
			// orm:"id, priority"
			// orm:"name, with:uid=id"
			array := strings.Split(value, ",")
			// json:",omitempty"
			trimmedTagName := strings.TrimSpace(array[0])
			tags = append(tags, trimmedTagName)
		}
	}
	tags = append(tags, field.Name)
	return tags
}

var (
	implUnmarshalText  = reflect.TypeOf((*iUnmarshalText)(nil)).Elem()
	implUnmarshalJson  = reflect.TypeOf((*iUnmarshalJSON)(nil)).Elem()
	implUnmarshalValue = reflect.TypeOf((*iUnmarshalValue)(nil)).Elem()
)

func checkTypeIsImplCommonInterface(field reflect.StructField) bool {
	isCommonInterface := false
	switch field.Type.String() {
	case "time.Time", "*time.Time":
	case "gtime.Time", "*gtime.Time":
		// default convert
	default:
		if field.Type.Kind() != reflect.Ptr {
			field.Type = reflect.PointerTo(field.Type)
		}
		switch {
		case field.Type.Implements(implUnmarshalText):
			isCommonInterface = true
		case field.Type.Implements(implUnmarshalJson):
			isCommonInterface = true
		case field.Type.Implements(implUnmarshalValue):
			isCommonInterface = true
		}
	}
	return isCommonInterface
}
