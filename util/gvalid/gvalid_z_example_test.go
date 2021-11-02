// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gvalid"
	"math"
	"reflect"
)

func ExampleCheckMap() {
	params := map[string]interface{}{
		"passport":  "",
		"password":  "123456",
		"password2": "1234567",
	}
	rules := []string{
		"passport@required|length:6,16#账号不能为空|账号长度应当在:min到:max之间",
		"password@required|length:6,16|same:password2#密码不能为空|密码长度应当在:min到:max之间|两次密码输入不相等",
		"password2@required|length:6,16#",
	}
	if e := gvalid.CheckMap(gctx.New(), params, rules); e != nil {
		fmt.Println(e.Map())
		fmt.Println(e.FirstItem())
		fmt.Println(e.FirstString())
	}
	// May Output:
	// map[required:账号不能为空 length:账号长度应当在6到16之间]
	// passport map[required:账号不能为空 length:账号长度应当在6到16之间]
	// 账号不能为空
}

func ExampleCheckMap2() {
	params := map[string]interface{}{
		"passport":  "",
		"password":  "123456",
		"password2": "1234567",
	}
	rules := []string{
		"passport@length:6,16#账号不能为空|账号长度应当在:min到:max之间",
		"password@required|length:6,16|same:password2#密码不能为空|密码长度应当在:min到:max之间|两次密码输入不相等",
		"password2@required|length:6,16#",
	}
	if e := gvalid.CheckMap(gctx.New(), params, rules); e != nil {
		fmt.Println(e.Map())
		fmt.Println(e.FirstItem())
		fmt.Println(e.FirstString())
	}
	// Output:
	// map[same:两次密码输入不相等]
	// password map[same:两次密码输入不相等]
	// 两次密码输入不相等
}

// Empty string attribute.
func ExampleCheckStruct() {
	type Params struct {
		Page      int    `v:"required|min:1         # page is required"`
		Size      int    `v:"required|between:1,100 # size is required"`
		ProjectId string `v:"between:1,10000        # project id must between :min, :max"`
	}
	obj := &Params{
		Page: 1,
		Size: 10,
	}
	err := gvalid.CheckStruct(gctx.New(), obj, nil)
	fmt.Println(err == nil)
	// Output:
	// true
}

// Empty pointer attribute.
func ExampleCheckStruct2() {
	type Params struct {
		Page      int       `v:"required|min:1         # page is required"`
		Size      int       `v:"required|between:1,100 # size is required"`
		ProjectId *gvar.Var `v:"between:1,10000        # project id must between :min, :max"`
	}
	obj := &Params{
		Page: 1,
		Size: 10,
	}
	err := gvalid.CheckStruct(gctx.New(), obj, nil)
	fmt.Println(err == nil)
	// Output:
	// true
}

// Empty integer attribute.
func ExampleCheckStruct3() {
	type Params struct {
		Page      int `v:"required|min:1         # page is required"`
		Size      int `v:"required|between:1,100 # size is required"`
		ProjectId int `v:"between:1,10000        # project id must between :min, :max"`
	}
	obj := &Params{
		Page: 1,
		Size: 10,
	}
	err := gvalid.CheckStruct(gctx.New(), obj, nil)
	fmt.Println(err)
	// Output:
	// project id must between 1, 10000
}

func ExampleRegisterRule() {
	type User struct {
		Id   int
		Name string `v:"required|unique-name # 请输入用户名称|用户名称已被占用"`
		Pass string `v:"required|length:6,18"`
	}
	user := &User{
		Id:   1,
		Name: "john",
		Pass: "123456",
	}

	rule := "unique-name"
	gvalid.RegisterRule(rule, func(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
		var (
			id   = data.(*User).Id
			name = gconv.String(value)
		)
		n, err := g.Model("user").Where("id != ? and name = ?", id, name).Count()
		if err != nil {
			return err
		}
		if n > 0 {
			return errors.New(message)
		}
		return nil
	})
	err := gvalid.CheckStruct(gctx.New(), user, nil)
	fmt.Println(err.Error())
	// May Output:
	// 用户名称已被占用
}

