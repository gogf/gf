// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtcp

import (
	"encoding/binary"
	"fmt"
	"time"
)

const (
	gPKG_DEFAULT_MAX_DATA_SIZE = 65535 // (Byte) Max package size.
	gPKG_DEFAULT_HEADER_SIZE   = 2     // Header size for simple package protocol.
	gPKG_MAX_HEADER_SIZE       = 4     // Max header size for simple package protocol.
)

// Package option for simple protocol.
type PkgOption struct {
	HeaderSize  int   // It's 2 bytes in default, max is 4 bytes.
	MaxDataSize int   // (Byte)data field size, it's 2 bytes in default, which means 65535 bytes.
	Retry       Retry // Retry policy.
}

// getPkgOption wraps and returns the PkgOption.
// If no option given, it returns a new option with default value.
func getPkgOption(option ...PkgOption) (*PkgOption, error) {
	pkgOption := PkgOption{}
	if len(option) > 0 {
		pkgOption = option[0]
	}
	if pkgOption.HeaderSize == 0 {
		pkgOption.HeaderSize = gPKG_DEFAULT_HEADER_SIZE
	}
	if pkgOption.MaxDataSize == 0 {
		pkgOption.MaxDataSize = gPKG_DEFAULT_MAX_DATA_SIZE
	} else if pkgOption.MaxDataSize > 0xFFFFFF {
		return nil, fmt.Errorf(`package size %d exceeds allowed max size %d`, pkgOption.MaxDataSize, 0xFFFFFF)
	}
	return &pkgOption, nil
}

// SendPkg send data using simple package protocol.
//
// Simple package protocol: DataLength(24bit)|DataField(variant)。
//
// Note that,
// 1. The DataLength is the length of DataField, which does not contain the header size 2 bytes.
// 2. The integer bytes of the package are encoded using BigEndian order.
func (c *Conn) SendPkg(data []byte, option ...PkgOption) error {
	pkgOption, err := getPkgOption(option...)
	if err != nil {
		return err
	}
	length := len(data)
	if length > pkgOption.MaxDataSize {
		return fmt.Errorf(`data size %d exceeds max pkg size %d`, length, pkgOption.MaxDataSize)
	}
	offset := gPKG_MAX_HEADER_SIZE - pkgOption.HeaderSize
	buffer := make([]byte, gPKG_MAX_HEADER_SIZE+len(data))
	binary.BigEndian.PutUint32(buffer[0:], uint32(length))
	copy(buffer[gPKG_MAX_HEADER_SIZE:], data)
	if pkgOption.Retry.Count > 0 {
		return c.Send(buffer[offset:], pkgOption.Retry)
	}
	//fmt.Println("SendPkg:", buffer[offset:])
	return c.Send(buffer[offset:])
}

// SendPkgWithTimeout writes data to connection with timeout using simple package protocol.
func (c *Conn) SendPkgWithTimeout(data []byte, timeout time.Duration, option ...PkgOption) (err error) {
	if err := c.SetSendDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}
	defer c.SetSendDeadline(time.Time{})
	err = c.SendPkg(data, option...)
	return
}

// SendRecvPkg writes data to connection and blocks reading response using simple package protocol.
func (c *Conn) SendRecvPkg(data []byte, option ...PkgOption) ([]byte, error) {
	if err := c.SendPkg(data, option...); err == nil {
		return c.RecvPkg(option...)
	} else {
		return nil, err
	}
}

// SendRecvPkgWithTimeout writes data to connection and reads response with timeout using simple package protocol.
func (c *Conn) SendRecvPkgWithTimeout(data []byte, timeout time.Duration, option ...PkgOption) ([]byte, error) {
	if err := c.SendPkg(data, option...); err == nil {
		return c.RecvPkgWithTimeout(timeout, option...)
	} else {
		return nil, err
	}
}

// Recv receives data from connection using simple package protocol.
func (c *Conn) RecvPkg(option ...PkgOption) (result []byte, err error) {
	var temp []byte
	var length int
	pkgOption, err := getPkgOption(option...)
	if err != nil {
		return nil, err
	}
	for {
		for {
			if len(c.buffer) >= pkgOption.HeaderSize {
				if length <= 0 {
					switch pkgOption.HeaderSize {
					case 1:
						// It fills with zero if the header size is lesser than 4 bytes (uint32).
						length = int(binary.BigEndian.Uint32([]byte{0, 0, 0, c.buffer[0]}))
					case 2:
						length = int(binary.BigEndian.Uint32([]byte{0, 0, c.buffer[0], c.buffer[1]}))
					case 3:
						length = int(binary.BigEndian.Uint32([]byte{0, c.buffer[0], c.buffer[1], c.buffer[2]}))
					default:
						length = int(binary.BigEndian.Uint32([]byte{c.buffer[0], c.buffer[1], c.buffer[2], c.buffer[3]}))
					}
				}
				// It here validates the size of the package.
				// It clears the buffer and returns error immediately if it validates failed.
				if length < 0 || length > pkgOption.MaxDataSize {
					c.buffer = c.buffer[:0]
					return nil, fmt.Errorf(`invalid package size %d`, length)
				}
				// It continues reading until it receives complete bytes of the package.
				if len(c.buffer) < length+pkgOption.HeaderSize {
					break
				}
				result = c.buffer[pkgOption.HeaderSize : pkgOption.HeaderSize+length]
				c.buffer = c.buffer[pkgOption.HeaderSize+length:]
				length = 0
				return
			} else {
				break
			}
		}
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

// RecvPkgWithTimeout reads data from connection with timeout using simple package protocol.
func (c *Conn) RecvPkgWithTimeout(timeout time.Duration, option ...PkgOption) (data []byte, err error) {
	if err := c.SetRecvDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}
	defer c.SetRecvDeadline(time.Time{})
	data, err = c.RecvPkg(option...)
	return
}
