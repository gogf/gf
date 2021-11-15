// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog_test

import (
	"sync"
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_Concurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := 1000
		l := glog.New()
		s := "@1234567890#"
		f := "test.log"
		p := gfile.TempDir(gtime.TimestampNanoStr())
		t.Assert(l.SetPath(p), nil)
		defer gfile.Remove(p)
		wg := sync.WaitGroup{}
		ch := make(chan struct{})
		for i := 0; i < c; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-ch
				l.File(f).Stdout(false).Print(ctx, s)
			}()
		}
		close(ch)
		wg.Wait()
		content := gfile.GetContents(gfile.Join(p, f))
		t.Assert(gstr.Count(content, s), c)
	})
}
