package gluster

import (
    "g/os/gconsole"
    "strings"
    "fmt"
    "g/net/ghttp"
    "g/encoding/gjson"
    "g/os/gfile"
)

// 查看集群节点
func cmd_getnode () {
    r := ghttp.Get(fmt.Sprintf("http://127.0.0.1:%d/node", gPORT_API))
    if r == nil {
        fmt.Println("connect to local gluster api failed")
        return
    }
    defer r.Close()
    fmt.Println(r.ReadAll())
}

// 添加集群节点
// 使用方式：gluster addnode IP1,IP2,IP3,...
func cmd_addnode () {
    nodes := gconsole.Value.Get(2)
    if nodes != "" {
        params := make([]string, 0)
        list   := strings.Split(strings.TrimSpace(nodes), ",")
        for _, v := range list {
            if v != "" {
                params = append(params, v)
            }
        }
        if len(params) > 0 {
            r := ghttp.Post(fmt.Sprintf("http://127.0.0.1:%d/node", gPORT_API), gjson.Encode(params))
            if r == nil {
                fmt.Println("connect to local gluster api failed")
                return
            }
            defer r.Close()
            content := r.ReadAll()
            data    := gjson.DecodeToJson(content)
            if data.GetInt("result") != 1 {
                fmt.Println(data.GetString("message"))
                return
            }
        }
    }
    fmt.Println("ok")
}

// 删除集群节点
// 使用方式：gluster delnode IP1,IP2,IP3,...
func cmd_delnode () {
    nodes := gconsole.Value.Get(2)
    if nodes != "" {
        params := make([]string, 0)
        list   := strings.Split(strings.TrimSpace(nodes), ",")
        for _, v := range list {
            if v != "" {
                params = append(params, v)
            }
        }
        if len(params) > 0 {
            r := ghttp.Delete(fmt.Sprintf("http://127.0.0.1:%d/node", gPORT_API), gjson.Encode(params))
            if r == nil {
                fmt.Println("connect to local gluster api failed")
                return
            }
            defer r.Close()
            content := r.ReadAll()
            data    := gjson.DecodeToJson(content)
            if data.GetInt("result") != 1 {
                fmt.Println(data.GetString("message"))
                return
            }
        }
    }
    fmt.Println("ok")
}

// 查询kv
// 查看所有键值数据：gluster
// 查看指定键名数据：gluster 键名
func cmd_getkv () {
    k := gconsole.Value.Get(2)
    r := ghttp.Get(fmt.Sprintf("http://127.0.0.1:%d/kv?k=%s", gPORT_API, k))
    if r == nil {
        fmt.Println("connect to local gluster api failed")
        return
    }
    defer r.Close()
    fmt.Println(r.ReadAll())
}

// 设置kv
// 使用方式：gluster 键名 键值
func cmd_addkv () {
    k := gconsole.Value.Get(2)
    v := gconsole.Value.Get(3)
    if k != "" && v != ""{
        m := map[string]string{k: v}
        r := ghttp.Post(fmt.Sprintf("http://127.0.0.1:%d/kv", gPORT_API), gjson.Encode(m))
        if r == nil {
            fmt.Println("connect to local gluster api failed")
            return
        }
        defer r.Close()
        content := r.ReadAll()
        data    := gjson.DecodeToJson(content)
        if data.GetInt("result") != 1 {
            fmt.Println(data.GetString("message"))
            return
        }
    }
    fmt.Println("ok")
}

// 删除
// 使用方式：gluster delkv 键名1,键名2,键名3,...
func cmd_delkv () {
    keys := gconsole.Value.Get(2)
    if keys != "" {
        params := make([]string, 0)
        list   := strings.Split(strings.TrimSpace(keys), ",")
        for _, v := range list {
            if v != "" {
                params = append(params, v)
            }
        }
        if len(params) > 0 {
            r := ghttp.Delete(fmt.Sprintf("http://127.0.0.1:%d/kv", gPORT_API), gjson.Encode(params))
            if r == nil {
                fmt.Println("connect to local gluster api failed")
                return
            }
            defer r.Close()
            content := r.ReadAll()
            data    := gjson.DecodeToJson(content)
            if data.GetInt("result") != 1 {
                fmt.Println(data.GetString("message"))
                return
            }
        }
    }
    fmt.Println("ok")
}

// 查看Service
// 使用方式：gluster getservice [Service名称]
func cmd_getservice () {
    name := gconsole.Value.Get(2)
    r    := ghttp.Get(fmt.Sprintf("http://127.0.0.1:%d/service?name=%s", gPORT_API, name))
    if r == nil {
        fmt.Println("connect to local gluster api failed")
        return
    }
    defer r.Close()
    fmt.Println(r.ReadAll())
}

// 添加Service
// 使用方式：gluster addservice Service文件路径
func cmd_addservice () {
    path := gconsole.Value.Get(2)
    if path == "" {
        fmt.Println("please sepecify the service config file path")
        return
    }
    if !gfile.Exists(path) {
        fmt.Println("service config file does not exist")
        return
    }
    r := ghttp.Post(fmt.Sprintf("http://127.0.0.1:%d/service", gPORT_API), gfile.GetContents(path))
    if r == nil {
        fmt.Println("connect to local gluster api failed")
        return
    }
    defer r.Close()
    content := r.ReadAll()
    data    := gjson.DecodeToJson(content)
    if data.GetInt("result") != 1 {
        fmt.Println(data.GetString("message"))
        return
    }
    fmt.Println("ok")
}

// 删除Service
// 使用方式：gluster delservice Service名称1,Service名称2,Service名称3,...
func cmd_delservice () {
    s := gconsole.Value.Get(2)
    if s != "" {
        params := make([]string, 0)
        list  := strings.Split(strings.TrimSpace(s), ",")
        for _, v := range list {
            if v != "" {
                params = append(params, v)
            }
        }
        if len(params) > 0 {
            r := ghttp.Delete(fmt.Sprintf("http://127.0.0.1:%d/service", gPORT_API), gjson.Encode(params))
            if r == nil {
                fmt.Println("connect to local gluster api failed")
                return
            }
            defer r.Close()
            content := r.ReadAll()
            data    := gjson.DecodeToJson(content)
            if data.GetInt("result") != 1 {
                fmt.Println(data.GetString("message"))
                return
            }
        }
    }
    fmt.Println("ok")
}