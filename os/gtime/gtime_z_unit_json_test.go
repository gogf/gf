// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtime_test

import (
	"encoding/json"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
	"testing"
)

func Test_Json_Pointer(t *testing.T) {
	// Marshal
	gtest.Case(t, func() {
		type T struct {
			Time *gtime.Time
		}
		t := new(T)
		s := "2006-01-02 15:04:05"
		t.Time = gtime.NewFromStr(s)
		j, err := json.Marshal(t)
		gtest.Assert(err, nil)
		gtest.Assert(j, `{"Time":"2006-01-02 15:04:05"}`)
	})
	// Marshal nil
	gtest.Case(t, func() {
		type T struct {
			Time *gtime.Time
		}
		t := new(T)
		j, err := json.Marshal(t)
		gtest.Assert(err, nil)
		gtest.Assert(j, `{"Time":null}`)
	})
	// Marshal nil omitempty
	gtest.Case(t, func() {
		type T struct {
			Time *gtime.Time `json:"time,omitempty"`
		}
		t := new(T)
		j, err := json.Marshal(t)
		gtest.Assert(err, nil)
		gtest.Assert(j, `{}`)
	})
	// Unmarshal
	gtest.Case(t, func() {
		var t gtime.Time
		s := []byte(`"2006-01-02 15:04:05"`)
		err := json.Unmarshal(s, &t)
		gtest.Assert(err, nil)
		gtest.Assert(t.String(), "2006-01-02 15:04:05")
	})
}

func Test_Json_Struct(t *testing.T) {
	// Marshal
	gtest.Case(t, func() {
		type T struct {
			Time gtime.Time
		}
		t := new(T)
		s := "2006-01-02 15:04:05"
		t.Time = *gtime.NewFromStr(s)
		j, err := json.Marshal(t)
		gtest.Assert(err, nil)
		gtest.Assert(j, `{"Time":"2006-01-02 15:04:05"}`)
	})
	// Marshal nil
	gtest.Case(t, func() {
		type T struct {
			Time gtime.Time
		}
		t := new(T)
		j, err := json.Marshal(t)
		gtest.Assert(err, nil)
		gtest.Assert(j, `{"Time":""}`)
	})
	// Marshal nil omitempty
	gtest.Case(t, func() {
		type T struct {
			Time gtime.Time `json:"time,omitempty"`
		}
		t := new(T)
		j, err := json.Marshal(t)
		gtest.Assert(err, nil)
		gtest.Assert(j, `{"time":""}`)
	})

}
