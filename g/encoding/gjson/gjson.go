// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// JSON解析/封装
package gjson

import (
    "strings"
    "strconv"
    "io/ioutil"
    "encoding/json"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/encoding/gxml"
    "gitee.com/johng/gf/g/encoding/gyaml"
)

// json解析结果存放数组
type Json struct {
    // 注意这是一个指针
    value *interface{}
}

// 一个json变量
type JsonVar interface{}

// 编码go变量为json字符串，并返回json字符串指针
func Encode (v interface{}) ([]byte, error) {
    return json.Marshal(v)
}

// 解码字符串为interface{}变量
func Decode (b []byte) (interface{}, error) {
    var v interface{}
    if err := DecodeTo(b, &v); err != nil {
        return nil, err
    } else {
        return v, nil
    }
}

// 解析json字符串为go变量，注意第二个参数为指针(任意结构的变量)
func DecodeTo (b []byte, v interface{}) error {
    return json.Unmarshal(b, v)
}

// 解析json字符串为gjson.Json对象，并返回操作对象指针
func DecodeToJson (b []byte) (*Json, error) {
    if v, err := Decode(b); err != nil {
        return nil, err
    } else {
        return &Json{&v}, nil
    }
}

// 支持多种配置文件类型转换为json格式内容并解析为gjson.Json对象
func Load (path string) (*Json, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }
    return LoadContent(data, gfile.Ext(path))
}

// 支持的配置文件格式：xml, json, yml
func LoadContent (data []byte, t string) (*Json, error) {
    var err    error
    var result interface{}
    switch t {
        case "xml":  fallthrough
        case ".xml":
            data, err = gxml.ToJson(data)
            if err != nil {
                return nil, err
            }
        case "yml":  fallthrough
        case "yaml": fallthrough
        case ".yml": fallthrough
        case ".yaml":
            data, err = gyaml.ToJson(data)
            if err != nil {
                return nil, err
            }
    }
    if err := json.Unmarshal(data, &result); err != nil {
        return nil, err
    }
    return &Json{ &result }, nil
}

// 将变量转换为Json对象进行处理，该变量至少应当是一个map或者array，否者转换没有意义
func NewJson(v *interface{}) *Json {
    return &Json{ v }
}

// 将指定的json内容转换为指定结构返回，查找失败或者转换失败，目标对象转换为nil
// 注意第二个参数需要给的是变量地址
func (p *Json) GetToVar(pattern string, v interface{}) error {
    r := p.Get(pattern)
    if r != nil {
        if t, err := Encode(r); err == nil {
            return DecodeTo(t, v)
        } else {
            return err
        }
    } else {
        v = nil
    }
    return nil
}

// 获得一个键值对关联数组/哈希表，方便操作，不需要自己做类型转换
// 注意，如果获取的值不存在，或者类型与json类型不匹配，那么将会返回nil
func (p *Json) GetMap(pattern string) map[string]interface{} {
    result := p.Get(pattern)
    if result != nil {
        if r, ok := result.(map[string]interface{}); ok {
            return r
        }
    }
    return nil
}

// 将检索值转换为Json对象指针返回
func (p *Json) GetJson(pattern string) *Json {
    result := p.Get(pattern)
    if result != nil {
        return &Json{&result}
    }
    return nil
}

// 获得一个数组[]interface{}，方便操作，不需要自己做类型转换
// 注意，如果获取的值不存在，或者类型与json类型不匹配，那么将会返回nil
func (p *Json) GetArray(pattern string) []interface{} {
    result := p.Get(pattern)
    if result != nil {
        if r, ok := result.([]interface{}); ok {
            return r
        }
    }
    return nil
}

// 返回指定json中的string
func (p *Json) GetString(pattern string) string {
    return gconv.String(p.Get(pattern))
}

// 返回指定json中的bool(false:"", 0, false, off)
func (p *Json) GetBool(pattern string) bool {
    return gconv.Bool(p.Get(pattern))
}

func (p *Json) GetInt(pattern string) int {
    return gconv.Int(p.Get(pattern))
}

func (p *Json) GetUint(pattern string) uint {
    return gconv.Uint(p.Get(pattern))
}

func (p *Json) GetFloat32(pattern string) float32 {
    return gconv.Float32(p.Get(pattern))
}

func (p *Json) GetFloat64(pattern string) float64 {
    return gconv.Float64(p.Get(pattern))
}

// 根据约定字符串方式访问json解析数据，参数形如： "items.name.first", "list.0"
// 返回的结果类型的interface{}，因此需要自己做类型转换
// 如果找不到对应节点的数据，返回nil
func (p *Json) Get(pattern string) interface{} {
    var result interface{}
    pointer  := p.value
    array    := strings.Split(pattern, ".")
    length   := len(array)
    for i:= 0; i < length; i++ {
        switch (*pointer).(type) {
            case map[string]interface{}:
                if v, ok := (*pointer).(map[string]interface{})[array[i]]; ok {
                    if i == length - 1 {
                        result = v
                    } else {
                        pointer = &v
                    }
                } else {
                    return nil
                }
            case []interface{}:
                if isNumeric(array[i]) {
                    n, err := strconv.Atoi(array[i])
                    if err == nil && len((*pointer).([]interface{})) > n {
                        if i == length - 1 {
                            result = (*pointer).([]interface{})[n]
                            break;
                        } else {
                            pointer = &(*pointer).([]interface{})[n]
                        }
                    }
                } else {
                    return nil
                }
            default:
                return nil
        }
    }
    return result
}

// 转换为map[string]interface{}类型,如果转换失败，返回nil
func (p *Json) ToMap() map[string]interface{} {
    switch (*(p.value)).(type) {
        case map[string]interface{}:
            return (*(p.value)).(map[string]interface{})
        default:
            return nil
    }
}

// 转换为[]interface{}类型,如果转换失败，返回nil
func (p *Json) ToArray() []interface{} {
    switch (*(p.value)).(type) {
        case []interface{}:
            return (*(p.value)).([]interface{})
        default:
            return nil
    }
}

func (p *Json) ToXml(rootTag...string) ([]byte, error) {
    return gxml.Encode(p.ToMap(), rootTag...)
}

func (p *Json) ToXmlIndent(rootTag...string) ([]byte, error) {
    return gxml.EncodeWithIndent(p.ToMap(), rootTag...)
}

func (p *Json) ToJson() ([]byte, error) {
    return Encode(*(p.value))
}

func (p *Json) ToJsonIndent() ([]byte, error) {
    return json.MarshalIndent(*(p.value), "", "\t")
}

func (p *Json) ToYaml() ([]byte, error) {
    return gyaml.Encode(*(p.value))
}

// 判断所给字符串是否为数字
func isNumeric(s string) bool  {
    for i := 0; i < len(s); i++ {
        if s[i] < byte('0') || s[i] > byte('9') {
            return false
        }
    }
    return true
}