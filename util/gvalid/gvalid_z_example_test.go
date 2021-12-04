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
	"github.com/gogf/gf/v2/i18n/gi18n"
	"math"
	"reflect"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gvalid"
)

func ExampleNew() {
	validator := gvalid.New()

	if err := validator.Data(16).Rules("min:18").Run(context.Background()); err != nil {
		fmt.Print(err)
	}

	// Output:
	// The value `16` must be equal or greater than 18
}

func ExampleValidator_Run() {
	// check value mode
	if err := g.Validator().Data(16).Rules("min:18").Run(context.Background()); err != nil {
		fmt.Println("check value err:", err)
	}
	// check map mode
	data := map[string]interface{}{
		"passport":  "",
		"password":  "123456",
		"password2": "1234567",
	}
	rules := map[string]string{
		"passport":  "required|length:6,16",
		"password":  "required|length:6,16|same:password2",
		"password2": "required|length:6,16",
	}
	if err := g.Validator().Data(data).Rules(rules).Run(context.Background()); err != nil {
		fmt.Println("check map err:", err)
	}
	// check struct mode
	type Params struct {
		Page      int    `v:"required|min:1"`
		Size      int    `v:"required|between:1,100"`
		ProjectId string `v:"between:1,10000"`
	}
	rules = map[string]string{
		"Page":      "required|min:1",
		"Size":      "required|between:1,100",
		"ProjectId": "between:1,10000",
	}
	obj := &Params{
		Page: 0,
		Size: 101,
	}
	if err := g.Validator().Data(obj).Run(context.Background()); err != nil {
		fmt.Println("check struct err:", err)
	}

	// May Output:
	// check value err: The value `16` must be equal or greater than 18
	// check map err: The passport field is required; The passport value `` length must be between 6 and 16; The password value `123456` must be the same as field password2
	// check struct err: The Page value `0` must be equal or greater than 1; The Size value `101` must be between 1 and 100
}

func ExampleValidator_Clone() {
	if err := g.Validator().Data(16).Rules("min:18").Run(context.Background()); err != nil {
		fmt.Println(err)
	}

	if err := g.Validator().Clone().Data(20).Run(context.Background()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Check Success!")
	}

	// Output:
	// The value `16` must be equal or greater than 18
	// Check Success!
}

func ExampleValidator_I18n() {
	var (
		i18nManager = gi18n.New()
		ctxCn       = gi18n.WithLanguage(context.Background(), "cn")
		validator   = gvalid.New()
	)

	validator = validator.Data(16).Rules("min:18")

	if err := validator.Run(context.Background()); err != nil {
		fmt.Println(err)
	}

	if err := validator.I18n(i18nManager).Run(ctxCn); err != nil {
		fmt.Println(err)
	}

	// Output:
	// The value `16` must be equal or greater than 18
	// 字段值`16`字段最小值应当为18
}

func ExampleValidator_Bail() {
	type BizReq struct {
		Account   string `v:"required|length:6,16|same:QQ"`
		QQ        string
		Password  string `v:"required|same:Password2"`
		Password2 string `v:"required"`
	}
	var (
		ctx = context.Background()
		req = BizReq{
			Account:   "gf",
			QQ:        "123456",
			Password:  "goframe.org",
			Password2: "goframe.org",
		}
	)

	if err := g.Validator().Bail().Data(req).Run(ctx); err != nil {
		fmt.Println("Use Bail Error:", err)
	}

	if err := g.Validator().Data(req).Run(ctx); err != nil {
		fmt.Println("Not Use Bail Error:", err)
	}

	// output:
	// Use Bail Error: The Account value `gf` length must be between 6 and 16
	// Not Use Bail Error: The Account value `gf` length must be between 6 and 16; The Account value `gf` must be the same as field QQ
}

func ExampleValidator_Ci() {

	type BizReq struct {
		Account   string `v:"required"`
		Password  string `v:"required|same:Password2"`
		Password2 string `v:"required"`
	}

	var (
		ctx = context.Background()
		req = BizReq{
			Account:   "gf",
			Password:  "Goframe.org", // Diff from Password2, but because of "ci", rule check passed
			Password2: "goframe.org",
		}
	)

	if err := g.Validator().Data(req).Run(ctx); err != nil {
		fmt.Println("Not Use CI Error:", err)
	}

	if err := g.Validator().Ci().Data(req).Run(ctx); err == nil {
		fmt.Println("Use CI Passed!")
	}

	// output:
	// Not Use CI Error: The Password value `Goframe.org` must be the same as field Password2
	// Use CI Passed!
}

