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
	gPKG_HEADER_SIZE_DEFAULT = 2 // Header size for simple package protocol.
	gPKG_HEADER_SIZE_MAX     = 4 // Max header size for simple package protocol.
)

// Package option for simple protocol.
type PkgOption struct {
	// HeaderSize is used to mark the data length for next data receiving.
	// It's 2 bytes in default, 4 bytes max, which stands for the max data length
	// from 65535 to 4294967295 bytes.
	HeaderSize int

	// MaxDataSize is the data field size in bytes for data length validation.
	// If it's not manually set, it'll automatically be set correspondingly with the HeaderSize.
	MaxDataSize int

	// Retry policy when operation fails.
	Retry Retry
}

// SendPkg send data using simple package protocol.
//
// Simple package protocol: DataLength(24bit)|DataField(variant)。
//
// Note that,
// 1. The DataLength is the length of DataField, which does not contain the header size.
// 2. The integer bytes of the package are encoded using BigEndian order.
func (c *Conn) SendPkg(data []byte, option ...PkgOption) error {
	pkgOption, err := getPkgOption(option...)
	if err != nil {
		return err
	}
	length := len(data)
	if length > pkgOption.MaxDataSize {
		return fmt.Errorf(
			`data too long, data size %d exceeds allowed max data size %d`,
			length, pkgOption.MaxDataSize,
		)
	}
	offset := gPKG_HEADER_SIZE_MAX - pkgOption.HeaderSize
	buffer := make([]byte, gPKG_HEADER_SIZE_MAX+len(data))
	binary.BigEndian.PutUint32(buffer[0:], uint32(length))
	copy(buffer[gPKG_HEADER_SIZE_MAX:], data)
	if pkgOption.Retry.Count > 0 {
		return c.Send(buffer[offset:], pkgOption.Retry)
	}
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

// RecvPkg receives data from connection using simple package protocol.
func (c *Conn) RecvPkg(option ...PkgOption) (result []byte, err error) {
	var buffer []byte
	var length int
	pkgOption, err := getPkgOption(option...)
	if err != nil {
		return nil, err
	}
	// Header field.
	buffer, err = c.Recv(pkgOption.HeaderSize, pkgOption.Retry)
	if err != nil {
		return nil, err
	}
	switch pkgOption.HeaderSize {
	case 1:
		// It fills with zero if the header size is lesser than 4 bytes (uint32).
		length = int(binary.BigEndian.Uint32([]byte{0, 0, 0, buffer[0]}))
	case 2:
		length = int(binary.BigEndian.Uint32([]byte{0, 0, buffer[0], buffer[1]}))
	case 3:
		length = int(binary.BigEndian.Uint32([]byte{0, buffer[0], buffer[1], buffer[2]}))
	default:
		length = int(binary.BigEndian.Uint32([]byte{buffer[0], buffer[1], buffer[2], buffer[3]}))
	}
	// It here validates the size of the package.
	// It clears the buffer and returns error immediately if it validates failed.
	if length < 0 || length > pkgOption.MaxDataSize {
		return nil, fmt.Errorf(`invalid package size %d`, length)
	}
	// Empty package.
	if length == 0 {
		return nil, nil
	}
	// Data field.
	return c.Recv(length, pkgOption.Retry)
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

// getPkgOption wraps and returns the PkgOption.
// If no option given, it returns a new option with default value.
func getPkgOption(option ...PkgOption) (*PkgOption, error) {
	pkgOption := PkgOption{}
	if len(option) > 0 {
		pkgOption = option[0]
	}
	if pkgOption.HeaderSize == 0 {
		pkgOption.HeaderSize = gPKG_HEADER_SIZE_DEFAULT
	}
	if pkgOption.HeaderSize > gPKG_HEADER_SIZE_MAX {
		return nil, fmt.Errorf(
			`package header size %d definition exceeds max header size %d`,
			pkgOption.HeaderSize, gPKG_HEADER_SIZE_MAX,
		)
	}
	if pkgOption.MaxDataSize == 0 {
		switch pkgOption.HeaderSize {
		case 1:
			pkgOption.MaxDataSize = 0xFF
		case 2:
			pkgOption.MaxDataSize = 0xFFFF
		case 3:
			pkgOption.MaxDataSize = 0xFFFFFF
		case 4:
			// math.MaxInt32 not math.MaxUint32
			pkgOption.MaxDataSize = 0x7FFFFFFF
		}
	}
	if pkgOption.MaxDataSize > 0x7FFFFFFF {
		return nil, fmt.Errorf(
			`package data size %d definition exceeds allowed max data size %d`,
			pkgOption.MaxDataSize, 0x7FFFFFFF,
		)
	}
	return &pkgOption, nil
}
