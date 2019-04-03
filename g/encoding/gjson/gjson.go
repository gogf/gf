// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gjson provides convenient API for JSON/XML/YAML/TOML data handling.
package gjson

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/gogf/gf/g/encoding/gtoml"
    "github.com/gogf/gf/g/encoding/gxml"
    "github.com/gogf/gf/g/encoding/gyaml"
    "github.com/gogf/gf/g/internal/rwmutex"
    "github.com/gogf/gf/g/os/gfcache"
    "github.com/gogf/gf/g/text/gregex"
    "github.com/gogf/gf/g/text/gstr"
    "github.com/gogf/gf/g/util/gconv"
    "reflect"
    "strconv"
    "strings"
    "time"
)

const (
    // Separator char for hierarchical data access.
    gDEFAULT_SPLIT_CHAR = '.'
)

// The customized JSON struct.
type Json struct {
    mu *rwmutex.RWMutex
    p  *interface{} // Pointer for hierarchical data access, it's the root of data in default.
    c  byte         // Char separator('.' in default).
    vc bool         // Violence Check(false in default), which is used to access data
                    // when the hierarchical data key contains separator char.
}

// New creates a Json object with any variable type of <data>,
// but <data> should be a map or slice for data access reason,
// or it will make no sense.
// The <unsafe> param specifies whether using this Json object
// in un-concurrent-safe context, which is false in default.
func New(data interface{}, unsafe...bool) *Json {
    j := (*Json)(nil)
    switch data.(type) {
        case map[string]interface{}, []interface{}, nil:
            j = &Json {
                p  : &data,
                c  : byte(gDEFAULT_SPLIT_CHAR),
                vc : false ,
            }
        case string, []byte:
            j, _ = LoadContent(gconv.Bytes(data))
        default:
            v := (interface{})(nil)
            if m := gconv.Map(data); m != nil {
                v = m
                j = &Json {
                    p  : &v,
                    c  : byte(gDEFAULT_SPLIT_CHAR),
                    vc : false,
                }
            } else {
                v = gconv.Interfaces(data)
                j = &Json {
                    p  : &v,
                    c  : byte(gDEFAULT_SPLIT_CHAR),
                    vc : false,
                }
            }
    }
    j.mu = rwmutex.New(unsafe...)
    return j
}

// NewUnsafe creates a un-concurrent-safe Json object.
func NewUnsafe(data...interface{}) *Json {
    if len(data) > 0 {
        return New(data[0], true)
    }
    return New(nil, true)
}

// Valid checks whether <data> is a valid JSON data type.
func Valid(data interface{}) bool {
    return json.Valid(gconv.Bytes(data))
}

// Encode encodes <value> to JSON data type of bytes.
func Encode(value interface{}) ([]byte, error) {
    return json.Marshal(value)
}

// Decode decodes <data>(string/[]byte) to golang variable.
func Decode(data interface{}) (interface{}, error) {
    var value interface{}
    if err := DecodeTo(gconv.Bytes(data), &value); err != nil {
        return nil, err
    } else {
        return value, nil
    }
}

// Decode decodes <data>(string/[]byte) to specified golang variable <v>.
// The <v> should be a pointer type.
func DecodeTo(data interface{}, v interface{}) error {
    decoder := json.NewDecoder(bytes.NewReader(gconv.Bytes(data)))
    decoder.UseNumber()
    return decoder.Decode(v)
}

// DecodeToJson codes <data>(string/[]byte) to a Json object.
func DecodeToJson(data interface{}, unsafe...bool) (*Json, error) {
    if v, err := Decode(gconv.Bytes(data)); err != nil {
        return nil, err
    } else {
        return New(v, unsafe...), nil
    }
}

// Load loads content from specified file <path>,
// and creates a Json object from its content.
func Load(path string, unsafe...bool) (*Json, error) {
    return LoadContent(gfcache.GetBinContents(path), unsafe...)
}

