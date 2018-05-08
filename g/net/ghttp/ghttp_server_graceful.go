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
    "context"
    "net/http"
    "crypto/tls"
    "gitee.com/johng/gf/g/os/glog"
)

// 优雅的Web Server对象封装
type gracefulServer struct {
    fd           uintptr
    addr         string
    httpServer   *http.Server
    listener     net.Listener
    shutdownChan chan bool
}

// 创建一个优雅的Http Server
func (s *Server) newGracefulServer(addr string, fd...int) *gracefulServer {
    gs := &gracefulServer {
        addr         : addr,
        httpServer   : s.newServer(addr),
        shutdownChan : make(chan bool),
    }
    if len(fd) > 0 && fd[0] > 0 {
        gs.fd = uintptr(fd[0])
    }
    return gs
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
    err := s.httpServer.Serve(s.listener)
    <-s.shutdownChan
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

// 执行请求优雅关闭
func (s *gracefulServer) shutdown() {
    if err := s.httpServer.Shutdown(context.Background()); err != nil {
        glog.Errorf("server %s shutdown error: %v\n", s.addr, err)
    } else {
        glog.Printf("server %s shutdown successfully\n", s.addr)
        s.shutdownChan <- true
    }
}

