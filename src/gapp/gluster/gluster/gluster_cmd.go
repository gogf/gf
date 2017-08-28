package gluster

import (
    "g/os/gconsole"
    "strings"
    "fmt"
    "g/net/ghttp"
    "g/encoding/gjson"
)

// 添加集群节点
func cmd_addnode () {
    nodes := gconsole.Value.Get(2)
    if nodes != "" {
        peers := make([]string, 0)
        list  := strings.Split(strings.Trim(nodes, " "), ",")
        for _, v := range list {
            if v != "" {
                peers = append(peers, v)
            }
        }
        if len(peers) > 0 {
            r := ghttp.Post(fmt.Sprintf("http://127.0.0.1:%d/node", gPORT_API), gjson.Encode(peers))
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
func cmd_delnode () {
    nodes := gconsole.Value.Get(2)
    if nodes != "" {
        peers := make([]string, 0)
        list  := strings.Split(strings.Trim(nodes, " "), ",")
        for _, v := range list {
            if v != "" {
                peers = append(peers, v)
            }
        }
        if len(peers) > 0 {
            r := ghttp.Delete(fmt.Sprintf("http://127.0.0.1:%d/node", gPORT_API), gjson.Encode(peers))
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

// 添加Service
func cmd_addservice () {
    s := gconsole.Value.Get(2)
    if s != "" {
        r := ghttp.Delete(fmt.Sprintf("http://127.0.0.1:%d/service", gPORT_REPL), s)
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

// 删除Service
func cmd_delservice () {
    nodes := gconsole.Value.Get(2)
    if nodes != "" {
        peers := make([]string, 0)
        list  := strings.Split(strings.Trim(nodes, " "), ",")
        for _, v := range list {
            if v != "" {
                peers = append(peers, v)
            }
        }
        if len(peers) > 0 {
            r := ghttp.Delete(fmt.Sprintf("http://127.0.0.1:%d/service", gPORT_REPL), gjson.Encode(peers))
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