// LoadContent creates a Json object from given content,
// it checks the data type of <content> automatically,
// supporting JSON, XML, YAML and TOML types of data.
func LoadContent(data interface{}, unsafe...bool) (*Json, error) {
    var err    error
    var result interface{}
    b := gconv.Bytes(data)
    t := "json"
    // auto check data type
    if json.Valid(b) {
        t = "json"
    } else if gregex.IsMatch(`^<.+>.*</.+>$`, b) {
        t = "xml"
    } else if gregex.IsMatch(`^[\s\t]*\w+\s*:\s*.+`, b) || gregex.IsMatch(`\n[\s\t]*\w+\s*:\s*.+`, b) {
        t = "yml"
    } else if gregex.IsMatch(`^[\s\t]*\w+\s*=\s*.+`, b) || gregex.IsMatch(`\n[\s\t]*\w+\s*=\s*.+`, b) {
        t = "toml"
    } else {
        return nil, errors.New("unsupported data type")
    }
    // convert to json type data
    switch t {
        case "json", ".json":
            // ok
        case "xml", ".xml":
            // TODO UseNumber
            b, err = gxml.ToJson(b)

        case "yml", "yaml", ".yml", ".yaml":
            // TODO UseNumber
            b, err = gyaml.ToJson(b)

        case "toml", ".toml":
            // TODO UseNumber
            b, err = gtoml.ToJson(b)

        default:
            err = errors.New("nonsupport type " + t)
    }
    if err != nil {
        return nil, err
    }
    if result == nil {
        decoder := json.NewDecoder(bytes.NewReader(b))
        decoder.UseNumber()
        if err := decoder.Decode(&result); err != nil {
            return nil, err
        }
        switch result.(type) {
            case string, []byte:
                return nil, fmt.Errorf(`json decoding failed for content: %s`, string(b))
        }
    }
    return New(result, unsafe...), nil
}

// SetSplitChar sets the separator char for hierarchical data access.
func (j *Json) SetSplitChar(char byte) {
    j.mu.Lock()
    j.c = char
    j.mu.Unlock()
}

// SetViolenceCheck enables/disables violence check for hierarchical data access.
func (j *Json) SetViolenceCheck(enabled bool) {
    j.mu.Lock()
    j.vc = enabled
    j.mu.Unlock()
}

// GetToVar gets the value by specified <pattern>,
// and converts it to specified golang variable <v>.
// The <v> should be a pointer type.
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

// GetMap gets the value by specified <pattern>,
// and converts it to map[string]interface{}.
func (j *Json) GetMap(pattern string) map[string]interface{} {
    result := j.Get(pattern)
    if result != nil {
        return gconv.Map(result)
    }
    return nil
}

// GetJson gets the value by specified <pattern>,
// and converts it to a Json object.
func (j *Json) GetJson(pattern string) *Json {
    result := j.Get(pattern)
    if result != nil {
        return New(result)
    }
    return nil
}

// GetJsons gets the value by specified <pattern>,
// and converts it to a slice of Json object.
func (j *Json) GetJsons(pattern string) []*Json {
    array := j.GetArray(pattern)
    if len(array) > 0 {
        jsons := make([]*Json, len(array))
        for i := 0; i < len(array); i++ {
            jsons[i] = New(array[i], !j.mu.IsSafe())
        }
        return jsons
    }
    return nil
}

// GetArray gets the value by specified <pattern>,
// and converts it to a slice of []interface{}.
func (j *Json) GetArray(pattern string) []interface{} {
    return gconv.Interfaces(j.Get(pattern))
}

// GetString gets the value by specified <pattern>,
// and converts it to string.
func (j *Json) GetString(pattern string) string {
    return gconv.String(j.Get(pattern))
}

// GetStrings gets the value by specified <pattern>,
// and converts it to a slice of []string.
func (j *Json) GetStrings(pattern string) []string {
    return gconv.Strings(j.Get(pattern))
}

// See GetArray.
func (j *Json) GetInterfaces(pattern string) []interface{} {
    return gconv.Interfaces(j.Get(pattern))
}

