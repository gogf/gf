// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

// Package gipv4 provides useful API for IPv4 address handling.
package gipv4

import (
	"encoding/binary"
	"fmt"
	"github.com/gogf/gf/g/text/gregex"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// 判断所给地址是否是一个IPv4地址
func Validate(ip string) bool {
	return gregex.IsMatchString(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`, ip)
}

// Get the IPv4 address corresponding to a given Internet host name.
func GetHostByName(hostname string) (string, error) {
	ips, err := net.LookupIP(hostname)
	if ips != nil {
		for _, v := range ips {
			if v.To4() != nil {
				return v.String(), nil
			}
		}
		return "", nil
	}
	return "", err
}

// Get a list of IPv4 addresses corresponding to a given Internet host name.
func GetHostsByName(hostname string) ([]string, error) {
	ips, err := net.LookupIP(hostname)
	if ips != nil {
		var ipStrs []string
		for _, v := range ips {
			if v.To4() != nil {
				ipStrs = append(ipStrs, v.String())
			}
		}
		return ipStrs, nil
	}
	return nil, err
}

// Get the Internet host name corresponding to a given IP address.
func GetNameByAddr(ipAddress string) (string, error) {
	names, err := net.LookupAddr(ipAddress)
	if names != nil {
		return strings.TrimRight(names[0], "."), nil
	}
	return "", err
}

// IP字符串转为整形.
func Ip2long(ipAddress string) uint32 {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return 0
	}
	return binary.BigEndian.Uint32(ip.To4())
}

// ip整形转为字符串
func Long2ip(properAddress uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, properAddress)
	return net.IP(ipByte).String()
}

// 获得ip的网段，例如：192.168.2.102 -> 192.168.2
func GetSegment(ip string) string {
	r := `^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$`
	reg, err := regexp.Compile(r)
	if err != nil {
		return ""
	}
	ips := reg.FindStringSubmatch(ip)
	if ips == nil {
		return ""
	}
	return fmt.Sprintf("%s.%s.%s", ips[1], ips[2], ips[3])
}

// 解析地址，形如：192.168.1.1:80 -> 192.168.1.1, 80
func ParseAddress(addr string) (string, int) {
	r := `^(.+):(\d+)$`
	reg, err := regexp.Compile(r)
	if err != nil {
		return "", 0
	}
	result := reg.FindStringSubmatch(addr)
	if result != nil {
		i, _ := strconv.Atoi(result[2])
		return result[1], i
	}
	return "", 0
}

// 获取本地局域网ip列表
func IntranetIP() (ips []string, err error) {
	ips = make([]string, 0)
	ifaces, e := net.Interfaces()
	if e != nil {
		return ips, e
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}

		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		// ignore warden bridge
		if strings.HasPrefix(iface.Name, "w-") {
			continue
		}

		addrs, e := iface.Addrs()
		if e != nil {
			return ips, e
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}

			ipStr := ip.String()
			if IsIntranet(ipStr) {
				ips = append(ips, ipStr)
			}
		}
	}
	return ips, nil
}

// 判断所给ip是否为局域网ip
// A类 10.0.0.0--10.255.255.255
// B类 172.16.0.0--172.31.255.255
// C类 192.168.0.0--192.168.255.255
func IsIntranet(ipStr string) bool {
	// ip协议保留的局域网ip
	if strings.HasPrefix(ipStr, "10.") || strings.HasPrefix(ipStr, "192.168.") {
		return true
	}
	if strings.HasPrefix(ipStr, "172.") {
		// 172.16.0.0 - 172.31.255.255
		arr := strings.Split(ipStr, ".")
		if len(arr) != 4 {
			return false
		}

		second, err := strconv.ParseInt(arr[1], 10, 64)
		if err != nil {
			return false
		}

		if second >= 16 && second <= 31 {
			return true
		}
	}

	return false
}
