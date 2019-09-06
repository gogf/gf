// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gkvdb

import (
	"time"

	"github.com/dgraph-io/badger"
)

// TX is the transaction object for db.
type TX struct {
	db  *DB
	txn *badger.Txn
}

// Begin starts the transaction for db.
func (db *DB) Begin(update bool) *TX {
	if err := db.init(); err != nil {
		return nil
	}
	return &TX{
		db:  db,
		txn: db.badger.NewTransaction(update),
	}
}

// Commit commits the changes and close the transaction.
func (tx *TX) Commit() error {
	return tx.txn.Commit()
}

// Rollback discards the changes and close the transaction.
func (tx *TX) Rollback() {
	tx.txn.Discard()
}

// Set sets <key>-<value> pair data to current db with <ttl> in this transaction.
// The <ttl> is optional, which is not expired in default.
func (tx *TX) Set(key []byte, value []byte, ttl ...time.Duration) error {
	if len(ttl) > 0 && ttl[0] > 0 {
		return tx.txn.SetEntry(badger.NewEntry(key, value).WithTTL(ttl[0]))
	}
	return tx.txn.Set(key, value)
}

// Get returns the value with given key in this transaction.
// It returns nil if <key> is not found in the db.
func (tx *TX) Get(key []byte) (value []byte) {
	item, err := tx.txn.Get(key)
	if err != nil {
		return nil
	}
	if item.IsDeletedOrExpired() {
		return nil
	}
	value, _ = item.ValueCopy(nil)
	return
}

// Delete removed data specified by <key> from current db in this transaction.
func (tx *TX) Delete(key []byte) error {
	return tx.txn.Delete(key)
}

// Iterate is alias of IterateAsc.
// See IterateAsc.
func (tx *TX) Iterate(prefix []byte, f func(key, value []byte) bool) {
	tx.IterateAsc(prefix, f)
}

// IteratorAsc iterates the db in ascending order
// with given callback function <f> starting from <prefix>.
// If <seek> is nil it iterates from the beginning of the db.
// If <f> returns true, then it continues iterating; or false to stop.
func (tx *TX) IterateAsc(prefix []byte, f func(key, value []byte) bool) {
	tx.doIterate(false, prefix, f)
}

// IteratorDesc iterates the db in descending order
// with given callback function <f> starting from <prefix>.
// If <seek> is nil it iterates from the beginning of the db.
// If <f> returns true, then it continues iterating; or false to stop.
func (tx *TX) IterateDesc(prefix []byte, f func(key, value []byte) bool) {
	tx.doIterate(true, prefix, f)
}

func (tx *TX) doIterate(reverse bool, prefix []byte, f func(key, value []byte) bool) {
	options := badger.DefaultIteratorOptions
	if prefix != nil {
		options.Prefix = prefix
	}
	options.Reverse = reverse
	options.PrefetchSize = 10
	options.PrefetchValues = true
	it := tx.txn.NewIterator(options)
	defer it.Close()
	var k, v []byte
	var err error
	var item *badger.Item
	for it.Rewind(); it.Valid(); it.Next() {
		item = it.Item()
		k = item.Key()
		v, err = item.ValueCopy(nil)
		if err != nil {
			return
		}
		if !f(k, v) {
			return
		}
	}
}
