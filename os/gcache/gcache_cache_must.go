// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
)

// MustGet acts like Get, but it panics if any error occurs.
func (c *Cache) MustGet(ctx context.Context, key interface{}) *gvar.Var {
	v, err := c.Get(ctx, key)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetOrSet acts like GetOrSet, but it panics if any error occurs.
func (c *Cache) MustGetOrSet(ctx context.Context, key interface{}, value interface{}, duration time.Duration) *gvar.Var {
	v, err := c.GetOrSet(ctx, key, value, duration)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetOrSetFunc acts like GetOrSetFunc, but it panics if any error occurs.
func (c *Cache) MustGetOrSetFunc(ctx context.Context, key interface{}, f Func, duration time.Duration) *gvar.Var {
	v, err := c.GetOrSetFunc(ctx, key, f, duration)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetOrSetFuncLock acts like GetOrSetFuncLock, but it panics if any error occurs.
func (c *Cache) MustGetOrSetFuncLock(ctx context.Context, key interface{}, f Func, duration time.Duration) *gvar.Var {
	v, err := c.GetOrSetFuncLock(ctx, key, f, duration)
	if err != nil {
		panic(err)
	}
	return v
}

// MustContains acts like Contains, but it panics if any error occurs.
func (c *Cache) MustContains(ctx context.Context, key interface{}) bool {
	v, err := c.Contains(ctx, key)
	if err != nil {
		panic(err)
	}
	return v
}

// MustGetExpire acts like GetExpire, but it panics if any error occurs.
func (c *Cache) MustGetExpire(ctx context.Context, key interface{}) time.Duration {
	v, err := c.GetExpire(ctx, key)
	if err != nil {
		panic(err)
	}
	return v
}

// MustSize acts like Size, but it panics if any error occurs.
func (c *Cache) MustSize(ctx context.Context) int {
	v, err := c.Size(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// MustData acts like Data, but it panics if any error occurs.
func (c *Cache) MustData(ctx context.Context) map[interface{}]interface{} {
	v, err := c.Data(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// MustKeys acts like Keys, but it panics if any error occurs.
func (c *Cache) MustKeys(ctx context.Context) []interface{} {
	v, err := c.Keys(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// MustKeyStrings acts like KeyStrings, but it panics if any error occurs.
func (c *Cache) MustKeyStrings(ctx context.Context) []string {
	v, err := c.KeyStrings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// MustValues acts like Values, but it panics if any error occurs.
func (c *Cache) MustValues(ctx context.Context) []interface{} {
	v, err := c.Values(ctx)
	if err != nil {
		panic(err)
	}
	return v
}
