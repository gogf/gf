// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gview_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gview"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Config(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		config := gview.Config{
			Paths: []string{gdebug.TestDataPath("config")},
			Data: g.Map{
				"name": "gf",
			},
			DefaultFile: "test.html",
			Delimiters:  []string{"${", "}"},
		}
		view := gview.New()
		err := view.SetConfig(config)
		t.Assert(err, nil)

		str := `hello ${.name},version:${.version}`
		view.Assigns(g.Map{"version": "1.7.0"})
		result, err := view.ParseContent(context.TODO(), str, nil)
		t.Assert(err, nil)
		t.Assert(result, "hello gf,version:1.7.0")

		result, err = view.ParseDefault(context.TODO())
		t.Assert(err, nil)
		t.Assert(result, "name:gf")
	})
}

func Test_ConfigWithMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		view := gview.New()
		err := view.SetConfigWithMap(g.Map{
			"Paths":       []string{gdebug.TestDataPath("config")},
			"DefaultFile": "test.html",
			"Delimiters":  []string{"${", "}"},
			"Data": g.Map{
				"name": "gf",
			},
		})
		t.Assert(err, nil)

		str := `hello ${.name},version:${.version}`
		view.Assigns(g.Map{"version": "1.7.0"})
		result, err := view.ParseContent(context.TODO(), str, nil)
		t.Assert(err, nil)
		t.Assert(result, "hello gf,version:1.7.0")

		result, err = view.ParseDefault(context.TODO())
		t.Assert(err, nil)
		t.Assert(result, "name:gf")
	})
}
