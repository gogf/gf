// Copyright 2017-2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfpool

import (
	"os"
	"time"

	"github.com/gogf/gf/container/gpool"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/os/gfsnotify"
)

// New creates and returns a file pointer pool with given file path, flag and opening permission.
//
// Note the expiration logic:
// ttl = 0 : not expired;
// ttl < 0 : immediate expired after use;
// ttl > 0 : timeout expired;
// It is not expired in default.
func New(path string, flag int, perm os.FileMode, ttl ...time.Duration) *Pool {
	var fpTTL time.Duration
	if len(ttl) > 0 {
		fpTTL = ttl[0]
	}
	p := &Pool{
		id:   gtype.NewInt(),
		ttl:  fpTTL,
		init: gtype.NewBool(),
	}
	p.pool = newFilePool(p, path, flag, perm, fpTTL)
	return p
}

// newFilePool creates and returns a file pointer pool with given file path, flag and opening permission.
func newFilePool(p *Pool, path string, flag int, perm os.FileMode, ttl time.Duration) *gpool.Pool {
	pool := gpool.New(ttl, func() (interface{}, error) {
		file, err := os.OpenFile(path, flag, perm)
		if err != nil {
			return nil, err
		}
		return &File{
			File: file,
			pid:  p.id.Val(),
			pool: p,
			flag: flag,
			perm: perm,
			path: path,
		}, nil
	}, func(i interface{}) {
		_ = i.(*File).File.Close()
	})
	return pool
}

// File retrieves file item from the file pointer pool and returns it. It creates one if
// the file pointer pool is empty.
// Note that it should be closed when it will never be used. When it's closed, it is not
// really closed the underlying file pointer but put back to the file pinter pool.
func (p *Pool) File() (*File, error) {
	if v, err := p.pool.Get(); err != nil {
		return nil, err
	} else {
		var err error
		f := v.(*File)
		f.stat, err = os.Stat(f.path)
		if f.flag&os.O_CREATE > 0 {
			if os.IsNotExist(err) {
				if f.File, err = os.OpenFile(f.path, f.flag, f.perm); err != nil {
					return nil, err
				} else {
					// Retrieve the state of the new created file.
					if f.stat, err = f.File.Stat(); err != nil {
						return nil, err
					}
				}
			}
		}
		if f.flag&os.O_TRUNC > 0 {
			if f.stat.Size() > 0 {
				if err = f.Truncate(0); err != nil {
					return nil, err
				}
			}
		}
		if f.flag&os.O_APPEND > 0 {
			if _, err = f.Seek(0, 2); err != nil {
				return nil, err
			}
		} else {
			if _, err = f.Seek(0, 0); err != nil {
				return nil, err
			}
		}
		// It firstly checks using !p.init.Val() for performance purpose.
		if !p.init.Val() && p.init.Cas(false, true) {
			_, _ = gfsnotify.Add(f.path, func(event *gfsnotify.Event) {
				// If teh file is removed or renamed, recreates the pool by increasing the pool id.
				if event.IsRemove() || event.IsRename() {
					// It drops the old pool.
					p.id.Add(1)
					// Clears the pool items staying in the pool.
					p.pool.Clear()
					// It uses another adding to drop the file items between the two adding.
					// Whenever the pool id changes, the pool will be recreated.
					p.id.Add(1)
				}
			}, false)
		}
		return f, nil
	}
}

// Close closes current file pointer pool.
func (p *Pool) Close() {
	p.pool.Close()
}
