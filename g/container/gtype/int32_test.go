package gtype_test

import (
	"github.com/gogf/gf/g/container/gtype"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func TestInt32_Clone(t *testing.T) {
	gtest.Case(t, func() {
		var number int32 = 32767

		oldNumber := gtype.NewInt32(number)
		newNumber := oldNumber.Clone()
		gtest.AssertEQ(oldNumber,newNumber)
		gtest.AssertEQ(newNumber.Val(),number)
	})
}

func TestInt32_Set(t *testing.T) {
	gtest.Case(t, func() {
		var number int32 = 32767
		myNumber := gtype.NewInt32()
		for i:=0;i<2 ;i++  {
			go func() {
				myNumber.Set(number)
			}()
		}
	})
}

func TestInt32_Val(t *testing.T) {
	gtest.Case(t, func() {
		var number int32 = 32767
		myNumber := gtype.NewInt32()
		myNumber.Set(number)
		gtest.AssertEQ(myNumber.Val(),number)
	})
}

func TestInt32_Add(t *testing.T) {
	gtest.Case(t, func() {
		var add int32 = 32767

		myNumber := gtype.NewInt32(add)
		myAdd := myNumber.Add(int32(0xa5a5))

		add += int32(0xa5a5)
		gtest.AssertEQ(myAdd,add)
	})
}

func TestInt32_Add2(t *testing.T) {
	gtest.Case(t, func() {
		myNumber := gtype.NewInt32(int32(32767))
		for i:=0;i<2 ;i++  {
			go func() {
				myNumber.Add(int32(32767))
			}()

		}
	})
}
