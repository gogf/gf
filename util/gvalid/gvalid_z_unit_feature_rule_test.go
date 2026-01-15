// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gtag"
)

var (
	ctx = gctx.New()
)

func Test_Check(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		rule := "abc:6,16"
		val1 := 0
		val2 := 7
		val3 := 20
		err1 := g.Validator().Data(val1).Rules(rule).Run(ctx)
		err2 := g.Validator().Data(val2).Rules(rule).Run(ctx)
		err3 := g.Validator().Data(val3).Rules(rule).Run(ctx)
		t.Assert(err1, "InvalidRules: abc:6,16")
		t.Assert(err2, "InvalidRules: abc:6,16")
		t.Assert(err3, "InvalidRules: abc:6,16")
	})
}

func Test_Array(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := g.Validator().Data("1").Rules("array").Run(ctx)
		t.Assert(err, "The value `1` is not of valid array type")
	})
	gtest.C(t, func(t *gtest.T) {
		err := g.Validator().Data("").Rules("array").Run(ctx)
		t.Assert(err, "The value `` is not of valid array type")
	})
	gtest.C(t, func(t *gtest.T) {
		err := g.Validator().Data("[1,2,3]").Rules("array").Run(ctx)
		t.Assert(err, "")
	})
	gtest.C(t, func(t *gtest.T) {
		err := g.Validator().Data("[]").Rules("array").Run(ctx)
		t.Assert(err, "")
	})
	gtest.C(t, func(t *gtest.T) {
		err := g.Validator().Data([]int{1, 2, 3}).Rules("array").Run(ctx)
		t.Assert(err, "")
	})
	gtest.C(t, func(t *gtest.T) {
		err := g.Validator().Data([]int{}).Rules("array").Run(ctx)
		t.Assert(err, "")
	})
}

func Test_Required(t *testing.T) {
	if m := g.Validator().Data("1").Rules("required").Run(ctx); m != nil {
		t.Error(m)
	}
	if m := g.Validator().Data("").Rules("required").Run(ctx); m == nil {
		t.Error(m)
	}
	if m := g.Validator().Data("").Assoc(map[string]any{"id": 1, "age": 19}).Rules("required-if: id,1,age,18").Run(ctx); m == nil {
		t.Error("Required校验失败")
	}
	if m := g.Validator().Data("").Assoc(map[string]any{"id": 2, "age": 19}).Rules("required-if: id,1,age,18").Run(ctx); m != nil {
		t.Error("Required校验失败")
	}
}

func Test_RequiredIf(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		rule := "required-if:id,1,age,18"
		t.AssertNE(g.Validator().Data("").Assoc(g.Map{"id": 1}).Rules(rule).Run(ctx), nil)
		t.Assert(g.Validator().Data("").Assoc(g.Map{"id": 0}).Rules(rule).Run(ctx), nil)
		t.AssertNE(g.Validator().Data("").Assoc(g.Map{"age": 18}).Rules(rule).Run(ctx), nil)
		t.Assert(g.Validator().Data("").Assoc(g.Map{"age": 20}).Rules(rule).Run(ctx), nil)
	})
}

func Test_RequiredIfAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		rule := "required-if-all:id,1,age,18"
		t.Assert(g.Validator().Data("").Assoc(g.Map{"id": 1}).Rules(rule).Run(ctx), nil)
		t.Assert(g.Validator().Data("").Assoc(g.Map{"age": 18}).Rules(rule).Run(ctx), nil)
		t.Assert(g.Validator().Data("").Assoc(g.Map{"id": 0, "age": 20}).Rules(rule).Run(ctx), nil)
		t.AssertNE(g.Validator().Data("").Assoc(g.Map{"id": 1, "age": 18}).Rules(rule).Run(ctx), nil)
	})
}

func Test_RequiredUnless(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		rule := "required-unless:id,1,age,18"
		t.Assert(g.Validator().Data("").Assoc(g.Map{"id": 1}).Rules(rule).Run(ctx), nil)
		t.AssertNE(g.Validator().Data("").Assoc(g.Map{"id": 0}).Rules(rule).Run(ctx), nil)
		t.Assert(g.Validator().Data("").Assoc(g.Map{"age": 18}).Rules(rule).Run(ctx), nil)
		t.AssertNE(g.Validator().Data("").Assoc(g.Map{"age": 20}).Rules(rule).Run(ctx), nil)
	})
}

