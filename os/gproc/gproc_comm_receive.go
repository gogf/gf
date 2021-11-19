// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gproc

import (
	"context"
	"fmt"
	"net"

	"github.com/gogf/gf/v2/container/gqueue"
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	// tcpListened marks whether the receiving listening service started.
	tcpListened = gtype.NewBool()
)

// Receive blocks and receives message from other process using local TCP listening.
// Note that, it only enables the TCP listening service when this function called.
func Receive(group ...string) *MsgRequest {
	// Use atomic operations to guarantee only one receiver goroutine listening.
	if tcpListened.Cas(false, true) {
		go receiveTcpListening()
	}
	var groupName string
	if len(group) > 0 {
		groupName = group[0]
	} else {
		groupName = defaultGroupNameForProcComm
	}
	queue := commReceiveQueues.GetOrSetFuncLock(groupName, func() interface{} {
		return gqueue.New(maxLengthForProcMsgQueue)
	}).(*gqueue.Queue)

	// Blocking receiving.
	if v := queue.Pop(); v != nil {
		return v.(*MsgRequest)
	}
	return nil
}

// receiveTcpListening scans local for available port and starts listening.
func receiveTcpListening() {
	var listen *net.TCPListener
	// Scan the available port for listening.
	for i := defaultTcpPortForProcComm; ; i++ {
		addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", i))
		if err != nil {
			continue
		}
		listen, err = net.ListenTCP("tcp", addr)
		if err != nil {
			continue
		}
		// Save the port to the pid file.
		if err := gfile.PutContents(getCommFilePath(Pid()), gconv.String(i)); err != nil {
			panic(err)
		}
		break
	}
	// Start listening.
	for {
		if conn, err := listen.Accept(); err != nil {
			glog.Error(context.TODO(), err)
		} else if conn != nil {
			go receiveTcpHandler(gtcp.NewConnByNetConn(conn))
		}
	}
}

// receiveTcpHandler is the connection handler for receiving data.
func receiveTcpHandler(conn *gtcp.Conn) {
	var (
		ctx      = context.TODO()
		result   []byte
		response MsgResponse
	)
	for {
		response.Code = 0
		response.Message = ""
		response.Data = nil
		buffer, err := conn.RecvPkg()
		if len(buffer) > 0 {
			// Package decoding.
			msg := new(MsgRequest)
			if err := json.UnmarshalUseNumber(buffer, msg); err != nil {
				// glog.Error(err)
				continue
			}
			if msg.RecvPid != Pid() {
				// Not mine package.
				response.Message = fmt.Sprintf("receiver pid not match, target: %d, current: %d", msg.RecvPid, Pid())
			} else if v := commReceiveQueues.Get(msg.Group); v == nil {
				// Group check.
				response.Message = fmt.Sprintf("group [%s] does not exist", msg.Group)
			} else {
				// Push to buffer queue.
				response.Code = 1
				v.(*gqueue.Queue).Push(msg)
			}
		} else {
			// Empty package.
			response.Message = "empty package"
		}
		if err == nil {
			result, err = json.Marshal(response)
			if err != nil {
				glog.Error(ctx, err)
			}
			if err := conn.SendPkg(result); err != nil {
				glog.Error(ctx, err)
			}
		} else {
			// Just close the connection if any error occurs.
			if err := conn.Close(); err != nil {
				glog.Error(ctx, err)
			}
			break
		}
	}
}
