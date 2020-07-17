// Copyright 2017-2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfpool

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// Open creates and returns a file item with given file path, flag and opening permission.
// It automatically creates an associated file pointer pool internally when it's called first time.
// It retrieves a file item from the file pointer pool after then.
func Open(path string, flag int, perm os.FileMode, ttl ...time.Duration) (file *File, err error) {
	var fpTTL time.Duration
	if len(ttl) > 0 {
		fpTTL = ttl[0]
	}
	// DO NOT search the path here wasting performance!
	// Leave following codes just for warning you.
	//
	//path, err = gfile.Search(path)
	//if err != nil {
	//	return nil, err
	//}
	pool := pools.GetOrSetFuncLock(
		fmt.Sprintf("%s&%d&%d&%d", path, flag, fpTTL, perm),
		func() interface{} {
			return New(path, flag, perm, fpTTL)
		},
	).(*Pool)

	return pool.File()
}

// Stat returns the FileInfo structure describing file.
func (f *File) Stat() (os.FileInfo, error) {
	if f.stat == nil {
		return nil, errors.New("file stat is empty")
	}
	return f.stat, nil
}

// Close puts the file pointer back to the file pointer pool.
func (f *File) Close() error {
	if f.pid == f.pool.id.Val() {
		return f.pool.pool.Put(f)
	}
	return nil
}