func ExampleValidator_Data() {
	type BizReq struct {
		Password1 string `v:"password"`
		Password2 string `v:"password"`
	}

	var (
		ctx = context.Background()
		req = BizReq{
			Password1: "goframe",
			Password2: "gofra", // error length between 6 and 18
		}
	)
	if err := g.Validator().Data(req).Run(ctx); err != nil {
		fmt.Print(err)
	}

	// Output:
	// The Password2 value `gofra` is not a valid password format
}

func ExampleValidator_Assoc() {

	type User struct {
		Name string `v:"required"`
		Type int    `v:"required"`
	}

	data := g.Map{
		"name": "john",
	}

	user := User{}

	if err := gconv.Scan(data, &user); err != nil {
		panic(err)
	}

	if err := g.Validator().Data(user).Assoc(data).Run(context.Background()); err != nil {
		fmt.Print(err)
	}

	// Output:
	// The Type field is required
}

func ExampleValidator_Rules() {

	if err := g.Validator().Data(16).Rules("min:18").Run(context.Background()); err != nil {
		fmt.Println(err)
	}

	// Output:
	// The value `16` must be equal or greater than 18
}

func ExampleValidator_Messages() {
	if err := g.Validator().Data(16).Rules("min:18").Messages("Can not regist, Age is less then 18!").Run(context.Background()); err != nil {
		fmt.Println(err)
	}

	// Output:
	// Can not regist, Age is less then 18!
}

func ExampleValidator_RuleFunc() {
	var (
		ctx             = context.Background()
		lenErrRuleName  = "LenErr"
		passErrRuleName = "PassErr"
		lenErrRuleFunc  = func(ctx context.Context, in gvalid.RuleFuncInput) error {
			pass := in.Value.String()
			if len(pass) != 6 {
				return errors.New(in.Message)
			}
			return nil
		}
		passErrRuleFunc = func(ctx context.Context, in gvalid.RuleFuncInput) error {
			pass := in.Value.String()
			if m := in.Data.Map(); m["data"] != pass {
				return errors.New(in.Message)
			}
			return nil
		}
	)

	type LenErrStruct struct {
		Value string `v:"uid@LenErr#Value Length Error!"`
		Data  string `p:"data"`
	}

	st := &LenErrStruct{
		Value: "123",
		Data:  "123456",
	}
	// single error sample
	if err := g.Validator().RuleFunc(lenErrRuleName, lenErrRuleFunc).Data(st).Run(ctx); err != nil {
		fmt.Println(err)
	}

	type MultiErrorStruct struct {
		Value string `v:"uid@LenErr|PassErr#Value Length Error!|Pass is not Same!"`
		Data  string `p:"data"`
	}

	multi := &MultiErrorStruct{
		Value: "123",
		Data:  "123456",
	}
	// multi error sample
	if err := g.Validator().RuleFunc(lenErrRuleName, lenErrRuleFunc).RuleFunc(passErrRuleName, passErrRuleFunc).Data(multi).Run(ctx); err != nil {
		fmt.Println(err)
	}

	// Output:
	// Value Length Error!
	// Value Length Error!; Pass is not Same!
}

func ExampleValidator_RuleFuncMap() {
	var (
		ctx             = context.Background()
		lenErrRuleName  = "LenErr"
		passErrRuleName = "PassErr"
		lenErrRuleFunc  = func(ctx context.Context, in gvalid.RuleFuncInput) error {
			pass := in.Value.String()
			if len(pass) != 6 {
				return errors.New(in.Message)
			}
			return nil
		}
		passErrRuleFunc = func(ctx context.Context, in gvalid.RuleFuncInput) error {
			pass := in.Value.String()
			if m := in.Data.Map(); m["data"] != pass {
				return errors.New(in.Message)
			}
			return nil
		}
		ruleMap = map[string]gvalid.RuleFunc{
			lenErrRuleName:  lenErrRuleFunc,
			passErrRuleName: passErrRuleFunc,
		}
	)

	type MultiErrorStruct struct {
		Value string `v:"uid@LenErr|PassErr#Value Length Error!|Pass is not Same!"`
		Data  string `p:"data"`
	}

	multi := &MultiErrorStruct{
		Value: "123",
		Data:  "123456",
	}

	if err := g.Validator().RuleFuncMap(ruleMap).Data(multi).Run(ctx); err != nil {
		fmt.Println(err)
	}

	// Output:
	// Value Length Error!; Pass is not Same!
}

