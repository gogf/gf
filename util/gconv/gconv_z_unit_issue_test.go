// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

// https://github.com/gogf/gf/issues/1227
func Test_Issue1227(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type StructFromIssue1227 struct {
			Name string `json:"n1"`
		}
		tests := []struct {
			name   string
			origin interface{}
			want   string
		}{
			{
				name:   "Case1",
				origin: `{"n1":"n1"}`,
				want:   "n1",
			},
			{
				name:   "Case2",
				origin: `{"name":"name"}`,
				want:   "",
			},
			{
				name:   "Case3",
				origin: `{"NaMe":"NaMe"}`,
				want:   "",
			},
			{
				name:   "Case4",
				origin: g.Map{"n1": "n1"},
				want:   "n1",
			},
			{
				name:   "Case5",
				origin: g.Map{"NaMe": "n1"},
				want:   "n1",
			},
		}
		for _, tt := range tests {
			p := StructFromIssue1227{}
			if err := gconv.Struct(tt.origin, &p); err != nil {
				t.Error(err)
			}
			t.Assert(p.Name, tt.want)
		}
	})

	// Chinese key.
	gtest.C(t, func(t *gtest.T) {
		type StructFromIssue1227 struct {
			Name string `json:"中文Key"`
		}
		tests := []struct {
			name   string
			origin interface{}
			want   string
		}{
			{
				name:   "Case1",
				origin: `{"中文Key":"n1"}`,
				want:   "n1",
			},
			{
				name:   "Case2",
				origin: `{"Key":"name"}`,
				want:   "",
			},
			{
				name:   "Case3",
				origin: `{"NaMe":"NaMe"}`,
				want:   "",
			},
			{
				name:   "Case4",
				origin: g.Map{"中文Key": "n1"},
				want:   "n1",
			},
			{
				name:   "Case5",
				origin: g.Map{"中文KEY": "n1"},
				want:   "",
			},
			{
				name:   "Case5",
				origin: g.Map{"KEY": "n1"},
				want:   "",
			},
		}
		for _, tt := range tests {
			p := StructFromIssue1227{}
			if err := gconv.Struct(tt.origin, &p); err != nil {
				t.Error(err)
			}
			//t.Log(tt)
			t.Assert(p.Name, tt.want)
		}
	})
}

// https://github.com/gogf/gf/issues/1607
type issue1607Float64 float64

func (f *issue1607Float64) UnmarshalValue(value interface{}) error {
	if v, ok := value.(*big.Rat); ok {
		f64, _ := v.Float64()
		*f = issue1607Float64(f64)
	}
	return nil
}

func Test_Issue1607(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Demo struct {
			B issue1607Float64
		}
		rat := &big.Rat{}
		rat.SetFloat64(1.5)

		var demos = make([]Demo, 1)
		err := gconv.Scan([]map[string]interface{}{
			{"A": 1, "B": rat},
		}, &demos)
		t.AssertNil(err)
		t.Assert(demos[0].B, 1.5)
	})
}

// https://github.com/gogf/gf/issues/1946
func Test_Issue1946(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			init *gtype.Bool
			Name string
		}
		type A struct {
			B *B
		}
		a := &A{
			B: &B{
				init: gtype.NewBool(true),
			},
		}
		err := gconv.Struct(g.Map{
			"B": g.Map{
				"Name": "init",
			},
		}, a)
		t.AssertNil(err)
		t.Assert(a.B.Name, "init")
		t.Assert(a.B.init.Val(), true)
	})
	// It cannot change private attribute.
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			init *gtype.Bool
			Name string
		}
		type A struct {
			B *B
		}
		a := &A{
			B: &B{
				init: gtype.NewBool(true),
			},
		}
		err := gconv.Struct(g.Map{
			"B": g.Map{
				"init": 0,
				"Name": "init",
			},
		}, a)
		t.AssertNil(err)
		t.Assert(a.B.Name, "init")
		t.Assert(a.B.init.Val(), true)
	})
	// It can change public attribute.
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			Init *gtype.Bool
			Name string
		}
		type A struct {
			B *B
		}
		a := &A{
			B: &B{
				Init: gtype.NewBool(),
			},
		}
		err := gconv.Struct(g.Map{
			"B": g.Map{
				"Init": 1,
				"Name": "init",
			},
		}, a)
		t.AssertNil(err)
		t.Assert(a.B.Name, "init")
		t.Assert(a.B.Init.Val(), true)
	})
}

