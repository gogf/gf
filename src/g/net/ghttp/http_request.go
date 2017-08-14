package ghttp

import (
    "io/ioutil"
    "g/encoding/gjson"
)

// 获得get参数
func (r *Request) GetQuery(k string) []string {
    if r.getvals == nil {
        values     := r.URL.Query()
        r.getvals = &values
    }
    if v, ok := (*r.getvals)[k]; ok {
        return v
    }
    return nil
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
    v := r.GetQuery(k)
    if v == nil {
        return nil
    } else {
        return v
    }
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值
func (r *Request) GetQueryMap(defaultMap map[string][]string) map[string][]string {
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
func (r *Request) GetPost(k string) []string {
    if v, ok := r.PostForm[k]; ok {
        return v
    }
    return nil
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
    v := r.GetPost(k)
    if v == nil {
        return nil
    } else {
        return v
    }
    return nil
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值
func (r *Request) GetPostMap(defaultMap map[string][]string) map[string][]string {
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

func (r *Request) GetRequestArray(k string) []string {
    v := r.GetRequest(k)
    if v == nil {
        return nil
    } else {
        return v
    }
    return nil
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值
func (r *Request) GetRequestMap(defaultMap map[string][]string) map[string][]string {
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
func (r *Request) GetRaw() string {
    result, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return ""
    } else {
        return string(result)
    }
}

// 获取原始请求输入字符串
func (r *Request) GetJson() *gjson.Json {
    data := r.GetRaw()
    if data != "" {
        return gjson.DecodeToJson(&data)
    }
    return nil
}