func ExampleCheckMap() {
	params := map[string]interface{}{
		"passport":  "",
		"password":  "123456",
		"password2": "1234567",
	}
	rules := []string{
		"passport@required|length:6,16#账号不能为空|账号长度应当在{min}到{max}之间",
		"password@required|length:6,16|same{password}2#密码不能为空|密码长度应当在{min}到{max}之间|两次密码输入不相等",
		"password2@required|length:6,16#",
	}
	if e := g.Validator().Data(params).Rules(rules).Run(gctx.New()); e != nil {
		fmt.Println(e.Map())
		fmt.Println(e.FirstItem())
		fmt.Println(e.FirstError())
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
		"passport@length:6,16#账号不能为空|账号长度应当在{min}到{max}之间",
		"password@required|length:6,16|same:password2#密码不能为空|密码长度应当在{min}到{max}之间|两次密码输入不相等",
		"password2@required|length:6,16#",
	}
	if e := g.Validator().Data(params).Rules(rules).Run(gctx.New()); e != nil {
		fmt.Println(e.Map())
		fmt.Println(e.FirstItem())
		fmt.Println(e.FirstError())
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
		ProjectId string `v:"between:1,10000        # project id must between {min}, {max}"`
	}
	obj := &Params{
		Page: 1,
		Size: 10,
	}
	err := g.Validator().Data(obj).Rules(nil).Run(gctx.New())
	fmt.Println(err == nil)
	// Output:
	// true
}

// Empty pointer attribute.
func ExampleCheckStruct2() {
	type Params struct {
		Page      int       `v:"required|min:1         # page is required"`
		Size      int       `v:"required|between:1,100 # size is required"`
		ProjectId *gvar.Var `v:"between:1,10000        # project id must between {min}, {max}"`
	}
	obj := &Params{
		Page: 1,
		Size: 10,
	}
	err := g.Validator().Data(obj).Rules(nil).Run(gctx.New())
	fmt.Println(err == nil)
	// Output:
	// true
}

// Empty integer attribute.
func ExampleCheckStruct3() {
	type Params struct {
		Page      int `v:"required|min:1         # page is required"`
		Size      int `v:"required|between:1,100 # size is required"`
		ProjectId int `v:"between:1,10000        # project id must between {min}, {max}"`
	}
	obj := &Params{
		Page: 1,
		Size: 10,
	}
	err := g.Validator().Data(obj).Rules(nil).Run(gctx.New())
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
	gvalid.RegisterRule(rule, func(ctx context.Context, in gvalid.RuleFuncInput) error {
		var (
			id   = in.Data.Val().(*User).Id
			name = gconv.String(in.Value)
		)
		n, err := g.Model("user").Where("id != ? and name = ?", id, name).Count()
		if err != nil {
			return err
		}
		if n > 0 {
			return errors.New(in.Message)
		}
		return nil
	})
	err := g.Validator().Data(user).Rules(nil).Run(gctx.New())
	fmt.Println(err.Error())
	// May Output:
	// 用户名称已被占用
}

func ExampleRegisterRule_OverwriteRequired() {
	rule := "required"
	gvalid.RegisterRule(rule, func(ctx context.Context, in gvalid.RuleFuncInput) error {
		reflectValue := reflect.ValueOf(in.Value.Val())
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
			return errors.New(in.Message)
		}
		return nil
	})
	fmt.Println(g.Validator().Data("").Rules("required").Messages("It's required").Run(gctx.New()))
	fmt.Println(g.Validator().Data(0).Rules("required").Messages("It's required").Run(gctx.New()))
	fmt.Println(g.Validator().Data(false).Rules("required").Messages("It's required").Run(gctx.New()))
	gvalid.DeleteRule(rule)
	fmt.Println("rule deleted")
	fmt.Println(g.Validator().Data("").Rules("required").Messages("It's required").Run(gctx.New()))
	fmt.Println(g.Validator().Data(0).Rules("required").Messages("It's required").Run(gctx.New()))
	fmt.Println(g.Validator().Data(false).Rules("required").Messages("It's required").Run(gctx.New()))
	// Output:
	// It's required
	// It's required
	// It's required
	// rule deleted
	// It's required
	// <nil>
	// <nil>
}

func ExampleValidator_CheckValue() {
	err := g.Validator().Rules("min:18").
		Messages("未成年人不允许注册哟").
		Data(16).Run(gctx.New())
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
		"passport": "账号不能为空|账号长度应当在{min}到{max}之间",
		"password": map[string]string{
			"required": "密码不能为空",
			"same":     "两次密码输入不相等",
		},
	}
	err := g.Validator().
		Messages(messages).
		Rules(rules).
		Data(params).Run(gctx.New())
	if err != nil {
		g.Dump(err.Maps())
	}

	// May Output:
	// {
	//	"passport": {
	//	"length": "账号长度应当在6到16之间",
	//		"required": "账号不能为空"
	// },
	//	"password": {
	//	"same": "两次密码输入不相等"
	// }
	// }
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
	err := g.Validator().Data(user).Assoc(data).Run(gctx.New())
	if err != nil {
		fmt.Println(err.Items())
	}

	// Output:
	// [map[Type:map[required:请选择用户类型]]]
}
