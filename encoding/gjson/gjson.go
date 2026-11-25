// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gjson provides convenient API for JSON/XML/INI/YAML/TOML data handling.
package gjson

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/reflection"
	"github.com/gogf/gf/v2/internal/rwmutex"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

type ContentType = string

const (
	ContentTypeJSON       ContentType = `json`
	ContentTypeJs         ContentType = `js`
	ContentTypeXML        ContentType = `xml`
	ContentTypeIni        ContentType = `ini`
	ContentTypeYaml       ContentType = `yaml`
	ContentTypeYml        ContentType = `yml`
	ContentTypeToml       ContentType = `toml`
	ContentTypeProperties ContentType = `properties`
)

const (
	// Separator char for hierarchical data access.
	defaultSplitChar = '.'
)

// Json is the customized JSON struct.
type Json struct {
	mu rwmutex.RWMutex

	// Pointer for hierarchical data access, it's the root of data in default.
	p *any

	// Char separator('.' in default).
	c byte

	// Violence Check(false in default),
	// which is used to access data when the hierarchical data key contains separator char.
	vc bool
}

// Options for Json object creating/loading.
type Options struct {
	// Mark this object is for in concurrent-safe usage. This is especially for Json object creating.
	Safe bool

	// Custom priority tags for decoding, eg: "json,yaml,MyTag".
	// This is specially for struct parsing into Json object.
	Tags string

	// Type specifies the data content type, eg: json, xml, yaml, toml, ini.
	Type ContentType

	// StrNumber causes the Decoder to unmarshal a number into an any as a string instead of as a float64.
	// This is specially for json content parsing into Json object.
	StrNumber bool
}

// iInterfaces is used for type assert api for Interfaces().
type iInterfaces interface {
	Interfaces() []any
}

// iMapStrAny is the interface support for converting struct parameter to map.
type iMapStrAny interface {
	MapStrAny() map[string]any
}

// iVal is the interface for underlying any retrieving.
type iVal interface {
	Val() any
}

// setValue sets `value` to `j` by `pattern`.
// Note:
// 1. If value is nil and removed is true, means deleting this value;
// 2. Missing nodes on the way are created: a slice for a numeric key and a map otherwise;
func (j *Json) setValue(pattern string, value any, removed bool) error {
	var (
		err   error
		array = strings.Split(pattern, string(j.c))
	)
	if value, err = j.convertValue(value); err != nil {
		return err
	}
	j.mu.Lock()
	defer j.mu.Unlock()
	// Initialization checks.
	if *j.p == nil {
		if gstr.IsNumeric(array[0]) {
			*j.p = make([]any, 0)
		} else {
			*j.p = make(map[string]any)
		}
	}
	if removed && value == nil {
		*j.p = j.removeNodeValue(*j.p, array)
		return nil
	}
	node, err := j.setNodeValue(*j.p, array, value)
	if err != nil {
		return err
	}
	*j.p = node
	return nil
}

// setNodeValue sets `value` to `node` by the hierarchical `keys` and returns the node.
// The returned node differs from `node` if `node` had to be grown or replaced, so that the
// caller must always store it back into its parent.
func (j *Json) setNodeValue(node any, keys []string, value any) (any, error) {
	if len(keys) == 0 {
		return value, nil
	}
	var key = keys[0]
	switch v := node.(type) {
	case map[string]any:
		if v == nil {
			v = make(map[string]any)
		}
		child, err := j.setNodeValue(v[key], keys[1:], value)
		if err != nil {
			return nil, err
		}
		v[key] = child
		return v, nil

	case []any:
		// A string key replaces the slice with a map below.
		if gstr.IsNumeric(key) {
			index, err := j.sliceIndex(key)
			if err != nil {
				return nil, err
			}
			if index >= len(v) {
				s := make([]any, index+1)
				copy(s, v)
				v = s
			}
			child, err := j.setNodeValue(v[index], keys[1:], value)
			if err != nil {
				return nil, err
			}
			v[index] = child
			return v, nil
		}
	}
	// The node does not exist or cannot hold `key`, so it is replaced by a new node.
	if gstr.IsNumeric(key) {
		index, err := j.sliceIndex(key)
		if err != nil {
			return nil, err
		}
		return j.setNodeValue(make([]any, index+1), keys, value)
	}
	return j.setNodeValue(make(map[string]any), keys, value)
}

// removeNodeValue deletes the value of the hierarchical `keys` from `node` and returns the node.
// Nothing is changed if the value does not exist.
func (j *Json) removeNodeValue(node any, keys []string) any {
	var key = keys[0]
	switch v := node.(type) {
	case map[string]any:
		child, ok := v[key]
		if !ok {
			return v
		}
		if len(keys) == 1 {
			delete(v, key)
		} else {
			v[key] = j.removeNodeValue(child, keys[1:])
		}
		return v

	case []any:
		if !gstr.IsNumeric(key) {
			return v
		}
		index, err := j.sliceIndex(key)
		if err != nil || index >= len(v) {
			return v
		}
		if len(keys) == 1 {
			return append(v[:index], v[index+1:]...)
		}
		v[index] = j.removeNodeValue(v[index], keys[1:])
		return v
	}
	return node
}

