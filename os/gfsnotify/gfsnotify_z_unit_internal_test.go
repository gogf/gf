// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfsnotify

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Watcher_RemoveCallback(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		file, err := os.CreateTemp("", "gfsnotify-remove-callback-*.txt")
		t.AssertNil(err)
		path := file.Name()
		t.AssertNil(file.Close())

		w, err := New()
		t.AssertNil(err)
		defer func() {
			w.Close()
			t.AssertNil(os.Remove(path))
		}()

		callback1, err := w.Add(path, func(event *Event) {})
		t.AssertNil(err)
		t.Assert(len(w.watcher.WatchList()), 1)

		callback2, err := w.Add(path, func(event *Event) {})
		t.AssertNil(err)
		t.Assert(len(w.watcher.WatchList()), 1)

		w.RemoveCallback(callback1.Id)
		t.Assert(len(w.watcher.WatchList()), 1)

		w.RemoveCallback(callback2.Id)
		t.Assert(len(w.watcher.WatchList()), 0)

		callback3, err := w.Add(path, func(event *Event) {})
		t.AssertNil(err)
		t.Assert(len(w.watcher.WatchList()), 1)

		w.RemoveCallback(callback3.Id)
		t.Assert(len(w.watcher.WatchList()), 0)
	})
}

func Test_Watcher_RemoveCallbackRecursiveParent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		rootPath, err := os.MkdirTemp("", "gfsnotify-recursive-*")
		t.AssertNil(err)
		subPath := filepath.Join(rootPath, "sub")
		t.AssertNil(os.Mkdir(subPath, 0755))

		w, err := New()
		t.AssertNil(err)
		defer func() {
			w.Close()
			t.AssertNil(os.RemoveAll(rootPath))
		}()

		eventCh := make(chan string, 10)
		_, err = w.Add(rootPath, func(event *Event) {
			eventCh <- event.Path
		})
		t.AssertNil(err)

		subCallback, err := w.Add(subPath, func(event *Event) {})
		t.AssertNil(err)
		w.RemoveCallback(subCallback.Id)
		t.Assert(len(w.watcher.WatchList()), 2)

		t.AssertNil(os.WriteFile(filepath.Join(subPath, "test.txt"), []byte("test"), 0644))
		select {
		case <-eventCh:
		case <-time.After(time.Second):
			t.Fatal("recursive callback did not receive the subdirectory event")
		}
	})
}

func Test_Watcher_RemoveCallbackRecursiveRoot(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		rootPath, err := os.MkdirTemp("", "gfsnotify-recursive-root-*")
		t.AssertNil(err)
		nestedPath := filepath.Join(rootPath, "sub", "nested")
		t.AssertNil(os.MkdirAll(nestedPath, 0755))

		w, err := New()
		t.AssertNil(err)
		defer func() {
			w.Close()
			t.AssertNil(os.RemoveAll(rootPath))
		}()

		callback, err := w.Add(rootPath, func(event *Event) {})
		t.AssertNil(err)
		t.Assert(len(w.watcher.WatchList()), 3)

		w.RemoveCallback(callback.Id)
		t.Assert(len(w.watcher.WatchList()), 0)
	})
}

func Test_Watcher_RemoveCallbackKeepsDescendantCallback(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		rootPath, err := os.MkdirTemp("", "gfsnotify-recursive-keep-*")
		t.AssertNil(err)
		subPath := filepath.Join(rootPath, "sub")
		nestedPath := filepath.Join(subPath, "nested")
		t.AssertNil(os.MkdirAll(nestedPath, 0755))

		w, err := New()
		t.AssertNil(err)
		defer func() {
			w.Close()
			t.AssertNil(os.RemoveAll(rootPath))
		}()

		rootCallback, err := w.Add(rootPath, func(event *Event) {})
		t.AssertNil(err)
		eventCh := make(chan string, 10)
		_, err = w.Add(subPath, func(event *Event) {
			eventCh <- event.Path
		})
		t.AssertNil(err)
		t.Assert(len(w.watcher.WatchList()), 3)

		w.RemoveCallback(rootCallback.Id)
		watched := w.watcher.WatchList()
		t.Assert(len(watched), 2)
		t.AssertIN(subPath, watched)
		t.AssertIN(nestedPath, watched)
		t.AssertNI(rootPath, watched)

		t.AssertNil(os.WriteFile(filepath.Join(nestedPath, "test.txt"), []byte("test"), 0644))
		select {
		case <-eventCh:
		case <-time.After(time.Second):
			t.Fatal("descendant callback did not receive the nested directory event")
		}
	})
}
