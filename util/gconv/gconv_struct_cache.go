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

type convertFieldInfoBase struct {
	// The index of a field may be a nested structure, so []int is required
	fieldIndex []int
	// All tags in the field include the field name, and the field name is the last one
	tags []string
	// 1.iUnmarshalValue
	// 2.iUnmarshalText
	// 3.iUnmarshalJSON
	// Implemented one of the three types of interfaces mentioned above,
	// except for [time.Time] and [gtime.Time]
	isCommonInterface bool
	isCustomConvert   bool
	structField       reflect.StructField
	// Field index used to store duplicate names
	// Generally used for nested structures
	otherFieldIndex [][]int
	// Cache type conversion function
	// For example:
	// if the type of the field is int, then directly cache the [gconv.Int] function
	convFunc func(from any, to reflect.Value)
}

type convertFieldInfo struct {
	*convertFieldInfoBase
	// lastFuzzKey string
	lastFuzzKey atomic.Value
	// The essence of this field is that there is only one unique key in the map,
	// field string `json:"name"`
	// map = {
	//	field:"f1",
	//	name :"n1",
	// }
	// In this case, name shall prevail
	isField                bool
	removeSymbolsFieldName string
}

func (c *convertFieldInfo) FieldName() string {
	return c.tags[len(c.tags)-1]
}

// Only serving value
func (c *convertFieldInfo) getFieldReflectValue(structValue reflect.Value) reflect.Value {
	if len(c.fieldIndex) == 1 {
		return structValue.Field(c.fieldIndex[0])
	}
	return fieldReflectValue(structValue, c.fieldIndex)
}

func (c *convertFieldInfo) getOtherFieldReflectValue(structValue reflect.Value, fieldLevel int) reflect.Value {
	fieldIndex := c.otherFieldIndex[fieldLevel]
	if len(fieldIndex) == 1 {
		return structValue.Field(fieldIndex[0])
	}
	return fieldReflectValue(structValue, fieldIndex)
}

type convertStructInfo struct {
	// key = field's name
	// Will save all field names and tags
	// for example：
	//	field string `json:"name"`
	// It will be stored twice
	fieldAndTagsMap map[string]*convertFieldInfo
	// It will only be stored according to the name of the field,
	// and only a few fields will be stored in the structure
	fieldNamesMap map[string]*convertFieldInfo
}

func (structInfo *convertStructInfo) NoFields() bool {
	return len(structInfo.fieldAndTagsMap) == 0
}

func (structInfo *convertStructInfo) GetFieldInfo(fieldName string) *convertFieldInfo {
	v := structInfo.fieldAndTagsMap[fieldName]
	return v
}

func (structInfo *convertStructInfo) AddField(field reflect.StructField, fieldIndex []int, priorityTags []string) {
	convFieldInfo, ok := structInfo.fieldAndTagsMap[field.Name]
	if !ok {
		baseInfo := &convertFieldInfoBase{
			isCommonInterface: checkTypeIsImplCommonInterface(field),
			structField:       field,
			fieldIndex:        fieldIndex,
			convFunc:          getFieldConvFunc(field.Type.String()),
			isCustomConvert:   fieldTypeIsCustomConvertType(field.Type),
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
				structInfo.fieldNamesMap[field.Name] = info
			}
		}
		return
	}
	if convFieldInfo.otherFieldIndex == nil {
		convFieldInfo.otherFieldIndex = make([][]int, 0, 2)
	}
	convFieldInfo.otherFieldIndex = append(convFieldInfo.otherFieldIndex, fieldIndex)
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

func fieldTypeIsCustomConvertType(fieldType reflect.Type) bool {
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
	cacheConvStructsInfo = sync.Map{}
)

func setCacheConvStructInfo(structType reflect.Type, info *convertStructInfo) {
	// Temporarily enabled as an experimental feature
	if convCacheExperiment {
		cacheConvStructsInfo.Store(structType, info)
	}
}

func getCacheConvStructInfo(structType reflect.Type) (*convertStructInfo, bool) {
	// Temporarily enabled as an experimental feature
	if convCacheExperiment {
		v, ok := cacheConvStructsInfo.Load(structType)
		if ok {
			return v.(*convertStructInfo), ok
		}
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
		fieldNamesMap:   make(map[string]*convertFieldInfo),
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
