package gtype_test

import (
	"github.com/gogf/gf/g/container/gtype"
	"github.com/gogf/gf/g/test/gtest"
	"math/rand"
	"testing"
	"time"
)

func TestByte_Clone(t *testing.T) {
	gtest.Case(t, func() {
		oldByte := gtype.NewByte(0xa5)
		newByte := oldByte.Clone()
		gtest.AssertEQ(oldByte,newByte)
		gtest.Assert(oldByte.Val(),0xa5)
	})
}

func TestByte_Set(t *testing.T) {
	gtest.Case(t, func() {
		rand := rand.New(rand.NewSource(time.Now().UnixNano()))
		myByte := gtype.NewByte()
		buf := [20]byte{}

		for i:=0;i<20 ;i++  {
			buf[i] = (byte)(rand.Int() %256)
		}

		// 根据随机值填充缓存区
		for i:=0 ;i<2 ; i++ {
			go func(i int) {
				for j:=0;j<10 ;j++  {
					myByte.Set(buf[i*10+j])
				}
			}(i)
		}
	})
}

func TestByte_Val(t *testing.T) {
	gtest.Case(t, func() {
		var testByte byte = 0xa5
		myByte := gtype.NewByte(testByte)
		gtest.AssertEQ(myByte.Val(),testByte)
		testByte = 0x5a
		myByte.Set(testByte)
		gtest.AssertEQ(myByte.Val(),testByte)
	})
}

func TestByte_Add(t *testing.T) {
	gtest.Case(t, func() {
		var add byte
		myByte := gtype.NewByte(11)
		myAdd := myByte.Add(58)
		add = 11+58
		gtest.Assert(myAdd,add)
	})
}

func TestByte_Add2(t *testing.T) {
	gtest.Case(t, func() {
		myByte := gtype.NewByte(11)
		for i:=0;i <2 ;i++  {
			go func() {
				myByte.Add(2)
			}()
		}
	})
}


//并发错误示范
/*
func TestByte_Add3(t *testing.T) {
	gtest.Case(t, func() {
		myByte := 32
		for i:=0;i <2 ;i++  {
			go func() {
				myByte += 2
			}()
		}
	})
}

 */