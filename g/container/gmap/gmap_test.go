// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*"

package gmap_test

import (
    "testing"
    "gitee.com/johng/gf/g/container/gmap"
    "strconv"
)


var ibm   = gmap.NewIntBoolMap()
var iim   = gmap.NewIntIntMap()
var iifm  = gmap.NewIntInterfaceMap()
var ism   = gmap.NewIntStringMap()
var ififm = gmap.NewInterfaceInterfaceMap()
var sbm   = gmap.NewStringBoolMap()
var sim   = gmap.NewStringIntMap()
var sifm  = gmap.NewStringInterfaceMap()
var ssm   = gmap.NewStringStringMap()
var uifm  = gmap.NewUintInterfaceMap()

func BenchmarkIntBoolMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ibm.Set(i, true)
    }
}

func BenchmarkIntIntMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        iim.Set(i, i)
    }
}

func BenchmarkIntInterfaceMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        iifm.Set(i, i)
    }
}

func BenchmarkIntStringMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ism.Set(i, strconv.Itoa(i))
    }
}

func BenchmarkInterfaceInterfaceMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ififm.Set(i, i)
    }
}

func BenchmarkStringBoolMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sbm.Set(strconv.Itoa(i), true)
    }
}

func BenchmarkStringIntMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sim.Set(strconv.Itoa(i), i)
    }
}

func BenchmarkStringInterfaceMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sifm.Set(strconv.Itoa(i), i)
    }
}

func BenchmarkStringStringMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ssm.Set(strconv.Itoa(i), strconv.Itoa(i))
    }
}

func BenchmarkUintInterfaceMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        uifm.Set(uint(i), i)
    }
}

