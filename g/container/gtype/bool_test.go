package gtype_test

import (
	"github.com/gogf/gf/g/container/gtype"
	"github.com/gogf/gf/g/test/gtest"
	"math/rand"
	"testing"
	"time"
)

func TestNewBool(t *testing.T) {
	gtest.Case(t, func() {
		myBool1 := gtype.NewBool(true)
		gtest.Assert(myBool1.Val(),true)

		myBool2 := gtype.NewBool(false)
		gtest.Assert(myBool2.Val(),false)
	})
}

func TestBool_Clone(t *testing.T) {
	gtest.Case(t, func() {
		oldBool := gtype.NewBool(true)
		newBool := oldBool.Clone()
		gtest.AssertEQ(oldBool,newBool)
	})
}


func TestBool_Set(t *testing.T) {
	gtest.Case(t, func() {
		rand := rand.New(rand.NewSource(time.Now().UnixNano()))
		myBool := gtype.NewBool()
		buf := [100]bool{}

		for i:=0;i<100 ;i++  {
			if rand.Int()%2 == 0 {
				buf[i] = false
			} else {
				buf[i] = true
			}
		}

		// 根据随机值填充缓存区
		for i:=0;i<10 ; i++ {
			go func(i int) {
				for j:=0;j<10 ;j++  {
					myBool.Set(buf[i*10+j])
				}
			}(i)
		}
	})
}

func TestBool_Val(t *testing.T) {
	gtest.Case(t, func() {
		myBool := gtype.NewBool(true)
		gtest.AssertEQ(myBool.Val(),true)
		myBool.Set(false)
		gtest.AssertEQ(myBool.Val(),false)
	})
}




