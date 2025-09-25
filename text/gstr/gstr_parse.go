// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

import (
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// Parse parses the string into map[string]any.
//
// v1=m&v2=n           -> map[v1:m v2:n]
// v[a]=m&v[b]=n       -> map[v:map[a:m b:n]]
// v[a][a]=m&v[a][b]=n -> map[v:map[a:map[a:m b:n]]]
// v[]=m&v[]=n         -> map[v:[m n]]
// v[a][]=m&v[a][]=n   -> map[v:map[a:[m n]]]
// v[][]=m&v[][]=n     -> map[v:[map[]]] // Currently does not support nested slice.
// v=m&v[a]=n          -> error
// a .[[b=c            -> map[a___[b:c]
func Parse(s string) (result map[string]any, err error) {
	if s == "" {
		return nil, nil
	}
	result = make(map[string]any)
	parts := strings.Split(s, "&")
	for _, part := range parts {
		pos := strings.Index(part, "=")
		if pos <= 0 {
			continue
		}
		key, err := url.QueryUnescape(part[:pos])
		if err != nil {
			err = gerror.Wrapf(err, `url.QueryUnescape failed for string "%s"`, part[:pos])
			return nil, err
		}

		for len(key) > 0 && key[0] == ' ' {
			key = key[1:]
		}

		if key == "" || key[0] == '[' {
			continue
		}
		value, err := url.QueryUnescape(part[pos+1:])
		if err != nil {
			err = gerror.Wrapf(err, `url.QueryUnescape failed for string "%s"`, part[pos+1:])
			return nil, err
		}
		// split into multiple keys
		var keys []string
		left := 0
		for i, k := range key {
			if k == '[' && left == 0 {
				left = i
			} else if k == ']' {
				if left > 0 {
					if len(keys) == 0 {
						keys = append(keys, key[:left])
					}
					keys = append(keys, key[left+1:i])
					left = 0
					if i+1 < len(key) && key[i+1] != '[' {
						break
					}
				}
			}
		}
		if len(keys) == 0 {
			keys = append(keys, key)
		}
		// first key
		first := ""
		for i, chr := range keys[0] {
			if chr == ' ' || chr == '.' || chr == '[' {
				first += "_"
			} else {
				first += string(chr)
			}
			if chr == '[' {
				first += keys[0][i+1:]
				break
			}
		}
		keys[0] = first

		// build nested map
		if err = build(result, keys, value); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// build nested map.
func build(result map[string]any, keys []string, value any) error {
	var (
		length = len(keys)
		key    = strings.Trim(keys[0], "'\"")
	)
	if length == 1 {
		result[key] = value
		return nil
	}

	// The end is slice. like f[], f[a][]
	if keys[1] == "" && length == 2 {
		// TODO nested slice
		if key == "" {
			return nil
		}
		val, ok := result[key]
		if !ok {
			result[key] = []any{value}
			return nil
		}
		children, ok := val.([]any)
		if !ok {
			return gerror.NewCodef(
				gcode.CodeInvalidParameter,
				"expected type '[]any' for key '%s', but got '%T'",
				key, val,
			)
		}
		result[key] = append(children, value)
		return nil
	}
	// The end is slice + map. like v[][a]
	if keys[1] == "" && length > 2 && keys[2] != "" {
		val, ok := result[key]
		if !ok {
			result[key] = []any{}
			val = result[key]
		}
		children, ok := val.([]any)
		if !ok {
			return gerror.NewCodef(
				gcode.CodeInvalidParameter,
				"expected type '[]any' for key '%s', but got '%T'",
				key, val,
			)
		}
		if l := len(children); l > 0 {
			if child, ok := children[l-1].(map[string]any); ok {
				if _, ok := child[keys[2]]; !ok {
					_ = build(child, keys[2:], value)
					return nil
				}
			}
		}
		child := map[string]any{}
		_ = build(child, keys[2:], value)
		result[key] = append(children, child)
		return nil
	}

	// map, like v[a], v[a][b]
	val, ok := result[key]
	if !ok {
		result[key] = map[string]any{}
		val = result[key]
	}
	children, ok := val.(map[string]any)
	if !ok {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"expected type 'map[string]any' for key '%s', but got '%T'",
			key, val,
		)
	}
	if err := build(children, keys[1:], value); err != nil {
		return err
	}
	return nil
}
