// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gmeta provides embedded meta data feature for struct.
package gmeta

import (
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/internal/structs"
	"strconv"
)

// Meta is used as an embedded attribute for struct to enabled meta data feature.
type Meta struct{}

const (
	metaAttributeName = "Meta"
)

var (
	metaDataCacheMap = gmap.NewStrAnyMap(true)
)

// Data retrieves and returns all meta data from `object`.
// It automatically parses the tag string from "Mata" attribute as its meta data.
func Data(object interface{}) map[string]interface{} {
	reflectType, err := structs.StructType(object)
	if err != nil {
		panic(err)
	}
	return metaDataCacheMap.GetOrSetFuncLock(reflectType.Signature(), func() interface{} {
		if field, ok := reflectType.FieldByName(metaAttributeName); ok {
			var (
				key  string
				tag  = field.Tag
				data = make(map[string]interface{})
			)
			for tag != "" {
				// Skip leading space.
				i := 0
				for i < len(tag) && tag[i] == ' ' {
					i++
				}
				tag = tag[i:]
				if tag == "" {
					break
				}
				// Scan to colon. A space, a quote or a control character is a syntax error.
				// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
				// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
				// as it is simpler to inspect the tag's bytes than the tag's runes.
				i = 0
				for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
					i++
				}
				if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
					break
				}
				key = string(tag[:i])
				tag = tag[i+1:]

				// Scan quoted string to find value.
				i = 1
				for i < len(tag) && tag[i] != '"' {
					if tag[i] == '\\' {
						i++
					}
					i++
				}
				if i >= len(tag) {
					break
				}
				quotedValue := string(tag[:i+1])
				tag = tag[i+1:]
				value, err := strconv.Unquote(quotedValue)
				if err != nil {
					panic(err)
				}
				data[key] = value
			}
			return data
		}
		return map[string]interface{}{}
	}).(map[string]interface{})
}

// Get retrieves and returns specified meta data by `key` from `object`.
func Get(object interface{}, key string) *gvar.Var {
	return gvar.New(Data(object)[key])
}
