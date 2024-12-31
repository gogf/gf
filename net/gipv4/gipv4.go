// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

// Package gipv4 provides useful API for IPv4 address handling.

package gipv4

import (
	"net"
	"strconv"
	"strings"
)

// Validate checks whether given `ip` a valid IPv4 address.
func Validate(ip string) bool {
	parsed := net.ParseIP(ip)
	return parsed != nil && parsed.To4() != nil
}

// ParseAddress parses `address` to its ip and port.
// Eg: 192.168.1.1:80 -> 192.168.1.1, 80
func ParseAddress(address string) (string, int) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return "", 0
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return "", 0
	}
	return host, portInt
}

// GetSegment returns the segment of given ip address.
// Eg: 192.168.2.102 -> 192.168.2
func GetSegment(ip string) string {
	if !Validate(ip) {
		return ""
	}
	segments := strings.Split(ip, ".")
	if len(segments) != 4 {
		return ""
	}
	return strings.Join(segments[:3], ".")
}
