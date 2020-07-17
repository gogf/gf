// Copyright 2020 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package gvar_test

import (
	"github.com/jin502437344/gf/container/gvar"
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/test/gtest"
	"testing"
)

func Test_MapToMap(t *testing.T) {
	// map[int]int -> map[string]string
	// empty original map.
	gtest.C(t, func(t *gtest.T) {
		m1 := g.MapIntInt{}
		m2 := g.MapStrStr{}
		t.Assert(gvar.New(m1).MapToMap(&m2), nil)
		t.Assert(len(m1), len(m2))
	})
	// map[int]int -> map[string]string
	gtest.C(t, func(t *gtest.T) {
		m1 := g.MapIntInt{
			1: 100,
			2: 200,
		}
		m2 := g.MapStrStr{}
		t.Assert(gvar.New(m1).MapToMap(&m2), nil)
		t.Assert(m2["1"], m1[1])
		t.Assert(m2["2"], m1[2])
	})
	// map[string]interface{} -> map[string]string
	gtest.C(t, func(t *gtest.T) {
		m1 := g.Map{
			"k1": "v1",
			"k2": "v2",
		}
		m2 := g.MapStrStr{}
		t.Assert(gvar.New(m1).MapToMap(&m2), nil)
		t.Assert(m2["k1"], m1["k1"])
		t.Assert(m2["k2"], m1["k2"])
	})
	// map[string]string -> map[string]interface{}
	gtest.C(t, func(t *gtest.T) {
		m1 := g.MapStrStr{
			"k1": "v1",
			"k2": "v2",
		}
		m2 := g.Map{}
		t.Assert(gvar.New(m1).MapToMap(&m2), nil)
		t.Assert(m2["k1"], m1["k1"])
		t.Assert(m2["k2"], m1["k2"])
	})
	// map[string]interface{} -> map[interface{}]interface{}
	gtest.C(t, func(t *gtest.T) {
		m1 := g.MapStrStr{
			"k1": "v1",
			"k2": "v2",
		}
		m2 := g.MapAnyAny{}
		t.Assert(gvar.New(m1).MapToMap(&m2), nil)
		t.Assert(m2["k1"], m1["k1"])
		t.Assert(m2["k2"], m1["k2"])
	})
}
