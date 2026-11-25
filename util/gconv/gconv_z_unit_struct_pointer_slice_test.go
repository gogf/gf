// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//
// This file tests conversion into struct fields containing pointer slices.

package gconv_test

import (
	"encoding/json"
	"testing"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

// TestStructPointerSliceFields covers collection and scalar input for basic pointer elements.
func TestStructPointerSliceFields(t *testing.T) {
	type count int
	type params struct {
		Ints    []*int
		Strings []*string
		Bools   []*bool
		Floats  []*float64
		Counts  []*count
	}
	for _, test := range []struct {
		name  string
		input map[string]any
		want  string
	}{
		{
			name: "slices",
			input: map[string]any{
				"Ints": []string{"1", "2"}, "Strings": []int{1, 2},
				"Bools": []string{"true", "false"}, "Floats": []string{"1.5", "2.5"},
				"Counts": []int{1, 2},
			},
			want: `{"Ints":[1,2],"Strings":["1","2"],"Bools":[true,false],"Floats":[1.5,2.5],"Counts":[1,2]}`,
		},
		{
			name: "arrays",
			input: map[string]any{
				"Ints": [2]int{1, 2}, "Strings": [2]string{"a", "b"},
				"Bools": [2]bool{true, false}, "Floats": [2]float64{1.5, 2.5},
				"Counts": [2]int{1, 2},
			},
			want: `{"Ints":[1,2],"Strings":["a","b"],"Bools":[true,false],"Floats":[1.5,2.5],"Counts":[1,2]}`,
		},
		{
			name: "scalars",
			input: map[string]any{
				"Ints": 1, "Strings": "word", "Bools": true, "Floats": 1.5, "Counts": 3,
			},
			want: `{"Ints":[1],"Strings":["word"],"Bools":[true],"Floats":[1.5],"Counts":[3]}`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				var target params
				err := gconv.NewConverter().Struct(test.input, &target)
				t.AssertNil(err)
				encoded, err := json.Marshal(target)
				t.AssertNil(err)
				t.Assert(string(encoded), test.want)
			})
		})
	}
}

// TestStructPointerRawMessageSlice preserves JSON encoding for dereferenced RawMessage elements.
func TestStructPointerRawMessageSlice(t *testing.T) {
	for _, input := range []any{
		map[string]any{"hello": "world"},
		[]any{map[string]any{"hello": "world"}},
	} {
		gtest.C(t, func(t *gtest.T) {
			var target struct{ Values []*json.RawMessage }
			err := gconv.NewConverter().Struct(map[string]any{"Values": input}, &target)
			t.AssertNil(err)
			t.Assert(len(target.Values), 1)
			t.AssertNE(target.Values[0], nil)
			t.Assert(string(*target.Values[0]), `{"hello":"world"}`)
		})
	}
}

// TestStructPointerSliceFieldErrors keeps failed conversions from replacing the original field.
func TestStructPointerSliceFieldErrors(t *testing.T) {
	for _, input := range []any{"invalid", []string{"1", "invalid"}} {
		gtest.C(t, func(t *gtest.T) {
			var (
				original = 7
				target   = struct{ Values []*int }{Values: []*int{&original}}
				data     = map[string]any{"Values": input}
			)
			err := gconv.NewConverter().Struct(data, &target)
			t.AssertNE(err, nil)
			t.AssertNE(gerror.Code(err), gcode.CodeInternalPanic)
			t.Assert(len(target.Values), 1)
			t.Assert(target.Values[0] == &original, true)
			t.Assert(original, 7)
			err = gconv.Struct(data, &target)
			t.AssertNil(err)
			t.Assert(len(target.Values), 1)
			t.Assert(target.Values[0] == &original, true)
			t.Assert(original, 7)
		})
	}
}
