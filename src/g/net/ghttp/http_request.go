package ghttp

// 获得get参数
func (r *Request) get(k string) interface{} {
    if r.getValues == nil {
        values     := r.URL.Query()
        r.getValues = &values
    }
    if v, ok := (*r.getValues)[k]; ok {
        return v
    }
    return nil
}

func (r *Request) Get(k string) string {
    v := r.get(k)
    if v == nil {
        return ""
    } else {
        return v.([]string)[0]
    }
}

func (r *Request) GetString(k string) string {
    return r.Get(k)
}

func (r *Request) GetArray(k string) []string {
    v := r.get(k)
    if v == nil {
        return nil
    } else {
        return v.([]string)
    }
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值
func (r *Request) GetMap(defaultMap map[string]interface{}) map[string]interface{} {
    m := make(map[string]interface{})
    for k, v := range defaultMap {
        if v2, ok := (*r.getValues)[k]; ok {
            m[k] = v2
        } else {
            m[k] = v
        }
    }
    return m
}

// 获得post参数
func (r *Request) GetPost(k string) interface{} {
    if v, ok := r.PostForm[k]; ok {
        return v
    }
    return nil
}

func (r *Request) GetPostString(k string) interface{} {
    v := r.GetPost(k)
    if v == nil {
        return ""
    } else {
        return v.(string)
    }
}

func (r *Request) GetPostArray(k string) []string {
    v := r.GetPost(k)
    if v == nil {
        return nil
    } else {
        return v.([]string)
    }
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值
func (r *Request) GetPostMap(defaultMap map[string]interface{}) map[string]interface{} {
    m := make(map[string]interface{})
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
func (r *Request) GetRequest(k string) interface{} {
    v := r.Get(k)
    if v == "" {
        return r.GetPost(k)
    }
    return v
}

func (r *Request) GetRequestString(k string) interface{} {
    v := r.GetRequest(k)
    if v == nil {
        return ""
    } else {
        return v.(string)
    }
}

func (r *Request) GetRequestArray(k string) []string {
    v := r.GetRequest(k)
    if v == nil {
        return nil
    } else {
        return v.([]string)
    }
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值
func (r *Request) GetRequestMap(defaultMap map[string]interface{}) map[string]interface{} {
    m := make(map[string]interface{})
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
