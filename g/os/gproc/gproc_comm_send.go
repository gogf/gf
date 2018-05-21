// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gproc

import (
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
}


// 获取指定进程的TCP通信对象
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
            return conn, nil
        } else {
            return nil, err
        }
    } else {
        return nil, errors.New(fmt.Sprintf("%d not found", pid))
    }

}

// 获取指定进程监听的端口号
func getPortByPid(pid int) int {
    path    := getCommFilePath(pid)
    content := gfile.GetContents(path)
    return gconv.Int(content)
}