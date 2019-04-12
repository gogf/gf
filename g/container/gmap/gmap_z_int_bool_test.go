package gmap_test

import (
	"github.com/gogf/gf/g/container/gmap"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func getBool() bool {
	return true
}
func intBoolCallBack(int, bool) bool {
	return true
}
func Test_IntBoolMap_Basic(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewIntBoolMap()
		m.Set(1, true)

		gtest.Assert(m.Get(1), true)
		gtest.Assert(m.Size(), 1)
		gtest.Assert(m.IsEmpty(), false)

		gtest.Assert(m.GetOrSet(2, false), false)
		gtest.Assert(m.SetIfNotExist(2, false), false)

		gtest.Assert(m.SetIfNotExist(3, false), true)

		gtest.Assert(m.Remove(2), false)
		gtest.Assert(m.Contains(2), false)

		gtest.AssertIN(3, m.Keys())
		gtest.AssertIN(1, m.Keys())

		m.Clear()
		gtest.Assert(m.Size(), 0)
		gtest.Assert(m.IsEmpty(), true)

		m2 := gmap.NewIntBoolMapFrom(map[int]bool{1: true, 2: false})
		gtest.Assert(m2.Map(), map[int]bool{1: true, 2: false})
		m3 := gmap.NewIntBoolMapFromArray([]int{1, 2}, []bool{true, false})
		gtest.Assert(m3.Map(), map[int]bool{1: true, 2: false})

	})
}
func Test_IntBoolMap_Set_Fun(t *testing.T) {
	m := gmap.NewIntBoolMap()

	m.GetOrSetFunc(1, getBool)
	m.GetOrSetFuncLock(2, getBool)
	gtest.Assert(m.Get(1), true)
	gtest.Assert(m.Get(2), true)
	gtest.Assert(m.SetIfNotExistFunc(1, getBool), false)
	gtest.Assert(m.SetIfNotExistFunc(4, getBool), true)

	gtest.Assert(m.SetIfNotExistFuncLock(2, getBool), false)
	gtest.Assert(m.SetIfNotExistFuncLock(3, getBool), true)

}

func Test_IntBoolMap_Batch(t *testing.T) {
	m := gmap.NewIntBoolMap()

	m.BatchSet(map[int]bool{1: true, 2: false, 3: true})
	gtest.Assert(m.Map(), map[int]bool{1: true, 2: false, 3: true})
	m.BatchRemove([]int{1, 2})
	gtest.Assert(m.Map(), map[int]bool{3: true})
}
func Test_IntBoolMap_Iterator(t *testing.T){
	m := gmap.NewIntBoolMapFrom(map[int]bool{1: true, 2: false})
	m.Iterator(intBoolCallBack)
}

func Test_IntBoolMap_Lock(t *testing.T){
	m := gmap.NewIntBoolMapFrom(map[int]bool{1: true, 2: false})
	m.LockFunc(func(m map[int]bool) {})
	m.RLockFunc(func(m map[int]bool) {})
}

func Test_IntBoolMap_Clone(t *testing.T) {
	//clone 方法是深克隆
	m := gmap.NewIntBoolMapFrom(map[int]bool{1: true, 2: false})

	m_clone := m.Clone()
	m.Remove(1)
	//修改原 map,clone 后的 map 不影响
	gtest.AssertIN(1, m_clone.Keys())

	m_clone.Remove(2)
	//修改clone map,原 map 不影响
	gtest.AssertIN(2, m.Keys())
}
func Test_IntBoolMap_Merge(t *testing.T) {
	m1 := gmap.NewIntBoolMap()
	m2 := gmap.NewIntBoolMap()
	m1.Set(1, true)
	m2.Set(2, false)
	m1.Merge(m2)
	gtest.Assert(m1.Map(), map[int]bool{1: true, 2: false})
}
