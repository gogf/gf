// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 进程管理.
package gproc

import (
    "os"
    "net"
    "fmt"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/net/gtcp"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/encoding/gbinary"
    "gitee.com/johng/gf/g/container/gqueue"
)

// 进程管理器
type Manager struct {
    processes *gmap.IntInterfaceMap // 所管理的子进程map
}

// 进程通信消息队列
var msgQueue = gqueue.New()

// 创建一个进程管理器
func New () *Manager {
    return &Manager{
        processes : gmap.NewIntInterfaceMap(),
    }
}

// 创建主进程与子进程的TCP通信监听服务
func (m *Manager) startTcpService() {
    go func() {
        var listen *net.TCPListener
        for i := gCOMMUNICATION_MAIN_PORT; i < gCOMMUNICATION_MAIN_PORT + 10000; i++ {
            addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("127.0.0.1:%d", i))
            if err != nil {
                continue
            }
            listen, err = net.ListenTCP("tcp", addr)
            if err != nil {
                continue
            }
        }
        for  {
            if conn, err := listen.Accept(); err != nil {
                glog.Error(err)
            } else if conn != nil {
                go tcpServiceHandler(conn)
            }
        }
    }()
}

// 创建一个进程(不执行)
func (m *Manager) NewProcess(path string, args []string, environment []string) *Process {
    env := make([]string, len(environment) + 1)
    for k, v := range environment {
        env[k] = v
    }
    env[len(env)] = gCHILD_PROCESS_ENV_STRING
    return &Process {
        pm   : m,
        path : path,
        args : args,
        attr : &os.ProcAttr {
            Env   : env,
            Files : []*os.File{ os.Stdin,os.Stdout,os.Stderr },
        },
    }
}

// 获取当前进程管理器中的一个进程
func (m *Manager) GetProcess(pid int) *Process {
    if v := m.processes.Get(pid); v != nil {
        return v.(*Process)
    }
    return nil
}

// 获取所有的进程对象，构成列表返回
func (m *Manager) Processes() []*Process {
    processes := make([]*Process, 0)
    m.processes.RLockFunc(func(m map[int]interface{}) {
        for _, v := range m {
            processes = append(processes, v.(*Process))
        }
    })
    return processes
}

// 等待所有子进程结束
func (m *Manager) WaitAll() {
    processes := m.Processes()
    if len(processes) > 0 {
        for _, p := range processes {
            p.Wait()
        }
    }
}

// 关闭所有的进程
func (m *Manager) KillAll() error {
    for _, p := range m.Processes() {
        if err := p.Kill(); err != nil {
            return err
        }
    }
    return nil
}

// 向所有进程发送信号量
func (m *Manager) SignalAll(sig os.Signal) error {
    for _, p := range m.Processes() {
        if err := p.Signal(sig); err != nil {
            return err
        }
    }
    return nil
}

// 当前进程总数
func (m *Manager) Size() int {
    return m.processes.Size()
}