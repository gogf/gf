// Copyright 2019 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package gtime_test

import (
	"github.com/jin502437344/gf/internal/json"
	"github.com/jin502437344/gf/os/gtime"
	"github.com/jin502437344/gf/test/gtest"
	"testing"
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
		err := json.Unmarshal(s, &t1)
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
