// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtcp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

const (
	// 默认允许最大的简单协议包大小(byte), 65535 byte
	gPKG_MAX_DATA_SIZE = 65535
	// 简单协议包头大小
	gPKG_HEADER_SIZE   = 3
)

// 数据读取选项
type PkgOption struct {
	MaxSize int    // (byte)数据读取的最大包大小，最大不能超过3字节(0xFFFFFF,15MB)，默认为65535byte
	Retry   Retry  // 失败重试
}

// getPkgOption wraps and returns the PkgOption.
// If no option given, it returns a new option with default value.
func getPkgOption(option...PkgOption) (*PkgOption, error) {
	pkgOption := PkgOption{}
	if len(option) > 0 {
		pkgOption = option[0]
	}
	if pkgOption.MaxSize == 0 {
		pkgOption.MaxSize = gPKG_MAX_DATA_SIZE
	} else if pkgOption.MaxSize > 0xFFFFFF {
		return nil, fmt.Errorf(`package size %d exceeds allowed max size %d`, pkgOption.MaxSize, 0xFFFFFF)
	}
	return &pkgOption, nil
}

// 根据简单协议发送数据包。
//
// 简单协议数据格式：数据长度(24bit)|数据字段(变长)。
//
// 注意：
// 1. "数据长度"仅为"数据字段"的长度，不包含头信息的长度字段3字节。
// 2. 由于"数据长度"为3字节，并且使用的BigEndian字节序，因此这里最后返回的buffer使用了buffer[1:]。
func (c *Conn) SendPkg(data []byte, option...PkgOption) error {
	pkgOption, err := getPkgOption(option...)
	if err != nil {
		return err
	}
	length := len(data)
	if length > pkgOption.MaxSize {
		return errors.New(fmt.Sprintf(`data size %d exceeds max pkg size %d`, length, gPKG_MAX_DATA_SIZE))
	}
	buffer := make([]byte, gPKG_HEADER_SIZE + 1 + len(data))
	binary.BigEndian.PutUint32(buffer[0 : ], uint32(length))
	copy(buffer[gPKG_HEADER_SIZE + 1 : ], data)
	if pkgOption.Retry.Count > 0 {
		return c.Send(buffer[1:], pkgOption.Retry)
	}
	//fmt.Println("SendPkg:", buffer[1:])
	return c.Send(buffer[1:])
}

// 简单协议: 带超时时间的数据发送
func (c *Conn) SendPkgWithTimeout(data []byte, timeout time.Duration, option...PkgOption) error {
	if err := c.SetSendDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}
	defer c.SetSendDeadline(time.Time{})
	return c.SendPkg(data, option...)
}

// 简单协议: 发送数据并等待接收返回数据
func (c *Conn) SendRecvPkg(data []byte, option...PkgOption) ([]byte, error) {
	if err := c.SendPkg(data, option...); err == nil {
		return c.RecvPkg(option...)
	} else {
		return nil, err
	}
}

// 简单协议: 发送数据并等待接收返回数据(带返回超时等待时间)
func (c *Conn) SendRecvPkgWithTimeout(data []byte, timeout time.Duration, option...PkgOption) ([]byte, error) {
	if err := c.SendPkg(data, option...); err == nil {
		return c.RecvPkgWithTimeout(timeout, option...)
	} else {
		return nil, err
	}
}

// 简单协议: 获取一个数据包。
func (c *Conn) RecvPkg(option...PkgOption) (result []byte, err error) {
	var temp   []byte
	var length int
	pkgOption, err := getPkgOption(option...)
	if err != nil {
		return nil, err
	}
	for {
		// 先根据对象的缓冲区数据进行计算
		for {
			if len(c.buffer) >= gPKG_HEADER_SIZE {
				// 注意"数据长度"为3个字节，不满足4个字节的uint32类型，因此这里"低位"补0
				length = int(binary.BigEndian.Uint32([]byte{0, c.buffer[0], c.buffer[1], c.buffer[2]}))
				// 解析的大小是否符合规范，清空从该连接接收到的所有数据包
				if length < 0 || length > pkgOption.MaxSize {
					c.buffer = c.buffer[:0]
					return nil, fmt.Errorf(`invalid package size %d`, length)
				}
				// 不满足包大小，需要继续读取
				if len(c.buffer) < length + gPKG_HEADER_SIZE {
					break
				}
				result   = c.buffer[gPKG_HEADER_SIZE : gPKG_HEADER_SIZE + length]
				c.buffer = c.buffer[gPKG_HEADER_SIZE + length: ]
				return
			} else {
				break
			}
		}
		// 读取系统socket缓冲区的完整数据
		temp, err = c.Recv(0, pkgOption.Retry)
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
func (c *Conn) RecvPkgWithTimeout(timeout time.Duration, option...PkgOption) ([]byte, error) {
	if err := c.SetRecvDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}
	defer c.SetRecvDeadline(time.Time{})
	return c.RecvPkg(option...)
}