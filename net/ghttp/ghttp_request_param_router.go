// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
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

// GetRouter retrieves and returns the router value with given key name <key>.
// It returns <def> if <key> does not exist.
func (r *Request) GetRouter(key string, def ...interface{}) *gvar.Var {
	if r.routerMap != nil {
		if v, ok := r.routerMap[key]; ok {
			return gvar.New(v)
		}
	}
	if len(def) > 0 {
		return gvar.New(def[0])
	}
	return nil
}
