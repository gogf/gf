// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package ghttp

import (
    "errors"
    "os"
    "fmt"
    "net"
    "context"
    "net/http"
    "crypto/tls"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gproc"
    "time"
)

// 优雅的Web Server对象封装
type gracefulServer struct {
    fd           uintptr
    addr         string
    httpServer   *http.Server
    rawListener  net.Listener // 原始listener
    listener     net.Listener // 接口化封装的listener
    isHttps      bool         // 是否HTTPS
    status       int          // 当前Server状态(关闭/运行)
}

// 创建一个优雅的Http Server
func (s *Server) newGracefulServer(addr string, fd...int) *gracefulServer {
    gs := &gracefulServer {
        addr         : addr,
        httpServer   : s.newHttpServer(addr),
    }
    // 是否有继承的文件描述符
    if len(fd) > 0 && fd[0] > 0 {
        gs.fd = uintptr(fd[0])
    }
    return gs
}

// 生成一个底层的Web Server对象
func (s *Server) newHttpServer(addr string) *http.Server {
    return &http.Server {
        Addr           : addr,
        Handler        : s.config.Handler,
        ReadTimeout    : s.config.ReadTimeout,
        WriteTimeout   : s.config.WriteTimeout,
        IdleTimeout    : s.config.IdleTimeout,
        MaxHeaderBytes : s.config.MaxHeaderBytes,
    }
}

// 执行HTTP监听
func (s *gracefulServer) ListenAndServe() error {
    addr    := s.httpServer.Addr
    ln, err := s.getNetListener(addr)
    if err != nil {
        return err
    }
    s.listener    = ln
    s.rawListener = ln
    return s.doServe()
}

// 获得文件描述符
func (s *gracefulServer) Fd() uintptr {
    if s.rawListener != nil {
        file, err := s.rawListener.(*net.TCPListener).File()
        if err == nil {
            return file.Fd()
        }
    }
    return 0
}

// 设置自定义fd
func (s *gracefulServer) setFd(fd int) {
    s.fd = uintptr(fd)
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
    err := error(nil)
    config.Certificates         = make([]tls.Certificate, 1)
    config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        return errors.New(fmt.Sprintf(`open cert file "%s","%s" failed: %s`, certFile, keyFile, err.Error()))
    }
    ln, err := s.getNetListener(addr)
    if err != nil {
        return err
    }

    s.listener    = tls.NewListener(ln, config)
    s.rawListener = ln
    return s.doServe()
}

// 获取服务协议字符串
func (s *gracefulServer) getProto() string {
    proto := "http"
    if s.isHttps {
        proto = "https"
    }
    return proto
}

// 开始执行Web Server服务处理
func (s *gracefulServer) doServe() error {
    action := "started"
    if s.fd != 0 {
        action = "reloaded"
    }
    glog.Printfln("%d: %s server %s listening on [%s]", gproc.Pid(), s.getProto(), action, s.addr)
    s.status = SERVER_STATUS_RUNNING
    err := s.httpServer.Serve(s.listener)
    s.status = SERVER_STATUS_STOPPED
    return err
}

// 自定义的net.Listener
func (s *gracefulServer) getNetListener(addr string) (net.Listener, error) {
    var ln net.Listener
    var err error
    if s.fd > 0 {
        f      := os.NewFile(s.fd, "")
        ln, err = net.FileListener(f)
        if err != nil {
            err = fmt.Errorf("%d: net.FileListener error: %v", gproc.Pid(), err)
            return nil, err
        }
    } else {
        // 如果监听失败，1秒后重试，最多重试3次
        for i := 0; i < 3; i++ {
            ln, err = net.Listen("tcp", addr)
            if err != nil {
                err = fmt.Errorf("%d: net.Listen error: %v", gproc.Pid(), err)
                time.Sleep(time.Second)
            } else {
                err = nil
                break
            }
        }
        if err != nil {
            return nil, err
        }
    }
    return ln, nil
}

// 执行请求优雅关闭
func (s *gracefulServer) shutdown() {
    if s.status == SERVER_STATUS_STOPPED {
        return
    }
    if err := s.httpServer.Shutdown(context.Background()); err != nil {
        glog.Errorfln("%d: %s server [%s] shutdown error: %v", gproc.Pid(), s.getProto(), s.addr, err)
    }
}

// 执行请求强制关闭
func (s *gracefulServer) close() {
    if s.status == SERVER_STATUS_STOPPED {
        return
    }
    if err := s.httpServer.Close(); err != nil {
        glog.Errorfln("%d: %s server [%s] closed error: %v", gproc.Pid(), s.getProto(), s.addr, err)
    }
}

