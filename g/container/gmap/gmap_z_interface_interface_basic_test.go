package gmap_test

import (
	"github.com/gogf/gf/g/container/gmap"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)
func getValue()interface{}{
	return 3
}

func Test_Map_Basic(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.New()
		m.Set("key1", "val1")
		gtest.Assert(m.Keys(),[]interface{}{"key1"})
		gtest.Assert(m.Values(),[]interface{}{"val1"})

		gtest.Assert(m.Get("key1"), "val1")
		m.BatchSet(map[interface{}]interface{}{1: 1, "key2": "val2", "key3": "val3"})
		gtest.Assert(m.Size(), 4)
		gtest.Assert(m.IsEmpty(), false)
		gtest.Assert(m.GetOrSet("key4", "val4"), "val4")
		gtest.Assert(m.SetIfNotExist("key4", "val4"), false)
		gtest.Assert(m.Remove("key2"), "val2")
		m.BatchRemove([]interface{}{"key1", 1})
		gtest.Assert(m.Contains("key3"), true)
		m.Flip()
		gtest.Assert(m.Map(), map[interface{}]interface{}{"val3": "key3", "val4": "key4"})
		m.GetOrSetFunc("fun",getValue)
		m.GetOrSetFuncLock("funlock",getValue)
		gtest.Assert(m.Get("funlock"),3)
		gtest.Assert(m.Get("fun"),3)
		m.GetOrSetFunc("fun",getValue)
		gtest.Assert(m.SetIfNotExistFunc("fun",getValue),false)
		gtest.Assert(m.SetIfNotExistFuncLock("funlock",getValue),false)

		m.Clear()
		gtest.Assert(m.Size(), 0)
		m2 := gmap.NewFrom(map[interface{}]interface{}{1: 1, "key1": "val1"})
		gtest.Assert(m2.Map(), map[interface{}]interface{}{1: 1, "key1": "val1"})
		m3 := gmap.NewFromArray([]interface{}{1, "key1"}, []interface{}{1, "val1"})
		gtest.Assert(m3.Map(), map[interface{}]interface{}{1: 1, "key1": "val1"})
		m4 := m3.Clone()
		gtest.Assert(m4.Map(), map[interface{}]interface{}{1: 1, "key1": "val1"})

	})
}
func Test_Map_Basic_Merge(t *testing.T) {
	m1 := gmap.New()
	m2 := gmap.New()
	m1.Set("key1", "val1")
	m2.Set("key2", "val2")
	m1.Merge(m2)
	gtest.Assert(m1.Map(), map[interface{}]interface{}{"key1": "val1", "key2": "val2"})
}
