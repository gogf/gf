// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtest_test

import (
	"errors"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

var (
	map1           = map[string]string{"k1": "v1"}
	map1Expect     = map[string]string{"k1": "v1"}
	map2           = map[string]string{"k2": "v2"}
	mapLong1       = map[string]string{"k1": "v1", "k2": "v2"}
	mapLong1Expect = map[string]string{"k2": "v2", "k1": "v1"}
)

func TestC(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(1, 1)
		t.AssertNE(1, 0)
		t.AssertEQ(float32(123.456), float32(123.456))
		t.AssertEQ(float32(123.456), float32(123.456))
		t.Assert(map[string]string{"1": "1"}, map[string]string{"1": "1"})
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, "[ASSERT] EXPECT 1 == 0")
			}
		}()
		t.Assert(1, 0)
	})
}

func TestCase(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(1, 1)
		t.AssertNE(1, 0)
		t.AssertEQ(float32(123.456), float32(123.456))
		t.AssertEQ(float32(123.456), float32(123.456))
	})
}

func TestAssert(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		var (
			nilChan chan struct{}
		)
		t.Assert(1, 1)
		t.Assert(nilChan, nil)
		t.Assert(map1, map1Expect)
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, `[ASSERT] EXPECT VALUE map["k2"]: == map["k2"]:v2
GIVEN : map[k1:v1]
EXPECT: map[k2:v2]`)
			}
		}()
		t.Assert(map1, map2)
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, `[ASSERT] EXPECT MAP LENGTH 2 == 1`)
			}
		}()
		t.Assert(mapLong1, map2)
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, `[ASSERT] EXPECT VALUE TO BE A MAP, BUT GIVEN "int"`)
			}
		}()
		t.Assert(0, map1)
	})
}

func TestAssertEQ(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		var (
			nilChan chan struct{}
		)
		t.AssertEQ(nilChan, nil)
		t.AssertEQ("0", "0")
		t.AssertEQ(float32(123.456), float32(123.456))
		t.AssertEQ(mapLong1, mapLong1Expect)
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, "[ASSERT] EXPECT 1 == 0")
			}
		}()
		t.AssertEQ(1, 0)
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, "[ASSERT] EXPECT TYPE 1[int] == 1[string]")
			}
		}()
		t.AssertEQ(1, "1")
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, `[ASSERT] EXPECT VALUE map["k2"]: == map["k2"]:v2
GIVEN : map[k1:v1]
EXPECT: map[k2:v2]`)
			}
		}()
		t.AssertEQ(map1, map2)
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, `[ASSERT] EXPECT MAP LENGTH 2 == 1`)
			}
		}()
		t.AssertEQ(mapLong1, map2)
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, `[ASSERT] EXPECT VALUE TO BE A MAP, BUT GIVEN "int"`)
			}
		}()
		t.AssertEQ(0, map1)
	})
}

func TestAssertNE(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			c = make(chan struct{}, 1)
		)
		t.AssertNE(nil, c)
		t.AssertNE("0", "1")
		t.AssertNE(float32(123.456), float32(123.4567))
		t.AssertNE(map1, map2)
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, "[ASSERT] EXPECT 1 != 1")
			}
		}()
		t.AssertNE(1, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, `[ASSERT] EXPECT map[k1:v1] != map[k1:v1]`)
			}
		}()
		t.AssertNE(map1, map1Expect)
	})
}

func TestAssertNQ(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertNQ(1, "0")
		t.AssertNQ(float32(123.456), float64(123.4567))
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, "[ASSERT] EXPECT 1 != 1")
			}
		}()
		t.AssertNQ(1, "1")
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, "[ASSERT] EXPECT TYPE 1[int] != 1[int]")
			}
		}()
		t.AssertNQ(1, 1)
	})
}

func TestAssertGT(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertGT("b", "a")
		t.AssertGT(1, -1)
		t.AssertGT(uint(1), uint(0))
		t.AssertGT(float32(123.45678), float32(123.4567))
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, "[ASSERT] EXPECT -1 > 1")
			}
		}()
		t.AssertGT(-1, 1)
	})
}

func TestAssertGE(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertGE("b", "a")
		t.AssertGE("a", "a")
		t.AssertGE(1, -1)
		t.AssertGE(1, 1)
		t.AssertGE(uint(1), uint(0))
		t.AssertGE(uint(0), uint(0))
		t.AssertGE(float32(123.45678), float32(123.4567))
		t.AssertGE(float32(123.456), float32(123.456))
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, "[ASSERT] EXPECT -1(int) >= 1(int)")
			}
		}()
		t.AssertGE(-1, 1)
	})
}

func TestAssertLT(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertLT("a", "b")
		t.AssertLT(-1, 1)
		t.AssertLT(uint(0), uint(1))
		t.AssertLT(float32(123.456), float32(123.4567))
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, "[ASSERT] EXPECT 1 < -1")
			}
		}()
		t.AssertLT(1, -1)
	})
}

func TestAssertLE(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertLE("a", "b")
		t.AssertLE("a", "a")
		t.AssertLE(-1, 1)
		t.AssertLE(1, 1)
		t.AssertLE(uint(0), uint(1))
		t.AssertLE(uint(0), uint(0))
		t.AssertLE(float32(123.456), float32(123.4567))
		t.AssertLE(float32(123.456), float32(123.456))
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, "[ASSERT] EXPECT 1 <= -1")
			}
		}()
		t.AssertLE(1, -1)
	})
}

func TestAssertIN(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertIN("a", []string{"a", "b", "c"})
		t.AssertIN(1, []int{1, 2, 3})
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, "[ASSERT] INVALID EXPECT VALUE TYPE: int")
			}
		}()
		t.AssertIN(0, 0)
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, "[ASSERT] EXPECT 4 IN [1 2 3]")
			}
		}()
		// t.AssertIN(0, []int{0, 1, 2, 3})
		// t.AssertIN(0, []int{ 1, 2, 3})
		t.AssertIN(4, []int{1, 2, 3})
	})
}

func TestAssertNI(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertNI("d", []string{"a", "b", "c"})
		t.AssertNI(4, []int{1, 2, 3})
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, "[ASSERT] INVALID EXPECT VALUE TYPE: int")
			}
		}()
		t.AssertNI(0, 0)
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, "[ASSERT] EXPECT 1 NOT IN [1 2 3]")
			}
		}()
		t.AssertNI(1, []int{1, 2, 3})
	})
}

func TestAssertNil(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			nilChan chan struct{}
		)
		t.AssertNil(nilChan)
		_, err := strconv.ParseInt("123", 10, 64)
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, "error")
			}
		}()
		t.AssertNil(errors.New("error"))
	})

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.AssertNE(err, nil)
			}
		}()
		t.AssertNil(1)
	})
}

func TestAssertError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer func() {
			if err := recover(); err != nil {
				t.Assert(err, "[ERROR] this is an error")
			}
		}()
		t.Error("this is an error")
	})
}

func TestDataPath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(filepath.ToSlash(gtest.DataPath("testdata.txt")), `./testdata/testdata.txt`)
	})
}

func TestDataContent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gtest.DataContent("testdata.txt"), `hello`)
		t.Assert(gtest.DataContent(""), "")
	})
}
