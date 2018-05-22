package gdes

import (
	"testing"
	"bytes"
	"encoding/hex"
	"fmt"
)

func TestDesECB(t *testing.T){
	{
		key := []byte("11111111")
		text := []byte("12345678")
		padding := NOPADDING
		cipherText, err := DesECBEncrypt(key, text, padding)
		if err != nil {
			t.Errorf("%v", err)
		}
  
		clearText, err := DesECBDecrypt(key, cipherText, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		if bytes.Equal(clearText, text) == false {
			t.Errorf("text:%v, clearText:%v", hex.EncodeToString(text), hex.EncodeToString(clearText))
		}
		fmt.Println("clearText:", hex.EncodeToString(clearText), "cipherText:", hex.EncodeToString(cipherText))

	}

	{
		key := []byte("11111111")
		text := []byte("12345678")
		padding := PKCS5PADDING
		cipherText, err := DesECBEncrypt(key, text, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		clearText, err := DesECBDecrypt(key, cipherText, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		if bytes.Equal(clearText, text) == false {
			t.Errorf("text:%v, clearText:%v", hex.EncodeToString(text), hex.EncodeToString(clearText))
		}
		fmt.Println("clearText:", hex.EncodeToString(clearText), "cipherText:", hex.EncodeToString(cipherText))
	}
}

func Test3DesECB(t *testing.T){
	{
		key := []byte("1111111111111234")
		text := []byte("1234567812345678")
		padding := NOPADDING
		cipherText, err := TripleDesECBEncrypt(key, text, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		clearText, err := TripleDesECBDecrypt(key, cipherText, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		if bytes.Equal(clearText, text) == false {
			t.Errorf("text:%v, clearText:%v", hex.EncodeToString(text), hex.EncodeToString(clearText))
		}
		fmt.Println("key:", hex.EncodeToString(key),"clearText:", hex.EncodeToString(clearText), "cipherText:", hex.EncodeToString(cipherText))

	}

	{
		key := []byte("111111111111123412345678")
		text := []byte("123456789")
		padding := PKCS5PADDING
		cipherText, err := TripleDesECBEncrypt(key, text, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		clearText, err := TripleDesECBDecrypt(key, cipherText, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		if bytes.Equal(clearText, text) == false {
			t.Errorf("text:%v, clearText:%v", hex.EncodeToString(text), hex.EncodeToString(clearText))
		}
		fmt.Println("key:", hex.EncodeToString(key),"clearText:", hex.EncodeToString(clearText), "cipherText:", hex.EncodeToString(cipherText))
	}
}

func TestDesCBC(t *testing.T){
	{
		key := []byte("11111111")
		text := []byte("1234567812345678")
		padding := NOPADDING
		iv := []byte("12345678")
		cipherText, err := DesCBCEncrypt(key, text, iv,padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		clearText, err := DesCBCDecrypt(key, cipherText, iv, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		if bytes.Equal(clearText, text) == false {
			t.Errorf("text:%v, clearText:%v", hex.EncodeToString(text), hex.EncodeToString(clearText))
		}
		fmt.Println("key:", hex.EncodeToString(key),"clearText:", hex.EncodeToString(clearText), "cipherText:", hex.EncodeToString(cipherText))

	}

	{
		key := []byte("11111111")
		text := []byte("12345678")
		padding := PKCS5PADDING
		iv := []byte("12345678")
		cipherText, err := DesCBCEncrypt(key, text, iv, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		clearText, err := DesCBCDecrypt(key, cipherText, iv, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		if bytes.Equal(clearText, text) == false {
			t.Errorf("text:%v, clearText:%v", hex.EncodeToString(text), hex.EncodeToString(clearText))
		}
		fmt.Println("key:", hex.EncodeToString(key),"clearText:", hex.EncodeToString(clearText), "cipherText:", hex.EncodeToString(cipherText))
	}
}

func Test3DesCBC(t *testing.T){
	{
		key := []byte("1111111112345678")
		text := []byte("1234567812345678")
		padding := NOPADDING
		iv := []byte("12345678")
		cipherText, err := TripleDesCBCEncrypt(key, text, iv,padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		clearText, err := TripleDesCBCDecrypt(key, cipherText, iv, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		if bytes.Equal(clearText, text) == false {
			t.Errorf("text:%v, clearText:%v", hex.EncodeToString(text), hex.EncodeToString(clearText))
		}
		fmt.Println("key:", hex.EncodeToString(key),"clearText:", hex.EncodeToString(clearText), "cipherText:", hex.EncodeToString(cipherText))

	}

	{
		key := []byte("111111111234567812345678")
		text := []byte("12345678")
		padding := PKCS5PADDING
		iv := []byte("12345678")
		cipherText, err := TripleDesCBCEncrypt(key, text, iv, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		clearText, err := TripleDesCBCDecrypt(key, cipherText, iv, padding)
		if err != nil {
			t.Errorf("%v", err)
		}

		if bytes.Equal(clearText, text) == false {
			t.Errorf("text:%v, clearText:%v", hex.EncodeToString(text), hex.EncodeToString(clearText))
		}
		fmt.Println("key:", hex.EncodeToString(key),"clearText:", hex.EncodeToString(clearText), "cipherText:", hex.EncodeToString(cipherText))
	}
}