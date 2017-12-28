package gxml

import (
    "fmt"
    "strings"
    "strconv"
    "io/ioutil"
    "encoding/xml"
)

// xml解析结果存放数组
type Xml struct {
    // 注意这是一个指针
    value *interface{}
}

// 一个xml变量
type XmlVar interface{}

// 编码go变量为xml字符串，并返回xml字符串指针
func Encode (v interface{}) ([]byte, error) {
    return xml.Marshal(v)
}

// 解码字符串为interface{}变量
func Decode (b []byte) (interface{}, error) {
    var v interface{}
    if err := DecodeTo(b, &v); err == nil {
        return nil, err
    } else {
        return v, nil
    }
}

// 解析xml字符串为go变量，注意第二个参数为指针
func DecodeTo (b []byte, v interface{}) error {
    return xml.Unmarshal(b, v)
}

// 解析xml字符串为gxml.Xml对象，并返回操作对象指针
func DecodeToXml (b []byte) (*Xml, error) {
    if v, err := Decode(b); err != nil {
        return &Xml{&v}, nil
    } else {
        return nil, err
    }
}

// 加载xml文件内容，并转换为xml对象
func Load (path string) (*Xml, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var result interface{}
    if err := xml.Unmarshal(data, &result); err != nil {
        return nil, err
    }
    return &Xml{ &result }, nil
}

// 将变量转换为Xml对象进行处理，该变量至少应当是一个map或者array，否者转换没有意义
func NewXml(v *interface{}) *Xml {
    return &Xml{ v }
}

// 将指定的xml内容转换为指定结构返回，查找失败或者转换失败，目标对象转换为nil
// 注意第二个参数需要给的是变量地址
func (p *Xml) GetToVar(pattern string, v interface{}) error {
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
// 注意，如果获取的值不存在，或者类型与xml类型不匹配，那么将会返回nil
func (p *Xml) GetMap(pattern string) map[string]interface{} {
    result := p.Get(pattern)
    if result != nil {
        if r, ok := result.(map[string]interface{}); ok {
            return r
        }
    }
    return nil
}

// 将检索值转换为Xml对象指针返回
func (p *Xml) GetXml(pattern string) *Xml {
    result := p.Get(pattern)
    if result != nil {
        return &Xml{&result}
    }
    return nil
}

// 获得一个数组[]interface{}，方便操作，不需要自己做类型转换
// 注意，如果获取的值不存在，或者类型与xml类型不匹配，那么将会返回nil
func (p *Xml) GetArray(pattern string) []interface{} {
    result := p.Get(pattern)
    if result != nil {
        if r, ok := result.([]interface{}); ok {
            return r
        }
    }
    return nil
}

// 返回指定xml中的string
func (p *Xml) GetString(pattern string) string {
    result := p.Get(pattern)
    if result != nil {
        if r, ok := result.(string); ok {
            return r
        }
    }
    return ""
}

// 返回指定xml中的bool
func (p *Xml) GetBool(pattern string) bool {
    result := p.Get(pattern)
    if result != nil {
        str := fmt.Sprintf("%v", result)
        if str != "" && str != "0" && str != "false" {
            return true
        }
    }
    return false
}

// 返回指定xml中的float64
func (p *Xml) GetFloat64(pattern string) float64 {
    result := p.Get(pattern)
    if result != nil {
        if r, ok := result.(float64); ok {
            return r
        }
    }
    return 0
}

// 返回指定xml中的float64->int
func (p *Xml) GetInt(pattern string) int {
    return int(p.GetFloat64(pattern))
}

// 返回指定xml中的float64->int64
func (p *Xml) GetInt64(pattern string) int64 {
    return int64(p.GetFloat64(pattern))
}

// 根据约定字符串方式访问xml解析数据，参数形如： "items.name.first", "list.0"
// 返回的结果类型的interface{}，因此需要自己做类型转换
// 如果找不到对应节点的数据，返回nil
func (p *Xml) Get(pattern string) interface{} {
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
func (p *Xml) ToMap() map[string]interface{} {
    pointer := p.value
    switch (*pointer).(type) {
        case map[string]interface{}:
            return (*pointer).(map[string]interface{})
        default:
            return nil
    }
}

// 转换为[]interface{}类型,如果转换失败，返回nil
func (p *Xml) ToArray() []interface{} {
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