// sliceIndex converts the numeric `key` to a slice index.
func (j *Json) sliceIndex(key string) (int, error) {
	index, err := strconv.Atoi(key)
	if err != nil {
		return 0, gerror.WrapCodef(gcode.CodeInvalidParameter, err, `strconv.Atoi failed for string "%s"`, key)
	}
	if index < 0 {
		return 0, gerror.NewCodef(gcode.CodeInvalidParameter, `invalid slice index "%s"`, key)
	}
	return index, nil
}

// convertValue converts `value` to map[string]any or []any,
// which can be supported for hierarchical data access.
func (j *Json) convertValue(value any) (convertedValue any, err error) {
	if value == nil {
		return
	}

	switch value.(type) {
	case map[string]any:
		return value, nil

	case []any:
		return value, nil

	default:
		var (
			reflectInfo = reflection.OriginValueAndKind(value)
		)
		switch reflectInfo.OriginKind {
		case reflect.Array:
			return gconv.Interfaces(value), nil

		case reflect.Slice:
			return gconv.Interfaces(value), nil

		case reflect.Map:
			return gconv.Map(value), nil

		case reflect.Struct:
			if v, ok := value.(iMapStrAny); ok {
				convertedValue = v.MapStrAny()
			}
			if utils.IsNil(convertedValue) {
				if v, ok := value.(iInterfaces); ok {
					convertedValue = v.Interfaces()
				}
			}
			if utils.IsNil(convertedValue) {
				convertedValue = gconv.Map(value)
			}
			if utils.IsNil(convertedValue) {
				err = gerror.NewCodef(gcode.CodeInvalidParameter, `unsupported value type "%s"`, reflect.TypeOf(value))
			}
			return

		default:
			return value, nil
		}
	}
}

// getPointerByPattern returns a pointer to the value by specified `pattern`.
func (j *Json) getPointerByPattern(pattern string) *any {
	if j.p == nil {
		return nil
	}
	if j.vc {
		return j.getPointerByPatternWithViolenceCheck(pattern)
	} else {
		return j.getPointerByPatternWithoutViolenceCheck(pattern)
	}
}

// getPointerByPatternWithViolenceCheck returns a pointer to the value of specified `pattern` with violence check.
func (j *Json) getPointerByPatternWithViolenceCheck(pattern string) *any {
	if !j.vc {
		return j.getPointerByPatternWithoutViolenceCheck(pattern)
	}

	// It returns nil if pattern is empty.
	if pattern == "" {
		return nil
	}
	// It returns all if pattern is ".".
	if pattern == "." {
		return j.p
	}

	var (
		index   = len(pattern)
		start   = 0
		length  = 0
		pointer = j.p
	)
	if index == 0 {
		return pointer
	}
	for {
		if r := j.checkPatternByPointer(pattern[start:index], pointer); r != nil {
			if length += index - start; start > 0 {
				length += 1
			}
			start = index + 1
			index = len(pattern)
			if length == len(pattern) {
				return r
			} else {
				pointer = r
			}
		} else {
			// Get the position for next separator char.
			index = strings.LastIndexByte(pattern[start:index], j.c)
			if index != -1 && length > 0 {
				index += length + 1
			}
		}
		if start >= index {
			break
		}
	}
	return nil
}

// getPointerByPatternWithoutViolenceCheck returns a pointer to the value of specified `pattern`, with no violence check.
func (j *Json) getPointerByPatternWithoutViolenceCheck(pattern string) *any {
	if j.vc {
		return j.getPointerByPatternWithViolenceCheck(pattern)
	}

	// It returns nil if pattern is empty.
	if pattern == "" {
		return nil
	}
	// It returns all if pattern is ".".
	if pattern == "." {
		return j.p
	}

	pointer := j.p
	if len(pattern) == 0 {
		return pointer
	}
	array := strings.Split(pattern, string(j.c))
	for k, v := range array {
		if r := j.checkPatternByPointer(v, pointer); r != nil {
			if k == len(array)-1 {
				return r
			} else {
				pointer = r
			}
		} else {
			break
		}
	}
	return nil
}

// checkPatternByPointer checks whether there's value by `key` in specified `pointer`.
// It returns a pointer to the value.
func (j *Json) checkPatternByPointer(key string, pointer *any) *any {
	switch (*pointer).(type) {
	case map[string]any:
		if v, ok := (*pointer).(map[string]any)[key]; ok {
			return &v
		}
	case []any:
		if gstr.IsNumeric(key) {
			n, err := strconv.Atoi(key)
			if err == nil && len((*pointer).([]any)) > n {
				return &(*pointer).([]any)[n]
			}
		}
	}
	return nil
}
