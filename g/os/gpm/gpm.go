// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 进程管理.
package gpm

import (
    "os"
    "gitee.com/johng/gf/g/container/gmap"
)

// 进程管理器
type Manager struct {
    processes *gmap.IntInterfaceMap // 所管理的子进程map
}

// 子进程
type Process struct {
    pm       *Manager        // 所属进程管理器
    path     string          // 可执行文件绝对路径
    args     []string        // 执行参数
    attr     *os.ProcAttr    // 进程属性
    process  *os.Process     // 底层进程对象
}

// 创建一个进程管理器
func New () *Manager {
    return &Manager{
        processes : gmap.NewIntInterfaceMap(),
    }
}

// 创建一个进程(不执行)
func (m *Manager) NewProcess(path string, args []string, env []string) *Process {
    attr := &os.ProcAttr {
        Env   : env,
        Files : []*os.File{ os.Stdin,os.Stdout,os.Stderr },
    }
    return &Process{
        pm   : m,
        path : path,
        args : args,
        attr : attr,
    }
}

// 获取一个进程
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