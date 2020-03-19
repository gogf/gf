// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gpool_test

import (
	"errors"
	"testing"
	"time"

	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf/container/gpool"
	"github.com/gogf/gf/test/gtest"
)

var nf gpool.NewFunc = func() (i interface{}, e error) {
	return "hello", nil
}

var assertIndex int = 0
var ef gpool.ExpireFunc = func(i interface{}) {
	assertIndex++
	t.Assert(i, assertIndex)
}

func Test_Gpool(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		//
		//expire = 0
		p1 := gpool.New(0, nf)
		p1.Put(1)
		p1.Put(2)
		time.Sleep(1 * time.Second)
		//test won't be timeout
		v1, err1 := p1.Get()
		t.Assert(err1, nil)
		t.AssertIN(v1, g.Slice{1, 2})
		//test clear
		p1.Clear()
		t.Assert(p1.Size(), 0)
		//test newFunc
		v1, err1 = p1.Get()
		t.Assert(err1, nil)
		t.Assert(v1, "hello")
		//put data again
		p1.Put(3)
		p1.Put(4)
		v1, err1 = p1.Get()
		t.Assert(err1, nil)
		t.AssertIN(v1, g.Slice{3, 4})
		//test close
		p1.Close()
		v1, err1 = p1.Get()
		t.Assert(err1, nil)
		t.Assert(v1, "hello")
	})

	gtest.C(t, func(t *gtest.T) {
		//
		//expire > 0
		p2 := gpool.New(2*time.Second, nil, ef)
		for index := 0; index < 10; index++ {
			p2.Put(index)
		}
		t.Assert(p2.Size(), 10)
		v2, err2 := p2.Get()
		t.Assert(err2, nil)
		t.Assert(v2, 0)
		//test timeout expireFunc
		time.Sleep(3 * time.Second)
		v2, err2 = p2.Get()
		t.Assert(err2, errors.New("pool is empty"))
		t.Assert(v2, nil)
		//test close expireFunc
		for index := 0; index < 10; index++ {
			p2.Put(index)
		}
		t.Assert(p2.Size(), 10)
		v2, err2 = p2.Get()
		t.Assert(err2, nil)
		t.Assert(v2, 0)
		assertIndex = 0
		p2.Close()
		time.Sleep(3 * time.Second)
	})

	gtest.C(t, func(t *gtest.T) {
		//
		//expire < 0
		p3 := gpool.New(-1, nil)
		v3, err3 := p3.Get()
		t.Assert(err3, errors.New("pool is empty"))
		t.Assert(v3, nil)
	})
}
