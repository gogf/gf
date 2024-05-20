// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func TestScanList(t *testing.T) {
	type EntityUser struct {
		Uid  int
		Name string
	}

	type EntityUserDetail struct {
		Uid     int
		Address string
	}

	type EntityUserScores struct {
		Id    int
		Uid   int
		Score int
	}

	// Test for struct attribute.
	gtest.C(t, func(t *gtest.T) {
		type Entity struct {
			User       EntityUser
			UserDetail EntityUserDetail
			UserScores []EntityUserScores
		}

		var (
			err         error
			entities    []Entity
			entityUsers = []EntityUser{
				{Uid: 1, Name: "name1"},
				{Uid: 2, Name: "name2"},
				{Uid: 3, Name: "name3"},
			}
			userDetails = []EntityUserDetail{
				{Uid: 1, Address: "address1"},
				{Uid: 2, Address: "address2"},
			}
			userScores = []EntityUserScores{
				{Id: 10, Uid: 1, Score: 100},
				{Id: 11, Uid: 1, Score: 60},
				{Id: 20, Uid: 2, Score: 99},
			}
		)
		err = gconv.ScanList(entityUsers, &entities, "User")
		t.AssertNil(err)

		err = gconv.ScanList(userDetails, &entities, "UserDetail", "User", "uid")
		t.AssertNil(err)

		err = gconv.ScanList(userScores, &entities, "UserScores", "User", "uid")
		t.AssertNil(err)

		t.Assert(len(entities), 3)
		t.Assert(entities[0].User, entityUsers[0])
		t.Assert(entities[1].User, entityUsers[1])
		t.Assert(entities[2].User, entityUsers[2])

		t.Assert(entities[0].UserDetail, userDetails[0])
		t.Assert(entities[1].UserDetail, userDetails[1])
		t.Assert(entities[2].UserDetail, EntityUserDetail{})

		t.Assert(len(entities[0].UserScores), 2)
		t.Assert(entities[0].UserScores[0], userScores[0])
		t.Assert(entities[0].UserScores[1], userScores[1])

		t.Assert(len(entities[1].UserScores), 1)
		t.Assert(entities[1].UserScores[0], userScores[2])

		t.Assert(len(entities[2].UserScores), 0)
	})

	// Test for pointer attribute.
	gtest.C(t, func(t *gtest.T) {
		type Entity struct {
			User       *EntityUser
			UserDetail *EntityUserDetail
			UserScores []*EntityUserScores
		}

		var (
			err         error
			entities    []*Entity
			entityUsers = []*EntityUser{
				{Uid: 1, Name: "name1"},
				{Uid: 2, Name: "name2"},
				{Uid: 3, Name: "name3"},
			}
			userDetails = []*EntityUserDetail{
				{Uid: 1, Address: "address1"},
				{Uid: 2, Address: "address2"},
			}
			userScores = []*EntityUserScores{
				{Id: 10, Uid: 1, Score: 100},
				{Id: 11, Uid: 1, Score: 60},
				{Id: 20, Uid: 2, Score: 99},
			}
		)
		err = gconv.ScanList(entityUsers, &entities, "User")
		t.AssertNil(err)

		err = gconv.ScanList(userDetails, &entities, "UserDetail", "User", "uid")
		t.AssertNil(err)

		err = gconv.ScanList(userScores, &entities, "UserScores", "User", "uid")
		t.AssertNil(err)

		t.Assert(len(entities), 3)
		t.Assert(entities[0].User, entityUsers[0])
		t.Assert(entities[1].User, entityUsers[1])
		t.Assert(entities[2].User, entityUsers[2])

		t.Assert(entities[0].UserDetail, userDetails[0])
		t.Assert(entities[1].UserDetail, userDetails[1])
		t.Assert(entities[2].UserDetail, nil)

		t.Assert(len(entities[0].UserScores), 2)
		t.Assert(entities[0].UserScores[0], userScores[0])
		t.Assert(entities[0].UserScores[1], userScores[1])

		t.Assert(len(entities[1].UserScores), 1)
		t.Assert(entities[1].UserScores[0], userScores[2])

		t.Assert(len(entities[2].UserScores), 0)
	})

	// Test struct embedded attribute.
	gtest.C(t, func(t *gtest.T) {
		type Entity struct {
			EntityUser
			UserDetail EntityUserDetail
			UserScores []EntityUserScores
		}

		var (
			err         error
			entities    []Entity
			entityUsers = []EntityUser{
				{Uid: 1, Name: "name1"},
				{Uid: 2, Name: "name2"},
				{Uid: 3, Name: "name3"},
			}
			userDetails = []EntityUserDetail{
				{Uid: 1, Address: "address1"},
				{Uid: 2, Address: "address2"},
			}
			userScores = []EntityUserScores{
				{Id: 10, Uid: 1, Score: 100},
				{Id: 11, Uid: 1, Score: 60},
				{Id: 20, Uid: 2, Score: 99},
			}
		)
		err = gconv.Scan(entityUsers, &entities)
		t.AssertNil(err)

		err = gconv.ScanList(userDetails, &entities, "UserDetail", "uid")
		t.AssertNil(err)

		err = gconv.ScanList(userScores, &entities, "UserScores", "uid")
		t.AssertNil(err)

		t.Assert(len(entities), 3)
		t.Assert(entities[0].EntityUser, entityUsers[0])
		t.Assert(entities[1].EntityUser, entityUsers[1])
		t.Assert(entities[2].EntityUser, entityUsers[2])

		t.Assert(entities[0].UserDetail, userDetails[0])
		t.Assert(entities[1].UserDetail, userDetails[1])
		t.Assert(entities[2].UserDetail, EntityUserDetail{})

		t.Assert(len(entities[0].UserScores), 2)
		t.Assert(entities[0].UserScores[0], userScores[0])
		t.Assert(entities[0].UserScores[1], userScores[1])

		t.Assert(len(entities[1].UserScores), 1)
		t.Assert(entities[1].UserScores[0], userScores[2])

		t.Assert(len(entities[2].UserScores), 0)
	})

	// Test struct embedded pointer attribute.
	gtest.C(t, func(t *gtest.T) {
		type Entity struct {
			*EntityUser
			UserDetail *EntityUserDetail
			UserScores []*EntityUserScores
		}

		var (
			err         error
			entities    []Entity
			entityUsers = []EntityUser{
				{Uid: 1, Name: "name1"},
				{Uid: 2, Name: "name2"},
				{Uid: 3, Name: "name3"},
			}
			userDetails = []EntityUserDetail{
				{Uid: 1, Address: "address1"},
				{Uid: 2, Address: "address2"},
			}
			userScores = []EntityUserScores{
				{Id: 10, Uid: 1, Score: 100},
				{Id: 11, Uid: 1, Score: 60},
				{Id: 20, Uid: 2, Score: 99},
			}
		)
		err = gconv.Scan(entityUsers, &entities)
		t.AssertNil(err)

		err = gconv.ScanList(userDetails, &entities, "UserDetail", "uid")
		t.AssertNil(err)

		err = gconv.ScanList(userScores, &entities, "UserScores", "uid")
		t.AssertNil(err)

		t.Assert(len(entities), 3)
		t.Assert(entities[0].EntityUser, entityUsers[0])
		t.Assert(entities[1].EntityUser, entityUsers[1])
		t.Assert(entities[2].EntityUser, entityUsers[2])

		t.Assert(entities[0].UserDetail, userDetails[0])
		t.Assert(entities[1].UserDetail, userDetails[1])
		t.Assert(entities[2].UserDetail, nil)

		t.Assert(len(entities[0].UserScores), 2)
		t.Assert(entities[0].UserScores[0], userScores[0])
		t.Assert(entities[0].UserScores[1], userScores[1])

		t.Assert(len(entities[1].UserScores), 1)
		t.Assert(entities[1].UserScores[0], userScores[2])

		t.Assert(len(entities[2].UserScores), 0)
	})

	// Test for special types.
	gtest.C(t, func(t *gtest.T) {
		type Entity struct {
			User       EntityUser
			UserDetail EntityUserDetail
			UserScores []EntityUserScores
		}

		var (
			err         error
			entities    []Entity
			entityUsers = []EntityUser{
				{Uid: 1, Name: "name1"},
				{Uid: 2, Name: "name2"},
				{Uid: 3, Name: "name3"},
			}
			userDetails = []EntityUserDetail{
				{Uid: 1, Address: "address1"},
				{Uid: 2, Address: "address2"},
			}
			//userScores = []EntityUserScores{
			//	{Id: 10, Uid: 1, Score: 100},
			//	{Id: 11, Uid: 1, Score: 60},
			//	{Id: 20, Uid: 2, Score: 99},
			//}
		)

		err = gconv.ScanList(nil, nil, "")
		t.AssertNil(err)

		err = gconv.ScanList(entityUsers, &entities, "")
		t.AssertNE(err, nil)

		err = gconv.ScanList(entityUsers, &entities, "User")
		t.AssertNil(err)

		err = gconv.ScanList(userDetails, entities, "User")
		t.AssertNE(err, nil)

		var a int = 1
		err = gconv.ScanList(userDetails, &a, "User")
		t.AssertNE(err, nil)
	})
}

