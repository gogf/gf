package gtype_test

import (
	"testing"

	"github.com/gogf/gf/g/container/gtype"
	"github.com/gogf/gf/g/test/gtest"
)

func TestUint64_Clone(t *testing.T) {
	gtest.Case(t, func() {
		var number uint64 = 32767

		oldNumber := gtype.NewUint64(number)
		newNumber := oldNumber.Clone()
		gtest.AssertEQ(oldNumber, newNumber)
		gtest.AssertEQ(newNumber.Val(), number)
	})
}

func TestUint64_Set(t *testing.T) {
	gtest.Case(t, func() {
		var number uint64 = 32767
		myNumber := gtype.NewUint64()
		for i := 0; i < 2; i++ {
			go func() {
				myNumber.Set(number)
			}()
		}
	})
}

func TestUint64_Val(t *testing.T) {
	gtest.Case(t, func() {
		var number uint64 = 32767
		myNumber := gtype.NewUint64()
		myNumber.Set(number)
		gtest.AssertEQ(myNumber.Val(), number)
	})
}

func TestUint64_Add(t *testing.T) {
	gtest.Case(t, func() {
		var add uint64 = 32767

		myNumber := gtype.NewUint64(add)
		myAdd := myNumber.Add(uint64(0xa5a5))

		add += uint64(0xa5a5)
		gtest.AssertEQ(myAdd, add)
	})
}

func TestUint64_Add2(t *testing.T) {
	gtest.Case(t, func() {
		myNumber := gtype.NewUint64(uint64(32767))
		for i := 0; i < 2; i++ {
			go func() {
				myNumber.Add(uint64(32767))
			}()

		}
	})
}
