// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gproc

import (
    "os"
    "fmt"
    "os/exec"
    "errors"
)

// 子进程
type Process struct {
    exec.Cmd
    Manager  *Manager // 所属进程管理器
    PPid     int      // 自定义关联的父进程ID
}

// 开始执行(非阻塞)
func (p *Process) Start() (int, error) {
    if p.Process != nil {
        return p.Pid(), nil
    }
    if p.PPid > 0 {
        p.Env = append(p.Env, fmt.Sprintf("%s=%d", gPROC_ENV_KEY_PPID_KEY, p.PPid))
    }
    if err := p.Cmd.Start(); err == nil {
        if p.Manager != nil {
            p.Manager.processes.Set(p.Process.Pid, p)
        }
        return p.Process.Pid, nil
    } else {
        return 0, err
    }
}

// 运行进程(阻塞等待执行完毕)
func (p *Process) Run() error {
    if _, err := p.Start(); err == nil {
        return p.Wait()
    } else {
        return err
    }
}

// PID
func (p *Process) Pid() int {
    if p.Process != nil {
        return p.Process.Pid
    }
    return 0
}

// 向进程发送消息
func (p *Process) Send(data []byte) error {
    if p.Process != nil {
        return Send(p.Process.Pid, data)
    }
    return errors.New("invalid process")
}


// Release releases any resources associated with the Process p,
// rendering it unusable in the future.
// Release only needs to be called if Wait is not.
func (p *Process) Release() error {
    return p.Process.Release()
}

// Kill causes the Process to exit immediately.
func (p *Process) Kill() error {
    if err := p.Process.Kill(); err == nil {
        if p.Manager != nil {
            p.Manager.processes.Remove(p.Pid())
        }
        return nil
    } else {
        return err
    }
}

// Signal sends a signal to the Process.
// Sending Interrupt on Windows is not implemented.
func (p *Process) Signal(sig os.Signal) error {
    return p.Process.Signal(sig)
}
