package ghttp

import (
    "io/ioutil"
    "gitee.com/johng/gf/g/encoding/gjson"
    "strconv"
)

// 获得指定名称的get参数列表
func (r *ClientRequest) GetQuery(k string) []string {
    if r.getvals == nil {
        values     := r.URL.Query()
        r.getvals = &values
    }
    if v, ok := (*r.getvals)[k]; ok {
        return v
    }
    return nil
}

// 获取指定名称的参数int类型
func (r *ClientRequest) GetQueryInt(k string) int {
    v := r.GetQuery(k)
    if v == nil {
        return -1
    } else {
        if i, err := strconv.Atoi(v[0]); err != nil {
            return -1
        } else {
            return i
        }
    }
}

func (r *ClientRequest) GetQueryString(k string) string {
    v := r.GetQuery(k)
    if v == nil {
        return ""
    } else {
        return v[0]
    }
}

func (r *ClientRequest) GetQueryArray(k string) []string {
    v := r.GetQuery(k)
    if v == nil {
        return nil
    } else {
        return v
    }
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值
func (r *ClientRequest) GetQueryMap(defaultMap map[string][]string) map[string][]string {
    m := make(map[string][]string)
    for k, v := range defaultMap {
        v2 := r.GetQueryArray(k)
        if v2 == nil {
            m[k] = v
        } else {
            m[k] = v2
        }
    }
    return m
}

// 获得post参数
func (r *ClientRequest) GetPost(k string) []string {
    if v, ok := r.PostForm[k]; ok {
        return v
    }
    return nil
}

func (r *ClientRequest) GetPostInt(k string) int {
    v := r.GetPost(k)
    if v == nil {
        return -1
    } else {
        if i, err := strconv.Atoi(v[0]); err != nil {
            return -1
        } else {
            return i
        }
    }
}

func (r *ClientRequest) GetPostString(k string) string {
    v := r.GetPost(k)
    if v == nil {
        return ""
    } else {
        return v[0]
    }
}

func (r *ClientRequest) GetPostArray(k string) []string {
    v := r.GetPost(k)
    if v == nil {
        return nil
    } else {
        return v
    }
    return nil
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值
func (r *ClientRequest) GetPostMap(defaultMap map[string][]string) map[string][]string {
    m := make(map[string][]string)
    for k, v := range defaultMap {
        if v2, ok := r.PostForm[k]; ok {
            m[k] = v2
        } else {
            m[k] = v
        }
    }
    return m
}

// 获得post或者get提交的参数，如果有同名参数，那么按照get->post优先级进行覆盖
func (r *ClientRequest) GetRequest(k string) []string {
    v := r.GetQuery(k)
    if v == nil {
        return r.GetPost(k)
    }
    return v
}

func (r *ClientRequest) GetRequestString(k string) string {
    v := r.GetRequest(k)
    if v == nil {
        return ""
    } else {
        return v[0]
    }
}

func (r *ClientRequest) GetRequestArray(k string) []string {
    v := r.GetRequest(k)
    if v == nil {
        return nil
    } else {
        return v
    }
    return nil
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值
func (r *ClientRequest) GetRequestMap(defaultMap map[string][]string) map[string][]string {
    m := make(map[string][]string)
    for k, v := range defaultMap {
        v2 := r.GetRequest(k)
        if v2 != nil {
            m[k] = v2
        } else {
            m[k] = v
        }
    }
    return m
}


// 获取原始请求输入字符串
func (r *ClientRequest) GetRaw() string {
    result, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return ""
    } else {
        return string(result)
    }
}

// 获取原始请求输入字符串
func (r *ClientRequest) GetJson() *gjson.Json {
    data := r.GetRaw()
    if data != "" {
        return gjson.DecodeToJson(data)
    }
    return nil
}



