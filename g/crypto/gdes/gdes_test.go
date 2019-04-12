package gdes_test

import (
	"encoding/hex"
	"testing"

	"github.com/gogf/gf/g/crypto/gdes"
	"github.com/gogf/gf/g/test/gtest"
)

var (
	errKey     = []byte("1111111111111234123456789")
	errIv      = []byte("123456789")
	errPadding = 5
)

func TestDesECB(t *testing.T) {
	gtest.Case(t, func() {
		key := []byte("11111111")
		text := []byte("12345678")
		padding := gdes.NOPADDING
		result := "858b176da8b12503"
		// encrypt test
		cipherText, err := gdes.DesECBEncrypt(key, text, padding)
		gtest.AssertEQ(err, nil)
		gtest.AssertEQ(hex.EncodeToString(cipherText),result)
		// decrypt test
		clearText, err := gdes.DesECBDecrypt(key, cipherText, padding)
		gtest.AssertEQ(err, nil)
		gtest.AssertEQ(string(clearText), "12345678")

		// encrypt err test. when throw exception,the err is not equal nil and the string is nil
		errEncrypt, err := gdes.DesECBEncrypt(key, text, errPadding)
		gtest.AssertNE(err, nil)
		gtest.AssertEQ(errEncrypt, nil)
		errEncrypt, err = gdes.DesECBEncrypt(errKey, text, padding)
		gtest.AssertNE(err, nil)
		gtest.AssertEQ(errEncrypt, nil)
		// err decrypt test.
		errDecrypt, err := gdes.DesECBDecrypt(errKey, cipherText, padding)
		gtest.AssertNE(err, nil)
		gtest.AssertEQ(errDecrypt, nil)
		errDecrypt, err = gdes.DesECBDecrypt(key, cipherText, errPadding)
		gtest.AssertNE(err, nil)
		gtest.AssertEQ(errDecrypt, nil)
	})

	gtest.Case(t, func() {
		key := []byte("11111111")
		text := []byte("12345678")
		padding := gdes.PKCS5PADDING
		errPadding := 5
		result := "858b176da8b12503ad6a88b4fa37833d"
		cipherText, err := gdes.DesECBEncrypt(key, text, padding)
		gtest.AssertEQ(err,nil)
		gtest.AssertEQ(hex.EncodeToString(cipherText),result)
		// decrypt test
		clearText, err := gdes.DesECBDecrypt(key, cipherText, padding)
		gtest.AssertEQ(err,nil)
		gtest.AssertEQ(string(clearText),"12345678")

		// err test
		errEncrypt, err := gdes.DesECBEncrypt(key, text, errPadding)
		gtest.AssertNE(err, nil)
		gtest.AssertEQ(errEncrypt, nil)
		errDecrypt, err := gdes.DesECBDecrypt(errKey, cipherText, padding)
		gtest.AssertNE(err, nil)
		gtest.AssertEQ(errDecrypt, nil)
	})
}

func Test3DesECB(t *testing.T) {
	gtest.Case(t, func() {
		key := []byte("1111111111111234")
		text := []byte("1234567812345678")
		padding := gdes.NOPADDING
		result := "a23ee24b98c26263a23ee24b98c26263"
		// encrypt test
		cipherText, err := gdes.TripleDesECBEncrypt(key, text, padding)
		gtest.AssertEQ(err,nil)
		gtest.AssertEQ(hex.EncodeToString(cipherText),result)
		// decrypt test
		clearText, err := gdes.TripleDesECBDecrypt(key, cipherText, padding)
		gtest.AssertEQ(err,nil)
		gtest.AssertEQ(string(clearText),"1234567812345678")
		// err test
		errEncrypt, err := gdes.DesECBEncrypt(key, text, errPadding)
		gtest.AssertNE(err, nil)
		gtest.AssertEQ(errEncrypt, nil)
	})

	gtest.Case(t, func() {
		key := []byte("111111111111123412345678")
		text := []byte("123456789")
		padding := gdes.PKCS5PADDING
		errPadding := 5
		result := "37989b1effc07a6d00ff89a7d052e79f"
		// encrypt test
		cipherText, err := gdes.TripleDesECBEncrypt(key, text, padding)
		gtest.AssertEQ(err,nil)
		gtest.AssertEQ(hex.EncodeToString(cipherText),result)
		// decrypt test
		clearText, err := gdes.TripleDesECBDecrypt(key, cipherText, padding)
		gtest.AssertEQ(err,nil)
		gtest.AssertEQ(string(clearText),"123456789")
		// err test, when key is err, but text and padding is right
		errEncrypt, err := gdes.TripleDesECBEncrypt(errKey, text, padding)
		gtest.AssertNE(err, nil)
		gtest.AssertEQ(errEncrypt, nil)
		// when padding is err,but key and text is right
		errEncrypt, err = gdes.TripleDesECBEncrypt(key, text, errPadding)
		gtest.AssertNE(err, nil)
		gtest.AssertEQ(errEncrypt, nil)
		// decrypt err test,when key is err
		errEncrypt, err = gdes.TripleDesECBDecrypt(errKey, text, padding)
		gtest.AssertNE(err, nil)
		gtest.AssertEQ(errEncrypt, nil)
	})
}

