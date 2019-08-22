// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gkvdb

import (
	"errors"
	"time"

	"github.com/gogf/gf/os/gfile"

	"github.com/dgraph-io/badger"
)

type DB struct {
	options Options
	badger  *badger.DB
}

// New creates and returns a new db object.
func New(options ...Options) *DB {
	if len(options) > 0 {
		return &DB{options: options[0]}
	}
	return &DB{options: DefaultOptions("")}
}

// init does lazy initialization for db.
func (db *DB) init() (err error) {
	if db.badger != nil {
		return nil
	}
	if !gfile.Exists(db.options.Dir) {
		err = gfile.Mkdir(db.options.Dir)
		if err != nil {
			return
		}
	}
	db.badger, err = badger.Open(db.options)
	return
}

// Options return the options of current db object.
func (db *DB) Options() *Options {
	return &db.options
}

// SetOptions sets the options for db.
func (db *DB) SetOptions(options Options) error {
	if db.badger != nil {
		return errors.New("options cannot be changed after db is initialized")
	}
	db.options = options
	return nil
}

// SetPath sets the storage folder path for db.
func (db *DB) SetPath(path string) error {
	if db.badger != nil {
		return errors.New("options cannot be changed after db is initialized")
	}
	db.options.Dir = path
	db.options.ValueDir = path
	return nil
}

// Size returns the data count of current db.
func (db *DB) Size() int64 {
	if err := db.init(); err != nil {
		return 0
	}
	lsm, vlog := db.badger.Size()
	return lsm + vlog
}

// Set sets <key>-<value> pair data to current db with <ttl>.
// The <ttl> is optional, which is not expired in default.
func (db *DB) Set(key []byte, value []byte, ttl ...time.Duration) (err error) {
	if err := db.init(); err != nil {
		return err
	}
	tx := db.Begin(true)
	defer tx.Commit()
	return tx.Set(key, value, ttl...)
}

// Get returns the value with given key.
// It returns nil if <key> is not found in the db.
func (db *DB) Get(key []byte) (value []byte) {
	if err := db.init(); err != nil {
		return
	}
	tx := db.Begin(false)
	defer tx.Rollback()
	return tx.Get(key)
}

// Delete removed data specified by <key> from current db.
func (db *DB) Delete(key []byte) error {
	if err := db.init(); err != nil {
		return err
	}
	tx := db.Begin(true)
	defer tx.Commit()
	return tx.Delete(key)
}

// Close closes the db.
func (db *DB) Close() error {
	if db.badger == nil {
		return nil
	}
	return db.badger.Close()
}

// Iterate is alias of IterateAsc.
// See IterateAsc.
func (db *DB) Iterate(prefix []byte, f func(key, value []byte) bool) {
	db.IterateAsc(prefix, f)
}

// IteratorAsc iterates the db in ascending order
// with given callback function <f> starting from <seek>.
// If <seek> is nil it iterates from the beginning of the db.
// If <f> returns true, then it continues iterating; or false to stop.
func (db *DB) IterateAsc(prefix []byte, f func(key, value []byte) bool) {
	if err := db.init(); err != nil {
		return
	}
	tx := db.Begin(false)
	defer tx.Rollback()
	tx.IterateAsc(prefix, f)
}

// IteratorDesc iterates the db in descending order
// with given callback function <f> starting from <seek>.
// If <prefix> is nil it iterates from the beginning of the db.
// If <f> returns true, then it continues iterating; or false to stop.
func (db *DB) IterateDesc(prefix []byte, f func(key, value []byte) bool) {
	if err := db.init(); err != nil {
		return
	}
	tx := db.Begin(false)
	defer tx.Rollback()
	tx.IterateDesc(prefix, f)
}
