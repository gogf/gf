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

func TestEncryptPKIX(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Generate PKCS#8 key pair (which uses PKIX public key format)
		privateKey8, publicKey8, err := grsa.GenerateKeyPairPKCS8(2048)
		t.AssertNil(err)

		plainText := []byte("Hello, PKIX World!")

		// Encrypt with PKIX public key
		cipherText, err := grsa.EncryptPKIX(plainText, publicKey8)
		t.AssertNil(err)
		t.AssertNE(cipherText, nil)

		// Decrypt with PKCS#8 private key
		decrypted, err := grsa.DecryptPKCS8(cipherText, privateKey8)
		t.AssertNil(err)
		t.Assert(string(decrypted), string(plainText))

		// Test with invalid public key
		_, err = grsa.EncryptPKIX(plainText, []byte("invalid key"))
		t.AssertNE(err, nil)

		// Test with PKCS#1 public key (should fail for EncryptPKIX)
		_, publicKey1, err := grsa.GenerateKeyPair(2048)
		t.AssertNil(err)
		_, err = grsa.EncryptPKIX(plainText, publicKey1)
		t.AssertNE(err, nil)

		// Test oversized plaintext
		oversizedPlainText := make([]byte, 300)
		_, err = grsa.EncryptPKIX(oversizedPlainText, publicKey8)
		t.AssertNE(err, nil)
	})
}

func TestEncryptPKCS8Alias(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Generate PKCS#8 key pair
		privateKey8, publicKey8, err := grsa.GenerateKeyPairPKCS8(2048)
		t.AssertNil(err)

		plainText := []byte("Hello, PKCS8 Alias!")

		// EncryptPKCS8 is an alias for EncryptPKIX
		cipherText, err := grsa.EncryptPKCS8(plainText, publicKey8)
		t.AssertNil(err)
		t.AssertNE(cipherText, nil)

		// Decrypt should work
		decrypted, err := grsa.DecryptPKCS8(cipherText, privateKey8)
		t.AssertNil(err)
		t.Assert(string(decrypted), string(plainText))
	})
}

func TestEncryptPKIXBase64(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Generate PKCS#8 key pair
		privateKey8, publicKey8, err := grsa.GenerateKeyPairPKCS8(2048)
		t.AssertNil(err)

		privateKeyBase64 := encodeToBase64(privateKey8)
		publicKeyBase64 := encodeToBase64(publicKey8)

		plainText := []byte("Hello, PKIX Base64!")

		// Encrypt with PKIX public key
		cipherTextBase64, err := grsa.EncryptPKIXBase64(plainText, publicKeyBase64)
		t.AssertNil(err)
		t.AssertNE(cipherTextBase64, "")

		// Decrypt with PKCS#8 private key
		decrypted, err := grsa.DecryptPKCS8Base64(cipherTextBase64, privateKeyBase64)
		t.AssertNil(err)
		t.Assert(string(decrypted), string(plainText))

		// Test with invalid base64 public key
		_, err = grsa.EncryptPKIXBase64(plainText, "not-valid-base64!!!")
		t.AssertNE(err, nil)

		// Test with invalid public key content
		_, err = grsa.EncryptPKIXBase64(plainText, encodeToBase64([]byte("invalid key")))
		t.AssertNE(err, nil)
	})
}

func TestEncryptPKCS8Base64Alias(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Generate PKCS#8 key pair
		privateKey8, publicKey8, err := grsa.GenerateKeyPairPKCS8(2048)
		t.AssertNil(err)

		privateKeyBase64 := encodeToBase64(privateKey8)
		publicKeyBase64 := encodeToBase64(publicKey8)

		plainText := []byte("Hello, PKCS8 Base64 Alias!")

		// EncryptPKCS8Base64 is an alias for EncryptPKIXBase64
		cipherTextBase64, err := grsa.EncryptPKCS8Base64(plainText, publicKeyBase64)
		t.AssertNil(err)
		t.AssertNE(cipherTextBase64, "")

		// Decrypt should work
		decrypted, err := grsa.DecryptPKCS8Base64(cipherTextBase64, privateKeyBase64)
		t.AssertNil(err)
		t.Assert(string(decrypted), string(plainText))
	})
}