func Test_RequiredWith(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		rule := "required-with:id,name"
		val1 := ""
		params1 := g.Map{
			"age": 18,
		}
		params2 := g.Map{
			"id": 100,
		}
		params3 := g.Map{
			"id":   100,
			"name": "john",
		}
		err1 := g.Validator().Data(val1).Assoc(params1).Rules(rule).Run(ctx)
		err2 := g.Validator().Data(val1).Assoc(params2).Rules(rule).Run(ctx)
		err3 := g.Validator().Data(val1).Assoc(params3).Rules(rule).Run(ctx)
		t.Assert(err1, nil)
		t.AssertNE(err2, nil)
		t.AssertNE(err3, nil)
	})
	// time.Time
	gtest.C(t, func(t *gtest.T) {
		rule := "required-with:id,time"
		val1 := ""
		params1 := g.Map{
			"age": 18,
		}
		params2 := g.Map{
			"id": 100,
		}
		params3 := g.Map{
			"time": time.Time{},
		}
		err1 := g.Validator().Data(val1).Assoc(params1).Rules(rule).Run(ctx)
		err2 := g.Validator().Data(val1).Assoc(params2).Rules(rule).Run(ctx)
		err3 := g.Validator().Data(val1).Assoc(params3).Rules(rule).Run(ctx)
		t.Assert(err1, nil)
		t.AssertNE(err2, nil)
		t.Assert(err3, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		rule := "required-with:id,time"
		val1 := ""
		params1 := g.Map{
			"age": 18,
		}
		params2 := g.Map{
			"id": 100,
		}
		params3 := g.Map{
			"time": time.Now(),
		}
		err1 := g.Validator().Data(val1).Assoc(params1).Rules(rule).Run(ctx)
		err2 := g.Validator().Data(val1).Assoc(params2).Rules(rule).Run(ctx)
		err3 := g.Validator().Data(val1).Assoc(params3).Rules(rule).Run(ctx)
		t.Assert(err1, nil)
		t.AssertNE(err2, nil)
		t.AssertNE(err3, nil)
	})
	// gtime.Time
	gtest.C(t, func(t *gtest.T) {
		type UserApiSearch struct {
			Uid       int64       `json:"uid"`
			Nickname  string      `json:"nickname" v:"required-with:Uid"`
			StartTime *gtime.Time `json:"start_time" v:"required-with:EndTime"`
			EndTime   *gtime.Time `json:"end_time" v:"required-with:StartTime"`
		}
		data := UserApiSearch{
			StartTime: nil,
			EndTime:   nil,
		}
		t.Assert(g.Validator().Data(data).Run(ctx), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		type UserApiSearch struct {
			Uid       int64       `json:"uid"`
			Nickname  string      `json:"nickname" v:"required-with:Uid"`
			StartTime *gtime.Time `json:"start_time" v:"required-with:EndTime"`
			EndTime   *gtime.Time `json:"end_time" v:"required-with:StartTime"`
		}
		data := UserApiSearch{
			StartTime: nil,
			EndTime:   gtime.Now(),
		}
		t.AssertNE(g.Validator().Data(data).Run(ctx), nil)
	})
}

func Test_RequiredWithAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		rule := "required-with-all:id,name"
		val1 := ""
		params1 := g.Map{
			"age": 18,
		}
		params2 := g.Map{
			"id": 100,
		}
		params3 := g.Map{
			"id":   100,
			"name": "john",
		}
		err1 := g.Validator().Data(val1).Assoc(params1).Rules(rule).Run(ctx)
		err2 := g.Validator().Data(val1).Assoc(params2).Rules(rule).Run(ctx)
		err3 := g.Validator().Data(val1).Assoc(params3).Rules(rule).Run(ctx)
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.AssertNE(err3, nil)
	})
}

func Test_RequiredWithOut(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		rule := "required-without:id,name"
		val1 := ""
		params1 := g.Map{
			"age": 18,
		}
		params2 := g.Map{
			"id": 100,
		}
		params3 := g.Map{
			"id":   100,
			"name": "john",
		}
		err1 := g.Validator().Data(val1).Assoc(params1).Rules(rule).Run(ctx)
		err2 := g.Validator().Data(val1).Assoc(params2).Rules(rule).Run(ctx)
		err3 := g.Validator().Data(val1).Assoc(params3).Rules(rule).Run(ctx)
		t.AssertNE(err1, nil)
		t.AssertNE(err2, nil)
		t.Assert(err3, nil)
	})
}

