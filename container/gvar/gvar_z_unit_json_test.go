// Copyright 2020 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package gvar_test

import (
	"github.com/jin502437344/gf/container/gvar"
	"github.com/jin502437344/gf/internal/json"
	"github.com/jin502437344/gf/test/gtest"
	"math"
	"testing"
)

func Test_Json(t *testing.T) {
	// Marshal
	gtest.C(t, func(t *gtest.T) {
		s := "i love gf"
		v := gvar.New(s)
		b1, err1 := json.Marshal(v)
		b2, err2 := json.Marshal(s)
		t.Assert(err1, err2)
		t.Assert(b1, b2)
	})

	gtest.C(t, func(t *gtest.T) {
		s := int64(math.MaxInt64)
		v := gvar.New(s)
		b1, err1 := json.Marshal(v)
		b2, err2 := json.Marshal(s)
		t.Assert(err1, err2)
		t.Assert(b1, b2)
	})

	// Unmarshal
	gtest.C(t, func(t *gtest.T) {
		s := "i love gf"
		v := gvar.New(nil)
		b, err := json.Marshal(s)
		t.Assert(err, nil)

		err = json.Unmarshal(b, v)
		t.Assert(err, nil)
		t.Assert(v.String(), s)
	})

	gtest.C(t, func(t *gtest.T) {
		var v gvar.Var
		s := "i love gf"
		b, err := json.Marshal(s)
		t.Assert(err, nil)

		err = json.Unmarshal(b, &v)
		t.Assert(err, nil)
		t.Assert(v.String(), s)
	})
}
