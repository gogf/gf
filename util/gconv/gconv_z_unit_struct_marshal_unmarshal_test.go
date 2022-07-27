// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/crypto/gcrc32"
	"github.com/gogf/gf/v2/encoding/gbinary"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
	"reflect"
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
		t.AssertNil(err)
		t.Assert(st.ServiceDate.Time.Format("2006-01-02 15:04:05"), "2020-10-10 12:00:01")
	})
	gtest.C(t, func(t *gtest.T) {
		st := &MyTimeSt{}
		err := gconv.Struct(g.Map{"ServiceDate": nil}, st)
		t.AssertNil(err)
		t.Assert(st.ServiceDate.Time.IsZero(), true)
	})
	gtest.C(t, func(t *gtest.T) {
		st := &MyTimeSt{}
		err := gconv.Struct(g.Map{"ServiceDate": "error"}, st)
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
		return gerror.New("invalid package length")
	}
	p.Length = gbinary.DecodeToUint16(b[:2])
	if len(b) < int(p.Length) {
		return gerror.New("invalid data length")
	}
	p.Crc32 = gbinary.DecodeToUint32(b[2:6])
	p.Data = b[6:]
	if gcrc32.Encrypt(p.Data) != p.Crc32 {
		return gerror.New("crc32 validation failed")
	}
	return nil
}

func Test_Struct_UnmarshalValue2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var p1, p2 *Pkg
		p1 = NewPkg([]byte("123"))
		err := gconv.Struct(p1.Marshal(), &p2)
		t.AssertNil(err)
		t.Assert(p1, p2)
	})
}

type BaseModel struct {
}

type SecretModel struct {
	BaseModel
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
}

func (m *BaseModel) UnmarshalValueWithTarget(target interface{}, value interface{}) (error, bool) {
	pointerReflectValue := reflect.ValueOf(target)
	pointerElemReflectValue := pointerReflectValue.Elem()
	if v, ok := value.(g.Map); ok {
		for fieldName, fieldVal := range v {
			structFieldValue := pointerElemReflectValue.FieldByName(fieldName)
			if !structFieldValue.IsValid() {
				continue
			}
			// CanSet checks whether attribute is public accessible.
			if !structFieldValue.CanSet() {
				continue
			}
			if fieldName == "CreatedAt" {
				fieldVal = gvar.New(fieldVal).Int64() + 10
			}
			convertedValue := gconv.Convert(fieldVal, structFieldValue.Type().String())
			structFieldValue.Set(reflect.ValueOf(convertedValue))
		}
		return nil, ok
	} else {
		return nil, ok
	}
}

func Test_Struct_UnmarshalValueWithTarget(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var now = gtime.Now().Unix()
		var secretModel = &SecretModel{}
		// g.Map->UnmarshalValueWithTarget struct转换
		var params = g.Map{"Id": 100, "Name": "TEST", "CreatedAt": now}
		if err := gconv.Struct(params, secretModel); err != nil {
			panic(err)
		} else {
			t.Assert(secretModel.Id, params["Id"])
			t.Assert(secretModel.Name, params["Name"])
			t.Assert(secretModel.CreatedAt, now+10) // UnmarshalValueWithTarget 转换后的值
		}
		// struct->UnmarshalValueWithTarget struct转换
		var paramStruct = struct {
			Id        int64  `json:"id"`
			Name      string `json:"name"`
			CreatedAt int64  `json:"created_at"`
		}{Id: 200, Name: "NEW", CreatedAt: now}
		var secretModel2 = &SecretModel{}
		if err := gconv.Struct(paramStruct, secretModel2); err != nil {
			panic(err)
		} else {
			t.Assert(secretModel2.Id, paramStruct.Id)
			t.Assert(secretModel2.Name, paramStruct.Name)
			t.Assert(secretModel2.CreatedAt, paramStruct.CreatedAt)
		}
		// g.Map->normal struct转换
		var paramStruct2 = &struct {
			Id        int64  `json:"id"`
			Name      string `json:"name"`
			CreatedAt int64  `json:"created_at"`
		}{}
		if err := gconv.Struct(params, paramStruct2); err != nil {
			panic(err)
		} else {
			t.Assert(paramStruct2.Id, params["Id"])
			t.Assert(paramStruct2.Name, params["Name"])
			t.Assert(paramStruct2.CreatedAt, now)
		}
	})
}
