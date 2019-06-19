// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gaes provides useful API for AES encryption/decryption algorithms.
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

// Encrypt is alias of EncryptCBC.
func Encrypt(plainText []byte, key []byte, iv ...[]byte) ([]byte, error) {
	return EncryptCBC(plainText, key, iv...)
}

// Decrypt is alias of DecryptCBC.
func Decrypt(cipherText []byte, key []byte, iv ...[]byte) ([]byte, error) {
	return DecryptCBC(cipherText, key, iv...)
}

// AES加密, 使用CBC模式，注意key必须为16/24/32位长度，iv初始化向量为非必需参数。
func EncryptCBC(plainText []byte, key []byte, iv ...[]byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	plainText = PKCS5Padding(plainText, blockSize)
	ivValue := ([]byte)(nil)
	if len(iv) > 0 {
		ivValue = iv[0]
	} else {
		ivValue = []byte(ivDefValue)
	}
	blockMode := cipher.NewCBCEncrypter(block, ivValue)
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)

	return cipherText, nil
}

// AES解密, 使用CBC模式，注意key必须为16/24/32位长度，iv初始化向量为非必需参数
func DecryptCBC(cipherText []byte, key []byte, iv ...[]byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(cipherText) < blockSize {
		return nil, errors.New("cipherText too short")
	}
	ivValue := ([]byte)(nil)
	if len(iv) > 0 {
		ivValue = iv[0]
	} else {
		ivValue = []byte(ivDefValue)
	}
	if len(cipherText)%blockSize != 0 {
		return nil, errors.New("cipherText is not a multiple of the block size")
	}
	blockModel := cipher.NewCBCDecrypter(block, ivValue)
	plainText := make([]byte, len(cipherText))
	blockModel.CryptBlocks(plainText, cipherText)
	plainText, e := PKCS5UnPadding(plainText, blockSize)
	if e != nil {
		return nil, e
	}
	return plainText, nil
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS5UnPadding(src []byte, blockSize int) ([]byte, error) {
	length := len(src)
	if blockSize <= 0 {
		return nil, errors.New("invalid blocklen")
	}

	if length%blockSize != 0 || length == 0 {
		return nil, errors.New("invalid data len")
	}

	unpadding := int(src[length-1])
	if unpadding > blockSize || unpadding == 0 {
		return nil, errors.New("invalid padding")
	}

	padding := src[length-unpadding:]
	for i := 0; i < unpadding; i++ {
		if padding[i] != byte(unpadding) {
			return nil, errors.New("invalid padding")
		}
	}

	return src[:(length - unpadding)], nil
}

// AES加密, 使用CFB模式。
// 注意key必须为16/24/32位长度，padding返回补位长度，iv初始化向量为非必需参数。
func EncryptCFB(plainText []byte, key []byte, padding *int, iv ...[]byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	plainText, *padding = ZeroPadding(plainText, blockSize) //补位0
	ivValue := ([]byte)(nil)
	if len(iv) > 0 {
		ivValue = iv[0]
	} else {
		ivValue = []byte(ivDefValue)
	}
	stream := cipher.NewCFBEncrypter(block, ivValue)
	cipherText := make([]byte, len(plainText))
	stream.XORKeyStream(cipherText, plainText)
	return cipherText, nil
}

// AES解密, 使用CFB模式。
// 注意key必须为16/24/32位长度，unpadding为去补位长度，iv初始化向量为非必需参数。
func DecryptCFB(cipherText []byte, key []byte, unpadding int, iv ...[]byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(cipherText) < aes.BlockSize {
		return nil, errors.New("cipherText too short")
	}
	ivValue := ([]byte)(nil)
	if len(iv) > 0 {
		ivValue = iv[0]
	} else {
		ivValue = []byte(ivDefValue)
	}
	stream := cipher.NewCFBDecrypter(block, ivValue)
	plainText := make([]byte, len(cipherText))
	stream.XORKeyStream(plainText, cipherText)
	plainText = ZeroUnPadding(plainText, unpadding) //去补位0
	return plainText, nil
}

func ZeroPadding(ciphertext []byte, blockSize int) ([]byte, int) {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(0)}, padding)
	return append(ciphertext, padtext...), padding
}

func ZeroUnPadding(plaintext []byte, unpadding int) []byte {
	length := len(plaintext)
	return plaintext[:(length - unpadding)]
}
