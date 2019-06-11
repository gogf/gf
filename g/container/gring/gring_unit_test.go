package gring_test

import (
	"container/ring"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/container/gring"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)
type Student struct {
	position int
	name    string
	upgrade bool
}

func TestRing_Val(t *testing.T) {
	gtest.Case(t, func() {
		//定义cap 为3的ring类型数据
		r := gring.New(3, true)
		//分别给3个元素初始化赋值
		r.Put(&Student{1,"jimmy", true})
		r.Put(&Student{2,"tom", true})
		r.Put(&Student{3,"alon", false})

		//元素取值并判断和预设值是否相等
		gtest.Assert(r.Val().(*Student).name,"jimmy")
		//从当前位置往后移两个元素
		r.Move(2)
		gtest.Assert(r.Val().(*Student).name,"alon")
		//更新元素值
		//测试 value == nil
		r.Set(nil)
		gtest.Assert(r.Val(),nil)
		//测试value != nil
		r.Set(&Student{3, "jack", true})
	})
}
func TestRing_CapLen(t *testing.T) {
	gtest.Case(t, func() {
		r := gring.New(10)
		r.Put("goframe")
		//cap长度 10
		gtest.Assert(r.Cap(), 10)
		//已有数据项 1
		gtest.Assert(r.Len(), 1)
	})
}

func TestRing_Position(t *testing.T) {
	gtest.Case(t, func() {
		r := gring.New(2)
		r.Put(1)
		r.Put(2)
		//往后移动1个元素
		r.Next()
		gtest.Assert(r.Val(),2)
		//往前移动1个元素
		r.Prev()
		gtest.Assert(r.Val(),1)

	})
}

func TestRing_Link(t *testing.T) {
	gtest.Case(t, func() {
		r := gring.New(3)
		r.Put(1)
		r.Put(2)
		r.Put(3)
		s := gring.New(2)
		s.Put("a")
		s.Put("b")

		rs := r.Link(s)
		gtest.Assert(rs.Move(2).Val(), "b")

	})
}

func TestRing_Unlink(t *testing.T) {
	gtest.Case(t, func() {
		r := gring.New(5)
		for i := 0; i< 5; i++  {
			r.Put(i+1)
		}
		// 1 2 3 4
		// 删除当前位置往后的2个数据，返回被删除的数据
		// 重新计算s len
		s := r.Unlink(2)		// 2 3
		gtest.Assert(s.Val(), 2)
		gtest.Assert(s.Len(), 1)
	})
}

func TestRing_Slice(t *testing.T) {
	gtest.Case(t, func() {
		ringLen := 5
		r := gring.New(ringLen)
		for i := 0; i< ringLen; i++  {
			r.Put(i+1)
		}
		r.Move(2)	// 3
		array := r.SliceNext()		// [3 4 5 1 2]
		gtest.Assert(array[0], 3)
		gtest.Assert(len(array), 5)

		//判断array是否等于[3 4 5 1 2]
		ra := []int{3,4,5,1,2}
		gtest.Assert(ra, array)

		//第3个元素设为nil
		r.Set(nil)
		array2 := r.SliceNext() 	//[4 5 1 2]
		//返回当前位置往后不为空的元素数组，长度为4
		gtest.Assert(array2, g.Slice{4,5,1,2})

		array3 := r.SlicePrev() 	//[2 1 5 4]
		gtest.Assert(array3, g.Slice{2,1,5,4})

		s := gring.New(ringLen)
		for i := 0; i< ringLen; i++  {
			s.Put(i+1)
		}
		array4 := s.SlicePrev()	// []
		gtest.Assert(array4, g.Slice{1,5,4,3,2})

	})
}

func TestRing_RLockIterator(t *testing.T) {
	gtest.Case(t, func() {
		ringLen := 5
		r := gring.New(ringLen)

		//ring不存在有值元素
		r.RLockIteratorNext(func(v interface{}) bool {
			gtest.Assert(v, nil)
			return false
		})
		r.RLockIteratorNext(func(v interface{}) bool {
			gtest.Assert(v, nil)
			return true
		})

		r.RLockIteratorPrev(func(v interface{}) bool {
			gtest.Assert(v, nil)
			return true
		})

		for i := 0; i< ringLen; i++  {
			r.Put(i+1)
		}

		//回调函数返回true,RLockIteratorNext遍历5次,期望值分别是1、2、3、4、5
		i := 0
		r.RLockIteratorNext(func(v interface{}) bool {
			gtest.Assert(v, i+1)
			i++;
			return true
		})

		//RLockIteratorPrev遍历1次返回 false,退出遍历
		r.RLockIteratorPrev(func(v interface{}) bool {
			gtest.Assert(v, 1)
			return false
		})

	})
}

func TestRing_LockIterator(t *testing.T) {
	gtest.Case(t, func() {
		ringLen := 5
		r := gring.New(ringLen)

		//不存在有值元素
		r.LockIteratorNext(func(item *ring.Ring) bool {
			gtest.Assert(item.Value, nil)
			return false
		})
		r.LockIteratorNext(func(item *ring.Ring) bool {
			gtest.Assert(item.Value, nil)
			return false
		})
		r.LockIteratorNext(func(item *ring.Ring) bool {
			gtest.Assert(item.Value, nil)
			return true
		})

		r.LockIteratorPrev(func(item *ring.Ring) bool {
			gtest.Assert(item.Value, nil)
			return false
		})
		r.LockIteratorPrev(func(item *ring.Ring) bool {
			gtest.Assert(item.Value, nil)
			return true
		})

		//ring初始化元素值
		for i := 0; i< ringLen; i++  {
			r.Put(i+1)
		}

		//往后遍历组成数据 [1,2,3,4,5]
		array1 := g.Slice{1,2,3,4,5}
		ii := 0
		r.LockIteratorNext(func(item *ring.Ring) bool {
			//校验每一次遍历取值是否是期望值
			gtest.Assert(item.Value, array1[ii])
			ii++;
			return true
		})

		//往后取3个元素组成数组
		//获得 [1,5,4]
		i := 0
		a := g.Slice{1,5,4}
		r.LockIteratorPrev(func(item *ring.Ring) bool {
			if i > 2 {
				return false
			}
			gtest.Assert(item.Value, a[i])
			i++;
			return true
		})


	})
}