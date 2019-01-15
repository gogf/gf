// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// @author: wenzi1<liyz23@qq.com>

// Package gdes provides useful API for DES encryption/decryption algorithms.
package gdes

import (
	"crypto/des"
	"crypto/cipher"
	"errors"
	"bytes"
)

const (
	NOPADDING     = iota
	PKCS5PADDING
)

//ECB模式DES加密
func DesECBEncrypt(key []byte, clearText []byte, padding int) ([]byte, error) {
	text, err := Padding(clearText, padding)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, len(text))
	
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	for i, count := 0, len(text)/blockSize; i < count; i++ {
		begin, end := i * blockSize, i * blockSize + blockSize
		block.Encrypt(cipherText[begin:end], text[begin:end])
	}
	return cipherText, nil
}

//ECB模式DES解密
func DesECBDecrypt(key []byte, cipherText []byte, padding int) ([]byte, error) {
	text := make([]byte, len(cipherText))
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	for i, count := 0, len(text)/blockSize; i < count; i++{
		begin, end := i * blockSize, i * blockSize + blockSize
		block.Decrypt(text[begin:end], cipherText[begin:end])
	}

	clearText, err := UnPadding(text, padding)
	if err != nil {
		return nil, err
	}
	return clearText, nil
}

//ECB模式3DES加密，密钥长度可以是16或24位长
func TripleDesECBEncrypt(key []byte, clearText []byte, padding int) ( []byte, error) {
	if len(key) != 16 && len(key) != 24 {
		return nil, errors.New("key length error")
	}

	text, err := Padding(clearText, padding)
	if err != nil {
		return nil, err
	}

	newKey := make([]byte, 0)
	if len(key) == 16 {
		newKey = append([]byte{}, key...)
		newKey = append(newKey, key[:8]...)
	} else {
		newKey = append([]byte{}, key...)
	}

	block, err := des.NewTripleDESCipher(newKey)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	cipherText := make([]byte, len(text))
	for i, count := 0, len(text) / blockSize; i < count; i++{
		begin, end := i * blockSize, i * blockSize + blockSize
		block.Encrypt(cipherText[begin:end], text[begin:end])
	}
	return cipherText, nil
}

//ECB模式3DES解密，密钥长度可以是16或24位长
func TripleDesECBDecrypt(key []byte, cipherText []byte, padding int) ([]byte,  error) {
	if len(key) != 16 && len(key) != 24 {
		return nil, errors.New("key length error")
	}

	newKey := make([]byte, 0)
	if len(key) == 16 {
		newKey = append([]byte{}, key...)
		newKey = append(newKey, key[:8]...)
	} else {
		newKey = append([]byte{}, key...)
	}

	block, err := des.NewTripleDESCipher(newKey)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	text := make([]byte, len(cipherText))
	for i, count := 0, len(text) / blockSize; i < count; i++ {
		begin, end := i * blockSize, i * blockSize + blockSize
		block.Decrypt(text[begin:end], cipherText[begin:end])
	}

	clearText, err := UnPadding(text, padding)
	if err != nil {
		return nil, err
	}
	return clearText, nil
}

//CBC模式DES加密
func DesCBCEncrypt(key []byte, clearText []byte, iv []byte, padding int) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(iv) != block.BlockSize() {
		return nil, errors.New("iv length invalid")
	}

	text, err := Padding(clearText, padding)
	if err != nil {
		return nil, err
	}
	cipherText := make([]byte, len(text))

	encrypter := cipher.NewCBCEncrypter(block, iv)
	encrypter.CryptBlocks(cipherText, text)

	return cipherText, nil
}

//CBC模式DES解密
func DesCBCDecrypt(key []byte, cipherText []byte, iv []byte, padding int) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(iv) != block.BlockSize() {
		return nil, errors.New("iv length invalid")
	}

	text := make([]byte, len(cipherText))
	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypter.CryptBlocks(text, cipherText)

	clearText, err := UnPadding(text, padding)
	if err != nil {
		return nil, err
	}

	return clearText, nil
}

//CBC模式3DES加密
func TripleDesCBCEncrypt(key []byte, clearText []byte, iv []byte, padding int) ([]byte, error) {
	if len(key) != 16 && len(key) != 24 {
		return nil, errors.New("key length invalid")
	}

	newKey := make([]byte, 0)
	if len(key) == 16 {
		newKey = append([]byte{}, key...)
		newKey = append(newKey, key[:8]...)
	} else {
		newKey = append([]byte{}, key...)
	}

	block, err := des.NewTripleDESCipher(newKey)
	if err != nil {
		return nil, err
	}

	if len(iv) != block.BlockSize() {
		return nil, errors.New("iv length invalid")
	}

	text, err := Padding(clearText, padding)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, len(text))
	encrypter := cipher.NewCBCEncrypter(block, iv)
	encrypter.CryptBlocks(cipherText, text)

	return cipherText, nil
}

//CBC模式3DES解密
func TripleDesCBCDecrypt(key []byte, cipherText []byte, iv []byte, padding int) ( []byte,  error) {
	if len(key) != 16 && len(key) != 24 {
		return nil, errors.New("key length invalid")
	}

	newKey := make([]byte, 0)
	if len(key) == 16 {
		newKey = append([]byte{}, key...)
		newKey = append(newKey, key[:8]...)
	} else {
		newKey = append([]byte{}, key...)
	}

	block, err := des.NewTripleDESCipher(newKey)
	if err != nil {
		return nil, err
	}

	if len(iv) != block.BlockSize() {
		return nil, errors.New("iv length invalid")
	}

	text := make([]byte, len(cipherText))
	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypter.CryptBlocks(text, cipherText)

	clearText, err := UnPadding(text, padding)
	if err != nil {
		return nil, err
	}

	return clearText, nil
}

//PKCS5补位
func PKCS5Padding(text []byte, blockSize int) []byte {
	padding := blockSize - len(text) % blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(text, padtext...)
}

//去除PKCS5补位
func PKCS5Unpadding(text []byte) []byte{
	length := len(text)
	padtext := int(text[length - 1])
	return text[:(length - padtext)]
}

//补位方法
func Padding(text []byte, padding int)([]byte, error) {
	switch padding {
	case NOPADDING:
		if len(text) % 8 != 0 {
			return nil, errors.New("text length invalid")
		}
	case PKCS5PADDING:
		return PKCS5Padding(text, 8), nil
	default:
		return nil, errors.New("padding type error")
	}

	return text, nil
}

//去除补位方法
func UnPadding(text []byte, padding int)([]byte, error) {
	switch padding {
	case NOPADDING:
		if len(text) % 8 != 0 {
			return nil, errors.New("text length invalid")
		}
	case PKCS5PADDING:
		return PKCS5Unpadding(text), nil
	default:
		return nil, errors.New("padding type error.")
	}
	return text, nil
}