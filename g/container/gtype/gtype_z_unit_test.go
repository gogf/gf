package gtype_test

import (
	"sync"
	"testing"

	"github.com/gogf/gf/g/container/gtype"
	"github.com/gogf/gf/g/test/gtest"
)

type Temp struct {
	Name string
	Age  int
}

func Test_Bool(t *testing.T) {
	gtest.Case(t, func() {
		i := gtype.NewBool(true)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(false), true)
		gtest.AssertEQ(iClone.Val(), false)

		i1 := gtype.NewBool(false)
		iClone1 := i1.Clone()
		gtest.AssertEQ(iClone1.Set(true), false)
		gtest.AssertEQ(iClone1.Val(), true)

		//空参测试
		i2 := gtype.NewBool()
		gtest.AssertEQ(i2.Val(), false)
	})
}

func Test_Byte(t *testing.T) {
	gtest.Case(t, func() {
		var wg sync.WaitGroup
		addTimes := 127
		i := gtype.NewByte(byte(0))
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(byte(1)), byte(0))
		gtest.AssertEQ(iClone.Val(), byte(1))
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		gtest.AssertEQ(byte(addTimes), i.Val())

		//空参测试
		i1 := gtype.NewByte()
		gtest.AssertEQ(i1.Val(), byte(0))
	})
}

func Test_Bytes(t *testing.T) {
	gtest.Case(t, func() {
		i := gtype.NewBytes([]byte("abc"))
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set([]byte("123")), []byte("abc"))
		gtest.AssertEQ(iClone.Val(), []byte("123"))

		//空参测试
		i1 := gtype.NewBytes()
		gtest.AssertEQ(i1.Val(), nil)
	})
}

func Test_String(t *testing.T) {
	gtest.Case(t, func() {
		i := gtype.NewString("abc")
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set("123"), "abc")
		gtest.AssertEQ(iClone.Val(), "123")

		//空参测试
		i1 := gtype.NewString()
		gtest.AssertEQ(i1.Val(), "")
	})
}

func Test_Interface(t *testing.T) {
	gtest.Case(t, func() {
		t := Temp{Name: "gf", Age: 18}
		t1 := Temp{Name: "gf", Age: 19}
		i := gtype.New(t)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(t1), t)
		gtest.AssertEQ(iClone.Val().(Temp), t1)

		//空参测试
		i1 := gtype.New()
		gtest.AssertEQ(i1.Val(), nil)
	})
}

func Test_Float32(t *testing.T) {
	gtest.Case(t, func() {
		//var wg sync.WaitGroup
		//addTimes := 100
		i := gtype.NewFloat32(0)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(0.1), float32(0))
		gtest.AssertEQ(iClone.Val(), float32(0.1))
		// for index := 0; index < addTimes; index++ {
		// 	wg.Add(1)
		// 	go func() {
		// 	defer wg.Done()
		// 	i.Add(0.2)
		// 	fmt.Println(i.Val())
		// 	}()
		// }
		// wg.Wait()
		// gtest.AssertEQ(100.0, i.Val())

		//空参测试
		i1 := gtype.NewFloat32()
		gtest.AssertEQ(i1.Val(), float32(0))
	})
}

func Test_Float64(t *testing.T) {
	gtest.Case(t, func() {
		//var wg sync.WaitGroup
		//addTimes := 100
		i := gtype.NewFloat64(0)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(0.1), float64(0))
		gtest.AssertEQ(iClone.Val(), float64(0.1))
		// for index := 0; index < addTimes; index++ {
		// 	wg.Add(1)
		// 	go func() {
		// 	defer wg.Done()
		// 	i.Add(0.1)
		// 	fmt.Println(i.Val())
		// 	}()
		// }
		// wg.Wait()
		// gtest.AssertEQ(100.0, i.Val())

		//空参测试
		i1 := gtype.NewFloat64()
		gtest.AssertEQ(i1.Val(), float64(0))
	})
}

func Test_Int(t *testing.T) {
	gtest.Case(t, func() {
		var wg sync.WaitGroup
		addTimes := 1000
		i := gtype.NewInt(0)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(1), 0)
		gtest.AssertEQ(iClone.Val(), 1)
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		gtest.AssertEQ(addTimes, i.Val())

		//空参测试
		i1 := gtype.NewInt()
		gtest.AssertEQ(i1.Val(), 0)
	})
}

func Test_Int32(t *testing.T) {
	gtest.Case(t, func() {
		var wg sync.WaitGroup
		addTimes := 1000
		i := gtype.NewInt32(0)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(1), int32(0))
		gtest.AssertEQ(iClone.Val(), int32(1))
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		gtest.AssertEQ(int32(addTimes), i.Val())

		//空参测试
		i1 := gtype.NewInt32()
		gtest.AssertEQ(i1.Val(), int32(0))
	})
}

func Test_Int64(t *testing.T) {
	gtest.Case(t, func() {
		var wg sync.WaitGroup
		addTimes := 1000
		i := gtype.NewInt64(0)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(1), int64(0))
		gtest.AssertEQ(iClone.Val(), int64(1))
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		gtest.AssertEQ(int64(addTimes), i.Val())

		//空参测试
		i1 := gtype.NewInt64()
		gtest.AssertEQ(i1.Val(), int64(0))
	})
}

func Test_Uint(t *testing.T) {
	gtest.Case(t, func() {
		var wg sync.WaitGroup
		addTimes := 1000
		i := gtype.NewUint(0)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(1), uint(0))
		gtest.AssertEQ(iClone.Val(), uint(1))
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		gtest.AssertEQ(uint(addTimes), i.Val())

		//空参测试
		i1 := gtype.NewUint()
		gtest.AssertEQ(i1.Val(), uint(0))
	})
}

func Test_Uint32(t *testing.T) {
	gtest.Case(t, func() {
		var wg sync.WaitGroup
		addTimes := 1000
		i := gtype.NewUint32(0)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(1), uint32(0))
		gtest.AssertEQ(iClone.Val(), uint32(1))
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		gtest.AssertEQ(uint32(addTimes), i.Val())

		//空参测试
		i1 := gtype.NewUint32()
		gtest.AssertEQ(i1.Val(), uint32(0))
	})
}

func Test_Uint64(t *testing.T) {
	gtest.Case(t, func() {
		var wg sync.WaitGroup
		addTimes := 1000
		i := gtype.NewUint64(0)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(1), uint64(0))
		gtest.AssertEQ(iClone.Val(), uint64(1))
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		gtest.AssertEQ(uint64(addTimes), i.Val())

		//空参测试
		i1 := gtype.NewUint64()
		gtest.AssertEQ(i1.Val(), uint64(0))
	})
}
