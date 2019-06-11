<<<<<<< HEAD
// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
=======
// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
>>>>>>> upstream/master

package ghttp

import (
<<<<<<< HEAD
    "io/ioutil"
    "net/http"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/encoding/gjson"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/util/gregx"
    "gitee.com/johng/gf/g/os/gtime"
=======
	"fmt"
	"github.com/gogf/gf/g/container/gvar"
    "github.com/gogf/gf/g/encoding/gjson"
    "github.com/gogf/gf/g/os/gtime"
    "github.com/gogf/gf/g/text/gregex"
    "github.com/gogf/gf/third/github.com/fatih/structs"
    "io/ioutil"
    "net/http"
    "strings"
>>>>>>> upstream/master
)

// 请求对象
type Request struct {
<<<<<<< HEAD
    http.Request
    parsedGet  *gtype.Bool         // GET参数是否已经解析
    parsedPost *gtype.Bool         // POST参数是否已经解析
    values     map[string][]string // GET参数
    exit       *gtype.Bool         // 是否退出当前请求流程执行
    Id         int                 // 请求id(唯一)
    Server     *Server             // 请求关联的服务器对象
    Cookie     *Cookie             // 与当前请求绑定的Cookie对象(并发安全)
    Session    *Session            // 与当前请求绑定的Session对象(并发安全)
    Response   *Response           // 对应请求的返回数据操作对象
    Router     *Router             // 匹配到的路由对象
    EnterTime  int64               // 请求进入时间(微秒)
    LeaveTime  int64               // 请求完成时间(微秒)
    Param      interface{}         // 开发者自定义参数
=======
    *http.Request
    parsedGet     bool                    // GET参数是否已经解析
    parsedPost    bool                    // POST参数是否已经解析
    queryVars     map[string][]string     // GET参数
    routerVars    map[string][]string     // 路由解析参数
    exit          bool                    // 是否退出当前请求流程执行
    Id            int                     // 请求id(唯一)
    Server        *Server                 // 请求关联的服务器对象
    Cookie        *Cookie                 // 与当前请求绑定的Cookie对象(并发安全)
    Session       *Session                // 与当前请求绑定的Session对象(并发安全)
    Response      *Response               // 对应请求的返回数据操作对象
    Router        *Router                 // 匹配到的路由对象
    EnterTime     int64                   // 请求进入时间(微秒)
    LeaveTime     int64                   // 请求完成时间(微秒)
    params        map[string]interface{}  // 开发者自定义参数(请求流程中有效)
    parsedHost    string                  // 解析过后不带端口号的服务器域名名称
    clientIp      string                  // 解析过后的客户端IP地址
    rawContent    []byte                  // 客户端提交的原始参数
    isFileRequest bool                    // 是否为静态文件请求(非服务请求，当静态文件存在时，优先级会被服务请求高，被识别为文件请求)
>>>>>>> upstream/master
}