func (j *Json) GetTime(pattern string, format ... string) time.Time {
    return gconv.Time(j.Get(pattern), format...)
}

func (j *Json) GetTimeDuration(pattern string) time.Duration {
    return gconv.TimeDuration(j.Get(pattern))
}

// GetBool gets the value by specified <pattern>,
// and converts it to bool.
// It returns false when value is: "", 0, false, off, nil;
// or returns true instead.
func (j *Json) GetBool(pattern string) bool {
    return gconv.Bool(j.Get(pattern))
}

func (j *Json) GetInt(pattern string) int {
    return gconv.Int(j.Get(pattern))
}

func (j *Json) GetInt8(pattern string) int8 {
    return gconv.Int8(j.Get(pattern))
}

func (j *Json) GetInt16(pattern string) int16 {
    return gconv.Int16(j.Get(pattern))
}

func (j *Json) GetInt32(pattern string) int32 {
    return gconv.Int32(j.Get(pattern))
}

func (j *Json) GetInt64(pattern string) int64 {
    return gconv.Int64(j.Get(pattern))
}

func (j *Json) GetInts(pattern string) []int {
    return gconv.Ints(j.Get(pattern))
}

func (j *Json) GetUint(pattern string) uint {
    return gconv.Uint(j.Get(pattern))
}

func (j *Json) GetUint8(pattern string) uint8 {
    return gconv.Uint8(j.Get(pattern))
}

func (j *Json) GetUint16(pattern string) uint16 {
    return gconv.Uint16(j.Get(pattern))
}

func (j *Json) GetUint32(pattern string) uint32 {
    return gconv.Uint32(j.Get(pattern))
}

func (j *Json) GetUint64(pattern string) uint64 {
    return gconv.Uint64(j.Get(pattern))
}

func (j *Json) GetFloat32(pattern string) float32 {
    return gconv.Float32(j.Get(pattern))
}

func (j *Json) GetFloat64(pattern string) float64 {
    return gconv.Float64(j.Get(pattern))
}

func (j *Json) GetFloats(pattern string) []float64 {
    return gconv.Floats(j.Get(pattern))
}

// GetToStruct gets the value by specified <pattern>,
// and converts it to specified object <objPointer>.
// The <objPointer> should be the pointer to an object.
func (j *Json) GetToStruct(pattern string, objPointer interface{}) error {
    return gconv.Struct(j.Get(pattern), objPointer)
}

// Set sets value with specified <pattern>.
// It supports hierarchical data access by char separator, which is '.' in default.
func (j *Json) Set(pattern string, value interface{}) error {
    return j.setValue(pattern, value, false)
}

// Remove deletes value with specified <pattern>.
// It supports hierarchical data access by char separator, which is '.' in default.
func (j *Json) Remove(pattern string) error {
    return j.setValue(pattern, nil, true)
}

