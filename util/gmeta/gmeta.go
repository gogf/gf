// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gmeta provides embedded meta data feature for struct.
package gmeta

import (
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/internal/structs"
)

// Meta is used as an embedded attribute for struct to enabled metadata feature.
type Meta struct{}

const (
	// metaAttributeName is the attribute name of metadata in struct.
	metaAttributeName = "Meta"
)

// Data retrieves and returns all metadata from `object`.
// It automatically parses and caches the tag string from "Mata" attribute as its metadata.
func Data(object interface{}) map[string]interface{} {
	reflectType, err := structs.StructType(object)
	if err != nil {
		panic(err)
	}
	if field, ok := reflectType.FieldByName(metaAttributeName); ok {
		var (
			tags = structs.ParseTag(string(field.Tag))
			data = make(map[string]interface{}, len(tags))
		)
		for k, v := range tags {
			data[k] = v
		}
		return data
	}
	return map[string]interface{}{}
}

// Get retrieves and returns specified metadata by `key` from `object`.
func Get(object interface{}, key string) *gvar.Var {
	v := Data(object)[key]
	if v == nil {
		return nil
	}
	return gvar.New(v)
}
