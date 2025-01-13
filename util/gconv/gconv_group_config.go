// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"github.com/gogf/gf/v2/util/gconv/internal/structcache"
)

type ConvertConfig = structcache.ConvertConfig

var defaultConfig = structcache.GetDefaultConfig()

func NewConvertConfig(name string) *ConvertConfig {
	return structcache.NewConvertConfig(name)
}

func ScanByConvertConfig(config *ConvertConfig, srcValue interface{}, dstPointer interface{}, paramKeyToAttrMap ...map[string]string) (err error) {
	return scan(config, srcValue, dstPointer, paramKeyToAttrMap...)
}

func StructByConfig(config *ConvertConfig, params interface{}, pointer interface{}, paramKeyToAttrMap ...map[string]string) (err error) {
	return scan(config, params, pointer, paramKeyToAttrMap...)
}

func StructTagByConfig(config *ConvertConfig, params interface{}, pointer interface{}, priorityTag string) (err error) {
	return doStruct(params, pointer, nil, priorityTag, config)
}
