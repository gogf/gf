package gtype_test

import (
	"github.com/gogf/gf/g/container/gtype"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func TestBytes_Clone(t *testing.T) {
	gtest.Case(t, func() {
		buf := make([]byte,100)
		for i:=0;i<100 ;i++  {
			buf[i] = byte(i)
		}
		oldBytes := gtype.NewBytes(buf)
		newBytes := oldBytes.Clone()
		gtest.AssertEQ(oldBytes,newBytes)
		gtest.AssertEQ(newBytes.Val(),buf)
	})
}

func TestBytes_Set(t *testing.T) {
	gtest.Case(t, func() {
		myBytes := gtype.NewBytes()
		buf := make([]byte,10)
		for i:=0;i<10 ;i++  {
			buf[i] = byte(i)
		}

		for i:=0 ;i<10 ; i++ {
			go func() {
				myBytes.Set(buf)
			}()
		}
	})
}


//并发安全错误示范
/*
func TestBytes_Set2(t *testing.T) {
	gtest.Case(t, func() {
		myBytes := []byte{}
		buf := make([]byte,10)
		for i:=0;i<10 ;i++  {
			buf[i] = byte(i)
		}

		for i:=0 ;i<2 ; i++ {
			go func() {
				myBytes = buf
			}()
		}

		gtest.AssertEQ(myBytes,buf)
	})
}
*/

func TestBytes_Val(t *testing.T) {
	gtest.Case(t, func() {
		buf := make([]byte,100)
		for i:=0;i<100 ;i++  {
			buf[i] = byte(i)
		}
		myBytes := gtype.NewBytes(buf)
		for i:=0;i<100 ;i++  {
			buf  = append(buf,byte(i))
		}
		myBytes.Set(buf)
		gtest.AssertEQ(myBytes.Val(),buf)
	})
}