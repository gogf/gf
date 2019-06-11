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

package gproc

import (
<<<<<<< HEAD
    "net"
    "gitee.com/johng/gf/g/net/gtcp"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/encoding/gbinary"
    "fmt"
    "errors"
)

const (
    gPROC_COMM_FAILURE_RETRY_COUNT   = 3  // 失败重试次数
)

// 向指定gproc进程发送数据
// 数据格式：总长度(32bit) | 发送进程PID(32bit) | 接收进程PID(32bit) | 校验(32bit) | 参数(变长)
func Send(pid int, data []byte) error {
    buffer := make([]byte, 0)
    buffer  = append(buffer, gbinary.EncodeInt32(int32(len(data) + 16))...)
    buffer  = append(buffer, gbinary.EncodeInt32(int32(Pid()))...)
    buffer  = append(buffer, gbinary.EncodeInt32(int32(pid))...)
    buffer  = append(buffer, gbinary.EncodeUint32(gtcp.Checksum(data))...)
    buffer  = append(buffer, data...)
    if conn, err := getConnByPid(pid); err == nil {
        for i := gPROC_COMM_FAILURE_RETRY_COUNT; i > 0; i-- {
            if err = gtcp.Send(conn, buffer); err != nil {
                conn.Close()
                if conn, err = newConnByPid(pid); err != nil {
                    return err
                }
            } else {
                //glog.Printfln("%d: sent to %d, %v", Pid(), pid, buffer)
                break
            }
        }
        return err
    } else {
        return err
    }
=======
    "bytes"
    "errors"
    "fmt"
    "github.com/gogf/gf/g/encoding/gbinary"
    "github.com/gogf/gf/g/net/gtcp"
    "github.com/gogf/gf/g/os/gfcache"
    "github.com/gogf/gf/g/os/glog"
    "github.com/gogf/gf/g/util/gconv"
    "io"
    "time"
)

const (
    gPROC_COMM_FAILURE_RETRY_COUNT   = 3    // 失败重试次数
    gPROC_COMM_FAILURE_RETRY_TIMEOUT = 1000 // (毫秒)失败重试间隔
    gPROC_COMM_SEND_TIMEOUT          = 5000 // (毫秒)发送超时时间
    gPROC_COMM_DEAFULT_GRUOP_NAME    = ""   // 默认分组名称
)

// 向指定gproc进程发送数据.
// 数据格式：总长度(24bit)|发送进程PID(24bit)|接收进程PID(24bit)|分组长度(8bit)|分组名称(变长)|校验(32bit)|参数(变长)
func Send(pid int, data []byte, group...string) error {
    groupName := gPROC_COMM_DEAFULT_GRUOP_NAME
    if len(group) > 0 {
        groupName = group[0]
    }
    buffer := make([]byte, 0)
    buffer  = append(buffer, gbinary.EncodeByLength(3, len(groupName) + len(data) + 14)...)
    buffer  = append(buffer, gbinary.EncodeByLength(3, Pid())...)
    buffer  = append(buffer, gbinary.EncodeByLength(3, pid)...)
    buffer  = append(buffer, gbinary.EncodeByLength(1, len(groupName))...)
    buffer  = append(buffer, []byte(groupName)...)
    buffer  = append(buffer, gbinary.EncodeUint32(gtcp.Checksum(data))...)
    buffer  = append(buffer, data...)
    // 执行发送流程
    var err  error
    var buf  []byte
    var conn *gtcp.Conn
    for i := gPROC_COMM_FAILURE_RETRY_COUNT; i > 0; i-- {
        if conn, err = getConnByPid(pid); err == nil {
            defer conn.Close()
            buf, err = conn.SendRecvWithTimeout(buffer, -1, gPROC_COMM_SEND_TIMEOUT*time.Millisecond)
            if len(buf) > 0 {
                // 如果有返回值，如果不是"ok"，那么表示是错误信息
                if !bytes.EqualFold(buf, []byte("ok")) {
                    err = errors.New(string(buf))
                    break
                }
            }
            // EOF不算异常错误
            if err == nil || err == io.EOF {
                break
            } else {
                glog.Error(err)
            }
        }
        time.Sleep(gPROC_COMM_FAILURE_RETRY_TIMEOUT*time.Millisecond)
    }
    return err
>>>>>>> upstream/master
}


// 获取指定进程的TCP通信对象
<<<<<<< HEAD
func getConnByPid(pid int) (net.Conn, error) {
    if v := commPidConnMap.Get(pid); v != nil {
        return v.(net.Conn), nil
    } else {
        return newConnByPid(pid)
    }
}

// 创建与指定进程的TCP通信对象
func newConnByPid(pid int) (net.Conn, error) {
    if port := getPortByPid(pid); port > 0 {
        if conn, err := gtcp.Conn("127.0.0.1", port); err == nil {
            commPidConnMap.Set(pid, conn)
=======
func getConnByPid(pid int) (*gtcp.Conn, error) {
    port := getPortByPid(pid)
    if port > 0 {
        if conn, err := gtcp.NewConn(fmt.Sprintf("127.0.0.1:%d", port)); err == nil {
>>>>>>> upstream/master
            return conn, nil
        } else {
            return nil, err
        }
<<<<<<< HEAD
    } else {
        return nil, errors.New(fmt.Sprintf("%d not found", pid))
    }

=======
    }
    return nil, errors.New(fmt.Sprintf("could not find port for pid: %d" , pid))
>>>>>>> upstream/master
}

// 获取指定进程监听的端口号
func getPortByPid(pid int) int {
    path    := getCommFilePath(pid)
<<<<<<< HEAD
    content := gfile.GetContents(path)
=======
    content := gfcache.GetContents(path)
>>>>>>> upstream/master
    return gconv.Int(content)
}