// https://github.com/gogf/gf/issues/2381
func Test_Issue2381(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Inherit struct {
			Id        int64       `json:"id"          description:"Id"`
			Flag      *gjson.Json `json:"flag"        description:"标签"`
			Title     string      `json:"title"       description:"标题"`
			CreatedAt *gtime.Time `json:"createdAt"   description:"创建时间"`
		}
		type Test1 struct {
			Inherit
		}
		type Test2 struct {
			Inherit
		}
		var (
			a1 Test1
			a2 Test2
		)

		a1 = Test1{
			Inherit{
				Id:        2,
				Flag:      gjson.New("[1, 2]"),
				Title:     "测试",
				CreatedAt: gtime.Now(),
			},
		}
		err := gconv.Scan(a1, &a2)
		t.AssertNil(err)
		t.Assert(a1.Id, a2.Id)
		t.Assert(a1.Title, a2.Title)
		t.Assert(a1.CreatedAt, a2.CreatedAt)
		t.Assert(a1.Flag.String(), a2.Flag.String())
	})
}

// https://github.com/gogf/gf/issues/2391
func Test_Issue2391(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Inherit struct {
			Ids   []int
			Ids2  []int64
			Flag  *gjson.Json
			Title string
		}

		type Test1 struct {
			Inherit
		}
		type Test2 struct {
			Inherit
		}

		var (
			a1 Test1
			a2 Test2
		)

		a1 = Test1{
			Inherit{
				Ids:   []int{1, 2, 3},
				Ids2:  []int64{4, 5, 6},
				Flag:  gjson.New("[\"1\", \"2\"]"),
				Title: "测试",
			},
		}

		err := gconv.Scan(a1, &a2)
		t.AssertNil(err)
		t.Assert(a1.Ids, a2.Ids)
		t.Assert(a1.Ids2, a2.Ids2)
		t.Assert(a1.Title, a2.Title)
		t.Assert(a1.Flag.String(), a2.Flag.String())
	})
}

// https://github.com/gogf/gf/issues/2395
func Test_Issue2395(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Test struct {
			Num int
		}
		var ()
		obj := Test{Num: 0}
		t.Assert(gconv.Interfaces(obj), []interface{}{obj})
	})
}

// https://github.com/gogf/gf/issues/2371
func Test_Issue2371(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			s = struct {
				Time time.Time `json:"time"`
			}{}
			jsonMap = map[string]interface{}{"time": "2022-12-15 16:11:34"}
		)

		err := gconv.Struct(jsonMap, &s)
		t.AssertNil(err)
		t.Assert(s.Time.UTC(), `2022-12-15 08:11:34 +0000 UTC`)
	})
}

// https://github.com/gogf/gf/issues/2901
func Test_Issue2901(t *testing.T) {
	type GameApp2 struct {
		ForceUpdateTime *time.Time
	}
	gtest.C(t, func(t *gtest.T) {
		src := map[string]interface{}{
			"FORCE_UPDATE_TIME": time.Now(),
		}
		m := GameApp2{}
		err := gconv.Scan(src, &m)
		t.AssertNil(err)
	})
}

// https://github.com/gogf/gf/issues/3006
func Test_Issue3006(t *testing.T) {
	type tFF struct {
		Val1 json.RawMessage            `json:"val1"`
		Val2 []json.RawMessage          `json:"val2"`
		Val3 map[string]json.RawMessage `json:"val3"`
	}

	gtest.C(t, func(t *gtest.T) {
		ff := &tFF{}
		var tmp = map[string]any{
			"val1": map[string]any{"hello": "world"},
			"val2": []any{map[string]string{"hello": "world"}},
			"val3": map[string]map[string]string{"val3": {"hello": "world"}},
		}

		err := gconv.Struct(tmp, ff)
		t.AssertNil(err)
		t.AssertNE(ff, nil)
		t.Assert(ff.Val1, []byte(`{"hello":"world"}`))
		t.AssertEQ(len(ff.Val2), 1)
		t.Assert(ff.Val2[0], []byte(`{"hello":"world"}`))
		t.AssertEQ(len(ff.Val3), 1)
		t.Assert(ff.Val3["val3"], []byte(`{"hello":"world"}`))
	})
}

// https://github.com/gogf/gf/issues/3731
func Test_Issue3731(t *testing.T) {
	type Data struct {
		Doc map[string]interface{} `json:"doc"`
	}

	gtest.C(t, func(t *gtest.T) {
		dataMap := map[string]any{
			"doc": map[string]any{
				"craft": nil,
			},
		}

		var args Data
		err := gconv.Struct(dataMap, &args)
		t.AssertNil(err)
		t.AssertEQ("<nil>", fmt.Sprintf("%T", args.Doc["craft"]))
	})
}

