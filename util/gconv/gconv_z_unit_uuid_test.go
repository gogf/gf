// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func TestBsToUUID(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		err := initConverter()
		t.Assert(err, nil)
		v, _ := uuid.NewV7()
		b := v[:]
		var u uuid.UUID

		err = gconv.Scan(v, &u)
		t.Assert(err, nil)
		t.Assert(v, u)

		err = gconv.Scan(b, &u)
		t.Assert(err, nil)
		t.Assert(v, u)

	})
}

func TestUUIDToBS(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		err := initConverter()
		t.Assert(err, nil)
		v, _ := uuid.NewV7()

		var bs []byte
		err = gconv.Scan(v, &bs)

		t.Assert(err, nil)
		t.Assert(v[:], bs)

		var u [16]byte

		err = gconv.Scan(v, &u)
		t.Assert(err, nil)
		t.Assert([16]byte(v), u)
	})
}

func initConverter() (err error) {
	for _, fn := range []any{
		convBsToUUID,
		convStrToUUID,
		convUUIDToBs,
		convUUIDToStr,
		convGvarToUUID,
		convUUIDToArray,
	} {
		if err = gconv.RegisterTypeConverterFunc(fn); err != nil {
			return
		}
	}
	return
}

func convGvarToUUID(v gvar.Var) (u *uuid.UUID, err error) {
	return convBsToUUID(v.Bytes())
}

func convBsToUUID(bs []byte) (uid *uuid.UUID, err error) {
	uid = new(uuid.UUID)
	if len(bs) == 16 {
		*uid, err = uuid.FromBytes(bs)
	} else {
		*uid, err = uuid.ParseBytes(bs)
	}
	return
}

func convUUIDToArray(uid uuid.UUID) (bs *[16]byte, err error) {
	bs = new([16]byte)
	*bs = uid
	return
}

func convUUIDToBs(uid uuid.UUID) (bs *[]byte, err error) {
	bs = new([]byte)
	*bs = uid[:]
	return
}

func convUUIDToStr(uid uuid.UUID) (str *string, err error) {
	str = new(string)
	*str = uid.String()
	return
}

func convStrToUUID(str string) (uid *uuid.UUID, err error) {
	return convBsToUUID([]byte(str))
}