func TestDecryptPKCS8Base64(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Generate PKCS#8 key pair
		privateKey8, publicKey8, err := grsa.GenerateKeyPairPKCS8(2048)
		t.AssertNil(err)

		privateKeyBase64 := encodeToBase64(privateKey8)
		publicKeyBase64 := encodeToBase64(publicKey8)

		plainText := []byte("Hello, DecryptPKCS8Base64!")

		// Encrypt
		cipherTextBase64, err := grsa.EncryptPKIXBase64(plainText, publicKeyBase64)
		t.AssertNil(err)

		// Decrypt
		decrypted, err := grsa.DecryptPKCS8Base64(cipherTextBase64, privateKeyBase64)
		t.AssertNil(err)
		t.Assert(string(decrypted), string(plainText))

		// Test with invalid base64 private key
		_, err = grsa.DecryptPKCS8Base64(cipherTextBase64, "not-valid-base64!!!")
		t.AssertNE(err, nil)

		// Test with invalid base64 ciphertext
		_, err = grsa.DecryptPKCS8Base64("not-valid-base64!!!", privateKeyBase64)
		t.AssertNE(err, nil)
	})
}

func TestGetPrivateKeyTypeBase64(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Generate both PKCS#1 and PKCS#8 key pairs
		privKey1, _, err := grsa.GenerateKeyPair(2048)
		t.AssertNil(err)

		privKey8, _, err := grsa.GenerateKeyPairPKCS8(2048)
		t.AssertNil(err)

		// Check types via base64
		keyType1, err := grsa.GetPrivateKeyTypeBase64(encodeToBase64(privKey1))
		t.AssertNil(err)
		t.Assert(keyType1, "PKCS#1")

		keyType8, err := grsa.GetPrivateKeyTypeBase64(encodeToBase64(privKey8))
		t.AssertNil(err)
		t.Assert(keyType8, "PKCS#8")

		// Test with invalid base64
		_, err = grsa.GetPrivateKeyTypeBase64("not-valid-base64!!!")
		t.AssertNE(err, nil)
	})
}

func TestExtractPKCS1PublicKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Generate PKCS#1 key pair
		privateKey, publicKey, err := grsa.GenerateKeyPair(2048)
		t.AssertNil(err)

		// Extract public key from private key
		extractedPublicKey, err := grsa.ExtractPKCS1PublicKey(privateKey)
		t.AssertNil(err)
		t.AssertNE(extractedPublicKey, nil)

		// The extracted public key should work for encryption
		plainText := []byte("Hello, Extracted Key!")
		cipherText, err := grsa.EncryptPKCS1(plainText, extractedPublicKey)
		t.AssertNil(err)

		// Decrypt with original private key
		decrypted, err := grsa.DecryptPKCS1(cipherText, privateKey)
		t.AssertNil(err)
		t.Assert(string(decrypted), string(plainText))

		// Compare extracted key with original (they should be equivalent)
		cipherText2, err := grsa.EncryptPKCS1(plainText, publicKey)
		t.AssertNil(err)
		decrypted2, err := grsa.DecryptPKCS1(cipherText2, privateKey)
		t.AssertNil(err)
		t.Assert(string(decrypted2), string(plainText))

		// Test with invalid private key
		_, err = grsa.ExtractPKCS1PublicKey([]byte("invalid key"))
		t.AssertNE(err, nil)

		// Test with PKCS#8 private key (should fail)
		privateKey8, _, err := grsa.GenerateKeyPairPKCS8(2048)
		t.AssertNil(err)
		_, err = grsa.ExtractPKCS1PublicKey(privateKey8)
		t.AssertNE(err, nil)
	})
}

