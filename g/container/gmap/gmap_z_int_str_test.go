// Copyright 2017-2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap_test

import (
	"github.com/gogf/gf/g/container/gmap"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func getStr() string {
	return "z"
}
func intStrCallBack(int, string) bool {
	return true
}
func Test_IntStrMap_Basic(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewIntStrMap()
		m.Set(1, "a")

		gtest.Assert(m.Get(1), "a")
		gtest.Assert(m.Size(), 1)
		gtest.Assert(m.IsEmpty(), false)

		gtest.Assert(m.GetOrSet(2, "b"), "b")
		gtest.Assert(m.SetIfNotExist(2, "b"), false)

		gtest.Assert(m.SetIfNotExist(3, "c"), true)

		gtest.Assert(m.Remove(2), "b")
		gtest.Assert(m.Contains(2), false)

		gtest.AssertIN(3, m.Keys())
		gtest.AssertIN(1, m.Keys())
		gtest.AssertIN("a", m.Values())
		gtest.AssertIN("c", m.Values())

		//反转之后不成为以下 map,flip 操作只是翻转原 map
		//gtest.Assert(m.Map(), map[string]int{"a": 1, "c": 3})
		m_f := gmap.NewIntStrMap()
		m_f.Set(1, "2")
		m_f.Flip()
		gtest.Assert(m_f.Map(), map[int]string{2: "1"})

		m.Clear()
		gtest.Assert(m.Size(), 0)
		gtest.Assert(m.IsEmpty(), true)

		m2 := gmap.NewIntStrMapFrom(map[int]string{1: "a", 2: "b"})
		gtest.Assert(m2.Map(), map[int]string{1: "a", 2: "b"})
	})
}
func Test_IntStrMap_Set_Fun(t *testing.T) {
	m := gmap.NewIntStrMap()

	m.GetOrSetFunc(1, getStr)
	m.GetOrSetFuncLock(2, getStr)
	gtest.Assert(m.Get(1), "z")
	gtest.Assert(m.Get(2), "z")
	gtest.Assert(m.SetIfNotExistFunc(1, getStr), false)
	gtest.Assert(m.SetIfNotExistFunc(3, getStr), true)

	gtest.Assert(m.SetIfNotExistFuncLock(2, getStr), false)
	gtest.Assert(m.SetIfNotExistFuncLock(4, getStr), true)

}

func Test_IntStrMap_Batch(t *testing.T) {
	m := gmap.NewIntStrMap()

	m.Sets(map[int]string{1: "a", 2: "b", 3: "c"})
	gtest.Assert(m.Map(), map[int]string{1: "a", 2: "b",3: "c"})
	m.Removes([]int{1, 2})
	gtest.Assert(m.Map(), map[int]interface{}{3: "c"})
}
func Test_IntStrMap_Iterator(t *testing.T){
	expect := map[int]string{1: "a", 2: "b"}
	m      := gmap.NewIntStrMapFrom(expect)
	m.Iterator(func(k int, v string) bool {
		gtest.Assert(expect[k], v)
		return true
	})
	// 断言返回值对遍历控制
	i := 0
	j := 0
	m.Iterator(func(k int, v string) bool {
		i++
		return true
	})
	m.Iterator(func(k int, v string) bool {
		j++
		return false
	})
	gtest.Assert(i, 2)
	gtest.Assert(j, 1)
}

func Test_IntStrMap_Lock(t *testing.T){

	expect := map[int]string{1: "a", 2: "b", 3: "c"}
	m      := gmap.NewIntStrMapFrom(expect)
	m.LockFunc(func(m map[int]string) {
		gtest.Assert(m, expect)
	})
	m.RLockFunc(func(m map[int]string) {
		gtest.Assert(m, expect)
	})

}
func Test_IntStrMap_Clone(t *testing.T) {
	//clone 方法是深克隆
	m := gmap.NewIntStrMapFrom(map[int]string{1: "a", 2: "b", 3: "c"})

	m_clone := m.Clone()
	m.Remove(1)
	//修改原 map,clone 后的 map 不影响
	gtest.AssertIN(1, m_clone.Keys())

	m_clone.Remove(2)
	//修改clone map,原 map 不影响
	gtest.AssertIN(2, m.Keys())
}
func Test_IntStrMap_Merge(t *testing.T) {
	m1 := gmap.NewIntStrMap()
	m2 := gmap.NewIntStrMap()
	m1.Set(1, "a")
	m2.Set(2, "b")
	m1.Merge(m2)
	gtest.Assert(m1.Map(), map[int]string{1: "a", 2: "b"})
}
