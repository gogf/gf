package gtype_test

import (
"github.com/gogf/gf/g/container/gtype"
"github.com/gogf/gf/g/test/gtest"
"testing"
)

func TestFloat64_Clone(t *testing.T) {
	gtest.Case(t, func() {
		var f64 float64
		f64 = 3.1415926
		oldFloat64 := gtype.NewFloat64(f64)
		newFloat64 := oldFloat64.Clone()
		gtest.AssertEQ(oldFloat64,newFloat64)
		gtest.AssertEQ(newFloat64.Val(),f64)
	})
}

func TestFloat64_Set(t *testing.T) {
	gtest.Case(t, func() {
		var f64 float64
		f64 = 3.1415926
		myFloat64 := gtype.NewFloat64()
		for i:=0;i<2 ;i++  {
			go func() {
				myFloat64.Set(f64)
			}()
		}
	})
}

func TestFloat64_Val(t *testing.T) {
	gtest.Case(t, func() {
		var f64 float64 = 3.1415926
		myF64 := gtype.NewFloat64()
		myF64.Set(f64)
		gtest.AssertEQ(myF64.Val(),f64)
	})
}

func TestFloat64_Add(t *testing.T) {
	gtest.Case(t, func() {
		var add float64
		add = 3.1415926

		myF32 := gtype.NewFloat64(add)
		myAdd := myF32.Add(float64(6.2951413))

		add += float64(6.2951413)
		gtest.AssertEQ(myAdd,add)
	})
}

func TestFloat64_Add2(t *testing.T) {
	gtest.Case(t, func() {
		myF64 := gtype.NewFloat64(float64(6.2951413))
		for i:=0;i<2 ;i++  {
			go func() {
				myF64.Add(float64(6.2951413))
			}()

		}
	})
}

