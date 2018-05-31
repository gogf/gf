// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 进程管理.
package gproc

import (
    "os"
    "strings"
    "os/exec"
    "gitee.com/johng/gf/g/container/gmap"
    "fmt"
)

// 进程管理器
type Manager struct {
    processes *gmap.IntInterfaceMap // 所管理的子进程map
}

// 创建一个进程管理器
func NewManager() *Manager {
    return &Manager{
        processes : gmap.NewIntInterfaceMap(),
    }
}

// 创建一个进程(不执行)
func NewProcess(path string, args []string, environment []string) *Process {
    env := make([]string, len(environment) + 1)
    for k, v := range environment {
        env[k] = v
    }
    env[len(env) - 1] = fmt.Sprintf("%s=%s", gPROC_TEMP_DIR_ENV_KEY, os.TempDir())
    p := &Process {
        Manager   : nil,
        PPid      : os.Getpid(),
        Cmd       : exec.Cmd {
            Args       : []string{path},
            Path       : path,
            Stdin      : os.Stdin,
            Stdout     : os.Stdout,
            Stderr     : os.Stderr,
            Env        : env,
            ExtraFiles : make([]*os.File, 0),
        },
    }
    // 当前工作目录
    if d, err := os.Getwd(); err == nil {
        p.Dir = d
    }
    if len(args) > 0 {
        start := 0
        if strings.EqualFold(path, args[0]) {
            start = 1
        }
        p.Args = append(p.Args, args[start : ]...)
    }
    return p
}

// 创建一个进程(不执行)
func (m *Manager) NewProcess(path string, args []string, environment []string) *Process {
    p := NewProcess(path, args, environment)
    p.Manager = m
    return p
}

// 获取当前进程管理器中的一个进程
func (m *Manager) GetProcess(pid int) *Process {
    if v := m.processes.Get(pid); v != nil {
        return v.(*Process)
    }
    return nil
}

// 添加一个已存在进程到进程管理器中
func (m *Manager) AddProcess(pid int) {
    if m.processes.Get(pid) == nil {
        if process, err := os.FindProcess(pid); err == nil {
            p := m.NewProcess("", nil, nil)
            p.Process = process
            m.processes.Set(pid, p)
        }
    }
}

// 移除进程管理器中的指定进程
func (m *Manager) RemoveProcess(pid int) {
    m.processes.Remove(pid)
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

// 获取所有的进程pid，构成列表返回
func (m *Manager) Pids() []int {
    return m.processes.Keys()
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

// 向所有进程发送消息
func (m *Manager) Send(data []byte) {
    for _, p := range m.Processes() {
        p.Send(data)
    }
}

// 向指定进程发送消息
func (m *Manager) SendTo(pid int, data []byte) error {
    return Send(pid, data)
}

// 清空管理器
func (m *Manager) Clear() {
    m.processes.Clear()
}

// 当前进程总数
func (m *Manager) Size() int {
    return m.processes.Size()
}