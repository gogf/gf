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
// 2. It's quite complicated in hierarchical data search, node creating and data assignment;
func (j *Json) setValue(pattern string, value any, removed bool) error {
	var (
		err   error
		array = strings.Split(pattern, string(j.c))
	)
	if value, err = j.convertValue(value); err != nil {
		return err
	}
	args := &setPathArgs{
		pattern:  pattern,
		value:    value,
		removing: removed && value == nil,
	}
	j.mu.Lock()
	defer j.mu.Unlock()
	newRoot, err := j.setPath(*j.p, array, args)
	if err != nil {
		return err
	}
	*j.p = newRoot
	return nil
}

// setPathArgs holds the state shared by one Set or Remove operation.
type setPathArgs struct {
	pattern  string
	value    any
	removing bool
}

// setPath resolves the remaining pattern segments and returns the new node value.
// It updates the current node only after the child node is resolved, so an error
// leaves the data unchanged.
func (j *Json) setPath(node any, array []string, args *setPathArgs) (any, error) {
	var (
		segment = array[0]
		rest    = array[1:]
	)
	switch v := node.(type) {
	case map[string]any:
		return j.setMapPath(v, segment, rest, args)
	case []any:
		return j.setArrayPath(v, segment, rest, args)
	default:
		return j.setDefaultPath(node, segment, rest, args)
	}
}

// setMapPath resolves one pattern segment against a map node.
func (j *Json) setMapPath(node map[string]any, segment string, rest []string, args *setPathArgs) (any, error) {
	if len(rest) == 0 {
		if args.removing {
			delete(node, segment)
			return node, nil
		}
		if node == nil {
			node = make(map[string]any)
		}
		node[segment] = args.value
		return node, nil
	}
	child, ok := node[segment]
	if !ok {
		if args.removing {
			return node, nil
		}
		var err error
		if child, err = newSetNode(rest[0], args.pattern); err != nil {
			return nil, err
		}
	}
	newChild, err := j.setPath(child, rest, args)
	if err != nil {
		return nil, err
	}
	if node == nil {
		node = make(map[string]any)
	}
	node[segment] = newChild
	return node, nil
}

// setArrayPath resolves one pattern segment against an array node.
func (j *Json) setArrayPath(node []any, segment string, rest []string, args *setPathArgs) (any, error) {
	if !gstr.IsNumeric(segment) {
		if args.removing {
			return node, nil
		}
		if len(rest) == 0 {
			return map[string]any{segment: args.value}, nil
		}
		return j.setPath(map[string]any{}, rest, args)
	}
	n, err := parseArrayIndex(segment, args.pattern)
	if err != nil {
		return nil, err
	}
	if len(rest) == 0 {
		if args.removing {
			if n < len(node) {
				return append(node[:n], node[n+1:]...), nil
			}
			return node, nil
		}
		if n < len(node) {
			node[n] = args.value
			return node, nil
		}
		extended := make([]any, n+1)
		copy(extended, node)
		extended[n] = args.value
		return extended, nil
	}
	var child any
	if n < len(node) {
		child = node[n]
	} else {
		if args.removing {
			return node, nil
		}
		if child, err = newSetNode(rest[0], args.pattern); err != nil {
			return nil, err
		}
	}
	newChild, err := j.setPath(child, rest, args)
	if err != nil {
		return nil, err
	}
	if n < len(node) {
		node[n] = newChild
		return node, nil
	}
	extended := make([]any, n+1)
	copy(extended, node)
	extended[n] = newChild
	return extended, nil
}

// setDefaultPath replaces a scalar node with the hierarchy built from the remaining segments.
func (j *Json) setDefaultPath(node any, segment string, rest []string, args *setPathArgs) (any, error) {
	if args.removing {
		return node, nil
	}
	if !gstr.IsNumeric(segment) {
		if len(rest) == 0 {
			return map[string]any{segment: args.value}, nil
		}
		child, err := j.setPath(nil, rest, args)
		if err != nil {
			return nil, err
		}
		return map[string]any{segment: child}, nil
	}
	n, err := parseArrayIndex(segment, args.pattern)
	if err != nil {
		return nil, err
	}
	nodeArray := make([]any, n+1)
	if len(rest) == 0 {
		nodeArray[n] = args.value
		return nodeArray, nil
	}
	child, err := j.setPath(nil, rest, args)
	if err != nil {
		return nil, err
	}
	nodeArray[n] = child
	return nodeArray, nil
}

// parseArrayIndex parses a non-negative array index from a pattern segment.
func parseArrayIndex(segment, pattern string) (int, error) {
	n, err := strconv.Atoi(segment)
	if err != nil {
		return 0, gerror.WrapCodef(gcode.CodeInvalidParameter, err, `strconv.Atoi failed for string "%s"`, segment)
	}
	if n < 0 {
		return 0, gerror.NewCodef(gcode.CodeInvalidParameter, `invalid array index "%s" in pattern "%s"`, segment, pattern)
	}
	return n, nil
}

// newSetNode creates an empty node matching the type implied by the next segment.
func newSetNode(segment, pattern string) (any, error) {
	if !gstr.IsNumeric(segment) {
		return map[string]any{}, nil
	}
	n, err := parseArrayIndex(segment, pattern)
	if err != nil {
		return nil, err
	}
	return make([]any, n+1), nil
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
			if err == nil && n >= 0 && len((*pointer).([]any)) > n {
				return &(*pointer).([]any)[n]
			}
		}
	}
	return nil
}