func ExampleRegisterRule_OverwriteRequired() {
	rule := "required"
	gvalid.RegisterRule(rule, func(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
		reflectValue := reflect.ValueOf(value)
		if reflectValue.Kind() == reflect.Ptr {
			reflectValue = reflectValue.Elem()
		}
		isEmpty := false
		switch reflectValue.Kind() {
		case reflect.Bool:
			isEmpty = !reflectValue.Bool()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			isEmpty = reflectValue.Int() == 0
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			isEmpty = reflectValue.Uint() == 0
		case reflect.Float32, reflect.Float64:
			isEmpty = math.Float64bits(reflectValue.Float()) == 0
		case reflect.Complex64, reflect.Complex128:
			c := reflectValue.Complex()
			isEmpty = math.Float64bits(real(c)) == 0 && math.Float64bits(imag(c)) == 0
		case reflect.String, reflect.Map, reflect.Array, reflect.Slice:
			isEmpty = reflectValue.Len() == 0
		}
		if isEmpty {
			return errors.New(message)
		}
		return nil
	})
	fmt.Println(gvalid.CheckValue(gctx.New(), "", "required", "It's required"))
	fmt.Println(gvalid.CheckValue(gctx.New(), 0, "required", "It's required"))
	fmt.Println(gvalid.CheckValue(gctx.New(), false, "required", "It's required"))
	gvalid.DeleteRule(rule)
	fmt.Println("rule deleted")
	fmt.Println(gvalid.CheckValue(gctx.New(), "", "required", "It's required"))
	fmt.Println(gvalid.CheckValue(gctx.New(), 0, "required", "It's required"))
	fmt.Println(gvalid.CheckValue(gctx.New(), false, "required", "It's required"))
	// Output:
	// It's required
	// It's required
	// It's required
	// rule deleted
	// It's required
	// <nil>
	// <nil>
}

func ExampleValidator_Rules() {
	data := g.Map{
		"password": "123",
	}
	err := g.Validator().Data(data).
		Rules("required-with:password").
		Messages("请输入确认密码").
		CheckValue(gctx.New(), "")
	fmt.Println(err.String())

	// Output:
	// 请输入确认密码
}

func ExampleValidator_CheckValue() {
	err := g.Validator().Rules("min:18").
		Messages("未成年人不允许注册哟").
		CheckValue(gctx.New(), 16)
	fmt.Println(err.String())

	// Output:
	// 未成年人不允许注册哟
}

func ExampleValidator_CheckMap() {
	params := map[string]interface{}{
		"passport":  "",
		"password":  "123456",
		"password2": "1234567",
	}
	rules := map[string]string{
		"passport":  "required|length:6,16",
		"password":  "required|length:6,16|same:password2",
		"password2": "required|length:6,16",
	}
	messages := map[string]interface{}{
		"passport": "账号不能为空|账号长度应当在:min到:max之间",
		"password": map[string]string{
			"required": "密码不能为空",
			"same":     "两次密码输入不相等",
		},
	}
	err := g.Validator().
		Messages(messages).
		Rules(rules).
		CheckMap(gctx.New(), params)
	if err != nil {
		g.Dump(err.Maps())
	}

	// May Output:
	//{
	//	"passport": {
	//	"length": "账号长度应当在6到16之间",
	//		"required": "账号不能为空"
	//},
	//	"password": {
	//	"same": "两次密码输入不相等"
	//}
	//}
}

func ExampleValidator_CheckStruct() {
	type User struct {
		Name string `v:"required#请输入用户姓名"`
		Type int    `v:"required#请选择用户类型"`
	}
	data := g.Map{
		"name": "john",
	}
	user := User{}
	if err := gconv.Scan(data, &user); err != nil {
		panic(err)
	}
	err := g.Validator().Data(data).CheckStruct(gctx.New(), user)
	if err != nil {
		fmt.Println(err.Items())
	}

	// Output:
	// [map[Type:map[required:请选择用户类型]]]
}

