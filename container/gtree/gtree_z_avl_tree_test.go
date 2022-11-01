// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gtree_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/container/gtree"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gutil"
)

func Test_AVLTree_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gtree.NewAVLTree(gutil.ComparatorString)
		m.Set("key1", "val1")
		t.Assert(m.Keys(), []interface{}{"key1"})

		t.Assert(m.Get("key1"), "val1")
		t.Assert(m.Size(), 1)
		t.Assert(m.IsEmpty(), false)

		t.Assert(m.GetOrSet("key2", "val2"), "val2")
		t.Assert(m.GetOrSet("key2", "val2"), "val2")
		t.Assert(m.SetIfNotExist("key2", "val2"), false)

		t.Assert(m.SetIfNotExist("key3", "val3"), true)

		t.Assert(m.Remove("key2"), "val2")
		t.Assert(m.Contains("key2"), false)

		t.AssertIN("key3", m.Keys())
		t.AssertIN("key1", m.Keys())
		t.AssertIN("val3", m.Values())
		t.AssertIN("val1", m.Values())

		m.Sets(map[interface{}]interface{}{"key3": "val3", "key1": "val1"})

		m.Flip()
		t.Assert(m.Map(), map[interface{}]interface{}{"val3": "key3", "val1": "key1"})

		m.Flip(gutil.ComparatorString)
		t.Assert(m.Map(), map[interface{}]interface{}{"key3": "val3", "key1": "val1"})

		m.Clear()
		t.Assert(m.Size(), 0)
		t.Assert(m.IsEmpty(), true)

		m2 := gtree.NewAVLTreeFrom(gutil.ComparatorString, map[interface{}]interface{}{1: 1, "key1": "val1"})
		t.Assert(m2.Map(), map[interface{}]interface{}{1: 1, "key1": "val1"})
	})
}

func Test_AVLTree_Set_Fun(t *testing.T) {
	//GetOrSetFunc lock or unlock
	gtest.C(t, func(t *gtest.T) {
		m := gtree.NewAVLTree(gutil.ComparatorString)
		t.Assert(m.GetOrSetFunc("fun", getValue), 3)
		t.Assert(m.GetOrSetFunc("fun", getValue), 3)
		t.Assert(m.GetOrSetFuncLock("funlock", getValue), 3)
		t.Assert(m.GetOrSetFuncLock("funlock", getValue), 3)
		t.Assert(m.Get("funlock"), 3)
		t.Assert(m.Get("fun"), 3)
	})
	//SetIfNotExistFunc lock or unlock
	gtest.C(t, func(t *gtest.T) {
		m := gtree.NewAVLTree(gutil.ComparatorString)
		t.Assert(m.SetIfNotExistFunc("fun", getValue), true)
		t.Assert(m.SetIfNotExistFunc("fun", getValue), false)
		t.Assert(m.SetIfNotExistFuncLock("funlock", getValue), true)
		t.Assert(m.SetIfNotExistFuncLock("funlock", getValue), false)
		t.Assert(m.Get("funlock"), 3)
		t.Assert(m.Get("fun"), 3)
	})

}

func Test_AVLTree_Get_Set_Var(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gtree.NewAVLTree(gutil.ComparatorString)
		t.AssertEQ(m.SetIfNotExist("key1", "val1"), true)
		t.AssertEQ(m.SetIfNotExist("key1", "val1"), false)
		t.AssertEQ(m.GetVarOrSet("key1", "val1"), gvar.New("val1", true))
		t.AssertEQ(m.GetVar("key1"), gvar.New("val1", true))
	})

	gtest.C(t, func(t *gtest.T) {
		m := gtree.NewAVLTree(gutil.ComparatorString)
		t.AssertEQ(m.GetVarOrSetFunc("fun", getValue), gvar.New(3, true))
		t.AssertEQ(m.GetVarOrSetFunc("fun", getValue), gvar.New(3, true))
		t.AssertEQ(m.GetVarOrSetFuncLock("funlock", getValue), gvar.New(3, true))
		t.AssertEQ(m.GetVarOrSetFuncLock("funlock", getValue), gvar.New(3, true))
	})
}

func Test_AVLTree_Batch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := gtree.NewAVLTree(gutil.ComparatorString)
		m.Sets(map[interface{}]interface{}{1: 1, "key1": "val1", "key2": "val2", "key3": "val3"})
		t.Assert(m.Map(), map[interface{}]interface{}{1: 1, "key1": "val1", "key2": "val2", "key3": "val3"})
		m.Removes([]interface{}{"key1", 1})
		t.Assert(m.Map(), map[interface{}]interface{}{"key2": "val2", "key3": "val3"})
	})
}

