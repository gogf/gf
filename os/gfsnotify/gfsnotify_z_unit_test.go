// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfsnotify_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gfsnotify"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func TestWatcher_AddOnce(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		value := gtype.New()
		path := gfile.Temp(gconv.String(gtime.TimestampNano()))
		err := gfile.PutContents(path, "init")
		t.AssertNil(err)
		defer gfile.Remove(path)

		time.Sleep(100 * time.Millisecond)
		callback1, err := gfsnotify.AddOnce("mywatch", path, func(event *gfsnotify.Event) {
			value.Set(1)
		})
		t.AssertNil(err)
		callback2, err := gfsnotify.AddOnce("mywatch", path, func(event *gfsnotify.Event) {
			value.Set(2)
		})
		t.AssertNil(err)
		t.Assert(callback2, nil)

		err = gfile.PutContents(path, "1")
		t.AssertNil(err)

		time.Sleep(100 * time.Millisecond)
		t.Assert(value, 1)

		err = gfsnotify.RemoveCallback(callback1.Id)
		t.AssertNil(err)

		err = gfile.PutContents(path, "2")
		t.AssertNil(err)

		time.Sleep(100 * time.Millisecond)
		t.Assert(value, 1)
	})
}

func TestWatcher_AddRemove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		path1 := gfile.Temp() + gfile.Separator + gconv.String(gtime.TimestampNano())
		path2 := gfile.Temp() + gfile.Separator + gconv.String(gtime.TimestampNano()) + "2"
		gfile.PutContents(path1, "1")
		defer func() {
			gfile.Remove(path1)
			gfile.Remove(path2)
		}()
		v := gtype.NewInt(1)
		callback, err := gfsnotify.Add(path1, func(event *gfsnotify.Event) {
			if event.IsWrite() {
				v.Set(2)
				return
			}
			if event.IsRename() {
				v.Set(3)
				gfsnotify.Exit()
				return
			}
		})
		t.AssertNil(err)
		t.AssertNE(callback, nil)

		gfile.PutContents(path1, "2")
		time.Sleep(100 * time.Millisecond)
		t.Assert(v.Val(), 2)

		gfile.Rename(path1, path2)
		time.Sleep(100 * time.Millisecond)
		t.Assert(v.Val(), 3)
	})

	gtest.C(t, func(t *gtest.T) {
		path1 := gfile.Temp() + gfile.Separator + gconv.String(gtime.TimestampNano())
		gfile.PutContents(path1, "1")
		defer func() {
			gfile.Remove(path1)
		}()
		v := gtype.NewInt(1)
		callback, err := gfsnotify.Add(path1, func(event *gfsnotify.Event) {
			if event.IsWrite() {
				v.Set(2)
				return
			}
			if event.IsRemove() {
				v.Set(4)
				return
			}
		})
		t.AssertNil(err)
		t.AssertNE(callback, nil)

		gfile.PutContents(path1, "2")
		time.Sleep(100 * time.Millisecond)
		t.Assert(v.Val(), 2)

		gfile.Remove(path1)
		time.Sleep(100 * time.Millisecond)
		t.Assert(v.Val(), 4)

		gfile.PutContents(path1, "1")
		time.Sleep(100 * time.Millisecond)
		t.Assert(v.Val(), 4)
	})
}

func TestWatcher_Callback1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		path1 := gfile.Temp(gtime.TimestampNanoStr())
		gfile.PutContents(path1, "1")
		defer func() {
			gfile.Remove(path1)
		}()
		v := gtype.NewInt(1)
		callback, err := gfsnotify.Add(path1, func(event *gfsnotify.Event) {
			if event.IsWrite() {
				v.Set(2)
				return
			}
		})
		t.AssertNil(err)
		t.AssertNE(callback, nil)

		gfile.PutContents(path1, "2")
		time.Sleep(100 * time.Millisecond)
		t.Assert(v.Val(), 2)

		v.Set(3)
		gfsnotify.RemoveCallback(callback.Id)
		gfile.PutContents(path1, "3")
		time.Sleep(100 * time.Millisecond)
		t.Assert(v.Val(), 3)
	})
}

func TestWatcher_Callback2(t *testing.T) {
	// multiple callbacks
	gtest.C(t, func(t *gtest.T) {
		path1 := gfile.Temp(gtime.TimestampNanoStr())
		t.Assert(gfile.PutContents(path1, "1"), nil)
		defer func() {
			gfile.Remove(path1)
		}()
		v1 := gtype.NewInt(1)
		v2 := gtype.NewInt(1)
		callback1, err1 := gfsnotify.Add(path1, func(event *gfsnotify.Event) {
			if event.IsWrite() {
				v1.Set(2)
				return
			}
		})
		callback2, err2 := gfsnotify.Add(path1, func(event *gfsnotify.Event) {
			if event.IsWrite() {
				v2.Set(2)
				return
			}
		})
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.AssertNE(callback1, nil)
		t.AssertNE(callback2, nil)

		t.Assert(gfile.PutContents(path1, "2"), nil)
		time.Sleep(100 * time.Millisecond)
		t.Assert(v1.Val(), 2)
		t.Assert(v2.Val(), 2)

		v1.Set(3)
		v2.Set(3)
		gfsnotify.RemoveCallback(callback1.Id)
		t.Assert(gfile.PutContents(path1, "3"), nil)
		time.Sleep(100 * time.Millisecond)
		t.Assert(v1.Val(), 3)
		t.Assert(v2.Val(), 2)
	})
}

func TestWatcher_WatchFolderWithoutRecursively(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err     error
			array   = garray.New(true)
			dirPath = gfile.Temp(gtime.TimestampNanoStr())
		)
		err = gfile.Mkdir(dirPath)
		t.AssertNil(err)

		_, err = gfsnotify.Add(dirPath, func(event *gfsnotify.Event) {
			// fmt.Println(event.String())
			array.Append(1)
		}, false)
		t.AssertNil(err)
		time.Sleep(time.Millisecond * 100)
		t.Assert(array.Len(), 0)

		f, err := gfile.Create(gfile.Join(dirPath, "1"))
		t.AssertNil(err)
		t.AssertNil(f.Close())
		time.Sleep(time.Millisecond * 100)
		t.Assert(array.Len(), 1)
	})
}
