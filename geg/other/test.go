package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
    "gitee.com/johng/gf/g/encoding/gbase64"
)

const (
	ivDefValue = "I Love Go Frame!"
)

func main() {
	v := "1234567890123456789012345678901234567890123456789012345678901234567890"
	k := "123456789012345 "
	e, err := AesEncrypt([]byte(v), []byte(k))
	fmt.Println(err)
	fmt.Println(len(e))
	fmt.Println(string(gbase64.Encode(string(e))))
	d, err := AesDecrypt([]byte(e), []byte(k))
	fmt.Println(err)
	fmt.Println(string(d))
}

func AesEncrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	plaintext = PKCS5Padding(plaintext, blockSize)
	iv := []byte(ivDefValue)
	blockMode := cipher.NewCBCEncrypter(block, iv)

	ciphertext := make([]byte, len(plaintext))
	blockMode.CryptBlocks(ciphertext, plaintext)

	return ciphertext, nil
}

func AesDecrypt(cipherText []byte, key []byte) ([]byte, error) {
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