// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package structcache

import (
	"sync"
)

var (
	poolUsedParamsKeyOrTagNameMap = &sync.Pool{
		New: func() any {
			return make(map[string]struct{})
		},
	}
)

// GetUsedParamsKeyOrTagNameMapFromPool retrieves and returns a map for storing params key or tag name.
func GetUsedParamsKeyOrTagNameMapFromPool() map[string]struct{} {
	return poolUsedParamsKeyOrTagNameMap.Get().(map[string]struct{})
}

// PutUsedParamsKeyOrTagNameMapToPool puts a map for storing params key or tag name into pool for re-usage.
func PutUsedParamsKeyOrTagNameMapToPool(m map[string]struct{}) {
	// need to be cleared before putting back into pool.
	for k := range m {
		delete(m, k)
	}
	poolUsedParamsKeyOrTagNameMap.Put(m)
}
