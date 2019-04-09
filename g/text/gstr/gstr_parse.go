// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

import (
    "fmt"
    "net/url"
    "strings"
)

// Parse parses the string into map[string]interface{}.
//
// f1=m&f2=n           -> map[f1:m f2:n]
// f[a]=m&f[b]=n       -> map[f:map[a:m b:n]]
// f[a][a]=m&f[a][b]=n -> map[f:map[a:map[a:m b:n]]]
// f[]=m&f[]=n         -> map[f:[m n]]
// f[a][]=m&f[a][]=n   -> map[f:map[a:[m n]]]
// f[][]=m&f[][]=n     -> map[f:[map[]]] // Currently does not support nested slice.
// f=m&f[a]=n          -> error
// a .[[b=c            -> map[a___[b:c]
//
func Parse(s string) (result map[string]interface{}, err error) {
    result = make(map[string]interface{})
    parts := strings.Split(s, "&")
    for _, part := range parts {
        pos := strings.Index(part, "=")
        if pos <= 0 {
            continue
        }
        key, err := url.QueryUnescape(part[:pos])
        if err != nil {
            return nil, err
        }
        for key[0] == ' ' {
            key = key[1:]
        }
        if key == "" || key[0] == '[' {
            continue
        }
        value, err := url.QueryUnescape(part[pos+1:])
        if err != nil {
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
                        keys = append(keys, key[ : left])
                    }
                    keys = append(keys, key[left + 1 : i])
                    left = 0
                    if i+1 < len(key) && key[i + 1] != '[' {
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
                first += keys[0][i + 1: ]
                break
            }
        }
        keys[0] = first

        // build nested map
        if err := build(result, keys, value); err != nil {
            return nil, err
        }
    }
    return result, nil
}

// build nested map.
func build(result map[string]interface{}, keys []string, value interface{}) error {
    length := len(keys)
    // trim ',"
    key := strings.Trim(keys[0], "'\"")
    if length == 1 {
        result[key] = value
        return nil
    }

    // The end is slice. like f[], f[a][]
    if keys[1] == "" && length == 2 {
        // todo nested slice
        if key == "" {
            return nil
        }
        val, ok := result[key]
        if !ok {
            result[key] = []interface{}{value}
            return nil
        }
        children, ok := val.([]interface{})
        if !ok {
            return fmt.Errorf("expected type '[]interface{}' for key '%s', but got '%T'", key, val)
        }
        result[key] = append(children, value)
        return nil
    }

    // The end is slice + map. like f[][a]
    if keys[1] == "" && length > 2 && keys[2] != "" {
        val, ok := result[key]
        if !ok {
            result[key] = []interface{}{}
            val = result[key]
        }
        children, ok := val.([]interface{})
        if !ok {
            return fmt.Errorf("expected type '[]interface{}' for key '%s', but got '%T'", key, val)
        }
        if l := len(children); l > 0 {
            if child, ok := children[l-1].(map[string]interface{}); ok {
                if _, ok := child[keys[2]]; !ok {
                    build(child, keys[2:], value)
                    return nil
                }
            }
        }
        child := map[string]interface{}{}
        build(child, keys[2:], value)
        result[key] = append(children, child)
        return nil
    }

    // map. like f[a], f[a][b]
    val, ok := result[key]
    if !ok {
        result[key] = map[string]interface{}{}
        val = result[key]
    }
    children, ok := val.(map[string]interface{})
    if !ok {
        return fmt.Errorf("expected type 'map[string]interface{}' for key '%s', but got '%T'", key, val)
    }
    if err := build(children, keys[1:], value); err != nil {
        return err
    }
    return nil
}