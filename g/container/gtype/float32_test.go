package gtype_test

import (
	"github.com/gogf/gf/g/container/gtype"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func TestFloat32_Clone(t *testing.T) {
	gtest.Case(t, func() {
		var f32 float32
		f32 = 3.1415926
		oldFloat32 := gtype.NewFloat32(f32)
		newFloat32 := oldFloat32.Clone()
		gtest.AssertEQ(oldFloat32,newFloat32)
		gtest.AssertEQ(newFloat32.Val(),f32)
	})
}

func TestFloat32_Set(t *testing.T) {
	gtest.Case(t, func() {
		var f32 float32
		f32 = 3.1415926
		myFloat32 := gtype.NewFloat32()
		myFloat32.Set(f32)
		for i:=0;i<2 ;i++  {
			go func() {
				myFloat32.Set(f32)
			}()
		}
	})
}

func TestFloat32_Val(t *testing.T) {
	gtest.Case(t, func() {
		var f32 float32 = 3.1415926
		myF32 := gtype.NewFloat32()
		myF32.Set(f32)
		gtest.AssertEQ(myF32.Val(),f32)
	})
}

func TestFloat32_Add(t *testing.T) {
	gtest.Case(t, func() {
		var add float32
		add = 3.1415926

		myF32 := gtype.NewFloat32(add)
		myAdd := myF32.Add(float32(6.2951413))

		add += float32(6.2951413)
		gtest.AssertEQ(myAdd,add)
	})
}

func TestFloat32_Add2(t *testing.T) {
	gtest.Case(t, func() {
		myF32 := gtype.NewFloat32(3.1415926)
		for i:=0;i<2 ;i++  {
			go func() {
				myF32.Add(6.2951413)
			}()

		}
	})
}

