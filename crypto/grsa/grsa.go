// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package grsa provides useful API for RSA encryption/decryption algorithms.
//
// This package includes functionality for:
// - Generating RSA key pairs in PKCS#1 and PKCS#8 formats
// - Encrypting and decrypting data with various key formats
// - Handling Base64 encoded keys
// - Detecting private key types
//
// # Security Considerations
//
// This package provides two padding schemes for RSA encryption:
//
// 1. PKCS#1 v1.5 (legacy): Used by Encrypt*, DecryptPKCS1*, DecryptPKCS8* functions.
// This padding scheme is considered less secure and vulnerable to padding oracle attacks.
// It is provided for backward compatibility with existing systems.
//
// 2. OAEP (recommended): Used by EncryptOAEP*, DecryptOAEP* functions.
// OAEP (Optimal Asymmetric Encryption Padding) is the recommended padding scheme
// for new applications as it provides better security guarantees.
//
// For new implementations, prefer using OAEP functions (EncryptOAEP, DecryptOAEP, etc.).
package grsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"hash"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

const (
	// DefaultRSAKeyBits is the default bit size for RSA key generation
	DefaultRSAKeyBits = 2048

	// KeyTypePKCS1 represents PKCS#1 format private key
	KeyTypePKCS1 = "PKCS#1"
	// KeyTypePKCS8 represents PKCS#8 format private key
	KeyTypePKCS8 = "PKCS#8"

	// PEM block types
	pemTypeRSAPrivateKey = "RSA PRIVATE KEY" // PKCS#1 private key
	pemTypePrivateKey    = "PRIVATE KEY"     // PKCS#8 private key
	pemTypeRSAPublicKey  = "RSA PUBLIC KEY"  // PKCS#1 public key
	pemTypePublicKey     = "PUBLIC KEY"      // PKIX public key
)

// Encrypt encrypts data with public key using PKCS#1 v1.5 padding (auto-detect format).
// The publicKey can be either PKCS#1 or PKCS#8 (PKIX) format.
//
// Note: RSA encryption has a size limit based on key size.
// For PKCS#1 v1.5 padding, max plaintext size = key_size_in_bytes - 11.
// For example, a 2048-bit key can encrypt at most 245 bytes.
//
// Security Warning: PKCS#1 v1.5 padding is vulnerable to padding oracle attacks.
// For new applications, consider using EncryptOAEP instead.
func Encrypt(plainText, publicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "invalid public key")
	}

	// Try PKCS#8 (PKIX) first
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		// Try PKCS#1
		pub, err = x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to parse public key")
		}
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "not an RSA public key")
	}

	// Validate plaintext size for PKCS#1 v1.5 padding
	maxSize := rsaPub.Size() - 11
	if len(plainText) > maxSize {
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter,
			"plaintext too long: max %d bytes for this key, got %d bytes", maxSize, len(plainText))
	}

	return rsa.EncryptPKCS1v15(rand.Reader, rsaPub, plainText)
}

// Decrypt decrypts data with private key using PKCS#1 v1.5 padding (auto-detect format).
// The privateKey can be either PKCS#1 or PKCS#8 format.
//
// Security Warning: PKCS#1 v1.5 padding is vulnerable to padding oracle attacks.
// For new applications, consider using DecryptOAEP instead.
func Decrypt(cipherText, privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "invalid private key")
	}

	// Try PKCS#8 first
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS#1
		priv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to parse private key")
		}
	}

	rsaPriv, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "not an RSA private key")
	}

	return rsa.DecryptPKCS1v15(rand.Reader, rsaPriv, cipherText)
}

// EncryptBase64 encrypts data with base64-encoded public key (auto-detect format)
// and returns base64-encoded result.
func EncryptBase64(plainText []byte, publicKeyBase64 string) (string, error) {
	publicKey, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return "", gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to decode public key")
	}

	encrypted, err := Encrypt(plainText, publicKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// DecryptBase64 decrypts base64-encoded data with base64-encoded private key (auto-detect format).
func DecryptBase64(cipherTextBase64, privateKeyBase64 string) ([]byte, error) {
	privateKey, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to decode private key")
	}

	cipherText, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to decode cipher text")
	}

	return Decrypt(cipherText, privateKey)
}

