// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtype_test

import (
	"encoding/json"
	"math"
	"sync"
	"testing"

	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/test/gtest"
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

	// Marshal
	gtest.Case(t, func() {
		i := gtype.NewBool(true)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)
	})
	gtest.Case(t, func() {
		i := gtype.NewBool(false)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)
	})
	// Unmarshal
	gtest.Case(t, func() {
		var err error
		i := gtype.NewBool()
		err = json.Unmarshal([]byte("true"), &i)
		gtest.Assert(err, nil)
		gtest.Assert(i.Val(), true)
		err = json.Unmarshal([]byte("false"), &i)
		gtest.Assert(err, nil)
		gtest.Assert(i.Val(), false)
		err = json.Unmarshal([]byte("1"), &i)
		gtest.Assert(err, nil)
		gtest.Assert(i.Val(), true)
		err = json.Unmarshal([]byte("0"), &i)
		gtest.Assert(err, nil)
		gtest.Assert(i.Val(), false)
	})

	gtest.Case(t, func() {
		i := gtype.NewBool(true)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)

		i2 := gtype.NewBool()
		err := json.Unmarshal(b2, &i2)
		gtest.Assert(err, nil)
		gtest.Assert(i2.Val(), i.Val())
	})
	gtest.Case(t, func() {
		i := gtype.NewBool(false)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)

		i2 := gtype.NewBool()
		err := json.Unmarshal(b2, &i2)
		gtest.Assert(err, nil)
		gtest.Assert(i2.Val(), i.Val())
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
	gtest.Case(t, func() {
		i := gtype.NewByte(49)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)
	})
	// Unmarshal
	gtest.Case(t, func() {
		var err error
		i := gtype.NewByte()
		err = json.Unmarshal([]byte("49"), &i)
		gtest.Assert(err, nil)
		gtest.Assert(i.Val(), "49")
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
	gtest.Case(t, func() {
		b := []byte("i love gf")
		i := gtype.NewBytes(b)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)

		i2 := gtype.NewBytes()
		err := json.Unmarshal(b2, &i2)
		gtest.Assert(err, nil)
		gtest.Assert(i2.Val(), b)
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
	gtest.Case(t, func() {
		s := "i love gf"
		i1 := gtype.NewString(s)
		b1, err1 := json.Marshal(i1)
		b2, err2 := json.Marshal(i1.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)

		i2 := gtype.NewString()
		err := json.Unmarshal(b2, &i2)
		gtest.Assert(err, nil)
		gtest.Assert(i2.Val(), s)
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
	gtest.Case(t, func() {
		s := "i love gf"
		i := gtype.New(s)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)

		i2 := gtype.New()
		err := json.Unmarshal(b2, &i2)
		gtest.Assert(err, nil)
		gtest.Assert(i2.Val(), s)
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
	gtest.Case(t, func() {
		v := float32(math.MaxFloat32)
		i := gtype.NewFloat32(v)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())

		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)

		i2 := gtype.NewFloat32()
		err := json.Unmarshal(b2, &i2)
		gtest.Assert(err, nil)
		gtest.Assert(i2.Val(), v)
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
	gtest.Case(t, func() {
		v := math.MaxFloat64
		i := gtype.NewFloat64(v)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)

		i2 := gtype.NewFloat64()
		err := json.Unmarshal(b2, &i2)
		gtest.Assert(err, nil)
		gtest.Assert(i2.Val(), v)
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
	gtest.Case(t, func() {
		v := 666
		i := gtype.NewInt(v)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)

		i2 := gtype.NewInt()
		err := json.Unmarshal(b2, &i2)
		gtest.Assert(err, nil)
		gtest.Assert(i2.Val(), v)
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
	gtest.Case(t, func() {
		v := int32(math.MaxInt32)
		i := gtype.NewInt32(v)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)

		i2 := gtype.NewInt32()
		err := json.Unmarshal(b2, &i2)
		gtest.Assert(err, nil)
		gtest.Assert(i2.Val(), v)
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
	gtest.Case(t, func() {
		i := gtype.NewInt64(math.MaxInt64)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)

		i2 := gtype.NewInt64()
		err := json.Unmarshal(b2, &i2)
		gtest.Assert(err, nil)
		gtest.Assert(i2.Val(), i)
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
	gtest.Case(t, func() {
		i := gtype.NewUint(666)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)

		i2 := gtype.NewUint()
		err := json.Unmarshal(b2, &i2)
		gtest.Assert(err, nil)
		gtest.Assert(i2.Val(), i)
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
	gtest.Case(t, func() {
		i := gtype.NewUint32(math.MaxUint32)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)

		i2 := gtype.NewUint32()
		err := json.Unmarshal(b2, &i2)
		gtest.Assert(err, nil)
		gtest.Assert(i2.Val(), i)
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
	gtest.Case(t, func() {
		i := gtype.NewUint64(math.MaxUint64)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)

		i2 := gtype.NewUint64()
		err := json.Unmarshal(b2, &i2)
		gtest.Assert(err, nil)
		gtest.Assert(i2.Val(), i)
	})
}
