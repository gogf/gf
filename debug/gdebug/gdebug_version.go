// Copyright 2019-2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdebug

import (
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/encoding/ghash"
	"io/ioutil"
	"strconv"
)

// BinVersion returns the version of current running binary.
// It uses ghash.BKDRHash+BASE36 algorithm to calculate the unique version of the binary.
func BinVersion() string {
	if binaryVersion == "" {
		binaryContent, _ := ioutil.ReadFile(selfPath)
		binaryVersion = strconv.FormatInt(
			int64(ghash.BKDRHash(binaryContent)),
			36,
		)
	}
	return binaryVersion
}

// BinVersionMd5 returns the version of current running binary.
// It uses MD5 algorithm to calculate the unique version of the binary.
func BinVersionMd5() string {
	if binaryVersionMd5 == "" {
		binaryVersionMd5, _ = gmd5.EncryptFile(selfPath)
	}
	return binaryVersionMd5
}
