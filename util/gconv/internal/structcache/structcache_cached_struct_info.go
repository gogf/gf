// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package structcache

import (
	"reflect"
	"strings"
	"time"

	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gtime"
)

// CachedStructInfo holds the cached info for certain struct.
type CachedStructInfo struct {
	// This map field is mainly used in the bindStructWithLoopParamsMap method
	// key = field's name
	// Will save all field names and PriorityTagAndFieldName
	// for example：
	//	field string `json:"name"`
	//
	// It will be stored twice, which keys are `name` and `field`.
	tagOrFiledNameToFieldInfoMap map[string]*CachedFieldInfo

	// All sub attributes field info slice.
	FieldConvertInfos []*CachedFieldInfo
}

func (csi *CachedStructInfo) HasNoFields() bool {
	return len(csi.tagOrFiledNameToFieldInfoMap) == 0
}

func (csi *CachedStructInfo) GetFieldInfo(fieldName string) *CachedFieldInfo {
	return csi.tagOrFiledNameToFieldInfoMap[fieldName]
}

func (csi *CachedStructInfo) AddField(field reflect.StructField, fieldIndexes []int, priorityTags []string) {
	alreadyExistFieldInfo, ok := csi.tagOrFiledNameToFieldInfoMap[field.Name]
	if !ok {
		priorityTagAndFieldName := csi.genPriorityTagAndFieldName(field, priorityTags)
		newFieldInfoBase := &CachedFieldInfoBase{
			IsCommonInterface:       checkTypeIsCommonInterface(field),
			StructField:             field,
			FieldIndexes:            fieldIndexes,
			ConvertFunc:             csi.genFieldConvertFunc(field.Type.String()),
			IsCustomConvert:         csi.checkTypeHasCustomConvert(field.Type),
			PriorityTagAndFieldName: priorityTagAndFieldName,
			RemoveSymbolsFieldName:  utils.RemoveSymbols(field.Name),
		}
		newFieldInfoBase.LastFuzzyKey.Store(field.Name)
		for _, tagOrFieldName := range priorityTagAndFieldName {
			newFieldInfo := &CachedFieldInfo{
				CachedFieldInfoBase: newFieldInfoBase,
				IsField:             tagOrFieldName == field.Name,
			}
			csi.tagOrFiledNameToFieldInfoMap[tagOrFieldName] = newFieldInfo
			if newFieldInfo.IsField {
				csi.FieldConvertInfos = append(csi.FieldConvertInfos, newFieldInfo)
			}
		}
		return
	}
	if alreadyExistFieldInfo.OtherSameNameFieldIndex == nil {
		alreadyExistFieldInfo.OtherSameNameFieldIndex = make([][]int, 0, 2)
	}
	alreadyExistFieldInfo.OtherSameNameFieldIndex = append(
		alreadyExistFieldInfo.OtherSameNameFieldIndex,
		fieldIndexes,
	)
	return
}

func (csi *CachedStructInfo) genFieldConvertFunc(fieldType string) (convertFunc func(from any, to reflect.Value)) {
	if fieldType[0] == '*' {
		convertFunc = csi.genFieldConvertFunc(fieldType[1:])
		if convertFunc == nil {
			return nil
		}
		return csi.genPtrConvertFunc(convertFunc)
	}
	switch fieldType {
	case "int", "int8", "int16", "int32", "int64":
		convertFunc = func(from any, to reflect.Value) {
			to.SetInt(localCommonConverter.Int64(from))
		}
	case "uint", "uint8", "uint16", "uint32", "uint64":
		convertFunc = func(from any, to reflect.Value) {
			to.SetUint(localCommonConverter.Uint64(from))
		}
	case "string":
		convertFunc = func(from any, to reflect.Value) {
			to.SetString(localCommonConverter.String(from))
		}
	case "float32":
		convertFunc = func(from any, to reflect.Value) {
			to.SetFloat(float64(localCommonConverter.Float32(from)))
		}
	case "float64":
		convertFunc = func(from any, to reflect.Value) {
			to.SetFloat(localCommonConverter.Float64(from))
		}
	case "Time", "time.Time":
		convertFunc = func(from any, to reflect.Value) {
			*to.Addr().Interface().(*time.Time) = localCommonConverter.Time(from)
		}
	case "GTime", "gtime.Time":
		convertFunc = func(from any, to reflect.Value) {
			v := localCommonConverter.GTime(from)
			if v == nil {
				v = gtime.New()
			}
			*to.Addr().Interface().(*gtime.Time) = *v
		}
	case "bool":
		convertFunc = func(from any, to reflect.Value) {
			to.SetBool(localCommonConverter.Bool(from))
		}
	case "[]byte":
		convertFunc = func(from any, to reflect.Value) {
			to.SetBytes(localCommonConverter.Bytes(from))
		}
	default:
		return nil
	}
	return convertFunc
}

func (csi *CachedStructInfo) genPriorityTagAndFieldName(
	field reflect.StructField, priorityTags []string,
) (priorityTagAndFieldName []string) {
	for _, tag := range priorityTags {
		value, ok := field.Tag.Lookup(tag)
		if ok {
			// If there's something else in the tag string,
			// it uses the first part which is split using char ','.
			// Example:
			// orm:"id, priority"
			// orm:"name, with:uid=id"
			tagValueItems := strings.Split(value, ",")
			// json:",omitempty"
			trimmedTagName := strings.TrimSpace(tagValueItems[0])
			if trimmedTagName != "" {
				priorityTagAndFieldName = append(priorityTagAndFieldName, trimmedTagName)
				break
			}
		}
	}
	priorityTagAndFieldName = append(priorityTagAndFieldName, field.Name)
	return
}

func (csi *CachedStructInfo) checkTypeHasCustomConvert(fieldType reflect.Type) bool {
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}
	_, ok := customConvertTypeMap[fieldType]
	return ok
}

func (csi *CachedStructInfo) genPtrConvertFunc(
	convertFunc func(from any, to reflect.Value),
) func(from any, to reflect.Value) {
	return func(from any, to reflect.Value) {
		if to.IsNil() {
			to.Set(reflect.New(to.Type().Elem()))
		}
		convertFunc(from, to.Elem())
	}
}
