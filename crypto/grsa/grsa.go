// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package grsa provides useful API for RSA encryption/decryption, sign/verify algorithms.
package grsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// Encrypt encrypts `data` with public key.
func Encrypt(data, publicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "failed to decode public key")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to parse public key")
	}
	pub := pubInterface.(*rsa.PublicKey)
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, pub, data)
	if err != nil {
		return nil, gerror.Wrap(err, "failed to encrypt data")
	}
	return cipherText, nil
}

// Decrypt decrypts `data` with private key.
func Decrypt(data, privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "failed to decode private key")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS8 format
		privInterface, err2 := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err2 != nil {
			return nil, gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to parse private key")
		}
		priv = privInterface.(*rsa.PrivateKey)
	}
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, priv, data)
	if err != nil {
		return nil, gerror.Wrap(err, "failed to decrypt data")
	}
	return plainText, nil
}

// GenerateKeyPair generates RSA key pair with specified bits.
// The bits can be 1024, 2048, 4096, etc.
func GenerateKeyPair(bits int) (privateKey, publicKey []byte, err error) {
	if bits < 512 {
		return nil, nil, gerror.NewCodef(gcode.CodeInvalidParameter, "bits must be at least 512, got %d", bits)
	}
	priKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, gerror.Wrap(err, "failed to generate private key")
	}
	derStream := x509.MarshalPKCS1PrivateKey(priKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	privateKey = pem.EncodeToMemory(block)

	pubKey := &priKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, nil, gerror.Wrap(err, "failed to marshal public key")
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	publicKey = pem.EncodeToMemory(block)

	return privateKey, publicKey, nil
}
