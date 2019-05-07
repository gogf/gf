package gtype_test

import (
	"github.com/gogf/gf/g/container/gtype"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func TestString_Clone(t *testing.T) {
	gtest.Case(t, func() {
		oldStr := gtype.NewString("golang")
		newStr := oldStr.Clone()
		gtest.AssertEQ(oldStr,newStr)
	})
}

func TestString_Set(t *testing.T) {
	gtest.Case(t, func() {
		str:=gtype.NewString()
		for i:=0;i<2;i++{
			go func() {
				str.Set("golang")
			}()
		}
	})
}

func TestString_Val(t *testing.T) {
	gtest.Case(t, func() {
		var str string = "golang"
		myStr := gtype.NewString()
		myStr.Set(str)
		gtest.AssertEQ(myStr.Val(),str)
	})
}