// 创建一个Request对象
func newRequest(s *Server, r *http.Request, w http.ResponseWriter) *Request {
<<<<<<< HEAD
    request := &Request{
        parsedGet  : gtype.NewBool(),
        parsedPost : gtype.NewBool(),
        values     : make(map[string][]string),
        exit       : gtype.NewBool(),
        Id         : s.servedCount.Add(1),
        Server     : s,
        Request    : *r,
        Response   : newResponse(w),
=======
    request := &Request {
        routerVars : make(map[string][]string),
        Id         : s.servedCount.Add(1),
        Server     : s,
        Request    : r,
        Response   : newResponse(s, w),
>>>>>>> upstream/master
        EnterTime  : gtime.Microsecond(),
    }
    // 会话处理
    request.Cookie           = GetCookie(request)
    request.Session          = GetSession(request)
    request.Response.request = request
    return request
}

<<<<<<< HEAD

// 初始化GET请求参数
func (r *Request) initGet() {
    if !r.parsedGet.Val() {
        if len(r.values) == 0 {
            r.values = r.URL.Query()
        } else {
            for k, v := range r.URL.Query() {
                r.values[k] = v
            }
        }
    }
}

// 初始化POST请求参数
func (r *Request) initPost() {
    if !r.parsedPost.Val() {
        // 快速保存，尽量避免并发问题
        r.parsedPost.Set(true)
        // MultiMedia表单请求解析允许最大使用内存：1GB
        r.ParseMultipartForm(1024*1024*1024)
    }
}

// 获得指定名称的参数字符串(GET/POST)，同 GetRequestString
// 这是常用方法的简化别名
func (r *Request) Get(k string) string {
    return r.GetRequestString(k)
}

// 获得指定名称的get参数列表
func (r *Request) GetQuery(k string) []string {
    r.initGet()
    if v, ok := r.values[k]; ok {
        return v
    }
    return nil
}

func (r *Request) GetQueryBool(k string) bool {
    return gconv.Bool(r.Get(k))
}

func (r *Request) GetQueryInt(k string) int {
    return gconv.Int(r.Get(k))
}

func (r *Request) GetQueryUint(k string) uint {
    return gconv.Uint(r.Get(k))
}

func (r *Request) GetQueryFloat32(k string) float32 {
    return gconv.Float32(r.Get(k))
}

func (r *Request) GetQueryFloat64(k string) float64 {
    return gconv.Float64(r.Get(k))
}

func (r *Request) GetQueryString(k string) string {
    v := r.GetQuery(k)
    if v == nil {
        return ""
    } else {
        return v[0]
    }
}

func (r *Request) GetQueryArray(k string) []string {
    return r.GetQuery(k)
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值
func (r *Request) GetQueryMap(defaultMap...map[string]string) map[string]string {
    r.initGet()
    m := make(map[string]string)
    if len(defaultMap) == 0 {
        for k, v := range r.values {
            m[k] = v[0]
        }
    } else {
        for k, v := range defaultMap[0] {
            v2 := r.GetQueryArray(k)
            if v2 == nil {
                m[k] = v
            } else {
                m[k] = v2[0]
            }
        }
    }
    return m
}

// 获得post参数
func (r *Request) GetPost(k string) []string {
    r.initPost()
    if v, ok := r.PostForm[k]; ok {
        return v
    }
    return nil
}

func (r *Request) GetPostBool(k string) bool {
    return gconv.Bool(r.GetPostString(k))
}

func (r *Request) GetPostInt(k string) int {
    return gconv.Int(r.GetPostString(k))
}

func (r *Request) GetPostUint(k string) uint {
    return gconv.Uint(r.GetPostString(k))
}

func (r *Request) GetPostFloat32(k string) float32 {
    return gconv.Float32(r.GetPostString(k))
}

func (r *Request) GetPostFloat64(k string) float64 {
    return gconv.Float64(r.GetPostString(k))
}

func (r *Request) GetPostString(k string) string {
    v := r.GetPost(k)
    if v == nil {
        return ""
    } else {
        return v[0]
    }
}

func (r *Request) GetPostArray(k string) []string {
    return r.GetPost(k)
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值
// 需要注意的是，如果其中一个字段为数组形式，那么只会返回第一个元素，如果需要获取全部的元素，请使用GetPostArray获取特定字段内容
func (r *Request) GetPostMap(defaultMap...map[string]string) map[string]string {
    r.initPost()
    m := make(map[string]string)
    if len(defaultMap) == 0 {
        for k, v := range r.PostForm {
            m[k] = v[0]
        }
    } else {
        for k, v := range defaultMap[0] {
            if v2, ok := r.PostForm[k]; ok {
                m[k] = v2[0]
            } else {
                m[k] = v
            }
        }
    }
    return m
}

// 获得post或者get提交的参数，如果有同名参数，那么按照get->post优先级进行覆盖
func (r *Request) GetRequest(k string) []string {
    v := r.GetQuery(k)
    if v == nil {
        return r.GetPost(k)
    }
    return v
}

func (r *Request) GetRequestString(k string) string {
    v := r.GetRequest(k)
    if v == nil {
        return ""
    } else {
        return v[0]
    }
}

func (r *Request) GetRequestBool(k string) bool {
    return gconv.Bool(r.GetRequestString(k))
}

func (r *Request) GetRequestInt(k string) int {
    return gconv.Int(r.GetRequestString(k))
}

func (r *Request) GetRequestUint(k string) uint {
    return gconv.Uint(r.GetRequestString(k))
}

func (r *Request) GetRequestFloat32(k string) float32 {
    return gconv.Float32(r.GetRequestString(k))
}

func (r *Request) GetRequestFloat64(k string) float64 {
    return gconv.Float64(r.GetRequestString(k))
}

func (r *Request) GetRequestArray(k string) []string {
    return r.GetRequest(k)
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值
// 需要注意的是，如果其中一个字段为数组形式，那么只会返回第一个元素，如果需要获取全部的元素，请使用GetRequestArray获取特定字段内容
func (r *Request) GetRequestMap(defaultMap...map[string]string) map[string]string {
    m := r.GetQueryMap()
    if len(defaultMap) == 0 {
        for k, v := range r.GetPostMap() {
            if _, ok := m[k]; !ok {
                m[k] = v
            }
        }
    } else {
        for k, v := range defaultMap[0] {
            v2 := r.GetRequest(k)
            if v2 != nil {
                m[k] = v2[0]
            } else {
                m[k] = v
            }
        }
    }
    return m
}

// 获取原始请求输入字符串，注意：只能获取一次，读完就没了
func (r *Request) GetRaw() []byte {
    result, _ := ioutil.ReadAll(r.Body)
    return result
}

// 获取原始json请求输入字符串，并解析为json对象
func (r *Request) GetJson() *gjson.Json {
    data := r.GetRaw()
    if data != nil {
        if j, err := gjson.DecodeToJson(data); err == nil {
            return j
        }
    }
    return nil
}

// 退出当前请求执行，原理是在Request.exit做标记，由服务逻辑流程做判断，自行停止
func (r *Request) Exit() {
    r.exit.Set(true)
}

// 判断当前请求是否停止执行
func (r *Request) IsExited() bool {
    return r.exit.Val()
}

// 获取请求的服务端IP/域名
func (r *Request) GetHost() string {
    array, _ := gregx.MatchString(`(.+):(\d+)`, r.Host)
    if len(array) > 1 {
        return array[1]
    }
    return r.Host
}

// 获取请求的客户端IP地址
func (r *Request) GetClientIp() string {
    array, _ := gregx.MatchString(`(.+):(\d+)`, r.RemoteAddr)
    if len(array) > 1 {
        return array[1]
    }
    return r.RemoteAddr
}

=======
// 获取Web Socket连接对象(如果是非WS请求会失败，注意检查返回的error结果)
func (r *Request) WebSocket() (*WebSocket, error) {
    if conn, err := wsUpgrader.Upgrade(r.Response.ResponseWriter.ResponseWriter, r.Request, nil); err == nil {
        return &WebSocket {
            conn,
        }, nil
    } else {
        return nil, err
    }
}

// 获得指定名称的参数字符串(Router/GET/POST)，同 GetRequestString
// 这是常用方法的简化别名
func (r *Request) Get(key string, def...interface{}) string {
    return r.GetRequestString(key, def...)
}

// 建议都用该参数替代参数获取
func (r *Request) GetVar(key string, def...interface{}) *gvar.Var {
    return r.GetRequestVar(key, def...)
}

// 获取原始请求输入二进制。
func (r *Request) GetRaw() []byte {
    err := error(nil)
    if r.rawContent == nil {
        r.rawContent, err = ioutil.ReadAll(r.Body)
        if err != nil {
            r.Error("error reading request body: ", err)
        }
    }
    return r.rawContent
}

// 获取原始请求输入字符串。
func (r *Request) GetRawString() string {
    return string(r.GetRaw())
}

// 获取原始json请求输入字符串，并解析为json对象
func (r *Request) GetJson() *gjson.Json {
    data := r.GetRaw()
    if len(data) > 0 {
        if j, err := gjson.DecodeToJson(data); err == nil {
            return j
        } else {
            r.Error(err, ": ", string(data))
        }
    }
    return nil
}

func (r *Request) GetString(key string, def...interface{}) string {
    return r.GetRequestString(key, def...)
}

func (r *Request) GetInt(key string, def...interface{}) int {
    return r.GetRequestInt(key, def...)
}

func (r *Request) GetInts(key string, def...interface{}) []int {
    return r.GetRequestInts(key, def...)
}

func (r *Request) GetUint(key string, def...interface{}) uint {
    return r.GetRequestUint(key, def...)
}

func (r *Request) GetFloat32(key string, def...interface{}) float32 {
    return r.GetRequestFloat32(key, def...)
}

func (r *Request) GetFloat64(key string, def...interface{}) float64 {
    return r.GetRequestFloat64(key, def...)
}

func (r *Request) GetFloats(key string, def...interface{}) []float64 {
    return r.GetRequestFloats(key, def...)
}

func (r *Request) GetArray(key string, def...interface{}) []string {
    return r.GetRequestArray(key, def...)
}

func (r *Request) GetStrings(key string, def...interface{}) []string {
    return r.GetRequestStrings(key, def...)
}

func (r *Request) GetInterfaces(key string, def...interface{}) []interface{} {
    return r.GetRequestInterfaces(key, def...)
}

func (r *Request) GetMap(def...map[string]string) map[string]string {
    return r.GetRequestMap(def...)
}

// 将所有的request参数映射到struct属性上，参数pointer应当为一个struct对象的指针,
// mapping为非必需参数，自定义参数与属性的映射关系
func (r *Request) GetToStruct(pointer interface{}, mapping...map[string]string) {
    r.GetRequestToStruct(pointer, mapping...)
}

// 仅退出当前逻辑执行函数, 如:服务函数、HOOK函数
func (r *Request) Exit() {
    panic(gEXCEPTION_EXIT)
}

// 退出当前请求执行，后续所有的服务逻辑流程(包括其他的HOOK)将不会执行
func (r *Request) ExitAll() {
    r.exit = true
    panic(gEXCEPTION_EXIT_ALL)
}

// 仅针对HOOK执行，默认情况下HOOK会按照优先级进行调用，当使用ExitHook后当前类型的后续HOOK将不会被调用
func (r *Request) ExitHook() {
    panic(gEXCEPTION_EXIT_HOOK)
}

// 判断当前请求是否停止执行
func (r *Request) IsExited() bool {
    return r.exit
}

// 获取请求的服务端IP/域名
func (r *Request) GetHost() string {
    if len(r.parsedHost) == 0 {
        array, _ := gregex.MatchString(`(.+):(\d+)`, r.Host)
        if len(array) > 1 {
            r.parsedHost = array[1]
        } else {
            r.parsedHost = r.Host
        }
    }
    return r.parsedHost
}

// 判断是否为静态文件请求
func (r *Request) IsFileRequest() bool {
    return r.isFileRequest
}

// 判断是否为AJAX请求
func (r *Request) IsAjaxRequest() bool {
    return strings.EqualFold(r.Header.Get("X-Requested-With"), "XMLHttpRequest")
}

// 获取请求的客户端IP地址
func (r *Request) GetClientIp() string {
    if len(r.clientIp) == 0 {
        if r.clientIp = r.Header.Get("X-Real-IP"); r.clientIp == "" {
            array, _ := gregex.MatchString(`(.+):(\d+)`, r.RemoteAddr)
            if len(array) > 1 {
                r.clientIp = array[1]
            } else {
                r.clientIp = r.RemoteAddr
            }
        }
    }
    return r.clientIp
}

// 获得当前请求URL地址
func (r *Request) GetUrl() string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf(`%s://%s%s`, scheme, r.Host, r.URL.String())
}

// 获得请求来源URL地址
func (r *Request) GetReferer() string {
    return r.Header.Get("Referer")
}

// 获得结构体对象的参数名称标签，构成map返回
func (r *Request) getStructParamsTagMap(pointer interface{}) map[string]string {
    tagMap := make(map[string]string)
    fields := structs.Fields(pointer)
    for _, field := range fields {
        if tag := field.Tag("params"); tag != "" {
            for _, v := range strings.Split(tag, ",") {
                tagMap[strings.TrimSpace(v)] = field.Name()
            }
        }
    }
    return tagMap
}
>>>>>>> upstream/master
