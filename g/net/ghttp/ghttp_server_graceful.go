<<<<<<< HEAD
// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
=======
// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
>>>>>>> upstream/master

package ghttp

import (
<<<<<<< HEAD
    "os"
    "fmt"
    "net"
    "context"
    "net/http"
    "crypto/tls"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gproc"
=======
    "context"
    "crypto/tls"
    "errors"
    "fmt"
    "github.com/gogf/gf/g/os/glog"
    "github.com/gogf/gf/g/os/gproc"
    "net"
    "net/http"
    "os"
    "time"
>>>>>>> upstream/master
)

// 优雅的Web Server对象封装
type gracefulServer struct {
<<<<<<< HEAD
    fd           uintptr
    addr         string
    httpServer   *http.Server
    rawListener  net.Listener // 原始listener
    listener     net.Listener // 接口化封装的listener
    isHttps      bool
    shutdownChan chan bool
=======
    fd           uintptr      // 热重启时传递的socket监听文件句柄
    addr         string       // 监听地址信息
    httpServer   *http.Server // 底层http.Server
    rawListener  net.Listener // 原始listener
    listener     net.Listener // 接口化封装的listener
    isHttps      bool         // 是否HTTPS
    status       int          // 当前Server状态(关闭/运行)
>>>>>>> upstream/master
}

// 创建一个优雅的Http Server
func (s *Server) newGracefulServer(addr string, fd...int) *gracefulServer {
    gs := &gracefulServer {
        addr         : addr,
        httpServer   : s.newHttpServer(addr),
<<<<<<< HEAD
        shutdownChan : make(chan bool),
=======
>>>>>>> upstream/master
    }
    // 是否有继承的文件描述符
    if len(fd) > 0 && fd[0] > 0 {
        gs.fd = uintptr(fd[0])
    }
    return gs
}

// 生成一个底层的Web Server对象
func (s *Server) newHttpServer(addr string) *http.Server {
<<<<<<< HEAD
    return &http.Server {
=======
    server := &http.Server {
>>>>>>> upstream/master
        Addr           : addr,
        Handler        : s.config.Handler,
        ReadTimeout    : s.config.ReadTimeout,
        WriteTimeout   : s.config.WriteTimeout,
        IdleTimeout    : s.config.IdleTimeout,
        MaxHeaderBytes : s.config.MaxHeaderBytes,
    }
<<<<<<< HEAD
=======
    server.SetKeepAlivesEnabled(s.config.KeepAlive)
    return server
>>>>>>> upstream/master
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
<<<<<<< HEAD
func (s *gracefulServer) ListenAndServeTLS(certFile, keyFile string) error {
    addr   := s.httpServer.Addr
    config := &tls.Config{}
    if s.httpServer.TLSConfig != nil {
=======
func (s *gracefulServer) ListenAndServeTLS(certFile, keyFile string, tlsConfig...*tls.Config) error {
    addr   := s.httpServer.Addr
    config := (*tls.Config)(nil)
    if len(tlsConfig) > 0 {
        config = tlsConfig[0]
    } else if s.httpServer.TLSConfig != nil {
>>>>>>> upstream/master
        *config = *s.httpServer.TLSConfig
    }
    if config.NextProtos == nil {
        config.NextProtos = []string{"http/1.1"}
    }
<<<<<<< HEAD
    var err error
    config.Certificates         = make([]tls.Certificate, 1)
    config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        return err
=======
    err := error(nil)
    if len(config.Certificates) == 0 {
        config.Certificates         = make([]tls.Certificate, 1)
        config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
    }
    if err != nil {
        return errors.New(fmt.Sprintf(`open cert file "%s","%s" failed: %s`, certFile, keyFile, err.Error()))
>>>>>>> upstream/master
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
<<<<<<< HEAD
    glog.Printfln("%d: %s server %s listening on [%s]", gproc.Pid(), s.getProto(), action, s.addr)
    err := s.httpServer.Serve(s.listener)
    <-s.shutdownChan
=======
    glog.Printf("%d: %s server %s listening on [%s]", gproc.Pid(), s.getProto(), action, s.addr)
    s.status = SERVER_STATUS_RUNNING
    err := s.httpServer.Serve(s.listener)
    s.status = SERVER_STATUS_STOPPED
>>>>>>> upstream/master
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
<<<<<<< HEAD
        ln, err = net.Listen("tcp", addr)
        if err != nil {
            err = fmt.Errorf("%d: net.Listen error: %v", gproc.Pid(), err)
=======
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
>>>>>>> upstream/master
            return nil, err
        }
    }
    return ln, nil
}

// 执行请求优雅关闭
func (s *gracefulServer) shutdown() {
<<<<<<< HEAD
    if err := s.httpServer.Shutdown(context.Background()); err != nil {
        glog.Errorfln("%d: %s server [%s] shutdown error: %v", gproc.Pid(), s.getProto(), s.addr, err)
    } else {
        //glog.Printfln("%d: %s server [%s] shutdown smoothly", gproc.Pid(), s.getProto(), s.addr)
        s.shutdownChan <- true
=======
    if s.status == SERVER_STATUS_STOPPED {
        return
    }
    if err := s.httpServer.Shutdown(context.Background()); err != nil {
        glog.Errorf("%d: %s server [%s] shutdown error: %v", gproc.Pid(), s.getProto(), s.addr, err)
>>>>>>> upstream/master
    }
}

// 执行请求强制关闭
func (s *gracefulServer) close() {
<<<<<<< HEAD
    if err := s.httpServer.Close(); err != nil {
        glog.Errorfln("%d: %s server [%s] closed error: %v", gproc.Pid(), s.getProto(), s.addr, err)
    } else {
        //glog.Printfln("%d: %s server [%s] closed smoothly", gproc.Pid(), s.getProto(), s.addr)
        s.shutdownChan <- true
=======
    if s.status == SERVER_STATUS_STOPPED {
        return
    }
    if err := s.httpServer.Close(); err != nil {
        glog.Errorf("%d: %s server [%s] closed error: %v", gproc.Pid(), s.getProto(), s.addr, err)
>>>>>>> upstream/master
    }
}

