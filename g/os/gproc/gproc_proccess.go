// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gproc

import (
    "os"
    "fmt"
    "gitee.com/johng/gf/g/os/glog"
    "net"
    "gitee.com/johng/gf/g/net/gtcp"
)

// 子进程
type Process struct {
    pm       *Manager        // 所属进程管理器
    path     string          // 可执行文件绝对路径
    args     []string        // 执行参数
    attr     *os.ProcAttr    // 进程属性
    process  *os.Process     // 底层进程对象
}

// 运行进程
func (p *Process) Run() (int, error) {
    if p.process != nil {
        return p.Pid(), nil
    }
    if process, err := os.StartProcess(p.path, p.args, p.attr); err == nil {
        p.process = process
        p.pm.processes.Set(process.Pid, p)
        return process.Pid, nil
    } else {
        return 0, err
    }
}

// 创建主进程与子进程的TCP通信监听服务
func (p *Process) startTcpService() {
    go func() {
        var listen *net.TCPListener
        for i := gCOMMUNICATION_CHILD_PORT; i < gCOMMUNICATION_CHILD_PORT + 10000; i++ {
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

func (p *Process) SetAttr(attr *os.ProcAttr) {
    p.attr = attr
}

func (p *Process) GetAttr() *os.ProcAttr {
    return p.attr
}

// PID
func (p *Process) Pid() int {
    if p.process != nil {
        return p.process.Pid
    }
    return 0
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
        p.pm.processes.Remove(p.Pid())
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