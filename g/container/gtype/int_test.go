package gtype_test

import (
	"github.com/gogf/gf/g/container/gtype"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func TestInt_Clone(t *testing.T) {
	gtest.Case(t, func() {
		var number int = 32767

		oldNumber := gtype.NewInt(number)
		newNumber := oldNumber.Clone()
		gtest.AssertEQ(oldNumber,newNumber)
		gtest.AssertEQ(newNumber.Val(),number)
	})
}

func TestInt_Set(t *testing.T) {
	gtest.Case(t, func() {
		var number int = 32767
		myNumber := gtype.NewInt()
		for i:=0;i<2 ;i++  {
			go func() {
				myNumber.Set(number)
			}()
		}
	})
}

func TestInt_Val(t *testing.T) {
	gtest.Case(t, func() {
		var number int = 32767
		myNumber := gtype.NewInt()
		myNumber.Set(number)
		gtest.AssertEQ(myNumber.Val(),number)
	})
}

func TestInt_Add(t *testing.T) {
	gtest.Case(t, func() {
		var add int = 32767

		myNumber := gtype.NewInt(add)
		myAdd := myNumber.Add(int(0xa5a5))

		add += int(0xa5a5)
		gtest.AssertEQ(myAdd,add)
	})
}

func TestInt_Add2(t *testing.T) {
	gtest.Case(t, func() {
		myNumber := gtype.NewInt(int(32767))
		for i:=0;i<2 ;i++  {
			go func() {
				myNumber.Add(int(32767))
			}()

		}
	})
}
