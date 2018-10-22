package mxj

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestStrictModeXml(t *testing.T) {
	fmt.Println("----------------- TestStrictModeXml ...")
	data := []byte(`<document> <name>Bill & Hallett</name> <salute>Duc &amp; 123xx</salute> <goes_by/> <lang>E</lang> </document>`)

	CustomDecoder = &xml.Decoder{Strict:false}
	m, err := NewMapXml(data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("m:",m)
}

func TestStrictModeXmlSeq(t *testing.T) {
	fmt.Println("----------------- TestStrictModeXmlSeq ...")
	data := []byte(`<document> <name>Bill & Hallett</name> <salute>Duc &amp; 123xx</salute> <goes_by/> <lang>E</lang> </document>`)

	CustomDecoder = &xml.Decoder{Strict:false}
	m, err := NewMapXmlSeq(data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("m:",m)
}

func TestStrictModeFail(t *testing.T) {
	fmt.Println("----------------- TestStrictFail ...")
	data := []byte(`<document> <name>Bill & Hallett</name> <salute>Duc &amp; 123xx</salute> <goes_by/> <lang>E</lang> </document>`)

	CustomDecoder = nil
	_, err := NewMapXml(data)
	if err == nil {
		t.Fatal("error not caught: NewMapXml")
	}
	_, err = NewMapXmlSeq(data)
	if err == nil {
		t.Fatal("error not caught: NewMapXmlSeq")
	}
	fmt.Println("OK")
}

