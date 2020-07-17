// Copyright 2018 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package gcache

import (
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/jin502437344/gf/os/gtimer"
)

// Cache struct.
type Cache struct {
	*memCache
}

// New creates and returns a new cache object.
func New(lruCap ...int) *Cache {
	c := &Cache{
		memCache: newMemCache(lruCap...),
	}
	gtimer.AddSingleton(time.Second, c.syncEventAndClearExpired)
	return c
}

// Clear clears all data of the cache.
func (c *Cache) Clear() {
	// atomic swap to ensure atomicity.
	old := atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&c.memCache)), unsafe.Pointer(newMemCache()))
	// close the old cache object.
	(*memCache)(old).Close()
}
