package gtype_test

import (
	"github.com/gogf/gf/g/container/gtype"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func TestInt64_Clone(t *testing.T) {
	gtest.Case(t, func() {
		var number int64 = 32767

		oldNumber := gtype.NewInt64(number)
		newNumber := oldNumber.Clone()
		gtest.AssertEQ(oldNumber,newNumber)
		gtest.AssertEQ(newNumber.Val(),number)
	})
}

func TestInt64_Set(t *testing.T) {
	gtest.Case(t, func() {
		var number int64 = 32767
		myNumber := gtype.NewInt64()
		for i:=0;i<2 ;i++  {
			go func() {
				myNumber.Set(number)
			}()
		}
	})
}

func TestInt64_Val(t *testing.T) {
	gtest.Case(t, func() {
		var number int64 = 32767
		myNumber := gtype.NewInt64()
		myNumber.Set(number)
		gtest.AssertEQ(myNumber.Val(),number)
	})
}

func TestInt64_Add(t *testing.T) {
	gtest.Case(t, func() {
		var add int64 = 32767

		myNumber := gtype.NewInt64(add)
		myAdd := myNumber.Add(int64(0xa5a5))

		add += int64(0xa5a5)
		gtest.AssertEQ(myAdd,add)
	})
}

func TestInt64_Add2(t *testing.T) {
	gtest.Case(t, func() {
		myNumber := gtype.NewInt64(int64(32767))
		for i:=0;i<2 ;i++  {
			go func() {
				myNumber.Add(int64(32767))
			}()

		}
	})
}