func TestDecryptPKCS1WithInvalidKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		privateKey, publicKey, err := grsa.GenerateKeyPair(2048)
		t.AssertNil(err)

		plainText := []byte("Hello, World!")
		cipherText, err := grsa.EncryptPKCS1(plainText, publicKey)
		t.AssertNil(err)

		// Test with invalid private key
		_, err = grsa.DecryptPKCS1(cipherText, []byte("invalid key"))
		t.AssertNE(err, nil)

		// Test with PKCS#8 private key (should fail for DecryptPKCS1)
		privateKey8, _, err := grsa.GenerateKeyPairPKCS8(2048)
		t.AssertNil(err)
		_, err = grsa.DecryptPKCS1(cipherText, privateKey8)
		t.AssertNE(err, nil)

		// Verify correct decryption works
		decrypted, err := grsa.DecryptPKCS1(cipherText, privateKey)
		t.AssertNil(err)
		t.Assert(string(decrypted), string(plainText))
	})
}

func TestDecryptPKCS8WithInvalidKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		privateKey8, publicKey8, err := grsa.GenerateKeyPairPKCS8(2048)
		t.AssertNil(err)

		plainText := []byte("Hello, World!")
		cipherText, err := grsa.EncryptPKIX(plainText, publicKey8)
		t.AssertNil(err)

		// Test with invalid private key
		_, err = grsa.DecryptPKCS8(cipherText, []byte("invalid key"))
		t.AssertNE(err, nil)

		// Verify correct decryption works
		decrypted, err := grsa.DecryptPKCS8(cipherText, privateKey8)
		t.AssertNil(err)
		t.Assert(string(decrypted), string(plainText))
	})
}

func TestEncryptPKCS1WithInvalidKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		plainText := []byte("Hello, World!")

		// Test with invalid public key
		_, err := grsa.EncryptPKCS1(plainText, []byte("invalid key"))
		t.AssertNE(err, nil)

		// Test with PKCS#8 public key (should fail for EncryptPKCS1)
		_, publicKey8, err := grsa.GenerateKeyPairPKCS8(2048)
		t.AssertNil(err)
		_, err = grsa.EncryptPKCS1(plainText, publicKey8)
		t.AssertNE(err, nil)
	})
}

func TestEncryptPKCS1WithOversizedPlaintext(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		_, publicKey, err := grsa.GenerateKeyPair(2048)
		t.AssertNil(err)

		// Create oversized plaintext
		oversizedPlainText := make([]byte, 300)
		_, err = grsa.EncryptPKCS1(oversizedPlainText, publicKey)
		t.AssertNE(err, nil)
	})
}

func TestEncryptPKCS1Base64WithInvalidInput(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		plainText := []byte("Hello, World!")

		// Test with invalid base64 public key
		_, err := grsa.EncryptPKCS1Base64(plainText, "not-valid-base64!!!")
		t.AssertNE(err, nil)

		// Test with invalid public key content
		_, err = grsa.EncryptPKCS1Base64(plainText, encodeToBase64([]byte("invalid key")))
		t.AssertNE(err, nil)
	})
}

func TestDecryptPKCS1Base64WithInvalidInput(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		privateKey, publicKey, err := grsa.GenerateKeyPair(2048)
		t.AssertNil(err)

		privateKeyBase64 := encodeToBase64(privateKey)
		publicKeyBase64 := encodeToBase64(publicKey)

		plainText := []byte("Hello, World!")
		cipherTextBase64, err := grsa.EncryptPKCS1Base64(plainText, publicKeyBase64)
		t.AssertNil(err)

		// Test with invalid base64 private key
		_, err = grsa.DecryptPKCS1Base64(cipherTextBase64, "not-valid-base64!!!")
		t.AssertNE(err, nil)

		// Test with invalid base64 ciphertext
		_, err = grsa.DecryptPKCS1Base64("not-valid-base64!!!", privateKeyBase64)
		t.AssertNE(err, nil)
	})
}

func TestEncryptWithNonRSAPublicKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create a PEM block that is valid but not an RSA key
		// This tests the "not an RSA public key" error path
		// We use a valid PEM structure but with invalid content
		invalidPEM := []byte(`-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE
-----END PUBLIC KEY-----`)

		plainText := []byte("Hello, World!")
		_, err := grsa.Encrypt(plainText, invalidPEM)
		t.AssertNE(err, nil)
	})
}

func TestDecryptWithNonRSAPrivateKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create a PEM block that is valid but not an RSA key
		invalidPEM := []byte(`-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQg
-----END PRIVATE KEY-----`)

		cipherText := []byte("some cipher text")
		_, err := grsa.Decrypt(cipherText, invalidPEM)
		t.AssertNE(err, nil)
	})
}

func TestEncryptBase64WithInvalidPublicKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		plainText := []byte("Hello, World!")

		// Test with valid base64 but invalid key content
		invalidKeyBase64 := encodeToBase64([]byte("invalid key"))
		_, err := grsa.EncryptBase64(plainText, invalidKeyBase64)
		t.AssertNE(err, nil)
	})
}

func TestGetPrivateKeyTypeWithNonRSAKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create a valid PKCS#8 PEM but with non-RSA content (EC key)
		// This tests the "not an RSA private key" error path in GetPrivateKeyType
		ecPrivateKeyPEM := []byte(`-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgVcB/UNPczVP6jE4Z
p7v6qYQXsKQZLJGBJKKnUWuHb6+hRANCAASYn3k2T4VqPt1HVAK5Rc7rMb6lGOzF
v0MVLfCgPKANNGdBvGPmaSLFIxGMNL0v1C2RRvqqEu/vL3POoaqfMJhw
-----END PRIVATE KEY-----`)

		_, err := grsa.GetPrivateKeyType(ecPrivateKeyPEM)
		t.AssertNE(err, nil)
	})
}

func TestDecryptPKCS8WithNonRSAKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create a valid PKCS#8 PEM but with non-RSA content (EC key)
		ecPrivateKeyPEM := []byte(`-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgVcB/UNPczVP6jE4Z
p7v6qYQXsKQZLJGBJKKnUWuHb6+hRANCAASYn3k2T4VqPt1HVAK5Rc7rMb6lGOzF
v0MVLfCgPKANNGdBvGPmaSLFIxGMNL0v1C2RRvqqEu/vL3POoaqfMJhw
-----END PRIVATE KEY-----`)

		cipherText := []byte("some cipher text")
		_, err := grsa.DecryptPKCS8(cipherText, ecPrivateKeyPEM)
		t.AssertNE(err, nil)
	})
}

func TestEncryptPKIXWithNonRSAKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create a valid PKIX PEM but with non-RSA content (EC key)
		ecPublicKeyPEM := []byte(`-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEmJ95Nk+Faj7dR1QCuUXO6zG+pRjs
xb9DFS3woDygDTRnQbxj5mkixSMRjDS9L9QtkUb6qhLv7y9zzqGqnzCYcA==
-----END PUBLIC KEY-----`)

		plainText := []byte("Hello, World!")
		_, err := grsa.EncryptPKIX(plainText, ecPublicKeyPEM)
		t.AssertNE(err, nil)
	})
}

func TestEncryptWithNonRSAPKIXKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create a valid PKIX PEM but with non-RSA content (EC key)
		// This tests the "not an RSA public key" error path in Encrypt
		ecPublicKeyPEM := []byte(`-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEmJ95Nk+Faj7dR1QCuUXO6zG+pRjs
xb9DFS3woDygDTRnQbxj5mkixSMRjDS9L9QtkUb6qhLv7y9zzqGqnzCYcA==
-----END PUBLIC KEY-----`)

		plainText := []byte("Hello, World!")
		_, err := grsa.Encrypt(plainText, ecPublicKeyPEM)
		t.AssertNE(err, nil)
	})
}

func TestDecryptWithNonRSAPKCS8Key(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create a valid PKCS#8 PEM but with non-RSA content (EC key)
		// This tests the "not an RSA private key" error path in Decrypt
		ecPrivateKeyPEM := []byte(`-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgVcB/UNPczVP6jE4Z
p7v6qYQXsKQZLJGBJKKnUWuHb6+hRANCAASYn3k2T4VqPt1HVAK5Rc7rMb6lGOzF
v0MVLfCgPKANNGdBvGPmaSLFIxGMNL0v1C2RRvqqEu/vL3POoaqfMJhw
-----END PRIVATE KEY-----`)

		cipherText := []byte("some cipher text")
		_, err := grsa.Decrypt(cipherText, ecPrivateKeyPEM)
		t.AssertNE(err, nil)
	})
}