// EncryptPKIX encrypts data with public key in PKIX (X.509) format.
// PKIX is the standard format for public keys, often referred to as "PKCS#8 public key".
//
// Note: RSA encryption has a size limit based on key size.
// For PKCS#1 v1.5 padding, max plaintext size = key_size_in_bytes - 11.
func EncryptPKIX(plainText, publicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "invalid public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to parse PKIX public key")
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "not an RSA public key")
	}

	// Validate plaintext size for PKCS#1 v1.5 padding
	maxSize := rsaPub.Size() - 11
	if len(plainText) > maxSize {
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter,
			"plaintext too long: max %d bytes for this key, got %d bytes", maxSize, len(plainText))
	}

	return rsa.EncryptPKCS1v15(rand.Reader, rsaPub, plainText)
}

// EncryptPKCS8 is an alias for EncryptPKIX for backward compatibility.
// Deprecated: Use EncryptPKIX instead. Public keys use PKIX format, not PKCS#8.
func EncryptPKCS8(plainText, publicKey []byte) ([]byte, error) {
	return EncryptPKIX(plainText, publicKey)
}

// EncryptPKCS1 encrypts data with public key in PKCS#1 format.
//
// Note: RSA encryption has a size limit based on key size.
// For PKCS#1 v1.5 padding, max plaintext size = key_size_in_bytes - 11.
func EncryptPKCS1(plainText, publicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "invalid public key")
	}

	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to parse PKCS#1 public key")
	}

	// Validate plaintext size for PKCS#1 v1.5 padding
	maxSize := pub.Size() - 11
	if len(plainText) > maxSize {
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter,
			"plaintext too long: max %d bytes for this key, got %d bytes", maxSize, len(plainText))
	}

	return rsa.EncryptPKCS1v15(rand.Reader, pub, plainText)
}

// DecryptPKCS8 decrypts data with private key by PKCS#8 format.
func DecryptPKCS8(cipherText, privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "invalid private key")
	}

	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to parse PKCS#8 private key")
	}

	rsaPriv, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "not an RSA private key")
	}

	return rsa.DecryptPKCS1v15(rand.Reader, rsaPriv, cipherText)
}

// DecryptPKCS1 decrypts data with private key by PKCS#1 format.
func DecryptPKCS1(cipherText, privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "invalid private key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to parse private key")
	}

	return rsa.DecryptPKCS1v15(rand.Reader, priv, cipherText)
}

// EncryptPKIXBase64 encrypts data with PKIX public key and returns base64-encoded result.
func EncryptPKIXBase64(plainText []byte, publicKeyBase64 string) (string, error) {
	publicKey, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return "", gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to decode public key")
	}

	encrypted, err := EncryptPKIX(plainText, publicKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// EncryptPKCS8Base64 is an alias for EncryptPKIXBase64 for backward compatibility.
// Deprecated: Use EncryptPKIXBase64 instead.
func EncryptPKCS8Base64(plainText []byte, publicKeyBase64 string) (string, error) {
	return EncryptPKIXBase64(plainText, publicKeyBase64)
}

// EncryptPKCS1Base64 encrypts data with PKCS#1 public key and returns base64-encoded result.
func EncryptPKCS1Base64(plainText []byte, publicKeyBase64 string) (string, error) {
	publicKey, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return "", gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to decode public key")
	}

	encrypted, err := EncryptPKCS1(plainText, publicKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// DecryptPKCS8Base64 decrypts data with private key by PKCS#8 format and decode base64 input.
func DecryptPKCS8Base64(cipherTextBase64, privateKeyBase64 string) ([]byte, error) {
	privateKey, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to decode private key")
	}

	cipherText, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to decode cipher text")
	}

	return DecryptPKCS8(cipherText, privateKey)
}

