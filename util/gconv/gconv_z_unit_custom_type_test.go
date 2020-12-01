// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gconv"
	"testing"
	"time"
)

type Duration time.Duration

// UnmarshalText unmarshal text to duration.
func (d *Duration) UnmarshalText(text []byte) error {
	tmp, err := time.ParseDuration(string(text))
	if err == nil {
		*d = Duration(tmp)
	}
	return err
}

func Test_Struct_CustomTimeDuration_Attribute(t *testing.T) {
	type A struct {
		Name    string
		Timeout Duration
	}
	gtest.C(t, func(t *gtest.T) {
		var a A
		err := gconv.Struct(g.Map{
			"name":    "john",
			"timeout": "1s",
		}, &a)
		t.Assert(err, nil)
	})
}
