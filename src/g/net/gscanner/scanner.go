// 局域网端口扫描
package gscanner

import (
    "net"
    "g/net/gip"
    "fmt"
    "errors"
    "sync"
)

// 扫描网段及端口，如果扫描的端口是打开的，那么将链接给定给回调函数进行调用
// 注意startIp和endIp需要是同一个网段，否则会报错,并且回调函数不会执行
func scan(network string, startIp string, endIp string, port int, callback func(net.Conn)) error {
    var waitGroup sync.WaitGroup
    startIplong := gip.Ip2long(startIp)
    endIplong   := gip.Ip2long(endIp)
    result      := endIplong - startIplong
    if startIplong == 0 || endIplong == 0 {
        return errors.New("invalid startip or endip: ipv4 string should be given")
    }
    if result < 0 || result > 255 {
        return errors.New("invalid startip and endip: startip and endip should be in the same ip segment")
    }
    if callback == nil {
        return errors.New("callback function should not be nil")
    }

    for i := startIplong; i <= endIplong; i++ {
        waitGroup.Add(1)
        ip := gip.Long2ip(i)
        go func() {
            //fmt.Println("scanning:", ip)
            conn, err := net.Dial(network, fmt.Sprintf("%s:%d", ip, port))
            if err == nil {
                callback(conn)
                conn.Close()
            }
            //fmt.Println("scanning:", ip, "done")
            waitGroup.Done()
        }()
    }
    waitGroup.Wait()
    return nil
}

// TCP扫描网段及端口，如果扫描的端口是打开的，那么将链接给定给回调函数进行调用
// 注意startIp和endIp需要是同一个网段，否则会报错,并且回调函数不会执行
func TcpScan(startIp string, endIp string, port int, callback func(net.Conn)) error {
    return scan("tcp", startIp, endIp, port, callback)
}

// TCP扫描网段及端口，如果扫描的端口是打开的，那么将链接给定给回调函数进行调用
// 注意startIp和endIp需要是同一个网段，否则会报错,并且回调函数不会执行
func UdpScan(startIp string, endIp string, port int, callback func(net.Conn)) error {
    return scan("udp", startIp, endIp, port, callback)
}

