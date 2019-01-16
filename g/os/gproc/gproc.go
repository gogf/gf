// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gproc implements communication and management of processes.
// 
// 进程管理/通信, 
// 本进程管理从syscall, os.StartProcess, exec.Cmd都使用过, 
// 最后采用了exec.Cmd来实现多进程管理，这是一个顶层的跨平台封装，兼容性更好，另外两个是偏底层的接口。
package gproc

import (
    "os"
    "time"
    "gitee.com/johng/gf/g/util/gconv"
    "strings"
    "bytes"
    "io"
    "runtime"
    "gitee.com/johng/gf/g/os/gfile"
)

const (
    gPROC_ENV_KEY_PPID_KEY = "GPROC_PPID"
    gPROC_TEMP_DIR_ENV_KEY = "GPROC_TEMP_DIR"
)

// 进程开始执行时间
var processStartTime = time.Now()

// 获取当前进程ID
func Pid() int {
    return os.Getpid()
}

// 获取父进程ID(gproc父进程，如果当前进程本身就是父进程，那么返回自身的pid，不存在时则使用系统父进程)
func PPid() int {
    if !IsChild() {
        return Pid()
    }
    // gPROC_ENV_KEY_PPID_KEY为gproc包自定义的父进程
    ppidValue := os.Getenv(gPROC_ENV_KEY_PPID_KEY)
    if ppidValue != "" && ppidValue != "0" {
        return gconv.Int(ppidValue)
    }
    return PPidOS()
}

// 获取父进程ID(系统父进程)
func PPidOS() int {
    return os.Getppid()
}

// 判断当前进程是否为gproc创建的子进程
func IsChild() bool {
    ppidValue := os.Getenv(gPROC_ENV_KEY_PPID_KEY)
    return ppidValue != "" && ppidValue != "0"
}

// 设置gproc父进程ID，当ppid为0时表示该进程为gproc主进程，否则为gproc子进程
func SetPPid(ppid int) {
    if ppid > 0 {
        os.Setenv(gPROC_ENV_KEY_PPID_KEY, gconv.String(ppid))
    } else {
        os.Unsetenv(gPROC_ENV_KEY_PPID_KEY)
    }
}

// 进程开始执行时间
func StartTime() time.Time {
    return processStartTime
}

// 进程已经运行的时间(毫秒)
func Uptime() int {
    return int(time.Now().UnixNano()/1e6 - processStartTime.UnixNano()/1e6)
}

// 阻塞执行shell指令，并给定输入输出对象
func Shell(cmd string, out io.Writer, in io.Reader) error {
    p := NewProcess(getShell(), []string{getShellOption(), cmd})
    p.Stdin  = in
    p.Stdout = out
    return p.Run()
}

// 阻塞执行shell指令，并输出结果当终端(如果需要异步，请使用goroutine)
func ShellRun(cmd string) error {
    p := NewProcess(getShell(), []string{getShellOption(), cmd})
    return p.Run()
}

// 阻塞执行shell指令，并返回输出结果(如果需要异步，请使用goroutine)
func ShellExec(cmd string) (string, error) {
    buf := bytes.NewBuffer(nil)
    p   := NewProcess(getShell(), []string{getShellOption(), cmd})
    p.Stdout = buf
    err := p.Run()
    return buf.String(), err
}

// 检测环境变量中是否已经存在指定键名
func checkEnvKey(env []string, key string) bool {
    for _, v := range env {
        if len(v) >= len(key) && strings.EqualFold(v[0 : len(key)], key) {
            return true
        }
    }
    return false
}

// 获取当前系统下的shell路径
func getShell() string {
    switch runtime.GOOS {
        case "windows":
            return searchBinFromEnvPath("cmd.exe")
        default:
            path := searchBinFromEnvPath("bash")
            if path == "" {
                path = searchBinFromEnvPath("sh")
            }
            return path
    }
}

// 获取当前系统默认shell执行指令的option参数
func getShellOption() string {
    switch runtime.GOOS {
        case "windows":
            return "/c"
        default:
            return "-c"
    }
}

// 从环境变量PATH中搜索可执行文件
func searchBinFromEnvPath(file string) string {
    // 如果是绝对路径，或者相对路径下存在，那么直接返回
    if gfile.Exists(file) {
        return file
    }
    array := ([]string)(nil)
    switch runtime.GOOS {
        case "windows":
            array = strings.Split(os.Getenv("Path"), ";")
            if gfile.Ext(file) != ".exe" {
                file += ".exe"
            }
        default:
            array = strings.Split(os.Getenv("PATH"), ":")
    }
    if len(array) > 0 {
        for _, v := range array {
            path := v + gfile.Separator + file
            if gfile.Exists(path) {
                return path
            }
        }
    }
    return ""
}

