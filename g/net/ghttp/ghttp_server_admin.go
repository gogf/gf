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
    "time"
    "runtime"
)

const (
    gADMIN_ACTION_INTERVAL_LIMIT = 2000 // (毫秒)服务开启后允许执行管理操作的间隔限制
    gADMIN_ACTION_RESTARTING     = 1
    gADMIN_ACTION_SHUTINGDOWN    = 2
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
                <p><a href="{{$.uri}}/restart">restart</a></p>
                <p><a href="{{$.uri}}/shutdown">shutdown</a></p>
            </body>
            </html>
    `, data)
    r.Response.Write(buffer)
}

// 服务重启
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

// 重启Web Server，参数支持自定义重启的可执行文件路径，不传递时默认和原有可执行文件路径一致。
// 针对*niux系统: 平滑重启
// 针对windows : 完整重启
func (s *Server) Restart(newExeFilePath...string) error {
    serverActionLocker.Lock()
    defer serverActionLocker.Unlock()
    if err := s.checkActionStatus(); err != nil {
        return err
    }
    if err := s.checkActionFrequence(); err != nil {
        return err
    }
    restartWebServers(newExeFilePath...)
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
    shutdownWebServers()
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
            case gADMIN_ACTION_RESTARTING:  return errors.New("server is restarting")
            case gADMIN_ACTION_SHUTINGDOWN: return errors.New("server is shutting down")
        }
    }
    return nil
}

// 平滑重启：创建一个子进程，通过环境变量传参
func forkReloadProcess(newExeFilePath...string) {
    path := os.Args[0]
    if len(newExeFilePath) > 0 {
        path = newExeFilePath[0]
    }
    p   := procManager.NewProcess(path, os.Args, os.Environ())
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
    p.Env = append(p.Env, gADMIN_ACTION_RELOAD_ENVKEY + "=" + string(buffer))
    if _, err := p.Start(); err != nil {
        glog.Errorfln("%d: fork process failed, error:%s, %s", gproc.Pid(), err.Error(), string(buffer))
    }
}

// 完整重启：创建一个新的子进程
func forkRestartProcess(newExeFilePath...string) {
    path := os.Args[0]
    if len(newExeFilePath) > 0 {
        path = newExeFilePath[0]
    }
    // 去掉平滑重启的环境变量参数
    os.Unsetenv(gADMIN_ACTION_RELOAD_ENVKEY)
    env := os.Environ()
    env  = append(env, gADMIN_ACTION_RESTART_ENVKEY + "=1")
    p := procManager.NewProcess(path, os.Args, env)
    if _, err := p.Start(); err != nil {
        glog.Errorfln("%d: fork process failed, error:%s", gproc.Pid(), err.Error())
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

// Web Server重启
func restartWebServers(newExeFilePath...string) {
    serverProcessStatus.Set(gADMIN_ACTION_RESTARTING)
    glog.Printfln("%d: server restarting", gproc.Pid())
    if runtime.GOOS == "windows" {
        // 异步1秒后再执行重启，目的是让接口能够正确返回结果，否则接口会报错(因为web server关闭了)
        gtime.SetTimeout(time.Second, func() {
            forcedlyCloseWebServers()
            forkRestartProcess(newExeFilePath...)
        })
    } else {
        forkReloadProcess(newExeFilePath...)
        go gracefulShutdownWebServers()
        doneChan <- struct{}{}
    }
}

// Web Server关闭服务
func shutdownWebServers() {
    serverProcessStatus.Set(gADMIN_ACTION_SHUTINGDOWN)
    glog.Printfln("%d: server shutting down", gproc.Pid())
    // 异步1秒后再执行重启，目的是让接口能够正确返回结果，否则接口会报错(因为web server关闭了)
    gtime.SetTimeout(time.Second, func() {
        forcedlyCloseWebServers()
        doneChan <- struct{}{}
    })
}

// 关优雅闭进程所有端口的Web Server服务
// 注意，只是关闭Web Server服务，并不是退出进程
func gracefulShutdownWebServers() {
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
func forcedlyCloseWebServers() {
    serverMapping.RLockFunc(func(m map[string]interface{}) {
        for _, v := range m {
            for _, s := range v.(*Server).servers {
                s.close()
            }
        }
    })
}