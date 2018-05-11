// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gproc

import (
    "os"
    "errors"
    "fmt"
    "syscall"
)

// 子进程
type Process struct {
    pm       *Manager        // 所属进程管理器
    path     string          // 可执行文件绝对路径
    args     []string        // 执行参数
    //attr     *os.ProcAttr    // 进程属性
    attr     *syscall.ProcAttr    // 进程属性
    ppid     int             // 自定义关联的父进程ID
    process  *os.Process     // 底层进程对象
}

// 运行进程
func (p *Process) Run() (int, error) {
    if p.process != nil {
        return p.Pid(), nil
    }
    p.attr.Env = append(p.attr.Env, fmt.Sprintf("%s=%d", gPROC_ENV_KEY_PPID_KEY, p.ppid))
    if pid, err := syscall.ForkExec(p.path, p.args, p.attr); err == nil {
        p.process, _ = os.FindProcess(pid)
        if p.pm != nil {
            p.pm.processes.Set(pid, p)
        }
        return pid, nil
    } else {
        return 0, err
    }
}

func (p *Process) SetManager(m *Manager) {
    p.pm = m
}

// 设置自定义的父进程ID
func (p *Process) SetPpid(ppid int) {
    p.ppid = ppid
}

func (p *Process) SetArgs(args []string) {
    p.args = args
}

func (p *Process) AddArgs(args []string) {
    for _, v := range args {
        p.args = append(p.args, v)
    }
}

func (p *Process) SetEnv(env []string) {
    p.attr.Env = env
}

func (p *Process) AddEnv(env []string) {
    for _, v := range env {
        p.attr.Env = append(p.attr.Env, v)
    }
}

func (p *Process) SetAttr(attr *syscall.ProcAttr) {
    p.attr = attr
}

func (p *Process) GetAttr() *syscall.ProcAttr {
    return p.attr
}

// PID
func (p *Process) Pid() int {
    if p.process != nil {
        return p.process.Pid
    }
    return 0
}

// 向进程发送消息
func (p *Process) Send(data interface{}) error {
    if p.process != nil {
        return Send(p.process.Pid, data)
    }
    return errors.New("process not running")
}


// Release releases any resources associated with the Process p,
// rendering it unusable in the future.
// Release only needs to be called if Wait is not.
func (p *Process) Release() error {
    return p.process.Release()
}

// Kill causes the Process to exit immediately.
func (p *Process) Kill() error {
    if err := p.process.Kill(); err == nil {
        if p.pm != nil {
            p.pm.processes.Remove(p.Pid())
        }
        return nil
    } else {
        return err
    }
}

// Wait waits for the Process to exit, and then returns a
// ProcessState describing its status and an error, if any.
// Wait releases any resources associated with the Process.
// On most operating systems, the Process must be a child
// of the current process or an error will be returned.
func (p *Process) Wait() (*os.ProcessState, error) {
    return p.process.Wait()
}

// Signal sends a signal to the Process.
// Sending Interrupt on Windows is not implemented.
func (p *Process) Signal(sig os.Signal) error {
    return p.process.Signal(sig)
}