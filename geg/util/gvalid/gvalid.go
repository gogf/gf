package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/util/gvalid"
)

func main() {
    //rule := "length:6,16"
    //if m := gvalid.Check("123456", rule, nil);  m != nil {
    //    fmt.Println(m)
    //}
    //if m := gvalid.Check("12345", rule, nil);  m != nil {
    //    fmt.Println(m)
    //    // map[length:字段长度为6到16个字符]
    //}

    //rule := "integer|between:6,16"
    //msgs := "请输入一个整数|参数大小不对啊老铁"
    //fmt.Println(gvalid.Check("5.66", rule, msgs))
    //// map[integer:请输入一个整数 between:参数大小不对啊老铁]

    //// 参数长度至少为6个数字或者6个字母，但是总长度不能超过16个字符
    //rule := `regex:\d{6,}|\D{6,}|max-length:16`
    //if m := gvalid.Check("123456", rule, nil);  m != nil {
    //    fmt.Println(m)
    //}
    //if m := gvalid.Check("abcde6", rule, nil);  m != nil {
    //    fmt.Println(m)
    //    // map[regex:字段值不合法]
    //}

    //params := map[string]string {
    //    "passport"  : "john",
    //    "password"  : "123456",
    //    "password2" : "1234567",
    //}
    //rules := map[string]string {
    //    "passport"  : "required|length:6,16",
    //    "password"  : "required|length:6,16|same:password2",
    //    "password2" : "required|length:6,16",
    //}
    //fmt.Println(gvalid.CheckMap(params, rules))
    //// map[passport:map[length:字段长度为6到16个字符] password:map[same:字段值不合法]]

    params := map[string]interface{} {
        "passport"  : "john",
        "password"  : "123456",
        "password2" : "1234567",
    }
    rules := map[string]string {
        "passport"  : "required|length:6,16",
        "password"  : "required|length:6,16|same:password2",
        "password2" : "required|length:6,16",
    }
    msgs  := map[string]interface{} {
        "passport" : "账号不能为空|账号长度应当在:min到:max之间",
        "password" : map[string]string {
            "required" : "密码不能为空",
            "same"     : "两次密码输入不相等",
        },
    }
    if e := gvalid.CheckMap(params, rules, msgs); e != nil {
        g.Dump(e.Maps())
    }
    // map[passport:map[length:账号长度应当在6到16之间] password:map[same:两次密码输入不相等]]
}
