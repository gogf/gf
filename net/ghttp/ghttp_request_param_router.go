// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import "github.com/gogf/gf/container/gvar"

// GetRouterMap retrieves and returns a copy of router map.
func (r *Request) GetRouterMap() map[string]string {
	if r.routerMap != nil {
		m := make(map[string]string, len(r.routerMap))
		for k, v := range r.routerMap {
			m[k] = v
		}
		return m
	}
	return nil
}

// GetRouterValue retrieves and returns the router value with given key name <key>.
// It returns <def> if <key> does not exist.
func (r *Request) GetRouterValue(key string, def ...interface{}) interface{} {
	if r.routerMap != nil {
		if v, ok := r.routerMap[key]; ok {
			return v
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}

// GetRouterVar retrieves and returns the router value as gvar.Var with given key name <key>.
// It returns <def> if <key> does not exist.
func (r *Request) GetRouterVar(key string, def ...interface{}) *gvar.Var {
	return gvar.New(r.GetRouterValue(key, def...))
}

// GetRouterString retrieves and returns the router value as string with given key name <key>.
// It returns <def> if <key> does not exist.
func (r *Request) GetRouterString(key string, def ...interface{}) string {
	return r.GetRouterVar(key, def...).String()
}
