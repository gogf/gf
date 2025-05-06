// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gfpool provides io-reusable pool for file pointer.
package gfpool

import (
	"os"
	"time"

	"github.com/gogf/gf/v3/container/gatomic"
	"github.com/gogf/gf/v3/container/gmap"
	"github.com/gogf/gf/v3/container/gpool"
)

// Pool pointer pool.
type Pool struct {
	id   *gatomic.Int  // Pool id, which is used to mark this pool whether recreated.
	pool *gpool.Pool   // Underlying pool.
	init *gatomic.Bool // Whether initialized, used for marking this file added to fsnotify, and it can only be added just once.
	ttl  time.Duration // Time to live for file pointer items.
}

// File is an item in the pool.
type File struct {
	*os.File             // Underlying file pointer.
	stat     os.FileInfo // State of current file pointer.
	pid      int         // Belonging pool id, which is set when file pointer created. It's used to check whether the pool is recreated.
	pool     *Pool       // Belonging ool.
	flag     int         // Flash for opening file.
	perm     os.FileMode // Permission for opening file.
	path     string      // Absolute path of the file.
}

var (
	// Global file pointer pool.
	pools = gmap.NewStrAnyMap(true)
)