func ExampleValidator_Required() {
	type BizReq struct {
		Id   uint   `v:"required"`
		Name string `v:"required"`
	}
	var (
		ctx = context.Background()
		req = BizReq{
			Id: 1,
		}
	)
	if err := g.Validator().CheckStruct(ctx, req); err != nil {
		fmt.Println(err)
	}

	// Output:
	// The Name field is required
}

func ExampleValidator_RequiredIf() {
	type BizReq struct {
		Id          uint   `v:"required" dc:"Your Id"`
		Name        string `v:"required" dc:"Your name"`
		Gender      uint   `v:"in:0,1,2" dc:"0:Secret;1:Male;2:Female"`
		WifeName    string `v:"required-if:gender,1"`
		HusbandName string `v:"required-if:gender,2"`
	}
	var (
		ctx = context.Background()
		req = BizReq{
			Id:     1,
			Name:   "test",
			Gender: 1,
		}
	)
	if err := g.Validator().CheckStruct(ctx, req); err != nil {
		fmt.Println(err)
	}

	// Output:
	// The WifeName field is required
}

func ExampleValidator_RequiredUnless() {
	type BizReq struct {
		Id          uint   `v:"required" dc:"Your Id"`
		Name        string `v:"required" dc:"Your name"`
		Gender      uint   `v:"in:0,1,2" dc:"0:Secret;1:Male;2:Female"`
		WifeName    string `v:"required-unless:gender,0,gender,2"`
		HusbandName string `v:"required-unless:id,0,gender,2"`
	}
	var (
		ctx = context.Background()
		req = BizReq{
			Id:     1,
			Name:   "test",
			Gender: 1,
		}
	)
	if err := g.Validator().CheckStruct(ctx, req); err != nil {
		fmt.Println(err)
	}

	// Output:
	// The WifeName field is required; The HusbandName field is required
}

func ExampleValidator_RequiredWith() {
	type BizReq struct {
		Id          uint   `v:"required" dc:"Your Id"`
		Name        string `v:"required" dc:"Your name"`
		Gender      uint   `v:"in:0,1,2" dc:"0:Secret;1:Male;2:Female"`
		WifeName    string
		HusbandName string `v:"required-with:WifeName"`
	}
	var (
		ctx = context.Background()
		req = BizReq{
			Id:       1,
			Name:     "test",
			Gender:   1,
			WifeName: "Ann",
		}
	)
	if err := g.Validator().CheckStruct(ctx, req); err != nil {
		fmt.Println(err)
	}

	// Output:
	// The HusbandName field is required
}

func ExampleValidator_RequiredWithAll() {
	type BizReq struct {
		Id          uint   `v:"required" dc:"Your Id"`
		Name        string `v:"required" dc:"Your name"`
		Gender      uint   `v:"in:0,1,2" dc:"0:Secret;1:Male;2:Female"`
		WifeName    string
		HusbandName string `v:"required-with-all:Id,Name,Gender,WifeName"`
	}
	var (
		ctx = context.Background()
		req = BizReq{
			Id:       1,
			Name:     "test",
			Gender:   1,
			WifeName: "Ann",
		}
	)
	if err := g.Validator().CheckStruct(ctx, req); err != nil {
		fmt.Println(err)
	}

	// Output:
	// The HusbandName field is required
}

func ExampleValidator_RequiredWithout() {
	type BizReq struct {
		Id          uint   `v:"required" dc:"Your Id"`
		Name        string `v:"required" dc:"Your name"`
		Gender      uint   `v:"in:0,1,2" dc:"0:Secret;1:Male;2:Female"`
		WifeName    string
		HusbandName string `v:"required-without:Id,WifeName"`
	}
	var (
		ctx = context.Background()
		req = BizReq{
			Id:     1,
			Name:   "test",
			Gender: 1,
		}
	)
	if err := g.Validator().CheckStruct(ctx, req); err != nil {
		fmt.Println(err)
	}

	// Output:
	// The HusbandName field is required
}