// https://github.com/gogf/gf/issues/3764
func Test_Issue3764(t *testing.T) {
	type T struct {
		True     bool  `json:"true"`
		False    bool  `json:"false"`
		TruePtr  *bool `json:"true_ptr"`
		FalsePtr *bool `json:"false_ptr"`
	}
	gtest.C(t, func(t *gtest.T) {
		trueValue := true
		falseValue := false
		m := g.Map{
			"true":      trueValue,
			"false":     falseValue,
			"true_ptr":  &trueValue,
			"false_ptr": &falseValue,
		}
		tt := &T{}
		err := gconv.Struct(m, &tt)
		t.AssertNil(err)
		t.AssertEQ(tt.True, true)
		t.AssertEQ(tt.False, false)
		t.AssertEQ(*tt.TruePtr, trueValue)
		t.AssertEQ(*tt.FalsePtr, falseValue)
	})
}

// https://github.com/gogf/gf/issues/3789
func Test_Issue3789(t *testing.T) {
	type ItemSecondThird struct {
		SecondID uint64 `json:"secondId,string"`
		ThirdID  uint64 `json:"thirdId,string"`
	}
	type ItemFirst struct {
		ID uint64 `json:"id,string"`
		ItemSecondThird
	}
	type ItemInput struct {
		ItemFirst
	}
	type HelloReq struct {
		g.Meta `path:"/hello" method:"GET"`
		ItemInput
	}
	gtest.C(t, func(t *gtest.T) {
		m := map[string]interface{}{
			"id":       1,
			"secondId": 2,
			"thirdId":  3,
		}
		var dest HelloReq
		err := gconv.Scan(m, &dest)
		t.AssertNil(err)
		t.Assert(dest.ID, uint64(1))
		t.Assert(dest.SecondID, uint64(2))
		t.Assert(dest.ThirdID, uint64(3))
	})
}

// https://github.com/gogf/gf/issues/3797
func Test_Issue3797(t *testing.T) {
	type Option struct {
		F1 int
		F2 string
	}
	type Rule struct {
		ID   int64     `json:"id"`
		Rule []*Option `json:"rule"`
	}
	type Res1 struct {
		g.Meta
		Rule
	}
	gtest.C(t, func(t *gtest.T) {
		var r = &Rule{
			ID: 100,
		}
		var res = &Res1{}
		for i := 0; i < 10000; i++ {
			err := gconv.Scan(r, res)
			t.AssertNil(err)
			t.Assert(res.ID, 100)
			t.AssertEQ(res.Rule.Rule, nil)
		}
	})
}

// https://github.com/gogf/gf/issues/3800
func Test_Issue3800(t *testing.T) {
	// might be random assignment in converting,
	// it here so runs multiple times to reproduce the issue.
	for i := 0; i < 1000; i++ {
		doTestIssue3800(t)
	}
}

