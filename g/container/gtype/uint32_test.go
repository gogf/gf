package gtype_test

import (
	"github.com/gogf/gf/g/container/gtype"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func TestUint32_Clone(t *testing.T) {
	gtest.Case(t, func() {
		var number uint32 = 32767

		oldNumber := gtype.NewUint32(number)
		newNumber := oldNumber.Clone()
		gtest.AssertEQ(oldNumber,newNumber)
		gtest.AssertEQ(newNumber.Val(),number)
	})
}

func TestUint32_Set(t *testing.T) {
	gtest.Case(t, func() {
		var number uint32 = 32767
		myNumber := gtype.NewUint32()
		for i:=0;i<2 ;i++  {
			go func() {
				myNumber.Set(number)
			}()
		}
	})
}

func TestUint32_Val(t *testing.T) {
	gtest.Case(t, func() {
		var number uint32 = 32767
		myNumber := gtype.NewUint32()
		myNumber.Set(number)
		gtest.AssertEQ(myNumber.Val(),number)
	})
}

func TestUint32_Add(t *testing.T) {
	gtest.Case(t, func() {
		var add uint32 = 32767

		myNumber := gtype.NewUint32(add)
		myAdd := myNumber.Add(uint32(0xa5a5))

		add += uint32(0xa5a5)
		gtest.AssertEQ(myAdd,add)
	})
}

func TestUint32_Add2(t *testing.T) {
	gtest.Case(t, func() {
		myNumber := gtype.NewUint32(uint32(32767))
		for i:=0;i<2 ;i++  {
			go func() {
				myNumber.Add(uint32(32767))
			}()

		}
	})
}
