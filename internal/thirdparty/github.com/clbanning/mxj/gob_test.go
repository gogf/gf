package mxj

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"
)

var gobData = map[string]interface{}{
	"one":   1,
	"two":   2.0001,
	"three": "tres",
	"four":  []int{1, 2, 3, 4},
	"five":  map[string]interface{}{"hi": "there"}}

func TestGobHeader(t *testing.T) {
	fmt.Println("\n----------------  gob_test.go ...")
}

func TestNewMapGob(t *testing.T) {
	var buf bytes.Buffer
	gob.Register(gobData)
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(gobData); err != nil {
		t.Fatal("enc.Encode err:", err.Error())
	}
	// decode 'buf' into a Map - map[string]interface{}
	m, err := NewMapGob(buf.Bytes())
	if err != nil {
		t.Fatal("NewMapGob err:", err.Error())
	}
	fmt.Printf("m: %v\n", m)
}

func TestMapGob(t *testing.T) {
	mv := Map(gobData)
	g, err := mv.Gob()
	if err != nil {
		t.Fatal("m.Gob err:", err.Error())
	}
	// decode 'g' into a map[string]interface{}
	m := make(map[string]interface{})
	r := bytes.NewReader(g)
	dec := gob.NewDecoder(r)
	if err := dec.Decode(&m); err != nil {
		t.Fatal("dec.Decode err:", err.Error())
	}
	fmt.Printf("m: %v\n", m)
}

func TestGobSymmetric(t *testing.T) {
	mv := Map(gobData)
	fmt.Printf("mv: %v\n", mv)
	g, err := mv.Gob()
	if err != nil {
		t.Fatal("m.Gob err:", err.Error())
	}
	m, err := NewMapGob(g)
	if err != nil {
		t.Fatal("NewMapGob err:", err.Error())
	}
	fmt.Printf("m : %v\n", m)
}
