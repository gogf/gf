// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap_test

import (
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gutil"
	"testing"
)

func Test_Scan(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		anyMap := gmap.NewAnyAnyMap()
		anyMap.Set("1", "2")
		anyMap.Set("3", "4")
		anyMap.Set("5", "6")
		gMap := gmap.ScanGMap[string, string](anyMap)
		t.Assert(gMap, anyMap.Map())

		strMap := gmap.NewStrAnyMap()
		strMap.Set("1", "2")
		strMap.Set("3", "4")
		strMap.Set("5", "6")
		strAny := gmap.ScanGStrAny[string](strMap)
		t.Assert(strAny, anyMap.MapStrAny())
		intMap := gmap.NewIntAnyMap()
		intMap.Set(1, "2")
		intMap.Set(3, "4")
		intMap.Set(5, "6")
		intAny := gmap.ScanGStrAny[string](strMap)
		t.Assert(intAny, intMap.MapStrAny())
		listMap := gmap.NewListMap()
		listMap.Set(1, "2")
		listMap.Set(3, "4")
		listMap.Set(5, "6")
		gListMap := gmap.ScanGListMap[int, string](listMap)
		t.Assert(gListMap, listMap.Map())
		treeMap := gmap.NewTreeMap(gutil.ComparatorString)
		treeMap.Set(1, "2")
		treeMap.Set(3, "4")
		treeMap.Set(5, "6")
		gTreeMap := gmap.ScanGTreeMap[int, string](treeMap)
		t.Assert(gTreeMap, treeMap.Map())
	})
}