func ExampleValidator_RequiredWithoutAll() {
	type BizReq struct {
		Id          uint   `v:"required" dc:"Your Id"`
		Name        string `v:"required" dc:"Your name"`
		Gender      uint   `v:"in:0,1,2" dc:"0:Secret;1:Male;2:Female"`
		WifeName    string
		HusbandName string `v:"required-without-all:Id,WifeName"`
	}
	var (
		ctx = context.Background()
		req = BizReq{
			Name:   "test",
			Gender: 1,
		}
	)
	if err := g.Validator().CheckStruct(ctx, req); err != nil {
		fmt.Println(err)
	}

	// Output:
	// The HusbandName field is required
}

func ExampleValidator_Same() {
	type BizReq struct {
		Name      string `v:"required"`
		Password  string `v:"required|same:Password2"`
		Password2 string `v:"required"`
	}
	var (
		ctx = context.Background()
		req = BizReq{
			Name:      "gf",
			Password:  "goframe.org",
			Password2: "goframe.net",
		}
	)
	if err := g.Validator().CheckStruct(ctx, req); err != nil {
		fmt.Println(err)
	}

	// Output:
	// The Password value must be the same as field Password2
}

func ExampleValidator_Different() {
	type BizReq struct {
		Name          string `v:"required"`
		MailAddr      string `v:"required"`
		OtherMailAddr string `v:"required|different:MailAddr"`
	}
	var (
		ctx = context.Background()
		req = BizReq{
			Name:          "gf",
			MailAddr:      "gf@goframe.org",
			OtherMailAddr: "gf@goframe.org",
		}
	)
	if err := g.Validator().CheckStruct(ctx, req); err != nil {
		fmt.Println(err)
	}

	// Output:
	// The OtherMailAddr value must be different from field MailAddr
}

func ExampleValidator_In() {
	type BizReq struct {
		Id     uint   `v:"required" dc:"Your Id"`
		Name   string `v:"required" dc:"Your name"`
		Gender uint   `v:"in:0,1,2" dc:"0:Secret;1:Male;2:Female"`
	}
	var (
		ctx = context.Background()
		req = BizReq{
			Id:     1,
			Name:   "test",
			Gender: 3,
		}
	)
	if err := g.Validator().CheckStruct(ctx, req); err != nil {
		fmt.Println(err)
	}

	// Output:
	// The Gender value is not in acceptable range
}

func ExampleValidator_NotIn() {
	type BizReq struct {
		Id           uint   `v:"required" dc:"Your Id"`
		Name         string `v:"required" dc:"Your name"`
		InvalidIndex uint   `v:"not-in:-1,0,1"`
	}
	var (
		ctx = context.Background()
		req = BizReq{
			Id:           1,
			Name:         "test",
			InvalidIndex: 1,
		}
	)
	if err := g.Validator().CheckStruct(ctx, req); err != nil {
		fmt.Println(err)
	}

	// Output:
	// The InvalidIndex value is not in acceptable range
}

func ExampleValidator_Date() {
	type BizReq struct {
		Date1 string `v:"date"`
		Date2 string `v:"date"`
		Date3 string `v:"date"`
		Date4 string `v:"date"`
		Date5 string `v:"date"`
	}

	var (
		ctx = context.Background()
		req = BizReq{
			Date1: "2021-10-31",
			Date2: "2021.10.31",
			Date3: "2021-Oct-31",
			Date4: "2021 Octa 31",
			Date5: "2021/Oct/31",
		}
	)
	if err := g.Validator().CheckStruct(ctx, req); err != nil {
		fmt.Print(err)
	}

	// Output:
	// The Date3 value is not a valid date; The Date4 value is not a valid date; The Date5 value is not a valid date
}

func ExampleValidator_Datetime() {
	type BizReq struct {
		Date1 string `v:"datetime"`
		Date2 string `v:"datetime"`
		Date3 string `v:"datetime"`
		Date4 string `v:"datetime"`
	}

	var (
		ctx = context.Background()
		req = BizReq{
			Date1: "2021-11-01 23:00:00",
			Date2: "2021-11-01 23:00",     // error
			Date3: "2021/11/01 23:00:00",  // error
			Date4: "2021/Dec/01 23:00:00", // error
		}
	)
	if err := g.Validator().CheckStruct(ctx, req); err != nil {
		fmt.Print(err)
	}

	// Output:
	// The Date2 value is not a valid datetime; The Date3 value is not a valid datetime; The Date4 value is not a valid datetime
}

