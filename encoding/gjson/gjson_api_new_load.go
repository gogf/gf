// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson

import (
	"reflect"

	"github.com/gogf/gf/v2/internal/reflection"
	"github.com/gogf/gf/v2/internal/rwmutex"
	"github.com/gogf/gf/v2/util/gconv"
)

// New creates a Json object with any variable type of `data`, but `data` should be a map
// or slice for data access reason, or it will make no sense.
//
// The parameter `safe` specifies whether using this Json object in concurrent-safe context,
// which is false in default.
func New(data any, safe ...bool) *Json {
	return NewWithTag(data, string(ContentTypeJSON), safe...)
}

// NewWithTag creates a Json object with any variable type of `data`, but `data` should be a map
// or slice for data access reason, or it will make no sense.
//
// The parameter `tags` specifies priority tags for struct conversion to map, multiple tags joined
// with char ','.
//
// The parameter `safe` specifies whether using this Json object in concurrent-safe context, which
// is false in default.
func NewWithTag(data any, tags string, safe ...bool) *Json {
	option := Options{
		Tags: tags,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return NewWithOptions(data, option)
}

// NewWithOptions creates a Json object with any variable type of `data`, but `data` should be a map
// or slice for data access reason, or it will make no sense.
func NewWithOptions(data any, options Options) *Json {
	var j *Json
	switch result := data.(type) {
	case []byte:
		if r, err := loadContentWithOptions(result, options); err == nil {
			j = r
			break
		}
		j = &Json{
			p:  &data,
			c:  byte(defaultSplitChar),
			vc: false,
		}
	case string:
		if r, err := loadContentWithOptions([]byte(result), options); err == nil {
			j = r
			break
		}
		j = &Json{
			p:  &data,
			c:  byte(defaultSplitChar),
			vc: false,
		}
	default:
		var (
			pointedData any
			reflectInfo = reflection.OriginValueAndKind(data)
		)
		switch reflectInfo.OriginKind {
		case reflect.Slice, reflect.Array:
			pointedData = gconv.Interfaces(data)

		case reflect.Map:
			pointedData = gconv.MapDeep(data, options.Tags)

		case reflect.Struct:
			if v, ok := data.(iVal); ok {
				return NewWithOptions(v.Val(), options)
			}
			pointedData = gconv.MapDeep(data, options.Tags)

		default:
			pointedData = data
		}
		j = &Json{
			p:  &pointedData,
			c:  byte(defaultSplitChar),
			vc: false,
		}
	}
	j.mu = rwmutex.Create(options.Safe)
	return j
}
