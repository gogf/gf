// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtime_test

import (
	"testing"

	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func Test_Json_Pointer(t *testing.T) {
	// Marshal
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Time *gtime.Time
		}
		t1 := new(T)
		s := "2006-01-02 15:04:05"
		t1.Time = gtime.NewFromStr(s)
		j, err := json.Marshal(t1)
		t.Assert(err, nil)
		t.Assert(j, `{"Time":"2006-01-02 15:04:05"}`)
	})
	// Marshal struct with embedded.
	gtest.C(t, func(t *gtest.T) {
		type Time struct {
			MyTime *gtime.Time
		}
		type T struct {
			Time
		}
		t1 := new(T)
		s := "2006-01-02 15:04:05"
		t1.MyTime = gtime.NewFromStr(s)
		j, err := json.Marshal(gconv.Map(t1))
		t.Assert(err, nil)
		t.Assert(j, `{"MyTime":"2006-01-02 15:04:05"}`)
	})
	// Marshal nil
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Time *gtime.Time
		}
		t1 := new(T)
		j, err := json.Marshal(t1)
		t.Assert(err, nil)
		t.Assert(j, `{"Time":null}`)
	})
	// Marshal nil omitempty
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Time *gtime.Time `json:"time,omitempty"`
		}
		t1 := new(T)
		j, err := json.Marshal(t1)
		t.Assert(err, nil)
		t.Assert(j, `{}`)
	})
	// Unmarshal
	gtest.C(t, func(t *gtest.T) {
		var t1 gtime.Time
		s := []byte(`"2006-01-02 15:04:05"`)
		err := json.UnmarshalUseNumber(s, &t1)
		t.Assert(err, nil)
		t.Assert(t1.String(), "2006-01-02 15:04:05")
	})
}

func Test_Json_Struct(t *testing.T) {
	// Marshal
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Time gtime.Time
		}
		t1 := new(T)
		s := "2006-01-02 15:04:05"
		t1.Time = *gtime.NewFromStr(s)
		j, err := json.Marshal(t1)
		t.Assert(err, nil)
		t.Assert(j, `{"Time":"2006-01-02 15:04:05"}`)
	})
	// Marshal nil
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Time gtime.Time
		}
		t1 := new(T)
		j, err := json.Marshal(t1)
		t.Assert(err, nil)
		t.Assert(j, `{"Time":""}`)
	})
	// Marshal nil omitempty
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Time gtime.Time `json:"time,omitempty"`
		}
		t1 := new(T)
		j, err := json.Marshal(t1)
		t.Assert(err, nil)
		t.Assert(j, `{"time":""}`)
	})

}
