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

	// anyToTypeConvertMap for custom type converting from any to its reflect.Value.
	anyToTypeConvertMap map[reflect.Type]AnyConvertFunc

	// typeConverterFuncMarkMap is used to store whether field types are registered to custom conversions,
	// for enhance converting performance in runtime.
	typeConverterFuncMarkMap map[reflect.Type]struct{}

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
func NewCachedStructInfo(
	typeConverterFuncMarkMap map[reflect.Type]struct{},
	anyToTypeConvertMap map[reflect.Type]AnyConvertFunc,
) *CachedStructInfo {
	return &CachedStructInfo{
		tagOrFiledNameToFieldInfoMap: make(map[string]*CachedFieldInfo),
		fieldConvertInfos:            make([]*CachedFieldInfo, 0),
		typeConverterFuncMarkMap:     typeConverterFuncMarkMap,
		anyToTypeConvertMap:          anyToTypeConvertMap,
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
		IsCustomConvert:         csi.checkTypeHasCustomConvert(field.Type),
		PriorityTagAndFieldName: csi.genPriorityTagAndFieldName(field, priorityTags),
		RemoveSymbolsFieldName:  utils.RemoveSymbols(field.Name),
	}
	base.LastFuzzyKey.Store(field.Name)
	return &CachedFieldInfo{
		CachedFieldInfoBase: base,
	}
}

func (csi *CachedStructInfo) genFieldConvertFunc(fieldType reflect.Type) (convertFunc AnyConvertFunc) {
	if v := csi.anyToTypeConvertMap[fieldType]; v != nil {
		return v
	}

	var fieldTypeKind = fieldType.Kind()
	if fieldTypeKind == reflect.Ptr {
		convertFunc = csi.genFieldConvertFunc(fieldType.Elem())
		if convertFunc == nil {
			return nil
		}
		return csi.genPtrConvertFunc(convertFunc)
	}
	return nil
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

func (csi *CachedStructInfo) genPtrConvertFunc(convertFunc AnyConvertFunc) AnyConvertFunc {
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
	_, ok := csi.typeConverterFuncMarkMap[fieldType]
	return ok
}
