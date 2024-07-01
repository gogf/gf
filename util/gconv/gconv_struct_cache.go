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
)

func poolGetUsedParamsKeyOrTagNameMap() map[string]struct{} {
	return poolUsedParamsKeyOrTagNameMap.Get().(map[string]struct{})
}

func poolPutUsedParamsKeyOrTagNameMap(m map[string]struct{}) {
	// 需要清空，不然会有bug
	for k := range m {
		delete(m, k)
	}
	poolUsedParamsKeyOrTagNameMap.Put(m)
}

type convertFieldInfo struct {
	// 字段的索引，有可能是一个嵌套的结构体，所以需要[]int
	fieldIndex []int
	// 包括字段的名字，且字段名是最后一个
	tags []string
	// iUnmarshalValue
	// iUnmarshalText
	// iUnmarshalJSON
	// 实现了以上三种接口之一的类型，除了[time.Time]和[gtime.Time]
	isCommonInterface bool
	fieldTypeName     string
	// 用来存储重名的字段索引
	// 一般用于嵌套的结构体
	otherFieldIndex [][]int
	// 缓存类型转换的函数
	// 比如字段的类型是int ，那么直接缓存gconv.Int函数
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

func getFieldConvFunc(fieldType string) (convFunc func(from any, to reflect.Value)) {
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
	if convFunc != nil {
		if fieldType[0] == '*' {
			return func(from any, to reflect.Value) {
				convFunc(from, to.Elem())
			}
		}
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

	implUnmarshalText  = reflect.TypeOf((*iUnmarshalText)(nil)).Elem()
	implUnmarshalJson  = reflect.TypeOf((*iUnmarshalJSON)(nil)).Elem()
	implUnmarshalValue = reflect.TypeOf((*iUnmarshalValue)(nil)).Elem()
)

func cacheConvStructInfo(structType reflect.Type, info *convertStructInfo) {
	cacheConvStructsInfo[structType] = info
}

func getConvStructInfo(structType reflect.Type, priorityTag string) *toBeConvertedStructInfo {
	if structType.Kind() != reflect.Struct {
		return nil
	}
	// key=fieldName
	toBeConvertedFieldNameToInfo := &toBeConvertedStructInfo{
		fields: make(map[string]toBeConvertedFieldInfo),
	}
	// 检查是否已经缓存，
	structInfo, ok := cacheConvStructsInfo[structType]
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
	cacheConvStructInfo(structType, structInfo)
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
	// TODO 查找缓存中是否已经缓存了结构体，如果缓存了，可以复用一些信息，但是需要重新设置[FieldIndex]，暂时不做实现，因为有些复杂
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
	// Value 不需要存储，或者单独使用两个结构体来
	// 存储 Value 和下面所有的字段
	Value any // Found value by tag name or field name from input.
	*convertFieldInfo
}

// 只为value服务
func (f *toBeConvertedFieldInfo) getFieldReflectValue(structValue reflect.Value) reflect.Value {
	if len(f.fieldIndex) == 1 {
		return structValue.Field(f.fieldIndex[0])
	}
	v := structValue
	for i, x := range f.fieldIndex {
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

func (f *toBeConvertedFieldInfo) getOtherFieldReflectValue(structValue reflect.Value, fieldLevel int) reflect.Value {
	fieldIndex := f.otherFieldIndex[fieldLevel]
	if len(fieldIndex) == 1 {
		return structValue.Field(fieldIndex[0])
	}
	v := structValue
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

func checkTypeIsImplCommonInterface(field reflect.StructField) bool {
	isCommonInterface := false
	switch field.Type.String() {
	case "time.Time", "*time.Time":
	case "gtime.Time", "*gtime.Time":
		// default convert
	default:
		// TODO slice 和 map 类型是否需要特殊处理
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
