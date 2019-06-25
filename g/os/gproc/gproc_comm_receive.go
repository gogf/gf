// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// "不要通过共享内存来通信，而应该通过通信来共享内存"

package gproc

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/g/container/gqueue"
	"github.com/gogf/gf/g/container/gtype"
	"github.com/gogf/gf/g/net/gtcp"
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/util/gconv"
	"net"
)

const (
	gPROC_DEFAULT_TCP_PORT     = 10000 // 默认开始监听的TCP端口号，如果占用则递增
	gPROC_MSG_QUEUE_MAX_LENGTH = 10000 // 进程消息队列最大长度(每个分组)
)

var (
	// 是否已开启TCP端口监听服务
	tcpListened = gtype.NewBool()
)

// 获取其他进程传递到当前进程的消息包，阻塞执行。
// 进程只有在执行该方法后才会打开请求端口，默认情况下不允许进程间通信。
func Receive(group ...string) *Msg {
	// 一个进程只能开启一个监听goroutine
	if !tcpListened.Val() && tcpListened.Set(true) == false {
		go startTcpListening()
	}
	queue := (*gqueue.Queue)(nil)
	groupName := gPROC_COMM_DEAFULT_GRUOP_NAME
	if len(group) > 0 {
		groupName = group[0]
	}
	if v := commReceiveQueues.Get(groupName); v == nil {
		commReceiveQueues.LockFunc(func(m map[string]interface{}) {
			if v, ok := m[groupName]; ok {
				queue = v.(*gqueue.Queue)
			} else {
				queue = gqueue.New(gPROC_MSG_QUEUE_MAX_LENGTH)
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

// 创建本地进程TCP通信服务
func startTcpListening() {
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
		if err := gfile.PutContents(getCommFilePath(Pid()), gconv.String(i)); err != nil {
			glog.Error(err)
		}
		break
	}
	for {
		if conn, err := listen.Accept(); err != nil {
			glog.Error(err)
		} else if conn != nil {
			go tcpServiceHandler(gtcp.NewConnByNetConn(conn))
		}
	}
}

// TCP数据通信处理回调函数
func tcpServiceHandler(conn *gtcp.Conn) {
	option := gtcp.PkgOption{
		Retry: gtcp.Retry{
			Count:    3,
			Interval: 10,
		},
	}
	for {
		var result []byte
		buffer, err := conn.RecvPkg(option)
		if len(buffer) > 0 {
			msg := new(Msg)
			if err := json.Unmarshal(buffer, msg); err != nil {
				glog.Error(err)
				continue
			}
			if v := commReceiveQueues.Get(msg.Group); v == nil {
				result = []byte(fmt.Sprintf("group [%s] does not exist", msg.Group))
				break
			} else {
				result = []byte("ok")
				if v := commReceiveQueues.Get(msg.Group); v != nil {
					v.(*gqueue.Queue).Push(msg)
				}
			}
		}
		if err == nil {
			if err := conn.SendPkg(result, option); err != nil {
				glog.Error(err)
			}
		} else {
			if err := conn.Close(); err != nil {
				glog.Error(err)
			}
			return
		}
	}
}