func TestScanListErr(t *testing.T) {
	type EntityUser struct {
		Uid  int
		Name string
	}

	type EntityUserDetail struct {
		Uid     int
		Address string
	}

	type EntityUserScores struct {
		Id    int
		Uid   int
		Score int
	}

	gtest.C(t, func(t *gtest.T) {
		type Entity struct {
			User       EntityUser
			UserDetail EntityUserDetail
			UserScores []EntityUserScores
		}

		var (
			err         error
			entities    []Entity
			entityUsers = []EntityUser{
				{Uid: 1, Name: "name1"},
				{Uid: 2, Name: "name2"},
				{Uid: 3, Name: "name3"},
			}
			userDetails = []EntityUserDetail{
				{Uid: 1, Address: "address1"},
				{Uid: 2, Address: "address2"},
			}
			userScores = []EntityUserScores{
				{Id: 10, Uid: 1, Score: 100},
				{Id: 11, Uid: 1, Score: 60},
				{Id: 20, Uid: 2, Score: 99},
			}
		)
		err = gconv.ScanList(entityUsers, &entities, "User")
		t.AssertNil(err)

		err = gconv.ScanList(userDetails, &entities, "UserDetail", "User", "uuid")
		t.AssertNE(err, nil)

		err = gconv.ScanList(userScores, &entities, "UserScores", "User", "uid")
		t.AssertNil(err)
	})
}
