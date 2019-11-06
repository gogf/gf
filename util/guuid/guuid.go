// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package guuid generates and inspects UUIDs.
//
// This package is a wrapper for most common used UUID package:
// https://github.com/google/uuid
package guuid

import (
	"github.com/google/uuid"
	"os"
)

// UUID representation compliant with specification
// described in RFC 4122.
type UUID = uuid.UUID

// UUID DCE domains.
const (
	DomainPerson = uuid.Domain(0)
	DomainGroup  = uuid.Domain(1)
	DomainOrg    = uuid.Domain(2)
)

// New creates a new random UUID or panics.
func New() UUID {
	return uuid.New()
}

// NewUUID returns a Version 1 UUID based on the current NodeID and clock
// sequence, and the current time.  If the NodeID has not been set by SetNodeID
// or SetNodeInterface then it will be set automatically.  If the NodeID cannot
// be set NewUUID returns nil.  If clock sequence has not been set by
// SetClockSequence then it will be set automatically.  If GetTime fails to
// return the current NewUUID returns nil and an error.
//
// In most cases, New should be used.
func NewUUID() (UUID, error) {
	return uuid.NewUUID()
}

// NewDCEGroup returns a DCE Security (Version 2) UUID in the group
// domain with the id returned by os.Getgid.
//
//  NewDCESecurity(Group, uint32(os.Getgid()))
func NewDCEGroup() (UUID, error) {
	return uuid.NewDCESecurity(DomainGroup, uint32(os.Getgid()))
}

// NewDCEPerson returns a DCE Security (Version 2) UUID in the person
// domain with the id returned by os.Getuid.
//
//  NewDCESecurity(Person, uint32(os.Getuid()))
func NewDCEPerson() (UUID, error) {
	return uuid.NewDCESecurity(DomainPerson, uint32(os.Getuid()))
}

// NewMD5 returns a new MD5 (Version 3) UUID based on the
// supplied name space and data.  It is the same as calling:
//
//  NewHash(md5.New(), space, data, 3)
func NewMD5(space UUID, data []byte) UUID {
	return uuid.NewMD5(space, data)
}

// NewRandom returns a Random (Version 4) UUID.
//
// The strength of the UUIDs is based on the strength of the crypto/rand
// package.
//
// A note about uniqueness derived from the UUID Wikipedia entry:
//
//  Randomly generated UUIDs have 122 random bits.  One's annual risk of being
//  hit by a meteorite is estimated to be one chance in 17 billion, that
//  means the probability is about 0.00000000006 (6 × 10−11),
//  equivalent to the odds of creating a few tens of trillions of UUIDs in a
//  year and having one duplicate.
func NewRandom() (UUID, error) {
	return uuid.NewRandom()
}

// NewSHA1 returns a new SHA1 (Version 5) UUID based on the
// supplied name space and data.  It is the same as calling:
//
//  NewHash(sha1.New(), space, data, 5)
func NewSHA1(space UUID, data []byte) UUID {
	return uuid.NewSHA1(space, data)
}
