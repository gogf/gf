// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"bytes"
	"context"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/text/gstr"
	"testing"
)

var (
	ctx = context.TODO()
)

func Test_Print(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Print(ctx, 1, 2, 3)
		l.Println(ctx, 1, 2, 3)
		l.Printf(ctx, "%d %d %d", 1, 2, 3)
		t.Assert(gstr.Count(w.String(), "["), 0)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 3)
	})
}

func Test_Debug(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Debug(ctx, 1, 2, 3)
		l.Debugf(ctx, "%d %d %d", 1, 2, 3)
		t.Assert(gstr.Count(w.String(), defaultLevelPrefixes[LEVEL_DEBU]), 2)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Info(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Info(ctx, 1, 2, 3)
		l.Infof(ctx, "%d %d %d", 1, 2, 3)
		t.Assert(gstr.Count(w.String(), defaultLevelPrefixes[LEVEL_INFO]), 2)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Notice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Notice(ctx, 1, 2, 3)
		l.Noticef(ctx, "%d %d %d", 1, 2, 3)
		t.Assert(gstr.Count(w.String(), defaultLevelPrefixes[LEVEL_NOTI]), 2)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Warning(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Warning(ctx, 1, 2, 3)
		l.Warningf(ctx, "%d %d %d", 1, 2, 3)
		t.Assert(gstr.Count(w.String(), defaultLevelPrefixes[LEVEL_WARN]), 2)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Error(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Error(ctx, 1, 2, 3)
		l.Errorf(ctx, "%d %d %d", 1, 2, 3)
		t.Assert(gstr.Count(w.String(), defaultLevelPrefixes[LEVEL_ERRO]), 2)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 2)
	})
}
