// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package guid_test

import (
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/util/guid"
	"sync"
	"testing"

	"github.com/gogf/gf/test/gtest"
)

func Test_S(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		set := gset.NewStrSet()
		for i := 0; i < 1000000; i++ {
			s := guid.S()
			t.Assert(set.AddIfNotExist(s), true)
			t.Assert(len(s), 32)
		}
	})
}

func Test_S_Data(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(len(guid.S([]byte("123"))), 32)
	})
}

func Test_I(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			set = gset.NewSet(true)
			wg  = sync.WaitGroup{}
			ch  = make(chan struct{})
		)
		wg.Add(10)
		for i := 0; i < 10; i++ {
			go func() {
				<-ch
				for i := 0; i < 100000; i++ {
					t.Assert(set.AddIfNotExist(guid.I()), true)
				}
				wg.Done()
			}()
		}
		close(ch)
		wg.Wait()
	})
}
