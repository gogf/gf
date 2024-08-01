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
	"sync/atomic"
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
	// need to be cleared, otherwise there will be a bug
	for k := range m {
		delete(m, k)
	}
	poolUsedParamsKeyOrTagNameMap.Put(m)
}

type convertFieldInfoBase struct {
	// 字段的索引，可能是匿名嵌套的结构体，所以是[]int
	fieldIndex []int
	// 字段的tag(可能是conv,param,p,c,json之类的),
	// tags 包含字段的名字，并且在最后一个
	tags []string
	// 1.iUnmarshalValue
	// 2.iUnmarshalText
	// 3.iUnmarshalJSON
	// 实现了以上3种接口的类型
	// 除了 [time.Time] and [gtime.Time]
	isCommonInterface bool
	// 注册自定义转换的时候，比如func(src *int)(dest *string,err error)
	// 当结构体字段类型为string的时候，isCustomConvert 字段会为true
	// 表示此次转换有可能会是自定义转换，具体还需要进一步确定
	isCustomConvert bool
	structField     reflect.StructField
	// type Name struct{
	//   LastName string
	//   FirstName string
	// }
	// type User struct{
	//  Name
	//  LastName string
	//  FirstName string
	// }
	// 当结构体可能是类似于User结构体这种情况时
	// 只会存储两个字段LastName, FirstName使用不同的索引来代表不同的字段
	// 对于 LastName 字段来说
	// fieldIndex      = []int{0,1}
	// otherSameNameFieldIndex = [][]int{[]int{1}}长度只有1，因为只有一个重复的,且索引为1
	// 在赋值时会对这两个索引{0,1}和{1}都赋同样的值
	// 目前对于重复的字段可以做以下3种可能
	// 1.只设置第一个，后面重名的不设置
	// 2.只设置最后一个
	// 3.全部设置 (目前的做法)
	otherSameNameFieldIndex [][]int
	// 直接缓存字段的转换函数,对于简单的类型来说,相当于直接调用gconv.Int
	convFunc func(from any, to reflect.Value)
}

type convertFieldInfo struct {
	*convertFieldInfoBase
	// 表示上次模糊匹配到的字段名字，可以缓存下来
	// 如果用户没有设置tag之类的条件
	// 而且字段名都匹配不上map的key时，缓存这个非常有用，可以省掉模糊匹配的开销
	// lastFuzzKey string
	lastFuzzKey atomic.Value
	// 这个字段主要用在 bindStructWithLoopParamsMap 方法中，
	// 当map中同时存在一个字段的fieldName和tag时需要用到这个字段
	// 例如为以下情况时
	// field string `json:"name"`
	// map = {
	//	field:"f1",
	//	name :"n1",
	// }
	// 这里应该以name为准,
	// 在 bindStructWithLoopParamsMap 方法中，由于map的无序性，可能会导致先遍历到field
	// 这个字段更多的是表示优先级，即name的优先级比field的优先级高，即便之前已经设置过了
	isField bool
	// removeSymbolsFieldName = utils.RemoveSymbols(fieldName)
	removeSymbolsFieldName string
}

func (c *convertFieldInfo) FieldName() string {
	return c.tags[len(c.tags)-1]
}

func (c *convertFieldInfo) getFieldReflectValue(structValue reflect.Value) reflect.Value {
	if len(c.fieldIndex) == 1 {
		return structValue.Field(c.fieldIndex[0])
	}
	return fieldReflectValue(structValue, c.fieldIndex)
}

func (c *convertFieldInfo) getOtherFieldReflectValue(structValue reflect.Value, fieldLevel int) reflect.Value {
	fieldIndex := c.otherSameNameFieldIndex[fieldLevel]
	if len(fieldIndex) == 1 {
		return structValue.Field(fieldIndex[0])
	}
	return fieldReflectValue(structValue, fieldIndex)
}

type convertStructInfo struct {
	// This map field is mainly used in the [bindStructWithLoopParamsMap] method
	// key = field's name
	// Will save all field names and tags
	// for example：
	//	field string `json:"name"`
	// It will be stored twice
	fieldAndTagsMap map[string]*convertFieldInfo
	// Using slices here can speed up the loop
	fieldConvertInfos []*convertFieldInfo
}

