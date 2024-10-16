// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gmeta provides embedded meta data feature for struct.
package gmeta

import (
	"sync"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/gstructs"
)

// Meta is used as an embedded attribute for struct to enabled metadata feature.
type Meta struct{}

const (
	metaAttributeName = "Meta"       // metaAttributeName is the attribute name of metadata in struct.
	metaTypeName      = "gmeta.Meta" // metaTypeName is for type string comparison.
)

// cachedMetadata stores the parsed metadata for struct types.
var cachedMetadata = sync.Map{}

// Data retrieves and returns all metadata from `object`.
func Data(object interface{}) map[string]string {
	reflectType, err := gstructs.StructType(object)
	if err != nil {
		return nil
	}

	if cachedData, ok := cachedMetadata.Load(reflectType.Type); ok {
		return cachedData.(map[string]string)
	}

	var metadata map[string]string
	if field, ok := reflectType.FieldByName(metaAttributeName); ok {
		if field.Type.String() == metaTypeName {
			metadata = gstructs.ParseTag(string(field.Tag))
		}
	}

	if metadata == nil {
		metadata = map[string]string{}
	}

	cachedMetadata.Store(reflectType.Type, metadata)
	return metadata
}

// Get retrieves and returns specified metadata by `key` from `object`.
func Get(object interface{}, key string) *gvar.Var {
	v, ok := Data(object)[key]
	if !ok {
		return nil
	}
	return gvar.New(v)
}
