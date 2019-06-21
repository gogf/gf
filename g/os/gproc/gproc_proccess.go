// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gproc

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// 子进程
type Process struct {
	exec.Cmd
	Manager *Manager // 所属进程管理器
	PPid    int      // 自定义关联的父进程ID
}

// 创建一个进程(不执行)
func NewProcess(path string, args []string, environment ...[]string) *Process {
	var env []string
	if len(environment) > 0 {
		env = make([]string, 0)
		for _, v := range environment[0] {
			env = append(env, v)
		}
	} else {
		env = os.Environ()
	}
	env = append(env, fmt.Sprintf("%s=%s", gPROC_TEMP_DIR_ENV_KEY, os.TempDir()))
	p := &Process{
		Manager: nil,
		PPid:    os.Getpid(),
		Cmd: exec.Cmd{
			Args:       []string{path},
			Path:       path,
			Stdin:      os.Stdin,
			Stdout:     os.Stdout,
			Stderr:     os.Stderr,
			Env:        env,
			ExtraFiles: make([]*os.File, 0),
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
		p.Args = append(p.Args, args[start:]...)
	}
	return p
}

// 开始执行(非阻塞)
func (p *Process) Start() (int, error) {
	if p.Process != nil {
		return p.Pid(), nil
	}
	p.Env = append(p.Env, fmt.Sprintf("%s=%d", gPROC_ENV_KEY_PPID_KEY, p.PPid))
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