// Set <value> by <pattern>.
// Notice:
// 1. If value is nil and removed is true, means deleting this value;
// 2. It's quite complicated in hierarchical data search, node creating and data assignment;
func (j *Json) setValue(pattern string, value interface{}, removed bool) error {
    array   := strings.Split(pattern, string(j.c))
    length  := len(array)
    value    = j.convertValue(value)
    // 初始化判断
    if *j.p == nil {
        if gstr.IsNumeric(array[0]) {
            *j.p = make([]interface{}, 0)
        } else {
            *j.p = make(map[string]interface{})
        }
    }
    var pparent *interface{} = nil // 父级元素项(设置时需要根据子级的内容确定数据类型，所以必须记录父级)
    var pointer *interface{} = j.p // 当前操作层级项
    j.mu.Lock()
    defer j.mu.Unlock()
    for i:= 0; i < length; i++ {
        switch (*pointer).(type) {
            case map[string]interface{}:
                if i == length - 1 {
                    if removed && value == nil {
                        // 删除map元素
                        delete((*pointer).(map[string]interface{}), array[i])
                    } else {
                        (*pointer).(map[string]interface{})[array[i]] = value
                    }
                } else {
                    // 当键名不存在的情况这里会进行处理
                    if v, ok := (*pointer).(map[string]interface{})[array[i]]; !ok {
                        if removed && value == nil {
                            goto done
                        }
                        // 创建新节点
                        if gstr.IsNumeric(array[i + 1]) {
                            // 创建array节点
                            n, _ := strconv.Atoi(array[i + 1])
                            var v interface{} = make([]interface{}, n + 1)
                            pparent = j.setPointerWithValue(pointer, array[i], v)
                            pointer = &v
                        } else {
                            // 创建map节点
                            var v interface{} = make(map[string]interface{})
                            pparent = j.setPointerWithValue(pointer, array[i], v)
                            pointer = &v
                        }
                    } else {
                        pparent = pointer
                        pointer = &v
                    }
                }

            case []interface{}:
                // 键名与当前指针类型不符合，需要执行**覆盖操作**
                if !gstr.IsNumeric(array[i]) {
                    if i == length - 1 {
                        *pointer = map[string]interface{}{ array[i] : value }
                    } else {
                        var v interface{} = make(map[string]interface{})
                        *pointer = v
                        pparent  = pointer
                        pointer  = &v
                    }
                    continue
                }

                valn, err := strconv.Atoi(array[i])
                if err != nil {
                    return err
                }
                // 叶子节点
                if i == length - 1 {
                    if len((*pointer).([]interface{})) > valn {
                        if removed && value == nil {
                            // 删除数据元素
                            j.setPointerWithValue(pparent, array[i - 1], append((*pointer).([]interface{})[ : valn], (*pointer).([]interface{})[valn + 1 : ]...))
                        } else {
                            (*pointer).([]interface{})[valn] = value
                        }
                    } else {
                        if removed && value == nil {
                            goto done
                        }
                        if pparent == nil {
                            // 表示根节点
                            j.setPointerWithValue(pointer, array[i], value)
                        } else {
                            // 非根节点
                            s := make([]interface{}, valn + 1)
                            copy(s, (*pointer).([]interface{}))
                            s[valn] = value
                            j.setPointerWithValue(pparent, array[i - 1], s)
                        }
                    }
                } else {
                    if gstr.IsNumeric(array[i + 1]) {
                        n, _ := strconv.Atoi(array[i + 1])
                        if len((*pointer).([]interface{})) > valn {
                            (*pointer).([]interface{})[valn] = make([]interface{}, n + 1)
                            pparent                          = pointer
                            pointer                          = &(*pointer).([]interface{})[valn]
                        } else {
                            if removed && value == nil {
                                goto done
                            }
                            var v interface{} = make([]interface{}, n + 1)
                            pparent = j.setPointerWithValue(pointer, array[i], v)
                            pointer = &v
                        }
                    } else {
                        var v interface{} = make(map[string]interface{})
                        pparent = j.setPointerWithValue(pointer, array[i], v)
                        pointer = &v
                    }
                }

            // 如果当前指针指向的变量不是引用类型的，
            // 那么修改变量必须通过父级进行修改，即 pparent
            default:
                if removed && value == nil {
                    goto done
                }
                if gstr.IsNumeric(array[i]) {
                    n, _    := strconv.Atoi(array[i])
                    s       := make([]interface{}, n + 1)
                    if i == length - 1 {
                        s[n] = value
                    }
                    if pparent != nil {
                        pparent = j.setPointerWithValue(pparent, array[i - 1], s)
                    } else {
                        *pointer = s
                        pparent  = pointer
                    }
                } else {
                    var v interface{} = make(map[string]interface{})
                    if i == length - 1 {
                        v = map[string]interface{}{
                            array[i] : value,
                        }
                    }
                    if pparent != nil {
                        pparent = j.setPointerWithValue(pparent, array[i - 1], v)
                    } else {
                        *pointer = v
                        pparent  = pointer
                    }
                    pointer = &v
                }
        }
    }
done:
    return nil
}

