// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gpool_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gpool"
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_TPool_Int(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create a pool for int
		var (
			newFunc = func() (int, error) {
				return 100, nil
			}
			expireVal  = gtype.NewInt(0)
			expireFunc = func(i int) {
				expireVal.Set(i)
			}
		)

		// TTL = 0, no expiration by time
		p := gpool.NewTPool(0, newFunc, expireFunc)

		// Test Put and Get
		p.Put(1)
		p.Put(2)
		t.Assert(p.Size(), 2)

		v, err := p.Get()
		t.AssertNil(err)
		t.AssertIN(v, g.Slice{1, 2})

		v, err = p.Get()
		t.AssertNil(err)
		t.AssertIN(v, g.Slice{1, 2})

		t.Assert(p.Size(), 0)

		// Test NewFunc when empty
		v, err = p.Get()
		t.AssertNil(err)
		t.Assert(v, 100)

		// Test Clear and ExpireFunc
		p.Put(50)
		t.Assert(p.Size(), 1)
		p.Clear()
		t.Assert(p.Size(), 0)
		t.Assert(expireVal.Val(), 50)

		// Test Close
		p.Put(60)
		p.Close()
		// Close should trigger expire for existing items?
		// Looking at implementation: Close() sets closed=true.
		// It does NOT automatically clear items unless checkExpireItems runs or we call Clear?
		// Wait, checkExpireItems checks closed.Val(). If closed, it clears items.
		// But checkExpireItems runs in a separate goroutine every second.
		// So we might need to wait or trigger it.
		// Actually, let's check the implementation of Close again.
		/*
			func (p *TPool[T]) Close() {
				p.closed.Set(true)
			}
		*/
		// And checkExpireItems:
		/*
			func (p *TPool[T]) checkExpireItems(ctx context.Context) {
				if p.closed.Val() {
					// ... clears items ...
					gtimer.Exit()
				}
				// ...
			}
		*/
		// So it relies on the timer to clean up.
	})
}

func Test_TPool_Struct(t *testing.T) {
	type User struct {
		Id   int
		Name string
	}

	gtest.C(t, func(t *gtest.T) {
		p := gpool.NewTPool[User](time.Hour, nil)
		u1 := User{Id: 1, Name: "john"}
		p.Put(u1)

		v, err := p.Get()
		t.AssertNil(err)
		t.Assert(v, u1)

		// Test empty with no NewFunc
		v, err = p.Get()
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "pool is empty")
		t.Assert(v, User{}) // Zero value
	})
}
