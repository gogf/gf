// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"errors"
	"github.com/gogf/gf/crypto/gcrc32"
	"github.com/gogf/gf/encoding/gbinary"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gconv"
	"testing"
	"time"
)

type MyTime struct {
	time.Time
}

type MyTimeSt struct {
	ServiceDate MyTime
}

func (st *MyTimeSt) UnmarshalValue(v interface{}) error {
	m := gconv.Map(v)
	t, err := gtime.StrToTime(gconv.String(m["ServiceDate"]))
	if err != nil {
		return err
	}
	st.ServiceDate = MyTime{t.Time}
	return nil
}

func Test_Struct_UnmarshalValue1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		st := &MyTimeSt{}
		err := gconv.Struct(g.Map{"ServiceDate": "2020-10-10 12:00:01"}, st)
		t.Assert(err, nil)
		t.Assert(st.ServiceDate.Time.Format("2006-01-02 15:04:05"), "2020-10-10 12:00:01")
	})
	gtest.C(t, func(t *gtest.T) {
		st := &MyTimeSt{}
		err := gconv.Struct(g.Map{"ServiceDate": nil}, st)
		t.AssertNE(err, nil)
	})
}

type Pkg struct {
	Length uint16 // Total length.
	Crc32  uint32 // CRC32.
	Data   []byte
}

// NewPkg creates and returns a package with given data.
func NewPkg(data []byte) *Pkg {
	return &Pkg{
		Length: uint16(len(data) + 6),
		Crc32:  gcrc32.Encrypt(data),
		Data:   data,
	}
}

// Marshal encodes the protocol struct to bytes.
func (p *Pkg) Marshal() []byte {
	b := make([]byte, 6+len(p.Data))
	copy(b, gbinary.EncodeUint16(p.Length))
	copy(b[2:], gbinary.EncodeUint32(p.Crc32))
	copy(b[6:], p.Data)
	return b
}

// UnmarshalValue decodes bytes to protocol struct.
func (p *Pkg) UnmarshalValue(v interface{}) error {
	b := gconv.Bytes(v)
	if len(b) < 6 {
		return errors.New("invalid package length")
	}
	p.Length = gbinary.DecodeToUint16(b[:2])
	if len(b) < int(p.Length) {
		return errors.New("invalid data length")
	}
	p.Crc32 = gbinary.DecodeToUint32(b[2:6])
	p.Data = b[6:]
	if gcrc32.Encrypt(p.Data) != p.Crc32 {
		return errors.New("crc32 validation failed")
	}
	return nil
}

func Test_Struct_UnmarshalValue2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var p1, p2 *Pkg
		p1 = NewPkg([]byte("123"))
		err := gconv.Struct(p1.Marshal(), &p2)
		t.Assert(err, nil)
		t.Assert(p1, p2)
	})
}