func doTestIssue3800(t *testing.T) {
	type NullID string

	type StructA struct {
		Superior    string `json:"superior"`
		UpdatedTick int    `json:"updated_tick"`
	}
	type StructB struct {
		Superior    *NullID `json:"superior"`
		UpdatedTick *int    `json:"updated_tick"`
	}

	type StructC struct {
		Superior    string `json:"superior"`
		UpdatedTick int    `json:"updated_tick"`
	}
	type StructD struct {
		StructC
		Superior    *NullID `json:"superior"`
		UpdatedTick *int    `json:"updated_tick"`
	}

	type StructE struct {
		Superior    string `json:"superior"`
		UpdatedTick int    `json:"updated_tick"`
	}
	type StructF struct {
		Superior    *NullID `json:"superior"`
		UpdatedTick *int    `json:"updated_tick"`
		StructE
	}

	type StructG struct {
		Superior    string `json:"superior"`
		UpdatedTick int    `json:"updated_tick"`
	}
	type StructH struct {
		Superior    *string `json:"superior"`
		UpdatedTick *int    `json:"updated_tick"`
		StructG
	}

	type StructI struct {
		Master struct {
			Superior    *NullID `json:"superior"`
			UpdatedTick int     `json:"updated_tick"`
		} `json:"master"`
	}
	type StructJ struct {
		StructA
		Superior    *NullID `json:"superior"`
		UpdatedTick *int    `json:"updated_tick"`
	}

	type StructK struct {
		Master struct {
			Superior    *NullID `json:"superior"`
			UpdatedTick int     `json:"updated_tick"`
		} `json:"master"`
	}
	type StructL struct {
		Superior    *NullID `json:"superior"`
		UpdatedTick *int    `json:"updated_tick"`
		StructA
	}

	// case 0
	// NullID should not be initialized.
	gtest.C(t, func(t *gtest.T) {
		structA := g.Map{
			"UpdatedTick": 10,
		}
		structB := StructB{}
		err := gconv.Scan(structA, &structB)
		t.AssertNil(err)
		t.AssertNil(structB.Superior)
		t.Assert(*structB.UpdatedTick, structA["UpdatedTick"])
	})

	// case 1
	gtest.C(t, func(t *gtest.T) {
		structA := StructA{
			Superior:    "superior100",
			UpdatedTick: 20,
		}
		structB := StructB{}
		err := gconv.Scan(structA, &structB)
		t.AssertNil(err)
		t.Assert(*structB.Superior, structA.Superior)
	})

	// case 2
	gtest.C(t, func(t *gtest.T) {
		structA1 := StructA{
			Superior:    "100",
			UpdatedTick: 20,
		}
		structB1 := StructB{}
		err := gconv.Scan(structA1, &structB1)
		t.AssertNil(err)
		t.Assert(*structB1.Superior, structA1.Superior)
		t.Assert(*structB1.UpdatedTick, structA1.UpdatedTick)
	})

	// case 3
	gtest.C(t, func(t *gtest.T) {
		structC := StructC{
			Superior:    "superior100",
			UpdatedTick: 20,
		}
		structD := StructD{}
		err := gconv.Scan(structC, &structD)
		t.AssertNil(err)
		t.Assert(structD.StructC.Superior, structC.Superior)
		t.Assert(*structD.Superior, structC.Superior)
		t.Assert(*structD.UpdatedTick, structC.UpdatedTick)
	})

	// case 4
	gtest.C(t, func(t *gtest.T) {
		structC1 := StructC{
			Superior:    "100",
			UpdatedTick: 20,
		}
		structD1 := StructD{}
		err := gconv.Scan(structC1, &structD1)
		t.AssertNil(err)
		t.Assert(structD1.StructC.Superior, structC1.Superior)
		t.Assert(structD1.StructC.UpdatedTick, structC1.UpdatedTick)
		t.Assert(*structD1.Superior, structC1.Superior)
		t.Assert(*structD1.UpdatedTick, structC1.UpdatedTick)
	})

	// case 5
	gtest.C(t, func(t *gtest.T) {
		structE := StructE{
			Superior:    "superior100",
			UpdatedTick: 20,
		}
		structF := StructF{}
		err := gconv.Scan(structE, &structF)
		t.AssertNil(err)
		t.Assert(structF.StructE.Superior, structE.Superior)
		t.Assert(structF.StructE.UpdatedTick, structE.UpdatedTick)
		t.Assert(*structF.Superior, structE.Superior)
		t.Assert(*structF.UpdatedTick, structE.UpdatedTick)
	})

	// case 6
	gtest.C(t, func(t *gtest.T) {
		structE1 := StructE{
			Superior:    "100",
			UpdatedTick: 20,
		}
		structF1 := StructF{}
		err := gconv.Scan(structE1, &structF1)
		t.AssertNil(err)
		t.Assert(*structF1.Superior, structE1.Superior)
		t.Assert(*structF1.UpdatedTick, structE1.UpdatedTick)
		t.Assert(structF1.StructE.Superior, structE1.Superior)
		t.Assert(structF1.StructE.UpdatedTick, structE1.UpdatedTick)
	})

	// case 7
	gtest.C(t, func(t *gtest.T) {
		structG := StructG{
			Superior:    "superior100",
			UpdatedTick: 20,
		}
		structH := StructH{}
		err := gconv.Scan(structG, &structH)
		t.AssertNil(err)
		t.Assert(*structH.Superior, structG.Superior)
		t.Assert(*structH.UpdatedTick, structG.UpdatedTick)
		t.Assert(structH.StructG.Superior, structG.Superior)
		t.Assert(structH.StructG.UpdatedTick, structG.UpdatedTick)
	})

	// case 8
	gtest.C(t, func(t *gtest.T) {
		structG1 := StructG{
			Superior:    "100",
			UpdatedTick: 20,
		}
		structH1 := StructH{}
		err := gconv.Scan(structG1, &structH1)
		t.AssertNil(err)
		t.Assert(*structH1.Superior, structG1.Superior)
		t.Assert(*structH1.UpdatedTick, structG1.UpdatedTick)
		t.Assert(structH1.StructG.Superior, structG1.Superior)
		t.Assert(structH1.StructG.UpdatedTick, structG1.UpdatedTick)
	})

	// case 9
	gtest.C(t, func(t *gtest.T) {
		structI := StructI{}
		xxx := NullID("superior100")
		structI.Master.Superior = &xxx
		structI.Master.UpdatedTick = 30
		structJ := StructJ{}
		err := gconv.Scan(structI.Master, &structJ)
		t.AssertNil(err)
		t.Assert(*structJ.Superior, structI.Master.Superior)
		t.Assert(*structJ.UpdatedTick, structI.Master.UpdatedTick)
		t.Assert(structJ.StructA.Superior, structI.Master.Superior)
		t.Assert(structJ.StructA.UpdatedTick, structI.Master.UpdatedTick)
	})

	// case 10
	gtest.C(t, func(t *gtest.T) {
		structK := StructK{}
		yyy := NullID("superior100")
		structK.Master.Superior = &yyy
		structK.Master.UpdatedTick = 40
		structL := StructL{}
		err := gconv.Scan(structK.Master, &structL)
		t.AssertNil(err)
		t.Assert(*structL.Superior, structK.Master.Superior)
		t.Assert(*structL.UpdatedTick, structK.Master.UpdatedTick)
		t.Assert(structL.StructA.Superior, structK.Master.Superior)
		t.Assert(structL.StructA.UpdatedTick, structK.Master.UpdatedTick)
	})
}

