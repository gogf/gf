// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gbase64 provides useful API for BASE64 encoding/decoding algorithm.
package gbase64

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"

	"github.com/gogf/gf/v2/errors/gerror"
)

const (
	EQ = byte(61)
)

// Encode encodes bytes with BASE64 algorithm.
func Encode(src []byte) []byte {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst
}

// UrlEncode encodes bytes with URL-secure Base64 algorithm.
// The noEqualSign parameter indicates that there is no equal sign after generation. default:true
func UrlEncode(src []byte, noEqualSign ...bool) []byte {
	dst := make([]byte, base64.URLEncoding.EncodedLen(len(src)))
	base64.URLEncoding.Encode(dst, src)
	if (len(noEqualSign) > 1 && noEqualSign[0]) || len(noEqualSign) == 0 {
		i := bytes.IndexByte(dst, EQ)
		if i != -1 {
			return dst[:i]
		}
	}
	return dst
}

// EncodeString encodes string with BASE64 algorithm.
func EncodeString(src string) string {
	return EncodeToString([]byte(src))
}

// UrlEncodeString encodes string with URL-secure Base64 algorithm.
// The noEqualSign parameter indicates that there is no equal sign after generation. default:true
func UrlEncodeString(src string, noEqualSign ...bool) string {
	return UrlEncodeToString([]byte(src), noEqualSign...)
}

// EncodeToString encodes bytes to string with BASE64 algorithm.
func EncodeToString(src []byte) string {
	return string(Encode(src))
}

// UrlEncodeToString encodes bytes to string with URL-secure BASE64 algorithm.
// The noEqualSign parameter indicates that there is no equal sign after generation. default:true
func UrlEncodeToString(src []byte, noEqualSign ...bool) string {
	return string(UrlEncode(src, noEqualSign...))
}

// EncodeFile encodes file content of `path` using BASE64 algorithms.
func EncodeFile(path string) ([]byte, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		err = gerror.Wrapf(err, `ioutil.ReadFile failed for filename "%s"`, path)
		return nil, err
	}
	return Encode(content), nil
}

// MustEncodeFile encodes file content of `path` using BASE64 algorithms.
// It panics if any error occurs.
func MustEncodeFile(path string) []byte {
	result, err := EncodeFile(path)
	if err != nil {
		panic(err)
	}
	return result
}

// EncodeFileToString encodes file content of `path` to string using BASE64 algorithms.
func EncodeFileToString(path string) (string, error) {
	content, err := EncodeFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// MustEncodeFileToString encodes file content of `path` to string using BASE64 algorithms.
// It panics if any error occurs.
func MustEncodeFileToString(path string) string {
	result, err := EncodeFileToString(path)
	if err != nil {
		panic(err)
	}
	return result
}

// Decode decodes bytes with BASE64 algorithm.
func Decode(data []byte) ([]byte, error) {
	var (
		src    = make([]byte, base64.StdEncoding.DecodedLen(len(data)))
		n, err = base64.StdEncoding.Decode(src, data)
	)
	if err != nil {
		err = gerror.Wrap(err, `base64.StdEncoding.Decode failed`)
	}
	return src[:n], err
}

// UrlDecode decodes bytes with URL-secure BASE64 algorithm.
// The noEqualSign parameter indicates that the string passed in has no equal sign. default:true
func UrlDecode(data []byte, noEqualSign ...bool) ([]byte, error) {
	var src []byte
	if (len(noEqualSign) > 1 && noEqualSign[0]) || len(noEqualSign) == 0 {
		length := len(data)
		if length%4 != 0 {
			length = length + 4 - (length % 4)
		}
		src = make([]byte, length)
		n := copy(src, data)
		for i := n; i < length; i++ {
			src[i] = EQ
		}
	} else {
		src = data
	}
	dstLen := base64.URLEncoding.DecodedLen(len(src))
	dst := make([]byte, dstLen)
	dstRealLen, err := base64.URLEncoding.Decode(dst, src)
	if err != nil {
		err = gerror.Wrap(err, `base64.URLEncoding.Decode failed`)
	}
	return dst[:dstRealLen], err
}

// MustDecode decodes bytes with BASE64 algorithm.
// It panics if any error occurs.
func MustDecode(data []byte) []byte {
	result, err := Decode(data)
	if err != nil {
		panic(err)
	}
	return result
}

// MustUrlDecode decodes bytes with URL-secure BASE64 algorithm.
// The noEqualSign parameter indicates that the string passed in has no equal sign. default:true
// It panics if any error occurs.
func MustUrlDecode(data []byte, noEqualSign ...bool) []byte {
	result, err := UrlDecode(data, noEqualSign...)
	if err != nil {
		panic(err)
	}
	return result
}

// UrlDecodeString decodes string with URL-secure BASE64 algorithm.
// The noEqualSign parameter indicates that the string passed in has no equal sign. default:true
func UrlDecodeString(data string, noEqualSign ...bool) ([]byte, error) {
	return UrlDecode([]byte(data), noEqualSign...)
}

// DecodeString decodes string with BASE64 algorithm.
func DecodeString(data string) ([]byte, error) {
	return Decode([]byte(data))
}

// MustUrlDecodeString decodes string with URL-secure BASE64 algorithm.
// The noEqualSign parameter indicates that the string passed in has no equal sign. default:true
// It panics if any error occurs.
func MustUrlDecodeString(data string, noEqualSign ...bool) []byte {
	result, err := UrlDecodeString(data, noEqualSign...)
	if err != nil {
		panic(err)
	}
	return result
}

// MustDecodeString decodes string with BASE64 algorithm.
// It panics if any error occurs.
func MustDecodeString(data string) []byte {
	result, err := DecodeString(data)
	if err != nil {
		panic(err)
	}
	return result
}

// UrlDecodeToString decodes string with URL-secure BASE64 algorithm.
// The noEqualSign parameter indicates that the string passed in has no equal sign. default:true
func UrlDecodeToString(data string, noEqualSign ...bool) (string, error) {
	b, err := UrlDecodeString(data, noEqualSign...)
	return string(b), err
}

// DecodeToString decodes string with BASE64 algorithm.
func DecodeToString(data string) (string, error) {
	b, err := DecodeString(data)
	return string(b), err
}

// MustDecodeToString decodes string with BASE64 algorithm.
// It panics if any error occurs.
func MustDecodeToString(data string) string {
	result, err := DecodeToString(data)
	if err != nil {
		panic(err)
	}
	return result
}

// MustUrlDecodeToString decodes string with URL-secure BASE64 algorithm.
// The noEqualSign parameter indicates that the string passed in has no equal sign. default:true
// It panics if any error occurs.
func MustUrlDecodeToString(data string, noEqualSign ...bool) string {
	result, err := UrlDecodeToString(data, noEqualSign...)
	if err != nil {
		panic(err)
	}
	return result
}
