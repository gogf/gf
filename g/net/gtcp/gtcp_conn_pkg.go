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
	"github.com/gogf/gf/g/crypto/gcrc32"
	"time"
)

const (
	// 默认允许最大的简单协议包大小(byte), 65535 byte
	gPKG_MAX_SIZE = 65535
)

// 数据读取选项
type Option struct {
	MaxSize    int    // (byte)数据读取的最大包大小，最大不能超过3字节(0xFFFFFF,15MB)，默认为65535byte
	Secret     []byte // (可选)安全通信密钥
	Retry      Retry  // 失败重试
}

// getPkgOption wraps and returns the option.
// If no option given, it returns a new option with default value.
func getPkgOption(option...Option) (*Option, error) {
	pkgOption := Option{}
	if len(option) > 0 {
		pkgOption = option[0]
	}
	if pkgOption.MaxSize == 0 {
		pkgOption.MaxSize = gPKG_MAX_SIZE
	} else if pkgOption.MaxSize > 0xFFFFFF {
		return nil, fmt.Errorf(`package size %d exceeds allowed max size %d`, pkgOption.MaxSize, 0xFFFFFF)
	}
	return &pkgOption, nil
}

// 根据简单协议发送数据包。
// 简单协议数据格式：总长度(24bit)|校验位(32bit,可选)|数据(变长)。
// 注意：
// 1. "总长度"包含自身3字节及"校验位"4字节(可选)。
// 2. 当Secret有提供时，"校验位"才会存在，否则该字段为空。
// 3. "校验位"提供简单的数据完整性及防篡改校验，默认没有开启。
// 4. 由于"总长度"为3字节，并且使用的BigEndian字节序，因此这里最后返回的buffer使用了buffer[1:]。
func (c *Conn) SendPkg(data []byte, option...Option) error {
	pkgOption, err := getPkgOption(option...)
	if err != nil {
		return err
	}
	headerSize := 3
	if len(pkgOption.Secret) > 0 {
		headerSize = 7
	}
	length := len(data)
	if length > pkgOption.MaxSize - headerSize {
		return errors.New(fmt.Sprintf(`data size %d exceeds max pkg size %d`, length, gPKG_MAX_SIZE - headerSize))
	}

	buffer := make([]byte, headerSize + 1 + len(data))
	copy(buffer[headerSize + 1 : ], data)
	binary.BigEndian.PutUint32(buffer[0 : ], uint32(headerSize + length))
	if len(pkgOption.Secret) > 0 {
		binary.BigEndian.PutUint32(buffer[4 : ], gcrc32.Encrypt(append(data, pkgOption.Secret...)))
	}
	if pkgOption.Retry.Count > 0 {
		c.Send(buffer[1:], pkgOption.Retry)
	}
	return c.Send(buffer[1:])
}

// 简单协议: 带超时时间的数据发送
func (c *Conn) SendPkgWithTimeout(data []byte, timeout time.Duration, option...Option) error {
	c.SetSendDeadline(time.Now().Add(timeout))
	defer c.SetSendDeadline(time.Time{})
	return c.SendPkg(data, option...)
}

// 简单协议: 发送数据并等待接收返回数据
func (c *Conn) SendRecvPkg(data []byte, option...Option) ([]byte, error) {
	if err := c.SendPkg(data, option...); err == nil {
		return c.RecvPkg(option...)
	} else {
		return nil, err
	}
}

// 简单协议: 发送数据并等待接收返回数据(带返回超时等待时间)
func (c *Conn) SendRecvPkgWithTimeout(data []byte, timeout time.Duration, option...Option) ([]byte, error) {
	if err := c.SendPkg(data, option...); err == nil {
		return c.RecvPkgWithTimeout(timeout, option...)
	} else {
		return nil, err
	}
}

// 简单协议: 获取一个数据包。
func (c *Conn) RecvPkg(option...Option) (result []byte, err error) {
	var temp   []byte
	var length int
	pkgOption, err := getPkgOption(option...)
	if err != nil {
		return nil, err
	}
	headerSize := 3
	if len(pkgOption.Secret) > 0 {
		headerSize = 7
	}
	for {
		// 先根据对象的缓冲区数据进行计算
		for {
			if len(c.buffer) >= headerSize {
				// 注意"总长度"为3个字节，不满足4个字节的uint32类型，因此这里"低位"补0
				length = int(binary.BigEndian.Uint32([]byte{0, c.buffer[0], c.buffer[1], c.buffer[2]}))
				// 解析的大小是否符合规范，清空从该连接接收到的所有数据包
				if length <= 0 || length + headerSize > pkgOption.MaxSize {
					c.buffer = c.buffer[:0]
					return nil, fmt.Errorf(`invalid package size %d`, length)
				}
				// 不满足包大小，需要继续读取
				if len(c.buffer) < length {
					break
				}
				// 数据校验，如果失败，丢弃该数据包
				receivedCrc32   := binary.BigEndian.Uint32(c.buffer[3 : headerSize])
				calculatedCrc32 := gcrc32.Encrypt(c.buffer[headerSize : length])
				if receivedCrc32 != calculatedCrc32 {
					c.buffer = c.buffer[length: ]
					return nil, fmt.Errorf(`data CRC32 validates failed, received %d, caculated %d`, receivedCrc32, calculatedCrc32)
				}
				result   = c.buffer[headerSize : length]
				c.buffer = c.buffer[length: ]
				return
			} else {
				break
			}
		}
		// 读取系统socket缓冲区的完整数据
		temp, err = c.Recv(-1, option...)
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
func (c *Conn) RecvPkgWithTimeout(timeout time.Duration, option...Option) ([]byte, error) {
	c.SetRecvDeadline(time.Now().Add(timeout))
	defer c.SetRecvDeadline(time.Time{})
	return c.RecvPkg(retry...)
}