// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// JSON解析/封装
package gjson

import (
    "sync"
    "errors"
    "strings"
    "strconv"
    "io/ioutil"
    "encoding/json"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/encoding/gxml"
    "gitee.com/johng/gf/g/encoding/gyaml"
    "gitee.com/johng/gf/g/encoding/gtoml"
)

// json解析结果存放数组
type Json struct {
    mu sync.RWMutex
    p  *interface{} // 注意这是一个指针
}

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
        return NewJson(v), nil
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

// 支持的配置文件格式：xml, json, yaml/yml, toml
func LoadContent (data []byte, t string) (*Json, error) {
    var err    error
    var result interface{}
    switch t {
        case  "xml":  fallthrough
        case ".xml":
            data, err = gxml.ToJson(data)
            if err != nil {
                return nil, err
            }
        case   "yml": fallthrough
        case  "yaml": fallthrough
        case  ".yml": fallthrough
        case ".yaml":
            data, err = gyaml.ToJson(data)
            if err != nil {
                return nil, err
            }

        case  "toml": fallthrough
        case ".toml":
            data, err = gtoml.ToJson(data)
            if err != nil {
                return nil, err
            }
    }
    if err := json.Unmarshal(data, &result); err != nil {
        return nil, err
    }
    return NewJson(result), nil
}

// 将变量转换为Json对象进行处理，该变量至少应当是一个map或者array，否者转换没有意义
func NewJson(v interface{}) *Json {
    return &Json{ p: &v }
}

