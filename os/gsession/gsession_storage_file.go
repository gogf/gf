// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsession

import (
	"encoding/json"
	"os"
	"time"

	"github.com/gogf/gf/crypto/gaes"

	"github.com/gogf/gf/os/gtimer"

	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/encoding/gbinary"

	"github.com/gogf/gf/os/gtime"

	"github.com/gogf/gf/os/glog"

	"github.com/gogf/gf/os/gfile"
)

// StorageFile implements the Session Storage interface with file system.
type StorageFile struct {
	ttl           time.Duration
	path          string
	updatingIdSet *gset.StrSet
}

var (
	DefaultStorageFilePath         = gfile.Join(gfile.TempDir(), "gsessions")
	DefaultStorageFileCryptoKey    = []byte("Session storage file crypto key!")
	DefaultStorageFileLoopInterval = 5 * time.Second
)

func init() {
	tmpPath := "/tmp"
	if gfile.Exists(tmpPath) && gfile.IsWritable(tmpPath) {
		DefaultStorageFilePath = gfile.Join(tmpPath, "gsessions")
	}
}

func NewStorageFile(ttl time.Duration, path ...string) *StorageFile {
	storagePath := DefaultStorageFilePath
	if len(path) > 0 && path[0] != "" {
		storagePath, _ = gfile.Search(path[0])
		if storagePath == "" {
			glog.Panicf("'%s' does not exist", path[0])
		}
		if !gfile.IsWritable(storagePath) {
			glog.Panicf("'%s' is not writable", path[0])
		}
	}
	if storagePath != "" {
		if err := gfile.Mkdir(storagePath); err != nil {
			glog.Panicf("mkdir '%s' failed: %v", path[0], err)
		}
	}
	s := &StorageFile{
		ttl:           ttl,
		path:          storagePath,
		updatingIdSet: gset.NewStrSet(true),
	}
	gtimer.AddSingleton(DefaultStorageFileLoopInterval, func() {
		s.updatingIdSet.Iterator(func(v string) bool {
			s.doUpdateTTL(v)
			return true
		})
	})
	return s
}

func (s *StorageFile) sessionFilePath(id string) string {
	return gfile.Join(s.path, id)
}

// Get return the session data for given session id.
func (s *StorageFile) Get(id string) map[string]interface{} {
	path := s.sessionFilePath(id)
	data := gfile.GetBytes(path)
	if len(data) > 8 {
		timestamp := gbinary.DecodeToInt64(data[:8])
		if timestamp+int64(s.ttl.Seconds()) < gtime.Second() {
			return nil
		}
		// Decrypt with AES.
		content, err := gaes.Decrypt(data[8:], DefaultStorageFileCryptoKey)
		if err != nil {
			return nil
		}
		var m map[string]interface{}
		if err := json.Unmarshal(content, &m); err != nil {
			return nil
		}
		return m
	}
	return nil
}

// Set updates the content for session id.
// Note that the parameter <content> is the serialized bytes for session map.
func (s *StorageFile) Set(id string, data map[string]interface{}) error {
	path := s.sessionFilePath(id)
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// Encrypt with AES.
	content, err = gaes.Encrypt(content, DefaultStorageFileCryptoKey)
	if err != nil {
		return err
	}
	file, err := gfile.OpenWithFlag(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return err
	}
	if _, err = file.Write(gbinary.EncodeInt64(gtime.Second())); err != nil {
		return err
	}
	if _, err = file.Write(content); err != nil {
		return err
	}
	return file.Close()
}

// UpdateTTL updates the TTL for session id.
// It just adds the session id to the async handling queue.
func (s *StorageFile) UpdateTTL(id string) error {
	s.updatingIdSet.Add(id)
	return nil
}

// doUpdateTTL updates the TTL for session id.
func (s *StorageFile) doUpdateTTL(id string) error {
	path := s.sessionFilePath(id)
	file, err := gfile.OpenWithFlag(path, os.O_WRONLY)
	if err != nil {
		return err
	}
	if _, err = file.Write(gbinary.EncodeInt64(gtime.Second())); err != nil {
		return err
	}
	return file.Close()
}
