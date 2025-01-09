// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package fieldvar

import "github.com/gogf/gf/v2/util/gconv"

// MapOption specifies the option for map converting.
type MapOption = gconv.MapOption

// Map converts and returns `v` as map[string]interface{}.
func (v *Var) Map(option ...MapOption) map[string]interface{} {
	return gconv.Map(v.Val(), option...)
}

// MapStrAny is like function Map, but implements the interface of MapStrAny.
func (v *Var) MapStrAny(option ...MapOption) map[string]interface{} {
	return v.Map(option...)
}

// MapStrStr converts and returns `v` as map[string]string.
func (v *Var) MapStrStr(option ...MapOption) map[string]string {
	return gconv.MapStrStr(v.Val(), option...)
}

// MapStrVar converts and returns `v` as map[string]Var.
func (v *Var) MapStrVar(option ...MapOption) map[string]*Var {
	m := v.Map(option...)
	if len(m) > 0 {
		vMap := make(map[string]*Var, len(m))
		for k, v := range m {
			vMap[k] = New(v)
		}
		return vMap
	}
	return nil
}

// Maps converts and returns `v` as map[string]string.
// See gconv.Maps.
func (v *Var) Maps(option ...MapOption) []map[string]interface{} {
	return gconv.Maps(v.Val(), option...)
}

// MapToMap converts any map type variable `params` to another map type variable `pointer`.
// See gconv.MapToMap.
func (v *Var) MapToMap(pointer interface{}, mapping ...map[string]string) (err error) {
	return gconv.MapToMap(v.Val(), pointer, mapping...)
}

// MapToMaps converts any map type variable `params` to another map type variable `pointer`.
// See gconv.MapToMaps.
func (v *Var) MapToMaps(pointer interface{}, mapping ...map[string]string) (err error) {
	return gconv.MapToMaps(v.Val(), pointer, mapping...)
}

// MapToMapsDeep converts any map type variable `params` to another map type variable
// `pointer` recursively.
// See gconv.MapToMapsDeep.
func (v *Var) MapToMapsDeep(pointer interface{}, mapping ...map[string]string) (err error) {
	return gconv.MapToMaps(v.Val(), pointer, mapping...)
}
