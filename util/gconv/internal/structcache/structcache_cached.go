// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package structcache

import (
	"reflect"
	"sync"
	"time"

	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gtag"
)

// CommonConverter holds some converting functions of common types for internal usage.
type CommonConverter struct {
	Int64   func(any interface{}) int64
	Uint64  func(any interface{}) uint64
	String  func(any interface{}) string
	Float32 func(any interface{}) float32
	Float64 func(any interface{}) float64
	Time    func(any interface{}, format ...string) time.Time
	GTime   func(any interface{}, format ...string) *gtime.Time
	Bytes   func(any interface{}) []byte
	Bool    func(any interface{}) bool
}

var (
	// map[reflect.Type]*CachedStructInfo
	cachedStructsInfoMap = sync.Map{}

	// localCommonConverter holds some converting functions of common types for internal usage.
	localCommonConverter CommonConverter
)

// RegisterCommonConverter registers the CommonConverter for local usage.
func RegisterCommonConverter(commonConverter CommonConverter) {
	localCommonConverter = commonConverter
}

// GetCachedStructInfo retrieves or parses and returns a cached info for certain struct type.
func GetCachedStructInfo(
	structType reflect.Type,
	priorityTag string,
) *CachedStructInfo {
	if structType.Kind() != reflect.Struct {
		return nil
	}
	// check if it has been cached.
	structInfo, ok := getCachedConvertStructInfo(structType)
	if ok {
		return structInfo
	}

	// it parses and generates a cache info for given struct type.
	structInfo = &CachedStructInfo{
		tagOrFiledNameToFieldInfoMap: make(map[string]*CachedFieldInfo),
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
	setCachedConvertStructInfo(structType, structInfo)
	return structInfo
}

func setCachedConvertStructInfo(structType reflect.Type, info *CachedStructInfo) {
	// Temporarily enabled as an experimental feature
	cachedStructsInfoMap.Store(structType, info)
}

func getCachedConvertStructInfo(structType reflect.Type) (*CachedStructInfo, bool) {
	// Temporarily enabled as an experimental feature
	v, ok := cachedStructsInfoMap.Load(structType)
	if ok {
		return v.(*CachedStructInfo), ok
	}
	return nil, false
}

func parseStruct(
	structType reflect.Type,
	fieldIndexes []int,
	structInfo *CachedStructInfo,
	priorityTagArray []string,
) {
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
		if fieldType.String() == "gmeta.Meta" {
			continue
		}
		copyFieldIndexes := make([]int, len(fieldIndexes))
		copy(copyFieldIndexes, fieldIndexes)

		// normal basic attributes.
		if structField.Anonymous {
			// handle struct attributes, it might be struct/*struct embedded..
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}
			if fieldType.Kind() != reflect.Struct {
				continue
			}
			if structField.Tag != "" {
				// Do not add anonymous structures without tags
				structInfo.AddField(structField, append(copyFieldIndexes, i), priorityTagArray)
			}
			parseStruct(fieldType, append(copyFieldIndexes, i), structInfo, priorityTagArray)
			continue
		}

		// Do not directly use append(fieldIndexes, i)
		// When the structure is nested deeply, it may lead to bugs,
		// which are caused by the slice expansion mechanism
		// So it is necessary to allocate a separate index for each field
		// See details https://github.com/gogf/gf/issues/3789
		structInfo.AddField(structField, append(copyFieldIndexes, i), priorityTagArray)
	}
}
