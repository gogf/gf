// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// "不要通过共享内存来通信，而应该通过通信来共享内存"


package gproc

import (
    "fmt"
    "net"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/net/gtcp"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/encoding/gbinary"
    "gitee.com/johng/gf/g/container/gqueue"
    "gitee.com/johng/gf/g/container/gtype"
)

const (
    gPROC_DEFAULT_TCP_PORT     = 10000 // 默认开始监听的TCP端口号，如果占用则递增
    gPROC_MSG_QUEUE_MAX_LENGTH = 10000 // 进程消息队列最大长度(每个分组)
)

var (
    // 是否已开启TCP端口监听服务(使用int而非bool，以便于使用原子操作判断是否开启)
    tcpListeningCount = gtype.NewInt()
)

// 创建本地进程TCP通信服务
func startTcpListening() {
    // 一个进程只能开启一个监听goroutine
    if tcpListeningCount.Add(1) != 1 {
        return
    }
    var listen *net.TCPListener
    for i := gPROC_DEFAULT_TCP_PORT; ; i++ {
        addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", i))
        if err != nil {
            continue
        }
        listen, err = net.ListenTCP("tcp", addr)
        if err != nil {
            continue
        }
        // 将监听的端口保存到通信文件中(字符串类型存放)
        gfile.PutContents(getCommFilePath(Pid()), gconv.String(i))
        //glog.Printfln("%d: gproc listening on [%s]", Pid(), addr)
        break
    }
    for  {
        if conn, err := listen.Accept(); err != nil {
            glog.Error(err)
        } else if conn != nil {
            go tcpServiceHandler(gtcp.NewConnByNetConn(conn))
        }
    }
}

// TCP数据通信处理回调函数
func tcpServiceHandler(conn *gtcp.Conn) {
    retry := gtcp.Retry {
        Count   : 3,
        Interval: 10,
    }
    for {
        var result []byte
        buffer, err := conn.Recv(-1, retry)
        if len(buffer) > 0 {
            var msgs []*Msg
            for _, msg := range bufferToMsgs(buffer) {
                if v := commReceiveQueues.Get(msg.Group); v != nil {
                    msgs = append(msgs, msg)
                } else {
                    result = []byte(fmt.Sprintf("group [%s] does not exist", msg.Group))
                    break
                }
            }
            // 成功时会返回ok给peer
            if len(result) == 0 {
                result = []byte("ok")
                for _, msg := range msgs {
                    if v := commReceiveQueues.Get(msg.Group); v != nil {
                        v.(*gqueue.Queue).Push(msg)
                    }
                }
            }
        }
        // 产生错误(或者对方已经关闭链接)时，退出接收循环
        if err == nil {
            conn.Send(result, retry)
        } else {
            conn.Close()
            return
        }
    }
}

// 数据解包，防止黏包
// 数据格式：总长度(24bit)|发送进程PID(24bit)|接收进程PID(24bit)|分组长度(8bit)|分组名称(变长)|校验(32bit)|参数(变长)
func bufferToMsgs(buffer []byte) []*Msg {
    s    := 0
    msgs := make([]*Msg, 0)
    for s < len(buffer) {
        // 长度解析及校验
        length := gbinary.DecodeToInt(buffer[s : s + 3])
        if length < 14 || length > len(buffer) {
            s++
            continue
        }
        // 分组信息解析
        groupLen  := gbinary.DecodeToInt(buffer[s + 9 : s + 10])
        // checksum校验(仅对参数做校验，提高校验效率)
        checksum1 := gbinary.DecodeToUint32(buffer[s + 10 + groupLen : s + 10 + groupLen + 4])
        checksum2 := gtcp.Checksum(buffer[s + 10 + groupLen + 4 : s + length])
        if checksum1 != checksum2 {
            s++
            continue
        }
        // 接收进程PID校验
        if Pid() ==  gbinary.DecodeToInt(buffer[s + 6 : s + 9]) {
            msgs = append(msgs, &Msg {
                Pid   : gbinary.DecodeToInt(buffer[s + 3 : s + 6]),
                Data  : buffer[s + 10 + groupLen + 4 : s + length],
                Group : string(buffer[s + 10 : s + 10 + groupLen]),
            })
        }
        s += length
    }
    return msgs
}

// 获取其他进程传递到当前进程的消息包，阻塞执行。
func Receive(group...string) *Msg {
    // 开启接收协程时才会开启端口监听
    go startTcpListening()

    var queue *gqueue.Queue
    groupName := gPROC_COMM_DEAFULT_GRUOP_NAME
    if len(group) > 0 {
        groupName = group[0]
    }

    if v := commReceiveQueues.Get(groupName); v == nil {
        commReceiveQueues.LockFunc(func(m map[string]interface{}) {
            if v, ok := m[groupName]; ok {
                queue        = v.(*gqueue.Queue)
            } else {
                queue        = gqueue.New(gPROC_MSG_QUEUE_MAX_LENGTH)
                m[groupName] = queue
            }
        })
    } else {
        queue = v.(*gqueue.Queue)
    }

    if v := queue.Pop(); v != nil {
        return v.(*Msg)
    }
    return nil
}