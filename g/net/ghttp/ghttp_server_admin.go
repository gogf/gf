// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// pprof封装.

package ghttp

import (
    "strings"
    "gitee.com/johng/gf/g/os/gview"
    "runtime"
    "gitee.com/johng/gf/g/os/gproc"
    "sync"
    "gitee.com/johng/gf/g/os/gtime"
    "errors"
    "fmt"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/os/glog"
)

const (
    gADMIN_ACTION_INTERVAL_LIMIT = 3000 // (毫秒)服务开启后允许执行管理操作的间隔限制
)

// 用于服务管理的对象
type utilAdmin struct {}

// (进程级别)用于Web Server管理操作的互斥锁，保证管理操作的原子性
var serverActionLocker sync.Mutex

// (进程级别)用于记录上一次操作的时间(毫秒)
var serverActionLastTime = gtype.NewInt64(gtime.Millisecond())

// 当前服务进程所处的互斥管理操作状态
// 1 : reload
// 2 : restart
// 4 : shutdown
var serverProcessStatus  = gtype.NewInt()

// 服务管理首页
func (p *utilAdmin) Index(r *Request) {
    data := map[string]interface{}{
        "uri" : strings.TrimRight(r.URL.Path, "/"),
    }
    buffer, _ := gview.ParseContent(`
            <html>
            <head>
                <title>gf ghttp admin</title>
            </head>
            <body>
                <p><a href="{{$.uri}}/reload">reload</a></p>
                <p><a href="{{$.uri}}/restart">restart</a></p>
                <p><a href="{{$.uri}}/shutdown">shutdown</a></p>
            </body>
            </html>
    `, data)
    r.Response.Write(buffer)
}

// 服务热重启
func (p *utilAdmin) Reload(r *Request) {
    if runtime.GOOS == "windows" {
        p.Restart(r)
    } else {
        if err := r.Server.Reload(); err == nil {
            r.Response.Write("server reloaded")
        } else {
            r.Response.Write(err.Error())
        }
    }
}

// 服务完整重启
func (p *utilAdmin) Restart(r *Request) {
    if err := r.Server.Restart(); err == nil {
        r.Response.Write("server restarted")
    } else {
        r.Response.Write(err.Error())
    }
}

// 服务关闭
func (p *utilAdmin) Shutdown(r *Request) {
    r.Server.Shutdown()
    if err := r.Server.Shutdown(); err == nil {
        r.Response.Write("server shutdown")
    } else {
        r.Response.Write(err.Error())
    }
}


// 开启服务管理支持
func (s *Server) EnableAdmin(pattern...string) {
    p := "/debug/admin"
    if len(pattern) > 0 {
        p = pattern[0]
    }
    s.BindObject(p, &utilAdmin{})
}

// 平滑重启Web Server
func (s *Server) Reload() error {
    serverActionLocker.Lock()
    defer serverActionLocker.Unlock()
    if err := s.checkActionStatus(); err != nil {
        return err
    }
    if err := s.checkActionFrequence(); err != nil {
        return err
    }
    glog.Printfln("%d: server reloading", gproc.Pid())
    sendProcessMsg(gproc.Pid(), gMSG_RELOAD, nil)
    return nil
}

// 完整重启Web Server
func (s *Server) Restart() error {
    serverActionLocker.Lock()
    defer serverActionLocker.Unlock()
    if err := s.checkActionStatus(); err != nil {
        return err
    }
    if err := s.checkActionFrequence(); err != nil {
        return err
    }
    glog.Printfln("%d: server restarting", gproc.Pid())
    sendProcessMsg(gproc.Pid(), gMSG_RESTART, nil)
    return nil
}

// 关闭Web Server
func (s *Server) Shutdown() error {
    serverActionLocker.Lock()
    defer serverActionLocker.Unlock()
    if err := s.checkActionStatus(); err != nil {
        return err
    }
    if err := s.checkActionFrequence(); err != nil {
        return err
    }
    glog.Printfln("%d: server shutting down", gproc.Pid())
    sendProcessMsg(gproc.PPid(), gMSG_SHUTDOWN, nil)
    return nil
}

// 检测当前操作的频繁度
func (s *Server) checkActionFrequence() error {
    interval := gtime.Millisecond() - serverActionLastTime.Val()
    if interval < gADMIN_ACTION_INTERVAL_LIMIT {
        return errors.New(fmt.Sprintf("too frequent action, please retry in %d ms", gADMIN_ACTION_INTERVAL_LIMIT - interval))
    }
    serverActionLastTime.Set(gtime.Millisecond())
    return nil
}

// 检查当前服务进程的状态
func (s *Server) checkActionStatus() error {
    status := serverProcessStatus.Val()
    if status > 0 {
        switch status {
            case 1:
                return errors.New("server is reloading")
            case 2:
                return errors.New("server is restarting")
            case 4:
                return errors.New("server is shutting down")
        }
    }
    return nil
}