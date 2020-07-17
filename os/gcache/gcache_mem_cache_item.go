// Copyright 2018 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package gcache

import (
	"github.com/jin502437344/gf/os/gtime"
)

// IsExpired checks whether <item> is expired.
func (item *memCacheItem) IsExpired() bool {
	// Note that it should use greater than or equal judgement here
	// imagining that the cache time is only 1 millisecond.
	if item.e >= gtime.TimestampMilli() {
		return false
	}
	return true
}
