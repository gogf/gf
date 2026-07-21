// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

// textArray mimics types like uuid.UUID ([16]byte with UnmarshalText).
type textArray [4]byte

func (a *textArray) UnmarshalText(text []byte) error {
	for i := 0; i < len(a) && i < len(text); i++ {
		a[i] = text[i]
	}
	return nil
}

func Test_Issue4786_SliceOfArrayWithUnmarshalText(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Req struct {
			IDs []textArray `json:"ids"`
		}
		var r Req
		err := gconv.Struct(map[string]any{
			"ids": []any{"abcd", "wxyz"},
		}, &r)
		t.AssertNil(err)
		t.Assert(len(r.IDs), 2)
		t.Assert(string(r.IDs[0][:]), "abcd")
		t.Assert(string(r.IDs[1][:]), "wxyz")
	})
}

func Test_Issue4786_SliceOfUUID(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Req struct {
			IDs []uuid.UUID `json:"ids"`
		}
		id1 := "550e8400-e29b-41d4-a716-446655440000"
		id2 := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
		var r Req
		err := gconv.Struct(map[string]any{
			"ids": []any{id1, id2},
		}, &r)
		t.AssertNil(err)
		t.Assert(len(r.IDs), 2)
		t.Assert(r.IDs[0].String(), id1)
		t.Assert(r.IDs[1].String(), id2)
	})
}

func Test_Issue4786_PointerSliceOfUUID(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Req struct {
			IDs []*uuid.UUID `json:"ids"`
		}
		id1 := "550e8400-e29b-41d4-a716-446655440000"
		var r Req
		err := gconv.Struct(map[string]any{
			"ids": []any{id1},
		}, &r)
		t.AssertNil(err)
		t.Assert(len(r.IDs), 1)
		t.AssertNE(r.IDs[0], nil)
		t.Assert(r.IDs[0].String(), id1)
	})
}

func Test_Issue4786_SingleUUIDStillWorks(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Req struct {
			ID uuid.UUID `json:"id"`
		}
		id1 := "550e8400-e29b-41d4-a716-446655440000"
		var r Req
		err := gconv.Struct(map[string]any{"id": id1}, &r)
		t.AssertNil(err)
		t.Assert(r.ID.String(), id1)
	})
}