func Test_AVLTree_Iterator(t *testing.T) {

	keys := []string{"1", "key1", "key2", "key3", "key4"}
	keyLen := len(keys)
	index := 0

	expect := map[interface{}]interface{}{"key4": "val4", 1: 1, "key1": "val1", "key2": "val2", "key3": "val3"}

	m := gtree.NewAVLTreeFrom(gutil.ComparatorString, expect)

	gtest.C(t, func(t *gtest.T) {
		m.Iterator(func(k interface{}, v interface{}) bool {
			t.Assert(k, keys[index])
			index++
			t.Assert(expect[k], v)
			return true
		})

		m.IteratorDesc(func(k interface{}, v interface{}) bool {
			index--
			t.Assert(k, keys[index])
			t.Assert(expect[k], v)
			return true
		})
	})

	m.Print()
	// 断言返回值对遍历控制
	gtest.C(t, func(t *gtest.T) {
		i := 0
		j := 0
		m.Iterator(func(k interface{}, v interface{}) bool {
			i++
			return true
		})
		m.Iterator(func(k interface{}, v interface{}) bool {
			j++
			return false
		})
		t.Assert(i, keyLen)
		t.Assert(j, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		i := 0
		j := 0
		m.IteratorDesc(func(k interface{}, v interface{}) bool {
			i++
			return true
		})
		m.IteratorDesc(func(k interface{}, v interface{}) bool {
			j++
			return false
		})
		t.Assert(i, keyLen)
		t.Assert(j, 1)
	})

}

func Test_AVLTree_IteratorFrom(t *testing.T) {
	m := make(map[interface{}]interface{})
	for i := 1; i <= 10; i++ {
		m[i] = i * 10
	}
	tree := gtree.NewAVLTreeFrom(gutil.ComparatorInt, m)

	gtest.C(t, func(t *gtest.T) {
		n := 5
		tree.IteratorFrom(5, true, func(key, value interface{}) bool {
			t.Assert(n, key)
			t.Assert(n*10, value)
			n++
			return true
		})

		i := 5
		tree.IteratorAscFrom(5, true, func(key, value interface{}) bool {
			t.Assert(i, key)
			t.Assert(i*10, value)
			i++
			return true
		})

		j := 5
		tree.IteratorDescFrom(5, true, func(key, value interface{}) bool {
			t.Assert(j, key)
			t.Assert(j*10, value)
			j--
			return true
		})
	})
}

func Test_AVLTree_Clone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		//clone 方法是深克隆
		m := gtree.NewAVLTreeFrom(gutil.ComparatorString, map[interface{}]interface{}{1: 1, "key1": "val1"})
		m_clone := m.Clone()
		m.Remove(1)
		//修改原 map,clone 后的 map 不影响
		t.AssertIN(1, m_clone.Keys())

		m_clone.Remove("key1")
		//修改clone map,原 map 不影响
		t.AssertIN("key1", m.Keys())
	})
}

func Test_AVLTree_LRNode(t *testing.T) {
	expect := map[interface{}]interface{}{"key4": "val4", "key1": "val1", "key2": "val2", "key3": "val3"}
	//safe
	gtest.C(t, func(t *gtest.T) {
		m := gtree.NewAVLTreeFrom(gutil.ComparatorString, expect)
		t.Assert(m.Left().Key, "key1")
		t.Assert(m.Right().Key, "key4")
	})
	//unsafe
	gtest.C(t, func(t *gtest.T) {
		m := gtree.NewAVLTreeFrom(gutil.ComparatorString, expect, true)
		t.Assert(m.Left().Key, "key1")
		t.Assert(m.Right().Key, "key4")
	})
}

func Test_AVLTree_CeilingFloor(t *testing.T) {
	expect := map[interface{}]interface{}{
		20: "val20",
		6:  "val6",
		10: "val10",
		12: "val12",
		1:  "val1",
		15: "val15",
		19: "val19",
		8:  "val8",
		4:  "val4"}
	//found and eq
	gtest.C(t, func(t *gtest.T) {
		m := gtree.NewAVLTreeFrom(gutil.ComparatorInt, expect)
		c, cf := m.Ceiling(8)
		t.Assert(cf, true)
		t.Assert(c.Value, "val8")
		f, ff := m.Floor(20)
		t.Assert(ff, true)
		t.Assert(f.Value, "val20")
	})
	//found and neq
	gtest.C(t, func(t *gtest.T) {
		m := gtree.NewAVLTreeFrom(gutil.ComparatorInt, expect)
		c, cf := m.Ceiling(9)
		t.Assert(cf, true)
		t.Assert(c.Value, "val10")
		f, ff := m.Floor(5)
		t.Assert(ff, true)
		t.Assert(f.Value, "val4")
	})
	//nofound
	gtest.C(t, func(t *gtest.T) {
		m := gtree.NewAVLTreeFrom(gutil.ComparatorInt, expect)
		c, cf := m.Ceiling(21)
		t.Assert(cf, false)
		t.Assert(c, nil)
		f, ff := m.Floor(-1)
		t.Assert(ff, false)
		t.Assert(f, nil)
	})
}

func Test_AVLTree_Remove(t *testing.T) {
	m := gtree.NewAVLTree(gutil.ComparatorInt)
	for i := 1; i <= 50; i++ {
		m.Set(i, fmt.Sprintf("val%d", i))
	}
	expect := m.Map()
	gtest.C(t, func(t *gtest.T) {
		for k, v := range expect {
			m1 := m.Clone()
			t.Assert(m1.Remove(k), v)
			t.Assert(m1.Remove(k), nil)
		}
	})
}