func Test_RequiredWithOutAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		rule := "required-without-all:id,name"
		val1 := ""
		params1 := g.Map{
			"age": 18,
		}
		params2 := g.Map{
			"id": 100,
		}
		params3 := g.Map{
			"id":   100,
			"name": "john",
		}
		err1 := g.Validator().Data(val1).Assoc(params1).Rules(rule).Run(ctx)
		err2 := g.Validator().Data(val1).Assoc(params2).Rules(rule).Run(ctx)
		err3 := g.Validator().Data(val1).Assoc(params3).Rules(rule).Run(ctx)
		t.AssertNE(err1, nil)
		t.Assert(err2, nil)
		t.Assert(err3, nil)
	})
}

func Test_Date(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"2010":       false,
			"201011":     false,
			"20101101":   true,
			"2010-11-01": true,
			"2010.11.01": true,
			"2010/11/01": true,
			"2010=11=01": false,
			"123":        false,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("date").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_Datetime(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"2010":                false,
			"2010.11":             false,
			"2010-11-01":          false,
			"2010-11-01 12:00":    false,
			"2010-11-01 12:00:00": true,
			"2010.11.01 12:00:00": false,
		}
		for k, v := range m {
			err := g.Validator().Rules(`datetime`).Data(k).Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_DateFormat(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrStr{
			"2010":                 "date-format:Y",
			"201011":               "date-format:Ym",
			"2010.11":              "date-format:Y.m",
			"201011-01":            "date-format:Ym-d",
			"2010~11~01":           "date-format:Y~m~d",
			"2010-11~01":           "date-format:Y-m~d",
			"2023-09-10T19:46:31Z": "date-format:2006-01-02\\T15:04:05Z07:00", // RFC3339
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules(v).Run(ctx)
			t.AssertNil(err)
		}
	})
	gtest.C(t, func(t *gtest.T) {
		errM := g.MapStrStr{
			"2010-11~01": "date-format:Y~m~d",
		}
		for k, v := range errM {
			err := g.Validator().Data(k).Rules(v).Run(ctx)
			t.AssertNE(err, nil)
		}
	})
	gtest.C(t, func(t *gtest.T) {
		t1 := gtime.Now()
		t2 := time.Time{}
		err1 := g.Validator().Data(t1).Rules("date-format:Y").Run(ctx)
		err2 := g.Validator().Data(t2).Rules("date-format:Y").Run(ctx)
		t.Assert(err1, nil)
		t.AssertNE(err2, nil)
	})
}