// Convert <value> to map[string]interface{} or []interface{},
// which can be supported for hierarchical data access.
func (j *Json) convertValue(value interface{}) interface{} {
    switch value.(type) {
        case map[string]interface{}:
            return value
        case []interface{}:
            return value
        default:
            rv   := reflect.ValueOf(value)
            kind := rv.Kind()
            if kind == reflect.Ptr {
                rv   = rv.Elem()
                kind = rv.Kind()
            }
            switch kind {
                case reflect.Array:  return gconv.Interfaces(value)
                case reflect.Slice:  return gconv.Interfaces(value)
                case reflect.Map:    return gconv.Map(value)
                case reflect.Struct: return gconv.Map(value)
                default:
                    // Use json decode/encode at last.
                    b, _ := Encode(value)
                    v, _ := Decode(b)
                    return v
            }
    }
}

// Set <key>:<value> to <pointer>, the <key> may be a map key or slice index.
// It returns the pointer to the new value set.
func (j *Json) setPointerWithValue(pointer *interface{}, key string, value interface{}) *interface{} {
    switch (*pointer).(type) {
        case map[string]interface{}:
            (*pointer).(map[string]interface{})[key] = value
            return &value
        case []interface{}:
            n, _ := strconv.Atoi(key)
            if len((*pointer).([]interface{})) > n {
                (*pointer).([]interface{})[n] = value
                return &(*pointer).([]interface{})[n]
            } else {
                s := make([]interface{}, n + 1)
                copy(s, (*pointer).([]interface{}))
                s[n] = value
                *pointer = s
                return &s[n]
            }
        default:
            *pointer = value
    }
    return pointer
}

// Get returns value by specified <pattern>.
// It returns all values of current Json object, if <pattern> is empty or not specified.
// It returns nil if no value found by <pattern>.
//
// We can also access slice item by its index number in <pattern>,
// eg: "items.name.first", "list.10".
func (j *Json) Get(pattern...string) interface{} {
    j.mu.RLock()
    defer j.mu.RUnlock()

    queryPattern := ""
    if len(pattern) > 0 {
        queryPattern = pattern[0]
    }
    var result *interface{}
    if j.vc {
        result = j.getPointerByPattern(queryPattern)
    } else {
        result = j.getPointerByPatternWithoutViolenceCheck(queryPattern)
    }
    if result != nil {
        return *result
    }
    return nil
}

// Contains checks whether the value by specified <pattern> exist.
func (j *Json) Contains(pattern...string) bool {
    return j.Get(pattern...) != nil
}

// Len returns the length/size of the value by specified <pattern>.
// The target value by <pattern> should be type of slice or map.
// It returns -1 if the target value is not found, or its type is invalid.
func (j *Json) Len(pattern string) int {
    p := j.getPointerByPattern(pattern)
    if p != nil {
        switch (*p).(type) {
            case map[string]interface{}:
                return len((*p).(map[string]interface{}))
            case []interface{}:
                return len((*p).([]interface{}))
            default:
                return -1
        }
    }
    return -1
}

// Append appends value to the value by specified <pattern>.
// The target value by <pattern> should be type of slice.
func (j *Json) Append(pattern string, value interface{}) error {
    p := j.getPointerByPattern(pattern)
    if p == nil {
        return j.Set(fmt.Sprintf("%s.0", pattern), value)
    }
    switch (*p).(type) {
        case []interface{}:
            return j.Set(fmt.Sprintf("%s.%d", pattern, len((*p).([]interface{}))), value)
    }
    return fmt.Errorf("invalid variable type of %s", pattern)
}

// Get a pointer to the value by specified <pattern>.
func (j *Json) getPointerByPattern(pattern string) *interface{} {
    if j.vc {
        return j.getPointerByPatternWithViolenceCheck(pattern)
    } else {
        return j.getPointerByPatternWithoutViolenceCheck(pattern)
    }
}

