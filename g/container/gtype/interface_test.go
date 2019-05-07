package gtype_test

import (
	"github.com/gogf/gf/g/container/gtype"
	"github.com/gogf/gf/g/test/gtest"
	"reflect"
	"testing"
)

type basicVar struct {
	num int
	str string
	flag bool
	myFloat float64
}

func TestInterface_Clone(t *testing.T) {
	gtest.Case(t, func() {
		bVar := basicVar{6553,"golang",false,3.141926}
		 oldV := gtype.NewInterface(bVar)
		 newV := oldV.Clone()

		 gtest.AssertEQ(oldV,newV)
	})
}

func TestInterface_Set(t *testing.T) {
	gtest.Case(t, func() {
		 myV := gtype.NewInterface(123)

		 for i:=0;i<2;i++{
		 	go func() {
		 		myV.Set("golang");
			}()
		 }
	})
}

func TestInterface_Val(t *testing.T) {
	gtest.Case(t, func() {
		bVar := basicVar{6553,"golang",false,3.141926}
		oldV := gtype.NewInterface(bVar)
		myValue := reflect.ValueOf(oldV.Val())
		gtest.Assert(bVar,myValue)
	})
}