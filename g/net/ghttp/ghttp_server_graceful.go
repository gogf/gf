// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package ghttp

import (
    "fmt"
    "net"
    "os"
    "syscall"
    "context"
    "net/http"
    "os/signal"
    "crypto/tls"
    "gitee.com/johng/gf/g/os/glog"
)

const (
    gGF_WEB_SERVER_GRACEFUL_ENVIRON_KEY    = "GF_WEB_SERVER_GRACEFUL"
    gGF_WEB_SERVER_GRACEFUL_ENVIRON_STRING = gGF_WEB_SERVER_GRACEFUL_ENVIRON_KEY + "=1"
    gGF_WEB_SERVER_GRACEFUL_LISTENER_FD    = 3
)

// 优雅的Web Server对象封装
type gracefulServer struct {
    httpServer   *http.Server
    listener     net.Listener
    isGraceful   bool
    signalChan   chan os.Signal
    shutdownChan chan bool
}

// 创建一个优雅的Http Server
func (s *Server) newGracefulServer(addr string) *gracefulServer {
    isGraceful := false
    if os.Getenv(gGF_WEB_SERVER_GRACEFUL_ENVIRON_KEY) != "" {
        isGraceful = true
    }
    return &gracefulServer {
        httpServer   : s.newServer(addr),
        isGraceful   : isGraceful,
        signalChan   : make(chan os.Signal),
        shutdownChan : make(chan bool),
    }
}

// 执行HTTP监听
func (s *gracefulServer) ListenAndServe() error {
    addr    := s.httpServer.Addr
    ln, err := s.getNetListener(addr)
    if err != nil {
        return err
    }
    s.listener = ln
    return s.doServe()
}

// 执行HTTPS监听
func (s *gracefulServer) ListenAndServeTLS(certFile, keyFile string) error {
    addr   := s.httpServer.Addr
    config := &tls.Config{}
    if s.httpServer.TLSConfig != nil {
        *config = *s.httpServer.TLSConfig
    }
    if config.NextProtos == nil {
        config.NextProtos = []string{"http/1.1"}
    }
    var err error
    config.Certificates         = make([]tls.Certificate, 1)
    config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        return err
    }
    ln, err := s.getNetListener(addr)
    if err != nil {
        return err
    }
    s.listener = tls.NewListener(ln, config)
    return s.doServe()
}

// 开始执行Web Server服务处理
func (s *gracefulServer) doServe() error {
    go s.handleSignals()
    err := s.httpServer.Serve(s.listener)
    <-s.shutdownChan
    return err
}

// 自定义的net.Listener
func (s *gracefulServer) getNetListener(addr string) (net.Listener, error) {
    var ln net.Listener
    var err error
    if s.isGraceful {
        //path   := fmt.Sprintf("%s%sgf.web.server.fd.%d", gfile.TempDir(), gfile.Separator,gtime.Nanosecond())
        //f, err := gfile.Open(path)
        //if err != nil {
        //    return nil, err
        //}
        f := os.NewFile(gGF_WEB_SERVER_GRACEFUL_LISTENER_FD, "")
        ln, err = net.FileListener(f)
        if err != nil {
            err = fmt.Errorf("net.FileListener error: %v", err)
            return nil, err
        }
    } else {
        ln, err = net.Listen("tcp", addr)
        if err != nil {
            err = fmt.Errorf("net.Listen error: %v", err)
            return nil, err
        }
    }
    return ln, nil
}

// 处理终端信号指令监听
func (s *gracefulServer) handleSignals() {
    var sig os.Signal

    signal.Notify(
        s.signalChan,
        syscall.SIGTERM,
        syscall.SIGUSR2,
    )

    for {
        sig = <-s.signalChan
        switch sig {
            case syscall.SIGTERM:
                glog.Println("received SIGTERM, graceful shutting down HTTP server")
                s.shutdown()

            case syscall.SIGUSR2:
                glog.Println("received SIGUSR2, graceful restarting HTTP server")
                if pid, err := s.startNewProcess(); err != nil {
                    glog.Printf("start new process failed: %v, continue serving\n", err)
                } else {
                    glog.Printf("start new process successfully, the new pid is %d\n", pid)
                    s.shutdown()
                }
            default:
        }
    }
}

// 执行请求优雅关闭
func (s *gracefulServer) shutdown() {
    if err := s.httpServer.Shutdown(context.Background()); err != nil {
        glog.Errorf("server shutdown error: %v\n", err)
    } else {
        glog.Println("server shutdown success")
        s.shutdownChan <- true
    }
}

// 创建子进程来监听并处理新的HTTP请求，与父进程使用的是同一个socket文件描述符
func (s *gracefulServer) startNewProcess() (uintptr, error) {
    listenerFd, err := s.getTCPListenerFd()
    if err != nil {
        return 0, fmt.Errorf("failed to get socket file descriptor: %v", err)
    }
    // 构造子进程的环境变量，并增加环境变量参数以标识该进程是graceful子进程
    env := make([]string, 0)
    for _, value := range os.Environ() {
        if value != gGF_WEB_SERVER_GRACEFUL_ENVIRON_STRING {
            env = append(env, value)
        }
    }
    env = append(env, gGF_WEB_SERVER_GRACEFUL_ENVIRON_STRING)
    fork, err := syscall.ForkExec(os.Args[0], os.Args, &syscall.ProcAttr {
        Env   : env,
        Files : []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), listenerFd},
    })
    if err != nil {
        return 0, fmt.Errorf("failed to forkexec: %v", err)
    }
    return uintptr(fork), nil
}

// 获得对应net.TCPListener的文件描述符文件ID
func (s *gracefulServer) getTCPListenerFd() (uintptr, error) {
    file, err := s.listener.(*net.TCPListener).File()
    if err != nil {
        return 0, err
    }
    return file.Fd(), nil
}