func Test_Email(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"m@johngcn":           false,
			"m@www@johngcn":       false,
			"m-m_m@mail.johng.cn": true,
			"m.m-m@johng.cn":      true,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("email").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_Phone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"1361990897":  false,
			"13619908979": true,
			"16719908979": true,
			"19719908989": true,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("phone").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_PhoneLoose(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"13333333333": true,
			"15555555555": true,
			"16666666666": true,
			"23333333333": false,
			"1333333333":  false,
			"10333333333": false,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("phone-loose").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_Telephone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"869265":       false,
			"028-869265":   false,
			"86292651":     true,
			"028-8692651":  true,
			"0830-8692651": true,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("telephone").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_Passport(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"123456":   false,
			"a12345-6": false,
			"aaaaa":    false,
			"aaaaaa":   true,
			"a123_456": true,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("passport").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_Password(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"12345":     false,
			"aaaaa":     false,
			"a12345-6":  true,
			">,/;'[09-": true,
			"a123_456":  true,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("password").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_Password2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"12345":     false,
			"Naaaa":     false,
			"a12345-6":  false,
			">,/;'[09-": false,
			"a123_456":  false,
			"Nant1986":  true,
			"Nant1986!": true,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("password2").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_Password3(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"12345":     false,
			"Naaaa":     false,
			"a12345-6":  false,
			">,/;'[09-": false,
			"a123_456":  false,
			"Nant1986":  false,
			"Nant1986!": true,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("password3").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_Postcode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"12345":  false,
			"610036": true,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("postcode").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_ResidentId(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"11111111111111":     false,
			"1111111111111111":   false,
			"311128500121201":    false,
			"510521198607185367": false,
			"51052119860718536x": true,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("resident-id").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_BankCard(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"6230514630000424470": false,
			"6230514630000424473": true,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("bank-card").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_QQ(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"100":       false,
			"1":         false,
			"10000":     true,
			"38996181":  true,
			"389961817": true,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("qq").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_Ip(t *testing.T) {
	if m := g.Validator().Data("10.0.0.1").Rules("ip").Run(ctx); m != nil {
		t.Error(m)
	}
	if m := g.Validator().Data("10.0.0.1").Rules("ipv4").Run(ctx); m != nil {
		t.Error(m)
	}
	if m := g.Validator().Data("0.0.0.0").Rules("ipv4").Run(ctx); m != nil {
		t.Error(m)
	}
	if m := g.Validator().Data("1920.0.0.0").Rules("ipv4").Run(ctx); m == nil {
		t.Error("ipv4校验失败")
	}
	if m := g.Validator().Data("1920.0.0.0").Rules("ip").Run(ctx); m == nil {
		t.Error("ipv4校验失败")
	}
	if m := g.Validator().Data("fe80::5484:7aff:fefe:9799").Rules("ipv6").Run(ctx); m != nil {
		t.Error(m)
	}
	if m := g.Validator().Data("fe80::5484:7aff:fefe:9799123").Rules("ipv6").Run(ctx); m == nil {
		t.Error(m)
	}
	if m := g.Validator().Data("fe80::5484:7aff:fefe:9799").Rules("ip").Run(ctx); m != nil {
		t.Error(m)
	}
	if m := g.Validator().Data("fe80::5484:7aff:fefe:9799123").Rules("ip").Run(ctx); m == nil {
		t.Error(m)
	}
}

func Test_IPv4(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"0.0.0":         false,
			"0.0.0.0":       true,
			"1.1.1.1":       true,
			"255.255.255.0": true,
			"127.0.0.1":     true,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("ipv4").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_IPv6(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"192.168.1.1": false,
			"CDCD:910A:2222:5498:8475:1111:3900:2020": true,
			"1030::C9B4:FF12:48AA:1A2B":               true,
			"2000:0:0:0:0:0:0:1":                      true,
			"0000:0000:0000:0000:0000:ffff:c0a8:5909": true,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("ipv6").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_MAC(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"192.168.1.1":       false,
			"44-45-53-54-00-00": true,
			"01:00:5e:00:00:00": true,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("mac").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_URL(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"127.0.0.1":             false,
			"https://www.baidu.com": true,
			"http://127.0.0.1":      true,
			"file:///tmp/test.txt":  true,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("url").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_Domain(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"localhost":     false,
			"baidu.com":     true,
			"www.baidu.com": true,
			"jn.np":         true,
			"www.jn.np":     true,
			"w.www.jn.np":   true,
			"127.0.0.1":     false,
			"www.360.com":   true,
			"www.360":       false,
			"360":           false,
			"my-gf":         false,
			"my-gf.com":     true,
			"my-gf.360.com": true,
		}
		var err error
		for k, v := range m {
			err = g.Validator().Data(k).Rules("domain").Run(ctx)
			if v {
				// fmt.Println(k)
				t.AssertNil(err)
			} else {
				// fmt.Println(k)
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_Length(t *testing.T) {
	rule := "length:6,16"
	if m := g.Validator().Data("123456").Rules(rule).Run(ctx); m != nil {
		t.Error(m)
	}
	if m := g.Validator().Data("12345").Rules(rule).Run(ctx); m == nil {
		t.Error("长度校验失败")
	}
}

func Test_MinLength(t *testing.T) {
	rule := "min-length:6"
	msgs := map[string]string{
		"min-length": "地址长度至少为{min}位",
	}
	if m := g.Validator().Data("123456").Rules(rule).Run(ctx); m != nil {
		t.Error(m)
	}
	if m := g.Validator().Data("12345").Rules(rule).Run(ctx); m == nil {
		t.Error("长度校验失败")
	}
	if m := g.Validator().Data("12345").Rules(rule).Messages(msgs).Run(ctx); m == nil {
		t.Error("长度校验失败")
	}

	rule2 := "min-length:abc"
	if m := g.Validator().Data("123456").Rules(rule2).Run(ctx); m == nil {
		t.Error("长度校验失败")
	}
}

func Test_MaxLength(t *testing.T) {
	rule := "max-length:6"
	msgs := map[string]string{
		"max-length": "地址长度至大为{max}位",
	}
	if m := g.Validator().Data("12345").Rules(rule).Run(ctx); m != nil {
		t.Error(m)
	}
	if m := g.Validator().Data("1234567").Rules(rule).Run(ctx); m == nil {
		t.Error("长度校验失败")
	}
	if m := g.Validator().Data("1234567").Rules(rule).Messages(msgs).Run(ctx); m == nil {
		t.Error("长度校验失败")
	}

	rule2 := "max-length:abc"
	if m := g.Validator().Data("123456").Rules(rule2).Run(ctx); m == nil {
		t.Error("长度校验失败")
	}
}

func Test_Size(t *testing.T) {
	rule := "size:5"
	if m := g.Validator().Data("12345").Rules(rule).Run(ctx); m != nil {
		t.Error(m)
	}
	if m := g.Validator().Data("123456").Rules(rule).Run(ctx); m == nil {
		t.Error("长度校验失败")
	}
}

func Test_Between(t *testing.T) {
	rule := "between:6.01, 10.01"
	if m := g.Validator().Data(10).Rules(rule).Run(ctx); m != nil {
		t.Error(m)
	}
	if m := g.Validator().Data(10.02).Rules(rule).Run(ctx); m == nil {
		t.Error("大小范围校验失败")
	}
	if m := g.Validator().Data("a").Rules(rule).Run(ctx); m == nil {
		t.Error("大小范围校验失败")
	}
}

func Test_Min(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"1":    false,
			"99":   false,
			"100":  true,
			"1000": true,
			"a":    false,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("min:100").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
	gtest.C(t, func(t *gtest.T) {
		err := g.Validator().Data("1").Rules("min:a").Run(ctx)
		t.AssertNE(err, nil)
	})
}

func Test_Max(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"1":    true,
			"99":   true,
			"100":  true,
			"1000": false,
			"a":    false,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("max:100").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
	gtest.C(t, func(t *gtest.T) {
		err := g.Validator().Data("1").Rules("max:a").Run(ctx)
		t.AssertNE(err, nil)
	})
}

func Test_Json(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"":                   false,
			".":                  false,
			"{}":                 true,
			"[]":                 true,
			"[1,2,3,4]":          true,
			`{"list":[1,2,3,4]}`: true,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("json").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_Integer(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"":          false,
			"1.0":       false,
			"001":       true,
			"1":         true,
			"100":       true,
			"999999999": true,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("integer").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_Float(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"":    false,
			"a":   false,
			"1":   true,
			"1.0": true,
			"1.1": true,
			"0.1": true,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("float").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_Boolean(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"a":    false,
			"-":    false,
			"":     true,
			"1":    true,
			"true": true,
			"off":  true,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("boolean").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_Same(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type testCase struct {
			params g.Map
			pass   bool
		}
		cases := []testCase{
			{g.Map{"age": 18}, false},
			{g.Map{"id": 100}, true},
			{g.Map{"id": 100, "name": "john"}, true},
		}
		for _, c := range cases {
			err := g.Validator().Data("100").Assoc(c.params).Rules("same:id").Run(ctx)
			if c.pass {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_Different(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type testCase struct {
			params g.Map
			pass   bool
		}
		cases := []testCase{
			{g.Map{"age": 18}, true},
			{g.Map{"id": 100}, false},
			{g.Map{"id": 100, "name": "john"}, false},
		}
		for _, c := range cases {
			err := g.Validator().Data("100").Assoc(c.params).Rules("different:id").Run(ctx)
			if c.pass {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_EQ(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type testCase struct {
			params g.Map
			pass   bool
		}
		cases := []testCase{
			{g.Map{"age": 18}, false},
			{g.Map{"id": 100}, true},
			{g.Map{"id": 100, "name": "john"}, true},
		}
		for _, c := range cases {
			err := g.Validator().Data("100").Assoc(c.params).Rules("eq:id").Run(ctx)
			if c.pass {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_Not_EQ(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type testCase struct {
			params g.Map
			pass   bool
		}
		cases := []testCase{
			{g.Map{"age": 18}, true},
			{g.Map{"id": 100}, false},
			{g.Map{"id": 100, "name": "john"}, false},
		}
		for _, c := range cases {
			err := g.Validator().Data("100").Assoc(c.params).Rules("not-eq:id").Run(ctx)
			if c.pass {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_In(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"":    false,
			"1":   false,
			"100": true,
			"200": true,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("in:100,200").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_NotIn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"":    true,
			"1":   true,
			"100": false,
			"200": true,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("not-in:100").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"":    true,
			"1":   true,
			"100": false,
			"200": false,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("not-in:100,200").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_Regex1(t *testing.T) {
	rule := `regex:\d{6}|\D{6}|length:6,16`
	if m := g.Validator().Data("123456").Rules(rule).Run(ctx); m != nil {
		t.Error(m)
	}
	if m := g.Validator().Data("abcde6").Rules(rule).Run(ctx); m == nil {
		t.Error("校验失败")
	}
}

func Test_Regex2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		rule := `required|min-length:6|regex:^data:image\/(jpeg|png);base64,`
		str1 := ""
		str2 := "data"
		str3 := "data:image/jpeg;base64,/9jrbattq22r"
		err1 := g.Validator().Data(str1).Rules(rule).Run(ctx)
		err2 := g.Validator().Data(str2).Rules(rule).Run(ctx)
		err3 := g.Validator().Data(str3).Rules(rule).Run(ctx)
		t.AssertNE(err1, nil)
		t.AssertNE(err2, nil)
		t.Assert(err3, nil)

		t.AssertNE(err1.Map()["required"], nil)
		t.AssertNE(err2.Map()["min-length"], nil)
	})
}

func Test_Not_Regex(t *testing.T) {
	rule := `not-regex:\d{6}|\D{6}|length:6,16`
	gtest.C(t, func(t *gtest.T) {
		err := g.Validator().Data("123456").Rules(rule).Run(ctx)
		t.Assert(err, "The value `123456` should not be in regex of: \\d{6}|\\D{6}")
	})
	gtest.C(t, func(t *gtest.T) {
		err := g.Validator().Data("abcde6").Rules(rule).Run(ctx)
		t.AssertNil(err)
	})
}

// issue: https://github.com/gogf/gf/issues/1077
func Test_InternalError_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type a struct {
			Name string `v:"hh"`
		}
		aa := a{Name: "2"}
		err := g.Validator().Data(&aa).Run(ctx)

		t.Assert(err.String(), "InvalidRules: hh")
		t.Assert(err.Strings(), g.Slice{"InvalidRules: hh"})
		t.Assert(err.FirstError(), "InvalidRules: hh")
		t.Assert(gerror.Current(err), "InvalidRules: hh")
	})
}

func Test_Code(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := g.Validator().Rules("required").Data("").Run(ctx)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeValidationFailed)
	})

	gtest.C(t, func(t *gtest.T) {
		err := g.Validator().Rules("none-exist-rule").Data("").Run(ctx)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeInternalError)
	})
}

func Test_Bail(t *testing.T) {
	// check value with no bail
	gtest.C(t, func(t *gtest.T) {
		err := g.Validator().
			Rules("required|min:1|between:1,100").
			Messages("|min number is 1|size is between 1 and 100").
			Data(-1).Run(ctx)
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "min number is 1; size is between 1 and 100")
	})

	// check value with bail
	gtest.C(t, func(t *gtest.T) {
		err := g.Validator().
			Rules("bail|required|min:1|between:1,100").
			Messages("||min number is 1|size is between 1 and 100").
			Data(-1).Run(ctx)
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "min number is 1")
	})

	// struct with no bail
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			Page int `v:"required|min:1"`
			Size int `v:"required|min:1|between:1,100 # |min number is 1|size is between 1 and 100"`
		}
		obj := &Params{
			Page: 1,
			Size: -1,
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "min number is 1; size is between 1 and 100")
	})
	// struct with bail
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			Page int `v:"required|min:1"`
			Size int `v:"bail|required|min:1|between:1,100 # ||min number is 1|size is between 1 and 100"`
		}
		obj := &Params{
			Page: 1,
			Size: -1,
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "min number is 1")
	})
}

func Test_After(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			T1 string `v:"after:T2"`
			T2 string
		}
		obj := &Params{
			T1: "2022-09-02",
			T2: "2022-09-01",
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			T1 string `v:"after:T2"`
			T2 string
		}
		obj := &Params{
			T1: "2022-09-01",
			T2: "2022-09-02",
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.Assert(err, "The T1 value `2022-09-01` must be after field T2 value `2022-09-02`")
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			T1 *gtime.Time `v:"after:T2"`
			T2 *gtime.Time
		}
		obj := &Params{
			T1: gtime.New("2022-09-02"),
			T2: gtime.New("2022-09-01"),
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			T1 *gtime.Time `v:"after:T2"`
			T2 *gtime.Time
		}
		obj := &Params{
			T1: gtime.New("2022-09-01"),
			T2: gtime.New("2022-09-02"),
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.Assert(err, "The T1 value `2022-09-01 00:00:00` must be after field T2 value `2022-09-02 00:00:00`")
	})
}

func Test_After_Equal(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			T1 string `v:"after-equal:T2"`
			T2 string
		}
		obj := &Params{
			T1: "2022-09-02",
			T2: "2022-09-01",
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			T1 string `v:"after-equal:T2"`
			T2 string
		}
		obj := &Params{
			T1: "2022-09-01",
			T2: "2022-09-02",
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.Assert(err, "The T1 value `2022-09-01` must be after or equal to field T2 value `2022-09-02`")
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			T1 *gtime.Time `v:"after-equal:T2"`
			T2 *gtime.Time
		}
		obj := &Params{
			T1: gtime.New("2022-09-02"),
			T2: gtime.New("2022-09-01"),
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			T1 *gtime.Time `v:"after-equal:T2"`
			T2 *gtime.Time
		}
		obj := &Params{
			T1: gtime.New("2022-09-01"),
			T2: gtime.New("2022-09-01"),
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			T1 *gtime.Time `v:"after-equal:T2"`
			T2 *gtime.Time
		}
		obj := &Params{
			T1: gtime.New("2022-09-01"),
			T2: gtime.New("2022-09-02"),
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.Assert(err, "The T1 value `2022-09-01 00:00:00` must be after or equal to field T2 value `2022-09-02 00:00:00`")
	})
}

func Test_Before(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			T1 string `v:"before:T2"`
			T2 string
		}
		obj := &Params{
			T1: "2022-09-01",
			T2: "2022-09-02",
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			T1 string `v:"before:T2"`
			T2 string
		}
		obj := &Params{
			T1: "2022-09-02",
			T2: "2022-09-01",
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.Assert(err, "The T1 value `2022-09-02` must be before field T2 value `2022-09-01`")
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			T1 *gtime.Time `v:"before:T2"`
			T2 *gtime.Time
		}
		obj := &Params{
			T1: gtime.New("2022-09-01"),
			T2: gtime.New("2022-09-02"),
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			T1 *gtime.Time `v:"before:T2"`
			T2 *gtime.Time
		}
		obj := &Params{
			T1: gtime.New("2022-09-02"),
			T2: gtime.New("2022-09-01"),
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.Assert(err, "The T1 value `2022-09-02 00:00:00` must be before field T2 value `2022-09-01 00:00:00`")
	})
}

func Test_Before_Equal(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			T1 string `v:"before-equal:T2"`
			T2 string
		}
		obj := &Params{
			T1: "2022-09-01",
			T2: "2022-09-02",
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			T1 string `v:"before-equal:T2"`
			T2 string
		}
		obj := &Params{
			T1: "2022-09-02",
			T2: "2022-09-01",
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.Assert(err, "The T1 value `2022-09-02` must be before or equal to field T2")
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			T1 *gtime.Time `v:"before-equal:T2"`
			T2 *gtime.Time
		}
		obj := &Params{
			T1: gtime.New("2022-09-01"),
			T2: gtime.New("2022-09-02"),
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			T1 *gtime.Time `v:"before-equal:T2"`
			T2 *gtime.Time
		}
		obj := &Params{
			T1: gtime.New("2022-09-01"),
			T2: gtime.New("2022-09-01"),
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			T1 *gtime.Time `v:"before-equal:T2"`
			T2 *gtime.Time
		}
		obj := &Params{
			T1: gtime.New("2022-09-02"),
			T2: gtime.New("2022-09-01"),
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.Assert(err, "The T1 value `2022-09-02 00:00:00` must be before or equal to field T2")
	})
}

func Test_GT(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			V1 string `v:"gt:V2"`
			V2 string
		}
		obj := &Params{
			V1: "1.2",
			V2: "1.1",
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			V1 string `v:"gt:V2"`
			V2 string
		}
		obj := &Params{
			V1: "1.1",
			V2: "1.2",
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.Assert(err, "The V1 value `1.1` must be greater than field V2 value `1.2`")
	})
}

func Test_GTE(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			V1 string `v:"gte:V2"`
			V2 string
		}
		obj := &Params{
			V1: "1.2",
			V2: "1.1",
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			V1 string `v:"gte:V2"`
			V2 string
		}
		obj := &Params{
			V1: "1.1",
			V2: "1.2",
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.Assert(err, "The V1 value `1.1` must be greater than or equal to field V2 value `1.2`")
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			V1 string `v:"gte:V2"`
			V2 string
		}
		obj := &Params{
			V1: "1.1",
			V2: "1.1",
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.AssertNil(err)
	})
}

func Test_LT(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			V1 string `v:"lt:V2"`
			V2 string
		}
		obj := &Params{
			V1: "1.1",
			V2: "1.2",
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			V1 string `v:"lt:V2"`
			V2 string
		}
		obj := &Params{
			V1: "1.2",
			V2: "1.1",
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.Assert(err, "The V1 value `1.2` must be lesser than field V2 value `1.1`")
	})
}