// DecryptPKCS1Base64 decrypts base64-encoded data with PKCS#1 private key.
func DecryptPKCS1Base64(cipherTextBase64, privateKeyBase64 string) ([]byte, error) {
	privateKey, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to decode private key")
	}

	cipherText, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to decode cipher text")
	}

	return DecryptPKCS1(cipherText, privateKey)
}

// GetPrivateKeyType detects the type of private key (PKCS#1 or PKCS#8).
// It attempts to parse the key in both formats to determine the actual type.
func GetPrivateKeyType(privateKey []byte) (string, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return "", gerror.NewCode(gcode.CodeInvalidParameter, "invalid private key")
	}

	// Try PKCS#1 first
	_, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return KeyTypePKCS1, nil
	}

	// Try PKCS#8
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err == nil {
		if _, ok := priv.(*rsa.PrivateKey); ok {
			return KeyTypePKCS8, nil
		}
		return "", gerror.NewCode(gcode.CodeInvalidParameter, "not an RSA private key")
	}

	return "", gerror.NewCode(gcode.CodeInvalidParameter, "unknown private key format")
}

// GetPrivateKeyTypeBase64 detects the type of base64 encoded private key (PKCS#1 or PKCS#8).
func GetPrivateKeyTypeBase64(privateKeyBase64 string) (string, error) {
	privateKey, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return "", gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to decode private key")
	}

	return GetPrivateKeyType(privateKey)
}

// GenerateKeyPair generates a new RSA key pair with the given bits.
func GenerateKeyPair(bits int) (privateKey, publicKey []byte, err error) {
	// Generate private key
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, gerror.WrapCode(gcode.CodeInternalError, err, "failed to generate rsa key")
	}

	// Validate private key
	err = privKey.Validate()
	if err != nil {
		return nil, nil, gerror.WrapCode(gcode.CodeInternalError, err, "failed to validate rsa key")
	}

	// Marshal private key to PKCS#1 format
	privKeyBytes := x509.MarshalPKCS1PrivateKey(privKey)
	privateKey = pem.EncodeToMemory(&pem.Block{
		Type:  pemTypeRSAPrivateKey,
		Bytes: privKeyBytes,
	})

	// Generate PKCS#1 public key
	pubKeyBytes := x509.MarshalPKCS1PublicKey(&privKey.PublicKey)
	publicKey = pem.EncodeToMemory(&pem.Block{
		Type:  pemTypeRSAPublicKey,
		Bytes: pubKeyBytes,
	})

	return privateKey, publicKey, nil
}

// GenerateKeyPairPKCS8 generates a new RSA key pair with the given bits in PKCS#8 format.
func GenerateKeyPairPKCS8(bits int) (privateKey, publicKey []byte, err error) {
	// Generate private key
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, gerror.WrapCode(gcode.CodeInternalError, err, "failed to generate rsa key")
	}

	// Validate private key
	err = privKey.Validate()
	if err != nil {
		return nil, nil, gerror.WrapCode(gcode.CodeInternalError, err, "failed to validate rsa key")
	}

	// Marshal private key to PKCS#8 format
	privKeyBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return nil, nil, gerror.WrapCode(gcode.CodeInternalError, err, "failed to marshal private key to PKCS#8")
	}

	privateKey = pem.EncodeToMemory(&pem.Block{
		Type:  pemTypePrivateKey,
		Bytes: privKeyBytes,
	})

	// Generate public key
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, nil, gerror.WrapCode(gcode.CodeInternalError, err, "failed to marshal public key")
	}

	publicKey = pem.EncodeToMemory(&pem.Block{
		Type:  pemTypePublicKey,
		Bytes: pubKeyBytes,
	})

	return privateKey, publicKey, nil
}

// GenerateDefaultKeyPair generates a new RSA key pair with default bits (2048).
func GenerateDefaultKeyPair() (privateKey, publicKey []byte, err error) {
	return GenerateKeyPair(DefaultRSAKeyBits)
}

