# GoFrame RSA Package

Package `grsa` provides useful API for RSA encryption/decryption algorithms within the GoFrame framework.

## Features

- Generating RSA key pairs in PKCS#1 and PKCS#8 formats
- Encrypting and decrypting data with various key formats
- Handling Base64 encoded keys
- Detecting private key types

## Installation

```bash
go get github.com/gogf/gf/v2/crypto/grsa
```

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

- `GenerateKeyPair(bits int)`: Generates a new RSA key pair with the given bits
- `GenerateKeyPairPKCS8(bits int)`: Generates a new RSA key pair with the given bits in PKCS#8 format
- `GenerateDefaultKeyPair()`: Generates a new RSA key pair with default bits (2048)

### General Encryption/Decryption

- `Encrypt(plainText, publicKey []byte)`: Encrypts data with public key (auto-detect format)
- `Decrypt(cipherText, privateKey []byte)`: Decrypts data with private key (auto-detect format)
- `EncryptBase64(plainText []byte, publicKeyBase64 string)`: Encrypts data with base64-encoded public key
- `DecryptBase64(cipherTextBase64, privateKeyBase64 string)`: Decrypts base64-encoded data with base64-encoded private key

### PKCS#1 Specific Functions

- `EncryptPKCS1(plainText, publicKey []byte)`: Encrypts data with public key by PKCS#1 format
- `DecryptPKCS1(cipherText, privateKey []byte)`: Decrypts data with private key by PKCS#1 format
- `EncryptPKCS1Base64(plainText []byte, publicKeyBase64 string)`: Encrypts data with public key by PKCS#1 format and encode result with base64
- `DecryptPKCS1Base64(cipherTextBase64, privateKeyBase64 string)`: Decrypts data with private key by PKCS#1 format and decode base64 input

### PKCS#8 Specific Functions

- `EncryptPKCS8(plainText, publicKey []byte)`: Encrypts data with public key by PKCS#8 format
- `DecryptPKCS8(cipherText, privateKey []byte)`: Decrypts data with private key by PKCS#8 format
- `EncryptPKCS8Base64(plainText []byte, publicKeyBase64 string)`: Encrypts data with public key by PKCS#8 format and encode result with base64
- `DecryptPKCS8Base64(cipherTextBase64, privateKeyBase64 string)`: Decrypts data with private key by PKCS#8 format and decode base64 input

### Utility Functions

- `GetPrivateKeyType(privateKey []byte)`: Detects the type of private key (PKCS#1 or PKCS#8)
- `GetPrivateKeyTypeBase64(privateKeyBase64 string)`: Detects the type of base64 encoded private key (PKCS#1 or PKCS#8)
- `ExtractPKCS1PublicKey(privateKey []byte)`: Extracts PKCS#1 public key from private key

## Key Formats

The package supports two popular RSA key formats:

1. **PKCS#1**: Traditional RSA key format
2. **PKCS#8**: More modern and flexible key format

Both formats are supported for encryption and decryption operations, with auto-detection capabilities for general functions.

## Error Handling

All functions return descriptive errors that can be handled using the GoFrame error package (`gerror`). Errors typically include:

- Invalid key format
- Failed key parsing
- Encryption/decryption failures

Always check for errors in production code to ensure robust handling of edge cases.

## Testing

Run the package tests with:

```bash
go test -v
```

## License

`GoFrame` is licensed under the [MIT License](LICENSE), 100% free and open-source, forever.