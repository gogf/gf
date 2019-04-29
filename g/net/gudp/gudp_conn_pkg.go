// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gudp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

const (
	// 允许最大的简单协议包大小(byte), 15MB
	PKG_MAX_SIZE    = 0xFFFFFF
	// 消息包头大小: "总长度"3字节+"校验码"4字节
	PKG_HEADER_SIZE = 7
)

// 根据简单协议发送数据包。
// 简单协议数据格式：总长度(24bit)|校验码(32bit)|数据(变长)。
// 注意：
// 1. "总长度"包含自身3字节及"校验码"4字节。
// 2. 由于"总长度"为3字节，并且使用的BigEndian字节序，因此最后返回的buffer使用了buffer[1:]。
func (c *Conn) SendPkg(data []byte, retry...Retry) error {
	length := uint32(len(data))
	if length > PKG_MAX_SIZE - PKG_HEADER_SIZE {
		return errors.New(fmt.Sprintf(`data size %d exceeds max pkg size %d`, length, PKG_MAX_SIZE - PKG_HEADER_SIZE))
	}
	buffer := make([]byte, PKG_HEADER_SIZE + 1 + len(data))
	copy(buffer[PKG_HEADER_SIZE + 1 : ], data)
	binary.BigEndian.PutUint32(buffer[0 : ], PKG_HEADER_SIZE + length)
	binary.BigEndian.PutUint32(buffer[4 : ], Checksum(data))
	//fmt.Println("SendPkg:", buffer[1:])
	return c.Send(buffer[1:], retry...)
}

// 简单协议: 带超时时间的数据发送
func (c *Conn) SendPkgWithTimeout(data []byte, timeout time.Duration, retry...Retry) error {
	c.SetSendDeadline(time.Now().Add(timeout))
	defer c.SetSendDeadline(time.Time{})
	return c.SendPkg(data, retry...)
}

// 简单协议: 发送数据并等待接收返回数据
func (c *Conn) SendRecvPkg(data []byte, retry...Retry) ([]byte, error) {
	if err := c.SendPkg(data, retry...); err == nil {
		return c.RecvPkg(retry...)
	} else {
		return nil, err
	}
}

// 简单协议: 发送数据并等待接收返回数据(带返回超时等待时间)
func (c *Conn) SendRecvPkgWithTimeout(data []byte, timeout time.Duration, retry...Retry) ([]byte, error) {
	if err := c.SendPkg(data, retry...); err == nil {
		return c.RecvPkgWithTimeout(timeout, retry...)
	} else {
		return nil, err
	}
}

// 简单协议: 获取一个数据包。
func (c *Conn) RecvPkg(retry...Retry) (result []byte, err error) {
	var temp   []byte
	var length uint32
	for {
		// 先根据对象的缓冲区数据进行计算
		for {
			if len(c.buffer) >= PKG_HEADER_SIZE {
				// 注意"总长度"为3个字节，不满足4个字节的uint32类型，因此这里"低位"补0
				length = binary.BigEndian.Uint32([]byte{0, c.buffer[0], c.buffer[1], c.buffer[2]})
				// 解析的大小是否符合规范
				if length == 0 || length + PKG_HEADER_SIZE > PKG_MAX_SIZE {
					c.buffer = c.buffer[1:]
					continue
				}
				// 不满足包大小，需要继续读取
				if uint32(len(c.buffer)) < length {
					break
				}
				// 数据校验
				if binary.BigEndian.Uint32(c.buffer[3 : PKG_HEADER_SIZE]) != Checksum(c.buffer[PKG_HEADER_SIZE : length]) {
					c.buffer = c.buffer[1:]
					continue
				}
				result   = c.buffer[PKG_HEADER_SIZE : length]
				c.buffer = c.buffer[length: ]
				return
			} else {
				break
			}
		}
		// 读取系统socket缓冲区的完整数据
		temp, err = c.Recv(-1, retry...)
		if err != nil {
			break
		}
		if len(temp) > 0 {
			c.buffer = append(c.buffer, temp...)
		}
		//fmt.Println("RecvPkg:", c.buffer)
	}
	return
}

// 简单协议: 带超时时间的消息包获取
func (c *Conn) RecvPkgWithTimeout(timeout time.Duration, retry...Retry) ([]byte, error) {
	c.SetRecvDeadline(time.Now().Add(timeout))
	defer c.SetRecvDeadline(time.Time{})
	return c.RecvPkg(retry...)
}