// ExtractPKCS1PublicKey extracts PKCS#1 public key from private key.
func ExtractPKCS1PublicKey(privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "invalid private key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to parse private key")
	}

	pubKeyBytes := x509.MarshalPKCS1PublicKey(&priv.PublicKey)
	return pem.EncodeToMemory(&pem.Block{
		Type:  pemTypeRSAPublicKey,
		Bytes: pubKeyBytes,
	}), nil
}

// ============================================================================
// OAEP Encryption/Decryption Functions (Recommended for new applications)
// ============================================================================

// EncryptOAEP encrypts data with public key using OAEP padding (auto-detect format).
// The publicKey can be either PKCS#1 or PKCS#8 (PKIX) format.
// Uses SHA-256 as the hash function by default.
//
// OAEP (Optimal Asymmetric Encryption Padding) is more secure than PKCS#1 v1.5
// and is recommended for new applications.
//
// Note: For OAEP with SHA-256, max plaintext size = key_size_in_bytes - 2*32 - 2.
// For a 2048-bit key, this is 190 bytes.
func EncryptOAEP(plainText, publicKey []byte) ([]byte, error) {
	return EncryptOAEPWithHash(plainText, publicKey, nil, sha256.New())
}

// EncryptOAEPWithHash encrypts data with public key using OAEP padding with custom hash.
// The publicKey can be either PKCS#1 or PKCS#8 (PKIX) format.
// The label parameter can be nil for most use cases.
// The hash parameter specifies the hash function to use (e.g., sha256.New()).
func EncryptOAEPWithHash(plainText, publicKey, label []byte, hash hash.Hash) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "invalid public key")
	}

	// Try PKCS#8 (PKIX) first
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		// Try PKCS#1
		pub, err = x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to parse public key")
		}
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "not an RSA public key")
	}

	// Validate plaintext size for OAEP padding
	// maxSize = keySize - 2*hashSize - 2
	maxSize := rsaPub.Size() - 2*hash.Size() - 2
	if len(plainText) > maxSize {
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter,
			"plaintext too long: max %d bytes for this key with OAEP, got %d bytes", maxSize, len(plainText))
	}

	return rsa.EncryptOAEP(hash, rand.Reader, rsaPub, plainText, label)
}

// DecryptOAEP decrypts data with private key using OAEP padding (auto-detect format).
// The privateKey can be either PKCS#1 or PKCS#8 format.
// Uses SHA-256 as the hash function by default.
func DecryptOAEP(cipherText, privateKey []byte) ([]byte, error) {
	return DecryptOAEPWithHash(cipherText, privateKey, nil, sha256.New())
}

// DecryptOAEPWithHash decrypts data with private key using OAEP padding with custom hash.
// The privateKey can be either PKCS#1 or PKCS#8 format.
// The label parameter must match the label used during encryption (nil if not used).
// The hash parameter must match the hash function used during encryption.
func DecryptOAEPWithHash(cipherText, privateKey, label []byte, hash hash.Hash) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "invalid private key")
	}

	// Try PKCS#8 first
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS#1
		priv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to parse private key")
		}
	}

	rsaPriv, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "not an RSA private key")
	}

	return rsa.DecryptOAEP(hash, rand.Reader, rsaPriv, cipherText, label)
}

// EncryptOAEPBase64 encrypts data with public key using OAEP padding
// and returns base64-encoded result.
func EncryptOAEPBase64(plainText []byte, publicKeyBase64 string) (string, error) {
	publicKey, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return "", gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to decode public key")
	}

	encrypted, err := EncryptOAEP(plainText, publicKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// DecryptOAEPBase64 decrypts base64-encoded data with private key using OAEP padding.
func DecryptOAEPBase64(cipherTextBase64, privateKeyBase64 string) ([]byte, error) {
	privateKey, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to decode private key")
	}

	cipherText, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInvalidParameter, err, "failed to decode cipher text")
	}

	return DecryptOAEP(cipherText, privateKey)
}
