// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grsa_test

import (
	"encoding/base64"
	"testing"

	"github.com/gogf/gf/v2/crypto/grsa"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestEncryptDecrypt(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Generate a key pair for testing
		privateKey, publicKey, err := grsa.GenerateDefaultKeyPair()
		t.AssertNil(err)
		t.AssertNE(privateKey, nil)
		t.AssertNE(publicKey, nil)

		// Test data to encrypt
		plainText := []byte("Hello, World!")

		// Encrypt with public key
		cipherText, err := grsa.Encrypt(plainText, publicKey)
		t.AssertNil(err)
		t.AssertNE(cipherText, nil)

		// Decrypt with private key
		decryptedText, err := grsa.Decrypt(cipherText, privateKey)
		t.AssertNil(err)
		t.Assert(string(decryptedText), string(plainText))
	})
}

func TestEncryptDecryptBase64(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Generate a key pair for testing
		privateKey, publicKey, err := grsa.GenerateDefaultKeyPair()
		t.AssertNil(err)
		t.AssertNE(privateKey, nil)
		t.AssertNE(publicKey, nil)

		// Encode keys to base64
		privateKeyBase64 := encodeToBase64(privateKey)
		publicKeyBase64 := encodeToBase64(publicKey)

		// Test data to encrypt
		plainText := []byte("Hello, Base64 World!")

		// Encrypt with public key
		cipherTextBase64, err := grsa.EncryptBase64(plainText, publicKeyBase64)
		t.AssertNil(err)
		t.AssertNE(cipherTextBase64, "")

		// Decrypt with private key
		decryptedText, err := grsa.DecryptBase64(cipherTextBase64, privateKeyBase64)
		t.AssertNil(err)
		t.Assert(string(decryptedText), string(plainText))
	})
}

func TestGenerateKeyPair(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test generating a 2048-bit RSA key pair
		privateKey, publicKey, err := grsa.GenerateKeyPair(2048)
		t.AssertNil(err)
		t.AssertNE(privateKey, nil)
		t.AssertNE(publicKey, nil)

		// Check if keys are in correct format
		privateKeyType, err := grsa.GetPrivateKeyType(privateKey)
		t.AssertNil(err)
		t.Assert(privateKeyType, "PKCS#1")

		// Test with 1024-bit key for faster test execution only.
		// Note: 1024-bit keys are NOT secure for production use.
		// Always use at least 2048-bit keys in production.
		privateKey, publicKey, err = grsa.GenerateKeyPair(1024)
		t.AssertNil(err)
		t.AssertNE(privateKey, nil)
		t.AssertNE(publicKey, nil)
	})
}

func TestGenerateKeyPairPKCS8(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test generating a 2048-bit RSA key pair in PKCS#8 format
		privateKey, publicKey, err := grsa.GenerateKeyPairPKCS8(2048)
		t.AssertNil(err)
		t.AssertNE(privateKey, nil)
		t.AssertNE(publicKey, nil)

		// Check if keys are in correct format
		privateKeyType, err := grsa.GetPrivateKeyType(privateKey)
		t.AssertNil(err)
		t.Assert(privateKeyType, "PKCS#8")

		// Test with 1024-bit key for faster test execution only.
		// Note: 1024-bit keys are NOT secure for production use.
		privateKey, publicKey, err = grsa.GenerateKeyPairPKCS8(1024)
		t.AssertNil(err)
		t.AssertNE(privateKey, nil)
		t.AssertNE(publicKey, nil)
	})
}

func TestEncryptAndDecryptPKCS(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Generate both types of key pairs for testing
		privateKey1, publicKey1, err := grsa.GenerateKeyPair(2048)
		t.AssertNil(err)

		privateKey8, publicKey8, err := grsa.GenerateKeyPairPKCS8(2048)
		t.AssertNil(err)

		// Test data to encrypt
		plainText := []byte("Hello, Mixed Formats!")

		// Test general encrypt/decrypt with PKCS#1 keys
		cipherText, err := grsa.Encrypt(plainText, publicKey1)
		t.AssertNil(err)
		t.AssertNE(cipherText, nil)

		decryptedText, err := grsa.Decrypt(cipherText, privateKey1)
		t.AssertNil(err)
		t.Assert(string(decryptedText), string(plainText))

		// Test general encrypt/decrypt with PKCS#8 keys
		cipherText8, err := grsa.Encrypt(plainText, publicKey8)
		t.AssertNil(err)
		t.AssertNE(cipherText8, nil)

		decryptedText8, err := grsa.Decrypt(cipherText8, privateKey8)
		t.AssertNil(err)
		t.Assert(string(decryptedText8), string(plainText))
	})
}

