// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package structcache

import (
	"reflect"
	"strings"

	"github.com/gogf/gf/v2/internal/utils"
)

// CachedStructInfo holds the cached info for certain struct.
type CachedStructInfo struct {
	// All sub attributes field info slice.
	fieldConvertInfos []*CachedFieldInfo

	converter *Converter

	// This map field is mainly used in the bindStructWithLoopParamsMap method
	// key = field's name
	// Will save all field names and PriorityTagAndFieldName
	// for example：
	//	field string `json:"name"`
	//
	// It will be stored twice, which keys are `name` and `field`.
	tagOrFiledNameToFieldInfoMap map[string]*CachedFieldInfo
}

// NewCachedStructInfo creates and returns a new CachedStructInfo object.
func NewCachedStructInfo(converter *Converter) *CachedStructInfo {
	return &CachedStructInfo{
		tagOrFiledNameToFieldInfoMap: make(map[string]*CachedFieldInfo),
		fieldConvertInfos:            make([]*CachedFieldInfo, 0),
		converter:                    converter,
	}
}

func (csi *CachedStructInfo) GetFieldConvertInfos() []*CachedFieldInfo {
	return csi.fieldConvertInfos
}

func (csi *CachedStructInfo) HasNoFields() bool {
	return len(csi.tagOrFiledNameToFieldInfoMap) == 0
}

func (csi *CachedStructInfo) GetFieldInfo(fieldName string) *CachedFieldInfo {
	return csi.tagOrFiledNameToFieldInfoMap[fieldName]
}

func (csi *CachedStructInfo) AddField(field reflect.StructField, fieldIndexes []int, priorityTags []string) {
	tagOrFieldNameArray := csi.genPriorityTagAndFieldName(field, priorityTags)
	for _, tagOrFieldName := range tagOrFieldNameArray {
		cachedFieldInfo, found := csi.tagOrFiledNameToFieldInfoMap[tagOrFieldName]
		newFieldInfo := csi.makeOrCopyCachedInfo(
			field, fieldIndexes, priorityTags, cachedFieldInfo, tagOrFieldName,
		)
		if newFieldInfo.IsField {
			csi.fieldConvertInfos = append(csi.fieldConvertInfos, newFieldInfo)
		}
		// if the field info by `tagOrFieldName` already cached,
		// it so adds this new field info to other same name field.
		if found {
			cachedFieldInfo.OtherSameNameField = append(cachedFieldInfo.OtherSameNameField, newFieldInfo)
		} else {
			csi.tagOrFiledNameToFieldInfoMap[tagOrFieldName] = newFieldInfo
		}
	}
}

func (csi *CachedStructInfo) makeOrCopyCachedInfo(
	field reflect.StructField, fieldIndexes []int, priorityTags []string,
	cachedFieldInfo *CachedFieldInfo,
	currTagOrFieldName string,
) (newFieldInfo *CachedFieldInfo) {
	if cachedFieldInfo == nil {
		// If the field is not cached, it creates a new one.
		newFieldInfo = csi.makeCachedFieldInfo(field, fieldIndexes, priorityTags)
		newFieldInfo.IsField = currTagOrFieldName == field.Name
		return
	}
	if cachedFieldInfo.StructField.Type != field.Type {
		// If the types are different, some information needs to be reset.
		newFieldInfo = csi.makeCachedFieldInfo(field, fieldIndexes, priorityTags)
	} else {
		// If the field types are the same.
		newFieldInfo = csi.copyCachedInfoWithFieldIndexes(cachedFieldInfo, fieldIndexes)
	}
	newFieldInfo.IsField = currTagOrFieldName == field.Name
	return
}

// copyCachedInfoWithFieldIndexes copies and returns a new CachedFieldInfo based on given CachedFieldInfo, but different
// FieldIndexes. Mainly used for copying fields with the same name and type.
func (csi *CachedStructInfo) copyCachedInfoWithFieldIndexes(
	cfi *CachedFieldInfo, fieldIndexes []int,
) *CachedFieldInfo {
	base := CachedFieldInfoBase{}
	base = *cfi.CachedFieldInfoBase
	base.FieldIndexes = fieldIndexes
	return &CachedFieldInfo{
		CachedFieldInfoBase: &base,
	}
}

func (csi *CachedStructInfo) makeCachedFieldInfo(
	field reflect.StructField, fieldIndexes []int, priorityTags []string,
) *CachedFieldInfo {
	base := &CachedFieldInfoBase{
		IsCommonInterface:       checkTypeIsCommonInterface(field),
		StructField:             field,
		FieldIndexes:            fieldIndexes,
		ConvertFunc:             csi.genFieldConvertFunc(field.Type),
		HasCustomConvert:        csi.checkTypeHasCustomConvert(field.Type),
		PriorityTagAndFieldName: csi.genPriorityTagAndFieldName(field, priorityTags),
		RemoveSymbolsFieldName:  utils.RemoveSymbols(field.Name),
	}
	base.LastFuzzyKey.Store(field.Name)
	return &CachedFieldInfo{
		CachedFieldInfoBase: base,
	}
}

func (csi *CachedStructInfo) genFieldConvertFunc(fieldType reflect.Type) (convertFunc AnyConvertFunc) {
	ptr := 0
	for fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
		ptr++
	}
	convertFunc = csi.converter.anyToTypeConvertMap[fieldType]
	if convertFunc == nil {
		// If the registered custom implementation cannot be found,
		// try to check if there is an implementation interface
		convertFunc = csi.converter.checkTypeImplInterface(fieldType)
	}
	// if the registered type is not found and
	// the corresponding interface is not implemented, return directly
	if convertFunc == nil {
		return nil
	}
	for i := 0; i < ptr; i++ {
		// If it is a pointer type, it needs to be packaged
		convertFunc = genPtrConvertFunc(convertFunc)
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

func genPtrConvertFunc(convertFunc AnyConvertFunc) AnyConvertFunc {
	return func(from any, to reflect.Value) error {
		if to.IsNil() {
			to.Set(reflect.New(to.Type().Elem()))
		}
		return convertFunc(from, to.Elem())
	}
}

func (csi *CachedStructInfo) checkTypeHasCustomConvert(fieldType reflect.Type) bool {
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}
	_, ok := csi.converter.typeConverterFuncMarkMap[fieldType]
	return ok
}
