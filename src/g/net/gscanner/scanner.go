// 局域网端口扫描
package gscanner

import (
    "net"
    "g/net/gip"
    "fmt"
    "errors"
    "sync"
    "time"
)

// TCP扫描网段及端口，如果扫描的端口是打开的，那么将链接给定给回调函数进行调用
// 注意startIp和endIp需要是同一个网段，否则会报错,并且回调函数不会执行
func TcpScan(startIp string, endIp string, port int, callback func(net.Conn)) error {
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
            // 这里必需设置超时时间，对于局域网异步端口扫描来讲，3秒已经足足够用
            conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), 3 * time.Second)
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