func TestGetPrivateKeyType(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Generate both PKCS#1 and PKCS#8 key pairs
		// Note: 1024-bit keys used here for faster test execution only.
		// NOT secure for production use.
		privKey1, _, err := grsa.GenerateKeyPair(1024)
		t.AssertNil(err)

		privKey8, _, err := grsa.GenerateKeyPairPKCS8(1024)
		t.AssertNil(err)

		// Check types
		keyType1, err := grsa.GetPrivateKeyType(privKey1)
		t.AssertNil(err)
		t.Assert(keyType1, "PKCS#1")

		keyType8, err := grsa.GetPrivateKeyType(privKey8)
		t.AssertNil(err)
		t.Assert(keyType8, "PKCS#8")
	})
}

func TestEncryptPKCS1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Generate a key pair for testing (PKCS#1 format)
		privateKey, publicKey, err := grsa.GenerateKeyPair(2048)
		t.AssertNil(err)
		t.AssertNE(privateKey, nil)
		t.AssertNE(publicKey, nil)

		// Test data to encrypt
		plainText := []byte("Hello, PKCS#1 World!")

		// Encrypt with public key using PKCS#1 format specifically
		cipherText, err := grsa.EncryptPKCS1(plainText, publicKey)
		t.AssertNil(err)
		t.AssertNE(cipherText, nil)

		// Decrypt with private key using PKCS#1 format specifically
		decryptedText, err := grsa.DecryptPKCS1(cipherText, privateKey)
		t.AssertNil(err)
		t.Assert(string(decryptedText), string(plainText))
	})
}

func TestEncryptPKCS1Base64(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Generate a key pair for testing
		privateKey, publicKey, err := grsa.GenerateKeyPair(2048)
		t.AssertNil(err)
		t.AssertNE(privateKey, nil)
		t.AssertNE(publicKey, nil)

		// Encode keys to base64
		privateKeyBase64 := encodeToBase64(privateKey)
		publicKeyBase64 := encodeToBase64(publicKey)

		// Test data to encrypt
		plainText := []byte("Hello, PKCS#1 Base64 World!")

		// Encrypt with public key using PKCS#1 format specifically
		cipherTextBase64, err := grsa.EncryptPKCS1Base64(plainText, publicKeyBase64)
		t.AssertNil(err)
		t.AssertNE(cipherTextBase64, "")

		// Decrypt with private key using PKCS#1 format specifically
		decryptedText, err := grsa.DecryptPKCS1Base64(cipherTextBase64, privateKeyBase64)
		t.AssertNil(err)
		t.Assert(string(decryptedText), string(plainText))
	})
}

// Helper function to encode to base64
func encodeToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func TestEncryptWithInvalidPublicKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		plainText := []byte("Hello, World!")

		// Test with invalid public key
		_, err := grsa.Encrypt(plainText, []byte("invalid key"))
		t.AssertNE(err, nil)

		// Test with empty public key
		_, err = grsa.Encrypt(plainText, []byte{})
		t.AssertNE(err, nil)

		// Test with nil public key
		_, err = grsa.Encrypt(plainText, nil)
		t.AssertNE(err, nil)
	})
}

func TestDecryptWithInvalidPrivateKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Generate a valid key pair and encrypt some data
		privateKey, publicKey, err := grsa.GenerateDefaultKeyPair()
		t.AssertNil(err)

		plainText := []byte("Hello, World!")
		cipherText, err := grsa.Encrypt(plainText, publicKey)
		t.AssertNil(err)

		// Test decryption with invalid private key
		_, err = grsa.Decrypt(cipherText, []byte("invalid key"))
		t.AssertNE(err, nil)

		// Test decryption with empty private key
		_, err = grsa.Decrypt(cipherText, []byte{})
		t.AssertNE(err, nil)

		// Test decryption with wrong private key
		wrongPrivKey, _, err := grsa.GenerateDefaultKeyPair()
		t.AssertNil(err)
		_, err = grsa.Decrypt(cipherText, wrongPrivKey)
		t.AssertNE(err, nil)

		// Verify correct decryption still works
		decrypted, err := grsa.Decrypt(cipherText, privateKey)
		t.AssertNil(err)
		t.Assert(string(decrypted), string(plainText))
	})
}

