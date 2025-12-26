# GoFrame RSA Package

Package `grsa` provides useful API for RSA encryption/decryption algorithms within the GoFrame framework.

## Features

- Generating RSA key pairs in PKCS#1 and PKCS#8 formats
- Encrypting and decrypting data with various key formats
- Handling Base64 encoded keys
- Detecting private key types
- Plaintext size validation

## Quick Start

### Basic Encryption/Decryption

```go
package main

import (
    "fmt"
    "github.com/gogf/gf/v2/crypto/grsa"
)

func main() {
    // Generate a default RSA key pair (2048 bits)
    privateKey, publicKey, err := grsa.GenerateDefaultKeyPair()
    if err != nil {
        panic(err)
    }

    // Data to encrypt
    plainText := []byte("Hello, World!")

    // Encrypt with public key
    cipherText, err := grsa.Encrypt(plainText, publicKey)
    if err != nil {
        panic(err)
    }

    // Decrypt with private key
    decryptedText, err := grsa.Decrypt(cipherText, privateKey)
    if err != nil {
        panic(err)
    }

    fmt.Println(string(decryptedText)) // Output: Hello, World!
}
```

### Working with Base64 Encoded Keys

```go
package main

import (
    "encoding/base64"
    "fmt"
    "github.com/gogf/gf/v2/crypto/grsa"
)

func main() {
    // Generate a key pair
    privateKey, publicKey, err := grsa.GenerateDefaultKeyPair()
    if err != nil {
        panic(err)
    }

    // Encode keys to Base64
    privateKeyBase64 := base64.StdEncoding.EncodeToString(privateKey)
    publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKey)

    // Data to encrypt
    plainText := []byte("Hello, Base64 World!")

    // Encrypt with Base64 encoded public key
    cipherTextBase64, err := grsa.EncryptBase64(plainText, publicKeyBase64)
    if err != nil {
        panic(err)
    }

    // Decrypt with Base64 encoded private key
    decryptedText, err := grsa.DecryptBase64(cipherTextBase64, privateKeyBase64)
    if err != nil {
        panic(err)
    }

    fmt.Println(string(decryptedText)) // Output: Hello, Base64 World!
}
```

## Functions

### Key Generation

- `GenerateKeyPair(bits int)`: Generates a new RSA key pair with the given bits in PKCS#1 format
- `GenerateKeyPairPKCS8(bits int)`: Generates a new RSA key pair with the given bits in PKCS#8 format
- `GenerateDefaultKeyPair()`: Generates a new RSA key pair with default bits (2048) in PKCS#1 format

### General Encryption/Decryption

- `Encrypt(plainText, publicKey []byte)`: Encrypts data with public key (auto-detect format)
- `Decrypt(cipherText, privateKey []byte)`: Decrypts data with private key (auto-detect format)
- `EncryptBase64(plainText []byte, publicKeyBase64 string)`: Encrypts data with base64-encoded public key and returns base64-encoded result
- `DecryptBase64(cipherTextBase64, privateKeyBase64 string)`: Decrypts base64-encoded data with base64-encoded private key

### PKCS#1 Specific Functions

- `EncryptPKCS1(plainText, publicKey []byte)`: Encrypts data with PKCS#1 format public key
- `DecryptPKCS1(cipherText, privateKey []byte)`: Decrypts data with PKCS#1 format private key
- `EncryptPKCS1Base64(plainText []byte, publicKeyBase64 string)`: Encrypts data with PKCS#1 public key and returns base64-encoded result
- `DecryptPKCS1Base64(cipherTextBase64, privateKeyBase64 string)`: Decrypts base64-encoded data with PKCS#1 private key

### PKIX Specific Functions

PKIX (X.509) is the standard format for public keys, used with PKCS#8 private keys.

- `EncryptPKIX(plainText, publicKey []byte)`: Encrypts data with PKIX format public key
- `EncryptPKIXBase64(plainText []byte, publicKeyBase64 string)`: Encrypts data with PKIX public key and returns base64-encoded result
- `DecryptPKCS8(cipherText, privateKey []byte)`: Decrypts data with PKCS#8 format private key
- `DecryptPKCS8Base64(cipherTextBase64, privateKeyBase64 string)`: Decrypts base64-encoded data with PKCS#8 private key

### Deprecated Functions

The following functions are deprecated and will be removed in future versions:

- `EncryptPKCS8(plainText, publicKey []byte)`: Use `EncryptPKIX` instead
- `EncryptPKCS8Base64(plainText []byte, publicKeyBase64 string)`: Use `EncryptPKIXBase64` instead

### Utility Functions

- `GetPrivateKeyType(privateKey []byte)`: Detects the type of private key (PKCS#1 or PKCS#8)
- `GetPrivateKeyTypeBase64(privateKeyBase64 string)`: Detects the type of base64 encoded private key
- `ExtractPKCS1PublicKey(privateKey []byte)`: Extracts PKCS#1 public key from PKCS#1 private key

## Key Formats

The package supports two popular RSA key formats:

1. **PKCS#1**: Traditional RSA key format
   - Private key PEM header: `-----BEGIN RSA PRIVATE KEY-----`
   - Public key PEM header: `-----BEGIN RSA PUBLIC KEY-----`

2. **PKCS#8/PKIX**: More modern and flexible key format
   - Private key PEM header: `-----BEGIN PRIVATE KEY-----`
   - Public key PEM header: `-----BEGIN PUBLIC KEY-----`

Both formats are supported for encryption and decryption operations, with auto-detection capabilities for general functions.

### Technical Background: PKCS#8 vs PKIX

**PKCS#8** is a standard for **private keys** only, not public keys. Public keys use the **PKIX (X.509 SubjectPublicKeyInfo)** format.

| Format | Private Key PEM Header | Public Key PEM Header |
|--------|------------------------|----------------------|
| PKCS#1 | `RSA PRIVATE KEY` | `RSA PUBLIC KEY` |
| PKCS#8/PKIX | `PRIVATE KEY` | `PUBLIC KEY` |

When we refer to a "PKCS#8 key pair", it actually means:
- **Private key**: PKCS#8 format (RFC 5208)
- **Public key**: PKIX/SubjectPublicKeyInfo format (RFC 5280, X.509)

This is why the Go standard library provides `x509.MarshalPKCS8PrivateKey` for private keys but `x509.MarshalPKIXPublicKey` for public keys — there is no `MarshalPKCS8PublicKey` function.

The deprecated `EncryptPKCS8` function was a misnomer because encryption uses public keys, and public keys are in PKIX format, not PKCS#8. The correct function name is `EncryptPKIX`.

## Plaintext Size Limit

RSA encryption has a size limit based on key size. For PKCS#1 v1.5 padding:

- **Max plaintext size = key_size_in_bytes - 11**
- For a 2048-bit key: max 245 bytes
- For a 4096-bit key: max 501 bytes

If you need to encrypt larger data, consider using hybrid encryption (RSA + AES).

## Error Handling

All functions return descriptive errors that can be handled using the GoFrame error package (`gerror`). Errors typically include:

- Invalid key format
- Failed key parsing
- Plaintext too long
- Encryption/decryption failures

Always check for errors in production code to ensure robust handling of edge cases.

## Testing

Run the package tests with:

```bash
go test -v
```