func ExampleValidator_DateFormat() {
	type BizReq struct {
		Date1 string `v:"date-format:Y-m-d"`
		Date2 string `v:"date-format:Y-m-d"`
		Date3 string `v:"date-format:Y-m-d H:i:s"`
		Date4 string `v:"date-format:Y-m-d H:i:s"`
	}

	var (
		ctx = context.Background()
		req = BizReq{
			Date1: "2021-11-01",
			Date2: "2021-11-01 23:00", // error
			Date3: "2021-11-01 23:00:00",
			Date4: "2021-11-01 23:00", // error
		}
	)
	if err := g.Validator().CheckStruct(ctx, req); err != nil {
		fmt.Print(err)
	}

	// Output:
	// The Date2 value does not match the format Y-m-d; The Date4 value does not match the format Y-m-d H:i:s
}

func ExampleValidator_Email() {
	type BizReq struct {
		MailAddr1 string `v:"email"`
		MailAddr2 string `v:"email"`
		MailAddr3 string `v:"email"`
		MailAddr4 string `v:"email"`
	}

	var (
		ctx = context.Background()
		req = BizReq{
			MailAddr1: "gf@goframe.org",
			MailAddr2: "gf@goframe", // error
			MailAddr3: "gf@goframe.org.cn",
			MailAddr4: "gf#goframe.org", // error
		}
	)
	if err := g.Validator().CheckStruct(ctx, req); err != nil {
		fmt.Print(err)
	}

	// Output:
	// The MailAddr2 value must be a valid email address; The MailAddr4 value must be a valid email address
}

func ExampleValidator_Phone() {
	type BizReq struct {
		PhoneNumber1 string `v:"phone"`
		PhoneNumber2 string `v:"phone"`
		PhoneNumber3 string `v:"phone"`
		PhoneNumber4 string `v:"phone"`
	}

	var (
		ctx = context.Background()
		req = BizReq{
			PhoneNumber1: "13578912345",
			PhoneNumber2: "11578912345", // error 11x not exist
			PhoneNumber3: "17178912345", // error 171 not exit
			PhoneNumber4: "1357891234",  // error len must be 11
		}
	)
	if err := g.Validator().CheckStruct(ctx, req); err != nil {
		fmt.Print(err)
	}

	// Output:
	// The PhoneNumber2 value must be a valid phone number; The PhoneNumber3 value must be a valid phone number; The PhoneNumber4 value must be a valid phone number
}

func ExampleValidator_PhoneLoose() {
	type BizReq struct {
		PhoneNumber1 string `v:"phone-loose"`
		PhoneNumber2 string `v:"phone-loose"`
		PhoneNumber3 string `v:"phone-loose"`
		PhoneNumber4 string `v:"phone-loose"`
	}

	var (
		ctx = context.Background()
		req = BizReq{
			PhoneNumber1: "13578912345",
			PhoneNumber2: "11578912345", // error 11x not exist
			PhoneNumber3: "17178912345",
			PhoneNumber4: "1357891234", // error len must be 11
		}
	)
	if err := g.Validator().CheckStruct(ctx, req); err != nil {
		fmt.Print(err)
	}

	// Output:
	// The PhoneNumber2 value must be a valid phone number; The PhoneNumber4 value must be a valid phone number
}

func ExampleValidator_Telephone() {
	type BizReq struct {
		Telephone1 string `v:"telephone"`
		Telephone2 string `v:"telephone"`
		Telephone3 string `v:"telephone"`
		Telephone4 string `v:"telephone"`
	}

	var (
		ctx = context.Background()
		req = BizReq{
			Telephone1: "010-77542145",
			Telephone2: "0571-77542145",
			Telephone3: "20-77542145", // error
			Telephone4: "775421451",   // error len must be 7 or 8
		}
	)
	if err := g.Validator().CheckStruct(ctx, req); err != nil {
		fmt.Print(err)
	}

	// Output:
	// The Telephone3 value must be a valid telephone number; The Telephone4 value must be a valid telephone number
}
