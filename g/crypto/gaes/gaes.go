// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// AES
package gaes

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "errors"
)

const (
    ivDefValue = "I Love Go Frame!"
)

func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    blockSize := block.BlockSize()
    plaintext = PKCS5Padding(plaintext, blockSize)
    iv := []byte(ivDefValue)
    blockMode  := cipher.NewCBCEncrypter(block, iv)
    ciphertext := make([]byte, len(plaintext))
    blockMode.CryptBlocks(ciphertext, plaintext)

    return ciphertext, nil
}

func Decrypt(cipherText []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    blockSize := block.BlockSize()
    if len(cipherText) < blockSize {
        return nil, errors.New("cipherText too short")
    }
    iv := []byte(ivDefValue)
    if len(cipherText)%blockSize != 0 {
        return nil, errors.New("cipherText is not a multiple of the block size")
    }
    blockModel := cipher.NewCBCDecrypter(block, iv)
    plaintext  := make([]byte, len(cipherText))
    blockModel.CryptBlocks(plaintext, cipherText)
    plaintext = PKCS5UnPadding(plaintext)

    return plaintext, nil
}

func PKCS5Padding(src []byte, blockSize int) []byte {
    padding := blockSize - len(src)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(src, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
    length := len(src)
    unpadding := int(src[length-1])
    return src[:(length - unpadding)]
}