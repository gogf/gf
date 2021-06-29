// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"bytes"
	"fmt"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/text/gstr"
	"testing"
	"time"
)

func Test_To(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		To(w).Error(1, 2, 3)
		To(w).Errorf("%d %d %d", 1, 2, 3)
		t.Assert(gstr.Count(w.String(), defaultLevelPrefixes[LEVEL_ERRO]), 2)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Path(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		path := gfile.TempDir(gtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gtime.TimestampNano())

		err := gfile.Mkdir(path)
		t.Assert(err, nil)
		defer gfile.Remove(path)

		Path(path).File(file).Stdout(false).Error(1, 2, 3)
		Path(path).File(file).Stdout(false).Errorf("%d %d %d", 1, 2, 3)
		content := gfile.GetContents(gfile.Join(path, file))
		t.Assert(gstr.Count(content, defaultLevelPrefixes[LEVEL_ERRO]), 2)
		t.Assert(gstr.Count(content, "1 2 3"), 2)
	})
}

func Test_Cat(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cat := "category"
		path := gfile.TempDir(gtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gtime.TimestampNano())

		err := gfile.Mkdir(path)
		t.Assert(err, nil)
		defer gfile.Remove(path)

		Path(path).File(file).Cat(cat).Stdout(false).Error(1, 2, 3)
		Path(path).File(file).Cat(cat).Stdout(false).Errorf("%d %d %d", 1, 2, 3)
		content := gfile.GetContents(gfile.Join(path, cat, file))
		t.Assert(gstr.Count(content, defaultLevelPrefixes[LEVEL_ERRO]), 2)
		t.Assert(gstr.Count(content, "1 2 3"), 2)
	})
}

func Test_Level(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		path := gfile.TempDir(gtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gtime.TimestampNano())

		err := gfile.Mkdir(path)
		t.Assert(err, nil)
		defer gfile.Remove(path)

		Path(path).File(file).Level(LEVEL_PROD).Stdout(false).Debug(1, 2, 3)
		Path(path).File(file).Level(LEVEL_PROD).Stdout(false).Debug("%d %d %d", 1, 2, 3)
		content := gfile.GetContents(gfile.Join(path, file))
		t.Assert(gstr.Count(content, defaultLevelPrefixes[LEVEL_DEBU]), 0)
		t.Assert(gstr.Count(content, "1 2 3"), 0)
	})
}

func Test_Skip(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		path := gfile.TempDir(gtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gtime.TimestampNano())

		err := gfile.Mkdir(path)
		t.Assert(err, nil)
		defer gfile.Remove(path)

		Path(path).File(file).Skip(10).Stdout(false).Error(1, 2, 3)
		Path(path).File(file).Stdout(false).Errorf("%d %d %d", 1, 2, 3)
		content := gfile.GetContents(gfile.Join(path, file))
		t.Assert(gstr.Count(content, defaultLevelPrefixes[LEVEL_ERRO]), 2)
		t.Assert(gstr.Count(content, "1 2 3"), 2)
		t.Assert(gstr.Count(content, "Stack"), 1)
	})
}

func Test_Stack(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		path := gfile.TempDir(gtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gtime.TimestampNano())

		err := gfile.Mkdir(path)
		t.Assert(err, nil)
		defer gfile.Remove(path)

		Path(path).File(file).Stack(false).Stdout(false).Error(1, 2, 3)
		Path(path).File(file).Stdout(false).Errorf("%d %d %d", 1, 2, 3)
		content := gfile.GetContents(gfile.Join(path, file))
		t.Assert(gstr.Count(content, defaultLevelPrefixes[LEVEL_ERRO]), 2)
		t.Assert(gstr.Count(content, "1 2 3"), 2)
		t.Assert(gstr.Count(content, "Stack"), 1)
	})
}