// 将指定的json内容转换为指定结构返回，查找失败或者转换失败，目标对象转换为nil
// 注意第二个参数需要给的是变量地址
func (j *Json) GetToVar(pattern string, v interface{}) error {
    r := j.Get(pattern)
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
func (j *Json) GetMap(pattern string) map[string]interface{} {
    result := j.Get(pattern)
    if result != nil {
        if r, ok := result.(map[string]interface{}); ok {
            return r
        }
    }
    return nil
}

// 将检索值转换为Json对象指针返回
func (j *Json) GetJson(pattern string) *Json {
    result := j.Get(pattern)
    if result != nil {
        return NewJson(result)
    }
    return nil
}

// 获得一个数组[]interface{}，方便操作，不需要自己做类型转换
// 注意，如果获取的值不存在，或者类型与json类型不匹配，那么将会返回nil
func (j *Json) GetArray(pattern string) []interface{} {
    result := j.Get(pattern)
    if result != nil {
        if r, ok := result.([]interface{}); ok {
            return r
        }
    }
    return nil
}

// 返回指定json中的string
func (j *Json) GetString(pattern string) string {
    return gconv.String(j.Get(pattern))
}

// 返回指定json中的bool(false:"", 0, false, off)
func (j *Json) GetBool(pattern string) bool {
    return gconv.Bool(j.Get(pattern))
}

func (j *Json) GetInt(pattern string) int {
    return gconv.Int(j.Get(pattern))
}

func (j *Json) GetUint(pattern string) uint {
    return gconv.Uint(j.Get(pattern))
}

func (j *Json) GetFloat32(pattern string) float32 {
    return gconv.Float32(j.Get(pattern))
}

func (j *Json) GetFloat64(pattern string) float64 {
    return gconv.Float64(j.Get(pattern))
}

// 根据pattern查找并设置数据
// 注意：
// 1、写入的时候"."符号只能表示层级，不能使用带"."符号的键名
// 2、写入的value为nil时，表示删除
func (j *Json) Set(pattern string, value interface{}) error {
    // 初始化判断
    if *j.p == nil {
        if isNumeric(pattern) {
            *j.p = make([]interface{}, 0)
        } else {
            *j.p = make(map[string]interface{})
        }
    }
    value  = j.convertValue(value)
    array := strings.Split(pattern, ".")
    // root节点
    if len(array) == 1 {
        return j.setRoot(pattern, value)
    }
    j.mu.Lock()
    pointer  := j.p
    length   := len(array)
    for i:= 0; i < length; i++ {
        switch (*pointer).(type) {
            case map[string]interface{}:
                if i == length - 1 {
                    if value == nil {
                        delete((*pointer).(map[string]interface{}), array[i])
                    } else {
                        (*pointer).(map[string]interface{})[array[i]] = value
                    }
                } else {
                    v, ok := (*pointer).(map[string]interface{})[array[i]]
                    if !ok {
                        // 如果父级键也不存在，并且是删除操作，那么直接返回
                        if value == nil {
                            goto done
                        }
                        if strings.Compare(array[i], "0") == 0 {
                            v = make([]interface{}, 0)
                        } else {
                            v = make(map[string]interface{})
                        }
                        (*pointer).(map[string]interface{})[array[i]] = v
                    }
                    pointer = &v
                }
            case []interface{}:
                if isNumeric(array[i]) {
                    if n, err := strconv.Atoi(array[i]); err == nil {
                        if i == length - 1 {
                            if cap((*pointer).([]interface{})) >= n {
                                if value == nil {
                                    j.mu.Unlock()
                                    return j.Set(strings.Join(array[0 : i], "."), append((*pointer).([]interface{})[ : n], (*pointer).([]interface{})[n + 1 : ]...))
                                } else {
                                    (*pointer).([]interface{})[n] = value
                                }
                            } else {
                                if value != nil {
                                    // 注意这里产生了临时变量和赋值拷贝
                                    s := make([]interface{}, n + 1)
                                    copy(s, (*pointer).([]interface{}))
                                    s[n] = value
                                    j.mu.Unlock()
                                    return j.Set(strings.Join(array[0 : i], "."), s)
                                }
                            }
                        } else {
                            pointer = &(*pointer).([]interface{})[n]
                        }
                    } else {
                        j.mu.Unlock()
                        return err
                    }
                } else {
                    j.mu.Unlock()
                    return errors.New("\"" + strings.Join(array[0:i], ".") + "\" is array, invalid index - \"" + array[i] + "\"")
                }
        }
    }
done:
    j.mu.Unlock()
    return nil
}

// 数据结构转换，map参数必须转换为map[string]interface{}，数组参数必须转换为[]interface{}
func (j *Json) convertValue(value interface{}) interface{} {
    switch value.(type) {
        case map[string]interface{}:
            return value
        case []interface{}:
            return value
        default:
            // 这里效率会比较低
            b, _ := Encode(value)
            v, _ := Decode(b)
            return v
    }
    return value
}

// 修改根节点数据
func (j *Json) setRoot(pattern string, value interface{}) error {
    switch (*j.p).(type) {
        case map[string]interface{}:
            if value == nil {
                delete((*j.p).(map[string]interface{}), pattern)
            } else {
                (*j.p).(map[string]interface{})[pattern] = value
            }
        case []interface{}:
            if isNumeric(pattern) {
                if n, err := strconv.Atoi(pattern); err != nil {
                    return err
                } else {
                    if cap((*j.p).([]interface{})) >= n {
                        if value == nil {
                            (*j.p) = append((*j.p).([]interface{})[ : n], (*j.p).([]interface{})[n + 1 : ]...)
                        } else {
                            (*j.p).([]interface{})[n] = value
                        }
                    } else {
                        if value != nil {
                            // 注意这里产生了临时变量和赋值拷贝
                            array := (*j.p).([]interface{})
                            array  = append(array, value)
                            *j.p   = array
                        }
                    }
                }
            }
    }
    return nil
}

// 根据约定字符串方式访问json解析数据，参数形如： "items.name.first", "list.0"
// 返回的结果类型的interface{}，因此需要自己做类型转换
// 如果找不到对应节点的数据，返回nil
func (j *Json) Get(pattern string) interface{} {
    j.mu.RLock()
    defer j.mu.RUnlock()
    if r := j.getPointerByPattern(pattern); r != nil {
        return *r
    }
    return nil
}

// 根据pattern层级查找变量指针
func (j *Json) getPointerByPattern(pattern string) *interface{} {
    start   := 0
    index   := len(pattern)
    length  := 0
    pointer := j.p
    for {
        if r := j.checkPatternByPointer(pattern[start:index], pointer); r != nil {
            length += index - start
            if start > 0 {
                length += 1
            }
            start = index + 1
            index = len(pattern)
            if length == len(pattern) {
                return r
            } else {
                pointer = r
            }
        } else {
            index = strings.LastIndex(pattern[start:index], ".")
            if index != -1 && length > 0 {
                index += length + 1
            }
        }
        if start >= index {
            break
        }
    }
    return nil
}

// 判断给定的pattern在当前的pointer下是否有值，并返回对应的pointer
// 注意这里返回的指针都是临时变量的内存地址
func (j *Json) checkPatternByPointer(pattern string, pointer *interface{}) *interface{} {
    switch (*pointer).(type) {
        case map[string]interface{}:
            if v, ok := (*pointer).(map[string]interface{})[pattern]; ok {
                return &v
            }
        case []interface{}:
            if isNumeric(pattern) {
                n, err := strconv.Atoi(pattern)
                if err == nil && len((*pointer).([]interface{})) > n {
                    return &(*pointer).([]interface{})[n]
                }
            }
    }
    return nil
}

// 转换为map[string]interface{}类型,如果转换失败，返回nil
func (j *Json) ToMap() map[string]interface{} {
    j.mu.RLock()
    defer j.mu.RUnlock()
    switch (*(j.p)).(type) {
        case map[string]interface{}:
            return (*(j.p)).(map[string]interface{})
        default:
            return nil
    }
}

// 转换为[]interface{}类型,如果转换失败，返回nil
func (j *Json) ToArray() []interface{} {
    j.mu.RLock()
    defer j.mu.RUnlock()
    switch (*(j.p)).(type) {
        case []interface{}:
            return (*(j.p)).([]interface{})
        default:
            return nil
    }
}

func (j *Json) ToXml(rootTag...string) ([]byte, error) {
    return gxml.Encode(j.ToMap(), rootTag...)
}

func (j *Json) ToXmlIndent(rootTag...string) ([]byte, error) {
    return gxml.EncodeWithIndent(j.ToMap(), rootTag...)
}

func (j *Json) ToJson() ([]byte, error) {
    j.mu.RLock()
    defer j.mu.RUnlock()
    return Encode(*(j.p))
}

func (j *Json) ToJsonIndent() ([]byte, error) {
    j.mu.RLock()
    defer j.mu.RUnlock()
    return json.MarshalIndent(*(j.p), "", "\t")
}

func (j *Json) ToYaml() ([]byte, error) {
    j.mu.RLock()
    defer j.mu.RUnlock()
    return gyaml.Encode(*(j.p))
}

func (j *Json) ToToml() ([]byte, error) {
    j.mu.RLock()
    defer j.mu.RUnlock()
    return gtoml.Encode(*(j.p))
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