func (structInfo *convertStructInfo) HasNoFields() bool {
	return len(structInfo.fieldAndTagsMap) == 0
}

func (structInfo *convertStructInfo) GetFieldInfo(fieldName string) *convertFieldInfo {
	return structInfo.fieldAndTagsMap[fieldName]
}

func (structInfo *convertStructInfo) AddField(field reflect.StructField, fieldIndex []int, priorityTags []string) {
	convFieldInfo, ok := structInfo.fieldAndTagsMap[field.Name]
	if !ok {
		baseInfo := &convertFieldInfoBase{
			isCommonInterface: checkTypeIsImplCommonInterface(field),
			structField:       field,
			fieldIndex:        fieldIndex,
			convFunc:          getFieldConvFunc(field.Type.String()),
			isCustomConvert:   checkTypeMaybeIsCustomConvert(field.Type),
			tags:              getFieldTags(field, priorityTags),
		}
		for _, tag := range baseInfo.tags {
			info := &convertFieldInfo{
				convertFieldInfoBase:   baseInfo,
				isField:                tag == field.Name,
				removeSymbolsFieldName: utils.RemoveSymbols(field.Name),
			}
			info.lastFuzzKey.Store(field.Name)
			structInfo.fieldAndTagsMap[tag] = info
			if info.isField {
				structInfo.fieldConvertInfos = append(structInfo.fieldConvertInfos, info)
			}
		}
		return
	}
	if convFieldInfo.otherSameNameFieldIndex == nil {
		convFieldInfo.otherSameNameFieldIndex = make([][]int, 0, 2)
	}
	convFieldInfo.otherSameNameFieldIndex = append(convFieldInfo.otherSameNameFieldIndex, fieldIndex)
	return
}

var (
	// Used to store whether field types are registered to custom conversions
	// For example:
	// func (src *TypeA) (dst *TypeB,err error)
	// This map will store TypeB for quick judgment during assignment
	customConvTypeMap = map[reflect.Type]struct{}{}
)

func registerCacheConvFieldCustomType(fieldType reflect.Type) {
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}
	customConvTypeMap[fieldType] = struct{}{}
}

func checkTypeMaybeIsCustomConvert(fieldType reflect.Type) bool {
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}
	_, ok := customConvTypeMap[fieldType]
	return ok
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
	case "int", "int8", "int16", "int32", "int64":
		convFunc = func(from any, to reflect.Value) {
			to.SetInt(Int64(from))
		}
	case "uint", "uint8", "uint16", "uint32", "uint64":
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
	case "bool":
		convFunc = func(from any, to reflect.Value) {
			to.SetBool(Bool(from))
		}
	case "[]byte":
		convFunc = func(from any, to reflect.Value) {
			to.SetBytes(Bytes(from))
		}
	default:
		return nil
	}
	return convFunc
}

var (
	// map[reflect.Type]*convertStructInfo
	cacheConvStructsInfo = sync.Map{}
)

func setCacheConvStructInfo(structType reflect.Type, info *convertStructInfo) {
	// Temporarily enabled as an experimental feature
	cacheConvStructsInfo.Store(structType, info)
}

func getCacheConvStructInfo(structType reflect.Type) (*convertStructInfo, bool) {
	// Temporarily enabled as an experimental feature
	v, ok := cacheConvStructsInfo.Load(structType)
	if ok {
		return v.(*convertStructInfo), ok
	}
	return nil, false
}

func getConvStructInfo(structType reflect.Type, priorityTag string) *convertStructInfo {
	if structType.Kind() != reflect.Struct {
		return nil
	}
	// Check if it has been cached
	structInfo, ok := getCacheConvStructInfo(structType)
	if ok {
		return structInfo
	}
	structInfo = &convertStructInfo{
		fieldAndTagsMap: make(map[string]*convertFieldInfo),
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
	return structInfo
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
		if structField.Tag != "" {
			structInfo.AddField(structField, append(parentIndex, i), priorityTagArray)
		}
		parseStruct(fieldType, append(parentIndex, i), structInfo, priorityTagArray)
	}
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
			if trimmedTagName != "" {
				tags = append(tags, trimmedTagName)
				break
			}
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
		// Implemented three types of interfaces that must be pointer types, otherwise it is meaningless
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