func Test_StackWithFilter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		path := gfile.TempDir(gtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gtime.TimestampNano())

		err := gfile.Mkdir(path)
		t.Assert(err, nil)
		defer gfile.Remove(path)

		Path(path).File(file).StackWithFilter("none").Stdout(false).Error(1, 2, 3)
		content := gfile.GetContents(gfile.Join(path, file))
		t.Assert(gstr.Count(content, defaultLevelPrefixes[LEVEL_ERRO]), 1)
		t.Assert(gstr.Count(content, "1 2 3"), 1)
		t.Assert(gstr.Count(content, "Stack"), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		path := gfile.TempDir(gtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gtime.TimestampNano())

		err := gfile.Mkdir(path)
		t.Assert(err, nil)
		defer gfile.Remove(path)

		Path(path).File(file).StackWithFilter("gogf").Stdout(false).Error(1, 2, 3)
		content := gfile.GetContents(gfile.Join(path, file))
		t.Assert(gstr.Count(content, defaultLevelPrefixes[LEVEL_ERRO]), 1)
		t.Assert(gstr.Count(content, "1 2 3"), 1)
		t.Assert(gstr.Count(content, "Stack"), 0)
	})
}

func Test_Header(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		path := gfile.TempDir(gtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gtime.TimestampNano())

		err := gfile.Mkdir(path)
		t.Assert(err, nil)
		defer gfile.Remove(path)

		Path(path).File(file).Header(true).Stdout(false).Error(1, 2, 3)
		content := gfile.GetContents(gfile.Join(path, file))
		t.Assert(gstr.Count(content, defaultLevelPrefixes[LEVEL_ERRO]), 1)
		t.Assert(gstr.Count(content, "1 2 3"), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		path := gfile.TempDir(gtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gtime.TimestampNano())

		err := gfile.Mkdir(path)
		t.Assert(err, nil)
		defer gfile.Remove(path)

		Path(path).File(file).Header(false).Stdout(false).Error(1, 2, 3)
		content := gfile.GetContents(gfile.Join(path, file))
		t.Assert(gstr.Count(content, defaultLevelPrefixes[LEVEL_ERRO]), 0)
		t.Assert(gstr.Count(content, "1 2 3"), 1)
	})
}

func Test_Line(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		path := gfile.TempDir(gtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gtime.TimestampNano())

		err := gfile.Mkdir(path)
		t.Assert(err, nil)
		defer gfile.Remove(path)

		Path(path).File(file).Line(true).Stdout(false).Debug(1, 2, 3)
		content := gfile.GetContents(gfile.Join(path, file))
		t.Assert(gstr.Count(content, defaultLevelPrefixes[LEVEL_DEBU]), 1)
		t.Assert(gstr.Count(content, "1 2 3"), 1)
		t.Assert(gstr.Count(content, ".go"), 1)
		t.Assert(gstr.Contains(content, gfile.Separator), true)
	})
	gtest.C(t, func(t *gtest.T) {
		path := gfile.TempDir(gtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gtime.TimestampNano())

		err := gfile.Mkdir(path)
		t.Assert(err, nil)
		defer gfile.Remove(path)

		Path(path).File(file).Line(false).Stdout(false).Debug(1, 2, 3)
		content := gfile.GetContents(gfile.Join(path, file))
		t.Assert(gstr.Count(content, defaultLevelPrefixes[LEVEL_DEBU]), 1)
		t.Assert(gstr.Count(content, "1 2 3"), 1)
		t.Assert(gstr.Count(content, ".go"), 1)
		t.Assert(gstr.Contains(content, gfile.Separator), false)
	})
}

func Test_Async(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		path := gfile.TempDir(gtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gtime.TimestampNano())

		err := gfile.Mkdir(path)
		t.Assert(err, nil)
		defer gfile.Remove(path)

		Path(path).File(file).Async().Stdout(false).Debug(1, 2, 3)
		content := gfile.GetContents(gfile.Join(path, file))
		t.Assert(content, "")
		time.Sleep(200 * time.Millisecond)

		content = gfile.GetContents(gfile.Join(path, file))
		t.Assert(gstr.Count(content, defaultLevelPrefixes[LEVEL_DEBU]), 1)
		t.Assert(gstr.Count(content, "1 2 3"), 1)
	})

	gtest.C(t, func(t *gtest.T) {
		path := gfile.TempDir(gtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gtime.TimestampNano())

		err := gfile.Mkdir(path)
		t.Assert(err, nil)
		defer gfile.Remove(path)

		Path(path).File(file).Async(false).Stdout(false).Debug(1, 2, 3)
		content := gfile.GetContents(gfile.Join(path, file))
		t.Assert(gstr.Count(content, defaultLevelPrefixes[LEVEL_DEBU]), 1)
		t.Assert(gstr.Count(content, "1 2 3"), 1)
	})
}
