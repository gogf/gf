package gcrc32

import (
    "hash/crc32"
)

func EncodeString(v string) uint32 {
    return crc32.ChecksumIEEE([]byte(v))
}

func EncodeBytes(v []byte) uint32 {
    return crc32.ChecksumIEEE(v)
}
