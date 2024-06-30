// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/util/gtag"
)

type convertFieldInfo struct {
	// 字段的索引，有可能是一个嵌套的结构体，所以需要[]int
	fieldIndex [][]int
	// 包括字段的名字，且字段名是最后一个
	tags []string
	// iUnmarshalValue
	// iUnmarshalText
	// iUnmarshalJSON
	// 实现了以上三种接口之一的类型，除了[time.Time]和[gtime.Time]
	isCommonInterface bool
	fieldTypeName     string
	fieldLevel        int
}

type convertedStructInfo struct {
	// key   = field's name
	// value = field's tag
	// TODO 在多层嵌套时字段可能会重复，需要另外的手段来区别
	fields map[string]*convertFieldInfo
}

type toBeConvertedStructInfo struct {
	// key   = field's name
	// value = field's tag
	// TODO 在多层嵌套时字段可能会重复，需要另外的手段来区别
	fields map[string]*toBeConvertedFieldInfo
}

func (structInfo *convertedStructInfo) AddField(field reflect.StructField, fieldIndex [][]int, priorityTags []string) *convertFieldInfo {
	_, ok := structInfo.fields[field.Name]
	if ok {
		panic(fmt.Sprintf("Field `%s` already exists", field.Name))
	}
	fieldInfo := &convertFieldInfo{
		isCommonInterface: checkTypeIsImplCommonInterface(field),
		fieldTypeName:     field.Type.String(),
		fieldIndex:        fieldIndex,
	}
	structInfo.fields[field.Name] = fieldInfo
	if field.Anonymous {
		// Since the main logic has already been recursive,
		// there is no need for it here
		return fieldInfo
	}
	fieldInfo.tags = getFieldTags(field, priorityTags)
	return fieldInfo
}

var (
	cacheConvertedStructsInfo = make(map[reflect.Type]*convertedStructInfo)

	implUnmarshalText  = reflect.TypeOf((*iUnmarshalText)(nil)).Elem()
	implUnmarshalJson  = reflect.TypeOf((*iUnmarshalJSON)(nil)).Elem()
	implUnmarshalValue = reflect.TypeOf((*iUnmarshalValue)(nil)).Elem()
)

func parseStruct(
	pointerElem reflect.Value,
	structType reflect.Type,
	priorityTag string, parentIndex []int,
	toBeConvertedFieldNameToInfo *toBeConvertedStructInfo) {

	if structType.Kind() != reflect.Struct {
		return
	}
	// 检查是否已经缓存，
	structInfo, ok := cacheConvertedStructsInfo[structType]
	if ok {
		// 如果缓存了，直接退出，不需要再次解析
		for k, v := range structInfo.fields {
			toBeConvertedFieldNameToInfo.fields[k] = &toBeConvertedFieldInfo{
				Value:            nil,
				convertFieldInfo: v,
			}
		}
		return
	}
	structInfo = &convertedStructInfo{
		fields: make(map[string]*convertFieldInfo),
	}

	var (
		priorityTagArray []string
		fieldName        string
		structField      reflect.StructField
		fieldType        reflect.Type
	)
	if priorityTag != "" {
		priorityTagArray = append(utils.SplitAndTrim(priorityTag, ","), gtag.StructTagPriority...)
	} else {
		priorityTagArray = gtag.StructTagPriority
	}
	for i := 0; i < structType.NumField(); i++ {
		structField = structType.Field(i)
		fieldType = structField.Type
		fieldName = structField.Name
		// Only do converting to public attributes.
		if !utils.IsLetterUpper(fieldName[0]) {
			continue
		}
		/////////可以用下面的[getFieldTags]来实现////////////////////////
		// var fieldTagName = getTagNameFromField(structField, priorityTagArray)
		//////////////////////////
		if structField.Anonymous == false {

			toBeConvertedFieldNameToInfo.fields[fieldName] = &toBeConvertedFieldInfo{
				convertFieldInfo: structInfo.AddField(structField, append(parentIndex, i), priorityTagArray),
			}
			continue
		}

		// Maybe it's struct/*struct embedded.
		// TODO 暂时不解析接口，使用特定的字段标识，直到赋值时才解析接口字段
		if fieldType.Kind() == reflect.Interface {
			// 如果是接口，不需要进入
			//fieldValue := pointerElem.FieldByIndex(append(parentIndex, i))
			//if fieldValue.IsValid() == false || fieldValue.IsNil() {
			//	// empty interface or nil
			//	continue
			//}
			//// interface => struct
			//fieldValue = fieldValue.Elem()
			//if fieldValue.Kind() == reflect.Ptr {
			//	fieldValue = fieldValue.Elem()
			//}
			//fieldType = fieldValue.Type()
		} else {
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}
		}

		if fieldType.Kind() != reflect.Struct {
			continue
		}

		// type Name struct {
		//    LastName  string `json:"lastName"`
		//    FirstName string `json:"firstName"`
		// }
		//
		// type User struct {
		//     Name `json:"name"`
		//     // ...
		// }
		//
		// It is only recorded if the name has a fieldTag
		// TODO: If it's an anonymous field with a tag, doesn't it need to be recursive?

		toBeConvertedFieldNameToInfo.fields[fieldName] = &toBeConvertedFieldInfo{
			Value:            nil,
			convertFieldInfo: structInfo.AddField(structField, append(parentIndex, i), priorityTagArray),
		}

		parseStruct(pointerElem, fieldType, priorityTag, append(parentIndex, i), toBeConvertedFieldNameToInfo)

	}
	/////////////////////////////////////////////////////
	cacheConvertedStructsInfo[fieldType] = structInfo
}

// 可以用下面的[getFieldTags]来实现
func getTagNameFromField(field reflect.StructField, priorityTags []string) string {
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
			return trimmedTagName
		}
	}
	return ""
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

// Holds the info for subsequent converting.
type toBeConvertedFieldInfo struct {
	// Value 不需要存储，或者单独使用两个结构体来
	// 存储 Value 和下面所有的字段
	Value any // Found value by tag name or field name from input.
	*convertFieldInfo
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
		// slice 和 map 类型是否需要特殊处理
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
