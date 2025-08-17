// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

// Package gipv4 provides useful API for IPv4 address handling.

package gipv4

import (
	"encoding/binary"
	"net"
)

// IpToLongBigEndian converts ip address to an uint32 integer with big endian.
func IpToLongBigEndian(ip string) uint32 {
	netIp := net.ParseIP(ip)
	if netIp == nil {
		return 0
	}
	return binary.BigEndian.Uint32(netIp.To4())
}

// LongToIpBigEndian converts an uint32 integer ip address to its string type address with big endian.
func LongToIpBigEndian(long uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, long)
	return net.IP(ipByte).String()
}

// IpToLongLittleEndian converts ip address to an uint32 integer with little endian.
func IpToLongLittleEndian(ip string) uint32 {
	netIp := net.ParseIP(ip)
	if netIp == nil {
		return 0
	}
	return binary.LittleEndian.Uint32(netIp.To4())
}

// LongToIpLittleEndian converts an uint32 integer ip address to its string type address with little endian.
func LongToIpLittleEndian(long uint32) string {
	ipByte := make([]byte, 4)
	binary.LittleEndian.PutUint32(ipByte, long)
	return net.IP(ipByte).String()
}

// Ip2long converts ip address to an uint32 integer.
// Deprecated: Use IpToLongBigEndian instead.
func Ip2long(ip string) uint32 {
	return IpToLongBigEndian(ip)
}

// Long2ip converts an uint32 integer ip address to its string type address.
// Deprecated: Use LongToIpBigEndian instead.
func Long2ip(long uint32) string {
	return LongToIpBigEndian(long)
}
