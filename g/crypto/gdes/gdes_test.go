package gdes_test

import (
	"github.com/gogf/gf/g/test/gtest"
	"testing"
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/gogf/gf/g/crypto/gdes"
)

var(
	errKey = []byte("1111111111111234123456789")
	errIv = []byte("12345678")
	errPadding = 5
)

func TestDesECB(t *testing.T){
	gtest.Case(t, func() {
		key := []byte("11111111")
		text := []byte("12345678")
		padding := gdes.NOPADDING
		cipherText, err := gdes.DesECBEncrypt(key, text, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		clearText, err := gdes.DesECBDecrypt(key, cipherText, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		if bytes.Equal(clearText, text) == false {
			t.Errorf("text:%v, clearText:%v", hex.EncodeToString(text), hex.EncodeToString(clearText))
		}
		fmt.Println("clearText:", hex.EncodeToString(clearText), "cipherText:", hex.EncodeToString(cipherText))

		// err test
		gdes.DesECBEncrypt(key, text, errPadding)
		gdes.DesECBDecrypt(errKey, cipherText, padding)
	})

	gtest.Case(t, func() {
		key := []byte("11111111")
		text := []byte("12345678")
		padding := gdes.PKCS5PADDING
		errPadding := 5
		cipherText, err := gdes.DesECBEncrypt(key, text, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		clearText, err := gdes.DesECBDecrypt(key, cipherText, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		if bytes.Equal(clearText, text) == false {
			t.Errorf("text:%v, clearText:%v", hex.EncodeToString(text), hex.EncodeToString(clearText))
		}
		fmt.Println("clearText:", hex.EncodeToString(clearText), "cipherText:", hex.EncodeToString(cipherText))

		// err test
		gdes.DesECBEncrypt(key, text, errPadding)
		gdes.DesECBDecrypt(errKey, cipherText, padding)
	})
}

func Test3DesECB(t *testing.T){
	{
		key := []byte("1111111111111234")
		text := []byte("1234567812345678")
		padding := gdes.NOPADDING
		cipherText, err := gdes.TripleDesECBEncrypt(key, text, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		clearText, err := gdes.TripleDesECBDecrypt(key, cipherText, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		if bytes.Equal(clearText, text) == false {
			t.Errorf("text:%v, clearText:%v", hex.EncodeToString(text), hex.EncodeToString(clearText))
		}
		fmt.Println("key:", hex.EncodeToString(key),"clearText:", hex.EncodeToString(clearText), "cipherText:", hex.EncodeToString(cipherText))
		// err test
		gdes.DesECBEncrypt(key, text, errPadding)
	}

	gtest.Case(t, func() {
		key := []byte("111111111111123412345678")
		text := []byte("123456789")
		padding := gdes.PKCS5PADDING
		errPadding := 5
		cipherText, err := gdes.TripleDesECBEncrypt(key, text, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		clearText, err := gdes.TripleDesECBDecrypt(key, cipherText, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		if bytes.Equal(clearText, text) == false {
			t.Errorf("text:%v, clearText:%v", hex.EncodeToString(text), hex.EncodeToString(clearText))
		}
		fmt.Println("key:", hex.EncodeToString(key),"clearText:", hex.EncodeToString(clearText), "cipherText:", hex.EncodeToString(cipherText))
		// err test
		gdes.TripleDesECBEncrypt(errKey, text, padding)
		gdes.TripleDesECBEncrypt(key, text, errPadding)

		gdes.TripleDesECBDecrypt(errKey, text, padding)
	})
}

func TestDesCBC(t *testing.T){
	gtest.Case(t, func() {
		key := []byte("11111111")
		text := []byte("1234567812345678")
		padding := gdes.NOPADDING
		iv := []byte("12345678")
		cipherText, err := gdes.DesCBCEncrypt(key, text, iv,padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		clearText, err := gdes.DesCBCDecrypt(key, cipherText, iv, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		if bytes.Equal(clearText, text) == false {
			t.Errorf("text:%v, clearText:%v", hex.EncodeToString(text), hex.EncodeToString(clearText))
		}
		fmt.Println("key:", hex.EncodeToString(key),"clearText:", hex.EncodeToString(clearText), "cipherText:", hex.EncodeToString(cipherText))
		// err test
		gdes.DesCBCEncrypt(errKey, text, iv,padding)
		gdes.DesCBCEncrypt(errKey, text, iv,errPadding)
	})

	gtest.Case(t, func() {
		key := []byte("11111111")
		text := []byte("12345678")
		padding := gdes.PKCS5PADDING
		iv := []byte("12345678")
		cipherText, err := gdes.DesCBCEncrypt(key, text, iv, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		clearText, err := gdes.DesCBCDecrypt(key, cipherText, iv, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		if bytes.Equal(clearText, text) == false {
			t.Errorf("text:%v, clearText:%v", hex.EncodeToString(text), hex.EncodeToString(clearText))
		}
		fmt.Println("key:", hex.EncodeToString(key),"clearText:", hex.EncodeToString(clearText), "cipherText:", hex.EncodeToString(cipherText))

		// err test
		gdes.DesCBCEncrypt(key, text, errIv, padding)
		gdes.DesCBCEncrypt(key, text, errIv, padding)
	})
}

func Test3DesCBC(t *testing.T){
	gtest.Case(t, func() {
		key := []byte("1111111112345678")
		text := []byte("1234567812345678")
		padding := gdes.NOPADDING
		iv := []byte("12345678")
		cipherText, err := gdes.TripleDesCBCEncrypt(key, text, iv,padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		clearText, err := gdes.TripleDesCBCDecrypt(key, cipherText, iv, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		if bytes.Equal(clearText, text) == false {
			t.Errorf("text:%v, clearText:%v", hex.EncodeToString(text), hex.EncodeToString(clearText))
		}
		fmt.Println("key:", hex.EncodeToString(key),"clearText:", hex.EncodeToString(clearText), "cipherText:", hex.EncodeToString(cipherText))

		gdes.TripleDesCBCEncrypt(errKey, text, iv,padding)
	})
	gtest.Case(t, func() {
		key := []byte("111111111234567812345678")
		text := []byte("12345678")
		padding := gdes.PKCS5PADDING
		iv := []byte("12345678")
		cipherText, err := gdes.TripleDesCBCEncrypt(key, text, iv, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		clearText, err := gdes.TripleDesCBCDecrypt(key, cipherText, iv, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		if bytes.Equal(clearText, text) == false {
			t.Errorf("text:%v, clearText:%v", hex.EncodeToString(text), hex.EncodeToString(clearText))
		}
		fmt.Println("key:", hex.EncodeToString(key),"clearText:", hex.EncodeToString(clearText), "cipherText:", hex.EncodeToString(cipherText))
	})

}