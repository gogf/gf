// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package guid

import (
	"errors"
	"github.com/gogf/gf/container/gtype"
	"net"
	"time"
)

const (
	BitLenTime      = 39                               // Bit length of time.
	BitLenSequence  = 8                                // Bit length of sequence number.
	BitLenMachineID = 63 - BitLenTime - BitLenSequence // Bit length of machine id.
)

var (
	// MachineId is 16 bits number composed with the last two parts of the first local private ip.
	// You can change it as your custom machine id in boot time, but do not change it in runtime.
	// Note that if there's no net card on the machine, it panics when the process boots.
	MachineId = getLower16BitPrivateIP()

	// sequenceNumber is used for internal concurrent-safe sequence number counting.
	sequenceNumber = gtype.NewInt()
)

// I creates and returns an uint64 id which using improved SnowFlake algorithm.
// An Improved SnowFlake ID is composed of:
//     39 bits for time in units of 10 msec
//     16 bits for a machine id
//      8 bits for a sequence number.
func I() uint64 {
	return uint64(time.Now().UnixNano())<<(BitLenMachineID+BitLenSequence) | uint64(MachineId<<BitLenSequence) | uint64(sequenceNumber.Add(1)%0xFF)
}

func getLower16BitPrivateIP() uint16 {
	if MachineId != 0 {
		return MachineId
	}
	ip, err := getPrivateIPv4()
	if err != nil {
		panic(err)
	}
	MachineId = uint16(ip[2])<<8 + uint16(ip[3])
	return MachineId
}

func getPrivateIPv4() (net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, a := range as {
		ipNet, ok := a.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() {
			continue
		}

		ip := ipNet.IP.To4()
		if isPrivateIPv4(ip) {
			return ip, nil
		}
	}
	return nil, errors.New("no private ip address")
}

func isPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}