func TestDesCBC(t *testing.T) {
	gtest.Case(t, func() {
		key := []byte("11111111")
		text := []byte("1234567812345678")
		padding := gdes.NOPADDING
		iv := []byte("12345678")
		result := "40826a5800608c87585ca7c9efabee47"
		// encrypt test
		cipherText, err := gdes.DesCBCEncrypt(key, text, iv, padding)
		gtest.AssertEQ(err,nil)
		gtest.AssertEQ(hex.EncodeToString(cipherText),result)
		// decrypt test
		clearText, err := gdes.DesCBCDecrypt(key, cipherText, iv, padding)
		gtest.AssertEQ(err,nil)
		gtest.AssertEQ(string(clearText),"1234567812345678")
		// encrypt err test.
		errEncrypt, err := gdes.DesCBCEncrypt(errKey, text, iv, padding)
		gtest.AssertNE(err, nil)
		gtest.AssertEQ(errEncrypt, nil)
		// the iv is err
		errEncrypt, err = gdes.DesCBCEncrypt(key, text, errIv, padding)
		//gtest.AssertNE(err,nil)
		gtest.AssertEQ(errEncrypt, nil)
		// the padding is err
		errEncrypt, err = gdes.DesCBCEncrypt(key, text, iv, errPadding)
		gtest.AssertNE(err, nil)
		gtest.AssertEQ(errEncrypt, nil)
		// decrypt err test. the key is err
		errDecrypt, err := gdes.DesCBCDecrypt(errKey, cipherText, iv, padding)
		gtest.AssertNE(err, nil)
		gtest.AssertEQ(errDecrypt, nil)
		// the iv is err
		errDecrypt, err = gdes.DesCBCDecrypt(key, cipherText, errIv, padding)
		gtest.AssertNE(err, nil)
		gtest.AssertEQ(errDecrypt, nil)
		// the padding is err
		errDecrypt, err = gdes.DesCBCDecrypt(key, cipherText, iv, errPadding)
		gtest.AssertNE(err, nil)
		gtest.AssertEQ(errDecrypt, nil)
	})

	gtest.Case(t, func() {
		key := []byte("11111111")
		text := []byte("12345678")
		padding := gdes.PKCS5PADDING
		iv := []byte("12345678")
		result := "40826a5800608c87100a25d86ac7c52c"
		// encrypt test
		cipherText, err := gdes.DesCBCEncrypt(key, text, iv, padding)
		gtest.AssertEQ(err,nil)
		gtest.AssertEQ(hex.EncodeToString(cipherText),result)
		// decrypt test
		clearText, err := gdes.DesCBCDecrypt(key, cipherText, iv, padding)
		gtest.AssertEQ(err,nil)
		gtest.AssertEQ(string(clearText),"12345678")
		// err test
		errEncrypt, err := gdes.DesCBCEncrypt(key, text, errIv, padding)
		gtest.AssertNE(err, nil)
		gtest.AssertEQ(errEncrypt, nil)
	})
}

func Test3DesCBC(t *testing.T) {
	gtest.Case(t, func() {
		key := []byte("1111111112345678")
		text := []byte("1234567812345678")
		padding := gdes.NOPADDING
		iv := []byte("12345678")
		result := "bfde1394e265d5f738d5cab170c77c88"
		// encrypt test
		cipherText, err := gdes.TripleDesCBCEncrypt(key, text, iv, padding)
		gtest.AssertEQ(err,nil)
		gtest.AssertEQ(hex.EncodeToString(cipherText),result)
		// decrypt test
		clearText, err := gdes.TripleDesCBCDecrypt(key, cipherText, iv, padding)
		gtest.AssertEQ(err,nil)
		gtest.AssertEQ(string(clearText),"1234567812345678")
		// encrypt err test
		errEncrypt, err := gdes.TripleDesCBCEncrypt(errKey, text, iv, padding)
		gtest.AssertNE(err, nil)
		gtest.AssertEQ(errEncrypt, nil)
		// the iv is err
		errEncrypt, err = gdes.TripleDesCBCEncrypt(key, text, errIv, padding)
		gtest.AssertNE(err, nil)
		gtest.AssertEQ(errEncrypt, nil)
		// the padding is err
		errEncrypt, err = gdes.TripleDesCBCEncrypt(key, text, iv, errPadding)
		gtest.AssertNE(err, nil)
		gtest.AssertEQ(errEncrypt, nil)
		// decrypt err test
		errDecrypt, err := gdes.TripleDesCBCDecrypt(errKey, cipherText, iv, padding)
		gtest.AssertNE(err, nil)
		gtest.AssertEQ(errDecrypt, nil)
		// the iv is err
		errDecrypt, err = gdes.TripleDesCBCDecrypt(key, cipherText, errIv, padding)
		gtest.AssertNE(err, nil)
		gtest.AssertEQ(errDecrypt, nil)
		// the padding is err
		errDecrypt, err = gdes.TripleDesCBCDecrypt(key, cipherText, iv, errPadding)
		gtest.AssertNE(err, nil)
		gtest.AssertEQ(errDecrypt, nil)
	})
	gtest.Case(t, func() {
		key := []byte("111111111234567812345678")
		text := []byte("12345678")
		padding := gdes.PKCS5PADDING
		iv := []byte("12345678")
		result := "40826a5800608c87100a25d86ac7c52c"
		// encrypt test
		cipherText, err := gdes.TripleDesCBCEncrypt(key, text, iv, padding)
		gtest.AssertEQ(err,nil)
		gtest.AssertEQ(hex.EncodeToString(cipherText),result)
		// decrypt test
		clearText, err := gdes.TripleDesCBCDecrypt(key, cipherText, iv, padding)
		gtest.AssertEQ(err,nil)
		gtest.AssertEQ(string(clearText),"12345678")
	})

}
