package gring_test

import (
	"container/ring"
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
		gtest.AssertEQ(r.Val().(*Student).name,"jimmy")
		//从当前位置往后移两个元素
		r.Move(2)
		gtest.AssertEQ(r.Val().(*Student).name,"alon")
		if r.Val().(*Student).upgrade == false {
			r.Val().(*Student).upgrade = true
			//更新元素值
			r.Set(&Student{3, "jack", true})
		}


	})
}
func TestRing_CapLen(t *testing.T) {
	gtest.Case(t, func() {
		r := gring.New(10)
		r.Put("goframe")
		//cap长度 10
		gtest.AssertEQ(r.Cap(), 10)
		//已有数据项 1
		gtest.AssertEQ(r.Len(), 1)
	})
}

func TestRing_Position(t *testing.T) {
	gtest.Case(t, func() {
		r := gring.New(2)
		r.Put(1)
		r.Put(2)
		//往后移动1个元素
		r.Next()
		gtest.AssertEQ(r.Val(),2)
		//往前移动1个元素
		r.Prev()
		gtest.AssertEQ(r.Val(),1)

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
		gtest.AssertEQ(rs.Move(2).Val(), "b")

	})
}

func TestRing_Unlink(t *testing.T) {
	gtest.Case(t, func() {
		r := gring.New(5)
		for i := 0; i< 5; i++  {
			r.Put(i+1)
		}
		// 1 2 3 4 5
		// 删除当前位置往后的2个数据，返回被删除的数据
		s := r.Unlink(2)		// 2 3
		gtest.AssertEQ(s.Val(), 2)
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
		gtest.AssertEQ(array[0], 3)
		gtest.AssertEQ(len(array), 5)

		//判断array是否等于[3 4 5 1 2]
		ra := []int{3,4,5,1,2}
		eq := true
		for i, v := range array {
			if v != ra[i] {
				eq = false
			}
		}
		gtest.AssertEQ(eq, true)

		//第3个元素设为nil
		r.Set(nil)
		array2 := r.SliceNext() 	//[4 5 1 2]
		//返回当前位置往后不为空的元素数组，长度为4
		gtest.AssertEQ(len(array2), 4)

		array3 := r.SlicePrev() 	//[2 1 5 4]
		//数组array3第一个元素为2
		gtest.AssertEQ(array3[0], 2)
		//数组长度4
		gtest.AssertEQ(len(array3), 4)

	})
}

func TestRing_RLockIterator(t *testing.T) {
	gtest.Case(t, func() {
		ringLen := 5
		r := gring.New(ringLen)
		for i := 0; i< ringLen; i++  {
			r.Put(i+1)
		}
		var i,j int
		//回调函数返回true,RLockIteratorNext遍历5次
		r.RLockIteratorNext(func(value interface{}) bool {
			i++;
			return true
		})
		gtest.AssertEQ(i, 5)

		//RLockIteratorPrev遍历3次返回 false,退出遍历
		r.RLockIteratorPrev(func(value interface{}) bool {
			if j++; j < 3 {
				return true
			}
			return false
		})
		gtest.AssertEQ(j, 3)
	})
}

func TestRing_LockIterator(t *testing.T) {
	gtest.Case(t, func() {
		ringLen := 5
		r := gring.New(ringLen)
		for i := 0; i< ringLen; i++  {
			r.Put(i+1)
		}
		var i,j int
		r.LockIteratorNext(func(item *ring.Ring) bool {
			i++;
			return true
		})
		gtest.AssertEQ(i, 5)

		r.LockIteratorPrev(func(item *ring.Ring) bool {
			if j++; j < 3 {
				return true
			}
			return false
		})

		gtest.AssertEQ(j, 3)
	})
}