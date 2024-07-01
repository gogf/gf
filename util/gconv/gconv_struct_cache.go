// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"reflect"
	"strings"

	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/util/gtag"
)

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
	// 用来存储重名的字段
	otherFieldIndex [][]int
}

func (c *convertFieldInfo) FieldName() string {
	return c.tags[len(c.tags)-1]
}

type convertedStructInfo struct {
	// key   = field's name
	// value = field's tag
	fields map[string]*convertFieldInfo
}

// MergeField
// 需要 parentIndex 参数的原因
// 当一个结构体类型A已经被缓存，如果他被嵌入在B结构体中时，
// 那么此时在注册B结构体时，由于A结构体已经被缓存，字段的索引也是相对于A结构体来说的
// 举例：
//
//	type A struct {
//		Name string  // index = 0
//		Age int      // index = 1
//	}
//
//	type B struct {
//		Id int        // index = 0
//		A
//		可以把A字段展开到B结构体里面，由于A结构体已经被缓存，所以它的字段索引需要重新设置
//		Name string  // index = 0 => {1,0}
//		Age int      // index = 1 => {1,1}
//	}
//func (structInfo *convertedStructInfo) MergeField(fieldName string, fieldInfo *convertFieldInfo, parentIndex []int) *convertFieldInfo {
//	convertInfo, ok := structInfo.fields[fieldName]
//
//	newFieldInfo := &convertFieldInfo{}
//	*newFieldInfo = *fieldInfo
//	newFieldInfo.fieldIndex = append(parentIndex, newFieldInfo.fieldIndex...)
//
//	if !ok {
//		structInfo.fields[fieldName] = newFieldInfo
//		return newFieldInfo
//	}
//	if convertInfo.otherFieldIndex == nil {
//		convertInfo.otherFieldIndex = make([][]int, 0, 2)
//	}
//	convertInfo.otherFieldIndex = append(convertInfo.otherFieldIndex, newFieldInfo.fieldIndex)
//	return convertInfo
//}

func (structInfo *convertedStructInfo) AddField(field reflect.StructField, fieldIndex []int, priorityTags []string) *convertFieldInfo {
	convFieldInfo, ok := structInfo.fields[field.Name]
	if !ok {
		fieldInfo := &convertFieldInfo{
			isCommonInterface: checkTypeIsImplCommonInterface(field),
			fieldTypeName:     field.Type.String(),
			fieldIndex:        fieldIndex,
		}
		structInfo.fields[field.Name] = fieldInfo
		//if field.Anonymous {
		//	// Since the main logic has already been recursive,
		//	// there is no need for it here
		//	return fieldInfo
		//}
		fieldInfo.tags = getFieldTags(field, priorityTags)
		return fieldInfo
	}
	if convFieldInfo.otherFieldIndex == nil {
		convFieldInfo.otherFieldIndex = make([][]int, 0, 2)
	}
	convFieldInfo.otherFieldIndex = append(convFieldInfo.otherFieldIndex, fieldIndex)

	return convFieldInfo
}

type toBeConvertedStructInfo struct {
	// key   = field's name
	// value = field's tag
	fields map[string]*toBeConvertedFieldInfo
}

func (t *toBeConvertedStructInfo) AddField(fieldName string, fieldInfo *convertFieldInfo) {
	convertInfo, ok := t.fields[fieldName]
	if !ok {
		t.fields[fieldName] = &toBeConvertedFieldInfo{
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
	cacheConvertedStructsInfo = make(map[reflect.Type]*convertedStructInfo)

	implUnmarshalText  = reflect.TypeOf((*iUnmarshalText)(nil)).Elem()
	implUnmarshalJson  = reflect.TypeOf((*iUnmarshalJSON)(nil)).Elem()
	implUnmarshalValue = reflect.TypeOf((*iUnmarshalValue)(nil)).Elem()
)

func parseStruct(structType reflect.Type, priorityTag string) *toBeConvertedStructInfo {

	if structType.Kind() != reflect.Struct {
		return nil
	}
	// key=fieldName
	toBeConvertedFieldNameToInfo := &toBeConvertedStructInfo{
		fields: make(map[string]*toBeConvertedFieldInfo),
	}

	// 检查是否已经缓存，
	structInfo, ok := cacheConvertedStructsInfo[structType]
	if ok {
		for k, v := range structInfo.fields {
			toBeConvertedFieldNameToInfo.AddField(k, v)
		}
		return toBeConvertedFieldNameToInfo
	}
	structInfo = &convertedStructInfo{
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

	parseStructField(structType, parentIndex, structInfo, priorityTagArray)

	cacheConvertedStructsInfo[structType] = structInfo

	for k, v := range structInfo.fields {
		toBeConvertedFieldNameToInfo.AddField(k, v)
	}
	return toBeConvertedFieldNameToInfo
}

func parseStructField(structType reflect.Type, parentIndex []int, structInfo *convertedStructInfo, priorityTagArray []string) *convertedStructInfo {
	var (
		fieldName   string
		structField reflect.StructField
		fieldType   reflect.Type
	)
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
		// It is only recorded if the name has a fieldTag
		// TODO: If it's an anonymous field with a tag, doesn't it need to be recursive?
		structInfo.AddField(structField, append(parentIndex, i), priorityTagArray)

		parseStructField(fieldType, append(parentIndex, i), structInfo, priorityTagArray)
		// TODO 如果解析的结构体之前已经解析过，需要重新设置FieldIndex
		//for k, v := range convStructInfo.fields {
		//	newFieldInfo := structInfo.MergeField(k, v, append([]int{}, i))
		//	_ = newFieldInfo
		//	// toBeConvertedFieldNameToInfo.AddField(k, newFieldInfo)
		//}
	}
	return structInfo
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
