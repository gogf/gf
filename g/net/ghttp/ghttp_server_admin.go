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
    "os"
    "gitee.com/johng/gf/g/encoding/gjson"
    "gitee.com/johng/gf/g/util/gconv"
)

const (
    gADMIN_ACTION_INTERVAL_LIMIT = 3000 // (毫秒)服务开启后允许执行管理操作的间隔限制
    gADMIN_ACTION_RELOADING      = 1
    gADMIN_ACTION_RESTARTING     = 2
    gADMIN_ACTION_SHUTINGDOWN    = 4
    gADMIN_ACTION_RELOAD_ENVKEY  = "gf.server.reload"
    gADMIN_ACTION_RESTART_ENVKEY = "gf.server.restart"
)

// 用于服务管理的对象
type utilAdmin struct {}

// (进程级别)用于Web Server管理操作的互斥锁，保证管理操作的原子性
var serverActionLocker sync.Mutex

// (进程级别)用于记录上一次操作的时间(毫秒)
var serverActionLastTime = gtype.NewInt64(gtime.Millisecond())

// 当前服务进程所处的互斥管理操作状态
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
    serverProcessStatus.Set(gADMIN_ACTION_RELOADING)
    glog.Printfln("%d: server reloading", gproc.Pid())
    forkReloadProcess()
    go shutdownWebServers()
    doneChan <- struct{}{}
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
    serverProcessStatus.Set(gADMIN_ACTION_RESTARTING)
    glog.Printfln("%d: server restarting", gproc.Pid())
    doneChan <- struct{}{}
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
    serverProcessStatus.Set(gADMIN_ACTION_SHUTINGDOWN)
    glog.Printfln("%d: server shutting down", gproc.Pid())
    go closeWebServers()
    doneChan <- struct{}{}
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
            case gADMIN_ACTION_RELOADING:
                return errors.New("server is reloading")
            case gADMIN_ACTION_RESTARTING:
                return errors.New("server is restarting")
            case gADMIN_ACTION_SHUTINGDOWN:
                return errors.New("server is shutting down")
        }
    }
    return nil
}

// 创建一个子进程，通过环境变量传参
func forkReloadProcess() {
    p   := procManager.NewProcess(os.Args[0], os.Args, os.Environ())
    // 创建新的服务进程，子进程自动从父进程复制文件描述来监听同样的端口
    sfm := getServerFdMap()
    // 将sfm中的fd按照子进程创建时的文件描述符顺序进行整理，以便子进程获取到正确的fd
    for name, m := range sfm {
        for fdk, fdv := range m {
            if len(fdv) > 0 {
                s := ""
                for _, item := range strings.Split(fdv, ",") {
                    array := strings.Split(item, "#")
                    fd    := uintptr(gconv.Uint(array[1]))
                    if fd > 0 {
                        s += fmt.Sprintf("%s#%d,", array[0], 3 + len(p.ExtraFiles))
                        p.ExtraFiles = append(p.ExtraFiles, os.NewFile(fd, ""))
                    } else {
                        s += fmt.Sprintf("%s#%d,", array[0], 0)
                    }
                }
                sfm[name][fdk] = strings.TrimRight(s, ",")
            }
        }
    }
    buffer, _ := gjson.Encode(sfm)
    p.Env = append(p.Env, fmt.Sprintf("%s=%s", gADMIN_ACTION_RELOAD_ENVKEY, string(buffer)))
    if _, err := p.Start(); err != nil {
        glog.Errorfln("%d: fork process failed, error:%s, %s", gproc.Pid(), err.Error(), string(buffer))
    }
}

// 获取所有Web Server的文件描述符map
func getServerFdMap() map[string]listenerFdMap {
    sfm := make(map[string]listenerFdMap)
    serverMapping.RLockFunc(func(m map[string]interface{}) {
        for k, v := range m {
            sfm[k] = v.(*Server).getListenerFdMap()
        }
    })
    return sfm
}

// 二进制转换为FdMap
func bufferToServerFdMap(buffer []byte) map[string]listenerFdMap {
    sfm := make(map[string]listenerFdMap)
    if len(buffer) > 0 {
        j, _ := gjson.LoadContent(buffer, "json")
        for k, _ := range j.ToMap() {
            m := make(map[string]string)
            for k, v := range j.GetMap(k) {
                m[k] = gconv.String(v)
            }
            sfm[k] = m
        }
    }
    return sfm
}

// 关优雅闭进程所有端口的Web Server服务
// 注意，只是关闭Web Server服务，并不是退出进程
func shutdownWebServers() {
    serverMapping.RLockFunc(func(m map[string]interface{}) {
        for _, v := range m {
            for _, s := range v.(*Server).servers {
                s.shutdown()
            }
        }
    })
}

// 强制关闭进程所有端口的Web Server服务
// 注意，只是关闭Web Server服务，并不是退出进程
func closeWebServers() {
    serverMapping.RLockFunc(func(m map[string]interface{}) {
        for _, v := range m {
            for _, s := range v.(*Server).servers {
                s.close()
            }
        }
    })
}