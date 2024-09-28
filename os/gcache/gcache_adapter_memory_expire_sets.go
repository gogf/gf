// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache

import (
	"sync"

	"github.com/gogf/gf/v2/container/gset"
)

type memoryExpireSets struct {
	// expireSetMu ensures the concurrent safety of expireSets map.
	mu sync.RWMutex
	// expireSets is the expiring timestamp in seconds to its key set mapping, which is used for quick indexing and deleting.
	expireSets map[int64]*gset.Set
}

func newMemoryExpireSets() *memoryExpireSets {
	return &memoryExpireSets{
		expireSets: make(map[int64]*gset.Set),
	}
}

func (d *memoryExpireSets) Get(key int64) (result *gset.Set) {
	d.mu.RLock()
	result = d.expireSets[key]
	d.mu.RUnlock()
	return
}

func (d *memoryExpireSets) GetOrNew(key int64) (result *gset.Set) {
	if result = d.Get(key); result != nil {
		return
	}
	d.mu.Lock()
	if es, ok := d.expireSets[key]; ok {
		result = es
	} else {
		result = gset.New(true)
		d.expireSets[key] = result
	}
	d.mu.Unlock()
	return
}

func (d *memoryExpireSets) Delete(key int64) {
	d.mu.Lock()
	delete(d.expireSets, key)
	d.mu.Unlock()
}
