// Copyright 2017 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package gjson_test

import (
	"bytes"
	"testing"

	"github.com/jin502437344/gf/encoding/gjson"
)

func Test_Set1(t *testing.T) {
	e := []byte(`{"k1":{"k11":[1,2,3]},"k2":"v2"}`)
	p := gjson.New(map[string]string{
		"k1": "v1",
		"k2": "v2",
	})
	p.Set("k1.k11", []int{1, 2, 3})
	if c, err := p.ToJson(); err == nil {

		if bytes.Compare(c, []byte(`{"k1":{"k11":[1,2,3]},"k2":"v2"}`)) != 0 {
			t.Error("expect:", string(e))
		}
	} else {
		t.Error(err)
	}
}

func Test_Set2(t *testing.T) {
	e := []byte(`[[null,1]]`)
	p := gjson.New([]string{"a"})
	p.Set("0.1", 1)
	if c, err := p.ToJson(); err == nil {

		if bytes.Compare(c, e) != 0 {
			t.Error("expect:", string(e))
		}
	} else {
		t.Error(err)
	}
}

func Test_Set3(t *testing.T) {
	e := []byte(`{"kv":{"k1":"v1"}}`)
	p := gjson.New([]string{"a"})
	p.Set("kv", map[string]string{
		"k1": "v1",
	})
	if c, err := p.ToJson(); err == nil {

		if bytes.Compare(c, e) != 0 {
			t.Error("expect:", string(e))
		}
	} else {
		t.Error(err)
	}
}

func Test_Set4(t *testing.T) {
	e := []byte(`["a",[{"k1":"v1"}]]`)
	p := gjson.New([]string{"a"})
	p.Set("1.0", map[string]string{
		"k1": "v1",
	})
	if c, err := p.ToJson(); err == nil {

		if bytes.Compare(c, e) != 0 {
			t.Error("expect:", string(e))
		}
	} else {
		t.Error(err)
	}
}

func Test_Set5(t *testing.T) {
	e := []byte(`[[[[[[[[[[[[[[[[[[[[[1,2,3]]]]]]]]]]]]]]]]]]]]]`)
	p := gjson.New([]string{"a"})
	p.Set("0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0", []int{1, 2, 3})
	if c, err := p.ToJson(); err == nil {

		if bytes.Compare(c, e) != 0 {
			t.Error("expect:", string(e))
		}
	} else {
		t.Error(err)
	}
}

func Test_Set6(t *testing.T) {
	e := []byte(`["a",[1,2,3]]`)
	p := gjson.New([]string{"a"})
	p.Set("1", []int{1, 2, 3})
	if c, err := p.ToJson(); err == nil {

		if bytes.Compare(c, e) != 0 {
			t.Error("expect:", string(e))
		}
	} else {
		t.Error(err)
	}
}

func Test_Set7(t *testing.T) {
	e := []byte(`{"0":[null,[1,2,3]],"k1":"v1","k2":"v2"}`)
	p := gjson.New(map[string]string{
		"k1": "v1",
		"k2": "v2",
	})
	p.Set("0.1", []int{1, 2, 3})
	if c, err := p.ToJson(); err == nil {

		if bytes.Compare(c, e) != 0 {
			t.Error("expect:", string(e))
		}
	} else {
		t.Error(err)
	}
}

func Test_Set8(t *testing.T) {
	e := []byte(`{"0":[[[[[[null,[1,2,3]]]]]]],"k1":"v1","k2":"v2"}`)
	p := gjson.New(map[string]string{
		"k1": "v1",
		"k2": "v2",
	})
	p.Set("0.0.0.0.0.0.1", []int{1, 2, 3})
	if c, err := p.ToJson(); err == nil {

		if bytes.Compare(c, e) != 0 {
			t.Error("expect:", string(e))
		}
	} else {
		t.Error(err)
	}
}

func Test_Set9(t *testing.T) {
	e := []byte(`{"k1":[null,[1,2,3]],"k2":"v2"}`)
	p := gjson.New(map[string]string{
		"k1": "v1",
		"k2": "v2",
	})
	p.Set("k1.1", []int{1, 2, 3})
	if c, err := p.ToJson(); err == nil {

		if bytes.Compare(c, e) != 0 {
			t.Error("expect:", string(e))
		}
	} else {
		t.Error(err)
	}
}

func Test_Set10(t *testing.T) {
	e := []byte(`{"a":{"b":{"c":1}}}`)
	p := gjson.New(nil)
	p.Set("a.b.c", 1)
	if c, err := p.ToJson(); err == nil {

		if bytes.Compare(c, e) != 0 {
			t.Error("expect:", string(e))
		}
	} else {
		t.Error(err)
	}
}

func Test_Set11(t *testing.T) {
	e := []byte(`{"a":{"b":{}}}`)
	p, _ := gjson.LoadContent([]byte(`{"a":{"b":{"c":1}}}`))
	p.Remove("a.b.c")
	if c, err := p.ToJson(); err == nil {

		if bytes.Compare(c, e) != 0 {
			t.Error("expect:", string(e))
		}
	} else {
		t.Error(err)
	}
}

func Test_Set12(t *testing.T) {
	e := []byte(`[0,1]`)
	p := gjson.New(nil)
	p.Set("0", 0)
	p.Set("1", 1)
	if c, err := p.ToJson(); err == nil {

		if bytes.Compare(c, e) != 0 {
			t.Error("expect:", string(e))
		}
	} else {
		t.Error(err)
	}
}

func Test_Set13(t *testing.T) {
	e := []byte(`{"array":[0,1]}`)
	p := gjson.New(nil)
	p.Set("array.0", 0)
	p.Set("array.1", 1)
	if c, err := p.ToJson(); err == nil {

		if bytes.Compare(c, e) != 0 {
			t.Error("expect:", string(e))
		}
	} else {
		t.Error(err)
	}
}

func Test_Set14(t *testing.T) {
	e := []byte(`{"f":{"a":1}}`)
	p := gjson.New(nil)
	p.Set("f", "m")
	p.Set("f.a", 1)
	if c, err := p.ToJson(); err == nil {

		if bytes.Compare(c, e) != 0 {
			t.Error("expect:", string(e))
		}
	} else {
		t.Error(err)
	}
}