// Get a pointer to the value of specified <pattern> with violence check.
func (j *Json) getPointerByPatternWithViolenceCheck(pattern string) *interface{} {
    if !j.vc {
        return j.getPointerByPatternWithoutViolenceCheck(pattern)
    }
    index   := len(pattern)
    start   := 0
    length  := 0
    pointer := j.p
    if index == 0 {
        return pointer
    }
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
            // Get the position for next separator char.
            index = strings.LastIndexByte(pattern[start:index], j.c)
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

// Get a pointer to the value of specified <pattern>, with no violence check.
func (j *Json) getPointerByPatternWithoutViolenceCheck(pattern string) *interface{} {
    if j.vc {
        return j.getPointerByPatternWithViolenceCheck(pattern)
    }
    pointer := j.p
    if len(pattern) == 0 {
        return pointer
    }
    array := strings.Split(pattern, string(j.c))
    for k, v := range array {
        if r := j.checkPatternByPointer(v, pointer); r != nil {
            if k == len(array) - 1 {
                return r
            } else {
                pointer = r
            }
        } else {
            break
        }
    }
    return nil
}

// Check whether there's value by <key> in specified <pointer>.
// It returns a pointer to the value.
func (j *Json) checkPatternByPointer(key string, pointer *interface{}) *interface{} {
    switch (*pointer).(type) {
        case map[string]interface{}:
            if v, ok := (*pointer).(map[string]interface{})[key]; ok {
                return &v
            }
        case []interface{}:
            if gstr.IsNumeric(key) {
                n, err := strconv.Atoi(key)
                if err == nil && len((*pointer).([]interface{})) > n {
                    return &(*pointer).([]interface{})[n]
                }
            }
    }
    return nil
}

// ToMap converts current Json object to map[string]interface{}.
// It returns nil if fails.
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

// ToArray converts current Json object to []interface{}.
// It returns nil if fails.
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

func (j *Json) ToXmlString(rootTag...string) (string, error) {
    b, e := j.ToXml(rootTag...)
    return string(b), e
}

func (j *Json) ToXmlIndent(rootTag...string) ([]byte, error) {
    return gxml.EncodeWithIndent(j.ToMap(), rootTag...)
}

func (j *Json) ToXmlIndentString(rootTag...string) (string, error) {
    b, e := j.ToXmlIndent(rootTag...)
    return string(b), e
}

func (j *Json) ToJson() ([]byte, error) {
    j.mu.RLock()
    defer j.mu.RUnlock()
    return Encode(*(j.p))
}

func (j *Json) ToJsonString() (string, error) {
    b, e := j.ToJson()
    return string(b), e
}

func (j *Json) ToJsonIndent() ([]byte, error) {
    j.mu.RLock()
    defer j.mu.RUnlock()
    return json.MarshalIndent(*(j.p), "", "\t")
}

func (j *Json) ToJsonIndentString() (string, error) {
    b, e := j.ToJsonIndent()
    return string(b), e
}

func (j *Json) ToYaml() ([]byte, error) {
    j.mu.RLock()
    defer j.mu.RUnlock()
    return gyaml.Encode(*(j.p))
}

func (j *Json) ToYamlString() (string, error) {
    b, e := j.ToYaml()
    return string(b), e
}

func (j *Json) ToToml() ([]byte, error) {
    j.mu.RLock()
    defer j.mu.RUnlock()
    return gtoml.Encode(*(j.p))
}

func (j *Json) ToTomlString() (string, error) {
    b, e := j.ToToml()
    return string(b), e
}

// ToStruct converts current Json object to specified object.
// The <objPointer> should be a pointer type.
func (j *Json) ToStruct(objPointer interface{}) error {
    j.mu.RLock()
    defer j.mu.RUnlock()
    return gconv.Struct(*(j.p), objPointer)
}

// Dump prints current Json object with more manually readable.
func (j *Json) Dump() error {
    j.mu.RLock()
    defer j.mu.RUnlock()
    if b, err := j.ToJsonIndent(); err != nil {
        return err
    } else {
        fmt.Println(string(b))
    }
    return nil
}