// https://github.com/gogf/gf/issues/3821
func Test_Issue3821(t *testing.T) {
	// Scan
	gtest.C(t, func(t *gtest.T) {
		var record = map[string]interface{}{
			`user_id`:   1,
			`user_name`: "teemo",
		}

		type DoubleInnerUser struct {
			UserId int64 `orm:"user_id"`
		}

		type InnerUser struct {
			UserId     int32   `orm:"user_id"`
			UserIdBool bool    `orm:"user_id"`
			Username   *string `orm:"user_name"`
			Username2  *string `orm:"user_name"`
			Username3  string  `orm:"username"`
			*DoubleInnerUser
		}

		type User struct {
			InnerUser
			UserId     int        `orm:"user_id"`
			UserIdBool gtype.Bool `orm:"user_id"`
			Username   string     `orm:"user_name"`
			Username2  string     `orm:"user_name"`
			Username3  *string    `orm:"user_name"`
			Username4  string     `orm:"username"` // empty string
		}
		var user = &User{}
		err := gconv.StructTag(record, user, "orm")

		t.AssertNil(err)
		t.AssertEQ(user.UserId, 1)
		t.AssertEQ(user.UserIdBool.Val(), true)
		t.AssertEQ(user.Username, "teemo")
		t.AssertEQ(user.Username2, "teemo")
		t.AssertEQ(*user.Username3, "teemo")
		t.AssertEQ(user.Username4, "")
		t.AssertEQ(user.InnerUser.UserId, int32(1))
		t.AssertEQ(user.InnerUser.UserIdBool, true)
		t.AssertEQ(*user.InnerUser.Username, "teemo")
		t.AssertEQ(*user.InnerUser.Username2, "teemo")
		t.AssertEQ(user.InnerUser.Username3, "")
		t.AssertEQ(user.DoubleInnerUser.UserId, int64(1))
	})
}

// https://github.com/gogf/gf/issues/3868
func Test_Issue3868(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Config struct {
			Enable   bool
			Spec     string
			PoolSize int
		}
		data := gjson.New(`[{"enable":false,"spec":"a"},{"enable":true,"poolSize":1}]`)
		for i := 0; i < 1000; i++ {
			var configs []*Config
			err := gconv.Structs(data, &configs)
			t.AssertNil(err)
			t.Assert(len(configs), 2)
			t.Assert(configs[0], &Config{
				Enable: false,
				Spec:   "a",
			})
			t.Assert(configs[1], &Config{
				Enable:   true,
				PoolSize: 1,
			})
		}
	})
}

// https://github.com/gogf/gf/issues/3903
func Test_Issue3903(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type TestA struct {
			UserId int `json:"UserId"   orm:"user_id"   `
		}
		type TestB struct {
			TestA
			UserId int `json:"NewUserId"  description:""`
		}
		var input = map[string]interface{}{
			"user_id": gvar.New(100, true),
		}
		var a TestB
		err := gconv.StructTag(input, &a, "orm")
		t.AssertNil(err)
		t.Assert(a.TestA.UserId, 100)
		t.Assert(a.UserId, 100)
	})
}