func Test_LTE(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			V1 string `v:"lte:V2"`
			V2 string
		}
		obj := &Params{
			V1: "1.1",
			V2: "1.2",
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			V1 string `v:"lte:V2"`
			V2 string
		}
		obj := &Params{
			V1: "1.2",
			V2: "1.1",
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.Assert(err, "The V1 value `1.2` must be lesser than or equal to field V2 value `1.1`")
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			V1 string `v:"lte:V2"`
			V2 string
		}
		obj := &Params{
			V1: "1.1",
			V2: "1.1",
		}
		err := g.Validator().Data(obj).Run(ctx)
		t.AssertNil(err)
	})
}

func Test_Enums(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type EnumsTest string
		const (
			EnumsTestA EnumsTest = "a"
			EnumsTestB EnumsTest = "b"
		)
		type Params struct {
			Id    int
			Enums EnumsTest `v:"enums"`
		}
		type PointerParams struct {
			Id    int
			Enums *EnumsTest `v:"enums"`
		}
		type SliceParams struct {
			Id    int
			Enums []EnumsTest `v:"foreach|enums"`
		}

		oldEnumsJson, err := gtag.GetGlobalEnums()
		t.AssertNil(err)
		defer t.AssertNil(gtag.SetGlobalEnums(oldEnumsJson))

		err = gtag.SetGlobalEnums(`{"github.com/gogf/gf/v2/util/gvalid_test.EnumsTest": ["a","b"]}`)
		t.AssertNil(err)

		err = g.Validator().Data(&Params{
			Id:    1,
			Enums: EnumsTestB,
		}).Run(ctx)
		t.AssertNil(err)

		err = g.Validator().Data(&Params{
			Id:    1,
			Enums: "c",
		}).Run(ctx)
		t.Assert(err, "The Enums value `c` should be in enums of: [\"a\",\"b\"]")

		var b EnumsTest = "b"
		err = g.Validator().Data(&PointerParams{
			Id:    1,
			Enums: &b,
		}).Run(ctx)
		t.AssertNil(err)

		err = g.Validator().Data(&SliceParams{
			Id:    1,
			Enums: []EnumsTest{EnumsTestA, EnumsTestB},
		}).Run(ctx)
		t.AssertNil(err)
	})
}
func Test_Alpha(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"abc":     true,
			"ABC":     true,
			"abcABC":  true,
			"abc123":  false,
			"abc-123": false,
			"abc_123": false,
			"123":     false,
			"":        false,
			"abc def": false,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("alpha").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_AlphaDash(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"abc":          true,
			"ABC":          true,
			"abc123":       true,
			"abc-123":      true,
			"abc_123":      true,
			"abc-_123":     true,
			"abc-_ABC-123": true,
			"abc 123":      false,
			"abc@123":      false,
			"":             false,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("alpha-dash").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_AlphaNum(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"abc":       true,
			"ABC":       true,
			"123":       true,
			"abc123":    true,
			"ABC123":    true,
			"abcABC123": true,
			"abc-123":   false,
			"abc_123":   false,
			"abc 123":   false,
			"":          false,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("alpha-num").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_Lowercase(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"abc":     true,
			"abcdef":  true,
			"ABC":     false,
			"Abc":     false,
			"aBc":     false,
			"abc123":  false,
			"abc-def": false,
			"abc_def": false,
			"":        false,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("lowercase").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_Numeric(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"0":         true,
			"123":       true,
			"0123":      true,
			"123456789": true,
			"1.23":      false,
			"abc":       false,
			"123abc":    false,
			"abc123":    false,
			"":          false,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("numeric").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}

func Test_Uppercase(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.MapStrBool{
			"ABC":     true,
			"ABCDEF":  true,
			"abc":     false,
			"Abc":     false,
			"AbC":     false,
			"ABC123":  false,
			"ABC-DEF": false,
			"ABC_DEF": false,
			"":        false,
		}
		for k, v := range m {
			err := g.Validator().Data(k).Rules("uppercase").Run(ctx)
			if v {
				t.AssertNil(err)
			} else {
				t.AssertNE(err, nil)
			}
		}
	})
}
