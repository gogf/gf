// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gcrc32 provides useful API for CRC32 encryption/decryption algorithms.
package gcrc32

import (
    "hash/crc32"
)

func EncryptString(v string) uint32 {
    return crc32.ChecksumIEEE([]byte(v))
}

func EncryptBytes(v []byte) uint32 {
    return crc32.ChecksumIEEE(v)
}
