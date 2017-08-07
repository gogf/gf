package gjson

import (
    "encoding/json"
    "errors"
    "strings"
    "strconv"
    "log"
)

// json解析结果存放数组
type Json struct {
    // 注意这是一个指针
    value *interface{}
}

// 一个json变量
type JsonVar interface{}

// 解析json字符串为go变量，并返回操作对象指针
func Decode (s *string) *Json {
    var result interface{}
    if err := json.Unmarshal([]byte(*s), &result); err != nil {
        log.Println("json unmarshaling failed: " + err.Error())
        return nil
    }
    return &Json{ &result }
}

// 解析json字符串为go变量，注意第二个参数为指针
func DecodeTo (s *string, v interface{}) error {
    if err := json.Unmarshal([]byte(*s), v); err != nil {
        return errors.New("json unmarshaling failed: " + err.Error())
    }
    return nil
}

// 解析go变量为json字符串，并返回json字符串指针
func Encode (v interface{}) *string {
    s, err := json.Marshal(v)
    if err != nil {
        log.Println("json marshaling failed: " + err.Error())
        return nil
    }
    r := string(s)
    return &r
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
    result := p.Get(pattern)
    if result != nil {
        if r, ok := result.(string); ok {
            return r
        }
    }
    return ""
}

// 返回指定json中的float64
func (p *Json) GetNumber(pattern string) float64 {
    result := p.Get(pattern)
    if result != nil {
        if r, ok := result.(float64); ok {
            return r
        }
    }
    return 0
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
    pointer := p.value
    switch (*pointer).(type) {
        case map[string]interface{}:
            return (*pointer).(map[string]interface{})
        default:
            return nil
    }
}

// 转换为[]interface{}类型,如果转换失败，返回nil
func (p *Json) ToArray() []interface{} {
    pointer := p.value
    switch (*pointer).(type) {
        case []interface{}:
            return (*pointer).([]interface{})
        default:
            return nil
    }
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