// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package grsa_test

import (
	"testing"

	"github.com/gogf/gf/v2/crypto/grsa"
	"github.com/gogf/gf/v2/test/gtest"
)

var (
	plainText = []byte("Hello GoFrame!")

	// Test keys will be generated dynamically
	privateKey []byte
	publicKey  []byte
)

func init() {
	// Generate test key pair once during initialization
	var err error
	privateKey, publicKey, err = grsa.GenerateKeyPair(1024)
	if err != nil {
		panic(err)
	}
}

func TestEncryptDecrypt(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Encrypt with public key
		cipherText, err := grsa.Encrypt(plainText, publicKey)
		t.AssertNil(err)
		t.AssertNE(cipherText, nil)

		// Decrypt with private key
		decrypted, err := grsa.Decrypt(cipherText, privateKey)
		t.AssertNil(err)
		t.Assert(decrypted, plainText)
	})
}

func TestEncryptErr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Invalid public key
		_, err := grsa.Encrypt(plainText, []byte("invalid key"))
		t.AssertNE(err, nil)

		// Empty public key
		_, err = grsa.Encrypt(plainText, []byte(""))
		t.AssertNE(err, nil)
	})
}

func TestDecryptErr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Invalid private key
		_, err := grsa.Decrypt(plainText, []byte("invalid key"))
		t.AssertNE(err, nil)

		// Empty private key
		_, err = grsa.Decrypt(plainText, []byte(""))
		t.AssertNE(err, nil)

		// Invalid cipher text
		_, err = grsa.Decrypt([]byte("invalid cipher"), privateKey)
		t.AssertNE(err, nil)
	})
}

func TestGenerateKeyPair(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Generate 1024 bit key pair
		privKey, pubKey, err := grsa.GenerateKeyPair(1024)
		t.AssertNil(err)
		t.AssertNE(privKey, nil)
		t.AssertNE(pubKey, nil)

		// Test encryption and decryption with generated keys
		testData := []byte("Test data with generated keys")
		cipherText, err := grsa.Encrypt(testData, pubKey)
		t.AssertNil(err)

		decrypted, err := grsa.Decrypt(cipherText, privKey)
		t.AssertNil(err)
		t.Assert(decrypted, testData)
	})

	gtest.C(t, func(t *gtest.T) {
		// Generate 2048 bit key pair
		privKey, pubKey, err := grsa.GenerateKeyPair(2048)
		t.AssertNil(err)
		t.AssertNE(privKey, nil)
		t.AssertNE(pubKey, nil)

		// Test encryption and decryption
		testData := []byte("Test with 2048 bit key")
		cipherText, err := grsa.Encrypt(testData, pubKey)
		t.AssertNil(err)

		decrypted, err := grsa.Decrypt(cipherText, privKey)
		t.AssertNil(err)
		t.Assert(decrypted, testData)
	})
}

func TestGenerateKeyPairErr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Invalid bits (too small)
		_, _, err := grsa.GenerateKeyPair(256)
		t.AssertNE(err, nil)

		// Invalid bits (zero)
		_, _, err = grsa.GenerateKeyPair(0)
		t.AssertNE(err, nil)

		// Invalid bits (negative)
		_, _, err = grsa.GenerateKeyPair(-1)
		t.AssertNE(err, nil)
	})
}
