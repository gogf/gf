package gtype_test

import (
	"github.com/gogf/gf/g/container/gtype"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func TestUint_Clone(t *testing.T) {
	gtest.Case(t, func() {
		var number uint = 32767

		oldNumber := gtype.NewUint(number)
		newNumber := oldNumber.Clone()
		gtest.AssertEQ(oldNumber,newNumber)
		gtest.AssertEQ(newNumber.Val(),number)
	})
}

func TestUint_Set(t *testing.T) {
	gtest.Case(t, func() {
		var number uint = 32767
		myNumber := gtype.NewUint()
		for i:=0;i<2 ;i++  {
			go func() {
				myNumber.Set(number)
			}()
		}
	})
}

func TestUint_Val(t *testing.T) {
	gtest.Case(t, func() {
		var number uint = 32767
		myNumber := gtype.NewUint()
		myNumber.Set(number)
		gtest.AssertEQ(myNumber.Val(),number)
	})
}

func TestUint_Add(t *testing.T) {
	gtest.Case(t, func() {
		var add uint = 32767

		myNumber := gtype.NewUint(add)
		myAdd := myNumber.Add(uint(0xa5a5))

		add += uint(0xa5a5)
		gtest.AssertEQ(myAdd,add)
	})
}

func TestUint_Add2(t *testing.T) {
	gtest.Case(t, func() {
		myNumber := gtype.NewUint(uint(32767))
		for i:=0;i<2 ;i++  {
			go func() {
				myNumber.Add(uint(32767))
			}()

		}
	})
}