func TestEncryptWithOversizedPlaintext(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Generate a 2048-bit key pair
		_, publicKey, err := grsa.GenerateDefaultKeyPair()
		t.AssertNil(err)

		// For 2048-bit key with PKCS#1 v1.5 padding, max size is 256 - 11 = 245 bytes
		// Create plaintext that exceeds this limit
		oversizedPlainText := make([]byte, 300)
		for i := range oversizedPlainText {
			oversizedPlainText[i] = 'A'
		}

		// Encryption should fail with oversized plaintext
		_, err = grsa.Encrypt(oversizedPlainText, publicKey)
		t.AssertNE(err, nil)

		// Verify that valid size plaintext works
		validPlainText := make([]byte, 200)
		for i := range validPlainText {
			validPlainText[i] = 'B'
		}
		_, err = grsa.Encrypt(validPlainText, publicKey)
		t.AssertNil(err)
	})
}

func TestDecryptWithCorruptedCiphertext(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		privateKey, publicKey, err := grsa.GenerateDefaultKeyPair()
		t.AssertNil(err)

		plainText := []byte("Hello, World!")
		cipherText, err := grsa.Encrypt(plainText, publicKey)
		t.AssertNil(err)

		// Corrupt the ciphertext
		corruptedCipherText := make([]byte, len(cipherText))
		copy(corruptedCipherText, cipherText)
		corruptedCipherText[0] ^= 0xFF
		corruptedCipherText[len(corruptedCipherText)-1] ^= 0xFF

		// Decryption should fail with corrupted ciphertext
		_, err = grsa.Decrypt(corruptedCipherText, privateKey)
		t.AssertNE(err, nil)
	})
}

func TestGetPrivateKeyTypeWithInvalidKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test with invalid key
		_, err := grsa.GetPrivateKeyType([]byte("invalid key"))
		t.AssertNE(err, nil)

		// Test with empty key
		_, err = grsa.GetPrivateKeyType([]byte{})
		t.AssertNE(err, nil)

		// Test with nil key
		_, err = grsa.GetPrivateKeyType(nil)
		t.AssertNE(err, nil)
	})
}

func TestBase64FunctionsWithInvalidInput(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		plainText := []byte("Hello, World!")

		// Test EncryptBase64 with invalid base64 public key
		_, err := grsa.EncryptBase64(plainText, "not-valid-base64!!!")
		t.AssertNE(err, nil)

		// Test DecryptBase64 with invalid base64 private key
		_, err = grsa.DecryptBase64("validbase64==", "not-valid-base64!!!")
		t.AssertNE(err, nil)

		// Test DecryptBase64 with invalid base64 ciphertext
		privateKey, _, err := grsa.GenerateDefaultKeyPair()
		t.AssertNil(err)
		privateKeyBase64 := encodeToBase64(privateKey)
		_, err = grsa.DecryptBase64("not-valid-base64!!!", privateKeyBase64)
		t.AssertNE(err, nil)
	})
}

func TestDecryptPKCS8WithPKCS1Key(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Generate PKCS#1 key pair
		privateKey1, publicKey1, err := grsa.GenerateKeyPair(2048)
		t.AssertNil(err)

		plainText := []byte("Hello, World!")
		cipherText, err := grsa.EncryptPKCS1(plainText, publicKey1)
		t.AssertNil(err)

		// DecryptPKCS8 should fail with PKCS#1 private key (no fallback)
		_, err = grsa.DecryptPKCS8(cipherText, privateKey1)
		t.AssertNE(err, nil)

		// DecryptPKCS1 should work
		decrypted, err := grsa.DecryptPKCS1(cipherText, privateKey1)
		t.AssertNil(err)
		t.Assert(string(decrypted), string(plainText))
	})
}
