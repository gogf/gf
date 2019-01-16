// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gscanner provides a port scanner for local intranet.
package gscanner

import (
    "net"
    "fmt"
    "sync"
    "time"
    "errors"
    "gitee.com/johng/gf/g/net/gipv4"
)

type scanner struct {
    timeout time.Duration
}

// 初始化一个扫描器
func New() *scanner {
    return &scanner{
        6*time.Second,
    }
}

// 设置超时时间，注意这个时间是每一次扫描的超时时间，而不是总共的超时时间
func (s *scanner) SetTimeout(t time.Duration) *scanner {
    s.timeout = t
    return s
}

// 异步TCP扫描网段及端口，如果扫描的端口是打开的，那么将链接给定给回调函数进行调用
// 注意startIp和endIp需要是同一个网段，否则会报错,并且回调函数不会执行
func (s *scanner) ScanIp(startIp string, endIp string, port int, callback func(net.Conn)) error {
    if callback == nil {
        return errors.New("callback function should not be nil")
    }
    var waitGroup sync.WaitGroup
    startIplong := gipv4.Ip2long(startIp)
    endIplong   := gipv4.Ip2long(endIp)
    result      := endIplong - startIplong
    if startIplong == 0 || endIplong == 0 {
        return errors.New("invalid startip or endip: ipv4 string should be given")
    }
    if result < 0 || result > 255 {
        return errors.New("invalid startip and endip: startip and endip should be in the same ip segment")
    }

    for i := startIplong; i <= endIplong; i++ {
        waitGroup.Add(1)
        go func(ip string) {
            //fmt.Println("scanning:", ip)
            // 这里必需设置超时时间
            conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), s.timeout)
            if err == nil {
                callback(conn)
                conn.Close()
            }
            //fmt.Println("scanning:", ip, "done")
            waitGroup.Done()
        }(gipv4.Long2ip(i))
    }
    waitGroup.Wait()
    return nil
}

// 扫描目标主机打开的端口列表
func (s *scanner) ScanPort(ip string, callback func(net.Conn)) error {
    if callback == nil {
        return errors.New("callback function should not be nil")
    }

    var waitGroup sync.WaitGroup
    for i := 0; i <= 65536; i++ {
        waitGroup.Add(1)
        //fmt.Println("scanning:", i)
        go func(port int) {
            conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), s.timeout)
            if err == nil {
                callback(conn)
                conn.Close()
            }
            waitGroup.Done()
        }(i)
    }
    waitGroup.Wait()
    return nil
}

